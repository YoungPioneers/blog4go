// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"fmt"
	"path"
	"strings"
)

// struct FileWriter defines a writer for multi-files writer with different
// message level
type FileWriter struct {
	// 日志等级
	level Level

	// file writers
	writers map[Level]*BaseFileWriter

	// 关闭标识
	closed bool
}

func NewFileWriter(baseDir string) (fileWriter *FileWriter, err error) {
	fileWriter = new(FileWriter)
	fileWriter.level = DEBUG
	fileWriter.closed = false

	fileWriter.writers = make(map[Level]*BaseFileWriter)
	for _, level := range Levels {
		fileName := fmt.Sprintf("%s.log", strings.ToLower(level.String()))
		writer, err := NewBaseFileWriter(path.Join(baseDir, fileName))
		if nil != err {
			return nil, err
		}
		fileWriter.writers[level] = writer
	}

	return
}

func (self *FileWriter) SetTimeRotated(timeRotated bool) {
	for _, fileWriter := range self.writers {
		fileWriter.SetTimeRotated(timeRotated)
	}
}

func (self *FileWriter) SetRotateSize(rotateSize ByteSize) {
	for _, fileWriter := range self.writers {
		fileWriter.SetRotateSize(rotateSize)
	}
}

func (self *FileWriter) SetRotateLines(rotateLines int) {
	for _, fileWriter := range self.writers {
		fileWriter.SetRotateLines(rotateLines)
	}
}

func (self *FileWriter) SetColored(colored bool) {
	for _, fileWriter := range self.writers {
		fileWriter.SetColored(colored)
	}
}

func (self *FileWriter) SetHook(hook Hook) {
	for _, fileWriter := range self.writers {
		fileWriter.SetHook(hook)
	}
}

func (self *FileWriter) SetLevel(level Level) *FileWriter {
	self.level = level
	for _, fileWriter := range self.writers {
		fileWriter.SetLevel(level)
	}
	return self
}

func (self *FileWriter) Level() Level {
	return self.level
}

func (self *FileWriter) Close() {
	for _, fileWriter := range self.writers {
		fileWriter.Close()
	}
	self.closed = true
}

func (self *FileWriter) Debug(format string) {
	if DEBUG < self.level {
		return
	}

	self.writers[DEBUG].write(DEBUG, format)
}

func (self *FileWriter) Debugf(format string, args ...interface{}) {
	if DEBUG < self.level {
		return
	}

	self.writers[DEBUG].writef(DEBUG, format, args...)
}

func (self *FileWriter) Trace(format string) {
	if TRACE < self.level {
		return
	}

	self.writers[TRACE].write(TRACE, format)
}

func (self *FileWriter) Tracef(format string, args ...interface{}) {
	if TRACE < self.level {
		return
	}

	self.writers[TRACE].writef(TRACE, format, args...)
}

func (self *FileWriter) Info(format string) {
	if INFO < self.level {
		return
	}

	self.writers[INFO].write(INFO, format)
}

func (self *FileWriter) Infof(format string, args ...interface{}) {
	if INFO < self.level {
		return
	}

	self.writers[INFO].writef(INFO, format, args...)
}

func (self *FileWriter) Warn(format string) {
	if WARNING < self.level {
		return
	}

	self.writers[WARNING].write(WARNING, format)
}

func (self *FileWriter) Warnf(format string, args ...interface{}) {
	if WARNING < self.level {
		return
	}

	self.writers[WARNING].writef(WARNING, format, args...)
}

func (self *FileWriter) Error(format string) {
	if ERROR < self.level {
		return
	}

	self.writers[ERROR].write(ERROR, format)
}

func (self *FileWriter) Errorf(format string, args ...interface{}) {
	if ERROR < self.level {
		return
	}

	self.writers[ERROR].writef(ERROR, format, args...)
}

func (self *FileWriter) Critical(format string) {
	if CRITICAL < self.level {
		return
	}

	self.writers[CRITICAL].write(CRITICAL, format)
}

func (self *FileWriter) Criticalf(format string, args ...interface{}) {
	if CRITICAL < self.level {
		return
	}

	self.writers[CRITICAL].writef(CRITICAL, format, args...)
}
