package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	HTTP struct {
		Host string
		Port int
	}
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "../config/config.yaml"
	}
	file, err := os.Open(configPath)

	if err != nil {
		panic(err)
	}

	defer file.Close()
	var config Config
	err = yaml.NewDecoder(file).Decode(&config)
	if err != nil {
		panic(err)
	}

	return &config
}
