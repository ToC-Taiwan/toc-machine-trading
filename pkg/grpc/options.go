// Package grpc package grpc
package grpc

import "time"

// Option -.
type Option func(*Connection)

// MaxPoolSize -.
func MaxPoolSize(size int) Option {
	return func(c *Connection) {
		c.maxPoolSize = size
	}
}

// ConnAttempts -.
func ConnAttempts(attempts int) Option {
	return func(c *Connection) {
		c.connAttempts = attempts
	}
}

// ConnTimeout -.
func ConnTimeout(timeout time.Duration) Option {
	return func(c *Connection) {
		c.connTimeout = timeout
	}
}
