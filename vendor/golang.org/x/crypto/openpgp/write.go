// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package openpgp

import (
	"crypto"
	"hash"
	"io"
	"strconv"
	"time"

	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/errors"
	"golang.org/x/crypto/openpgp/packet"
	"golang.org/x/crypto/openpgp/s2k"
)

// DetachSign signs message with the private key from signer (which must
// already have been decrypted) and writes the signature to w.
// If config is nil, sensible defaults will be used.
func DetachSign(w io.Writer, signer *Entity, message io.Reader, config *packet.Config) error ***REMOVED***
	return detachSign(w, signer, message, packet.SigTypeBinary, config)
***REMOVED***

// ArmoredDetachSign signs message with the private key from signer (which
// must already have been decrypted) and writes an armored signature to w.
// If config is nil, sensible defaults will be used.
func ArmoredDetachSign(w io.Writer, signer *Entity, message io.Reader, config *packet.Config) (err error) ***REMOVED***
	return armoredDetachSign(w, signer, message, packet.SigTypeBinary, config)
***REMOVED***

// DetachSignText signs message (after canonicalising the line endings) with
// the private key from signer (which must already have been decrypted) and
// writes the signature to w.
// If config is nil, sensible defaults will be used.
func DetachSignText(w io.Writer, signer *Entity, message io.Reader, config *packet.Config) error ***REMOVED***
	return detachSign(w, signer, message, packet.SigTypeText, config)
***REMOVED***

// ArmoredDetachSignText signs message (after canonicalising the line endings)
// with the private key from signer (which must already have been decrypted)
// and writes an armored signature to w.
// If config is nil, sensible defaults will be used.
func ArmoredDetachSignText(w io.Writer, signer *Entity, message io.Reader, config *packet.Config) error ***REMOVED***
	return armoredDetachSign(w, signer, message, packet.SigTypeText, config)
***REMOVED***

func armoredDetachSign(w io.Writer, signer *Entity, message io.Reader, sigType packet.SignatureType, config *packet.Config) (err error) ***REMOVED***
	out, err := armor.Encode(w, SignatureType, nil)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	err = detachSign(out, signer, message, sigType, config)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	return out.Close()
***REMOVED***

func detachSign(w io.Writer, signer *Entity, message io.Reader, sigType packet.SignatureType, config *packet.Config) (err error) ***REMOVED***
	if signer.PrivateKey == nil ***REMOVED***
		return errors.InvalidArgumentError("signing key doesn't have a private key")
	***REMOVED***
	if signer.PrivateKey.Encrypted ***REMOVED***
		return errors.InvalidArgumentError("signing key is encrypted")
	***REMOVED***

	sig := new(packet.Signature)
	sig.SigType = sigType
	sig.PubKeyAlgo = signer.PrivateKey.PubKeyAlgo
	sig.Hash = config.Hash()
	sig.CreationTime = config.Now()
	sig.IssuerKeyId = &signer.PrivateKey.KeyId

	h, wrappedHash, err := hashForSignature(sig.Hash, sig.SigType)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	io.Copy(wrappedHash, message)

	err = sig.Sign(h, signer.PrivateKey, config)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	return sig.Serialize(w)
***REMOVED***

// FileHints contains metadata about encrypted files. This metadata is, itself,
// encrypted.
type FileHints struct ***REMOVED***
	// IsBinary can be set to hint that the contents are binary data.
	IsBinary bool
	// FileName hints at the name of the file that should be written. It's
	// truncated to 255 bytes if longer. It may be empty to suggest that the
	// file should not be written to disk. It may be equal to "_CONSOLE" to
	// suggest the data should not be written to disk.
	FileName string
	// ModTime contains the modification time of the file, or the zero time if not applicable.
	ModTime time.Time
***REMOVED***

// SymmetricallyEncrypt acts like gpg -c: it encrypts a file with a passphrase.
// The resulting WriteCloser must be closed after the contents of the file have
// been written.
// If config is nil, sensible defaults will be used.
func SymmetricallyEncrypt(ciphertext io.Writer, passphrase []byte, hints *FileHints, config *packet.Config) (plaintext io.WriteCloser, err error) ***REMOVED***
	if hints == nil ***REMOVED***
		hints = &FileHints***REMOVED******REMOVED***
	***REMOVED***

	key, err := packet.SerializeSymmetricKeyEncrypted(ciphertext, passphrase, config)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	w, err := packet.SerializeSymmetricallyEncrypted(ciphertext, config.Cipher(), key, config)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	literaldata := w
	if algo := config.Compression(); algo != packet.CompressionNone ***REMOVED***
		var compConfig *packet.CompressionConfig
		if config != nil ***REMOVED***
			compConfig = config.CompressionConfig
		***REMOVED***
		literaldata, err = packet.SerializeCompressed(w, algo, compConfig)
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	var epochSeconds uint32
	if !hints.ModTime.IsZero() ***REMOVED***
		epochSeconds = uint32(hints.ModTime.Unix())
	***REMOVED***
	return packet.SerializeLiteral(literaldata, hints.IsBinary, hints.FileName, epochSeconds)
