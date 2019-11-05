package auth

import (
	"context"

	"github.com/giant-tech/go-service/base/mail/consts"
)

// TokenAuth 实现credentials.PerRPCCredentials接口
type TokenAuth struct {
	Token string
}

// GetRequestMetadata 获得请求的元数据
func (t TokenAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		consts.AuthHeader: t.Token,
	}, nil
}

// RequireTransportSecurity 启用TLS
func (TokenAuth) RequireTransportSecurity() bool {
	return true
}
