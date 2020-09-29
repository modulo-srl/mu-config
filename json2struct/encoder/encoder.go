package encoder

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"gitlab.com/c0b/go-ordered-json"
)

func JSONToStruct(filename string, rootStructName string, force64 bool) (string, error) {
	m, err := jsonToMap(filename)
	if err != nil {
		return "", err
	}

	return mapToStruct(m, rootStructName, force64), nil
}

func JSONToValues(filename string, funcName, rootStructName string, force64 bool) (string, error) {
	m, err := jsonToMap(filename)
	if err != nil {
		return "", err
	}

	return mapToValues(m, funcName, rootStructName, force64), nil
}

// WriteFile writes string to file.
func WriteFile(filename, content string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	_, err = f.WriteString(content)
	if err != nil {
		f.Close()
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

// converts JSON file (with comments) to generic map[string]interface{}, but preserving fields order.
func jsonToMap(filename string) (*ordered.OrderedMap, error) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if len(raw) == 0 {
		return nil, errors.New("empty JSON")
	}

	if raw[0] == '[' {
		return nil, errors.New("JSON does not starts with '{'")
	}

	raw = stripComments(raw)

	m := ordered.NewOrderedMap()
	//m := make(map[string]interface{})

	buf := bytes.NewReader(raw)
	dec := json.NewDecoder(buf)
	dec.UseNumber() // numbers as json.Number (the original string representation of the number)

	err = dec.Decode(&m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// mapToStruct converts map to (annidated) struct declaration.
func mapToStruct(m *ordered.OrderedMap, rootStructName string, force64 bool) string {
	out := mapToStructRecursive(m, rootStructName, 0, force64) + "\n"

	return out
}

// mapToValues converts map to struct values initialization.
func mapToValues(m *ordered.OrderedMap, funcName, rootStructName string, force64 bool) string {
	if funcName == "" {
		funcName = "init" + strings.Title(rootStructName)
	}

	rootStructName, _ = jsonToGoField(rootStructName)

	out := "func " + funcName + "() " + rootStructName + " {\n" +
		"\tvar s " + rootStructName + "\n\n"

	out += mapToValuesRecursive(m, "s", 0, force64) +
		"\n\treturn s\n" +
		"}\n\n"

	return out
}

// Remove all comments using pattern " //.......<EOL>"
func stripComments(stream []byte) []byte {
	reg := regexp.MustCompile(`\s//.*`)
	res := reg.ReplaceAll(stream, nil)
	return res
}

func jsonNumberToGoType(num string, force64 bool) string {
	var err error

	_, err = strconv.ParseInt(num, 10, 64)
	if err == nil {
		if force64 {
			return "int64"
		} else {
			return "int"
		}
	}

	_, err = strconv.ParseFloat(num, 64)
	if err == nil {
		return "float64"
	}

	return ""
}

func jsonToGoField(jsonField string) (goField, goTag string) {
	if unicode.IsLower([]rune(jsonField)[0]) {
		goField = strings.Title(jsonField)
	} else {
		goField = jsonField
	}

	goField = strings.ReplaceAll(goField, "-", "_")

	if jsonField != goField {
		goTag = "\t`json:\"" + jsonField + "\"`"
	}

	return
}

func escapeString(s string) string {
	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func jsonValueToGoTypeValue(value interface{}, force64, forceAutoInt bool) (goType, goValue string) {
	valueType := reflect.TypeOf(value)

	if valueType.String() == "*ordered.OrderedMap" {
		return "omap", ""
	}

	switch valueType.Kind() {

	//case reflect.Map:
	//	return "map", ""

	case reflect.String:
		if valueType.String() == "json.Number" {
			// int, float
			s := value.(json.Number).String()
			return jsonNumberToGoType(s, force64), s
		}

		if forceAutoInt {
			// check for int, float
			t := jsonNumberToGoType(value.(string), force64)
			if t != "" {
				return t, value.(string)
			}
		}

		return "string", escapeString(value.(string))

	case reflect.Bool:
		return "bool", fmt.Sprintf("%v", value)

	case reflect.Array:
	case reflect.Slice:
		arr := value.([]interface{})
		if len(arr) == 0 {
			return "[]string", "[]string{}"
		}

		values := ""
		valuesS := ""
		arrayType := ""
		oldItemType := ""
		mixedTypes := false

		for k, v := range value.([]interface{}) {
			itemType, itemValue := jsonValueToGoTypeValue(v, force64, true)
			if itemType != oldItemType && k > 0 {
				switch {
				// Mixed int and float, fallback to float
				case oldItemType == "int" && (itemType == "float64" || itemType == "float32"):
				case (oldItemType == "float64" || oldItemType == "float32") && itemType == "int":
					itemType = oldItemType
				// Mixed types, fallback to string
				default:
					mixedTypes = true
				}
			}
			oldItemType = itemType
			arrayType = itemType

			if k > 0 {
				values += ", "
				valuesS += ", "
			}
			values += itemValue
			if itemType == "string" {
				valuesS += itemValue
			} else {
				valuesS += escapeString(itemValue)
			}
		}

		if mixedTypes {
			return "[]string", "[]string{" + valuesS + "}"
		} else {
			return "[]" + arrayType, "[]" + arrayType + "{" + values + "}"
		}

	}

	return "", ""
}

func mapToStructRecursive(m *ordered.OrderedMap, rootName string, depth int, force64 bool) string {
	indent := strings.Repeat("\t", depth)

	rootField, rootTag := jsonToGoField(rootName)
	if depth == 0 {
		rootTag = ""
	}

	out := indent
	if depth == 0 {
		out += "type "
	}
	out += rootField + " struct {\n"

	wroteField := false
	wroteStruct := false

	//for key, value := range m {
	iter := m.EntriesIter()
	for {
		pair, ok := iter()
		if !ok {
			break
		}
		key := pair.Key
		value := pair.Value

		fieldType, _ := jsonValueToGoTypeValue(value, force64, false)

		if fieldType == "omap" {
			if wroteField || wroteStruct {
				out += "\n"
				wroteField = false
				wroteStruct = false
			}
			//out += mapToStructRecursive(value.(map[string]interface{}), key, depth+1, force64)
			out += mapToStructRecursive(value.(*ordered.OrderedMap), key, depth+1, force64)
			wroteStruct = true
			continue
		}

		fieldName, fieldTag := jsonToGoField(key)
		if wroteStruct {
			out += "\n"
			wroteStruct = false
		}
		out += indent + "\t" + fieldName + "\t" + fieldType + fieldTag + "\n"
		wroteField = true
	}

	out += indent + "}" + rootTag + "\n"

	return out
}

func mapToValuesRecursive(m *ordered.OrderedMap, fullItemName string, depth int, force64 bool) string {
	out := ""
	indent := strings.Repeat("\t", depth)

	/*if depth == 0 {
		fullItemName, _ = jsonToGoField(fullItemName)
	}*/

	wrote := false

	//for key, value := range m {
	iter := m.EntriesIter()
	for {
		pair, ok := iter()
		if !ok {
			break
		}
		key := pair.Key
		value := pair.Value

		fieldName, _ := jsonToGoField(key)
		fieldType, fieldValue := jsonValueToGoTypeValue(value, force64, false)

		if fieldType == "omap" {
			if wrote {
				out += "\n"
			}
			//out += mapToValuesRecursive(value.(map[string]interface{}), fullItemName+"."+fieldName, depth, force64)
			out += mapToValuesRecursive(value.(*ordered.OrderedMap), fullItemName+"."+fieldName, depth, force64)
			wrote = true
			continue
		}

		out += indent + "\t" + fullItemName + "." + fieldName + " = " + fieldValue + "\n"
		wrote = true
	}

	return out
}
