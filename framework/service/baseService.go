package service

import (
	"github.com/spf13/viper"
	"github.com/giant-tech/go-service/framework/entity"
	"github.com/giant-tech/go-service/framework/idata"
	"github.com/giant-tech/go-service/framework/iserver"
	"github.com/giant-tech/go-service/framework/msghandler"

	//dbservice "github.com/giant-tech/go-service/logic/logicredis"
	dbservice "github.com/giant-tech/go-service/base/redisservice"
)

// BaseService 上层服务依赖的基础服务
type BaseService struct {
	msghandler.IRPCHandlers
	iserver.IEntities
	*entity.ProtoType
	serviceName string             //服务名
	serviceInfo *idata.ServiceInfo //服务信息
}

// InitBaseService 初始化
func (base *BaseService) InitBaseService(serviceName string, serviceType idata.ServiceType, ilocal iserver.ILocalService) error {
	id, err := dbservice.GetIDGenerator().GetGlobalID()
	if err != nil {
		return err
	}

	base.serviceName = serviceName
	base.serviceInfo = idata.NewServiceInfo(id, serviceType, iserver.GetApp().GetAppID())

	base.IRPCHandlers = msghandler.NewRPCHandlers()

	base.ProtoType = entity.NewProtoType()

	isMultiThread := viper.GetBool(serviceName + ".EntityMultiThread")
	base.IEntities = entity.NewEntities(isMultiThread, ilocal)

	return nil
}

// GetSID 获取服务ID
func (base *BaseService) GetSID() uint64 {
	return base.serviceInfo.ServiceID
}

// GetSType 获取服务type
func (base *BaseService) GetSType() idata.ServiceType {
	return base.serviceInfo.Type
}

// GetSName 获取服务名
func (base *BaseService) GetSName() string {
	return base.serviceName
}

// GetServiceInfo 获取服务信息
func (base *BaseService) GetServiceInfo() *idata.ServiceInfo {
	return base.serviceInfo
}

// SetMetadata 设置元数据
func (base *BaseService) SetMetadata(key, value string) {
	base.serviceInfo.Metadata[key] = value
}

// GetMetadata 获取元数据
func (base *BaseService) GetMetadata(key string) string {
	return base.serviceInfo.Metadata[key]
}
