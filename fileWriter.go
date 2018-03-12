// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"fmt"
	"path"
	"strings"
	"sync"
)

// NewFileWriter initialize a file writer
// baseDir must be base directory of log files
// rotate determine if it will logrotate
func NewFileWriter(baseDir string, rotate bool) (err error) {
	singltonLock.Lock()
	defer singltonLock.Unlock()
	if nil != blog {
		return ErrAlreadyInit
	}

	fileWriter := new(MultiWriter)
	fileWriter.lock = new(sync.RWMutex)
	fileWriter.level = DEBUG
	fileWriter.closed = false

	fileWriter.writers = make(map[LevelType]Writer)
	for _, level := range Levels {
		fileName := fmt.Sprintf("%s.log", strings.ToLower(level.String()))
		writer, err := newBaseFileWriter(path.Join(baseDir, fileName), rotate)
		if nil != err {
			return err
		}
		fileWriter.writers[level] = writer
	}

	// log hook
	fileWriter.hook = nil
	fileWriter.hookLevel = DEBUG
	fileWriter.hookAsync = true

	blog = fileWriter
	return
}
