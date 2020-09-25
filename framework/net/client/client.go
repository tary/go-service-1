package client

import (
	"fmt"
	"net"

	"github.com/cihub/seelog"
	"github.com/giant-tech/go-service/framework/net/inet"
	kcp "github.com/xtaci/kcp-go"
)

// Dial 创建一个连接
func Dial(protocol string, addr string) (inet.ISession, error) {
	var conn net.Conn
	var err error

	seelog.Debug("Dial protocol: ", protocol, ", addr: ", addr)

	if protocol == "tcp" || protocol == "idip" {
		if conn, err = net.Dial("tcp", addr); err != nil {
			return nil, err
		}
	} else if protocol == "kcp" {
		var kcpConn *kcp.UDPSession
		if kcpConn, err = kcp.DialWithOptions(addr, nil, 3, 2); err != nil {
			return nil, err
		}

		kcpConn.SetStreamMode(false)
		kcpConn.SetNoDelay(1, 10, 2, 1)
		kcpConn.SetDSCP(46)
		kcpConn.SetACKNoDelay(true)

		conn = kcpConn
	} else {
		return nil, fmt.Errorf("unknown network protocol '%s'", protocol)
	}

	isIdip := false
	if protocol == "idip" {
		isIdip = true
	}

	sess := NewSession(conn, isIdip)

	return sess, nil
}
