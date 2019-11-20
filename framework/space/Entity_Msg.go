package space

import (
	"github.com/giant-tech/go-service/base/linmath"
	"github.com/giant-tech/go-service/framework/entity"
	"github.com/giant-tech/go-service/framework/idata"
	"github.com/giant-tech/go-service/framework/iserver"
	"github.com/giant-tech/go-service/framework/msgdef"
)

// AsyncCallAIOClient 发送消息给所有关注我的客户端
func (e *Entity) AsyncCallAIOClient(methodName string, args ...interface{}) {
	msg := entity.MakeRPC(idata.ServiceClient, 0, methodName, args...)
	e.SendAIOMsg(msg, true)
}

// AsyncCallAOIClientExceptMe  发送给除了自己外的其它人
func (e *Entity) AsyncCallAOIClientExceptMe(methodName string, args ...interface{}) {
	msg := entity.MakeRPC(idata.ServiceClient, 0, methodName, args...)
	e.SendAIOMsg(msg, false)
}

// AsyncCallCenter 遍历center为中心，半径范围radius内的所有实体
func (e *Entity) AsyncCallCenter(center *linmath.Vector3, radius int, methodName string, args ...interface{}) {
	msg := entity.MakeRPC(idata.ServiceClient, 0, methodName, args...)

	e.GetSpace().TravsalCenter(center, radius, func(ia iserver.ICoordEntity) {
		if ise, ok := ia.(iserver.IEntityStateGetter); ok {
			if ise.GetEntityState() != iserver.EntityStateLoop {
				return
			}

			if ie, ok := ia.(iserver.IEntity); ok {
				ie.Send(msg)
			}
		}
	})
}

// AsyncCallCenterExceptMe 遍历center为中心，半径范围radius内的所有实体
func (e *Entity) AsyncCallCenterExceptMe(center *linmath.Vector3, radius int, methodName string, args ...interface{}) {
	msg := entity.MakeRPC(idata.ServiceClient, 0, methodName, args...)

	e.GetSpace().TravsalCenter(center, radius, func(ia iserver.ICoordEntity) {
		if ise, ok := ia.(iserver.IEntityStateGetter); ok {
			if ise.GetEntityState() != iserver.EntityStateLoop {
				return
			}

			if ie, ok := ia.(iserver.IEntity); ok && ie.GetEntityID() != e.GetEntityID() {
				ie.Send(msg)
			}
		}
	})
}

// AsyncCallTowerCenter 遍历center所在的Tower，在该Tower内的center为中心，半径范围radius内的所有实体
func (e *Entity) AsyncCallTowerCenter(center *linmath.Vector3, radius int, methodName string, args ...interface{}) {
	msg := entity.MakeRPC(idata.ServiceClient, 0, methodName, args...)

	e.GetSpace().TravsalTowerCenter(center, radius, func(ia iserver.ICoordEntity) {
		if ise, ok := ia.(iserver.IEntityStateGetter); ok {
			if ise.GetEntityState() != iserver.EntityStateLoop {
				return
			}

			if ie, ok := ia.(iserver.IEntity); ok {
				ie.Send(msg)
			}
		}
	})
}

// AsyncCallTowerCenterExceptMe 遍历center所在的Tower，在该Tower内的center为中心，半径范围radius内的所有实体
func (e *Entity) AsyncCallTowerCenterExceptMe(center *linmath.Vector3, radius int, methodName string, args ...interface{}) {
	msg := entity.MakeRPC(idata.ServiceClient, 0, methodName, args...)

	e.GetSpace().TravsalTowerCenter(center, radius, func(ia iserver.ICoordEntity) {
		if ise, ok := ia.(iserver.IEntityStateGetter); ok {
			if ise.GetEntityState() != iserver.EntityStateLoop {
				return
			}

			if ie, ok := ia.(iserver.IEntity); ok && ie.GetEntityID() != e.GetEntityID() {
				ie.Send(msg)
			}
		}
	})
}

// SendAIOMsg 发送AOI消息
func (e *Entity) SendAIOMsg(msg *msgdef.CallMsg, isCastToMe bool) {
	e.GetSpace().TravsalAOI(e, func(ia iserver.ICoordEntity) {
		if ise, ok := ia.(iserver.IEntityStateGetter); ok {
			if ise.GetEntityState() != iserver.EntityStateLoop {
				return
			}

			if ie, ok := ia.(iserver.IEntity); ok && (e.GetEntityID() != ie.GetEntityID() || isCastToMe) {
				ie.Send(msg)
			}
		}
	})

	e.TravsalExtWatchs(func(o *extWatchEntity) {
		if ise, ok := o.entity.(iserver.IEntityStateGetter); ok {
			if ise.GetEntityState() != iserver.EntityStateLoop {
				return
			}

			if ie, ok := o.entity.(iserver.IEntity); ok {
				ie.Send(msg)
			}
		}
	})
}
