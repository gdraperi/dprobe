// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"bytes"
	"encoding/hex"
	"io"
	"io/ioutil"
	"testing"
)

func TestSymmetricKeyEncrypted(t *testing.T) ***REMOVED***
	buf := readerFromHex(symmetricallyEncryptedHex)
	packet, err := Read(buf)
	if err != nil ***REMOVED***
		t.Errorf("failed to read SymmetricKeyEncrypted: %s", err)
		return
	***REMOVED***
	ske, ok := packet.(*SymmetricKeyEncrypted)
	if !ok ***REMOVED***
		t.Error("didn't find SymmetricKeyEncrypted packet")
		return
	***REMOVED***
	key, cipherFunc, err := ske.Decrypt([]byte("password"))
	if err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***

	packet, err = Read(buf)
	if err != nil ***REMOVED***
		t.Errorf("failed to read SymmetricallyEncrypted: %s", err)
		return
	***REMOVED***
	se, ok := packet.(*SymmetricallyEncrypted)
	if !ok ***REMOVED***
		t.Error("didn't find SymmetricallyEncrypted packet")
		return
	***REMOVED***
	r, err := se.Decrypt(cipherFunc, key)
	if err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***

	contents, err := ioutil.ReadAll(r)
	if err != nil && err != io.EOF ***REMOVED***
		t.Error(err)
		return
	***REMOVED***

	expectedContents, _ := hex.DecodeString(symmetricallyEncryptedContentsHex)
	if !bytes.Equal(expectedContents, contents) ***REMOVED***
		t.Errorf("bad contents got:%x want:%x", contents, expectedContents)
	***REMOVED***
***REMOVED***

const symmetricallyEncryptedHex = "8c0d04030302371a0b38d884f02060c91cf97c9973b8e58e028e9501708ccfe618fb92afef7fa2d80ddadd93cf"
const symmetricallyEncryptedContentsHex = "cb1062004d14c4df636f6e74656e74732e0a"

func TestSerializeSymmetricKeyEncryptedCiphers(t *testing.T) ***REMOVED***
	tests := [...]struct ***REMOVED***
		cipherFunc CipherFunction
		name       string
	***REMOVED******REMOVED***
		***REMOVED***Cipher3DES, "Cipher3DES"***REMOVED***,
		***REMOVED***CipherCAST5, "CipherCAST5"***REMOVED***,
		***REMOVED***CipherAES128, "CipherAES128"***REMOVED***,
		***REMOVED***CipherAES192, "CipherAES192"***REMOVED***,
		***REMOVED***CipherAES256, "CipherAES256"***REMOVED***,
	***REMOVED***

	for _, test := range tests ***REMOVED***
		var buf bytes.Buffer
		passphrase := []byte("testing")
		config := &Config***REMOVED***
			DefaultCipher: test.cipherFunc,
		***REMOVED***

		key, err := SerializeSymmetricKeyEncrypted(&buf, passphrase, config)
		if err != nil ***REMOVED***
			t.Errorf("cipher(%s) failed to serialize: %s", test.name, err)
			continue
		***REMOVED***

		p, err := Read(&buf)
		if err != nil ***REMOVED***
			t.Errorf("cipher(%s) failed to reparse: %s", test.name, err)
			continue
		***REMOVED***

		ske, ok := p.(*SymmetricKeyEncrypted)
		if !ok ***REMOVED***
			t.Errorf("cipher(%s) parsed a different packet type: %#v", test.name, p)
			continue
		***REMOVED***

		if ske.CipherFunc != config.DefaultCipher ***REMOVED***
			t.Errorf("cipher(%s) SKE cipher function is %d (expected %d)", test.name, ske.CipherFunc, config.DefaultCipher)
		***REMOVED***
		parsedKey, parsedCipherFunc, err := ske.Decrypt(passphrase)
		if err != nil ***REMOVED***
			t.Errorf("cipher(%s) failed to decrypt reparsed SKE: %s", test.name, err)
			continue
		***REMOVED***
		if !bytes.Equal(key, parsedKey) ***REMOVED***
			t.Errorf("cipher(%s) keys don't match after Decrypt: %x (original) vs %x (parsed)", test.name, key, parsedKey)
		***REMOVED***
		if parsedCipherFunc != test.cipherFunc ***REMOVED***
			t.Errorf("cipher(%s) cipher function doesn't match after Decrypt: %d (original) vs %d (parsed)",
				test.name, test.cipherFunc, parsedCipherFunc)
		***REMOVED***
	***REMOVED***
***REMOVED***
