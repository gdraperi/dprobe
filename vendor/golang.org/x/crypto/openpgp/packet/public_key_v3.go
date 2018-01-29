// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"crypto"
	"crypto/md5"
	"crypto/rsa"
	"encoding/binary"
	"fmt"
	"hash"
	"io"
	"math/big"
	"strconv"
	"time"

	"golang.org/x/crypto/openpgp/errors"
)

// PublicKeyV3 represents older, version 3 public keys. These keys are less secure and
// should not be used for signing or encrypting. They are supported here only for
// parsing version 3 key material and validating signatures.
// See RFC 4880, section 5.5.2.
type PublicKeyV3 struct ***REMOVED***
	CreationTime time.Time
	DaysToExpire uint16
	PubKeyAlgo   PublicKeyAlgorithm
	PublicKey    *rsa.PublicKey
	Fingerprint  [16]byte
	KeyId        uint64
	IsSubkey     bool

	n, e parsedMPI
***REMOVED***

// newRSAPublicKeyV3 returns a PublicKey that wraps the given rsa.PublicKey.
// Included here for testing purposes only. RFC 4880, section 5.5.2:
// "an implementation MUST NOT generate a V3 key, but MAY accept it."
func newRSAPublicKeyV3(creationTime time.Time, pub *rsa.PublicKey) *PublicKeyV3 ***REMOVED***
	pk := &PublicKeyV3***REMOVED***
		CreationTime: creationTime,
		PublicKey:    pub,
		n:            fromBig(pub.N),
		e:            fromBig(big.NewInt(int64(pub.E))),
	***REMOVED***

	pk.setFingerPrintAndKeyId()
	return pk
***REMOVED***

func (pk *PublicKeyV3) parse(r io.Reader) (err error) ***REMOVED***
	// RFC 4880, section 5.5.2
	var buf [8]byte
	if _, err = readFull(r, buf[:]); err != nil ***REMOVED***
		return
	***REMOVED***
	if buf[0] < 2 || buf[0] > 3 ***REMOVED***
		return errors.UnsupportedError("public key version")
	***REMOVED***
	pk.CreationTime = time.Unix(int64(uint32(buf[1])<<24|uint32(buf[2])<<16|uint32(buf[3])<<8|uint32(buf[4])), 0)
	pk.DaysToExpire = binary.BigEndian.Uint16(buf[5:7])
	pk.PubKeyAlgo = PublicKeyAlgorithm(buf[7])
	switch pk.PubKeyAlgo ***REMOVED***
	case PubKeyAlgoRSA, PubKeyAlgoRSAEncryptOnly, PubKeyAlgoRSASignOnly:
		err = pk.parseRSA(r)
	default:
		err = errors.UnsupportedError("public key type: " + strconv.Itoa(int(pk.PubKeyAlgo)))
	***REMOVED***
	if err != nil ***REMOVED***
		return
	***REMOVED***

	pk.setFingerPrintAndKeyId()
	return
***REMOVED***

func (pk *PublicKeyV3) setFingerPrintAndKeyId() ***REMOVED***
	// RFC 4880, section 12.2
	fingerPrint := md5.New()
	fingerPrint.Write(pk.n.bytes)
	fingerPrint.Write(pk.e.bytes)
	fingerPrint.Sum(pk.Fingerprint[:0])
	pk.KeyId = binary.BigEndian.Uint64(pk.n.bytes[len(pk.n.bytes)-8:])
***REMOVED***

