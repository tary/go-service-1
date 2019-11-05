package servicedef

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	log "github.com/cihub/seelog"
)

// ServiceDefs 所有的服务定义
type ServiceDefs struct {
	servicedefs map[string]*ServiceDef
}

// inst ServiceDefs实例
var inst *ServiceDefs

// GetServiceDefs 获得service def
func GetServiceDefs() *ServiceDefs {
	return inst
}

// initServiceDefs  初始化服务定义
func initServiceDefs() {
	inst = &ServiceDefs{
		servicedefs: make(map[string]*ServiceDef),
	}

	inst.Init()
}

// Init 初始化文件定义结构
func (defs *ServiceDefs) Init() {
	err := filepath.Walk("../res/servicedef", func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return nil
		}

		if f.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".json") {
			raw, err := ioutil.ReadFile(path)
			if err != nil {
				log.Error("read service def file error ", err)
				return nil
			}

			jsonInfo := make(map[string]interface{})
			err = json.Unmarshal(raw, &jsonInfo)
			if err != nil {
				log.Error("parse service def file error ", err)
				return nil
			}

			def := newServiceDef()
			if err := def.fill(jsonInfo); err != nil {
				log.Error("fill serive def error ", err)
				return nil
			}

			defs.servicedefs[jsonInfo["name"].(string)] = def
		}

		return nil
	})

	log.Debug("def servicedefs= ", defs.servicedefs)
	if err != nil {
		log.Error("walk service def file error")
	}
}

// GetDef 根据名字获得service定义
func (defs *ServiceDefs) GetDef(name string) *ServiceDef {
	d, ok := defs.servicedefs[name]
	if !ok {
		return nil
	}

	return d
}

////////////////////////////////////////////////////////////////

// ServiceDef 单个service定义
type ServiceDef struct {
	Name          string
	Methods       map[string]*MethodDef
	MethodsParams map[string][]string
}

// newServiceDef 新建service def
func newServiceDef() *ServiceDef {
	return &ServiceDef{
		Methods:       make(map[string]*MethodDef),
		MethodsParams: make(map[string][]string),
	}
}

// servicedef 填充servicedef
func (def *ServiceDef) fill(jsonInfo map[string]interface{}) error {

	def.Name = jsonInfo["name"].(string)

	jsonMethods := jsonInfo["methods"].(map[string]interface{})

	for methodName, methodInfo := range jsonMethods {

		//method的参数没有读取， 暂时不会检测参数的一致性，只检测方法有没有实现。
		serviced := newMehodDef()
		serviced.Name = methodName
		def.Methods[methodName] = serviced

		jsonMethod := methodInfo.(map[string]interface{})
		var params []string
		for _, v := range jsonMethod {
			params = append(params, v.(string))
		}
		def.MethodsParams[methodName] = params

	}

	log.Debug("MethodsParams=", def.MethodsParams)

	return nil
}

// getTypeByStr 通过字符串拿类型
func (def *ServiceDef) getTypeByStr(st string) reflect.Type {

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

// MethodDef 字段定义
type MethodDef struct {
	Name   string
	Params []string //param
}

// newMehodDef 新建方法定义
func newMehodDef() *MethodDef {
	return &MethodDef{}
}
