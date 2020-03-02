package logrus

import (
	"context"
	"io"

	"github.com/giant-tech/go-service/base/itf/ilog"
	"github.com/sirupsen/logrus"
)

type formatterKey struct{}
type levelKey struct{}
type outKey struct{}
type hooksKey struct{}
type reportCallerKey struct{}
type exitKey struct{}

type Options struct {
	ilog.Options
}

func WithTextTextFormatter(formatter *logrus.TextFormatter) ilog.OptionFunc {
	return setOption(formatterKey{}, formatter)
}

func WithJSONFormatter(formatter *logrus.JSONFormatter) ilog.OptionFunc {
	return setOption(formatterKey{}, formatter)
}

func WithLevel(lvl ilog.Level) ilog.OptionFunc {
	return setOption(levelKey{}, lvl)
}

func WithOut(out io.Writer) ilog.OptionFunc {
	return setOption(outKey{}, out)
}

func WithLevelHooks(hooks logrus.LevelHooks) ilog.OptionFunc {
	return setOption(hooksKey{}, hooks)
}

// warning to use this option. because logrus doest not open CallerDepth option
// this will only print this package
func WithReportCaller(reportCaller bool) ilog.OptionFunc {
	return setOption(reportCallerKey{}, reportCaller)
}

func WithExitFunc(exit func(int)) ilog.OptionFunc {
	return setOption(exitKey{}, exit)
}

func setOption(k, v interface{}) ilog.OptionFunc {
	return func(o *ilog.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}
