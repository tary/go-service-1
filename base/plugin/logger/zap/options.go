package zap

import (
	"github.com/giant-tech/go-service/base/itf/ilog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Options struct {
	logger.Options
}

type configKey struct{}

// WithConfig pass zap.Config to logger
func WithConfig(c zap.Config) ilog.OptionFunc {
	return setOption(configKey{}, c)
}

type encoderConfigKey struct{}

// WithEncoderConfig pass zapcore.EncoderConfig to logger
func WithEncoderConfig(c zapcore.EncoderConfig) ilog.OptionFunc {
	return setOption(encoderConfigKey{}, c)
}

type levelKey struct{}

// WithLevel pass log level
func WithLevel(l ilog.Level) ilog.OptionFunc {
	return setOption(levelKey{}, l)
}
