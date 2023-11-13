package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/modulo-srl/mu-config/examples/full/config"
)

func main() {
	fconfig := flag.String("config-file", "settings", "Config file (absolute or relative path, without extension will try for .json, .jsonc, .yaml, .toml)")
	fdefaultConfig := flag.String("default-config", "", "Show default config (\"json\", \"toml\" or \"yaml\")")
	flag.Parse()

	if *fdefaultConfig != "" {
		s := config.GetDefaultConfig(*fdefaultConfig)
		if s == "" {
			flag.Usage()
			os.Exit(1)
		}
		fmt.Println(s)
		os.Exit(0)
	}

	configFilename := *fconfig

	err := config.Load(configFilename)
	if err != nil {
		panic(err)
	}

	fmt.Println(config.Cfg)
}
