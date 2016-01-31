// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

// struct BaseFileWriter defines a writer for single file.
// It suppurts partially write while formatting message, logging level filtering,
// logrotate, user defined hook for every logging action, change configuration
// on the fly and logging with colors.
type BaseFileWriter struct {
	// logging level
	// every message level exceed this level will be written
	level Level

	// configuration about file
	// full path of the file
	fileName string
	// the file object
	file *os.File
	// bufio.Writer object of the file
	writer *bufio.Writer

	// exclusive lock while calling write function of bufio.Writer
	lock *sync.Mutex

	// close sign, default false
	// set this tag true if writer is closed
	closed bool

	// configuration about logrotate
	// exclusive lock use in logrotate
	rotateLock *sync.Mutex

	// configuration about time base logrotate
	// sign of time base logrotate, default false
	// set this tag true if logrotate in time base mode
	timeRotated bool
	// signal send when time base rotate needed
	timeRotateSig chan bool

	// configuration about size && line base logrotate
	// sign of line base logrotate, default false
	// set this tag true if logrotate in line base mode
	lineRotated bool
	// line base logrotate threshold
	rotateLines int
	// total lines written from last size && line base logrotate
	currentLines int
	// sign of size base logrotate, default false
	// set this tag true if logrotate in size base mode
	sizeRotated bool
	// size rotate按行数、大小rotate, 后缀 xxx.1, xxx.2
	// signal send when size && line base logrotate
	sizeRotateSig chan bool
	// size base logrotate threshold
	rotateSize ByteSize
	// total size written after last size && line logrotate
	currentSize ByteSize
	// times of size && line base logrotate executions
	sizeRotateTimes int
	// channel used to sum up sizes written from last logrotate
	logSizeChan chan int

	// sign decided logging with colors or not, default false
	colored bool

	// configuration about user defined logging hook
	// actual hook instance
	hook Hook
	// hook is called when message level exceed level of logging action
	hookLevel Level
}

// NewBaseFileWriter create a single file writer instance and return the poionter
// of it. When any errors happened during creation, a null writer and appropriate
// will be returned.
// fileName must be an absolute path to the destination log file
func NewBaseFileWriter(fileName string) (fileWriter *BaseFileWriter, err error) {
	fileWriter = new(BaseFileWriter)
	fileWriter.fileName = fileName
	fileWriter.level = DEBUG
	fileWriter.lock = new(sync.Mutex)
	fileWriter.closed = false

	// about logrotate
	fileWriter.rotateLock = new(sync.Mutex)
	fileWriter.timeRotated = false
	fileWriter.timeRotateSig = make(chan bool)
	fileWriter.sizeRotateSig = make(chan bool)
	fileWriter.logSizeChan = make(chan int, 4096)

	fileWriter.lineRotated = false
	fileWriter.rotateSize = DefaultRotateSize
	fileWriter.currentSize = 0

	fileWriter.sizeRotated = false
	fileWriter.rotateLines = DefaultRotateLines
	fileWriter.currentLines = 0

	fileWriter.colored = false

	// log hook
	fileWriter.hook = nil
	fileWriter.hookLevel = DEBUG

	// open file target file
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	if nil != err {
		return nil, err
	}
	fileWriter.file = file
	fileWriter.writer = bufio.NewWriterSize(file, DefaultBufferSize)

	go fileWriter.daemon()

	return fileWriter, nil
}