***REMOVED***

// intersectPreferences mutates and returns a prefix of a that contains only
// the values in the intersection of a and b. The order of a is preserved.
func intersectPreferences(a []uint8, b []uint8) (intersection []uint8) ***REMOVED***
	var j int
	for _, v := range a ***REMOVED***
		for _, v2 := range b ***REMOVED***
			if v == v2 ***REMOVED***
				a[j] = v
				j++
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return a[:j]
***REMOVED***

func hashToHashId(h crypto.Hash) uint8 ***REMOVED***
	v, ok := s2k.HashToHashId(h)
	if !ok ***REMOVED***
		panic("tried to convert unknown hash")
	***REMOVED***
	return v
***REMOVED***

// Encrypt encrypts a message to a number of recipients and, optionally, signs
// it. hints contains optional information, that is also encrypted, that aids
// the recipients in processing the message. The resulting WriteCloser must
// be closed after the contents of the file have been written.
// If config is nil, sensible defaults will be used.
func Encrypt(ciphertext io.Writer, to []*Entity, signed *Entity, hints *FileHints, config *packet.Config) (plaintext io.WriteCloser, err error) ***REMOVED***
	var signer *packet.PrivateKey
	if signed != nil ***REMOVED***
		signKey, ok := signed.signingKey(config.Now())
		if !ok ***REMOVED***
			return nil, errors.InvalidArgumentError("no valid signing keys")
		***REMOVED***
		signer = signKey.PrivateKey
		if signer == nil ***REMOVED***
			return nil, errors.InvalidArgumentError("no private key in signing key")
		***REMOVED***
		if signer.Encrypted ***REMOVED***
			return nil, errors.InvalidArgumentError("signing key must be decrypted")
		***REMOVED***
	***REMOVED***

	// These are the possible ciphers that we'll use for the message.
	candidateCiphers := []uint8***REMOVED***
		uint8(packet.CipherAES128),
		uint8(packet.CipherAES256),
		uint8(packet.CipherCAST5),
	***REMOVED***
	// These are the possible hash functions that we'll use for the signature.
	candidateHashes := []uint8***REMOVED***
		hashToHashId(crypto.SHA256),
		hashToHashId(crypto.SHA512),
		hashToHashId(crypto.SHA1),
		hashToHashId(crypto.RIPEMD160),
	***REMOVED***
	// In the event that a recipient doesn't specify any supported ciphers
	// or hash functions, these are the ones that we assume that every
	// implementation supports.
	defaultCiphers := candidateCiphers[len(candidateCiphers)-1:]
	defaultHashes := candidateHashes[len(candidateHashes)-1:]

	encryptKeys := make([]Key, len(to))
	for i := range to ***REMOVED***
		var ok bool
		encryptKeys[i], ok = to[i].encryptionKey(config.Now())
		if !ok ***REMOVED***
			return nil, errors.InvalidArgumentError("cannot encrypt a message to key id " + strconv.FormatUint(to[i].PrimaryKey.KeyId, 16) + " because it has no encryption keys")
		***REMOVED***

		sig := to[i].primaryIdentity().SelfSignature

		preferredSymmetric := sig.PreferredSymmetric
		if len(preferredSymmetric) == 0 ***REMOVED***
			preferredSymmetric = defaultCiphers
		***REMOVED***
		preferredHashes := sig.PreferredHash
		if len(preferredHashes) == 0 ***REMOVED***
			preferredHashes = defaultHashes
		***REMOVED***
		candidateCiphers = intersectPreferences(candidateCiphers, preferredSymmetric)
		candidateHashes = intersectPreferences(candidateHashes, preferredHashes)
	***REMOVED***

	if len(candidateCiphers) == 0 || len(candidateHashes) == 0 ***REMOVED***
		return nil, errors.InvalidArgumentError("cannot encrypt because recipient set shares no common algorithms")
	***REMOVED***

	cipher := packet.CipherFunction(candidateCiphers[0])
	// If the cipher specified by config is a candidate, we'll use that.
	configuredCipher := config.Cipher()
	for _, c := range candidateCiphers ***REMOVED***
		cipherFunc := packet.CipherFunction(c)
		if cipherFunc == configuredCipher ***REMOVED***
			cipher = cipherFunc
			break
		***REMOVED***
	***REMOVED***

	var hash crypto.Hash
	for _, hashId := range candidateHashes ***REMOVED***
		if h, ok := s2k.HashIdToHash(hashId); ok && h.Available() ***REMOVED***
			hash = h
			break
		***REMOVED***
	***REMOVED***

	// If the hash specified by config is a candidate, we'll use that.
	if configuredHash := config.Hash(); configuredHash.Available() ***REMOVED***
		for _, hashId := range candidateHashes ***REMOVED***
			if h, ok := s2k.HashIdToHash(hashId); ok && h == configuredHash ***REMOVED***
				hash = h
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if hash == 0 ***REMOVED***
		hashId := candidateHashes[0]
		name, ok := s2k.HashIdToString(hashId)
		if !ok ***REMOVED***
			name = "#" + strconv.Itoa(int(hashId))
		***REMOVED***
		return nil, errors.InvalidArgumentError("cannot encrypt because no candidate hash functions are compiled in. (Wanted " + name + " in this case.)")
	***REMOVED***

	symKey := make([]byte, cipher.KeySize())
	if _, err := io.ReadFull(config.Random(), symKey); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	for _, key := range encryptKeys ***REMOVED***
		if err := packet.SerializeEncryptedKey(ciphertext, key.PublicKey, cipher, symKey, config); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	encryptedData, err := packet.SerializeSymmetricallyEncrypted(ciphertext, cipher, symKey, config)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	if signer != nil ***REMOVED***
		ops := &packet.OnePassSignature***REMOVED***
			SigType:    packet.SigTypeBinary,
			Hash:       hash,
			PubKeyAlgo: signer.PubKeyAlgo,
			KeyId:      signer.KeyId,
			IsLast:     true,
		***REMOVED***
		if err := ops.Serialize(encryptedData); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	if hints == nil ***REMOVED***
		hints = &FileHints***REMOVED******REMOVED***
	***REMOVED***

	w := encryptedData
	if signer != nil ***REMOVED***
		// If we need to write a signature packet after the literal
		// data then we need to stop literalData from closing
		// encryptedData.
		w = noOpCloser***REMOVED***encryptedData***REMOVED***

	***REMOVED***
	var epochSeconds uint32
	if !hints.ModTime.IsZero() ***REMOVED***
		epochSeconds = uint32(hints.ModTime.Unix())
	***REMOVED***
	literalData, err := packet.SerializeLiteral(w, hints.IsBinary, hints.FileName, epochSeconds)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if signer != nil ***REMOVED***
		return signatureWriter***REMOVED***encryptedData, literalData, hash, hash.New(), signer, config***REMOVED***, nil
	***REMOVED***
	return literalData, nil
