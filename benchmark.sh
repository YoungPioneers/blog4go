#/bin/sh
go test -test.bench=".*"
#go test -test.bench=".*" -gcflags=-m
