package service

import (
	"fmt"

	"github.com/giant-tech/go-service/base/imsg"
	"github.com/giant-tech/go-service/framework/idata"
	"github.com/giant-tech/go-service/framework/net/inet"
)

// SNilProxy 空服务代理
type SNilProxy struct {
}

// GetSID 获取服务ID
func (s *SNilProxy) GetSID() uint64 {
	return 0
}

// GetSType 获取服务类型
func (s *SNilProxy) GetSType() idata.ServiceType {
	return idata.ServiceType(0)
}

// GetAppID 获取服务所属APP ID
func (s *SNilProxy) GetAppID() uint64 {
	return 0
}

// GetServiceInfo 获取serviceInfo
func (s *SNilProxy) GetServiceInfo() *idata.ServiceInfo {
	return nil
}

// GetMetadata 获取metadata, 传入key返回value
func (s *SNilProxy) GetMetadata(key string) string {
	return ""
}

//GetSess  获取服务代理sess
func (s *SNilProxy) GetSess() inet.ISession {
	return nil
}

// SyncCall 同步调用，等待返回
func (s *SNilProxy) SyncCall(retPtr interface{}, methodName string, args ...interface{}) error {
	return fmt.Errorf("this is nil proxy")
}

// AsyncCall 异步调用，立即返回
func (s *SNilProxy) AsyncCall(methodName string, args ...interface{}) error {
	return fmt.Errorf("this is nil proxy")
}

// SendMsg 发送消息给自己的服务器
func (s *SNilProxy) SendMsg(msg imsg.IMsg) error {
	return fmt.Errorf("this is nil proxy")
}

// IsLocal 是否为本地服务
func (s *SNilProxy) IsLocal() bool {
	return false
}

// IsValid 是否有效
func (s *SNilProxy) IsValid() bool {
	return false
}
