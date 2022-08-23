// Package postgres implements postgres connection.
package postgres

import (
	"context"
	"fmt"
	"time"

	"tmt/pkg/logger"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var log = logger.Get()

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// Postgres -.
type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Builder squirrel.StatementBuilderType
	pool    *pgxpool.Pool
}

// New -.
func New(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(pg)
	}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConnIdleTime = time.Second
	poolConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connAttempts > 0 {
		pg.pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Infof("Postgres trying connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, err
	}

	return pg, nil
}

// Close -.
func (p *Postgres) Close() {
	if p.pool != nil {
		p.pool.Close()
	}
}

// Pool -.
func (p *Postgres) Pool() *pgxpool.Pool {
	return p.pool
}

// BeginTransaction -.
func (p *Postgres) BeginTransaction() (pgx.Tx, error) {
	tx, err := p.pool.Begin(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Postgres: error on begin transaction: %s", err)
	}
	return tx, nil
}

// EndTransaction -.
func (p *Postgres) EndTransaction(tx pgx.Tx, err error) {
	if err != nil {
		rollErr := tx.Rollback(context.Background())
		if rollErr != nil {
			log.Errorf("Postgres: error on rollback transaction: %s", rollErr)
		}
	} else {
		commitErr := tx.Commit(context.Background())
		if commitErr != nil {
			log.Errorf("Postgres: error on commit transaction: %s", commitErr)
		}
	}
}
