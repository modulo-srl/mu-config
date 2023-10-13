package main

import (
	"fmt"

	"github.com/modulo-srl/mu-config/settings"
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
	Email string
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
				Email: "john@email",
			},
			{
				Name: "Smith",
			},
		},
	}
}

func main() {
	cfg := defaultSettings()

	// Carica i settaggi da file jsonc.
	err := settings.LoadFile("mysettings.jsonc", &cfg)
	if err != nil {
		panic(err)
	}

	// Carica ulteriori settaggi in override da file yaml.
	err = settings.LoadFile("mysettings.yaml", &cfg)
	if err != nil {
		panic(err)
	}

	// Carica override di eventuali valori secret da systemd.
	err = settings.LoadSystemdCredentials("mysettings.toml", &cfg)
	if err != nil {
		panic(err)
	}

	// Manipola i settaggi.
	cfg.Main.ParamInt = 99

	// Salva i settaggi su file toml (ad esclusione dei valori di default)
	err = settings.SaveFile("mysettings.toml", cfg, defaultSettings())
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", cfg)
}
