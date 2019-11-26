package space

import (
	"github.com/giant-tech/go-service/framework/iserver"
)

type iSpaceCtrl interface {
	updateCoord(iserver.ICellEntity)
}

func (s *Space) refreshEntityState() {

	s.TravsalEntity("Player", func(entity iserver.IEntity) {
		iw, ok := entity.(IWatcher)
		if ok {
			iw.reflushStateChangeMsg()
		}
	})
}
