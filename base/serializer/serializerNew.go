package serializer

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/giant-tech/go-service/base/imsg"
	"github.com/giant-tech/go-service/base/stream"

	"github.com/cihub/seelog"
)

// SerializeNew 序列化
func SerializeNew(args ...interface{}) []byte {
	size := GetSizeNew(args...)
	data := make([]byte, size)
	bw := stream.NewByteStream(data)
	var err error

	for _, arg := range args {
		val := reflect.ValueOf(arg)
		err = serialize(val, bw)

		if err != nil {
			panic(err)
		}
	}

	return data
}

// SerializeNewWithBuff 带buff的序列化
func SerializeNewWithBuff(buff []byte, args ...interface{}) []byte {
	if buff == nil {
		size := GetSizeNew(args...)
		buff = make([]byte, size)
	}

	bw := stream.NewByteStream(buff)
	var err error

	for _, arg := range args {
		val := reflect.ValueOf(arg)
		err = serialize(val, bw)

		if err != nil {
			panic(err)
		}
	}

	return buff
}

// UnSerializeNew 根据所传数据反序列化, arg 非nil数据
func UnSerializeNew(arg interface{}, data []byte) error {
	if arg == nil {
		panic("UnSerializeStruct val is nil")
	}

	val := reflect.ValueOf(arg)
	if val.Kind() == reflect.Struct {
		panic("UnSerializeNew: arg can't be struct, please change to Ptr")
	}

	br := stream.NewByteStream(data)

	return unserialize(val, br)
}

// UnSerializeByFunc 根据函数参数反序列化, val 为函数的Value
func UnSerializeByFunc(val reflect.Value, data []byte) ([]reflect.Value, error) {
	ret := make([]reflect.Value, 0, 1)
	br := stream.NewByteStream(data)

	mt := val.Type()
	for i := 0; i < mt.NumIn(); i++ {
		pt := mt.In(i)
		elem := newRealValue(pt)

		err := unserialize(elem, br)
		if err != nil {
			return ret, err
		}

		ret = append(ret, elem)
	}

	if len(ret) == 0 {
		return nil, nil
	}

	return ret, nil
}

// UnSerializeJSONByFunc 根据函数参数反序列化Json, val 为函数的Value
func UnSerializeJSONByFunc(val reflect.Value, data []byte) ([]reflect.Value, error) {
	var ret []reflect.Value

	mt := val.Type()
	if mt.NumIn() == 0 {
		return ret, nil
	}

	pt := mt.In(0)
	elem := newRealValue(pt)

	json.Unmarshal(data, elem.Interface())
	ret = append(ret, elem)

	return ret, nil
}

// UnSerializGMByFunc 根据函数参数反序列化, val 为函数的Value
func UnSerializGMByFunc(val reflect.Value, params []string) ([]reflect.Value, error) {
	var ret []reflect.Value
	mt := val.Type()

	if mt.NumIn() != len(params) {
		return ret, nil
	}

	for i := 0; i < mt.NumIn(); i++ {
		pt := mt.In(i)
		elem := newRealValue(pt)

		unserializeString(elem, params[i])

		ret = append(ret, elem)
	}

	return ret, nil
}

// GetSizeNew 获取参数占用的总字节数
func GetSizeNew(args ...interface{}) int {
	var len int
	for _, arg := range args {
		val := reflect.ValueOf(arg)
		if val.Kind() == reflect.Struct {
			panic("GetSizeNew: arg can't be struct, please change to Ptr")
		}

		getValueSize(val, &len)
	}

	return len
}

