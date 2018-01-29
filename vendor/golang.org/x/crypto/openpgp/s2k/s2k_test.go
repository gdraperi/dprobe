// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package s2k

import (
	"bytes"
	"crypto"
	_ "crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"
	"encoding/hex"
	"testing"

	_ "golang.org/x/crypto/ripemd160"
)

var saltedTests = []struct ***REMOVED***
	in, out string
***REMOVED******REMOVED***
	***REMOVED***"hello", "10295ac1"***REMOVED***,
	***REMOVED***"world", "ac587a5e"***REMOVED***,
	***REMOVED***"foo", "4dda8077"***REMOVED***,
	***REMOVED***"bar", "bd8aac6b9ea9cae04eae6a91c6133b58b5d9a61c14f355516ed9370456"***REMOVED***,
	***REMOVED***"x", "f1d3f289"***REMOVED***,
	***REMOVED***"xxxxxxxxxxxxxxxxxxxxxxx", "e00d7b45"***REMOVED***,
***REMOVED***

func TestSalted(t *testing.T) ***REMOVED***
	h := sha1.New()
	salt := [4]byte***REMOVED***1, 2, 3, 4***REMOVED***

	for i, test := range saltedTests ***REMOVED***
		expected, _ := hex.DecodeString(test.out)
		out := make([]byte, len(expected))
		Salted(out, h, []byte(test.in), salt[:])
		if !bytes.Equal(expected, out) ***REMOVED***
			t.Errorf("#%d, got: %x want: %x", i, out, expected)
		***REMOVED***
	***REMOVED***
***REMOVED***

var iteratedTests = []struct ***REMOVED***
	in, out string
***REMOVED******REMOVED***
	***REMOVED***"hello", "83126105"***REMOVED***,
	***REMOVED***"world", "6fa317f9"***REMOVED***,
	***REMOVED***"foo", "8fbc35b9"***REMOVED***,
	***REMOVED***"bar", "2af5a99b54f093789fd657f19bd245af7604d0f6ae06f66602a46a08ae"***REMOVED***,
	***REMOVED***"x", "5a684dfe"***REMOVED***,
	***REMOVED***"xxxxxxxxxxxxxxxxxxxxxxx", "18955174"***REMOVED***,
***REMOVED***

func TestIterated(t *testing.T) ***REMOVED***
	h := sha1.New()
	salt := [4]byte***REMOVED***4, 3, 2, 1***REMOVED***

	for i, test := range iteratedTests ***REMOVED***
		expected, _ := hex.DecodeString(test.out)
		out := make([]byte, len(expected))
		Iterated(out, h, []byte(test.in), salt[:], 31)
		if !bytes.Equal(expected, out) ***REMOVED***
			t.Errorf("#%d, got: %x want: %x", i, out, expected)
		***REMOVED***
	***REMOVED***
***REMOVED***

var parseTests = []struct ***REMOVED***
	spec, in, out string
***REMOVED******REMOVED***
	/* Simple with SHA1 */
	***REMOVED***"0002", "hello", "aaf4c61d"***REMOVED***,
	/* Salted with SHA1 */
	***REMOVED***"01020102030405060708", "hello", "f4f7d67e"***REMOVED***,
	/* Iterated with SHA1 */
	***REMOVED***"03020102030405060708f1", "hello", "f2a57b7c"***REMOVED***,
***REMOVED***

func TestParse(t *testing.T) ***REMOVED***
	for i, test := range parseTests ***REMOVED***
		spec, _ := hex.DecodeString(test.spec)
		buf := bytes.NewBuffer(spec)
		f, err := Parse(buf)
		if err != nil ***REMOVED***
			t.Errorf("%d: Parse returned error: %s", i, err)
			continue
		***REMOVED***

		expected, _ := hex.DecodeString(test.out)
		out := make([]byte, len(expected))
		f(out, []byte(test.in))
		if !bytes.Equal(out, expected) ***REMOVED***
			t.Errorf("%d: output got: %x want: %x", i, out, expected)
		***REMOVED***
		if testing.Short() ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSerialize(t *testing.T) ***REMOVED***
	hashes := []crypto.Hash***REMOVED***crypto.MD5, crypto.SHA1, crypto.RIPEMD160,
		crypto.SHA256, crypto.SHA384, crypto.SHA512, crypto.SHA224***REMOVED***
	testCounts := []int***REMOVED***-1, 0, 1024, 65536, 4063232, 65011712***REMOVED***
	for _, h := range hashes ***REMOVED***
		for _, c := range testCounts ***REMOVED***
			testSerializeConfig(t, &Config***REMOVED***Hash: h, S2KCount: c***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***

func testSerializeConfig(t *testing.T, c *Config) ***REMOVED***
	t.Logf("Running testSerializeConfig() with config: %+v", c)

	buf := bytes.NewBuffer(nil)
	key := make([]byte, 16)
	passphrase := []byte("testing")
	err := Serialize(buf, key, rand.Reader, passphrase, c)
	if err != nil ***REMOVED***
		t.Errorf("failed to serialize: %s", err)
		return
	***REMOVED***

	f, err := Parse(buf)
	if err != nil ***REMOVED***
		t.Errorf("failed to reparse: %s", err)
		return
	***REMOVED***
	key2 := make([]byte, len(key))
	f(key2, passphrase)
	if !bytes.Equal(key2, key) ***REMOVED***
		t.Errorf("keys don't match: %x (serialied) vs %x (parsed)", key, key2)
	***REMOVED***
***REMOVED***
