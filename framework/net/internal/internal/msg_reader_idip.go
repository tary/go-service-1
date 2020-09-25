package internal

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/giant-tech/go-service/framework/net/inet"
	"github.com/giant-tech/go-service/framework/net/internal/internal/consts"

	assert "github.com/aurelien-rainone/assertgo"
)

const (
	idipMsgIDSize = consts.IdipMsgIDSize
)

// DoZlibCompress 进行zlib压缩
func DoZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

// DoZlibUnCompress 进行zlib解压缩
func DoZlibUnCompress(compressSrc []byte) ([]byte, error) {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(&out, r)
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

// readARQIdipMsg 读取一条消息.
// 返回 msgID, rawMsgBuf, err
// rawMsgBuf中是不含消息ID和消息长度的消息体，已解密解压。
// 无数据则返回 0, nil, nil
// 0 为无效ID
func readARQIdipMsg(conn net.Conn) (msgID inet.MsgID, msgBuf []byte, err error) {
	if conn == nil {
		return 0, nil, errors.New("无效连接")
	}

	msgHead := make([]byte, consts.IdipMsgHeadSize-idipMsgIDSize)
	if _, err := io.ReadFull(conn, msgHead); err != nil {
		return 0, nil, err
	}

	msgSize := (int(msgHead[0]) | int(msgHead[1])<<8)
	if msgSize > maxUDPPacket || msgSize < idipMsgIDSize {
		return 0, nil, fmt.Errorf("收到的数据长度非法:%d", msgSize)
	}

	msgData := make([]byte, msgSize)
	if _, err := io.ReadFull(conn, msgData); err != nil {
		return 0, nil, err
	}

	assert.True(len(msgData) >= idipMsgIDSize, "data len is too short")
	// var msgBody []byte
	// if len(msgData) >= msgIDSize {
	// 	msgBody = msgData[msgIDSize:]
	// }

	flag := msgHead[3]
	// encryptFlag := flag & 0x2
	// if encryptFlag > 0 {
	// 	msgBody = crypt.DecryptData(msgBody)
	// }

	//zlab压缩
	compressFlag := flag & 0x40
	if compressFlag == 0 {
		msgBody := msgData[idipMsgIDSize:]
		msgID = inet.MsgID(getIdipMsgID(msgData))
		return msgID, msgBody, nil
	}

	// 解压
	//buf := make([]byte, consts.MaxMsgBuffer)
	msgBody, err := DoZlibUnCompress(msgData)
	if err != nil {
		return 0, nil, err
	}

	msgBuf = msgBody[idipMsgIDSize:]
	msgID = inet.MsgID(getIdipMsgID(msgBody))
	return msgID, msgBuf, err
}
