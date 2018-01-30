// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bpf_test

import (
	"testing"

	"golang.org/x/net/bpf"
)

func TestVMJumpOne(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 1,
		***REMOVED***,
		bpf.Jump***REMOVED***
			Skip: 1,
		***REMOVED***,
		bpf.RetConstant***REMOVED***
			Val: 0,
		***REMOVED***,
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
		1,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 1, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMJumpOutOfProgram(t *testing.T) ***REMOVED***
	_, _, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.Jump***REMOVED***
			Skip: 1,
		***REMOVED***,
		bpf.RetA***REMOVED******REMOVED***,
	***REMOVED***)
	if errStr(err) != "cannot jump 1 instructions; jumping past program bounds" ***REMOVED***
		t.Fatalf("unexpected error: %v", err)
	***REMOVED***
***REMOVED***

func TestVMJumpIfTrueOutOfProgram(t *testing.T) ***REMOVED***
	_, _, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.JumpIf***REMOVED***
			Cond:     bpf.JumpEqual,
			SkipTrue: 2,
		***REMOVED***,
		bpf.RetA***REMOVED******REMOVED***,
	***REMOVED***)
	if errStr(err) != "cannot jump 2 instructions in true case; jumping past program bounds" ***REMOVED***
		t.Fatalf("unexpected error: %v", err)
	***REMOVED***
***REMOVED***

func TestVMJumpIfFalseOutOfProgram(t *testing.T) ***REMOVED***
	_, _, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.JumpIf***REMOVED***
			Cond:      bpf.JumpEqual,
			SkipFalse: 3,
		***REMOVED***,
		bpf.RetA***REMOVED******REMOVED***,
	***REMOVED***)
	if errStr(err) != "cannot jump 3 instructions in false case; jumping past program bounds" ***REMOVED***
		t.Fatalf("unexpected error: %v", err)
	***REMOVED***
***REMOVED***

func TestVMJumpIfEqual(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 1,
		***REMOVED***,
		bpf.JumpIf***REMOVED***
			Cond:     bpf.JumpEqual,
			Val:      1,
			SkipTrue: 1,
		***REMOVED***,
		bpf.RetConstant***REMOVED***
			Val: 0,
		***REMOVED***,
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
		1,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 1, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMJumpIfNotEqual(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 1,
		***REMOVED***,
		bpf.JumpIf***REMOVED***
			Cond:      bpf.JumpNotEqual,
			Val:       1,
			SkipFalse: 1,
		***REMOVED***,
		bpf.RetConstant***REMOVED***
			Val: 0,
		***REMOVED***,
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
		1,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 1, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMJumpIfGreaterThan(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 4,
		***REMOVED***,
		bpf.JumpIf***REMOVED***
			Cond:     bpf.JumpGreaterThan,
			Val:      0x00010202,
			SkipTrue: 1,
		***REMOVED***,
		bpf.RetConstant***REMOVED***
			Val: 0,
		***REMOVED***,
		bpf.RetConstant***REMOVED***
			Val: 12,
		***REMOVED***,
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
	if want, got := 4, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMJumpIfLessThan(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 4,
		***REMOVED***,
		bpf.JumpIf***REMOVED***
			Cond:     bpf.JumpLessThan,
			Val:      0xff010203,
			SkipTrue: 1,
		***REMOVED***,
		bpf.RetConstant***REMOVED***
			Val: 0,
		***REMOVED***,
		bpf.RetConstant***REMOVED***
			Val: 12,
		***REMOVED***,
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
	if want, got := 4, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMJumpIfGreaterOrEqual(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 4,
		***REMOVED***,
		bpf.JumpIf***REMOVED***
			Cond:     bpf.JumpGreaterOrEqual,
			Val:      0x00010203,
			SkipTrue: 1,
		***REMOVED***,
		bpf.RetConstant***REMOVED***
			Val: 0,
		***REMOVED***,
		bpf.RetConstant***REMOVED***
			Val: 12,
		***REMOVED***,
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
	if want, got := 4, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMJumpIfLessOrEqual(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 4,
		***REMOVED***,
		bpf.JumpIf***REMOVED***
			Cond:     bpf.JumpLessOrEqual,
			Val:      0xff010203,
			SkipTrue: 1,
		***REMOVED***,
		bpf.RetConstant***REMOVED***
			Val: 0,
		***REMOVED***,
		bpf.RetConstant***REMOVED***
			Val: 12,
		***REMOVED***,
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
	if want, got := 4, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMJumpIfBitsSet(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 2,
		***REMOVED***,
		bpf.JumpIf***REMOVED***
			Cond:     bpf.JumpBitsSet,
			Val:      0x1122,
			SkipTrue: 1,
		***REMOVED***,
		bpf.RetConstant***REMOVED***
			Val: 0,
		***REMOVED***,
		bpf.RetConstant***REMOVED***
			Val: 10,
		***REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("failed to load BPF program: %v", err)
	***REMOVED***
	defer done()

	out, err := vm.Run([]byte***REMOVED***
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0x01, 0x02,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 2, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMJumpIfBitsNotSet(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 2,
		***REMOVED***,
		bpf.JumpIf***REMOVED***
			Cond:     bpf.JumpBitsNotSet,
			Val:      0x1221,
			SkipTrue: 1,
		***REMOVED***,
		bpf.RetConstant***REMOVED***
			Val: 0,
		***REMOVED***,
		bpf.RetConstant***REMOVED***
			Val: 10,
		***REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("failed to load BPF program: %v", err)
	***REMOVED***
	defer done()

	out, err := vm.Run([]byte***REMOVED***
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0x01, 0x02,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 2, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***
