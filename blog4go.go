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

	// 时间前缀的格式
	PrefixTimeFormat = "[2006/01/02:15:04:05]"
	DateFormat       = "2006-01-02"

	// 换行符
	EOL = "\n"
	// 转移符
	ESCAPE = '\\'
)

var (
	ErrInvalidFormat = errors.New("Invalid format type.")
)

// 时间格式化的cache
type timeFormatCacheType struct {
	now time.Time
	// 日期
	date string
	// 时间格式化结果
	format []byte
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
	closed    bool
	closedSig chan bool

	// logrotate
	rotated bool
	// size rotate按行数、大小rotate, 后缀 xxx.1, xxx.2
	timeRotateSig chan bool
	sizeRotateSig chan bool

	sizeRotates  int // 当前按size rotate次数
	maxLines     int
	currentLines int
	maxSize      int
	currentSize  int
	// 记录每次log的size
	logSizeChan chan int
}

// 包初始化函数
func init() {
	timeCache.now = time.Now()
	timeCache.date = timeCache.now.Format(DateFormat)
	timeCache.format = []byte(timeCache.now.Format(PrefixTimeFormat))
}

// 创建file writer
func NewFileLogWriter(filename string, rotated bool) (fileWriter *FileLogWriter, err error) {
	fileWriter = new(FileLogWriter)
	fileWriter.filename = filename

	fileWriter.lock = new(sync.Mutex)
	fileWriter.closed = false
	fileWriter.rotated = rotated
	fileWriter.timeRotateSig = make(chan bool, 0)
	fileWriter.sizeRotateSig = make(chan bool, 0)
	fileWriter.logSizeChan = make(chan int, 0)

	// 打开文件描述符
	if rotated {
		go fileWriter.rotate()
		filename = fmt.Sprintf("%s.%s", filename, timeCache.date)
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	if nil != err {
		return nil, err
	}
	fileWriter.file = file
	fileWriter.writer = bufio.NewWriterSize(file, DefaultBufferSize)

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
	self.closedSig <- true
	self.writer = nil
	self.closed = true
	self.lock.Unlock()
}

func (self *FileLogWriter) Flush() {
	self.writer.Flush()
}

// 常驻goroutine, 更新格式化的时间
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
			timeCache.date = now.Format(DateFormat)
			timeCache.format = []byte(now.Format(PrefixTimeFormat))

			date := now.Format(DateFormat)
			if self.rotated && date != timeCache.date {
				// 需要rotate
				self.timeRotateSig <- true
			}
		}
	}
}

func (self *FileLogWriter) rotate() {
RotateLoop:
	for {
		select {
		// 按size轮询
		case <-self.sizeRotateSig:
			if self.closed {
				break RotateLoop
			}
			self.sizeRotates++
			filename := fmt.Sprintf("%s.%s.%d", self.filename, timeCache.date, self.sizeRotates)
			file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))

			if nil != err {
				// 如果创建文件失败怎么做？
			}
			self.file.Close()
			self.file = file
			self.writer.Reset(file)

		// 按日期轮询
		case <-self.timeRotateSig:
			self.sizeRotates = 0

			filename := fmt.Sprintf("%s.%s", self.filename, timeCache.date)
			file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))

			if nil != err {
				// 如果创建文件失败怎么做？
				// 重试？
				//self.timeRotateSig <- true
				//continue
			}

			self.file.Close()
			self.file = file
			self.writer.Reset(file)
		// 统计log write size
		case size := <-self.logSizeChan:
			self.currentSize += size
		// 程序退出
		case <-self.closedSig:
			break RotateLoop
		}
	}
}

func (self *FileLogWriter) write(level Level, format string, args ...interface{}) (err error) {
	if level < self.level {
		return
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	if self.closed {
		return
	}

	self.writer.Write(timeCache.format)
	self.writer.WriteString(level.Prefix())
	self.writer.WriteString(format)
	self.writer.WriteString(EOL)

	// 不logrotate退出
	if !self.rotated {
		return
	}
	return
}

// 格式化构造message
// 边解析边输出
// 使用 % 作占位符
func (self *FileLogWriter) writef(level Level, format string, args ...interface{}) (err error) {
	if level < self.level {
		return
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	if self.closed {
		return
	}

	self.writer.Write(timeCache.format)
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
			// 还没想好怎么兼容int, int32, int64
			case 'd':
				if escape {
					escape = false
				}

				// 判断数据类型
				self.writer.WriteString(fmt.Sprintf(format[tagPos:i+1], args[n]))
				n++
				last = i + 1
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

	// 不logrotate退出
	if !self.rotated {
		return
	}
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