// parseRSA parses RSA public key material from the given Reader. See RFC 4880,
// section 5.5.2.
func (pk *PublicKeyV3) parseRSA(r io.Reader) (err error) ***REMOVED***
	if pk.n.bytes, pk.n.bitLength, err = readMPI(r); err != nil ***REMOVED***
		return
	***REMOVED***
	if pk.e.bytes, pk.e.bitLength, err = readMPI(r); err != nil ***REMOVED***
		return
	***REMOVED***

	// RFC 4880 Section 12.2 requires the low 8 bytes of the
	// modulus to form the key id.
	if len(pk.n.bytes) < 8 ***REMOVED***
		return errors.StructuralError("v3 public key modulus is too short")
	***REMOVED***
	if len(pk.e.bytes) > 3 ***REMOVED***
		err = errors.UnsupportedError("large public exponent")
		return
	***REMOVED***
	rsa := &rsa.PublicKey***REMOVED***N: new(big.Int).SetBytes(pk.n.bytes)***REMOVED***
	for i := 0; i < len(pk.e.bytes); i++ ***REMOVED***
		rsa.E <<= 8
		rsa.E |= int(pk.e.bytes[i])
	***REMOVED***
	pk.PublicKey = rsa
	return
***REMOVED***

// SerializeSignaturePrefix writes the prefix for this public key to the given Writer.
// The prefix is used when calculating a signature over this public key. See
// RFC 4880, section 5.2.4.
func (pk *PublicKeyV3) SerializeSignaturePrefix(w io.Writer) ***REMOVED***
	var pLength uint16
	switch pk.PubKeyAlgo ***REMOVED***
	case PubKeyAlgoRSA, PubKeyAlgoRSAEncryptOnly, PubKeyAlgoRSASignOnly:
		pLength += 2 + uint16(len(pk.n.bytes))
		pLength += 2 + uint16(len(pk.e.bytes))
	default:
		panic("unknown public key algorithm")
	***REMOVED***
	pLength += 6
	w.Write([]byte***REMOVED***0x99, byte(pLength >> 8), byte(pLength)***REMOVED***)
	return
***REMOVED***

func (pk *PublicKeyV3) Serialize(w io.Writer) (err error) ***REMOVED***
	length := 8 // 8 byte header

	switch pk.PubKeyAlgo ***REMOVED***
	case PubKeyAlgoRSA, PubKeyAlgoRSAEncryptOnly, PubKeyAlgoRSASignOnly:
		length += 2 + len(pk.n.bytes)
		length += 2 + len(pk.e.bytes)
	default:
		panic("unknown public key algorithm")
	***REMOVED***

	packetType := packetTypePublicKey
	if pk.IsSubkey ***REMOVED***
		packetType = packetTypePublicSubkey
	***REMOVED***
	if err = serializeHeader(w, packetType, length); err != nil ***REMOVED***
		return
	***REMOVED***
	return pk.serializeWithoutHeaders(w)
***REMOVED***

// serializeWithoutHeaders marshals the PublicKey to w in the form of an
// OpenPGP public key packet, not including the packet header.
func (pk *PublicKeyV3) serializeWithoutHeaders(w io.Writer) (err error) ***REMOVED***
	var buf [8]byte
	// Version 3
	buf[0] = 3
	// Creation time
	t := uint32(pk.CreationTime.Unix())
	buf[1] = byte(t >> 24)
	buf[2] = byte(t >> 16)
	buf[3] = byte(t >> 8)
	buf[4] = byte(t)
	// Days to expire
	buf[5] = byte(pk.DaysToExpire >> 8)
	buf[6] = byte(pk.DaysToExpire)
	// Public key algorithm
	buf[7] = byte(pk.PubKeyAlgo)

	if _, err = w.Write(buf[:]); err != nil ***REMOVED***
		return
	***REMOVED***

	switch pk.PubKeyAlgo ***REMOVED***
	case PubKeyAlgoRSA, PubKeyAlgoRSAEncryptOnly, PubKeyAlgoRSASignOnly:
		return writeMPIs(w, pk.n, pk.e)
	***REMOVED***
	return errors.InvalidArgumentError("bad public-key algorithm")
***REMOVED***

// CanSign returns true iff this public key can generate signatures
func (pk *PublicKeyV3) CanSign() bool ***REMOVED***
	return pk.PubKeyAlgo != PubKeyAlgoRSAEncryptOnly
***REMOVED***

