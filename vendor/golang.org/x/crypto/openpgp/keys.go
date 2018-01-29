// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package openpgp

import (
	"crypto/rsa"
	"io"
	"time"

	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/errors"
	"golang.org/x/crypto/openpgp/packet"
)

// PublicKeyType is the armor type for a PGP public key.
var PublicKeyType = "PGP PUBLIC KEY BLOCK"

// PrivateKeyType is the armor type for a PGP private key.
var PrivateKeyType = "PGP PRIVATE KEY BLOCK"

// An Entity represents the components of an OpenPGP key: a primary public key
// (which must be a signing key), one or more identities claimed by that key,
// and zero or more subkeys, which may be encryption keys.
type Entity struct ***REMOVED***
	PrimaryKey  *packet.PublicKey
	PrivateKey  *packet.PrivateKey
	Identities  map[string]*Identity // indexed by Identity.Name
	Revocations []*packet.Signature
	Subkeys     []Subkey
***REMOVED***

// An Identity represents an identity claimed by an Entity and zero or more
// assertions by other entities about that claim.
type Identity struct ***REMOVED***
	Name          string // by convention, has the form "Full Name (comment) <email@example.com>"
	UserId        *packet.UserId
	SelfSignature *packet.Signature
	Signatures    []*packet.Signature
***REMOVED***

// A Subkey is an additional public key in an Entity. Subkeys can be used for
// encryption.
type Subkey struct ***REMOVED***
	PublicKey  *packet.PublicKey
	PrivateKey *packet.PrivateKey
	Sig        *packet.Signature
***REMOVED***

// A Key identifies a specific public key in an Entity. This is either the
// Entity's primary key or a subkey.
type Key struct ***REMOVED***
	Entity        *Entity
	PublicKey     *packet.PublicKey
	PrivateKey    *packet.PrivateKey
	SelfSignature *packet.Signature
***REMOVED***

// A KeyRing provides access to public and private keys.
type KeyRing interface ***REMOVED***
	// KeysById returns the set of keys that have the given key id.
	KeysById(id uint64) []Key
	// KeysByIdAndUsage returns the set of keys with the given id
	// that also meet the key usage given by requiredUsage.
	// The requiredUsage is expressed as the bitwise-OR of
	// packet.KeyFlag* values.
	KeysByIdUsage(id uint64, requiredUsage byte) []Key
	// DecryptionKeys returns all private keys that are valid for
	// decryption.
	DecryptionKeys() []Key
***REMOVED***

// primaryIdentity returns the Identity marked as primary or the first identity
// if none are so marked.
func (e *Entity) primaryIdentity() *Identity ***REMOVED***
	var firstIdentity *Identity
	for _, ident := range e.Identities ***REMOVED***
		if firstIdentity == nil ***REMOVED***
			firstIdentity = ident
		***REMOVED***
		if ident.SelfSignature.IsPrimaryId != nil && *ident.SelfSignature.IsPrimaryId ***REMOVED***
			return ident
		***REMOVED***
	***REMOVED***
	return firstIdentity
***REMOVED***

