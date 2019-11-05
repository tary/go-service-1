package msgenc

import (
	"encoding/binary"

	"github.com/GA-TECH-SERVER/zeus/base/net/inet"
	"github.com/GA-TECH-SERVER/zeus/base/net/internal/internal/consts"
	"github.com/GA-TECH-SERVER/zeus/base/serializer"

	"github.com/golang/snappy"
)

// IdipMsgHeadSize idip msg头长度
const IdipMsgHeadSize = consts.IdipMsgHeadSize

// EncodeIdipMsg 序列化消息.
// 返回的slice带头部长度和消息ID, 但是长度待设置.
func EncodeIdipMsg(msg interface{}, msgID inet.IdipMsgID) ([]byte, error) {
	data := serializer.SerializeNew(msg)
	// if err != nil {
	// 	return data, err
	// }

	buff := make([]byte, IdipMsgHeadSize+len(data))
	binary.LittleEndian.PutUint16(buff[4:], uint16(msgID))
	copy(buff[IdipMsgHeadSize:], data[:])

	return buff, nil
}

// CompressAndEncryptIdip 压缩和加密已序列化消息.
// 输入消息缓冲区无压缩和加密，带头部长度和消息ID。
func CompressAndEncryptIdip(buf []byte, forceNoCompress bool, encryptEnabled bool) ([]byte, error) {
	// 设置头部长度
	setIdipMsgBufLen(buf)
	return buf, nil
}

func compressIdip(buf []byte, forceNoCompress bool) ([]byte, error) {
	msgSize := len(buf) - IdipMsgHeadSize
	if forceNoCompress || msgSize < minCompressSize {
		return buf, nil
	}

	maxLen := snappy.MaxEncodedLen(msgSize) // 不压缩2字节的ID
	p := make([]byte, IdipMsgHeadSize+maxLen)
	mbuff := snappy.Encode(p[IdipMsgHeadSize:], buf[IdipMsgHeadSize:])

	// p[0..2]长度暂不设置，仅设置p[3]标志位
	p[3] = buf[3] | 0x1 // 压缩标志

	// MsgID 2 字节
	p[4] = buf[4]
	p[5] = buf[5]

	return p[:len(mbuff)+IdipMsgHeadSize], nil
}

func setIdipMsgBufLen(msgBuf []byte) {
	bufLen := len(msgBuf)
	cmdSize := bufLen - 4 // 去除长度和标志共4字节

	binary.LittleEndian.PutUint16(msgBuf[0:], uint16(cmdSize))
}
