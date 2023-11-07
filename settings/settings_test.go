package settings

import (
	"strings"
	"testing"

	"github.com/modulo-srl/mu-config/settings/parsers"
)

type MySettings struct {
	Main  settingsMain
	Users []settingsUsersItem
}

type settingsMain struct {
	ParamString string
	ParamBool   bool
	ParamInt    int
	ParamFloat  float64
}
type settingsUsersItem struct {
	Name  string
	EMail string
}

func defaultSettings() MySettings {
	return MySettings{
		Main: settingsMain{
			ParamString: "ParamValue \n \"test\"",
			ParamBool:   true,
			ParamInt:    12,
			ParamFloat:  1.234,
		},
		Users: []settingsUsersItem{
			{
				Name:  "John",
				EMail: "john@email",
			},
			{
				Name: "Smith",
			},
		},
	}
}

type loadFunc func([]byte, interface{}) error

func testLoad(f loadFunc, src string, t *testing.T) {
	data := defaultSettings()

	err := f([]byte(src), &data)
	if err != nil {
		t.Fatal(err)
	}

	if !data.Main.ParamBool {
		t.Fatal("defaults non persistono")
	}
	if data.Main.ParamInt != 13 {
		t.Fatal("override non settati")
	}
}

func TestLoadJson(t *testing.T) {
	js := `
{ 
	"Main": {
		"paramint": 13
	} 
}
`
	testLoad(parsers.LoadJson, js, t)
}

func TestLoadJsonc(t *testing.T) {
	jsc := `
{ 
	"Main": {
		// Override
		"ParamInt": 13
	} 
}
`
	testLoad(parsers.LoadJsonc, jsc, t)
}

func TestLoadYaml(t *testing.T) {
	yaml := `
main:
  Paramint: 13
`

	testLoad(parsers.LoadYaml, yaml, t)
}

func TestLoadToml(t *testing.T) {
	toml := `
[main]
Paramint = 13
`

	testLoad(parsers.LoadToml, toml, t)
}

func TestSave(t *testing.T) {
	defaults := settingsUsersItem{
		Name:  "foo",
		EMail: "foo",
	}

	cfg := settingsUsersItem{
		Name:  "foo",
		EMail: "bar",
	}

	mapToSave, err := diff(defaults, cfg)
	if err != nil {
		t.Fatal(err)
	}

	bb, err := parsers.SaveJson(mapToSave)
	if err != nil {
		t.Fatal(err)
	}

	js := `
{
	"EMail": "bar"
}
`
	if strings.TrimSpace(string(bb)) != strings.TrimSpace(js) {
		t.Fatal("mismatch")
	}
}

func TestSaveJson(t *testing.T) {
	config := settingsUsersItem{
		Name:  "foo",
		EMail: "bar",
	}

	bb, err := parsers.SaveJson(config)
	if err != nil {
		t.Fatal(err)
	}

	js := `
{
	"Name": "foo",
	"EMail": "bar"
}
`
	if strings.TrimSpace(string(bb)) != strings.TrimSpace(js) {
		t.Fatal("mismatch")
	}
}

func TestSaveYaml(t *testing.T) {
	config := settingsUsersItem{
		Name:  "foo",
		EMail: "bar",
	}

	bb, err := parsers.SaveYaml(config)
	if err != nil {
		t.Fatal(err)
	}

	yaml := `
name: foo
email: bar
`
	if strings.TrimSpace(string(bb)) != strings.TrimSpace(yaml) {
		t.Fatal("mismatch")
	}
}

func TestSaveToml(t *testing.T) {
	config := settingsUsersItem{
		Name:  "foo",
		EMail: "bar",
	}

	bb, err := parsers.SaveToml(config)
	if err != nil {
		t.Fatal(err)
	}

	toml := `
Name = 'foo'
EMail = 'bar'
`
	if strings.TrimSpace(string(bb)) != strings.TrimSpace(toml) {
		t.Fatal("mismatch")
	}
}

func TestMultiLoad(t *testing.T) {
	yaml := `
main:
  paramint: 13
`

	data := defaultSettings()

	err := parsers.LoadYaml([]byte(yaml), &data)
	if err != nil {
		t.Fatal(err)
	}

	yaml2 := `
main:
  paramint: 14
  paramstring: asd
`

	err = parsers.LoadYaml([]byte(yaml2), &data)
	if err != nil {
		t.Fatal(err)
	}

	if !data.Main.ParamBool {
		t.Fatal("defaults non persistono")
	}
	if data.Main.ParamInt != 14 || data.Main.ParamString != "asd" {
		t.Fatal("override non settati")
	}
}
