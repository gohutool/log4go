package log4go

import (
	"encoding/json"
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

type AppenderProperty struct {
	Name  string `xml:"name,attr" json:"name"`
	Value string `xml:",chardata" json:"value"`
}

type AppenderConfig struct {
	Enabled  string             `xml:"enabled,attr" json:"enabled"`
	Name     string             `xml:"name,attr" json:"name"`
	Type     string             `xml:"type" json:"type"`
	Pattern  string             `xml:"pattern" json:"pattern"`
	Property []AppenderProperty `xml:"property" json:"properties"`
}

type AppenderRef struct {
	Ref string `xml:"ref,attr" json:"appender"`
}

type RootLoggerConfig struct {
	Level    string        `xml:"level" json:"level"`
	Appender []AppenderRef `xml:"appender-ref" json:"appender-refs"`
}

type LoggerConfig struct {
	Name     string        `xml:"name,attr" json:"name"`
	Level    string        `xml:"level" json:"level"`
	Appender []AppenderRef `xml:"appender-ref" json:"appender-refs"`
}

type LoggerConfiguration struct {
	Appender []AppenderConfig `xml:"appender" json:"appenders"`
	Root     RootLoggerConfig `xml:"root" json:"root"`
	Logger   []LoggerConfig   `xml:"logger" json:"loggers"`
}

func LoadXML(content string) LoggerConfiguration {
	b := []byte(content)

	xc := new(LoggerConfiguration)

	if err := xml.Unmarshal(b, xc); err != nil {
		msg := fmt.Sprintf("LoadConfiguration: Error: Could not parse XML configuration in %q: %s\n", content, err)
		fmt.Fprintln(os.Stderr, msg)
		panic(msg)
	}

	return *xc
}

func LoadXMLConfigurationProperties(filename string) LoggerConfiguration {
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

	xc := new(LoggerConfiguration)
	if err := xml.Unmarshal(contents, xc); err != nil {
		msg := fmt.Sprintf("LoadConfiguration: Error: Could not parse XML configuration in %q: %s\n", filename, err)
		fmt.Fprintln(os.Stderr, msg)
		panic(msg)
	}

	return *xc
}

func LoadJson(content string) LoggerConfiguration {
	b := []byte(content)

	xc := new(LoggerConfiguration)

	if err := json.Unmarshal(b, xc); err != nil {
		msg := fmt.Sprintf("LoadConfiguration: Error: Could not parse Json configuration in %q: %s\n", content, err)
		fmt.Fprintln(os.Stderr, msg)
		panic(msg)
	}

	return *xc
}

func LoadJsonConfigurationProperties(filename string) LoggerConfiguration {
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

	xc := new(LoggerConfiguration)
	if err := json.Unmarshal(contents, xc); err != nil {
		msg := fmt.Sprintf("LoadConfiguration: Error: Could not parse Json configuration in %q: %s\n", filename, err)
		fmt.Fprintln(os.Stderr, msg)
		panic(msg)
	}

	return *xc
}
