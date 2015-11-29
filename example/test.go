package main

import (
	"blog4go"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	writer, err := blog4go.NewFileLogWriter("output.log")
	defer writer.Close()
	if nil != err {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	writer.Debug("test")
	writer.Debugf("haha %s. en\\en, always %d and %.4f", "eddie", 18, 3.1415)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)

	for i := 1; i < 20; i++ {
		go logging(writer)
	}

	for {
		select {
		case <-c:
			fmt.Println("Exit..")
			writer.Close()
			return
		}
	}

}

func logging(writer *blog4go.FileLogWriter) {
	for {
		writer.Debug("test")
		writer.Debugf("haha %s. en\\en, always %d and %.4f", "eddie", 18, 3.1415)
	}
}
