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
	hookLevel LevelType
	hookAsync bool
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
	consoleWriter.hookAsync = true

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

func (writer *ConsoleWriter) write(level LevelType, args ...interface{}) {
	if writer.closed {
		return
	}

	defer func() {
		if nil != writer.hook && !(level < writer.hookLevel) {
			if writer.hookAsync {
				go func(level LevelType, args ...interface{}) {
					writer.hook.Fire(level, args...)
				}(level, args...)

			} else {
				writer.hook.Fire(level, args...)
			}
		}
	}()

	writer.blog.write(level, args...)
}

func (writer *ConsoleWriter) writef(level LevelType, format string, args ...interface{}) {
	if writer.closed {
		return
	}

	defer func() {

		if nil != writer.hook && !(level < writer.hookLevel) {
			if writer.hookAsync {
				go func(level LevelType, format string, args ...interface{}) {
					writer.hook.Fire(level, fmt.Sprintf(format, args...))
				}(level, format, args...)

			} else {
				writer.hook.Fire(level, fmt.Sprintf(format, args...))

			}
		}
	}()

	writer.blog.writef(level, format, args...)
}

// Level get level
func (writer *ConsoleWriter) Level() LevelType {
	return writer.blog.Level()
}

// SetLevel set logger level
func (writer *ConsoleWriter) SetLevel(level LevelType) {
	writer.blog.SetLevel(level)
}

// Colored get Colored
func (writer *ConsoleWriter) Colored() bool {
	return writer.colored
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

// SetHookAsync set hook async for base file writer
func (writer *ConsoleWriter) SetHookAsync(async bool) {
	writer.hookAsync = async
}

// SetHookLevel set when hook will be called
func (writer *ConsoleWriter) SetHookLevel(level LevelType) {
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

// TimeRotated do nothing
func (writer *ConsoleWriter) TimeRotated() bool {
	return false
}

// SetTimeRotated do nothing
func (writer *ConsoleWriter) SetTimeRotated(timeRotated bool) {
	return
}

// Retentions do nothing
func (writer *ConsoleWriter) Retentions() int64 {
	return 0
}

// SetRetentions do nothing
func (writer *ConsoleWriter) SetRetentions(retentions int64) {
	return
}

// RotateSize do nothing
func (writer *ConsoleWriter) RotateSize() int64 {
	return 0
}

// SetRotateSize do nothing
func (writer *ConsoleWriter) SetRotateSize(rotateSize int64) {
	return
}

// RotateLines do nothing
func (writer *ConsoleWriter) RotateLines() int {
	return 0
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