// serialize
func serialize(val reflect.Value, bw *stream.ByteStream) error {
	var err error
	switch val.Kind() {
	case reflect.Bool:
		err = bw.WriteBool(val.Bool())
	case reflect.Int8:
		err = bw.WriteInt8(int8(val.Int()))
	case reflect.Int16:
		err = bw.WriteInt16(int16(val.Int()))
	case reflect.Int32:
		err = bw.WriteInt32(int32(val.Int()))
	case reflect.Int64:
		err = bw.WriteInt64(val.Int())
	case reflect.Uint8:
		err = bw.WriteByte(uint8(val.Uint()))
	case reflect.Uint16:
		err = bw.WriteUInt16(uint16(val.Uint()))
	case reflect.Uint32:
		err = bw.WriteUInt32(uint32(val.Uint()))
	case reflect.Uint64:
		err = bw.WriteUInt64(val.Uint())
	case reflect.Float32:
		err = bw.WriteFloat32(float32(val.Float()))
	case reflect.Float64:
		err = bw.WriteFloat64(val.Float())
	case reflect.String:
		err = bw.WriteStr(val.String())
	case reflect.Ptr:
		//proto特殊判断
		if protoMsg, ok := val.Interface().(imsg.IProtoMsg); ok {
			tmpData := make([]byte, protoMsg.Size())
			_, msErr := protoMsg.MarshalTo(tmpData)
			if msErr != nil {
				panic(msErr)
			}
			err = bw.WriteBytes(tmpData)
		} else {
			err = serialize(val.Elem(), bw)
		}
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			err = serialize(val.Field(i), bw)
		}
	case reflect.Slice:
		//[]byte优化
		if byteSlice, ok := val.Interface().([]byte); ok {
			err = bw.WriteBytes(byteSlice)
		} else {
			len := val.Len()
			err = bw.WriteInt16(int16(len))

			for i := 0; i < len; i++ {
				err = serialize(val.Index(i), bw)
			}
		}

	case reflect.Array:
		len := val.Len()
		for i := 0; i < len; i++ {
			err = serialize(val.Index(i), bw)
		}
	case reflect.Map:
		err = bw.WriteInt16(int16(val.Len()))
		for _, key := range val.MapKeys() {
			err = serialize(key, bw)
			val2 := val.MapIndex(key)
			err = serialize(val2, bw)
		}
	default:
		panic("serialize: unknow type: %v" + val.Kind().String())
	}

	return err
}

// unserialize
func unserialize(val reflect.Value, bw *stream.ByteStream) error {
	var errRet error
	var v interface{}
	switch val.Kind() {
	case reflect.Bool:
		v, errRet = bw.ReadBool()
		val.SetBool(v.(bool))
	case reflect.Int8:
		v, errRet = bw.ReadInt8()
		val.SetInt(int64(v.(int8)))
	case reflect.Int16:
		v, errRet = bw.ReadInt16()
		val.SetInt(int64(v.(int16)))
	case reflect.Int32:
		v, errRet = bw.ReadInt32()
		val.SetInt(int64(v.(int32)))
	case reflect.Int64:
		v, errRet = bw.ReadInt64()
		val.SetInt(int64(v.(int64)))
	case reflect.Uint8:
		v, errRet = bw.ReadByte()
		val.SetUint(uint64(v.(uint8)))
	case reflect.Uint16:
		v, errRet = bw.ReadUInt16()
		val.SetUint(uint64(v.(uint16)))
	case reflect.Uint32:
		v, errRet = bw.ReadUInt32()
		val.SetUint(uint64(v.(uint32)))
	case reflect.Uint64:
		v, errRet = bw.ReadUInt64()
		val.SetUint(uint64(v.(uint64)))
	case reflect.Float32:
		v, errRet = bw.ReadFloat32()
		val.SetFloat(float64(v.(float32)))
	case reflect.Float64:
		v, errRet = bw.ReadFloat64()
		val.SetFloat(v.(float64))
	case reflect.String:
		v, errRet = bw.ReadStr()
		val.SetString(v.(string))
	case reflect.Ptr:
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}

		//proto特殊判断
		if protoMsg, ok := val.Interface().(imsg.IProtoMsg); ok {
			buff, err := bw.ReadBytes()
			if err == nil {
				errRet = protoMsg.Unmarshal(buff)
			} else {
				errRet = err
			}
		} else {
			errRet = unserialize(val.Elem(), bw)
		}
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			errRet = unserialize(val.Field(i), bw)
			if errRet != nil {
				break
			}
		}
	case reflect.Slice:
		//[]byte优化
		if byteSlice, ok := val.Interface().([]byte); ok {
			byteSlice, errRet = bw.ReadBytes()
			val.Set(reflect.ValueOf(byteSlice))
		} else {
			var len uint16
			len, errRet = bw.ReadUInt16()
			sliceVal := reflect.MakeSlice(val.Type(), int(len), int(len))
			val.Set(sliceVal)

			for i := 0; i < int(len); i++ {
				errRet = unserialize(val.Index(i), bw)
				if errRet != nil {
					break
				}
			}
		}
	case reflect.Array:
		len := val.Len()
		for i := 0; i < len; i++ {
			valElem := val.Index(i)
			errRet = unserialize(valElem, bw)
		}
	case reflect.Map:
		var len uint16
		len, errRet = bw.ReadUInt16()
		if errRet != nil {
			break
		}

		newMap := reflect.MakeMapWithSize(val.Type(), int(len))
		for i := uint16(0); i < len; i++ {
			keyV := newRealValue(val.Type().Key())
			errRet = unserialize(keyV, bw)
			valV := newRealValue(val.Type().Elem())
			errRet = unserialize(valV, bw)

			newMap.SetMapIndex(keyV, valV)
		}

		val.Set(newMap)

	default:
		errRet = fmt.Errorf("unserialize: unknow type: " + val.Kind().String())
		//panic("serialize: unknow type: %v" + val.Kind().String())
	}

	if errRet != nil {
		seelog.Error("unserialize error: ", errRet)
		//panic("unserialize error: " + val.Kind().String())
	}

	return errRet
}

