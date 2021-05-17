package encoder

import (
	"fmt"
	"io/ioutil"
	"os"
)

type dataField struct {
	name string

	typ   fieldType
	typeS string

	tag string

	value string
}

type fieldType int

const (
	typeUnknown fieldType = iota
	typeString
	typeBool
	typeInt
	typeFloat
	typeArray
	typeMap
)

// JSONToStructs validates and dumps the Go structs declarations of a JSON file.
func JSONToStructs(filename string, rootStructName string, force64 bool) (string, error) {
	m, err := jsonToMap(filename)
	if err != nil {
		return "", err
	}

	out, err := mapToStructs(m, rootStructName, force64)

	return out, err
}

// JSONToValues render an init function that returns the values of a JSON file.
func JSONToValues(filename string, funcName, rootStructName string, force64 bool) (string, error) {
	m, err := jsonToMap(filename)
	if err != nil {
		return "", err
	}

	out, err := mapToValues(m, funcName, rootStructName, force64)

	return out, err
}

// IncludeRaw embeds original file content.
func IncludeRaw(filename, constName string) (string, error) {
	out := "const " + constName + " = \""

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	for _, c := range content {
		out += fmt.Sprintf("\\x%02x", c)
	}

	out += "\"\n\n"

	return out, nil
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
