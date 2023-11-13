package config

import (
	"fmt"
	"path/filepath"

	"github.com/modulo-srl/mu-config/settings"
)

// Configurazione globale.
var Cfg = MySettingsDefaults()

// Ritorna il contenuto del file di configurazione di default precedentemente generato.
// - format: "json", "toml", "yaml"
func GetDefaultConfig(format string) string {
	switch format {
	case "json":
		fallthrough
	case "jsonc":
		return defaultsRawJsonc
	case "toml":
		return defaultsRawToml
	case "yaml":
		return defaultsRawYaml
	default:
		return ""
	}
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
