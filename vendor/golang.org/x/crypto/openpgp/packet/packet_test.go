// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/openpgp/errors"
	"io"
	"io/ioutil"
	"testing"
)

func TestReadFull(t *testing.T) ***REMOVED***
	var out [4]byte

	b := bytes.NewBufferString("foo")
	n, err := readFull(b, out[:3])
	if n != 3 || err != nil ***REMOVED***
		t.Errorf("full read failed n:%d err:%s", n, err)
	***REMOVED***

	b = bytes.NewBufferString("foo")
	n, err = readFull(b, out[:4])
	if n != 3 || err != io.ErrUnexpectedEOF ***REMOVED***
		t.Errorf("partial read failed n:%d err:%s", n, err)
	***REMOVED***

	b = bytes.NewBuffer(nil)
	n, err = readFull(b, out[:3])
	if n != 0 || err != io.ErrUnexpectedEOF ***REMOVED***
		t.Errorf("empty read failed n:%d err:%s", n, err)
	***REMOVED***
***REMOVED***

func readerFromHex(s string) io.Reader ***REMOVED***
	data, err := hex.DecodeString(s)
	if err != nil ***REMOVED***
		panic("readerFromHex: bad input")
	***REMOVED***
	return bytes.NewBuffer(data)
***REMOVED***

var readLengthTests = []struct ***REMOVED***
	hexInput  string
	length    int64
	isPartial bool
	err       error
***REMOVED******REMOVED***
	***REMOVED***"", 0, false, io.ErrUnexpectedEOF***REMOVED***,
	***REMOVED***"1f", 31, false, nil***REMOVED***,
	***REMOVED***"c0", 0, false, io.ErrUnexpectedEOF***REMOVED***,
	***REMOVED***"c101", 256 + 1 + 192, false, nil***REMOVED***,
	***REMOVED***"e0", 1, true, nil***REMOVED***,
	***REMOVED***"e1", 2, true, nil***REMOVED***,
	***REMOVED***"e2", 4, true, nil***REMOVED***,
	***REMOVED***"ff", 0, false, io.ErrUnexpectedEOF***REMOVED***,
	***REMOVED***"ff00", 0, false, io.ErrUnexpectedEOF***REMOVED***,
	***REMOVED***"ff0000", 0, false, io.ErrUnexpectedEOF***REMOVED***,
	***REMOVED***"ff000000", 0, false, io.ErrUnexpectedEOF***REMOVED***,
	***REMOVED***"ff00000000", 0, false, nil***REMOVED***,
	***REMOVED***"ff01020304", 16909060, false, nil***REMOVED***,
***REMOVED***

func TestReadLength(t *testing.T) ***REMOVED***
	for i, test := range readLengthTests ***REMOVED***
		length, isPartial, err := readLength(readerFromHex(test.hexInput))
		if test.err != nil ***REMOVED***
			if err != test.err ***REMOVED***
				t.Errorf("%d: expected different error got:%s want:%s", i, err, test.err)
			***REMOVED***
			continue
		***REMOVED***
		if err != nil ***REMOVED***
			t.Errorf("%d: unexpected error: %s", i, err)
			continue
		***REMOVED***
		if length != test.length || isPartial != test.isPartial ***REMOVED***
			t.Errorf("%d: bad result got:(%d,%t) want:(%d,%t)", i, length, isPartial, test.length, test.isPartial)
		***REMOVED***
	***REMOVED***
***REMOVED***

var partialLengthReaderTests = []struct ***REMOVED***
	hexInput  string
	err       error
	hexOutput string
***REMOVED******REMOVED***
	***REMOVED***"e0", io.ErrUnexpectedEOF, ""***REMOVED***,
	***REMOVED***"e001", io.ErrUnexpectedEOF, ""***REMOVED***,
	***REMOVED***"e0010102", nil, "0102"***REMOVED***,
	***REMOVED***"ff00000000", nil, ""***REMOVED***,
	***REMOVED***"e10102e1030400", nil, "01020304"***REMOVED***,
	***REMOVED***"e101", io.ErrUnexpectedEOF, ""***REMOVED***,
***REMOVED***

func TestPartialLengthReader(t *testing.T) ***REMOVED***
	for i, test := range partialLengthReaderTests ***REMOVED***
		r := &partialLengthReader***REMOVED***readerFromHex(test.hexInput), 0, true***REMOVED***
		out, err := ioutil.ReadAll(r)
		if test.err != nil ***REMOVED***
			if err != test.err ***REMOVED***
				t.Errorf("%d: expected different error got:%s want:%s", i, err, test.err)
			***REMOVED***
			continue
		***REMOVED***
		if err != nil ***REMOVED***
			t.Errorf("%d: unexpected error: %s", i, err)
			continue
		***REMOVED***

		got := fmt.Sprintf("%x", out)
		if got != test.hexOutput ***REMOVED***
			t.Errorf("%d: got:%s want:%s", i, test.hexOutput, got)
		***REMOVED***
	***REMOVED***
***REMOVED***

var readHeaderTests = []struct ***REMOVED***
	hexInput        string
	structuralError bool
	unexpectedEOF   bool
	tag             int
	length          int64
	hexOutput       string
