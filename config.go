package log4go

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : config.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/4 22:08
* 修改历史 : 1. [2022/4/4 22:08] 创建文件 by NST
*/

type xmlProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

type xmlAppenderProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

type xmlAppender struct {
	Enabled  string        `xml:"enabled,attr"`
	Name     string        `xml:"name,attr"`
	Type     string        `xml:"type"`
	Pattern  string        `xml:"pattern"`
	Property []xmlProperty `xml:"property"`
}

type xmlAppenderRef struct {
	Ref string `xml:"ref,attr"`
}

type xmlRoot struct {
	Level    string           `xml:"level"`
	Appender []xmlAppenderRef `xml:"appender-ref"`
}

type xmlLogger struct {
	Name     string           `xml:"name,attr"`
	Level    string           `xml:"level"`
	Appender []xmlAppenderRef `xml:"appender-ref"`
}

type xmlLoggerConfig struct {
	Appender []xmlAppender `xml:"appender"`
	Root     xmlRoot       `xml:"root"`
	Logger   []xmlLogger   `xml:"logger"`
}

func loadConfigurationProperties(filename string) xmlLoggerConfig {
	fd, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not open %q for reading: %s\n", filename, err)
		os.Exit(1)
	}
	contents, err := ioutil.ReadAll(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not read %q: %s\n", filename, err)
		os.Exit(1)
	}

	xc := new(xmlLoggerConfig)
	if err := xml.Unmarshal(contents, xc); err != nil {
		msg := fmt.Sprintf("LoadConfiguration: Error: Could not parse XML configuration in %q: %s\n", filename, err)
		fmt.Fprintln(os.Stderr, msg)
		panic(msg)
	}

	return *xc
}
