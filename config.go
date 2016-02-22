// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

const (
	// TypeTimeBaseRotate is time base logrotate tag
	TypeTimeBaseRotate = "time"
	// TypeSizeBaseRotate is size base logrotate tag
	TypeSizeBaseRotate = "size"
)

// Config struct define the config struct used for file wirter
type Config struct {
	//XMLName  xml.Name `xml:"blog4go"`
	Filters  []filter `xml:"filter"`
	MinLevel string   `xml:"minlevel,attr"`
}

// log filter
type filter struct {
	XMLName    xml.Name   `xml:"filter"`
	Levels     string     `xml:"levels,attr"`
	Colored    bool       `xml:"colored,attr"`
	File       file       `xml:"file"`
	RotateFile rotateFile `xml:"rotatefile"`
	//Console    console    `xml:"console"`
	Socket socket `xml:"socket"`
}

type file struct {
	XMLName xml.Name `xml:"file"`
	Path    string   `xml:"path,attr"`
}

type rotateFile struct {
	XMLName     xml.Name `xml:"rotatefile"`
	Path        string   `xml:"path,attr"`
	Type        string   `xml:"type,attr"`
	RotateLines int      `xml:"rotateLines,attr"`
	RotateSize  ByteSize `xml:"rotateSize,attr"`
}

//type console struct {
//XMLName xml.Name `xml:"console"`
//}

type socket struct {
	XMLName xml.Name `xml:"socket"`
	Network string   `xml:"network,attr"`
	Address string   `xml:"address,attr"`
}

// check if config is valid
func (config *Config) valid() bool {
	return true
}

// read config from a xml file
func readConfig(fileName string) (*Config, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	in, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err

	}

	config := new(Config)
	err = xml.Unmarshal(in, config)
	if err != nil {
		return nil, err
	}

	return config, err
}
