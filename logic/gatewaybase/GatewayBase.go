package gatewaybase

import (
	"github.com/giant-tech/go-service/framework/iserver"
	"github.com/giant-tech/go-service/framework/msgdef"
	"github.com/giant-tech/go-service/framework/net/server"
	"github.com/giant-tech/go-service/framework/service"
	"github.com/giant-tech/go-service/logic/gatewaybase/igateway"
	"github.com/giant-tech/go-service/logic/gatewaybase/proc"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

// GatewayBase 大厅服务器
type GatewayBase struct {
	service.BaseService
	ClientServer *server.Server
}

// OnInit 初始化
func (srv *GatewayBase) OnInit(lobbyServer interface{}) error {
	log.Info("GatewayBase OnInit")

	sName := srv.GetSName()

	var err error

	//创建监听服务
	LobbyListenAddr := viper.GetString(sName + ".SvrListenAddr")
	srvMaxConn := viper.GetInt(sName + ".SvrMaxConn")
	srv.ClientServer, err = server.New("tcp", LobbyListenAddr, srvMaxConn)
	if err != nil {
		log.Error("gensvr failed, err: ", err, ", srvListenAddr: ", LobbyListenAddr)
		return err
	}

	// 添加MsgProc, 这样新连接创建时会注册处理函数
	srv.ClientServer.AddMsgProc(&proc.PGatewayServer{IServiceBase: lobbyServer.(iserver.IServiceBase)})

	srv.ClientServer.SetVerifyMsgID(msgdef.LoginReqMsgID)

	//是否加密
	isEncrypt := viper.GetBool(sName + ".IsEncrypt")
	srv.ClientServer.SetEncrypt(isEncrypt)

	outerIP := viper.GetString(sName + ".OuterIP")
	outerAddr := outerIP + ":" + srv.ClientServer.GetListenPort()
	log.Debug("outer addr: ", outerAddr)

	go srv.ClientServer.Run()

	return nil
}

// OnDestroy 退出时调用
func (srv *GatewayBase) OnDestroy() {
	log.Info("GatewayBase OnDestroy")

	if srv.ClientServer != nil {
		srv.ClientServer.Close()
	}

	srv.TravsalEntity("Player", func(e iserver.IEntity) {

		iclose, ok := e.(igateway.ICloseHandler)
		if ok {
			e.PostFunction(func() { iclose.OnClose() })
		}
	})

	//清空IEntities
	//srv.Destroy()
}
