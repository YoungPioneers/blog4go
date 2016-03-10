// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"testing"
)

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
		//_, _, err = conn.ReadFromUDP(bytes)
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
