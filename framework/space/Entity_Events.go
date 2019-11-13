package space

import (
	"github.com/giant-tech/go-service/framework/msgdef"

	"github.com/giant-tech/go-service/base/serializer"
)

// BroadcastEvent 广播事件
func (e *Entity) BroadcastEvent(event string, args ...interface{}) {
	msg := &msgdef.EntityEvent{}
	msg.SrcEntityID = e.GetEntityID()
	msg.EventName = event
	msg.Data = serializer.Serialize(args...)

	e.CastMsgToAllClient(msg)
}

// BroadcastEventExceptMe 广播事件
func (e *Entity) BroadcastEventExceptMe(event string, args ...interface{}) {
	msg := &msgdef.EntityEvent{}
	msg.SrcEntityID = e.GetEntityID()
	msg.EventName = event
	msg.Data = serializer.Serialize(args...)

	e.CastMsgToAllClientExceptMe(msg)
}
