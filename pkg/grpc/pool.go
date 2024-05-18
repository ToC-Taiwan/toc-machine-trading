// Package grpc package grpc
package grpc

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	_defaultMaxPoolSize  = 10
	_defaultConnAttempts = 10
	_defaultConnTimeout  = 10 * time.Second
)

// ConnPool -.
type ConnPool struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	pool      []*grpc.ClientConn
	readyConn chan *grpc.ClientConn

	logger Logger
}

// New -.
func New(url string, opts ...Option) (*ConnPool, error) {
	conn := &ConnPool{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(conn)
	}

	conn.readyConn = make(chan *grpc.ClientConn, conn.maxPoolSize)
	for conn.connAttempts > 0 {
		if len(conn.pool) == conn.maxPoolSize {
			break
		}
		newConn, err := grpc.NewClient(
			url,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(1024*1024*1024),
				grpc.MaxCallSendMsgSize(1024*1024*1024),
			),
		)
		if err != nil {
			conn.connAttempts--
			if errors.Is(err, context.DeadlineExceeded) {
				conn.Infof("gRPC connection timeout, attempts left: %d\n", conn.connAttempts)
				continue
			}
			return nil, err
		}
		conn.connAttempts = _defaultConnAttempts
		conn.pool = append(conn.pool, newConn)
		conn.readyConn <- newConn
	}

	if len(conn.pool) != conn.maxPoolSize {
		return nil, fmt.Errorf("gRPC connection failed")
	}

	return conn, nil
}

// Get -.
func (s *ConnPool) Get() *grpc.ClientConn {
	return <-s.readyConn
}

// Put -.
func (s *ConnPool) Put(conn *grpc.ClientConn) {
	s.readyConn <- conn
}

func (s *ConnPool) Infof(format string, args ...interface{}) {
	if s.logger != nil {
		s.logger.Infof(strings.ReplaceAll(format, "\n", ""), args...)
	} else {
		fmt.Printf(format, args...)
	}
}

func (s *ConnPool) Errorf(format string, args ...interface{}) {
	if s.logger != nil {
		s.logger.Errorf(strings.ReplaceAll(format, "\n", ""), args...)
	} else {
		fmt.Printf(format, args...)
	}
}
