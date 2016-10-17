// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"testing"
)

func TestConfigValidation(t *testing.T) {
	config := new(Config)

	// min level test
	config.MinLevel = "something"
	if err := config.valid(); ErrConfigBadAttributes != err {
		t.Errorf("config minlevel validation failed. MinLevel: %s", config.MinLevel)
	}

	config.MinLevel = "debug"
	if err := config.valid(); ErrConfigFiltersNotFound != err {
		t.Error("config filter length check failed.")
	}

	// levels test
	f := filter{
		File: file{
			Path: "/tmp/test.log",
		},
	}
	config.Filters = make([]filter, 0)
	config.Filters = append(config.Filters, f)

	if err := config.valid(); ErrConfigLevelsNotFound != err {
		t.Error("config file levels check failed.")
	}

	// filter check
	f = filter{
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

	// rotate file check
	// file path check
	f = filter{
		Levels: "debug",
		RotateFile: rotateFile{
			Type: "time",
			Path: "",
		},
	}
	config.Filters = make([]filter, 0)
	config.Filters = append(config.Filters, f)

	if err := config.valid(); ErrConfigFilePathNotFound != err {
		t.Error("config rotate file filter check failed.")
	}

	// rotate type check
	f = filter{
		Levels: "debug",
		RotateFile: rotateFile{
			Type: "",
			Path: "/tmp/test.log",
		},
	}
	config.Filters = make([]filter, 0)
	config.Filters = append(config.Filters, f)

	if err := config.valid(); ErrConfigFileRotateTypeNotFound != err {
		t.Error("config rotate file filter check failed.")
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

	// socket check
	// address check
	f = filter{
		Levels: "debug",
		Socket: socket{
			Network: "udp",
			Address: "",
		},
	}
	config.Filters = make([]filter, 0)
	config.Filters = append(config.Filters, f)

	if err := config.valid(); ErrConfigSocketAddressNotFound != err {
		t.Error("config socket filter check failed.")
	}

	// network check
	f = filter{
		Levels: "debug",
		Socket: socket{
			Network: "",
			Address: "127.0.0.1:4567",
		},
	}
	config.Filters = make([]filter, 0)
	config.Filters = append(config.Filters, f)

	if err := config.valid(); ErrConfigSocketNetworkNotFound != err {
		t.Error("config socket filter check failed.")
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
