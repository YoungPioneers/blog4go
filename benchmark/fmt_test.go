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

type timeFormatCacheType struct {
	now    time.Time
	format string
}

func BenchmarkFmtFormat(b *testing.B) {
	b.StopTimer()
	file, err := os.OpenFile("output.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	if nil != err {
		fmt.Println(err.Error())
	}
	writer := bufio.NewWriterSize(file, 4096)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		writer.WriteString(fmt.Sprintf("%s [%s]haha %s. en\\en, always %d and %.4f", time.Now().Format("2006-01-02 15:04:05"), "DEBUG", "eddie", 18, 3.1415))
	}
}

func BenchmarkFmtWithTimecacheFormat(b *testing.B) {
	b.StopTimer()
	now := time.Now()

	timeCache := timeFormatCacheType{
		now:    now,
		format: now.Format("2006-01-02 15:04:05"),
	}

	file, err := os.OpenFile("output.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	if nil != err {
		fmt.Println(err.Error())
	}
	writer := bufio.NewWriterSize(file, 4096)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		now := time.Now()
		if now != timeCache.now {
			timeCache.now = now
			timeCache.format = now.Format("[2006/01/02:15:04:05]")
		}
		writer.WriteString(fmt.Sprintf("%s [%s]haha %s. en\\en, always %d and %.4f", timeCache.format, "DEBUG", "eddie", 18, 3.1415))
	}
}
