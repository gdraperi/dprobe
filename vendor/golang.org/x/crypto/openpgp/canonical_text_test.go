// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package openpgp

import (
	"bytes"
	"testing"
)

type recordingHash struct ***REMOVED***
	buf *bytes.Buffer
***REMOVED***

func (r recordingHash) Write(b []byte) (n int, err error) ***REMOVED***
	return r.buf.Write(b)
***REMOVED***

func (r recordingHash) Sum(in []byte) []byte ***REMOVED***
	return append(in, r.buf.Bytes()...)
***REMOVED***

func (r recordingHash) Reset() ***REMOVED***
	panic("shouldn't be called")
***REMOVED***

func (r recordingHash) Size() int ***REMOVED***
	panic("shouldn't be called")
***REMOVED***

func (r recordingHash) BlockSize() int ***REMOVED***
	panic("shouldn't be called")
***REMOVED***

func testCanonicalText(t *testing.T, input, expected string) ***REMOVED***
	r := recordingHash***REMOVED***bytes.NewBuffer(nil)***REMOVED***
	c := NewCanonicalTextHash(r)
	c.Write([]byte(input))
	result := c.Sum(nil)
	if expected != string(result) ***REMOVED***
		t.Errorf("input: %x got: %x want: %x", input, result, expected)
	***REMOVED***
***REMOVED***

func TestCanonicalText(t *testing.T) ***REMOVED***
	testCanonicalText(t, "foo\n", "foo\r\n")
	testCanonicalText(t, "foo", "foo")
	testCanonicalText(t, "foo\r\n", "foo\r\n")
	testCanonicalText(t, "foo\r\nbar", "foo\r\nbar")
	testCanonicalText(t, "foo\r\nbar\n\n", "foo\r\nbar\r\n\r\n")
***REMOVED***
