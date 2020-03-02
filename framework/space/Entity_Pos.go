package space

import (
	"github.com/giant-tech/go-service/base/linmath"
	"github.com/giant-tech/go-service/framework/iserver"
)

// SetPos 设置位置
func (e *Entity) SetPos(pos linmath.Vector3) {
	if e.pos.IsEqual(pos) {
		return
	}

	origPos := e.pos
	e.pos = pos

	e.determineAOIFlag()

	// e.updatePosCoord(pos)
	iPC, ok := e.GetRealPtr().(iserver.IPosChange)
	if ok {
		iPC.OnPosChange(e.pos, origPos)
	}
}

// SetRota 设置旋转
func (e *Entity) SetRota(rota linmath.Vector3) {
	e.rota = rota
}

// GetPos 获取位置
func (e *Entity) GetPos() linmath.Vector3 {
	return e.pos
}

// GetPosPtr 获取位置指针
func (e *Entity) GetPosPtr() *linmath.Vector3 {
	return &e.pos
}

// GetRota 获取旋转
func (e *Entity) GetRota() linmath.Vector3 {
	return e.rota
}
