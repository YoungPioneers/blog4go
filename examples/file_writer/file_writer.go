package main

import (
	"fmt"
	log "github.com/YoungPioneers/blog4go"
	"os"
	"time"
)

// MyHook .
type MyHook struct {
	something string
}

// Fire .
func (hook *MyHook) Fire(level log.LevelType, args ...interface{}) {
	fmt.Println(args)
}

// T .
type T struct {
	A int
	B string
}

func main() {
	hook := new(MyHook)

	err := log.NewFileWriter("./", true)
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
