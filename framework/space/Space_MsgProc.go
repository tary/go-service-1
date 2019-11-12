package space

import (
	"zeus/iserver"
	"zeus/msgdef"
	"zeus/serializer"
)

// SpaceMsgProc 消息处理函数
type SpaceMsgProc struct {
	space *Space
}

func (proc *SpaceMsgProc) MsgProc_EnterSpaceReq(content msgdef.IMsg) {
	msg := content.(*msgdef.EnterSpaceReq)

	params := serializer.UnSerialize(msg.InitParam)
	if len(params) < 1 {
		proc.space.Error("Unmarshal initparam error ", msg.InitParam)
		return
	}

	err := proc.space.AddEntity(msg.EntityType, msg.EntityID, msg.DBID, params[0], false, false)
	if err != nil {
		proc.space.Error("Add entity error ", err, msg)
		return
	}
}

func (proc *SpaceMsgProc) MsgProc_LeaveSpaceReq(content msgdef.IMsg) {
	msg := content.(*msgdef.LeaveSpaceReq)

	if err := proc.space.RemoveEntity(msg.EntityID); err != nil {
		proc.space.Error("Remove entity error ", err, msg)
	}
}

func (proc *SpaceMsgProc) MsgProc_SpaceUserSess(sess iserver.ISess) {
	ie := proc.space.GetEntityByDBID("Player", sess.GetID())
	if ie == nil {
		proc.space.Error("there is no player in space ui = ", sess.GetID())
		sess.Close()
		return
	}

	ise, ok := ie.(iserver.ISpaceEntity)
	if !ok {
		proc.space.Error("conert to ispaceentity error , strange!! ")
		return
	}

	ise.SetClient(sess)
	if err := ise.Post(iserver.ServerTypeClient, &msgdef.SpaceUserConnectSucceedRet{}); err != nil {
		proc.space.Error("Send SpaceUserConnectSucceedRet failed ", err)
	}

	sess.FlushBacklog()

	// if iP, ok := ie.(iPropsSender); ok {
	// 	if err := iP.SendFullProps(); err != nil {
	// 		log.Error(err)
	// 	}
	// } else {
	// 	log.Warn("Cant send props")
	// }

	// if iS, ok := ie.(iStateSender); ok {
	// 	iS.SendFullState()
	// } else {
	// 	log.Warn("Cant send state")
	// }

	// if iA, ok := ie.(iAOISender); ok {
	// 	if err := iA.SendFullAOIs(); err != nil {
	// 		log.Error(err)
	// 	}
	// } else {
	// 	log.Warn("Cant send aois")
	// }

	// if iG, ok := ie.(iGameStateSender); ok {
	// 	iG.SendFullGameState()
	// }
}
