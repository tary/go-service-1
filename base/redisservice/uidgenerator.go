package dbservice

import (
	"math"
	"sync"

	"github.com/garyburd/redigo/redis"
)

const (
	// GlobalIDStr Player实体ID
	GlobalIDStr = "GlobalID"
	// EntityIDStr 实体唯一ID
	EntityIDStr = "EntityTempID"
)

var (
	// 仅初始化一次
	onceIDGenerator sync.Once
	idGeneratorInst *IDGenerator
)

func initIDGenerator() {
	idGeneratorInst = &IDGenerator{}

	c := GetPersistenceConn()
	defer c.Close()
	// 初始化EntityTempID 初始值
	redis.Uint64(c.Do("HSETNX", "uidgenerator", EntityIDStr, math.MaxUint32+1))
}

// IDGenerator IDGenerator
type IDGenerator struct {
}

// GetIDGenerator ID生成器
func GetIDGenerator() *IDGenerator {
	onceIDGenerator.Do(initIDGenerator)
	return idGeneratorInst
}

// GetGlobalID 获取ID, 保证全局唯一, 从DB库中生成
func (util *IDGenerator) GetGlobalID() (uint64, error) {
	c := GetPersistenceConn()
	defer c.Close()

	return redis.Uint64(c.Do("HINCRBY", "uidgenerator", EntityIDStr, 1))
}

// GetNewPlayerDBID 获取一个新的玩家ID（由于要兼容线上版本，暂时使用原GlobalID）
func (util *IDGenerator) GetNewPlayerDBID() (uint64, error) {

	c := GetPersistenceConn()
	defer c.Close()

	return redis.Uint64(c.Do("HINCRBY", "uidgenerator", GlobalIDStr, 1))
}

// Get 获取指定类型的UID, 保证全局唯一, 从DB库中生成
func (util *IDGenerator) Get(field string) (uint64, error) {
	c := GetPersistenceConn()
	defer c.Close()

	return redis.Uint64(c.Do("HINCRBY", "uidgenerator", field, 1))
}
