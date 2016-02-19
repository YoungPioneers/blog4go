Introduction
=======

BLog4go is authorized by [YOUMI](https://www.youmi.net/). It is an efficient logging library written in the [Go](http://golang.org/) programming language, providing logging hook, log rotate, filtering and formatting log message. 

[![Build Status](https://travis-ci.org/YoungPioneers/blog4go.svg?branch=master)](https://travis-ci.org/YoungPioneers/blog4go)


Features
------------------
* *Partially write* to the [bufio.Writer](https://golang.org/pkg/bufio/#Writer) as soon as posible while formatting message to improve performance
* Support different logging output file for different logging level
* Support configure with files in xml format
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
	// init a file write using xml config file
	err := blog4go.NewFileWriterFromConfigAsFile("config.xml")
	// init a file writer just give it a logging base directory
	// err := blog4go.blog4go.NewBaseFileWriter("output.log")
	if nil != err {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer blog4go.Close()

	blog4go.SetHook(hook) // writersFromConfig can be replaced with writers
	blog4go.SetHookLevel(blog4go.INFO)
	blog4go.Debugf("Good morning, %s", "eddie")

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

I do some benchmark on a HDD disk comparing amoung fmt,blog4go,seelog,logrus. [Benchmark Code](https://github.com/YoungPioneers/blog4go/tree/master/benchmark)

```
BenchmarkBlog4goFormat-4         	 1000000	      1056 ns/op
BenchmarkFmtFormat-4             	  500000	      2491 ns/op
BenchmarkFmtWithTimecacheFormat-4	  300000	      4566 ns/op
BenchmarkLogrus-4                	  100000	     14216 ns/op
BenchmarkLogrusWithTimecache-4   	  100000	     14541 ns/op
BenchmarkSeelogFormat-4          	   50000	     34426 ns/op
```

It shows that blog4go can write log very fast~


Documentation
------------------

TODO


Examples
---------------

[EXAMPLES](https://github.com/YoungPioneers/blog4go/tree/master/example)


Changelog
------------------

[CHANGELOG](https://raw.githubusercontent.com/YoungPioneers/blog4go/master/CHANGELOG)
