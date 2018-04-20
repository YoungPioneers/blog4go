Introduction
=======

BLog4go is an efficient logging library written in the [Go](http://golang.org/) programming language, providing logging hook, log rotate, filtering and formatting log message.

BLog4go 是高性能日志库。创新地使用“边解析边输出”方法进行日志输出，同时支持回调函数、日志淘汰和配置文件。可以解决高并发，调用日志函数频繁的情境下，日志库造成的性能问题。

[![Build Status](https://travis-ci.org/YoungPioneers/blog4go.svg?branch=master)](https://travis-ci.org/YoungPioneers/blog4go)
[![CircleCI](https://circleci.com/gh/YoungPioneers/blog4go.svg?style=svg)](https://circleci.com/gh/YoungPioneers/blog4go)
[![Coverage Status](https://coveralls.io/repos/github/YoungPioneers/blog4go/badge.svg?branch=master)](https://coveralls.io/github/YoungPioneers/blog4go?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/YoungPioneers/blog4go)](https://goreportcard.com/report/github.com/YoungPioneers/blog4go)
[![GoDoc](https://godoc.org/github.com/YoungPioneers/blog4go?status.svg)](https://godoc.org/github.com/YoungPioneers/blog4go)

Features
------------------
* *Partially write* to the [bufio.Writer](https://golang.org/pkg/bufio/#Writer) as soon as posible while formatting message to improve performance
* Support different logging output file for different logging level
* Support configure with files in xml format
* Configurable logrotate strategy
* Call user defined hook in asynchronous mode for every logging action
* Adjustable message formatting
* Configurable logging behavier when logging *on the fly* without restarting
* Suit configuration to the environment when logging start
* Try best to get every done in background
* File writer can be configured according to given config file
* Different output writers
	* Console writer
	* File writer
	* Socket writer

Quick-start
------------------

```
package main

import (
	log "github.com/YoungPioneers/blog4go"
	"fmt"
	"os"
)

// optionally set user defined hook for logging
type MyHook struct {
	something string
}

// when log-level exceed level, call the hook
// level is the level associate with that logging action.
// message is the formatted string already written.
func (self *MyHook) Fire(level log.LevelType, args ...interface{}) {
	fmt.Println(args...)
}

func main() {
	// init a file write using xml config file
	// log.SetBufferSize(0) // close buffer for in time logging when debugging
	err := log.NewWriterFromConfigAsFile("config.xml")
	if nil != err {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer log.Close()

	// initialize your hook instance
	hook := new(MyHook)
	log.SetHook(hook) // writersFromConfig can be replaced with writers
	log.SetHookLevel(log.INFO)
	log.SetHookAsync(true) // hook will be called in async mode

	// optionally set output colored
	log.SetColored(true)

	log.Debugf("Good morning, %s", "eddie")
	log.Warn("It's time to have breakfast")
}
```

config.xml
```xml
<blog4go minlevel="info">
	<filter levels="trace">
		<rotatefile path="trace.log" type="time"></rotatefile>
	</filter>
	<filter levels="debug,info" colored="true">
		<file path="debug.log"></file>
	</filter>
	<filter levels="error,critical">
		<rotatefile path="error.log" type="size" rotateSize="50000000" rotateLines="8000000"></rotatefile>
	</filter>
</blog4go>
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

I do some benchmark on a SSD disk with my macbook pro comparing amoung fmt,blog4go,seelog,logrus. [Benchmark Code](https://github.com/YoungPioneers/blog4go-benchmark)

```
BenchmarkBlog4goSingleGoroutine-4                        300000          4981 ns/op
BenchmarkBlog4goMultiGoroutine-4                           3000        554542 ns/op
BenchmarkFmtFormatSingleGoroutine-4                      300000          3727 ns/op
BenchmarkFmtFormatWithTimecacheSingleGoroutine-4         500000          2951 ns/op
BenchmarkFmtFormatWithTimecacheMultiGoroutine-4            3000        421204 ns/op
BenchmarkLogrusSingleGoroutine-4                         100000         18652 ns/op
BenchmarkLogrusWithTimecacheSingleGoroutine-4            100000         16024 ns/op
BenchmarkLogrusWithTimecacheMultiGoroutine-4                500       2238614 ns/op
BenchmarkSeelogSingleGoroutine-4                          50000         23476 ns/op
BenchmarkSeelogMultiGoroutine-4                             500       2722851 ns/op
```

It shows that blog4go can write log very fast, especially in situation with multi goroutines running at the same time~


Documentation
------------------

TODO


Examples
---------------

Full examples please view [EXAMPLES](https://github.com/YoungPioneers/blog4go-examples)


Changelog
------------------

[CHANGELOG](https://raw.githubusercontent.com/YoungPioneers/blog4go/master/CHANGELOG)
