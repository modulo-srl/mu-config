package config

import "github.com/modulo-srl/mu-config/settings"

// Configurazione globale.
var Cfg = MySettingsDefaults()

// Ritorna il contenuto del file di configurazione di default precedentemente generato.
func GetDefaultConfig() string {
	return defaultsRaw
}

// Carica la configurazione da file.
func Load(filename string) error {
	return settings.LoadFile(filename, &Cfg)
}
