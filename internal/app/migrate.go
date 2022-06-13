package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"toc-machine-trading/pkg/config"
	"toc-machine-trading/pkg/postgres"

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

// MigrateDB -.
func MigrateDB(cfg *config.Config) {
	createErr := tryCreateDB(cfg.Postgres.DBName)
	if createErr != nil {
		log.Panic(createErr)
	}

	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	dbPath := fmt.Sprintf("%s%s%s", cfg.Postgres.URL, cfg.Postgres.DBName, "?sslmode=disable")
	for attempts > 0 {
		m, err = migrate.New("file://migrations", dbPath)
		if err == nil {
			break
		}

		log.Infof("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("Migrate: postgres connect error: %s", err)
	}

	err = m.Up()
	defer func() {
		_, _ = m.Close()
	}()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Infof("Migrate: up error: %s", err)
		return
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Info("Migrate: no change")
		return
	}

	log.Info("Migrate: up success")
}

func tryCreateDB(dbName string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	pg, err := postgres.New(cfg.Postgres.URL, postgres.MaxPoolSize(cfg.Postgres.PoolMax))
	if err != nil {
		return err
	}
	defer pg.Close()

	var name string
	err = pg.Pool.QueryRow(context.Background(), "SELECT datname FROM pg_catalog.pg_database WHERE datname = $1", dbName).Scan(&name)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	if name == "" {
		_, err := pg.Pool.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return err
		}
	}
	return nil
}
