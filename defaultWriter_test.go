// Copyright (c) 2023, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"testing"
)

func TestDefaultWriterBasicOperation(t *testing.T) {
	blog = &DefaultWriter{}
	defer blog.Close()

	// test basic operations
	blog.SetTags(map[string]string{"tagName": "tagValue"})
	blog.Tags()

	blog.Debug("Debug", 1)
	blog.Debugf("%s", "Debug")
	blog.Trace("Trace", 2)
	blog.Tracef("%s", "Trace")
	blog.Info("Info", 3)
	blog.Infof("%s", "Info")
	blog.Warn("Warn", 4)
	blog.Warnf("%s", "Warn")
	blog.Error("Error", 5)
	blog.Errorf("%s", "Error")
	blog.Critical("Critical", 6)
	blog.Criticalf("%s", "Critical")
	blog.flush()

	blog.SetHookAsync(true)
	blog.Colored()
	blog.SetColored(true)
	blog.TimeRotated()
	blog.SetTimeRotated(true)
	blog.Level()
	blog.SetLevel(CRITICAL)
	blog.Retentions()
	blog.SetRetentions(0)
	blog.SetRetentions(7)
	blog.RotateLines()
	blog.SetRotateLines(0)
	blog.SetRotateLines(100000)
	blog.RotateSize()
	blog.SetRotateSize(0)
	blog.SetRotateSize(1024 * 1024 * 500)

	blog.Debug("Debug", 1)
	blog.Debugf("%s\\", "Debug")
	blog.Trace("Trace", 2)
	blog.Tracef("%s", "Trace")
	blog.Info("Info", 3)
	blog.Infof("%s", "Info")
	blog.Warn("Warn", 4)
	blog.Warnf("%s", "Warn")
	blog.Error("Error", 5)
	blog.Errorf("%s", "Error")
	blog.Critical("Critical", 6)
	blog.Criticalf("%s", "Critical")
}
