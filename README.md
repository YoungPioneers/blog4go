Introduction
=======

BLog4go is authorized by [YOUMI](https://www.youmi.net/). It is an efficient logging library written in the [Go](http://golang.org/) programming language, providing logging hook, log rotate, filtering and formatting log message. 

[![Build Status](https://travis-ci.org/YoungPioneers/blog4go.svg?branch=master)](https://travis-ci.org/YoungPioneers/blog4go)


Features
------------------
* *Partially write* to the [bufio.Writer](https://golang.org/pkg/bufio/#Writer) as soon as posible while formatting message to improve performance
* Support different logging output file for different logging level
* Configurable logrotate strategy
* Call user defined hook in asynchronous mode for every logging action
* Adjustable message formatting
* Configurable logging behavier when looging *on the fly* without restarting
* Suit configuration to the environment when logging start
* Try best to get every done in background
* File writer can be configured according to given config file
* Different output writers
	* Console writer
	* File writer
	* Socket writer 


Quick-start
------------------

```go
package main

import (
	"github.com/YoungPioneers/blog4go"
	"fmt"
	"os"
)

// optionally set user defined hook for logging
type MyHook struct {
	something string
}

// when log-level exceed level, call the hook
func (self *MyHook) Fire(level blog4go.Level, message string) {
	fmt.Println(message)
}

func main() {
	writer, err := blog4go.NewBaseFileWriter("output.log")
	if nil != err {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer writer.Close()
	
	// optionally set logrotate every day
	writer.SetTimeRotated(true)
	
	// optionally set hook for logging
	hook := new(MyHook)
	writer.SetHook(hook)
	writer.SetHookLevel(blog4go.INFO)
	writer.Debugf("Good morning, %s", "eddie")	
	
	
	// init a file write using xml config file
	writersFromConfig, err := blog4go.NewFileWriterFromConfigAsFile("config.xml")
	if nil != err {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer writersFromConfig.Close()
		
	
	// init a file writer just give it a logging base directory
	writers, err := blog4go.blog4go.NewBaseFileWriter("output.log")
	if nil != err {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	
	defer writersFromConfig.Close()
	writersFromConfig.SetHook(hook) // writersFromConfig can be replaced with writers
	writersFromConfig.SetHookLevel(blog4go.INFO)
	writersFromConfig.Debugf("Good morning, %s", "eddie")	
	
}
```


Installation
------------------

If you don't have the Go development environment installed, visit the 
[Getting Started](http://golang.org/doc/install.html) document and follow the instructions. Once you're ready, execute the following command:

```
go get -u github.com/YoungPioneers/blog4go
```

Benchmark
------------------

TODO

Documentation
------------------

TODO


Examples
---------------

[EXAMPLES](https://github.com/YoungPioneers/blog4go/tree/master/example)


Changelog
------------------

[CHANGELOG](https://raw.githubusercontent.com/YoungPioneers/blog4go/master/CHANGELOG)
