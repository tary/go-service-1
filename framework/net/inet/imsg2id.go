package inet

import "github.com/giant-tech/go-service/base/imsg"

// IMsg2ID  取ID接口.
type IMsg2ID interface {
	GetMsgID(msg imsg.IMsg) MsgID
	RegMsg2ID(msg imsg.IMsg, msgID MsgID)
}
