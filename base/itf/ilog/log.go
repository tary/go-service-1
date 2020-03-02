package ilog

import (
	"fmt"
	"os"
	"sync/atomic"
)

// level is a log level
type Level int32

const (
	TraceLevel Level = iota - 1
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

var (
	// the local defaultLog
	defaultLog ILogger

	// default log level is info
	level = InfoLevel

	// prefix for all messages
	prefix string
)

// ILogger is a generic logging interface
type ILogger interface {
	// Init initialises options
	Init(options ...OptionFunc) error
	// Level returns the logging level
	Level() Level
	// Log inserts a log entry.  Arguments may be handled in the manner
	// of fmt.Print, but the underlying defaultLog may also decide to handle
	// them differently.
	Log(level Level, v ...interface{})
	// Logf insets a log entry.  Arguments are handled in the manner of
	// fmt.Printf.
	Logf(level Level, format string, v ...interface{})
	// Fields set fields to always be logged
	//Fields(fields map[string]interface{}) Logger
	// Error set `error` field to be logged
	//Error(err error) Logger
	// SetLevel updates the logging level.
	SetLevel(Level)
	// String returns the name of defaultLog
	String() string
}

func init() {
	switch os.Getenv("SERVICE_LOG_LEVEL") {
	case "trace":
		level = TraceLevel
	case "debug":
		level = DebugLevel
	case "warn":
		level = WarnLevel
	case "info":
		level = InfoLevel
	case "error":
		level = ErrorLevel
	case "fatal":
		level = FatalLevel
	}
}

func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case WarnLevel:
		return "warn"
	case InfoLevel:
		return "info"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	default:
		return "unknown"
	}
}

// WithLevel logs with the level specified
func WithLevel(l Level, v ...interface{}) {
	if l > level {
		return
	}
	defaultLog.Log(l, v...)
}

// WithLevel logs with the level specified
func WithLevelf(l Level, format string, v ...interface{}) {
	if l > level {
		return
	}
	defaultLog.Logf(l, format, v...)
}

// Trace provides trace level logging
func Trace(v ...interface{}) {
	WithLevel(TraceLevel, v...)
}

// Tracef provides trace level logging
func Tracef(format string, v ...interface{}) {
	WithLevelf(TraceLevel, format, v...)
}

// Debug provides debug level logging
func Debug(v ...interface{}) {
	WithLevel(DebugLevel, v...)
}

// Debugf provides debug level logging
func Debugf(format string, v ...interface{}) {
	WithLevelf(DebugLevel, format, v...)
}

// Warn provides warn level logging
func Warn(v ...interface{}) {
	WithLevel(WarnLevel, v...)
}

// Warnf provides warn level logging
func Warnf(format string, v ...interface{}) {
	WithLevelf(WarnLevel, format, v...)
}

// Info provides info level logging
func Info(v ...interface{}) {
	WithLevel(InfoLevel, v...)
}

// Infof provides info level logging
func Infof(format string, v ...interface{}) {
	WithLevelf(InfoLevel, format, v...)
}

// Error provides warn level logging
func Error(v ...interface{}) {
	WithLevel(ErrorLevel, v...)
}

// Errorf provides warn level logging
func Errorf(format string, v ...interface{}) {
	WithLevelf(ErrorLevel, format, v...)
}

// Fatal logs with Log and then exits with os.Exit(1)
func Fatal(v ...interface{}) {
	WithLevel(FatalLevel, v...)
	os.Exit(1)
}

// Fatalf logs with Logf and then exits with os.Exit(1)
func Fatalf(format string, v ...interface{}) {
	WithLevelf(FatalLevel, format, v...)
	os.Exit(1)
}

// SetLogger sets the local defaultLog
func SetLogger(l ILogger) {
	defaultLog = l
}

// GetLogger returns the local defaultLog
func GetLogger() ILogger {
	return defaultLog
}

// SetLevel sets the log level
func SetLevel(l Level) {
	atomic.StoreInt32((*int32)(&level), int32(l))
}

// GetLevel returns the current level
func GetLevel() Level {
	return level
}

// SetPrefix Set a prefix for the defaultLog
func SetPrefix(p string) {
	prefix = p
}

// Name Set service name
func Name(name string) {
	prefix = fmt.Sprintf("[%s]", name)
}
