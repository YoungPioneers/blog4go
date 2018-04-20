// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// ConsoleWriter is a console logger
type ConsoleWriter struct {
	blog *BLog
	// for stderr
	errblog *BLog

	redirected bool

	closed bool

	colored bool

	// log hook
	hook      Hook
	hookLevel LevelType
	hookAsync bool

	lock *sync.RWMutex
}

// NewConsoleWriter initialize a console writer, singlton
func NewConsoleWriter(redirected bool) (err error) {
	singltonLock.Lock()
	defer singltonLock.Unlock()
	if nil != blog {
		return ErrAlreadyInit
	}

	consoleWriter, err := newConsoleWriter(redirected)
	if nil != err {
		return err
	}

	blog = consoleWriter
	go consoleWriter.daemon()
	return nil
}

// newConsoleWriter initialize a console writer, not singlton
// if redirected, stderr will be redirected to stdout
func newConsoleWriter(redirected bool) (consoleWriter *ConsoleWriter, err error) {
	consoleWriter = new(ConsoleWriter)
	consoleWriter.blog = NewBLog(os.Stdout)
	consoleWriter.redirected = redirected
	if !redirected {
		consoleWriter.errblog = NewBLog(os.Stderr)
	}

	consoleWriter.closed = false
	consoleWriter.colored = false

	// log hook
	consoleWriter.hook = nil
	consoleWriter.hookLevel = DEBUG
	consoleWriter.hookAsync = true

	consoleWriter.lock = new(sync.RWMutex)

	go consoleWriter.daemon()

	blog = consoleWriter
	return consoleWriter, nil
}

func (writer *ConsoleWriter) daemon() {
	f := time.Tick(1 * time.Second)

DaemonLoop:
	for {
		select {
		case <-f:
			if writer.Closed() {
				break DaemonLoop
			}

			writer.flush()
		}
	}
}

func (writer *ConsoleWriter) write(level LevelType, args ...interface{}) {
	writer.lock.RLock()
	defer writer.lock.RUnlock()

	if writer.closed {
		return
	}

	if !writer.redirected && level >= WARNING {
		writer.errblog.write(level, args...)
	} else {
		writer.blog.write(level, args...)
	}

	if nil != writer.hook && !(level < writer.hookLevel) && !writer.closed {
		if writer.hookAsync {
			go writer.hook.Fire(level, writer.blog.Tags(), args...)

		} else {
			writer.hook.Fire(level, writer.blog.Tags(), args...)
		}
	}
}

func (writer *ConsoleWriter) writef(level LevelType, format string, args ...interface{}) {
	writer.lock.RLock()
	defer writer.lock.RUnlock()

	if writer.closed {
		return
	}

	if !writer.redirected && level >= WARNING {
		writer.errblog.writef(level, format, args...)
	} else {
		writer.blog.writef(level, format, args...)
	}

	if nil != writer.hook && !(level < writer.hookLevel) && !writer.closed {
		if writer.hookAsync {
			go writer.hook.Fire(level, writer.blog.Tags(), fmt.Sprintf(format, args...))
		} else {
			writer.hook.Fire(level, writer.blog.Tags(), fmt.Sprintf(format, args...))
		}
	}
}

// Closed get writer status
func (writer *ConsoleWriter) Closed() bool {
	writer.lock.RLock()
	writer.lock.RUnlock()

	return writer.closed
}

// Level get level
func (writer *ConsoleWriter) Level() LevelType {
	writer.lock.RLock()
	defer writer.lock.RUnlock()

	return writer.blog.Level()
}

// SetLevel set logger level
func (writer *ConsoleWriter) SetLevel(level LevelType) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	writer.blog.SetLevel(level)
}

// Tags return logging tags
func (writer *ConsoleWriter) Tags() map[string]string {
	writer.lock.RLock()
	defer writer.lock.RUnlock()

	return writer.blog.Tags()
}

// SetTags set logging tags
func (writer *ConsoleWriter) SetTags(tags map[string]string) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	writer.blog.SetTags(tags)
}

// Colored get Colored
func (writer *ConsoleWriter) Colored() bool {
	writer.lock.RLock()
	defer writer.lock.RUnlock()

	return writer.colored
}

// SetColored set logging color
func (writer *ConsoleWriter) SetColored(colored bool) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	if colored == writer.colored {
		return
	}

	writer.colored = colored

	initPrefix(colored)
}

// SetHook set hook for logging action
func (writer *ConsoleWriter) SetHook(hook Hook) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	writer.hook = hook
}

