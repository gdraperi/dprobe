// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"io"
	"math/big"

	"golang.org/x/crypto/curve25519"
)

const (
	kexAlgoDH1SHA1          = "diffie-hellman-group1-sha1"
	kexAlgoDH14SHA1         = "diffie-hellman-group14-sha1"
	kexAlgoECDH256          = "ecdh-sha2-nistp256"
	kexAlgoECDH384          = "ecdh-sha2-nistp384"
	kexAlgoECDH521          = "ecdh-sha2-nistp521"
	kexAlgoCurve25519SHA256 = "curve25519-sha256@libssh.org"
)

// kexResult captures the outcome of a key exchange.
type kexResult struct ***REMOVED***
	// Session hash. See also RFC 4253, section 8.
	H []byte

	// Shared secret. See also RFC 4253, section 8.
	K []byte

	// Host key as hashed into H.
	HostKey []byte

	// Signature of H.
	Signature []byte

	// A cryptographic hash function that matches the security
	// level of the key exchange algorithm. It is used for
	// calculating H, and for deriving keys from H and K.
	Hash crypto.Hash

	// The session ID, which is the first H computed. This is used
	// to derive key material inside the transport.
	SessionID []byte
***REMOVED***

// handshakeMagics contains data that is always included in the
// session hash.
type handshakeMagics struct ***REMOVED***
	clientVersion, serverVersion []byte
	clientKexInit, serverKexInit []byte
***REMOVED***

func (m *handshakeMagics) write(w io.Writer) ***REMOVED***
	writeString(w, m.clientVersion)
	writeString(w, m.serverVersion)
	writeString(w, m.clientKexInit)
	writeString(w, m.serverKexInit)
***REMOVED***

// kexAlgorithm abstracts different key exchange algorithms.
type kexAlgorithm interface ***REMOVED***
	// Server runs server-side key agreement, signing the result
	// with a hostkey.
	Server(p packetConn, rand io.Reader, magics *handshakeMagics, s Signer) (*kexResult, error)

	// Client runs the client-side key agreement. Caller is
	// responsible for verifying the host key signature.
	Client(p packetConn, rand io.Reader, magics *handshakeMagics) (*kexResult, error)
***REMOVED***

// dhGroup is a multiplicative group suitable for implementing Diffie-Hellman key agreement.
type dhGroup struct ***REMOVED***
	g, p, pMinus1 *big.Int
***REMOVED***

func (group *dhGroup) diffieHellman(theirPublic, myPrivate *big.Int) (*big.Int, error) ***REMOVED***
	if theirPublic.Cmp(bigOne) <= 0 || theirPublic.Cmp(group.pMinus1) >= 0 ***REMOVED***
		return nil, errors.New("ssh: DH parameter out of bounds")
	***REMOVED***
	return new(big.Int).Exp(theirPublic, myPrivate, group.p), nil
***REMOVED***

func (group *dhGroup) Client(c packetConn, randSource io.Reader, magics *handshakeMagics) (*kexResult, error) ***REMOVED***
	hashFunc := crypto.SHA1

	var x *big.Int
	for ***REMOVED***
		var err error
		if x, err = rand.Int(randSource, group.pMinus1); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if x.Sign() > 0 ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	X := new(big.Int).Exp(group.g, x, group.p)
	kexDHInit := kexDHInitMsg***REMOVED***
		X: X,
	***REMOVED***
	if err := c.writePacket(Marshal(&kexDHInit)); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	packet, err := c.readPacket()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var kexDHReply kexDHReplyMsg
	if err = Unmarshal(packet, &kexDHReply); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ki, err := group.diffieHellman(kexDHReply.Y, x)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	h := hashFunc.New()
	magics.write(h)
	writeString(h, kexDHReply.HostKey)
	writeInt(h, X)
	writeInt(h, kexDHReply.Y)
	K := make([]byte, intLength(ki))
	marshalInt(K, ki)
	h.Write(K)

	return &kexResult***REMOVED***
		H:         h.Sum(nil),
		K:         K,
		HostKey:   kexDHReply.HostKey,
		Signature: kexDHReply.Signature,
		Hash:      crypto.SHA1,
	***REMOVED***, nil
***REMOVED***

