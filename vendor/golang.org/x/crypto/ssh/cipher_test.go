// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"testing"
)

func TestDefaultCiphersExist(t *testing.T) ***REMOVED***
	for _, cipherAlgo := range supportedCiphers ***REMOVED***
		if _, ok := cipherModes[cipherAlgo]; !ok ***REMOVED***
			t.Errorf("supported cipher %q is unknown", cipherAlgo)
		***REMOVED***
	***REMOVED***
	for _, cipherAlgo := range preferredCiphers ***REMOVED***
		if _, ok := cipherModes[cipherAlgo]; !ok ***REMOVED***
			t.Errorf("preferred cipher %q is unknown", cipherAlgo)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestPacketCiphers(t *testing.T) ***REMOVED***
	defaultMac := "hmac-sha2-256"
	defaultCipher := "aes128-ctr"
	for cipher := range cipherModes ***REMOVED***
		t.Run("cipher="+cipher,
			func(t *testing.T) ***REMOVED*** testPacketCipher(t, cipher, defaultMac) ***REMOVED***)
	***REMOVED***
	for mac := range macModes ***REMOVED***
		t.Run("mac="+mac,
			func(t *testing.T) ***REMOVED*** testPacketCipher(t, defaultCipher, mac) ***REMOVED***)
	***REMOVED***
***REMOVED***

func testPacketCipher(t *testing.T, cipher, mac string) ***REMOVED***
	kr := &kexResult***REMOVED***Hash: crypto.SHA1***REMOVED***
	algs := directionAlgorithms***REMOVED***
		Cipher:      cipher,
		MAC:         mac,
		Compression: "none",
	***REMOVED***
	client, err := newPacketCipher(clientKeys, algs, kr)
	if err != nil ***REMOVED***
		t.Fatalf("newPacketCipher(client, %q, %q): %v", cipher, mac, err)
	***REMOVED***
	server, err := newPacketCipher(clientKeys, algs, kr)
	if err != nil ***REMOVED***
		t.Fatalf("newPacketCipher(client, %q, %q): %v", cipher, mac, err)
	***REMOVED***

	want := "bla bla"
	input := []byte(want)
	buf := &bytes.Buffer***REMOVED******REMOVED***
	if err := client.writePacket(0, buf, rand.Reader, input); err != nil ***REMOVED***
		t.Fatalf("writePacket(%q, %q): %v", cipher, mac, err)
	***REMOVED***

	packet, err := server.readPacket(0, buf)
	if err != nil ***REMOVED***
		t.Fatalf("readPacket(%q, %q): %v", cipher, mac, err)
	***REMOVED***

	if string(packet) != want ***REMOVED***
		t.Errorf("roundtrip(%q, %q): got %q, want %q", cipher, mac, packet, want)
	***REMOVED***
***REMOVED***

func TestCBCOracleCounterMeasure(t *testing.T) ***REMOVED***
	kr := &kexResult***REMOVED***Hash: crypto.SHA1***REMOVED***
	algs := directionAlgorithms***REMOVED***
		Cipher:      aes128cbcID,
		MAC:         "hmac-sha1",
		Compression: "none",
	***REMOVED***
	client, err := newPacketCipher(clientKeys, algs, kr)
	if err != nil ***REMOVED***
		t.Fatalf("newPacketCipher(client): %v", err)
	***REMOVED***

	want := "bla bla"
	input := []byte(want)
	buf := &bytes.Buffer***REMOVED******REMOVED***
	if err := client.writePacket(0, buf, rand.Reader, input); err != nil ***REMOVED***
		t.Errorf("writePacket: %v", err)
	***REMOVED***

	packetSize := buf.Len()
	buf.Write(make([]byte, 2*maxPacket))

	// We corrupt each byte, but this usually will only test the
	// 'packet too large' or 'MAC failure' cases.
	lastRead := -1
	for i := 0; i < packetSize; i++ ***REMOVED***
		server, err := newPacketCipher(clientKeys, algs, kr)
		if err != nil ***REMOVED***
			t.Fatalf("newPacketCipher(client): %v", err)
		***REMOVED***

		fresh := &bytes.Buffer***REMOVED******REMOVED***
		fresh.Write(buf.Bytes())
		fresh.Bytes()[i] ^= 0x01

		before := fresh.Len()
		_, err = server.readPacket(0, fresh)
		if err == nil ***REMOVED***
			t.Errorf("corrupt byte %d: readPacket succeeded ", i)
			continue
		***REMOVED***
		if _, ok := err.(cbcError); !ok ***REMOVED***
			t.Errorf("corrupt byte %d: got %v (%T), want cbcError", i, err, err)
			continue
		***REMOVED***

		after := fresh.Len()
		bytesRead := before - after
		if bytesRead < maxPacket ***REMOVED***
			t.Errorf("corrupt byte %d: read %d bytes, want more than %d", i, bytesRead, maxPacket)
			continue
		***REMOVED***

		if i > 0 && bytesRead != lastRead ***REMOVED***
			t.Errorf("corrupt byte %d: read %d bytes, want %d bytes read", i, bytesRead, lastRead)
		***REMOVED***
		lastRead = bytesRead
	***REMOVED***
***REMOVED***
