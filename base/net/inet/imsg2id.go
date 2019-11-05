package inet

// IMsg2ID  取ID接口.
type IMsg2ID interface {
	GetMsgID(msg IMsg) MsgID
	RegMsg2ID(msg IMsg, msgID MsgID)
}
