package msgdef

// InitBytesMsg 初始字节流消息
func (def *MsgDef) InitBytesMsg() {
	def.RegMsg(CallMsgID, new(CallMsg))
	def.RegMsg(CallRespMsgID, new(CallRespMsg))

	def.RegMsg(ForwardToClientMsgID, new(ForwardToClientMsg))

	def.RegMsg(EntitySrvInfoNotifyID, new(EntitySrvInfoNotify))
	//def.RegMsg(SessionNotifyID, new(SessionNotify))
	//def.RegMsg(SessionCreateID, new(SessionCreate))
	//def.RegMsg(ClientVertifySucceedMsgID, new(ClientVertifySucceed))
	def.RegMsg(ProtoSyncMsgID, new(ProtoSync))
	def.RegMsg(ClientFrameMsgID, new(ClientFrameMsg))
	def.RegMsg(ServerFrameMsgID, new(ServerFrameMsg))
	def.RegMsg(FramesMsgID, new(FramesMsg))
	def.RegMsg(RequireFramesMsgID, new(RequireFramesMsg))
	def.RegMsg(PropsSyncMsgID, new(PropsSync))
	def.RegMsg(PropsSyncClientMsgID, new(PropsSyncClient))
	def.RegMsg(MRolePropsSyncClientMsgID, new(MRolePropsSyncClient))
	def.RegMsg(EnterCellReqMsgID, new(EnterCellReq))
	def.RegMsg(LeaveCellReqMsgID, new(LeaveCellReq))

	//def.RegMsg(AOIPosChangeMsgID, new(AOIPosChange))

	//def.RegMsg(UserMoveMsgID, new(UserMove))
	//def.RegMsg(EntityPosSetMsgID, new(EntityPosSet))
	def.RegMsg(CellEntityMsgID, new(CellEntityMsg))
	def.RegMsg(ClientTransportMsgID, new(ClientTransport))

	//def.RegMsg(CreateBaseReqMsgID, new(CreateEntityReq))
	//def.RegMsg(CreateBaseRetMsgID, new(CreateEntityRet))
	//def.RegMsg(DestroyBaseReqMsgID, new(DestroyEntityReq))
	//def.RegMsg(DestroyBaseRetMsgID, new(DestroyEntityRet))

	def.RegMsg(CreateCellEntityReqMsgID, new(CreateEntityReq))
	def.RegMsg(CreateCellEntityRetMsgID, new(CreateEntityRet))
	def.RegMsg(DestroyCellEntityReqMsgID, new(DestroyEntityReq))
	def.RegMsg(DestroyCellEntityRetMsgID, new(DestroyEntityRet))

	def.RegMsg(EntityMsgTransportMsgID, new(EntityMsgTransport))
	def.RegMsg(EntityMsgChangeMsgID, new(EntityMsgChange))
	def.RegMsg(SrvMsgTransportMsgID, new(SrvMsgTransport))

	//def.RegMsg(RPCMsgID, new(RPCMsg))
	def.RegMsg(UserDuplicateLoginNotifyMsgID, new(UserDuplicateLoginNotify))

	def.RegMsg(SyncClockMsgID, new(SyncClock))

	def.RegMsg(SpaceUserConnectMsgID, new(SpaceUserConnect))
	def.RegMsg(SpaceUserConnectSucceedRetMsgID, new(SpaceUserConnectSucceedRet))
	def.RegMsg(SyncUserStateMsgID, new(SyncUserState))
	def.RegMsg(AOISyncUserStateMsgID, new(AOISyncUserState))
	def.RegMsg(AdjustUserStateMsgID, new(AdjustUserState))
	def.RegMsg(EntityEventMsgID, new(EntityEvent))

	def.RegMsg(CreateGhostReqMsgID, new(CreateGhostReq))
	def.RegMsg(DeleteGhostReqMsgID, new(DeleteGhostReq))
	def.RegMsg(TransferRealReqMsgID, new(TransferRealReq))
	def.RegMsg(NewRealNotifyMsgID, new(NewRealNotify))

	def.RegMsg(LoginReqMsgID, new(LoginReq))
	def.RegMsg(LoginRespMsgID, new(LoginResp))
	def.RegMsg(PingMsgID, new(Ping))
	def.RegMsg(PongMsgID, new(Pong))

	def.RegMsg(EnterReqMsgID, new(EnterReq))
	def.RegMsg(EnterRespMsgID, new(EnterResp))
	def.RegMsg(ActMsgMsgID, new(ActMsg))
	def.RegMsg(TickDataMsgID, new(TickData))

	def.RegMsg(VideoListReqMsgID, new(VideoListReq))
	def.RegMsg(WatchVideoReqMsgID, new(WatchVideoReq))
	def.RegMsg(CancelWatchReqMsgID, new(CancelWatchReq))

	def.RegMsg(CreateEntityNotifyMsgID, new(CreateEntityNotify))

	def.RegMsg(ClientVerifyReqMsgID, new(ClientVerifyReq))
	def.RegMsg(ClientVerifyRespMsgID, new(ClientVerifyResp))

	// def.RegMsg(StartMatchReqMsgID, new(StartMatchReq))
	// def.RegMsg(StartMatchRespMsgID, new(StartMatchResp))
	// def.RegMsg(StopMatchReqMsgID, new(StopMatchReq))
	// def.RegMsg(StopMatchRespMsgID, new(StopMatchResp))
	// def.RegMsg(CreateMatchRoomNotifyMsgID, new(CreateMatchRoomNotify))
	// def.RegMsg(LeaveMatchRoomNotifyMsgID, new(LeaveMatchRoomNotify))
	// def.RegMsg(JoinMatchRoomNotifyMsgID, new(JoinMatchRoomNotify))
	// def.RegMsg(MatchFinishNotifyMsgID, new(MatchFinishNotify))
}
