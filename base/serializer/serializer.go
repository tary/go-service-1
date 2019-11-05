package serializer

import (
	"fmt"
	"reflect"

	"github.com/GA-TECH-SERVER/zeus/base/stream"

	log "github.com/cihub/seelog"
	"github.com/gogo/protobuf/proto"
)

const (
	_           = iota
	typeUint8   = 1
	typeUint16  = 2
	typeUint32  = 3
	typeUint64  = 4
	typeInt8    = 5
	typeInt16   = 6
	typeInt32   = 7
	typeInt64   = 8
	typeFloat32 = 9
	typeFloat64 = 10
	typeString  = 11
	typeBytes   = 12
	typeBool    = 13
	typeProto   = 14
)

// Serialize 序列化
func Serialize(args ...interface{}) []byte {
	size := getSize(args...)
	data := make([]byte, size)
	bw := stream.NewByteStream(data)

	for _, arg := range args {
		var err error

		switch paramType := arg.(type) {
		case uint8:
			err = bw.WriteByte(typeUint8)
			err = bw.WriteByte(arg.(uint8))
		case uint16:
			err = bw.WriteByte(typeUint16)
			err = bw.WriteUInt16(arg.(uint16))
		case uint32:
			err = bw.WriteByte(typeUint32)
			err = bw.WriteUInt32(arg.(uint32))
		case uint64:
			err = bw.WriteByte(typeUint64)
			err = bw.WriteUInt64(arg.(uint64))
		case int8:
			err = bw.WriteByte(typeInt8)
			err = bw.WriteInt8(arg.(int8))
		case int16:
			err = bw.WriteByte(typeInt16)
			err = bw.WriteInt16(arg.(int16))
		case int32:
			err = bw.WriteByte(typeInt32)
			err = bw.WriteInt32(arg.(int32))
		case int64:
			err = bw.WriteByte(typeInt64)
			err = bw.WriteInt64(arg.(int64))
		case float32:
			err = bw.WriteByte(typeFloat32)
			err = bw.WriteFloat32(arg.(float32))
		case float64:
			err = bw.WriteByte(typeFloat64)
			err = bw.WriteFloat64(arg.(float64))
		case string:
			err = bw.WriteByte(typeString)
			err = bw.WriteStr(arg.(string))
		case []byte:
			err = bw.WriteByte(typeBytes)
			err = bw.WriteBytes(arg.([]byte))
		case bool:
			err = bw.WriteByte(typeBool)
			err = bw.WriteBool(arg.(bool))
		case proto.Message:
			err = bw.WriteByte(typeProto)
			err = bw.WriteStr(proto.MessageName(paramType))
			buf, err := proto.Marshal(paramType)
			if err != nil {
				panic(err)
			}

			err = bw.WriteBytes(buf)
		default:
			err = fmt.Errorf("Serialize unsupport type: %T", arg)
		}

		if err != nil {
			panic(err)
		}
	}

	return data
}

// UnSerialize 反序列化
func UnSerialize(data []byte) []interface{} {
	ret := make([]interface{}, 0, 1)
	br := stream.NewByteStream(data)

	for br.ReadEnd() {
		var err error
		var typ uint8
		var v interface{}
		if typ, err = br.ReadByte(); err != nil {
			panic(err)
		}

		switch typ {
		case typeUint8:
			v, err = br.ReadByte()
		case typeUint16:
			v, err = br.ReadUInt16()
		case typeUint32:
			v, err = br.ReadUInt32()
		case typeUint64:
			v, err = br.ReadUInt64()
		case typeInt8:
			v, err = br.ReadInt8()
		case typeInt16:
			v, err = br.ReadInt16()
		case typeInt32:
			v, err = br.ReadInt32()
		case typeInt64:
			v, err = br.ReadInt64()
		case typeFloat32:
			v, err = br.ReadFloat32()
		case typeFloat64:
			v, err = br.ReadFloat64()
		case typeString:
			v, err = br.ReadStr()
		case typeBytes:
			v, err = br.ReadBytes()
		case typeBool:
			v, err = br.ReadBool()
		case typeProto:
			name, err := br.ReadStr()
			if err != nil {
				panic(err)
			}

			buff, err := br.ReadBytes()
			mt := proto.MessageType(name)
			if mt == nil {
				err = fmt.Errorf("unknown message, name: %s", name)
				log.Error("unknown message, name: ", name)
				break
			}

			elem := reflect.New(mt.Elem())
			v = elem.Interface()
			if err := proto.Unmarshal(buff, v.(proto.Message)); err != nil {
				panic(err)
			}
		default:
			err = fmt.Errorf("UnSerialize unsupport type: %d", typ)
		}

		if err == nil {
			ret = append(ret, v)
		} else {
			panic(err)
		}
	}

	if len(ret) == 0 {
		return nil
	}

	return ret
}

// 每个参数需要固定1个字节表示数据类型
// 每种类型的长度不固定
func getSize(args ...interface{}) int {
	size := 0
	for _, arg := range args {
		//类型1个字节
		size++

		switch paramType := arg.(type) {
		case uint8, int8, bool:
			size++
		case uint16, int16:
			size += 2
		case uint32, int32, float32:
			size += 4
		case uint64, int64, float64:
			size += 8
		case string:
			// 字符串需要2个字节标识长度+本身的长度
			size += 2
			size += len(arg.(string))
		case []byte:
			size += 2
			size += len(arg.([]byte))
		case proto.Message:
			size += 2
			size += len(proto.MessageName(paramType))

			size += 2
			size += proto.Size(paramType)
		default:
			panic(fmt.Errorf("getSize unsupport type: %T", arg))
		}
	}

	return size
}

//PointerToValue 取出指针所指向的值
func PointerToValue(v interface{}) interface{} {
	out := reflect.ValueOf(v)
	if out.Kind() == reflect.Ptr {
		out = out.Elem()
	}

	return out.Interface()
}
