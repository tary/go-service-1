package zap

import (
	"context"

	"github.com/giant-tech/go-service/base/itf/ilog"
)

// setOption returns a function to setup a context with given value
func setOption(k, v interface{}) ilog.OptionFunc {
	return func(o *ilog.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}
