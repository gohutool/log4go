package log4go

import (
	"encoding/json"
	"net"
	"strconv"
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
	LoggerAppenderFactory.RegistryType("socket", &SocketLoggerAppender{})
}

type SocketLoggerAppender struct {
	pattern  string
	endpoint string
	protocol string
	tryTimes int
	retry    int
	check    bool
	interval int

	rec    chan LogRecord
	reConn chan bool

	conn net.Conn
	quit bool
}

func (sla *SocketLoggerAppender) Init(pattern string, property []AppenderProperty) error {
	LoggerManager.debug("[Log4go]SocketLoggerAppender init with ", property, pattern)

	dFormat := "[%T %D %m] [%L][%l] (%S) %M"

	if len(pattern) > 0 {
		sla.pattern = pattern
	} else {
		sla.pattern = dFormat
	}

	sla.rec = make(chan LogRecord, LogBufferLength)
	sla.reConn = make(chan bool)

	endpoint := ""
	protocol := "udp"
	check := false
	tryTimes := 0
	retry := 0
	interval := 30

	for _, prop := range property {
		switch strings.ToLower(prop.Name) {
		case "endpoint":
			endpoint = strings.Trim(prop.Value, " \r\n")
		case "protocol":
			protocol = strings.Trim(prop.Value, " \r\n")
		case "check":
			check = strings.ToLower(strings.TrimSpace(prop.Value)) == "true"
		case "interval":
			if len(strings.TrimSpace(prop.Value)) > 0 {
				v, err := strconv.Atoi(strings.TrimSpace(prop.Value))
				if err != nil {
					panic("interval must be a integer")
				} else {
					interval = v
				}
			}
		case "retry":
			///retry = strings.ToLower(strings.TrimSpace(prop.Value)) != "false"
			if len(strings.TrimSpace(prop.Value)) > 0 {
				try1, err := strconv.Atoi(strings.TrimSpace(prop.Value))
				if err != nil {
					panic("try must be a integer")
				} else {
					retry = try1
				}
			}
		}
	}

	sla.endpoint = endpoint
	sla.retry = retry
	sla.check = check
	sla.tryTimes = tryTimes
	sla.protocol = protocol
	sla.interval = interval

	return nil
}

func (sla *SocketLoggerAppender) Start() error {
	LoggerManager.debug("[Log4go]Start SocketLoggerAppender")

	sla.initConn()

	go func() {
		defer func() {
			if sla.conn != nil && sla.protocol == "tcp" {
				sla.conn.Close()
			}
		}()

		for !sla.quit {
			select {
			case rec, ok := <-sla.rec:
				if !ok {
					continue
				}
				// 是否需要进行通道失败后的缓存，如果缓存担心，恢复过程中由于，chan的阻塞导致，外部的协程被阻塞住了所以失败后的，选择进行忽略
				// 此处还可以失败后，进行缓存，恢复后，缓存在发送
				// 还可以进行压缩算法，多次一起进行压缩，然后在发送， 目前就基本实现
				if sla.conn == nil {
					continue
				}

				// Marshall into JSON
				js, err := json.Marshal(rec)

				if err != nil {
					LoggerManager.error(BuildFormatString("SocketLoggerAppender(%q/%s): %s", sla.endpoint,
						sla.protocol, err))
					continue
				}

				_, err = sla.conn.Write(js)

				if err != nil {
					LoggerManager.error(BuildFormatString("SocketLoggerAppender(%q/%s): %s", sla.endpoint,
						sla.protocol, err))
					sla.conn = nil
					sla.reConn <- true
				}
			case <-sla.reConn:
				sla.initConn()
			default:
				time.Sleep(10 * time.Microsecond)
			}
		}
	}()

	return nil
}

func (sla *SocketLoggerAppender) initConn() {
	sla.tryTimes = sla.tryTimes + 1

	sock, err := net.Dial(sla.protocol, sla.endpoint)

	if err != nil {
		if sla.check {
			panic(BuildFormatString("SocketLoggerAppender(%s/%s) %s ", sla.endpoint,
				sla.protocol, err))
		}

		go func() {
			for sock == nil && (sla.tryTimes < sla.retry || sla.retry < 0) {

				time.Sleep(time.Duration(sla.interval) * time.Second)

				sock, err = net.Dial(sla.protocol, sla.endpoint)
				/*if sla.check {
					LoggerManager.warning(BuildFormatString("SocketLoggerAppender(%s/%s) %s ", sla.endpoint,
						sla.protocol, err))
				}*/

			}

			sla.conn = sock
			sla.tryTimes = 0
		}()

	} else {
		sla.conn = sock
		sla.tryTimes = 0
	}

}

// LogWrite
//This will be called to log a LogRecord message.
func (sla *SocketLoggerAppender) LogWrite(rec LogRecord) error {
	//fmt.Printf("[Socket] %+v\n", rec)
	sla.rec <- rec
	return nil
}

// Close
// This should clean up anything lingering about the LogWriter, as it is called before
// the LogWriter is removed.  LogWrite should not be called after Close.
func (sla *SocketLoggerAppender) Close() error {
	LoggerManager.debug("[Log4go]Close SocketLoggerAppender")

	sla.quit = true

	if sla.conn != nil {
		sla.conn.Close()
	}

	if sla.rec != nil {
		close(sla.rec)
	}

	return nil
}
