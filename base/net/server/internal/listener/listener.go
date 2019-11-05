package listener

import (
	"fmt"

	"github.com/giant-tech/go-service/base/net/server/internal/listener/internal"
)

// IListener 网络监听
type IListener interface {
	Run(internal.IConnHandler)
	Close()
	GetPort() string
}

// NewListener 创建并开始监听.
// protocal 支持："kcp", "tcp", "tcp+kcp".
// addr 形如：":80", "1.2.3.4:80"
// maxConns 是最大连接数
func NewListener(protocal string, addr string, maxConns int) (IListener, error) {
	if protocal == "tcp" || protocal == "kcp" || protocal == "tcp+kcp" {
		return internal.NewARQListener(protocal, addr, maxConns)
		// } else if protocal == "udp" {
		// 	return newUDPNetSrv()
	}
	return nil, fmt.Errorf("illegal protocal: '%s'", protocal)
}
