package session

import (
	"net"

	"github.com/giant-tech/go-service/base/net/inet"
	"github.com/giant-tech/go-service/base/net/internal"

	"go.uber.org/atomic"
	"golang.org/x/time/rate"
)

// Session 封装客户端服务器通用会话，并提供服务器专用功能。
type Session struct {
	_IInternalSession

	// 可以保存任意用户数据
	userData atomic.Value

	// 会话事件接收器
	evtSink inet.ISessEvtSink
}

// _IInternalSession 代表一个客户端服务器通用会话.
type _IInternalSession interface {
	RegMsgProcFunc(msgID inet.MsgID, procFunc func(inet.IMsg))
	RegMsgProc(interface{})

	Send(inet.IMsg) error
	SendRaw([]byte) error
	EncodeMsg(inet.IMsg) ([]byte, error)
	CompressAndEncrypt(buf []byte) ([]byte, error)

	Start()
	Close()
	IsClosed() bool

	RemoteAddr() string

	ResetHb()

	SetVerifyMsgID(inet.MsgID)
	SetVerified()

	SetEncrypt(isEncrypt bool)

	SetOnClosed(func())
	AddOnClosed(func()) // 注册断开时动作，可注册多个依次调用

	SetBytePerSecLimiter(r rate.Limit, b int)
	SetQueryPerSecLimiter(r rate.Limit, b int)

	//只有idip才能调用，同征途消息结构匹配
	RegIdipMsgProcFunc(msgID inet.IdipMsgID, procFunc func(interface{}))
	SendIdip(msg interface{})

	GetID() uint64

	On(evt string, f func(interface{}))
	Emit(evt string, p interface{})
}

// New 新建
func New(conn net.Conn, encryEnabled bool, sessEvtSink inet.ISessEvtSink, isIdip bool) *Session {
	result := &Session{
		_IInternalSession: internal.NewSession(conn, encryEnabled, false, isIdip),
		evtSink:           sessEvtSink,
	}
	result.SetOnClosed(result.onClosed)

	if sessEvtSink != nil {
		sessEvtSink.OnConnected(result)
	}
	return result
}

// onClosed session关闭回调
func (s *Session) onClosed() {
	if s.evtSink == nil {
		return
	}
	s.evtSink.OnClosed(s)
}

// GetUserData 获得用户数据
func (s *Session) GetUserData() interface{} {
	return s.userData.Load()
}

// SetUserData 设置用户数据
func (s *Session) SetUserData(data interface{}) {
	s.userData.Store(data)
}
