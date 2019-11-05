package svrdb

import (
	dbservice "github.com/GA-TECH-SERVER/zeus/base/mongodbservice"

	"github.com/globalsign/mgo/bson"
)

//PlayerTableName 玩家的表名
var PlayerTableName = /*string*/ "player"

//PlayerDBData 玩家数据库初始数据
type PlayerDBData struct {
	ID   bson.ObjectId `bson:"_id"`
	DBID uint64        `bson:"dbid"`
}

//NewPlayerDBData 创建带有默认值的数据
func NewPlayerDBData() *PlayerDBData {
	return &PlayerDBData{}
}

//QueryPlayerData 查询全部数据
func QueryPlayerData(dbid uint64, ret interface{}) {
	dbservice.MongoDBQueryOne(GameDBName, PlayerTableName, bson.M{"dbid": dbid}, ret)
}

//QueryPlayerPartData 查询部分数据
func QueryPlayerPartData(dbid uint64, props DBMap, ret interface{}) {
	dbservice.MongoDBQueryOneWithSelect(GameDBName, PlayerTableName, bson.M{"dbid": dbid}, props, ret)
}

//InsertPlayer 插入玩家数据
func InsertPlayer(p *PlayerDBData) {
	if p == nil {
		return
	}

	p.ID = bson.NewObjectId()
	dbservice.MongoDBInsert(GameDBName, PlayerTableName, p)
}

//UpdatePlayer 更新玩家数据
func UpdatePlayer(dbid uint64, updateData DBMap) {
	dbservice.MongoDBUpdate(GameDBName, PlayerTableName, bson.M{"dbid": dbid}, bson.M{"$set": updateData})
}
