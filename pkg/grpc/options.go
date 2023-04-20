// Package grpc package grpc
package grpc

import "time"

// Option -.
type Option func(*ConnPool)

// MaxPoolSize -.
func MaxPoolSize(size int) Option {
	return func(c *ConnPool) {
		c.maxPoolSize = size
	}
}

// ConnAttempts -.
func ConnAttempts(attempts int) Option {
	return func(c *ConnPool) {
		c.connAttempts = attempts
	}
}

// ConnTimeout -.
func ConnTimeout(timeout time.Duration) Option {
	return func(c *ConnPool) {
		c.connTimeout = timeout
	}
}

// AddLogger -.
func AddLogger(logger Logger) Option {
	return func(c *ConnPool) {
		c.logger = logger
	}
}
