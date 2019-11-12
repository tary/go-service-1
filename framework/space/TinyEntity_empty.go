package space

import (
	"zeus/dbservice"
	"zeus/iserver"
	"zeus/linmath"
	"zeus/msgdef"
)

func (e *TinyEntity) GetDBID() uint64 {
	return 0
}

func (e *TinyEntity) GetCellID() uint64 {
	if e.space == nil {
		return 0
	}

	return e.space.GetID()
}

func (e *TinyEntity) GetSrvIDS() map[uint8]*dbservice.EntitySrvInfo {
	return nil
}

func (e *TinyEntity) IsOwnerSpaceEntity() bool {
	return true
}

func (e *TinyEntity) IsSpaceEntity() bool {
	return true
}

func (e *TinyEntity) Post(srvType uint8, msg msgdef.IMsg) error {
	return nil
}

func (e *TinyEntity) RPC(srvType uint8, methodName string, args ...interface{}) error {
	return nil
}

/*func (e *TinyEntity) RPCOther(srvType uint8, srcEntityID uint64, methodName string, args ...interface{}) error {
	return nil
}*/

func (e *TinyEntity) EnterSpace(spaceID uint64) {

}

func (e *TinyEntity) LeaveSpace() {

}

func (e *TinyEntity) GetProxy() iserver.IEntityProxy {
	return nil
}

func (e *TinyEntity) LinkToTarget(t iserver.ICoordEntity) {}

func (e *TinyEntity) Unlink() {}

func (e *TinyEntity) LinkedBy(o iserver.ICoordEntity) {}

func (e *TinyEntity) UnlinkedBy(o iserver.ICoordEntity) {}

func (e *TinyEntity) GetLinkedEntity() map[uint64]iserver.ICoordEntity {
	return nil
}

func (e *TinyEntity) IsLinked() bool {
	return false
}

func (e *TinyEntity) IsBeLinked() bool {
	return false
}

func (e *TinyEntity) SetCoordPos(pos linmath.Vector3) {}

func (e *TinyEntity) GetBaseProps() []byte {
	return nil
}
