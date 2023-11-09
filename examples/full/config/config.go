package config

import (
	"fmt"
	"path/filepath"

	"github.com/modulo-srl/mu-config/settings"
)

// Configurazione globale.
var Cfg = MySettingsDefaults()

// Ritorna il contenuto del file di configurazione di default precedentemente generato.
func GetDefaultConfig() string {
	return defaultsRaw
}

// Carica la configurazione da file.
func Load(filename string) error {
	loadedFilename, err := settings.LoadFile(filename, &Cfg, true)
	if err != nil {
		return err
	}
	fmt.Println("settings loaded from: " + loadedFilename)

	// Carica la configurazione da systemd, per qualsiasi formato.
	systemdFilename := filepath.Base(filename)
	loadedFilename, err = settings.LoadSystemdCredentials(systemdFilename, &Cfg, false)
	if err != nil {
		return err
	}
	if loadedFilename != "" {
		fmt.Println("systemd settings loaded from: " + loadedFilename)
	}

	return err
}
