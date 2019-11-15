package iserver

import "github.com/giant-tech/go-service/base/linmath"

// ISpace space接口
type ISpace interface {
	IEntityGroup
	ICoord

	GetTimeStamp() uint32

	AddEntity(entityType string, entityID uint64, initParam interface{}, syncInit bool) (IEntity, error)
	RemoveEntity(entityID uint64) error

	AddTinyEntity(entityType string, entityID uint64, initParam interface{}) error
	RemoveTinyEntity(entityID uint64) error

	IsMapLoaded() bool
	FindPath(srcPos, destPos linmath.Vector3) ([]linmath.Vector3, error)
	Raycast(origin, direction linmath.Vector3, length float32, mask int32) (float32, linmath.Vector3, int32, bool, error)
	CapsuleRaycast(head, foot linmath.Vector3, radius float32, origin, direction linmath.Vector3, length float32) (float32, bool, error)
	SphereRaycast(center linmath.Vector3, radius float32, origin, direction linmath.Vector3, length float32) (float32, bool, error)
	GetHeight(x, z float32) (float32, error)
	IsWater(x, z float32, waterlevel float32) (bool, error)
}
