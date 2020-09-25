package inet

import (
	"github.com/giant-tech/go-service/base/imsg"
)

// ISessionBase session基础
type ISessionBase interface {
	Send(imsg.IMsg) error
	//有其它需求再加
}
