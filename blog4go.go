// Copyright 2015
// Author: huangjunwei@youmi.net

package blog4go

import ()

type LogWriter interface {
	// 用于写log的方法
	Write()
	// 关闭log writer的处理方法
	Close()
}
