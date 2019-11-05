package loginclt

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/GA-TECH-SERVER/zeus/framework/iserver"

	"github.com/GA-TECH-SERVER/zeus/framework/login"
	"github.com/GA-TECH-SERVER/zeus/logic/gatewaybase/sbase"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

var loginMgr *LoginCliMgr

// Init 包初始化
func Init(svrAddr string, sBase iserver.IServiceBase) {
	if loginMgr == nil {
		loginMgr = &LoginCliMgr{
			clientList:  list.New(),
			svrAddr:     svrAddr,
			serviceBase: sBase,
		}

		sbase.SetLoginMgr(loginMgr)

		go loginMgr.loopReportSvrLoad()
	}
}

// GetLoginCliMgr 获取登录管理器
func GetLoginCliMgr() *LoginCliMgr {
	return loginMgr
}

// LoginCliMgr 链接管理
type LoginCliMgr struct {
	serviceBase iserver.IServiceBase
	clientList  *list.List
	svrAddr     string
	listLock    sync.RWMutex
}

// AddLoginSession 添加登录sesison
func (lcm *LoginCliMgr) AddLoginSession( /*protocol,*/ addr string) error {
	lcm.listLock.Lock()
	defer lcm.listLock.Unlock()
	for l := lcm.clientList.Front(); l != nil; l = l.Next() {
		cli := l.Value.(LoginClient)
		if addr == cli.loginAddr {
			log.Errorf("Same login address :%s", addr)
			return fmt.Errorf("Same login address :%s ", addr)
		}
	}
	mc := &LoginClient{loginAddr: addr}
	lcm.clientList.PushBack(mc)
	return nil
}

// reportSvrLoad 报告服务器负载
func (lcm *LoginCliMgr) reportSvrLoad(lobbyAddr string, svrLoad uint32) error {

	loginAddr := lcm.GetRandLoginSessionAddr()
	if loginAddr == "" {
		return errors.New("login server not found")
	}

	msg := &login.UpdateLoadReq{
		OuterAddr: lobbyAddr,
		Load:      int(svrLoad),
	}
	data, err := json.Marshal(msg)
	if err != nil {
		log.Info(err)
		return err
	}

	url := "http://" + loginAddr + "/update_load"
	resp, err := http.Post(url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		log.Error(err)
		return err
	}
	resp.Body.Close()

	return nil
}

// PlayerLogout 玩家登出
func (lcm *LoginCliMgr) PlayerLogout(uid uint64) error {

	loginAddr := lcm.GetRandLoginSessionAddr()
	if loginAddr == "" {
		return errors.New("login server not found")
	}

	msg := &login.UserLogoutReq{
		OuterAddr: lcm.svrAddr,
		UID:       uid,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		log.Info(err)
		return err
	}

	url := "http://" + loginAddr + "/logout"
	resp, err := http.Post(url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		log.Error(err)
		return err
	}
	resp.Body.Close()

	return nil
}

// loopReportSvrLoad 循环报告负载
func (lcm *LoginCliMgr) loopReportSvrLoad() {

	defer func() {
		if err := recover(); err != nil {
			log.Error("loopReportSvrLoad panic:", err, ", Stack: ", string(debug.Stack()))
			if viper.GetString("Config.Recover") == "0" {
				panic(err)
			}
		}
	}()

	log.Info("begin loopReportSvrLoad with ", lcm.svrAddr)
	ticker := time.NewTicker(time.Duration(5) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := lcm.reportSvrLoad(lcm.svrAddr, lcm.serviceBase.EntityCount())
			if err != nil {
				log.Info(err)
			}
		}
	}
}

// GetRandLoginSessionAddr 拿随机登录session addr
func (lcm *LoginCliMgr) GetRandLoginSessionAddr() string {
	lcm.listLock.RLock()
	defer lcm.listLock.RUnlock()
	if lcm.clientList.Len() == 0 {
		return ""
	}

	var mc *LoginClient
	for l := lcm.clientList.Front(); l != nil; l = l.Next() {
		mc = l.Value.(*LoginClient)
		break
	}

	return mc.loginAddr
}

// CheckToken 监测token
func (lcm *LoginCliMgr) CheckToken(id uint64, token string) error {
	logoinAddr := lcm.GetRandLoginSessionAddr()
	msg := &login.LVerifyReq{
		Token: token,
		UID:   id,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		log.Info(err)
		return err
	}

	url := "http://" + logoinAddr + "/verify"
	resp, err := http.Post(url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		log.Error(err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	respMsg := &login.LVerifyAck{}
	json.Unmarshal(body, &respMsg)

	if err != nil {
		log.Error(err)
		return err
	}

	if respMsg.Result != 0 {
		log.Error(respMsg.Result)
		return errors.New("login fail")
	}

	return nil
}

// LoginClient 登录客户端
type LoginClient struct {
	//inet.ISession
	loginAddr string
}
