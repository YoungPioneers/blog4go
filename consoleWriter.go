// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

// console logger
type ConsoleWriter struct {
	// 日志登记
	level Level

	// log文件
	writer *bufio.Writer

	// 互斥锁，用于互斥调用bufio
	lock *sync.Mutex

	// writer 关闭标识
	closed bool

	// 日志等级是否带颜色输出
	// 默认false
	colored bool

	// log hook
	hook      Hook
	hookLevel Level
}

// 创建console writer
func NewConsoleWriter() (consoleWriter *ConsoleWriter, err error) {
	consoleWriter = new(ConsoleWriter)
	consoleWriter.level = DEBUG

	consoleWriter.lock = new(sync.Mutex)
	consoleWriter.closed = false

	// 日志等级颜色输出
	consoleWriter.colored = false

	// log hook
	consoleWriter.hook = nil
	consoleWriter.hookLevel = DEBUG

	consoleWriter.writer = bufio.NewWriterSize(os.Stdout, DefaultBufferSize)

	go consoleWriter.daemon()

	return consoleWriter, nil
}

// 常驻goroutine, 更新格式化的时间
func (self *ConsoleWriter) daemon() {
	// 每秒刷新时间缓存，并判断是否需要logrotate
	t := time.Tick(1 * time.Second)
	// 10秒钟自动flush一次bufio
	f := time.Tick(10 * time.Second)

DaemonLoop:
	for {
		select {
		case <-f:
			if self.closed {
				break DaemonLoop
			}

			self.flush()
		case <-t:
			if self.closed {
				break DaemonLoop
			}

			now := time.Now()
			timeCache.now = now
			timeCache.format = []byte(now.Format(PrefixTimeFormat))

		}
	}
}

func (self *ConsoleWriter) write(level Level, format string) {
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

	self.writer.Write(timeCache.format)
	self.writer.WriteString(level.Prefix())
	self.writer.WriteString(format)
	self.writer.WriteByte(EOL)
}

// 格式化构造message
// 边解析边输出
// 使用 % 作占位符
func (self *ConsoleWriter) writef(level Level, format string, args ...interface{}) {
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

	// 识别占位符标记
	var tag bool = false
	var tagPos int = 0
	// 转义字符标记
	var escape bool = false
	// 在处理的args 下标
	var n int = 0
	// 未输出的，第一个普通字符位置
	var last int = 0

	self.writer.Write(timeCache.format)
	self.writer.WriteString(level.Prefix())

	for i, v := range format {
		if tag {
			switch v {
			case 'd', 'f', 'v', 'b', 'o', 'x', 'X', 'c', 'p', 't', 's', 'T', 'q', 'U', 'e', 'E', 'g', 'G':
				if escape {
					escape = false
				}

				self.writer.WriteString(fmt.Sprintf(format[tagPos:i+1], args[n]))
				n++
				last = i + 1
				tag = false
			//转义符
			case ESCAPE:
				if escape {
					self.writer.WriteByte(ESCAPE)
				}
				escape = !escape
			//默认
			default:

			}

		} else {
			// 占位符，百分号
			if PLACEHOLDER == format[i] && !escape {
				tag = true
				tagPos = i
				self.writer.WriteString(format[last:i])
				escape = false
			}
		}
	}
	self.writer.WriteString(format[last:])
	self.writer.WriteByte(EOL)
}

func (self *ConsoleWriter) Level() Level {
	return self.level
}

func (self *ConsoleWriter) SetLevel(level Level) *ConsoleWriter {
	self.level = level
	return self
}

func (self *ConsoleWriter) Colored() bool {
	return self.colored
}

func (self *ConsoleWriter) SetColored(colored bool) {
	if colored == self.colored {
		return
	}

	self.colored = colored
	self.lock.Lock()
	defer self.lock.Unlock()

	initPrefix(colored)
}

func (self *ConsoleWriter) SetHook(hook Hook) {
	self.hook = hook
}

func (self *ConsoleWriter) SetHookLevel(level Level) {
	self.hookLevel = level
}

func (self *ConsoleWriter) Close() {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.closed {
		return
	}

	self.writer.Flush()
	self.writer = nil
	self.closed = true
}

func (self *ConsoleWriter) flush() {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.writer.Flush()
}

func (self *ConsoleWriter) Debug(format string) {
	if DEBUG < self.level {
		return
	}

	self.write(DEBUG, format)
}

func (self *ConsoleWriter) Debugf(format string, args ...interface{}) {
	if DEBUG < self.level {
		return
	}

	self.writef(DEBUG, format, args...)
}

func (self *ConsoleWriter) Trace(format string) {
	if TRACE < self.level {
		return
	}

	self.write(TRACE, format)
}

func (self *ConsoleWriter) Tracef(format string, args ...interface{}) {
	if TRACE < self.level {
		return
	}

	self.writef(TRACE, format, args...)
}

func (self *ConsoleWriter) Info(format string) {
	if INFO < self.level {
		return
	}

	self.write(INFO, format)
}

func (self *ConsoleWriter) Infof(format string, args ...interface{}) {
	if INFO < self.level {
		return
	}

	self.writef(INFO, format, args...)
}

func (self *ConsoleWriter) Error(format string) {
	if ERROR < self.level {
		return
	}

	self.write(ERROR, format)
}

func (self *ConsoleWriter) Errorf(format string, args ...interface{}) {
	if ERROR < self.level {
		return
	}

	self.writef(ERROR, format, args...)
}

func (self *ConsoleWriter) Warn(format string) {
	if WARNING < self.level {
		return
	}

	self.write(WARNING, format)
}

func (self *ConsoleWriter) Warnf(format string, args ...interface{}) {
	if WARNING < self.level {
		return
	}

	self.writef(WARNING, format, args...)
}

func (self *ConsoleWriter) Critical(format string) {
	if CRITICAL < self.level {
		return
	}

	self.write(CRITICAL, format)
}

func (self *ConsoleWriter) Criticalf(format string, args ...interface{}) {
	if CRITICAL < self.level {
		return
	}

	self.writef(CRITICAL, format, args...)
}
