package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string     `yaml:"env"`
	Http_server HTTPServer `yaml:"http_server"`
	Db          DB         `yaml:"db"`
}

type HTTPServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type DB struct {
	Username string `yaml:"username"`
	DBName   string `yaml:"dbname"`
	Password string `yaml:"password"`
}

func MustLoad() *Config {
	configPath := "../../configs/local.yaml"

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg

	// panic("not implemented")
}
