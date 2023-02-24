// Package grpc package grpc
package grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	_defaultMaxPoolSize  = 10
	_defaultConnAttempts = 10
	_defaultConnTimeout  = 30 * time.Second
)

// Connection -.
type Connection struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	pool      []*grpc.ClientConn
	ReadyConn chan *grpc.ClientConn
}

// New -.
func New(url string, opts ...Option) (*Connection, error) {
	conn := &Connection{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(conn)
	}

	conn.ReadyConn = make(chan *grpc.ClientConn, conn.maxPoolSize)
	for conn.connAttempts > 0 {
		if len(conn.pool) == conn.maxPoolSize {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), conn.connTimeout)
		newConn, err := grpc.DialContext(
			ctx,
			url,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(1024*1024*1024),
				grpc.MaxCallSendMsgSize(1024*1024*1024),
			),
		)
		cancel()
		if err != nil {
			conn.connAttempts--
			if errors.Is(err, context.DeadlineExceeded) {
				fmt.Printf("gRPC connection timeout, attempts left: %d\n", conn.connAttempts)
				continue
			}
			return nil, err
		}

		if newConn.GetState() == connectivity.Ready {
			conn.connAttempts = _defaultConnAttempts
			conn.pool = append(conn.pool, newConn)
			conn.ReadyConn <- newConn
		}
	}

	if len(conn.pool) != conn.maxPoolSize {
		return nil, fmt.Errorf("gRPC connection failed")
	}

	return conn, nil
}

// GetReadyConn -.
func (s *Connection) GetReadyConn() *grpc.ClientConn {
	var conn *grpc.ClientConn
	for {
		r := <-s.ReadyConn
		if r != nil {
			conn = r
			break
		}
	}
	return conn
}

// PutReadyConn -.
func (s *Connection) PutReadyConn(conn *grpc.ClientConn) {
	s.ReadyConn <- conn
}
