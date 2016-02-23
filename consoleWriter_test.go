// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"fmt"
	"testing"
)

func BenchmarkConsoleWriter(b *testing.B) {
	b.StopTimer()
	_, err := NewConsoleWriter()
	defer blog.Close()
	if nil != err {
		fmt.Println(err.Error())
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		blog.Debugf("haha %s. en\\en, always %d and %f", "eddie", 18, 3.1415)
	}
}
