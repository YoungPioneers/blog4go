// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"errors"
	"fmt"
)

var (
	// ErrFilePathNotFound file path not found
	ErrFilePathNotFound = errors.New("File Path must be defined.")
	// ErrInvalidLevel invalid level string
	ErrInvalidLevel = errors.New("Invalid level string.")
	// ErrInvalidRotateType invalid logrotate type
	ErrInvalidRotateType = errors.New("Invalid log rotate type.")
)

// MultiWriter struct defines an instance for multi writers with different message level
type MultiWriter struct {
	level Level

	// file writers
	writers map[Level]Writer

	closed bool

	// configuration about user defined logging hook
	// actual hook instance
	hook Hook
	// hook is called when message level exceed level of logging action
	hookLevel Level
}

// SetTimeRotated toggle time base logrotate
func (writer *MultiWriter) SetTimeRotated(timeRotated bool) {
	for _, fileWriter := range writer.writers {
		fileWriter.SetTimeRotated(timeRotated)
	}
}

// SetRotateSize set size when logroatate
func (writer *MultiWriter) SetRotateSize(rotateSize ByteSize) {
	for _, fileWriter := range writer.writers {
		fileWriter.SetRotateSize(rotateSize)
	}
}

// SetRotateLines set line number when logrotate
func (writer *MultiWriter) SetRotateLines(rotateLines int) {
	for _, fileWriter := range writer.writers {
		fileWriter.SetRotateLines(rotateLines)
	}
}

// SetColored set logging color
func (writer *MultiWriter) SetColored(colored bool) {
	for _, fileWriter := range writer.writers {
		fileWriter.SetColored(colored)
	}
}

// SetHook set hook for every logging actions
func (writer *MultiWriter) SetHook(hook Hook) {
	writer.hook = hook
}

// SetHookLevel set when hook will be called
func (writer *MultiWriter) SetHookLevel(level Level) {
	writer.hookLevel = level
}

// SetLevel set logging level threshold
func (writer *MultiWriter) SetLevel(level Level) {
	writer.level = level
	for _, fileWriter := range writer.writers {
		fileWriter.SetLevel(level)
	}
}

// Level return logging level threshold
func (writer *MultiWriter) Level() Level {
	return writer.level
}

// Close close file writer
func (writer *MultiWriter) Close() {
	for _, fileWriter := range writer.writers {
		fileWriter.Close()
	}
	writer.closed = true
}

func (writer *MultiWriter) write(level Level, format string) {
	defer func() {
		// 异步调用log hook
		if nil != writer.hook && !(level < writer.hookLevel) {
			go func(level Level, format string) {
				writer.hook.Fire(level, format)
			}(level, format)
		}
	}()

	writer.writers[level].write(level, format)
}

func (writer *MultiWriter) writef(level Level, format string, args ...interface{}) {
	defer func() {
		// 异步调用log hook
		if nil != writer.hook && !(level < writer.hookLevel) {
			go func(level Level, format string, args ...interface{}) {
				writer.hook.Fire(level, fmt.Sprintf(format, args...))
			}(level, format, args...)
		}
	}()

	writer.writers[level].writef(level, format, args...)
}

// Debug debug
func (writer *MultiWriter) Debug(format string) {
	_, ok := writer.writers[DEBUG]
	if !ok || DEBUG < writer.level {
		return
	}

	writer.write(DEBUG, format)
}

// Debugf debugf
func (writer *MultiWriter) Debugf(format string, args ...interface{}) {
	_, ok := writer.writers[DEBUG]
	if !ok || DEBUG < writer.level {
		return
	}

	writer.writef(DEBUG, format, args...)
}

// Trace trace
func (writer *MultiWriter) Trace(format string) {
	_, ok := writer.writers[TRACE]
	if !ok || TRACE < writer.level {
		return
	}

	writer.write(TRACE, format)
}

// Tracef tracef
func (writer *MultiWriter) Tracef(format string, args ...interface{}) {
	_, ok := writer.writers[TRACE]
	if !ok || TRACE < writer.level {
		return
	}

	writer.writef(TRACE, format, args...)
}

// Info info
func (writer *MultiWriter) Info(format string) {
	_, ok := writer.writers[INFO]
	if !ok || INFO < writer.level {
		return
	}

	writer.write(INFO, format)
}

// Infof infof
func (writer *MultiWriter) Infof(format string, args ...interface{}) {
	_, ok := writer.writers[INFO]
	if !ok || INFO < writer.level {
		return
	}

	writer.writef(INFO, format, args...)
}

// Warn warn
func (writer *MultiWriter) Warn(format string) {
	_, ok := writer.writers[WARNING]
	if !ok || WARNING < writer.level {
		return
	}

	writer.write(WARNING, format)
}

// Warnf warnf
func (writer *MultiWriter) Warnf(format string, args ...interface{}) {
	_, ok := writer.writers[WARNING]
	if !ok || WARNING < writer.level {
		return
	}

	writer.writef(WARNING, format, args...)
}

// Error error
func (writer *MultiWriter) Error(format string) {
	_, ok := writer.writers[ERROR]
	if !ok || ERROR < writer.level {
		return
	}

	writer.write(ERROR, format)
}

// Errorf error
func (writer *MultiWriter) Errorf(format string, args ...interface{}) {
	_, ok := writer.writers[ERROR]
	if !ok || ERROR < writer.level {
		return
	}

	writer.writef(ERROR, format, args...)
}

// Critical critical
func (writer *MultiWriter) Critical(format string) {
	_, ok := writer.writers[CRITICAL]
	if !ok || CRITICAL < writer.level {
		return
	}

	writer.write(CRITICAL, format)
}

// Criticalf criticalf
func (writer *MultiWriter) Criticalf(format string, args ...interface{}) {
	_, ok := writer.writers[CRITICAL]
	if !ok || CRITICAL < writer.level {
		return
	}

	writer.writef(CRITICAL, format, args...)
}
