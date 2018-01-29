// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkcs12

import (
	"bytes"
	"encoding/hex"
	"testing"
)

var bmpStringTests = []struct ***REMOVED***
	in          string
	expectedHex string
	shouldFail  bool
***REMOVED******REMOVED***
	***REMOVED***"", "0000", false***REMOVED***,
	// Example from https://tools.ietf.org/html/rfc7292#appendix-B.
	***REMOVED***"Beavis", "0042006500610076006900730000", false***REMOVED***,
	// Some characters from the "Letterlike Symbols Unicode block".
	***REMOVED***"\u2115 - Double-struck N", "21150020002d00200044006f00750062006c0065002d00730074007200750063006b0020004e0000", false***REMOVED***,
	// any character outside the BMP should trigger an error.
	***REMOVED***"\U0001f000 East wind (Mahjong)", "", true***REMOVED***,
***REMOVED***

func TestBMPString(t *testing.T) ***REMOVED***
	for i, test := range bmpStringTests ***REMOVED***
		expected, err := hex.DecodeString(test.expectedHex)
		if err != nil ***REMOVED***
			t.Fatalf("#%d: failed to decode expectation", i)
		***REMOVED***

		out, err := bmpString(test.in)
		if err == nil && test.shouldFail ***REMOVED***
			t.Errorf("#%d: expected to fail, but produced %x", i, out)
			continue
		***REMOVED***

		if err != nil && !test.shouldFail ***REMOVED***
			t.Errorf("#%d: failed unexpectedly: %s", i, err)
			continue
		***REMOVED***

		if !test.shouldFail ***REMOVED***
			if !bytes.Equal(out, expected) ***REMOVED***
				t.Errorf("#%d: expected %s, got %x", i, test.expectedHex, out)
				continue
			***REMOVED***

			roundTrip, err := decodeBMPString(out)
			if err != nil ***REMOVED***
				t.Errorf("#%d: decoding output gave an error: %s", i, err)
				continue
			***REMOVED***

			if roundTrip != test.in ***REMOVED***
				t.Errorf("#%d: decoding output resulted in %q, but it should have been %q", i, roundTrip, test.in)
				continue
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
