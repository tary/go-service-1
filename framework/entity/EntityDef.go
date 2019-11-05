package entity

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	log "github.com/cihub/seelog"
)

// Defs 所有的实体定义
type Defs struct {
	defs map[string]*Def
}

var inst *Defs

// GetDefs 获取实体定义信息
func GetDefs() *Defs {
	return inst
}

func initDefs() {
	inst = &Defs{
		defs: make(map[string]*Def),
	}

	inst.Init()
}

// Init 初始化文件定义结构
func (defs *Defs) Init() {
	err := filepath.Walk("../res/entitydef", func(path string, f os.FileInfo, err error) error {

		if f == nil {
			return nil
		}
		if f.Name() != "alias.json" {

			if f.IsDir() {
				return nil
			}

			if strings.HasSuffix(path, ".json") {
				raw, err := ioutil.ReadFile(path)
				if err != nil {
					log.Error("read entity def file error ", err)
					return nil
				}

				jsonInfo := make(map[string]interface{})
				err = json.Unmarshal(raw, &jsonInfo)
				if err != nil {
					log.Error("parse entity def file error ", err)
					return nil
				}

				def := newDef()
				if err := def.fill(jsonInfo); err != nil {
					log.Error("fill def error ", err)
					return nil
				}

				defs.defs[jsonInfo["name"].(string)] = def
			}

		}
		return nil
	})

	if err != nil {
		log.Error("walk entity def file error")
	}
}

// GetDef 获取一个entity定义
func (defs *Defs) GetDef(name string) *Def {
	d, ok := defs.defs[name]
	if !ok {
		return nil
	}

	return d
}

////////////////////////////////////////////////////////////////

// Def 实体定义
type Def struct {
	Name  string
	Props map[string]*PropDef
}

// 新建实体定义
func newDef() *Def {
	return &Def{
		Props: make(map[string]*PropDef),
	}
}

// 填充实体属性描述等信息
func (def *Def) fill(jsonInfo map[string]interface{}) error {

	def.Name = jsonInfo["name"].(string)

	jsonProps := jsonInfo["props"].(map[string]interface{})

	sync := jsonInfo["sync"].(map[string]interface{})

	for propName, propInfo := range jsonProps {
		jsonProp := propInfo.(map[string]interface{})

		prop := newPropDef()
		def.Props[propName] = prop
		prop.Name = propName
		prop.Desc = jsonProp["desc"].(string)
		prop.Type = def.getTypeByStr(jsonProp["type"].(string))
		prop.TypeName = jsonProp["type"].(string)

		defautV := jsonProp["default"]
		if defautV != nil {
			prop.DefaultValue = defautV.(string)
		}

		if strings.Contains(prop.TypeName, "protoMsg") {
			prop.TypeName = prop.TypeName[10:]
		}
		prop.Persistence = true
		if persistence, ok := jsonProp["save"].(string); ok {
			if persistence == "0" {
				prop.Persistence = false
			}
		}

		for typeString, syncInfo := range sync {
			//syncType := typeString
			propList, ok := syncInfo.(map[string]interface{})["props"].([]interface{})
			if !ok {
				log.Error("分析服务关心的实体属性列表失败", typeString, syncInfo)
				continue
			}

			for _, srvPropName := range propList {
				if srvPropName.(string) == propName {
					syncType, _ := strconv.Atoi(typeString)
					prop.Sync = append(prop.Sync, uint32(syncType))
					break
				}
			}
		}
	}

	return nil
}

func (def *Def) getTypeByStr(st string) reflect.Type {

	var t reflect.Type

	switch st {
	case "int8":
		t = reflect.TypeOf(int8(0))
	case "int16":
		t = reflect.TypeOf(int16(0))
	case "int32":
		t = reflect.TypeOf(int32(0))
	case "int64":
		t = reflect.TypeOf(int64(0))
	case "byte":
		t = reflect.TypeOf(byte(0))
	case "uint8":
		t = reflect.TypeOf(uint8(0))
	case "uint16":
		t = reflect.TypeOf(uint16(0))
	case "uint32":
		t = reflect.TypeOf(uint32(0))
	case "uint64":
		t = reflect.TypeOf(uint64(0))
	case "float32":
		t = reflect.TypeOf(float32(0))
	case "float64":
		t = reflect.TypeOf(float64(0))
	case "string":
		t = reflect.TypeOf("")
	case "bool":
		t = reflect.TypeOf(false)
	default:
		return nil
	}

	return t
}

////////////////////////////////////////////////////////////////

// PropDef 字段定义
type PropDef struct {
	Name         string
	Desc         string
	Type         reflect.Type
	TypeName     string
	DefaultValue string   //默认值
	Persistence  bool     //是否需要持久化, 默认true
	Sync         []uint32 //同步给谁
}

func newPropDef() *PropDef {
	return &PropDef{}
}

// CreateInst 创建该属性实例
func (pd *PropDef) CreateInst() (interface{}, error) {

	if pd.Type != nil {
		return reflect.New(pd.Type), nil
	}

	return nil, nil
}

// IsValidValue 当前值是否能设置到该属性上
func (pd *PropDef) IsValidValue(value interface{}) bool {

	if value == nil {
		return true
	}

	//log.Debug("bbbbbbb, pd.Type = ", pd.Type, ",  reflect.TypeOf(value).String() : ", reflect.TypeOf(value).String(), ",  pd.TypeName: ", pd.TypeName)
	if pd.Type != nil {
		return pd.Type == reflect.TypeOf(value)
	}

	// 因map类型没法反射做验证, 暂时注释
	//return strings.Contains(reflect.TypeOf(value).String(), pd.TypeName)
	return true
}
