package parsers

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Nota: Yaml è case insensitive ma il decoder si comporta male: https://github.com/go-yaml/yaml/issues/123
// Si è preferito comunque uscire leggermente dalla RFC
// permettendo il case insensitive nei nomi variabili struct,
// in modo da essere allineati con il comportamento dei decoder di Json e Toml.

func LoadYamlFile(filename string, data interface{}) error {
	bb, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return LoadYaml(bb, data)
}

func LoadYaml(bb []byte, data interface{}) error {
	r := bytes.NewReader(lowerKeys(bb))

	d := yaml.NewDecoder(r)
	d.KnownFields(true)

	err := d.Decode(data)
	if err != nil {
		return err
	}

	return nil
}

func SaveYamlFile(filename string, data interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := SaveYaml(data)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, b, 0666)

	return err
}

func SaveYaml(data interface{}) ([]byte, error) {
	bb, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bb, nil
}

// Converte in lowercase tutti i nomi variabili delle strutture,
// favorendo un comportamento case insensitive in fase di unmarshal.
func lowerKeys(in []byte) []byte {
	lines := []string{}

	for _, line := range strings.Split(string(in), "\n") {
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			line = fmt.Sprintf("%s:%s", strings.ToLower(parts[0]), parts[1])
		}
		lines = append(lines, line)
	}

	return []byte(strings.Join(lines, "\n"))
}
