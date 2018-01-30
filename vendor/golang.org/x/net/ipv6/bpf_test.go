// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv6_test

import (
	"net"
	"runtime"
	"testing"
	"time"

	"golang.org/x/net/bpf"
	"golang.org/x/net/ipv6"
)

func TestBPF(t *testing.T) ***REMOVED***
	if runtime.GOOS != "linux" ***REMOVED***
		t.Skipf("not supported on %s", runtime.GOOS)
	***REMOVED***
	if !supportsIPv6 ***REMOVED***
		t.Skip("ipv6 is not supported")
	***REMOVED***

	l, err := net.ListenPacket("udp6", "[::1]:0")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer l.Close()

	p := ipv6.NewPacketConn(l)

	// This filter accepts UDP packets whose first payload byte is
	// even.
	prog, err := bpf.Assemble([]bpf.Instruction***REMOVED***
		// Load the first byte of the payload (skipping UDP header).
		bpf.LoadAbsolute***REMOVED***Off: 8, Size: 1***REMOVED***,
		// Select LSB of the byte.
		bpf.ALUOpConstant***REMOVED***Op: bpf.ALUOpAnd, Val: 1***REMOVED***,
		// Byte is even?
		bpf.JumpIf***REMOVED***Cond: bpf.JumpEqual, Val: 0, SkipFalse: 1***REMOVED***,
		// Accept.
		bpf.RetConstant***REMOVED***Val: 4096***REMOVED***,
		// Ignore.
		bpf.RetConstant***REMOVED***Val: 0***REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("compiling BPF: %s", err)
	***REMOVED***

	if err = p.SetBPF(prog); err != nil ***REMOVED***
		t.Fatalf("attaching filter to Conn: %s", err)
	***REMOVED***

	s, err := net.Dial("udp6", l.LocalAddr().String())
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer s.Close()
	go func() ***REMOVED***
		for i := byte(0); i < 10; i++ ***REMOVED***
			s.Write([]byte***REMOVED***i***REMOVED***)
		***REMOVED***
	***REMOVED***()

	l.SetDeadline(time.Now().Add(2 * time.Second))
	seen := make([]bool, 5)
	for ***REMOVED***
		var b [512]byte
		n, _, err := l.ReadFrom(b[:])
		if err != nil ***REMOVED***
			t.Fatalf("reading from listener: %s", err)
		***REMOVED***
		if n != 1 ***REMOVED***
			t.Fatalf("unexpected packet length, want 1, got %d", n)
		***REMOVED***
		if b[0] >= 10 ***REMOVED***
			t.Fatalf("unexpected byte, want 0-9, got %d", b[0])
		***REMOVED***
		if b[0]%2 != 0 ***REMOVED***
			t.Fatalf("got odd byte %d, wanted only even bytes", b[0])
		***REMOVED***
		seen[b[0]/2] = true

		seenAll := true
		for _, v := range seen ***REMOVED***
			if !v ***REMOVED***
				seenAll = false
				break
			***REMOVED***
		***REMOVED***
		if seenAll ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
***REMOVED***
