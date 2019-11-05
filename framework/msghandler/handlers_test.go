package msghandler

import (
	"encoding/json"
	"fmt"
	"testing"
)

type SubParam struct {
	SubName string `json:"SubName" doc:"子名字"`
	SubID   uint64 `json:"SubID" doc:"子ID"`
}

type Param struct {
	Name          string               `json:"Name" doc:"测试string"`
	IDuint8       uint8                `doc:"测试uint8"`
	IDuint16      uint16               `doc:"测试uint16"`
	IDuint32      uint32               `doc:"测试uint32"`
	IDuint64      uint64               `doc:"测试uint64"`
	IDint8        int8                 `doc:"测试int8"`
	IDint16       int16                `doc:"测试int16"`
	IDint32       int32                `doc:"测试int32"`
	IDint64       int64                `doc:"测试int64"`
	IDfloat32     float32              `doc:"测试float32"`
	IDfloat64     float64              `doc:"测试float64"`
	TestBool      bool                 `doc:"测试bool"`
	TestStruct    SubParam             `doc:"测试struct"`
	TestMap1      map[uint64]string    `doc:"测试map int string"`
	TestMap2      map[string]*SubParam `doc:"测试map string struct"`
	TestSlice1    []SubParam           `doc:"测试slice struct"`
	TestSlice2    []string             `doc:"测试slice string"`
	TestSliceByte []byte               `doc:"测试slice byte"`
}

type MyClass struct {
}

// IDIPGetLevel 获取等级
func (m *MyClass) IDIPGetLevel(p *Param) {

}

func (m *MyClass) IDIPGetName(a int64) {

}

func (m *MyClass) IDIPGetData() {

}

func TestIDIPHandler(t *testing.T) {

	ihandler := NewIDIPHandlers()
	ihandler.RegIDIPMsg(&MyClass{})

	jsonStr, err := ihandler.GetFuncJSON()

	fmt.Print("jsonStr: ", jsonStr, ", err: ", err)
}

func TestStruct(t *testing.T) {
	p := &Param{}
	p.IDfloat64 = 3.5555

	p.TestMap1 = make(map[uint64]string)
	p.TestMap2 = make(map[string]*SubParam)
	p.TestSlice1 = append(p.TestSlice1, SubParam{})
	p.TestSlice1 = append(p.TestSlice1, SubParam{})
	p.TestSlice1 = append(p.TestSlice1, SubParam{})

	p.TestSlice2 = append(p.TestSlice2, "s1")
	p.TestSlice2 = append(p.TestSlice2, "ss2")
	p.TestSlice2 = append(p.TestSlice2, "sss3")

	p.TestSliceByte = []byte("testtesttest")

	p.TestMap1[1] = "test1"
	p.TestMap1[2] = "test2"
	p.TestMap1[3] = "test3"

	p.TestMap2["t1"] = &SubParam{}
	p.TestMap2["t2"] = &SubParam{}

	data, err := json.Marshal(p)

	fmt.Println("err: ", err, ", data: ", string(data))
}

// IDIPReqMsg idip请求消息
type IDIPReqHead struct {
	Cmdid        string
	Seqid        uint32
	ServiceName  string
	SendTime     uint32
	Version      string
	Authenticate string
}

// IDIPReqMsg idip请求消息
type IDIPReqMsg struct {
	Head IDIPReqHead `json:"head"`
	Body string      `json:"body"`
}

func TestIDIPReqMsg(t *testing.T) {
	p := &IDIPReqMsg{}

	data, _ := json.Marshal(p)

	fmt.Println("data: ", string(data))
}
