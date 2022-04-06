package log4go

import (
	"fmt"
	"strings"
)

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
type LevelText string

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

const (
	FINEST_TXT   LevelText = "FINEST"
	FINE_TXT     LevelText = "FINE"
	DEBUG_TXT    LevelText = "DEBUG"
	TRACE_TXT    LevelText = "TRACE"
	INFO_TXT     LevelText = "INFO"
	WARNING_TXT  LevelText = "WARNING"
	ERROR_TXT    LevelText = "ERROR"
	CRITICAL_TXT LevelText = "CRITICAL"
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

func (l Level) levelText() LevelText {
	switch int(l) {
	case int(FINEST):
		return FINEST_TXT
	case int(FINE):
		return FINE_TXT
	case int(DEBUG):
		return DEBUG_TXT
	case int(TRACE):
		return TRACE_TXT
	case int(INFO):
		return INFO_TXT
	case int(WARNING):
		return WARNING_TXT
	case int(ERROR):
		return ERROR_TXT
	case int(CRITICAL):
		return CRITICAL_TXT
	}

	return INFO_TXT
}

func (l LevelText) String() string {
	return string(l)
}

func str2LevelText(level string) Level {
	switch strings.ToUpper(level) {
	case "FINEST":
		return FINEST
	case "FINE":
		return FINE
	case "DEBUG":
		return DEBUG
	case "TRACE":
		return TRACE
	case "INFO":
		return INFO
	case "WARNING":
		return WARNING
	case "ERROR":
		return ERROR
	case "CRITICAL":
		return CRITICAL
	}

	return INFO
}

func (l LevelText) Level() Level {
	return str2LevelText(l.String())
	/*switch strings.ToUpper(l.String()) {
	case "FINEST":
		return FINEST
	case "FINE":
		return FINE
	case "DEBUG":
		return DEBUG
	case "TRACE":
		return TRACE
	case "INFO":
		return INFO
	case "WARNING":
		return WARNING
	case "ERROR":
		return ERROR
	case "CRITICAL":
		return CRITICAL
	}

	return INFO_TXT*/
}

var LoggerManager = &loggerManager{}

type loggerManager struct {
	laf        *loggerAppenderFactory
	rootLogger LoggerAppenderReference
	loggerMap  map[string]LoggerAppenderReference
}

type LoggerAppenderReference struct {
	level    Level
	appender []LoggerAppender
}

func (lm *loggerManager) GetLogger(name string) Logger {
	logger := Logger{Name: name}
	return logger
}

func (lm *loggerManager) InitWithDefaultConfig() error {
	content := `<?xml version="1.0" encoding="UTF-8"?>
<configuration>
  <appender enabled="true" name="console">
    <type>console</type>
    <pattern>[%D %T] [%L] (%S) %M</pattern>
    <!-- level is (:?FINEST|FINE|DEBUG|TRACE|INFO|WARNING|ERROR) -->
  </appender>
  
  <!-- 输出级别是info级别及以上的日志，下面的ref关联的两个appender没有filter设置，所以，info及以上的日志都是会输出到这2个appender的 -->
  <root>
    <level>info</level>
    <appender-ref ref="console" />
  </root>

</configuration>
`
	lc := LoadXML(content)
	return lm.InitWithConfig(lc)
}

func (lm *loggerManager) InitWithConfig(configuration LoggerConfiguration) error {
	lm.laf = LoggerAppenderFactory.new()
	lm.loggerMap = make(map[string]LoggerAppenderReference)
	lm.rootLogger = LoggerAppenderReference{}

	enabledMap := make(map[string]bool)

	if configuration.Appender != nil && len(configuration.Appender) > 0 {
		// init all of appender at first
		for _, appender := range configuration.Appender {
			enabledMap[strings.ToLower(appender.Name)] = strings.TrimSpace(strings.ToLower(appender.Enabled)) != "false"
			_, _, e := lm.laf.registerLoggerAppender(appender.Name, appender.Type, appender.Pattern, appender.Property)
			if e != nil {
				panic(fmt.Sprintf("Add Appender[%v] %v", appender.Name, e.Error()))
			}
		}
	}

	// init root
	rc := configuration.Root
	if len(rc.Level) == 0 {
		rc.Level = INFO_TXT.String()
	}

	lm.rootLogger.level = str2LevelText(rc.Level)
	lm.rootLogger.appender = make([]LoggerAppender, 0, 10)

	if rc.Appender != nil && len(rc.Appender) > 0 {
		for _, appenderRef := range rc.Appender {

			enable, ok := enabledMap[strings.ToLower(appenderRef.Ref)]

			if ok && enable {
				ap, e := lm.laf.getAppenderRefByName(appenderRef.Ref)
				if e == nil {
					lm.rootLogger.appender = append(lm.rootLogger.appender, ap)
				} else {
					fmt.Println("[Warning] ", e.Error())
				}
			}

		}
	}
	// init root end

	// init logger
	if configuration.Logger != nil && len(configuration.Logger) > 0 {
		for _, logger := range configuration.Logger {
			oneLogger := LoggerAppenderReference{}
			oneLogger.level = str2LevelText(logger.Level)
			oneLogger.appender = make([]LoggerAppender, 0, 10)

			for _, appenderRef := range logger.Appender {
				enable, ok := enabledMap[strings.ToLower(appenderRef.Ref)]

				if ok && enable {
					ap, e := lm.laf.getAppenderRefByName(appenderRef.Ref)
					if e == nil {
						oneLogger.appender = append(oneLogger.appender, ap)
					} else {
						fmt.Println("[Warning] ", e.Error())
					}
				}
			}

			if len(oneLogger.appender) > 0 {
				lm.loggerMap[logger.Name] = oneLogger
			}
		}
	}
	// init logger end

	return nil
}

func (lm *loggerManager) InitWithXML(xmlFile string) error {
	lc := LoadXMLConfigurationProperties(xmlFile)
	return lm.InitWithConfig(lc)
}

func (lm *loggerManager) InitWithJson(jsonFile string) error {
	lc := LoadJsonConfigurationProperties(jsonFile)
	return lm.InitWithConfig(lc)
}

func (lm *loggerManager) InitWithJsonConfig(jsonContent string) error {
	lc := LoadJson(jsonContent)
	return lm.InitWithConfig(lc)
}
