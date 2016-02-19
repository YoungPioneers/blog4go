// Copyright 2015
// Author: huangjunwei@youmi.net

package blog4go

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"time"
)

func BenchmarkFmtFormat(b *testing.B) {
	b.StopTimer()
	file, err := os.OpenFile("fmt_output.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	if nil != err {
		fmt.Println(err.Error())
	}
	writer := bufio.NewWriterSize(file, 4096)

	t := T{123, "test"}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		writer.WriteString(fmt.Sprintf("%s [%s]haha %s. en\\en, always %d and %f, %t, %+v", time.Now().Format("2006-01-02 15:04:05"), "DEBUG", "eddie", 18, 3.1415, true, t))
	}
	writer.Flush()
}

func BenchmarkFmtWithTimecacheFormat(b *testing.B) {
	b.StopTimer()
	now := time.Now()

	timeCache := timeFormatCacheType{
		now:    now,
		format: now.Format("2006-01-02 15:04:05"),
	}

	file, err := os.OpenFile("fmt_timeCache_output.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	if nil != err {
		fmt.Println(err.Error())
	}
	writer := bufio.NewWriterSize(file, 4096)

	t := T{123, "test"}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		now := time.Now()
		if now != timeCache.now {
			timeCache.now = now
			timeCache.format = now.Format("[2006/01/02:15:04:05]")
		}
		writer.WriteString(fmt.Sprintf("%s [%s]haha %s. en\\en, always %d and %f, %t, %+v", timeCache.format, "INFO", "eddie", 18, 3.1415, true, t))
		writer.WriteString(fmt.Sprintf("%s [%s]haha %s. en\\en, always %d and %f, %t, %+v", timeCache.format, "ERROR", "eddie", 18, 3.1415, true, t))
	}
	writer.Flush()
}
