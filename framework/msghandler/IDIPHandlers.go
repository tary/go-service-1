package msghandler

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"

	"github.com/giant-tech/go-service/base/serializer"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

/*
	MsgHandler 作为底层通讯层与上层应用层之间逻辑传递的桥梁
*/

// IIDIPHandlers 消息处理模块的接口
type IIDIPHandlers interface {
	RegIDIPMsg(proc interface{})
	GetIDIPHandlers() (*sync.Map, error)
	GetFuncJSON() (string, error)
	DoIDIPMsg(string, []byte) ([]byte, error)
}

// NewIDIPHandlers 创建一个新的消息处理器
func NewIDIPHandlers() IIDIPHandlers {
	return &IDIPHandlers{}
}

// IDIPHandlers 消息处理中心
type IDIPHandlers struct {
	idipFuncs sync.Map
}

// RegIDIPMsg 注册idip消息处理对象
// 其中 proc 是一个对象，包含是类似于 IDIPXXXXX的一系列函数，分别用来处理不同的IDIP消息
func (handlers *IDIPHandlers) RegIDIPMsg(proc interface{}) {
	v := reflect.ValueOf(proc)
	t := reflect.TypeOf(proc)

	for i := 0; i < t.NumMethod(); i++ {
		methodName := t.Method(i).Name

		// 判断是否是IDIP处理函数
		msgName, msgHandler, err := handlers.getIDIPHandler(methodName, v.MethodByName(methodName))
		if err == nil {
			handlers.addIDIPHandler(msgName, msgHandler)
		}
	}
}

// DoIDIPMsg 异步IDIP
func (handlers *IDIPHandlers) DoIDIPMsg(methodName string, data []byte) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("DoIDIPMsg panic:", err, ", methodName: ", methodName, ", ", string(debug.Stack()))
			if viper.GetString("Config.Recover") == "0" {
				panic(err)
			}
		}
	}()

	ifunc, ok := handlers.idipFuncs.Load(methodName)
	if !ok {
		log.Error("DoAsyncIDIPMsg, Can't Find Method: ", methodName)
		return nil, fmt.Errorf("Method %s can't find", methodName)
	}

	idipFunc, ok := ifunc.(reflect.Value)
	args, err := serializer.UnSerializeJSONByFunc(idipFunc, data)
	if err != nil {
		log.Error("UnSerializeJSONByFunc err, methodName: ", methodName)
		return nil, err
	}

	ret := idipFunc.Call(args)
	if len(ret) < 2 {
		log.Error("return len error: ", len(ret))
		return nil, fmt.Errorf("return len error")
	}

	if ret[1].Interface() != nil {
		if err, ok = ret[1].Interface().(error); !ok {
			log.Error("return second value is not error type")
			return nil, fmt.Errorf("return second value is not error type")
		}
	}

	if err != nil {
		return nil, err
	}

	if data, ok := ret[0].Interface().([]byte); ok {
		return data, nil
	}

	return nil, fmt.Errorf("return error")
}

// addIDIPHandler 添加idip处理
func (handlers *IDIPHandlers) addIDIPHandler(msgName string, msgHandler reflect.Value) {
	//log.Debug("msgName= ", msgName, " handlers = ", handlers)
	_, ok := handlers.idipFuncs.Load(msgName)
	if ok {
		log.Error("addIDIPHandler err, msgName already registered: ", msgName)
	}

	handlers.idipFuncs.Store(msgName, msgHandler)
	//log.Debug(" store msgName= ", msgName, "func= ", funcs)
}

//getIDIPHandler 获取idip处理
func (handlers *IDIPHandlers) getIDIPHandler(methodName string, v reflect.Value) (string, reflect.Value, error) {
	if strings.HasPrefix(methodName, "IDIP") {
		return methodName[4:], v, nil
	}

	return "", reflect.ValueOf(nil), fmt.Errorf("")
}

//GetIDIPHandlers 获取idip处理
func (handlers *IDIPHandlers) GetIDIPHandlers() (*sync.Map, error) {

	return &handlers.idipFuncs, nil
}

