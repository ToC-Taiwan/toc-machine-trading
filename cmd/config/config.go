// Package config package config
package config

import (
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	singleton *Config
	once      sync.Once
)

// Config -.
type Config struct {
	Simulation   bool         `yaml:"simulation"`
	ManualTrade  bool         `yaml:"manual_trade"`
	History      History      `yaml:"history"`
	Quota        Quota        `yaml:"quota"`
	TargetStock  TargetStock  `yaml:"target_stock"`
	AnalyzeStock AnalyzeStock `yaml:"analyze_stock"`
	TradeStock   TradeStock   `yaml:"trade_stock"`
	TradeFuture  TradeFuture  `yaml:"trade_future"`

	// env must be the last field
	EnvConfig
}

// Get -.
func Get() *Config {
	if singleton != nil {
		return singleton
	}

	once.Do(func() {
		filePath := "configs/config.yml"
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

	if singleton.TradeStock.AllowTrade && !singleton.TradeStock.Subscribe {
		panic("stock trade switch allow trade but not subscribe")
	}

	if singleton.TradeFuture.AllowTrade && !singleton.TradeFuture.Subscribe {
		panic("stock trade switch allow trade but not subscribe")
	}

	return singleton
}
