package settings

import (
	"io/ioutil"
	"strings"
	"testing"
)

const testSettingsFilename = "settings_test.json"

type testSettingsData struct {
	Foo int
	Bar string
}

var testData testSettingsData
var testDefaultData testSettingsData

var testSet *Settings

func initTests(t *testing.T) error {
	testData = testSettingsData{}
	testDefaultData = testSettingsData{Foo: 1, Bar: "bar"}

	var err error
	testSet, err = New(testSettingsFilename, &testData, testDefaultData, false)

	if err != nil {
		if t != nil {
			t.Error("Error init settings:", err)
		}
		return err
	}

	return nil
}

func TestSave(t *testing.T) {
	err := initTests(t)
	if err != nil {
		return
	}

	err = testSet.SaveSettings()
	if err != nil {
		t.Error("Error saving settings:", err)
	}
}

func TestLoad(t *testing.T) {
	err := initTests(t)
	if err != nil {
		return
	}

	err = testSet.LoadSettings(true)
	if err != nil {
		t.Error("Error loading settings:", err)
	}
}

func TestLoadSave(t *testing.T) {
	err := initTests(t)
	if err != nil {
		return
	}

	testData = testSettingsData{
		Foo: -1,
		Bar: "mh",
	}

	err = testSet.SaveSettings()
	if err != nil {
		if t != nil {
			t.Error("Error saving settings:", err)
		}
		return
	}

	testData2 := testData
	testData = testSettingsData{}

	err = testSet.LoadSettings(true)
	if err != nil {
		if t != nil {
			t.Error("Error loading settings:", err)
		}
		return
	}

	if testData != testData2 {
		if t != nil {
			t.Error("Loaded data mismatch: ", testData2, testData)
		}
		return
	}
}

func TestDefault(t *testing.T) {
	readfile := func() string {
		content, _ := ioutil.ReadFile(testSet.GetFilename())
		return string(content)
	}

	err := initTests(t)
	if err != nil {
		return
	}

	testData = testSettingsData{Foo: 1, Bar: "bar"}
	testSet.SaveSettings()

	content := readfile()
	if content != "{}" {
		t.Error("Mismatch 1", content)
	}

	testData = testSettingsData{Foo: 2, Bar: "bar2"}
	testSet.SaveSettings()

	content = strings.ReplaceAll(readfile(), " ", "")
	content = strings.ReplaceAll(content, "\n", "")
	content = strings.ReplaceAll(content, "\t", "")
	if content != `{"Foo":2,"Bar":"bar2"}` {
		t.Error("Mismatch 2", content)
	}
}

func BenchmarkLoadSave(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TestLoadSave(nil)
	}
}
