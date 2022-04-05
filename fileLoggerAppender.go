package log4go

import (
	"fmt"
	"os"
)

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
	LoggerAppenderFactory.RegistryType("file", &FileLoggerAppender{})
}

type FileLoggerAppender struct {
	pattern string
	// The opened file
	filename string
	file     *os.File

	// The logging format
	format string

	// File header/trailer
	header, trailer string

	// Rotate at linecount
	maxlines          int
	maxlines_curlines int

	// Rotate at size
	maxsize         int
	maxsize_cursize int

	// Rotate daily
	daily          bool
	daily_opendate int

	// Keep old logfiles (.001, .002, etc)
	rotate    bool
	maxbackup int
}

func (cla *FileLoggerAppender) Init(pattern string, property []AppenderProperty) error {
	return nil
}

func (cla *FileLoggerAppender) Start() error {
	return nil
}

// LogWrite
//This will be called to log a LogRecord message.
func (cla *FileLoggerAppender) LogWrite(rec LogRecord) error {
	fmt.Printf("[File] %+v\n", rec)
	return nil
}

// Close
// This should clean up anything lingering about the LogWriter, as it is called before
// the LogWriter is removed.  LogWrite should not be called after Close.
func (cla *FileLoggerAppender) Close() error {
	return nil
}
