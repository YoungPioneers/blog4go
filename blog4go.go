// Copyright 2015
// Author: huangjunwei@youmi.net

package blog4go

import (
	"bufio"
	"errors"
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

var DefaultFileLogWriter *FileLogWriter = new(FileLogWriter)

const (
	// 好像buffer size 调大点benchmark效果更好
	DefaultBufferSize = 4096
	// 浮点数默认精确值
	DefaultPrecise = 2
)

// 装逼的logger
type FileLogWriter struct {
	level Level

	// log文件
	filename string
	file     *os.File
	writer   *bufio.Writer

	// TODO logrotate
	rotate bool

	lock   *sync.RWMutex
	closed bool
}

// 时间格式化的cache
type timeFormatCacheType struct {
	now    time.Time
	format string
}

// 包初始化函数
func init() {
	timeCache.now = time.Now()
	timeCache.format = timeCache.now.Format("[2006/01/02:15:04:05]")
}

// 创建file writer
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
	fileWriter.lock = new(sync.RWMutex)
	fileWriter.closed = false

	go fileWriter.daemon()

	return fileWriter, nil
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

func (self *FileLogWriter) SetLevel(level Level) *FileLogWriter {
	self.level = level
	return self
}

func (self *FileLogWriter) GetLevel() Level {
	return self.level
}

var timeCache = timeFormatCacheType{}

// 常驻goroutine, 更新格式化的时间及logrotate
func (self *FileLogWriter) daemon() {
	t := time.Tick(1 * time.Second)
	for {
		select {
		case <-t:
			now := time.Now()
			timeCache.now = now
			timeCache.format = now.Format("[2006/01/02:15:04:05]")

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
	self.writer.WriteString(" [")
	self.writer.WriteString(level.String())
	self.writer.WriteString("] ")
	self.writer.WriteString(format)
	self.writer.WriteString("\n")
}

func (self *FileLogWriter) writef(level Level, format string, args ...interface{}) (err error) {
	if level < self.level {
		return
	}
	// 格式化构造message
	// 使用 % 作占位符

	self.lock.Lock()
	defer self.lock.Unlock()

	if self.closed {
		return
	}

	self.writer.WriteString(timeCache.format)
	self.writer.WriteString(" [")
	self.writer.WriteString(level.String())
	self.writer.WriteString("] ")

	// 识别占位符标记
	var tag bool = false
	// 转义字符标记
	var escape bool = false
	// 在处理的args 下标
	var n int = 0
	// 未输出的，第一个普通字符位置
	var last int = 0 //

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
			// %.xf
			case 'f':
				if escape {

					escape = false
				}

				if f, ok := args[n].(float64); ok {
					// 获取精确度
					prec, err := strconv.ParseInt(format[i-1:i], 10, 0)
					if nil != err {
						// 如果f前不是数字，是%
						prec = DefaultPrecise
					}

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
