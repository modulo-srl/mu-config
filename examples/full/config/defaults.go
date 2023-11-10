package config

// Depends on: go install github.com/modulo-srl/mu-config/go2cfg@latest
//go:generate go2cfg -type=MySettings -out=defaults.jsonc

import (
	_ "embed"
)

type MySettings struct {
	// Main configuration parameters.
	Main SettingsMain

	// Users list.
	Users []SettingsUserItem
}

// Sub type example.
type SubType struct {
	ParamMap   map[string]string
	ParamArray []int
}

type SettingsMain struct {
	ParamSub    SubType
	ParamString string
	ParamBool   bool
	ParamInt    int
	ParamFloat  float64 // Float value, default 1.234
}

type SettingsUserItem struct {
	Name  string // User name.
	Email string // User e-mail.
}

func MySettingsDefaults() *MySettings {
	return &MySettings{
		Main: SettingsMain{
			ParamString: "ParamValue \n \"test\" 😊 世界",
			ParamBool:   true,
			ParamInt:    12,
			ParamFloat:  1.234,
			ParamSub: SubType{
				ParamMap:   map[string]string{"key 1": "/path/to/file", "key2": "c:\\path\\to\\file"},
				ParamArray: []int{1, 2, 3},
			},
		},
		Users: []SettingsUserItem{
			{
				Name:  "John",
				Email: "john@email",
			},
			{
				Name: "Smith",
			},
		},
	}
}

//go:embed defaults.jsonc
var defaultsRaw string
