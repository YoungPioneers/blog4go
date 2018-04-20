// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	// unit of sizes

	_ = iota // ignore first value by assigning to blank identifier
	// KB unit of kilobyte
	KB int64 = 1 << (10 * iota)
	// MB unit of megabyte
	MB
	// GB unit of gigabyte
	GB

	// default logrotate condition

	// DefaultRotateSize is default size when size base logrotate needed
	DefaultRotateSize = 500 * MB
	// DefaultRotateLines is default lines when lines base logrotate needed
	DefaultRotateLines = 2000000 // 2 million

	// DefaultLogRetentionCount is the default days of logs to be keeped
	DefaultLogRetentionCount = 7
)

// baseFileWriter defines a writer for single file.
// It suppurts partially write while formatting message, logging level filtering,
// logrotate, user defined hook for every logging action, change configuration
// on the fly and logging with colors.
type baseFileWriter struct {
	// configuration about file
	// full path of the file, the same as configuration
	fileName string
	// current file name of the writer, may be changed with logrotate
	currentFileName string
	// the file object
	file *os.File

	// the BLog
	blog *BLog

	// close sign, default false
	// set this tag true if writer is closed
	closed bool

	// configuration about user defined logging hook
	// actual hook instance
	hook Hook
	// hook is called when message level exceed level of logging action
	hookLevel LevelType
	// it determines whether hook is called async, default true
	hookAsync bool

	// configuration about logrotate
	// exclusive lock use in logrotate
	lock *sync.RWMutex

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
	rotateSize int64
	// total size written after last size && line logrotate
	currentSize int64
	// channel used to sum up sizes written from last logrotate
	logSizeChan chan int

	// number of logs retention when time base logrotate or size base logrotate
	retentions int64

	// sign decided logging with colors or not, default false
	colored bool
}

// NewBaseFileWriter initialize a base file writer
func NewBaseFileWriter(fileName string, timeRotated bool) (err error) {
	singltonLock.Lock()
	defer singltonLock.Unlock()

	if nil != blog {
		return ErrAlreadyInit
	}

	baseFileWriter, err := newBaseFileWriter(fileName, timeRotated)
	if nil != err {
		return err
	}

	blog = baseFileWriter
	return err
}

// newbaseFileWriter create a single file writer instance and return the poionter
// of it. When any errors happened during creation, a null writer and appropriate
// will be returned.
// fileName must be an absolute path to the destination log file
// rotate determine if it will logrotate
func newBaseFileWriter(fileName string, timeRotated bool) (fileWriter *baseFileWriter, err error) {
	fileWriter = new(baseFileWriter)

	fileWriter.fileName = fileName
	// open file target file
	if timeRotated {
		fileName = fmt.Sprintf("%s.%s", fileName, timeCache.Date())
	}
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	fileWriter.file = file
	fileWriter.currentFileName = fileName
	if nil != err {
		return nil, err
	}
	fileWriter.blog = NewBLog(file)

	fileWriter.closed = false

	// about logrotate
	fileWriter.lock = new(sync.RWMutex)
	fileWriter.timeRotated = timeRotated
	fileWriter.timeRotateSig = make(chan bool)
	fileWriter.sizeRotateSig = make(chan bool)
	fileWriter.logSizeChan = make(chan int, 8192)

	fileWriter.lineRotated = false
	fileWriter.rotateSize = DefaultRotateSize
	fileWriter.currentSize = 0

	fileWriter.sizeRotated = false
	fileWriter.rotateLines = DefaultRotateLines
	fileWriter.currentLines = 0
	fileWriter.retentions = DefaultLogRetentionCount

	fileWriter.colored = false

	// log hook
	fileWriter.hook = nil
	fileWriter.hookLevel = DEBUG
	fileWriter.hookAsync = true

	go fileWriter.daemon()

	return fileWriter, nil
}

