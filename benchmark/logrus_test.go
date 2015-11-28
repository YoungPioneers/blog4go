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
	file, err := os.OpenFile("output.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	if nil != err {
		fmt.Println(err.Error())
	}
	defer file.Close()
	log.SetOutput(file)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		log.Printf("%s [%s] haha %s. en\\en, always %d and %.4f\n", time.Now().Format("2006-01-02 15:04:05"), "DEBUG", "eddie", 18, 3.1415)
	}
}
