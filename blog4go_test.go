// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"os/exec"
	"testing"
	"time"
)

func TestGlobalOperation(t *testing.T) {
	err := NewFileWriter("/tmp", false)
	defer func() {
		Close()

		// clean logs
		_, err = exec.Command("/bin/sh", "-c", "/bin/rm /tmp/*.log*").Output()
		if nil != err {
			t.Errorf("clean files failed. err: %s", err.Error())
		}
	}()

	// test file writer hook
	hook := new(MyHook)
	hook.cnt = 0

	SetHook(hook)
	SetHookLevel(INFO)

	Debug("something")
	// wait for hook called
	time.Sleep(1 * time.Millisecond)
	if 0 != hook.cnt {
		t.Error("hook called not valid")
	}

	if DEBUG == hook.level || "something" == hook.message {
		t.Errorf("hook parameters wrong. level: %s, message: %s", hook.level.String(), hook.message)
	}

	Info("yes")
	// wait for hook called
	time.Sleep(1 * time.Millisecond)
	if 1 != hook.cnt {
		t.Error("hook not called")
	}

	if INFO != hook.level || "yes" != hook.message {
		t.Errorf("hook parameters wrong. level: %d, message: %s", hook.level, hook.message)
	}

	// test basic operations
	Debug("Debug", 1)
	Debugf("%s", "Debug")
	Trace("Trace", 2)
	Tracef("%s", "Trace")
	Info("Info", 3)
	Infof("%s", "Info")
	Warn("Warn", 4)
	Warnf("%s", "Warn")
	Error("Error", 5)
	Errorf("%s", "Error")
	Critical("Critical", 6)
	Criticalf("%s", "Critical")
	Flush()

	SetColored(true)
	SetColored(true)
	SetTimeRotated(true)
	SetLevel(CRITICAL)
	SetRetentions(0)
	SetRetentions(7)
	SetRotateLines(0)
	SetRotateLines(100000)
	SetRotateSize(0)
	SetRotateSize(1024 * 1024 * 500)

	Debug("Debug", 1)
	Debugf("%s", "Debug")
	Trace("Trace", 2)
	Tracef("%s", "Trace")
	Info("Info", 3)
	Infof("%s", "Info")
	Warn("Warn", 4)
	Warnf("%s", "Warn")
	Error("Error", 5)
	Errorf("%s", "Error")
	Critical("Critical", 6)
	Criticalf("%s", "Critical")
}
