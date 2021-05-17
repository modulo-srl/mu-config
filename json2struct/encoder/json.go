package encoder

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/modulo-srl/mu-config/jsonc"
	"gitlab.com/c0b/go-ordered-json"
)

// jsonToMap converts JSON/JSONC file to generic map[string]interface{},
// preserving fields order.
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

	if filepath.Ext(filename) == ".jsonc" {
		//raw = stripComments(raw)
		raw = jsonc.ToJSON(raw)
	}

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

// jsonKeyToGoField converts a JSON property name to a valid Go struct field name
// ("json-field" => "JsonField").
func jsonKeyToFieldName(key string) string {
	name := ""

	nextUpper := false
	for i, char := range key {
		switch {
		case char >= 'a' && char <= 'z':
			if nextUpper || i == 0 {
				name += string(char - 32)
			} else {
				name += string(char)
			}
			nextUpper = false

		case char >= 'A' && char <= 'Z':
			name += string(char)
			nextUpper = false

		case char >= '0' && char <= '9':
			if i == 0 {
				name += "p"
			}

			name += string(char)
			nextUpper = false

		default:
			nextUpper = true
		}
	}

	return name
}

// jsonKeyToTag return the struct field tag.
func jsonKeyToTag(key string) string {
	return "`json:\"" + key + "\"`"
}

// jsonNumberToGoType returns the Go type (int/int64/float64) suitable for a JSON number value.
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

// jsonValueType detects the field type of a JSON value.
func jsonValueType(value interface{}, force64 bool) (fieldType, string) {
	valueType := reflect.TypeOf(value)

	if valueType.String() == "*ordered.OrderedMap" {
		return typeMap, "map"
	}

	switch valueType.Kind() {
	case reflect.Map:
		return typeMap, "map"

	case reflect.String:
		if valueType.String() == "json.Number" {
			ntype := jsonNumberToGoType(value.(json.Number).String(), force64)
			if ntype[:3] == "int" {
				return typeInt, ntype
			}
			return typeFloat, "float64"
		}

		return typeString, "string"

	case reflect.Bool:
		return typeBool, "bool"

	case reflect.Array:
	case reflect.Slice:
		return typeArray, "[]"
	}

	return typeUnknown, ""
}

// jsonItemToField convert JSON key,value to dataField.
func jsonItemToField(key string, value interface{}, structName string, force64 bool) (dataField, error) {
	typ, typeS := jsonValueType(value, force64)

	field := dataField{
		name:  jsonKeyToFieldName(key),
		tag:   jsonKeyToTag(key),
		typ:   typ,
		typeS: typeS,
	}

	switch field.typ {

	case typeInt, typeFloat, typeString, typeBool:
		field.value = castJsonValue(value, field.typ)
		return field, nil

	case typeMap:
		obj := value.(*ordered.OrderedMap)

		_, ok := obj.EntriesIter()()
		if !ok {
			return field, fmt.Errorf("\"%s\" empty object", key)
		}

		subName := structName + "_" + jsonKeyToFieldName(key)
		field.typeS = subName

		return field, nil

	case typeArray:
		arr := value.([]interface{})

		if len(arr) == 0 {
			return field, fmt.Errorf("\"%s\" empty array", key)
		}

		itemType, itemTypeS, err := detectJsonArrayType(arr, force64)
		if err != nil {
			return field, fmt.Errorf("\"%s\" %s", key, err)
		}

		if itemType == typeMap {
			itemName := structName + "_" + jsonKeyToFieldName(key) + "Item"
			field.typeS = "[]" + itemName
		} else {
			field.typeS = "[]" + itemTypeS
			field.value = "[]" + itemTypeS + "{"
			for i, v := range arr {
				if i > 0 {
					field.value += ", "
				}
				field.value += castJsonValue(v, itemType)
			}
			field.value += "}"
		}

		return field, nil
	}

	return field, fmt.Errorf("\"%s\" type unknown", key)
}

// castJsonValue convert JSON value to kwnown int/float/string/bool field type value.
func castJsonValue(value interface{}, typ fieldType) string {
	switch typ {
	case typeInt:
		return value.(json.Number).String()

	case typeFloat:
		return value.(json.Number).String()

	case typeString:
		return escapeString(value.(string))

	case typeBool:
		return fmt.Sprintf("%v", value)
	}

	return ""
}

// detectJsonArrayType validates and detect array items type.
func detectJsonArrayType(arr []interface{}, force64 bool) (typ fieldType, typeS string, err error) {
	var arrayType fieldType
	var arrayTypeS string

	for i, value := range arr {
		t, ts := jsonValueType(value, force64)

		if i == 0 {
			arrayType = t
			arrayTypeS = ts
		} else {
			if t != arrayType {
				return typeUnknown, "", errors.New("array items type mismatch")
			}
		}
	}

	if arrayType == typeArray {
		return typeArray, "", errors.New("array of array unimplemented")
	}

	return arrayType, arrayTypeS, nil
}

// escapeString escapes a string for Go.
func escapeString(s string) string {
	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	return string(b)
}

// stripComments removes all comments using pattern " //.......<EOL>"
/*
func stripComments(stream []byte) []byte {
	reg := regexp.MustCompile(`\s//.*`)
	res := reg.ReplaceAll(stream, nil)
	return res
}
*/
