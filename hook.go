// Copyright 2015
// Author: huangjunwei@youmi.net

package blog4go

// 用户可定义hook函数，每次log之后都会被调用
type Hook interface {
	Fire(level Level, message string)
}
