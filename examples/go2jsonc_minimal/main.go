package main

// Depends: go install github.com/marco-sacchi/go2jsonc/cmd/go2jsonc
//go:generate go2jsonc -type=MySettings -out=default-settings.jsonc

import (
	"fmt"
	"github.com/modulo-srl/mu-config/settings"
)

var cfg MySettings

var configCtrl *settings.Settings

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
}