***REMOVED******REMOVED***
	***REMOVED***"", false, false, 0, 0, ""***REMOVED***,
	***REMOVED***"7f", true, false, 0, 0, ""***REMOVED***,

	// Old format headers
	***REMOVED***"80", false, true, 0, 0, ""***REMOVED***,
	***REMOVED***"8001", false, true, 0, 1, ""***REMOVED***,
	***REMOVED***"800102", false, false, 0, 1, "02"***REMOVED***,
	***REMOVED***"81000102", false, false, 0, 1, "02"***REMOVED***,
	***REMOVED***"820000000102", false, false, 0, 1, "02"***REMOVED***,
	***REMOVED***"860000000102", false, false, 1, 1, "02"***REMOVED***,
	***REMOVED***"83010203", false, false, 0, -1, "010203"***REMOVED***,

	// New format headers
	***REMOVED***"c0", false, true, 0, 0, ""***REMOVED***,
	***REMOVED***"c000", false, false, 0, 0, ""***REMOVED***,
	***REMOVED***"c00102", false, false, 0, 1, "02"***REMOVED***,
	***REMOVED***"c0020203", false, false, 0, 2, "0203"***REMOVED***,
	***REMOVED***"c00202", false, true, 0, 2, ""***REMOVED***,
	***REMOVED***"c3020203", false, false, 3, 2, "0203"***REMOVED***,
***REMOVED***

func TestReadHeader(t *testing.T) ***REMOVED***
	for i, test := range readHeaderTests ***REMOVED***
		tag, length, contents, err := readHeader(readerFromHex(test.hexInput))
		if test.structuralError ***REMOVED***
			if _, ok := err.(errors.StructuralError); ok ***REMOVED***
				continue
			***REMOVED***
			t.Errorf("%d: expected StructuralError, got:%s", i, err)
			continue
		***REMOVED***
		if err != nil ***REMOVED***
			if len(test.hexInput) == 0 && err == io.EOF ***REMOVED***
				continue
			***REMOVED***
			if !test.unexpectedEOF || err != io.ErrUnexpectedEOF ***REMOVED***
				t.Errorf("%d: unexpected error from readHeader: %s", i, err)
			***REMOVED***
			continue
		***REMOVED***
		if int(tag) != test.tag || length != test.length ***REMOVED***
			t.Errorf("%d: got:(%d,%d) want:(%d,%d)", i, int(tag), length, test.tag, test.length)
			continue
		***REMOVED***

		body, err := ioutil.ReadAll(contents)
		if err != nil ***REMOVED***
			if !test.unexpectedEOF || err != io.ErrUnexpectedEOF ***REMOVED***
				t.Errorf("%d: unexpected error from contents: %s", i, err)
			***REMOVED***
			continue
		***REMOVED***
		if test.unexpectedEOF ***REMOVED***
			t.Errorf("%d: expected ErrUnexpectedEOF from contents but got no error", i)
			continue
		***REMOVED***
		got := fmt.Sprintf("%x", body)
		if got != test.hexOutput ***REMOVED***
			t.Errorf("%d: got:%s want:%s", i, got, test.hexOutput)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSerializeHeader(t *testing.T) ***REMOVED***
	tag := packetTypePublicKey
	lengths := []int***REMOVED***0, 1, 2, 64, 192, 193, 8000, 8384, 8385, 10000***REMOVED***

	for _, length := range lengths ***REMOVED***
		buf := bytes.NewBuffer(nil)
		serializeHeader(buf, tag, length)
		tag2, length2, _, err := readHeader(buf)
		if err != nil ***REMOVED***
			t.Errorf("length %d, err: %s", length, err)
		***REMOVED***
		if tag2 != tag ***REMOVED***
			t.Errorf("length %d, tag incorrect (got %d, want %d)", length, tag2, tag)
		***REMOVED***
		if int(length2) != length ***REMOVED***
			t.Errorf("length %d, length incorrect (got %d)", length, length2)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestPartialLengths(t *testing.T) ***REMOVED***
	buf := bytes.NewBuffer(nil)
	w := new(partialLengthWriter)
	w.w = noOpCloser***REMOVED***buf***REMOVED***

	const maxChunkSize = 64

	var b [maxChunkSize]byte
	var n uint8
	for l := 1; l <= maxChunkSize; l++ ***REMOVED***
		for i := 0; i < l; i++ ***REMOVED***
			b[i] = n
			n++
		***REMOVED***
		m, err := w.Write(b[:l])
		if m != l ***REMOVED***
			t.Errorf("short write got: %d want: %d", m, l)
		***REMOVED***
		if err != nil ***REMOVED***
			t.Errorf("error from write: %s", err)
		***REMOVED***
	***REMOVED***
	w.Close()

	want := (maxChunkSize * (maxChunkSize + 1)) / 2
	copyBuf := bytes.NewBuffer(nil)
	r := &partialLengthReader***REMOVED***buf, 0, true***REMOVED***
	m, err := io.Copy(copyBuf, r)
	if m != int64(want) ***REMOVED***
		t.Errorf("short copy got: %d want: %d", m, want)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("error from copy: %s", err)
	***REMOVED***

	copyBytes := copyBuf.Bytes()
	for i := 0; i < want; i++ ***REMOVED***
		if copyBytes[i] != uint8(i) ***REMOVED***
			t.Errorf("bad pattern in copy at %d", i)
			break
		***REMOVED***
	***REMOVED***
***REMOVED***
