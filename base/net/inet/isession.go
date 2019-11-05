package inet

import (
	"golang.org/x/time/rate"
)

// ISession session接口
type ISession interface {
	RegMsgProcFunc(msgID MsgID, procFunc func(IMsg))
	RegMsgProc(interface{})

	Send(IMsg)
	SendRaw([]byte)
	EncodeMsg(IMsg) ([]byte, error)
	CompressAndEncrypt([]byte) ([]byte, error)

	Close()
	IsClosed() bool
	AddOnClosed(func()) // 注册断开时动作，可注册多个依次调用

	RemoteAddr() string

	ResetHb()
	SetOnClosed(func())

	SetVerifyMsgID(MsgID)
	SetVerified()

	SetEncrypt(isEncrypt bool)

	GetID() uint64

	//GetUserData() interface{}
	//SetUserData(data interface{})

	SetBytePerSecLimiter(r rate.Limit, b int)
	SetQueryPerSecLimiter(r rate.Limit, b int)

	On(evt string, f func(interface{}))
	Emit(evt string, p interface{})

	//只有idip才能调用，同征途消息结构匹配
	RegIdipMsgProcFunc(msgID IdipMsgID, procFunc func(interface{}))
	SendIdip(msg interface{})

	Start()
}
