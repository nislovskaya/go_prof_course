package cmd

import (
	"flag"
	"log"
	"os"

	"github.com/nislovskaya/go_prof_course/hw12_13_14_15_calendar/internal/configs"
)

var ConfigFile string

func init() {
	flag.StringVar(&ConfigFile, "config", "configs/config.yml", "Path to configuration file")
}

func GetConfig(configFile string) *configs.Config {
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatalf("failed to open config file: %s", err)
	}

	cfg, err := configs.NewConfig(file)
	if err != nil {
		log.Fatalf("failed to parse config file: %s", err)
	}

	return cfg
}
