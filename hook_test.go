// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"fmt"
	"os/exec"
	"sync"
	"testing"
	"time"
)

type MyHook struct {
	cnt     int
	level   LevelType
	message string

	l *sync.RWMutex
}

func NewMyHook() (hook *MyHook) {
	hook = new(MyHook)
	hook.cnt = 0
	hook.level = TRACE
	hook.message = ""
	hook.l = new(sync.RWMutex)

	return
}

func (hook *MyHook) Add() {
	hook.l.Lock()
	defer hook.l.Unlock()
	hook.cnt++
}

func (hook *MyHook) Cnt() int {
	hook.l.RLock()
	defer hook.l.RUnlock()
	return hook.cnt
}

func (hook *MyHook) Level() LevelType {
	hook.l.RLock()
	defer hook.l.RUnlock()
	return hook.level
}

func (hook *MyHook) SetLevel(level LevelType) {
	hook.l.Lock()
	defer hook.l.Unlock()
	hook.level = level
}

func (hook *MyHook) Message() string {
	hook.l.RLock()
	defer hook.l.RUnlock()
	return hook.message
}

func (hook *MyHook) SetMessage(message string) {
	hook.l.Lock()
	defer hook.l.Unlock()
	hook.message = message
}

func (hook *MyHook) Fire(level LevelType, tags map[string]string, args ...interface{}) {
	hook.Add()
	hook.SetLevel(level)
	hook.SetMessage(fmt.Sprint(args...))
}

func TestHook(t *testing.T) {
	hook := NewMyHook()

	err := NewFileWriter("/tmp", false)
	defer Close()
	if nil != err {
		t.Errorf("initialize file writer faied. err: %s", err.Error())
	}

	blog.SetHook(hook)
	blog.SetHookLevel(INFO)

	blog.Debug("something")
	// wait for hook called
	time.Sleep(1 * time.Millisecond)
	if 0 != hook.Cnt() {
		t.Error("hook called not valid")
	}

	if DEBUG == hook.Level() || "something" == hook.Message() {
		t.Errorf("hook parameters wrong. level: %s, message: %s", hook.level.String(), hook.message)
	}

	blog.Info("yes")
	// wait for hook called
	time.Sleep(1 * time.Millisecond)
	if 1 != hook.Cnt() {
		t.Error("hook not called")
	}

	if INFO != hook.Level() || "yes" != hook.Message() {
		t.Errorf("hook parameters wrong. level: %d, message: %s", hook.level, hook.message)
	}

	// clean logs
	_, err = exec.Command("/bin/sh", "-c", "/bin/rm /tmp/*.log*").Output()
	if nil != err {
		t.Errorf("clean files failed. err: %s", err.Error())
	}
}
