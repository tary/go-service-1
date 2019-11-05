package entity

import log "github.com/cihub/seelog"
import "github.com/giant-tech/go-service/base/stream"

// RPCSyncEntitySrvInfo 同步实体的服务信息
func (e *Entity) RPCSyncEntitySrvInfo() {
	log.Debug("RPCSyncSrvInfo, ID: ", e.GetEntityID())

	e.RefreshSrvIDS()
}

// RPCSyncProps sync props
func (e *Entity) RPCSyncProps(num uint32, data []byte) {
	log.Debug("RPCSyncProps, ID: ", e.GetEntityID())

	bs := stream.NewByteStream(data)
	for i := uint32(0); i < num; i++ {
		name, err := bs.ReadStr()
		if err != nil {
			e.Error("read prop name fail ", err)
			return
		}

		prop, ok := e.props[name]
		if !ok {
			e.Error("target entity not own prop ", name)
			return
		}

		err = prop.ReadValueFromStream(bs)
		if err != nil {
			e.Error("read prop from stream failed ", name, err)
			return
		}
	}
}
