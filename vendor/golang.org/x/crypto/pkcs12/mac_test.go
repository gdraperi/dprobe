// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkcs12

import (
	"encoding/asn1"
	"testing"
)

func TestVerifyMac(t *testing.T) ***REMOVED***
	td := macData***REMOVED***
		Mac: digestInfo***REMOVED***
			Digest: []byte***REMOVED***0x18, 0x20, 0x3d, 0xff, 0x1e, 0x16, 0xf4, 0x92, 0xf2, 0xaf, 0xc8, 0x91, 0xa9, 0xba, 0xd6, 0xca, 0x9d, 0xee, 0x51, 0x93***REMOVED***,
		***REMOVED***,
		MacSalt:    []byte***REMOVED***1, 2, 3, 4, 5, 6, 7, 8***REMOVED***,
		Iterations: 2048,
	***REMOVED***

	message := []byte***REMOVED***11, 12, 13, 14, 15***REMOVED***
	password, _ := bmpString("")

	td.Mac.Algorithm.Algorithm = asn1.ObjectIdentifier([]int***REMOVED***1, 2, 3***REMOVED***)
	err := verifyMac(&td, message, password)
	if _, ok := err.(NotImplementedError); !ok ***REMOVED***
		t.Errorf("err: %v", err)
	***REMOVED***

	td.Mac.Algorithm.Algorithm = asn1.ObjectIdentifier([]int***REMOVED***1, 3, 14, 3, 2, 26***REMOVED***)
	err = verifyMac(&td, message, password)
	if err != ErrIncorrectPassword ***REMOVED***
		t.Errorf("Expected incorrect password, got err: %v", err)
	***REMOVED***

	password, _ = bmpString("Sesame open")
	err = verifyMac(&td, message, password)
	if err != nil ***REMOVED***
		t.Errorf("err: %v", err)
	***REMOVED***

***REMOVED***
