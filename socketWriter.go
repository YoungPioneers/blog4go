// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"bytes"
	"fmt"
	"net"
	"sync"
)

// socket logger
type SocketWriter struct {
	// 日志登记
	level Level

	// writer 关闭标识
	closed bool

	// log hook
	hook      Hook
	hookLevel Level

	// socket
	writer net.Conn

	// 互斥锁，用于互斥调用socket writer
	lock *sync.Mutex
}

// 创建socket writer
func NewSocketWriter(network string, address string) (socketWriter *SocketWriter, err error) {
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

	return
}

func (self *SocketWriter) write(level Level, format string) {
	self.lock.Lock()

	defer func() {
		self.lock.Unlock()
		// 异步调用log hook
		if nil != self.hook && !(level < self.hookLevel) {
			go func(level Level, format string) {
				self.hook.Fire(level, format)
			}(level, format)
		}
	}()

	if self.closed {
		return
	}

	buffer := bytes.NewBuffer(timeCache.format)
	buffer.WriteString(level.Prefix())
	buffer.WriteString(format)
	self.writer.Write(buffer.Bytes())
}

// 格式化构造message
// 边解析边输出
// 使用 % 作占位符
func (self *SocketWriter) writef(level Level, format string, args ...interface{}) {
	self.lock.Lock()

	defer func() {
		self.lock.Unlock()

		// 异步调用log hook
		if nil != self.hook && !(level < self.hookLevel) {
			go func(level Level, format string, args ...interface{}) {
				self.hook.Fire(level, fmt.Sprintf(format, args...))
			}(level, format, args...)
		}
	}()

	if self.closed {
		return
	}

	buffer := bytes.NewBuffer(timeCache.format)
	buffer.WriteString(level.Prefix())
	buffer.WriteString(fmt.Sprintf(format, args...))
	self.writer.Write(buffer.Bytes())
}

func (self *SocketWriter) Level() Level {
	return self.level
}

func (self *SocketWriter) SetLevel(level Level) *SocketWriter {
	self.level = level
	return self
}

func (self *SocketWriter) SetHook(hook Hook) {
	self.hook = hook
}

func (self *SocketWriter) SetHookLevel(level Level) {
	self.hookLevel = level
}

func (self *SocketWriter) Close() {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.closed {
		return
	}

	self.writer.Close()
	self.writer = nil
	self.closed = true
}

func (self *SocketWriter) Debug(format string) {
	if DEBUG < self.level {
		return
	}

	self.write(DEBUG, format)
}

func (self *SocketWriter) Debugf(format string, args ...interface{}) {
	if DEBUG < self.level {
		return
	}

	self.writef(DEBUG, format, args...)
}

func (self *SocketWriter) Trace(format string) {
	if TRACE < self.level {
		return
	}

	self.write(TRACE, format)
}

func (self *SocketWriter) Tracef(format string, args ...interface{}) {
	if TRACE < self.level {
		return
	}

	self.writef(TRACE, format, args...)
}

func (self *SocketWriter) Info(format string) {
	if INFO < self.level {
		return
	}

	self.write(INFO, format)
}

func (self *SocketWriter) Infof(format string, args ...interface{}) {
	if INFO < self.level {
		return
	}

	self.writef(INFO, format, args...)
}

func (self *SocketWriter) Error(format string) {
	if ERROR < self.level {
		return
	}

	self.write(ERROR, format)
}

func (self *SocketWriter) Errorf(format string, args ...interface{}) {
	if ERROR < self.level {
		return
	}

	self.writef(ERROR, format, args...)
}

func (self *SocketWriter) Warn(format string) {
	if WARNING < self.level {
		return
	}

	self.write(WARNING, format)
}

func (self *SocketWriter) Warnf(format string, args ...interface{}) {
	if WARNING < self.level {
		return
	}

	self.writef(WARNING, format, args...)
}

func (self *SocketWriter) Critical(format string) {
	if CRITICAL < self.level {
		return
	}

	self.write(CRITICAL, format)
}

func (self *SocketWriter) Criticalf(format string, args ...interface{}) {
	if CRITICAL < self.level {
		return
	}

	self.writef(CRITICAL, format, args...)
}
