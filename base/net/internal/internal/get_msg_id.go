package internal

import (
	"github.com/GA-TECH-SERVER/zeus/base/net/inet"
	"github.com/GA-TECH-SERVER/zeus/base/net/internal/internal/consts"
)

// MsgID MSGid
type MsgID = inet.MsgID

// IdipMsgID IdipMsgID
type IdipMsgID = inet.IdipMsgID

// getMsgID 获取消息ID
func getMsgID(buf []byte) MsgID {
	if len(buf) < consts.MsgIDSize {
		return 0
	}

	return MsgID(buf[0]) | MsgID(buf[1])<<8
}

// getIdipMsgID 获取消息ID
func getIdipMsgID(buf []byte) IdipMsgID {
	if len(buf) < consts.MsgIDSize {
		return 0
	}

	//fmt.Println("getIdipMsgID:", IdipMsgID(buf[0])|IdipMsgID(buf[1])<<8)

	return IdipMsgID(buf[0]) | IdipMsgID(buf[1])<<8
}
