package config

// Depends on: go install github.com/modulo-srl/mu-config/go2cfg
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

type SettingsMain struct {
	ParamString string
	ParamBool   bool
	ParamInt    int
	ParamFloat  float64 // Float value, defalt 1.234
}

type SettingsUserItem struct {
	Name  string // User name.
	Email string // User e-mail.
}

func MySettingsDefaults() *MySettings {
	return &MySettings{
		Main: SettingsMain{
			ParamString: "ParamValue \n \"test\"",
			ParamBool:   true,
			ParamInt:    12,
			ParamFloat:  1.234,
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
