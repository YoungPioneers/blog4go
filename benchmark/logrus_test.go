// Copyright 2015
// Author: huangjunwei@youmi.net

package blog4go

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"sync"
	"testing"
	"time"
)

func BenchmarkLogrusSingleGoroutine(b *testing.B) {
	file, err := os.OpenFile("output_logrus.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	if nil != err {
		fmt.Println(err.Error())
	}
	defer file.Close()
	log.SetOutput(file)

	t := T{123, "test"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Printf("%s [%s] haha %s. en\\en, always %d and %f, %t, %+v\n", time.Now().Format("2006-01-02 15:04:05"), "INFO", "eddie", 18, 3.1415, true, t)
		log.Printf("%s [%s] haha %s. en\\en, always %d and %f, %t, %+v\n", time.Now().Format("2006-01-02 15:04:05"), "ERROR", "eddie", 18, 3.1415, true, t)
	}
}

func BenchmarkLogrusWithTimecacheSingleGoroutine(b *testing.B) {
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

	file, err := os.OpenFile("output_logrus.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	if nil != err {
		fmt.Println(err.Error())
	}
	defer file.Close()
	log.SetOutput(file)

	t := T{123, "test"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Printf("%s [%s] haha %s. en\\en, always %d and %f, %t, %+v\n", timeCache.format, "INFO", "eddie", 18, 3.1415, true, t)
		log.Printf("%s [%s] haha %s. en\\en, always %d and %f, %t, %+v\n", timeCache.format, "ERROR", "eddie", 18, 3.1415, true, t)
	}
}

func BenchmarkLogrusWithTimecacheMultiGoroutine(b *testing.B) {
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

	file, err := os.OpenFile("output_logrus.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	if nil != err {
		fmt.Println(err.Error())
	}
	defer file.Close()
	log.SetOutput(file)

	t := T{123, "test"}

	var wg sync.WaitGroup
	var beginWg sync.WaitGroup

	f := func() {
		defer wg.Done()
		beginWg.Wait()
		for i := 0; i < b.N; i++ {
			log.Printf("%s [%s] haha %s. en\\en, always %d and %f, %t, %+v\n", timeCache.format, "INFO", "eddie", 18, 3.1415, true, t)
			log.Printf("%s [%s] haha %s. en\\en, always %d and %f, %t, %+v\n", timeCache.format, "ERROR", "eddie", 18, 3.1415, true, t)
		}
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		beginWg.Add(1)
	}

	b.ResetTimer()
	for i := 0; i < 100; i++ {
		go f()
		beginWg.Done()
	}

	wg.Wait()
}
