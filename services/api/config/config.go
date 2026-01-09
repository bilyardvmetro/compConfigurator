package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPConfig struct {
	Address string        `yaml:"address" env:"API_ADDRESS" env-default:"localhost:80"`
	Timeout time.Duration `yaml:"timeout" env:"API_TIMEOUT" env-default:"5s"`
}

type Config struct {
	LogLevel      string        `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
	HTTPConfig    HTTPConfig    `yaml:"api_server"`
	UpdateAddress string        `yaml:"update_address" env:"UPDATE_ADDRESS" env-default:"update:82"`
	SearchAddress string        `yaml:"search_address" env:"SEARCH_ADDRESS" env-default:"search:83"`
	TokenTTL      time.Duration `yaml:"token_ttl" env:"TOKEN_TTL" env-default:"24h"`
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
