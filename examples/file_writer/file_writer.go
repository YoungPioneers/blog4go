package main

import (
	"fmt"
	log "github.com/YoungPioneers/blog4go"
	"os"
	"time"
)

type MyHook struct {
	something string
}

func (self *MyHook) Fire(level log.Level, message string) {
	fmt.Println(message)
}

type T struct {
	A int
	B string
}

func main() {
	hook := new(MyHook)

	err := log.NewFileWriter("./")
	if nil != err {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer log.Close()

	// optionally define your logging hook
	log.SetHook(hook)
	log.SetHookLevel(log.INFO)

	// optionally set output colored
	log.SetColored(true)

	log.Debug("Debug")
	log.Trace("Trace")
	log.Info("Info")
	log.Warn("Warn")
	log.Error("Error")
	log.Critical("Critical")

	// wait for hook runs
	time.Sleep(1 * time.Second)
}
