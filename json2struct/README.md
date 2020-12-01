# JSON2Struct
Converts JSON files to Go structs and initializers functions.

* supports JSONC format with comments

## Installation
`go get -u github.com/modulo-srl/mu-config/json2struct`

## Example usage with Go-generate
Generate structs declarations and initializers, using default parameters:
`//go:generate json2struct -out=config_data.go ../bin/default.jsonc`

Generate (to stdout) package "settings" and "settingsData" struct declaraction, without warning, initializers and raw content:
`//go:generate json2struct -pkg="settings" -declare="settingsData" -values="" -warn=false --raw="" ../bin/default.jsonc`