// encryptionKey returns the best candidate Key for encrypting a message to the
// given Entity.
func (e *Entity) encryptionKey(now time.Time) (Key, bool) ***REMOVED***
	candidateSubkey := -1

	// Iterate the keys to find the newest key
	var maxTime time.Time
	for i, subkey := range e.Subkeys ***REMOVED***
		if subkey.Sig.FlagsValid &&
			subkey.Sig.FlagEncryptCommunications &&
			subkey.PublicKey.PubKeyAlgo.CanEncrypt() &&
			!subkey.Sig.KeyExpired(now) &&
			(maxTime.IsZero() || subkey.Sig.CreationTime.After(maxTime)) ***REMOVED***
			candidateSubkey = i
			maxTime = subkey.Sig.CreationTime
		***REMOVED***
	***REMOVED***

	if candidateSubkey != -1 ***REMOVED***
		subkey := e.Subkeys[candidateSubkey]
		return Key***REMOVED***e, subkey.PublicKey, subkey.PrivateKey, subkey.Sig***REMOVED***, true
	***REMOVED***

	// If we don't have any candidate subkeys for encryption and
	// the primary key doesn't have any usage metadata then we
	// assume that the primary key is ok. Or, if the primary key is
	// marked as ok to encrypt to, then we can obviously use it.
	i := e.primaryIdentity()
	if !i.SelfSignature.FlagsValid || i.SelfSignature.FlagEncryptCommunications &&
		e.PrimaryKey.PubKeyAlgo.CanEncrypt() &&
		!i.SelfSignature.KeyExpired(now) ***REMOVED***
		return Key***REMOVED***e, e.PrimaryKey, e.PrivateKey, i.SelfSignature***REMOVED***, true
	***REMOVED***

	// This Entity appears to be signing only.
	return Key***REMOVED******REMOVED***, false
***REMOVED***

// signingKey return the best candidate Key for signing a message with this
// Entity.
func (e *Entity) signingKey(now time.Time) (Key, bool) ***REMOVED***
	candidateSubkey := -1

	for i, subkey := range e.Subkeys ***REMOVED***
		if subkey.Sig.FlagsValid &&
			subkey.Sig.FlagSign &&
			subkey.PublicKey.PubKeyAlgo.CanSign() &&
			!subkey.Sig.KeyExpired(now) ***REMOVED***
			candidateSubkey = i
			break
		***REMOVED***
	***REMOVED***

	if candidateSubkey != -1 ***REMOVED***
		subkey := e.Subkeys[candidateSubkey]
		return Key***REMOVED***e, subkey.PublicKey, subkey.PrivateKey, subkey.Sig***REMOVED***, true
	***REMOVED***

	// If we have no candidate subkey then we assume that it's ok to sign
	// with the primary key.
	i := e.primaryIdentity()
	if !i.SelfSignature.FlagsValid || i.SelfSignature.FlagSign &&
		!i.SelfSignature.KeyExpired(now) ***REMOVED***
		return Key***REMOVED***e, e.PrimaryKey, e.PrivateKey, i.SelfSignature***REMOVED***, true
	***REMOVED***

	return Key***REMOVED******REMOVED***, false
***REMOVED***

// An EntityList contains one or more Entities.
type EntityList []*Entity

// KeysById returns the set of keys that have the given key id.
func (el EntityList) KeysById(id uint64) (keys []Key) ***REMOVED***
	for _, e := range el ***REMOVED***
		if e.PrimaryKey.KeyId == id ***REMOVED***
			var selfSig *packet.Signature
			for _, ident := range e.Identities ***REMOVED***
				if selfSig == nil ***REMOVED***
					selfSig = ident.SelfSignature
				***REMOVED*** else if ident.SelfSignature.IsPrimaryId != nil && *ident.SelfSignature.IsPrimaryId ***REMOVED***
					selfSig = ident.SelfSignature
					break
				***REMOVED***
			***REMOVED***
			keys = append(keys, Key***REMOVED***e, e.PrimaryKey, e.PrivateKey, selfSig***REMOVED***)
		***REMOVED***

		for _, subKey := range e.Subkeys ***REMOVED***
			if subKey.PublicKey.KeyId == id ***REMOVED***
				keys = append(keys, Key***REMOVED***e, subKey.PublicKey, subKey.PrivateKey, subKey.Sig***REMOVED***)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// KeysByIdAndUsage returns the set of keys with the given id that also meet
