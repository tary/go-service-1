package dbservice

import (

	//这里要改成mongodb.go的文件存放路径

	"fmt"
	"log"
	"testing"
	"time"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type User struct {
	ID        bson.ObjectId `bson:"_id"`
	Name      string        `bson:"name"`
	Age       int           `bson:"age"`
	JoinedAt  time.Time     `bson:"joned_at"`
	Interests []string      `bson:"interests"`
}

type Person struct {
	Name  string
	Phone string
}

func TestMongodbInsert(t *testing.T) {

	//configPath := "../../../res/config/server.toml"
	configPath := "./mongodb_test.toml"

	// 设置配置文件, 放在第一个测试用例里, 要保证能连上mongodb,否则这些测试用例会报错。
	setConfig(configPath)

	user := &User{ID: bson.NewObjectId(),
		Name:      "yekoufeng",
		Age:       33,
		JoinedAt:  time.Now(),
		Interests: []string{"Develop", "Movie"}}

	MongoDBInsert("game", "people", user)

	user.ID = bson.NewObjectId()
	user.Name = "yekoufeng2"
	user.Age = 32
	MongoDBInsert("game", "people", user)

	user.ID = bson.NewObjectId()
	user.Name = "yekoufeng3"
	user.Age = 31
	MongoDBInsert("game", "people", user)

}

/*单条件查询

=($eq)
c.Find(bson.M{"name": "Jimmy Kuu"}).All(&users)
!=($ne)
c.Find(bson.M{"name": bson.M{"$ne": "Jimmy Kuu"}}).All(&users)
>($gt)
c.Find(bson.M{"age": bson.M{"$gt": 32}}).All(&users)
<($lt)
c.Find(bson.M{"age": bson.M{"$lt": 32}}).All(&users)
>=($gte)
c.Find(bson.M{"age": bson.M{"$gte": 33}}).All(&users)
<=($lte)
c.Find(bson.M{"age": bson.M{"$lte": 31}}).All(&users)
in($in)
c.Find(bson.M{"name": bson.M{"$in": []string{"Jimmy Kuu", "Tracy Yu"}}}).All(&users)
多条件查询

and($and)
c.Find(bson.M{"name": "Jimmy Kuu", "age": 33}).All(&users)
or($or)
c.Find(bson.M{"$or": []bson.M{bson.M{"name": "Jimmy Kuu"}, bson.M{"age": 31}}}).All(&users)

*/

func TestMongodbQueryAll(t *testing.T) {

	var users []User
	MongoDBQueryAll("game", "people", bson.M{"name": "yekoufeng", "age": 33}, &users)
	fmt.Println(users)
}

func TestMongodbQueryOne(t *testing.T) {

	var user User
	MongoDBQueryOne("game", "people", bson.M{"name": "yekoufeng", "age": 33}, &user)
	fmt.Println(user)

}
func TestMongodbInsertAndQueryAll(t *testing.T) {
	session := CloneSession() //调用这个获得session
	defer session.Close()     //一定要记得释放

	c := session.DB("game").C("people")
	err := c.Insert(&User{ID: bson.NewObjectId(),
		Name:      "yekoufeng",
		Age:       33,
		JoinedAt:  time.Now(),
		Interests: []string{"Develop", "Movie"}})

	if err != nil {
		panic(err)
	}

	var users []User
	err = c.Find(nil).Limit(5).All(&users)
	if err != nil {
		panic(err)
	}
	fmt.Println(users)

}

func TestMongodbInsertAndQueryOne(t *testing.T) {
	session := CloneSession() //调用这个获得session
	defer session.Close()     //一定要记得释放

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true) //读模式，与副本集有关，详情参考https://docs.mongodb.com/manual/reference/read-preference/ & https://docs.mongodb.com/manual/replication/

	c := session.DB("game").C("people")

	val1 := &Person{"xxx", "123456"}
	val2 := &Person{"ykf", "13621730307"}

	err := c.Insert(val1,
		val2)
	if err != nil {
		log.Fatal(err)
	}

	result := Person{}
	err = c.Find(bson.M{"name": "yekoufeng"}).One(&result) //如果查询失败，返回“not found”
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Phone:", result.Phone)
}

func TestMongodbUpdate(t *testing.T) {
	session := CloneSession() //调用这个获得session
	defer session.Close()     //一定要记得释放

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true) //读模式，与副本集有关，详情参考https://docs.mongodb.com/manual/reference/read-preference/ & https://docs.mongodb.com/manual/replication/

	c := session.DB("game").C("people")

	/*c.Update(bson.M{"_id": bson.ObjectIdHex("5bc0683142a97042e0fd4fb9")},
	bson.M{"$set": bson.M{
		"name": "yekoufeng",
		"age":  84,
	}})
	*/

	MongoDBUpdate("game", "people", bson.M{"_id": bson.ObjectIdHex("5bc0683142a97042e0fd4fb9")}, bson.M{"$set": bson.M{
		"name": "yekoufeng",
		"age":  94,
	}})

	var users []User
	err := c.Find(bson.M{"name": "yekoufng", "age": 94}).All(&users)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Phone:", users)
}

func TestMongodbDelete(t *testing.T) {
	//session, err := mgo.Dial("server1.example.com,server2.example.com") //传入数据库的地址，可以传入多个，具体请看接口文档
	session := CloneSession() //调用这个获得session
	defer session.Close()     //一定要记得释放

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true) //读模式，与副本集有关，详情参考https://docs.mongodb.com/manual/reference/read-preference/ & https://docs.mongodb.com/manual/replication/

	MongoDBDelete("game", "people", bson.M{"age": 94})

	/*
			c := session.DB("game").C("people")

			_, err := c.RemoveAll(bson.M{"name": "Ale"})

			if err != nil {
				log.Fatal(err)
			}


		var users []User
		err = c.Find(bson.M{"name": "Jimmy Kuu", "age": 33}).All(&users)

		if err != nil {
			log.Fatal(err)
		}

	fmt.Println("Phone:", users)*/
}

func TestEnsureIndex(t *testing.T) {
	session := CloneSession() //调用这个获得session
	defer session.Close()     //一定要记得释放

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	//hash index
	/*idx := mgo.Index{
		Key: []string{"$hashed:_id"},
	}*/

	//complex index
	/*idx := mgo.Index{
		Key: []string{"name", "age"},
	}*/
	//unique index
	idx := mgo.Index{
		Key:    []string{"name"},
		Unique: true,
	}
	MongoDBEnsureIndex("game", "people", idx)

}

func TestDropIndex(t *testing.T) {

	session := CloneSession() //调用这个获得session
	defer session.Close()     //一定要记得释放

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	idx := "name"
	MongoDBDropIndex("game", "people", idx)
}

func TestDropIndexName(t *testing.T) {

	session := CloneSession() //调用这个获得session
	defer session.Close()     //一定要记得释放

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	//删除组合索引
	idx := "name_1_age_1"
	MongoDBDropIndexName("game", "people", idx)
}
