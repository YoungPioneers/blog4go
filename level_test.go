// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"testing"
)

func TestLevelValidation(t *testing.T) {
	if LevelType(-1).valid() {
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

	if " level=\"DEBUG\" " != DEBUG.prefix() {
		t.Error("DEBUG Level to wrong prefix string format.")
	}

	if " level=\"TRACE\" " != TRACE.prefix() {
		t.Error("TRACE Level to wrong prefix string format.")
	}

	if " level=\"INFO\" " != INFO.prefix() {
		t.Error("INFO Level to wrong prefix string format.")
	}

	if " level=\"WARN\" " != WARNING.prefix() {
		t.Error("WARN Level to wrong prefix string format.")
	}

	if " level=\"ERROR\" " != ERROR.prefix() {
		t.Error("ERROR Level to wrong prefix string format.")
	}

	if " level=\"CRITICAL\" " != CRITICAL.prefix() {
		t.Error("CRITICAL Level to wrong prefix string format.")
	}

	if "UNKNOWN" != LevelType(-1).String() {
		t.Error("Wrong Level to wrong string format.")
	}

	initPrefix(true)

	if " level=\"\x1b[37mTRACE\x1b[0m\" " != TRACE.prefix() {
		t.Error("TRACE Level with color to wrong prefix string format.")
	}

	if " level=\"\x1b[32mDEBUG\x1b[0m\" " != DEBUG.prefix() {
		t.Error("DEBUG Level with color to wrong prefix string format.")
	}

	if " level=\"\x1b[34mINFO\x1b[0m\" " != INFO.prefix() {
		t.Error("INFO Level with color to wrong prefix string format.")
	}

	if " level=\"\x1b[33mWARN\x1b[0m\" " != WARNING.prefix() {
		t.Error("WARN Level with color to wrong prefix string format.")
	}

	if " level=\"\x1b[31mERROR\x1b[0m\" " != ERROR.prefix() {
		t.Error("ERROR Level with color to wrong prefix string format.")
	}

	if " level=\"\x1b[31mCRITICAL\x1b[0m\" " != CRITICAL.prefix() {
		t.Error("CRITICAL Level with color to wrong prefix string format.")
	}
}

func TestStringToLevel(t *testing.T) {
	str := "debug"
	if DEBUG != LevelFromString(str) {
		t.Errorf("String to level failed. str: %s", str)
	}

	str = "something"
	if LevelFromString(str).valid() {
		t.Errorf("String to level invalid. str: %s", str)
	}

	str = ""
	if LevelFromString(str).valid() {
		t.Error("Empty string to level invalid.")
	}
}
