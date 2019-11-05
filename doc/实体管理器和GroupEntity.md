---
title: 实体管理器
---

完整代码示例请参考：http://github.com/tech/public/go-service/examples/tree/master/entity-call-entity

# 实体管理器
实体管理器可以管理多个实体，

创建实体管理器的接口
```
// NewEntities 创建一个新的Entities
func NewEntities(isMultilThread bool, ilocal iserver.ILocalService) *Entities {
	...
}
```
其中isMultilThread参数表示此管理器中的实体是否会开启自己的协程。
如果isMultilThread为true，则实体需要实现 ```Run()```接口，Run()会在新协程中开启```go Run()```。
如果isMultilThread为true，则说明实体的每次tick操作均交给服务协程。

每一个服务均为一个实体管理器，可以在配置文件中设置服务中的实体是否开启自己的协程。
在entity-call-entity示例中，gateway service配置的isMultilThread为true，team service中
配置为false。


# GroupEntity
GroupEntity是一个特殊的实体，是实体和实体管理器的合集。Team就是一个GroupEntity，Team自己是个实体，而且还管理着多个TeamMember。

GroupEntity的定义如下：
```
type GroupEntity struct {
	Entity
	*Entities
}
```

如果要实现一个队伍，需要包含一个GroupEntity，且在初始化和析构时需要调用
```OnGroupInit``` 和 ```OnGroupDestroy```

```
// Team 队伍
type Team struct {
	entity.GroupEntity
}


// OnInit 初始化
func (t *Team) OnInit(initData interface{}) error {
	t.GroupEntity.OnGroupInit()
	return nil
}

// OnLoop 每帧调用
func (t *Team) OnLoop() {
	t.GroupEntity.OnGroupLoop()
}

// OnDestroy 销毁
func (t *Team) OnDestroy() {
	t.GroupEntity.OnGroupDestroy()
}
```
