// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"fmt"
	"path"
	"strings"
)

// FileWriter struct defines a writer for multi-files writer with different
// message level
type FileWriter struct {
	level Level

	// file writers
	writers map[Level]*BaseFileWriter

	closed bool
}

// NewFileWriter initialize a file writer
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

// SetTimeRotated toggle time base logrotate
func (writer *FileWriter) SetTimeRotated(timeRotated bool) {
	for _, fileWriter := range writer.writers {
		fileWriter.SetTimeRotated(timeRotated)
	}
}

// SetRotateSize set size when logroatate
func (writer *FileWriter) SetRotateSize(rotateSize ByteSize) {
	for _, fileWriter := range writer.writers {
		fileWriter.SetRotateSize(rotateSize)
	}
}

// SetRotateLines set line number when logrotate
func (writer *FileWriter) SetRotateLines(rotateLines int) {
	for _, fileWriter := range writer.writers {
		fileWriter.SetRotateLines(rotateLines)
	}
}

// SetColored set logging color
func (writer *FileWriter) SetColored(colored bool) {
	for _, fileWriter := range writer.writers {
		fileWriter.SetColored(colored)
	}
}

// SetHook set hook for every logging actions
func (writer *FileWriter) SetHook(hook Hook) {
	for _, fileWriter := range writer.writers {
		fileWriter.SetHook(hook)
	}
}

// SetLevel set logging level threshold
func (writer *FileWriter) SetLevel(level Level) *FileWriter {
	writer.level = level
	for _, fileWriter := range writer.writers {
		fileWriter.SetLevel(level)
	}
	return writer
}

// Level return logging level threshold
func (writer *FileWriter) Level() Level {
	return writer.level
}

// Close close file writer
func (writer *FileWriter) Close() {
	for _, fileWriter := range writer.writers {
		fileWriter.Close()
	}
	writer.closed = true
}

// Debug debug
func (writer *FileWriter) Debug(format string) {
	if DEBUG < writer.level {
		return
	}

	writer.writers[DEBUG].write(DEBUG, format)
}

// Debugf debugf
func (writer *FileWriter) Debugf(format string, args ...interface{}) {
	if DEBUG < writer.level {
		return
	}

	writer.writers[DEBUG].writef(DEBUG, format, args...)
}

// Trace trace
func (writer *FileWriter) Trace(format string) {
	if TRACE < writer.level {
		return
	}

	writer.writers[TRACE].write(TRACE, format)
}

// Tracef tracef
func (writer *FileWriter) Tracef(format string, args ...interface{}) {
	if TRACE < writer.level {
		return
	}

	writer.writers[TRACE].writef(TRACE, format, args...)
}

// Info info
func (writer *FileWriter) Info(format string) {
	if INFO < writer.level {
		return
	}

	writer.writers[INFO].write(INFO, format)
}

// Infof infof
func (writer *FileWriter) Infof(format string, args ...interface{}) {
	if INFO < writer.level {
		return
	}

	writer.writers[INFO].writef(INFO, format, args...)
}

// Warn warn
func (writer *FileWriter) Warn(format string) {
	if WARNING < writer.level {
		return
	}

	writer.writers[WARNING].write(WARNING, format)
}

// Warnf warnf
func (writer *FileWriter) Warnf(format string, args ...interface{}) {
	if WARNING < writer.level {
		return
	}

	writer.writers[WARNING].writef(WARNING, format, args...)
}

// Error error
func (writer *FileWriter) Error(format string) {
	if ERROR < writer.level {
		return
	}

	writer.writers[ERROR].write(ERROR, format)
}

// Errorf errorf
func (writer *FileWriter) Errorf(format string, args ...interface{}) {
	if ERROR < writer.level {
		return
	}

	writer.writers[ERROR].writef(ERROR, format, args...)
}

// Critical critical
func (writer *FileWriter) Critical(format string) {
	if CRITICAL < writer.level {
		return
	}

	writer.writers[CRITICAL].write(CRITICAL, format)
}

// Criticalf criticalf
func (writer *FileWriter) Criticalf(format string, args ...interface{}) {
	if CRITICAL < writer.level {
		return
	}

	writer.writers[CRITICAL].writef(CRITICAL, format, args...)
}
