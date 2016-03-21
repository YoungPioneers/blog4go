// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"fmt"
	"os"
	"path"
	"strings"
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

	// configuration about logrotate
	// exclusive lock use in logrotate
	lock *sync.Mutex

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
	// channel used to sum up sizes written from last logrotate
	logSizeChan chan int

	// number of logs retention when time base logrotate or size base logrotate
	retentions int64

	// sign decided logging with colors or not, default false
	colored bool
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
		fileName = fmt.Sprintf("%s.%s", fileName, timeCache.date)
	}
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	fileWriter.file = file
	fileWriter.currentFileName = fileName
	if nil != err {
		return fileWriter, err
	}
	fileWriter.blog = NewBLog(file)

	fileWriter.closed = false

	// about logrotate
	fileWriter.lock = new(sync.Mutex)
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

	go fileWriter.daemon()

	return fileWriter, nil
}

// NewFileWriter initialize a file writer
// baseDir must be base directory of log files
// rotate determine if it will logrotate
func NewFileWriter(baseDir string, rotate bool) (err error) {
	singltonLock.Lock()
	defer singltonLock.Unlock()
	if nil != blog {
		return ErrAlreadyInit
	}

	fileWriter := new(MultiWriter)
	fileWriter.level = DEBUG
	fileWriter.closed = false

	fileWriter.writers = make(map[Level]Writer)
	for _, level := range Levels {
		fileName := fmt.Sprintf("%s.log", strings.ToLower(level.String()))
		writer, err := newBaseFileWriter(path.Join(baseDir, fileName), rotate)
		if nil != err {
			return err
		}
		fileWriter.writers[level] = writer
	}

	// log hook
	fileWriter.hook = nil
	fileWriter.hookLevel = DEBUG

	blog = fileWriter
	return
}

// daemon run in background as NewbaseFileWriter called.
// It flushes writer buffer every 1 second.
// It decides whether a time base when logrotate is needed.
// It sums up lines && sizes already written. Alse it does the lines &&
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
			if writer.closed {
				break DaemonLoop
			}

			writer.blog.flush()
		case <-t:
			if writer.closed {
				break DaemonLoop
			}

			if writer.timeRotated {
				// if fileName not equal to currentFileName, it needs a time base logrotate
				if fileName := fmt.Sprintf("%s.%s", writer.fileName, timeCache.date); writer.currentFileName != fileName {
					// lock at this place may cause logrotate not accurate, but reduce lock acquire
					// TODO have any better solution?
					// use func to ensure writer.lock will be released
					writer.lock.Lock()
					writer.resetFile()
					writer.currentFileName = fileName
					writer.lock.Unlock()

					// when it needs to expire logs
					if writer.retentions > 0 {
						// format the expired log file name
						date := timeCache.now.Add(time.Duration(-24*(writer.retentions+1)) * time.Hour).Format(DateFormat)
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
			if writer.closed {
				break DaemonLoop
			}

			if !writer.sizeRotated && !writer.lineRotated {
				continue
			}

			writer.lock.Lock()

			writer.currentSize += ByteSize(size)
			writer.currentLines++

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
					writer.currentSize = 0
					writer.currentLines = 0
				}
			}
			writer.lock.Unlock()
		}
	}
}

// resetFile reset current writing file
func (writer *baseFileWriter) resetFile() {
	fileName := writer.fileName
	if writer.timeRotated {
		fileName = fmt.Sprintf("%s.%s", fileName, timeCache.date)
	}
	file, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	writer.file.Close()
	writer.blog.resetFile(file)
	writer.file = file
}

// write writes pure message with specific level
func (writer *baseFileWriter) write(level Level, args ...interface{}) {
	var size = 0
	defer func() {
		// logrotate
		if writer.sizeRotated || writer.lineRotated {
			writer.logSizeChan <- size
		}
	}()

	if writer.closed {
		return
	}

	size = writer.blog.write(level, args...)
}

// write formats message with specific level and write it
func (writer *baseFileWriter) writef(level Level, format string, args ...interface{}) {
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
	}()

	if writer.closed {
		return
	}

	size = writer.blog.writef(level, format, args...)
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

// SetTimeRotated toggle time base logrotate on the fly
func (writer *baseFileWriter) SetTimeRotated(timeRotated bool) {
	writer.timeRotated = timeRotated
}

// SetExpiredDays set how many days of logs will keep
func (writer *baseFileWriter) SetRetentions(retentions int64) {
	if retentions < 1 {
		return
	}
	writer.retentions = retentions
}

// SetRotateSize set size when logroatate
func (writer *baseFileWriter) SetRotateSize(rotateSize ByteSize) {
	if rotateSize > ByteSize(0) {
		writer.sizeRotated = true
		writer.rotateSize = rotateSize
	} else {
		writer.sizeRotated = false
	}
}

// SetRotateLines set line number when logrotate
func (writer *baseFileWriter) SetRotateLines(rotateLines int) {
	if rotateLines > 0 {
		writer.lineRotated = true
		writer.rotateLines = rotateLines
	} else {
		writer.lineRotated = false
	}
}

// SetColored set logging color
func (writer *baseFileWriter) SetColored(colored bool) {
	if colored == writer.colored {
		return
	}

	writer.colored = colored
	initPrefix(colored)
}

// SetLevel set logging level threshold
func (writer *baseFileWriter) SetLevel(level Level) {
	writer.blog.SetLevel(level)
}

// SetHook do nothing
func (writer *baseFileWriter) SetHook(hook Hook) {
	return
}

// SetHookLevel do nothing
func (writer *baseFileWriter) SetHookLevel(level Level) {
	return
}

// flush flush logs to disk
func (writer *baseFileWriter) flush() {
	writer.blog.flush()
}

// Trace do nothing
func (writer *baseFileWriter) Trace(args ...interface{}) {
	return
}

// Tracef do nothing
func (writer *baseFileWriter) Tracef(format string, args ...interface{}) {
	return
}

// Debug do nothing
func (writer *baseFileWriter) Debug(args ...interface{}) {
	return
}

// Debugf do nothing
func (writer *baseFileWriter) Debugf(format string, args ...interface{}) {
	return
}

// Info do nothing
func (writer *baseFileWriter) Info(args ...interface{}) {
	return
}

// Infof do nothing
func (writer *baseFileWriter) Infof(format string, args ...interface{}) {
	return
}

// Warn do nothing
func (writer *baseFileWriter) Warn(args ...interface{}) {
	return
}

// Warnf do nothing
func (writer *baseFileWriter) Warnf(format string, args ...interface{}) {
	return
}

// Error do nothing
func (writer *baseFileWriter) Error(args ...interface{}) {
	return
}

// Errorf do nothing
func (writer *baseFileWriter) Errorf(format string, args ...interface{}) {
	return
}

// Critical do nothing
func (writer *baseFileWriter) Critical(args ...interface{}) {
	return
}

// Criticalf do nothing
func (writer *baseFileWriter) Criticalf(format string, args ...interface{}) {
	return
}
