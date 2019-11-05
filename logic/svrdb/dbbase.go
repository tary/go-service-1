package svrdb

import (
	"github.com/GA-TECH-SERVER/zeus/framework/entity"

	log "github.com/cihub/seelog"
	"github.com/globalsign/mgo/bson"
)

// DBMap dbmap
type DBMap bson.M

// GameDBName game db name
var GameDBName = "game"

// Init svrdb初始化
func Init(dbName string) {
	if dbName != "" {
		GameDBName = dbName
	}

	entity.PropDBName = GameDBName
	entity.PropTableName = PlayerTableName
}

// GetIntValue 获得整形值
func GetIntValue(dbmap DBMap, key string, defaultValue int64) int64 {
	val := dbmap[key]
	if val == nil {
		return defaultValue
	}

	int64Val, ok := val.(int64)
	if ok {
		return int64Val
	}

	intVal, ok := val.(int)
	if ok {
		return int64(intVal)
	}
	//不是整型数据
	log.Error("GetDBValue, key: ", key, " is not int value, ", val)

	return defaultValue
}

// GetStringValue 获得string值
func GetStringValue(dbmap DBMap, key string, defaultValue string) string {
	val := dbmap[key]
	if val == nil {
		return defaultValue
	}

	strVal, ok := val.(string)
	if ok {
		return strVal
	}
	//不是整型数据
	log.Error("GetDBValue, key: ", key, " is not string value, ", val)

	return defaultValue
}

// GetFloatValue 获得float值
func GetFloatValue(dbmap DBMap, key string, defaultValue float64) float64 {
	val := dbmap[key]
	if val == nil {
		return defaultValue
	}

	floatVal, ok := val.(float64)
	if ok {
		return floatVal
	}
	//不是float数据
	log.Error("GetDBValue, key: ", key, " is not float value, ", val)

	return defaultValue
}

// GetBoolValue 获得bool值
func GetBoolValue(dbmap DBMap, key string, defaultValue bool) bool {
	val := dbmap[key]
	if val == nil {
		return defaultValue
	}

	boolVal, ok := val.(bool)
	if ok {
		return boolVal
	}
	//不是bool数据
	log.Error("GetDBValue, key: ", key, " is not bool value, ", val)

	return defaultValue
}
