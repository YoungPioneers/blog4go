// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

// Hook Interface determine types of functions should be declared and
// implemented when user offers user defined function call before every
// logging action end.
// users may use this hook as a callback function when something happen.
// Fire function received two parameters.
// level is the level associate with that logging action.
// message is the formatted string already written.
type Hook interface {
	Fire(level LevelType, tags map[string]string, args ...interface{})
}
