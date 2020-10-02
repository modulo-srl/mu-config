package settings

import "encoding/json"

// DiffMaps returns a map as difference mapChild - mapParent, returning only fields with different values.
// Fields with the same name on two maps must to be of the same type.
func DiffMaps(mapParent, mapChild interface{}) (diffedMap map[string]interface{}, err error) {
	marshUnmarsh := func(m interface{}) (out map[string]interface{}, err error) {
		b, err := json.Marshal(&m)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(b, &out)
		if err != nil {
			return nil, err
		}

		return
	}

	var m1, m2 map[string]interface{}

	m1, err = marshUnmarsh(&mapParent)
	if err != nil {
		return nil, err
	}
	m2, err = marshUnmarsh(&mapChild)
	if err != nil {
		return nil, err
	}

	return diffMaps(m1, m2), nil
}

func diffMaps(mapParent, mapChild map[string]interface{}) map[string]interface{} {
	mapOut := map[string]interface{}{}

	for k, v1 := range mapParent {
		v2, ok := mapChild[k]
		if !ok {
			continue
		}

		vd := diffFields(v1, v2)
		if vd != nil {
			mapOut[k] = vd
		}
	}

	return mapOut
}

func diffArrays(arr1, arr2 []interface{}) []interface{} {
	if len(arr1) != len(arr2) {
		return arr2
	}

	for i := range arr1 {
		f := diffFields(arr1[i], arr2[i])
		if f != nil {
			return arr2
		}
	}

	return nil
}

func diffFields(field1, field2 interface{}) interface{} {
	// Map
	if m1, ok := field1.(map[string]interface{}); ok {
		m2, ok := field2.(map[string]interface{})
		if !ok || m2 == nil {
			return nil
		}

		return diffMaps(m1, m2)
	}

	// Array
	if a1, ok := field1.([]interface{}); ok {
		a2, ok := field2.([]interface{})
		if !ok {
			return nil
		}

		return diffArrays(a1, a2)
	}

	// Boolean
	if b1, ok := field1.(bool); ok {
		b2, ok := field2.(bool)
		if !ok || b1 == b2 {
			return nil
		}
		return field2
	}

	// Integer
	if i1, ok := field1.(int); ok {
		i2, ok := field2.(int)
		if !ok || i1 == i2 {
			return nil
		}
		return field2
	}

	// Float64
	if f1, ok := field1.(float64); ok {
		f2, ok := field2.(float64)
		if !ok || f1 == f2 {
			return nil
		}
		return field2
	}

	// Float32
	if f1, ok := field1.(float32); ok {
		f2, ok := field2.(float32)
		if !ok || f1 == f2 {
			return nil
		}
		return field2
	}

	// String
	if s1, ok := field1.(string); ok {
		s2, ok := field2.(string)
		if !ok || s1 == s2 {
			return nil
		}
		return field2
	}

	return nil
}

func cloneData(from, to interface{}) error {
	bout, err := json.Marshal(&from)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bout, &to)
	return err
}
