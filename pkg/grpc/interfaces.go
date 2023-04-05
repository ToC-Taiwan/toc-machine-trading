// Package grpc package grpc
package grpc

type GRPCLogger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}
