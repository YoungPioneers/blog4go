// Copyright 2015
// Author: huangjunwei@youmi.net

package blog4go

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

func BenchmarkFmtFormatSingleGoroutine(b *testing.B) {
	file, err := os.OpenFile("fmt_output.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	if nil != err {
		fmt.Println(err.Error())
	}
	writer := bufio.NewWriterSize(file, 4096)

	t := T{123, "test"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		writer.WriteString(fmt.Sprintf("%s [%s] haha %s. en\\en, always %d and %f, %t, %+v", time.Now().Format("2006-01-02 15:04:05"), "DEBUG", "eddie", 18, 3.1415, true, t))
		writer.WriteString(fmt.Sprintf("%s [%s] haha %s. en\\en, always %d and %f, %t, %+v", time.Now().Format("2006-01-02 15:04:05"), "ERROR", "eddie", 18, 3.1415, true, t))
	}
	writer.Flush()
}

func BenchmarkFmtFormatWithTimecacheSingleGoroutine(b *testing.B) {
	now := time.Now()

	timeCache := timeFormatCacheType{
		now:    now,
		format: now.Format("2006-01-02 15:04:05"),
	}

	// update timeCache every seconds
	go func() {
		// tick every seconds
		t := time.Tick(1 * time.Second)

		//UpdateTimeCacheLoop:
		for {
			select {
			case <-t:
				now := time.Now()
				timeCache.now = now
				timeCache.format = now.Format("[2006-01-02 15:04:05]")
			}
		}
	}()

	file, err := os.OpenFile("fmt_output_timeCache.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	if nil != err {
		fmt.Println(err.Error())
	}
	writer := bufio.NewWriterSize(file, 4096)

	t := T{123, "test"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		writer.WriteString(fmt.Sprintf("%s %s haha %s. en\\en, always %d and %f, %t, %+v", timeCache.format, "INFO", "eddie", 18, 3.1415, true, t))
		writer.WriteString(fmt.Sprintf("%s %s haha %s. en\\en, always %d and %f, %t, %+v", timeCache.format, "ERROR", "eddie", 18, 3.1415, true, t))
	}
	writer.Flush()
}

func BenchmarkFmtFormatWithTimecacheMultiGoroutine(b *testing.B) {
	now := time.Now()

	timeCache := timeFormatCacheType{
		now:    now,
		format: now.Format("2006-01-02 15:04:05"),
	}

	var lock *sync.Mutex
	lock = new(sync.Mutex)

	// update timeCache every seconds
	go func() {
		// tick every seconds
		t := time.Tick(1 * time.Second)

		//UpdateTimeCacheLoop:
		for {
			select {
			case <-t:
				now := time.Now()
				timeCache.now = now
				timeCache.format = now.Format("[2006-01-02 15:04:05]")
			}
		}
	}()

	file, err := os.OpenFile("fmt_output_timeCache.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	if nil != err {
		fmt.Println(err.Error())
	}
	writer := bufio.NewWriterSize(file, 4096)

	t := T{123, "test"}

	var wg sync.WaitGroup
	var beginWg sync.WaitGroup

	f := func() {
		defer wg.Done()
		beginWg.Wait()
		for i := 0; i < b.N; i++ {
			lock.Lock()
			writer.WriteString(fmt.Sprintf("%s %s haha %s. en\\en, always %d and %f, %t, %+v", timeCache.format, "INFO", "eddie", 18, 3.1415, true, t))
			writer.WriteString(fmt.Sprintf("%s %s haha %s. en\\en, always %d and %f, %t, %+v", timeCache.format, "ERROR", "eddie", 18, 3.1415, true, t))
			lock.Unlock()
		}
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		beginWg.Add(1)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
	}

	for i := 0; i < 100; i++ {
		go f()
		beginWg.Done()
	}

	writer.Flush()
	wg.Wait()
}
