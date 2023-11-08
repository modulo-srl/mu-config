package settings

import (
	"encoding/json"

	"gitlab.com/c0b/go-ordered-json"
)

// Effettua la differenza tra due entità dati ritornando i soli campi con valori differenti (out = child - parent).
// Campi con lo stesso nome nelle due struct devono essere dello stesso tipo.
// Ritorna una map ordinata in modo da mantenere il medesimo ordine dei campi originali.
func diff(parent, child interface{}) (diffedMap *ordered.OrderedMap, err error) {
	marshUnmarsh := func(m interface{}) (out *ordered.OrderedMap, err error) {
		b, err := json.Marshal(&m)
		if err != nil {
			return nil, err
		}

		out = ordered.NewOrderedMap()

		err = json.Unmarshal(b, out)
		if err != nil {
			return nil, err
		}

		return
	}

	var m1, m2 *ordered.OrderedMap

	m1, err = marshUnmarsh(&parent)
	if err != nil {
		return nil, err
	}
	m2, err = marshUnmarsh(&child)
	if err != nil {
		return nil, err
	}

	return diffMaps(m1, m2), nil
}

func diffMaps(mapParent, mapChild *ordered.OrderedMap) *ordered.OrderedMap {
	mapOut := ordered.NewOrderedMap()

	iter := mapParent.EntriesIter()
	for {
		pair, ok := iter()
		if !ok {
			break
		}
		k := pair.Key
		v1 := pair.Value

		v2, ok := mapChild.GetValue(k)
		if !ok {
			continue
		}

		vd := diffFields(v1, v2)
		if vd != nil {
			mapOut.Set(k, vd)
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
	if m1, ok := field1.(*ordered.OrderedMap); ok {
		m2, ok := field2.(*ordered.OrderedMap)
		if !ok || m2 == nil {
			return nil
		}

		dm := diffMaps(m1, m2)
		i := dm.EntriesIter()
		_, ok = i()
		if !ok {
			// empty map
			return nil
		}

		return dm
	}

	// Array
	if a1, ok := field1.([]interface{}); ok {
		a2, ok := field2.([]interface{})
		if !ok {
			return nil
		}

		i := diffArrays(a1, a2)
		if i == nil {
			return nil // evita il cast a slice vuoto
		}
		return i
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

	// JSON Number
	if b1, ok := field1.(json.Number); ok {
		b2, ok := field2.(json.Number)
		if !ok || b1 == b2 {
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

/*func cloneData(from, to interface{}) error {
	bout, err := json.Marshal(&from)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bout, &to)
	return err
}*/
