// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"os/exec"
	"testing"
	"time"
)

func TestBaseFileWriterBasicOperation(t *testing.T) {
	err := NewBaseFileWriter("/tmp/mylog.log", true)
	if nil != err {
		t.Errorf("Failed when initializing base file writer. err: %s", err.Error())
	}
	defer func() {
		Close()

		// clean logs
		_, err = exec.Command("/bin/sh", "-c", "/bin/rm /tmp/*.log*").Output()
		if nil != err {
			t.Errorf("clean files failed. err: %s", err.Error())
		}
	}()

	// duplicate init
	err = NewBaseFileWriter("/tmp/mylog.log", true)
	if ErrAlreadyInit != err {
		t.Errorf("Duplicate initialization check failed. err: %s", err.Error())
	}

	// test file writer hook
	hook := NewMyHook()

	blog.SetHook(hook)
	blog.SetHookLevel(INFO)

	blog.Debug("something")
	blog.Debugf("%s", "something")
	// sync
	if 0 != hook.Cnt() {
		t.Error("hook called not valid")
	}

	if DEBUG == hook.Level() || "something" == hook.Message() {
		t.Errorf("hook parameters wrong. level: %s, message: %s", hook.Level().String(), hook.Message())
	}

	blog.SetHookAsync(true)
	blog.Info("yes")
	// async
	// wait for hook called
	time.Sleep(1 * time.Millisecond)
	if 1 != hook.Cnt() {
		t.Error("hook not called")
	}

	if INFO != hook.Level() || "yes" != hook.Message() {
		t.Errorf("hook parameters wrong. level: %d, message: %s", hook.Level(), hook.Message())
	}

	// sync
	blog.SetHookAsync(false)
	blog.Warn("warn")
	if 2 != hook.Cnt() {
		t.Error("hook not called")
	}

	if WARNING != hook.Level() || "warn" != hook.Message() {
		t.Errorf("hook parameters wrong. level: %d, message: %s", hook.Level(), hook.Message())
	}

	// test basic operations
	blog.SetTags(map[string]string{"tagName": "tagValue"})
	blog.Tags()

	blog.Debug("Debug", 1)
	blog.Debugf("%s", "Debug")
	blog.Trace("Trace", 2)
	blog.Tracef("%s", "Trace")
	blog.Info("Info", 3)
	blog.Infof("%s", "Info")
	blog.Warn("Warn", 4)
	blog.Warnf("%s", "Warn")
	blog.Error("Error", 5)
	blog.Errorf("%s", "Error")
	blog.Critical("Critical", 6)
	blog.Criticalf("%s", "Critical")
	blog.flush()

	blog.SetHookAsync(true)
	blog.Colored()
	blog.SetColored(true)
	blog.TimeRotated()
	blog.SetTimeRotated(true)
	blog.Level()
	blog.SetLevel(CRITICAL)
	blog.Retentions()
	blog.SetRetentions(0)
	blog.SetRetentions(7)
	blog.RotateLines()
	blog.SetRotateLines(0)
	blog.SetRotateLines(100000)
	blog.RotateSize()
	blog.SetRotateSize(0)
	blog.SetRotateSize(1024 * 1024 * 500)

	blog.Debug("Debug", 1)
	blog.Debugf("%s\\", "Debug")
	blog.Trace("Trace", 2)
	blog.Tracef("%s", "Trace")
	blog.Info("Info", 3)
	blog.Infof("%s", "Info")
	blog.Warn("Warn", 4)
	blog.Warnf("%s", "Warn")
	blog.Error("Error", 5)
	blog.Errorf("%s", "Error")
	blog.Critical("Critical", 6)
	blog.Criticalf("%s", "Critical")

	// wait for timeRotate run
	time.Sleep(1 * time.Second)

	blog.Close()
	time.Sleep(1 * time.Second)
	blog.Debug("Debug", 1)
	blog.Debugf("%s", "Debug")
}
