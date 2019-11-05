package entity

import (
	"fmt"

	"github.com/cihub/seelog"
)

// String Entity基础信息
func (e *Entity) String() string {
	return fmt.Sprintf("[Type:%s EntityID:%d]",
		e.entityType, e.entityID)
}

// Info 日志
func (e *Entity) Info(v ...interface{}) {
	params := []interface{}{e.String()}
	params = append(params, v...)
	seelog.Info(params...)
}

// Infof 日志
func (e *Entity) Infof(format string, v ...interface{}) {
	ff := "%s " + format
	params := []interface{}{e.String()}
	params = append(params, v...)
	seelog.Infof(ff, params...)
}

// Warn 日志
func (e *Entity) Warn(v ...interface{}) {
	params := []interface{}{e.String()}
	params = append(params, v...)
	seelog.Warn(params...)
}

// Warnf 日志
func (e *Entity) Warnf(format string, v ...interface{}) {
	ff := "%s " + format
	params := []interface{}{e.String()}
	params = append(params, v...)
	seelog.Warnf(ff, params...)
}

// Error 日志
func (e *Entity) Error(v ...interface{}) {
	params := []interface{}{e.String()}
	params = append(params, v...)
	seelog.Error(params...)
}

// Errorf 日志
func (e *Entity) Errorf(format string, v ...interface{}) {
	ff := "%s " + format
	params := []interface{}{e.String()}
	params = append(params, v...)
	seelog.Errorf(ff, params...)
}

// Debug 日志
func (e *Entity) Debug(v ...interface{}) {
	params := []interface{}{e.String()}
	params = append(params, v...)
	seelog.Debug(params...)
}

// Debugf 日志
func (e *Entity) Debugf(format string, v ...interface{}) {
	ff := "%s " + format
	params := []interface{}{e.String()}
	params = append(params, v...)
	seelog.Debugf(ff, params...)
}
