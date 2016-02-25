// Copyright 2015
// Author: huangjunwei@youmi.net

package blog4go

import (
	"fmt"
	log "github.com/cihub/seelog"
	"sync"
	"testing"
)

func BenchmarkSeelogSingleGoroutine(b *testing.B) {
	logger, err := log.LoggerFromConfigAsFile("seelog_config.xml")
	if nil != err {
		fmt.Println(err.Error())
	}

	t := T{123, "test"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Infof("haha %s. en\\en, always %d and %f, %t, %+v", "eddie", 18, 3.1415, true, t)
		logger.Errorf("haha %s. en\\en, always %d and %f, %t, %+v", "eddie", 18, 3.1415, true, t)
	}
}

func BenchmarkSeelogMultiGoroutine(b *testing.B) {
	logger, err := log.LoggerFromConfigAsFile("seelog_config.xml")
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
			logger.Infof("haha %s. en\\en, always %d and %f, %t, %+v", "eddie", 18, 3.1415, true, t)
			logger.Errorf("haha %s. en\\en, always %d and %f, %t, %+v", "eddie", 18, 3.1415, true, t)
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
