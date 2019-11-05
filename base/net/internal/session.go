package internal

import (
	"net"

	"github.com/GA-TECH-SERVER/zeus/base/net/inet"
	"github.com/GA-TECH-SERVER/zeus/base/net/internal/internal"
)

// NewSession 新建session
func NewSession(conn net.Conn, encryEnabled bool, isClient bool, isIdip bool) inet.ISession {
	return internal.NewSession(conn, encryEnabled, isClient, isIdip)
}
