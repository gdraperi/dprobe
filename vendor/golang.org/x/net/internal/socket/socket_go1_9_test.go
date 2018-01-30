// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.9
// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package socket_test

import (
	"bytes"
	"fmt"
	"net"
	"runtime"
	"testing"

	"golang.org/x/net/internal/nettest"
	"golang.org/x/net/internal/socket"
)

type mockControl struct ***REMOVED***
	Level int
	Type  int
	Data  []byte
***REMOVED***

func TestControlMessage(t *testing.T) ***REMOVED***
	for _, tt := range []struct ***REMOVED***
		cs []mockControl
	***REMOVED******REMOVED***
		***REMOVED***
			[]mockControl***REMOVED***
				***REMOVED***Level: 1, Type: 1***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			[]mockControl***REMOVED***
				***REMOVED***Level: 2, Type: 2, Data: []byte***REMOVED***0xfe***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			[]mockControl***REMOVED***
				***REMOVED***Level: 3, Type: 3, Data: []byte***REMOVED***0xfe, 0xff, 0xff, 0xfe***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			[]mockControl***REMOVED***
				***REMOVED***Level: 4, Type: 4, Data: []byte***REMOVED***0xfe, 0xff, 0xff, 0xfe, 0xfe, 0xff, 0xff, 0xfe***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			[]mockControl***REMOVED***
				***REMOVED***Level: 4, Type: 4, Data: []byte***REMOVED***0xfe, 0xff, 0xff, 0xfe, 0xfe, 0xff, 0xff, 0xfe***REMOVED******REMOVED***,
				***REMOVED***Level: 2, Type: 2, Data: []byte***REMOVED***0xfe***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		var w []byte
		var tailPadLen int
		mm := socket.NewControlMessage([]int***REMOVED***0***REMOVED***)
		for i, c := range tt.cs ***REMOVED***
			m := socket.NewControlMessage([]int***REMOVED***len(c.Data)***REMOVED***)
			l := len(m) - len(mm)
			if i == len(tt.cs)-1 && l > len(c.Data) ***REMOVED***
				tailPadLen = l - len(c.Data)
			***REMOVED***
			w = append(w, m...)
		***REMOVED***

		var err error
		ww := make([]byte, len(w))
		copy(ww, w)
		m := socket.ControlMessage(ww)
		for _, c := range tt.cs ***REMOVED***
			if err = m.MarshalHeader(c.Level, c.Type, len(c.Data)); err != nil ***REMOVED***
				t.Fatalf("(%v).MarshalHeader() = %v", tt.cs, err)
			***REMOVED***
			copy(m.Data(len(c.Data)), c.Data)
			m = m.Next(len(c.Data))
		***REMOVED***
		m = socket.ControlMessage(w)
		for _, c := range tt.cs ***REMOVED***
			m, err = m.Marshal(c.Level, c.Type, c.Data)
			if err != nil ***REMOVED***
				t.Fatalf("(%v).Marshal() = %v", tt.cs, err)
			***REMOVED***
		***REMOVED***
		if !bytes.Equal(ww, w) ***REMOVED***
			t.Fatalf("got %#v; want %#v", ww, w)
		***REMOVED***

		ws := [][]byte***REMOVED***w***REMOVED***
		if tailPadLen > 0 ***REMOVED***
			// Test a message with no tail padding.
			nopad := w[:len(w)-tailPadLen]
			ws = append(ws, [][]byte***REMOVED***nopad***REMOVED***...)
		***REMOVED***
		for _, w := range ws ***REMOVED***
			ms, err := socket.ControlMessage(w).Parse()
			if err != nil ***REMOVED***
				t.Fatalf("(%v).Parse() = %v", tt.cs, err)
			***REMOVED***
			for i, m := range ms ***REMOVED***
				lvl, typ, dataLen, err := m.ParseHeader()
				if err != nil ***REMOVED***
					t.Fatalf("(%v).ParseHeader() = %v", tt.cs, err)
				***REMOVED***
				if lvl != tt.cs[i].Level || typ != tt.cs[i].Type || dataLen != len(tt.cs[i].Data) ***REMOVED***
					t.Fatalf("%v: got %d, %d, %d; want %d, %d, %d", tt.cs[i], lvl, typ, dataLen, tt.cs[i].Level, tt.cs[i].Type, len(tt.cs[i].Data))
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestUDP(t *testing.T) ***REMOVED***
	c, err := nettest.NewLocalPacketListener("udp")
	if err != nil ***REMOVED***
		t.Skipf("not supported on %s/%s: %v", runtime.GOOS, runtime.GOARCH, err)
	***REMOVED***
	defer c.Close()
	cc, err := socket.NewConn(c.(net.Conn))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	t.Run("Message", func(t *testing.T) ***REMOVED***
		data := []byte("HELLO-R-U-THERE")
		wm := socket.Message***REMOVED***
			Buffers: bytes.SplitAfter(data, []byte("-")),
			Addr:    c.LocalAddr(),
		***REMOVED***
		if err := cc.SendMsg(&wm, 0); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		b := make([]byte, 32)
		rm := socket.Message***REMOVED***
			Buffers: [][]byte***REMOVED***b[:1], b[1:3], b[3:7], b[7:11], b[11:]***REMOVED***,
		***REMOVED***
		if err := cc.RecvMsg(&rm, 0); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if !bytes.Equal(b[:rm.N], data) ***REMOVED***
			t.Fatalf("got %#v; want %#v", b[:rm.N], data)
		***REMOVED***
	***REMOVED***)
	switch runtime.GOOS ***REMOVED***
	case "android", "linux":
		t.Run("Messages", func(t *testing.T) ***REMOVED***
			data := []byte("HELLO-R-U-THERE")
			wmbs := bytes.SplitAfter(data, []byte("-"))
			wms := []socket.Message***REMOVED***
				***REMOVED***Buffers: wmbs[:1], Addr: c.LocalAddr()***REMOVED***,
				***REMOVED***Buffers: wmbs[1:], Addr: c.LocalAddr()***REMOVED***,
			***REMOVED***
			n, err := cc.SendMsgs(wms, 0)
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			if n != len(wms) ***REMOVED***
				t.Fatalf("got %d; want %d", n, len(wms))
			***REMOVED***
			b := make([]byte, 32)
			rmbs := [][][]byte***REMOVED******REMOVED***b[:len(wmbs[0])]***REMOVED***, ***REMOVED***b[len(wmbs[0]):]***REMOVED******REMOVED***
			rms := []socket.Message***REMOVED***
				***REMOVED***Buffers: rmbs[0]***REMOVED***,
				***REMOVED***Buffers: rmbs[1]***REMOVED***,
			***REMOVED***
			n, err = cc.RecvMsgs(rms, 0)
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			if n != len(rms) ***REMOVED***
				t.Fatalf("got %d; want %d", n, len(rms))
			***REMOVED***
			nn := 0
			for i := 0; i < n; i++ ***REMOVED***
				nn += rms[i].N
			***REMOVED***
			if !bytes.Equal(b[:nn], data) ***REMOVED***
				t.Fatalf("got %#v; want %#v", b[:nn], data)
			***REMOVED***
		***REMOVED***)
	***REMOVED***

	// The behavior of transmission for zero byte paylaod depends
	// on each platform implementation. Some may transmit only
	// protocol header and options, other may transmit nothing.
	// We test only that SendMsg and SendMsgs will not crash with
	// empty buffers.
	wm := socket.Message***REMOVED***
		Buffers: [][]byte***REMOVED******REMOVED******REMOVED******REMOVED***,
		Addr:    c.LocalAddr(),
	***REMOVED***
	cc.SendMsg(&wm, 0)
	wms := []socket.Message***REMOVED***
		***REMOVED***Buffers: [][]byte***REMOVED******REMOVED******REMOVED******REMOVED***, Addr: c.LocalAddr()***REMOVED***,
	***REMOVED***
	cc.SendMsgs(wms, 0)
