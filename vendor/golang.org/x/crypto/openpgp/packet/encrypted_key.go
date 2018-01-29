// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"crypto/rsa"
	"encoding/binary"
	"io"
	"math/big"
	"strconv"

	"golang.org/x/crypto/openpgp/elgamal"
	"golang.org/x/crypto/openpgp/errors"
)

const encryptedKeyVersion = 3

// EncryptedKey represents a public-key encrypted session key. See RFC 4880,
// section 5.1.
type EncryptedKey struct ***REMOVED***
	KeyId      uint64
	Algo       PublicKeyAlgorithm
	CipherFunc CipherFunction // only valid after a successful Decrypt
	Key        []byte         // only valid after a successful Decrypt

	encryptedMPI1, encryptedMPI2 parsedMPI
***REMOVED***

func (e *EncryptedKey) parse(r io.Reader) (err error) ***REMOVED***
	var buf [10]byte
	_, err = readFull(r, buf[:])
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if buf[0] != encryptedKeyVersion ***REMOVED***
		return errors.UnsupportedError("unknown EncryptedKey version " + strconv.Itoa(int(buf[0])))
	***REMOVED***
	e.KeyId = binary.BigEndian.Uint64(buf[1:9])
	e.Algo = PublicKeyAlgorithm(buf[9])
	switch e.Algo ***REMOVED***
	case PubKeyAlgoRSA, PubKeyAlgoRSAEncryptOnly:
		e.encryptedMPI1.bytes, e.encryptedMPI1.bitLength, err = readMPI(r)
	case PubKeyAlgoElGamal:
		e.encryptedMPI1.bytes, e.encryptedMPI1.bitLength, err = readMPI(r)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		e.encryptedMPI2.bytes, e.encryptedMPI2.bitLength, err = readMPI(r)
	***REMOVED***
	_, err = consumeAll(r)
	return
***REMOVED***

func checksumKeyMaterial(key []byte) uint16 ***REMOVED***
	var checksum uint16
	for _, v := range key ***REMOVED***
		checksum += uint16(v)
	***REMOVED***
	return checksum
***REMOVED***

// Decrypt decrypts an encrypted session key with the given private key. The
// private key must have been decrypted first.
// If config is nil, sensible defaults will be used.
func (e *EncryptedKey) Decrypt(priv *PrivateKey, config *Config) error ***REMOVED***
	var err error
	var b []byte

	// TODO(agl): use session key decryption routines here to avoid
	// padding oracle attacks.
	switch priv.PubKeyAlgo ***REMOVED***
	case PubKeyAlgoRSA, PubKeyAlgoRSAEncryptOnly:
		b, err = rsa.DecryptPKCS1v15(config.Random(), priv.PrivateKey.(*rsa.PrivateKey), e.encryptedMPI1.bytes)
	case PubKeyAlgoElGamal:
		c1 := new(big.Int).SetBytes(e.encryptedMPI1.bytes)
		c2 := new(big.Int).SetBytes(e.encryptedMPI2.bytes)
		b, err = elgamal.Decrypt(priv.PrivateKey.(*elgamal.PrivateKey), c1, c2)
	default:
		err = errors.InvalidArgumentError("cannot decrypted encrypted session key with private key of type " + strconv.Itoa(int(priv.PubKeyAlgo)))
	***REMOVED***

	if err != nil ***REMOVED***
		return err
	***REMOVED***

	e.CipherFunc = CipherFunction(b[0])
	e.Key = b[1 : len(b)-2]
	expectedChecksum := uint16(b[len(b)-2])<<8 | uint16(b[len(b)-1])
	checksum := checksumKeyMaterial(e.Key)
	if checksum != expectedChecksum ***REMOVED***
		return errors.StructuralError("EncryptedKey checksum incorrect")
	***REMOVED***

	return nil
***REMOVED***

