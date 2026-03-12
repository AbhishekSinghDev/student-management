package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Address string
}

type Config struct {
	Env string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}

func MustLoad() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	// check configPath is missing in env
	if configPath == "" {
		// check weither the configPath provided via cmd flags
		configFlagPtr := flag.String("config", "", "path to the configuration file")
		configPath = *configFlagPtr

		if configPath == "" {
			// log fatal, we cannot start server without configuration file
			log.Fatal("Config path is not set")
		}
	}

	// check the config file actually exists on the provided path
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exists on this path: %s", configPath)
	}

	// now we can proceed with the serialization of the config file into var
	var config Config

	err := cleanenv.ReadConfig(configPath, &config)
	if err != nil {
		log.Fatalf("cannot read config file: %s", err.Error())
	}

	return &config
}
