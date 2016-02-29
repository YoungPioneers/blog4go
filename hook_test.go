// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"testing"
	"time"
)

type MyHook struct {
	cnt     int
	level   Level
	message string
}

func (hook *MyHook) add() {
	hook.cnt++
}

func (hook *MyHook) Fire(level Level, message string) {
	hook.add()
	hook.level = level
	hook.message = message
}

func TestHook(t *testing.T) {
	hook := new(MyHook)
	hook.cnt = 0

	err := NewFileWriter("/tmp", false)
	defer blog.Close()
	if nil != err {
		t.Errorf("initialize file writer faied. err: %s", err.Error())
	}

	blog.SetHook(hook)
	blog.SetHookLevel(INFO)

	blog.Trace("something")
	// wait for hook called
	time.Sleep(10 * time.Millisecond)
	if 0 != hook.cnt {
		t.Error("hook called not valid")
	}

	if TRACE == hook.level || "something" == hook.message {
		t.Errorf("hook parameters wrong. level: %d, message: %s", hook.level, hook.message)
	}

	blog.Info("yes")
	// wait for hook called
	time.Sleep(10 * time.Millisecond)
	if 1 != hook.cnt {
		t.Error("hook not called")
	}

	if INFO != hook.level || "yes" != hook.message {
		t.Errorf("hook parameters wrong. level: %d, message: %s", hook.level, hook.message)
	}
}
