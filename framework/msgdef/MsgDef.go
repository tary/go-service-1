package msgdef

import (
	"fmt"
	"reflect"
	"runtime/debug"

	"github.com/giant-tech/go-service/base/net/inet"
	"github.com/giant-tech/go-service/base/net/msgid"

	log "github.com/cihub/seelog"
)

// MsgInfo 消息信息
type MsgInfo struct {
	id      inet.MsgID
	name    string
	msgType reflect.Type
}

// MsgDef 消息定义管理器，包含是消息与消息号的映射结构，客户端初始化的时候，由服务器下发下去
type MsgDef struct {
	id2Info     map[inet.MsgID]*MsgInfo
	name2Info   map[string]*MsgInfo
	msgTypeToID map[reflect.Type]*MsgInfo
}

// Init 创建MsgDef并初始化
func Init() {
	if msgDefInst == nil {
		msgDefInst = &MsgDef{
			make(map[inet.MsgID]*MsgInfo),
			make(map[string]*MsgInfo),
			make(map[reflect.Type]*MsgInfo),
		}
		msgDefInst.Init()
	}
}

// Init 初始消息定义管理器
func (def *MsgDef) Init() {
	def.InitBytesMsg()
}

// GetMsgInfo 根据消息号，获得消息的类型，名称，如果是protobuf消息，获得proto消息的容器
func (def *MsgDef) GetMsgInfo(msgID inet.MsgID) (msgName string, msgContent inet.IMsg, err error) {

	//此处会被多线程调用，不确定会不会有问题
	info, ok := def.id2Info[msgID]

	if !ok {
		return "", nil, fmt.Errorf("不存在消息号 ID: %d", msgID)
	}

	return info.name, reflect.New(info.msgType.Elem()).Interface().(inet.IMsg), nil
}

// IsMsgExist 消息是否存在
func (def *MsgDef) IsMsgExist(msgID inet.MsgID) bool {
	_, ok := def.id2Info[msgID]
	return ok
}

// GetMsgIDByName 通过名字获取ID号
func (def *MsgDef) GetMsgIDByName(msgName string) (inet.MsgID, error) {

	info, ok := def.name2Info[msgName]
	if !ok {
		log.Errorf("未注册的消息名 : %s", msgName)
		return 0, fmt.Errorf("未注册的消息名 : %s", msgName)
	}

	return info.id, nil
}

// GetMsgIDByType 从消息类型获取ID.
func (def *MsgDef) GetMsgIDByType(msg inet.IMsg) (inet.MsgID, error) {
	info, ok := def.msgTypeToID[reflect.TypeOf(msg)]
	if !ok {
		log.Debug(string(debug.Stack()))
		return 0, fmt.Errorf("不存在消息 Msg: %v", reflect.TypeOf(msg))
	}

	return info.id, nil
}

// RegMsg 注册消息
func (def *MsgDef) RegMsg(msgID inet.MsgID, msg inet.IMsg) {
	if msgID == 0 {
		msgID = msgid.GetMsgID(reflect.TypeOf(msg).Elem().Name())
	}

	_, ok := def.id2Info[msgID]

	if ok {
		log.Warn("消息ID已经存在 ", msgID)
		return
	}

	msgName := reflect.TypeOf(msg).Elem().Name()

	_, ok = def.name2Info[msgName]

	if ok {
		log.Warn("消息名称已经存在  ", msgName)
		return
	}

	info := &MsgInfo{
		msgID,
		msgName,
		reflect.TypeOf(msg),
	}

	def.id2Info[msgID] = info
	def.name2Info[msgName] = info
	def.msgTypeToID[reflect.TypeOf(msg)] = info
}

var msgDefInst *MsgDef

// GetMsgDef 获取消息定义对象的全局实例
func GetMsgDef() *MsgDef {

	if msgDefInst == nil {
		msgDefInst = &MsgDef{
			make(map[inet.MsgID]*MsgInfo),
			make(map[string]*MsgInfo),
			make(map[reflect.Type]*MsgInfo),
		}
		msgDefInst.Init()
	}

	return msgDefInst
}
