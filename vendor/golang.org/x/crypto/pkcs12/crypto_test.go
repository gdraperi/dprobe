// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkcs12

import (
	"bytes"
	"crypto/x509/pkix"
	"encoding/asn1"
	"testing"
)

var sha1WithTripleDES = asn1.ObjectIdentifier([]int***REMOVED***1, 2, 840, 113549, 1, 12, 1, 3***REMOVED***)

func TestPbDecrypterFor(t *testing.T) ***REMOVED***
	params, _ := asn1.Marshal(pbeParams***REMOVED***
		Salt:       []byte***REMOVED***1, 2, 3, 4, 5, 6, 7, 8***REMOVED***,
		Iterations: 2048,
	***REMOVED***)
	alg := pkix.AlgorithmIdentifier***REMOVED***
		Algorithm: asn1.ObjectIdentifier([]int***REMOVED***1, 2, 3***REMOVED***),
		Parameters: asn1.RawValue***REMOVED***
			FullBytes: params,
		***REMOVED***,
	***REMOVED***

	pass, _ := bmpString("Sesame open")

	_, _, err := pbDecrypterFor(alg, pass)
	if _, ok := err.(NotImplementedError); !ok ***REMOVED***
		t.Errorf("expected not implemented error, got: %T %s", err, err)
	***REMOVED***

	alg.Algorithm = sha1WithTripleDES
	cbc, blockSize, err := pbDecrypterFor(alg, pass)
	if err != nil ***REMOVED***
		t.Errorf("unexpected error from pbDecrypterFor %v", err)
	***REMOVED***
	if blockSize != 8 ***REMOVED***
		t.Errorf("unexpected block size %d, wanted 8", blockSize)
	***REMOVED***

	plaintext := []byte***REMOVED***1, 2, 3, 4, 5, 6, 7, 8***REMOVED***
	expectedCiphertext := []byte***REMOVED***185, 73, 135, 249, 137, 1, 122, 247***REMOVED***
	ciphertext := make([]byte, len(plaintext))
	cbc.CryptBlocks(ciphertext, plaintext)

	if bytes.Compare(ciphertext, expectedCiphertext) != 0 ***REMOVED***
		t.Errorf("bad ciphertext, got %x but wanted %x", ciphertext, expectedCiphertext)
	***REMOVED***
***REMOVED***

var pbDecryptTests = []struct ***REMOVED***
	in            []byte
	expected      []byte
	expectedError error
***REMOVED******REMOVED***
	***REMOVED***
		[]byte("\x33\x73\xf3\x9f\xda\x49\xae\xfc\xa0\x9a\xdf\x5a\x58\xa0\xea\x46"), // 7 padding bytes
		[]byte("A secret!"),
		nil,
	***REMOVED***,
	***REMOVED***
		[]byte("\x33\x73\xf3\x9f\xda\x49\xae\xfc\x96\x24\x2f\x71\x7e\x32\x3f\xe7"), // 8 padding bytes
		[]byte("A secret"),
		nil,
	***REMOVED***,
	***REMOVED***
		[]byte("\x35\x0c\xc0\x8d\xab\xa9\x5d\x30\x7f\x9a\xec\x6a\xd8\x9b\x9c\xd9"), // 9 padding bytes, incorrect
		nil,
		ErrDecryption,
	***REMOVED***,
	***REMOVED***
		[]byte("\xb2\xf9\x6e\x06\x60\xae\x20\xcf\x08\xa0\x7b\xd9\x6b\x20\xef\x41"), // incorrect padding bytes: [ ... 0x04 0x02 ]
		nil,
		ErrDecryption,
	***REMOVED***,
***REMOVED***

func TestPbDecrypt(t *testing.T) ***REMOVED***
	for i, test := range pbDecryptTests ***REMOVED***
		decryptable := testDecryptable***REMOVED***
			data: test.in,
			algorithm: pkix.AlgorithmIdentifier***REMOVED***
				Algorithm: sha1WithTripleDES,
				Parameters: pbeParams***REMOVED***
					Salt:       []byte("\xf1\xf2\xf3\xf4\xf5\xf6\xf7\xf8"),
					Iterations: 4096,
				***REMOVED***.RawASN1(),
			***REMOVED***,
		***REMOVED***
		password, _ := bmpString("sesame")

		plaintext, err := pbDecrypt(decryptable, password)
		if err != test.expectedError ***REMOVED***
			t.Errorf("#%d: got error %q, but wanted %q", i, err, test.expectedError)
			continue
		***REMOVED***

		if !bytes.Equal(plaintext, test.expected) ***REMOVED***
			t.Errorf("#%d: got %x, but wanted %x", i, plaintext, test.expected)
		***REMOVED***
	***REMOVED***
***REMOVED***

type testDecryptable struct ***REMOVED***
	data      []byte
	algorithm pkix.AlgorithmIdentifier
***REMOVED***

func (d testDecryptable) Algorithm() pkix.AlgorithmIdentifier ***REMOVED*** return d.algorithm ***REMOVED***
func (d testDecryptable) Data() []byte                        ***REMOVED*** return d.data ***REMOVED***

func (params pbeParams) RawASN1() (raw asn1.RawValue) ***REMOVED***
	asn1Bytes, err := asn1.Marshal(params)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	_, err = asn1.Unmarshal(asn1Bytes, &raw)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return
***REMOVED***
