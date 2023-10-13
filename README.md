# muConfig

Libreria leggera per la gestione dei settings.

* Supporta i formati Json, Jsonc (json con commenti), Yaml, Toml;

* Supporta il load multiplo a mo' di override.

* Salva le sole differenze rispetto ai valori di default.

## Note

* Go:
  i nomi delle variabili nelle struct di configurazione
  devono obbligatoriamente iniziare per lettera maiuscola.

* Json/c, Yaml, Toml:
  i nomi variabili nelle struct di configurazione sono case insensitive.

* L'utilizzo dei tag `json:".."`, `toml:".."`, `yaml:".."`
  è sconsigliabile sia per la scomodità di doverli esprimere sempre tutti,
  sia per questioni di mantenibilità nel caso di subentro di nuovi parser.

* I decoder sono configurati in modalità _strict_,
  ovvero ritornano errore nel caso di field presenti nei file di configurazione
  ma mancanti nella struct destinataria in Go.

## Utilizzo

Vedi `examples/`.

## TODO

* go2yaml

* go2toml

* testare il caricamento con systemd
