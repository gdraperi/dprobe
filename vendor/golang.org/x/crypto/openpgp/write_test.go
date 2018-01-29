// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package openpgp

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
	"time"

	"golang.org/x/crypto/openpgp/packet"
)

func TestSignDetached(t *testing.T) ***REMOVED***
	kring, _ := ReadKeyRing(readerFromHex(testKeys1And2PrivateHex))
	out := bytes.NewBuffer(nil)
	message := bytes.NewBufferString(signedInput)
	err := DetachSign(out, kring[0], message, nil)
	if err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***

	testDetachedSignature(t, kring, out, signedInput, "check", testKey1KeyId)
***REMOVED***

func TestSignTextDetached(t *testing.T) ***REMOVED***
	kring, _ := ReadKeyRing(readerFromHex(testKeys1And2PrivateHex))
	out := bytes.NewBuffer(nil)
	message := bytes.NewBufferString(signedInput)
	err := DetachSignText(out, kring[0], message, nil)
	if err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***

	testDetachedSignature(t, kring, out, signedInput, "check", testKey1KeyId)
***REMOVED***

func TestSignDetachedDSA(t *testing.T) ***REMOVED***
	kring, _ := ReadKeyRing(readerFromHex(dsaTestKeyPrivateHex))
	out := bytes.NewBuffer(nil)
	message := bytes.NewBufferString(signedInput)
	err := DetachSign(out, kring[0], message, nil)
	if err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***

	testDetachedSignature(t, kring, out, signedInput, "check", testKey3KeyId)
***REMOVED***

func TestSignDetachedP256(t *testing.T) ***REMOVED***
	kring, _ := ReadKeyRing(readerFromHex(p256TestKeyPrivateHex))
	kring[0].PrivateKey.Decrypt([]byte("passphrase"))

	out := bytes.NewBuffer(nil)
	message := bytes.NewBufferString(signedInput)
	err := DetachSign(out, kring[0], message, nil)
	if err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***

	testDetachedSignature(t, kring, out, signedInput, "check", testKeyP256KeyId)
***REMOVED***

