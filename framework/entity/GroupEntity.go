package entity

// GroupEntity 本身为一个实体，且能管理多个实体（如队伍，房间等）
type GroupEntity struct {
	Entity
	*Entities
}

// OnEntityInit 初始化
func (g *GroupEntity) OnEntityInit() error {
	g.Entity.OnEntityInit()
	g.Entities = NewEntities(false, g.Entity.GetIEntities().GetLocalService())
	return nil
}

// OnEntityLoop loop
func (g *GroupEntity) OnEntityLoop() {
	g.Entity.OnEntityLoop()
	g.Entities.Loop()
}

// OnEntityDestroy group destroy
func (g *GroupEntity) OnEntityDestroy() {
	g.Entity.OnEntityDestroy()
	//销毁所有成员
	g.Destroy()
}
