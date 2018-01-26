// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkcs12

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"

	"golang.org/x/crypto/pkcs12/internal/rc2"
)

var (
	oidPBEWithSHAAnd3KeyTripleDESCBC = asn1.ObjectIdentifier([]int***REMOVED***1, 2, 840, 113549, 1, 12, 1, 3***REMOVED***)
	oidPBEWithSHAAnd40BitRC2CBC      = asn1.ObjectIdentifier([]int***REMOVED***1, 2, 840, 113549, 1, 12, 1, 6***REMOVED***)
)

// pbeCipher is an abstraction of a PKCS#12 cipher.
type pbeCipher interface ***REMOVED***
	// create returns a cipher.Block given a key.
	create(key []byte) (cipher.Block, error)
	// deriveKey returns a key derived from the given password and salt.
	deriveKey(salt, password []byte, iterations int) []byte
	// deriveKey returns an IV derived from the given password and salt.
	deriveIV(salt, password []byte, iterations int) []byte
***REMOVED***

type shaWithTripleDESCBC struct***REMOVED******REMOVED***

func (shaWithTripleDESCBC) create(key []byte) (cipher.Block, error) ***REMOVED***
	return des.NewTripleDESCipher(key)
***REMOVED***

func (shaWithTripleDESCBC) deriveKey(salt, password []byte, iterations int) []byte ***REMOVED***
	return pbkdf(sha1Sum, 20, 64, salt, password, iterations, 1, 24)
***REMOVED***

func (shaWithTripleDESCBC) deriveIV(salt, password []byte, iterations int) []byte ***REMOVED***
	return pbkdf(sha1Sum, 20, 64, salt, password, iterations, 2, 8)
***REMOVED***

type shaWith40BitRC2CBC struct***REMOVED******REMOVED***

func (shaWith40BitRC2CBC) create(key []byte) (cipher.Block, error) ***REMOVED***
	return rc2.New(key, len(key)*8)
***REMOVED***

func (shaWith40BitRC2CBC) deriveKey(salt, password []byte, iterations int) []byte ***REMOVED***
	return pbkdf(sha1Sum, 20, 64, salt, password, iterations, 1, 5)
***REMOVED***

func (shaWith40BitRC2CBC) deriveIV(salt, password []byte, iterations int) []byte ***REMOVED***
	return pbkdf(sha1Sum, 20, 64, salt, password, iterations, 2, 8)
***REMOVED***

type pbeParams struct ***REMOVED***
	Salt       []byte
	Iterations int
***REMOVED***

func pbDecrypterFor(algorithm pkix.AlgorithmIdentifier, password []byte) (cipher.BlockMode, int, error) ***REMOVED***
	var cipherType pbeCipher

	switch ***REMOVED***
	case algorithm.Algorithm.Equal(oidPBEWithSHAAnd3KeyTripleDESCBC):
		cipherType = shaWithTripleDESCBC***REMOVED******REMOVED***
	case algorithm.Algorithm.Equal(oidPBEWithSHAAnd40BitRC2CBC):
		cipherType = shaWith40BitRC2CBC***REMOVED******REMOVED***
	default:
		return nil, 0, NotImplementedError("algorithm " + algorithm.Algorithm.String() + " is not supported")
	***REMOVED***

	var params pbeParams
	if err := unmarshal(algorithm.Parameters.FullBytes, &params); err != nil ***REMOVED***
		return nil, 0, err
	***REMOVED***

	key := cipherType.deriveKey(params.Salt, password, params.Iterations)
	iv := cipherType.deriveIV(params.Salt, password, params.Iterations)

	block, err := cipherType.create(key)
	if err != nil ***REMOVED***
		return nil, 0, err
	***REMOVED***

	return cipher.NewCBCDecrypter(block, iv), block.BlockSize(), nil
***REMOVED***

func pbDecrypt(info decryptable, password []byte) (decrypted []byte, err error) ***REMOVED***
	cbc, blockSize, err := pbDecrypterFor(info.Algorithm(), password)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	encrypted := info.Data()
	if len(encrypted) == 0 ***REMOVED***
		return nil, errors.New("pkcs12: empty encrypted data")
	***REMOVED***
	if len(encrypted)%blockSize != 0 ***REMOVED***
		return nil, errors.New("pkcs12: input is not a multiple of the block size")
	***REMOVED***
	decrypted = make([]byte, len(encrypted))
	cbc.CryptBlocks(decrypted, encrypted)

	psLen := int(decrypted[len(decrypted)-1])
	if psLen == 0 || psLen > blockSize ***REMOVED***
		return nil, ErrDecryption
	***REMOVED***

	if len(decrypted) < psLen ***REMOVED***
		return nil, ErrDecryption
	***REMOVED***
	ps := decrypted[len(decrypted)-psLen:]
	decrypted = decrypted[:len(decrypted)-psLen]
	if bytes.Compare(ps, bytes.Repeat([]byte***REMOVED***byte(psLen)***REMOVED***, psLen)) != 0 ***REMOVED***
		return nil, ErrDecryption
	***REMOVED***

	return
***REMOVED***

// decryptable abstracts a object that contains ciphertext.
type decryptable interface ***REMOVED***
	Algorithm() pkix.AlgorithmIdentifier
	Data() []byte
***REMOVED***
