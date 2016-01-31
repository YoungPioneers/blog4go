// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

// Package blog4go provide an efficient and easy-to-use writers library for
// logging into files, console or sockets. Writers suports formatting
// string filtering and calling user defined hook in asynchronous mode.
package blog4go

// Struct BLog4go is just a placeholder now. It may be of important use
// in the near future.
type BLog4go struct {
}

// Interface Writer is a common definition of any writers in this package.
// Any struct implements Writer interface must implement functions below.
// Close is used for close the writer and free any elements if needed.
// write is an internal function that write pure message with specific
// logging level.
// writef is an internal function that formatting message with specific
// logging level. Placeholders in the format string will be replaced with
// args given.
// Both write and writef may have an asynchronous call of user defined
// function before write and writef function end..
type Writer interface {
	Close() // do anything end before program end

	write(level Level, format string)                       // write pure string
	writef(level Level, format string, args ...interface{}) // format string and write it
}
