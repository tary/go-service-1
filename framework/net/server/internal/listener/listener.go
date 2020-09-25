package listener

import (
	"fmt"

	"github.com/giant-tech/go-service/framework/net/server/internal/listener/internal"
)

// IListener 网络监听
type IListener interface {
	Run(internal.IConnHandler)
	Close()
	GetPort() string
}

// NewListener 创建并开始监听.
// protocol 支持："kcp", "tcp", "tcp+kcp".
// addr 形如：":80", "1.2.3.4:80"
// maxConns 是最大连接数
func NewListener(protocol string, addr string, maxConns int) (IListener, error) {
	if protocol == "tcp" || protocol == "kcp" || protocol == "tcp+kcp" {
		return internal.NewARQListener(protocol, addr, maxConns)
		// } else if protocol == "udp" {
		// 	return newUDPNetSrv()
	}
	return nil, fmt.Errorf("illegal protocol: '%s'", protocol)
}
