package parsers

import (
	"bytes"
	"encoding/json"
	"os"
)

func LoadJsonFile(filename string, data interface{}) error {
	bb, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return LoadJson(bb, data)
}

func LoadJson(bb []byte, data interface{}) error {
	r := bytes.NewReader(bb)

	d := json.NewDecoder(r)
	d.DisallowUnknownFields()

	err := d.Decode(data)
	if err != nil {
		return err
	}

	return nil
}

func SaveJsonFile(filename string, data interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := SaveJson(data)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, b, 0666)

	return err
}

func SaveJson(data interface{}) ([]byte, error) {
	bb, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return nil, err
	}

	return bb, nil
}
