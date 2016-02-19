// Copyright 2015
// Author: huangjunwei@youmi.net

package blog4go

import (
	"fmt"
	log "github.com/cihub/seelog"
	"testing"
)

func BenchmarkSeelogFormat(b *testing.B) {
	b.StopTimer()
	logger, err := log.LoggerFromConfigAsFile("seelog_config.xml")
	if nil != err {
		fmt.Println(err.Error())
	}

	t := T{123, "test"}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		logger.Infof("haha %s. en\\en, always %d and %f, %t, %+v", "eddie", 18, 3.1415, true, t)
		logger.Errorf("haha %s. en\\en, always %d and %f, %t, %+v", "eddie", 18, 3.1415, true, t)
	}
}
