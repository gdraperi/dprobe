// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"golang.org/x/crypto/openpgp/errors"
	"io"
	"io/ioutil"
	"testing"
)

// TestReader wraps a []byte and returns reads of a specific length.
type testReader struct ***REMOVED***
	data   []byte
	stride int
***REMOVED***

func (t *testReader) Read(buf []byte) (n int, err error) ***REMOVED***
	n = t.stride
	if n > len(t.data) ***REMOVED***
		n = len(t.data)
	***REMOVED***
	if n > len(buf) ***REMOVED***
		n = len(buf)
	***REMOVED***
	copy(buf, t.data)
	t.data = t.data[n:]
	if len(t.data) == 0 ***REMOVED***
		err = io.EOF
	***REMOVED***
	return
***REMOVED***

func testMDCReader(t *testing.T) ***REMOVED***
	mdcPlaintext, _ := hex.DecodeString(mdcPlaintextHex)

	for stride := 1; stride < len(mdcPlaintext)/2; stride++ ***REMOVED***
		r := &testReader***REMOVED***data: mdcPlaintext, stride: stride***REMOVED***
		mdcReader := &seMDCReader***REMOVED***in: r, h: sha1.New()***REMOVED***
		body, err := ioutil.ReadAll(mdcReader)
		if err != nil ***REMOVED***
			t.Errorf("stride: %d, error: %s", stride, err)
			continue
		***REMOVED***
		if !bytes.Equal(body, mdcPlaintext[:len(mdcPlaintext)-22]) ***REMOVED***
			t.Errorf("stride: %d: bad contents %x", stride, body)
			continue
		***REMOVED***

		err = mdcReader.Close()
		if err != nil ***REMOVED***
			t.Errorf("stride: %d, error on Close: %s", stride, err)
		***REMOVED***
	***REMOVED***

	mdcPlaintext[15] ^= 80

	r := &testReader***REMOVED***data: mdcPlaintext, stride: 2***REMOVED***
	mdcReader := &seMDCReader***REMOVED***in: r, h: sha1.New()***REMOVED***
	_, err := ioutil.ReadAll(mdcReader)
	if err != nil ***REMOVED***
		t.Errorf("corruption test, error: %s", err)
		return
	***REMOVED***
	err = mdcReader.Close()
	if err == nil ***REMOVED***
		t.Error("corruption: no error")
	***REMOVED*** else if _, ok := err.(*errors.SignatureError); !ok ***REMOVED***
		t.Errorf("corruption: expected SignatureError, got: %s", err)
	***REMOVED***
***REMOVED***

const mdcPlaintextHex = "a302789c3b2d93c4e0eb9aba22283539b3203335af44a134afb800c849cb4c4de10200aff40b45d31432c80cb384299a0655966d6939dfdeed1dddf980"

func TestSerialize(t *testing.T) ***REMOVED***
	buf := bytes.NewBuffer(nil)
	c := CipherAES128
	key := make([]byte, c.KeySize())

	w, err := SerializeSymmetricallyEncrypted(buf, c, key, nil)
	if err != nil ***REMOVED***
		t.Errorf("error from SerializeSymmetricallyEncrypted: %s", err)
		return
	***REMOVED***

	contents := []byte("hello world\n")

	w.Write(contents)
	w.Close()

	p, err := Read(buf)
	if err != nil ***REMOVED***
		t.Errorf("error from Read: %s", err)
		return
	***REMOVED***

	se, ok := p.(*SymmetricallyEncrypted)
	if !ok ***REMOVED***
		t.Errorf("didn't read a *SymmetricallyEncrypted")
		return
	***REMOVED***

	r, err := se.Decrypt(c, key)
	if err != nil ***REMOVED***
		t.Errorf("error from Decrypt: %s", err)
		return
	***REMOVED***

	contentsCopy := bytes.NewBuffer(nil)
	_, err = io.Copy(contentsCopy, r)
	if err != nil ***REMOVED***
		t.Errorf("error from io.Copy: %s", err)
		return
	***REMOVED***
	if !bytes.Equal(contentsCopy.Bytes(), contents) ***REMOVED***
		t.Errorf("contents not equal got: %x want: %x", contentsCopy.Bytes(), contents)
	***REMOVED***
***REMOVED***