// daemon run in background as NewBaseFileWriter called.
// It flushes writer buffer every 10 seconds.
// It update timeCache every 1  seconds. Also it decides whether a time base
// logrotate is needed. When it is needed, it just run time base logrotate.
// It analyses lines && sizes already written. Alse it does the lines &&
// size base logrotate
func (self *BaseFileWriter) daemon() {
	// tick every seconds
	// update time cache && time base logrotate
	t := time.Tick(1 * time.Second)
	// tick every 10 seconds
	// auto flush writer buffer
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

			self.rotateLock.Lock()

			now := time.Now()
			timeCache.now = now
			timeCache.format = []byte(now.Format(PrefixTimeFormat))
			date := now.Format(DateFormat)

			if self.timeRotated && date != timeCache.date {
				// need time base logrotate
				self.sizeRotateTimes = 0

				fileName := fmt.Sprintf("%s.%s", self.fileName, timeCache.date_yesterday)
				// update file descriptor of the writer
				self.lock.Lock()
				self.flush()
				os.Rename(self.fileName, fileName)
				err := self.resetFile()
				if nil == err {
					timeCache.date_yesterday = timeCache.date
					timeCache.date = now.Format(DateFormat)
				}
				self.lock.Unlock()
			}

			self.rotateLock.Unlock()
		// analyse lines && size written
		// do lines && size base logrotate
		case size := <-self.logSizeChan:
			if self.closed {
				break DaemonLoop
			}

			if !self.sizeRotated && !self.lineRotated {
				continue
			}

			self.rotateLock.Lock()

			self.currentSize += ByteSize(size)
			self.currentLines++

			if (self.sizeRotated && self.currentSize > self.rotateSize) || (self.lineRotated && self.currentLines > self.rotateLines) {
				// need lines && size base logrotate

				fileName := fmt.Sprintf("%s.%d", self.fileName, self.sizeRotateTimes+1)
				if self.timeRotated {
					fileName = fmt.Sprintf("%s.%s.%d", self.fileName, timeCache.date, self.sizeRotateTimes+1)

				}

				// update file descriptor of the writer
				self.lock.Lock()
				self.flush()
				os.Rename(self.fileName, fileName)

				err := self.resetFile()
				if nil == err {
					self.sizeRotateTimes++
					self.currentSize = 0
					self.currentLines = 0
				}
				self.lock.Unlock()
			}
			self.rotateLock.Unlock()
		}
	}
}

// resetFile resets file descriptor of the writer
func (self *BaseFileWriter) resetFile() (err error) {
	file, err := os.OpenFile(self.fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))

	if nil != err {
		// 如果创建文件失败怎么做？
		// 重试？
		//file, err = os.OpenFile(self.fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	}

	self.file.Close()
	self.file = file
	self.writer.Reset(file)

	return
}

