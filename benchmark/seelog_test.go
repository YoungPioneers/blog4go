// Copyright 2015
// Author: huangjunwei@youmi.net

package blog4go

import (
	"fmt"
	log "github.com/cihub/seelog"
	"testing"
)

func BenchmarkFormat(b *testing.B) {
	b.StopTimer()
	logger, err := log.LoggerFromConfigAsFile("log_config.xml")
	if nil != err {
		fmt.Println(err.Error())
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		logger.Debugf("haha %s. en\\en, always %d and %.4f", "eddie", 18, 3.1415)
	}
}
