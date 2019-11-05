package entity

// GroupEntity 本身为一个实体，且能管理多个实体（如队伍，房间等）
type GroupEntity struct {
	Entity
	*Entities
}

// OnGroupInit 初始化
func (g *GroupEntity) OnGroupInit() error {
	g.Entities = NewEntities(false, g.Entity.GetIEntities().GetLocalService())
	return nil
}

// OnGroupLoop loop
func (g *GroupEntity) OnGroupLoop() {
	g.Loop()
}

// OnGroupDestroy group destroy
func (g *GroupEntity) OnGroupDestroy() {
	//销毁所有成员
	g.Destroy()
}
