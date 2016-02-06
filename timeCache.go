// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"time"
)

// 时间格式化的cache
type timeFormatCacheType struct {
	// 当前
	now time.Time
	// 当前日期
	date string
	// 当前时间格式化结果
	// bufio write bytes会比write string效率高
	format []byte
	// 昨日日期
	dateYesterday string
}

// 用全局的timeCache好像比较好
// 实例的timeCache没那么好统一更新
var timeCache = timeFormatCacheType{}

// 包初始化函数
func init() {
	timeCache.now = time.Now()
	timeCache.date = timeCache.now.Format(DateFormat)
	timeCache.format = []byte(timeCache.now.Format(PrefixTimeFormat))
	timeCache.dateYesterday = timeCache.now.Add(-24 * time.Hour).Format(DateFormat)
}
