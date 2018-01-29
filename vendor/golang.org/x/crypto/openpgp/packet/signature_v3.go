// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"crypto"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"time"

	"golang.org/x/crypto/openpgp/errors"
	"golang.org/x/crypto/openpgp/s2k"
)

// SignatureV3 represents older version 3 signatures. These signatures are less secure
// than version 4 and should not be used to create new signatures. They are included
// here for backwards compatibility to read and validate with older key material.
// See RFC 4880, section 5.2.2.
type SignatureV3 struct ***REMOVED***
	SigType      SignatureType
	CreationTime time.Time
	IssuerKeyId  uint64
	PubKeyAlgo   PublicKeyAlgorithm
	Hash         crypto.Hash
	HashTag      [2]byte

	RSASignature     parsedMPI
	DSASigR, DSASigS parsedMPI
***REMOVED***

func (sig *SignatureV3) parse(r io.Reader) (err error) ***REMOVED***
	// RFC 4880, section 5.2.2
	var buf [8]byte
	if _, err = readFull(r, buf[:1]); err != nil ***REMOVED***
		return
	***REMOVED***
	if buf[0] < 2 || buf[0] > 3 ***REMOVED***
		err = errors.UnsupportedError("signature packet version " + strconv.Itoa(int(buf[0])))
		return
	***REMOVED***
	if _, err = readFull(r, buf[:1]); err != nil ***REMOVED***
		return
	***REMOVED***
	if buf[0] != 5 ***REMOVED***
		err = errors.UnsupportedError(
			"invalid hashed material length " + strconv.Itoa(int(buf[0])))
		return
	***REMOVED***

	// Read hashed material: signature type + creation time
	if _, err = readFull(r, buf[:5]); err != nil ***REMOVED***
		return
	***REMOVED***
	sig.SigType = SignatureType(buf[0])
	t := binary.BigEndian.Uint32(buf[1:5])
	sig.CreationTime = time.Unix(int64(t), 0)

	// Eight-octet Key ID of signer.
	if _, err = readFull(r, buf[:8]); err != nil ***REMOVED***
		return
	***REMOVED***
	sig.IssuerKeyId = binary.BigEndian.Uint64(buf[:])

	// Public-key and hash algorithm
	if _, err = readFull(r, buf[:2]); err != nil ***REMOVED***
		return
	***REMOVED***
	sig.PubKeyAlgo = PublicKeyAlgorithm(buf[0])
	switch sig.PubKeyAlgo ***REMOVED***
	case PubKeyAlgoRSA, PubKeyAlgoRSASignOnly, PubKeyAlgoDSA:
	default:
		err = errors.UnsupportedError("public key algorithm " + strconv.Itoa(int(sig.PubKeyAlgo)))
		return
	***REMOVED***
	var ok bool
	if sig.Hash, ok = s2k.HashIdToHash(buf[1]); !ok ***REMOVED***
		return errors.UnsupportedError("hash function " + strconv.Itoa(int(buf[2])))
	***REMOVED***

	// Two-octet field holding left 16 bits of signed hash value.
	if _, err = readFull(r, sig.HashTag[:2]); err != nil ***REMOVED***
		return
	***REMOVED***

	switch sig.PubKeyAlgo ***REMOVED***
	case PubKeyAlgoRSA, PubKeyAlgoRSASignOnly:
		sig.RSASignature.bytes, sig.RSASignature.bitLength, err = readMPI(r)
	case PubKeyAlgoDSA:
		if sig.DSASigR.bytes, sig.DSASigR.bitLength, err = readMPI(r); err != nil ***REMOVED***
			return
		***REMOVED***
		sig.DSASigS.bytes, sig.DSASigS.bitLength, err = readMPI(r)
	default:
		panic("unreachable")
	***REMOVED***
	return
***REMOVED***

// Serialize marshals sig to w. Sign, SignUserId or SignKey must have been
// called first.
func (sig *SignatureV3) Serialize(w io.Writer) (err error) ***REMOVED***
	buf := make([]byte, 8)

	// Write the sig type and creation time
	buf[0] = byte(sig.SigType)
	binary.BigEndian.PutUint32(buf[1:5], uint32(sig.CreationTime.Unix()))
	if _, err = w.Write(buf[:5]); err != nil ***REMOVED***
		return
	***REMOVED***

	// Write the issuer long key ID
	binary.BigEndian.PutUint64(buf[:8], sig.IssuerKeyId)
	if _, err = w.Write(buf[:8]); err != nil ***REMOVED***
		return
	***REMOVED***

	// Write public key algorithm, hash ID, and hash value
	buf[0] = byte(sig.PubKeyAlgo)
	hashId, ok := s2k.HashToHashId(sig.Hash)
	if !ok ***REMOVED***
		return errors.UnsupportedError(fmt.Sprintf("hash function %v", sig.Hash))
	***REMOVED***
	buf[1] = hashId
	copy(buf[2:4], sig.HashTag[:])
	if _, err = w.Write(buf[:4]); err != nil ***REMOVED***
		return
	***REMOVED***

	if sig.RSASignature.bytes == nil && sig.DSASigR.bytes == nil ***REMOVED***
		return errors.InvalidArgumentError("Signature: need to call Sign, SignUserId or SignKey before Serialize")
	***REMOVED***

	switch sig.PubKeyAlgo ***REMOVED***
	case PubKeyAlgoRSA, PubKeyAlgoRSASignOnly:
		err = writeMPIs(w, sig.RSASignature)
	case PubKeyAlgoDSA:
		err = writeMPIs(w, sig.DSASigR, sig.DSASigS)
	default:
		panic("impossible")
	***REMOVED***
	return
***REMOVED***
