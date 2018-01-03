// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"testing"
	"time"
)

func TestConsoleWriterBasicOperation(t *testing.T) {
	err := NewConsoleWriter(false)
	defer Close()
	if nil != err {
		t.Error(err.Error())
	}

	// test console writer hook
	hook := NewMyHook()

	blog.SetHook(hook)
	blog.SetHookLevel(INFO)

	blog.Debug("something")
	blog.Debugf("%s", "something")
	// wait for hook called
	time.Sleep(1 * time.Millisecond)
	if 0 != hook.Cnt() {
		t.Error("hook called not valid")
	}

	if DEBUG == hook.Level() || "something" == hook.Message() {
		t.Errorf("hook parameters wrong. level: %s, message: %s", hook.level.String(), hook.Message())
	}

	// async
	blog.Info("yes")
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
	// wait for hook called
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
	blog.SetRetentions(7)
	blog.RotateLines()
	blog.SetRotateLines(100000)
	blog.RotateSize()
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

	blog.Close()
	time.Sleep(1 * time.Second)
	blog.Debug("Debug", 1)
	blog.Debugf("%s", "Debug")
}

func TestRedirectedConsoleWriterBasicOperation(t *testing.T) {
	err := NewConsoleWriter(true)
	defer Close()
	if nil != err {
		t.Error(err.Error())
	}

	// test console writer hook
	hook := NewMyHook()

	blog.SetHook(hook)
	blog.SetHookLevel(INFO)

	blog.Debug("something")
	blog.Debugf("%s", "something")
	// wait for hook called
	time.Sleep(1 * time.Millisecond)
	if 0 != hook.Cnt() {
		t.Error("hook called not valid")
	}

	if DEBUG == hook.Level() || "something" == hook.Message() {
		t.Errorf("hook parameters wrong. level: %s, message: %s", hook.level.String(), hook.Message())
	}

	// async
	blog.Info("yes")
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
	// wait for hook called
	if 2 != hook.Cnt() {
		t.Error("hook not called")
	}

	if WARNING != hook.Level() || "warn" != hook.Message() {
		t.Errorf("hook parameters wrong. level: %d, message: %s", hook.Level(), hook.Message())
	}

	// test basic operations
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

	blog.Colored()
	blog.SetColored(true)
	blog.TimeRotated()
	blog.SetTimeRotated(true)
	blog.Level()
	blog.SetLevel(CRITICAL)
	blog.Retentions()
	blog.SetRetentions(7)
	blog.RotateLines()
	blog.SetRotateLines(100000)
	blog.RotateSize()
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

	blog.Close()
	time.Sleep(1 * time.Second)
	blog.Debug("Debug", 1)
	blog.Debugf("%s", "Debug")
}

func TestSingleConsoleWriter(t *testing.T) {
	err := NewConsoleWriter(true)
	defer Close()
	if nil != err {
		t.Error(err.Error())
	}

	// duplicate init check
	err = NewConsoleWriter(true)
	defer Close()
	if ErrAlreadyInit != err {
		t.Error("duplicate init check fail")
	}
}

func BenchmarkConsoleWriter(b *testing.B) {
	err := NewConsoleWriter(true)
	defer Close()
	if nil != err {
		b.Error(err.Error())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		blog.Debugf("haha %s. en\\en, always %d and %f", "eddie", 18, 3.1415)
	}
}
