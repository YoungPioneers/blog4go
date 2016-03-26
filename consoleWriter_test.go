// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"testing"
	"time"
)

func TestConsoleWriterBasicOperation(t *testing.T) {
	err := NewConsoleWriter()
	defer Close()
	if nil != err {
		t.Error(err.Error())
	}

	// test console writer hook
	hook := new(MyHook)
	hook.cnt = 0

	blog.SetHook(hook)
	blog.SetHookLevel(INFO)

	blog.Debug("something")
	// wait for hook called
	time.Sleep(1 * time.Millisecond)
	if 0 != hook.cnt {
		t.Error("hook called not valid")
	}

	if DEBUG == hook.level || "something" == hook.message {
		t.Errorf("hook parameters wrong. level: %s, message: %s", hook.level.String(), hook.message)
	}

	blog.Info("yes")
	// wait for hook called
	time.Sleep(1 * time.Millisecond)
	if 1 != hook.cnt {
		t.Error("hook not called")
	}

	if INFO != hook.level || "yes" != hook.message {
		t.Errorf("hook parameters wrong. level: %d, message: %s", hook.level, hook.message)
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

	blog.SetColored(true)
	blog.SetTimeRotated(true)
	blog.SetLevel(CRITICAL)
	blog.SetRetentions(7)
	blog.SetRotateLines(100000)
	blog.SetRotateSize(1024 * 1024 * 500)

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
}

func TestSingleConsoleWriter(t *testing.T) {
	err := NewConsoleWriter()
	defer Close()
	if nil != err {
		t.Error(err.Error())
	}

	// duplicate init check
	err = NewConsoleWriter()
	defer Close()
	if ErrAlreadyInit != err {
		t.Error("duplicate init check fail")
	}
}

func BenchmarkConsoleWriter(b *testing.B) {
	err := NewConsoleWriter()
	defer Close()
	if nil != err {
		b.Error(err.Error())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		blog.Debugf("haha %s. en\\en, always %d and %f", "eddie", 18, 3.1415)
	}
}
