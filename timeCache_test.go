// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"testing"
	"time"
)

func TestTimeCacheUpdate(t *testing.T) {
	// first time
	// time
	if (time.Time{}) == timeCache.Now() {
		t.Error("time cache not correct when init, time wrong")
	}

	// date
	if timeCache.Date() != time.Now().Format(DateFormat) {
		t.Error("time cache not correct when init, date wrong")
	}

	// dateYesterday
	if timeCache.DateYesterday() != time.Now().Add(-24*time.Hour).Format(DateFormat) {
		t.Error("time cache not correct when init, dateYesterday wrong")
	}

	time.Sleep(1500 * time.Millisecond)

	// when updated
	// date
	if timeCache.Date() != time.Now().Format(DateFormat) {
		t.Error("time cache not correct when updated, date wrong")
	}

	// dateYesterday
	if timeCache.DateYesterday() != time.Now().Add(-24*time.Hour).Format(DateFormat) {
		t.Error("time cache not correct when updated, dateYesterday wrong")
	}
}
