package space

import (
	"errors"
)

var errEntityExisted = errors.New("entityExisted")
var errEntityNoExisted = errors.New("entityNoExisted")
var errCreateEntityFail = errors.New("createTinyEntity error")

// AddTinyEntity 增加简单物件
func (s *Space) AddTinyEntity(entityType string, entityID uint64, initParam interface{}) error {

	if _, ok := s.tinyEntities[entityID]; ok {

		return errEntityExisted
	}

	// ie := iserver.GetSrvInst().NewEntityByProtoType(entityType).(ITinyEntity)
	// ie.onEntityCreated(entityID, entityType, s.GetRealPtr().(iserver.ISpace), initParam, ie)

	//ie.onInit()
	//s.tinyEntities[entityID] = ie

	return nil
}

// RemoveTinyEntity 删除简单物件
func (s *Space) RemoveTinyEntity(entityID uint64) error {

	var ie ITinyEntity
	var ok bool
	if ie, ok = s.tinyEntities[entityID]; !ok {

		return errEntityNoExisted
	}

	ie.onDestroy()

	delete(s.tinyEntities, entityID)
	return nil
}

// GetTinyEntity 获取简单实体
func (s *Space) GetTinyEntity(entityID uint64) ITinyEntity {
	var ie ITinyEntity
	var ok bool
	if ie, ok = s.tinyEntities[entityID]; !ok {
		return nil
	}
	return ie
}
