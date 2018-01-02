// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"fmt"
	"strings"
)

// LevelType type defined for logging level
// just use int
type LevelType int

const (
	// level enum

	// TRACE trace level
	TRACE LevelType = iota
	// DEBUG debug level
	DEBUG
	// INFO info level
	INFO
	// WARNING warn level
	WARNING
	// ERROR error level
	ERROR
	// CRITICAL critical level
	CRITICAL
	// UNKNOWN unknown level
	UNKNOWN = "UNKNOWN"

	// DefaultLevel default level for writers
	DefaultLevel = TRACE

	// PrefixFormat is the level format ahead every message
	PrefixFormat = " level=\"%s\" " // pure format
	// ColoredPrefixFormat is the colored level format adhead every message
	ColoredPrefixFormat = " level=\"\x1b[%dm%s\x1b[0m\" " // colored format

	// color enum used in formating color bytes

	// NOCOLOR no color
	NOCOLOR = 0
	// RED red color
	RED = 31
	// GREEN green color
	GREEN = 32
	// YELLOW yellow color
	YELLOW = 33
	// BLUE blue color
	BLUE = 34
	// GRAY gray color
	GRAY = 37
)

var (
	// LevelStrings is string present for each level
	LevelStrings = [...]string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "CRITICAL"}

	// StringLevels is map, level strings to levels
	StringLevels = map[string]LevelType{"TRACE": TRACE, "DEBUG": DEBUG, "INFO": INFO, "WARN": WARNING, "ERROR": ERROR, "CRITICAL": CRITICAL}

	// Levels is a slice consist of all levels
	Levels = [...]LevelType{TRACE, DEBUG, INFO, WARNING, ERROR, CRITICAL}

	// Prefix is preformatted level prefix string
	// help reduce string formatted burden in realtime logging
	Prefix = make(map[LevelType]string)
)

func init() {
	initPrefix(false) // preformat level prefix string
}

// initPrefix is designed to preformat level prefix string for each level.
// colored decide whether preformat in colored format or not.
// if colored is true, preformat level prefix string in colored format
func initPrefix(colored bool) {
	if colored {
		Prefix[TRACE] = fmt.Sprintf(ColoredPrefixFormat, GRAY, TRACE.String())
		Prefix[DEBUG] = fmt.Sprintf(ColoredPrefixFormat, GREEN, DEBUG.String())
		Prefix[INFO] = fmt.Sprintf(ColoredPrefixFormat, BLUE, INFO.String())
		Prefix[WARNING] = fmt.Sprintf(ColoredPrefixFormat, YELLOW, WARNING.String())
		Prefix[ERROR] = fmt.Sprintf(ColoredPrefixFormat, RED, ERROR.String())
		Prefix[CRITICAL] = fmt.Sprintf(ColoredPrefixFormat, RED, CRITICAL.String())
	} else {
		Prefix[TRACE] = fmt.Sprintf(PrefixFormat, TRACE.String())
		Prefix[DEBUG] = fmt.Sprintf(PrefixFormat, DEBUG.String())
		Prefix[INFO] = fmt.Sprintf(PrefixFormat, INFO.String())
		Prefix[WARNING] = fmt.Sprintf(PrefixFormat, WARNING.String())
		Prefix[ERROR] = fmt.Sprintf(PrefixFormat, ERROR.String())
		Prefix[CRITICAL] = fmt.Sprintf(PrefixFormat, CRITICAL.String())
	}
}

// valid determines whether a Level instance is valid or not
func (level LevelType) valid() bool {
	if TRACE > level || CRITICAL < level {
		return false
	}
	return true
}

// String return string format associate with a Level instance
func (level LevelType) String() string {
	if !level.valid() {
		return UNKNOWN
	}
	return LevelStrings[level]
}

// prefix return formatted prefix string associate with a Level instance
func (level LevelType) prefix() string {
	return Prefix[level]
}

// LevelFromString return Level according to given string
func LevelFromString(str string) LevelType {
	level, ok := StringLevels[strings.ToUpper(str)]
	if !ok {
		return LevelType(-1)
	}
	return level
}
