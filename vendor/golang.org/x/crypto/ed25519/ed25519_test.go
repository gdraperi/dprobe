// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ed25519

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto"
	"crypto/rand"
	"encoding/hex"
	"os"
	"strings"
	"testing"

	"golang.org/x/crypto/ed25519/internal/edwards25519"
)

type zeroReader struct***REMOVED******REMOVED***

func (zeroReader) Read(buf []byte) (int, error) ***REMOVED***
	for i := range buf ***REMOVED***
		buf[i] = 0
	***REMOVED***
	return len(buf), nil
***REMOVED***

func TestUnmarshalMarshal(t *testing.T) ***REMOVED***
	pub, _, _ := GenerateKey(rand.Reader)

	var A edwards25519.ExtendedGroupElement
	var pubBytes [32]byte
	copy(pubBytes[:], pub)
	if !A.FromBytes(&pubBytes) ***REMOVED***
		t.Fatalf("ExtendedGroupElement.FromBytes failed")
	***REMOVED***

	var pub2 [32]byte
	A.ToBytes(&pub2)

	if pubBytes != pub2 ***REMOVED***
		t.Errorf("FromBytes(%v)->ToBytes does not round-trip, got %x\n", pubBytes, pub2)
	***REMOVED***
***REMOVED***

func TestSignVerify(t *testing.T) ***REMOVED***
	var zero zeroReader
	public, private, _ := GenerateKey(zero)

	message := []byte("test message")
	sig := Sign(private, message)
	if !Verify(public, message, sig) ***REMOVED***
		t.Errorf("valid signature rejected")
	***REMOVED***

	wrongMessage := []byte("wrong message")
	if Verify(public, wrongMessage, sig) ***REMOVED***
		t.Errorf("signature of different message accepted")
	***REMOVED***
***REMOVED***

func TestCryptoSigner(t *testing.T) ***REMOVED***
	var zero zeroReader
	public, private, _ := GenerateKey(zero)

	signer := crypto.Signer(private)

	publicInterface := signer.Public()
	public2, ok := publicInterface.(PublicKey)
	if !ok ***REMOVED***
		t.Fatalf("expected PublicKey from Public() but got %T", publicInterface)
	***REMOVED***

	if !bytes.Equal(public, public2) ***REMOVED***
		t.Errorf("public keys do not match: original:%x vs Public():%x", public, public2)
	***REMOVED***

	message := []byte("message")
	var noHash crypto.Hash
	signature, err := signer.Sign(zero, message, noHash)
	if err != nil ***REMOVED***
		t.Fatalf("error from Sign(): %s", err)
	***REMOVED***

	if !Verify(public, message, signature) ***REMOVED***
		t.Errorf("Verify failed on signature from Sign()")
	***REMOVED***
***REMOVED***

func TestGolden(t *testing.T) ***REMOVED***
	// sign.input.gz is a selection of test cases from
	// https://ed25519.cr.yp.to/python/sign.input
	testDataZ, err := os.Open("testdata/sign.input.gz")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer testDataZ.Close()
	testData, err := gzip.NewReader(testDataZ)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer testData.Close()

	scanner := bufio.NewScanner(testData)
	lineNo := 0

	for scanner.Scan() ***REMOVED***
		lineNo++

		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 5 ***REMOVED***
			t.Fatalf("bad number of parts on line %d", lineNo)
		***REMOVED***

		privBytes, _ := hex.DecodeString(parts[0])
		pubKey, _ := hex.DecodeString(parts[1])
		msg, _ := hex.DecodeString(parts[2])
		sig, _ := hex.DecodeString(parts[3])
		// The signatures in the test vectors also include the message
		// at the end, but we just want R and S.
		sig = sig[:SignatureSize]

		if l := len(pubKey); l != PublicKeySize ***REMOVED***
			t.Fatalf("bad public key length on line %d: got %d bytes", lineNo, l)
		***REMOVED***

		var priv [PrivateKeySize]byte
		copy(priv[:], privBytes)
		copy(priv[32:], pubKey)

		sig2 := Sign(priv[:], msg)
		if !bytes.Equal(sig, sig2[:]) ***REMOVED***
			t.Errorf("different signature result on line %d: %x vs %x", lineNo, sig, sig2)
		***REMOVED***

		if !Verify(pubKey, msg, sig2) ***REMOVED***
			t.Errorf("signature failed to verify on line %d", lineNo)
		***REMOVED***
	***REMOVED***

	if err := scanner.Err(); err != nil ***REMOVED***
		t.Fatalf("error reading test data: %s", err)
	***REMOVED***
***REMOVED***

func BenchmarkKeyGeneration(b *testing.B) ***REMOVED***
	var zero zeroReader
	for i := 0; i < b.N; i++ ***REMOVED***
		if _, _, err := GenerateKey(zero); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkSigning(b *testing.B) ***REMOVED***
	var zero zeroReader
	_, priv, err := GenerateKey(zero)
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	message := []byte("Hello, world!")
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		Sign(priv, message)
	***REMOVED***
***REMOVED***

func BenchmarkVerification(b *testing.B) ***REMOVED***
	var zero zeroReader
	pub, priv, err := GenerateKey(zero)
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	message := []byte("Hello, world!")
	signature := Sign(priv, message)
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		Verify(pub, message, signature)
	***REMOVED***
***REMOVED***
