package space

import (
	"zeus/iserver"
	"zeus/msgdef"
	"zeus/msghandler"

	"github.com/cihub/seelog"
)

type SpaceSessesMsgProc struct {
	srv *SpaceSesses
}

func (proc *SpaceSessesMsgProc) MsgProc_SessVertified(content interface{}) {

	uid := content.(uint64)

	sess := proc.srv.clientSrv.GetSession(uid)
	if sess == nil {
		seelog.Error("couldn't found sess ", uid)
		return
	}

	// Source 为1 代表是 SpaceSess
	sess.Send(&msgdef.ClientVertifySucceedRet{
		Source:   1,
		UID:      uid,
		SourceID: iserver.GetSrvInst().GetSrvID(),
		Type:     0,
	})

	seelog.Debug("space sess establish !! ", content)
}

func (proc *SpaceSessesMsgProc) MsgProc_SpaceUserConnect(content interface{}) {
	msg := content.(*msgdef.SpaceUserConnect)
	sess := proc.srv.clientSrv.GetSession(msg.UID)
	if sess == nil {
		seelog.Error("space user connected but not find sess ", msg.UID)
		return
	}

	space := iserver.GetSrvInst().GetEntity(msg.SpaceID)
	if space == nil {
		seelog.Error("couldn't find space ", msg.SpaceID)
		return
	}

	imh, ok := space.(msghandler.IMsgHandlers)
	if !ok {
		seelog.Error("this not go happen")
		return
	}

	imh.FireMsg("SpaceUserSess", sess)
}
