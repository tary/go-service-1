package zerolog

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/giant-tech/go-service/base/itf/ilog"
	"github.com/rs/zerolog"
)

func TestName(t *testing.T) {
	l := NewLogger()

	if l.String() != "zerolog" {
		t.Errorf("error: name expected 'zerolog' actual: %s", l.String())
	}

	t.Logf("testing logger name: %s", l.String())
}

// func ExampleWithOut() {
// 	l := NewLogger(WithOut(os.Stdout), WithProductionMode())

// 	l.Logf(logger.InfoLevel, "testing: %s", "logf")

// 	// Output:
// 	// {"level":"info","time":"2020-02-14T22:15:36-08:00","message":"testing: logf"}
// }

func TestSetLevel(t *testing.T) {
	l := NewLogger()

	l.SetLevel(ilog.DebugLevel)
	l.Logf(ilog.DebugLevel, "test show debug: %s", "debug msg")

	l.SetLevel(ilog.InfoLevel)
	l.Logf(ilog.DebugLevel, "test non-show debug: %s", "debug msg")
}

func TestWithReportCaller(t *testing.T) {
	l := NewLogger(ReportCaller())

	l.Logf(ilog.InfoLevel, "testing: %s", "WithReportCaller")
}

func TestWithOut(t *testing.T) {
	l := NewLogger(WithOut(os.Stdout))

	l.Logf(ilog.InfoLevel, "testing: %s", "WithOut")
}

func TestWithDevelopmentMode(t *testing.T) {
	l := NewLogger(WithDevelopmentMode(), WithTimeFormat(time.Kitchen))

	l.Logf(ilog.InfoLevel, "testing: %s", "DevelopmentMode")
}

func TestWithFields(t *testing.T) {
	l := NewLogger()

	l.Fields(map[string]interface{}{
		"sumo":  "demo",
		"human": true,
		"age":   99,
	}).Logf(ilog.InfoLevel, "testing: %s", "WithFields")
}

func TestWithError(t *testing.T) {
	l := NewLogger()

	l.Error(errors.New("I am Error")).Logf(ilog.ErrorLevel, "testing: %s", "WithError")
}

func TestWithHooks(t *testing.T) {
	simpleHook := zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, msg string) {
		e.Bool("has_level", level != zerolog.NoLevel)
		e.Str("test", "logged")
	})

	l := NewLogger(WithHooks([]zerolog.Hook{simpleHook}))

	l.Logf(ilog.InfoLevel, "testing: %s", "WithHooks")
}
