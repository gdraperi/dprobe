// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"bytes"
	"crypto"
	"crypto/cipher"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha1"
	"io"
	"io/ioutil"
	"math/big"
	"strconv"
	"time"

	"golang.org/x/crypto/openpgp/elgamal"
	"golang.org/x/crypto/openpgp/errors"
	"golang.org/x/crypto/openpgp/s2k"
)

// PrivateKey represents a possibly encrypted private key. See RFC 4880,
// section 5.5.3.
type PrivateKey struct ***REMOVED***
	PublicKey
	Encrypted     bool // if true then the private key is unavailable until Decrypt has been called.
	encryptedData []byte
	cipher        CipherFunction
	s2k           func(out, in []byte)
	PrivateKey    interface***REMOVED******REMOVED*** // An ****REMOVED***rsa|dsa|ecdsa***REMOVED***.PrivateKey or a crypto.Signer.
	sha1Checksum  bool
	iv            []byte
***REMOVED***

func NewRSAPrivateKey(currentTime time.Time, priv *rsa.PrivateKey) *PrivateKey ***REMOVED***
	pk := new(PrivateKey)
	pk.PublicKey = *NewRSAPublicKey(currentTime, &priv.PublicKey)
	pk.PrivateKey = priv
	return pk
***REMOVED***

func NewDSAPrivateKey(currentTime time.Time, priv *dsa.PrivateKey) *PrivateKey ***REMOVED***
	pk := new(PrivateKey)
	pk.PublicKey = *NewDSAPublicKey(currentTime, &priv.PublicKey)
	pk.PrivateKey = priv
	return pk
***REMOVED***

func NewElGamalPrivateKey(currentTime time.Time, priv *elgamal.PrivateKey) *PrivateKey ***REMOVED***
	pk := new(PrivateKey)
	pk.PublicKey = *NewElGamalPublicKey(currentTime, &priv.PublicKey)
	pk.PrivateKey = priv
	return pk
***REMOVED***

func NewECDSAPrivateKey(currentTime time.Time, priv *ecdsa.PrivateKey) *PrivateKey ***REMOVED***
	pk := new(PrivateKey)
	pk.PublicKey = *NewECDSAPublicKey(currentTime, &priv.PublicKey)
	pk.PrivateKey = priv
	return pk
***REMOVED***

// NewSignerPrivateKey creates a sign-only PrivateKey from a crypto.Signer that
// implements RSA or ECDSA.
func NewSignerPrivateKey(currentTime time.Time, signer crypto.Signer) *PrivateKey ***REMOVED***
	pk := new(PrivateKey)
	switch pubkey := signer.Public().(type) ***REMOVED***
	case rsa.PublicKey:
		pk.PublicKey = *NewRSAPublicKey(currentTime, &pubkey)
		pk.PubKeyAlgo = PubKeyAlgoRSASignOnly
	case ecdsa.PublicKey:
		pk.PublicKey = *NewECDSAPublicKey(currentTime, &pubkey)
	default:
		panic("openpgp: unknown crypto.Signer type in NewSignerPrivateKey")
	***REMOVED***
	pk.PrivateKey = signer
	return pk
***REMOVED***

func (pk *PrivateKey) parse(r io.Reader) (err error) ***REMOVED***
	err = (&pk.PublicKey).parse(r)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	var buf [1]byte
	_, err = readFull(r, buf[:])
	if err != nil ***REMOVED***
		return
	***REMOVED***

	s2kType := buf[0]

	switch s2kType ***REMOVED***
	case 0:
		pk.s2k = nil
		pk.Encrypted = false
	case 254, 255:
		_, err = readFull(r, buf[:])
		if err != nil ***REMOVED***
			return
		***REMOVED***
		pk.cipher = CipherFunction(buf[0])
		pk.Encrypted = true
		pk.s2k, err = s2k.Parse(r)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		if s2kType == 254 ***REMOVED***
			pk.sha1Checksum = true
		***REMOVED***
	default:
		return errors.UnsupportedError("deprecated s2k function in private key")
	***REMOVED***

	if pk.Encrypted ***REMOVED***
		blockSize := pk.cipher.blockSize()
		if blockSize == 0 ***REMOVED***
			return errors.UnsupportedError("unsupported cipher in private key: " + strconv.Itoa(int(pk.cipher)))
		***REMOVED***
		pk.iv = make([]byte, blockSize)
		_, err = readFull(r, pk.iv)
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	pk.encryptedData, err = ioutil.ReadAll(r)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	if !pk.Encrypted ***REMOVED***
		return pk.parsePrivateKey(pk.encryptedData)
	***REMOVED***

	return
