---
title: 如何使用RPC函数
---
完整代码示例：http://github.com/tech/public/go-service/examples/tree/master/client-call-serviceA


# RPC函数
RPC是指名字以"RPC"开头的函数，比如定义一个名为Hello的函数：
```
func (s *NewService) RPCHello(str string) string {

}
```

不管是服务还是实体，都对外提供0到多个RPC方法给外部调用。


## RPC函数的参数
### 参数支持的类型

类型|说明
--|:--
Bool|
Int8|
Int16|
Int32|
Int64|
Uint8|
Uint16|
Uint32|
Uint64|
Float32|
Float64|
String|
Struct| 如果参数是结构体，需要设置为指针
Slice|
Array|目前C#客户端不支持，服务器之间支持
Map|
Ptr|上面所有类型的指针

### 序列化和反序列化
我们提供了一套默认的序列化和反序列化机制。
同时支持protobuff以及自定义的序列化和反序列化，需要满足如下接口
```
// IProtoMsg proto消息接口
type IProtoMsg interface {
	MarshalTo(data []byte) (n int, err error)
	Unmarshal(data []byte) error
	Size() (n int)
}
```

## RPC函数的返回值
RPC函数的返回值只有在同步调用的时候才有作用，同步调用允许一个返回值或者没有返回值，不允许多个返回值。

## RPC函数的调用
RPC函数的调用支持同步调用和异步调用。

服务的RPC函数调用接口

```
// AsyncCall 异步调用，立即返回
AsyncCall(methodName string, args ...interface{}) error

// SyncCall 同步调用，等待返回
SyncCall(retPtr interface{}, methodName string, args ...interface{}) error

```

实体的RPC函数调用接口，比服务的RPC调用多了一个服务类型，表示被调用实体所在的服务类型。
```
// AsyncCall 异步调用，立即返回
AsyncCall(stype idata.ServiceType, methodName string, args ...interface{}) error

// SyncCall 同步调用，等待返回
SyncCall(stype idata.ServiceType, retPtr interface{}, methodName string, args ...interface{}) error
```
