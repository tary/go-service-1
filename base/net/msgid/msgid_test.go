package msgid

import (
	"fmt"
	"testing"
)

func TestMsgID(t *testing.T) {
	msgID := GetMsgID("ClientVerifyReq")
	fmt.Println(msgID)

	msgID = GetMsgID("ClientVerifyResp")
	fmt.Println(msgID)

	msgID = GetMsgID("LoginReq")
	fmt.Println(msgID)

}
