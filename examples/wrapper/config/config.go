package config

// Depends: go get -u github.com/modulo-srl/mu-config/json2struct
//go:generate json2struct -out=data.go default.jsonc

import (
	"time"

	"github.com/modulo-srl/mu-config/settings"
)

var cfg Data

var configCtrl *settings.Settings

const saveAfter = 5 * time.Second

// GetConfigData returns config data struct.
func GetConfigData(forceReload bool) *Data {
	if forceReload {
		err := configCtrl.LoadSettings(false)
		if err != nil {
			panic(err)
		}
	}

	return &cfg
}

// SaveConfig can be called frequently (uses internal antibounce).
func SaveConfig() {
	configCtrl.SaveSettingsDebounce(saveAfter)
}

// SaveConfigNow saves immediately.
func SaveConfigNow() error {
	return configCtrl.SaveSettings()
}

// LoadConfig loads the settings.
func LoadConfig() error {
	err := configCtrl.LoadSettings(false)
	if err != nil {
		return err
	}

	return nil
}

// GetFilename returns the current config full path.
func GetFilename() string {
	return configCtrl.GetFilename()
}

// GetDefaultConfig returns the original config file content.
func GetDefaultConfig() string {
	return raw
}

// Init initialize settings.
func Init(filename string) error {
	var err error

	configCtrl, err = settings.New(filename, &cfg, defaultData(), false)
	if err != nil {
		return err
	}

	return nil
}
