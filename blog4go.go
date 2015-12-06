// Copyright 2015
// Author: huangjunwei@youmi.net

package blog4go

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

// 各种日志结构接口
type LogWriter interface {
	// 关闭log writer的处理方法
	// 善后
	Close()

	// 用于内部写log的方法
	write(level Level, format string)
	writef(level Level, format string, args ...interface{})

	// 供用户强制刷日志到输出
	Flush()
}

const (
	// 好像buffer size 调大点benchmark效果更好
	DefaultBufferSize = 4096
	// 浮点数默认精确值
	DefaultPrecise = 2

	// 换行符
	// End Of Line
	EOL = "\n"

	// 时间前缀的格式
	PrefixTimeFormat = "[2006/01/02:15:04:05]"
)

var (
	ErrInvalidFormat = errors.New("Invalid format type.")
)

// 时间格式化的cache
type timeFormatCacheType struct {
	now time.Time
	// 时间格式化结果
	format string
}

// 用全局的timeCache好像比较好
// 实例的timeCache没那么好统一更新
var timeCache = timeFormatCacheType{}

// 装逼的logger
type FileLogWriter struct {
	level Level

	// log文件
	filename string
	file     *os.File
	writer   *bufio.Writer

	// 互斥锁，用户互斥调用bufio
	lock *sync.Mutex

	// writer 关闭标识
	closed bool
}

// 包初始化函数
func init() {
	timeCache.now = time.Now()
	timeCache.format = timeCache.now.Format(PrefixTimeFormat)
}

// 创建file writer
func NewFileLogWriter(filename string) (fileWriter *FileLogWriter, err error) {
	fileWriter = new(FileLogWriter)
	fileWriter.filename = filename

	// 打开文件描述符
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	if nil != err {
		return nil, err
	}
	fileWriter.file = file
	fileWriter.writer = bufio.NewWriterSize(file, DefaultBufferSize)
	fileWriter.lock = new(sync.Mutex)
	fileWriter.closed = false

	go fileWriter.daemon()

	return fileWriter, nil
}

func (self *FileLogWriter) SetLevel(level Level) *FileLogWriter {
	self.level = level
	return self
}

func (self *FileLogWriter) Level() Level {
	return self.level
}

func (self *FileLogWriter) Close() {
	self.lock.Lock()
	self.writer.Flush()
	self.file.Close()
	self.writer = nil
	self.closed = true
	self.lock.Unlock()
}

func (self *FileLogWriter) Flush() {
	self.writer.Flush()
}

// 常驻goroutine, 更新格式化的时间及logrotate
func (self *FileLogWriter) daemon() {
	t := time.Tick(1 * time.Second)

DaemonLoop:
	for {
		select {
		case <-t:
			if self.closed {
				break DaemonLoop
			}

			now := time.Now()
			timeCache.now = now
			timeCache.format = now.Format(PrefixTimeFormat)
		}
	}
}

func (self *FileLogWriter) write(level Level, format string, args ...interface{}) {
	if level < self.level {
		return
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	if self.closed {
		return
	}

	self.writer.WriteString(timeCache.format)
	self.writer.WriteString(level.Prefix())
	self.writer.WriteString(format)
	self.writer.WriteString(EOL)
}

// 格式化构造message
// 使用 % 作占位符
// 好像写bytes会比较快，待改进跟测试
// http://hashrocket.com/blog/posts/go-performance-observations
func (self *FileLogWriter) writef(level Level, format string, args ...interface{}) (err error) {
	if level < self.level {
		return
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	if self.closed {
		return
	}

	self.writer.WriteString(timeCache.format)
	self.writer.WriteString(level.Prefix())

	// 识别占位符标记
	var tag bool = false
	var tagPos int = 0
	// 转义字符标记
	var escape bool = false
	// 在处理的args 下标
	var n int = 0
	// 未输出的，第一个普通字符位置
	var last int = 0 //

	for i, v := range format {
		if tag {
			switch v {
			// 类型检查/ 特殊字符处理
			// 占位符，有意义部分
			// 字符串
			// %s
			case 's':
				if escape {
					escape = false
				}

				if str, ok := args[n].(string); ok {
					self.writer.WriteString(str)
					n++
					last = i + 1
				} else {
					return ErrInvalidFormat
				}
				tag = false
			// 整型
			// %d
			case 'd':
				if escape {
					escape = false
				}

				// 判断数据类型
				if number, ok := args[n].(int); ok {
					self.writer.WriteString(strconv.Itoa(number))
					n++
					last = i + 1
				} else {
					return ErrInvalidFormat
				}
				tag = false
			// 浮点型
			// %.xf
			case 'f':
				if escape {
					escape = false
				}

				// 还没想到好的解决方案，先用fmt自带的
				self.writer.WriteString(fmt.Sprintf(format[tagPos:i+1], args[n]))
				n++
				last = i + 1
				tag = false
			// Value
			// {xxx:xxx}
			case 'v':
				if escape {
					escape = false
				}

				// 还没想到好的解决方案，先用fmt自带的
				self.writer.WriteString(fmt.Sprintf(format[tagPos:i+1], args[n]))
				n++
				last = i + 1
				tag = false
			// 布尔型
			// %t
			case 't':
				if escape {
					escape = false
				}

				if b, ok := args[n].(bool); ok {
					self.writer.WriteString(strconv.FormatBool(b))
					n++
					last = i + 1
				} else {
					return ErrInvalidFormat
				}
				tag = false
			//转义符
			case '\\':
				if escape {
					self.writer.WriteString("\\")
				}
				escape = !escape
			//默认
			default:

			}

		} else {
			// 占位符，百分号
			if '%' == format[i] && !escape {
				tag = true
				tagPos = i
				self.writer.WriteString(format[last:i])
				escape = false
			}
		}
	}
	self.writer.WriteString(format[last:])
	self.writer.WriteString(EOL)

	return
}

func (self *FileLogWriter) Debug(format string) {
	self.write(DEBUG, format)
}

func (self *FileLogWriter) Debugf(format string, args ...interface{}) {
	self.writef(DEBUG, format, args...)
}

func (self *FileLogWriter) Trace(format string) {
	self.write(TRACE, format)
}

func (self *FileLogWriter) Tracef(format string, args ...interface{}) {
	self.writef(TRACE, format, args...)
}

func (self *FileLogWriter) Info(format string) {
	self.write(INFO, format)
}

func (self *FileLogWriter) Infof(format string, args ...interface{}) {
	self.writef(INFO, format, args...)
}

func (self *FileLogWriter) Warn(format string) {
	self.write(WARNING, format)
}

func (self *FileLogWriter) Warnf(format string, args ...interface{}) {
	self.writef(WARNING, format, args...)
}

func (self *FileLogWriter) Error(format string) {
	self.write(ERROR, format)
}

func (self *FileLogWriter) Errorf(format string, args ...interface{}) {
	self.writef(ERROR, format, args...)
}

func (self *FileLogWriter) Critical(format string) {
	self.write(CRITICAL, format)
}

func (self *FileLogWriter) Criticalf(format string, args ...interface{}) {
	self.writef(CRITICAL, format, args...)
}