// the key usage given by requiredUsage.  The requiredUsage is expressed as
// the bitwise-OR of packet.KeyFlag* values.
func (el EntityList) KeysByIdUsage(id uint64, requiredUsage byte) (keys []Key) ***REMOVED***
	for _, key := range el.KeysById(id) ***REMOVED***
		if len(key.Entity.Revocations) > 0 ***REMOVED***
			continue
		***REMOVED***

		if key.SelfSignature.RevocationReason != nil ***REMOVED***
			continue
		***REMOVED***

		if key.SelfSignature.FlagsValid && requiredUsage != 0 ***REMOVED***
			var usage byte
			if key.SelfSignature.FlagCertify ***REMOVED***
				usage |= packet.KeyFlagCertify
			***REMOVED***
			if key.SelfSignature.FlagSign ***REMOVED***
				usage |= packet.KeyFlagSign
			***REMOVED***
			if key.SelfSignature.FlagEncryptCommunications ***REMOVED***
				usage |= packet.KeyFlagEncryptCommunications
			***REMOVED***
			if key.SelfSignature.FlagEncryptStorage ***REMOVED***
				usage |= packet.KeyFlagEncryptStorage
			***REMOVED***
			if usage&requiredUsage != requiredUsage ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***

		keys = append(keys, key)
	***REMOVED***
	return
***REMOVED***

// DecryptionKeys returns all private keys that are valid for decryption.
func (el EntityList) DecryptionKeys() (keys []Key) ***REMOVED***
	for _, e := range el ***REMOVED***
		for _, subKey := range e.Subkeys ***REMOVED***
			if subKey.PrivateKey != nil && (!subKey.Sig.FlagsValid || subKey.Sig.FlagEncryptStorage || subKey.Sig.FlagEncryptCommunications) ***REMOVED***
				keys = append(keys, Key***REMOVED***e, subKey.PublicKey, subKey.PrivateKey, subKey.Sig***REMOVED***)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// ReadArmoredKeyRing reads one or more public/private keys from an armor keyring file.
func ReadArmoredKeyRing(r io.Reader) (EntityList, error) ***REMOVED***
	block, err := armor.Decode(r)
	if err == io.EOF ***REMOVED***
		return nil, errors.InvalidArgumentError("no armored data found")
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if block.Type != PublicKeyType && block.Type != PrivateKeyType ***REMOVED***
		return nil, errors.InvalidArgumentError("expected public or private key block, got: " + block.Type)
	***REMOVED***

	return ReadKeyRing(block.Body)
***REMOVED***

// ReadKeyRing reads one or more public/private keys. Unsupported keys are
// ignored as long as at least a single valid key is found.
func ReadKeyRing(r io.Reader) (el EntityList, err error) ***REMOVED***
	packets := packet.NewReader(r)
	var lastUnsupportedError error

	for ***REMOVED***
		var e *Entity
		e, err = ReadEntity(packets)
		if err != nil ***REMOVED***
			// TODO: warn about skipped unsupported/unreadable keys
			if _, ok := err.(errors.UnsupportedError); ok ***REMOVED***
				lastUnsupportedError = err
				err = readToNextPublicKey(packets)
			***REMOVED*** else if _, ok := err.(errors.StructuralError); ok ***REMOVED***
				// Skip unreadable, badly-formatted keys
				lastUnsupportedError = err
				err = readToNextPublicKey(packets)
			***REMOVED***
			if err == io.EOF ***REMOVED***
				err = nil
				break
			***REMOVED***
			if err != nil ***REMOVED***
				el = nil
				break
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			el = append(el, e)
		***REMOVED***
	***REMOVED***

	if len(el) == 0 && err == nil ***REMOVED***
		err = lastUnsupportedError
	***REMOVED***
	return
***REMOVED***

