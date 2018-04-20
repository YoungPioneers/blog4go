// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"bytes"
	"fmt"
	"net"
	"sync"
)

// SocketWriter is a socket logger
type SocketWriter struct {
	level LevelType

	closed bool

	// log hook
	hook      Hook
	hookLevel LevelType
	hookAsync bool

	// socket
	writer net.Conn

	lock *sync.RWMutex

	// tags
	tags   map[string]string
	tagStr string
}

// NewSocketWriter creates a socket writer, singlton
func NewSocketWriter(network string, address string) (err error) {
	singltonLock.Lock()
	defer singltonLock.Unlock()
	if nil != blog {
		return ErrAlreadyInit
	}

	socketWriter, err := newSocketWriter(network, address)
	if nil != err {
		return err
	}

	blog = socketWriter
	return nil
}

// newSocketWriter creates a socket writer, not singlton
func newSocketWriter(network string, address string) (socketWriter *SocketWriter, err error) {
	socketWriter = new(SocketWriter)
	socketWriter.level = DEBUG
	socketWriter.closed = false
	socketWriter.lock = new(sync.RWMutex)

	// log hook
	socketWriter.hook = nil
	socketWriter.hookLevel = DEBUG

	conn, err := net.Dial(network, address)
	if nil != err {
		return nil, err
	}
	socketWriter.writer = conn

	blog = socketWriter
	return socketWriter, nil
}

func (writer *SocketWriter) write(level LevelType, args ...interface{}) {
	writer.lock.RLock()
	defer writer.lock.RUnlock()

	if writer.closed {
		return
	}

	buffer := bytes.NewBuffer(timeCache.Format())
	buffer.WriteString(level.prefix())
	buffer.WriteString(writer.tagStr)
	buffer.WriteString(fmt.Sprintf("msg=\"%s\" ", fmt.Sprint(args...)))
	buffer.WriteByte(EOL)
	writer.writer.Write(buffer.Bytes())

	// call log hook
	if nil != writer.hook && !(level < writer.hookLevel) {
		if writer.hookAsync {
			go writer.hook.Fire(level, writer.Tags(), args...)
		} else {
			writer.hook.Fire(level, writer.Tags(), args...)
		}
	}
}

func (writer *SocketWriter) writef(level LevelType, format string, args ...interface{}) {
	writer.lock.RLock()
	defer writer.lock.RUnlock()

	if writer.closed {
		return
	}

	buffer := bytes.NewBuffer(timeCache.Format())
	buffer.WriteString(level.prefix())
	buffer.WriteString(writer.tagStr)
	buffer.WriteString(fmt.Sprintf("msg=\"%s\" ", fmt.Sprintf(format, args...)))
	buffer.WriteByte(EOL)
	writer.writer.Write(buffer.Bytes())

	// call log hook
	if nil != writer.hook && !(level < writer.hookLevel) {
		if writer.hookAsync {
			go writer.hook.Fire(level, writer.Tags(), fmt.Sprintf(format, args...))
		} else {
			writer.hook.Fire(level, writer.Tags(), fmt.Sprintf(format, args...))
		}
	}
}

// Level get level
func (writer *SocketWriter) Level() LevelType {
	return writer.level
}

// SetLevel set logger level
func (writer *SocketWriter) SetLevel(level LevelType) {
	writer.level = level
}

// Tags return logging tags
func (writer *SocketWriter) Tags() map[string]string {
	writer.lock.RLock()
	defer writer.lock.RUnlock()
	return writer.tags
}

// SetTags set logging tags
func (writer *SocketWriter) SetTags(tags map[string]string) {
	writer.lock.Lock()
	defer writer.lock.Unlock()
	writer.tags = tags

	var tagStr string
	for tagName, tagValue := range writer.tags {
		tagStr = fmt.Sprintf("%s%s=\"%s\" ", tagStr, tagName, tagValue)
	}

	writer.tagStr = tagStr
}

// SetHook set hook for logging action
func (writer *SocketWriter) SetHook(hook Hook) {
	writer.hook = hook
}

