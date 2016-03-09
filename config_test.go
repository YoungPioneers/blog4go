// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"testing"
)

func TestConfigValidation(t *testing.T) {
	config := new(Config)

	config.MinLevel = "something"
	if err := config.valid(); ErrConfigBadAttributes != err {
		t.Errorf("config minlevel validation failed. MinLevel: %s", config.MinLevel)
	}

	config.MinLevel = "debug"
	if err := config.valid(); ErrConfigFiltersNotFound != err {
		t.Error("config filter length check failed.")
	}

	f := filter{
		Levels: "debug",
		File: file{
			Path: "/tmp/test.log",
		},
	}
	config.Filters = make([]filter, 0)
	config.Filters = append(config.Filters, f)

	if err := config.valid(); ErrConfigLevelsNotFound == err || ErrConfigFilePathNotFound == err {
		t.Error("config file filter check failed.")
	}

	f = filter{
		Levels: "debug",
		RotateFile: rotateFile{
			Type: "time",
			Path: "/tmp/test.log",
		},
	}
	config.Filters = make([]filter, 0)
	config.Filters = append(config.Filters, f)

	if err := config.valid(); ErrConfigLevelsNotFound == err || ErrConfigFilePathNotFound == err || ErrConfigFileRotateTypeNotFound == err {
		t.Errorf("config rotate file filter check failed. err: %+v", config.Filters)
	}

	f = filter{
		Levels: "debug",
		Socket: socket{
			Network: "udp",
			Address: "127.0.0.1:4567",
		},
	}
	config.Filters = make([]filter, 0)
	config.Filters = append(config.Filters, f)

	if err := config.valid(); ErrConfigLevelsNotFound == err || ErrConfigSocketAddressNotFound == err || ErrConfigSocketNetworkNotFound == err {
		t.Error("config socket filter check failed.")
	}
}
