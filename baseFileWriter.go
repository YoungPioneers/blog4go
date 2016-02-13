// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"fmt"
	"sync"
	"time"
)

// BaseFileWriter defines a writer for single file.
// It suppurts partially write while formatting message, logging level filtering,
// logrotate, user defined hook for every logging action, change configuration
// on the fly and logging with colors.
type BaseFileWriter struct {
	// the BLog
	blog *BLog

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
	blog, err := NewBLog(fileName)
	if nil != err {
		return nil, err
	}
	fileWriter.blog = blog
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

	go fileWriter.daemon()

	return fileWriter, nil
}

// daemon run in background as NewBaseFileWriter called.
// It flushes writer buffer every 10 seconds.
// It update timeCache every 1  seconds. Also it decides whether a time base
// logrotate is needed. When it is needed, it just run time base logrotate.
// It analyses lines && sizes already written. Alse it does the lines &&
// size base logrotate
func (writer *BaseFileWriter) daemon() {
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
			if writer.closed {
				break DaemonLoop
			}

			writer.flush()
		case <-t:
			if writer.closed {
				break DaemonLoop
			}

			writer.rotateLock.Lock()

			now := time.Now()
			timeCache.now = now
			timeCache.format = []byte(now.Format(PrefixTimeFormat))
			date := now.Format(DateFormat)

			if writer.timeRotated && date != timeCache.date {
				// need time base logrotate
				writer.sizeRotateTimes = 0

				fileName := fmt.Sprintf("%s.%s", writer.blog.FileName(), timeCache.dateYesterday)
				writer.blog.resetFileWithName(fileName)

				timeCache.dateYesterday = timeCache.date
				timeCache.date = now.Format(DateFormat)
			}

			writer.rotateLock.Unlock()
		// analyse lines && size written
		// do lines && size base logrotate
		case size := <-writer.logSizeChan:
			if writer.closed {
				break DaemonLoop
			}

			if !writer.sizeRotated && !writer.lineRotated {
				continue
			}

			writer.rotateLock.Lock()

			writer.currentSize += ByteSize(size)
			writer.currentLines++

			if (writer.sizeRotated && writer.currentSize > writer.rotateSize) || (writer.lineRotated && writer.currentLines > writer.rotateLines) {
				// need lines && size base logrotate

				fileName := fmt.Sprintf("%s.%d", writer.blog.FileName(), writer.sizeRotateTimes+1)
				if writer.timeRotated {
					fileName = fmt.Sprintf("%s.%s.%d", writer.blog.FileName(), timeCache.date, writer.sizeRotateTimes+1)
				}
				writer.blog.resetFileWithName(fileName)
				writer.sizeRotateTimes++
				writer.currentSize = 0
				writer.currentLines = 0
			}
			writer.rotateLock.Unlock()
		}
	}
}

// write writes pure message with specific level
func (writer *BaseFileWriter) write(level Level, format string) {
	var size = 0
	writer.lock.Lock()
	defer func() {
		writer.lock.Unlock()
		// logrotate
		if writer.sizeRotated || writer.lineRotated {
			writer.logSizeChan <- size
		}

		// 异步调用log hook
		if nil != writer.hook && !(level < writer.hookLevel) {
			go func(level Level, format string) {
				writer.hook.Fire(level, format)
			}(level, format)
		}
	}()

	if writer.closed {
		return
	}

	size = writer.blog.write(level, format)
}

// write formats message with specific level and write it
func (writer *BaseFileWriter) writef(level Level, format string, args ...interface{}) {
	// 格式化构造message
	// 边解析边输出
	// 使用 % 作占位符

	writer.lock.Lock()
	// 统计日志size
	var size = 0

	defer func() {
		writer.lock.Unlock()
		// logrotate
		if writer.sizeRotated || writer.lineRotated {
			writer.logSizeChan <- size
		}

		// 异步调用log hook
		if nil != writer.hook && !(level < writer.hookLevel) {
			go func(level Level, format string, args ...interface{}) {
				writer.hook.Fire(level, fmt.Sprintf(format, args...))
			}(level, format, args...)
		}
	}()

	if writer.closed {
		return
	}

	size = writer.blog.writef(level, format, args...)
}

// SetTimeRotated toggle time base logrotate on the fly
func (writer *BaseFileWriter) SetTimeRotated(timeRotated bool) {
	writer.timeRotated = timeRotated
}

// RotateSize return size threshold when logrotate
func (writer *BaseFileWriter) RotateSize() ByteSize {
	return writer.rotateSize
}

