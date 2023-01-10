package app

import (
	"context"
	"errors"
	"fmt"
	"log"
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

// MigrateDB -.
func MigrateDB(cfg *config.Config) {
	createErr := tryCreateDB(cfg.Database.DBName)
	if createErr != nil {
		logger.Panic(createErr)
	}

	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	dbPath := fmt.Sprintf("%s%s%s", cfg.Database.URL, cfg.Database.DBName, "?sslmode=disable")
	for attempts > 0 {
		m, err = migrate.New("file://migrations", dbPath)
		if err == nil {
			break
		}

		logger.Infof("Migrate: postgres is trying to connect, attempts left: %d", attempts)
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
		logger.Infof("Migrate: up error: %s", err)
		return
	}

	if errors.Is(err, migrate.ErrNoChange) {
		logger.Info("Migrate: no change")
		return
	}

	logger.Info("Migrate: up success")
}

func tryCreateDB(dbName string) error {
	cfg := config.GetConfig()
	pg, err := postgres.New(cfg.Database.URL, postgres.MaxPoolSize(cfg.Database.PoolMax))
	if err != nil {
		return err
	}
	defer pg.Close()

	var name string
	err = pg.Pool().QueryRow(context.Background(), "SELECT datname FROM pg_catalog.pg_database WHERE datname = $1", dbName).Scan(&name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			_, err = pg.Pool().Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", dbName))
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}
