package renderers

import (
	"fmt"
	"github.com/modulo-srl/mu-config/go2cfg/distiller"
	"github.com/modulo-srl/mu-config/go2cfg/ordered"
	"github.com/pelletier/go-toml/v2"
	"go/types"
	"regexp"
	"sort"
	"strings"
)

// Toml renders Toml code from distiller info.
type Toml struct {
	docTypesMode DocTypesMode
	path         string
	inArray      bool
	indented     bool
}

// NewToml creates a new Toml renderer.
// mode controls the rendering of field types in Toml comments.
func NewToml(mode DocTypesMode, indented bool) *Toml {
	return &Toml{
		docTypesMode: mode,
		path:         "",
		inArray:      false,
		indented:     indented,
	}
}

func (t *Toml) RenderStruct(info *distiller.StructInfo, defaults interface{}, indent string, embedded bool, _ []string) (string, error) {
	var builder strings.Builder

	if !embedded {
		if len(t.path) > 0 {
			if t.inArray {
				builder.WriteString(fmt.Sprintf("[[%s]]\n", t.path))
			} else {
				builder.WriteString(fmt.Sprintf("[%s]\n", t.path))
			}
		}
	}

	sorted, sortedDefaults, err := t.sortFields(info.Fields, defaults)
	if err != nil {
		return "", err
	}

	newline := ""

	for _, field := range sorted {
		name := field.Name

		builder.WriteString(newline)

		if jsonName, ok := field.Tags["json"]; ok {
			name = t.renderKey(jsonName)
		}

		var value interface{}
		ok := false
		if sortedDefaults != nil {
			value, ok = sortedDefaults[field.Name]
		}

		consts := distiller.LookupTypedConsts(field.Type.String())

		renderType := t.docTypesMode == AllFields

		simple := t.isSimpleField(field)

		fieldIndent := indent
		if t.indented && !simple && len(t.path) > 0 {
			fieldIndent += "\t"
		}

		parent := t.path
		if len(t.path) > 0 {
			t.path += "."
		}
		t.path += t.renderKey(name)

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
					value, err = t.renderString(value)
					if err != nil {
						return "", err
					}
				}
			}
		} else {
			switch field.Layout {
			case distiller.LayoutSingle:
				if isNamed && consts == nil {
					subInfo := distiller.LookupStruct(field.Type.String())
					if subInfo == nil {
						return "", fmt.Errorf("cannot lookup structure %s", field.Type.String())
					}

					value, err = t.RenderStruct(subInfo, value, fieldIndent, field.IsEmbedded, nil)

					if err != nil {
						return "", err
					}
				} else {
					// No special handling required for basic types.
					var basicT *types.Basic
					basicT, ok = field.Type.(*types.Basic)
					if ok && basicT.Kind() == types.String {
						value, err = t.renderString(value)
						if err != nil {
							return "", err
						}
					}
					renderType = renderType || (t.docTypesMode == BasicFields)
				}

			case distiller.LayoutArray:
				if value == nil {
					// Add an example item in case of nil array.
					value, err = t.RenderArray(field, []interface{}{nil}, fieldIndent)
				} else {
					value, err = t.RenderArray(field, value.([]interface{}), fieldIndent)
				}

			case distiller.LayoutMap:
				value, err = t.RenderMap(field, value.(*ordered.Map), fieldIndent)
			}

			if err != nil {
				return "", err
			}
		}

		if field.IsEmbedded {
			builder.WriteString(fmt.Sprintf("%v", value))
		} else {
			doc := renderDoc(field, fieldIndent, "#", renderType)
			if !simple && field != sorted[0] {
				doc = "\n" + doc
			}
			builder.WriteString(doc)

			if simple {
				builder.WriteString(fmt.Sprintf("%s%s = %v", fieldIndent, name, value))
			} else if field.Layout == distiller.LayoutMap {
				builder.WriteString(fmt.Sprintf("%s[%s]\n%v", fieldIndent, t.path, value))
			} else {
				builder.WriteString(fmt.Sprintf("%s%v", fieldIndent, value))
			}
		}

		t.path = parent

		newline = "\n"
	}

	return builder.String(), nil
}

func (t *Toml) RenderArray(field *distiller.FieldInfo, value []interface{}, indent string) (string, error) {
	if len(value) == 0 {
		return "", nil
	}

	simple := t.isSimpleField(field)
	code := ""
	separator := "\n\n"
	eltsIndent := indent
	if simple {
		code = "[\n"
		separator = ",\n"
		eltsIndent += "\t"
	}

	t.inArray = true
	for _, elt := range value {
		literal, err := t.RenderElement(field.EltType, elt, eltsIndent)
		if err != nil {
			return "", err
		}

		code += eltsIndent + literal + separator
	}
	t.inArray = false

	if simple {
		code = strings.TrimSuffix(code, separator) + "\n" + indent + "]"
	}

	return code, nil
}

