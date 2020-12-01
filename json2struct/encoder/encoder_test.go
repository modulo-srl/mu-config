package encoder

import (
	"encoding/hex"
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

	out, err := JSONToStruct(testDataPath+"settings.default.jsonc", "testData", false)
	if err != nil {
		t.Error(err)
		return
	}

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

	// Serve lo stesso, dato che json2struct genera sempre e comunque la definizione della struttura.
	out1, err := JSONToStruct(testDataPath+"settings.default.jsonc", "testData", false)
	if err != nil {
		t.Error(err)
		return
	}

	out2, err := JSONToValues(testDataPath+"settings.default.jsonc", "init", "testData", false)
	if err != nil {
		t.Error(err)
		return
	}

	if strings.TrimSpace(out1+out2) != strings.TrimSpace(compare) {
		t.Error("Output and compare mismatch, see " + failFilename)
		fmt.Println("********************* Compare: \n", compare)
		fmt.Println("********************* Generated: \n", out1+out2)

		WriteFile(failFilename, out1+out2)
	}
}

func TestRaw(t *testing.T) {
	raw, err := IncludeRaw(testDataPath+"settings.default.jsonc", "raw")
	if err != nil {
		t.Error(err)
		return
	}
	raw = strings.TrimSpace(raw)
	raw = raw[13 : len(raw)-1]
	raw = strings.ReplaceAll(raw, `\x`, "")

	rawB, err := hex.DecodeString(raw)
	if err != nil {
		t.Error(err)
		return
	}

	compare, err := loadFile(t, testDataPath+"settings.default.jsonc")
	if err != nil {
		return
	}

	if string(rawB) != compare {
		t.Error("Output and compare mismatch")
		fmt.Println("********************* Compare: \n", compare)
		fmt.Println("********************* Generated: \n", string(rawB))
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
