# go2cfg

go2cfg è un programma autonomo, un generatore go e una libreria che crea
file jsonc/toml/yaml a partire da una struttura in go, compresi blocchi
di documentazione, commenti e valori predefiniti, facilitando, ad esempio,
la manutenzione dei template di file di configurazione.

Qualsiasi struttura incorporata o nidificata in quella iniziale sarà inclusa,
anche proveniente da pacchetti diversi.

Le costanti tipizzate vengono risolte e incluse nei commenti generati
per semplificare la verifica e la modifica dei valori.

go2cfg è scritto in modo tale che il parsing e l'estrazione dei dati necessari
da AST e dei tipi sono disaccoppiati dalla generazione del file
di configurazione, rendendo possibile scrivere praticamente qualsiasi formato.

## Valori di default

Se nel pacchetto per cui si vuole generare il file di config
esiste una funzione con questo tipo di firma:

```go
func StructTypeNameDefaults() *StructTypeName
```

allora verrà parsata ed estratti i valori di default per ogni elemento.  
Per semplicità di parsing il corpo della funzione deve presentare la sintassi
come di seguito illustrato.

## Esempio

Codice:

```go
package multipkg

import (
    "github.com/modulo-srl/mu-config/go2cfg/testdata/multipkg/network"
    alias "github.com/modulo-srl/mu-config/go2cfg/testdata/multipkg/stats"
)

//go:generate go2cfg -type MultiPackage -out multi_package.jsonc

// MultiPackage tests the multi-package and import aliasing case.
type MultiPackage struct {
    NetStatus  network.Status // Network status.
    alias.Info                // Statistics info.
}

func MultiPackageDefaults() *MultiPackage {
    return &MultiPackage{
        NetStatus: network.Status{
            Connected: true,
            State:     network.StateDisconnected,
        },
        Info: alias.Info{
            PacketLoss:    32 * 2,
            RoundTripTime: 123,
        },
    }
}
```

```go
package network

type ConnState int

const (
    // StateDisconnected signals the Disconnected state.
    StateDisconnected ConnState = iota
    // StateConnecting signals the connection-pending state.
    StateConnecting
    // StateConnected signals the Connected state.
    StateConnected
)

const (
    // StateFailed signals the Failed state.
    StateFailed ConnState = iota + 5
    // StateReconnecting signals the Reconnecting state.
    StateReconnecting
)

// Status reports connection status.
type Status struct {
    Connected bool      // Connected flag comment.
    State     ConnState // Connection state comment.
}
```

```go
package stats

// Info reports statistical info.
type Info struct {
    // PacketLoss documentation block.
    PacketLoss    int `json:"packet_loss"`     // Packet loss comment.
    RoundTripTime int `json:"round_trip_time"` // Round-trip time in milliseconds.
}
```

Output generato:

```json5
{
    // Network status.
    "NetStatus": {
        // Connected flag comment.
        "Connected": true,
      
        // Connection state comment.
        // Allowed values:
        // StateDisconnected = 0  StateDisconnected signals the Disconnected state.
        // StateConnecting   = 1  StateConnecting signals the connection-pending state.
        // StateConnected    = 2  StateConnected signals the Connected state.
        // StateFailed       = 5  StateFailed signals the Failed state.
        // StateReconnecting = 6  StateReconnecting signals the Reconnecting state.
        "State": 0
    },

    // PacketLoss documentation block.
    // Packet loss comment.
    "packet_loss": 64,

    // Round-trip time in milliseconds.
    "round_trip_time": 123
}
```

## Installazione

```shell
go install github.com/modulo-srl/mu-config/go2cfg@latest
```

## Utilizzo (CLI)

```shell
go2cfg -type <type-name> [-doc-types doc] [-out filename] [package-dir]
```

- `-doc-types` - `string`: Flag per generare anche il tipo dei valori nei commenti
- `-out` - `string`: nome file di uscita output filepath; se omesso l'output
  sarà verso `stdout` in jsonc
- `-type` - `string`: nome tipo struttura da cui generare l'output (obbligatorio)
- `package-dir`: directory che contiene il file go dove il nome tipo struttura
  è definito; se omesso utilizza la dir corrente

Valori permessi per `-doc-types`:

- `all`: genera tutti i tipi per tutti gli elementi
- `basic`: genera i tipi dei soli elementi di tipo di base (int, float, bool, string)

## Utilizzo (generatore)

Aggiungere il commento speciale `go:generate` dove desiderato:

```go
//go:generate go2cfg -type type-name [-doc-types doc] -out outfile.ext [package-dir]
```

`package-dir` può essere tranquillamente omesso in questo caso,
poiché verrà usata la directory dove il file corrente risiede.

## Come libreria

go2cfg contiene due pacchetti:

- go2cfg/generator
- go2cfg/distiller

Basta importare il primo e usare la funzione:

```go
func Generate(dir, typeName string, mode DocTypesMode) (string, error)
```

Oppure importare il secondo per estrarre facilmente le info da AST
per generare magari altri formati desiderati.

## Limitazioni

go2jcfg supporta mappe, ma con le seguenti limitazioni:

- le chiavi devono essere di tipo stringa;
- i valori possono essere di qualunque tipo tranne che `interface{}`,
ciò significa che tutti i valori della mappa devono essere del medesimo tipo;
- le mappe non possono essere annidate.

Queste scelte sono state fatte considerando che l'aggiunta dei tipi di mappe
non rende possibile integrare la documentazione della struttura e dei tipi richiesti,
annullando completamente l'intento di jsonc o di altri formati che necessitano
commenti per essere facilmente comprensibili.
