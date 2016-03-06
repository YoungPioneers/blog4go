// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"testing"
)

func TestConfigValidation(t *testing.T) {
	config := new(Config)

	if !config.valid() {
		t.Error("config validation failed.")
	}
}
