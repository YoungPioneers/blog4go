// Copyright 2015
// Author: huangjunwei@youmi.net

package blog4go

import (
	"testing"
)

func BenchmarkFormat(b *testing.B) {
	b.StopTimer()
	writer, _ := NewFileLogWriter("output.log", false)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		writer.Debugf("haha %s. en\\en, always %d and %.4f", "eddie", 18, 3.1415)
	}
	writer.Close()
}
