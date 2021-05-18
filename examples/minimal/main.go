package main

// Depends: go get -u github.com/modulo-srl/mu-config/json2struct
//go:generate json2struct -out=mySettings.go -declare=mySettings -pkg=main -raw= default-settings.jsonc

import (
	"fmt"

	"github.com/modulo-srl/mu-config/settings"
)

var cfg MySettings

var configCtrl *settings.Settings

func main() {
	var err error

	configCtrl, err = settings.New("settings.json", &cfg, defaultMySettings(), false)
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
