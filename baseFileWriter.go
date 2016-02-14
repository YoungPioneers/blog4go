// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// ByteSize is type of sizes
type ByteSize int64

const (
	// unit of sizes

	_ = iota // ignore first value by assigning to blank identifier
	// KB unit of kilobyte
	KB ByteSize = 1 << (10 * iota)
	// MB unit of megabyte
	MB
	// GB unit of gigabyte
	GB

	// default logrotate condition

	// DefaultRotateSize is default size when size base logrotate needed
	DefaultRotateSize = 500 * MB
	// DefaultRotateLines is default lines when lines base logrotate needed
	DefaultRotateLines = 2000000 // 2 million
)

// BaseFileWriter defines a writer for single file.
// It suppurts partially write while formatting message, logging level filtering,
// logrotate, user defined hook for every logging action, change configuration
// on the fly and logging with colors.
type BaseFileWriter struct {
	// configuration about file
	// full path of the file
	fileName string
	// the file object
	file *os.File

	// the BLog
	blog *BLog

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
	// open file target file
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	fileWriter.file = file
	if nil != err {
		return nil, err
	}
	fileWriter.blog = NewBLog(file)

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
// It decides whether a time base when logrotate is needed.
// It sums up lines && sizes already written. Alse it does the lines &&
// size base logrotate
func (writer *BaseFileWriter) daemon() {
	// tick every seconds
	// time base logrotate
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

			writer.blog.flush()
		case <-t:
			if writer.closed {
				break DaemonLoop
			}

			writer.rotateLock.Lock()

			date := time.Now().Format(DateFormat)

			if writer.timeRotated && date != timeCache.date {
				// need time base logrotate
				writer.sizeRotateTimes = 0

				fileName := fmt.Sprintf("%s.%s", writer.fileName, timeCache.dateYesterday)
				file, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))

				writer.file.Close()
				writer.blog.resetFile(file)
				writer.file = file
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
				writer.sizeRotateTimes++
				writer.currentSize = 0
				writer.currentLines = 0

				fileName := fmt.Sprintf("%s.%d", writer.fileName, writer.sizeRotateTimes+1)
				if writer.timeRotated {
					fileName = fmt.Sprintf("%s.%s.%d", writer.fileName, timeCache.date, writer.sizeRotateTimes+1)
				}
				file, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))

				writer.file.Close()
				writer.blog.resetFile(file)
				writer.file = file
			}
			writer.rotateLock.Unlock()
		}
	}
}

// write writes pure message with specific level
func (writer *BaseFileWriter) write(level Level, format string) {
	var size = 0
	defer func() {
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

	// 统计日志size
	var size = 0

	defer func() {
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
	if writer.closed {
		return
	}

	writer.blog.Close()
	writer.blog = nil
	writer.closed = true
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
