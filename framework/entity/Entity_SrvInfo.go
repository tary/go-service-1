package entity

import (
	"github.com/cihub/seelog"
	"github.com/giant-tech/go-service/base/imsg"
	"github.com/giant-tech/go-service/framework/idata"
	dbservice "github.com/giant-tech/go-service/framework/logicredis"
)

// isEntityExisted isEntityExisted
func (e *Entity) isEntityExisted(srvType uint8) bool {

	e.srvIDSMux.RLock()
	defer e.srvIDSMux.RUnlock()

	_, ok := e.srvIDS[srvType]

	return ok
}

// RefreshSrvIDS 从redis上更新 当前 entity所有的分布式信息
func (e *Entity) RefreshSrvIDS() {
	e.srvIDSMux.Lock()
	defer e.srvIDSMux.Unlock()

	srvIDs, err := dbservice.GetEntitySrvUtil(e.entityID).GetSrvIDs()
	if err != nil {
		seelog.Error("Get entity srv info failed")
		return
	}

	e.srvIDS = srvIDs

	//ghost不更新cellID
	/*if e.IsGhost() {
		return
	}

	e.cellID = 0
	e.cellSrvType = 0

	for srvType, info := range e.srvIDS {

		//e.Debug("RefreshSrvIDS,  CellID: ", info.CellID, "srvType: ", srvType)

		if info.CellID != 0 {
			e.cellID = info.CellID
			e.cellSrvType = srvType
		}
	}*/

}

// RegSrvID 将当前部分的Entity注册到Redis上
func (e *Entity) RegSrvID() {
	if e.IsGhost() {
		e.RefreshSrvIDS()
	} else {
		if err := dbservice.GetEntitySrvUtil(e.entityID).RegSrvID(
			uint8(e.GetIEntities().GetLocalService().GetSType()),
			e.GetIEntities().GetLocalService().GetSID(),
			e.GetGroupID(),
			e.entityType); err != nil {
			e.Error("Reg SrvID failed ", err)
			return
		}

		if e.GetIEntities().GetLocalService().GetSType() == idata.ServiceGateway {
			if err := dbservice.GetEntitySrvUtil(e.entityID).RegSrvID(
				uint8(idata.ServiceClient),
				e.GetIEntities().GetLocalService().GetSID(),
				e.GetGroupID(),
				e.entityType); err != nil {
				e.Error("Reg SrvID failed ", err)
				return
			}
		}

		e.RefreshSrvIDS()
		e.broadcastSrvInfo()
	}
}

// UnregSrvID 将当前部分的Entity从Redis上删除
func (e *Entity) UnregSrvID() {
	if err := dbservice.GetEntitySrvUtil(e.entityID).UnRegSrvID(
		uint8(e.GetIEntities().GetLocalService().GetSType()),
		e.GetIEntities().GetLocalService().GetSID(),
		e.GetGroupID()); err != nil {
		e.Error("Unreg SrvID failed ", err)
	}

	if e.GetIEntities().GetLocalService().GetSType() == idata.ServiceGateway {
		if err := dbservice.GetEntitySrvUtil(e.entityID).UnRegSrvID(
			uint8(idata.ServiceClient),
			e.GetIEntities().GetLocalService().GetSID(),
			e.GetGroupID()); err != nil {
			e.Error("Unreg SrvID failed ", err)
		}
	}

	e.broadcastSrvInfo()
}

// broadcastSrvInfo 广播服务器信息
func (e *Entity) broadcastSrvInfo() {
	e.srvIDSMux.RLock()
	defer e.srvIDSMux.RUnlock()

	for srvType := range e.srvIDS {
		if srvType != uint8(idata.ServiceClient) && srvType != uint8(e.GetIEntities().GetLocalService().GetSType()) {
			e.AsyncCall(idata.ServiceType(srvType), "SyncEntitySrvInfo")
		}
	}
}

// MsgProcEntitySrvInfoNotify MsgProcEntitySrvInfoNotify
func (e *Entity) MsgProcEntitySrvInfoNotify(imsg.IMsg) {
	e.RefreshSrvIDS()
}
