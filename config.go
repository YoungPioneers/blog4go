// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

// Config struct define the config struct used for file wirter
type Config struct {
	XMLName  xml.Name `xml:"blog4go"`
	Filters  []filter `xml:"filter"`
	MinLevel string   `xml:"minlevel,attr"`
}

type filter struct {
	XMLName    xml.Name   `xml:"filter"`
	Levels     string     `xml:"levels,attr"`
	Colored    bool       `xml:"colored,attr"`
	File       file       `xml:"file"`
	RotateFile rotateFile `xml:"rotatefile"`
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

	var config Config
	err = xml.Unmarshal(in, &config)
	if err != nil {
		return nil, err
	}

	return &config, err
}
