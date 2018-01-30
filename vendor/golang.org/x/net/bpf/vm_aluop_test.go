// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bpf_test

import (
	"testing"

	"golang.org/x/net/bpf"
)

func TestVMALUOpAdd(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 1,
		***REMOVED***,
		bpf.ALUOpConstant***REMOVED***
			Op:  bpf.ALUOpAdd,
			Val: 3,
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
		8, 2, 3,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 3, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMALUOpSub(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 1,
		***REMOVED***,
		bpf.TAX***REMOVED******REMOVED***,
		bpf.ALUOpX***REMOVED***
			Op: bpf.ALUOpSub,
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
		1, 2, 3,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 0, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMALUOpMul(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 1,
		***REMOVED***,
		bpf.ALUOpConstant***REMOVED***
			Op:  bpf.ALUOpMul,
			Val: 2,
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
		6, 2, 3, 4,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 4, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMALUOpDiv(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 1,
		***REMOVED***,
		bpf.ALUOpConstant***REMOVED***
			Op:  bpf.ALUOpDiv,
			Val: 2,
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
		20, 2, 3, 4,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 2, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMALUOpDivByZeroALUOpConstant(t *testing.T) ***REMOVED***
	_, _, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.ALUOpConstant***REMOVED***
			Op:  bpf.ALUOpDiv,
			Val: 0,
		***REMOVED***,
		bpf.RetA***REMOVED******REMOVED***,
	***REMOVED***)
	if errStr(err) != "cannot divide by zero using ALUOpConstant" ***REMOVED***
		t.Fatalf("unexpected error: %v", err)
	***REMOVED***
***REMOVED***

func TestVMALUOpDivByZeroALUOpX(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		// Load byte 0 into X
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 1,
		***REMOVED***,
		bpf.TAX***REMOVED******REMOVED***,
		// Load byte 1 into A
		bpf.LoadAbsolute***REMOVED***
			Off:  9,
			Size: 1,
		***REMOVED***,
		// Attempt to perform 1/0
		bpf.ALUOpX***REMOVED***
			Op: bpf.ALUOpDiv,
		***REMOVED***,
		// Return 4 bytes if program does not terminate
		bpf.LoadConstant***REMOVED***
			Val: 12,
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
		0, 1, 3, 4,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 0, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMALUOpOr(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 2,
		***REMOVED***,
		bpf.ALUOpConstant***REMOVED***
			Op:  bpf.ALUOpOr,
			Val: 0x01,
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
		0x00, 0x10, 0x03, 0x04,
		0x05, 0x06, 0x07, 0x08,
		0x09, 0xff,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 9, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMALUOpAnd(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 2,
		***REMOVED***,
		bpf.ALUOpConstant***REMOVED***
			Op:  bpf.ALUOpAnd,
			Val: 0x0019,
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
		0xaa, 0x09,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 1, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMALUOpShiftLeft(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 1,
		***REMOVED***,
		bpf.ALUOpConstant***REMOVED***
			Op:  bpf.ALUOpShiftLeft,
			Val: 0x01,
		***REMOVED***,
		bpf.JumpIf***REMOVED***
			Cond:     bpf.JumpEqual,
			Val:      0x02,
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
		0x01, 0xaa,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 1, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMALUOpShiftRight(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 1,
		***REMOVED***,
		bpf.ALUOpConstant***REMOVED***
			Op:  bpf.ALUOpShiftRight,
			Val: 0x01,
		***REMOVED***,
		bpf.JumpIf***REMOVED***
			Cond:     bpf.JumpEqual,
			Val:      0x04,
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
		0x08, 0xff, 0xff,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 1, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMALUOpMod(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 1,
		***REMOVED***,
		bpf.ALUOpConstant***REMOVED***
			Op:  bpf.ALUOpMod,
			Val: 20,
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
		30, 0, 0,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 2, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMALUOpModByZeroALUOpConstant(t *testing.T) ***REMOVED***
	_, _, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 1,
		***REMOVED***,
		bpf.ALUOpConstant***REMOVED***
			Op:  bpf.ALUOpMod,
			Val: 0,
		***REMOVED***,
		bpf.RetA***REMOVED******REMOVED***,
	***REMOVED***)
	if errStr(err) != "cannot divide by zero using ALUOpConstant" ***REMOVED***
		t.Fatalf("unexpected error: %v", err)
	***REMOVED***
***REMOVED***

func TestVMALUOpModByZeroALUOpX(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		// Load byte 0 into X
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 1,
		***REMOVED***,
		bpf.TAX***REMOVED******REMOVED***,
		// Load byte 1 into A
		bpf.LoadAbsolute***REMOVED***
			Off:  9,
			Size: 1,
		***REMOVED***,
		// Attempt to perform 1%0
		bpf.ALUOpX***REMOVED***
			Op: bpf.ALUOpMod,
		***REMOVED***,
		// Return 4 bytes if program does not terminate
		bpf.LoadConstant***REMOVED***
			Val: 12,
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
		0, 1, 3, 4,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 0, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMALUOpXor(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 1,
		***REMOVED***,
		bpf.ALUOpConstant***REMOVED***
			Op:  bpf.ALUOpXor,
			Val: 0x0a,
		***REMOVED***,
		bpf.JumpIf***REMOVED***
			Cond:     bpf.JumpEqual,
			Val:      0x01,
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
		0x0b, 0x00, 0x00, 0x00,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 1, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMALUOpUnknown(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 1,
		***REMOVED***,
		bpf.ALUOpConstant***REMOVED***
			Op:  bpf.ALUOpAdd,
			Val: 1,
		***REMOVED***,
		// Verify that an unknown operation is a no-op
		bpf.ALUOpConstant***REMOVED***
			Op: 100,
		***REMOVED***,
		bpf.JumpIf***REMOVED***
			Cond:     bpf.JumpEqual,
			Val:      0x02,
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
