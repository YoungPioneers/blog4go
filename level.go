// Copyright 2015
// Author: huangjunwei@youmi.net

package blog4go

import (
	"fmt"
)

type Level int

const (
	DEBUG Level = iota
	TRACE
	INFO
	WARNING
	ERROR
	CRITICAL

	PrefixFormat        = " [%s] "
	ColoredPrefixFormat = " [\x1b[%dm%s\x1b[0m] "
	UNKNOWN             = "UNKNOWN"

	NOCOLOR = 0
	RED     = 31
	GREEN   = 32
	YELLOW  = 33
	BLUE    = 34
	GRAY    = 37
)

var (
	LevelStrings = [...]string{"DEBUG", "TRACE", "INFO", "WARN", "ERROR", "CRITAL"}

	// 定义一些日志格式的前缀，减少字符串拼接操作
	Prefix map[Level]string = make(map[Level]string)
)

func init() {
	initPrefix(false)
}

func initPrefix(colored bool) {
	if colored {
		Prefix[DEBUG] = fmt.Sprintf(ColoredPrefixFormat, GRAY, DEBUG.String())
		Prefix[TRACE] = fmt.Sprintf(ColoredPrefixFormat, GREEN, TRACE.String())
		Prefix[INFO] = fmt.Sprintf(ColoredPrefixFormat, BLUE, INFO.String())
		Prefix[WARNING] = fmt.Sprintf(ColoredPrefixFormat, YELLOW, WARNING.String())
		Prefix[ERROR] = fmt.Sprintf(ColoredPrefixFormat, RED, ERROR.String())
		Prefix[CRITICAL] = fmt.Sprintf(ColoredPrefixFormat, RED, CRITICAL.String())

	} else {
		Prefix[DEBUG] = fmt.Sprintf(PrefixFormat, DEBUG.String())
		Prefix[TRACE] = fmt.Sprintf(PrefixFormat, TRACE.String())
		Prefix[INFO] = fmt.Sprintf(PrefixFormat, INFO.String())
		Prefix[WARNING] = fmt.Sprintf(PrefixFormat, WARNING.String())
		Prefix[ERROR] = fmt.Sprintf(PrefixFormat, ERROR.String())
		Prefix[CRITICAL] = fmt.Sprintf(PrefixFormat, CRITICAL.String())

	}

}

// 有效性判断好像必要性不大
func (self Level) valid() bool {
	if DEBUG > self || CRITICAL < self {
		return false
	}
	return true
}

func (self Level) String() string {
	if !self.valid() {
		return UNKNOWN
	}
	return LevelStrings[self]
}

func (self Level) Prefix() string {
	return Prefix[self]
}
