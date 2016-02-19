// Copyright 2015
// Author: huangjunwei@youmi.net

package blog4go

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"testing"
	"time"
)

func BenchmarkLogrus(b *testing.B) {
	b.StopTimer()
	file, err := os.OpenFile("output_logrus.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	if nil != err {
		fmt.Println(err.Error())
	}
	defer file.Close()
	log.SetOutput(file)

	t := T{123, "test"}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		log.Printf("%s [%s] haha %s. en\\en, always %d and %f, %t, %+v\n", time.Now().Format("2006-01-02 15:04:05"), "INFO", "eddie", 18, 3.1415, true, t)
		log.Printf("%s [%s] haha %s. en\\en, always %d and %f, %t, %+v\n", time.Now().Format("2006-01-02 15:04:05"), "ERROR", "eddie", 18, 3.1415, true, t)
	}
}
func BenchmarkLogrusWithTimecache(b *testing.B) {
	b.StopTimer()
	now := time.Now()

	timeCache := timeFormatCacheType{
		now:    now,
		format: now.Format("2006-01-02 15:04:05"),
	}

	file, err := os.OpenFile("output_logrus.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	if nil != err {
		fmt.Println(err.Error())
	}
	defer file.Close()
	log.SetOutput(file)

	t := T{123, "test"}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		now := time.Now()
		if now != timeCache.now {
			timeCache.now = now
			timeCache.format = now.Format("[2006/01/02:15:04:05]")
		}
		log.Printf("%s [%s] haha %s. en\\en, always %d and %f, %t, %+v\n", time.Now().Format("2006-01-02 15:04:05"), "INFO", "eddie", 18, 3.1415, true, t)
		log.Printf("%s [%s] haha %s. en\\en, always %d and %f, %t, %+v\n", time.Now().Format("2006-01-02 15:04:05"), "ERROR", "eddie", 18, 3.1415, true, t)
	}
}
