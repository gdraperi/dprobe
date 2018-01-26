// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x509

// RFC 1423 describes the encryption of PEM blocks. The algorithm used to
// generate a key from the password was derived by looking at the OpenSSL
// implementation.

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"io"
	"strings"
)

type PEMCipher int

// Possible values for the EncryptPEMBlock encryption algorithm.
const (
	_ PEMCipher = iota
	PEMCipherDES
	PEMCipher3DES
	PEMCipherAES128
	PEMCipherAES192
	PEMCipherAES256
)

// rfc1423Algo holds a method for enciphering a PEM block.
type rfc1423Algo struct ***REMOVED***
	cipher     PEMCipher
	name       string
	cipherFunc func(key []byte) (cipher.Block, error)
	keySize    int
	blockSize  int
***REMOVED***

// rfc1423Algos holds a slice of the possible ways to encrypt a PEM
// block.  The ivSize numbers were taken from the OpenSSL source.
var rfc1423Algos = []rfc1423Algo***REMOVED******REMOVED***
	cipher:     PEMCipherDES,
	name:       "DES-CBC",
	cipherFunc: des.NewCipher,
	keySize:    8,
	blockSize:  des.BlockSize,
***REMOVED***, ***REMOVED***
	cipher:     PEMCipher3DES,
	name:       "DES-EDE3-CBC",
	cipherFunc: des.NewTripleDESCipher,
	keySize:    24,
	blockSize:  des.BlockSize,
***REMOVED***, ***REMOVED***
	cipher:     PEMCipherAES128,
	name:       "AES-128-CBC",
	cipherFunc: aes.NewCipher,
	keySize:    16,
	blockSize:  aes.BlockSize,
***REMOVED***, ***REMOVED***
	cipher:     PEMCipherAES192,
	name:       "AES-192-CBC",
	cipherFunc: aes.NewCipher,
	keySize:    24,
	blockSize:  aes.BlockSize,
***REMOVED***, ***REMOVED***
	cipher:     PEMCipherAES256,
	name:       "AES-256-CBC",
	cipherFunc: aes.NewCipher,
	keySize:    32,
	blockSize:  aes.BlockSize,
***REMOVED***,
***REMOVED***

// deriveKey uses a key derivation function to stretch the password into a key
// with the number of bits our cipher requires. This algorithm was derived from
// the OpenSSL source.
func (c rfc1423Algo) deriveKey(password, salt []byte) []byte ***REMOVED***
	hash := md5.New()
	out := make([]byte, c.keySize)
	var digest []byte

	for i := 0; i < len(out); i += len(digest) ***REMOVED***
		hash.Reset()
		hash.Write(digest)
		hash.Write(password)
		hash.Write(salt)
		digest = hash.Sum(digest[:0])
		copy(out[i:], digest)
	***REMOVED***
	return out
***REMOVED***

// IsEncryptedPEMBlock returns if the PEM block is password encrypted.
func IsEncryptedPEMBlock(b *pem.Block) bool ***REMOVED***
	_, ok := b.Headers["DEK-Info"]
	return ok
***REMOVED***

// IncorrectPasswordError is returned when an incorrect password is detected.
var IncorrectPasswordError = errors.New("x509: decryption password incorrect")

