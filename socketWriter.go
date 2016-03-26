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
	level Level

	closed bool

	// log hook
	hook      Hook
	hookLevel Level

	// socket
	writer net.Conn

	lock *sync.Mutex
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
	socketWriter.lock = new(sync.Mutex)

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

func (writer *SocketWriter) write(level Level, args ...interface{}) {
	writer.lock.Lock()

	defer func() {
		writer.lock.Unlock()
		// call log hook
		if nil != writer.hook && !(level < writer.hookLevel) {
			go func(level Level, args ...interface{}) {
				writer.hook.Fire(level, args...)
			}(level, args...)
		}
	}()

	if writer.closed {
		return
	}

	buffer := bytes.NewBuffer(timeCache.format)
	buffer.WriteString(level.prefix())
	buffer.WriteString(fmt.Sprint(args...))
	writer.writer.Write(buffer.Bytes())
}

func (writer *SocketWriter) writef(level Level, format string, args ...interface{}) {
	writer.lock.Lock()

	defer func() {
		writer.lock.Unlock()

		// call log hook
		if nil != writer.hook && !(level < writer.hookLevel) {
			go func(level Level, format string, args ...interface{}) {
				writer.hook.Fire(level, fmt.Sprintf(format, args...))
			}(level, format, args...)
		}
	}()

	if writer.closed {
		return
	}

	buffer := bytes.NewBuffer(timeCache.format)
	buffer.WriteString(level.prefix())
	buffer.WriteString(fmt.Sprintf(format, args...))
	writer.writer.Write(buffer.Bytes())
}

// SetLevel set logger level
func (writer *SocketWriter) SetLevel(level Level) {
	writer.level = level
}

// SetHook set hook for logging action
func (writer *SocketWriter) SetHook(hook Hook) {
	writer.hook = hook
}

// SetHookLevel set when hook will be called
func (writer *SocketWriter) SetHookLevel(level Level) {
	writer.hookLevel = level
}

// SetTimeRotated do nothing
func (writer *SocketWriter) SetTimeRotated(timeRotated bool) {
	return
}

// SetRetentions do nothing
func (writer *SocketWriter) SetRetentions(retentions int64) {
	return
}

// SetRotateSize do nothing
func (writer *SocketWriter) SetRotateSize(rotateSize int64) {
	return
}

// SetRotateLines do nothing
func (writer *SocketWriter) SetRotateLines(rotateLines int) {
	return
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
	if TRACE < writer.level {
		return
	}

	writer.write(TRACE, args...)
}

// Tracef tracef
func (writer *SocketWriter) Tracef(format string, args ...interface{}) {
	if TRACE < writer.level {
		return
	}

	writer.writef(TRACE, format, args...)
}

// Debug debug
func (writer *SocketWriter) Debug(args ...interface{}) {
	if DEBUG < writer.level {
		return
	}

	writer.write(DEBUG, args...)
}

// Debugf debugf
func (writer *SocketWriter) Debugf(format string, args ...interface{}) {
	if DEBUG < writer.level {
		return
	}

	writer.writef(DEBUG, format, args...)
}

// Info info
func (writer *SocketWriter) Info(args ...interface{}) {
	if INFO < writer.level {
		return
	}

	writer.write(INFO, args...)
}

// Infof infof
func (writer *SocketWriter) Infof(format string, args ...interface{}) {
	if INFO < writer.level {
		return
	}

	writer.writef(INFO, format, args...)
}

// Warn warn
func (writer *SocketWriter) Warn(args ...interface{}) {
	if WARNING < writer.level {
		return
	}

	writer.write(WARNING, args...)
}

// Warnf warnf
func (writer *SocketWriter) Warnf(format string, args ...interface{}) {
	if WARNING < writer.level {
		return
	}

	writer.writef(WARNING, format, args...)
}

// Error error
func (writer *SocketWriter) Error(args ...interface{}) {
	if ERROR < writer.level {
		return
	}

	writer.write(ERROR, args...)
}

// Errorf error
func (writer *SocketWriter) Errorf(format string, args ...interface{}) {
	if ERROR < writer.level {
		return
	}

	writer.writef(ERROR, format, args...)
}

// Critical critical
func (writer *SocketWriter) Critical(args ...interface{}) {
	if CRITICAL < writer.level {
		return
	}

	writer.write(CRITICAL, args...)
}

// Criticalf criticalf
func (writer *SocketWriter) Criticalf(format string, args ...interface{}) {
	if CRITICAL < writer.level {
		return
	}

	writer.writef(CRITICAL, format, args...)
}
