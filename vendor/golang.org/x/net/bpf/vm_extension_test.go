// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bpf_test

import (
	"testing"

	"golang.org/x/net/bpf"
)

func TestVMLoadExtensionNotImplemented(t *testing.T) ***REMOVED***
	_, _, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadExtension***REMOVED***
			Num: 100,
		***REMOVED***,
		bpf.RetA***REMOVED******REMOVED***,
	***REMOVED***)
	if errStr(err) != "extension 100 not implemented" ***REMOVED***
		t.Fatalf("unexpected error: %v", err)
	***REMOVED***
***REMOVED***

func TestVMLoadExtensionExtLen(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadExtension***REMOVED***
			Num: bpf.ExtLen,
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
	if want, got := 4, out; want != got ***REMOVED***
		t.Fatalf("unexpected number of output bytes:\n- want: %d\n-  got: %d",
			want, got)
	***REMOVED***
***REMOVED***
