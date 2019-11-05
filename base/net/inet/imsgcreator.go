package inet

// IMsgCreator 消息创建
type IMsgCreator interface {
	NewMsg(MsgID) IMsg
	NewIdipMsg(IdipMsgID) interface{}

	RegMsgCreator(msgID MsgID, msg IMsg)
	RegIdipMsgCreator(msgID IdipMsgID, IdipMsgCreator func() interface{})

	IsEmpty() bool
}