// daemon run in background as NewbaseFileWriter called.
// It flushes writer buffer every 1 second.
// It decides whether a time base when logrotate is needed.
// It sums up lines && sizes already written. Also it supports the lines &&
// size base logrotate
func (writer *baseFileWriter) daemon() {
	// tick every seconds
	// time base logrotate
	t := time.Tick(1 * time.Second)
	// tick every second
	// auto flush writer buffer
	f := time.Tick(1 * time.Second)

DaemonLoop:
	for {
		select {
		case <-f:
			if writer.Closed() {
				break DaemonLoop
			}

			writer.Flush()
		case <-t:
			if writer.Closed() {
				break DaemonLoop
			}

			if writer.timeRotated {
				// if fileName not equal to currentFileName, it needs a time base logrotate
				if fileName := fmt.Sprintf("%s.%s", writer.fileName, timeCache.Date()); writer.currentFileName != fileName {
					writer.resetFile()
					writer.currentFileName = fileName

					// when it needs to expire logs
					if writer.retentions > 0 {
						// format the expired log file name
						date := timeCache.Now().Add(time.Duration(-24*(writer.retentions+1)) * time.Hour).Format(DateFormat)
						expiredFileName := fmt.Sprintf("%s.%s", writer.fileName, date)
						// check if expired log exists
						if _, err := os.Stat(expiredFileName); nil == err {
							os.Remove(expiredFileName)
						}
					}
				}
			}

		// analyse lines && size written
		// do lines && size base logrotate
		case size := <-writer.logSizeChan:
			if writer.Closed() {
				break DaemonLoop
			}

			if !writer.sizeRotated && !writer.lineRotated {
				continue
			}

			// TODO have any better solution?
			// use func to ensure writer.lock will be released
			writer.lock.Lock()
			writer.currentSize += int64(size)
			writer.currentLines++
			writer.lock.Unlock()

			if (writer.sizeRotated && writer.currentSize >= writer.rotateSize) || (writer.lineRotated && writer.currentLines >= writer.rotateLines) {
				// need lines && size base logrotate
				var oldName, newName string
				oldName = fmt.Sprintf("%s.%d", writer.currentFileName, writer.retentions)
				// check if expired log exists
				if _, err := os.Stat(oldName); os.IsNotExist(err) {
					os.Remove(oldName)
				}
				if writer.retentions > 0 {

					for i := writer.retentions - 1; i > 0; i-- {
						oldName = fmt.Sprintf("%s.%d", writer.currentFileName, i)
						newName = fmt.Sprintf("%s.%d", writer.currentFileName, i+1)
						os.Rename(oldName, newName)
					}
					os.Rename(writer.currentFileName, oldName)

					writer.resetFile()
				}
			}
		}
	}
}

// resetFile reset current writing file
func (writer *baseFileWriter) resetFile() {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	fileName := writer.fileName
	if writer.timeRotated {
		fileName = fmt.Sprintf("%s.%s", fileName, timeCache.Date())
	}
	file, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	writer.blog.resetFile(file)
	writer.file.Close()
	writer.file = file

	writer.currentSize = 0
	writer.currentLines = 0
}

// write writes pure message with specific level
func (writer *baseFileWriter) write(level LevelType, args ...interface{}) {
	writer.lock.RLock()
	defer writer.lock.RUnlock()
	var size = 0

	if writer.closed {
		return
	}

	size = writer.blog.write(level, args...)

	// logrotate
	if writer.sizeRotated || writer.lineRotated {
		writer.logSizeChan <- size
	}

	if nil != writer.hook && !(level < writer.hookLevel) {
		if writer.hookAsync {
			// 异步调用log hook
			go writer.hook.Fire(level, writer.blog.Tags(), args...)
		} else {
			writer.hook.Fire(level, writer.blog.Tags(), args...)
		}
	}
}