func (t *Toml) RenderMap(field *distiller.FieldInfo, value *ordered.Map, indent string) (string, error) {
	if field.IsEmbedded {
		return "", fmt.Errorf("field of slice or map type cannot be embedded")
	}

	if value.Len() == 0 {
		return "\n", nil
	}

	eltsIndent := indent
	code := ""

	var err error
	value.Iterate(func(key string, elt interface{}) bool {
		var literal string
		literal, err = t.RenderElement(field.EltType, elt, eltsIndent)
		if err != nil {
			return false
		}

		code += eltsIndent + fmt.Sprintf("%s = %s", t.renderKey(key), literal) + "\n"
		return true
	})

	if err != nil {
		return "", err
	}

	return code, nil
}

func (t *Toml) RenderElement(itemType types.Type, item interface{}, indent string) (string, error) {
	basicT, ok := itemType.(*types.Basic)
	if ok || distiller.LookupTypedConsts(itemType.String()) != nil {
		if basicT.Kind() == types.String {
			return t.renderString(fmt.Sprintf("%v", item))
		}
		return fmt.Sprintf("%v", item), nil
	}

	subInfo := distiller.LookupStruct(itemType.String())
	if subInfo == nil {
		return "", fmt.Errorf("cannot lookup structure %s", itemType.String())
	}

	return t.RenderStruct(subInfo, item, indent, false, nil)
}

// sortFields sorts the fields by putting those that have basic types or that are
// slices or arrays of basic types first, taking shadowing and embedding into account.
// It also returns the defaults to simplify access in case of embedding.
func (t *Toml) sortFields(fields []*distiller.FieldInfo, defaults interface{}) ([]*distiller.FieldInfo, map[string]interface{}, error) {
	var sorted []*distiller.FieldInfo
	fieldsDefaults := make(map[string]interface{})

	for _, field := range fields {
		if !field.IsEmbedded {
			if i := fieldsSlice(sorted).indexOf(field.Name); i != -1 {
				sorted = append(sorted[0:i], sorted[i+1:]...)
			}
			sorted = append(sorted, field)
			if defaults != nil {
				if value, ok := defaults.(map[string]interface{})[field.Name]; ok {
					fieldsDefaults[field.Name] = value
				}
			}
			continue
		}

		key := field.Type.String()
		if pathEnd := strings.LastIndex(key, "/"); pathEnd >= 0 {
			key = key[pathEnd+strings.Index(key[pathEnd+1:], ".")+2:]
		}

		subInfo := distiller.LookupStruct(field.Type.String())
		if subInfo == nil {
			return nil, nil, fmt.Errorf("cannot lookup structure %s", field.Type.String())
		}

		var defaultsMap interface{}
		if defaults != nil {
			defaultsMap = defaults.(map[string]interface{})[key]
		}
		subFields, subDefaults, err := t.sortFields(subInfo.Fields, defaultsMap)
		if err != nil {
			return nil, nil, err
		}

		for _, subField := range subFields {
			if i := fieldsSlice(sorted).indexOf(subField.Name); i != -1 {
				sorted = append(sorted[0:i], sorted[i+1:]...)
			}
			sorted = append(sorted, subField)

			if value, ok := subDefaults[subField.Name]; ok {
				fieldsDefaults[subField.Name] = value
			}
		}
	}

	sort.SliceStable(sorted, func(i, j int) bool {
		return t.isSimpleField(sorted[i]) && !t.isSimpleField(sorted[j])
	})

	return sorted, fieldsDefaults, nil
}

// isSimpleField verifies that a field is simple, i.e. that it is of native type or an array or slice of native types.
func (t *Toml) isSimpleField(field *distiller.FieldInfo) bool {
	_, isNamed := field.Type.(*types.Named)
	consts := distiller.LookupTypedConsts(field.Type.String())
	switch field.Layout {
	case distiller.LayoutSingle:
		return consts != nil || !isNamed

	case distiller.LayoutArray:
		_, isNamed = field.EltType.(*types.Named)
		consts = distiller.LookupTypedConsts(field.EltType.String())
		return consts != nil || !isNamed
	}

	return false
}

// renderKey renders a TOML key surrounding it with quotes when needed.
func (t *Toml) renderKey(key string) string {
	key = strings.TrimSuffix(strings.TrimPrefix(key, "\""), "\"")

	if match, _ := regexp.MatchString(`[^\w-]+`, key); match {
		return "'" + strings.Trim(key, "\"") + "'"
	}

	// bare key.
	return key
}

// renderString renders a TOML string with the simpler format allowed by the content.
func (t *Toml) renderString(v interface{}) (string, error) {
	p, err := toml.Marshal(unescapeString(v))
	return string(p), err
}
