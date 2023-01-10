// Package config package config
package config

import (
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	once      sync.Once
	singleton *Config
)

// Config -.
type Config struct {
	Simulation bool

	Database Database
	Server   Server
	Sinopac  Sinopac
	RabbitMQ RabbitMQ

	StockTradeSwitch  StockTradeSwitch  `yaml:"stock_trade_switch"`
	FutureTradeSwitch FutureTradeSwitch `yaml:"future_trade_switch"`
	History           History           `yaml:"history"`
	Quota             Quota             `yaml:"quota"`
	TargetCond        TargetCond        `yaml:"target_cond"`
	StockAnalyze      StockAnalyze      `yaml:"stock_analyze"`
	FutureAnalyze     FutureAnalyze     `yaml:"future_analyze"`
}

// GetConfig -.
func GetConfig() *Config {
	if singleton != nil {
		return singleton
	}

	once.Do(func() {
		filePath := "./configs/config.yml"
		fileStat, err := os.Stat(filePath)
		if err != nil || fileStat.IsDir() {
			panic(err)
		}

		newConfig := Config{}
		if fileStat.Size() > 0 {
			err := cleanenv.ReadConfig(filePath, &newConfig)
			if err != nil {
				panic(err)
			}
		}

		if err := cleanenv.ReadEnv(&newConfig); err != nil {
			panic(err)
		}

		singleton = &newConfig
	})

	if singleton.StockTradeSwitch.AllowTrade && !singleton.StockTradeSwitch.Subscribe {
		panic("stock trade switch allow trade but not subscribe")
	}

	if singleton.FutureTradeSwitch.AllowTrade && !singleton.FutureTradeSwitch.Subscribe {
		panic("stock trade switch allow trade but not subscribe")
	}

	return singleton
}