// VerifySignatureV3 returns nil iff sig is a valid signature, made by this
// public key, of the data hashed into signed. signed is mutated by this call.
func (pk *PublicKeyV3) VerifySignatureV3(signed hash.Hash, sig *SignatureV3) (err error) ***REMOVED***
	if !pk.CanSign() ***REMOVED***
		return errors.InvalidArgumentError("public key cannot generate signatures")
	***REMOVED***

	suffix := make([]byte, 5)
	suffix[0] = byte(sig.SigType)
	binary.BigEndian.PutUint32(suffix[1:], uint32(sig.CreationTime.Unix()))
	signed.Write(suffix)
	hashBytes := signed.Sum(nil)

	if hashBytes[0] != sig.HashTag[0] || hashBytes[1] != sig.HashTag[1] ***REMOVED***
		return errors.SignatureError("hash tag doesn't match")
	***REMOVED***

	if pk.PubKeyAlgo != sig.PubKeyAlgo ***REMOVED***
		return errors.InvalidArgumentError("public key and signature use different algorithms")
	***REMOVED***

	switch pk.PubKeyAlgo ***REMOVED***
	case PubKeyAlgoRSA, PubKeyAlgoRSASignOnly:
		if err = rsa.VerifyPKCS1v15(pk.PublicKey, sig.Hash, hashBytes, sig.RSASignature.bytes); err != nil ***REMOVED***
			return errors.SignatureError("RSA verification failure")
		***REMOVED***
		return
	default:
		// V3 public keys only support RSA.
		panic("shouldn't happen")
	***REMOVED***
***REMOVED***

// VerifyUserIdSignatureV3 returns nil iff sig is a valid signature, made by this
// public key, that id is the identity of pub.
func (pk *PublicKeyV3) VerifyUserIdSignatureV3(id string, pub *PublicKeyV3, sig *SignatureV3) (err error) ***REMOVED***
	h, err := userIdSignatureV3Hash(id, pk, sig.Hash)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return pk.VerifySignatureV3(h, sig)
***REMOVED***

// VerifyKeySignatureV3 returns nil iff sig is a valid signature, made by this
// public key, of signed.
func (pk *PublicKeyV3) VerifyKeySignatureV3(signed *PublicKeyV3, sig *SignatureV3) (err error) ***REMOVED***
	h, err := keySignatureHash(pk, signed, sig.Hash)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return pk.VerifySignatureV3(h, sig)
***REMOVED***

// userIdSignatureV3Hash returns a Hash of the message that needs to be signed
// to assert that pk is a valid key for id.
func userIdSignatureV3Hash(id string, pk signingKey, hfn crypto.Hash) (h hash.Hash, err error) ***REMOVED***
	if !hfn.Available() ***REMOVED***
		return nil, errors.UnsupportedError("hash function")
	***REMOVED***
	h = hfn.New()

	// RFC 4880, section 5.2.4
	pk.SerializeSignaturePrefix(h)
	pk.serializeWithoutHeaders(h)

	h.Write([]byte(id))

	return
***REMOVED***

// KeyIdString returns the public key's fingerprint in capital hex
// (e.g. "6C7EE1B8621CC013").
func (pk *PublicKeyV3) KeyIdString() string ***REMOVED***
	return fmt.Sprintf("%X", pk.KeyId)
***REMOVED***

// KeyIdShortString returns the short form of public key's fingerprint
// in capital hex, as shown by gpg --list-keys (e.g. "621CC013").
func (pk *PublicKeyV3) KeyIdShortString() string ***REMOVED***
	return fmt.Sprintf("%X", pk.KeyId&0xFFFFFFFF)
***REMOVED***

// BitLength returns the bit length for the given public key.
func (pk *PublicKeyV3) BitLength() (bitLength uint16, err error) ***REMOVED***
	switch pk.PubKeyAlgo ***REMOVED***
	case PubKeyAlgoRSA, PubKeyAlgoRSAEncryptOnly, PubKeyAlgoRSASignOnly:
		bitLength = pk.n.bitLength
	default:
		err = errors.InvalidArgumentError("bad public-key algorithm")
	***REMOVED***
	return
***REMOVED***
