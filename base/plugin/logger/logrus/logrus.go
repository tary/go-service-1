package logrus

import (
	"io"
	"os"

	"github.com/giant-tech/go-service/base/itf/ilog"
	"github.com/sirupsen/logrus"
)

var (
	formatter    logrus.Formatter = new(logrus.TextFormatter)
	lvl                           = logrus.InfoLevel
	out          io.Writer        = os.Stderr
	hooks                         = make(logrus.LevelHooks)
	reportCaller                  = false
	exit                          = os.Exit
)

type logger struct {
	*logrus.Logger
}

func (l *logger) Fields(fields map[string]interface{}) ilog.ILogger {
	// shall we need pool here?
	// but logrus already has pool for its entry.
	return &logger{logrus.WithFields(fields).Logger}
}

func (l *logger) Error(err error) ilog.ILogger {
	return &logger{logrus.WithError(err).Logger}
}

func (l *logger) Init(opts ...ilog.OptionFunc) error {
	options := &Options{}
	for _, o := range opts {
		o(&options.Options)
	}

	if options.Context != nil {
		f, ok := options.Context.Value(formatterKey{}).(logrus.Formatter)
		if ok {
			formatter = f
		}

		l, ok := options.Context.Value(levelKey{}).(ilog.Level)
		if ok {
			lvl = convertToLogrusLevel(l)
		}

		o, ok := options.Context.Value(outKey{}).(io.Writer)
		if ok {
			out = o
		}

		h, ok := options.Context.Value(hooksKey{}).(logrus.LevelHooks)
		if ok {
			hooks = h
		}

		r, ok := options.Context.Value(reportCallerKey{}).(bool)
		if ok {
			if r == true {
			}
			reportCaller = r
		}

		e, ok := options.Context.Value(exitKey{}).(func(int))
		if ok {
			exit = e
		}
	}

	l.Logger = &logrus.Logger{
		Out:          out,
		Formatter:    formatter,
		Hooks:        hooks,
		Level:        lvl,
		ExitFunc:     exit,
		ReportCaller: reportCaller,
	}

	return nil
}

func (l *logger) SetLevel(level ilog.Level) {
	l.Logger.SetLevel(convertToLogrusLevel(level))
}

func (l *logger) Level() ilog.Level {
	return convertLevel(l.Logger.Level)
}

func (l *logger) Log(level ilog.Level, args ...interface{}) {
	//l.Logger.Log(convertToLogrusLevel(level), args...)
}

func (l *logger) Logf(level ilog.Level, format string, args ...interface{}) {
	//l.Logger.Logf(convertToLogrusLevel(level), format, args...)
}

func (l *logger) String() string {
	return "logrus"
}

// New builds a new logger based on options
func NewLogger(opts ...ilog.OptionFunc) ilog.ILogger {
	l := &logger{}
	_ = l.Init(opts...)
	return l
}

func convertToLogrusLevel(level ilog.Level) logrus.Level {
	switch level {
	case ilog.TraceLevel:
		return logrus.TraceLevel
	case ilog.DebugLevel:
		return logrus.DebugLevel
	case ilog.InfoLevel:
		return logrus.InfoLevel
	case ilog.WarnLevel:
		return logrus.WarnLevel
	case ilog.ErrorLevel:
		return logrus.ErrorLevel
	case ilog.PanicLevel:
		return logrus.PanicLevel
	case ilog.FatalLevel:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}

func convertLevel(level logrus.Level) ilog.Level {
	switch level {
	case logrus.TraceLevel:
		return ilog.TraceLevel
	case logrus.DebugLevel:
		return ilog.DebugLevel
	case logrus.InfoLevel:
		return ilog.InfoLevel
	case logrus.WarnLevel:
		return ilog.WarnLevel
	case logrus.ErrorLevel:
		return ilog.ErrorLevel
	case logrus.PanicLevel:
		return ilog.PanicLevel
	case logrus.FatalLevel:
		return ilog.FatalLevel
	default:
		return ilog.InfoLevel
	}
}
