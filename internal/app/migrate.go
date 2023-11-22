package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"tmt/pkg/log"

	"tmt/internal/config"
	"tmt/pkg/postgres"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v4"

	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

func InitDB() {
	cfg := config.Get().Database
	logger := log.Get()
	if e := createDB(cfg); e != nil {
		logger.Fatal(fmt.Errorf("postgres create db error: %s", e.Error()))
	}
	migrateScheme(cfg)
}

func createDB(cfg config.Database) error {
	pg, err := postgres.New(
		cfg.URL,
		postgres.MaxPoolSize(cfg.PoolMax),
		postgres.AddLogger(log.Get()),
	)
	if err != nil {
		return err
	}
	defer pg.Close()

	var name string
	if err := pg.Pool().QueryRow(context.Background(),
		"SELECT datname FROM pg_catalog.pg_database WHERE datname = $1", cfg.DBName).
		Scan(&name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			_, err = pg.Pool().Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}

func migrateScheme(dbConfig config.Database) {
	m := &migrate.Migrate{}
	logger := log.Get()

	path := fmt.Sprintf("%s%s%s", dbConfig.URL, dbConfig.DBName, "?sslmode=disable")
	attempts := _defaultAttempts
	var err error
	for attempts > 0 {
		m, err = migrate.New("file://migrations", path)
		if err == nil {
			break
		}

		logger.Infof("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		logger.Fatal(fmt.Errorf("postgres connect error in migrate: %s", err))
	}

	defer func() {
		_, _ = m.Close()
	}()
	err = m.Up()
	if err != nil {
		switch err {
		case migrate.ErrNoChange:
			logger.Info("Migrate: no change")
		default:
			logger.Errorf("Migrate: up error: %s", err)
		}
		return
	}
	logger.Info("Migrate: up success")
}
