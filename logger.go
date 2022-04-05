package log4go

import (
	"errors"
	"runtime"
	"strings"
	"time"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : logger.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/5 20:59
* 修改历史 : 1. [2022/4/5 20:59] 创建文件 by NST
*/

var LogMetrics = &logMetrics{LogCount: 0, IgnoreCount: 0, AppendCount: 0}

type logMetrics struct {
	LogCount     int64
	IgnoreCount  int64
	AppendCount  int64
	InvalidCount int64
	// LogCount = IgnoreCount + AppendCount + InvalidCount
}

// A LogRecord contains all the pertinent information for each message
type LogRecord struct {
	Level   Level     // The log level
	Created time.Time // The time at which the log message was created (nanoseconds)
	Source  string    // The message source
	Message string    // The log message
}

type Logger struct {
	Name string
}

func (c Logger) Critical(arg0 any, args ...any) error {
	return c.Log(CRITICAL, arg0, args...)
}

func (c Logger) Error(arg0 any, args ...any) error {
	return c.Log(ERROR, arg0, args...)
}

func (c Logger) Warning(arg0 any, args ...any) error {
	return c.Log(WARNING, arg0, args...)
}

func (c Logger) Info(arg0 any, args ...any) error {
	return c.Log(INFO, arg0, args...)
}

func (c Logger) Trace(arg0 any, args ...any) error {
	return c.Log(TRACE, arg0, args...)
}

func (c Logger) Debug(arg0 any, args ...any) error {
	return c.Log(DEBUG, arg0, args...)
}

func (c Logger) Fine(arg0 any, args ...any) error {
	return c.Log(FINE, arg0, args...)
}

func (c Logger) Finest(arg0 any, args ...any) error {
	return c.Log(FINEST, arg0, args...)
}

func (c Logger) Log(level Level, arg0 any, args ...any) error {

	var msg string
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		msg = BuildString(first, args...)
	case func(...any) string:
		// Log the closure (no other arguments used)
		msg = first(args...)
	default:
		// Build a format string so that it will be similar to Sprint
		msg = BuildString(BuildString("%v", first)+strings.Repeat(" %v", len(args)), args...)
	}

	pc, _, lineno, ok := runtime.Caller(2)
	src := ""
	if ok {
		src = BuildString("%v:%v", runtime.FuncForPC(pc).Name(), lineno)
	}

	rec := &LogRecord{
		Level:   level,
		Created: time.Now(),
		Source:  src,
		Message: msg,
	}

	go c.intLogf(level, rec)

	return nil

	if int(level) >= int(ERROR) {
		return errors.New(msg)
	} else {
		return nil
	}
	//
}

func (c Logger) appendLog(lvl Level, rec *LogRecord, appenderRef LoggerAppenderReference) {
	LogMetrics.LogCount++

	if appenderRef.appender == nil || len(appenderRef.appender) == 0 {
		LogMetrics.InvalidCount++

		return
	}

	// if the log level is greater than logger's level, append the log, otherwise do nothing
	if int(lvl) >= int(appenderRef.level) {
		/*pc, _, lineno, ok := runtime.Caller(2)
		src := ""
		if ok {
			src = BuildString("%v:%d", runtime.FuncForPC(pc).Name(), lineno)
		}*/
		// Make the log record

		for _, appender := range appenderRef.appender {
			appender.LogWrite(*rec)
		}

		LogMetrics.AppendCount++
	} else {
		LogMetrics.IgnoreCount++
	}
}

/******* Logging *******/
// Send a formatted log message internally
func (c Logger) intLogf(lvl Level, rec *LogRecord) {

	// append for rootLogger
	c.appendLog(lvl, rec, LoggerManager.rootLogger)

	// append for logger
	if LoggerManager.loggerMap != nil && len(LoggerManager.loggerMap) > 0 {
		for name, loggerAppenderReference := range LoggerManager.loggerMap {
			if strings.Index(name, c.Name) == 0 {
				c.appendLog(lvl, rec, loggerAppenderReference)
			}
		}
	}

}
