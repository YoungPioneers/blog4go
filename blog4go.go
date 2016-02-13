// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

// Package blog4go provide an efficient and easy-to-use writers library for
// logging into files, console or sockets. Writers suports formatting
// string filtering and calling user defined hook in asynchronous mode.
package blog4go

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

const (
	// EOL end of a line
	EOL = '\n'
	// ESCAPE escape character
	ESCAPE = '\\'
	// PLACEHOLDER placeholder
	PLACEHOLDER = '%'
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
	Close() // do anything end before program end

	write(level Level, format string)                       // write pure string
	writef(level Level, format string, args ...interface{}) // format string and write it
}

// BLog struct is a threadsafe log writer inherit bufio.Writer
type BLog struct {
	// logging level
	// every message level exceed this level will be written
	level Level

	// configuration about file
	// full path of the file
	fileName string
	// the file object
	file *os.File
	// bufio.Writer object of the file
	writer *bufio.Writer

	// exclusive lock while calling write function of bufio.Writer
	lock *sync.Mutex
}

// NewBLog create a BLog instance and return the pointer of it.
// fileName must be an absolute path to the destination log file
func NewBLog(fileName string) (blog *BLog, err error) {
	blog = new(BLog)
	blog.fileName = fileName
	blog.level = DEBUG
	blog.lock = new(sync.Mutex)

	// open file target file
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	if nil != err {
		return nil, err
	}
	blog.file = file
	blog.writer = bufio.NewWriterSize(file, DefaultBufferSize)
	return
}

// write writes pure message with specific level
func (log BLog) write(level Level, format string) int {
	// 统计日志size
	var size = 0

	log.lock.Lock()
	defer log.lock.Unlock()

	log.writer.Write(timeCache.format)
	log.writer.WriteString(level.Prefix())
	log.writer.WriteString(format)
	log.writer.WriteByte(EOL)

	size = len(timeCache.format) + len(level.Prefix()) + len(format) + 1
	return size
}

// write formats message with specific level and write it
func (log *BLog) writef(level Level, format string, args ...interface{}) int {
	// 格式化构造message
	// 边解析边输出
	// 使用 % 作占位符
	log.lock.Lock()
	defer log.lock.Unlock()

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

	log.writer.Write(timeCache.format)
	log.writer.WriteString(level.Prefix())

	size += len(timeCache.format) + len(level.Prefix())

	for i, v := range format {
		if tag {
			switch v {
			case 'd', 'f', 'v', 'b', 'o', 'x', 'X', 'c', 'p', 't', 's', 'T', 'q', 'U', 'e', 'E', 'g', 'G':
				if escape {
					escape = false
				}

				s, _ = log.writer.WriteString(fmt.Sprintf(format[tagPos:i+1], args[n]))
				size += s
				n++
				last = i + 1
				tag = false
			//转义符
			case ESCAPE:
				if escape {
					log.writer.WriteByte(ESCAPE)
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
				s, _ = log.writer.WriteString(format[last:i])
				size += s
				escape = false
			}
		}
	}
	log.writer.WriteString(format[last:])
	log.writer.WriteByte(EOL)

	size += len(format[last:]) + 1
	return size
}

// flush buffer to disk
func (blog *BLog) Flush() {
	blog.lock.Lock()
	defer blog.lock.Unlock()
	blog.writer.Flush()
}

// Close close file writer
func (blog *BLog) Close() {
	blog.lock.Lock()
	defer blog.lock.Unlock()

	blog.writer.Flush()
	blog.file.Close()
	blog.writer = nil
}

// FileName return file name
func (blog *BLog) FileName() string {
	return blog.fileName
}

// Level return logging level threshold
func (blog *BLog) Level() Level {
	return blog.level
}

// SetLevel set logging level threshold
func (blog *BLog) SetLevel(level Level) *BLog {
	blog.level = level
	return blog
}

// resetFile resets file descriptor of the writer
func (blog *BLog) resetFile() (err error) {
	blog.lock.Lock()
	defer blog.lock.Unlock()
	blog.writer.Flush()

	file, err := os.OpenFile(blog.fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))

	if nil != err {
		// 如果创建文件失败怎么做？
		// 重试？
		//file, err = os.OpenFile(blog.fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
		return
	}

	blog.file.Close()
	blog.file = file
	blog.writer.Reset(file)

	return
}

// resetFile resets file descriptor of the writer with specific file name
// can be used in logrotate
func (blog *BLog) resetFileWithName(fileName string) (err error) {
	blog.lock.Lock()
	defer blog.lock.Unlock()
	blog.writer.Flush()

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))

	if nil != err {
		// 如果创建文件失败怎么做？
		// 重试？
		//file, err = os.OpenFile(blog.fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
		return
	}

	blog.file.Close()
	blog.file = file
	blog.writer.Reset(file)

	return
}
