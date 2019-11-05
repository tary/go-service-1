package iserver

// IEntityPrototype 实体原型
type IEntityPrototype interface {
	RegProtoType(name string, protoType IEntity, autoCreate bool)
	NewEntityByProtoType(entityType string) interface{}
}
