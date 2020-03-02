package zerolog

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime/debug"

	
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type Mode uint8

const (
	Production Mode = iota + 1
	Development
)

var (
	//  It's common to set this to a file, or leave it default which is `os.Stderr`
	out io.Writer = os.Stderr
	// Function to exit the application, defaults to `os.Exit()`
	exitFunc = os.Exit
	// Flag for whether to log caller info (off by default)
	reportCaller = false
	// use this logger as system wide default logger
	useAsDefault = false
	// The logging level the logger should log at.
	// This defaults to 100 means not explicitly set by user
	level      ilog.Level = 100
	fields     map[string]interface{}
	hooks      []zerolog.Hook
	timeFormat string
	// default Production (1)
	mode Mode = Production
)

type zeroLogger struct {
	nativelogger zerolog.Logger
}

func (l *zeroLogger) Fields(fields map[string]interface{}) ilog.ILogger {
	return &zeroLogger{l.nativelogger.With().Fields(fields).Logger()}
}

func (l *zeroLogger) Error(err error) ilog.ILogger {
	return &zeroLogger{
		l.nativelogger.With().Fields(map[string]interface{}{zerolog.ErrorFieldName: err}).Logger(),
	}
}

func (l *zeroLogger) Init(opts ...ilog.OptionFunc) error {

	options := &Options{ilog.Options{Context: context.Background()}}
	for _, o := range opts {
		o(&options.Options)
	}

	if o, ok := options.Context.Value(outKey{}).(io.Writer); ok {
		out = o
	}
	if hs, ok := options.Context.Value(hooksKey{}).([]zerolog.Hook); ok {
		hooks = hs
	}
	if flds, ok := options.Context.Value(fieldsKey{}).(map[string]interface{}); ok {
		fields = flds
	}
	if lvl, ok := options.Context.Value(levelKey{}).(ilog.Level); ok {
		level = lvl
	}
	if tf, ok := options.Context.Value(timeFormatKey{}).(string); ok {
		timeFormat = tf
	}
	if exitFunction, ok := options.Context.Value(exitKey{}).(func(int)); ok {
		exitFunc = exitFunction
	}
	if caller, ok := options.Context.Value(reportCallerKey{}).(bool); ok && caller {
		reportCaller = caller
	}
	if useDefault, ok := options.Context.Value(useAsDefaultKey{}).(bool); ok && useDefault {
		useAsDefault = useDefault
	}
	if devMode, ok := options.Context.Value(developmentModeKey{}).(bool); ok && devMode {
		mode = Development
	}
	if prodMode, ok := options.Context.Value(productionModeKey{}).(bool); ok && prodMode {
		mode = Production
	}

	switch mode {
	case Development:
		zerolog.ErrorStackMarshaler = func(err error) interface{} {
			fmt.Println(string(debug.Stack()))
			return nil
		}
		consOut := zerolog.NewConsoleWriter(
			func(w *zerolog.ConsoleWriter) {
				if len(timeFormat) > 0 {
					w.TimeFormat = timeFormat
				}
				w.Out = out
				w.NoColor = false
			},
		)
		level = ilog.DebugLevel
		l.nativelogger = zerolog.New(consOut).
			Level(zerolog.DebugLevel).
			With().Timestamp().Stack().Logger()
	default: // Production
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		l.nativelogger = zerolog.New(out).
			Level(zerolog.InfoLevel).
			With().Timestamp().Stack().Logger()
	}

	// Change  Writer if not default
	if out != os.Stderr {
		l.nativelogger = l.nativelogger.Output(out)
	}

	// Set log Level if not default
	if level != 100 {
		//zerolog.SetGlobalLevel(loggerToZerologLevel(level))
		l.nativelogger = l.nativelogger.Level(loggerToZerologLevel(level))
	}

	// Adding hooks if exist
	if reportCaller {
		l.nativelogger = l.nativelogger.With().Caller().Logger()
	}
	for _, hook := range hooks {
		l.nativelogger = l.nativelogger.Hook(hook)
	}

	// Setting timeFormat
	if len(timeFormat) > 0 {
		zerolog.TimeFieldFormat = timeFormat
	}

	// Adding seed fields if exist
	if fields != nil {
		l.nativelogger = l.nativelogger.With().Fields(fields).Logger()
	}

	// Also set it as zerolog's Default logger
	if useAsDefault {
		zlog.Logger = l.nativelogger
	}

	return nil
}

func (l *zeroLogger) SetLevel(level ilog.Level) {
	//zerolog.SetGlobalLevel(loggerToZerologLevel(level))
	l.nativelogger = l.nativelogger.Level(loggerToZerologLevel(level))
}

func (l *zeroLogger) Level() ilog.Level {
	return ZerologToLoggerLevel(l.nativelogger.GetLevel())
}

func (l *zeroLogger) Log(level ilog.Level, args ...interface{}) {
	msg := fmt.Sprintf("%s", args)
	l.nativelogger.WithLevel(loggerToZerologLevel(level)).Msg(msg[1 : len(msg)-1])
	// Invoke os.Exit because unlike zerolog.Logger.Fatal zerolog.Logger.WithLevel won't stop the execution.
	if level == ilog.FatalLevel {
		exitFunc(1)
	}
}

func (l *zeroLogger) Logf(level ilog.Level, format string, args ...interface{}) {
	l.nativelogger.WithLevel(loggerToZerologLevel(level)).Msgf(format, args...)
	// Invoke os.Exit because unlike zerolog.Logger.Fatal zerolog.Logger.WithLevel won't stop the execution.
	if level == ilog.FatalLevel {
		exitFunc(1)
	}
}

func (l *zeroLogger) String() string {
	return "zerolog"
}

// NewLogger builds a new logger based on options
func NewLogger(opts ...ilog.OptionFunc) ilog.ILogger {
	l := &zeroLogger{}
	_ = l.Init(opts...)
	return l
}

// ParseLevel converts a level string into a logger Level value.
// returns an error if the input string does not match known values.
func ParseLevel(levelStr string) (lvl ilog.Level, err error) {
	if zLevel, err := zerolog.ParseLevel(levelStr); err == nil {
		return ZerologToLoggerLevel(zLevel), err
	} else {
		return lvl, fmt.Errorf("Unknown Level String: '%s' %w", levelStr, err)
	}
}

func loggerToZerologLevel(level ilog.Level) zerolog.Level {
	switch level {
	case ilog.TraceLevel:
		return zerolog.TraceLevel
	case ilog.DebugLevel:
		return zerolog.DebugLevel
	case ilog.InfoLevel:
		return zerolog.InfoLevel
	case ilog.WarnLevel:
		return zerolog.WarnLevel
	case ilog.ErrorLevel:
		return zerolog.ErrorLevel
	case ilog.PanicLevel:
		return zerolog.PanicLevel
	case ilog.FatalLevel:
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

func ZerologToLoggerLevel(level zerolog.Level) ilog.Level {
	switch level {
	case zerolog.TraceLevel:
		return ilog.TraceLevel
	case zerolog.DebugLevel:
		return ilog.DebugLevel
	case zerolog.InfoLevel:
		return ilog.InfoLevel
	case zerolog.WarnLevel:
		return ilog.WarnLevel
	case zerolog.ErrorLevel:
		return ilog.ErrorLevel
	case zerolog.PanicLevel:
		return ilog.PanicLevel
	case zerolog.FatalLevel:
		return ilog.FatalLevel
	default:
		return ilog.InfoLevel
	}
}
