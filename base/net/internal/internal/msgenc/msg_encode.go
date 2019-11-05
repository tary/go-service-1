package msgenc

import (
	"encoding/binary"

	"github.com/GA-TECH-SERVER/zeus/base/net/inet"
	"github.com/GA-TECH-SERVER/zeus/base/net/internal/internal/consts"
	"github.com/GA-TECH-SERVER/zeus/base/net/internal/internal/crypt"
	"github.com/GA-TECH-SERVER/zeus/base/serializer"

	assert "github.com/aurelien-rainone/assertgo"
	"github.com/golang/snappy"
)

// minCompressSize 最小压缩大小
const minCompressSize = 100

// MsgHeadSize 消息头大小
const MsgHeadSize = consts.MsgHeadSize

// MsgIDSize 消息ID大小
const MsgIDSize = consts.MsgIDSize

// compressFlag 压缩标记
const compressFlag byte = 0x01

// encryptFlag 加密标记
const encryptFlag byte = 0x02

// EncodeMsg 序列化消息.
// 返回的slice带头部长度和消息ID, 但是长度待设置.
func EncodeMsg(msg inet.IMsg, msgID inet.MsgID) ([]byte, error) {
	assert.True(msgID != 0)
	buf := make([]byte, MsgHeadSize+MsgIDSize+serializer.GetSizeNew(msg))
	binary.LittleEndian.PutUint16(buf[4:], uint16(msgID))

	data := serializer.SerializeNewWithBuff(buf[MsgHeadSize+MsgIDSize:], msg)
	n := len(data)

	assert.True(n <= serializer.GetSizeNew(msg))
	return buf[:n+MsgHeadSize+MsgIDSize], nil
}

// CompressAndEncrypt 压缩和加密已序列化消息.
// 输入消息缓冲区无压缩和加密，带头部长度和消息ID。
func CompressAndEncrypt(buf []byte, forceNoCompress bool, encryptEnabled bool, isClient bool) ([]byte, error) {
	// 压缩后会返回新的buf, 如果没压缩就返回原buf
	msgBuf, err := compress(buf, forceNoCompress)
	if err != nil {
		return nil, err
	}

	if encryptEnabled && msgBuf[3]&encryptFlag == 0 {
		data := msgBuf[MsgHeadSize:]
		crypt.EncryptData(data, isClient)
		msgBuf[3] = msgBuf[3] | encryptFlag
	}

	// 设置头部长度
	setMsgBufLen(msgBuf)
	return msgBuf, nil
}

// compress 压缩
func compress(buf []byte, forceNoCompress bool) ([]byte, error) {
	//如果已经压缩，直接返回
	if buf[3]&compressFlag > 0 {
		return buf, nil
	}

	msgSize := len(buf) - MsgHeadSize
	if forceNoCompress || msgSize < minCompressSize {
		return buf, nil
	}

	maxLen := snappy.MaxEncodedLen(msgSize)
	p := make([]byte, MsgHeadSize+maxLen)
	mbuff := snappy.Encode(p[MsgHeadSize:], buf[MsgHeadSize:])

	// p[0..2]长度暂不设置，仅设置p[3]标志位
	p[3] = buf[3] | compressFlag // 压缩标志

	return p[:len(mbuff)+MsgHeadSize], nil
}

// setMsgBufLen 设置消息buflen
func setMsgBufLen(msgBuf []byte) {
	//设置过长度，直接返回
	if msgBuf[0] > 0 || msgBuf[1] > 0 || msgBuf[2] > 0 {
		return
	}

	bufLen := len(msgBuf)
	cmdSize := bufLen - 4 // 去除长度和标志共4字节
	msgBuf[0] = byte(cmdSize)
	msgBuf[1] = byte(cmdSize >> 8)
	msgBuf[2] = byte(cmdSize >> 16)
}
