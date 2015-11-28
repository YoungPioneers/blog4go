// Copyright 2015
// Author: huangjunwei@youmi.net

package blog4go

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func BenchmarkFmtFormat(b *testing.B) {
	b.StopTimer()
	file, err := os.OpenFile("output.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	if nil != err {
		fmt.Println(err.Error())
	}
	writer := bufio.NewWriterSize(file, 4096)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		writer.WriteString(fmt.Sprintf("haha %s. en\\en, always %d and %.4f", "eddie", 18, 3.1415))
	}
}
