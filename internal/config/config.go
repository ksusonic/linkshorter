package config

import (
	"log"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	Address     string `env:"ADDRESS" envDefault:"localhost:8080"`
	DatabaseDsn string `env:"DATABASE_DSN"`
}

func NewConfig() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	return &cfg
}
