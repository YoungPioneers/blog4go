// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"errors"
	"fmt"
	"path"
	"strings"
)

const (
	// TypeTimeBaseRotate is time base logrotate tag
	TypeTimeBaseRotate = "time"
	// TypeSizeBaseRotate is size base logrotate tag
	TypeSizeBaseRotate = "size"
)

var (
	// ErrFilePathNotFound file path not found
	ErrFilePathNotFound = errors.New("File Path must be defined.")
	// ErrInvalidLevel invalid level string
	ErrInvalidLevel = errors.New("Invalid level string.")
	// ErrInvalidRotateType invalid logrotate type
	ErrInvalidRotateType = errors.New("Invalid log rotate type.")
)

// FileWriter struct defines a writer for multi-files writer with different message level
type FileWriter struct {
	level Level

	// file writers
	writers map[Level]*BaseFileWriter

	closed bool
}

// NewFileWriter initialize a file writer
// baseDir must be base directory of log files
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

	blog = fileWriter
	return
}

// NewFileWriterFromConfigAsFile initialize a file writer according to given config file
// configFile must be the path to the config file
func NewFileWriterFromConfigAsFile(configFile string) (fileWriter *FileWriter, err error) {
	// read config from file
	config, err := readConfig(configFile)
	if nil != err {
		return nil, err
	}

	fileWriter = new(FileWriter)

	fileWriter.level = DEBUG
	if level := LevelFromString(config.MinLevel); level.valid() {

		fileWriter.level = level
	}
	fileWriter.closed = false
	fileWriter.writers = make(map[Level]*BaseFileWriter)

	for _, filter := range config.Filters {
		var rotate = false
		// get file path
		var filePath string
		if nil != &filter.File && "" != filter.File.Path {
			filePath = filter.File.Path
			rotate = false
		} else if nil != &filter.RotateFile && "" != filter.RotateFile.Path {
			filePath = filter.RotateFile.Path
			rotate = true
		} else {
			// config error
			return nil, ErrFilePathNotFound
		}

		// init a base file writer
		writer, err := NewBaseFileWriter(filePath)
		if nil != err {
			return nil, err
		}

		levels := strings.Split(filter.Levels, ",")
		for _, levelStr := range levels {
			var level Level
			if level = LevelFromString(levelStr); !level.valid() {
				return nil, ErrInvalidLevel
			}

			if rotate {
				// set logrotate strategy
				switch filter.RotateFile.Type {
				case TypeTimeBaseRotate:
					writer.SetTimeRotated(true)
				case TypeSizeBaseRotate:
					writer.SetRotateSize(filter.RotateFile.RotateSize)
					writer.SetRotateLines(filter.RotateFile.RotateLines)
				default:
					return nil, ErrInvalidRotateType
				}
			}

			// set color
			fileWriter.SetColored(filter.Colored)
			fileWriter.writers[level] = writer
		}
	}

	blog = fileWriter
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

// SetHookLevel set hook level for every logging actions
func (writer *FileWriter) SetHookLevel(level Level) {
	for _, fileWriter := range writer.writers {
		fileWriter.SetHookLevel(level)
	}
}

// SetLevel set logging level threshold
func (writer *FileWriter) SetLevel(level Level) {
	writer.level = level
	for _, fileWriter := range writer.writers {
		fileWriter.SetLevel(level)
	}
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
	_, ok := writer.writers[DEBUG]
	if !ok || DEBUG < writer.level {
		return
	}

	writer.writers[DEBUG].write(DEBUG, format)
}

// Debugf debugf
func (writer *FileWriter) Debugf(format string, args ...interface{}) {
	_, ok := writer.writers[DEBUG]
	if !ok || DEBUG < writer.level {
		return
	}

	writer.writers[DEBUG].writef(DEBUG, format, args...)
}

// Trace trace
func (writer *FileWriter) Trace(format string) {
	_, ok := writer.writers[TRACE]
	if !ok || TRACE < writer.level {
		return
	}

	writer.writers[TRACE].write(TRACE, format)
}

// Tracef tracef
func (writer *FileWriter) Tracef(format string, args ...interface{}) {
	_, ok := writer.writers[TRACE]
	if !ok || TRACE < writer.level {
		return
	}

	writer.writers[TRACE].writef(TRACE, format, args...)
}

// Info info
func (writer *FileWriter) Info(format string) {
	_, ok := writer.writers[INFO]
	if !ok || INFO < writer.level {
		return
	}

	writer.writers[INFO].write(INFO, format)
}

// Infof infof
func (writer *FileWriter) Infof(format string, args ...interface{}) {
	_, ok := writer.writers[INFO]
	if !ok || INFO < writer.level {
		return
	}

	writer.writers[INFO].writef(INFO, format, args...)
}

// Warn warn
func (writer *FileWriter) Warn(format string) {
	_, ok := writer.writers[WARNING]
	if !ok || WARNING < writer.level {
		return
	}

	writer.writers[WARNING].write(WARNING, format)
}

// Warnf warnf
func (writer *FileWriter) Warnf(format string, args ...interface{}) {
	_, ok := writer.writers[WARNING]
	if !ok || WARNING < writer.level {
		return
	}

	writer.writers[WARNING].writef(WARNING, format, args...)
}

// Error error
func (writer *FileWriter) Error(format string) {
	_, ok := writer.writers[ERROR]
	if !ok || ERROR < writer.level {
		return
	}

	writer.writers[ERROR].write(ERROR, format)
}

// Errorf errorf
func (writer *FileWriter) Errorf(format string, args ...interface{}) {
	_, ok := writer.writers[ERROR]
	if !ok || ERROR < writer.level {
		return
	}

	writer.writers[ERROR].writef(ERROR, format, args...)
}

// Critical critical
func (writer *FileWriter) Critical(format string) {
	_, ok := writer.writers[CRITICAL]
	if !ok || CRITICAL < writer.level {
		return
	}

	writer.writers[CRITICAL].write(CRITICAL, format)
}

// Criticalf criticalf
func (writer *FileWriter) Criticalf(format string, args ...interface{}) {
	_, ok := writer.writers[CRITICAL]
	if !ok || CRITICAL < writer.level {
		return
	}

	writer.writers[CRITICAL].writef(CRITICAL, format, args...)
}
