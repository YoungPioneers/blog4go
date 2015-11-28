// Copyright 2015
// Author: huangjunwei@youmi.net

package blog4go

type Level int

const (
	DEBUG Level = iota
	TRACE
	INFO
	WARNING
	ERROR
	CRITICAL
)

var (
	levelStrings = [...]string{"DEBUG", "TRAC", "INFO", "WARN", "ERROR", "CRITAL"}
)

func (self Level) valid() bool {
	if DEBUG > self || CRITICAL < self {
		return false
	}
	return true
}

func (self Level) String() string {
	if !self.valid() {
		return "UNKNOWN"
	}
	return levelStrings[self]
}