// write formats message with specific level and write it
func (writer *baseFileWriter) writef(level LevelType, format string, args ...interface{}) {
	// 格式化构造message
	// 边解析边输出
	// 使用 % 作占位符

	// 统计日志size
	var size = 0

	if writer.closed {
		return
	}

	size = writer.blog.writef(level, format, args...)

	// logrotate
	if writer.sizeRotated || writer.lineRotated {
		writer.logSizeChan <- size
	}

	if nil != writer.hook && !(level < writer.hookLevel) {
		if writer.hookAsync {
			// 异步调用log hook
			go writer.hook.Fire(level, writer.blog.Tags(), fmt.Sprintf(format, args...))
		} else {
			writer.hook.Fire(level, writer.blog.Tags(), fmt.Sprintf(format, args...))
		}
	}
}

// Closed get writer status
func (writer *baseFileWriter) Closed() bool {
	writer.lock.RLock()
	defer writer.lock.RUnlock()

	return writer.closed
}

// Close close file writer
func (writer *baseFileWriter) Close() {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	if writer.closed {
		return
	}

	writer.closed = true
	writer.blog.flush()
	writer.blog.Close()
	writer.blog = nil
	writer.file.Close()
	close(writer.logSizeChan)
	close(writer.timeRotateSig)
	close(writer.sizeRotateSig)
}

// TimeRotated get timeRotated
func (writer *baseFileWriter) TimeRotated() bool {
	writer.lock.RLock()
	defer writer.lock.RUnlock()
	return writer.timeRotated
}

// SetTimeRotated toggle time base logrotate on the fly
func (writer *baseFileWriter) SetTimeRotated(timeRotated bool) {
	writer.lock.Lock()
	defer writer.lock.Unlock()
	writer.timeRotated = timeRotated
}

// Retentions get log retention days
func (writer *baseFileWriter) Retentions() int64 {
	writer.lock.RLock()
	defer writer.lock.RUnlock()
	return writer.retentions
}

// SetExpiredDays set how many days of logs will keep
func (writer *baseFileWriter) SetRetentions(retentions int64) {
	writer.lock.Lock()
	defer writer.lock.Unlock()
	if retentions < 1 {
		return
	}
	writer.retentions = retentions
}

// RotateSize get log rotate size
func (writer *baseFileWriter) RotateSize() int64 {
	writer.lock.RLock()
	defer writer.lock.RUnlock()
	return writer.rotateSize
}

// SetRotateSize set size when logroatate
func (writer *baseFileWriter) SetRotateSize(rotateSize int64) {
	writer.lock.Lock()
	defer writer.lock.Unlock()
	if rotateSize > 0 {
		writer.sizeRotated = true
		writer.rotateSize = rotateSize
	} else {
		writer.sizeRotated = false
	}
}

// RotateLines get log rotate lines
func (writer *baseFileWriter) RotateLines() int {
	writer.lock.RLock()
	defer writer.lock.RUnlock()
	return writer.rotateLines
}

// SetRotateLines set line number when logrotate
func (writer *baseFileWriter) SetRotateLines(rotateLines int) {
	writer.lock.Lock()
	defer writer.lock.Unlock()
	if rotateLines > 0 {
		writer.lineRotated = true
		writer.rotateLines = rotateLines
	} else {
		writer.lineRotated = false
	}
}

// Colored get whether it is log with colored
func (writer *baseFileWriter) Colored() bool {
	writer.lock.RLock()
	defer writer.lock.RUnlock()
	return writer.colored
}

// SetColored set logging color
func (writer *baseFileWriter) SetColored(colored bool) {
	writer.lock.Lock()
	defer writer.lock.Unlock()
	if colored == writer.colored {
		return
	}

	writer.colored = colored
	initPrefix(colored)
}

// Level get log level
func (writer *baseFileWriter) Level() LevelType {
	writer.lock.RLock()
	defer writer.lock.RUnlock()
	return writer.blog.Level()
}

// SetLevel set logging level threshold
func (writer *baseFileWriter) SetLevel(level LevelType) {
	writer.lock.Lock()
	defer writer.lock.Unlock()
	writer.blog.SetLevel(level)
}

