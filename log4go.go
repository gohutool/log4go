package log4go

import "time"

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : log4go.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/4 22:04
* 修改历史 : 1. [2022/4/4 22:04] 创建文件 by NST
*/

const (
	L4G_VERSION = "log4go-v1.0.1"
	L4G_MAJOR   = 1
	L4G_MINOR   = 0
	L4G_BUILD   = 1
)

type Level int

const (
	FINEST Level = iota
	FINE
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	CRITICAL
)

var (
	levelStrings = [...]string{"FNST", "FINE", "DEBG", "TRAC", "INFO", "WARN", "EROR", "CRIT"}
)

func (l Level) String() string {
	if l < 0 || int(l) > len(levelStrings) {
		return "UNKNOWN"
	}
	return levelStrings[int(l)]
}

// A LogRecord contains all the pertinent information for each message
type LogRecord struct {
	Level   Level     // The log level
	Created time.Time // The time at which the log message was created (nanoseconds)
	Source  string    // The message source
	Message string    // The log message
}

type Logger struct {
}

func (c Logger) Critical(arg0 interface{}, args ...interface{}) error {
	return nil
}

func (c Logger) Error(arg0 interface{}, args ...interface{}) error {
	return nil
}

func (c Logger) Warning(arg0 interface{}, args ...interface{}) error {
	return nil
}

func (c Logger) Info(arg0 interface{}, args ...interface{}) error {
	return nil
}

func (c Logger) Trace(arg0 interface{}, args ...interface{}) error {
	return nil
}

func (c Logger) Debug(arg0 interface{}, args ...interface{}) error {
	return nil
}

func (c Logger) Fine(arg0 interface{}, args ...interface{}) error {
	return nil
}

func (c Logger) Finest(arg0 interface{}, args ...interface{}) error {
	return nil
}
