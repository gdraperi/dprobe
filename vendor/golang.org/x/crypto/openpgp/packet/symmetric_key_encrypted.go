// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"bytes"
	"crypto/cipher"
	"io"
	"strconv"

	"golang.org/x/crypto/openpgp/errors"
	"golang.org/x/crypto/openpgp/s2k"
)

// This is the largest session key that we'll support. Since no 512-bit cipher
// has even been seriously used, this is comfortably large.
const maxSessionKeySizeInBytes = 64

// SymmetricKeyEncrypted represents a passphrase protected session key. See RFC
// 4880, section 5.3.
type SymmetricKeyEncrypted struct ***REMOVED***
	CipherFunc   CipherFunction
	s2k          func(out, in []byte)
	encryptedKey []byte
***REMOVED***

const symmetricKeyEncryptedVersion = 4

func (ske *SymmetricKeyEncrypted) parse(r io.Reader) error ***REMOVED***
	// RFC 4880, section 5.3.
	var buf [2]byte
	if _, err := readFull(r, buf[:]); err != nil ***REMOVED***
		return err
	***REMOVED***
	if buf[0] != symmetricKeyEncryptedVersion ***REMOVED***
		return errors.UnsupportedError("SymmetricKeyEncrypted version")
	***REMOVED***
	ske.CipherFunc = CipherFunction(buf[1])

	if ske.CipherFunc.KeySize() == 0 ***REMOVED***
		return errors.UnsupportedError("unknown cipher: " + strconv.Itoa(int(buf[1])))
	***REMOVED***

	var err error
	ske.s2k, err = s2k.Parse(r)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	encryptedKey := make([]byte, maxSessionKeySizeInBytes)
	// The session key may follow. We just have to try and read to find
	// out. If it exists then we limit it to maxSessionKeySizeInBytes.
	n, err := readFull(r, encryptedKey)
	if err != nil && err != io.ErrUnexpectedEOF ***REMOVED***
		return err
	***REMOVED***

	if n != 0 ***REMOVED***
		if n == maxSessionKeySizeInBytes ***REMOVED***
			return errors.UnsupportedError("oversized encrypted session key")
		***REMOVED***
		ske.encryptedKey = encryptedKey[:n]
	***REMOVED***

	return nil
***REMOVED***

// Decrypt attempts to decrypt an encrypted session key and returns the key and
// the cipher to use when decrypting a subsequent Symmetrically Encrypted Data
// packet.
func (ske *SymmetricKeyEncrypted) Decrypt(passphrase []byte) ([]byte, CipherFunction, error) ***REMOVED***
	key := make([]byte, ske.CipherFunc.KeySize())
	ske.s2k(key, passphrase)

	if len(ske.encryptedKey) == 0 ***REMOVED***
		return key, ske.CipherFunc, nil
	***REMOVED***

	// the IV is all zeros
	iv := make([]byte, ske.CipherFunc.blockSize())
	c := cipher.NewCFBDecrypter(ske.CipherFunc.new(key), iv)
	plaintextKey := make([]byte, len(ske.encryptedKey))
	c.XORKeyStream(plaintextKey, ske.encryptedKey)
	cipherFunc := CipherFunction(plaintextKey[0])
	if cipherFunc.blockSize() == 0 ***REMOVED***
		return nil, ske.CipherFunc, errors.UnsupportedError("unknown cipher: " + strconv.Itoa(int(cipherFunc)))
	***REMOVED***
	plaintextKey = plaintextKey[1:]
	if l, cipherKeySize := len(plaintextKey), cipherFunc.KeySize(); l != cipherFunc.KeySize() ***REMOVED***
		return nil, cipherFunc, errors.StructuralError("length of decrypted key (" + strconv.Itoa(l) + ") " +
			"not equal to cipher keysize (" + strconv.Itoa(cipherKeySize) + ")")
	***REMOVED***
	return plaintextKey, cipherFunc, nil
***REMOVED***

// SerializeSymmetricKeyEncrypted serializes a symmetric key packet to w. The
// packet contains a random session key, encrypted by a key derived from the
// given passphrase. The session key is returned and must be passed to
// SerializeSymmetricallyEncrypted.
// If config is nil, sensible defaults will be used.
func SerializeSymmetricKeyEncrypted(w io.Writer, passphrase []byte, config *Config) (key []byte, err error) ***REMOVED***
	cipherFunc := config.Cipher()
	keySize := cipherFunc.KeySize()
	if keySize == 0 ***REMOVED***
		return nil, errors.UnsupportedError("unknown cipher: " + strconv.Itoa(int(cipherFunc)))
	***REMOVED***

	s2kBuf := new(bytes.Buffer)
	keyEncryptingKey := make([]byte, keySize)
	// s2k.Serialize salts and stretches the passphrase, and writes the
	// resulting key to keyEncryptingKey and the s2k descriptor to s2kBuf.
	err = s2k.Serialize(s2kBuf, keyEncryptingKey, config.Random(), passphrase, &s2k.Config***REMOVED***Hash: config.Hash(), S2KCount: config.PasswordHashIterations()***REMOVED***)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	s2kBytes := s2kBuf.Bytes()

	packetLength := 2 /* header */ + len(s2kBytes) + 1 /* cipher type */ + keySize
	err = serializeHeader(w, packetTypeSymmetricKeyEncrypted, packetLength)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	var buf [2]byte
	buf[0] = symmetricKeyEncryptedVersion
	buf[1] = byte(cipherFunc)
	_, err = w.Write(buf[:])
	if err != nil ***REMOVED***
		return
	***REMOVED***
	_, err = w.Write(s2kBytes)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	sessionKey := make([]byte, keySize)
	_, err = io.ReadFull(config.Random(), sessionKey)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	iv := make([]byte, cipherFunc.blockSize())
	c := cipher.NewCFBEncrypter(cipherFunc.new(keyEncryptingKey), iv)
	encryptedCipherAndKey := make([]byte, keySize+1)
	c.XORKeyStream(encryptedCipherAndKey, buf[1:])
	c.XORKeyStream(encryptedCipherAndKey[1:], sessionKey)
	_, err = w.Write(encryptedCipherAndKey)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	key = sessionKey
	return
***REMOVED***
