package log4go

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : loggerAppender.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/5 10:56
* 修改历史 : 1. [2022/4/5 10:56] 创建文件 by NST
*/

var (
	defaultAppender = map[string]interface{}{
		"file":    nil,
		"socket":  nil,
		"console": nil,
	}
)

var LoggerAppenderFactory = &loggerAppenderFactory{}

/****** LogWriter ******/
// This is an interface for anything that should be able to write logs

type LoggerAppender interface {
	Init(pattern string, property []Property) error

	Start() error

	// LogWrite
	//This will be called to log a LogRecord message.
	LogWrite(rec *LogRecord) error

	// Close
	// This should clean up anything lingering about the LogWriter, as it is called before
	// the LogWriter is removed.  LogWrite should not be called after Close.
	Close() error
}

type loggerAppenderFactory struct {
	appender     map[string]LoggerAppender
	appenderType map[string]any
	typeLock     sync.RWMutex
}

func (laf *loggerAppenderFactory) new() *loggerAppenderFactory {
	fmt.Println("New LoggerAppenderFactory")

	laf.typeLock.RLock()
	defer laf.typeLock.RUnlock()

	for _, la := range laf.appender {
		la.Close()
	}

	fmt.Println("Close LoggerAppenderFactory's loggerAppender")

	fmt.Println("init LoggerAppenderFactory's loggerAppender")
	laf.appender = make(map[string]LoggerAppender)
	laf.appenderType = make(map[string]any)
	return laf
}

func (laf *loggerAppenderFactory) RegistryType(typename string, typeClz any) *loggerAppenderFactory {
	if laf.appender == nil {
		laf = laf.new()
	}
	//fmt.Println("registry one type ", typename, " ", reflect.ValueOf(typeClz))

	if app, ok := typeClz.(LoggerAppender); ok {
		fmt.Println("registry one type ", typename, " ", reflect.TypeOf(app))
	} else {
		panic(reflect.ValueOf(typeClz).String() + " is not a implementation of LoggerAppender")
	}

	laf.typeLock.RLock()
	defer laf.typeLock.RUnlock()

	laf.appenderType[strings.ToLower(typename)] = typeClz

	return laf
}

func (laf *loggerAppenderFactory) getInterfaceByType(typename string) (any, error) {
	if laf.appenderType == nil {
		laf = laf.new()
	}

	rtn, ok := laf.appenderType[strings.ToLower(typename)]

	if ok == true {
		return rtn, nil
	} else {
		return nil, errors.New("Not found " + typename)
	}
}

func (laf *loggerAppenderFactory) LoggerAppender(name string) (LoggerAppender, error) {
	return nil, nil
}

func (laf *loggerAppenderFactory) DefaultLoggerAppender() (LoggerAppender, error) {
	return laf.LoggerAppender("default")
}

func (laf *loggerAppenderFactory) registerLoggerAppender(name, typename, pattern string,
	properties []AppenderProperty) (*loggerAppenderFactory, interface{}, error) {
	laf.typeLock.RLock()
	defer laf.typeLock.RUnlock()

	itf, err := laf.getInterfaceByType(typename)

	if err != nil {
		return laf, nil, err
	}
	fmt.Println("interface : ", itf)
	return laf, nil, nil
}
