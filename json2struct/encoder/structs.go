package encoder

import (
	"fmt"

	"gitlab.com/c0b/go-ordered-json"
)

type dataStruct struct {
	name   string
	fields []dataField
}

// mapToStructs converts map to structs declaration.
func mapToStructs(m *ordered.OrderedMap, rootStructName string, force64 bool) (string, error) {
	dataStructs := make([]*dataStruct, 0, 32)

	err := mapToStructRecursive(m, jsonKeyToFieldName(rootStructName), force64, true, &dataStructs)
	if err != nil {
		return "", err
	}

	// Print the structs declaration.

	out := ""

	for _, stru := range dataStructs {
		out += "type " + stru.name + " struct {\n"

		for _, field := range stru.fields {
			out += "\t" + field.name + "\t" + field.typeS
			if field.tag != "" {
				out += "\t" + field.tag
			}
			out += "\n"
		}

		out += "}\n\n"
	}

	return out, nil
}

func mapToStructRecursive(m *ordered.OrderedMap, structName string, force64 bool, declare bool, structs *[]*dataStruct) error {
	dataItem := dataStruct{
		name:   structName,
		fields: make([]dataField, 0, 32),
	}

	if declare {
		*structs = append(*structs, &dataItem)
	}

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
			return err
		}

		switch field.typ {
		case typeMap:
			err := mapToStructRecursive(value.(*ordered.OrderedMap), field.typeS, force64, true, structs)
			if err != nil {
				return err
			}

		case typeArray:
			v := value.([]interface{})[0]
			t, _ := jsonValueType(v, force64)

			if t != typeMap {
				break
			}

			itemName := field.typeS[2:]

			itemStruct, err := jsonArrayToDataStruct(value.([]interface{}), itemName, force64)
			if err != nil {
				return fmt.Errorf("\"%s\" %s", key, err.Error())
			}
			*structs = append(*structs, itemStruct)

			for _, item := range value.([]interface{}) {
				err := mapToStructRecursive(item.(*ordered.OrderedMap), itemName, force64, false, structs)
				if err != nil {
					return err
				}
			}
		}

		dataItem.fields = append(dataItem.fields, field)
	}

	return nil
}

// mapsToStruct make a dataStruct by merging the fields of a list of JSON objects.
func jsonArrayToDataStruct(arr []interface{}, structName string, force64 bool) (*dataStruct, error) {
	dataItem := dataStruct{
		name:   structName,
		fields: make([]dataField, 0, 32),
	}

	fmap := make(map[string]*dataField)

	for _, item := range arr {
		iter := item.(*ordered.OrderedMap).EntriesIter()
		for {
			pair, ok := iter()
			if !ok {
				break
			}
			key := pair.Key
			value := pair.Value

			field, err := jsonItemToField(key, value, structName, force64)
			if err != nil {
				return nil, err
			}

			entry, ok := fmap[field.name]
			if ok {
				if entry.typeS != field.typeS {
					return nil, fmt.Errorf("object field \"%s\" type collision", key)
				}
			} else {
				fmap[field.name] = &field
			}
		}
	}

	for _, f := range fmap {
		dataItem.fields = append(dataItem.fields, *f)
	}

	return &dataItem, nil
}