// Serialize writes the encrypted key packet, e, to w.
func (e *EncryptedKey) Serialize(w io.Writer) error ***REMOVED***
	var mpiLen int
	switch e.Algo ***REMOVED***
	case PubKeyAlgoRSA, PubKeyAlgoRSAEncryptOnly:
		mpiLen = 2 + len(e.encryptedMPI1.bytes)
	case PubKeyAlgoElGamal:
		mpiLen = 2 + len(e.encryptedMPI1.bytes) + 2 + len(e.encryptedMPI2.bytes)
	default:
		return errors.InvalidArgumentError("don't know how to serialize encrypted key type " + strconv.Itoa(int(e.Algo)))
	***REMOVED***

	serializeHeader(w, packetTypeEncryptedKey, 1 /* version */ +8 /* key id */ +1 /* algo */ +mpiLen)

	w.Write([]byte***REMOVED***encryptedKeyVersion***REMOVED***)
	binary.Write(w, binary.BigEndian, e.KeyId)
	w.Write([]byte***REMOVED***byte(e.Algo)***REMOVED***)

	switch e.Algo ***REMOVED***
	case PubKeyAlgoRSA, PubKeyAlgoRSAEncryptOnly:
		writeMPIs(w, e.encryptedMPI1)
	case PubKeyAlgoElGamal:
		writeMPIs(w, e.encryptedMPI1, e.encryptedMPI2)
	default:
		panic("internal error")
	***REMOVED***

	return nil
***REMOVED***

// SerializeEncryptedKey serializes an encrypted key packet to w that contains
// key, encrypted to pub.
// If config is nil, sensible defaults will be used.
func SerializeEncryptedKey(w io.Writer, pub *PublicKey, cipherFunc CipherFunction, key []byte, config *Config) error ***REMOVED***
	var buf [10]byte
	buf[0] = encryptedKeyVersion
	binary.BigEndian.PutUint64(buf[1:9], pub.KeyId)
	buf[9] = byte(pub.PubKeyAlgo)

	keyBlock := make([]byte, 1 /* cipher type */ +len(key)+2 /* checksum */)
	keyBlock[0] = byte(cipherFunc)
	copy(keyBlock[1:], key)
	checksum := checksumKeyMaterial(key)
	keyBlock[1+len(key)] = byte(checksum >> 8)
	keyBlock[1+len(key)+1] = byte(checksum)

	switch pub.PubKeyAlgo ***REMOVED***
	case PubKeyAlgoRSA, PubKeyAlgoRSAEncryptOnly:
		return serializeEncryptedKeyRSA(w, config.Random(), buf, pub.PublicKey.(*rsa.PublicKey), keyBlock)
	case PubKeyAlgoElGamal:
		return serializeEncryptedKeyElGamal(w, config.Random(), buf, pub.PublicKey.(*elgamal.PublicKey), keyBlock)
	case PubKeyAlgoDSA, PubKeyAlgoRSASignOnly:
		return errors.InvalidArgumentError("cannot encrypt to public key of type " + strconv.Itoa(int(pub.PubKeyAlgo)))
	***REMOVED***

	return errors.UnsupportedError("encrypting a key to public key of type " + strconv.Itoa(int(pub.PubKeyAlgo)))
***REMOVED***

func serializeEncryptedKeyRSA(w io.Writer, rand io.Reader, header [10]byte, pub *rsa.PublicKey, keyBlock []byte) error ***REMOVED***
	cipherText, err := rsa.EncryptPKCS1v15(rand, pub, keyBlock)
	if err != nil ***REMOVED***
		return errors.InvalidArgumentError("RSA encryption failed: " + err.Error())
	***REMOVED***

	packetLen := 10 /* header length */ + 2 /* mpi size */ + len(cipherText)

	err = serializeHeader(w, packetTypeEncryptedKey, packetLen)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = w.Write(header[:])
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return writeMPI(w, 8*uint16(len(cipherText)), cipherText)
***REMOVED***

func serializeEncryptedKeyElGamal(w io.Writer, rand io.Reader, header [10]byte, pub *elgamal.PublicKey, keyBlock []byte) error ***REMOVED***
	c1, c2, err := elgamal.Encrypt(rand, pub, keyBlock)
	if err != nil ***REMOVED***
		return errors.InvalidArgumentError("ElGamal encryption failed: " + err.Error())
	***REMOVED***

	packetLen := 10 /* header length */
	packetLen += 2 /* mpi size */ + (c1.BitLen()+7)/8
	packetLen += 2 /* mpi size */ + (c2.BitLen()+7)/8

	err = serializeHeader(w, packetTypeEncryptedKey, packetLen)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = w.Write(header[:])
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = writeBig(w, c1)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return writeBig(w, c2)
***REMOVED***
