// Package config package config
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"tmt/pkg/grpc"
	"tmt/pkg/log"
	"tmt/pkg/postgres"
	"tmt/pkg/rabbitmq"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

var (
	singleton *Config
	logger    = log.Get()
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

	dbPool      *postgres.Postgres
	sinopacPool *grpc.ConnPool
	fuglePool   *grpc.ConnPool
}

// Get -.
func Get() *Config {
	if singleton != nil {
		return singleton
	}

	filePath := "configs/config.yml"
	fileStat, err := os.Stat(filePath)
	if err != nil || fileStat.IsDir() {
		logger.Fatalf("config file not found: %v", err)
	}

	newConfig := Config{}
	if fileStat.Size() > 0 {
		err = cleanenv.ReadConfig(filePath, &newConfig)
		if err != nil {
			logger.Fatalf("config file read error: %v", err)
		}
	}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	err = godotenv.Load(filepath.Join(filepath.Dir(ex), ".env"))
	if err != nil {
		panic(err)
	}

	if err := cleanenv.ReadEnv(&newConfig); err != nil {
		logger.Fatalf("config env read error: %v", err)
	}

	if newConfig.TradeStock.AllowTrade && !newConfig.TradeStock.Subscribe {
		logger.Fatal("stock trade switch allow trade but not subscribe")
	}

	if newConfig.TradeFuture.AllowTrade && !newConfig.TradeFuture.Subscribe {
		logger.Fatal("stock trade switch allow trade but not subscribe")
	}

	if e := newConfig.setupCronJob(); e != nil {
		logger.Fatal(e)
	}

	singleton = &newConfig
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
		postgres.AddLogger(logger),
	)
	if err != nil {
		logger.Fatal(err)
	}
	c.dbPool = pg
	return pg
}

func (c *Config) GetSinopacPool() *grpc.ConnPool {
	if c.sinopacPool != nil {
		return c.sinopacPool
	}
	logger.Info("Connecting to sinopac gRPC server")
	sc, err := grpc.New(
		c.Sinopac.URL,
		grpc.MaxPoolSize(c.Sinopac.PoolMax),
		grpc.AddLogger(logger),
	)
	if err != nil {
		logger.Fatal(err)
	}
	c.sinopacPool = sc
	return sc
}

func (c *Config) GetFuglePool() *grpc.ConnPool {
	if c.fuglePool != nil {
		return c.fuglePool
	}
	logger.Info("Connecting to fugle gRPC server")
	fg, err := grpc.New(
		c.Fugle.URL,
		grpc.MaxPoolSize(c.Fugle.PoolMax),
		grpc.AddLogger(logger),
	)
	if err != nil {
		logger.Fatal(err)
	}
	c.fuglePool = fg
	return fg
}

func (c *Config) GetRabbitConn() *rabbitmq.Connection {
	conn, err := rabbitmq.New(
		c.RabbitMQ.Exchange, c.RabbitMQ.URL,
		rabbitmq.Attempts(c.RabbitMQ.Attempts),
		rabbitmq.WaitTime(int(c.RabbitMQ.WaitTime)),
		rabbitmq.AddLogger(logger),
	)
	if err != nil {
		logger.Fatal(err)
	}
	return conn
}

func (c *Config) CloseDB() {
	if c.dbPool != nil {
		c.dbPool.Close()
		logger.Warn("TMT is shutting down")
	}
}
