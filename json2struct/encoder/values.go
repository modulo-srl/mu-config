package encoder

import (
	"strings"

	"gitlab.com/c0b/go-ordered-json"
)

// mapToValues converts map to struct values initialization.
func mapToValues(m *ordered.OrderedMap, funcName, rootStructName string, force64 bool) (string, error) {
	if funcName == "" {
		funcName = "init"
	}
	funcName += strings.Title(rootStructName)

	rootStructName = jsonKeyToFieldName(rootStructName)

	out := "func " + funcName + "() " + rootStructName + " {\n" +
		"\treturn " + rootStructName + "{\n"

	s, err := mapToValuesRecursive(m, rootStructName, force64, 1)
	if err != nil {
		return "", err

	}

	out += s +
		"\t}\n" +
		"}\n\n"

	return out, nil
}

func mapToValuesRecursive(m *ordered.OrderedMap, structName string, force64 bool, depth int) (string, error) {
	out := ""
	indent := strings.Repeat("\t", depth)

	iter := m.EntriesIter()
	for {
		pair, ok := iter()
		if !ok {
			break
		}
		key := pair.Key
		value := pair.Value

		field, err := jsonItemToField(key, value, structName, force64)
		if err != nil {
			return "", err
		}

		switch field.typ {
		case typeMap:
			out += "\n" + indent + "\t" + field.name + ": "

			out += field.typeS + "{\n"

			s, err := mapToValuesRecursive(value.(*ordered.OrderedMap), field.typeS, force64, depth+1)
			if err != nil {
				return "", err
			}

			out += s + indent + "\t}"

		case typeArray:
			out += indent + "\t" + field.name + ": "

			v := value.([]interface{})[0]
			t, _ := jsonValueType(v, force64)

			if t != typeMap {
				out += field.value
				break
			}

			out += field.typeS + "{\n"

			itemName := field.typeS[2:]

			for _, item := range value.([]interface{}) {
				out += indent + "\t\t{\n"

				s, err := mapToValuesRecursive(item.(*ordered.OrderedMap), itemName, force64, depth+2)
				if err != nil {
					return "", err
				}

				out += s + indent + "\t\t},\n"
			}

			out += indent + "\t}"

		default:
			out += indent + "\t" + field.name + ": "

			out += field.value
		}

		out += ",\n"
	}

	return out, nil
}
