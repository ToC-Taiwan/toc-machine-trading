// Package postgres implements postgres connection.
package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	_defaultMaxPoolSize     = 1
	_defaultConnAttempts    = 10
	_defaultConnTimeout     = 3 * time.Second
	_defaultMaxConnIdleTime = 3 * time.Minute
)

// Postgres -.
type Postgres struct {
	maxPoolSize     int
	connAttempts    int
	connTimeout     time.Duration
	maxConnIdleTime time.Duration

	Builder squirrel.StatementBuilderType
	pool    *pgxpool.Pool

	logger PGLogger
}

// New -.
func New(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:     _defaultMaxPoolSize,
		connAttempts:    _defaultConnAttempts,
		connTimeout:     _defaultConnTimeout,
		maxConnIdleTime: _defaultMaxConnIdleTime,
	}

	for _, opt := range opts {
		opt(pg)
	}

	// builder
	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// pool
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	poolConfig.MaxConnIdleTime = pg.maxConnIdleTime
	poolConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connAttempts > 0 {
		pg.pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if err == nil {
			return pg, nil
		}

		pg.Infof("Postgres trying connect, attempts left: %d\n", pg.connAttempts)

		pg.connAttempts--
		time.Sleep(pg.connTimeout)
	}
	return nil, errors.New("Postgres: connection attempts exceeded")
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
	tx, err := p.Pool().Begin(context.Background())
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
			p.Errorf("Postgres: error on rollback transaction: %s\n", rollErr)
		}
	} else {
		commitErr := tx.Commit(context.Background())
		if commitErr != nil {
			p.Errorf("Postgres: error on commit transaction: %s\n", commitErr)
		}
	}
}

func (p *Postgres) Infof(format string, args ...interface{}) {
	if p.logger != nil {
		p.logger.Infof(format, args...)
	} else {
		fmt.Printf(format, args...)
	}
}

func (p *Postgres) Errorf(format string, args ...interface{}) {
	if p.logger != nil {
		p.logger.Errorf(format, args...)
	} else {
		fmt.Printf(format, args...)
	}
}
