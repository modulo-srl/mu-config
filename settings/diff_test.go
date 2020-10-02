package settings

import (
	"encoding/json"
	"testing"
)

func TestDiff(t *testing.T) {
	var err error

	json1 := `{ "a": "1", "b": { "c": 2, "d": true,  "e": 2.0, "f": { "g": { "h": 100, "i": 102 }, "j": true, "k": { "l": 1 } } } }`
	json2 := `{ "a": "",  "b": { "c": 2, "d": false, "f": { "g": { "h": 101 } } } }`
	jsonExpected := `{"a":"","b":{"d":false,"f":{"g":{"h":101}}}}`

	var map1, map2 map[string]interface{}

	err = json.Unmarshal([]byte(json1), &map1)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal([]byte(json2), &map2)
	if err != nil {
		t.Fatal(err)
	}

	var mapDiff map[string]interface{}

	mapDiff = diffMaps(map1, map2)

	j, err := json.Marshal(&mapDiff)
	if err != nil {
		t.Fatal(err)
	}

	if string(j) != jsonExpected {
		t.Fatal("Mismatch output\nDiff:     " + string(j) + "\nExpected: " + jsonExpected + "\n")
	}
}
