package renderers

import (
	"encoding/json"
	"fmt"
	"github.com/modulo-srl/mu-config/go2cfg/distiller"
	"github.com/modulo-srl/mu-config/go2cfg/ordered"
	"go/types"
	"regexp"
	"strings"
)

// Yaml renders Yaml code from distiller info.
type Yaml struct {
	docTypesMode DocTypesMode
	path         string
	inArray      bool
	indent       string
}

// NewYaml creates a new Yaml renderer.
// mode controls the rendering of field types in Yaml comments.
func NewYaml(mode DocTypesMode, indentSize int) *Yaml {
	return &Yaml{
		docTypesMode: mode,
		inArray:      false,
		indent:       strings.Repeat(" ", indentSize),
	}
}

func (y *Yaml) RenderStruct(info *distiller.StructInfo, defaults interface{}, indent string, embedded bool, parentShadowing []string) (string, error) {
	var builder strings.Builder

	newline := ""

	var shadowing []string
	for _, field := range info.Fields {
		if !field.IsEmbedded {
			shadowing = append(shadowing, field.Name)
		}
	}

	for i, field := range info.Fields {
		name := field.Name

		// This field will be shadowed by another one, so skip it.
		if (!field.IsEmbedded && lastIndexOf(shadowing, name) > i) ||
			(embedded && lastIndexOf(parentShadowing, name) != -1) {
			continue
		}

		builder.WriteString(newline)

		var value interface{}
		ok := false
		if defaults != nil {
			key := field.Name
			if field.IsEmbedded {
				key = field.Type.String()
				if pathEnd := strings.LastIndex(key, "/"); pathEnd >= 0 {
					key = key[pathEnd+strings.Index(key[pathEnd+1:], ".")+2:]
				}
			}

			value, ok = defaults.(map[string]interface{})[key]
		}

		if jsonName, ok := field.Tags["json"]; ok {
			name = jsonName
		}

		name = y.renderKey(name)

		consts := distiller.LookupTypedConsts(field.Type.String())

		renderType := y.docTypesMode == AllFields

		simple := y.isSimpleField(field)

		fieldIndent := indent

		parent := y.path
		if len(y.path) > 0 {
			y.path += "."
		}

		// No default defined for this field, if named (struct) or array will be rendered below.
		_, isNamed := field.Type.(*types.Named)
		if !ok && field.Layout == distiller.LayoutSingle && (consts != nil || !isNamed) {
			if consts != nil {
				value = consts[0].Value
			} else {
				value = typeZero(field)
				var basicT *types.Basic
				basicT, ok = field.Type.(*types.Basic)
				if ok && basicT.Kind() == types.String {
					var err error
					value, err = y.renderString(value)
					if err != nil {
						return "", err
					}
				}
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

					structIndent := fieldIndent
					if !field.IsEmbedded {
						structIndent += y.indent
					}

					value, err = y.RenderStruct(subInfo, value, structIndent, field.IsEmbedded, shadowing[i:])

					if err != nil {
						return "", err
					}
				} else {
					// No special handling required for basic types.
					var basicT *types.Basic
					basicT, ok = field.Type.(*types.Basic)
					if ok && basicT.Kind() == types.String {
						value, err = y.renderString(value)
						if err != nil {
							return "", err
						}
					}
					renderType = renderType || (y.docTypesMode == BasicFields)
				}

			case distiller.LayoutArray:
				if value == nil {
					// Add an example item in case of nil array.
					value, err = y.RenderArray(field, []interface{}{nil}, fieldIndent)
				} else {
					value, err = y.RenderArray(field, value.([]interface{}), fieldIndent)
				}

			case distiller.LayoutMap:
				value, err = y.RenderMap(field, value.(*ordered.Map), fieldIndent)
			}

			if err != nil {
				return "", err
			}
		}

		if field.IsEmbedded {
			builder.WriteString(fmt.Sprintf("%v", value))
		} else {
			docIndent := fieldIndent
			if y.inArray {
				if i == 0 {
					fieldIndent += fmt.Sprintf("%-*s", len(y.indent), "-")
				} else {
					fieldIndent += y.indent
					docIndent += y.indent
				}
			}

			doc := renderDoc(field, docIndent, "#", renderType)
			builder.WriteString(doc)

			if simple {
				builder.WriteString(fmt.Sprintf("%s%s: %v", fieldIndent, name, value))
			} else {
				builder.WriteString(fmt.Sprintf("%s%s:\n%v", fieldIndent, name, value))
			}
		}

		y.path = parent

		newline = "\n"
	}

	return builder.String(), nil
}

