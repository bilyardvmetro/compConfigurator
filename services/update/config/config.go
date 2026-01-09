package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel  string `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
	Address   string `yaml:"update_address" env:"UPDATE_ADDRESS" env-default:"localhost:80"`
	DBAddress string `yaml:"db_address" env:"DB_ADDRESS" env-default:"postgres://postgres:password@postgres:5432/postgres"`
}

func MustLoad(configPath string) Config {
	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			log.Fatalf("cannot read env vars: %s", err)
		}
		return cfg
	}

	if err := cleanenv.UpdateEnv(&cfg); err != nil {
		log.Printf("warning: cannot update from env: %s", err)
	}

	return cfg
}
