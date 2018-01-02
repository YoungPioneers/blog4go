// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"sync"
	"time"
)

const (
	// PrefixTimeFormat const time format prefix
	PrefixTimeFormat = "time=\"2006-01-02 15:04:05\""

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

	// lock for read && write
	lock *sync.RWMutex
}

// global time cache instance used for every log writer
var timeCache = timeFormatCacheType{}

func init() {
	timeCache.lock = new(sync.RWMutex)

	timeCache.now = time.Now()
	timeCache.date = timeCache.now.Format(DateFormat)
	timeCache.format = []byte(timeCache.now.Format(PrefixTimeFormat))
	timeCache.dateYesterday = timeCache.now.Add(-24 * time.Hour).Format(DateFormat)

	// update timeCache every seconds
	go func() {
		// tick every seconds
		t := time.Tick(1 * time.Millisecond)

		//UpdateTimeCacheLoop:
		for {
			select {
			case <-t:
				timeCache.fresh()
			}
		}
	}()
}

// Now now
func (timeCache *timeFormatCacheType) Now() time.Time {
	timeCache.lock.RLock()
	defer timeCache.lock.RUnlock()

	return timeCache.now
}

// Date date
func (timeCache *timeFormatCacheType) Date() string {
	timeCache.lock.RLock()
	defer timeCache.lock.RUnlock()

	return timeCache.date
}

// DateYesterday date
func (timeCache *timeFormatCacheType) DateYesterday() string {
	timeCache.lock.RLock()
	defer timeCache.lock.RUnlock()

	return timeCache.dateYesterday
}

// Format format
func (timeCache *timeFormatCacheType) Format() []byte {
	timeCache.lock.RLock()
	defer timeCache.lock.RUnlock()

	return timeCache.format
}

// fresh data in timeCache
func (timeCache *timeFormatCacheType) fresh() {
	timeCache.lock.Lock()
	defer timeCache.lock.Unlock()

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
