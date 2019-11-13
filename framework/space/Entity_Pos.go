package space

import (
	"github.com/giant-tech/go-service/base/linmath"
	"github.com/giant-tech/go-service/framework/iserver"
)

// SetPos 设置位置
func (e *Entity) SetCoordPos(pos linmath.Vector3) {
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
