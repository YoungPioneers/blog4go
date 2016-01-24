// Copyright 2015
// Author: huangjunwei@youmi.net

// TODO 支持JSON, CSV等不同格式输出
// TODO 分离下代码文件
// TODO 支持多种输出方式, console, file, socket
// TODO 支持多文件输出
package blog4go

type BLog4go struct {
}

// 各种日志结构接口
type Writer interface {
	// 关闭log writer的处理方法
	// 善后
	Close()

	// 用于内部写log的方法
	write(level Level, format string)
	writef(level Level, format string, args ...interface{})
}
