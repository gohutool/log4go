package examples

import (
	"fmt"
	"log4go"
	"testing"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : configuration_test.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/4 22:35
* 修改历史 : 1. [2022/4/4 22:35] 创建文件 by NST
*/

func init() {
	log4go.LoggerAppenderFactory.RegistryType("sample", &SampleLoggerAppender{})
}

type SampleLoggerAppender struct {
	pattern string
}

func (cla SampleLoggerAppender) Init(pattern string, property []log4go.Property) error {
	return nil
}

func (cla SampleLoggerAppender) Start() error {
	return nil
}

// LogWrite
//This will be called to log a LogRecord message.
func (cla SampleLoggerAppender) LogWrite(rec *log4go.LogRecord) error {
	return nil
}

// Close
// This should clean up anything lingering about the LogWriter, as it is called before
// the LogWriter is removed.  LogWrite should not be called after Close.
func (cla SampleLoggerAppender) Close() error {
	return nil
}

func TestInitExample(t *testing.T) {
	fmt.Println(log4go.LoggerAppenderFactory.LoggerAppender(""))
}