func (group *dhGroup) Server(c packetConn, randSource io.Reader, magics *handshakeMagics, priv Signer) (result *kexResult, err error) ***REMOVED***
	hashFunc := crypto.SHA1
	packet, err := c.readPacket()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	var kexDHInit kexDHInitMsg
	if err = Unmarshal(packet, &kexDHInit); err != nil ***REMOVED***
		return
	***REMOVED***

	var y *big.Int
	for ***REMOVED***
		if y, err = rand.Int(randSource, group.pMinus1); err != nil ***REMOVED***
			return
		***REMOVED***
		if y.Sign() > 0 ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	Y := new(big.Int).Exp(group.g, y, group.p)
	ki, err := group.diffieHellman(kexDHInit.X, y)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	hostKeyBytes := priv.PublicKey().Marshal()

	h := hashFunc.New()
	magics.write(h)
	writeString(h, hostKeyBytes)
	writeInt(h, kexDHInit.X)
	writeInt(h, Y)

	K := make([]byte, intLength(ki))
	marshalInt(K, ki)
	h.Write(K)

	H := h.Sum(nil)

	// H is already a hash, but the hostkey signing will apply its
	// own key-specific hash algorithm.
	sig, err := signAndMarshal(priv, randSource, H)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	kexDHReply := kexDHReplyMsg***REMOVED***
		HostKey:   hostKeyBytes,
		Y:         Y,
		Signature: sig,
	***REMOVED***
	packet = Marshal(&kexDHReply)

	err = c.writePacket(packet)
	return &kexResult***REMOVED***
		H:         H,
		K:         K,
		HostKey:   hostKeyBytes,
		Signature: sig,
		Hash:      crypto.SHA1,
	***REMOVED***, nil
***REMOVED***

// ecdh performs Elliptic Curve Diffie-Hellman key exchange as
// described in RFC 5656, section 4.
type ecdh struct ***REMOVED***
	curve elliptic.Curve
***REMOVED***

func (kex *ecdh) Client(c packetConn, rand io.Reader, magics *handshakeMagics) (*kexResult, error) ***REMOVED***
	ephKey, err := ecdsa.GenerateKey(kex.curve, rand)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	kexInit := kexECDHInitMsg***REMOVED***
		ClientPubKey: elliptic.Marshal(kex.curve, ephKey.PublicKey.X, ephKey.PublicKey.Y),
	***REMOVED***

	serialized := Marshal(&kexInit)
	if err := c.writePacket(serialized); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	packet, err := c.readPacket()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var reply kexECDHReplyMsg
	if err = Unmarshal(packet, &reply); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	x, y, err := unmarshalECKey(kex.curve, reply.EphemeralPubKey)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// generate shared secret
	secret, _ := kex.curve.ScalarMult(x, y, ephKey.D.Bytes())

	h := ecHash(kex.curve).New()
	magics.write(h)
	writeString(h, reply.HostKey)
	writeString(h, kexInit.ClientPubKey)
	writeString(h, reply.EphemeralPubKey)
	K := make([]byte, intLength(secret))
	marshalInt(K, secret)
	h.Write(K)

	return &kexResult***REMOVED***
		H:         h.Sum(nil),
		K:         K,
		HostKey:   reply.HostKey,
		Signature: reply.Signature,
		Hash:      ecHash(kex.curve),
	***REMOVED***, nil
***REMOVED***

// unmarshalECKey parses and checks an EC key.
func unmarshalECKey(curve elliptic.Curve, pubkey []byte) (x, y *big.Int, err error) ***REMOVED***
	x, y = elliptic.Unmarshal(curve, pubkey)
	if x == nil ***REMOVED***
		return nil, nil, errors.New("ssh: elliptic.Unmarshal failure")
	***REMOVED***
	if !validateECPublicKey(curve, x, y) ***REMOVED***
		return nil, nil, errors.New("ssh: public key not on curve")
	***REMOVED***
	return x, y, nil
***REMOVED***

// validateECPublicKey checks that the point is a valid public key for
// the given curve. See [SEC1], 3.2.2
func validateECPublicKey(curve elliptic.Curve, x, y *big.Int) bool ***REMOVED***
	if x.Sign() == 0 && y.Sign() == 0 ***REMOVED***
		return false
	***REMOVED***

	if x.Cmp(curve.Params().P) >= 0 ***REMOVED***
		return false
	***REMOVED***

	if y.Cmp(curve.Params().P) >= 0 ***REMOVED***
		return false
	***REMOVED***

	if !curve.IsOnCurve(x, y) ***REMOVED***
		return false
	***REMOVED***

	// We don't check if N * PubKey == 0, since
	//
	// - the NIST curves have cofactor = 1, so this is implicit.
	// (We don't foresee an implementation that supports non NIST
	// curves)
	//
	// - for ephemeral keys, we don't need to worry about small
	// subgroup attacks.
	return true
