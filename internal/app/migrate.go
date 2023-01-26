package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"tmt/cmd/config"
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

// InitDB -.
func InitDB(cfg config.Database) {
	TryCreateDB(cfg)
	MigrateDB(cfg)
}

// MigrateDB -.
func MigrateDB(dbConfig config.Database) {
	m := &migrate.Migrate{}
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

	err = m.Up()
	defer func() {
		_, _ = m.Close()
	}()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Infof("Migrate: up error: %s", err)
		return
	}

	if errors.Is(err, migrate.ErrNoChange) {
		logger.Info("Migrate: no change")
		return
	}

	logger.Info("Migrate: up success")
}

// TryCreateDB -.
func TryCreateDB(cfg config.Database) {
	pg, err := postgres.New(cfg.URL, postgres.MaxPoolSize(cfg.PoolMax))
	if err != nil {
		logger.Fatal(err)
	}
	defer pg.Close()

	var name string
	if err := pg.Pool().QueryRow(context.Background(), "SELECT datname FROM pg_catalog.pg_database WHERE datname = $1", cfg.DBName).Scan(&name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			_, err = pg.Pool().Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
			if err != nil {
				logger.Fatal(err)
			}
			return
		}
		logger.Fatal(err)
	}
}