// SetHook set hook for the base file writer
func (writer *baseFileWriter) SetHook(hook Hook) {
	writer.lock.Lock()
	defer writer.lock.Unlock()
	writer.hook = hook
}

// Tags return logging tags
func (writer *baseFileWriter) Tags() map[string]string {
	writer.lock.RLock()
	defer writer.lock.RUnlock()
	return writer.blog.Tags()
}

// SetTags set logging tags
func (writer *baseFileWriter) SetTags(tags map[string]string) {
	writer.lock.Lock()
	defer writer.lock.Unlock()
	writer.blog.SetTags(tags)
}

// SetHookAsync set hook async for base file writer
func (writer *baseFileWriter) SetHookAsync(async bool) {
	writer.lock.Lock()
	defer writer.lock.Unlock()
	writer.hookAsync = async
}

// SetHookLevel set when hook will be called
func (writer *baseFileWriter) SetHookLevel(level LevelType) {
	writer.lock.Lock()
	defer writer.lock.Unlock()
	writer.hookLevel = level
}

// Flush flush logs to disk
func (writer *baseFileWriter) Flush() {
	writer.lock.RLock()
	defer writer.lock.RUnlock()

	writer.flush()
}

// flush flush logs to disk
func (writer *baseFileWriter) flush() {
	writer.blog.flush()
}

// Trace trace
func (writer *baseFileWriter) Trace(args ...interface{}) {
	if nil == writer.blog || TRACE < writer.blog.Level() {
		return
	}

	writer.write(TRACE, args...)
}

// Tracef tracef
func (writer *baseFileWriter) Tracef(format string, args ...interface{}) {
	if nil == writer.blog || TRACE < writer.blog.Level() {
		return
	}

	writer.writef(TRACE, format, args...)
}

// Debug debug
func (writer *baseFileWriter) Debug(args ...interface{}) {
	if nil == writer.blog || DEBUG < writer.blog.Level() {
		return
	}

	writer.write(DEBUG, args...)
}

// Debugf debugf
func (writer *baseFileWriter) Debugf(format string, args ...interface{}) {
	if nil == writer.blog || DEBUG < writer.blog.Level() {
		return
	}

	writer.writef(DEBUG, format, args...)
}

// Info info
func (writer *baseFileWriter) Info(args ...interface{}) {
	if nil == writer.blog || INFO < writer.blog.Level() {
		return
	}

	writer.write(INFO, args...)
}

// Infof infof
func (writer *baseFileWriter) Infof(format string, args ...interface{}) {
	if nil == writer.blog || INFO < writer.blog.Level() {
		return
	}

	writer.writef(INFO, format, args...)
}

// Warn warn
func (writer *baseFileWriter) Warn(args ...interface{}) {
	if nil == writer.blog || WARNING < writer.blog.Level() {
		return
	}

	writer.write(WARNING, args...)
}

// Warnf warn
func (writer *baseFileWriter) Warnf(format string, args ...interface{}) {
	if nil == writer.blog || WARNING < writer.blog.Level() {
		return
	}

	writer.writef(WARNING, format, args...)
}

// Error error
func (writer *baseFileWriter) Error(args ...interface{}) {
	if nil == writer.blog || ERROR < writer.blog.Level() {
		return
	}

	writer.write(ERROR, args...)
}

// Errorf errorf
func (writer *baseFileWriter) Errorf(format string, args ...interface{}) {
	if nil == writer.blog || ERROR < writer.blog.Level() {
		return
	}

	writer.writef(ERROR, format, args...)
}

// Critical critical
func (writer *baseFileWriter) Critical(args ...interface{}) {
	if nil == writer.blog || CRITICAL < writer.blog.Level() {
		return
	}

	writer.write(CRITICAL, args...)
}

// Criticalf criticalf
func (writer *baseFileWriter) Criticalf(format string, args ...interface{}) {
	if nil == writer.blog || CRITICAL < writer.blog.Level() {
		return
	}

	writer.writef(CRITICAL, format, args...)
}
