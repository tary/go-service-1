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

# How to start
- Environment preparation
	* Need to install go.1.12 and above, no need to configure GOPATH
	* The editor can be vscode, or a normal text editor
	* If you create your own project, you can refer to the following example to build the project, the batch processing inside will be automatically downloaded, compiled and run.


# Code example
## Get service and call service method
Asynchronously call the SetName method provided by TeamService and set the name to "NewName"
```
randProxy := iserver.GetServiceProxyMgr().GetRandService(servicetype.TeamService)
err := randProxy.AsyncCall("SetName", "NewName")
```
Synchronously call the GetName method provided by TeamService, and assign the returned result to name
```
randProxy := iserver.GetServiceProxyMgr().GetRandService(servicetype.TeamService)
//Used to save the return value of a synchronous call
var name string
err := randProxy.SyncCall(&name, "GetName")
```

## Rpc call inside the entity (entity is an entity object)
```
  entity.AsyncCall(lobby, "SetName", "NewName")
```

## Get entity agent and call entity method

Asynchronously call the SetName method of the Team entity, and set the name of the team whose team ID is "TeamID" to "NewTeamName"
```
  err := entity.NewEntityProxy(TeamID).AsyncCall(servicetype.TeamService, "SetName", "NewTeamName")
```

Synchronously call the GetName method of the team entity to obtain the name of the team whose team ID is "TeamID"
```
  var teamName string
  err := entity.NewEntityProxy(TeamID).SyncCall(servicetype.TeamService, &teamName, "GetName")
```

## Complete code example
- Unity client (c#) calls the server

	https://github.com/giant-tech/go-service-examples/tree/master/client-service-demo-csharp

- go client calls server

	https://github.com/giant-tech/go-service-examples/tree/master/client-call-serviceA

- how to call between entities

	https://github.com/giant-tech/go-service-examples/tree/master/entity-call-entity

- how to call between two services

	https://github.com/giant-tech/go-service-examples/tree/master/serviceA-call-serviceB

## qq(Welcome to discuss)
- 942711528
## If the github download is slow, you can go to gitee to download, download address:
- https://gitee.com/yekoufeng/go-service
- https://gitee.com/yekoufeng/go-service-examples
