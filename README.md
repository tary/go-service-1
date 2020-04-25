[![Build Status](https://travis-ci.org/atotto/travisci-golang-example.png)](https://travis-ci.org/atotto/travisci-golang-example)

go-service是用Golang开发的一款分布式服务器框架。

# Overview

go-service提供了一套分布式组件（分布式服务和分布式实体），以及它们所提供的RPC方法。我们可以在任何地方调用这些服务和实体所提供的方法，而不用关心这些服务和实体是不是与调用者在一个进程或机器。框架的三个核心概念为：App，Service，Entity。
- App：为一个进程，一台机器可以开启多个App，每个App可以根据配置加载一到多个Service。
- Service：开发者所开发的逻辑功能的集合（比如聊天服，网关服等），可以管理多个实体，也可以是无实体的服务。
- Entity：服务中管理的对象，比如Player、Team等，可以分布在一到多个Service中。
三者关系如下：

<img src="https://github.com/giant-tech/go-service/blob/master/resources/app-service-entity.jpg" />

Service和Entity的消息传递都交给App处理，Service之间并不建立连接。App会处理好服务发现和服务管理，以及本进程和跨进程的通讯。Service和Entity会提供一些RPC方法供他人调用，Entity的属性可以在不同Service之间进行自动同步。Service和Entity的RPC方法均支持同步调用和异步调用两种方式。

# 功能及特色
- app, service, entity架构，一个app可以包含一个或多个service, 每个service可以包含零到多个entity
- 基于redis的app发现，service发现。
- 支持service的同步和异步远程调用
- 支持entity的属性同步以及entity方法的同步和异步远程调用
- 网络库，支持tcp, kcp
- grpc支持
- elo算法支持
- db支持
	* mongodb mgo库支持，能够访问mongodb分片集群
	* redis库支持
- log库
- 可用vscode单步调试的分布式服务端，N变1
  * 一般来说，分布式服务端要启动很多进程，一旦进程多了，单步调试就变得非常困难，导致服务端开发基本上靠打log来查找问题。平常开发游戏逻辑也得开启一大堆进程，不仅启动慢，而且查找问题及其不方便，要在一堆堆日志里面查问题，这感觉非常糟糕，zeus框架使用了service设计，所有服务端内容都拆成了一个个service，启动时根据服务器类型挂载自己所需要的service。
- 随意可拆分功能的分布式服务端，1变N
  * 分布式服务端要开发多种类型的服务器进程，比如Login server，gate server，battle server，chat server friend server等等一大堆各种server，传统开发方式需要预先知道当前的功能要放在哪个服务器上，当功能越来越多的时候，比如聊天功能之前在一个中心服务器上，之后需要拆出来单独做成一个服务器，这时会牵扯到大量迁移代码的工作，烦不胜烦。zeus框架在平常开发的时候根本不太需要关心当前开发的这个功能会放在什么server上，只用一个进程进行开发，功能开发成组件的形式。发布的时候使用一份多进程的配置即可发布成多进程的形式，是不是很方便呢？随便你怎么拆分服务器。只需要修改极少的代码就可以进行拆分。不同的server挂上不同的组件就行了嘛！
- go语言天生跨平台
  * 提供windows ,linux的一键运行脚本
- go mod支持，能够自动下载需要的三方库，和框架的底层库

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

## 实体内部的rpc调用
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
