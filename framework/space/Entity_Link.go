package space

import (
	"errors"
	"zeus/iserver"
	"zeus/linmath"
)

// 实体链接相关代码

// TODO
// 链接状态下的时候, 状态快照中的位置和AOI坐标系中的位置不一致

var errLinkFailedNil = errors.New("Param Nil")

// LinkToTarget 链接到目标Entity, 链接成功后, 本身不再更新位置信息, 本身的位置信息被目标Entity信息替代
func (e *Entity) LinkToTarget(target iserver.ICoordEntity) {
	if target == nil {
		return
	}

	if e.linkTarget != nil {
		e.Unlink()
	}

	e.SetPos(target.GetPos())
	e.SetRota(target.GetRota())

	target.LinkedBy(e)
	e.linkTarget = target

	e.setBasePropsDirty(true)
}

// Unlink 取消链接
func (e *Entity) Unlink() {
	if e.linkTarget == nil {
		return
	}

	target := e.linkTarget

	e.linkTarget.UnlinkedBy(e)
	e.linkTarget = nil

	e.SetPos(target.GetPos())
	e.SetRota(target.GetRota())

	e.setBasePropsDirty(true)
}

// LinkedBy 被其他Entity链接
func (e *Entity) LinkedBy(o iserver.ICoordEntity) {
	if o == nil {
		return
	}

	if e.linkerList == nil {
		e.linkerList = make(map[uint64]iserver.ICoordEntity)
	}

	e.linkerList[o.GetID()] = o
	e.setBasePropsDirty(true)
}

// UnlinkedBy 被o取消链接
func (e *Entity) UnlinkedBy(o iserver.ICoordEntity) {
	if o == nil || e.linkerList == nil {
		return
	}

	delete(e.linkerList, o.GetID())
	e.setBasePropsDirty(true)
}

// GetLinkedEntity 获取链接到自己身上的Entity列表
func (e *Entity) GetLinkedEntity() map[uint64]iserver.ICoordEntity {
	return e.linkerList
}

// IsLinked 是否链接到其他Entity
func (e *Entity) IsLinked() bool {
	return e.linkTarget != nil
}

// IsBeLinked 是否被其他Entity链接
func (e *Entity) IsBeLinked() bool {
	return len(e.linkerList) != 0
}

func (e *Entity) updateLinkerPos(pos linmath.Vector3) {
	for _, o := range e.linkerList {
		o.SetCoordPos(pos)
	}
}
