// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"bytes"
	"encoding/hex"
	"testing"
	"time"
)

var pubKeyV3Test = struct ***REMOVED***
	hexFingerprint string
	creationTime   time.Time
	pubKeyAlgo     PublicKeyAlgorithm
	keyId          uint64
	keyIdString    string
	keyIdShort     string
***REMOVED******REMOVED***
	"103BECF5BD1E837C89D19E98487767F7",
	time.Unix(779753634, 0),
	PubKeyAlgoRSA,
	0xDE0F188A5DA5E3C9,
	"DE0F188A5DA5E3C9",
	"5DA5E3C9"***REMOVED***

func TestPublicKeyV3Read(t *testing.T) ***REMOVED***
	i, test := 0, pubKeyV3Test
	packet, err := Read(v3KeyReader(t))
	if err != nil ***REMOVED***
		t.Fatalf("#%d: Read error: %s", i, err)
	***REMOVED***
	pk, ok := packet.(*PublicKeyV3)
	if !ok ***REMOVED***
		t.Fatalf("#%d: failed to parse, got: %#v", i, packet)
	***REMOVED***
	if pk.PubKeyAlgo != test.pubKeyAlgo ***REMOVED***
		t.Errorf("#%d: bad public key algorithm got:%x want:%x", i, pk.PubKeyAlgo, test.pubKeyAlgo)
	***REMOVED***
	if !pk.CreationTime.Equal(test.creationTime) ***REMOVED***
		t.Errorf("#%d: bad creation time got:%v want:%v", i, pk.CreationTime, test.creationTime)
	***REMOVED***
	expectedFingerprint, _ := hex.DecodeString(test.hexFingerprint)
	if !bytes.Equal(expectedFingerprint, pk.Fingerprint[:]) ***REMOVED***
		t.Errorf("#%d: bad fingerprint got:%x want:%x", i, pk.Fingerprint[:], expectedFingerprint)
	***REMOVED***
	if pk.KeyId != test.keyId ***REMOVED***
		t.Errorf("#%d: bad keyid got:%x want:%x", i, pk.KeyId, test.keyId)
	***REMOVED***
	if g, e := pk.KeyIdString(), test.keyIdString; g != e ***REMOVED***
		t.Errorf("#%d: bad KeyIdString got:%q want:%q", i, g, e)
	***REMOVED***
	if g, e := pk.KeyIdShortString(), test.keyIdShort; g != e ***REMOVED***
		t.Errorf("#%d: bad KeyIdShortString got:%q want:%q", i, g, e)
	***REMOVED***
***REMOVED***

func TestPublicKeyV3Serialize(t *testing.T) ***REMOVED***
	//for i, test := range pubKeyV3Tests ***REMOVED***
	i := 0
	packet, err := Read(v3KeyReader(t))
	if err != nil ***REMOVED***
		t.Fatalf("#%d: Read error: %s", i, err)
	***REMOVED***
	pk, ok := packet.(*PublicKeyV3)
	if !ok ***REMOVED***
		t.Fatalf("#%d: failed to parse, got: %#v", i, packet)
	***REMOVED***
	var serializeBuf bytes.Buffer
	if err = pk.Serialize(&serializeBuf); err != nil ***REMOVED***
		t.Fatalf("#%d: failed to serialize: %s", i, err)
	***REMOVED***

	if packet, err = Read(bytes.NewBuffer(serializeBuf.Bytes())); err != nil ***REMOVED***
		t.Fatalf("#%d: Read error (from serialized data): %s", i, err)
	***REMOVED***
	if pk, ok = packet.(*PublicKeyV3); !ok ***REMOVED***
		t.Fatalf("#%d: failed to parse serialized data, got: %#v", i, packet)
	***REMOVED***
***REMOVED***
