// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bpf_test

import (
	"testing"

	"golang.org/x/net/bpf"
)

func TestVMStoreScratchInvalidScratchRegisterTooSmall(t *testing.T) ***REMOVED***
	_, _, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.StoreScratch***REMOVED***
			Src: bpf.RegA,
			N:   -1,
		***REMOVED***,
		bpf.RetA***REMOVED******REMOVED***,
	***REMOVED***)
	if errStr(err) != "assembling instruction 1: invalid scratch slot -1" ***REMOVED***
		t.Fatalf("unexpected error: %v", err)
	***REMOVED***
***REMOVED***

func TestVMStoreScratchInvalidScratchRegisterTooLarge(t *testing.T) ***REMOVED***
	_, _, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.StoreScratch***REMOVED***
			Src: bpf.RegA,
			N:   16,
		***REMOVED***,
		bpf.RetA***REMOVED******REMOVED***,
	***REMOVED***)
	if errStr(err) != "assembling instruction 1: invalid scratch slot 16" ***REMOVED***
		t.Fatalf("unexpected error: %v", err)
	***REMOVED***
***REMOVED***

func TestVMStoreScratchUnknownSourceRegister(t *testing.T) ***REMOVED***
	_, _, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.StoreScratch***REMOVED***
			Src: 100,
			N:   0,
		***REMOVED***,
		bpf.RetA***REMOVED******REMOVED***,
	***REMOVED***)
	if errStr(err) != "assembling instruction 1: invalid source register 100" ***REMOVED***
		t.Fatalf("unexpected error: %v", err)
	***REMOVED***
***REMOVED***

func TestVMLoadScratchInvalidScratchRegisterTooSmall(t *testing.T) ***REMOVED***
	_, _, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadScratch***REMOVED***
			Dst: bpf.RegX,
			N:   -1,
		***REMOVED***,
		bpf.RetA***REMOVED******REMOVED***,
	***REMOVED***)
	if errStr(err) != "assembling instruction 1: invalid scratch slot -1" ***REMOVED***
		t.Fatalf("unexpected error: %v", err)
	***REMOVED***
***REMOVED***

func TestVMLoadScratchInvalidScratchRegisterTooLarge(t *testing.T) ***REMOVED***
	_, _, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadScratch***REMOVED***
			Dst: bpf.RegX,
			N:   16,
		***REMOVED***,
		bpf.RetA***REMOVED******REMOVED***,
	***REMOVED***)
	if errStr(err) != "assembling instruction 1: invalid scratch slot 16" ***REMOVED***
		t.Fatalf("unexpected error: %v", err)
	***REMOVED***
***REMOVED***

func TestVMLoadScratchUnknownDestinationRegister(t *testing.T) ***REMOVED***
	_, _, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadScratch***REMOVED***
			Dst: 100,
			N:   0,
		***REMOVED***,
		bpf.RetA***REMOVED******REMOVED***,
	***REMOVED***)
	if errStr(err) != "assembling instruction 1: invalid target register 100" ***REMOVED***
		t.Fatalf("unexpected error: %v", err)
	***REMOVED***
***REMOVED***

func TestVMStoreScratchLoadScratchOneValue(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		// Load byte 255
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 1,
		***REMOVED***,
		// Copy to X and store in scratch[0]
		bpf.TAX***REMOVED******REMOVED***,
		bpf.StoreScratch***REMOVED***
			Src: bpf.RegX,
			N:   0,
		***REMOVED***,
		// Load byte 1
		bpf.LoadAbsolute***REMOVED***
			Off:  9,
			Size: 1,
		***REMOVED***,
		// Overwrite 1 with 255 from scratch[0]
		bpf.LoadScratch***REMOVED***
			Dst: bpf.RegA,
			N:   0,
		***REMOVED***,
		// Return 255
		bpf.RetA***REMOVED******REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("failed to load BPF program: %v", err)
	***REMOVED***
	defer done()

	out, err := vm.Run([]byte***REMOVED***
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		255, 1, 2,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 3, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***

func TestVMStoreScratchLoadScratchMultipleValues(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		// Load byte 10
		bpf.LoadAbsolute***REMOVED***
			Off:  8,
			Size: 1,
		***REMOVED***,
		// Store in scratch[0]
		bpf.StoreScratch***REMOVED***
			Src: bpf.RegA,
			N:   0,
		***REMOVED***,
		// Load byte 20
		bpf.LoadAbsolute***REMOVED***
			Off:  9,
			Size: 1,
		***REMOVED***,
		// Store in scratch[1]
		bpf.StoreScratch***REMOVED***
			Src: bpf.RegA,
			N:   1,
		***REMOVED***,
		// Load byte 30
		bpf.LoadAbsolute***REMOVED***
			Off:  10,
			Size: 1,
		***REMOVED***,
		// Store in scratch[2]
		bpf.StoreScratch***REMOVED***
			Src: bpf.RegA,
			N:   2,
		***REMOVED***,
		// Load byte 1
		bpf.LoadAbsolute***REMOVED***
			Off:  11,
			Size: 1,
		***REMOVED***,
		// Store in scratch[3]
		bpf.StoreScratch***REMOVED***
			Src: bpf.RegA,
			N:   3,
		***REMOVED***,
		// Load in byte 10 to X
		bpf.LoadScratch***REMOVED***
			Dst: bpf.RegX,
			N:   0,
		***REMOVED***,
		// Copy X -> A
		bpf.TXA***REMOVED******REMOVED***,
		// Verify value is 10
		bpf.JumpIf***REMOVED***
			Cond:     bpf.JumpEqual,
			Val:      10,
			SkipTrue: 1,
		***REMOVED***,
		// Fail test if incorrect
		bpf.RetConstant***REMOVED***
			Val: 0,
		***REMOVED***,
		// Load in byte 20 to A
		bpf.LoadScratch***REMOVED***
			Dst: bpf.RegA,
			N:   1,
		***REMOVED***,
		// Verify value is 20
		bpf.JumpIf***REMOVED***
			Cond:     bpf.JumpEqual,
			Val:      20,
			SkipTrue: 1,
		***REMOVED***,
		// Fail test if incorrect
		bpf.RetConstant***REMOVED***
			Val: 0,
		***REMOVED***,
		// Load in byte 30 to A
		bpf.LoadScratch***REMOVED***
			Dst: bpf.RegA,
			N:   2,
		***REMOVED***,
		// Verify value is 30
		bpf.JumpIf***REMOVED***
			Cond:     bpf.JumpEqual,
			Val:      30,
			SkipTrue: 1,
		***REMOVED***,
		// Fail test if incorrect
		bpf.RetConstant***REMOVED***
			Val: 0,
		***REMOVED***,
		// Return first two bytes on success
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
		10, 20, 30, 1,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
	if want, got := 2, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***
