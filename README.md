# log4go
a logkit like as log4j with go language

# Introduce
This project is come from an unmaintained fork, the branch can see http://log4go.googlecode.com

- left only so it doesn't break imports.
- some enhancement if to need 

Usage:
- Add log4go with the following import
```
import "github.com/gohutool/log4go"
```
- Support Xml configuration
```
./examples/example.xml

<?xml version="1.0" encoding="UTF-8"?>
<configuration>
  <appender enabled="true" name="console">
    <type>console</type>
    <pattern>[%D %T] [%L] (%S) %M</pattern>
    <!-- level is (:?FINEST|FINE|DEBUG|TRACE|INFO|WARNING|ERROR) -->
  </appender>
  <appender enabled="true" name="file">
    <type>file</type>
    <pattern>[%D %T] [%L] (%S) %M</pattern>
    <property name="filename">test.log</property>
    <!--
       %T - Time (15:04:05 MST)
       %t - Time (15:04)
       %D - Date (2006/01/02)
       %d - Date (01/02/06)
       %L - Level (FNST, FINE, DEBG, TRAC, WARN, EROR, CRIT)
       %S - Source
       %M - Message
       It ignores unknown format strings (and removes them)
       Recommended: "[%D %T] [%L] (%S) %M"
    -->
    <property name="rotate">false</property> 
    <!-- true enables log rotation, otherwise append -->
    <property name="maxsize">0M</property> 
    <!-- \d+[KMG]? Suffixes are in terms of 2**10 -->
    <property name="maxlines">0K</property> 
    <!-- \d+[KMG]? Suffixes are in terms of thousands -->
    <property name="daily">true</property> 
    <!-- Automatically rotates when a log message is written after midnight -->
  </appender>
  <appender enabled="true" name="testfile">
    <type>file</type>
    <pattern>[%D %T] [%L] (%S) %M</pattern>
    <property name="filename">trace.xml</property>
    <property name="rotate">false</property> 
    <!-- true enables log rotation, otherwise append -->
    <property name="maxsize">100M</property> 
    <!-- \d+[KMG]? Suffixes are in terms of 2**10 -->
    <property name="maxrecords">6K</property> 
    <!-- \d+[KMG]? Suffixes are in terms of thousands -->
    <property name="daily">false</property> 
    <!-- Automatically rotates when a log message is written after midnight -->
  </appender>
  <!-- enabled=false means this logger won't actually be created -->
  <appender enabled="false" name="logstash">
    <type>socket</type>
    <pattern>[%D %T] [%L] (%S) %M</pattern>
    <property name="endpoint">192.168.1.255:12124</property> 
    <!-- recommend UDP broadcast -->
    <property name="protocol">udp</property> 
    <!-- tcp or udp -->
  </appender>
  
  <!-- 这个就是自定义的一个logger -->
  <!-- 输出级别是info级别及以上的日志，不要怕，StreamOperateFile已经过滤，只输出info级别的日志 -->
  <!-- additivity 这个默认是true，即继承父类 root logger -->
  <!-- 也就是说，你的这个日志也会在root的logger里面输出的，我这里配置false，就是不继承，各走各的。 -->
  <!-- appender-ref 也就是说这个logger的输出目的地是哪里，ref就是关联到上面声明的一个file，一个console -->
  <logger name="com.ginghan">
    <level>info</level>
    <appender-ref ref="file" />
    <appender-ref ref="console" />
  </logger>

  <logger name="com.hello">
    <level>info</level>
    <appender-ref ref="console" />
    <appender-ref ref="file" />
  </logger>

  <!-- 输出级别是info级别及以上的日志，下面的ref关联的两个appender没有filter设置，
  所以，info及以上的日志都是会输出到这2个appender的 -->
  <root>
    <level>info</level>
    <appender-ref ref="testfile" />
    <appender-ref ref="logstash" />
  </root>

</configuration>
```
- Support Json configuration
```
{
	"appenders": [
		{
			"enabled": "true",
			"name": "console",
			"type": "console",
			"pattern": "[%D %T] [%L] (%S) %M",
			"properties": null
		},
		{
			"enabled": "true",
			"name": "file",
			"type": "file",
			"pattern": "[%D %T] [%L] (%S) %M",
			"properties": [
				{
					"name": "filename",
					"value": "test.log"
				},
				{
					"name": "rotate",
					"value": "false"
				},
				{
					"name": "maxsize",
					"value": "0M"
				},
				{
					"name": "maxlines",
					"value": "0K"
				},
				{
					"name": "daily",
					"value": "true"
				}
			]
		},
		{
			"enabled": "true",
			"name": "testfile",
			"type": "file",
			"pattern": "[%D %T] [%L] (%S) %M",
			"properties": [
				{
					"name": "filename",
					"value": "trace.xml"
				},
				{
					"name": "rotate",
					"value": "false"
				},
				{
					"name": "maxsize",
					"value": "100M"
				},
				{
					"name": "maxrecords",
					"value": "6K"
				},
				{
					"name": "daily",
					"value": "false"
				}
			]
		},
		{
			"enabled": "false",
			"name": "logstash",
			"type": "socket",
			"pattern": "[%D %T] [%L] (%S) %M",
			"properties": [
				{
					"name": "endpoint",
					"value": "192.168.1.255:12124"
				},
				{
					"name": "protocol",
					"value": "udp"
				}
			]
		}
	],
	"root": {
		"level": "info",
		"appender-refs": [
			{
				"appender": "testfile"
			},
			{
				"appender": "logstash"
			}
		]
	},
	"loggers": [
		{
			"name": "Sample",
			"level": "info",
			"appender-refs": [
				{
					"appender": "file"
				},
				{
					"appender": "console"
				}
			]
		},
		{
			"name": "Demo",
			"level": "info",
			"appender-refs": [
				{
					"appender": "console"
				},
				{
					"appender": "file"
				}
			]
		}
	]
}
```
- log4go sample
```

import (
	"github.com/gohutool/log4go"
	"testing"
	"time"
)

var logger = log4go.LoggerManager.GetLogger("com.hello")

func TestLoggerExample(t *testing.T) {
	log4go.LoggerManager.InitWithXML("./example.xml")
	logger.Info("hello")
	logger.Info("hello")
	time.Sleep(1 * time.Second)
	logger.Error("hello")

	time.Sleep(3 * time.Second)
}
```

- sample output

```
=== RUN   TestLoggerExample
[2022/04/06 15:17:59 CST 299] [INFO][com.hello] (github.com/gohutool/log4go/examples.TestLoggerExample:73) hello
[2022/04/06 15:17:59 CST 299] [INFO][com.hello] (github.com/gohutool/log4go/examples.TestLoggerExample:72) hello
[2022/04/06 15:18:00 CST 299] [EROR][com.hello] (github.com/gohutool/log4go/examples.TestLoggerExample:75) hello
--- PASS: TestLoggerExample (4.00s)
PASS

Debugger finished with the exit code 0
```
- FileLoggerAppender output

![File log output snapshot](misc/filelog-snapshot.png "File log output snapshot")

Acknowledgements:
- pomack
  For providing awesome patches to bring log4go up to the latest Go spec