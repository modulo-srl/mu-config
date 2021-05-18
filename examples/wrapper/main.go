package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/modulo-srl/mu-config/examples/wrapper/config"
)

func main() {
	fdefaultConfig := flag.Bool("default-config", false, "Show default config")
	fconfig := flag.String("config-file", "settings.json", "Config file's (absolute or relative path)")

	if *fdefaultConfig {
		fmt.Print(config.GetDefaultConfig())
		os.Exit(0)
	}

	configFilename := *fconfig
	err := config.Init(configFilename)
	if err != nil {
		panic(err)
	}

	err = config.LoadConfig()
	if err != nil {
		panic(err)
	}

	fmt.Println("Config loaded by", config.GetFilename())

	cfg := config.GetConfigData(false)
	fmt.Println("Current config:", cfg)

	// Get config forcing reload (simulating other task changing the config file)
	cfg = config.GetConfigData(true)

	fmt.Println("\nCurrent reloaded config:", cfg)

	fmt.Println("\nChanging settings...")
	rand.Seed(time.Now().UnixNano())
	cfg.Main.ParamString = fmt.Sprintf("Changed to %d", rand.Int())
	cfg.Users[1].Email = fmt.Sprintf("new@email-%d.com", rand.Int())

	config.SaveConfigNow() // Saves immediately

	config.LoadConfig()
	fmt.Println("\nCurrent reloaded config:", cfg)
}
