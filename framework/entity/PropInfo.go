package entity

import (
	"fmt"

	"github.com/giant-tech/go-service/base/stream"

	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"github.com/globalsign/mgo/bson"
	"github.com/giant-tech/go-service/base/utility"
)

// PropInfo 属性相关
type PropInfo struct {
	value    interface{}
	syncFlag bool
	dbFlag   bool
	def      *PropDef
}

func newPropInfo(def *PropDef) *PropInfo {
	prop := &PropInfo{
		value:    nil,
		syncFlag: false,
		dbFlag:   false,
		def:      def,
	}
	prop.init()
	return prop
}

func (p *PropInfo) init() {
	if p.def == nil {
		log.Error("属性初始化失败, Def为空")
		return
	}

	switch p.def.TypeName {
	case "bool":
		p.value = utility.Atob(p.def.DefaultValue)
	case "int8":
		p.value = int8(utility.Atoi(p.def.DefaultValue))
	case "int16":
		p.value = int16(utility.Atoi(p.def.DefaultValue))
	case "int32":
		p.value = int32(utility.Atoi(p.def.DefaultValue))
	case "int64":
		p.value = int64(utility.Atoi(p.def.DefaultValue))
	case "uint8":
		p.value = uint8(utility.Atoi(p.def.DefaultValue))
	case "uint16":
		p.value = uint16(utility.Atoi(p.def.DefaultValue))
	case "uint32":
		p.value = uint32(utility.Atoi(p.def.DefaultValue))
	case "uint64":
		p.value = uint64(utility.Atoi(p.def.DefaultValue))
	case "float32":
		p.value = float32(utility.Atof(p.def.DefaultValue))
	case "float64":
		p.value = float64(utility.Atof(p.def.DefaultValue))
	case "string":
		p.value = p.def.DefaultValue
	default:
		var err error
		p.value, err = p.def.CreateInst()
		if err != nil {
			log.Error(err, p.def)
			return
		}
	}
}

// GetValue 获取属性的value
func (p *PropInfo) GetValue() interface{} {
	return p.value
}

// GetValueStreamSize 获取某个属性需要的尺寸
func (p *PropInfo) GetValueStreamSize() int {
	s := 0
	st := p.def.TypeName

	switch st {
	case "int8", "uint8", "bool":
		s = 1
	case "int16", "uint16":
		s = 2
	case "int32", "uint32":
		s = 4
	case "int64", "uint64":
		s = 8
	case "float32":
		s = 4
	case "float64":
		s = 8
	case "string":
		s = len(p.value.(string)) + 2
	default:
		log.Error("Convert proto struct failed", st)
	}

	return s + len(st) + 2
}

// WriteValueToStream 把属性值加入到ByteStream中
func (p *PropInfo) WriteValueToStream(bs *stream.ByteStream) error {
	err := bs.WriteStr(p.def.TypeName)
	if err != nil {
		return err
	}

	st := p.def.TypeName

	switch st {
	case "bool":
		err = bs.WriteBool(bool(p.value.(bool)))
	case "int8":
		err = bs.WriteByte(byte(p.value.(int8)))
	case "int16":
		err = bs.WriteUInt16(uint16(p.value.(int16)))
	case "int32":
		err = bs.WriteUInt32(uint32(p.value.(int32)))
	case "int64":
		err = bs.WriteUInt64(uint64(p.value.(int64)))
	case "uint8":
		err = bs.WriteByte(p.value.(byte))
	case "uint16":
		err = bs.WriteUInt16(p.value.(uint16))
	case "uint32":
		err = bs.WriteUInt32(p.value.(uint32))
	case "uint64":
		err = bs.WriteUInt64(p.value.(uint64))
	case "string":
		err = bs.WriteStr(p.value.(string))
	case "float32":
		err = bs.WriteFloat32(p.value.(float32))
	case "float64":
		err = bs.WriteFloat64(p.value.(float64))
	default:
		log.Error("Convert proto struct failed", st)

	}

	return err
}

