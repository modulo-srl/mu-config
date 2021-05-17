package encoder

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"gitlab.com/c0b/go-ordered-json"
)

const testDataPath = "../test-data/"

func TestInvalid(t *testing.T) {
	m, err := jsonToMap(testDataPath + "should-fail.jsonc")
	if err != nil {
		t.Error(err)
		return
	}

	iter := m.EntriesIter()
	for {
		pair, ok := iter()
		if !ok {
			break
		}
		key := pair.Key
		value := pair.Value

		fmt.Printf("Test %s ... ", key)

		_, err := mapToStructs(value.(*ordered.OrderedMap), "test", false)
		if err == nil {
			fmt.Println("FAILED")
			t.Fail()
		} else {
			fmt.Printf("OK (%s)\n", err.Error())
		}
	}
}

func TestOutput(t *testing.T) {
	failFilename := testDataPath + "test-fail.txt"
	os.Remove(failFilename)

	compare, err := loadFile(t, testDataPath+"should-ok.out")
	if err != nil {
		t.Error(err)
		return
	}

	out, err := JSONToStructs(testDataPath+"../example.jsonc", "test", false)
	if err != nil {
		t.Error(err)
		return
	}
	s, err := JSONToValues(testDataPath+"../example.jsonc", "default", "test", false)
	if err != nil {
		t.Error(err)
		return
	}
	out += s

	if strings.TrimSpace(out) != strings.TrimSpace(compare) {
		t.Error("Output and compare mismatch. For current (failed) output, see " + failFilename)
		//fmt.Println("********************* Compare: \n", compare)
		//fmt.Println("********************* Generated: \n", out)

		WriteFile(failFilename, out)
	}
}

func TestRaw(t *testing.T) {
	raw, err := IncludeRaw(testDataPath+"../example.jsonc", "raw")
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

	compare, err := loadFile(t, testDataPath+"../example.jsonc")
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
