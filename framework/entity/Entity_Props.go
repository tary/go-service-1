package entity

import (
	"fmt"
	"reflect"

	dbservice "github.com/GA-TECH-SERVER/zeus/base/mongodbservice"
	"github.com/GA-TECH-SERVER/zeus/base/stream"

	log "github.com/cihub/seelog"
	"github.com/globalsign/mgo/bson"

	"github.com/GA-TECH-SERVER/zeus/framework/idata"

	"github.com/spf13/viper"
)

// PropDBName PropDBName
var PropDBName string

// PropTableName PropTableName
var PropTableName string

// InitProp 初始化属性列表
func (e *Entity) InitProp(def *Def, loadFromDB bool) {

	if def == nil {
		return
	}

	e.def = def
	for _, p := range e.def.Props {
		for _, sType := range p.Sync {
			if uint32(e.GetIEntities().GetLocalService().GetSType()) == sType {
				e.addProp(p)
				break
			}
		}
	}

	// 读取server.toml里的 mongodb配置
	PropDBName = viper.GetString("MongoDB.GameDBName")
	PropTableName = viper.GetString("MongoDB.GameTableName")

	if loadFromDB {
		e.loadFromDB()
	}
}

// addProp 添加属性
func (e *Entity) addProp(prop *PropDef) {
	e.props[prop.Name] = newPropInfo(prop)
}

// SetProp 设置一个属性的值
func (e *Entity) SetProp(name string, v interface{}) {

	p := e.props[name]

	//log.Debug("SetProp name begin: ", name, " p.value: ", p.value)
	if e.def == nil {
		panic(fmt.Errorf("no def file exist %s", name))
	}

	if !p.def.IsValidValue(v) {
		panic(fmt.Errorf("prop type error %s", name))
	}

	if reflect.DeepEqual(p.value, v) {
		//log.Debug("SetProp DeepEqual return: ", name)
		return
	}

	p.value = v

	//log.Debug("SetProp name end: ", name)
	e.PropDirty(name)
}

// GetProp 获取一个属性的值
func (e *Entity) GetProp(name string) interface{} {
	return e.props[name].value
}

// PropDirty 设置某属性为Dirty
func (e *Entity) PropDirty(name string) {
	p := e.props[name]

	if p.syncFlag == false {
		e.dirtyPropList = append(e.dirtyPropList, p)
	}
	p.syncFlag = true

	if p.def.Persistence {
		if p.dbFlag == false {
			e.dirtySaveProps[name] = p
		}
		p.dbFlag = true
	}
}

// FlushDirtyProp 处理脏属性
func (e *Entity) FlushDirtyProp() {
	e.SyncProps()
	e.SavePropsToDB()
}

func (e *Entity) getDirtyProps(syncType uint32) []*PropInfo {
	m, ok := e.dirtySyncProps[syncType]
	if !ok {
		m = make([]*PropInfo, 0, 10)
		e.dirtySyncProps[syncType] = m
	}
	return m
}

// SendPropsSyncMsg 发送属性同步消息
func (e *Entity) SendPropsSyncMsg(propMap map[uint32][]*PropInfo) error {
	//msg := &msgdef.SyncProps{}
	//msg.EntityID = e.entityID

	for syncType, v := range propMap {
		//msg.PropNum = uint32(len(v))
		//msg.Data = PackPropsToBytes(v)

		servicebase := e.GetIEntities().GetLocalService()
		stype := servicebase.GetSType()
		if syncType != uint32(idata.ServiceClient) && syncType != uint32(stype) {

			log.Debug("SendPropsSyncMsg, syncType: ", syncType)
			e.AsyncCall(idata.ServiceType(syncType), "SyncProps", uint32(len(v)), PackPropsToBytes(v))

			/*s := iserver.GetServiceProxyMgr().GetServiceByID(sid)
			if !s.IsValid() {
				log.Error("AsyncCall service not exist: ", sid, ",  sType: ", sType, " ,methodName: ", methodName)
				return fmt.Errorf("service not exist: %d", sid)
			}



			sid, _, err := e.GetEntitySrvID(uint8(syncType))
			if err != nil {
				log.Error("SendPropsSyncMsg GetEntitySrvID error: ", err, ", syncType: ", syncType)
				return err
			}

			s := iserver.GetServiceProxyMgr().GetServiceByID(sid)

			if !s.IsValid() {
				log.Error("SendPropsSyncMsg service not exist: ", sid, ",  syncType: ", syncType)
				return fmt.Errorf("service not exist: %d", sid)
			}

			if s.IsLocal() {
				//直接发送
				is := iserver.GetLocalServiceMgr().GetLocalService(sid)
				if is == nil {
					log.Error("SendPropsSyncMsg error, service is local, but not found, sid: ", sid)
					return fmt.Errorf("SendPropsSyncMsg error, service is local, but not found")
				}

				return is.PostCallMsg(msg)
			}

			return s.SendMsg(msg)
			*/
		}
	}

	return nil
}

