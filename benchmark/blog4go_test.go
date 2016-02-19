// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"fmt"
	"github.com/YoungPioneers/blog4go"
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

func BenchmarkBlog4goFormat(b *testing.B) {
	b.StopTimer()
	err := blog4go.NewFileWriterFromConfigAsFile("blog4go_config.xml")
	defer blog4go.Close()
	if nil != err {
		fmt.Println(err.Error())
	}

	t := T{123, "test"}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		blog4go.Infof("haha %s. en\\en, always %d and %f, %t, %+v", "eddie", 18, 3.1415, true, t)
		blog4go.Errorf("haha %s. en\\en, always %d and %f, %t, %+v", "eddie", 18, 3.1415, true, t)
	}
}