//GetFuncJSON 获取func json
func (handlers *IDIPHandlers) GetFuncJSON() (string, error) {
	var err error
	allStr := `{ "Funcs":[`
	allPre := ``

	structMap := make(map[string]string)

	handlers.idipFuncs.Range(func(k, v interface{}) bool {
		methodName, _ := k.(string)
		idipFunc, _ := v.(reflect.Value)
		mt := idipFunc.Type()
		if mt.NumIn() > 1 {
			err = fmt.Errorf("param num error: %d", mt.NumIn())
			log.Error("param num error: ", mt.NumIn())
			return false
		}

		paramName := ""

		if mt.NumIn() > 0 {
			pt := mt.In(0)
			paramName = GetTypeName(pt)
			_, err = handlers.GetTypeJSON(structMap, pt)
			if err != nil {
				log.Error("GetTypeJSON error: ", err)
				return false
			}
		}

		allStr += allPre
		allStr += `{`
		allStr += `"FuncName":"` + methodName + `",`
		allStr += `"TypeName":"` + paramName + `"`
		allStr += `}`
		allPre = `,`

		return true
	})

	allStr += `],`

	allStr += `"DataTypes":` + "["

	preStr := ""
	for _, v := range structMap {
		allStr += preStr
		allStr += v
		preStr = ","
	}

	allStr += "]}"

	return allStr, err
}

// GetTypeName 获取类型名
func GetTypeName(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		return GetTypeName(t.Elem())
	}

	return t.String()
}

// GetType 获得类型
func GetType(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		return GetType(t.Elem())
	}

	return strconv.Itoa(int(t.Kind()))
}

// GetTypeJSON 获得类型json
func (handlers *IDIPHandlers) GetTypeJSON(typeMap map[string]string, pt reflect.Type) (string, error) {
	var jsonStr string
	var errRet error

	switch pt.Kind() {
	case reflect.Bool:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
	case reflect.Float32:
	case reflect.Float64:
	case reflect.String:
		break

	case reflect.Ptr:
		jsonStr, errRet = handlers.GetTypeJSON(typeMap, pt.Elem())
	case reflect.Struct:
		preStr := ""
		structStr := `{`
		structStr += `"Name": "` + GetTypeName(pt) + `",`
		structStr += `"Type": ` + GetType(pt) + `,`
		structStr += `"Fields":[`
		for i := 0; i < pt.NumField(); i++ {
			f := pt.Field(i)

			tag := f.Tag.Get("json")

			if tag == "-" {
				continue
			}

			jsonName, _ := parseTag(tag)

			fieldStr := `{`
			//Type
			fieldStr += `"TypeName":"` + GetTypeName(f.Type) + `"`

			//JsonName
			if jsonName == "" {
				jsonName = f.Name
			}

			fieldStr += `,` + `"JsonName":"` + jsonName + `"`

			//DocName
			docName := f.Tag.Get("doc")
			if docName == "" {
				docName = jsonName
			}

			fieldStr += `,` + `"DocName":"` + docName + `"`

			fieldStr += `}`

			structStr += preStr
			structStr += fieldStr
			preStr = `,`

			handlers.GetTypeJSON(typeMap, pt.Field(i).Type)
		}
		structStr += `]`
		structStr += `}`
		typeMap[GetTypeName(pt)] = structStr

	case reflect.Slice:
		sliceStr := `{`
		sliceStr += `"Name": "` + GetTypeName(pt) + `",`
		sliceStr += `"Type": ` + GetType(pt) + `,`
		sliceStr += `"Value":"` + GetTypeName(pt.Elem()) + `"`
		sliceStr += `}`

		typeMap[GetTypeName(pt)] = sliceStr

		handlers.GetTypeJSON(typeMap, pt.Elem())

	case reflect.Map:
		mapStr := `{`
		mapStr += `"Name": "` + GetTypeName(pt) + `",`
		mapStr += `"Type": ` + GetType(pt) + `,`
		mapStr += `"Key":"` + GetTypeName(pt.Key()) + `",`
		mapStr += `"Value":"` + GetTypeName(pt.Elem()) + `"`
		mapStr += `}`

		typeMap[GetTypeName(pt)] = mapStr

		handlers.GetTypeJSON(typeMap, pt.Key())
		handlers.GetTypeJSON(typeMap, pt.Elem())

	default:
		errRet = fmt.Errorf("unserialize: unknow type: " + pt.Kind().String())
	}

	return jsonStr, errRet
}

// tagOptions is the string following a comma in a struct field's "json"
// tag, or the empty string. It does not include the leading comma.
type tagOptions string

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(tag[idx+1:])
	}
	return tag, tagOptions("")
}
