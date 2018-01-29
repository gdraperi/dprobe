// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chacha20poly1305

import (
	"bytes"
	cr "crypto/rand"
	"encoding/hex"
	mr "math/rand"
	"testing"
)

func TestVectors(t *testing.T) ***REMOVED***
	for i, test := range chacha20Poly1305Tests ***REMOVED***
		key, _ := hex.DecodeString(test.key)
		nonce, _ := hex.DecodeString(test.nonce)
		ad, _ := hex.DecodeString(test.aad)
		plaintext, _ := hex.DecodeString(test.plaintext)

		aead, err := New(key)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		ct := aead.Seal(nil, nonce, plaintext, ad)
		if ctHex := hex.EncodeToString(ct); ctHex != test.out ***REMOVED***
			t.Errorf("#%d: got %s, want %s", i, ctHex, test.out)
			continue
		***REMOVED***

		plaintext2, err := aead.Open(nil, nonce, ct, ad)
		if err != nil ***REMOVED***
			t.Errorf("#%d: Open failed", i)
			continue
		***REMOVED***

		if !bytes.Equal(plaintext, plaintext2) ***REMOVED***
			t.Errorf("#%d: plaintext's don't match: got %x vs %x", i, plaintext2, plaintext)
			continue
		***REMOVED***

		if len(ad) > 0 ***REMOVED***
			alterAdIdx := mr.Intn(len(ad))
			ad[alterAdIdx] ^= 0x80
			if _, err := aead.Open(nil, nonce, ct, ad); err == nil ***REMOVED***
				t.Errorf("#%d: Open was successful after altering additional data", i)
			***REMOVED***
			ad[alterAdIdx] ^= 0x80
		***REMOVED***

		alterNonceIdx := mr.Intn(aead.NonceSize())
		nonce[alterNonceIdx] ^= 0x80
		if _, err := aead.Open(nil, nonce, ct, ad); err == nil ***REMOVED***
			t.Errorf("#%d: Open was successful after altering nonce", i)
		***REMOVED***
		nonce[alterNonceIdx] ^= 0x80

		alterCtIdx := mr.Intn(len(ct))
		ct[alterCtIdx] ^= 0x80
		if _, err := aead.Open(nil, nonce, ct, ad); err == nil ***REMOVED***
			t.Errorf("#%d: Open was successful after altering ciphertext", i)
		***REMOVED***
		ct[alterCtIdx] ^= 0x80
	***REMOVED***
***REMOVED***

func TestRandom(t *testing.T) ***REMOVED***
	// Some random tests to verify Open(Seal) == Plaintext
	for i := 0; i < 256; i++ ***REMOVED***
		var nonce [12]byte
		var key [32]byte

		al := mr.Intn(128)
		pl := mr.Intn(16384)
		ad := make([]byte, al)
		plaintext := make([]byte, pl)
		cr.Read(key[:])
		cr.Read(nonce[:])
		cr.Read(ad)
		cr.Read(plaintext)

		aead, err := New(key[:])
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		ct := aead.Seal(nil, nonce[:], plaintext, ad)

		plaintext2, err := aead.Open(nil, nonce[:], ct, ad)
		if err != nil ***REMOVED***
			t.Errorf("Random #%d: Open failed", i)
			continue
		***REMOVED***

		if !bytes.Equal(plaintext, plaintext2) ***REMOVED***
			t.Errorf("Random #%d: plaintext's don't match: got %x vs %x", i, plaintext2, plaintext)
			continue
		***REMOVED***

		if len(ad) > 0 ***REMOVED***
			alterAdIdx := mr.Intn(len(ad))
			ad[alterAdIdx] ^= 0x80
			if _, err := aead.Open(nil, nonce[:], ct, ad); err == nil ***REMOVED***
				t.Errorf("Random #%d: Open was successful after altering additional data", i)
			***REMOVED***
			ad[alterAdIdx] ^= 0x80
		***REMOVED***

		alterNonceIdx := mr.Intn(aead.NonceSize())
		nonce[alterNonceIdx] ^= 0x80
		if _, err := aead.Open(nil, nonce[:], ct, ad); err == nil ***REMOVED***
			t.Errorf("Random #%d: Open was successful after altering nonce", i)
		***REMOVED***
		nonce[alterNonceIdx] ^= 0x80

		alterCtIdx := mr.Intn(len(ct))
		ct[alterCtIdx] ^= 0x80
		if _, err := aead.Open(nil, nonce[:], ct, ad); err == nil ***REMOVED***
			t.Errorf("Random #%d: Open was successful after altering ciphertext", i)
		***REMOVED***
		ct[alterCtIdx] ^= 0x80
	***REMOVED***
***REMOVED***

func benchamarkChaCha20Poly1305Seal(b *testing.B, buf []byte) ***REMOVED***
	b.SetBytes(int64(len(buf)))

	var key [32]byte
	var nonce [12]byte
	var ad [13]byte
	var out []byte

	aead, _ := New(key[:])
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		out = aead.Seal(out[:0], nonce[:], buf[:], ad[:])
	***REMOVED***
***REMOVED***

func benchamarkChaCha20Poly1305Open(b *testing.B, buf []byte) ***REMOVED***
	b.SetBytes(int64(len(buf)))

	var key [32]byte
	var nonce [12]byte
	var ad [13]byte
	var ct []byte
	var out []byte

	aead, _ := New(key[:])
	ct = aead.Seal(ct[:0], nonce[:], buf[:], ad[:])

	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		out, _ = aead.Open(out[:0], nonce[:], ct[:], ad[:])
	***REMOVED***
***REMOVED***

func BenchmarkChacha20Poly1305Open_64(b *testing.B) ***REMOVED***
	benchamarkChaCha20Poly1305Open(b, make([]byte, 64))
***REMOVED***

func BenchmarkChacha20Poly1305Seal_64(b *testing.B) ***REMOVED***
	benchamarkChaCha20Poly1305Seal(b, make([]byte, 64))
***REMOVED***

func BenchmarkChacha20Poly1305Open_1350(b *testing.B) ***REMOVED***
	benchamarkChaCha20Poly1305Open(b, make([]byte, 1350))
***REMOVED***

func BenchmarkChacha20Poly1305Seal_1350(b *testing.B) ***REMOVED***
	benchamarkChaCha20Poly1305Seal(b, make([]byte, 1350))
***REMOVED***

func BenchmarkChacha20Poly1305Open_8K(b *testing.B) ***REMOVED***
	benchamarkChaCha20Poly1305Open(b, make([]byte, 8*1024))
***REMOVED***

func BenchmarkChacha20Poly1305Seal_8K(b *testing.B) ***REMOVED***
	benchamarkChaCha20Poly1305Seal(b, make([]byte, 8*1024))
***REMOVED***
