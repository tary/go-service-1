package internal

import (
	"fmt"
	"reflect"

	"github.com/giant-tech/go-service/framework/msgdef"
	"github.com/giant-tech/go-service/framework/net/inet"
	"github.com/giant-tech/go-service/framework/net/internal/internal/msgenc"

	log "github.com/cihub/seelog"
)

// SendIdip 发送IDIP消息，立即返回.
func (sess *Session) SendIdip(msg interface{}) {
	if msg == nil {
		return
	}

	if sess.IsClosed() {
		log.Warnf("Send after sess close %s %s %s",
			sess.conn.RemoteAddr(), reflect.TypeOf(msg),
			fmt.Sprintf("%s", msg)) // log是异步的，所以 msg 必须复制下。
		return
	}

	msgBuf, err := sess.EncodeIdipMsg(msg)
	if err != nil {
		log.Error("Encode message error in Send(): ", err)
		return
	}

	sess.sendBuf.Put(msgBuf)
}

// EncodeIdipMsg 编码idip消息
func (sess *Session) EncodeIdipMsg(msg interface{}) ([]byte, error) {
	msgID, err := msgdef.GetMsgDefIDIP().GetMsgIDByType(msg)
	if err != nil {
		log.Errorf("message '%s' is not registered", reflect.TypeOf(msg))
		return nil, err
	}

	msgBuf, err := msgenc.EncodeIdipMsg(msg, inet.IdipMsgID(msgID))
	if err != nil {
		log.Error("Encode message error in EncodeIdipMsg(): ", err)
		return nil, err
	}

	return msgBuf, nil
}

// RegIdipMsgProcFunc 注册消息处理函数.
func (sess *Session) RegIdipMsgProcFunc(msgID inet.IdipMsgID, procFun func(interface{})) {
	sess.msgHdlr.RegIdipMsgProcFunc(msgID, procFun)
}

// msgBuf 发送心跳
func (sess *Session) sendServerHB() {
	sess.SendIdip(&msgdef.ServerHB{})
}