// readToNextPublicKey reads packets until the start of the entity and leaves
// the first packet of the new entity in the Reader.
func readToNextPublicKey(packets *packet.Reader) (err error) ***REMOVED***
	var p packet.Packet
	for ***REMOVED***
		p, err = packets.Next()
		if err == io.EOF ***REMOVED***
			return
		***REMOVED*** else if err != nil ***REMOVED***
			if _, ok := err.(errors.UnsupportedError); ok ***REMOVED***
				err = nil
				continue
			***REMOVED***
			return
		***REMOVED***

		if pk, ok := p.(*packet.PublicKey); ok && !pk.IsSubkey ***REMOVED***
			packets.Unread(p)
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// ReadEntity reads an entity (public key, identities, subkeys etc) from the
// given Reader.
func ReadEntity(packets *packet.Reader) (*Entity, error) ***REMOVED***
	e := new(Entity)
	e.Identities = make(map[string]*Identity)

	p, err := packets.Next()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var ok bool
	if e.PrimaryKey, ok = p.(*packet.PublicKey); !ok ***REMOVED***
		if e.PrivateKey, ok = p.(*packet.PrivateKey); !ok ***REMOVED***
			packets.Unread(p)
			return nil, errors.StructuralError("first packet was not a public/private key")
		***REMOVED***
		e.PrimaryKey = &e.PrivateKey.PublicKey
	***REMOVED***

	if !e.PrimaryKey.PubKeyAlgo.CanSign() ***REMOVED***
		return nil, errors.StructuralError("primary key cannot be used for signatures")
	***REMOVED***

	var current *Identity
	var revocations []*packet.Signature
EachPacket:
	for ***REMOVED***
		p, err := packets.Next()
		if err == io.EOF ***REMOVED***
			break
		***REMOVED*** else if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		switch pkt := p.(type) ***REMOVED***
		case *packet.UserId:
			current = new(Identity)
			current.Name = pkt.Id
			current.UserId = pkt
			e.Identities[pkt.Id] = current

			for ***REMOVED***
				p, err = packets.Next()
				if err == io.EOF ***REMOVED***
					return nil, io.ErrUnexpectedEOF
				***REMOVED*** else if err != nil ***REMOVED***
					return nil, err
				***REMOVED***

				sig, ok := p.(*packet.Signature)
				if !ok ***REMOVED***
					return nil, errors.StructuralError("user ID packet not followed by self-signature")
				***REMOVED***

				if (sig.SigType == packet.SigTypePositiveCert || sig.SigType == packet.SigTypeGenericCert) && sig.IssuerKeyId != nil && *sig.IssuerKeyId == e.PrimaryKey.KeyId ***REMOVED***
					if err = e.PrimaryKey.VerifyUserIdSignature(pkt.Id, e.PrimaryKey, sig); err != nil ***REMOVED***
						return nil, errors.StructuralError("user ID self-signature invalid: " + err.Error())
					***REMOVED***
					current.SelfSignature = sig
					break
				***REMOVED***
				current.Signatures = append(current.Signatures, sig)
			***REMOVED***
		case *packet.Signature:
			if pkt.SigType == packet.SigTypeKeyRevocation ***REMOVED***
				revocations = append(revocations, pkt)
			***REMOVED*** else if pkt.SigType == packet.SigTypeDirectSignature ***REMOVED***
				// TODO: RFC4880 5.2.1 permits signatures
				// directly on keys (eg. to bind additional
				// revocation keys).
			***REMOVED*** else if current == nil ***REMOVED***
				return nil, errors.StructuralError("signature packet found before user id packet")
			***REMOVED*** else ***REMOVED***
				current.Signatures = append(current.Signatures, pkt)
			***REMOVED***
		case *packet.PrivateKey:
			if pkt.IsSubkey == false ***REMOVED***
				packets.Unread(p)
				break EachPacket
			***REMOVED***
			err = addSubkey(e, packets, &pkt.PublicKey, pkt)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		case *packet.PublicKey:
			if pkt.IsSubkey == false ***REMOVED***
				packets.Unread(p)
				break EachPacket
			***REMOVED***
			err = addSubkey(e, packets, pkt, nil)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		default:
			// we ignore unknown packets
		***REMOVED***
	***REMOVED***

	if len(e.Identities) == 0 ***REMOVED***
		return nil, errors.StructuralError("entity without any identities")
	***REMOVED***

	for _, revocation := range revocations ***REMOVED***
		err = e.PrimaryKey.VerifyRevocationSignature(revocation)
		if err == nil ***REMOVED***
			e.Revocations = append(e.Revocations, revocation)
		***REMOVED*** else ***REMOVED***
			// TODO: RFC 4880 5.2.3.15 defines revocation keys.
			return nil, errors.StructuralError("revocation signature signed by alternate key")
		***REMOVED***
	***REMOVED***

	return e, nil
