// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"bufio"
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

	// test file writer hook
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
	err := NewWriterFromConfigAsFile("examples/writer_from_configfile/config.example.xml")
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

	blog.Debug("Debug")
	blog.Trace("Trace")
	blog.Info("Info")
	blog.Warn("Warn")
	blog.Error("Error")
	blog.Critical("Critical")
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

	if 100*100 != lines {
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
		arrs := strings.Split(lineStr, "[INFO] ")
		if len(arrs) != 2 {
			t.Errorf("line %d detect inconsistent line. not formatted. lineStr: %s", line, lineStr)
		}

		if "haha eddie. en\\en, always 18 and 3.141500, true, {A:123 B:test}" != arrs[1] {
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
	blog.SetRotateSize(2)

	blog.Info("1")
	Flush()
	time.Sleep(1 * time.Millisecond)

	if _, err = os.Stat("/tmp/info.log.1"); os.IsExist(err) {
		t.Errorf("size base logrotate failed, log should not exist. err: %s", err.Error())
	}

	blog.Info("2")
	Flush()
	time.Sleep(1 * time.Millisecond)

	blog.Info("3")
	Flush()
	time.Sleep(1 * time.Millisecond)

	if _, err = os.Stat("/tmp/info.log.1"); os.IsNotExist(err) {
		t.Errorf("size base logrotate failed., err: %s", err.Error())
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

	if _, err = os.Stat("/tmp/info.log.1"); os.IsExist(err) {
		t.Errorf("line base logrotate failed, log should not exist. err: %s", err.Error())
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

	if _, err = os.Stat("/tmp/info.log.1"); os.IsExist(err) {
		t.Errorf("logrotate retention failed, log should not exist. err: %s", err.Error())
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
	if _, err = os.Stat("/tmp/info.log.2"); os.IsExist(err) {
		t.Errorf("logrotate retention failed. err: %s", err.Error())
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
