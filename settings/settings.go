package settings

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/modulo-srl/mu-config/settings/parsers"
)

type ErrorFileNotFound struct {
	filename string
}

func (e *ErrorFileNotFound) Error() string {
	return "file not found: " + e.filename
}

// Carica la configurazione da file.
//   - filename: se non ha percorso o lo ha relativo, sarà rispetto alla directory corrente;
//     se ha percorso assoluto può anche iniziare per '~'.
//   - cfg: PUNTATORE a struttura configurazione da popolare.
//
// Ritorna ErrorFileNotFound se il file non esiste.
func LoadFile(filename string, cfg interface{}) error {
	fullpathFile, err := GetFileFullPath(filename)
	if err != nil {
		return err
	}

	return loadFile(fullpathFile, cfg)
}

// Carica la configurazione da Systemd.
// ref: https://systemd.io/CREDENTIALS/
//
// - filename: deve essere un nome file, sprovvisto di percorso, situato in $CREDENTIALS_DIRECTORY.
// - cfg: PUNTATORE a struttura configurazione da popolare.
//
// Ritorna ErrorFileNotFound se il file non esiste.
func LoadSystemdCredentials(filename string, cfg interface{}) error {
	path := os.Getenv("CREDENTIALS_DIRECTORY")
	if path == "" {
		return errors.New("systemd credential directory not found")
	}

	fullpathFile := path + "/" + filename

	return loadFile(fullpathFile, cfg)
}

// Funzione interna per caricare la configurazione da file.
//   - filename: nome file con percorso assoluto.
//   - cfg: PUNTATORE a struttura configurazione da popolare.
//
// Ritorna ErrorFileNotFound se il file non esiste.
func loadFile(filename string, cfg interface{}) error {
	if !fileExists(filename) {
		return &ErrorFileNotFound{filename: filename}
	}

	var err error
	var data interface{}

	// Parsa il file.
	switch filepath.Ext(filename) {
	case ".json":
		err = parsers.LoadJsonFile(filename, &data)
	case ".jsonc":
		err = parsers.LoadJsoncFile(filename, &data)
	case ".yaml":
		err = parsers.LoadYamlFile(filename, &data)
	case ".toml":
		err = parsers.LoadTomlFile(filename, &data)
	}

	if err != nil {
		return fmt.Errorf("cannot parse %s: %s", filename, err)
	}

	return nil
}

// Salva la configurazione su file.
//   - filename: se non ha percorso o lo ha relativo, sarà rispetto alla directory corrente;
//     se ha percorso assoluto può anche iniziare per '~'.
//   - cfg: struttura configurazione da salvare.
//   - defaults: (opzionale) struttura configurazione di default.
//     se passata il file conterrà i soli valori che differiscono da questa struttura.
func SaveFile(filename string, cfg interface{}, defaults interface{}) error {
	if cfg == nil {
		return errors.New("config data cannot be nil")
	}

	var mapToSave interface{}

	if defaults == nil {
		var err error

		mapToSave, err = diff(defaults, cfg)
		if err != nil {
			return err
		}
	} else {
		mapToSave = cfg
	}

	var err error
	ext := filepath.Ext(filename)

	switch ext {
	case ".json":
		err = parsers.SaveJsonFile(filename, mapToSave)
	case ".yaml":
		err = parsers.SaveYamlFile(filename, mapToSave)
	case ".toml":
		err = parsers.SaveTomlFile(filename, mapToSave)
	default:
		err = fmt.Errorf("no encoder for %s extension", ext)
	}

	if err != nil {
		return fmt.Errorf("cannot save to %s: %s", filename, err)
	}

	return nil
}

// Ritorna il nome file completo di percorso assoluto.
//   - filename: se non ha percorso o lo ha relativo, sarà rispetto alla directory corrente;
//     se ha percorso assoluto può anche iniziare per '~'.
func GetFileFullPath(filename string) (string, error) {
	var fullpath string

	switch filename[0] {
	case '/':
		// Percorso assoluto.
		fullpath = filename

	case '~':
		// Percorso assoluto alla home.
		homeDir, err := homedir.Dir()
		if err != nil {
			return "", errors.New("can't find homedir: " + err.Error())
		}
		fullpath = path.Join(homeDir, filename[1:])

	default:
		// Percorso relativo: prefissa con la directory corrente.
		currDir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		fullpath = path.Join(currDir, filename)
	}

	return fullpath, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
