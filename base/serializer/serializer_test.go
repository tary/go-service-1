package serializer

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestSerializer(t *testing.T) {
	data := Serialize(float32(3.14))
	fmt.Printf("%x\n", data)
	fmt.Println(data)

	ret := UnSerialize(data)
	fmt.Println(ret)
	for _, v := range ret {
		fmt.Println(reflect.TypeOf(v), v)
	}
}

type TestStruct struct {
	TestString  string
	TestInt8    int8
	TestInt16   int16
	TestInt32   int32
	TestInt64   int64
	TestUint8   uint8
	TestUint16  uint16
	TestUint32  uint32
	TestUint64  uint64
	TestFloat32 float32
	TestFloat64 float64
	//TestProto   spb.CreateRoomReq

	PTestString  *string
	PTestInt8    *int8
	PTestInt16   *int16
	PTestInt32   *int32
	PTestInt64   *int64
	PTestUint8   *uint8
	PTestUint16  *uint16
	PTestUint32  *uint32
	PTestUint64  *uint64
	PTestFloat32 *float32
	PTestFloat64 *float64
	//PTestProto   *spb.CreateRoomReq

	TestMap   map[string]string
	TestSlice []string

	PTestMap   *map[string]string
	PTestSlice *[]string

	TestStru  TestSubStruct
	PTestStru *TestSubStruct

	TestPointMap map[string]*TestSubStruct
}

type TestSubStruct struct {
	TestString  string
	TestInt8    int8
	TestInt16   int16
	TestInt32   int32
	TestInt64   int64
	TestUint8   uint8
	TestUint16  uint16
	TestUint32  uint32
	TestUint64  uint64
	TestFloat32 float32
	TestFloat64 float64

	PTestString  *string
	PTestInt8    *int8
	PTestInt16   *int16
	PTestInt32   *int32
	PTestInt64   *int64
	PTestUint8   *uint8
	PTestUint16  *uint16
	PTestUint32  *uint32
	PTestUint64  *uint64
	PTestFloat32 *float32
	PTestFloat64 *float64

	TestSubMap map[string]string
	TestSlice  []string

	PTestSubMap *map[string]string
	PTestSlice  *[]string
}

type TestMapMethod struct {
}

func (*TestMapMethod) Say(testMap *TestMap) {
	fmt.Println("testMap: ", *testMap)
	fmt.Println("testSubMap of testMap: ", testMap.TestStru)
}

type TestMap struct {
	TestArray [7]string
	//TestMap   map[string]string
	//PTestMap  *map[string]string
	TestStru TestSubMap
	//PTestStru *TestSubMap
}

type TestSubMap struct {
	TestSubMap  map[string]string
	PTestSubMap *map[string]string
}

type Test struct {
}

// func (*Test) Say(testProto *spb.CreateRoomReq, testI uint16, testF float32, testStr *string, testStruct *TestStruct) {
// 	fmt.Println("Say: ", testI, testF, *testStr, testStruct)

// 	val := testStruct.TestPointMap["testKey"]
// 	fmt.Println("Say: ", *val)
// }

func call(args ...interface{}) []byte {
	return SerializeNew(args...)
}

func TestMyMap(te *testing.T) {
	subStruct := TestSubMap{}
	subStruct.TestSubMap = make(map[string]string)
	subStruct.TestSubMap["testSubKey1"] = "testSubValue1"
	subStruct.TestSubMap["testSubKey2"] = "testSubValue2"
	subStruct.PTestSubMap = &subStruct.TestSubMap

	testStruct := TestMap{}
	//testStruct.TestMap = make(map[string]string)
	//testStruct.PTestMap = &testStruct.TestMap
	//testStruct.TestMap["testKey"] = "testValue"
	testStruct.TestStru = subStruct
	//testStruct.PTestStru = &subStruct

	testStruct.TestArray[0] = "00000"
	testStruct.TestArray[1] = "11111"
	testStruct.TestArray[2] = "22222"
	testStruct.TestArray[3] = "33333"
	testStruct.TestArray[4] = "44444"
	testStruct.TestArray[5] = "55555"
	testStruct.TestArray[6] = "66666"

	testMap := make(map[string]string)
	testMap["testKey"] = "testValue"
	var TestSlice []string
	TestSlice = append(TestSlice, "testSlice")

	data := call(&testStruct)

	test := &TestMapMethod{}
	v := reflect.ValueOf(test)
	t := reflect.TypeOf(test)

	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		methodValue := v.MethodByName(method.Name)
		args, err := UnSerializeByFunc(methodValue, data)
		if err == nil {
			methodValue.Call(args)
		}
	}
}

