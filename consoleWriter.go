// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"fmt"
	"os"
	"time"
)

// ConsoleWriter is a console logger
type ConsoleWriter struct {
	blog *BLog

	closed bool

	colored bool

	// log hook
	hook      Hook
	hookLevel Level
}

// NewConsoleWriter initialize a console writer, singlton
func NewConsoleWriter() (err error) {
	singltonLock.Lock()
	defer singltonLock.Unlock()
	if nil != blog {
		return ErrAlreadyInit
	}

	consoleWriter, err := newConsoleWriter()
	if nil != err {
		return err
	}

	blog = consoleWriter
	go consoleWriter.daemon()
	return nil
}

// newConsoleWriter initialize a console writer, not singlton
func newConsoleWriter() (consoleWriter *ConsoleWriter, err error) {
	consoleWriter = new(ConsoleWriter)
	consoleWriter.blog = NewBLog(os.Stdout)

	consoleWriter.closed = false

	consoleWriter.colored = false

	// log hook
	consoleWriter.hook = nil
	consoleWriter.hookLevel = DEBUG

	go consoleWriter.daemon()

	blog = consoleWriter
	return consoleWriter, nil
}

func (writer *ConsoleWriter) daemon() {
	f := time.Tick(10 * time.Second)

DaemonLoop:
	for {
		select {
		case <-f:
			if writer.closed {
				break DaemonLoop
			}

			writer.flush()
		}
	}
}

func (writer *ConsoleWriter) write(level Level, args ...interface{}) {
	if writer.closed {
		return
	}

	defer func() {
		if nil != writer.hook && !(level < writer.hookLevel) {
			go func(level Level, args ...interface{}) {
				writer.hook.Fire(level, args...)
			}(level, args...)
		}
	}()

	writer.blog.write(level, args...)
}

func (writer *ConsoleWriter) writef(level Level, format string, args ...interface{}) {
	if writer.closed {
		return
	}

	defer func() {

		if nil != writer.hook && !(level < writer.hookLevel) {
			go func(level Level, format string, args ...interface{}) {
				writer.hook.Fire(level, fmt.Sprintf(format, args...))
			}(level, format, args...)
		}
	}()

	writer.blog.writef(level, format, args...)
}

// SetLevel set logger level
func (writer *ConsoleWriter) SetLevel(level Level) {
	writer.blog.SetLevel(level)
}

// SetColored set logging color
func (writer *ConsoleWriter) SetColored(colored bool) {
	if colored == writer.colored {
		return
	}

	writer.colored = colored

	initPrefix(colored)
}

// SetHook set hook for logging action
func (writer *ConsoleWriter) SetHook(hook Hook) {
	writer.hook = hook
}

// SetHookLevel set when hook will be called
func (writer *ConsoleWriter) SetHookLevel(level Level) {
	writer.hookLevel = level
}

// Close close console writer
func (writer *ConsoleWriter) Close() {
	if writer.closed {
		return
	}

	writer.blog.flush()
	writer.blog = nil
	writer.closed = true
}

// SetTimeRotated do nothing
func (writer *ConsoleWriter) SetTimeRotated(timeRotated bool) {
	return
}

// SetRetentions do nothing
func (writer *ConsoleWriter) SetRetentions(retentions int64) {
	return
}

// SetRotateSize do nothing
func (writer *ConsoleWriter) SetRotateSize(rotateSize int64) {
	return
}

// SetRotateLines do nothing
func (writer *ConsoleWriter) SetRotateLines(rotateLines int) {
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
