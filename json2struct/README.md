# JSON2Struct
Converts JSON files to Go structs and initializers functions.

* supports JSONC format with comments

## Installation
`go get -u github.com/modulo-srl/mu-config/json2struct`

## Example usage with Go-generate
Generate structs declarations:
`//go:generate json2struct -declare -warn -struct=configData -pkg=config -out=structs.go ../bin/default.jsonc`

Generate structs initializer:
`//go:generate json2struct -values  -warn -struct=configData -pkg=config -out=defaults.go ../bin/default.jsonc`
