// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

// Package blog4go provide an efficient and easy-to-use writers library for
// logging into files, console or sockets. Writers suports formatting
// string filtering and calling user defined hook in asynchronous mode.
package blog4go

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

const (
	// EOL end of a line
	EOL = '\n'
	// ESCAPE escape character
	ESCAPE = '\\'
	// PLACEHOLDER placeholder
	PLACEHOLDER = '%'
	// QUOTE quote character
	QUOTE = '"'
	// SPACE space character
	SPACE = ' '
)

var (
	// blog is the singleton instance use for blog.write/writef
	blog Writer

	// global mutex log used for singlton
	singltonLock *sync.RWMutex

	// DefaultBufferSize bufio buffer size
	DefaultBufferSize = 4096 // default memory page size
	// ErrInvalidFormat invalid format error
	ErrInvalidFormat = errors.New("Invalid format type")
	// ErrAlreadyInit show that blog is already initialized once
	ErrAlreadyInit = errors.New("blog4go has been already initialized")
)

// Writer interface is a common definition of any writers in this package.
// Any struct implements Writer interface must implement functions below.
// Close is used for close the writer and free any elements if needed.
// write is an internal function that write pure message with specific
// logging level.
// writef is an internal function that formatting message with specific
// logging level. Placeholders in the format string will be replaced with
// args given.
// Both write and writef may have an asynchronous call of user defined
// function before write and writef function end..
type Writer interface {
	// Close do anything end before program end
	Close()

	// SetLevel set logging level threshold
	SetLevel(level LevelType)
	// Level get log level
	Level() LevelType

	// write/writef functions with different levels
	write(level LevelType, args ...interface{})
	writef(level LevelType, format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Critical(args ...interface{})
	Criticalf(format string, args ...interface{})

	// flush log to disk
	flush()

	// hook
	SetHook(hook Hook)
	SetHookLevel(level LevelType)
	SetHookAsync(async bool)

	// logrotate
	SetTimeRotated(timeRotated bool)
	TimeRotated() bool
	SetRotateSize(rotateSize int64)
	RotateSize() int64
	SetRotateLines(rotateLines int)
	RotateLines() int
	SetRetentions(retentions int64)
	Retentions() int64
	SetColored(colored bool)
	Colored() bool

	// tags
	SetTags(tags map[string]string)
	Tags() map[string]string
}

func init() {
	singltonLock = new(sync.RWMutex)
	DefaultBufferSize = os.Getpagesize()
}

// NewWriterFromConfigAsFile initialize a writer according to given config file
// configFile must be the path to the config file
func NewWriterFromConfigAsFile(configFile string) (err error) {
	singltonLock.Lock()
	defer singltonLock.Unlock()
	if nil != blog {
		return ErrAlreadyInit
	}

	// read config from file
	config, err := readConfig(configFile)
	if nil != err {
		return
	}

	if err = config.valid(); nil != err {
		return
	}

	multiWriter := new(MultiWriter)
	multiWriter.lock = new(sync.RWMutex)

	multiWriter.level = DEBUG
	if level := LevelFromString(config.MinLevel); level.valid() {
		multiWriter.level = level
	}

	multiWriter.closed = false
	multiWriter.writers = make(map[LevelType]Writer)

	for _, filter := range config.Filters {
		var rotate = false
		var timeRotate = false
		var isSocket = false
		var isConsole = false

		// get file path
		var filePath string
		if (file{}) != filter.File {
			// file do not need logrotate
			filePath = filter.File.Path
			rotate = false
		} else if (rotateFile{}) != filter.RotateFile {
			// file need logrotate
			filePath = filter.RotateFile.Path
			rotate = true
			timeRotate = TypeTimeBaseRotate == filter.RotateFile.Type
		} else if (socket{}) != filter.Socket {
			isSocket = true
		} else {
			// use console writer as default
			isConsole = true
		}

		levels := strings.Split(filter.Levels, ",")
		for _, levelStr := range levels {
			var level LevelType
			if level = LevelFromString(levelStr); !level.valid() {
				return ErrInvalidLevel
			}

			if isConsole {
				// console writer
				writer, err := newConsoleWriter(filter.Console.Redirect)
				if nil != err {
					return err
				}

				multiWriter.writers[level] = writer
				continue
			}

			if isSocket {
				// socket writer
				writer, err := newSocketWriter(filter.Socket.Network, filter.Socket.Address)
				if nil != err {
					return err
				}

				multiWriter.writers[level] = writer
				continue
			}

			// init a base file writer
			writer, err := newBaseFileWriter(filePath, timeRotate)
			if nil != err {
				return err
			}

			if rotate {
				// set logrotate strategy
				if TypeTimeBaseRotate == filter.RotateFile.Type {
					writer.SetTimeRotated(true)
					writer.SetRetentions(filter.RotateFile.Retentions)
				} else if TypeSizeBaseRotate == filter.RotateFile.Type {
					writer.SetRotateSize(filter.RotateFile.RotateSize)
					writer.SetRotateLines(filter.RotateFile.RotateLines)
					writer.SetRetentions(filter.RotateFile.Retentions)
				} else {
					return ErrInvalidRotateType
				}
			}

			// set color
			multiWriter.SetColored(filter.Colored)
			multiWriter.writers[level] = writer
		}
	}

	blog = multiWriter
	return
}

// BLog struct is a threadsafe log writer inherit bufio.Writer
type BLog struct {
	// logging level
	// every message level exceed this level will be written
	level LevelType

	// input io
	in io.Writer

	// bufio.Writer object of the input io
	writer *bufio.Writer

	// exclusive lock while calling write function of bufio.Writer
	lock *sync.RWMutex

	// tags
	tags   map[string]string
	tagStr string

	// closed tag
	closed bool
}

// NewBLog create a BLog instance and return the pointer of it.
// fileName must be an absolute path to the destination log file
func NewBLog(in io.Writer) (blog *BLog) {
	blog = new(BLog)
	blog.in = in
	blog.level = TRACE
	blog.lock = new(sync.RWMutex)
	blog.closed = false

	blog.writer = bufio.NewWriterSize(in, DefaultBufferSize)
	return
}

// write writes pure message with specific level
func (blog *BLog) write(level LevelType, args ...interface{}) int {
	blog.lock.Lock()
	defer blog.lock.Unlock()

	// 统计日志size
	var size = 0

	format := fmt.Sprintf("msg=\"%s\" ", fmt.Sprint(args...))

	blog.writer.Write(timeCache.Format())
	blog.writer.WriteString(level.prefix())
	blog.writer.WriteString(blog.tagStr)
	blog.writer.WriteString(format)
	blog.writer.WriteByte(EOL)

	size = len(timeCache.Format()) + len(level.prefix()) + len(blog.tagStr) + len(format) + 1
	return size
}

// write formats message with specific level and write it
func (blog *BLog) writef(level LevelType, format string, args ...interface{}) int {
	// 格式化构造message
	// 边解析边输出
	// 使用 % 作占位符
	blog.lock.Lock()
	defer blog.lock.Unlock()

	// 统计日志size
	var size = 0

	// 识别占位符标记
	var tag = false
	var tagPos int
	// 转义字符标记
	var escape = false
	// 在处理的args 下标
	var n int
	// 未输出的，第一个普通字符位置
	var last int
	var s int

	blog.writer.Write(timeCache.Format())
	blog.writer.WriteString(level.prefix())
	blog.writer.WriteString(blog.tagStr)
	blog.writer.WriteString("msg=\"")

	size += len(timeCache.Format()) + len(level.prefix()) + len(blog.tagStr)

	for i, v := range format {
		if tag {
			switch v {
			case 'd', 'f', 'v', 'b', 'o', 'x', 'X', 'c', 'p', 't', 's', 'T', 'q', 'U', 'e', 'E', 'g', 'G':
				if escape {
					escape = false
				}

				s, _ = blog.writer.WriteString(fmt.Sprintf(format[tagPos:i+1], args[n]))
				size += s
				n++
				last = i + 1
				tag = false
			//转义符
			case ESCAPE:
				if escape {
					blog.writer.WriteByte(ESCAPE)
					size++
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
				s, _ = blog.writer.WriteString(format[last:i])
				size += s
				escape = false
			}
		}
	}
	blog.writer.WriteString(format[last:])
	blog.writer.WriteByte(QUOTE)
	blog.writer.WriteByte(SPACE)
	blog.writer.WriteByte(EOL)

	size += len(format[last:]) + 1
	return size
}

// Flush flush buffer to disk
func (blog *BLog) flush() {
	blog.lock.Lock()
	defer blog.lock.Unlock()

	if blog.closed {
		return
	}

	blog.writer.Flush()
}

// Close close file writer
func (blog *BLog) Close() {
	blog.lock.Lock()
	defer blog.lock.Unlock()

	if nil == blog || blog.closed {
		return
	}

	blog.closed = true
	blog.writer.Flush()
	blog.writer = nil
}

// In return the input io.Writer
func (blog *BLog) In() io.Writer {
	return blog.in
}

// Level return logging level threshold
func (blog *BLog) Level() LevelType {
	return blog.level
}

// SetLevel set logging level threshold
func (blog *BLog) SetLevel(level LevelType) *BLog {
	blog.level = level
	return blog
}

// Tags return logging tags
func (blog *BLog) Tags() map[string]string {
	blog.lock.RLock()
	defer blog.lock.RUnlock()
	return blog.tags
}

// SetTags set logging tags
func (blog *BLog) SetTags(tags map[string]string) *BLog {
	blog.lock.Lock()
	defer blog.lock.Unlock()
	blog.tags = tags

	var tagStr string
	for tagName, tagValue := range blog.tags {
		tagStr = fmt.Sprintf("%s%s=\"%s\" ", tagStr, tagName, tagValue)
	}
	blog.tagStr = tagStr

	return blog
}

// resetFile resets file descriptor of the writer with specific file name
func (blog *BLog) resetFile(in io.Writer) (err error) {
	blog.lock.Lock()
	defer blog.lock.Unlock()

	blog.writer.Flush()

	blog.in = in
	blog.writer.Reset(in)

	return
}

// SetBufferSize set bufio buffer size in bytes
func SetBufferSize(size int) {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	DefaultBufferSize = size
}

// Level get log level
func Level() LevelType {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	return blog.Level()
}

// SetLevel set level for logging action
func SetLevel(level LevelType) {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	blog.SetLevel(level)
}

// Tags return logging tags
func Tags() map[string]string {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	return blog.Tags()
}

// SetTags set logging tags
func SetTags(tags map[string]string) {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	blog.SetTags(tags)
}

// SetHook set hook for logging action
func SetHook(hook Hook) {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	blog.SetHook(hook)
}

// SetHookLevel set when hook will be called
func SetHookLevel(level LevelType) {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	blog.SetHookLevel(level)
}

// SetHookAsync set whether hook is called async
func SetHookAsync(async bool) {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	blog.SetHookAsync(async)
}

// Colored get whether it is log with colored
func Colored() bool {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	return blog.Colored()
}

// SetColored set logging color
func SetColored(colored bool) {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	blog.SetColored(colored)
}

// TimeRotated get timeRotated
func TimeRotated() bool {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	return blog.TimeRotated()
}

// SetTimeRotated toggle time base logrotate on the fly
func SetTimeRotated(timeRotated bool) {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	blog.SetTimeRotated(timeRotated)
}

// Retentions get retentions
func Retentions() int64 {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	return blog.Retentions()
}

// SetRetentions set how many logs will keep after logrotate
func SetRetentions(retentions int64) {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	blog.SetRetentions(retentions)
}

// RotateSize get rotateSize
func RotateSize() int64 {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	return blog.RotateSize()
}

// SetRotateSize set size when logroatate
func SetRotateSize(rotateSize int64) {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	blog.SetRotateSize(rotateSize)
}

// RotateLines get rotateLines
func RotateLines() int {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	return blog.RotateLines()
}

// SetRotateLines set line number when logrotate
func SetRotateLines(rotateLines int) {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	blog.SetRotateLines(rotateLines)
}

// Flush flush logs to disk
func Flush() {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	blog.flush()
}

// Trace static function for Trace
func Trace(args ...interface{}) {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	blog.Trace(args...)
}

// Tracef static function for Tracef
func Tracef(format string, args ...interface{}) {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	blog.Tracef(format, args...)
}

// Debug static function for Debug
func Debug(args ...interface{}) {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	blog.Debug(args...)
}

// Debugf static function for Debugf
func Debugf(format string, args ...interface{}) {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	blog.Debugf(format, args...)
}

// Info static function for Info
func Info(args ...interface{}) {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	blog.Info(args...)
}

// Infof static function for Infof
func Infof(format string, args ...interface{}) {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	blog.Infof(format, args...)
}

// Warn static function for Warn
func Warn(args ...interface{}) {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	blog.Warn(args...)
}

// Warnf static function for Warnf
func Warnf(format string, args ...interface{}) {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	blog.Warnf(format, args...)
}

// Error static function for Error
func Error(args ...interface{}) {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	blog.Error(args...)
}

// Errorf static function for Errorf
func Errorf(format string, args ...interface{}) {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	blog.Errorf(format, args...)
}

// Critical static function for Critical
func Critical(args ...interface{}) {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	blog.Critical(args...)
}

// Criticalf static function for Criticalf
func Criticalf(format string, args ...interface{}) {
	singltonLock.RLock()
	defer singltonLock.RUnlock()

	blog.Criticalf(format, args...)
}

// Close close the logger
func Close() {
	singltonLock.Lock()
	defer singltonLock.Unlock()

	if nil == blog {
		return
	}

	blog.Close()
	blog = nil
}
