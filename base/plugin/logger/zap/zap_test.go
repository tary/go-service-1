package zap

import (
	"testing"

	"github.com/giant-tech/go-service/base/itf/ilog"
)

func TestName(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}

	if l.String() != "zap" {
		t.Errorf("name is error %s", l.String())
	}

	t.Logf("test logger name: %s", l.String())
}

func TestLogf(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}

	l.Logf(ilog.InfoLevel, "test logf: %s", "name")
}

func TestSetLevel(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}

	l.SetLevel(ilog.DebugLevel)
	l.Logf(ilog.DebugLevel, "test show debug: %s", "debug msg")

	l.SetLevel(ilog.InfoLevel)
	l.Logf(ilog.DebugLevel, "test non-show debug: %s", "debug msg")
}
