// Package config package config
package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

// Config -.
type Config struct {
	HTTP       `yaml:"http"`
	PG         `yaml:"postgres"`
	Sinopac    `yaml:"sinopac"`
	Deployment string `env-required:"true" env:"DEPLOYMENT"`
}

// HTTP -.
type HTTP struct {
	Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
}

// PG -.
type PG struct {
	PoolMax int    `env-required:"true" yaml:"pool_max"`
	URL     string `env-required:"true" env:"PG_URL"`
	DBName  string `env-required:"true" env:"DB_NAME"`
}

// Sinopac -.
type Sinopac struct {
	PoolMax int    `env-required:"true" yaml:"pool_max"`
	URL     string `env-required:"true" env:"SINOPAC_URL"`
}

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./configs/config.yml", cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
