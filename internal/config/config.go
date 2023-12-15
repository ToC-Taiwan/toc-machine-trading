// Package config package config
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"tmt/pkg/grpc"
	"tmt/pkg/log"
	"tmt/pkg/postgres"
	"tmt/pkg/rabbitmq"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config -.
type Config struct {
	Simulation   bool         `yaml:"Simulation"`
	ManualTrade  bool         `yaml:"ManualTrade"`
	TradeStock   TradeStock   `yaml:"TradeStock"`
	History      History      `yaml:"History"`
	Quota        Quota        `yaml:"Quota"`
	TargetStock  TargetStock  `yaml:"TargetStock"`
	AnalyzeStock AnalyzeStock `yaml:"AnalyzeStock"`
	TradeFuture  TradeFuture  `yaml:"TradeFuture"`

	dbPool      *postgres.Postgres
	sinopacPool *grpc.ConnPool
	fuglePool   *grpc.ConnPool
	EnvConfig

	logger   *log.Log
	vp       *viper.Viper
	basePath string
}

var (
	singleton *Config
	once      sync.Once
)

func newConfig() *Config {
	logger := log.Get()
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	if err = godotenv.Load(filepath.Join(filepath.Dir(ex), ".env")); err != nil {
		if os.IsNotExist(err) {
			logger.Warn("No .env file found, using default env values")
		}
	}
	return &Config{
		logger:   logger,
		vp:       viper.New(),
		basePath: filepath.Dir(ex),
	}
}

func (c *Config) readConfig() {
	c.vp.SetConfigName("config")
	c.vp.SetConfigType("yml")
	c.vp.AddConfigPath(filepath.Join(c.basePath, "configs"))
	err := c.vp.ReadInConfig()
	if err != nil {
		c.logger.Fatalf("fatal error config file: %s", err)
	}
	err = c.vp.Unmarshal(c)
	if err != nil {
		c.logger.Fatalf("fatal error config file: %s", err)
	}
}

func (c *Config) readEnv() {
	c.vp.SetDefault("HTTP", "26670")
	c.vp.SetDefault("DISABLE_SWAGGER_HTTP_HANDLER", "")

	c.vp.SetDefault("LOG_LEVEL", "info")
	c.vp.SetDefault("LOG_FORMAT", "text")
	c.vp.SetDefault("LOG_NEED_CALLER", false)
	c.vp.SetDefault("LOG_TIME_FORMAT", "")
	c.vp.SetDefault("LOG_LINK_SLACK", false)
	c.vp.SetDefault("LOG_DISABLE_CONSOLE", false)
	c.vp.SetDefault("LOG_DISABLE_FILE", false)

	c.vp.SetDefault("SINOPAC_URL", "127.0.0.1:56666")
	c.vp.SetDefault("SINOPAC_POOL_MAX", 20)
	c.vp.SetDefault("FUGLE_URL", "127.0.0.1:58888")
	c.vp.SetDefault("FUGLE_POOL_MAX", 20)

	c.vp.SetDefault("DB_URL", "postgres://postgres:asdf0000@127.0.0.1:5432/")
	c.vp.SetDefault("DB_NAME", "machine_trade")
	c.vp.SetDefault("DB_POOL_MAX", 80)

	c.vp.SetDefault("RABBITMQ_URL", "amqp://admin:password@127.0.0.1:5672/%2f?heartbeat=0")
	c.vp.SetDefault("RABBITMQ_EXCHANGE", "toc")
	c.vp.SetDefault("RABBITMQ_WAIT_TIME", 5)
	c.vp.SetDefault("RABBITMQ_ATTEMPTS", 10)

	c.vp.SetDefault("SLACK_TOKEN", "")
	c.vp.SetDefault("SLACK_CHANNEL_ID", "")
	c.vp.SetDefault("SLACK_LOG_LEVEL", "warn")
	c.vp.AutomaticEnv()
	env := EnvConfig{
		Database: Database{
			DBName:  c.vp.GetString("DB_NAME"),
			URL:     c.vp.GetString("DB_URL"),
			PoolMax: c.vp.GetInt("DB_POOL_MAX"),
		},
		Server: Server{
			HTTP:                      c.vp.GetString("HTTP"),
			DisableSwaggerHTTPHandler: c.vp.GetString("DISABLE_SWAGGER_HTTP_HANDLER"),
		},
		Sinopac: Sinopac{
			PoolMax: c.vp.GetInt("SINOPAC_POOL_MAX"),
			URL:     c.vp.GetString("SINOPAC_URL"),
		},
		Fugle: Fugle{
			PoolMax: c.vp.GetInt("FUGLE_POOL_MAX"),
			URL:     c.vp.GetString("FUGLE_URL"),
		},
		RabbitMQ: RabbitMQ{
			URL:      c.vp.GetString("RABBITMQ_URL"),
			Exchange: c.vp.GetString("RABBITMQ_EXCHANGE"),
			WaitTime: c.vp.GetInt64("RABBITMQ_WAIT_TIME"),
			Attempts: c.vp.GetInt("RABBITMQ_ATTEMPTS"),
		},
	}
	c.EnvConfig = env
}

func (c *Config) checkValid() {
	if c.TradeStock.AllowTrade && !c.TradeStock.Subscribe {
		c.logger.Fatalf("stock trade switch allow trade but not subscribe")
	}

	if c.TradeFuture.AllowTrade && !c.TradeFuture.Subscribe {
		c.logger.Fatalf("stock trade switch allow trade but not subscribe")
	}
}

func Init() {
	once.Do(func() {
		data := newConfig()
		data.readConfig()
		data.readEnv()
		data.checkValid()
		data.setPostgresPool()
		data.setSinopacPool()
		data.setFuglePool()
		singleton = data
	})
}

// Get -.
func Get() *Config {
	if singleton == nil {
		once.Do(Init)
		return Get()
	}
	return singleton
}

func (c *Config) setPostgresPool() {
	pg, err := postgres.New(
		fmt.Sprintf("%s%s", c.Database.URL, c.Database.DBName),
		postgres.MaxPoolSize(c.Database.PoolMax),
		postgres.AddLogger(c.logger),
	)
	if err != nil {
		c.logger.Fatal(err)
	}
	c.dbPool = pg
}

func (c *Config) GetPostgresPool() *postgres.Postgres {
	if c.dbPool == nil {
		c.logger.Fatal("postgres not connected")
	}
	return c.dbPool
}

func (c *Config) setSinopacPool() {
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
}

func (c *Config) GetSinopacPool() *grpc.ConnPool {
	if c.sinopacPool == nil {
		c.logger.Fatal("sinopac gRPC server not connected")
	}
	return c.sinopacPool
}

func (c *Config) setFuglePool() {
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
}

func (c *Config) GetFuglePool() *grpc.ConnPool {
	if c.fuglePool == nil {
		c.logger.Fatal("fugle gRPC server not connected")
	}
	return c.fuglePool
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
	}
}
