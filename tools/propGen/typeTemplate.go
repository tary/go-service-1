package main

import (
	"fmt"
)

// TypeTemplate 类型结构代码生成模版
type TypeTemplate struct {
	name         string
	typeInfo     map[string]interface{}
	interfaceStr string
}

var tt *TypeTemplate

// NewTypeTemplate 生成新的模版工具
func NewTypeTemplate(name string) *TypeTemplate {
	tt = &TypeTemplate{}
	tt.name = name
	tt.typeInfo = make(map[string]interface{})
	return tt
}

// GetTypeTemplate 获得类型模板
func GetTypeTemplate() *TypeTemplate {

	return tt
}

// AddType 增加类型
func (t *TypeTemplate) AddType(typename string, typinfo interface{}) {
	t.typeInfo[typename] = typinfo
}

// IsExistTypeTemplate 判断类型是否存在
func (t *TypeTemplate) IsExistTypeTemplate(typname string) bool {

	_, isIn := t.typeInfo[typname]

	return isIn
}

// genhead 产生文件头
func (t *TypeTemplate) genhead() string {
	var baseStr string

	baseStr += "package entitydef\n"
	baseStr += "\n"

	return baseStr
}

// genString 产生字符串
func (t *TypeTemplate) genString() string {
	str := t.genhead()
	//t.genint32Array()
	//t.genUint32Array()
	//t.genint64Array()
	//t.genUint64Array()

	for key, val := range t.typeInfo {

		valMap := val.(map[string]interface{})
		if valMap["type"] == "struct" {
			str += t.genStruct(key, valMap)
		} else if valMap["type"] == "map" {
			str += t.genMap(key, valMap)
		}
	}
	return str
}

/*type FriendsInfo struct {
	MyFriendsDbid    []uint64 `bson:"MyFriendsDbid"`
	ApplyFriendsDbid []uint64 `bson:"ApplyFriendsDbid"`
}

type HEROINFO struct {
	HeroName string `bson:"HeroName"`
	HeroID   int32  `bson:"HeroID"`
}

type HEROS = map[string]HEROINFO
*/

// genStruct 产生struct
func (t *TypeTemplate) genStruct(name string, valMap map[string]interface{}) string {

	var baseStr string

	baseStr += fmt.Sprintf("type %s struct {\n", name)

	//valMap := val.(map[string]interface{})

	// name为每个结构体的字段
	for name, val := range valMap["fields"].(map[string]interface{}) {

		typMap := val.(map[string]interface{})

		tval := typMap["type"].(string)
		//struct里是否有map成员变量的需求？
		/*	if typMap["type"] == "arrayuint64" || typMap["type"] == "arrayuint32" || typMap["type"] == "arrayint64" || typMap["type"] == "arrayint32" {
				baseStr += fmt.Sprintf("%s	[]%s `bson:\"%s\"` \n", name, typMap["type"][5:], name)
			} else {
				baseStr += fmt.Sprintf("%s	%s `bson:\"%s\"` \n", name, typMap["type"], name)
			}
		*/

		if tval == "arrayuint64" || tval == "arrayuint32" || tval == "arrayint64" || tval == "arrayint32" {
			baseStr += fmt.Sprintf("	%s	[]%s `bson:\"%s\"` \n", name, tval[5:], name)
		} else {
			baseStr += fmt.Sprintf("	%s	%s `bson:\"%s\"` \n", name, tval, name)
		}
	}

	baseStr += "}\n"

	return baseStr
}

// genMap 产生map
func (t *TypeTemplate) genMap(name string, valMap map[string]interface{}) string {

	var baseStr string

	baseStr += fmt.Sprintf("type %s = map[", name)

	//valMap := val.(map[string]interface{})

	// name为每个结构体的字段
	for name, val := range valMap["fields"].(map[string]interface{}) {

		typMap := val.(map[string]interface{})
		tval := typMap["type"].(string)

		if name == "index" {
			baseStr += fmt.Sprintf("%s]", tval)
		} else if name == "value" {
			baseStr += fmt.Sprintf("%s\n", tval)
		}
	}

	baseStr += "\n"
	return baseStr
}
