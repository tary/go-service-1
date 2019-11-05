package service

import (
	"fmt"
	"time"

	"github.com/GA-TECH-SERVER/zeus/base/net/inet"
	"github.com/GA-TECH-SERVER/zeus/base/serializer"
	"github.com/GA-TECH-SERVER/zeus/framework/idata"
	"github.com/GA-TECH-SERVER/zeus/framework/iserver"
	"github.com/GA-TECH-SERVER/zeus/framework/msgdef"

	"github.com/cihub/seelog"
)

// SProxy 服务代理
type SProxy struct {
	serviceInfo *idata.ServiceInfo //服务信息
	Sess        inet.ISession
	isLocal     bool // 是否为本进程服务
}

// CreateServiceProxy 创建服务
func CreateServiceProxy(sinfo *idata.ServiceInfo, islocal bool) *SProxy {
	return &SProxy{
		serviceInfo: sinfo,
		isLocal:     islocal,
	}
}

// SyncCall 同步调用，等待返回
func (s *SProxy) SyncCall(retPtr interface{}, methodName string, args ...interface{}) error {
	msg := &msgdef.CallMsg{}

	msg.SID = s.serviceInfo.ServiceID
	msg.MethodName = methodName
	msg.IsSync = true
	msg.Params = serializer.SerializeNew(args...)

	if s.isLocal {
		//直接发送
		is := GetLocalServiceMgr().GetLocalService(s.serviceInfo.ServiceID)
		if is == nil {
			seelog.Error("SProxy.SyncCall, GetLocalService , SID not exist: ", s.serviceInfo.ServiceID)
			return fmt.Errorf(" SID not exist: %d", s.serviceInfo.ServiceID)
		}

		retData := is.PostCallMsgAndWait(msg)
		if retData.Err != nil {
			seelog.Error("SProxy.SyncCall, PostCallMsgAndWait , Err: ", retData.Err)
			return retData.Err
		}

		if retPtr != nil {
			if err := serializer.UnSerializeNew(retPtr, retData.Ret); err != nil {
				return err
			}
		}
	} else {
		msg.Seq = iserver.GetApp().GetSeq()
		s.SendMsg(msg)

		//加入到pending中
		call := &idata.PendingCall{}
		call.RetChan = make(chan *idata.RetData, 1)
		call.Seq = msg.Seq
		call.MethodName = methodName
		call.Reply = retPtr
		call.ToServiceID = s.serviceInfo.ServiceID
		call.StartTime = time.Now().Unix()
		iserver.GetApp().AddPendingCall(call)

		retData := <-call.RetChan
		if retData.Err != nil {
			seelog.Error("SProxy.SyncCall, remote retData.Err: ", retData.Err)
			return retData.Err
		}

		if retPtr != nil {
			if err := serializer.UnSerializeNew(retPtr, retData.Ret); err != nil {
				return err
			}
		}
	}

	return nil
}

// AsyncCall 异步调用，立即返回
func (s *SProxy) AsyncCall(methodName string, args ...interface{}) error {
	msg := &msgdef.CallMsg{}

	msg.IsSync = false
	msg.SID = s.serviceInfo.ServiceID
	msg.MethodName = methodName
	msg.Params = serializer.SerializeNew(args...)

	if s.isLocal {
		//直接发送
		is := GetLocalServiceMgr().GetLocalService(s.serviceInfo.ServiceID)
		if is == nil {
			seelog.Error("")
			return fmt.Errorf("")
		}

		return is.PostCallMsg(msg)
	}
	s.SendMsg(msg)

	return nil
}

// SendMsg 发送消息给自己的服务器
func (s *SProxy) SendMsg(msg inet.IMsg) error {
	if s.Sess == nil {
		seelog.Error("SProxy.SendMsg, Sess is nil")
		return fmt.Errorf("Sess is nil")
	}

	s.Sess.Send(msg)

	return nil
}

// GetSID 获取服务ID
func (s *SProxy) GetSID() uint64 {
	return s.serviceInfo.ServiceID
}

// GetSType 获取服务类型
func (s *SProxy) GetSType() idata.ServiceType {
	return s.serviceInfo.Type
}

// GetAppID 获取服务所属APP ID
func (s *SProxy) GetAppID() uint64 {
	return s.serviceInfo.AppID
}

// GetServiceInfo  获取ServiceInfo
func (s *SProxy) GetServiceInfo() *idata.ServiceInfo {
	return s.serviceInfo
}

// GetMetadata 获取metadata, 传入key返回value
func (s *SProxy) GetMetadata(key string) string {
	if val, ok := s.serviceInfo.Metadata[key]; ok {
		return val
	}

	return ""
}

// GetSess 获取session
func (s *SProxy) GetSess() inet.ISession {
	return s.Sess
}

// IsLocal 是否为本地服务
func (s *SProxy) IsLocal() bool {
	return s.isLocal
}

// IsValid 是否有效
func (s *SProxy) IsValid() bool {
	return true
}

/*func (s *SProxy) OnClosed() {

	//一些逻辑上的善后处理
}
*/
