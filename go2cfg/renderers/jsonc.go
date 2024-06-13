package renderers

import (
	"encoding/json"
	"fmt"
	"github.com/modulo-srl/mu-config/go2cfg/distiller"
	"github.com/modulo-srl/mu-config/go2cfg/ordered"
	"go/types"
	"strings"
)

// Jsonc renders JSONC code from distiller info.
type Jsonc struct {
	docTypesMode DocTypesMode
}

// NewJsonc creates a new JSONC renderer.
// mode controls the rendering of field types in JSONC comments.
func NewJsonc(mode DocTypesMode) *Jsonc {
	return &Jsonc{docTypesMode: mode}
}

func (j *Jsonc) RenderStruct(info *distiller.StructInfo, defaults interface{}, indent string,
	embedded bool, parentShadowing []string) (string, error) {
	var builder strings.Builder

	if !embedded {
		builder.WriteString("{\n")
		indent += "\t"
	}

	var shadowing []string
	for _, field := range info.Fields {
		if !field.IsEmbedded {
			shadowing = append(shadowing, field.Name)
		}
	}

	comma := ""
	blockSpacing := false
	for i, field := range info.Fields {
		name := field.Name

		// This field will be shadowed by another one, so skip it.
		if (!field.IsEmbedded && lastIndexOf(shadowing, name) > i) ||
			(embedded && lastIndexOf(parentShadowing, name) != -1) {
			continue
		}

		builder.WriteString(comma)

		if jsonName, ok := field.Tags["json"]; ok {
			name = jsonName
		}

		key := field.Name
		if field.IsEmbedded {
			key = field.Type.String()
			if pathEnd := strings.LastIndex(key, "/"); pathEnd >= 0 {
				key = key[pathEnd+strings.Index(key[pathEnd+1:], ".")+2:]
			}
		}

		var value interface{}
		ok := false
		if defaults != nil {
			value, ok = defaults.(map[string]interface{})[key]
		}

		consts := distiller.LookupTypedConsts(field.Type.String())

		renderType := j.docTypesMode == AllFields

		// No default defined for this field, if named (struct) or array will be rendered below.
		_, isNamed := field.Type.(*types.Named)
		if !ok && field.Layout == distiller.LayoutSingle && (consts != nil || !isNamed) {
			if consts != nil {
				value = consts[0].Value
			} else {
				value = typeZero(field)
			}
		} else {
			var err error
			switch field.Layout {
			case distiller.LayoutSingle:
				if isNamed && consts == nil {
					subInfo := distiller.LookupStruct(field.Type.String())
					if subInfo == nil {
						return "", fmt.Errorf("cannot lookup structure %s", field.Type.String())
					}

					value, err = j.RenderStruct(subInfo, value, indent, field.IsEmbedded, shadowing[i:])
					if err != nil {
						return "", err
					}
				} else {
					basicT, ok := field.Type.(*types.Basic)
					if ok && basicT.Kind() == types.String {
						var b []byte
						b, err = json.Marshal(unescapeString(value))
						if err != nil {
							return "", err
						}
						value = string(b)
					}

					// No special handling required for basic types.
					renderType = renderType || (j.docTypesMode == BasicFields)
				}

			case distiller.LayoutArray:
				if value == nil {
					// Add an example item in case of nil array.
					value, err = j.RenderArray(field, []interface{}{nil}, indent)
				} else {
					value, err = j.RenderArray(field, value.([]interface{}), indent)
				}

			case distiller.LayoutMap:
				value, err = j.RenderMap(field, value.(*ordered.Map), indent)
			}

			if err != nil {
				return "", err
			}
		}

		if field.IsEmbedded {
			builder.WriteString(fmt.Sprintf("%v", value))
		} else {
			doc := renderDoc(field, indent, "//", renderType)
			if doc != "" {
				// Adds a blank line when the comment block is present.
				if !blockSpacing && (comma != "") {
					builder.WriteString("\n")
				}
				blockSpacing = true

				builder.WriteString(doc)
			} else {
				blockSpacing = false
			}

			builder.WriteString(fmt.Sprintf("%s\"%s\": %v", indent, name, value))
		}

		comma = ",\n"
		if blockSpacing {
			comma += "\n"
		}
	}

	if !embedded {
		if comma != "" {
			builder.WriteString("\n")
		}

		builder.WriteString(indent[:len(indent)-1] + "}")
	}

	return builder.String(), nil
}

func (j *Jsonc) RenderArray(field *distiller.FieldInfo, value []interface{}, indent string) (string, error) {
	if len(value) == 0 {
		return "[]", nil
	}

	eltsIdent := indent + "\t"
	code := "[\n"
	for _, elt := range value {
		literal, err := j.RenderElement(field.EltType, elt, eltsIdent)
		if err != nil {
			return "", err
		}

		code += eltsIdent + literal + ",\n"
	}
	code = strings.TrimRight(code, ",\n") + "\n" + indent + "]"

	return code, nil
}

func (j *Jsonc) RenderMap(field *distiller.FieldInfo, value *ordered.Map, indent string) (string, error) {
	if field.IsEmbedded {
		return "", fmt.Errorf("field of slice or map type cannot be embedded")
	}

	if value.Len() == 0 {
		return "{}", nil
	}

	eltsIndent := indent + "\t"
	code := "{\n"

	var err error
	value.Iterate(func(key string, elt interface{}) bool {
		var literal string
		literal, err = j.RenderElement(field.EltType, elt, eltsIndent)
		if err != nil {
			return false
		}

		code += eltsIndent + fmt.Sprintf("%s: %s", key, literal) + ",\n"
		return true
	})

	if err != nil {
		return "", err
	}

	code = strings.TrimRight(code, ",\n") + "\n" + indent + "}"

	return code, nil
}

func (j *Jsonc) RenderElement(itemType types.Type, item interface{}, indent string) (string, error) {
	_, ok := itemType.(*types.Basic)
	if ok || distiller.LookupTypedConsts(itemType.String()) != nil {
		return fmt.Sprintf("%v", item), nil
	}

	subInfo := distiller.LookupStruct(itemType.String())
	if subInfo == nil {
		return "", fmt.Errorf("cannot lookup structure %s", itemType.String())
	}

	return j.RenderStruct(subInfo, item, indent, false, nil)
}
