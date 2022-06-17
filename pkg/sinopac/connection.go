// Package sinopac package sinopac
package sinopac

import (
	"time"

	"toc-machine-trading/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var log = logger.Get()

const (
	_defaultMaxPoolSize  = 10
	_defaultConnAttempts = 10
	_defaultConnTimeout  = 3 * time.Second
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

	// Custom options
	for _, opt := range opts {
		opt(conn)
	}

	conn.ReadyConn = make(chan *grpc.ClientConn, conn.maxPoolSize)

	var newConn *grpc.ClientConn
	var err error
	for conn.connAttempts > 0 {
		if len(conn.pool) == conn.maxPoolSize {
			err = nil
			break
		}
		newConn, err = grpc.Dial(
			url,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(1024*1024*1024),
				grpc.MaxCallSendMsgSize(1024*1024*1024),
			),
		)
		if err == nil && newConn != nil {
			conn.pool = append(conn.pool, newConn)
			conn.ReadyConn <- newConn
			continue
		}

		log.Infof("gRPC trying connect, attempts left: %d", conn.connAttempts)
		time.Sleep(conn.connTimeout)
		conn.connAttempts--
	}

	if err != nil {
		return nil, err
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
