// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"testing"
	"time"
)

func TestTimeCacheUpdate(t *testing.T) {
	// first time
	if timeCache.date != time.Now().Format(DateFormat) {
		t.Error("time cache not correct when init")
	}

	time.Sleep(1500 * time.Millisecond)

	// when updated
	if timeCache.date != time.Now().Format(DateFormat) {
		t.Error("time cache not correct when updated")
	}
}