// ReadValueFromStream 从Stream中读取属性
func (p *PropInfo) ReadValueFromStream(bs *stream.ByteStream) error {
	st, err := bs.ReadStr()
	if err != nil {
		return err
	}

	switch st {
	case "bool":
		v, err := bs.ReadBool()
		if err != nil {
			return err
		}
		p.value = bool(v)
	case "int8":
		v, err := bs.ReadByte()
		if err != nil {
			return err
		}
		p.value = int8(v)
	case "int16":
		v, err := bs.ReadUInt16()
		if err != nil {
			return err
		}
		p.value = int16(v)
	case "int32":
		v, err := bs.ReadUInt32()
		if err != nil {
			return err
		}
		p.value = int32(v)
	case "int64":
		v, err := bs.ReadUInt64()
		if err != nil {
			return err
		}
		p.value = int64(v)
	case "uint8":
		v, err := bs.ReadByte()
		if err != nil {
			return err
		}
		p.value = v
	case "uint16":
		v, err := bs.ReadUInt16()
		if err != nil {
			return err
		}
		p.value = v
	case "uint32":
		v, err := bs.ReadUInt32()
		if err != nil {
			return err
		}
		p.value = v
	case "uint64":
		v, err := bs.ReadUInt64()
		if err != nil {
			return err
		}
		p.value = v
	case "float32":
		v, err := bs.ReadFloat32()
		if err != nil {
			return err
		}
		p.value = v
	case "float64":
		v, err := bs.ReadFloat64()
		if err != nil {
			return err
		}
		p.value = v
	case "string":
		v, err := bs.ReadStr()
		if err != nil {
			return err
		}
		p.value = v
	default:

	}

	return nil
}

// PackValue 打包value 给Redis用
func (p *PropInfo) PackValue() []byte {
	switch p.def.TypeName {
	case "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "float32", "float64", "bool", "string":
		return []byte(fmt.Sprintf("%v", p.value))
	default:

		return nil
	}
}

// UnPackValue 从Redis中恢复Value
func (p *PropInfo) UnPackValue(data interface{}) {
	switch p.def.TypeName {
	case "bool":
		v, _ := redis.Bool(data, nil)
		p.value = bool(v)
	case "int8":
		v, _ := redis.Int(data, nil)
		p.value = int8(v)
	case "int16":
		v, _ := redis.Int(data, nil)
		p.value = int16(v)
	case "int32":
		v, _ := redis.Int(data, nil)
		p.value = int32(v)
	case "int64":
		p.value, _ = redis.Int64(data, nil)
	case "uint8":
		v, _ := redis.Uint64(data, nil)
		p.value = uint8(v)
	case "uint16":
		v, _ := redis.Uint64(data, nil)
		p.value = uint16(v)
	case "uint32":
		v, _ := redis.Uint64(data, nil)
		p.value = uint32(v)
	case "uint64":
		p.value, _ = redis.Uint64(data, nil)
	case "float32":
		v, _ := redis.Float64(data, nil)
		p.value = float32(v)
	case "float64":
		p.value, _ = redis.Float64(data, nil)
	case "string":
		p.value, _ = redis.String(data, nil)
	default:
		log.Warn("Unsupport prop type", p.def.TypeName)

	}
}

// UnPackMongoValue 从mongodb中恢复Value
func (p *PropInfo) UnPackMongoValue(data interface{}, elems []bson.RawDocElem) {
	if data == nil {
		return
	}

	switch p.def.TypeName {
	case "bool":
		p.value = data.(bool)
	case "int8":
		p.value = int8(data.(int))
	case "int16":
		p.value = int16(data.(int))
	case "int32":
		p.value = int32(data.(int))
	case "int64":
		p.value = data.(int64)
	case "uint8":
		p.value = uint64(data.(int64))
	case "uint16":
		p.value = uint16(data.(int))
	case "uint32":
		p.value = uint32(data.(int))
	case "uint64":
		p.value = uint64(data.(int64))
	case "float32":
		p.value = float32(data.(float64))
	case "float64":
		p.value = float64(data.(float64))
	case "string":
		p.value = data.(string)
	default:

		for _, elem := range elems {
			if elem.Name == p.def.Name {
				//err := bson.Unmarshal(elem.Value.Data, data)
				p.value = elem.Value.Data
				//log.Debug("set prop val: ", p.value, " name: ", p.def.Name, " typename: ", p.def.TypeName, " elems: ", elems)
				break
			}
		}
		log.Warn("prop complex type: ", p.def.TypeName)

	}
}
