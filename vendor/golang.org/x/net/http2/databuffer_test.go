// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.7

package http2

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func fmtDataChunk(chunk []byte) string ***REMOVED***
	out := ""
	var last byte
	var count int
	for _, c := range chunk ***REMOVED***
		if c != last ***REMOVED***
			if count > 0 ***REMOVED***
				out += fmt.Sprintf(" x %d ", count)
				count = 0
			***REMOVED***
			out += string([]byte***REMOVED***c***REMOVED***)
			last = c
		***REMOVED***
		count++
	***REMOVED***
	if count > 0 ***REMOVED***
		out += fmt.Sprintf(" x %d", count)
	***REMOVED***
	return out
***REMOVED***

func fmtDataChunks(chunks [][]byte) string ***REMOVED***
	var out string
	for _, chunk := range chunks ***REMOVED***
		out += fmt.Sprintf("***REMOVED***%q***REMOVED***", fmtDataChunk(chunk))
	***REMOVED***
	return out
***REMOVED***

func testDataBuffer(t *testing.T, wantBytes []byte, setup func(t *testing.T) *dataBuffer) ***REMOVED***
	// Run setup, then read the remaining bytes from the dataBuffer and check
	// that they match wantBytes. We use different read sizes to check corner
	// cases in Read.
	for _, readSize := range []int***REMOVED***1, 2, 1 * 1024, 32 * 1024***REMOVED*** ***REMOVED***
		t.Run(fmt.Sprintf("ReadSize=%d", readSize), func(t *testing.T) ***REMOVED***
			b := setup(t)
			buf := make([]byte, readSize)
			var gotRead bytes.Buffer
			for ***REMOVED***
				n, err := b.Read(buf)
				gotRead.Write(buf[:n])
				if err == errReadEmpty ***REMOVED***
					break
				***REMOVED***
				if err != nil ***REMOVED***
					t.Fatalf("error after %v bytes: %v", gotRead.Len(), err)
				***REMOVED***
			***REMOVED***
			if got, want := gotRead.Bytes(), wantBytes; !bytes.Equal(got, want) ***REMOVED***
				t.Errorf("FinalRead=%q, want %q", fmtDataChunk(got), fmtDataChunk(want))
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestDataBufferAllocation(t *testing.T) ***REMOVED***
	writes := [][]byte***REMOVED***
		bytes.Repeat([]byte("a"), 1*1024-1),
		[]byte("a"),
		bytes.Repeat([]byte("b"), 4*1024-1),
		[]byte("b"),
		bytes.Repeat([]byte("c"), 8*1024-1),
		[]byte("c"),
		bytes.Repeat([]byte("d"), 16*1024-1),
		[]byte("d"),
		bytes.Repeat([]byte("e"), 32*1024),
	***REMOVED***
	var wantRead bytes.Buffer
	for _, p := range writes ***REMOVED***
		wantRead.Write(p)
	***REMOVED***

	testDataBuffer(t, wantRead.Bytes(), func(t *testing.T) *dataBuffer ***REMOVED***
		b := &dataBuffer***REMOVED******REMOVED***
		for _, p := range writes ***REMOVED***
			if n, err := b.Write(p); n != len(p) || err != nil ***REMOVED***
				t.Fatalf("Write(%q x %d)=%v,%v want %v,nil", p[:1], len(p), n, err, len(p))
			***REMOVED***
		***REMOVED***
		want := [][]byte***REMOVED***
			bytes.Repeat([]byte("a"), 1*1024),
			bytes.Repeat([]byte("b"), 4*1024),
			bytes.Repeat([]byte("c"), 8*1024),
			bytes.Repeat([]byte("d"), 16*1024),
			bytes.Repeat([]byte("e"), 16*1024),
			bytes.Repeat([]byte("e"), 16*1024),
		***REMOVED***
		if !reflect.DeepEqual(b.chunks, want) ***REMOVED***
			t.Errorf("dataBuffer.chunks\ngot:  %s\nwant: %s", fmtDataChunks(b.chunks), fmtDataChunks(want))
		***REMOVED***
		return b
	***REMOVED***)
***REMOVED***

func TestDataBufferAllocationWithExpected(t *testing.T) ***REMOVED***
	writes := [][]byte***REMOVED***
		bytes.Repeat([]byte("a"), 1*1024), // allocates 16KB
		bytes.Repeat([]byte("b"), 14*1024),
		bytes.Repeat([]byte("c"), 15*1024), // allocates 16KB more
		bytes.Repeat([]byte("d"), 2*1024),
		bytes.Repeat([]byte("e"), 1*1024), // overflows 32KB expectation, allocates just 1KB
	***REMOVED***
	var wantRead bytes.Buffer
	for _, p := range writes ***REMOVED***
		wantRead.Write(p)
	***REMOVED***

	testDataBuffer(t, wantRead.Bytes(), func(t *testing.T) *dataBuffer ***REMOVED***
		b := &dataBuffer***REMOVED***expected: 32 * 1024***REMOVED***
		for _, p := range writes ***REMOVED***
			if n, err := b.Write(p); n != len(p) || err != nil ***REMOVED***
				t.Fatalf("Write(%q x %d)=%v,%v want %v,nil", p[:1], len(p), n, err, len(p))
			***REMOVED***
		***REMOVED***
		want := [][]byte***REMOVED***
			append(bytes.Repeat([]byte("a"), 1*1024), append(bytes.Repeat([]byte("b"), 14*1024), bytes.Repeat([]byte("c"), 1*1024)...)...),
			append(bytes.Repeat([]byte("c"), 14*1024), bytes.Repeat([]byte("d"), 2*1024)...),
			bytes.Repeat([]byte("e"), 1*1024),
		***REMOVED***
		if !reflect.DeepEqual(b.chunks, want) ***REMOVED***
			t.Errorf("dataBuffer.chunks\ngot:  %s\nwant: %s", fmtDataChunks(b.chunks), fmtDataChunks(want))
		***REMOVED***
		return b
	***REMOVED***)
***REMOVED***

func TestDataBufferWriteAfterPartialRead(t *testing.T) ***REMOVED***
	testDataBuffer(t, []byte("cdxyz"), func(t *testing.T) *dataBuffer ***REMOVED***
		b := &dataBuffer***REMOVED******REMOVED***
		if n, err := b.Write([]byte("abcd")); n != 4 || err != nil ***REMOVED***
			t.Fatalf("Write(\"abcd\")=%v,%v want 4,nil", n, err)
		***REMOVED***
		p := make([]byte, 2)
		if n, err := b.Read(p); n != 2 || err != nil || !bytes.Equal(p, []byte("ab")) ***REMOVED***
			t.Fatalf("Read()=%q,%v,%v want \"ab\",2,nil", p, n, err)
		***REMOVED***
		if n, err := b.Write([]byte("xyz")); n != 3 || err != nil ***REMOVED***
			t.Fatalf("Write(\"xyz\")=%v,%v want 3,nil", n, err)
		***REMOVED***
		return b
	***REMOVED***)
***REMOVED***