// SyncProps 同步属性给其他人
func (e *Entity) SyncProps() {
	//log.Debug("SyncProps len1: ", len(e.dirtyPropList))
	if len(e.dirtyPropList) == 0 {
		return
	}

	log.Debug("SyncProps len: ", len(e.dirtyPropList))

	for _, p := range e.dirtyPropList {
		for _, s := range p.def.Sync {
			m := e.getDirtyProps(s)
			e.dirtySyncProps[s] = append(m, p)
		}

		p.syncFlag = false
	}

	e.dirtyPropList = e.dirtyPropList[0:0]

	if len(e.dirtySyncProps) > 0 {
		e.SendPropsSyncMsg(e.dirtySyncProps)
	}
}

// UpdateFromMsg 从消息中更新属性
func (e *Entity) UpdateFromMsg(num int, data []byte) {
	bs := stream.NewByteStream(data)
	for i := 0; i < num; i++ {
		name, err := bs.ReadStr()
		if err != nil {
			//e.Error("read prop name fail ", err)
			return
		}

		prop, ok := e.props[name]
		if !ok {
			//e.Error("target entity not own prop ", name)
			return
		}

		err = prop.ReadValueFromStream(bs)
		if err != nil {
			//e.Error("read prop from stream failed ", name, err)
			return
		}
	}
}

// loadFromDB 从数据库中恢复
func (e *Entity) loadFromDB() {
	if e.GetEntityID() == 0 {
		return
	}

	selectProps := bson.M{}
	for k, v := range e.props {
		if v.def.Persistence {
			selectProps[k] = 1
		}
	}

	retRawData := bson.Raw{}
	retM := bson.M{}
	var tempDBElems []bson.RawDocElem

	log.Debug("loadFromDB MongoDBQueryOneWithSelect: ", " , PropDBName: ", PropDBName, " , PropTableName: ", PropTableName, " , selectProps", selectProps)

	dbservice.MongoDBQueryOneWithSelect(PropDBName, PropTableName, bson.M{"dbid": e.GetEntityID()}, selectProps, &retRawData)

	bson.Unmarshal(retRawData.Data, &retM)
	bson.Unmarshal(retRawData.Data, &tempDBElems)

	//log.Info("select props: ", selectProps, "props return: ", ret)

	for k, v := range retM {
		if k == "_id" {
			continue
		}

		info, ok := e.props[k]
		if ok {
			log.Debug("loadFromDB, key is : ", k)
			info.UnPackMongoValue(v, tempDBElems)
		} else {
			log.Error("loadFromDB, prop not exist: ", k)
		}
	}

	/*for _, elem := range tempDBElems {

		if elem.Name == "_id" {
			continue
		}

		info, ok := e.props[elem.Name]
		if ok {
			info.UnPackMongoValue(elem.Value.Data)

		} else {

		}
	}

	*/
}

// SavePropsToDB 属性保存到db
func (e *Entity) SavePropsToDB() {
	if e.GetEntityID() == 0 {
		return
	}

	if len(e.dirtySaveProps) == 0 {
		return
	}

	log.Debug("SyncProps len: ", len(e.dirtySaveProps))

	saveMap := bson.M{}
	for n, p := range e.dirtySaveProps {
		saveMap[n] = p.GetValue()

		p.dbFlag = false
	}

	log.Info("prop saveMap: ", saveMap)

	dbservice.MongoDBUpdate(PropDBName, PropTableName, bson.M{"dbid": e.GetEntityID()}, bson.M{"$set": saveMap})

	e.dirtySaveProps = make(map[string]*PropInfo)
}

// PackProps 打包属性
func (e *Entity) PackProps(syncType uint32) (uint32, []byte) {

	props := make([]*PropInfo, 0, 10)

	for _, prop := range e.props {
		for _, sync := range prop.def.Sync {
			if sync == syncType {
				props = append(props, prop)
				break
			}
		}
	}

	return uint32(len(props)), PackPropsToBytes(props)
}

// PackPropsToBytes 把属列表打包成bytes
func PackPropsToBytes(props []*PropInfo) []byte {
	size := 0

	for _, prop := range props {
		size = size + len(prop.def.Name) + 2
		size = size + prop.GetValueStreamSize()
	}

	if size == 0 {
		return nil
	}

	buf := make([]byte, size)
	bs := stream.NewByteStream(buf)

	for _, prop := range props {
		if err := bs.WriteStr(prop.def.Name); err != nil {
			//e.Error("Pack props failed ", err, prop.def.Name)
		}
		if err := prop.WriteValueToStream(bs); err != nil {
			//e.Error("Pack props failed ", err, prop.def.Name)
		}
	}

	return buf
}