***REMOVED***

func (kex *ecdh) Server(c packetConn, rand io.Reader, magics *handshakeMagics, priv Signer) (result *kexResult, err error) ***REMOVED***
	packet, err := c.readPacket()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var kexECDHInit kexECDHInitMsg
	if err = Unmarshal(packet, &kexECDHInit); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	clientX, clientY, err := unmarshalECKey(kex.curve, kexECDHInit.ClientPubKey)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// We could cache this key across multiple users/multiple
	// connection attempts, but the benefit is small. OpenSSH
	// generates a new key for each incoming connection.
	ephKey, err := ecdsa.GenerateKey(kex.curve, rand)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	hostKeyBytes := priv.PublicKey().Marshal()

	serializedEphKey := elliptic.Marshal(kex.curve, ephKey.PublicKey.X, ephKey.PublicKey.Y)

	// generate shared secret
	secret, _ := kex.curve.ScalarMult(clientX, clientY, ephKey.D.Bytes())

	h := ecHash(kex.curve).New()
	magics.write(h)
	writeString(h, hostKeyBytes)
	writeString(h, kexECDHInit.ClientPubKey)
	writeString(h, serializedEphKey)

	K := make([]byte, intLength(secret))
	marshalInt(K, secret)
	h.Write(K)

	H := h.Sum(nil)

	// H is already a hash, but the hostkey signing will apply its
	// own key-specific hash algorithm.
	sig, err := signAndMarshal(priv, rand, H)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	reply := kexECDHReplyMsg***REMOVED***
		EphemeralPubKey: serializedEphKey,
		HostKey:         hostKeyBytes,
		Signature:       sig,
	***REMOVED***

	serialized := Marshal(&reply)
	if err := c.writePacket(serialized); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &kexResult***REMOVED***
		H:         H,
		K:         K,
		HostKey:   reply.HostKey,
		Signature: sig,
		Hash:      ecHash(kex.curve),
	***REMOVED***, nil
***REMOVED***

var kexAlgoMap = map[string]kexAlgorithm***REMOVED******REMOVED***

func init() ***REMOVED***
	// This is the group called diffie-hellman-group1-sha1 in RFC
	// 4253 and Oakley Group 2 in RFC 2409.
	p, _ := new(big.Int).SetString("FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE649286651ECE65381FFFFFFFFFFFFFFFF", 16)
	kexAlgoMap[kexAlgoDH1SHA1] = &dhGroup***REMOVED***
		g:       new(big.Int).SetInt64(2),
		p:       p,
		pMinus1: new(big.Int).Sub(p, bigOne),
	***REMOVED***

	// This is the group called diffie-hellman-group14-sha1 in RFC
	// 4253 and Oakley Group 14 in RFC 3526.
	p, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3DC2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F83655D23DCA3AD961C62F356208552BB9ED529077096966D670C354E4ABC9804F1746C08CA18217C32905E462E36CE3BE39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9DE2BCBF6955817183995497CEA956AE515D2261898FA051015728E5A8AACAA68FFFFFFFFFFFFFFFF", 16)

	kexAlgoMap[kexAlgoDH14SHA1] = &dhGroup***REMOVED***
		g:       new(big.Int).SetInt64(2),
		p:       p,
		pMinus1: new(big.Int).Sub(p, bigOne),
	***REMOVED***

	kexAlgoMap[kexAlgoECDH521] = &ecdh***REMOVED***elliptic.P521()***REMOVED***
	kexAlgoMap[kexAlgoECDH384] = &ecdh***REMOVED***elliptic.P384()***REMOVED***
	kexAlgoMap[kexAlgoECDH256] = &ecdh***REMOVED***elliptic.P256()***REMOVED***
	kexAlgoMap[kexAlgoCurve25519SHA256] = &curve25519sha256***REMOVED******REMOVED***
***REMOVED***

// curve25519sha256 implements the curve25519-sha256@libssh.org key
// agreement protocol, as described in
// https://git.libssh.org/projects/libssh.git/tree/doc/curve25519-sha256@libssh.org.txt
type curve25519sha256 struct***REMOVED******REMOVED***

