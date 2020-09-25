package msgprocset

import "github.com/giant-tech/go-service/framework/net/inet"

// IMsgProc 消息处理接口
type IMsgProc interface {
	RegisterMsgProcFunctions(sess inet.ISession) interface{}
}

// MsgProcSet 消息处理器集合.
// 非线程安全,要求在服务器Run()之前添加所有MsgProc。
type MsgProcSet struct {
	// MsgProc集合：IMsgProc -> true
	msgProcSet []IMsgProc
}

// New 新建
func New() *MsgProcSet {
	return &MsgProcSet{}
}

// AddMsgProc 添加消息处理
func (m *MsgProcSet) AddMsgProc(msgProc IMsgProc) {
	m.msgProcSet = append(m.msgProcSet, msgProc)
}

// RegisterToSession 注册到sess
func (m *MsgProcSet) RegisterToSession(sess inet.ISession) {
	for _, msgProc := range m.msgProcSet {
		msgProc.RegisterMsgProcFunctions(sess)
	}
}

// IsEmpty 是否空
func (m *MsgProcSet) IsEmpty() bool {
	return len(m.msgProcSet) == 0
}
