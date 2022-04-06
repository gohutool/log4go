package log4go

import (
	"fmt"
	"os"
	"strings"
	"time"
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
	rec chan LogRecord
	rot chan bool

	pattern string
	// The opened file
	filename string
	file     *os.File

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

func (fla *FileLoggerAppender) Init(pattern string, property []AppenderProperty) error {
	file := ""
	format := "[%T %D %m] [%L][%l] (%S) %M"
	maxlines := 0
	maxsize := 0
	daily := false
	rotate := false

	if len(pattern) > 0 {
		fla.pattern = pattern
	} else {
		fla.pattern = format
	}

	// Parse properties
	for _, prop := range property {
		switch prop.Name {
		case "filename":
			file = strings.Trim(prop.Value, " \r\n")
		case "maxlines":
			maxlines = strToNumSuffix(strings.Trim(prop.Value, " \r\n"), 1000)
		case "maxsize":
			maxsize = strToNumSuffix(strings.Trim(prop.Value, " \r\n"), 1024)
		case "daily":
			daily = strings.Trim(prop.Value, " \r\n") != "false"
		case "rotate":
			rotate = strings.Trim(prop.Value, " \r\n") != "false"
		default:
			LoggerManager.error(BuildFormatString("LoadConfiguration: Warning: Unknown property \"%s\"", prop.Name))
		}
	}

	// Check properties
	if len(file) == 0 {
		panic("LoadConfiguration: Error: Required property \"filename\" missing")
	}

	fla.rec = make(chan LogRecord, LogBufferLength)
	fla.rot = make(chan bool)
	fla.filename = file
	fla.pattern = format
	fla.rotate = rotate

	fla.maxbackup = 999
	fla.maxlines = maxlines
	fla.maxsize = maxsize
	fla.daily = daily

	// open the file for the first time
	if err := fla.intRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", fla.filename, err)
		panic(fmt.Sprintf("FileLogWriter(%q): %s\n", fla.filename, err))
	}

	return nil
}

func (fla *FileLoggerAppender) Start() error {
	LoggerManager.debug("[Log4go]Start FileLoggerAppender")

	go func() {
		defer func() {
			if fla.file != nil {
				fmt.Fprint(fla.file, getDefaultPatternConvert().FormatLogRecord(fla.trailer, LogRecord{Created: time.Now()}))
				fla.file.Close()
			}
		}()

		for {
			select {
			case <-fla.rot:
				if err := fla.intRotate(); err != nil {
					LoggerManager.error(fmt.Sprintf("FileLogWriter(%q): %s\n", fla.filename, err))
					return
				}
			case rec, ok := <-fla.rec:
				if !ok {
					return
				}
				now := time.Now()
				if (fla.maxlines > 0 && fla.maxlines_curlines >= fla.maxlines) ||
					(fla.maxsize > 0 && fla.maxsize_cursize >= fla.maxsize) ||
					(fla.daily && now.Day() != fla.daily_opendate) {
					// if with Rotate to put bool to rot chan, it will be raised a bug when concurrent environment
					if err := fla.intRotate(); err != nil {
						LoggerManager.error(fmt.Sprintf("FileLogWriter(%q): %s\n", fla.filename, err))
						return
					}
				}

				// Perform the write
				n, err := fmt.Fprint(fla.file, getDefaultPatternConvert().FormatLogRecord(fla.pattern, rec))
				if err != nil {
					LoggerManager.error(fmt.Sprintf("FileLogWriter(%q): %s\n", fla.filename, err))
					return
				}

				// Update the counts
				fla.maxlines_curlines++
				fla.maxsize_cursize += n
			}
		}

		if error := fla.Close(); error != nil {
			LoggerManager.debug("fileAppender("+fla.file.Name()+") error and close now.", error)
		}
	}()

	return nil
}

// LogWrite
//This will be called to log a LogRecord message.
func (fla *FileLoggerAppender) LogWrite(rec LogRecord) error {
	fla.rec <- rec
	return nil
}

// Close
// This should clean up anything lingering about the LogWriter, as it is called before
// the LogWriter is removed.  LogWrite should not be called after Close.
func (fla *FileLoggerAppender) Close() error {

	LoggerManager.debug("[Log4go]Close FileLoggerAppender")

	close(fla.rec)
	fla.file.Sync()

	return nil
}

// Rotate
// Request that the logs rotate
func (fla *FileLoggerAppender) Rotate() {
	fla.rot <- true
}

// If this is called in a threaded context, it MUST be synchronized
func (fla *FileLoggerAppender) intRotate() error {
	// Close any log file that may be open
	if fla.file != nil {
		fmt.Fprint(fla.file, getDefaultPatternConvert().FormatLogRecord(fla.trailer, LogRecord{Created: time.Now()}))
		fla.file.Close()
	}

	// If we are keeping log files, move it to the next available number
	if fla.rotate {
		_, err := os.Lstat(fla.filename)
		if err == nil { // file exists
			// Find the next available number
			num := 1
			fname := ""
			if fla.daily && time.Now().Day() != fla.daily_opendate {
				yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

				for ; err == nil && num <= fla.maxbackup; num++ {
					fname = fla.filename + fmt.Sprintf(".%s.%03d", yesterday, num)
					_, err = os.Lstat(fname)
				}
				// return error if the last file checked still existed
				if err == nil {
					return fmt.Errorf("Rotate: Cannot find free log number to rename %s\n", fla.filename)
				}
			} else {
				num = fla.maxbackup - 1
				for ; num >= 1; num-- {
					fname = fla.filename + fmt.Sprintf(".%d", num)
					nfname := fla.filename + fmt.Sprintf(".%d", num+1)
					_, err = os.Lstat(fname)
					if err == nil {
						os.Rename(fname, nfname)
					}
				}
			}

			fla.file.Close()
			// Rename the file to its newfound home
			err = os.Rename(fla.filename, fname)
			if err != nil {
				return fmt.Errorf("Rotate: %s\n", err)
			}
		}
	}

	// Open the log file
	fd, err := os.OpenFile(fla.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		panic(BuildString("open FileAppender ", fla.filename, " ", err.Error()))
	}

	fla.file = fd

	now := time.Now()
	fmt.Fprint(fla.file, getDefaultPatternConvert().FormatLogRecord(fla.header, LogRecord{Created: now}))

	// Set the daily open date to the current date
	fla.daily_opendate = now.Day()

	// initialize rotation values
	fla.maxlines_curlines = 0
	fla.maxsize_cursize = 0

	return nil
}
