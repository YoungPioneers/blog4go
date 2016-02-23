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

type T struct {
	A int
	B string
}

func (self *MyHook) Fire(level log.Level, message string) {
	fmt.Println(message)
}

func main() {
	hook := new(MyHook)

	// nc -u -l 12124 , to receive udp data
	_, err := log.NewSocketWriter("udp", "127.0.0.1:12124")
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

	log.Debug("Debug\n")
	log.Trace("Trace\n")
	log.Info("Info\n")
	log.Warn("Warn\n")
	log.Error("Error\n")
	log.Critical("Critical\n")

	// wait for hook runs
	time.Sleep(1 * time.Second)
}
