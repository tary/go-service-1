package msghdl

import (
	"fmt"

	"github.com/GA-TECH-SERVER/zeus/base/net/inet"
	"github.com/GA-TECH-SERVER/zeus/base/serializer"
	"github.com/GA-TECH-SERVER/zeus/framework/msgdef"

	assert "github.com/aurelien-rainone/assertgo"
	log "github.com/cihub/seelog"
)

// tIdipMsgID msgid
type tIdipMsgID = inet.IdipMsgID

// tIdipMsgProcFunc 处理器函数，如 (*MsgProc).MsgProcTest(msg interface{})
type tIdipMsgProcFunc func(interface{})

// IdipMessageHandler 非线程安全。仅被 session.recvLoop() 协程使用。
type IdipMessageHandler struct {
	mapIDToFunc map[tIdipMsgID]tIdipMsgProcFunc
}

// NewIdipHandler 新建idip处理
func NewIdipHandler() *IdipMessageHandler {
	return &IdipMessageHandler{
		mapIDToFunc: make(map[tIdipMsgID]tIdipMsgProcFunc),
	}
	// 需要由调用者来 RegMsgProcFunc() 注册所有处理函数。
}

// HandleRawMsg 处理原始消息.
// 输入为去头解密解压后的数据。
// 应该允许不同版本的客户端，所以忽略所有版本不同造成的错误。
func (m *IdipMessageHandler) HandleRawMsg(msgID tMsgID, rawMsgBuf []byte) {
	idipMsgID := tIdipMsgID(msgID)
	_, msg, err := msgdef.GetMsgDefIDIP().GetMsgInfo(inet.MsgID(idipMsgID))
	if err != nil {
		log.Debugf("unknown IDIP message ID %d", idipMsgID)
		return
	}

	if err := serializer.UnSerializeNew(msg, rawMsgBuf); err != nil {
		log.Debugf("illegal message: %s", err)
		return
	}

	f, ok := m.mapIDToFunc[idipMsgID]
	if !ok {
		log.Debugf("no handler for msg ID %d", idipMsgID)
		return
	}
	f(msg)
}

// RegMsgProcFunc 注册消息处理函数
func (m *IdipMessageHandler) RegMsgProcFunc(msgID tMsgID, msgProcFunc tMsgProcFunc) {
	panic("can't call IdipMessageHandler.RegMsgProcFunc")
}

// RegMsgProc 注册消息处理
func (m *IdipMessageHandler) RegMsgProc(proc interface{}) {
	panic("can't call IdipMessageHandler.RegMsgProcFunc")
}

// RegIdipMsgProcFunc 注册一个消息处理函数.
// Session 创建时会注册所有的 MsgProc.
func (m *IdipMessageHandler) RegIdipMsgProcFunc(msgID tIdipMsgID, msgProcFunc tIdipMsgProcFunc) {
	log.Debugf("RegIdipMsgProcFunc IDIP message ID %d", msgID)

	assert.True(msgID != 0, "message ID 0 is illegal")
	assert.True(nil != msgProcFunc, "msgProc is nil")
	assert.False(m.isIdipRegistered(msgID),
		fmt.Sprintf("message ID is already registered: %d", msgID))

	m.mapIDToFunc[msgID] = msgProcFunc
}

// isRegistered 是否注册
func (m *IdipMessageHandler) isRegistered(msgID tMsgID) bool {
	panic("can't call IdipMessageHandler.isRegistered")
}

// isIdipRegistered 是否idip注册
func (m *IdipMessageHandler) isIdipRegistered(msgID tIdipMsgID) bool {
	idipMsgID := tIdipMsgID(msgID)
	_, ok := m.mapIDToFunc[idipMsgID]
	return ok
}
