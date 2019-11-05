package svrdb

import (
	//"fmt" //这里要改成mongodb.go的文件存放路径
	//"reflect"
	"testing" //mgo "gopkg.in/mgo.v2"

	//dbservice "github.com/GA-TECH-SERVER/zeus/base/mongodbservice"

	//"github.com/cihub/seelog"
	"github.com/globalsign/mgo/bson"
	//"gopkg.in/mgo.v2/bson"
)

type TotalStruct struct {
	DBID  uint64 `bson:"dbid"`
	Name  string
	Part1 PartStruct1
}

type PartStruct1 struct {
	Card  uint64
	Level uint32 `bson:"level"`
}

func TestInsert(t *testing.T) {
	/*	total := NewPlayerDBData()
		total.DBID = 1
		//total.Name = "test"
		dbservice.MongoDBInsert(GameDBName, PlayerTableName, total)*/
}

type TestDBStruct struct {
	ID    bson.ObjectId `bson:"_id"`
	Name  string        `bson:"name"`
	Level uint32        `bson:"level"`
}

func TestQueryPart2(t *testing.T) {
	/*ret := &TestDBStruct{}
	dbservice.MongoDBQueryOneWithSelect(GameDBName, PlayerTableName, DBMap{"dbid": 21}, DBMap{"level": 1, "name": 1}, ret)
	fmt.Print(ret)*/
}

func TestQueryPart3(t *testing.T) {
	/*ret := &TestDBStruct{}
	dbservice.MongoDBQueryOne(GameDBName, PlayerTableName, DBMap{"dbid": 21}, ret)
	fmt.Print(ret)*/
}

type SliceEle struct {
	Itemid   uint32
	Itemname string
}

type MapValue struct {
	Itemid   uint32
	Itemname string
}

type TestDBData struct {
	ID     bson.ObjectId `bson:"_id"`
	DBID   uint64        `bson:"dbid"`
	Name   string        `bson:"name"`
	Level  uint32        `bson:"level"`
	Rating int32         `bson:"rating"`
	Test1  bool
	Test2  int16
	Test3  float32
	Test4  float64
	Test5  uint8
	Test6  byte
	Test7  interface{}
	Test8  interface{}

	TestSlice []SliceEle
	TestMap   map[string]MapValue
}

func TestArrayDB(t *testing.T) {
	/*	data := &TestDBData{}
		data.Id_ = bson.NewObjectId()
		data.DBID = 2
		data.Name = "test2"
		data.Test1 = true
		data.Test2 = 2
		data.Test3 = 3.14
		data.Test4 = 3.14
		data.Test5 = 5
		data.Test6 = '6'
		data.Test7 = "test7"
		data.Test8 = int32(8)

		dbservice.MongoDBInsert(GameDBName, PlayerTableName, data)

		//cardMap := DBMap{"cards": DBMap{"id": 1, "name": "haha"}}
		//dbservice.MongoDBUpdate(GameDBName, PlayerTableName, DBMap{"dbid": 2}, DBMap{"$push": cardMap})

		ret := DBMap{}
		dbservice.MongoDBQueryOne(GameDBName, PlayerTableName, DBMap{"dbid": 2}, ret)

		dbid := ret["dbid"].(int64)
		name := ret["name"].(string)
		rating := ret["rating"].(int)
		level := ret["level"].(int)

		test1 := ret["test1"].(bool)
		test2 := ret["test2"].(int)
		test3 := ret["test3"].(float64)
		test4 := ret["test4"].(float64)
		test5 := ret["test5"].(int)
		test6 := ret["test6"].(int)
		test7 := ret["test7"].(string)
		test8 := ret["test8"].(int)

		notExist := ret["notExist"]
		if notExist == nil {
			fmt.Print("return nil", notExist)
		}

		fmt.Print(dbid, name, rating, level)

		fmt.Print(test1, test2, test3, test4, test5, test6, test7, test8, notExist)

		tt := reflect.TypeOf(dbid).Kind()

		fmt.Print(tt)
		cards := ret["card"]
		aa := reflect.TypeOf(cards).Elem()

		fmt.Print(ret)
		fmt.Print(cards, aa)*/
}

func TestDB(t *testing.T) {
	/*data := TestDBData{}
	data.Id_ = bson.NewObjectId()
	data.DBID = 11
	data.Name = "test"
	data.TestMap = make(map[string]MapValue)

	data.TestMap["first1"] = MapValue{1, "bb first1"}
	data.TestMap["first2"] = MapValue{2, "bb first2"}

	TestMapUpdate := make(map[string]MapValue)
	TestMapUpdate["update1"] = MapValue{1, "aa update1"}
	TestMapUpdate["update2"] = MapValue{2, "aa update2"}

	data.TestSlice = append(data.TestSlice, SliceEle{1, "t1"})

	//dbservice.MongoDBInsert(GameDBName, PlayerTableName, data)

	elems := []bson.RawDocElem{}
	dbmap := make(DBMap)

	raw := &bson.Raw{}

	dbservice.MongoDBQueryOne(GameDBName, PlayerTableName, DBMap{"dbid": 11}, raw)

	bson.Unmarshal(raw.Data, &elems)
	fmt.Print(elems)
	bson.Unmarshal(raw.Data, dbmap)
	fmt.Print(dbmap)
	testMap := make(map[string]MapValue)

	for _, elem := range elems {
		if elem.Name == "testmap" {
			bson.Unmarshal(elem.Value.Data, testMap)
			fmt.Print(testMap)
			break
		}
	}

	// myMap := make(map[string]MapValue)
	// err := bson.Unmarshal(reflect.ValueOf(testMap).Data, myMap)
	// if err == nil {
	// 	return
	// } else {
	// 	fmt.Print(myMap)
	// }

	dbservice.MongoDBUpdate(GameDBName, PlayerTableName, DBMap{"dbid": 11}, DBMap{"$set": DBMap{"testmap": TestMapUpdate}})

	var updateSlice []SliceEle
	updateSlice = append(updateSlice, SliceEle{1, "t1"})
	updateSlice = append(updateSlice, SliceEle{2, "t2"})
	updateSlice = append(updateSlice, SliceEle{3, "t3"})
	dbservice.MongoDBUpdate(GameDBName, PlayerTableName, DBMap{"dbid": 11}, DBMap{"$set": DBMap{"testslice": updateSlice}})*/
}

func TestQueryPart(t *testing.T) {
	/*ret := make(DBMap)
	dbservice.MongoDBQueryOneWithSelect(GameDBName, PlayerTableName, DBMap{"dbid": 1}, DBMap{"dbid": 1, "name": 1}, ret)
	fmt.Print(ret)

	id := ret["dbid"]
	name := ret["name"]
	fmt.Print("id: ", id, ", name: ", name)*/
}

func TestUpdatePart(t *testing.T) {
	/*UpdatePlayer(21, DBMap{"name": "test", "level": 2})*/
}

type SubData struct {
	SubMap map[string]int32 `bson:"SubMap"`
}

type Test struct {
	BaseMap map[string]SubData `bson:"BaseMap"`
}

func TestMapMap(t *testing.T) {
	/*Test := &Test{BaseMap: make(map[string]SubData)}

	subData := SubData{SubMap: make(map[string]int32)}
	subData.SubMap["test1"] = 1
	subData.SubMap["test2"] = 2

	Test.BaseMap["test1"] = subData
	data, _ := bson.Marshal(Test)
	seelog.Debug("data: ", string(data), ", test:", *Test)

	UpdatePlayer(2, DBMap{"test": Test})*/
}
