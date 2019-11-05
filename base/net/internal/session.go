package internal

import (
	"net"

	"github.com/giant-tech/go-service/base/net/inet"
	"github.com/giant-tech/go-service/base/net/internal/internal"
)

// NewSession 新建session
func NewSession(conn net.Conn, encryEnabled bool, isClient bool, isIdip bool) inet.ISession {
	return internal.NewSession(conn, encryEnabled, isClient, isIdip)
}
