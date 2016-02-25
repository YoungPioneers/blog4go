// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"fmt"
	log "github.com/YoungPioneers/blog4go"
	"sync"
	"testing"
	"time"
)

type T struct {
	A int
	B string
}

type timeFormatCacheType struct {
	now    time.Time
	format string
}

func BenchmarkBlog4goSingleGoroutine(b *testing.B) {
	err := log.NewWriterFromConfigAsFile("blog4go_config.xml")
	defer log.Close()
	if nil != err {
		fmt.Println(err.Error())
	}

	t := T{123, "test"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Infof("haha %s. en\\en, always %d and %f, %t, %+v", "eddie", 18, 3.1415, true, t)
		log.Errorf("haha %s. en\\en, always %d and %f, %t, %+v", "eddie", 18, 3.1415, true, t)
	}
}

func BenchmarkBlog4goMultiGoroutine(b *testing.B) {
	err := log.NewWriterFromConfigAsFile("blog4go_config.xml")
	defer log.Close()
	if nil != err {
		fmt.Println(err.Error())
	}

	t := T{123, "test"}

	var wg sync.WaitGroup
	var beginWg sync.WaitGroup

	f := func() {
		defer wg.Done()
		beginWg.Wait()
		for i := 0; i < b.N; i++ {
			log.Infof("haha %s. en\\en, always %d and %f, %t, %+v", "eddie", 18, 3.1415, true, t)
			log.Errorf("haha %s. en\\en, always %d and %f, %t, %+v", "eddie", 18, 3.1415, true, t)
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
