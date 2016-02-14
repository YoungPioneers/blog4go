package main

import (
	"fmt"
	"github.com/YoungPioneers/blog4go"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

type MyHook struct {
	something string
}

func (self *MyHook) Fire(level blog4go.Level, message string) {
	if level > blog4go.ERROR {
		fmt.Println(message)
	}
}

func main() {
	runtime.GOMAXPROCS(4)

	writers, err := blog4go.NewFileWriterFromConfigAsFile("config.xml")
	if nil != err {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer writers.Close()
	writers.Debug("Debug")
	writers.Trace("Trace")
	writers.Info("Info")
	writers.Warn("Warn")
	writers.Error("Error")
	writers.Critical("Critical")

	// blog
	//writer, err := blog4go.NewBaseFileWriter("output.log")
	//if nil != err {
	//fmt.Println(err.Error())
	//os.Exit(1)
	//}
	//defer writer.Close()

	// test rotate line
	//writer.SetRotateLines(100)

	// test hook
	//hook := new(MyHook)
	//writer.SetHook(hook)

	//for i := 1; i < 5; i++ {
	//go logging(writer)
	//}

	// blog writers
	//writers, err := blog4go.NewFileWriter("./")
	//if nil != err {
	//fmt.Println(err.Error())
	//os.Exit(1)
	//}
	//defer writers.Close()
	//writers.Debug("Debug")
	//writers.Trace("Trace")
	//writers.Info("Info")
	//writers.Warn("Warn")
	//writers.Error("Error")
	//writers.Critical("Critical")

	// socket writer
	// nc -u -l 12124
	//socketWriter, err := blog4go.NewSocketWriter("udp", "127.0.0.1:12124")
	//defer socketWriter.Close()
	//socketWriter.Debug("debug")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		select {
		case <-c:
			fmt.Println("Exit..")
			time.Sleep(1 * time.Second)
			return
		}
	}
}

type T struct {
	A int
	B string
}

// blog
func logging(writer *blog4go.BaseFileWriter) {
	t := T{123, "test"}
	d := int64(18)
	for {
		writer.Debug("test_debug")
		writer.Trace("test_trace")
		writer.Info("test_info")
		writer.Warn("test_warn")
		writer.Error("test_error")
		writer.Critical("test_critical")
		writer.Debugf("haha %s. en\\en, always %d and %5.4f, %t, %+v", "eddie", d, 3.14159, true, t)
		time.Sleep(2 * time.Second)
	}
}
