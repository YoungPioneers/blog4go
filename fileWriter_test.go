// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

type T struct {
	A int
	B string
}

func TestFileWriterBasicOperation(t *testing.T) {
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

	blog.Close()
	blog.Debug("Debug", 1)
	blog.Debugf("%s", "Debug")
}

func TestSingleFileWriter(t *testing.T) {
	err := NewFileWriter("/tmp", false)
	defer func() {
		Close()

		// clean logs
		_, err = exec.Command("/bin/sh", "-c", "/bin/rm /tmp/*.log*").Output()
		if nil != err {
			t.Errorf("clean files failed. err: %s", err.Error())
		}
	}()

	// duplicate initialization test
	err = NewFileWriter("/tmp", false)
	if ErrAlreadyInit != err {
		t.Error("duplicate init check fail")
	}

	// should be closed
	Close()
	if nil != blog {
		t.Error("blog should be closed.")
	}
}

func TestFileWriterAsConfigFile(t *testing.T) {
	err := NewWriterFromConfigAsFile("config.example.xml")
	defer func() {
		Close()

		// clean logs
		_, err = exec.Command("/bin/sh", "-c", "/bin/rm /tmp/*.log*").Output()
		if nil != err {
			t.Errorf("clean files failed. err: %s", err.Error())
		}
	}()

	if nil != err {
		t.Error(err.Error())
	}

	// duplicate init
	err = NewWriterFromConfigAsFile("config.example.xml")
	if ErrAlreadyInit != err {
		t.Errorf("Duplicate initialization check failed. err: %s", err.Error())
	}

	blog.Debug("Debug")
	blog.Trace("Trace")
	blog.Info("Info")
	blog.Warn("Warn")
	blog.Error("Error")
	blog.Critical("Critical")

	blog.Close()
	blog.Debug("Debug", 1)
	blog.Debugf("%s", "Debug")
}

