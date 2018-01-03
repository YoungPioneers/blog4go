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
	hook := NewMyHook()

	SetHook(hook)
	SetHookLevel(INFO)

	Debug("something")
	Debugf("%s", "something")
	// wait for hook called
	time.Sleep(1 * time.Millisecond)
	if 0 != hook.cnt {
		t.Error("hook called not valid")
	}

	if DEBUG == hook.Level() || "something" == hook.Message() {
		t.Errorf("hook parameters wrong. level: %s, message: %s", hook.level.String(), hook.Message())
	}

	// async
	Info("yes")
	// wait for hook called
	time.Sleep(1 * time.Millisecond)
	if 1 != hook.Cnt() {
		t.Error("hook not called")
	}

	if INFO != hook.Level() || "yes" != hook.Message() {
		t.Errorf("hook parameters wrong. level: %d, message: %s", hook.level, hook.Message())
	}

	// sync
	SetHookAsync(false)
	Warn("warn")
	if 2 != hook.Cnt() {
		t.Error("hook not called")
	}

	if WARNING != hook.Level() || "warn" != hook.Message() {
		t.Errorf("hook parameters wrong. level: %d, message: %s", hook.level, hook.Message())
	}

	// test basic operations
	SetTags(map[string]string{"tagName": "tagValue"})
	Tags()

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

	SetHookAsync(true)
	Colored()
	SetColored(true)
	TimeRotated()
	SetTimeRotated(true)
	Level()
	SetLevel(CRITICAL)
	Retentions()
	SetRetentions(0)
	SetRetentions(7)
	RotateLines()
	SetRotateLines(0)
	SetRotateLines(100000)
	RotateSize()
	SetRotateSize(0)
	SetRotateSize(1024 * 1024 * 500)

	Debug("Debug", 1)
	Debugf("%s\\", "Debug")
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

	SetBufferSize(0)
}
