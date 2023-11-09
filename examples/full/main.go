package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/modulo-srl/mu-config/examples/full/config"
)

func main() {
	fdefaultConfig := flag.Bool("default-config", false, "Show default config")
	fconfig := flag.String("config-file", "settings", "Config file (absolute or relative path, without extension will try for .json, .jsonc, .yaml, .toml)")

	if *fdefaultConfig {
		fmt.Print(config.GetDefaultConfig())
		os.Exit(0)
	}

	configFilename := *fconfig

	err := config.Load(configFilename)
	if err != nil {
		panic(err)
	}

	fmt.Println(config.Cfg)
}
