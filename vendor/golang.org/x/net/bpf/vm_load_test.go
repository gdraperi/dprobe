// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bpf_test

import (
	"net"
	"testing"

	"golang.org/x/net/bpf"
	"golang.org/x/net/ipv4"
)

func TestVMLoadAbsoluteOffsetOutOfBounds(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Off:  100,
			Size: 2,
		***REMOVED***,
		bpf.RetA***REMOVED******REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("failed to load BPF program: %v", err)
	***REMOVED***
	defer done()

	out, err := vm.Run([]byte***REMOVED***
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0, 1, 2, 3,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 0, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMLoadAbsoluteOffsetPlusSizeOutOfBounds(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 2,
		***REMOVED***,
		bpf.RetA***REMOVED******REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("failed to load BPF program: %v", err)
	***REMOVED***
	defer done()

	out, err := vm.Run([]byte***REMOVED***
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 0, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMLoadAbsoluteBadInstructionSize(t *testing.T) ***REMOVED***
	_, _, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Size: 5,
		***REMOVED***,
		bpf.RetA***REMOVED******REMOVED***,
	***REMOVED***)
	if errStr(err) != "assembling instruction 1: invalid load byte length 0" ***REMOVED***
		t.Fatalf("unexpected error: %v", err)
	***REMOVED***
***REMOVED***

func TestVMLoadConstantOK(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadConstant***REMOVED***
			Dst: bpf.RegX,
			Val: 9,
		***REMOVED***,
		bpf.TXA***REMOVED******REMOVED***,
		bpf.RetA***REMOVED******REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("failed to load BPF program: %v", err)
	***REMOVED***
	defer done()

	out, err := vm.Run([]byte***REMOVED***
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 1, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMLoadIndirectOutOfBounds(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadIndirect***REMOVED***
			Off:  100,
			Size: 1,
		***REMOVED***,
		bpf.RetA***REMOVED******REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("failed to load BPF program: %v", err)
	***REMOVED***
	defer done()

	out, err := vm.Run([]byte***REMOVED***
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 0, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMLoadMemShiftOutOfBounds(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadMemShift***REMOVED***
			Off: 100,
		***REMOVED***,
		bpf.RetA***REMOVED******REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("failed to load BPF program: %v", err)
	***REMOVED***
	defer done()

	out, err := vm.Run([]byte***REMOVED***
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 0, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

const (
	dhcp4Port = 53
)

func TestVMLoadMemShiftLoadIndirectNoResult(t *testing.T) ***REMOVED***
	vm, in, done := testDHCPv4(t)
	defer done()

	// Append mostly empty UDP header with incorrect DHCPv4 port
	in = append(in, []byte***REMOVED***
		0, 0,
		0, dhcp4Port + 1,
		0, 0,
		0, 0,
	***REMOVED***...)

	out, err := vm.Run(in)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 0, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMLoadMemShiftLoadIndirectOK(t *testing.T) ***REMOVED***
	vm, in, done := testDHCPv4(t)
	defer done()

	// Append mostly empty UDP header with correct DHCPv4 port
	in = append(in, []byte***REMOVED***
		0, 0,
		0, dhcp4Port,
		0, 0,
		0, 0,
	***REMOVED***...)

	out, err := vm.Run(in)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := len(in)-8, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func testDHCPv4(t *testing.T) (virtualMachine, []byte, func()) ***REMOVED***
	// DHCPv4 test data courtesy of David Anderson:
	// https://github.com/google/netboot/blob/master/dhcp4/conn_linux.go#L59-L70
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		// Load IPv4 packet length
		bpf.LoadMemShift***REMOVED***Off: 8***REMOVED***,
		// Get UDP dport
		bpf.LoadIndirect***REMOVED***Off: 8 + 2, Size: 2***REMOVED***,
		// Correct dport?
		bpf.JumpIf***REMOVED***Cond: bpf.JumpEqual, Val: dhcp4Port, SkipFalse: 1***REMOVED***,
		// Accept
		bpf.RetConstant***REMOVED***Val: 1500***REMOVED***,
		// Ignore
		bpf.RetConstant***REMOVED***Val: 0***REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("failed to load BPF program: %v", err)
	***REMOVED***

	// Minimal requirements to make a valid IPv4 header
	h := &ipv4.Header***REMOVED***
		Len: ipv4.HeaderLen,
		Src: net.IPv4(192, 168, 1, 1),
		Dst: net.IPv4(192, 168, 1, 2),
	***REMOVED***
	hb, err := h.Marshal()
	if err != nil ***REMOVED***
		t.Fatalf("failed to marshal IPv4 header: %v", err)
	***REMOVED***

	hb = append([]byte***REMOVED***
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
	***REMOVED***, hb...)

	return vm, hb, done
***REMOVED***
