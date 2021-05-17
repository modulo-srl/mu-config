# JSON2Struct

Converts JSON file to Go structs and initializer function.

* supports JSONC format with comments

## Installation

`go get -u github.com/modulo-srl/mu-config/json2struct`

## Example usage with Go-generate

Generate structs declarations and initializers, using default parameters:
`//go:generate json2struct -out=config-data.go ../bin/settings-default.jsonc`

Generate (to stdout) package "settings" and "settingsData" struct declaraction, without warning header, initializers and raw content:
`//go:generate json2struct -pkg="settings" -declare="settingsData" -values="" -warn=false --raw="" ../bin/settings-default.jsonc`