// unserialize
func unserializeString(val reflect.Value, str string) error {
	var err error
	switch val.Kind() {
	case reflect.Bool:
		v, errRet := strconv.ParseBool(str)
		val.SetBool(v)
		err = errRet
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, errRet := strconv.Atoi(str)
		val.SetInt(int64(v))
		err = errRet
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, errRet := strconv.Atoi(str)
		val.SetUint(uint64(v))
		err = errRet
	case reflect.Float32, reflect.Float64:
		v, errRet := strconv.ParseFloat(str, 64)
		val.SetFloat(v)
		err = errRet
	case reflect.String:
		val.SetString(str)
	case reflect.Ptr:
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}

		err = unserializeString(val.Elem(), str)

	default:
		err = fmt.Errorf("unserialize: unknow type: " + val.Kind().String())
		//panic("serialize: unknow type: %v" + val.Kind().String())
	}

	if err != nil {
		seelog.Error("unserialize error: ", err)
		//panic("unserialize error: " + val.Kind().String())
	}

	return err
}

// getValueSize 获取占用字节数
func getValueSize(val reflect.Value, size *int) error {
	var err error
	switch val.Kind() {
	case reflect.Bool, reflect.Int8, reflect.Uint8:
		*size++
	case reflect.Int16, reflect.Uint16:
		*size += 2
	case reflect.Int32, reflect.Uint32, reflect.Float32:
		*size += 4
	case reflect.Int64, reflect.Uint64, reflect.Float64:
		*size += 8
	case reflect.String:
		*size += 2
		*size += val.Len()
	case reflect.Ptr:
		if protoMsg, ok := val.Interface().(imsg.IProtoMsg); ok {
			*size += 2
			*size += protoMsg.Size()
		} else {
			getValueSize(val.Elem(), size)
		}
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			getValueSize(val.Field(i), size)
		}
	case reflect.Slice:
		*size += 2
		len := val.Len()
		for i := 0; i < len; i++ {
			getValueSize(val.Index(i), size)
		}
	case reflect.Array:
		len := val.Len()
		for i := 0; i < len; i++ {
			getValueSize(val.Index(i), size)
		}
	case reflect.Map:
		*size += 2
		for _, key := range val.MapKeys() {
			getValueSize(key, size)
			val2 := val.MapIndex(key)
			getValueSize(val2, size)
		}

	default:
		panic("getValueSize: unknow type: " + val.Kind().String())
	}

	return err
}

//newRealValue 如果是指针
func newRealValue(typ reflect.Type) reflect.Value {
	if typ.Kind() == reflect.Ptr {
		return reflect.New(typ.Elem())
	}
	return reflect.New(typ).Elem()

}

//getRealValue 根据typ获取相匹配的val
func getRealValue(typ reflect.Type, val reflect.Value) reflect.Value {
	if typ.Kind() == reflect.Ptr {
		return val
	}
	return reflect.Indirect(val)

}
