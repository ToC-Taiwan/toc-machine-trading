// Package grpc package grpc
package grpc

// Option -.
type Option func(*ConnPool)

// MaxPoolSize -.
func MaxPoolSize(size int) Option {
	return func(c *ConnPool) {
		c.maxPoolSize = size
	}
}

// AddLogger -.
func AddLogger(logger Logger) Option {
	return func(c *ConnPool) {
		c.logger = logger
	}
}
