package examples

import (
	"fmt"
	"github.com/gohutool/log4go"
	"testing"
	"time"
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

	//log4go.LoggerManager.SetDebug(true)
	log4go.LoggerAppenderFactory.RegistryType("sample", &SampleLoggerAppender{})
}

type SampleLoggerAppender struct {
	pattern string
}

func (cla SampleLoggerAppender) Init(pattern string, property []log4go.AppenderProperty) error {
	return nil
}

func (cla SampleLoggerAppender) Start() error {
	return nil
}

// LogWrite
//This will be called to log a LogRecord message.
func (cla SampleLoggerAppender) LogWrite(rec log4go.LogRecord) error {
	fmt.Printf("[Sample] %+v\n", rec)
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

func TestLogger(t *testing.T) {
	log4go.LoggerManager.InitWithDefaultConfig()
	logger.Info("hello")
	logger.Info("hello")
	time.Sleep(1 * time.Second)
	logger.Error("hello")

	time.Sleep(10 * time.Second)
}

var logger = log4go.LoggerManager.GetLogger("com.hello")

func TestLoggerExample(t *testing.T) {
	log4go.LoggerManager.InitWithXML("./example.xml")
	logger.Info("hello")
	logger.Info("hello")
	logger.Info("hello")
	logger.Info("hello")
	time.Sleep(1 * time.Second)
	logger.Error("hello")

	time.Sleep(3 * time.Second)

	select {}
}
