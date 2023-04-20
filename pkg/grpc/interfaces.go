// Package grpc package grpc
package grpc

type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}
