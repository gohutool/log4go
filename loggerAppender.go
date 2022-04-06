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
	Init(pattern string, property []AppenderProperty) error

	Start() error

	// LogWrite
	//This will be called to log a LogRecord message.
	LogWrite(rec LogRecord) error

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
	LoggerManager.debug("New LoggerAppenderFactory")

	laf.typeLock.RLock()
	defer laf.typeLock.RUnlock()

	if laf.appender != nil && len(laf.appender) > 0 {
		for _, la := range laf.appender {
			la.Close()
		}
	}

	laf.appender = make(map[string]LoggerAppender)

	if laf.appenderType == nil {
		laf.appenderType = make(map[string]any)

		LoggerManager.debug("init LoggerAppenderFactory's loggerAppender")
	}

	return laf
}

func (laf *loggerAppenderFactory) RegistryType(typename string, typeClz any) *loggerAppenderFactory {

	if reflect.TypeOf(typeClz).Kind() != reflect.Ptr {
		panic(reflect.TypeOf(typeClz).String() + " is not a pointer of LoggerAppender")
	}

	laf.typeLock.RLock()
	defer laf.typeLock.RUnlock()

	if laf.appenderType == nil {
		laf.appenderType = make(map[string]any)
	}
	//fmt.Println("registry one type ", typename, " ", reflect.ValueOf(typeClz))

	if app, ok := typeClz.(LoggerAppender); ok {
		LoggerManager.debug("registry one type ", typename, " ", reflect.TypeOf(app))
	} else {
		panic(reflect.ValueOf(typeClz).String() + " is not a implementation of LoggerAppender")
	}

	laf.appenderType[strings.ToLower(typename)] = typeClz

	return laf
}

func (laf *loggerAppenderFactory) getInterfaceByType(typename string) (any, error) {

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

func (laf *loggerAppenderFactory) getAppenderRefByName(name string) (LoggerAppender, error) {
	name = strings.ToLower(name)
	if laf.appender == nil {
		return nil, errors.New("LoggerAppenderFactory is not init")
	}

	if ap, ok := laf.appender[name]; !ok {
		panic("LoggerAppender named '" + name + "' is not found.")
	} else {
		return ap, nil
	}
}

func (laf *loggerAppenderFactory) registerLoggerAppender(name, typename, pattern string,
	properties []AppenderProperty) (*loggerAppenderFactory, interface{}, error) {
	laf.typeLock.RLock()
	defer laf.typeLock.RUnlock()

	name = strings.ToLower(name)

	if laf.appender == nil {
		laf.appender = make(map[string]LoggerAppender)
	}

	if laf.appenderType == nil || len(laf.appenderType) == 0 {
		panic("AppenderType is empty, please RegistryType at first")
	}

	itf, err := laf.getInterfaceByType(typename)

	if err != nil {
		panic("LoggerAppender type '" + typename + "' is not registry.")
	}

	if _, ok := laf.appender[name]; ok {
		panic("LoggerAppender named '" + name + "' is registry already.")
	}

	newObj := reflect.New(reflect.TypeOf(itf).Elem()).Interface().(LoggerAppender)
	newObj.Init(pattern, properties)
	newObj.Start()

	laf.appender[name] = newObj

	LoggerManager.debug(fmt.Sprintf("LoggerAppender(%v)[%v] is register with %q", reflect.TypeOf(newObj), &newObj, name))

	//LoggerManager.debug(fmt.Sprintf("%v ========== %v \n", itf, newObj))
	LoggerManager.debug(fmt.Sprintf("%v ========== %v ", &itf, &newObj))
	//LoggerManager.debug(fmt.Sprintf("%v ========== %v \n", reflect.TypeOf(&itf), reflect.TypeOf(&newObj)))
	LoggerManager.debug(fmt.Sprintf("%v ========== %v ", reflect.TypeOf(itf), reflect.TypeOf(newObj)))

	//LoggerManager.debug("interface : ", itf)
	return laf, nil, nil
}
