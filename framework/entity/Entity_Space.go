package entity

// EnterCell EnterCell
func (e *Entity) EnterCell(spaceID uint64) {

	// srvID, err := dbservice.CellUtil(spaceID).GetSrvID()
	// if err != nil {
	// 	//e.Error("Get srvID error ", err)
	// 	return
	// }

	// util := dbservice.GetEntitySrvUtil(e.GetID())
	// _, oldSpaceID, err := util.GetCellInfo()
	// if err != nil {
	// 	e.Error("Get space info failed: ", err)
	// 	return
	// }

	// if oldSpaceID != 0 {
	// 	if oldSpaceID == spaceID {
	// 		e.Warn("oldSpaceID == spaceID")
	// 		return
	// 	}
	// }

	// msg := &msgdef.EnterCellReq{
	// 	SrvID:      srvID,
	// 	CellID:     spaceID,
	// 	EntityType: e.entityType,
	// 	EntityID:   e.entityID,
	// 	DBID:       e.Dbid,
	// 	InitParam:  serializer.Serialize(e.initParam),
	// 	OldSrvID:   0,
	// 	OldCellID:  0,
	// }

	// if err := iserver.GetSrvInst().PostMsgToCell(srvID, spaceID, msg); err != nil {
	// 	e.Error("Enter space failed: ", err)
	// }
}

// LeaveCell 离开场景
func (e *Entity) LeaveCell() {
	// if e.IsCell() {
	// 	e.Warn("Space entity couldn't move into space")
	// 	return
	// }

	// if e.GetCellID() == 0 {
	// 	e.Warn("Entity not in space")
	// 	return
	// }

	// srvID, err := dbservice.CellUtil(e.GetCellID()).GetSrvID()
	// if err != nil {
	// 	e.Error("Get srvID error ", err)
	// 	return
	// }

	// msg := &msgdef.LeaveCellReq{
	// 	EntityID: e.entityID,
	// }

	// if err := iserver.GetSrvInst().PostMsgToCell(srvID, e.GetCellID(), msg); err != nil {
	// 	e.Error("Leave space failed ", err)
	// }
}

// GetCellID
/*func (e *Entity) GetCellID() uint64 {
	return e.cellID
}

// IsOwnerCellEntity 是否拥有SpaceEntity的部分
func (e *Entity) IsOwnerCellEntity() bool {
	return e.cellID != 0
}

// IsCell 是否是空间
func (e *Entity) IsCell() bool {
	return false
}

// IsCellEntity 是否是个空间实体
func (e *Entity) IsCellEntity() bool {
	return false
}
*/
