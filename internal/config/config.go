// Package config package config
package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/toc-taiwan/toc-machine-trading/pkg/grpc"
	"github.com/toc-taiwan/toc-machine-trading/pkg/log"
	"github.com/toc-taiwan/toc-machine-trading/pkg/postgres"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Config -.
type Config struct {
	Simulation   bool         `json:"Simulation" yaml:"Simulation"`
	ManualTrade  bool         `json:"ManualTrade" yaml:"ManualTrade"`
	TradeStock   TradeStock   `json:"TradeStock" yaml:"TradeStock"`
	History      History      `json:"History" yaml:"History"`
	Quota        Quota        `json:"Quota" yaml:"Quota"`
	AnalyzeStock AnalyzeStock `json:"AnalyzeStock" yaml:"AnalyzeStock"`
	TradeFuture  TradeFuture  `json:"TradeFuture" yaml:"TradeFuture"`

	dbPool      *postgres.Postgres `json:"-" yaml:"-"`
	sinopacPool *grpc.ConnPool     `json:"-" yaml:"-"`
	EnvConfig   `json:"-" yaml:"-"`

	logger   *log.Log     `json:"-" yaml:"-"`
	vp       *viper.Viper `json:"-" yaml:"-"`
	basePath string       `json:"-" yaml:"-"`
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
	c.vp.SetDefault("DB_NAME", "machine_trade")
	c.vp.SetDefault("DB_URL", "postgres://postgres:password@127.0.0.1:5432/")
	c.vp.SetDefault("DB_POOL_MAX", 80)

	c.vp.SetDefault("HTTP", "26670")

	c.vp.SetDefault("SINOPAC_POOL_MAX", 20)
	c.vp.SetDefault("SINOPAC_URL", "127.0.0.1:56666")

	c.vp.AutomaticEnv()
	env := EnvConfig{
		Database: Database{
			DBName:  c.vp.GetString("DB_NAME"),
			URL:     c.vp.GetString("DB_URL"),
			PoolMax: c.vp.GetInt("DB_POOL_MAX"),
		},
		Server: Server{
			HTTP: c.vp.GetString("HTTP"),
		},
		Sinopac: Sinopac{
			PoolMax: c.vp.GetInt("SINOPAC_POOL_MAX"),
			URL:     c.vp.GetString("SINOPAC_URL"),
		},
		SMTP: SMTP{
			Host:     c.vp.GetString("SMTP_HOST"),
			Port:     c.vp.GetInt("SMTP_PORT"),
			Username: c.vp.GetString("SMTP_USERNAME"),
			Password: c.vp.GetString("SMTP_PASSWORD"),
		},
	}
	c.EnvConfig = env
}

func Init() {
	once.Do(func() {
		data := newConfig()
		data.readConfig()
		data.readEnv()
		data.setPostgresPool()
		data.setSinopacPool()
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

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

func (c *Config) createDB() {
	pg, err := postgres.New(
		c.Database.URL,
		postgres.MaxPoolSize(c.Database.PoolMax),
		postgres.AddLogger(c.logger),
	)
	if err != nil {
		c.logger.Fatalf("postgres create db error: %s", err)
	}
	defer pg.Close()

	var name string
	if err := pg.Pool().QueryRow(context.Background(),
		"SELECT datname FROM pg_catalog.pg_database WHERE datname = $1", c.Database.DBName).
		Scan(&name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			_, err = pg.Pool().Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", c.Database.DBName))
			if err != nil {
				c.logger.Fatalf("postgres create db error: %s", err)
			}
			return
		}
		c.logger.Fatalf("postgres create db error: %s", err)
	}
}

func (c *Config) migrateScheme() {
	m := &migrate.Migrate{}

	path := fmt.Sprintf("%s%s%s", c.Database.URL, c.Database.DBName, "?sslmode=disable")
	attempts := _defaultAttempts
	var err error
	for attempts > 0 {
		m, err = migrate.New("file://migrations", path)
		if err == nil {
			break
		}

		c.logger.Infof("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		c.logger.Fatal(fmt.Errorf("postgres connect error in migrate: %s", err))
	}

	defer func() {
		_, _ = m.Close()
	}()
	err = m.Up()
	if err != nil {
		switch err {
		case migrate.ErrNoChange:
			c.logger.Info("Migrate: no change")
		default:
			c.logger.Errorf("Migrate: up error: %s", err)
		}
		return
	}
	c.logger.Info("Migrate: up success")
}

func (c *Config) setPostgresPool() {
	c.createDB()
	c.migrateScheme()

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

func (c *Config) CloseDB() {
	if c.dbPool != nil {
		c.dbPool.Close()
	}
}
