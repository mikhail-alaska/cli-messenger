package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address  string `yaml:"address" env-default:"localhost:8080"`
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
}

func MustLoad() *Config {
	configPath := "./config/local.yaml"
	if configPath == "" {
		log.Fatal("We have a problem with config, CONFIG_PATH not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %v doesnt exist", configPath)
	}
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config, %s", err)
	}

	return &cfg
}
