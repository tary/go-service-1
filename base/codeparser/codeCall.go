package codeparser

import (
	"fmt"
	"reflect"
	"runtime/debug"

	"github.com/giant-tech/go-service/base/serializer"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

// CodeNode 属性或者函数
type CodeNode struct {
	NodeName string
	Params   []string
	IsFunc   bool
}

// ExcuteCode 执行代码
func ExcuteCode(val reflect.Value, nodes []*CodeNode) (retStr string, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error("ExcuteCode panic:", e, string(debug.Stack()))
			if viper.GetString("Config.Recover") == "0" {
				panic(e)
			}

			str := fmt.Sprintf("%+v", e)
			err = fmt.Errorf(str)
		}
	}()

	for idx, node := range nodes {
		//属性的处理
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		//如果是struct，可以继续解析
		if val.Kind() == reflect.Struct {
			if node.IsFunc {
				if val.Kind() != reflect.Ptr {
					val = val.Addr()
				}

				log.Debug("ExcuteCode isFunc, val: ", val)

				method := val.MethodByName(node.NodeName)
				if !method.IsValid() {
					log.Debug("ExcuteCode return 1: ", retStr)
					return retStr, fmt.Errorf("functionName not found: %s", node.NodeName)
				}

				args, err := serializer.UnSerializGMByFunc(method, node.Params)
				if err != nil {
					log.Debug("ExcuteCode return 2: ", retStr)
					return retStr, err
				}

				ret := method.Call(args)
				if len(ret) != 1 {
					log.Debug("return len: ", len(ret))

					for _, v := range ret {
						retStr += fmt.Sprintf("%#v, ", v)
					}

					if len(nodes) != idx+1 {
						log.Debug("ExcuteCode return 3: ", retStr)
						return retStr, fmt.Errorf("return len error")
					}

					return retStr, nil
				}

				val = ret[0]
			} else {
				//属性的处理
				if val.Kind() == reflect.Ptr {
					val = val.Elem()
				}

				//log.Debugf("fieldName: %s, field: %#v", node.NodeName, val)
				val = val.FieldByName(node.NodeName)
				log.Debugf("fieldName: %s, field: %#v", node.NodeName, val)
				if !val.IsValid() {
					log.Debug("ExcuteCode return 6: ", retStr)
					return retStr, fmt.Errorf("propName not found: %s", node.NodeName)
				}
			}
		} else {
			retStr = fmt.Sprintf("%#v", val)
			log.Debug("ExcuteCode return 8: ", retStr)
			return retStr, nil
		}
	}

	retStr = fmt.Sprintf("%#v", val)
	return retStr, nil
}
