package msghandler

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/giant-tech/go-service/base/serializer"
	"github.com/giant-tech/go-service/framework/idata"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

/*
	MsgHandler 作为底层通讯层与上层应用层之间逻辑传递的桥梁
*/

// IRPCHandlers 消息处理模块的接口
type IRPCHandlers interface {
	RegRPCMsg(proc interface{})
	GetRPCHandlers() (*sync.Map, error)
	GetRPCMethodParams() (*sync.Map, error)
	DoRPCMsg(methodName string, data []byte, chanRet chan *idata.RetData) error
	ClearRPCData()
}

// NewRPCHandlers 创建一个新的消息处理器
func NewRPCHandlers() IRPCHandlers {
	return &RPCHandlers{}
}

// RPCHandlers 消息处理中心
type RPCHandlers struct {
	rpcFuncs     sync.Map
	methodParams sync.Map //函数名，函数的参数列表
}

// RegRPCMsg 注册rpc消息处理对象
// 其中 proc 是一个对象，包含是类似于 RPCXXXXX的一系列函数，分别用来处理不同的RPC消息
func (handlers *RPCHandlers) RegRPCMsg(proc interface{}) {
	v := reflect.ValueOf(proc)
	t := reflect.TypeOf(proc)

	for i := 0; i < t.NumMethod(); i++ {
		methodName := t.Method(i).Name

		// 判断是否是RPC处理函数
		msgName, msgHandler, err := handlers.getRPCHandler(methodName, v.MethodByName(methodName))
		if err == nil {
			handlers.addRPCHandler(msgName, msgHandler)
			//加入参数列表
			/*for i := 0; i < msgHandler.NumField(); i++ {
				handlers.methodParams.Store(msgName, v.Field(i))
			}
			*/
			var params []string
			mt := msgHandler.Type()
			for i := 0; i < mt.NumIn(); i++ {
				pt := mt.In(i)
				//log.Debug("methodname = ", methodName, " ,kind = ", pt.Kind())
				if pt.Kind() == reflect.Ptr && pt.Elem().Kind() == reflect.Struct {
					params = append(params, "struct")
				} else {
					params = append(params, pt.Name())
				}

			}
			//log.Debug("methodname = ", methodName, " ,params = ", params)
			handlers.methodParams.Store(msgName, params)
		}
	}
}

// DoRPCMsg 处理rpc消
// chanRet 为nil则是异步调用，否则为同步调用
func (handlers *RPCHandlers) DoRPCMsg(methodName string, data []byte, chanRet chan *idata.RetData) error {
	//异步调用
	if chanRet == nil {
		return handlers.doAsyncRPCMsg(methodName, data)
	}

	//同步调用
	return handlers.doSyncRPCMsg(methodName, data, chanRet)
}

// doAsyncRPCMsg 异步RPC
func (handlers *RPCHandlers) doAsyncRPCMsg(methodName string, data []byte) error {
	defer func() {
		if err := recover(); err != nil {
			log.Error("doAsyncRPCMsg panic:", err, ", methodName: ", methodName, ", ", string(debug.Stack()))
			if viper.GetString("Config.Recover") == "0" {
				panic(err)
			}
		}
	}()

	ifunc, ok := handlers.rpcFuncs.Load(methodName)
	if !ok {
		log.Error("doAsyncRPCMsg, Can't Find Method: ", methodName)
		return fmt.Errorf("Method %s can't find", methodName)
	}

	rpcFunc, ok := ifunc.(reflect.Value)

	args, err := serializer.UnSerializeByFunc(rpcFunc, data)
	if err != nil {
		log.Error("doAsyncRPCMsg, method: ", methodName, ", err: ", err)
		return err
	}

	rpcFunc.Call(args)

	return err
}

// doSyncRPCMsg 同步RPC
func (handlers *RPCHandlers) doSyncRPCMsg(methodName string, data []byte, chanRet chan *idata.RetData) error {

	defer func() {
		if err := recover(); err != nil {
			log.Error("doSyncRPCMsg panic:", err, ", methodName: ", methodName, ", ", string(debug.Stack()))
			if viper.GetString("Config.Recover") == "0" {
				panic(err)
			}
		}
	}()

	retData := &idata.RetData{}
	var err error
	ifunc, ok := handlers.rpcFuncs.Load(methodName)
	if !ok {
		//log.Error("Method ", methodName, " Can't Find", " handlers = ", handlers)
		err = fmt.Errorf("Method %s can't find", methodName)
	}

	if err == nil {
		rpcFunc, _ := ifunc.(reflect.Value)

		args, err1 := serializer.UnSerializeByFunc(rpcFunc, data)
		if err1 == nil {
			ret := rpcFunc.Call(args)

			if len(ret) == 1 {
				//TODO: 为Center定制的临时补丁，后面会删除
				if _, ok := ret[0].Interface().(*idata.RetData); ok {
					retData = ret[0].Interface().(*idata.RetData)
				} else {
					retData.Ret = serializer.SerializeNew(ret[0].Interface())
				}
			} else if len(ret) == 0 {
			} else {
				//暂时不支持多参数
				err = fmt.Errorf("Return params error, len: %d", len(ret))
			}
		}
	}

	if err != nil {
		retData.Err = err
	}

	chanRet <- retData

	return nil
}

// addRPCHandler 添加rpc处理
func (handlers *RPCHandlers) addRPCHandler(msgName string, msgHandler reflect.Value) {
	//log.Debug("msgName= ", msgName, " handlers = ", handlers)
	_, ok := handlers.rpcFuncs.Load(msgName)
	if ok {
		log.Error("addRPCHandler err, msgName already registered: ", msgName)
	}

	handlers.rpcFuncs.Store(msgName, msgHandler)
	//log.Debug(" store msgName= ", msgName, "func= ", funcs)
}

// getRPCHandler 获得rpc处理
func (handlers *RPCHandlers) getRPCHandler(methodName string, v reflect.Value) (string, reflect.Value, error) {
	if strings.HasPrefix(methodName, "OnDisconnected") {
		return methodName[2:], v, nil
	}

	if strings.HasPrefix(methodName, "OnConnected") {
		return methodName[2:], v, nil
	}

	if strings.HasPrefix(methodName, "RPC") {
		return methodName[3:], v, nil
	}

	return "", reflect.ValueOf(nil), fmt.Errorf("")
}

// GetRPCHandlers 获得rpc处理
func (handlers *RPCHandlers) GetRPCHandlers() (*sync.Map, error) {

	return &handlers.rpcFuncs, nil
}

// GetRPCMethodParams 获得rpc方法参数
func (handlers *RPCHandlers) GetRPCMethodParams() (*sync.Map, error) {

	return &handlers.methodParams, nil
}

// ClearRPCData 清除RPCHandlers的注册内容
func (handlers *RPCHandlers) ClearRPCData() {
	handlers.rpcFuncs = sync.Map{}
	handlers.methodParams = sync.Map{}
}
