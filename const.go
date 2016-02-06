// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"errors"
	"os"
)

// ByteSize is type of sizes
type ByteSize int

const (
	// unit of sizes

	_ = iota // ignore first value by assigning to blank identifier
	// KB unit of kilobyte
	KB ByteSize = 1 << (10 * iota)
	// MB unit of megabyte
	MB
	// GB unit of gigabyte
	GB

	// default logrotate condition

	// DefaultRotateSize is default size when size base logrotate needed
	DefaultRotateSize = 500 * MB
	// DefaultRotateLines is default lines when lines base logrotate needed
	DefaultRotateLines = 2000000 // 2 million

	// PrefixTimeFormat const time format prefix
	PrefixTimeFormat = "[2006/01/02:15:04:05]"
	// DateFormat date format
	DateFormat = "2006-01-02"

	// EOL end of a line
	EOL = '\n'
	// ESCAPE escape character
	ESCAPE = '\\'
	// PLACEHOLDER placeholder
	PLACEHOLDER = '%'
)

var (
	// DefaultBufferSize bufio buffer size
	DefaultBufferSize = 4096 // default memory page size

	// ErrInvalidFormat invalid format error
	ErrInvalidFormat = errors.New("Invalid format type.")
)

func init() {
	DefaultBufferSize = os.Getpagesize()
}