func TestNewEntity(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		return
	***REMOVED***

	// Check bit-length with no config.
	e, err := NewEntity("Test User", "test", "test@example.com", nil)
	if err != nil ***REMOVED***
		t.Errorf("failed to create entity: %s", err)
		return
	***REMOVED***
	bl, err := e.PrimaryKey.BitLength()
	if err != nil ***REMOVED***
		t.Errorf("failed to find bit length: %s", err)
	***REMOVED***
	if int(bl) != defaultRSAKeyBits ***REMOVED***
		t.Errorf("BitLength %v, expected %v", int(bl), defaultRSAKeyBits)
	***REMOVED***

	// Check bit-length with a config.
	cfg := &packet.Config***REMOVED***RSABits: 1024***REMOVED***
	e, err = NewEntity("Test User", "test", "test@example.com", cfg)
	if err != nil ***REMOVED***
		t.Errorf("failed to create entity: %s", err)
		return
	***REMOVED***
	bl, err = e.PrimaryKey.BitLength()
	if err != nil ***REMOVED***
		t.Errorf("failed to find bit length: %s", err)
	***REMOVED***
	if int(bl) != cfg.RSABits ***REMOVED***
		t.Errorf("BitLength %v, expected %v", bl, cfg.RSABits)
	***REMOVED***

	w := bytes.NewBuffer(nil)
	if err := e.SerializePrivate(w, nil); err != nil ***REMOVED***
		t.Errorf("failed to serialize entity: %s", err)
		return
	***REMOVED***
	serialized := w.Bytes()

	el, err := ReadKeyRing(w)
	if err != nil ***REMOVED***
		t.Errorf("failed to reparse entity: %s", err)
		return
	***REMOVED***

	if len(el) != 1 ***REMOVED***
		t.Errorf("wrong number of entities found, got %d, want 1", len(el))
	***REMOVED***

	w = bytes.NewBuffer(nil)
	if err := e.SerializePrivate(w, nil); err != nil ***REMOVED***
		t.Errorf("failed to serialize entity second time: %s", err)
		return
	***REMOVED***

	if !bytes.Equal(w.Bytes(), serialized) ***REMOVED***
		t.Errorf("results differed")
	***REMOVED***
***REMOVED***

func TestSymmetricEncryption(t *testing.T) ***REMOVED***
	buf := new(bytes.Buffer)
	plaintext, err := SymmetricallyEncrypt(buf, []byte("testing"), nil, nil)
	if err != nil ***REMOVED***
		t.Errorf("error writing headers: %s", err)
		return
	***REMOVED***
	message := []byte("hello world\n")
	_, err = plaintext.Write(message)
	if err != nil ***REMOVED***
		t.Errorf("error writing to plaintext writer: %s", err)
	***REMOVED***
	err = plaintext.Close()
	if err != nil ***REMOVED***
		t.Errorf("error closing plaintext writer: %s", err)
	***REMOVED***

	md, err := ReadMessage(buf, nil, func(keys []Key, symmetric bool) ([]byte, error) ***REMOVED***
		return []byte("testing"), nil
	***REMOVED***, nil)
	if err != nil ***REMOVED***
		t.Errorf("error rereading message: %s", err)
	***REMOVED***
	messageBuf := bytes.NewBuffer(nil)
	_, err = io.Copy(messageBuf, md.UnverifiedBody)
	if err != nil ***REMOVED***
		t.Errorf("error rereading message: %s", err)
	***REMOVED***
	if !bytes.Equal(message, messageBuf.Bytes()) ***REMOVED***
		t.Errorf("recovered message incorrect got '%s', want '%s'", messageBuf.Bytes(), message)
	***REMOVED***
***REMOVED***

var testEncryptionTests = []struct ***REMOVED***
	keyRingHex string
	isSigned   bool
***REMOVED******REMOVED***
	***REMOVED***
		testKeys1And2PrivateHex,
		false,
	***REMOVED***,
	***REMOVED***
		testKeys1And2PrivateHex,
		true,
	***REMOVED***,
	***REMOVED***
		dsaElGamalTestKeysHex,
		false,
	***REMOVED***,
	***REMOVED***
		dsaElGamalTestKeysHex,
		true,
	***REMOVED***,
***REMOVED***

func TestEncryption(t *testing.T) ***REMOVED***
	for i, test := range testEncryptionTests ***REMOVED***
		kring, _ := ReadKeyRing(readerFromHex(test.keyRingHex))

		passphrase := []byte("passphrase")
		for _, entity := range kring ***REMOVED***
			if entity.PrivateKey != nil && entity.PrivateKey.Encrypted ***REMOVED***
				err := entity.PrivateKey.Decrypt(passphrase)
				if err != nil ***REMOVED***
					t.Errorf("#%d: failed to decrypt key", i)
				***REMOVED***
			***REMOVED***
			for _, subkey := range entity.Subkeys ***REMOVED***
				if subkey.PrivateKey != nil && subkey.PrivateKey.Encrypted ***REMOVED***
					err := subkey.PrivateKey.Decrypt(passphrase)
					if err != nil ***REMOVED***
						t.Errorf("#%d: failed to decrypt subkey", i)
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***

		var signed *Entity
		if test.isSigned ***REMOVED***
			signed = kring[0]
		***REMOVED***

		buf := new(bytes.Buffer)
		w, err := Encrypt(buf, kring[:1], signed, nil /* no hints */, nil)
		if err != nil ***REMOVED***
			t.Errorf("#%d: error in Encrypt: %s", i, err)
			continue
		***REMOVED***

		const message = "testing"
		_, err = w.Write([]byte(message))
		if err != nil ***REMOVED***
			t.Errorf("#%d: error writing plaintext: %s", i, err)
			continue
		***REMOVED***
		err = w.Close()
		if err != nil ***REMOVED***
			t.Errorf("#%d: error closing WriteCloser: %s", i, err)
			continue
		***REMOVED***

		md, err := ReadMessage(buf, kring, nil /* no prompt */, nil)
		if err != nil ***REMOVED***
			t.Errorf("#%d: error reading message: %s", i, err)
			continue
		***REMOVED***

		testTime, _ := time.Parse("2006-01-02", "2013-07-01")
		if test.isSigned ***REMOVED***
			signKey, _ := kring[0].signingKey(testTime)
			expectedKeyId := signKey.PublicKey.KeyId
			if md.SignedByKeyId != expectedKeyId ***REMOVED***
				t.Errorf("#%d: message signed by wrong key id, got: %v, want: %v", i, *md.SignedBy, expectedKeyId)
			***REMOVED***
			if md.SignedBy == nil ***REMOVED***
				t.Errorf("#%d: failed to find the signing Entity", i)
			***REMOVED***
		***REMOVED***

		plaintext, err := ioutil.ReadAll(md.UnverifiedBody)
		if err != nil ***REMOVED***
			t.Errorf("#%d: error reading encrypted contents: %s", i, err)
			continue
		***REMOVED***

		encryptKey, _ := kring[0].encryptionKey(testTime)
		expectedKeyId := encryptKey.PublicKey.KeyId
		if len(md.EncryptedToKeyIds) != 1 || md.EncryptedToKeyIds[0] != expectedKeyId ***REMOVED***
			t.Errorf("#%d: expected message to be encrypted to %v, but got %#v", i, expectedKeyId, md.EncryptedToKeyIds)
		***REMOVED***

		if string(plaintext) != message ***REMOVED***
			t.Errorf("#%d: got: %s, want: %s", i, string(plaintext), message)
		***REMOVED***

		if test.isSigned ***REMOVED***
			if md.SignatureError != nil ***REMOVED***
				t.Errorf("#%d: signature error: %s", i, md.SignatureError)
			***REMOVED***
			if md.Signature == nil ***REMOVED***
				t.Error("signature missing")
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
