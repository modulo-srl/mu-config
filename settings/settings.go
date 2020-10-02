package settings

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/modulo-srl/mu-config/jsonc"
)

// New - restituisce un gestore inizializzato per dati custom
// - filename: se il path è relativo, allora lo cerca nel medesimo path dell'eseguibile. Accetta eventualmente "~/...".
// Se esiste un file con estensione ".jsonc" preferisce quello.
// - data: passare il puntatore ad una struct contenente i dati da caricare e salvare.
// - defaultData: passare una struttura dati popolata con i dati di default,
// che non verranno mai salvati nel file di configurazione.
func New(filename string, data, defaultData interface{}, verbose bool) (*Settings, error) {
	ctrl := Settings{}
	err := ctrl.init(filename, data, defaultData, verbose)
	if err != nil {
		return nil, err
	}

	return &ctrl, nil
}

// GetFilename returns full path filename
func (set *Settings) GetFilename() string {
	return set.filename
}

// LoadSettings - carica/ricarica le impostazioni
func (set *Settings) LoadSettings() error {
	if set.Data == nil {
		return errors.New("settings data struct not set")
	}

	set.lock.Lock()
	defer set.lock.Unlock()

	if set.verbose {
		log.Println("Load settings")
	}

	fileContent, err := ioutil.ReadFile(set.filename)
	if err != nil {
		return err
	}

	if filepath.Ext(set.filename) == ".jsonc" {
		fileContent = jsonc.ToJSON(fileContent)
	}

	// Reimposta i dati di default, evitando riferimenti di qualsiasi tipo.
	cloneData(set.defaultData, set.Data)

	data, isMap := set.Data.(map[string]interface{})
	if isMap {
		// I settings sono in una mappa generica:
		// è costretto a deserializzare ogni singola struttura, anche se in modo inefficiente.
		var temp map[string]interface{}
		err = json.Unmarshal([]byte(fileContent), &temp)
		if err != nil {
			return err
		}

		for key, item := range temp {
			tempj, _ := json.Marshal(item)

			err := json.Unmarshal(tempj, data[key])
			if err != nil {
				return err
			}
		}
	} else {
		err = json.Unmarshal([]byte(fileContent), &set.Data)
		if err != nil {
			return err
		}
	}

	//fmt.Println(set)

	if set.Data == nil {
		return errors.New("null settings in file")
	}

	return nil
}

// SaveSettings - salva le impostazioni
func (set *Settings) SaveSettings() error {
	if set.Data == nil {
		return errors.New("data not set")
	}

	set.lock.Lock()
	defer set.lock.Unlock()

	if set.verbose {
		log.Println("Save settings")
	}

	f, err := os.Create(set.filename)
	if err != nil {
		return err
	}
	defer f.Close()

	diffMap, err := DiffMaps(set.defaultData, set.Data)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(diffMap, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(set.filename, b, 0666)

	return err
}

// SaveSettingsDebounce - salva le impostazioni dopo un certo ritardo dall'ultima invocazione.
// Ogni invocazione resetta il il conteggio del timeout.
func (set *Settings) SaveSettingsDebounce(saveAfter time.Duration) {
	if set.timerSave != nil {
		set.timerSave.Reset(saveAfter)
	} else {
		saveFunc := func() {
			set.SaveSettings()
		}

		set.timerSave = time.AfterFunc(saveAfter, saveFunc)
	}
}

// Settings - struttura oggetto direttamente utilizzabile.
type Settings struct {
	verbose  bool
	filename string

	lock      sync.Mutex
	timerSave *time.Timer

	Data        interface{} // Punta ad una struttura dati custom
	defaultData interface{} // Contiene i dati di default
}

func (set *Settings) init(filename string, data, defaultData interface{}, verbose bool) error {
	switch {
	case filename[0] == '~':
		homeDir, err := homedir.Dir()
		if err != nil {
			return errors.New("Can't find homedir")
		}

		set.filename = homeDir + filename[1:]

	case filename[0] != '/':
		// Se il path è relativo allora prefissa con il path dell'eseguibile
		ex, err := os.Executable()
		if err != nil {
			return err
		}
		set.filename = filepath.Dir(ex) + "/" + filename

	default:
		set.filename = filename
	}

	if filepath.Ext(set.filename) == ".json" && fileExists(set.filename+"c") {
		// filename.jsonc exists, so use it instead.
		set.filename += "c"
	}

	set.Data = data
	set.defaultData = defaultData

	set.verbose = verbose

	if verbose {
		log.Println("Settings filename", set.filename)
	}

	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
