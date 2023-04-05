// Package rabbitmq package rabbitmq
package rabbitmq

import "time"

// Option -.
type Option func(*Connection)

func WaitTime(waitTime int) Option {
	return func(c *Connection) {
		c.waitTime = time.Duration(waitTime) * time.Second
	}
}

func Attempts(attempts int) Option {
	return func(c *Connection) {
		c.attempts = attempts
	}
}

func Logger(logger MQLogger) Option {
	return func(c *Connection) {
		c.logger = logger
	}
}
