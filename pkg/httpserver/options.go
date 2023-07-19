package httpserver

import (
	"net"
	"time"
)

// Option -.
type Option func(*Server)

// Port -.
func Port(port string) Option {
	return func(s *Server) {
		s.srv.Addr = net.JoinHostPort("", port)
	}
}

// ReadTimeout -.
func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.srv.ReadTimeout = timeout
	}
}

// WriteTimeout -.
func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.srv.WriteTimeout = timeout
	}
}

func KeyPath(keyPath string) Option {
	return func(s *Server) {
		s.keyPath = keyPath
	}
}

func CertPath(certPath string) Option {
	return func(s *Server) {
		s.certPath = certPath
	}
}

func AddLogger(logger Logger) Option {
	return func(c *Server) {
		c.logger = logger
	}
}
