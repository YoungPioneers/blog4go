// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"time"
)

const (
	// PrefixTimeFormat const time format prefix
	PrefixTimeFormat = "[2006/01/02:15:04:05]"

	// DateFormat date format
	DateFormat = "2006-01-02"
)

// timeFormatCacheType is a time formated cache
type timeFormatCacheType struct {
	// current time
	now time.Time
	// current date
	date string
	// current formated date
	format []byte
	// yesterdate
	dateYesterday string
}

// global time cache instance used for every log writer
var timeCache = timeFormatCacheType{}

func init() {
	timeCache.now = time.Now()
	timeCache.date = timeCache.now.Format(DateFormat)
	timeCache.format = []byte(timeCache.now.Format(PrefixTimeFormat))
	timeCache.dateYesterday = timeCache.now.Add(-24 * time.Hour).Format(DateFormat)

	// update timeCache every seconds
	go func() {
		// tick every seconds
		t := time.Tick(1 * time.Second)

		//UpdateTimeCacheLoop:
		for {
			select {
			case <-t:
				// get current time and update timeCache
				now := time.Now()
				timeCache.now = now
				timeCache.format = []byte(now.Format(PrefixTimeFormat))
				date := now.Format(DateFormat)
				if date != timeCache.date {
					timeCache.dateYesterday = timeCache.date
					timeCache.date = now.Format(DateFormat)
				}
			}
		}
	}()
}
