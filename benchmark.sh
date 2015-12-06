#/bin/sh
go test -test.bench=".*" -gcflags=-m
