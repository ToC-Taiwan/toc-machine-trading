// Package config package config
package config

import (
	"sync"

	"toc-machine-trading/pkg/logger"

	"github.com/ilyakaznacheev/cleanenv"
)

var log = logger.Get()

var (
	globalConfig *Config
	once         sync.Once
)

// GetConfig -.
func GetConfig() (*Config, error) {
	if globalConfig != nil {
		return globalConfig, nil
	}

	once.Do(parseConfigFile)
	return globalConfig, nil
}

func parseConfigFile() {
	newConfig := Config{}
	err := cleanenv.ReadConfig("./configs/config.yml", &newConfig)
	if err != nil {
		log.Panic(err)
	}

	err = cleanenv.ReadEnv(&newConfig)
	if err != nil {
		log.Panic(err)
	}

	globalConfig = &newConfig
}

// Config -.
type Config struct {
	HTTP     `env-required:"true" yaml:"http"`
	Postgres `env-required:"true" yaml:"postgres"`
	Sinopac  `env-required:"true" yaml:"sinopac"`
	RabbitMQ `env-required:"true" yaml:"rabbitmq"`

	Deployment string `env-required:"true" env:"DEPLOYMENT"`
}

// HTTP -.
type HTTP struct {
	Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
}

// Postgres -.
type Postgres struct {
	PoolMax int `env-required:"true" yaml:"pool_max"`

	URL    string `env-required:"true" env:"PG_URL"`
	DBName string `env-required:"true" env:"DB_NAME"`
}

// Sinopac -.
type Sinopac struct {
	PoolMax int `env-required:"true" yaml:"pool_max"`

	URL string `env-required:"true" env:"SINOPAC_URL"`
}

// RabbitMQ -.
type RabbitMQ struct {
	URL      string `env-required:"true" env:"RABBITMQ_URL"`
	Exchange string `env-required:"true" env:"RABBITMQ_EXCHANGE"`

	WaitTime int64 `env-required:"true" yaml:"wait_time"`
	Attempts int   `env-required:"true" yaml:"attempts"`
}
