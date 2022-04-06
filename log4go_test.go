package log4go

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"testing"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : log4go_test.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/4 22:36
* 修改历史 : 1. [2022/4/4 22:36] 创建文件 by NST
*/

func TestConfiguration(t *testing.T) {
	xc := LoadXMLConfigurationProperties("./examples/example.xml")
	b, _ := json.Marshal(xc)
	fmt.Println(string(b))
	fmt.Printf("%+v", xc)
}

func TestJsonConfiguration(t *testing.T) {
	xc := LoadJsonConfigurationProperties("./examples/example.json")
	b, _ := xml.Marshal(xc)
	fmt.Println(string(b))
	fmt.Printf("%+v", xc)
}

func TestLoggerAppenderFactoryGetType(t *testing.T) {
	fmt.Println(LoggerAppenderFactory.getInterfaceByType("Console"))

	LoggerAppenderFactory.registerLoggerAppender("default", "console", "", nil)
	LoggerAppenderFactory.registerLoggerAppender("default", "console", "", nil)
}

func TestLoadXMLConfiguration(t *testing.T) {
	xc := LoggerManager.InitWithXML("./examples/example.xml")
	fmt.Printf("%+v\n", xc)
}

func TestLoadDefaultConfiguration(t *testing.T) {
	xc := LoggerManager.InitWithDefaultConfig()
	fmt.Printf("%+v\n", xc)

	xc = LoggerManager.InitWithXML("./examples/example.xml")
	fmt.Printf("%+v\n", xc)
}
