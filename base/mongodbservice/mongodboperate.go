package dbservice

import (
	"errors"

	log "github.com/cihub/seelog"
	"github.com/globalsign/mgo"
)

//MongoDBQueryAll mongodb查询所有记录
func MongoDBQueryAll(db string, collection string, queryCondition interface{}, retval interface{}) error {

	session := CloneSession()
	if session != nil {
		defer session.Close() //一定要记得释放
		c := session.DB(db).C(collection)

		err := c.Find(queryCondition).All(retval)
		if err != nil {
			log.Error("query all failed", err.Error())
			return err
		}
	} else {
		log.Debug("session nil, not connect or auth")
	}
	return nil
}

//MongoDBQueryOne mongodb查询单个记录
func MongoDBQueryOne(db string, collection string, queryCondition interface{}, retval interface{}) error {

	session := CloneSession()

	if session != nil {
		defer session.Close() //一定要记得释放
		c := session.DB(db).C(collection)

		err := c.Find(queryCondition).One(retval)
		if err != nil {
			log.Error("MongoDBQueryOne failed, err: ", err.Error(), ", condition: ", queryCondition, ", ret: ", retval)
			return err
		}
	} else {
		log.Error("session nil, not connect or auth")
	}
	return nil
}

//MongoDBQuery mongodb查询
func MongoDBQuery(db string, collection string, queryCondition interface{}, retval interface{}) error {

	session := CloneSession()
	if session != nil {
		defer session.Close() //一定要记得释放
		c := session.DB(db).C(collection)

		//var users []User
		//err := c.Find(nil).Limit(5).All(&users)
		err := c.Find(queryCondition).All(&retval)
		if err != nil {
			log.Error("MongoDBQuery failed, err: ", err.Error(), ", condition: ", queryCondition, ", ret: ", retval)
			return err
		}
	} else {
		log.Error("session nil, not connect or auth")
	}
	return nil
}

//MongoDBQueryOneWithSelect mongodb查询单个记录
func MongoDBQueryOneWithSelect(db string, collection string, queryCondition interface{}, qselect interface{}, retval interface{}) error {
	if queryCondition == nil || qselect == nil || retval == nil {
		log.Error("MongoDBQueryOneWithSelect, params is nil")
		return errors.New("param nil")
	}

	session := CloneSession()
	if session != nil {
		defer session.Close() //一定要记得释放
		c := session.DB(db).C(collection)

		err := c.Find(queryCondition).Select(qselect).One(retval)
		if err != nil {
			log.Error("MongoDBQueryOneWithSelect from db: ", db, ", collection: ", collection, " failed, err: ", err.Error(), ", queryCondition: ", queryCondition, ", select: ", qselect, ", ret: ", retval)
			return err
		}

	} else {
		log.Error("session nil, not connect or auth")
	}
	return nil
}

//MongoDBInsert mongodb插入
func MongoDBInsert(db string, collection string, data interface{}) error {
	session := CloneSession()
	if session != nil {
		defer session.Close() //一定要记得释放
		c := session.DB(db).C(collection)
		err := c.Insert(data)

		if err != nil {
			log.Error("insert failed", err.Error())
			return err
		}
	} else {
		log.Error("session nil, not connect or auth")
	}
	return nil
}

//MongoDBUpdate mongodb更新
func MongoDBUpdate(db string, collection string, updateCondition interface{}, data interface{}) error {
	session := CloneSession()
	if session != nil {
		defer session.Close() //一定要记得释放
		c := session.DB(db).C(collection)

		err := c.Update(updateCondition, data)

		if err != nil {
			log.Error("update failed, ", err.Error())
			return err
		}
	} else {
		log.Error("session nil, not connect or auth")
	}
	return nil

}

// MongoDBDelete mongodb删除
func MongoDBDelete(db string, collection string, removeCondition interface{}) error {

	session := CloneSession()
	if session != nil {
		defer session.Close() //一定要记得释放
		c := session.DB(db).C(collection)
		_, err := c.RemoveAll(removeCondition)

		if err != nil {
			log.Error("delete failed", err.Error())
			return err
		}
	} else {
		log.Error("session nil, not connect or auth")
	}
	return nil
}

//MongoDBEnsureIndex mongodb创建索引
func MongoDBEnsureIndex(db string, collection string, idx mgo.Index) error {

	session := CloneSession()
	if session != nil {
		defer session.Close() //一定要记得释放
		c := session.DB(db).C(collection)
		err := c.EnsureIndex(idx)
		if err != nil {
			log.Error("insert failed", err.Error())
			return err
		}
	} else {
		log.Error("session nil, not connect or auth")
	}

	return nil
}

//MongoDBListIndexs mongodb列出索引
func MongoDBListIndexs(db string, collection string, outIndexes *[]mgo.Index) error {

	session := CloneSession()
	if session != nil {
		defer session.Close() //一定要记得释放
		c := session.DB(db).C(collection)
		idxes, err := c.Indexes()
		if err != nil {
			log.Error("list index failed", err.Error())
			return err
		}
		*outIndexes = idxes
	} else {
		log.Error("session nil, not connect or auth")
	}
	return nil
}

//MongoDBDropIndex mongodb删除索引
func MongoDBDropIndex(db string, collection string, idx string) error {
	session := CloneSession()
	if session != nil {
		defer session.Close() //一定要记得释放
		c := session.DB(db).C(collection)
		err := c.DropIndex(idx)

		if err != nil {
			//println("Drop index failed", err.Error())
			log.Error("Drop index failed", err.Error())
			return err
		}
	}
	return nil
}

//MongoDBDropIndexName mongodb删除索引名字
func MongoDBDropIndexName(db string, collection string, idxname string) error {

	session := CloneSession()
	if session != nil {
		defer session.Close() //一定要记得释放
		c := session.DB(db).C(collection)

		err := c.DropIndexName(idxname)

		if err != nil {
			//println("Drop index failed", err.Error())
			log.Error("Drop index name failed", err.Error())
			return err
		}
	}

	return nil
}

//MongoDBFindAndModify 原子操作
func MongoDBFindAndModify(db string, collection string, updateCondition interface{}, data interface{}) error {
	session := CloneSession()
	if session != nil {
		defer session.Close() //一定要记得释放
		c := session.DB(db).C(collection)
		change := mgo.Change{
			Update: data,
		}

		changedInfo, err := c.Find(updateCondition).Apply(change, nil)
		if err != nil {
			log.Error("findAndModify failed", err.Error())
			return err
		}
		if changedInfo == nil {
			log.Error("findAndModify not match condition:", updateCondition)
		} else {
			log.Error("findAndModify changed:", changedInfo)
		}
	} else {
		log.Error("session nil, not connect or auth")
	}
	return nil

}
