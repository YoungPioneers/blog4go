// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"errors"
	"os"
)

type ByteSize int

const (
	// 大小单位
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB

	// 默认logrotate条件
	DefaultRotateSize  = 500 * MB
	DefaultRotateLines = 2000000 // 200w

	// 时间前缀的格式
	PrefixTimeFormat = "[2006/01/02:15:04:05]"
	// 日期格式
	DateFormat = "2006-01-02"

	// 换行符
	EOL = '\n'
	// 转移符
	ESCAPE = '\\'
	// 占位符
	PLACEHOLDER = '%'
)

var (
	// bufio buffer size
	// 好像buffer size 调大点benchmark效果更好
	// 默认使用内存页大小
	DefaultBufferSize = 4096

	ErrInvalidFormat = errors.New("Invalid format type.")
)

// 包初始化函数
func init() {
	DefaultBufferSize = os.Getpagesize()
}
