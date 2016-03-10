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
	SetColored(true)
	SetTimeRotated(true)
	SetLevel(INFO)
	SetRetentions(7)
	SetRotateLines(100000)
	SetRotateSize(ByteSize(1024 * 1024 * 500))

	Debug("Debug")
	Debugf("%s", "Debug")
	Trace("Trace")
	Tracef("%s", "Trace")
	Info("Info")
	Infof("%s", "Info")
	Warn("Warn")
	Warnf("%s", "Warn")
	Error("Error")
	Errorf("%s", "Error")
	Critical("Critical")
	Criticalf("%s", "Critical")
}
