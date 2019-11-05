package msghandler

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/GA-TECH-SERVER/zeus/base/serializer"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

// IGmGmMsgHandlers 消息处理模块的接口
type IGmGmMsgHandlers interface {
	RegGmMsg(proc interface{})

	DoGmMsg(string, []byte) ([]byte, error)
}

// NewGmMsgHandlers 创建一个新的消息处理器
func NewGmMsgHandlers() IGmGmMsgHandlers {
	return &GmMsgHandlers{}
}

// GmMsgHandlers 消息处理中心
type GmMsgHandlers struct {
	gmFuncs sync.Map
}

// RegGmMsg 注册rpc消息处理对象
// 其中 proc 是一个对象，包含是类似于 GmXXXXX的一系列函数，分别用来处理不同的gm消息
func (handlers *GmMsgHandlers) RegGmMsg(proc interface{}) {
	v := reflect.ValueOf(proc)
	t := reflect.TypeOf(proc)

	for i := 0; i < t.NumMethod(); i++ {
		methodName := t.Method(i).Name

		// 判断是否是RPC处理函数
		msgName, msgHandler, err := handlers.getGMHandler(methodName, v.MethodByName(methodName))
		if err == nil {
			handlers.addGMHandler(msgName, msgHandler)
		}
	}
}

// DoGmMsg gm执行
func (handlers *GmMsgHandlers) DoGmMsg(methodName string, data []byte) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("DoGmMsg panic:", err, ", methodName: ", methodName, ", ", string(debug.Stack()))
			if viper.GetString("Config.Recover") == "0" {
				panic(err)
			}
		}
	}()

	ifunc, ok := handlers.gmFuncs.Load(methodName)
	if !ok {
		log.Error("Method ", methodName, " Can't Find")
		return nil, fmt.Errorf("Method %s can't find", methodName)
	}

	gmFunc := ifunc.(reflect.Value)
	args, err := serializer.UnSerializeByFunc(gmFunc, data)
	if err != nil {
		return nil, err
	}

	ret := gmFunc.Call(args)
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

func (handlers *GmMsgHandlers) addGMHandler(msgName string, msgHandler reflect.Value) {
	_, ok := handlers.gmFuncs.Load(msgName)
	if ok {
		return
	}

	handlers.gmFuncs.Store(msgName, msgHandler)
}

func (handlers *GmMsgHandlers) getGMHandler(methodName string, v reflect.Value) (string, reflect.Value, error) {
	if strings.HasPrefix(methodName, "Gm") {
		return methodName[2:], v, nil
	}

	return "", reflect.ValueOf(nil), fmt.Errorf("")
}
