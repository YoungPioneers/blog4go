// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
)

const (
	// TypeTimeBaseRotate is time base logrotate tag
	TypeTimeBaseRotate = "time"
	// TypeSizeBaseRotate is size base logrotate tag
	TypeSizeBaseRotate = "size"
)

var (
	// ErrConfigFiltersNotFound not found filters
	ErrConfigFiltersNotFound = errors.New("Please define at least one filter")
	// ErrConfigBadAttributes wrong attribute
	ErrConfigBadAttributes = errors.New("Bad attributes setting")
	// ErrConfigLevelsNotFound not found levels
	ErrConfigLevelsNotFound = errors.New("Please define levels attribution")
	// ErrConfigFilePathNotFound not found file path
	ErrConfigFilePathNotFound = errors.New("Please define the file path")
	// ErrConfigFileRotateTypeNotFound not found rotate type
	ErrConfigFileRotateTypeNotFound = errors.New("Please define the file rotate type")
	// ErrConfigSocketAddressNotFound not found socket address
	ErrConfigSocketAddressNotFound = errors.New("Please define a socket address")
	// ErrConfigSocketNetworkNotFound not found socket port
	ErrConfigSocketNetworkNotFound = errors.New("Please define a socket network type")
)

// Config struct define the config struct used for file wirter
type Config struct {
	Filters  []filter `xml:"filter"`
	MinLevel string   `xml:"minlevel,attr"`
}

// log filter
type filter struct {
	Levels     string     `xml:"levels,attr"`
	Colored    bool       `xml:"colored,attr"`
	File       file       `xml:"file"`
	RotateFile rotateFile `xml:"rotatefile"`
	Console    console    `xml:"console"`
	Socket     socket     `xml:"socket"`
}

type file struct {
	Path string `xml:"path,attr"`
}

type rotateFile struct {
	Path        string `xml:"path,attr"`
	Type        string `xml:"type,attr"`
	RotateLines int    `xml:"rotateLines,attr"`
	RotateSize  int64  `xml:"rotateSize,attr"`
	Retentions  int64  `xml:"retentions,attr"`
}

type console struct {
	// redirect stderr to stdout
	Redirect bool `xml:"redirect"`
}

type socket struct {
	Network string `xml:"network,attr"`
	Address string `xml:"address,attr"`
}

// check if config is valid
func (config *Config) valid() error {
	// check minlevel validation
	if "" != config.MinLevel && !LevelFromString(config.MinLevel).valid() {
		return ErrConfigBadAttributes
	}

	// check filters len
	if len(config.Filters) < 1 {
		return ErrConfigFiltersNotFound
	}

	// check filter one by one
	for _, filter := range config.Filters {
		if "" == filter.Levels {
			return ErrConfigLevelsNotFound
		}

		if (file{}) != filter.File {
			// seem not needed now
			//if "" == filter.File.Path {
			//return ErrConfigFilePathNotFound
			//}
		} else if (rotateFile{}) != filter.RotateFile {
			if "" == filter.RotateFile.Path {
				return ErrConfigFilePathNotFound
			}

			if "" == filter.RotateFile.Type {
				return ErrConfigFileRotateTypeNotFound
			}
		} else if (socket{}) != filter.Socket {
			if "" == filter.Socket.Address {
				return ErrConfigSocketAddressNotFound
			}

			if "" == filter.Socket.Network {
				return ErrConfigSocketNetworkNotFound
			}
		}
	}

	return nil
}

// read config from a xml file
func readConfig(fileName string) (*Config, error) {
	file, err := os.Open(fileName)
	if nil != err {
		return nil, err
	}
	defer file.Close()

	in, err := ioutil.ReadAll(file)
	if nil != err {
		return nil, err
	}

	config := new(Config)
	err = xml.Unmarshal(in, config)
	if nil != err {
		return nil, err
	}

	return config, err
}
