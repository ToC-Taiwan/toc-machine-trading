// Package rabbitmq package rabbitmq
package rabbitmq

type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}
