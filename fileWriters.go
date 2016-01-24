// Copyright 2015
// Author: huangjunwei@youmi.net
package blog4go

import (
	"fmt"
	"path"
	"strings"
)

type FileWriters struct {
	// 日志登记
	level Level

	// file writers
	writers map[Level]*FileWriter

	// 关闭标识
	closed bool
}

func NewFileWriters(baseDir string) (fileWriters *FileWriters, err error) {
	fileWriters = new(FileWriters)
	fileWriters.level = DEBUG
	fileWriters.closed = false

	fileWriters.writers = make(map[Level]*FileWriter)
	for _, level := range Levels {
		fileName := fmt.Sprintf("%s.log", strings.ToLower(level.String()))
		writer, err := NewFileWriter(path.Join(baseDir, fileName))
		if nil != err {
			return nil, err
		}
		fileWriters.writers[level] = writer
	}

	return
}

func (self *FileWriters) SetRotateSize(rotateSize ByteSize) {
	for _, fileWriter := range self.writers {
		fileWriter.SetRotateSize(rotateSize)
	}
}

func (self *FileWriters) SetRotateLines(rotateLines int) {
	for _, fileWriter := range self.writers {
		fileWriter.SetRotateLines(rotateLines)
	}
}

func (self *FileWriters) SetColored(colored bool) {
	for _, fileWriter := range self.writers {
		fileWriter.SetColored(colored)
	}
}

func (self *FileWriters) SetHook(hook Hook) {
	for _, fileWriter := range self.writers {
		fileWriter.SetHook(hook)
	}
}

func (self *FileWriters) SetLevel(level Level) *FileWriters {
	self.level = level
	for _, fileWriter := range self.writers {
		fileWriter.SetLevel(level)
	}
	return self
}

func (self *FileWriters) Level() Level {
	return self.level
}

func (self *FileWriters) Close() {
	for _, fileWriter := range self.writers {
		fileWriter.Close()
	}
	self.closed = true
}

func (self *FileWriters) Debug(format string) {
	if DEBUG < self.level {
		return
	}

	self.writers[DEBUG].write(DEBUG, format)
}

func (self *FileWriters) Debugf(format string, args ...interface{}) {
	if DEBUG < self.level {
		return
	}

	self.writers[DEBUG].writef(DEBUG, format, args...)
}

func (self *FileWriters) Trace(format string) {
	if TRACE < self.level {
		return
	}

	self.writers[TRACE].write(TRACE, format)
}

func (self *FileWriters) Tracef(format string, args ...interface{}) {
	if TRACE < self.level {
		return
	}

	self.writers[TRACE].writef(TRACE, format, args...)
}

func (self *FileWriters) Info(format string) {
	if INFO < self.level {
		return
	}

	self.writers[INFO].write(INFO, format)
}

func (self *FileWriters) Infof(format string, args ...interface{}) {
	if INFO < self.level {
		return
	}

	self.writers[INFO].writef(INFO, format, args...)
}

func (self *FileWriters) Warn(format string) {
	if WARNING < self.level {
		return
	}

	self.writers[WARNING].write(WARNING, format)
}

func (self *FileWriters) Warnf(format string, args ...interface{}) {
	if WARNING < self.level {
		return
	}

	self.writers[WARNING].writef(WARNING, format, args...)
}

func (self *FileWriters) Error(format string) {
	if ERROR < self.level {
		return
	}

	self.writers[ERROR].write(ERROR, format)
}

func (self *FileWriters) Errorf(format string, args ...interface{}) {
	if ERROR < self.level {
		return
	}

	self.writers[ERROR].writef(ERROR, format, args...)
}

func (self *FileWriters) Critical(format string) {
	if CRITICAL < self.level {
		return
	}

	self.writers[CRITICAL].write(CRITICAL, format)
}

func (self *FileWriters) Criticalf(format string, args ...interface{}) {
	if CRITICAL < self.level {
		return
	}

	self.writers[CRITICAL].writef(CRITICAL, format, args...)
}
