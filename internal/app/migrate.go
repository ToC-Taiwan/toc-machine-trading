package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"toc-machine-trading/pkg/config"
	"toc-machine-trading/pkg/logger"
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
func MigrateDB() {
	databaseURL, ok := os.LookupEnv("PG_URL")
	if !ok || len(databaseURL) == 0 {
		logger.Get().Fatal("environment variable not declared: PG_URL")
	}

	databaseName, ok := os.LookupEnv("DB_NAME")
	if !ok || len(databaseURL) == 0 {
		logger.Get().Fatal("environment variable not declared: DB_NAME")
	}

	createErr := createDB(databaseName)
	if createErr != nil {
		logger.Get().Panic(createErr)
	}

	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	dbPath := fmt.Sprintf("%s%s%s", databaseURL, databaseName, "?sslmode=disable")
	for attempts > 0 {
		m, err = migrate.New("file://migrations", dbPath)
		if err == nil {
			break
		}

		logger.Get().Infof("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		logger.Get().Fatalf("Migrate: postgres connect error: %s", err)
	}

	err = m.Up()
	defer func() {
		_, _ = m.Close()
	}()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Get().Infof("Migrate: up error: %s", err)
		return
	}

	if errors.Is(err, migrate.ErrNoChange) {
		logger.Get().Info("Migrate: no change")
		return
	}

	logger.Get().Info("Migrate: up success")
}

func createDB(dbName string) error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
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
