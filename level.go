// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"fmt"
)

// type defined for logging level
// just use int
type Level int

const (
	// level enum
	DEBUG Level = iota
	TRACE
	INFO
	WARNING
	ERROR
	CRITICAL
	UNKNOWN = "UNKNOWN"

	DefaultLevel = DEBUG // default level for writers

	// level format ahead every message
	PrefixFormat        = " [%s] "                // pure format
	ColoredPrefixFormat = " [\x1b[%dm%s\x1b[0m] " // colored format

	// color enum used in formating color bytes
	NOCOLOR = 0
	RED     = 31
	GREEN   = 32
	YELLOW  = 33
	BLUE    = 34
	GRAY    = 37
)

var (
	// string present for each level
	LevelStrings = [...]string{"DEBUG", "TRACE", "INFO", "WARN", "ERROR", "CRITICAL"}

	// a slice consist of all levels
	Levels = [...]Level{DEBUG, TRACE, INFO, WARNING, ERROR, CRITICAL}

	// preformatted level prefix string
	// help reduce string formatted burden in realtime logging
	Prefix map[Level]string = make(map[Level]string)
)

func init() {
	initPrefix(false) // preformat level prefix string
}

// initPrefix is designed to preformat level prefix string for each level.
// colored decide whether preformat in colored format or not.
// if colored is true, preformat level prefix string in colored format
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

// valid determines whether a Level instance is valid or not
func (self Level) valid() bool {
	if DEBUG > self || CRITICAL < self {
		return false
	}
	return true
}

// String return string format associate with a Level instance
func (self Level) String() string {
	if !self.valid() {
		return UNKNOWN
	}
	return LevelStrings[self]
}

// Prefix return formatted prefix string associate with a Level instance
func (self Level) Prefix() string {
	return Prefix[self]
}
