package encoder

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const testDataPath = "../test_data/"

func TestDefine(t *testing.T) {
	failFilename := testDataPath + "generated_struct.fail.txt"
	os.Remove(failFilename)

	compare, err := loadFile(t, testDataPath+"settings.default.def.txt")
	if err != nil {
		return
	}

	out, err := JSONToStruct(testDataPath+"settings.default.json", "testData", false)
	if err != nil {
		t.Error(err)
		return
	}

	out = "package test\n\n" + out

	if strings.TrimSpace(out) != strings.TrimSpace(compare) {
		t.Error("Output and compare mismatch, see " + failFilename)
		fmt.Println("********************* Compare: \n", compare)
		fmt.Println("********************* Generated: \n", out)

		WriteFile(failFilename, out)
	}
}

func TestValues(t *testing.T) {
	failFilename := testDataPath + "generated_values.fail.txt"
	os.Remove(failFilename)

	compare, err := loadFile(t, testDataPath+"settings.default.init.txt")
	if err != nil {
		return
	}

	out, err := JSONToValues(testDataPath+"settings.default.json", "initTestData", "testData", false)
	if err != nil {
		t.Error(err)
		return
	}

	out = "package test\n\n" + out

	if strings.TrimSpace(out) != strings.TrimSpace(compare) {
		t.Error("Output and compare mismatch, see " + failFilename)
		fmt.Println("********************* Compare: \n", compare)
		fmt.Println("********************* Generated: \n", out)

		WriteFile(failFilename, out)
	}
}

func loadFile(t *testing.T, filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		if t != nil {
			t.Error("Error loading "+filename+": ", err)
			t.Error(err)
		}
		return "", err
	}

	return string(content), nil
}
