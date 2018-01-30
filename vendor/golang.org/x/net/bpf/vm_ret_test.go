// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bpf_test

import (
	"testing"

	"golang.org/x/net/bpf"
)

func TestVMRetA(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
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
		9,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 1, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMRetALargerThanInput(t *testing.T) ***REMOVED***
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
		0, 255,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 2, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMRetConstant(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.RetConstant***REMOVED***
			Val: 9,
		***REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("failed to load BPF program: %v", err)
	***REMOVED***
	defer done()

	out, err := vm.Run([]byte***REMOVED***
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0, 1,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 1, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMRetConstantLargerThanInput(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.RetConstant***REMOVED***
			Val: 16,
		***REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("failed to load BPF program: %v", err)
	***REMOVED***
	defer done()

	out, err := vm.Run([]byte***REMOVED***
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0, 1,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 2, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***
