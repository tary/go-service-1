package msgid

import (
	"fmt"
	"strconv"

	"github.com/GA-TECH-SERVER/zeus/base/net/inet"
)

// internalMSgID 网络内部保留消息id
const internalMSgID inet.MsgID = 10

// GetMsgID  BKDR Hash Function
func GetMsgID(msgName string) inet.MsgID {
	var seed uint32 = 131 // 31 131 1313 13131 131313 etc..
	var hash uint32

	for _, v := range msgName {
		hash = hash*seed + uint32(v)
	}

	return inet.MsgID(hash)
}

// GenMsgMap GenMsgMap
func GenMsgMap(msgVec []string) map[string]string {
	msgMap := make(map[string]string)
	for _, v := range msgVec {
		id := GetMsgID(v)
		if id <= internalMSgID {
			panic(fmt.Errorf("协议名与内部保留id冲突 %s, id: %d", v, id))
		}

		msgID := strconv.Itoa(int(id))
		value, ok := msgMap[msgID]
		if ok && value != v {
			panic(fmt.Errorf("协议名冲突 %s : %s", v, value))
		}

		msgMap[msgID] = v
	}

	return msgMap
}
