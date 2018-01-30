// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bpf_test

import (
	"fmt"
	"testing"

	"golang.org/x/net/bpf"
)

var _ bpf.Instruction = unknown***REMOVED******REMOVED***

type unknown struct***REMOVED******REMOVED***

func (unknown) Assemble() (bpf.RawInstruction, error) ***REMOVED***
	return bpf.RawInstruction***REMOVED******REMOVED***, nil
***REMOVED***

func TestVMUnknownInstruction(t *testing.T) ***REMOVED***
	vm, done, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadConstant***REMOVED***
			Dst: bpf.RegA,
			Val: 100,
		***REMOVED***,
		// Should terminate the program with an error immediately
		unknown***REMOVED******REMOVED***,
		bpf.RetA***REMOVED******REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error: %v", err)
	***REMOVED***
	defer done()

	_, err = vm.Run([]byte***REMOVED***
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0x00, 0x00,
	***REMOVED***)
	if errStr(err) != "unknown Instruction at index 1: bpf_test.unknown" ***REMOVED***
		t.Fatalf("unexpected error while running program: %v", err)
	***REMOVED***
***REMOVED***

func TestVMNoReturnInstruction(t *testing.T) ***REMOVED***
	_, _, err := testVM(t, []bpf.Instruction***REMOVED***
		bpf.LoadConstant***REMOVED***
			Dst: bpf.RegA,
			Val: 1,
		***REMOVED***,
	***REMOVED***)
	if errStr(err) != "BPF program must end with RetA or RetConstant" ***REMOVED***
		t.Fatalf("unexpected error: %v", err)
	***REMOVED***
***REMOVED***

func TestVMNoInputInstructions(t *testing.T) ***REMOVED***
	_, _, err := testVM(t, []bpf.Instruction***REMOVED******REMOVED***)
	if errStr(err) != "one or more Instructions must be specified" ***REMOVED***
		t.Fatalf("unexpected error: %v", err)
	***REMOVED***
***REMOVED***

// ExampleNewVM demonstrates usage of a VM, using an Ethernet frame
// as input and checking its EtherType to determine if it should be accepted.
func ExampleNewVM() ***REMOVED***
	// Offset | Length | Comment
	// -------------------------
	//   00   |   06   | Ethernet destination MAC address
	//   06   |   06   | Ethernet source MAC address
	//   12   |   02   | Ethernet EtherType
	const (
		etOff = 12
		etLen = 2

		etARP = 0x0806
	)

	// Set up a VM to filter traffic based on if its EtherType
	// matches the ARP EtherType.
	vm, err := bpf.NewVM([]bpf.Instruction***REMOVED***
		// Load EtherType value from Ethernet header
		bpf.LoadAbsolute***REMOVED***
			Off:  etOff,
			Size: etLen,
		***REMOVED***,
		// If EtherType is equal to the ARP EtherType, jump to allow
		// packet to be accepted
		bpf.JumpIf***REMOVED***
			Cond:     bpf.JumpEqual,
			Val:      etARP,
			SkipTrue: 1,
		***REMOVED***,
		// EtherType does not match the ARP EtherType
		bpf.RetConstant***REMOVED***
			Val: 0,
		***REMOVED***,
		// EtherType matches the ARP EtherType, accept up to 1500
		// bytes of packet
		bpf.RetConstant***REMOVED***
			Val: 1500,
		***REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		panic(fmt.Sprintf("failed to load BPF program: %v", err))
	***REMOVED***

	// Create an Ethernet frame with the ARP EtherType for testing
	frame := []byte***REMOVED***
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0x00, 0x11, 0x22, 0x33, 0x44, 0x55,
		0x08, 0x06,
		// Payload omitted for brevity
	***REMOVED***

	// Run our VM's BPF program using the Ethernet frame as input
	out, err := vm.Run(frame)
	if err != nil ***REMOVED***
		panic(fmt.Sprintf("failed to accept Ethernet frame: %v", err))
	***REMOVED***

	// BPF VM can return a byte count greater than the number of input
	// bytes, so trim the output to match the input byte length
	if out > len(frame) ***REMOVED***
		out = len(frame)
	***REMOVED***

	fmt.Printf("out: %d bytes", out)

	// Output:
	// out: 14 bytes
***REMOVED***

// errStr returns the string representation of an error, or
// "<nil>" if it is nil.
func errStr(err error) string ***REMOVED***
	if err == nil ***REMOVED***
		return "<nil>"
	***REMOVED***

	return err.Error()
***REMOVED***
