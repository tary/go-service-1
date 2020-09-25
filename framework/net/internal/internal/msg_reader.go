package internal

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/giant-tech/go-service/framework/net/inet"
	"github.com/giant-tech/go-service/framework/net/internal/internal/consts"
	"github.com/giant-tech/go-service/framework/net/internal/internal/crypt"

	assert "github.com/aurelien-rainone/assertgo"
	"github.com/golang/snappy"
)

const (
	msgIDSize    = consts.MsgIDSize
	maxUDPPacket = 62 * 1024
)

// ReadARQMsg 读取一条消息.
// 返回 msgID, rawMsgBuf, err
// rawMsgBuf中是不含消息ID和消息长度的消息体，已解密解压。
// 无数据则返回 0, nil, nil
// 0 为无效ID
func readARQMsg(conn net.Conn, isClient bool) (msgID inet.MsgID, msgBuf []byte, err error) {
	if conn == nil {
		return 0, nil, errors.New("无效连接")
	}

	msgHead := make([]byte, consts.MsgHeadSize)
	if _, err := io.ReadFull(conn, msgHead); err != nil {
		return 0, nil, err
	}

	msgSize := (int(msgHead[0]) | int(msgHead[1])<<8 | int(msgHead[2])<<16)
	if msgSize > maxUDPPacket || msgSize < msgIDSize {
		return 0, nil, fmt.Errorf("收到的数据长度非法:%d", msgSize)
	}

	msgData := make([]byte, msgSize)
	if _, err := io.ReadFull(conn, msgData); err != nil {
		return 0, nil, err
	}
	assert.True(len(msgData) >= msgIDSize, "data len is too short")

	flag := msgHead[3]
	encryptFlag := flag & 0x2
	if encryptFlag > 0 {
		msgData = crypt.DecryptData(msgData, isClient)
	}

	compressFlag := flag & 0x1
	if compressFlag == 0 {
		msgID = getMsgID(msgData)
		return msgID, msgData[msgIDSize:], nil
	}

	// 解压
	buf := make([]byte, consts.MaxMsgBuffer)
	msgData, err = snappy.Decode(buf, msgData)
	if err != nil {
		return 0, nil, err
	}

	msgID = getMsgID(msgData)
	return msgID, msgData[msgIDSize:], err
}
