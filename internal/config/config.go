package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP struct {
		Port            string        `env:"HTTP_PORT" env-default:"8080"`
		ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"5s"`
		WriteTimeout    time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"10s"`
		ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" env-default:"10s"`
	}

	DB struct {
		DSN           string        `env:"POSTGRES_DSN" env-required:"true"`
		MaxConns      int32         `env:"DB_MAX_CONNS" env-default:"10"`
		MinConns      int32         `env:"DB_MIN_CONNS" env-default:"0"`
		ConnMaxIdle   time.Duration `env:"DB_CONN_MAX_IDLE" env-default:"5m"`
		ConnMaxLife   time.Duration `env:"DB_CONN_MAX_LIFE" env-default:"30m"`
		HealthTimeout time.Duration `env:"DB_HEALTH_TIMEOUT" env-default:"2s"`
	}

	Auth struct {
		JWTSecret string        `env:"JWT_SECRET" env-required:"true"`
		JWTTTL    time.Duration `env:"JWT_TTL" env-default:"24h"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" env-default:"info"`
	}
}

func Load() (Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
