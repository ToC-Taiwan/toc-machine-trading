// Package rabbitmq package rabbitmq
package rabbitmq

type MQLogger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}
