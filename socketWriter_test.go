// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestSocketWriterBasicOperation(t *testing.T) {
	_, err := NewSocketWriter("udp", "127.0.0.1:12124")
	defer Close()
	if nil != err {
		t.Error(err.Error())
	}

	// test socket writer hook
	hook := new(MyHook)
	hook.cnt = 0

	blog.SetHook(hook)
	blog.SetHookLevel(INFO)

	blog.Debugf("%s", "something")
	// wait for hook called
	time.Sleep(1 * time.Millisecond)
	if 0 != hook.cnt {
		t.Error("hook called not valid")
	}

	if DEBUG == hook.level || "something" == hook.message {
		t.Errorf("hook parameters wrong. level: %s, message: %s", hook.level.String(), hook.message)
	}

	blog.Info("yes")
	// wait for hook called
	time.Sleep(1 * time.Millisecond)
	if 1 != hook.cnt {
		t.Error("hook not called")
	}

	if INFO != hook.level || "yes" != hook.message {
		t.Errorf("hook parameters wrong. level: %d, message: %s", hook.level, hook.message)
	}

	// test basic operations
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

	blog.SetColored(true)
	blog.SetTimeRotated(true)
	blog.SetLevel(CRITICAL)
	blog.SetRetentions(7)
	blog.SetRotateLines(100000)
	blog.SetRotateSize(ByteSize(1024 * 1024 * 500))

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
}

func TestSignleSocketWriter(t *testing.T) {
	_, err := NewSocketWriter("udp", "127.0.0.1:12124")
	defer Close()
	if nil != err {
		t.Error(err.Error())
	}

	var wg sync.WaitGroup
	var wgListen sync.WaitGroup

	wg.Add(1)
	wgListen.Add(1)
	go func() {
		defer wg.Done()

		initPrefix(false)

		// begin listen udp packages on 127.0.0.1:12124
		serverAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:12124")
		conn, err := net.ListenUDP("udp", serverAddr)
		if nil != err {
			fmt.Println(err.Error())
			t.Error(err.Error())
		}
		wgListen.Done()

		var bytes = make([]byte, 1024)
		_, err = conn.Read(bytes)
		if nil != err {
			fmt.Println(err.Error())
			t.Error(err.Error())
		}

		str := string(bytes)
		arrs := strings.Split(str, "[DEBUG] ")
		if len(arrs) != 2 {
			t.Errorf("udp message format wrong. str: %s", str)
			return
		}

		// FIXME this may not be accurate
		if arrs[1][:4] != "haha" {
			t.Errorf("udp message content wrong. str: %s", arrs[1][:4])
			return
		}
	}()

	wgListen.Wait()
	blog.Debug("haha")
	wg.Wait()

	// chekc init socket writer multi time
	_, err = NewSocketWriter("udp", "127.0.0.1:12124")
	defer Close()
	if ErrAlreadyInit != err {
		t.Error("duplicate init check fail")
	}
}

func BenchmarkSocketWriter(b *testing.B) {
	_, err := NewSocketWriter("udp", "127.0.0.1:12124")
	defer Close()
	if nil != err {
		b.Error(err.Error())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		blog.Debugf("haha %s. en\\en, always %d and %f", "eddie", 18, 3.1415)
	}
}
