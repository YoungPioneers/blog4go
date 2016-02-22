// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"fmt"
	"path"
	"strings"
)

// NewFileWriter initialize a file writer
// baseDir must be base directory of log files
func NewFileWriter(baseDir string) (err error) {
	singltonLock.Lock()
	defer singltonLock.Unlock()
	if nil != blog {
		return
	}

	fmt.Println("here")
	fileWriter := new(FileWriter)
	fileWriter.level = DEBUG
	fileWriter.closed = false

	fileWriter.writers = make(map[Level]*baseFileWriter)
	for _, level := range Levels {
		fileName := fmt.Sprintf("%s.log", strings.ToLower(level.String()))
		writer, err := newBaseFileWriter(path.Join(baseDir, fileName))
		if nil != err {
			return err
		}
		fileWriter.writers[level] = writer
	}

	// log hook
	fileWriter.hook = nil
	fileWriter.hookLevel = DEBUG

	blog = fileWriter
	return
}
