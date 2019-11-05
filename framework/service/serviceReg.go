package service

import (
	"fmt"

	"github.com/GA-TECH-SERVER/zeus/framework/idata"
	"github.com/GA-TECH-SERVER/zeus/framework/iserver"

	"github.com/cihub/seelog"
)

// SData 服务注册数据
type SData struct {
	ServiceCtrl   iserver.IServiceCtrl
	ServiceName   string
	ServiceTypeID idata.ServiceType
}

var serviceNameMap map[string]*SData
var serviceTypeMap map[idata.ServiceType]*SData

// init 初始化
func init() {
	serviceNameMap = make(map[string]*SData)
	serviceTypeMap = make(map[idata.ServiceType]*SData)
}

//RegService 注册服务
func RegService(typeID idata.ServiceType, name string, is iserver.IServiceCtrl) error {
	if _, ok := serviceNameMap[name]; ok {
		seelog.Error("service name already register: %s", name)
		return fmt.Errorf("service name already register: %s", name)
	}

	if _, ok := serviceTypeMap[typeID]; ok {
		seelog.Error("service type already register: %d", typeID)
		return fmt.Errorf("service type already register: %d", typeID)
	}

	data := &SData{
		ServiceCtrl:   is,
		ServiceName:   name,
		ServiceTypeID: typeID,
	}

	serviceNameMap[name] = data
	serviceTypeMap[typeID] = data

	return nil
}

// GetServiceByName 根据服务名获取服务信息
func GetServiceByName(name string) *SData {
	if is, ok := serviceNameMap[name]; ok {
		return is
	}

	return nil
}

func getServiceByTypeID(typeID idata.ServiceType) *SData {
	if is, ok := serviceTypeMap[typeID]; ok {
		return is
	}

	return nil
}
