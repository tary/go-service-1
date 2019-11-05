package msghdl

import (
	"fmt"
	"reflect"
	"runtime/debug"

	"github.com/GA-TECH-SERVER/zeus/base/net/inet"
	"github.com/GA-TECH-SERVER/zeus/base/serializer"
	"github.com/GA-TECH-SERVER/zeus/framework/msgdef"

	assert "github.com/aurelien-rainone/assertgo"
	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

// IMessgeHandler 消息处理
type IMessgeHandler interface {
	HandleRawMsg(msgID tMsgID, rawMsgBuf []byte)
	RegMsgProcFunc(msgID tMsgID, msgProcFunc tMsgProcFunc)
	RegIdipMsgProcFunc(msgID tIdipMsgID, msgProcFunc tIdipMsgProcFunc)
	isRegistered(msgID tMsgID) bool
	isIdipRegistered(msgID tIdipMsgID) bool

	RegMsgProc(proc interface{})
}

// tMsgID msgID
type tMsgID = inet.MsgID

// tMsgProcFunc 处理器函数，如 (*MsgProc).MsgProcTest(msg inet.IMsg)
type tMsgProcFunc func(inet.IMsg)

// MessageHandler 非线程安全。仅被 session.recvLoop() 协程使用。
type MessageHandler struct {
	mapIDToFunc map[tMsgID]reflect.Value
}

// New 新建
func New() *MessageHandler {
	return &MessageHandler{
		mapIDToFunc: make(map[tMsgID]reflect.Value),
	}
	// 需要由调用者来 RegMsgProcFunc() 注册所有处理函数。
}

// HandleRawMsg 处理原始消息.
// 输入为去头解密解压后的数据。
// 应该允许不同版本的客户端，所以忽略所有版本不同造成的错误。
func (m *MessageHandler) HandleRawMsg(msgID tMsgID, rawMsgBuf []byte) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("handle msg is panic:", msgID, err, string(debug.Stack()))
			if viper.GetString("Config.Recover") == "0" {
				panic(err)
			}
		}
	}()

	var msg inet.IMsg
	var err error
	_, msg, err = msgdef.GetMsgDef().GetMsgInfo(msgID)
	if err != nil {
		log.Debugf("unknown message ID %d", msgID)
		return
	}

	if err := serializer.UnSerializeNew(msg, rawMsgBuf); err != nil {
		log.Debugf("illegal message: %s", err)
		return
	}

	f, ok := m.mapIDToFunc[msgID]
	if !ok {
		log.Debugf("no handler for msg ID %d", msgID)
		return
	}

	f.Call([]reflect.Value{reflect.ValueOf(msg)})
}

// RegMsgProc 注册消息处理对象
// 其中 proc 是一个对象，包含是类似于 MsgProcXXXXX的一系列函数，分别用来处理不同的消息
func (m *MessageHandler) RegMsgProc(proc interface{}) {
	v := reflect.ValueOf(proc)
	t := reflect.TypeOf(proc)

	for i := 0; i < t.NumMethod(); i++ {
		methodName := t.Method(i).Name
		msgName, msgHandler, err := m.getMsgHandler(methodName, v.MethodByName(methodName))
		if err == nil {
			msgid, err1 := msgdef.GetMsgDef().GetMsgIDByName(msgName)
			if err1 != nil {
				continue
			}

			m.RegMsgProcValue(msgid, msgHandler)
			continue
		}
	}
}

// RegMsgProcFunc 注册一个消息处理函数.
// Session 创建时会注册所有的 MsgProc.
func (m *MessageHandler) RegMsgProcFunc(msgID tMsgID, msgProcFunc tMsgProcFunc) {
	assert.True(msgID != 0, "message ID 0 is illegal")
	assert.True(nil != msgProcFunc, "msgProc is nil")
	m.RegMsgProcValue(msgID, reflect.ValueOf(msgProcFunc))
}

// RegMsgProcValue 注册一个消息处理函数.
// Session 创建时会注册所有的 MsgProc.
func (m *MessageHandler) RegMsgProcValue(msgID tMsgID, msgProcValue reflect.Value) {
	assert.True(msgID != 0, "message ID 0 is illegal")
	assert.False(m.isRegistered(msgID),
		fmt.Sprintf("message ID is already registered: %d", msgID))

	m.mapIDToFunc[msgID] = msgProcValue
}

// RegIdipMsgProcFunc RegIdipMsgProcFunc
func (m *MessageHandler) RegIdipMsgProcFunc(msgID tIdipMsgID, msgProcFunc tIdipMsgProcFunc) {
	panic("can't call MessageHandler.RegIdipMsgProcFunc")
}

// isRegistered isRegistered
func (m *MessageHandler) isRegistered(msgID tMsgID) bool {
	_, ok := m.mapIDToFunc[msgID]
	return ok
}

// isIdipRegistered  isIdipRegistered
func (m *MessageHandler) isIdipRegistered(msgID tIdipMsgID) bool {
	panic("can't call MessageHandler.isIdipRegistered")
}

// getMsgHandler getMsgHandler
func (m *MessageHandler) getMsgHandler(methodName string, v reflect.Value) (string, reflect.Value, error) {
	methodHead := "MsgProc"
	methodHeadLen := len(methodHead)

	if len(methodName) < methodHeadLen+1 {
		return "", reflect.ValueOf(nil), fmt.Errorf("")
	}

	if methodName[0:methodHeadLen] != methodHead {
		return "", reflect.ValueOf(nil), fmt.Errorf("")
	}

	msgName := methodName[methodHeadLen:]

	//此处应该检查该函数是否是MsgHanderFunc类型的参数
	return msgName, v, nil
}