// SetHookAsync set hook async for base file writer
func (writer *ConsoleWriter) SetHookAsync(async bool) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	writer.hookAsync = async
}

// SetHookLevel set when hook will be called
func (writer *ConsoleWriter) SetHookLevel(level LevelType) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	writer.hookLevel = level
}

// Close close console writer
func (writer *ConsoleWriter) Close() {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	if writer.closed {
		return
	}

	writer.blog.flush()
	writer.blog = nil
	writer.closed = true
}

// TimeRotated do nothing
func (writer *ConsoleWriter) TimeRotated() bool {
	writer.lock.RLock()
	defer writer.lock.RUnlock()

	return false
}

// SetTimeRotated do nothing
func (writer *ConsoleWriter) SetTimeRotated(timeRotated bool) {
	return
}

// Retentions do nothing
func (writer *ConsoleWriter) Retentions() int64 {
	writer.lock.RLock()
	defer writer.lock.RUnlock()

	return 0
}

// SetRetentions do nothing
func (writer *ConsoleWriter) SetRetentions(retentions int64) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	return
}

// RotateSize do nothing
func (writer *ConsoleWriter) RotateSize() int64 {
	writer.lock.RLock()
	defer writer.lock.RUnlock()

	return 0
}

// SetRotateSize do nothing
func (writer *ConsoleWriter) SetRotateSize(rotateSize int64) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	return
}

// RotateLines do nothing
func (writer *ConsoleWriter) RotateLines() int {
	writer.lock.RLock()
	defer writer.lock.RUnlock()

	return 0
}

// SetRotateLines do nothing
func (writer *ConsoleWriter) SetRotateLines(rotateLines int) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	return
}

// flush buffer to disk
func (writer *ConsoleWriter) flush() {
	writer.blog.flush()
}

// Trace trace
func (writer *ConsoleWriter) Trace(args ...interface{}) {
	if nil == writer.blog || TRACE < writer.blog.Level() {
		return
	}

	writer.write(TRACE, args...)
}

// Tracef tracef
func (writer *ConsoleWriter) Tracef(format string, args ...interface{}) {
	if nil == writer.blog || TRACE < writer.blog.Level() {
		return
	}

	writer.writef(TRACE, format, args...)
}

// Debug debug
func (writer *ConsoleWriter) Debug(args ...interface{}) {
	if nil == writer.blog || DEBUG < writer.blog.Level() {
		return
	}

	writer.write(DEBUG, args...)
}

// Debugf debugf
func (writer *ConsoleWriter) Debugf(format string, args ...interface{}) {
	if nil == writer.blog || DEBUG < writer.blog.Level() {
		return
	}

	writer.writef(DEBUG, format, args...)
}

// Info info
func (writer *ConsoleWriter) Info(args ...interface{}) {
	if nil == writer.blog || INFO < writer.blog.Level() {
		return
	}

	writer.write(INFO, args...)
}

// Infof infof
func (writer *ConsoleWriter) Infof(format string, args ...interface{}) {
	if nil == writer.blog || INFO < writer.blog.Level() {
		return
	}

	writer.writef(INFO, format, args...)
}

// Warn warn
func (writer *ConsoleWriter) Warn(args ...interface{}) {
	if nil == writer.blog || WARNING < writer.blog.Level() {
		return
	}

	writer.write(WARNING, args...)
}

// Warnf warnf
func (writer *ConsoleWriter) Warnf(format string, args ...interface{}) {
	if nil == writer.blog || WARNING < writer.blog.Level() {
		return
	}

	writer.writef(WARNING, format, args...)
}

// Error error
func (writer *ConsoleWriter) Error(args ...interface{}) {
	if nil == writer.blog || ERROR < writer.blog.Level() {
		return
	}

	writer.write(ERROR, args...)
}

// Errorf errorf
func (writer *ConsoleWriter) Errorf(format string, args ...interface{}) {
	if nil == writer.blog || ERROR < writer.blog.Level() {
		return
	}

	writer.writef(ERROR, format, args...)
}

// Critical critical
func (writer *ConsoleWriter) Critical(args ...interface{}) {
	if nil == writer.blog || CRITICAL < writer.blog.Level() {
		return
	}

	writer.write(CRITICAL, args...)
}

// Criticalf criticalf
func (writer *ConsoleWriter) Criticalf(format string, args ...interface{}) {
	if nil == writer.blog || CRITICAL < writer.blog.Level() {
		return
	}

	writer.writef(CRITICAL, format, args...)
}
