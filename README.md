[![Build Status](https://travis-ci.org/atotto/travisci-golang-example.png)](https://travis-ci.org/atotto/travisci-golang-example)

go-service is a distributed server framework developed with Golang

# Overview

go-service provides a set of distributed components (distributed services and distributed entities) and the RPC methods they provide. We can call the methods provided by these services and entities anywhere, regardless of whether these services and entities are in the same process or machine as the caller. The three core concepts of the framework are: App, Service, Entity
- App：As a process, a machine can open multiple apps, and each app can load one or more services according to the configuration
- Service：A collection of logical functions developed by developers (such as chat service, gateway service, etc.), which can manage multiple entities or non-entity services
- Entity：The objects managed in the service, such as Player, Team, etc., can be distributed in one or more Services
The relationship between the three is as follows：

<img src="https://github.com/giant-tech/go-service/blob/master/resources/app-service-entity.jpg" />

The messaging of Service and Entity is handed over to App for processing, and no connection is established between Services. App will handle service discovery and service management, as well as the communication between this process and cross-process. Service and Entity will provide some RPC methods for others to call, and the attributes of Entity can be automatically synchronized between different Services. Both Service and Entity RPC methods support synchronous and asynchronous calls

# Functions and features
- App, service, entity architecture, an app can contain one or more services, and each service can contain zero to multiple entities
- Redis-based app discovery, service discovery
- Support synchronous and asynchronous remote invocation of service
- Support entity attribute synchronization and entity's synchronous and asynchronous remote calls
- Network library supports tcp, kcp, grpc
- db 
	* Mongodb mgo library support, able to access mongodb sharded cluster
	* redis library support
- Distributed server that can be debugged with vscode single step, N changes to 1
  * Generally speaking, a distributed server needs to start many processes. Once there are more processes, single-step debugging becomes very difficult. As a result, server development basically depends on logging to find problems. Normally developing game logic also has to open a lot of processes. Not only is it slow to start, but it is also inconvenient to find the problem. It feels very bad to check the problem in a pile of logs. The go-service framework uses the service design, all logic Split into many services. Mount the service you need according to the server type at startup.
- Distributed server with functions that can be split at will, 1 becomes N
  * The distributed server needs to develop multiple types of server processes, such as Login server，gate server，battle server，chat server， friend server  and so on，Traditional development methods need to know in advance which server the current function will be placed on. When there are more and more functions, for example, the chat function was previously on a central server, and then it needs to be disassembled and made into a separate server. This will involve The work of migrating a lot of code is annoying. In normal development, the go-service framework does not need to care about what server the currently developed function will be placed on. It only uses one process for development, and the function is developed into the form of a component. Is it very convenient to use a multi-process configuration to publish in a multi-process form when publishing? Split the server as you like. It can be split with very little code modification. Just hang different components on different servers
- Go language is inherently cross-platform
  * Provide windows, linux one-click running script
- Go mod support, can automatically download the required tripartite library, and the framework's underlying library

# 如何开始
- 环境准备
	* 需要安装go.1.12及以上版本,无需配置GOPATH
	* 编辑器可以采用vscode,或者普通的文本编辑器
	* 如果新建自己的工程，可以参照下面的例子建工程，里面的批处理会自动下载、编译跑起来。

# 代码示例
## 服务的获取及服务方法的调用
异步调用TeamService提供的SetName方法，把名字设置为"NewName"
```
randProxy := iserver.GetServiceProxyMgr().GetRandService(servicetype.TeamService)
err := randProxy.AsyncCall("SetName", "NewName")
```
同步调用TeamService提供的GetName方法，并把返回的结果赋值给name
```
randProxy := iserver.GetServiceProxyMgr().GetRandService(servicetype.TeamService)
//用于保存同步调用的返回值
var name string
err := randProxy.SyncCall(&name, "GetName")
```

## 实体内部的rpc调用,entity为实体对象
```
  entity.AsyncCall(lobby, "SetName", "NewName")
```

## 实体代理的获取及实体方法的调用

异步调用Team实体的SetName方法，把队伍ID为"TeamID"的队伍名字设置为"NewTeamName"。
```
  err := entity.NewEntityProxy(TeamID).AsyncCall(servicetype.TeamService, "SetName", "NewTeamName")
```

同步调用队伍实体的GetName方法，获取队伍ID为"TeamID"的队伍的名字。
```
  var teamName string
  err := entity.NewEntityProxy(TeamID).SyncCall(servicetype.TeamService, &teamName, "GetName")
```

## 完整代码示例
- unity客户端(c#)调用服务器

	https://github.com/giant-tech/go-service-examples/tree/master/client-service-demo-csharp

- go客户端调用服务器

	https://github.com/giant-tech/go-service-examples/tree/master/client-call-serviceA

- entity之间如何调用

	https://github.com/giant-tech/go-service-examples/tree/master/entity-call-entity

- 两个service之间如何调用

	https://github.com/giant-tech/go-service-examples/tree/master/serviceA-call-serviceB

## qq讨论群
- 942711528
## github下载慢的，可到码云下载，下载地址：
- https://gitee.com/yekoufeng/go-service
- https://gitee.com/yekoufeng/go-service-examples
