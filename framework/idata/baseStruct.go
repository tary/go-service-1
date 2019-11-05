package idata

import (
	"github.com/GA-TECH-SERVER/zeus/framework/msgdef"
)

// RetData 同步调用返回的数据
type RetData struct {
	Ret []byte
	Err error
}

// CallData 调用数据
type CallData struct {
	Msg     *msgdef.CallMsg //函数调用的消息
	ChanRet chan *RetData   //存放返回值的管道
	Func    func()          //可以直接调用的函数，如果不为nil则忽略其他属性
}

// PendingCall 等待回应的同步调用
type PendingCall struct {
	Seq         uint64      // 序列号
	MethodName  string      // The name of the service and method to call.
	Reply       interface{} // The reply from the function (*struct).
	ToServiceID uint64
	StartTime   int64 // 开始时间
	RetChan     chan *RetData
}

// ServiceInfo 服务信息
type ServiceInfo struct {
	ServiceID uint64            `json:"serviceid"` //service info
	Type      ServiceType       `json:"type"`      //service类型
	AppID     uint64            `json:"appid"`     //appID
	Metadata  map[string]string `json:"metadata"`  //metadata，存储服务额外信息
}

// NewServiceInfo 创建一个ServiceInfo
func NewServiceInfo(ServiceID uint64, Type ServiceType, AppID uint64) *ServiceInfo {
	return &ServiceInfo{
		ServiceID: ServiceID,
		Type:      Type,
		AppID:     AppID,
		Metadata:  make(map[string]string),
	}
}