***REMOVED***

func addSubkey(e *Entity, packets *packet.Reader, pub *packet.PublicKey, priv *packet.PrivateKey) error ***REMOVED***
	var subKey Subkey
	subKey.PublicKey = pub
	subKey.PrivateKey = priv
	p, err := packets.Next()
	if err == io.EOF ***REMOVED***
		return io.ErrUnexpectedEOF
	***REMOVED***
	if err != nil ***REMOVED***
		return errors.StructuralError("subkey signature invalid: " + err.Error())
	***REMOVED***
	var ok bool
	subKey.Sig, ok = p.(*packet.Signature)
	if !ok ***REMOVED***
		return errors.StructuralError("subkey packet not followed by signature")
	***REMOVED***
	if subKey.Sig.SigType != packet.SigTypeSubkeyBinding && subKey.Sig.SigType != packet.SigTypeSubkeyRevocation ***REMOVED***
		return errors.StructuralError("subkey signature with wrong type")
	***REMOVED***
	err = e.PrimaryKey.VerifyKeySignature(subKey.PublicKey, subKey.Sig)
	if err != nil ***REMOVED***
		return errors.StructuralError("subkey signature invalid: " + err.Error())
	***REMOVED***
	e.Subkeys = append(e.Subkeys, subKey)
	return nil
***REMOVED***

const defaultRSAKeyBits = 2048

// NewEntity returns an Entity that contains a fresh RSA/RSA keypair with a
// single identity composed of the given full name, comment and email, any of
// which may be empty but must not contain any of "()<>\x00".
// If config is nil, sensible defaults will be used.
func NewEntity(name, comment, email string, config *packet.Config) (*Entity, error) ***REMOVED***
	currentTime := config.Now()

	bits := defaultRSAKeyBits
	if config != nil && config.RSABits != 0 ***REMOVED***
		bits = config.RSABits
	***REMOVED***

	uid := packet.NewUserId(name, comment, email)
	if uid == nil ***REMOVED***
		return nil, errors.InvalidArgumentError("user id field contained invalid characters")
	***REMOVED***
	signingPriv, err := rsa.GenerateKey(config.Random(), bits)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	encryptingPriv, err := rsa.GenerateKey(config.Random(), bits)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	e := &Entity***REMOVED***
		PrimaryKey: packet.NewRSAPublicKey(currentTime, &signingPriv.PublicKey),
		PrivateKey: packet.NewRSAPrivateKey(currentTime, signingPriv),
		Identities: make(map[string]*Identity),
	***REMOVED***
	isPrimaryId := true
	e.Identities[uid.Id] = &Identity***REMOVED***
		Name:   uid.Name,
		UserId: uid,
		SelfSignature: &packet.Signature***REMOVED***
			CreationTime: currentTime,
			SigType:      packet.SigTypePositiveCert,
			PubKeyAlgo:   packet.PubKeyAlgoRSA,
			Hash:         config.Hash(),
			IsPrimaryId:  &isPrimaryId,
			FlagsValid:   true,
			FlagSign:     true,
			FlagCertify:  true,
			IssuerKeyId:  &e.PrimaryKey.KeyId,
		***REMOVED***,
	***REMOVED***

	// If the user passes in a DefaultHash via packet.Config,
	// set the PreferredHash for the SelfSignature.
	if config != nil && config.DefaultHash != 0 ***REMOVED***
		e.Identities[uid.Id].SelfSignature.PreferredHash = []uint8***REMOVED***hashToHashId(config.DefaultHash)***REMOVED***
	***REMOVED***

	e.Subkeys = make([]Subkey, 1)
	e.Subkeys[0] = Subkey***REMOVED***
		PublicKey:  packet.NewRSAPublicKey(currentTime, &encryptingPriv.PublicKey),
		PrivateKey: packet.NewRSAPrivateKey(currentTime, encryptingPriv),
		Sig: &packet.Signature***REMOVED***
			CreationTime:              currentTime,
			SigType:                   packet.SigTypeSubkeyBinding,
			PubKeyAlgo:                packet.PubKeyAlgoRSA,
			Hash:                      config.Hash(),
			FlagsValid:                true,
			FlagEncryptStorage:        true,
			FlagEncryptCommunications: true,
			IssuerKeyId:               &e.PrimaryKey.KeyId,
		***REMOVED***,
	***REMOVED***
	e.Subkeys[0].PublicKey.IsSubkey = true
	e.Subkeys[0].PrivateKey.IsSubkey = true

	return e, nil
