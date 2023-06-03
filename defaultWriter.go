// Copyright (c) 2023, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

// DefaultWriter default empty logger
type DefaultWriter struct{}

// Close .
func (writer *DefaultWriter) Close() {}

// SetLevel set logging level threshold
func (writer *DefaultWriter) SetLevel(level LevelType) {}

// Level get log level
func (writer *DefaultWriter) Level() LevelType {
	return TRACE
}

// write/writef functions with different levels
func (writer *DefaultWriter) write(level LevelType, args ...interface{})                 {}
func (writer *DefaultWriter) writef(level LevelType, format string, args ...interface{}) {}

// Debug .
func (writer *DefaultWriter) Debug(args ...interface{}) {}

// Debugf .
func (writer *DefaultWriter) Debugf(format string, args ...interface{}) {}

// Trace .
func (writer *DefaultWriter) Trace(args ...interface{}) {}

// Tracef .
func (writer *DefaultWriter) Tracef(format string, args ...interface{}) {}

// Info .
func (writer *DefaultWriter) Info(args ...interface{}) {}

// Infof .
func (writer *DefaultWriter) Infof(format string, args ...interface{}) {}

// Warn .
func (writer *DefaultWriter) Warn(args ...interface{}) {}

// Warnf .
func (writer *DefaultWriter) Warnf(format string, args ...interface{}) {}

// Error .
func (writer *DefaultWriter) Error(args ...interface{}) {}

// Errorf .
func (writer *DefaultWriter) Errorf(format string, args ...interface{}) {}

// Critical .
func (writer *DefaultWriter) Critical(args ...interface{}) {}

// Criticalf .
func (writer *DefaultWriter) Criticalf(format string, args ...interface{}) {}

// flush log to disk
func (writer *DefaultWriter) flush() {}

// SetHook .
func (writer *DefaultWriter) SetHook(hook Hook) {}

// SetHookLevel .
func (writer *DefaultWriter) SetHookLevel(level LevelType) {}

// SetHookAsync .
func (writer *DefaultWriter) SetHookAsync(async bool) {}

// SetTimeRotated .
func (writer *DefaultWriter) SetTimeRotated(timeRotated bool) {}

// TimeRotated .
func (writer *DefaultWriter) TimeRotated() bool {
	return false
}
func (writer *DefaultWriter) SetRotateSize(rotateSize int64) {}
func (writer *DefaultWriter) RotateSize() int64 {
	return 0
}

// SetRotateLines .
func (writer *DefaultWriter) SetRotateLines(rotateLines int) {}

// RotateLines .
func (writer *DefaultWriter) RotateLines() int {
	return 0
}

// SetRetentions .
func (writer *DefaultWriter) SetRetentions(retentions int64) {}

// Retentions .
func (writer *DefaultWriter) Retentions() int64 {
	return 0
}

// SetColored .
func (writer *DefaultWriter) SetColored(colored bool) {}

// Colored .
func (writer *DefaultWriter) Colored() bool {
	return false
}

// SetTags .
func (writer *DefaultWriter) SetTags(tags map[string]string) {}

// Tags .
func (writer *DefaultWriter) Tags() map[string]string {
	return map[string]string{}
}
