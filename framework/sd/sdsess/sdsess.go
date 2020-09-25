package sdsess

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/giant-tech/go-service/framework/net/inet"

	"github.com/cihub/seelog"
)

var (
	sessMap     *sync.Map
	typeSessMap *sync.Map
)

// init 包初始化
func init() {
	sessMap = &sync.Map{}
	typeSessMap = &sync.Map{}
}

// GetServiceNewCheck 判断一个连接是否是新连接
/*func GetServiceNewCheck() servicediscovery.ICheckNewService {
	return &_check{}
}
*/

// serverSessionInfo serverSessionInfo
type serverSessionInfo struct {
	serverID   uint64
	serverType int32
	sess       inet.ISessionBase
}

// GetSession 获取连接
func GetSession(serverID uint64) (inet.ISessionBase, error) {
	srv, ok := sessMap.Load(serverID)
	if ok {
		return srv.(*serverSessionInfo).sess, nil
	}

	return nil, fmt.Errorf("serverID %d not exist", serverID)
}

// GetRandSession 随机一个某类型的连接
func GetRandSession(serverType int32) (uint64, inet.ISessionBase, error) {
	v, _ := typeSessMap.LoadOrStore(serverType, []*serverSessionInfo{})
	lst := v.([]*serverSessionInfo)
	num := int32(len(lst))

	if num > 0 {
		info := lst[rand.Int31n(num)]
		return info.serverID, info.sess, nil
	}
	return 0, nil, fmt.Errorf("server list is empty %d", uint32(serverType))
}

// AddSession 添加连接，认证后添加
func AddSession(serverID uint64, serverType int32, sess inet.ISessionBase) {
	info := &serverSessionInfo{
		serverID:   serverID,
		serverType: serverType,
		sess:       sess,
	}
	sessMap.Store(serverID, info)

	updateTypeSession(info)
	seelog.Info("AddSession:", serverID, " ", serverType)
}

// DeleteSession 删除连接，断开连接时删除
func DeleteSession(serverID uint64) {
	srv, ok := sessMap.Load(serverID)
	if !ok {
		return
	}
	sessMap.Delete(serverID)

	info := srv.(*serverSessionInfo)
	deleteTypeSession(info)
	seelog.Info("delete server session ", serverID, ", type ", info.serverType)
}

// updateTypeSession updateTypeSession
func updateTypeSession(info *serverSessionInfo) {
	v, _ := typeSessMap.LoadOrStore(info.serverType, []*serverSessionInfo{})
	lst := v.([]*serverSessionInfo)

	for _, srv := range lst {
		if srv.serverID == info.serverID {
			srv.sess = info.sess
			return
		}
	}
	typeSessMap.Store(info.serverType, append(lst, info))
}

// deleteTypeSession deleteTypeSession
func deleteTypeSession(info *serverSessionInfo) {
	v, _ := typeSessMap.LoadOrStore(info.serverType, []*serverSessionInfo{})
	lst := v.([]*serverSessionInfo)

	for i, srv := range lst {
		if srv.serverID == info.serverID {
			lst = append(lst[0:i], lst[i+1:]...)
			typeSessMap.Store(info.serverType, lst)
			return
		}
	}
}

/*type check struct {
}

func (*check) IsNewService(serverID uint64, serverType uint8) bool {
	_, err := GetSession(serverID)
	return err == nil
}
*/