***REMOVED***

func mod64kHash(d []byte) uint16 ***REMOVED***
	var h uint16
	for _, b := range d ***REMOVED***
		h += uint16(b)
	***REMOVED***
	return h
***REMOVED***

func (pk *PrivateKey) Serialize(w io.Writer) (err error) ***REMOVED***
	// TODO(agl): support encrypted private keys
	buf := bytes.NewBuffer(nil)
	err = pk.PublicKey.serializeWithoutHeaders(buf)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	buf.WriteByte(0 /* no encryption */)

	privateKeyBuf := bytes.NewBuffer(nil)

	switch priv := pk.PrivateKey.(type) ***REMOVED***
	case *rsa.PrivateKey:
		err = serializeRSAPrivateKey(privateKeyBuf, priv)
	case *dsa.PrivateKey:
		err = serializeDSAPrivateKey(privateKeyBuf, priv)
	case *elgamal.PrivateKey:
		err = serializeElGamalPrivateKey(privateKeyBuf, priv)
	case *ecdsa.PrivateKey:
		err = serializeECDSAPrivateKey(privateKeyBuf, priv)
	default:
		err = errors.InvalidArgumentError("unknown private key type")
	***REMOVED***
	if err != nil ***REMOVED***
		return
	***REMOVED***

	ptype := packetTypePrivateKey
	contents := buf.Bytes()
	privateKeyBytes := privateKeyBuf.Bytes()
	if pk.IsSubkey ***REMOVED***
		ptype = packetTypePrivateSubkey
	***REMOVED***
	err = serializeHeader(w, ptype, len(contents)+len(privateKeyBytes)+2)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	_, err = w.Write(contents)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	_, err = w.Write(privateKeyBytes)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	checksum := mod64kHash(privateKeyBytes)
	var checksumBytes [2]byte
	checksumBytes[0] = byte(checksum >> 8)
	checksumBytes[1] = byte(checksum)
	_, err = w.Write(checksumBytes[:])

	return
***REMOVED***

func serializeRSAPrivateKey(w io.Writer, priv *rsa.PrivateKey) error ***REMOVED***
	err := writeBig(w, priv.D)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = writeBig(w, priv.Primes[1])
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = writeBig(w, priv.Primes[0])
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return writeBig(w, priv.Precomputed.Qinv)
***REMOVED***

func serializeDSAPrivateKey(w io.Writer, priv *dsa.PrivateKey) error ***REMOVED***
	return writeBig(w, priv.X)
***REMOVED***

func serializeElGamalPrivateKey(w io.Writer, priv *elgamal.PrivateKey) error ***REMOVED***
	return writeBig(w, priv.X)
***REMOVED***

func serializeECDSAPrivateKey(w io.Writer, priv *ecdsa.PrivateKey) error ***REMOVED***
	return writeBig(w, priv.D)
***REMOVED***

// Decrypt decrypts an encrypted private key using a passphrase.
func (pk *PrivateKey) Decrypt(passphrase []byte) error ***REMOVED***
	if !pk.Encrypted ***REMOVED***
		return nil
	***REMOVED***

	key := make([]byte, pk.cipher.KeySize())
	pk.s2k(key, passphrase)
	block := pk.cipher.new(key)
	cfb := cipher.NewCFBDecrypter(block, pk.iv)

	data := make([]byte, len(pk.encryptedData))
	cfb.XORKeyStream(data, pk.encryptedData)

	if pk.sha1Checksum ***REMOVED***
		if len(data) < sha1.Size ***REMOVED***
			return errors.StructuralError("truncated private key data")
		***REMOVED***
		h := sha1.New()
		h.Write(data[:len(data)-sha1.Size])
		sum := h.Sum(nil)
		if !bytes.Equal(sum, data[len(data)-sha1.Size:]) ***REMOVED***
			return errors.StructuralError("private key checksum failure")
		***REMOVED***
		data = data[:len(data)-sha1.Size]
	***REMOVED*** else ***REMOVED***
		if len(data) < 2 ***REMOVED***
			return errors.StructuralError("truncated private key data")
		***REMOVED***
		var sum uint16
		for i := 0; i < len(data)-2; i++ ***REMOVED***
			sum += uint16(data[i])
		***REMOVED***
		if data[len(data)-2] != uint8(sum>>8) ||
			data[len(data)-1] != uint8(sum) ***REMOVED***
			return errors.StructuralError("private key checksum failure")
		***REMOVED***
		data = data[:len(data)-2]
	***REMOVED***

	return pk.parsePrivateKey(data)
