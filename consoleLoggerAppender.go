package log4go

import "fmt"

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : consoleLoggerAppender.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/5 13:06
* 修改历史 : 1. [2022/4/5 13:06] 创建文件 by NST
*/

func init() {
	LoggerAppenderFactory.RegistryType("console", &ConsoleLoggerAppender{})
}

type ConsoleLoggerAppender struct {
	pattern string
}

func (cla *ConsoleLoggerAppender) Init(pattern string, property []AppenderProperty) error {

	fmt.Println("ConsoleLoggerAppender init with ", property, pattern)

	return nil
}

func (cla *ConsoleLoggerAppender) Start() error {
	return nil
}

// LogWrite
//This will be called to log a LogRecord message.
func (cla *ConsoleLoggerAppender) LogWrite(rec *LogRecord) error {
	return nil
}

// Close
// This should clean up anything lingering about the LogWriter, as it is called before
// the LogWriter is removed.  LogWrite should not be called after Close.
func (cla *ConsoleLoggerAppender) Close() error {
	return nil
}
