// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"testing"
)

func BenchmarkConsoleWriter(b *testing.B) {
	_, err := NewConsoleWriter()
	defer blog.Close()
	if nil != err {
		b.Error(err.Error())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		blog.Debugf("haha %s. en\\en, always %d and %f", "eddie", 18, 3.1415)
	}
}