func TestMechodAll(te *testing.T) {
	subStr := "subTestString"
	str := "TestString"

	subStruct := TestSubStruct{}
	subStruct.TestString = subStr
	subStruct.TestInt8 = int8(1)
	subStruct.TestInt16 = int16(1)
	subStruct.TestInt32 = int32(1)
	subStruct.TestInt64 = int64(1)
	subStruct.TestUint8 = uint8(1)
	subStruct.TestUint16 = uint16(1)
	subStruct.TestUint32 = uint32(1)
	subStruct.TestUint64 = uint64(1)
	subStruct.TestFloat32 = float32(1)
	subStruct.TestFloat64 = float64(1)
	subStruct.PTestString = &subStruct.TestString
	subStruct.PTestInt8 = &subStruct.TestInt8
	subStruct.PTestInt16 = &subStruct.TestInt16
	subStruct.PTestInt32 = &subStruct.TestInt32
	subStruct.PTestInt64 = &subStruct.TestInt64
	subStruct.PTestUint8 = &subStruct.TestUint8
	subStruct.PTestUint16 = &subStruct.TestUint16
	subStruct.PTestUint32 = &subStruct.TestUint32
	subStruct.PTestUint64 = &subStruct.TestUint64
	subStruct.PTestFloat32 = &subStruct.TestFloat32
	subStruct.PTestFloat64 = &subStruct.TestFloat64
	subStruct.TestSubMap = make(map[string]string)
	subStruct.TestSubMap["testSubKey1"] = "testSubValue1"
	subStruct.TestSubMap["testSubKey2"] = "testSubValue2"
	subStruct.PTestSubMap = &subStruct.TestSubMap
	subStruct.TestSlice = append(subStruct.TestSlice, "testSubSlice")
	subStruct.PTestSlice = &subStruct.TestSlice

	testStruct := TestStruct{}
	testStruct.TestString = str
	testStruct.TestInt8 = int8(1)
	testStruct.TestInt16 = int16(1)
	testStruct.TestInt32 = int32(1)
	testStruct.TestInt64 = int64(1)
	testStruct.TestUint8 = uint8(1)
	testStruct.TestUint16 = uint16(1)
	testStruct.TestUint32 = uint32(1)
	testStruct.TestUint64 = uint64(1)
	testStruct.TestFloat32 = float32(1)
	testStruct.TestFloat64 = float64(1)
	// testProto := &spb.CreateRoomReq{}
	// testProto.TickRate = 12
	// testProto.RandSeed = 33
	// testStruct.TestProto = *testProto

	testStruct.PTestString = &testStruct.TestString
	testStruct.PTestInt8 = &testStruct.TestInt8
	testStruct.PTestInt16 = &testStruct.TestInt16
	testStruct.PTestInt32 = &testStruct.TestInt32
	testStruct.PTestInt64 = &testStruct.TestInt64
	testStruct.PTestUint8 = &testStruct.TestUint8
	testStruct.PTestUint16 = &testStruct.TestUint16
	testStruct.PTestUint32 = &testStruct.TestUint32
	testStruct.PTestUint64 = &testStruct.TestUint64
	testStruct.PTestFloat32 = &testStruct.TestFloat32
	testStruct.PTestFloat64 = &testStruct.TestFloat64
	//	testStruct.PTestProto = testProto

	testStruct.TestMap = make(map[string]string)
	testStruct.PTestMap = &testStruct.TestMap
	testStruct.TestMap["testKey"] = "testValue"
	testStruct.TestSlice = append(testStruct.TestSlice, "testSlice")
	testStruct.PTestSlice = &testStruct.TestSlice
	testStruct.TestStru = subStruct
	testStruct.PTestStru = &subStruct
	testStruct.TestPointMap = make(map[string]*TestSubStruct)
	testStruct.TestPointMap["testKey"] = &subStruct

	testMap := make(map[string]string)
	testMap["testKey"] = "testValue"
	var TestSlice []string
	TestSlice = append(TestSlice, "testSlice")

	data := call( /*testProto, */ uint16(11), float32(33), &str, &testStruct)

	test := &Test{}
	v := reflect.ValueOf(test)
	t := reflect.TypeOf(test)

	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		methodValue := v.MethodByName(method.Name)
		args, err := UnSerializeByFunc(methodValue, data)
		if err == nil {
			methodValue.Call(args)
		}
	}
}

