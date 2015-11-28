package main

import (
	"blog4go"
	"fmt"
	"os"
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
}