// SetHookAsync set hook async for base file writer
func (writer *SocketWriter) SetHookAsync(async bool) {
	writer.hookAsync = async
}

// SetHookLevel set when hook will be called
func (writer *SocketWriter) SetHookLevel(level LevelType) {
	writer.hookLevel = level
}

// TimeRotated do nothing
func (writer *SocketWriter) TimeRotated() bool {
	return false
}

// SetTimeRotated do nothing
func (writer *SocketWriter) SetTimeRotated(timeRotated bool) {
	return
}

// Retentions do nothing
func (writer *SocketWriter) Retentions() int64 {
	return 0
}

// SetRetentions do nothing
func (writer *SocketWriter) SetRetentions(retentions int64) {
	return
}

// RotateSize do nothing
func (writer *SocketWriter) RotateSize() int64 {
	return 0
}

// SetRotateSize do nothing
func (writer *SocketWriter) SetRotateSize(rotateSize int64) {
	return
}

// RotateLines do nothing
func (writer *SocketWriter) RotateLines() int {
	return 0
}

// SetRotateLines do nothing
func (writer *SocketWriter) SetRotateLines(rotateLines int) {
	return
}

// Colored do nothing
func (writer *SocketWriter) Colored() bool {
	return false
}

// SetColored do nothing
func (writer *SocketWriter) SetColored(colored bool) {
	return
}

// Close will close the writer
func (writer *SocketWriter) Close() {
	writer.lock.Lock()
	defer writer.lock.Unlock()
	if writer.closed {
		return
	}

	writer.writer.Close()
	writer.writer = nil
	writer.closed = true
}

// flush do nothing
func (writer *SocketWriter) flush() {
	return
}

// Trace trace
func (writer *SocketWriter) Trace(args ...interface{}) {
	if nil == writer.writer || TRACE < writer.level {
		return
	}

	writer.write(TRACE, args...)
}

// Tracef tracef
func (writer *SocketWriter) Tracef(format string, args ...interface{}) {
	if nil == writer.writer || TRACE < writer.level {
		return
	}

	writer.writef(TRACE, format, args...)
}

// Debug debug
func (writer *SocketWriter) Debug(args ...interface{}) {
	if nil == writer.writer || DEBUG < writer.level {
		return
	}

	writer.write(DEBUG, args...)
}

// Debugf debugf
func (writer *SocketWriter) Debugf(format string, args ...interface{}) {
	if nil == writer.writer || DEBUG < writer.level {
		return
	}

	writer.writef(DEBUG, format, args...)
}

// Info info
func (writer *SocketWriter) Info(args ...interface{}) {
	if nil == writer.writer || INFO < writer.level {
		return
	}

	writer.write(INFO, args...)
}

// Infof infof
func (writer *SocketWriter) Infof(format string, args ...interface{}) {
	if nil == writer.writer || INFO < writer.level {
		return
	}

	writer.writef(INFO, format, args...)
}

// Warn warn
func (writer *SocketWriter) Warn(args ...interface{}) {
	if nil == writer.writer || WARNING < writer.level {
		return
	}

	writer.write(WARNING, args...)
}

// Warnf warnf
func (writer *SocketWriter) Warnf(format string, args ...interface{}) {
	if nil == writer.writer || WARNING < writer.level {
		return
	}

	writer.writef(WARNING, format, args...)
}

// Error error
func (writer *SocketWriter) Error(args ...interface{}) {
	if nil == writer.writer || ERROR < writer.level {
		return
	}

	writer.write(ERROR, args...)
}

// Errorf error
func (writer *SocketWriter) Errorf(format string, args ...interface{}) {
	if nil == writer.writer || ERROR < writer.level {
		return
	}

	writer.writef(ERROR, format, args...)
}

// Critical critical
func (writer *SocketWriter) Critical(args ...interface{}) {
	if nil == writer.writer || CRITICAL < writer.level {
		return
	}

	writer.write(CRITICAL, args...)
}

// Criticalf criticalf
func (writer *SocketWriter) Criticalf(format string, args ...interface{}) {
	if nil == writer.writer || CRITICAL < writer.level {
		return
	}

	writer.writef(CRITICAL, format, args...)
}
