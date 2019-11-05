package entity

import (
	"github.com/giant-tech/go-service/base/serializer"
	"github.com/giant-tech/go-service/framework/idata"
	"github.com/giant-tech/go-service/framework/msgdef"
)

//RpcType rpc类型
//type RpcType uint32

// const (
// 	RpcTypeClient RpcType = 1 // 发送给客户端
// 	RpcTypeLobby  RpcType = 2 // 发送给lobby
// 	RpcTypeCenter RpcType = 3 // 发送给Center
// 	RpcTypeRoom   RpcType = 6 // 发送给center中的room
// 	RpcTypeTeam   RpcType = 7 // 发送给centre中的team
// )

// // RPC 发送rpc消息给客户端
// func (e *Entity) RPC(rpcType RpcType, methodName string, args ...interface{}) error {
// 	// data := serializer.SerializeNew(args...)
// 	// msg := &msgdef.RpcMsg{
// 	// 	RpcType:      uint32(rpcType),
// 	// 	FromEntityID: e.GetEntityID(),
// 	// 	EntityID:     e.GetEntityID(),
// 	// 	MethodName:  methodName,
// 	// 	Data:         data,
// 	// }

// 	// if iserver.GetSrvInst().GetServerType() == uint8(idata.ServiceGateway) {
// 	// 	if rpcType == RpcTypeClient {
// 	// 		cli := e.GetClientSess()
// 	// 		if cli != nil {
// 	// 			cli.Send(msg)
// 	// 			return nil
// 	// 		} else {
// 	// 			log.Errorf("RPC failed: ClientSess is null, serverType: %d, rpcType: %d", iserver.GetSrvInst().GetServerType(), rpcType)
// 	// 			return fmt.Errorf("RPC failed: ClientSess is null, serverType: %d, rpcType: %d", iserver.GetSrvInst().GetServerType(), rpcType)
// 	// 		}
// 	// 	} else if rpcType == RpcTypeLobby {
// 	// 		log.Errorf("lobby to lobby RPC failed: serverType: %d, rpcType: %d", iserver.GetSrvInst().GetServerType(), rpcType)
// 	// 		return fmt.Errorf("lobby to lobby RPC failed: serverType: %d, rpcType: %d", iserver.GetSrvInst().GetServerType(), rpcType)
// 	// 	} else {
// 	// 		//转发到相应服务器
// 	// 		err := iserver.GetSrvInst().ForwardRpcMsg(msg)
// 	// 		if err != nil {
// 	// 			log.Debug("iserver.GetSrvInst().ForwardRpcMsg failed:", err)
// 	// 			return err
// 	// 		}
// 	// 	}
// 	// } else if iserver.GetSrvInst().GetServerType() == uint8(spb.ServerType_Center) {
// 	// 	if rpcType == RpcTypeLobby || rpcType == RpcTypeClient {
// 	// 		//转发到相应服务器
// 	// 		err := iserver.GetSrvInst().ForwardRpcMsg(msg)
// 	// 		if err != nil {
// 	// 			return err
// 	// 		}
// 	// 	} else {
// 	// 		log.Errorf("RPC failed: serverType: %d, rpcType: %d", iserver.GetSrvInst().GetServerType(), rpcType)
// 	// 		return fmt.Errorf("RPC failed: serverType: %d, rpcType: %d", iserver.GetSrvInst().GetServerType(), rpcType)
// 	// 	}
// 	// } else {
// 	// 	log.Errorf("RPC failed: serverType: %d, rpcType: %d", iserver.GetSrvInst().GetServerType(), rpcType)
// 	// 	return fmt.Errorf("RPC failed: serverType: %d, rpcType: %d", iserver.GetSrvInst().GetServerType(), rpcType)
// 	// }

// 	return nil
// }

// MakeRPC 构造一个RPC
func MakeRPC(sType idata.ServiceType, toEntityID uint64, methodName string, args ...interface{}) *msgdef.CallMsg {

	data := serializer.SerializeNew(args...)

	msg := &msgdef.CallMsg{}
	msg.SType = uint8(sType)
	msg.EntityID = toEntityID
	msg.MethodName = methodName
	msg.Params = data
	return msg
}
