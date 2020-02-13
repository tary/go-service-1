package internal

import (
	"context"
	"net"
	"reflect"
	"runtime/debug"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
	kcp "github.com/xtaci/kcp-go"
	"golang.org/x/net/netutil"
)

// ARQListener 网络服务器监听
type ARQListener struct {
	protocol    string
	listener    net.Listener
	kcpListener net.Listener
	ctx         context.Context
	ctxCancel   context.CancelFunc
	listenPort  string

	maxConns int
}

// NewARQListener 新建arq监听
func NewARQListener(protocol string, addr string, maxConns int) (*ARQListener, error) {
	l := &ARQListener{
		protocol:    protocol,
		listener:    nil,
		kcpListener: nil,
		maxConns:    maxConns,
	}
	l.ctx, l.ctxCancel = context.WithCancel(context.Background())

	err := l.listen(addr)
	return l, err
}

// listen 监听
func (a *ARQListener) listen(addr string) error {
	var err error
	var kcpListener *kcp.Listener

	switch a.protocol {
	case "tcp":
		a.listener, err = net.Listen("tcp", addr)
	case "kcp":
		kcpListener, err = kcp.ListenWithOptions(addr, nil, 3, 2)
		kcpListener.SetDSCP(46)
		a.listener = kcpListener
	case "tcp+kcp":
		a.listener, err = net.Listen("tcp", addr)
		if err != nil {
			return err
		}

		newAddr, _ := getNewAddr(addr, a.listener.Addr().String())
		if err != nil {
			return err
		}

		kcpListener, err = kcp.ListenWithOptions(newAddr, nil, 3, 2)
		kcpListener.SetDSCP(46)
		a.kcpListener = kcpListener
	default:
		panic("WRONG PROTOCOL")
	}
	if err != nil {
		return err
	}

	_, a.listenPort, _ = net.SplitHostPort(a.listener.Addr().String())

	if a.maxConns > 0 {
		a.listener = netutil.LimitListener(a.listener, a.maxConns)
		//todo ？ maxConns优化
		if a.kcpListener != nil {
			a.kcpListener = netutil.LimitListener(a.kcpListener, a.maxConns)
		}
	}

	return nil
}

// Close 关闭
func (a *ARQListener) Close() {
	a.ctxCancel()
	a.listener.Close()
	if a.kcpListener != nil {
		a.kcpListener.Close()
	}
}

// Run 运行
func (a *ARQListener) Run(connHandler IConnHandler) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("ARQListener.Run panic:", err, string(debug.Stack()))
			if viper.GetString("Config.Recover") == "0" {
				panic(err)
			}
		}
	}()

	// 同时开启kcp
	if a.kcpListener != nil {
		go func() {
			for {
				select {
				case <-a.ctx.Done():
					return
				default:
					{
						conn, err := a.kcpListener.Accept()
						if err != nil {
							log.Error("accept connection error ", err)
							continue
						}

						var kcpConn *kcp.UDPSession
						if a.maxConns > 0 {
							// 获取limitListenerConn的net.Conn
							// type limitListenerConn struct {
							// 	net.Conn
							// 	releaseOnce sync.Once
							// 	release     func()
							// }
							v := reflect.ValueOf(conn)
							if v.Kind() == reflect.Ptr {
								v = v.Elem()
							}

							f := v.Field(0)
							kcpConn = f.Interface().(*kcp.UDPSession)
						} else {
							kcpConn = conn.(*kcp.UDPSession)
						}

						kcpConn.SetStreamMode(false)
						kcpConn.SetNoDelay(1, 10, 2, 1)
						kcpConn.SetDSCP(46)
						kcpConn.SetACKNoDelay(true)

						go connHandler.HandleConn(conn)
					}
				}
			}
		}()
	}

	for {
		select {
		case <-a.ctx.Done():
			return
		default:
			{
				conn, err := a.listener.Accept()
				if err != nil {
					log.Error("accept connection error ", err)
					continue
				}
				go connHandler.HandleConn(conn)
			}
		}
	}
}

//GetPort 获取监听端口
func (a *ARQListener) GetPort() string {
	return a.listenPort
}

//getNewAddr 获取监听地址， addr1为配置的监听地址， addr2为获取的地址，返回新的监听地址
func getNewAddr(addr1, addr2 string) (string, error) {
	addr, port, err := net.SplitHostPort(addr1)
	if err != nil {
		return "", err
	}

	_, port2, err2 := net.SplitHostPort(addr2)
	if err2 != nil {
		return "", err2
	}

	if port != port2 {
		return addr + ":" + port2, nil
	}
	return addr1, nil

}
