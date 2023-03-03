package postgres

import "time"

// Option -.
type Option func(*Postgres)

// MaxPoolSize -.
func MaxPoolSize(size int) Option {
	return func(c *Postgres) {
		c.maxPoolSize = size
	}
}

// ConnAttempts -.
func ConnAttempts(attempts int) Option {
	return func(c *Postgres) {
		c.connAttempts = attempts
	}
}

// ConnTimeout -.
func ConnTimeout(timeout time.Duration) Option {
	return func(c *Postgres) {
		c.connTimeout = timeout
	}
}

// MaxConnIdleTime -.
func MaxConnIdleTime(timeout time.Duration) Option {
	return func(c *Postgres) {
		c.maxConnIdleTime = timeout
	}
}

// Logger -.
func Logger(logger PGLogger) Option {
	return func(c *Postgres) {
		c.logger = logger
	}
}
