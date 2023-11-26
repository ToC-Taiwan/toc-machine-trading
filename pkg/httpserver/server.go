// Package httpserver implements HTTP server.
package httpserver

import (
	"fmt"
	"io"
	olog "log"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	localHost                 = "127.0.0.1"
	_defaultPort              = "80"
	_defaultReadTimeout       = 5 * time.Second
	_defaultReadHeaderTimeout = 5 * time.Second
	_defaultWriteTimeout      = 5 * time.Minute
	_defaultShutdownTimeout   = 3 * time.Second
)

// Server -.
type Server struct {
	srv      *http.Server
	logger   Logger
	keyPath  string
	certPath string
}

// New -.
func New(handler http.Handler, opts ...Option) *Server {
	s := &Server{
		srv: &http.Server{
			ErrorLog:          olog.New(io.Discard, "", 0),
			Handler:           handler,
			Addr:              net.JoinHostPort("", _defaultPort),
			ReadHeaderTimeout: _defaultReadHeaderTimeout,
			ReadTimeout:       _defaultReadTimeout,
			WriteTimeout:      _defaultWriteTimeout,
		},
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) Infof(format string, args ...interface{}) {
	if s.logger != nil {
		s.logger.Infof(strings.ReplaceAll(format, "\n", ""), args...)
	} else {
		fmt.Printf(format, args...)
	}
}

func (s *Server) Errorf(format string, args ...interface{}) {
	if s.logger != nil {
		s.logger.Errorf(strings.ReplaceAll(format, "\n", ""), args...)
	} else {
		fmt.Printf(format, args...)
	}
}

func (s *Server) Fatalf(format string, args ...interface{}) {
	if s.logger != nil {
		s.logger.Fatalf(strings.ReplaceAll(format, "\n", ""), args...)
	} else {
		panic(fmt.Errorf(format, args...))
	}
}

func (s *Server) Start() error {
	if err := s.tryListen(); err != nil {
		return err
	}
	return nil
}

func (s *Server) tryListen() error {
	errChan := make(chan error)
	go func() {
		if s.certPath == "" || s.keyPath == "" {
			err := s.srv.ListenAndServe()
			if err != nil {
				errChan <- err
			}
			return
		}
		err := s.srv.ListenAndServeTLS(s.certPath, s.keyPath)
		if err != nil {
			errChan <- err
		}
	}()

	for {
		select {
		case err := <-errChan:
			return err
		case <-time.After(1 * time.Second):
			if s.getPortIsUsed(localHost, s.srv.Addr[1:]) {
				s.Infof("API Server On %v", s.srv.Addr)
				return nil
			}
		}
	}
}

// GetPortIsUsed GetPortIsUsed
func (s *Server) getPortIsUsed(host string, port string) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), 500*time.Millisecond)
	if err != nil && conn != nil {
		return false
	}

	if conn != nil {
		defer func() {
			if err = conn.Close(); err != nil {
				s.logger.Errorf("getPortIsUsed: %s", err.Error())
				return
			}
		}()
		return true
	}

	ln, err := net.Listen("tcp", net.JoinHostPort("", port))
	if err != nil {
		return true
	}

	if ln != nil {
		defer func() {
			if err = ln.Close(); err != nil {
				s.logger.Errorf("ln close: %s", err.Error())
				return
			}
		}()
	}
	return false
}