// test if log lose in multi goroutine mode
func TestFileWriterMultiGoroutine(t *testing.T) {
	err := NewFileWriter("/tmp", false)
	defer func() {
		Close()

		// clean logs
		_, err = exec.Command("/bin/sh", "-c", "/bin/rm /tmp/*.log*").Output()
		if nil != err {
			t.Errorf("clean files failed. err: %s", err.Error())
		}
	}()
	temp := T{123, "test"}

	if nil != err {
		t.Errorf("initialize file writer faied. err: %s", err.Error())
	}

	var wg sync.WaitGroup
	var beginWg sync.WaitGroup

	f := func() {
		defer wg.Done()
		beginWg.Wait()
		for i := 0; i < 100; i++ {
			blog.Infof("haha %s. en\\en, always %d and %f, %t, %+v", "eddie", 18, 3.1415, true, temp)
			blog.Info("test for not formated")
		}
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		beginWg.Add(1)
	}

	// write 100 * 100 lines to /tmp/info.log
	for i := 0; i < 100; i++ {
		go f()
		beginWg.Done()
	}

	wg.Wait()
	Flush()

	out, err := exec.Command("/bin/sh", "-c", "/usr/bin/wc -l /tmp/info.log").Output()
	if nil != err {
		t.Errorf("count file lines failed. err: %s", err.Error())
	}

	arr := strings.Split(string(out), " ")
	intStr := arr[len(arr)-2]
	lines, err := strconv.Atoi(intStr)
	if nil != err {
		t.Errorf("line string convert to int failed. err: %s", err.Error())
	}

	if 100*100*2 != lines {
		t.Errorf("it loses %d lines.", 100*100-lines)
	}

	// check log message line by line
	file, err := os.Open("/tmp/info.log")
	if err != nil {
		t.Error(err.Error())

	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	line := 0
	for scanner.Scan() {
		line++
		lineStr := scanner.Text()
		arrs := strings.Split(lineStr, "msg=\"")
		if len(arrs) != 2 {
			t.Errorf("line %d detect inconsistent line. not formatted. lineStr: %s", line, lineStr)
		}

		if "haha eddie. en\\en, always 18 and 3.141500, true, {A:123 B:test}\" " != arrs[1] && "test for not formated\" " != arrs[1] {
			t.Errorf("line %d detect inconsistent line. message not correct. lineStr: %s", line, lineStr)
		}
	}

	if err := scanner.Err(); err != nil {
		t.Error(err.Error())
	}
}

func TestFileWriterSizeBaseLogrotate(t *testing.T) {
	err := NewFileWriter("/tmp", false)
	defer func() {
		Close()

		// clean logs
		_, err = exec.Command("/bin/sh", "-c", "/bin/rm /tmp/*.log*").Output()
		if nil != err {
			t.Errorf("clean files failed. err: %s", err.Error())
		}
	}()

	if nil != err {
		t.Errorf("initialize file writer faied. err: %s", err.Error())
	}
	blog.SetRotateSize(60)
	blog.SetTimeRotated(true)

	blog.Info("1")
	Flush()
	time.Sleep(1 * time.Millisecond)

	if _, err = os.Stat("/tmp/info.log.1"); nil == err {
		t.Error("size base logrotate failed, log should not exist.")
	}

	blog.Info("2")
	Flush()
	time.Sleep(1 * time.Millisecond)

	blog.Info("3")
	Flush()
	time.Sleep(1 * time.Millisecond)

	if _, err = os.Stat("/tmp/info.log.1"); os.IsNotExist(err) {
		t.Errorf("size base logrotate failed. err: %s", err.Error())
	}
}

// TODO how to test time base logratate ?
func TestFileWriterTimeBaseLogrotate(t *testing.T) {
	err := NewFileWriter("/tmp", true)
	defer func() {
		Close()

		// clean logs
		_, err = exec.Command("/bin/sh", "-c", "/bin/rm /tmp/*.log*").Output()
		if nil != err {
			t.Errorf("clean files failed. err: %s", err.Error())
		}
	}()

	SetRetentions(1)

	// check if formatted file name exist
	if _, err = os.Stat(fmt.Sprintf("/tmp/trace.log.%s", timeCache.date)); os.IsNotExist(err) {
		t.Error("time base logrotate formatted file name incorrect.")
	}
	if _, err = os.Stat(fmt.Sprintf("/tmp/debug.log.%s", timeCache.date)); os.IsNotExist(err) {
		t.Error("time base logrotate formatted file name incorrect.")
	}
	if _, err = os.Stat(fmt.Sprintf("/tmp/info.log.%s", timeCache.date)); os.IsNotExist(err) {
		t.Error("time base logrotate formatted file name incorrect.")
	}
	if _, err = os.Stat(fmt.Sprintf("/tmp/warn.log.%s", timeCache.date)); os.IsNotExist(err) {
		t.Error("time base logrotate formatted file name incorrect.")
	}
	if _, err = os.Stat(fmt.Sprintf("/tmp/error.log.%s", timeCache.date)); os.IsNotExist(err) {
		t.Error("time base logrotate formatted file name incorrect.")
	}
	if _, err = os.Stat(fmt.Sprintf("/tmp/critical.log.%s", timeCache.date)); os.IsNotExist(err) {
		t.Error("time base logrotate formatted file name incorrect.")
	}
}

func TestFileWriterLinesBaseLogrotate(t *testing.T) {
	err := NewFileWriter("/tmp", false)
	defer func() {
		Close()

		// clean logs
		_, err = exec.Command("/bin/sh", "-c", "/bin/rm /tmp/*.log*").Output()
		if nil != err {
			t.Errorf("clean files failed. err: %s", err.Error())
		}
	}()

	if nil != err {
		t.Errorf("initialize file writer faied. err: %s", err.Error())
	}

	blog.SetRotateLines(2)
	blog.Info("some")
	Flush()
	time.Sleep(1 * time.Millisecond)

	if _, err = os.Stat("/tmp/info.log.1"); nil == err {
		t.Error("line base logrotate failed, log should not exist.")
	}

	blog.Info("some")
	Flush()
	time.Sleep(1 * time.Millisecond)

	blog.Info("some")
	Flush()
	time.Sleep(1 * time.Millisecond)

	if _, err = os.Stat("/tmp/info.log.1"); os.IsNotExist(err) {
		t.Errorf("line base logrotate failed. err: %s", err.Error())
	}
}

func TestFileWriterLogrorateRetentionCount(t *testing.T) {
	err := NewFileWriter("/tmp", false)
	defer func() {
		Close()

		// clean logs
		_, err = exec.Command("/bin/sh", "-c", "/bin/rm /tmp/*.log*").Output()
		if nil != err {
			t.Errorf("clean files failed. err: %s", err.Error())
		}
	}()

	if nil != err {
		t.Errorf("initialize file writer faied. err: %s", err.Error())
	}

	blog.SetRotateLines(2)
	blog.SetRetentions(1)

	blog.Info("1")
	Flush()
	time.Sleep(1 * time.Millisecond)

	if _, err = os.Stat("/tmp/info.log.1"); nil == err {
		t.Error("logrotate retention failed, log should not exist.")
	}

	blog.Info("2")
	Flush()
	time.Sleep(1 * time.Millisecond)

	blog.Info("3")
	Flush()
	time.Sleep(1 * time.Millisecond)

	blog.Info("4")
	Flush()
	time.Sleep(1 * time.Millisecond)

	blog.Info("5")
	Flush()
	time.Sleep(1 * time.Millisecond)
	if _, err = os.Stat("/tmp/info.log.2"); nil == err {
		t.Error("logrotate retention failed.")
	}
}

func TestFileWriterSinglton(t *testing.T) {
	err := NewFileWriter("/tmp", false)
	defer func() {
		Close()

		// clean logs
		_, err = exec.Command("/bin/sh", "-c", "/bin/rm /tmp/*.log*").Output()
		if nil != err {
			t.Errorf("clean files failed. err: %s", err.Error())
		}
	}()

	if nil != err {
		t.Errorf("initialize file writer faied. err: %s", err.Error())
	}

	err = NewFileWriter("/tmp", false)
	if ErrAlreadyInit != err {
		t.Errorf("file writer singlton failed. err: %s", err.Error())
	}
}

func BenchmarkFileWriters(b *testing.B) {
	err := NewFileWriter("/tmp", false)
	defer func() {
		Close()

		// clean logs
		_, err = exec.Command("/bin/sh", "-c", "/bin/rm /tmp/*.log*").Output()
		if nil != err {
			b.Errorf("clean files failed. err: %s", err.Error())
		}
	}()

	if nil != err {
		b.Error(err.Error())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		blog.Debugf("haha %s. en\\en, always %d and %f", "eddie", 18, 3.1415)
	}
	Flush()
}
