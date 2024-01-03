package httpserver

import (
	"net"
)

// Option -.
type Option func(*Server)

// Port -.
func Port(port string) Option {
	return func(s *Server) {
		s.srv.Addr = net.JoinHostPort("", port)
	}
}

func AddLogger(logger Logger) Option {
	return func(c *Server) {
		c.logger = logger
	}
}
