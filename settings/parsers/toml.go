package parsers

import (
	"bytes"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func LoadTomlFile(filename string, data interface{}) error {
	bb, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return LoadToml(bb, data)
}

func LoadToml(bb []byte, data interface{}) error {
	r := bytes.NewReader(bb)

	d := toml.NewDecoder(r)
	d.DisallowUnknownFields()

	err := d.Decode(data)
	if err != nil {
		return err
	}

	return nil
}

func SaveTomlFile(filename string, data interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := SaveToml(data)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, b, 0666)

	return err
}

func SaveToml(data interface{}) ([]byte, error) {
	bb, err := toml.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bb, nil
}