// write writes pure message with specific level
func (self *BaseFileWriter) write(level Level, format string) {
	self.lock.Lock()
	defer func() {
		self.lock.Unlock()
		// logrotate
		if self.sizeRotated || self.lineRotated {
			self.logSizeChan <- len(timeCache.format) + len(level.Prefix()) + len(format) + 1
		}

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

// write formats message with specific level and write it
func (self *BaseFileWriter) writef(level Level, format string, args ...interface{}) {
	// 格式化构造message
	// 边解析边输出
	// 使用 % 作占位符

	self.lock.Lock()
	// 统计日志size
	var size int = 0

	defer func() {
		self.lock.Unlock()
		// logrotate
		if self.sizeRotated || self.lineRotated {
			self.logSizeChan <- size
		}

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
	var s int = 0

	self.writer.Write(timeCache.format)
	self.writer.WriteString(level.Prefix())

	// logrotate
	if self.sizeRotated || self.lineRotated {
		size += len(timeCache.format) + len(level.Prefix())
	}

	for i, v := range format {
		if tag {
			switch v {
			case 'd', 'f', 'v', 'b', 'o', 'x', 'X', 'c', 'p', 't', 's', 'T', 'q', 'U', 'e', 'E', 'g', 'G':
				if escape {
					escape = false
				}

				s, _ = self.writer.WriteString(fmt.Sprintf(format[tagPos:i+1], args[n]))
				// logrotate
				if self.sizeRotated || self.lineRotated {
					size += s
				}
				n++
				last = i + 1
				tag = false
			//转义符
			case ESCAPE:
				if escape {
					self.writer.WriteByte(ESCAPE)
					// logrotate
					if self.sizeRotated || self.lineRotated {
						size += 1
					}
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
				s, _ = self.writer.WriteString(format[last:i])
				size += s
				escape = false
			}
		}
	}
	self.writer.WriteString(format[last:])
	self.writer.WriteByte(EOL)

	if self.sizeRotated || self.lineRotated {
		size += len(format[last:]) + 1
	}
}

func (self *BaseFileWriter) SetTimeRotated(timeRotated bool) {
	self.timeRotated = timeRotated
}

func (self *BaseFileWriter) RotateSize() ByteSize {
	return self.rotateSize
}

func (self *BaseFileWriter) SetRotateSize(rotateSize ByteSize) {
	if rotateSize > ByteSize(0) {
		self.sizeRotated = true
		self.rotateSize = rotateSize
	} else {
		self.sizeRotated = false
	}
}

func (self *BaseFileWriter) RotateLine() int {
	return self.rotateLines
}

func (self *BaseFileWriter) SetRotateLines(rotateLines int) {
	if rotateLines > 0 {
		self.lineRotated = true
		self.rotateLines = rotateLines
	} else {
		self.lineRotated = false
	}
}

func (self *BaseFileWriter) Colored() bool {
	return self.colored
}

func (self *BaseFileWriter) SetColored(colored bool) {
	if colored == self.colored {
		return
	}

	self.colored = colored
	self.lock.Lock()
	defer self.lock.Unlock()

	initPrefix(colored)
}

func (self *BaseFileWriter) SetHook(hook Hook) {
	self.hook = hook
}

func (self *BaseFileWriter) SetHookLevel(level Level) {
	self.hookLevel = level
}

func (self *BaseFileWriter) Level() Level {
	return self.level
}

func (self *BaseFileWriter) SetLevel(level Level) *BaseFileWriter {
	self.level = level
	return self
}

func (self *BaseFileWriter) Close() {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.closed {
		return
	}

	self.writer.Flush()
	self.file.Close()
	self.writer = nil
	self.closed = true
}

func (self *BaseFileWriter) flush() {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.writer.Flush()
}

func (self *BaseFileWriter) Debug(format string) {
	if DEBUG < self.level {
		return
	}

	self.write(DEBUG, format)
}

func (self *BaseFileWriter) Debugf(format string, args ...interface{}) {
	if DEBUG < self.level {
		return
	}

	self.writef(DEBUG, format, args...)
}

func (self *BaseFileWriter) Trace(format string) {
	if TRACE < self.level {
		return
	}

	self.write(TRACE, format)
}

func (self *BaseFileWriter) Tracef(format string, args ...interface{}) {
	if TRACE < self.level {
		return
	}

	self.writef(TRACE, format, args...)
}

func (self *BaseFileWriter) Info(format string) {
	if INFO < self.level {
		return
	}

	self.write(INFO, format)
}

func (self *BaseFileWriter) Infof(format string, args ...interface{}) {
	if INFO < self.level {
		return
	}

	self.writef(INFO, format, args...)
}

func (self *BaseFileWriter) Error(format string) {
	if ERROR < self.level {
		return
	}

	self.write(ERROR, format)
}

func (self *BaseFileWriter) Errorf(format string, args ...interface{}) {
	if ERROR < self.level {
		return
	}

	self.writef(ERROR, format, args...)
}

func (self *BaseFileWriter) Warn(format string) {
	if WARNING < self.level {
		return
	}

	self.write(WARNING, format)
}

func (self *BaseFileWriter) Warnf(format string, args ...interface{}) {
	if WARNING < self.level {
		return
	}

	self.writef(WARNING, format, args...)
}

func (self *BaseFileWriter) Critical(format string) {
	if CRITICAL < self.level {
		return
	}

	self.write(CRITICAL, format)
}

func (self *BaseFileWriter) Criticalf(format string, args ...interface{}) {
	if CRITICAL < self.level {
		return
	}

	self.writef(CRITICAL, format, args...)
}
