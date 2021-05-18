package main

// Depends: go get -u github.com/modulo-srl/mu-config/json2struct
//go:generate json2struct -out=mySettings.go -declare=mySettings -pkg=main -raw= default-settings.jsonc

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/modulo-srl/mu-config/settings"
)

var cfg MySettings
var config *settings.Settings

func showConfig() {
	fmt.Println("Param string:", cfg.Main.ParamString)
	fmt.Println("Param bool:", cfg.Main.ParamBool)
	fmt.Println("Param int:", cfg.Main.ParamInt)
	fmt.Println("Param float:", cfg.Main.ParamFloat)

	fmt.Println("Param users:", cfg.Users)
}

func main() {
	var err error

	config, err = settings.New("settings.json", &cfg, defaultMySettings(), false)
	if err != nil {
		panic(err)
	}

	err = config.LoadSettings(false)
	if err != nil {
		panic(err)
	}

	fmt.Println("Current settings...")
	showConfig()

	fmt.Println("\nChanging settings...")
	rand.Seed(time.Now().UnixNano())
	cfg.Main.ParamString = fmt.Sprintf("Changed to %d", rand.Int())
	cfg.Users[1].Email = fmt.Sprintf("new@email-%d.com", rand.Int())

	config.SaveSettings() // Saves immediately

	config.LoadSettings(true)
	showConfig()

	// Changing settings frequently
	fmt.Println("\nChanging settings frequently...")

	var saveAfter = time.Second

	cfg.Main.ParamString = "Changed to 1"
	config.SaveSettingsDebounce(saveAfter)
	cfg.Main.ParamString = "Changed to 2"
	config.SaveSettingsDebounce(saveAfter)
	cfg.Main.ParamString = "Changed to 3"
	config.SaveSettingsDebounce(saveAfter)

	time.Sleep(saveAfter * 2)

	config.LoadSettings(true)
	showConfig()
}
