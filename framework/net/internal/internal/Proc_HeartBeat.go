package internal

import (
	"github.com/giant-tech/go-service/base/imsg"
	"github.com/giant-tech/go-service/framework/msgdef"
)

// NewHeartBeatProc new心跳处理
func NewHeartBeatProc(sess *Session) *ProcHeartBeat {
	return &ProcHeartBeat{
		sess: sess,
	}
}

// ProcHeartBeat 心跳相关消息处理
type ProcHeartBeat struct {
	sess *Session // 一般都需要包含session对象
}

// MsgProcPing 心跳请求
func (p *ProcHeartBeat) MsgProcPing(msg imsg.IMsg) {

	retMsg := &msgdef.Pong{}
	p.sess.Send(retMsg)
}

// MsgProcPong 心跳响应
func (p *ProcHeartBeat) MsgProcPong(msg imsg.IMsg) {

}
