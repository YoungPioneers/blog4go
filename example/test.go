package main

import (
	"blog4go"
	"fmt"
	//log "github.com/cihub/seelog"
	//log "github.com/Sirupsen/logrus"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	runtime.GOMAXPROCS(4)

	// blog
	writer, err := blog4go.NewFileLogWriter("output.log", true)
	if nil != err {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer writer.Close()

	for i := 1; i < 10; i++ {
		//logging(writer)
		go logging(writer)
	}

	// seelog
	//logger, err := log.LoggerFromConfigAsFile("log_config.xml")
	//if nil != err {
	//fmt.Println(err.Error())
	//}

	//for i := 1; i < 10; i++ {
	//go logging1(logger)
	//}

	// logrus
	//file, err := os.OpenFile("output.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0644))
	//if nil != err {
	//fmt.Println(err.Error())
	//}
	//defer file.Close()
	//log.SetOutput(file)

	//for i := 1; i < 10; i++ {
	//go logging2()
	//}

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
func logging(writer *blog4go.FileLogWriter) {
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

// seelog
//func logging1(writer log.LoggerInterface) {
//for {
//writer.Debug("test")
//writer.Debugf("haha %s. en\\en, always %d and %.4f", "eddie", 18, 3.1415)
//}
//}

// logrus
//func logging2() {
//for {
//log.Print("test\n")
//log.Printf("%s [%s] haha %s. en\\en, always %d and %.4f\n", time.Now().Format("2006-01-02 15:04:05"), "DEBUG", "eddie", 18, 3.1415)
//}
//}