***REMOVED***

// SerializePrivate serializes an Entity, including private key material, to
// the given Writer. For now, it must only be used on an Entity returned from
// NewEntity.
// If config is nil, sensible defaults will be used.
func (e *Entity) SerializePrivate(w io.Writer, config *packet.Config) (err error) ***REMOVED***
	err = e.PrivateKey.Serialize(w)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	for _, ident := range e.Identities ***REMOVED***
		err = ident.UserId.Serialize(w)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		err = ident.SelfSignature.SignUserId(ident.UserId.Id, e.PrimaryKey, e.PrivateKey, config)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		err = ident.SelfSignature.Serialize(w)
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	for _, subkey := range e.Subkeys ***REMOVED***
		err = subkey.PrivateKey.Serialize(w)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		err = subkey.Sig.SignKey(subkey.PublicKey, e.PrivateKey, config)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		err = subkey.Sig.Serialize(w)
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Serialize writes the public part of the given Entity to w. (No private
// key material will be output).
func (e *Entity) Serialize(w io.Writer) error ***REMOVED***
	err := e.PrimaryKey.Serialize(w)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, ident := range e.Identities ***REMOVED***
		err = ident.UserId.Serialize(w)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		err = ident.SelfSignature.Serialize(w)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		for _, sig := range ident.Signatures ***REMOVED***
			err = sig.Serialize(w)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for _, subkey := range e.Subkeys ***REMOVED***
		err = subkey.PublicKey.Serialize(w)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		err = subkey.Sig.Serialize(w)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// SignIdentity adds a signature to e, from signer, attesting that identity is
// associated with e. The provided identity must already be an element of
// e.Identities and the private key of signer must have been decrypted if
// necessary.
// If config is nil, sensible defaults will be used.
func (e *Entity) SignIdentity(identity string, signer *Entity, config *packet.Config) error ***REMOVED***
	if signer.PrivateKey == nil ***REMOVED***
		return errors.InvalidArgumentError("signing Entity must have a private key")
	***REMOVED***
	if signer.PrivateKey.Encrypted ***REMOVED***
		return errors.InvalidArgumentError("signing Entity's private key must be decrypted")
	***REMOVED***
	ident, ok := e.Identities[identity]
	if !ok ***REMOVED***
		return errors.InvalidArgumentError("given identity string not found in Entity")
	***REMOVED***

	sig := &packet.Signature***REMOVED***
		SigType:      packet.SigTypeGenericCert,
		PubKeyAlgo:   signer.PrivateKey.PubKeyAlgo,
		Hash:         config.Hash(),
		CreationTime: config.Now(),
		IssuerKeyId:  &signer.PrivateKey.KeyId,
	***REMOVED***
	if err := sig.SignUserId(identity, e.PrimaryKey, signer.PrivateKey, config); err != nil ***REMOVED***
		return err
	***REMOVED***
	ident.Signatures = append(ident.Signatures, sig)
	return nil
***REMOVED***
