package ilog

import (
	"context"
)

// Option for load profiles maybe
// eg. yml
// micro:
//   logger:
//     name:
//     dialect: zap/default/logrus
//     zap:
//       xxx:
//     logrus:
//       xxx:
type OptionFunc func(*Options)

type Options struct {
	Context context.Context
}