type curve25519KeyPair struct ***REMOVED***
	priv [32]byte
	pub  [32]byte
***REMOVED***

func (kp *curve25519KeyPair) generate(rand io.Reader) error ***REMOVED***
	if _, err := io.ReadFull(rand, kp.priv[:]); err != nil ***REMOVED***
		return err
	***REMOVED***
	curve25519.ScalarBaseMult(&kp.pub, &kp.priv)
	return nil
***REMOVED***

// curve25519Zeros is just an array of 32 zero bytes so that we have something
// convenient to compare against in order to reject curve25519 points with the
// wrong order.
var curve25519Zeros [32]byte

func (kex *curve25519sha256) Client(c packetConn, rand io.Reader, magics *handshakeMagics) (*kexResult, error) ***REMOVED***
	var kp curve25519KeyPair
	if err := kp.generate(rand); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := c.writePacket(Marshal(&kexECDHInitMsg***REMOVED***kp.pub[:]***REMOVED***)); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	packet, err := c.readPacket()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var reply kexECDHReplyMsg
	if err = Unmarshal(packet, &reply); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(reply.EphemeralPubKey) != 32 ***REMOVED***
		return nil, errors.New("ssh: peer's curve25519 public value has wrong length")
	***REMOVED***

	var servPub, secret [32]byte
	copy(servPub[:], reply.EphemeralPubKey)
	curve25519.ScalarMult(&secret, &kp.priv, &servPub)
	if subtle.ConstantTimeCompare(secret[:], curve25519Zeros[:]) == 1 ***REMOVED***
		return nil, errors.New("ssh: peer's curve25519 public value has wrong order")
	***REMOVED***

	h := crypto.SHA256.New()
	magics.write(h)
	writeString(h, reply.HostKey)
	writeString(h, kp.pub[:])
	writeString(h, reply.EphemeralPubKey)

	ki := new(big.Int).SetBytes(secret[:])
	K := make([]byte, intLength(ki))
	marshalInt(K, ki)
	h.Write(K)

	return &kexResult***REMOVED***
		H:         h.Sum(nil),
		K:         K,
		HostKey:   reply.HostKey,
		Signature: reply.Signature,
		Hash:      crypto.SHA256,
	***REMOVED***, nil
***REMOVED***

func (kex *curve25519sha256) Server(c packetConn, rand io.Reader, magics *handshakeMagics, priv Signer) (result *kexResult, err error) ***REMOVED***
	packet, err := c.readPacket()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	var kexInit kexECDHInitMsg
	if err = Unmarshal(packet, &kexInit); err != nil ***REMOVED***
		return
	***REMOVED***

	if len(kexInit.ClientPubKey) != 32 ***REMOVED***
		return nil, errors.New("ssh: peer's curve25519 public value has wrong length")
	***REMOVED***

	var kp curve25519KeyPair
	if err := kp.generate(rand); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var clientPub, secret [32]byte
	copy(clientPub[:], kexInit.ClientPubKey)
	curve25519.ScalarMult(&secret, &kp.priv, &clientPub)
	if subtle.ConstantTimeCompare(secret[:], curve25519Zeros[:]) == 1 ***REMOVED***
		return nil, errors.New("ssh: peer's curve25519 public value has wrong order")
	***REMOVED***

	hostKeyBytes := priv.PublicKey().Marshal()

	h := crypto.SHA256.New()
	magics.write(h)
	writeString(h, hostKeyBytes)
	writeString(h, kexInit.ClientPubKey)
	writeString(h, kp.pub[:])

	ki := new(big.Int).SetBytes(secret[:])
	K := make([]byte, intLength(ki))
	marshalInt(K, ki)
	h.Write(K)

	H := h.Sum(nil)

	sig, err := signAndMarshal(priv, rand, H)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	reply := kexECDHReplyMsg***REMOVED***
		EphemeralPubKey: kp.pub[:],
		HostKey:         hostKeyBytes,
		Signature:       sig,
	***REMOVED***
	if err := c.writePacket(Marshal(&reply)); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &kexResult***REMOVED***
		H:         H,
		K:         K,
		HostKey:   hostKeyBytes,
		Signature: sig,
		Hash:      crypto.SHA256,
	***REMOVED***, nil
***REMOVED***
