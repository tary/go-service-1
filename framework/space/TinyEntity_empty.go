package space

import (
	"github.com/giant-tech/go-service/base/net/inet"
	logicredis "github.com/giant-tech/go-service/framework/logicredis"
)

func (e *TinyEntity) GetDBID() uint64 {
	return 0
}

func (e *TinyEntity) GetCellID() uint64 {
	if e.space == nil {
		return 0
	}

	return e.space.GetEntityID()
}

func (e *TinyEntity) GetSrvIDS() map[uint8]*logicredis.EntitySrvInfo {
	return nil
}

func (e *TinyEntity) IsOwnerSpaceEntity() bool {
	return true
}

func (e *TinyEntity) IsSpaceEntity() bool {
	return true
}

func (e *TinyEntity) Post(srvType uint8, msg inet.IMsg) error {
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

func (e *TinyEntity) GetBaseProps() []byte {
	return nil
}