***REMOVED***

// signatureWriter hashes the contents of a message while passing it along to
// literalData. When closed, it closes literalData, writes a signature packet
// to encryptedData and then also closes encryptedData.
type signatureWriter struct ***REMOVED***
	encryptedData io.WriteCloser
	literalData   io.WriteCloser
	hashType      crypto.Hash
	h             hash.Hash
	signer        *packet.PrivateKey
	config        *packet.Config
***REMOVED***

func (s signatureWriter) Write(data []byte) (int, error) ***REMOVED***
	s.h.Write(data)
	return s.literalData.Write(data)
***REMOVED***

func (s signatureWriter) Close() error ***REMOVED***
	sig := &packet.Signature***REMOVED***
		SigType:      packet.SigTypeBinary,
		PubKeyAlgo:   s.signer.PubKeyAlgo,
		Hash:         s.hashType,
		CreationTime: s.config.Now(),
		IssuerKeyId:  &s.signer.KeyId,
	***REMOVED***

	if err := sig.Sign(s.h, s.signer, s.config); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := s.literalData.Close(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := sig.Serialize(s.encryptedData); err != nil ***REMOVED***
		return err
	***REMOVED***
	return s.encryptedData.Close()
***REMOVED***

// noOpCloser is like an ioutil.NopCloser, but for an io.Writer.
// TODO: we have two of these in OpenPGP packages alone. This probably needs
// to be promoted somewhere more common.
type noOpCloser struct ***REMOVED***
	w io.Writer
***REMOVED***

func (c noOpCloser) Write(data []byte) (n int, err error) ***REMOVED***
	return c.w.Write(data)
***REMOVED***

func (c noOpCloser) Close() error ***REMOVED***
	return nil
***REMOVED***
