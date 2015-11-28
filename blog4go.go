// Copyright 2015
// Author: huangjunwei@youmi.net

package blog4go

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"strconv"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	TRACE
	INFO
	WARNING
	ERROR
	CRITICAL
)

var (
	levelStrings = [...]string{"DEBUG", "TRAC", "INFO", "WARN", "ERROR", "CRITAL"}
)

func (self Level) valid() bool {
	if DEBUG > self || CRITICAL < self {
		return false
	}
	return true
}

func (self Level) String() string {
	if !self.valid() {
		return "UNKNOWN"
	}
	return levelStrings[self]
}

// 单条日志记录结构体
type LogRecord struct {
	level   Level
	message string
}

func (self *LogRecord) String() string {
	var b bytes.Buffer
	now := time.Now().Format("2006-01-02 15:04:05")
	b.WriteString(now)
	b.WriteString(" [" + self.level.String() + "] ")
	b.WriteString(self.message)

	return b.String()
}

// 各种日志结构接口
type LogWriter interface {
	// 关闭log writer的处理方法
	// 善后
	Close()

	// 用于内部写log的方法
	write(level Level, format string, args ...interface{})

	// 供用户强制刷日志到输出
	Flush()
}

// DefaultFileLogWriter.c 为无缓冲channel
var DefaultFileLogWriter *FileLogWriter = new(FileLogWriter)

const (
	DefaultBufferSize = 4096
)

// 装逼的logger
type FileLogWriter struct {
	level Level

	// log文件
	filename string
	file     *os.File
	writer   *bufio.Writer

	// logrotate
	rotate bool
}

// 包初始化函数
func init() {
	timeCache.now = time.Now()
	timeCache.format = timeCache.now.Format("[2006/01/02:15:04:05]")
}

func NewFileLogWriter(filename string) (fileWriter *FileLogWriter, err error) {
	fileWriter = new(FileLogWriter)

	// 打开文件描述符
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	if nil != err {
		return nil, err
	}
	fileWriter.filename = filename
	fileWriter.file = file
	fileWriter.writer = bufio.NewWriterSize(file, DefaultBufferSize)

	return fileWriter, nil
}

func (self *FileLogWriter) Close() {
	self.writer.Flush()
	self.file.Close()
}

func (self *FileLogWriter) Flush() {
	self.writer.Flush()
}

func (self *FileLogWriter) SetLevel(level Level) *FileLogWriter {
	self.level = level
	return self
}

func (self *FileLogWriter) GetLevel() Level {
	return self.level
}

// 时间格式化的cache
type timeFormatCacheType struct {
	now    time.Time
	format string
}

var timeCache = timeFormatCacheType{}

func (self *FileLogWriter) write(level Level, format string, args ...interface{}) {
	if level < self.level {
		return
	}

	// 尝试缓存time format, 高并发的时候或许有用
	// 可以尝试独立goroutine每秒修改timeCache
	now := time.Now()
	if now != timeCache.now {
		timeCache.now = now
		timeCache.format = now.Format("[2006/01/02:15:04:05]")
	}

	self.writer.WriteString(timeCache.format)
	self.writer.WriteString(" [" + level.String() + "] ")
	self.writer.WriteString(format + "\n")

	//self.writer.Flush()
}

func (self *FileLogWriter) writef(level Level, format string, args ...interface{}) (err error) {
	if level < self.level {
		return
	}
	// 格式化构造message
	// 使用 % 作占位符

	// 尝试缓存time format, 高并发的时候或许有用
	// 可以尝试独立goroutine每秒修改timeCache
	now := time.Now()
	if now != timeCache.now {
		timeCache.now = now
		timeCache.format = now.Format("[2006/01/02:15:04:05]")
	}

	self.writer.WriteString(timeCache.format)
	self.writer.WriteString(" [" + level.String() + "] ")

	// 识别占位符标记
	var tag bool = false
	// 转义字符标记
	var escape bool = false

	var n int = 0
	var last int = 0

	for i := 0; i < len(format); i++ {
		if tag {
			switch format[i] {
			// 类型检查/ 特殊字符处理
			// 占位符，有意义部分
			// 字符串
			case 's':
				if escape {

					escape = false
				}

				if str, ok := args[n].(string); ok {
					self.writer.WriteString(str)
					n++
					last = i + 1
				} else {
					return errors.New("Wrong format type.")
				}
				tag = false
			// 整型
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
					return errors.New("Wrong format type.")
				}
				tag = false
			// 浮点型
			// 暂时不支持
			case 'f':
				if escape {

					escape = false
				}

				if f, ok := args[n].(float64); ok {
					// 获取精确度

					prec, _ := strconv.ParseInt(format[i-1:i], 10, 0)
					// %.xf
					self.writer.WriteString(strconv.FormatFloat(f, 'f', int(prec), 64))
					n++
					last = i + 1
				} else {
					return errors.New("Wrong format type.")
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
				self.writer.WriteString(format[last:i])
			}
		}
	}
	self.writer.WriteString(format[last:])
	self.writer.WriteString("\n")

	return
}

func (self *FileLogWriter) Debug(format string) {
	self.write(DEBUG, format)
}

func (self *FileLogWriter) Debugf(format string, args ...interface{}) {
	self.writef(DEBUG, format, args...)
}
