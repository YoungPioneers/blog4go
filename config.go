// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"strings"
)

const (
	// TypeTimeBaseRotate is time base logrotate tag
	TypeTimeBaseRotate = "time"
	// TypeSizeBaseRotate is size base logrotate tag
	TypeSizeBaseRotate = "size"
)

// Config struct define the config struct used for file wirter
type Config struct {
	XMLName  xml.Name `xml:"blog4go"`
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
	Console    console    `xml:"console"`
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

type console struct {
	XMLName xml.Name `xml:"console"`
}

type socket struct {
	XMLName xml.Name `xml:"socket"`
	Network string   `xml:"network,attr"`
	Address string   `xml"address,attr"`
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

// NewWriterFromConfigAsFile initialize a writer according to given config file
// configFile must be the path to the config file
func NewWriterFromConfigAsFile(configFile string) (err error) {
	singltonLock.Lock()
	defer singltonLock.Unlock()
	if nil != blog {
		return
	}

	// read config from file
	config, err := readConfig(configFile)
	if nil != err {
		return
	}

	fileWriter := new(FileWriter)

	fileWriter.level = DEBUG
	if level := LevelFromString(config.MinLevel); level.valid() {

		fileWriter.level = level
	}
	fileWriter.closed = false
	fileWriter.writers = make(map[Level]*baseFileWriter)

	for _, filter := range config.Filters {
		var rotate = false
		// get file path
		var filePath string
		if nil != &filter.File && "" != filter.File.Path {
			// single file
			filePath = filter.File.Path
			rotate = false
		} else if nil != &filter.RotateFile && "" != filter.RotateFile.Path {
			// multi files
			filePath = filter.RotateFile.Path
			rotate = true
		} else if nil != &filter.Console {
			// console writer

		} else {
			// config error
			return ErrFilePathNotFound
		}

		// init a base file writer
		writer, err := newBaseFileWriter(filePath)
		if nil != err {
			return err
		}

		levels := strings.Split(filter.Levels, ",")
		for _, levelStr := range levels {
			var level Level
			if level = LevelFromString(levelStr); !level.valid() {
				return ErrInvalidLevel
			}

			if rotate {
				// set logrotate strategy
				switch filter.RotateFile.Type {
				case TypeTimeBaseRotate:
					writer.SetTimeRotated(true)
				case TypeSizeBaseRotate:
					writer.SetRotateSize(filter.RotateFile.RotateSize)
					writer.SetRotateLines(filter.RotateFile.RotateLines)
				default:
					return ErrInvalidRotateType
				}
			}

			// set color
			fileWriter.SetColored(filter.Colored)
			fileWriter.writers[level] = writer
		}
	}

	// log hook
	fileWriter.hook = nil
	fileWriter.hookLevel = DEBUG

	blog = fileWriter
	return
}
