// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

// ConsoleWriter is a console logger
type ConsoleWriter struct {
	level Level

	writer *bufio.Writer

	lock *sync.Mutex

	closed bool

	colored bool

	// log hook
	hook      Hook
	hookLevel Level
}

// NewConsoleWriter initialize a console writer
func NewConsoleWriter() (consoleWriter *ConsoleWriter, err error) {
	consoleWriter = new(ConsoleWriter)
	consoleWriter.level = DEBUG

	consoleWriter.lock = new(sync.Mutex)
	consoleWriter.closed = false

	consoleWriter.colored = false

	// log hook
	consoleWriter.hook = nil
	consoleWriter.hookLevel = DEBUG

	consoleWriter.writer = bufio.NewWriterSize(os.Stdout, DefaultBufferSize)

	go consoleWriter.daemon()

	return consoleWriter, nil
}

func (writer *ConsoleWriter) daemon() {
	t := time.Tick(1 * time.Second)
	f := time.Tick(10 * time.Second)

DaemonLoop:
	for {
		select {
		case <-f:
			if writer.closed {
				break DaemonLoop
			}

			writer.flush()
		case <-t:
			if writer.closed {
				break DaemonLoop
			}

			now := time.Now()
			timeCache.now = now
			timeCache.format = []byte(now.Format(PrefixTimeFormat))

		}
	}
}

func (writer *ConsoleWriter) write(level Level, format string) {
	writer.lock.Lock()

	defer func() {
		writer.lock.Unlock()
		if nil != writer.hook && !(level < writer.hookLevel) {
			go func(level Level, format string) {
				writer.hook.Fire(level, format)
			}(level, format)
		}
	}()

	if writer.closed {
		return
	}

	writer.writer.Write(timeCache.format)
	writer.writer.WriteString(level.Prefix())
	writer.writer.WriteString(format)
	writer.writer.WriteByte(EOL)
}

func (writer *ConsoleWriter) writef(level Level, format string, args ...interface{}) {
	writer.lock.Lock()

	defer func() {
		writer.lock.Unlock()

		if nil != writer.hook && !(level < writer.hookLevel) {
			go func(level Level, format string, args ...interface{}) {
				writer.hook.Fire(level, fmt.Sprintf(format, args...))
			}(level, format, args...)
		}
	}()

	if writer.closed {
		return
	}

	var tag = false
	var tagPos int
	var escape = false
	var n int
	var last int

	writer.writer.Write(timeCache.format)
	writer.writer.WriteString(level.Prefix())

	for i, v := range format {
		if tag {
			switch v {
			case 'd', 'f', 'v', 'b', 'o', 'x', 'X', 'c', 'p', 't', 's', 'T', 'q', 'U', 'e', 'E', 'g', 'G':
				if escape {
					escape = false
				}

				writer.writer.WriteString(fmt.Sprintf(format[tagPos:i+1], args[n]))
				n++
				last = i + 1
				tag = false
			case ESCAPE:
				if escape {
					writer.writer.WriteByte(ESCAPE)
				}
				escape = !escape
			default:

			}

		} else {
			if PLACEHOLDER == format[i] && !escape {
				tag = true
				tagPos = i
				writer.writer.WriteString(format[last:i])
				escape = false
			}
		}
	}
	writer.writer.WriteString(format[last:])
	writer.writer.WriteByte(EOL)
}

// Level return logging level threshold
func (writer *ConsoleWriter) Level() Level {
	return writer.level
}

// SetLevel set logger level
func (writer *ConsoleWriter) SetLevel(level Level) *ConsoleWriter {
	writer.level = level
	return writer
}

// Colored return whether writer log with color
func (writer *ConsoleWriter) Colored() bool {
	return writer.colored
}

// SetColored set logging color
func (writer *ConsoleWriter) SetColored(colored bool) {
	if colored == writer.colored {
		return
	}

	writer.colored = colored
	writer.lock.Lock()
	defer writer.lock.Unlock()

	initPrefix(colored)
}

// SetHook set hook for logging action
func (writer *ConsoleWriter) SetHook(hook Hook) {
	writer.hook = hook
}

// SetHookLevel set when hook will be called
func (writer *ConsoleWriter) SetHookLevel(level Level) {
	writer.hookLevel = level
}

// Close close console writer
func (writer *ConsoleWriter) Close() {
	writer.lock.Lock()
	defer writer.lock.Unlock()
	if writer.closed {
		return
	}

	writer.writer.Flush()
	writer.writer = nil
	writer.closed = true
}

// flush buffer to disk
func (writer *ConsoleWriter) flush() {
	writer.lock.Lock()
	defer writer.lock.Unlock()
	writer.writer.Flush()
}

// Debug debug
func (writer *ConsoleWriter) Debug(format string) {
	if DEBUG < writer.level {
		return
	}

	writer.write(DEBUG, format)
}

// Debugf debugf
func (writer *ConsoleWriter) Debugf(format string, args ...interface{}) {
	if DEBUG < writer.level {
		return
	}

	writer.writef(DEBUG, format, args...)
}

// Trace trace
func (writer *ConsoleWriter) Trace(format string) {
	if TRACE < writer.level {
		return
	}

	writer.write(TRACE, format)
}

// Tracef tracef
func (writer *ConsoleWriter) Tracef(format string, args ...interface{}) {
	if TRACE < writer.level {
		return
	}

	writer.writef(TRACE, format, args...)
}

// Info info
func (writer *ConsoleWriter) Info(format string) {
	if INFO < writer.level {
		return
	}

	writer.write(INFO, format)
}

// Infof infof
func (writer *ConsoleWriter) Infof(format string, args ...interface{}) {
	if INFO < writer.level {
		return
	}

	writer.writef(INFO, format, args...)
}

// Error error
func (writer *ConsoleWriter) Error(format string) {
	if ERROR < writer.level {
		return
	}

	writer.write(ERROR, format)
}

// Errorf errorf
func (writer *ConsoleWriter) Errorf(format string, args ...interface{}) {
	if ERROR < writer.level {
		return
	}

	writer.writef(ERROR, format, args...)
}

// Warn warn
func (writer *ConsoleWriter) Warn(format string) {
	if WARNING < writer.level {
		return
	}

	writer.write(WARNING, format)
}

// Warnf warnf
func (writer *ConsoleWriter) Warnf(format string, args ...interface{}) {
	if WARNING < writer.level {
		return
	}

	writer.writef(WARNING, format, args...)
}

// Critical critical
func (writer *ConsoleWriter) Critical(format string) {
	if CRITICAL < writer.level {
		return
	}

	writer.write(CRITICAL, format)
}

// Criticalf criticalf
func (writer *ConsoleWriter) Criticalf(format string, args ...interface{}) {
	if CRITICAL < writer.level {
		return
	}

	writer.writef(CRITICAL, format, args...)
}
