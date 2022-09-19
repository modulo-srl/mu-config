package main

// Depends on: go install github.com/marco-sacchi/go2jsonc/cmd/go2jsonc
//go:generate go2jsonc -type=MySettings -out=default-settings.jsonc

import (
	_ "embed"
	"fmt"
	"github.com/modulo-srl/mu-config/settings"
)

var cfg MySettings

var configCtrl *settings.Settings

// Incorporates the configuration into the executable so that it can be viewed on request.
// go2jsonc needs the sources to generate the jsonc file.
//go:embed default-settings.jsonc
var configDefaultsRaw string

func main() {
	var err error

	configCtrl, err = settings.New("settings.jsonc", &cfg, MySettingsDefaults(), false)
	if err != nil {
		panic(err)
	}

	err = configCtrl.LoadSettings(false)
	if err != nil {
		panic(err)
	}

	fmt.Println("Param string:", cfg.Main.ParamString)
	fmt.Println("Param bool:", cfg.Main.ParamBool)
	fmt.Println("Param int:", cfg.Main.ParamInt)
	fmt.Println("Param float:", cfg.Main.ParamFloat)

	fmt.Println("Param users:", cfg.Users)

	fmt.Println("\nDefault configuration:")
	fmt.Println(configDefaultsRaw)
}
