package internal

import "net"

// IConnHandler 连接处理接口
type IConnHandler interface {
	HandleConn(net.Conn)
}