func (y *Yaml) RenderArray(field *distiller.FieldInfo, value []interface{}, indent string) (string, error) {
	if len(value) == 0 {
		return "", nil
	}

	code := ""
	separator := "\n"
	simple := y.isSimpleType(field.EltType)

	// Array character "-" rendering for non-simple elements
	// is handled inside specific renderers to correctly align documentation.
	if simple {
		indent += y.indent + fmt.Sprintf("%-*s", len(y.indent), "-")
	}

	y.inArray = true
	for _, elt := range value {
		literal, err := y.RenderElement(field.EltType, elt, indent+y.indent)
		if err != nil {
			return "", err
		}

		if simple {
			code += indent + literal + separator
		} else {
			code += literal + separator
		}
	}
	y.inArray = false

	return code[:len(code)-1], nil
}

func (y *Yaml) RenderMap(field *distiller.FieldInfo, value *ordered.Map, indent string) (string, error) {
	if field.IsEmbedded {
		return "", fmt.Errorf("field of slice or map type cannot be embedded")
	}

	if value.Len() == 0 {
		return "\n", nil
	}

	eltsIndent := indent + y.indent
	code := ""

	var err error
	value.Iterate(func(key string, elt interface{}) bool {
		var literal string
		literal, err = y.RenderElement(field.EltType, elt, eltsIndent)
		if err != nil {
			return false
		}

		code += eltsIndent + fmt.Sprintf("%s: %s", y.renderKey(key), literal) + "\n"
		return true
	})

	if err != nil {
		return "", err
	}

	return code[:len(code)-1], nil
}

func (y *Yaml) RenderElement(itemType types.Type, item interface{}, indent string) (string, error) {
	basicT, ok := itemType.(*types.Basic)
	if ok || distiller.LookupTypedConsts(itemType.String()) != nil {
		if basicT.Kind() == types.String {
			return y.renderString(fmt.Sprintf("%v", item))
		}
		return fmt.Sprintf("%v", item), nil
	}

	subInfo := distiller.LookupStruct(itemType.String())
	if subInfo == nil {
		return "", fmt.Errorf("cannot lookup structure %s", itemType.String())
	}

	return y.RenderStruct(subInfo, item, indent, false, nil)
}

// isSimpleField verifies that a field is simple, i.e. that it is of native type, not arrays/slices or maps.
func (y *Yaml) isSimpleField(field *distiller.FieldInfo) bool {
	_, isNamed := field.Type.(*types.Named)
	consts := distiller.LookupTypedConsts(field.Type.String())
	switch field.Layout {
	case distiller.LayoutSingle:
		return consts != nil || !isNamed
	}

	return false
}

// isSimpleType verifies that a type is simple, i.e. that it is of native type, not arrays/slices or maps.
func (y *Yaml) isSimpleType(t types.Type) bool {
	_, ok := t.(*types.Basic)
	return ok || distiller.LookupTypedConsts(t.String()) != nil
}

// renderKey renders a Yaml key surrounding it with quotes when needed.
func (y *Yaml) renderKey(key string) string {
	key = strings.TrimSuffix(strings.TrimPrefix(key, "\""), "\"")

	if match, _ := regexp.MatchString(`[^\w-]+`, key); match {
		return "'" + strings.Trim(key, "\"") + "'"
	}

	// bare key.
	return key
}

// renderString renders a Yaml string with the simpler format allowed by the content.
func (y *Yaml) renderString(v interface{}) (string, error) {
	s := fmt.Sprintf("%v", v)

	err := json.Unmarshal([]byte(s), &s)
	if err != nil {
		return "", err
	}

	var needQuotes bool
	needQuotes, err = regexp.MatchString(`[\t\n\r"']+`, s)
	if err != nil {
		return "", err
	}

	if needQuotes {
		return fmt.Sprintf("%+q", s), nil
	}

	return "'" + s + "'", nil
}
