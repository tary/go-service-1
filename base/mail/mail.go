package mail

import (
	"errors"
	"fmt"
	"sync"

	"github.com/prometheus/common/log"
	"github.com/giant-tech/go-service/base/mail/auth"
	"github.com/giant-tech/go-service/base/mail/consts"
	"github.com/giant-tech/go-service/base/mail/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	client            pb.MailClient
	subMailChan       chan *pb.SubscribeRequest // 玩家上下线，订阅/取消新邮件通知
	destroyChan       chan string               // 关服，释放资源
	newMailHandlerMap *sync.Map
)

// Init 初始化
func Init(serverAddr string, token string) error {
	if client != nil {
		return nil // 只需初始化一次
	}
	// Create the client TLS credentials
	creds, err := credentials.NewClientTLSFromFile("../res/mail-server-key/"+consts.Cert, "gamemail.ztgame.com")
	if err != nil {
		panic(fmt.Errorf("could not load tls cert: %s", err))
	}

	serviceURL := serverAddr

	// We don't need to error here, as this creates a pool and connections
	// will happen later
	conn, _ := grpc.Dial(
		serviceURL,
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(auth.TokenAuth{
			Token: token,
		}),
	)

	client = pb.NewMailClient(conn)
	destroyChan = make(chan string, 0)

	return nil
}

// StartNewMailNotify 启用新邮件通知协程
// Register/Unregister 依赖该方法调用
func StartNewMailNotify() error {
	if subMailChan != nil {
		return nil // StartNewMailNotify 只需执行一次
	}
	subMailChan = make(chan *pb.SubscribeRequest, 1)
	newMailHandlerMap = &sync.Map{}

	go newMailNotifyGoroutine()
	return nil
}

// Destroy 释放资源，退出协程
func Destroy() {
	destroyChan <- "stop"
}

// newMailNotifyGoroutine 新邮件通知协程
func newMailNotifyGoroutine() {
	newMailNotifyClient, err := client.GetNewMailNotifications(context.Background())
	if err != nil {
		panic(err)
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("mail newMailNotifyGoroutine panic:", err)
			} else {
				log.Error("mail newMailNotifyGoroutine exited.")
			}
		}()
		deliverMail := func(newMail *pb.MailHeaderWithID, isBroadcast bool) {
			if isBroadcast {
				newMailHandlerMap.Range(func(_, v interface{}) bool {
					handler := v.(func(*pb.MailHeaderWithID, bool))
					handler(newMail, isBroadcast)
					return true
				})
			} else {
				userid := newMail.Header.To
				v, ok := newMailHandlerMap.Load(string(userid))
				if ok {
					handler := v.(func(*pb.MailHeaderWithID, bool))
					handler(newMail, isBroadcast)
				}
			}
		}

		for {
			recvedMails, err := newMailNotifyClient.Recv()
			if err != nil {
				log.Error("mail newMailNotifyGoroutine recv error:", err)
				return
			}

			for _, newMail := range recvedMails.BroadcastMails {
				deliverMail(newMail, true)
			}
			for _, newMail := range recvedMails.UserMails {
				deliverMail(newMail, false)
			}
		}
	}()

	for {
		select {
		case req := <-subMailChan:
			newMailNotifyClient.Send(req)
		case <-destroyChan:
			newMailNotifyClient.CloseSend()
			return // 退出协程
		}
	}
}

// Register 玩家上线后注册自己
func Register(userid string, newMailHandler func(*pb.MailHeaderWithID, bool)) error {

	newMailHandlerMap.Store(userid, newMailHandler)

	req := &pb.SubscribeRequest{
		IsSubscribe: true,
		To:          []byte(userid),
	}
	subMailChan <- req
	return nil
}

// Unregister 玩家下线时反注册自己
func Unregister(userid string) error {
	req := &pb.SubscribeRequest{
		IsSubscribe: false,
		To:          []byte(userid),
	}
	subMailChan <- req

	newMailHandlerMap.Delete(userid)
	return nil
}

// GetMailCount 获取当前邮件状态，是否有未读邮件(total,unread)
func GetMailCount(userid string, fromTime int64) (uint32, uint32) {

	req := &pb.GetMailCountRequest{
		To:       []byte(userid),
		FromTime: fromTime,
	}

	resp, err := client.GetMailCount(context.Background(), req)
	if err != nil {
		return 0, 0
	}
	unread := uint32(0)
	if resp.Total > resp.Read {
		unread = resp.Total - resp.Read
	}
	return resp.Total, unread
}

// Send 发送邮件
func Send(req *pb.SendRequest) (*pb.SendResponse, error) {
	return client.Send(context.Background(), req)
}

// Broadcast 广播邮件
func Broadcast(req *pb.BroadcastRequest) (*pb.BroadcastResponse, error) {
	return client.Broadcast(context.Background(), req)
}

// ListMails 获取邮件列表
func ListMails(userid string, fromTime int64, isBroadcast bool) []*pb.MailHeaderWithID {

	// 获取广播邮件
	startTime := &pb.ListRequest_TimeAndID{
		Time: fromTime,
	}
	req := &pb.ListRequest{
		IsBroadcast: isBroadcast,
		To:          []byte(userid),
		Start:       startTime,
	}
	resp, err := client.List(context.Background(), req)
	if err != nil {
		return nil
	}
	return resp.Mails
}

// Get 获取指定邮件内容
func Get(mailIdx *pb.MailIndex) (*pb.MailBody, error) {

	req := &pb.GetRequest{
		MailIndex: mailIdx,
	}
	resp, err := client.Get(context.Background(), req)
	if err != nil {
		return nil, err
	}
	if !resp.Result.Ok {
		return nil, errors.New(resp.Result.Error)
	}
	MarkAsRead(mailIdx)
	return resp.Body, nil
}

// Delete 删除指定邮件
func Delete(mailIdx *pb.MailIndex) error {
	req := &pb.DeleteRequest{
		MailIndex: mailIdx,
	}
	resp, err := client.Delete(context.Background(), req)
	if err != nil {
		return err
	}
	if !resp.Result.Ok {
		return errors.New(resp.Result.Error)
	}
	return nil
}

// MarkAsRead 标记为已读文件（Get后自动调用）
func MarkAsRead(mailIdx *pb.MailIndex) error {
	req := &pb.MarkAsReadRequest{
		MailIndex: mailIdx,
	}
	resp, err := client.MarkAsRead(context.Background(), req)
	if err != nil {
		return err
	}
	if !resp.Result.Ok {
		return errors.New(resp.Result.Error)
	}
	return nil
}

// MarkAttachmentsAsReceived 领取附件
func MarkAttachmentsAsReceived(mailIdx *pb.MailIndex) error {
	req := &pb.MarkAttachmentsAsReceivedRequest{
		MailIndex: mailIdx,
	}
	resp, err := client.MarkAttachmentsAsReceived(context.Background(), req)
	if err != nil {
		return err
	}
	if resp.CheckFailed {
		return errors.New("already received")
	}
	if !resp.Result.Ok {
		return errors.New(resp.Result.Error)
	}
	return nil
}
