// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"errors"
	"fmt"
	"sync"
)

var (
	// ErrFilePathNotFound file path not found
	ErrFilePathNotFound = errors.New("File Path must be defined")
	// ErrInvalidLevel invalid level string
	ErrInvalidLevel = errors.New("Invalid level string")
	// ErrInvalidRotateType invalid logrotate type
	ErrInvalidRotateType = errors.New("Invalid log rotate type")
)

// MultiWriter struct defines an instance for multi writers with different message level
type MultiWriter struct {
	level LevelType

	// file writers
	writers map[LevelType]Writer

	colored bool

	closed bool

	// configuration about user defined logging hook
	// actual hook instance
	hook Hook
	// hook is called when message level exceed level of logging action
	hookLevel LevelType
	// it determines whether hook is called async, default true
	hookAsync bool

	// logrotate
	timeRotated bool
	retentions  int64
	rotateSize  int64
	rotateLines int

	// tags
	tags map[string]string

	// lock
	lock *sync.RWMutex
}

// TimeRotated get timeRotated
func (writer *MultiWriter) TimeRotated() bool {
	return writer.timeRotated
}

// SetTimeRotated toggle time base logrotate
func (writer *MultiWriter) SetTimeRotated(timeRotated bool) {
	writer.timeRotated = timeRotated
	for _, fileWriter := range writer.writers {
		fileWriter.SetTimeRotated(timeRotated)
	}
}

// Retentions get retentions
func (writer *MultiWriter) Retentions() int64 {
	return writer.retentions
}

// SetRetentions set how many logs will keep after logrotate
func (writer *MultiWriter) SetRetentions(retentions int64) {
	if retentions < 1 {
		return
	}

	writer.retentions = retentions
	for _, fileWriter := range writer.writers {
		fileWriter.SetRetentions(retentions)
	}
}

// RotateSize get rotateSize
func (writer *MultiWriter) RotateSize() int64 {
	return writer.rotateSize
}

// SetRotateSize set size when logroatate
func (writer *MultiWriter) SetRotateSize(rotateSize int64) {
	writer.rotateSize = rotateSize
	for _, fileWriter := range writer.writers {
		fileWriter.SetRotateSize(rotateSize)
	}
}

// RotateLines get rotateLines
func (writer *MultiWriter) RotateLines() int {
	return writer.rotateLines
}

// SetRotateLines set line number when logrotate
func (writer *MultiWriter) SetRotateLines(rotateLines int) {
	writer.rotateLines = rotateLines
	for _, fileWriter := range writer.writers {
		fileWriter.SetRotateLines(rotateLines)
	}
}

// Colored get colored
func (writer *MultiWriter) Colored() bool {
	return writer.colored
}

// SetColored set logging color
func (writer *MultiWriter) SetColored(colored bool) {
	writer.colored = colored
	for _, fileWriter := range writer.writers {
		fileWriter.SetColored(colored)
	}
}

// SetHook set hook for every logging actions
func (writer *MultiWriter) SetHook(hook Hook) {
	writer.hook = hook
}

// SetHookAsync set hook async for base file writer
func (writer *MultiWriter) SetHookAsync(async bool) {
	writer.hookAsync = async
}

// SetHookLevel set when hook will be called
func (writer *MultiWriter) SetHookLevel(level LevelType) {
	writer.hookLevel = level
}

// SetLevel set logging level threshold
func (writer *MultiWriter) SetLevel(level LevelType) {
	writer.level = level
	for _, fileWriter := range writer.writers {
		fileWriter.SetLevel(level)
	}
}

// Tags return logging tags
func (writer *MultiWriter) Tags() map[string]string {
	writer.lock.RLock()
	defer writer.lock.RUnlock()
	return writer.tags
}

// SetTags set logging tags
func (writer *MultiWriter) SetTags(tags map[string]string) {
	writer.lock.Lock()
	defer writer.lock.Unlock()
	writer.tags = tags

	for _, singleWriter := range writer.writers {
		singleWriter.SetTags(tags)
	}
}

// Level return logging level threshold
func (writer *MultiWriter) Level() LevelType {
	return writer.level
}

// Close close file writer
func (writer *MultiWriter) Close() {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	for _, fileWriter := range writer.writers {
		fileWriter.Close()
	}
	writer.closed = true
}

func (writer *MultiWriter) write(level LevelType, args ...interface{}) {
	writer.writers[level].write(level, args...)

	if nil != writer.hook && !(level < writer.hookLevel) {
		if writer.hookAsync {
			// 异步调用log hook
			go writer.hook.Fire(level, writer.Tags(), args...)
		} else {
			writer.hook.Fire(level, writer.Tags(), args...)
		}
	}
}

func (writer *MultiWriter) writef(level LevelType, format string, args ...interface{}) {
	writer.writers[level].writef(level, format, args...)

	if nil != writer.hook && !(level < writer.hookLevel) {
		if writer.hookAsync {
			// 异步调用log hook
			go writer.hook.Fire(level, writer.Tags(), fmt.Sprintf(format, args...))
		} else {
			writer.hook.Fire(level, writer.Tags(), fmt.Sprintf(format, args...))

		}
	}
}

// flush flush logs to disk
func (writer *MultiWriter) flush() {
	for _, writer := range writer.writers {
		writer.flush()
	}
}

// Trace trace
func (writer *MultiWriter) Trace(args ...interface{}) {
	_, ok := writer.writers[TRACE]
	if !ok || TRACE < writer.level {
		return
	}

	writer.write(TRACE, args...)
}

// Tracef tracef
func (writer *MultiWriter) Tracef(format string, args ...interface{}) {
	_, ok := writer.writers[TRACE]
	if !ok || TRACE < writer.level {
		return
	}

	writer.writef(TRACE, format, args...)
}

// Debug debug
func (writer *MultiWriter) Debug(args ...interface{}) {
	_, ok := writer.writers[DEBUG]
	if !ok || DEBUG < writer.level {
		return
	}

	writer.write(DEBUG, args...)
}

// Debugf debugf
func (writer *MultiWriter) Debugf(format string, args ...interface{}) {
	_, ok := writer.writers[DEBUG]
	if !ok || DEBUG < writer.level {
		return
	}

	writer.writef(DEBUG, format, args...)
}

// Info info
func (writer *MultiWriter) Info(args ...interface{}) {
	_, ok := writer.writers[INFO]
	if !ok || INFO < writer.level {
		return
	}

	writer.write(INFO, args...)
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
func (writer *MultiWriter) Warn(args ...interface{}) {
	_, ok := writer.writers[WARNING]
	if !ok || WARNING < writer.level {
		return
	}

	writer.write(WARNING, args...)
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
func (writer *MultiWriter) Error(args ...interface{}) {
	_, ok := writer.writers[ERROR]
	if !ok || ERROR < writer.level {
		return
	}

	writer.write(ERROR, args...)
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
func (writer *MultiWriter) Critical(args ...interface{}) {
	_, ok := writer.writers[CRITICAL]
	if !ok || CRITICAL < writer.level {
		return
	}

	writer.write(CRITICAL, args...)
}

// Criticalf criticalf
func (writer *MultiWriter) Criticalf(format string, args ...interface{}) {
	_, ok := writer.writers[CRITICAL]
	if !ok || CRITICAL < writer.level {
		return
	}

	writer.writef(CRITICAL, format, args...)
}