***REMOVED***

func (pk *PrivateKey) parsePrivateKey(data []byte) (err error) ***REMOVED***
	switch pk.PublicKey.PubKeyAlgo ***REMOVED***
	case PubKeyAlgoRSA, PubKeyAlgoRSASignOnly, PubKeyAlgoRSAEncryptOnly:
		return pk.parseRSAPrivateKey(data)
	case PubKeyAlgoDSA:
		return pk.parseDSAPrivateKey(data)
	case PubKeyAlgoElGamal:
		return pk.parseElGamalPrivateKey(data)
	case PubKeyAlgoECDSA:
		return pk.parseECDSAPrivateKey(data)
	***REMOVED***
	panic("impossible")
***REMOVED***

func (pk *PrivateKey) parseRSAPrivateKey(data []byte) (err error) ***REMOVED***
	rsaPub := pk.PublicKey.PublicKey.(*rsa.PublicKey)
	rsaPriv := new(rsa.PrivateKey)
	rsaPriv.PublicKey = *rsaPub

	buf := bytes.NewBuffer(data)
	d, _, err := readMPI(buf)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	p, _, err := readMPI(buf)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	q, _, err := readMPI(buf)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	rsaPriv.D = new(big.Int).SetBytes(d)
	rsaPriv.Primes = make([]*big.Int, 2)
	rsaPriv.Primes[0] = new(big.Int).SetBytes(p)
	rsaPriv.Primes[1] = new(big.Int).SetBytes(q)
	if err := rsaPriv.Validate(); err != nil ***REMOVED***
		return err
	***REMOVED***
	rsaPriv.Precompute()
	pk.PrivateKey = rsaPriv
	pk.Encrypted = false
	pk.encryptedData = nil

	return nil
***REMOVED***

func (pk *PrivateKey) parseDSAPrivateKey(data []byte) (err error) ***REMOVED***
	dsaPub := pk.PublicKey.PublicKey.(*dsa.PublicKey)
	dsaPriv := new(dsa.PrivateKey)
	dsaPriv.PublicKey = *dsaPub

	buf := bytes.NewBuffer(data)
	x, _, err := readMPI(buf)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	dsaPriv.X = new(big.Int).SetBytes(x)
	pk.PrivateKey = dsaPriv
	pk.Encrypted = false
	pk.encryptedData = nil

	return nil
***REMOVED***

func (pk *PrivateKey) parseElGamalPrivateKey(data []byte) (err error) ***REMOVED***
	pub := pk.PublicKey.PublicKey.(*elgamal.PublicKey)
	priv := new(elgamal.PrivateKey)
	priv.PublicKey = *pub

	buf := bytes.NewBuffer(data)
	x, _, err := readMPI(buf)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	priv.X = new(big.Int).SetBytes(x)
	pk.PrivateKey = priv
	pk.Encrypted = false
	pk.encryptedData = nil

	return nil
***REMOVED***

func (pk *PrivateKey) parseECDSAPrivateKey(data []byte) (err error) ***REMOVED***
	ecdsaPub := pk.PublicKey.PublicKey.(*ecdsa.PublicKey)

	buf := bytes.NewBuffer(data)
	d, _, err := readMPI(buf)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	pk.PrivateKey = &ecdsa.PrivateKey***REMOVED***
		PublicKey: *ecdsaPub,
		D:         new(big.Int).SetBytes(d),
	***REMOVED***
	pk.Encrypted = false
	pk.encryptedData = nil

	return nil
***REMOVED***