***REMOVED***

func BenchmarkUDP(b *testing.B) ***REMOVED***
	c, err := nettest.NewLocalPacketListener("udp")
	if err != nil ***REMOVED***
		b.Skipf("not supported on %s/%s: %v", runtime.GOOS, runtime.GOARCH, err)
	***REMOVED***
	defer c.Close()
	cc, err := socket.NewConn(c.(net.Conn))
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	data := []byte("HELLO-R-U-THERE")
	wm := socket.Message***REMOVED***
		Buffers: [][]byte***REMOVED***data***REMOVED***,
		Addr:    c.LocalAddr(),
	***REMOVED***
	rm := socket.Message***REMOVED***
		Buffers: [][]byte***REMOVED***make([]byte, 128)***REMOVED***,
		OOB:     make([]byte, 128),
	***REMOVED***

	for M := 1; M <= 1<<9; M = M << 1 ***REMOVED***
		b.Run(fmt.Sprintf("Iter-%d", M), func(b *testing.B) ***REMOVED***
			for i := 0; i < b.N; i++ ***REMOVED***
				for j := 0; j < M; j++ ***REMOVED***
					if err := cc.SendMsg(&wm, 0); err != nil ***REMOVED***
						b.Fatal(err)
					***REMOVED***
					if err := cc.RecvMsg(&rm, 0); err != nil ***REMOVED***
						b.Fatal(err)
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***)
		switch runtime.GOOS ***REMOVED***
		case "android", "linux":
			wms := make([]socket.Message, M)
			for i := range wms ***REMOVED***
				wms[i].Buffers = [][]byte***REMOVED***data***REMOVED***
				wms[i].Addr = c.LocalAddr()
			***REMOVED***
			rms := make([]socket.Message, M)
			for i := range rms ***REMOVED***
				rms[i].Buffers = [][]byte***REMOVED***make([]byte, 128)***REMOVED***
				rms[i].OOB = make([]byte, 128)
			***REMOVED***
			b.Run(fmt.Sprintf("Batch-%d", M), func(b *testing.B) ***REMOVED***
				for i := 0; i < b.N; i++ ***REMOVED***
					if _, err := cc.SendMsgs(wms, 0); err != nil ***REMOVED***
						b.Fatal(err)
					***REMOVED***
					if _, err := cc.RecvMsgs(rms, 0); err != nil ***REMOVED***
						b.Fatal(err)
					***REMOVED***
				***REMOVED***
			***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***
