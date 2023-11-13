package settings

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/modulo-srl/mu-config/settings/parsers"
)

// Carica la configurazione da file.
//   - filename: se non ha percorso o lo ha relativo, sarà rispetto alla directory corrente;
//     se ha percorso assoluto può anche iniziare per '~'.
//     Se sprovvisto di estensione tenta il caricamento di qualsiasi formato conosciuto.
//   - cfg: PUNTATORE a struttura configurazione da popolare.
//
// - errorWhenNotFound: true per generare un errore se il file non viene trovato.
func LoadFile(filename string, cfg interface{}, errorWhenNotFound bool) (loadedFilename string, err error) {
	fullpathFile, err := GetFileFullPath(filename)
	if err != nil {
		return "", err
	}

	return loadFile(fullpathFile, cfg, errorWhenNotFound)
}

// Carica la configurazione da Systemd.
// ref: https://systemd.io/CREDENTIALS/
//
//   - filename: deve essere un nome file, sprovvisto di percorso assoluto, situato in $CREDENTIALS_DIRECTORY.
//     L'estensione viene ignorata, tentando il carimento di qualsiasi formato conosciuto.
//   - cfg: PUNTATORE a struttura configurazione da popolare.
//   - errorWhenNotFound: true per generare un errore se il file non viene trovato o se $CREDENTIALS_DIRECTORY non è settato.
func LoadSystemdCredentials(filename string, cfg interface{}, errorWhenNotFound bool) (loadedFilename string, err error) {
	path := os.Getenv("CREDENTIALS_DIRECTORY")
	if path == "" {
		if errorWhenNotFound {
			return "", errors.New("systemd credential directory not found")
		}
		return "", nil
	}

	// Rimuove l'eventuale estensione, permettendo un override di qualsiasi formato.
	ext := filepath.Ext(filename)
	switch ext {
	case ".json":
		fallthrough
	case ".jsonc":
		fallthrough
	case ".yaml":
		fallthrough
	case ".toml":
		filename = strings.TrimSuffix(filename, ext)
	}

	fullpathFile := path + "/" + filename

	return loadFile(fullpathFile, cfg, errorWhenNotFound)
}

// Funzione interna per caricare la configurazione da file.
//   - filename: nome file con percorso assoluto.
//     se senza estensione cerca di caricare .json, .jsonc, .yaml, .toml
//     Se sprovvisto di estensione tenta il caricamento di qualsiasi formato conosciuto.
//   - cfg: PUNTATORE a struttura configurazione da popolare.
//   - errorWhenNotFound: true per generare un errore se il file non viene trovato.
func loadFile(filename string, cfg interface{}, errorWhenNotFound bool) (loadedFilename string, err error) {
	ext := filepath.Ext(filename)

	switch ext {
	case ".json":
		fallthrough
	case ".jsonc":
		fallthrough
	case ".yaml":
		fallthrough
	case ".toml":
		if !fileExists(filename) {
			if errorWhenNotFound {
				return "", errors.New("file not found: " + filename)
			}
			return "", nil
		}

	default:
		if fileExists(filename + ".json") {
			ext = ".json"
		} else if fileExists(filename + ".jsonc") {
			ext = ".jsonc"
		} else if fileExists(filename + ".yaml") {
			ext = ".yaml"
		} else if fileExists(filename + ".toml") {
			ext = ".toml"
		} else {
			if errorWhenNotFound {
				return "", errors.New("file not found: " + filename + ".json/.jsonc/.yaml/.toml")
			}
			return "", nil
		}
		filename += ext
	}

	// Parsa il file.
	switch ext {
	case ".json":
		fallthrough
	case ".jsonc":
		err = parsers.LoadJsoncFile(filename, &cfg)
	case ".yaml":
		err = parsers.LoadYamlFile(filename, &cfg)
	case ".toml":
		err = parsers.LoadTomlFile(filename, &cfg)
	}

	if err != nil {
		return "", fmt.Errorf("cannot parse %s: %s", filename, err)
	}

	return filename, nil
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
		fallthrough
	case ".jsonc":
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
