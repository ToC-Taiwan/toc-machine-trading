// Package utils package utils
package utils

import (
	"net"
	"time"
)

func GetPortIsUsed(host string, port string) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), 500*time.Millisecond)
	if err != nil && conn != nil {
		return false
	}

	if conn != nil {
		defer func() {
			if err = conn.Close(); err != nil {
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
				return
			}
		}()
	}
	return false
}
