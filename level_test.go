// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"testing"
)

func TestLevelValidation(t *testing.T) {
	if Level(-1).valid() {
		t.Errorf("Level Validation Failed. level: %d", -1)
	}

	if !DEBUG.valid() {
		t.Error("DEBUG Level Validation Failed.")
	}
}
func TestLevelStringFormat(t *testing.T) {
	if "DEBUG" != DEBUG.String() {
		t.Error("DEBUG Level to wrong string format.")
	}

	if "TRACE" != TRACE.String() {
		t.Error("TRACE Level to wrong string format.")
	}

	if "INFO" != INFO.String() {
		t.Error("INFO Level to wrong string format.")
	}

	if "WARN" != WARNING.String() {
		t.Error("WARN Level to wrong string format.")
	}

	if "ERROR" != ERROR.String() {
		t.Error("ERROR Level to wrong string format.")
	}

	if "CRITICAL" != CRITICAL.String() {
		t.Error("CRITICAL Level to wrong string format.")
	}

	if " [CRITICAL] " != CRITICAL.Prefix() {
		t.Error("CRITICAL Level to wrong prefix string format.")
	}

	if "UNKNOWN" != Level(-1).String() {
		t.Error("Wrong Level to wrong string format.")
	}
}
