# log4go
a logkit like as log4j with go language

# Introduce
This project is come from an unmaintained fork, the branch can see http://log4go.googlecode.com

- left only so it doesn't break imports.
- some enhancement if to need 


Usage:
- Add the following import:
import l4g "github.com/gohutool/log4go"
- log4go configuration sample

./examples/example.xml

```
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
    <property name="rotate">false</property> <!-- true enables log rotation, otherwise append -->
    <property name="maxsize">0M</property> <!-- \d+[KMG]? Suffixes are in terms of 2**10 -->
    <property name="maxlines">0K</property> <!-- \d+[KMG]? Suffixes are in terms of thousands -->
    <property name="daily">true</property> <!-- Automatically rotates when a log message is written after midnight -->
  </appender>
  <appender enabled="true" name="testfile">
    <type>file</type>
    <pattern>[%D %T] [%L] (%S) %M</pattern>
    <property name="filename">trace.xml</property>
    <property name="rotate">false</property> <!-- true enables log rotation, otherwise append -->
    <property name="maxsize">100M</property> <!-- \d+[KMG]? Suffixes are in terms of 2**10 -->
    <property name="maxrecords">6K</property> <!-- \d+[KMG]? Suffixes are in terms of thousands -->
    <property name="daily">false</property> <!-- Automatically rotates when a log message is written after midnight -->
  </appender>
  <appender enabled="false" name="logstash"><!-- enabled=false means this logger won't actually be created -->
    <type>socket</type>
    <pattern>[%D %T] [%L] (%S) %M</pattern>
    <property name="endpoint">192.168.1.255:12124</property> <!-- recommend UDP broadcast -->
    <property name="protocol">udp</property> <!-- tcp or udp -->
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

```
LoggerManager.InitWithXML("./examples/example.xml")
var logger = log4go.LoggerManager.GetLogger("com.hello")
logger.Info("hello")
logger.Info("hello")
```

Acknowledgements:
- pomack
  For providing awesome patches to bring log4go up to the latest Go spec