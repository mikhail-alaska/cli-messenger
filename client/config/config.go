package config

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPServer `yaml:"http_server" env-required:"true"`
	Login      string `yaml:"login" env-required:"true"`
	Password   string `yaml:"password" env-required:"true"`
	OpenKey    int    `yaml:"openkey" env-required:"true"`
	ClosedKey  int    `yaml:"closedkey" env-required:"true"`
}

type HTTPServer struct {
	Address string `yaml:"address" env-default:"localhost:8080"`
}

func (c Config) Check() bool {
	r, _ := regexp.Compile("^[a-zA-Z0-9]+([_ -]?[a-zA-Z0-9])*$")

	l := r.MatchString(c.Login)
	p := r.MatchString(c.Password)

	if !l {
		fmt.Println("this username can not be used")
		return false
	}

	if !p {
		fmt.Println("this password can not be used")
		return false
	}

	return true
}

func MustLoad() *Config {
	configPath := "./config/local.yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %v doesnt exist", configPath)
	}
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config, %s", err)
	}

	return &cfg
}