// DecryptPEMBlock takes a password encrypted PEM block and the password used to
// encrypt it and returns a slice of decrypted DER encoded bytes. It inspects
// the DEK-Info header to determine the algorithm used for decryption. If no
// DEK-Info header is present, an error is returned. If an incorrect password
// is detected an IncorrectPasswordError is returned.
func DecryptPEMBlock(b *pem.Block, password []byte) ([]byte, error) ***REMOVED***
	dek, ok := b.Headers["DEK-Info"]
	if !ok ***REMOVED***
		return nil, errors.New("x509: no DEK-Info header in block")
	***REMOVED***

	idx := strings.Index(dek, ",")
	if idx == -1 ***REMOVED***
		return nil, errors.New("x509: malformed DEK-Info header")
	***REMOVED***

	mode, hexIV := dek[:idx], dek[idx+1:]
	ciph := cipherByName(mode)
	if ciph == nil ***REMOVED***
		return nil, errors.New("x509: unknown encryption mode")
	***REMOVED***
	iv, err := hex.DecodeString(hexIV)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(iv) != ciph.blockSize ***REMOVED***
		return nil, errors.New("x509: incorrect IV size")
	***REMOVED***

	// Based on the OpenSSL implementation. The salt is the first 8 bytes
	// of the initialization vector.
	key := ciph.deriveKey(password, iv[:8])
	block, err := ciph.cipherFunc(key)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	data := make([]byte, len(b.Bytes))
	dec := cipher.NewCBCDecrypter(block, iv)
	dec.CryptBlocks(data, b.Bytes)

	// Blocks are padded using a scheme where the last n bytes of padding are all
	// equal to n. It can pad from 1 to blocksize bytes inclusive. See RFC 1423.
	// For example:
	//	[x y z 2 2]
	//	[x y 7 7 7 7 7 7 7]
	// If we detect a bad padding, we assume it is an invalid password.
	dlen := len(data)
	if dlen == 0 || dlen%ciph.blockSize != 0 ***REMOVED***
		return nil, errors.New("x509: invalid padding")
	***REMOVED***
	last := int(data[dlen-1])
	if dlen < last ***REMOVED***
		return nil, IncorrectPasswordError
	***REMOVED***
	if last == 0 || last > ciph.blockSize ***REMOVED***
		return nil, IncorrectPasswordError
	***REMOVED***
	for _, val := range data[dlen-last:] ***REMOVED***
		if int(val) != last ***REMOVED***
			return nil, IncorrectPasswordError
		***REMOVED***
	***REMOVED***
	return data[:dlen-last], nil
***REMOVED***

// EncryptPEMBlock returns a PEM block of the specified type holding the
// given DER-encoded data encrypted with the specified algorithm and
// password.
func EncryptPEMBlock(rand io.Reader, blockType string, data, password []byte, alg PEMCipher) (*pem.Block, error) ***REMOVED***
	ciph := cipherByKey(alg)
	if ciph == nil ***REMOVED***
		return nil, errors.New("x509: unknown encryption mode")
	***REMOVED***
	iv := make([]byte, ciph.blockSize)
	if _, err := io.ReadFull(rand, iv); err != nil ***REMOVED***
		return nil, errors.New("x509: cannot generate IV: " + err.Error())
	***REMOVED***
	// The salt is the first 8 bytes of the initialization vector,
	// matching the key derivation in DecryptPEMBlock.
	key := ciph.deriveKey(password, iv[:8])
	block, err := ciph.cipherFunc(key)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	enc := cipher.NewCBCEncrypter(block, iv)
	pad := ciph.blockSize - len(data)%ciph.blockSize
	encrypted := make([]byte, len(data), len(data)+pad)
	// We could save this copy by encrypting all the whole blocks in
	// the data separately, but it doesn't seem worth the additional
	// code.
	copy(encrypted, data)
	// See RFC 1423, section 1.1
	for i := 0; i < pad; i++ ***REMOVED***
		encrypted = append(encrypted, byte(pad))
	***REMOVED***
	enc.CryptBlocks(encrypted, encrypted)

	return &pem.Block***REMOVED***
		Type: blockType,
		Headers: map[string]string***REMOVED***
			"Proc-Type": "4,ENCRYPTED",
			"DEK-Info":  ciph.name + "," + hex.EncodeToString(iv),
		***REMOVED***,
		Bytes: encrypted,
	***REMOVED***, nil
***REMOVED***

func cipherByName(name string) *rfc1423Algo ***REMOVED***
	for i := range rfc1423Algos ***REMOVED***
		alg := &rfc1423Algos[i]
		if alg.name == name ***REMOVED***
			return alg
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func cipherByKey(key PEMCipher) *rfc1423Algo ***REMOVED***
	for i := range rfc1423Algos ***REMOVED***
		alg := &rfc1423Algos[i]
		if alg.cipher == key ***REMOVED***
			return alg
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
