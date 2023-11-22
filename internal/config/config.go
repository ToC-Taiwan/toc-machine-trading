// Package config package config
package config

import (
	"fmt"
	"os"
	"sync"

	"tmt/pkg/grpc"
	"tmt/pkg/log"
	"tmt/pkg/postgres"
	"tmt/pkg/rabbitmq"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/robfig/cron/v3"
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

	dbPool      *postgres.Postgres
	sinopacPool *grpc.ConnPool
	fuglePool   *grpc.ConnPool

	logger *log.Log

	// env must be the last field
	EnvConfig
}

var (
	singleton *Config
	initOnce  sync.Once
)

// Get -.
func Get() *Config {
	if singleton != nil {
		return singleton
	}

	initOnce.Do(func() {
		cfg := Config{
			logger: log.Get(),
		}

		filePath := "configs/config.yml"
		fileStat, err := os.Stat(filePath)
		if err != nil || fileStat.IsDir() {
			cfg.logger.Fatalf("config file not found: %v", err)
		}

		if fileStat.Size() > 0 {
			err = cleanenv.ReadConfig(filePath, &cfg)
			if err != nil {
				cfg.logger.Fatalf("config file read error: %v", err)
			}
		}

		if err := cleanenv.ReadEnv(&cfg); err != nil {
			cfg.logger.Fatalf("config env read error: %v", err)
		}

		if cfg.TradeStock.AllowTrade && !cfg.TradeStock.Subscribe {
			cfg.logger.Fatalf("stock trade switch allow trade but not subscribe")
		}

		if cfg.TradeFuture.AllowTrade && !cfg.TradeFuture.Subscribe {
			cfg.logger.Fatalf("stock trade switch allow trade but not subscribe")
		}

		if e := cfg.setupCronJob(); e != nil {
			cfg.logger.Fatalf(e.Error())
		}

		singleton = &cfg
	})

	return singleton
}

func (c *Config) setupCronJob() error {
	job := cron.New()
	if _, e := job.AddFunc("20 8 * * *", c.exit); e != nil {
		return e
	}
	if _, e := job.AddFunc("40 14 * * *", c.exit); e != nil {
		return e
	}
	job.Start()
	return nil
}

func (c *Config) exit() {
	os.Exit(0)
}

func (c *Config) GetPostgresPool() *postgres.Postgres {
	if c.dbPool != nil {
		return c.dbPool
	}
	pg, err := postgres.New(
		fmt.Sprintf("%s%s", c.Database.URL, c.Database.DBName),
		postgres.MaxPoolSize(c.Database.PoolMax),
		postgres.AddLogger(c.logger),
	)
	if err != nil {
		c.logger.Fatal(err)
	}
	c.dbPool = pg
	return pg
}

func (c *Config) GetSinopacPool() *grpc.ConnPool {
	if c.sinopacPool != nil {
		return c.sinopacPool
	}
	c.logger.Info("Connecting to sinopac gRPC server")
	sc, err := grpc.New(
		c.Sinopac.URL,
		grpc.MaxPoolSize(c.Sinopac.PoolMax),
		grpc.AddLogger(c.logger),
	)
	if err != nil {
		c.logger.Fatal(err)
	}
	c.sinopacPool = sc
	return sc
}

func (c *Config) GetFuglePool() *grpc.ConnPool {
	if c.fuglePool != nil {
		return c.fuglePool
	}
	c.logger.Info("Connecting to fugle gRPC server")
	fg, err := grpc.New(
		c.Fugle.URL,
		grpc.MaxPoolSize(c.Fugle.PoolMax),
		grpc.AddLogger(c.logger),
	)
	if err != nil {
		c.logger.Fatal(err)
	}
	c.fuglePool = fg
	return fg
}

func (c *Config) GetRabbitConn() *rabbitmq.Connection {
	conn, err := rabbitmq.New(
		c.RabbitMQ.Exchange, c.RabbitMQ.URL,
		rabbitmq.Attempts(c.RabbitMQ.Attempts),
		rabbitmq.WaitTime(int(c.RabbitMQ.WaitTime)),
		rabbitmq.AddLogger(c.logger),
	)
	if err != nil {
		c.logger.Fatal(err)
	}
	return conn
}

func (c *Config) CloseDB() {
	if c.dbPool != nil {
		c.dbPool.Close()
		c.logger.Warn("TMT is shutting down")
	}
}