type Father struct {
	MySon Son
	FName string
}

func (*Father) FatherSay() {
	fmt.Println("FatherSay ")
}

type Son struct {
	*Father
	SName string
}

func (*Son) SonSay() {
	fmt.Println("SonSay ")
}

func TestMechod(te *testing.T) {
	father := &Father{}
	father.FName = "FName"
	//father.MySon = &Son{}
	father.MySon.SName = "SName"
	father.MySon.Father = father

	v := reflect.ValueOf(father)
	//t := reflect.TypeOf(father)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	son := v.FieldByName("MySon")

	if son.Kind() != reflect.Ptr {
		son = son.Addr()
	}

	fmt.Println("father: ", v)
	fmt.Println("son: ", son)

	for i := 0; i < son.NumMethod(); i++ {
		method := son.Type().Method(i)
		fmt.Println("method.Name: ", method.Name)
	}

	if son.Kind() == reflect.Ptr {
		son = son.Elem()
	}

	for i := 0; i < son.NumField(); i++ {
		field := son.Field(i)
		fmt.Println("field: ", field)
	}
}

func TestSlice(te *testing.T) {
	var strVec []string
	strVec = append(strVec, "aa")
	strVec = append(strVec, "bb")
	strVec = append(strVec, "cc")

	data := SerializeNew(strVec)
	fmt.Println("data: ", data)

	strVec2 := &[]string{}
	UnSerializeNew(strVec2, data)

	fmt.Println("strVec2: ", strVec2)
}

type TestByteSlice struct {
	ByteSlice []uint8
}

func TestTypeSlice(te *testing.T) {
	byteSlice := []uint8("helloworld!")

	data := SerializeNew(&TestByteSlice{ByteSlice: []byte("helloworld!")})
	fmt.Println("data: ", string(byteSlice))

	strVec2 := &TestByteSlice{}
	UnSerializeNew(strVec2, data)

	fmt.Println("strVec2: ", string(strVec2.ByteSlice))
}

type TestData struct {
	//Msg spb.CreateRoomReq
}

func TestProto(te *testing.T) {
	// msg0 := &spb.CreateRoomReq{TickRate: 3.1, RandSeed: 345}

	// data := SerializeNew(msg0)
	// fmt.Println("data: ", data)

	// msg00 := &TestData{}
	// data00 := SerializeNew(msg00)
	// fmt.Println("data00: ", data00)

	// msg2 := &spb.CreateRoomReq{}
	// UnSerializeNew(msg2, data)

	// fmt.Println("msg2: ", msg2)
}

type TestParam struct {
	Name string
	Val  uint32
}
type TestGM struct {
	Name    string
	Val     uint32
	Param   TestParam
	TestMap map[string]uint32
}

// WriteInt8 写Int8
func (t *TestGM) Test(v *TestParam) {

}

func TestGMFunc(te *testing.T) {
	test := &TestGM{Name: "TestGM", Val: 321}
	test.TestMap = make(map[string]uint32)
	test.TestMap["hahha"] = 1
	test.TestMap["heihei"] = 33
	data, err := json.Marshal(test)
	if err != nil {
		return
	}

	fmt.Println("data: ", string(data))

	v := reflect.ValueOf(test)
	t := reflect.TypeOf(test)
	for i := 0; i < t.NumMethod(); i++ {
		methodName := t.Method(i).Name
		mv := v.MethodByName(methodName)

		mt := mv.Type()
		for i := 0; i < mt.NumIn(); i++ {
			pt := mt.In(i)
			fmt.Println("ptName: ", pt.Name(), ", typeStr: ", pt.String())
		}
	}
}
