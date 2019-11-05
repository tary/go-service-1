package iserver

// IEntityGroup 同时具有实体和实体管理的能力
type IEntityGroup interface {
	IEntity
	IEntities
}