// SetRotateSize set size when logroatate
func (writer *BaseFileWriter) SetRotateSize(rotateSize ByteSize) {
	if rotateSize > ByteSize(0) {
		writer.sizeRotated = true
		writer.rotateSize = rotateSize
	} else {
		writer.sizeRotated = false
	}
}

// RotateLine return line threshold when logrotate
func (writer *BaseFileWriter) RotateLine() int {
	return writer.rotateLines
}

// SetRotateLines set line number when logrotate
func (writer *BaseFileWriter) SetRotateLines(rotateLines int) {
	if rotateLines > 0 {
		writer.lineRotated = true
		writer.rotateLines = rotateLines
	} else {
		writer.lineRotated = false
	}
}

// Colored return whether writer log with color
func (writer *BaseFileWriter) Colored() bool {
	return writer.colored
}

// SetColored set logging color
func (writer *BaseFileWriter) SetColored(colored bool) {
	if colored == writer.colored {
		return
	}

	writer.colored = colored
	writer.lock.Lock()
	defer writer.lock.Unlock()

	initPrefix(colored)
}

// SetHook set hook for every logging actions
func (writer *BaseFileWriter) SetHook(hook Hook) {
	writer.hook = hook
}

// SetHookLevel set when hook will be called
func (writer *BaseFileWriter) SetHookLevel(level Level) {
	writer.hookLevel = level
}

// Level return logging level threshold
func (writer *BaseFileWriter) Level() Level {
	return writer.blog.Level()
}

// SetLevel set logging level threshold
func (writer *BaseFileWriter) SetLevel(level Level) *BaseFileWriter {
	writer.blog.SetLevel(level)
	return writer
}

// Close close file writer
func (writer *BaseFileWriter) Close() {
	writer.lock.Lock()
	defer writer.lock.Unlock()
	if writer.closed {
		return
	}

	writer.blog.Flush()
	writer.blog.Close()
	writer.closed = true
}

// flush buffer to disk
func (writer *BaseFileWriter) flush() {
	writer.lock.Lock()
	defer writer.lock.Unlock()
	writer.blog.Flush()
}

// Debug debug
func (writer *BaseFileWriter) Debug(format string) {
	if DEBUG < writer.blog.Level() {
		return
	}

	writer.write(DEBUG, format)
}

// Debugf debugf
func (writer *BaseFileWriter) Debugf(format string, args ...interface{}) {
	if DEBUG < writer.blog.Level() {
		return
	}

	writer.writef(DEBUG, format, args...)
}

// Trace trace
func (writer *BaseFileWriter) Trace(format string) {
	if TRACE < writer.blog.Level() {
		return
	}

	writer.write(TRACE, format)
}

// Tracef tracef
func (writer *BaseFileWriter) Tracef(format string, args ...interface{}) {
	if TRACE < writer.blog.Level() {
		return
	}

	writer.writef(TRACE, format, args...)
}

// Info info
func (writer *BaseFileWriter) Info(format string) {
	if INFO < writer.blog.Level() {
		return
	}

	writer.write(INFO, format)
}

// Infof infof
func (writer *BaseFileWriter) Infof(format string, args ...interface{}) {
	if INFO < writer.blog.Level() {
		return
	}

	writer.writef(INFO, format, args...)
}

// Error error
func (writer *BaseFileWriter) Error(format string) {
	if ERROR < writer.blog.Level() {
		return
	}

	writer.write(ERROR, format)
}

// Errorf errorf
func (writer *BaseFileWriter) Errorf(format string, args ...interface{}) {
	if ERROR < writer.blog.Level() {
		return
	}

	writer.writef(ERROR, format, args...)
}

// Warn warn
func (writer *BaseFileWriter) Warn(format string) {
	if WARNING < writer.blog.Level() {
		return
	}

	writer.write(WARNING, format)
}

// Warnf warnf
func (writer *BaseFileWriter) Warnf(format string, args ...interface{}) {
	if WARNING < writer.blog.Level() {
		return
	}

	writer.writef(WARNING, format, args...)
}

// Critical critical
func (writer *BaseFileWriter) Critical(format string) {
	if CRITICAL < writer.blog.Level() {
		return
	}

	writer.write(CRITICAL, format)
}

// Criticalf criticalf
func (writer *BaseFileWriter) Criticalf(format string, args ...interface{}) {
	if CRITICAL < writer.blog.Level() {
		return
	}

	writer.writef(CRITICAL, format, args...)
}
