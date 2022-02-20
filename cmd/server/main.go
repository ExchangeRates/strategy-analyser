package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/ExchangeRates/strategy-analyser/internal"
	"github.com/ExchangeRates/strategy-analyser/internal/config"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config/service.toml", "path to config file")
}

func main() {
	flag.Parse()

	configuration := config.NewConfig()
	_, err := toml.DecodeFile(configPath, configuration)
	if err != nil {
		log.Fatal(err)
	}

	if err := internal.Start(configuration); err != nil {
		log.Fatal(err)
	}
}
