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
}

var (
	// 好像buffer size 调大点benchmark效果更好
	// 默认使用内存页大小
	DefaultBufferSize = 4096

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
	closed bool

	// logrotate
	rotated bool
	// size rotate按行数、大小rotate, 后缀 xxx.1, xxx.2
	timeRotateSig chan bool
	sizeRotateSig chan bool

	// 按行rotate
	lineRotated  bool
	rotateLines  int
	currentLines int

	// 按大小rotate
	sizeRotated bool
	rotateSize  ByteSize
	currentSize ByteSize

	sizeRotateTimes int // 当前按size,line rotate次数

	// 记录每次log的size
	logSizeChan chan int

	// 日志等级是否带颜色输出
	colored bool
}

// 包初始化函数
func init() {
	DefaultBufferSize = os.Getpagesize()
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

	// 日志轮询
	fileWriter.rotated = rotated
	fileWriter.timeRotateSig = make(chan bool)
	fileWriter.sizeRotateSig = make(chan bool)
	fileWriter.logSizeChan = make(chan int, 10)

	fileWriter.lineRotated = false
	fileWriter.rotateSize = DefaultRotateSize
	fileWriter.currentSize = 0

	fileWriter.sizeRotated = false
	fileWriter.rotateLines = DefaultRotateLines
	fileWriter.currentLines = 0

	// 日志等级颜色输出
	fileWriter.colored = true

	// 打开文件描述符
	// TODO 文件名或许可以改成rotate之后才加后缀
	if rotated {
		filename = fmt.Sprintf("%s.%s", filename, timeCache.date)
		go fileWriter.rotate()
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

func (self *FileLogWriter) RotateSize() ByteSize {
	return self.rotateSize
}

func (self *FileLogWriter) SetRotateSize(rotateSize ByteSize) {
	if rotateSize > ByteSize(0) {
		self.sizeRotated = true
		self.rotateSize = rotateSize
	}
}

func (self *FileLogWriter) Colored() bool {
	return self.colored
}

func (self *FileLogWriter) SetColored(colored bool) {
	if colored == self.colored {
		return
	}

	self.colored = colored
	self.lock.Lock()
	defer self.lock.Unlock()

	initPrefix(colored)
}

func (self *FileLogWriter) RotateLine() int {
	return self.rotateLines
}

func (self *FileLogWriter) SetRotateLines(rotateLines int) {
	if rotateLines > 0 {
		self.lineRotated = true
		self.rotateLines = rotateLines
	}
}

func (self *FileLogWriter) Close() {
	self.lock.Lock()
	if self.closed {
		return
	}

	self.flush()
	self.file.Close()
	self.writer = nil
	self.closed = true
	self.lock.Unlock()
}

func (self *FileLogWriter) flush() {
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
			timeCache.format = []byte(now.Format(PrefixTimeFormat))

			date := now.Format(DateFormat)
			if self.rotated && date != timeCache.date {
				// 需要rotate
				self.timeRotateSig <- true
			}
			timeCache.date = now.Format(DateFormat)
		// 统计log write size
		case size := <-self.logSizeChan:
			if self.closed {
				break DaemonLoop
			}

			if !self.sizeRotated && !self.lineRotated {
				continue
			}

			self.currentSize += ByteSize(size)
			self.currentLines++
			if self.currentSize > self.rotateSize || self.currentLines > self.rotateLines {
				self.sizeRotateSig <- true
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
			self.lock.Lock()
			if self.closed {
				break RotateLoop
			}

			self.sizeRotateTimes++
			fileName := fmt.Sprintf("%s.%s.%d", self.filename, timeCache.date, self.sizeRotateTimes)
			self.resetFile(fileName)
			self.currentSize = 0
			self.currentLines = 0
			self.lock.Unlock()

		// 按日期轮询
		case <-self.timeRotateSig:
			self.lock.Lock()
			if self.closed {
				break RotateLoop
			}

			self.sizeRotateTimes = 0

			fileName := fmt.Sprintf("%s.%s", self.filename, timeCache.date)
			self.resetFile(fileName)
			self.lock.Unlock()
		}
	}
}

func (self *FileLogWriter) resetFile(fileName string) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))

	if nil != err {
		// 如果创建文件失败怎么做？
		// 重试？
		//self.timeRotateSig <- true
	}

	self.file.Close()
	self.file = file
	self.writer.Reset(file)

}

func (self *FileLogWriter) write(level Level, format string, args ...interface{}) {
	if level < self.level {
		return
	}

	self.lock.Lock()
	defer func() {
		self.lock.Unlock()
		// logrotate
		if self.rotated {
			self.logSizeChan <- len(timeCache.format) + len(level.Prefix()) + len(format) + 1
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

// 格式化构造message
// 边解析边输出
// 使用 % 作占位符
func (self *FileLogWriter) writef(level Level, format string, args ...interface{}) {
	if level < self.level {
		return
	}

	self.lock.Lock()
	// 统计日志size
	var size int = 0

	defer func() {
		self.lock.Unlock()
		// logrotate
		if self.rotated {
			self.logSizeChan <- size
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
	size += len(timeCache.format) + len(level.Prefix())

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
					s, _ = self.writer.WriteString(str)
					size += s
					n++
					last = i + 1
				} else {
					//return ErrInvalidFormat
				}
				tag = false
			// 整型
			// %d
			// 还没想好怎么兼容int, int32, int64
			case 'd':
				if escape {
					escape = false
				}

				s, _ = self.writer.WriteString(fmt.Sprintf(format[tagPos:i+1], args[n]))
				size += s
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
				s, _ = self.writer.WriteString(fmt.Sprintf(format[tagPos:i+1], args[n]))
				size += s
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
				s, _ = self.writer.WriteString(fmt.Sprintf(format[tagPos:i+1], args[n]))
				size += s
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
					s, _ = self.writer.WriteString(strconv.FormatBool(b))
					size += s
					n++
					last = i + 1
				} else {
					//return ErrInvalidFormat
				}
				tag = false
			//转义符
			case ESCAPE:
				if escape {
					self.writer.WriteByte(ESCAPE)
					size += 1
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
				s, _ = self.writer.WriteString(format[last:i])
				size += s
				escape = false
			}
		}
	}
	self.writer.WriteString(format[last:])
	self.writer.WriteByte(EOL)
	size += len(format[last:]) + 1
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
