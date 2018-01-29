// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"crypto"
	"crypto/rand"
	"fmt"
	"io"
	"math"
	"sync"

	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"
)

// These are string constants in the SSH protocol.
const (
	compressionNone = "none"
	serviceUserAuth = "ssh-userauth"
	serviceSSH      = "ssh-connection"
)

// supportedCiphers lists ciphers we support but might not recommend.
var supportedCiphers = []string***REMOVED***
	"aes128-ctr", "aes192-ctr", "aes256-ctr",
	"aes128-gcm@openssh.com",
	chacha20Poly1305ID,
	"arcfour256", "arcfour128", "arcfour",
	aes128cbcID,
	tripledescbcID,
***REMOVED***

// preferredCiphers specifies the default preference for ciphers.
var preferredCiphers = []string***REMOVED***
	"aes128-gcm@openssh.com",
	chacha20Poly1305ID,
	"aes128-ctr", "aes192-ctr", "aes256-ctr",
***REMOVED***

// supportedKexAlgos specifies the supported key-exchange algorithms in
// preference order.
var supportedKexAlgos = []string***REMOVED***
	kexAlgoCurve25519SHA256,
	// P384 and P521 are not constant-time yet, but since we don't
	// reuse ephemeral keys, using them for ECDH should be OK.
	kexAlgoECDH256, kexAlgoECDH384, kexAlgoECDH521,
	kexAlgoDH14SHA1, kexAlgoDH1SHA1,
***REMOVED***

// supportedHostKeyAlgos specifies the supported host-key algorithms (i.e. methods
// of authenticating servers) in preference order.
var supportedHostKeyAlgos = []string***REMOVED***
	CertAlgoRSAv01, CertAlgoDSAv01, CertAlgoECDSA256v01,
	CertAlgoECDSA384v01, CertAlgoECDSA521v01, CertAlgoED25519v01,

	KeyAlgoECDSA256, KeyAlgoECDSA384, KeyAlgoECDSA521,
	KeyAlgoRSA, KeyAlgoDSA,

	KeyAlgoED25519,
***REMOVED***

// supportedMACs specifies a default set of MAC algorithms in preference order.
// This is based on RFC 4253, section 6.4, but with hmac-md5 variants removed
// because they have reached the end of their useful life.
var supportedMACs = []string***REMOVED***
	"hmac-sha2-256-etm@openssh.com", "hmac-sha2-256", "hmac-sha1", "hmac-sha1-96",
***REMOVED***

var supportedCompressions = []string***REMOVED***compressionNone***REMOVED***

// hashFuncs keeps the mapping of supported algorithms to their respective
// hashes needed for signature verification.
var hashFuncs = map[string]crypto.Hash***REMOVED***
	KeyAlgoRSA:          crypto.SHA1,
	KeyAlgoDSA:          crypto.SHA1,
	KeyAlgoECDSA256:     crypto.SHA256,
	KeyAlgoECDSA384:     crypto.SHA384,
	KeyAlgoECDSA521:     crypto.SHA512,
	CertAlgoRSAv01:      crypto.SHA1,
	CertAlgoDSAv01:      crypto.SHA1,
	CertAlgoECDSA256v01: crypto.SHA256,
	CertAlgoECDSA384v01: crypto.SHA384,
	CertAlgoECDSA521v01: crypto.SHA512,
***REMOVED***

// unexpectedMessageError results when the SSH message that we received didn't
// match what we wanted.
func unexpectedMessageError(expected, got uint8) error ***REMOVED***
	return fmt.Errorf("ssh: unexpected message type %d (expected %d)", got, expected)
***REMOVED***

// parseError results from a malformed SSH message.
func parseError(tag uint8) error ***REMOVED***
	return fmt.Errorf("ssh: parse error in message type %d", tag)
***REMOVED***

func findCommon(what string, client []string, server []string) (common string, err error) ***REMOVED***
	for _, c := range client ***REMOVED***
		for _, s := range server ***REMOVED***
			if c == s ***REMOVED***
				return c, nil
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return "", fmt.Errorf("ssh: no common algorithm for %s; client offered: %v, server offered: %v", what, client, server)
***REMOVED***

type directionAlgorithms struct ***REMOVED***
	Cipher      string
	MAC         string
	Compression string
***REMOVED***

// rekeyBytes returns a rekeying intervals in bytes.
func (a *directionAlgorithms) rekeyBytes() int64 ***REMOVED***
	// According to RFC4344 block ciphers should rekey after
	// 2^(BLOCKSIZE/4) blocks. For all AES flavors BLOCKSIZE is
	// 128.
	switch a.Cipher ***REMOVED***
	case "aes128-ctr", "aes192-ctr", "aes256-ctr", gcmCipherID, aes128cbcID:
		return 16 * (1 << 32)

	***REMOVED***

	// For others, stick with RFC4253 recommendation to rekey after 1 Gb of data.
	return 1 << 30
***REMOVED***

type algorithms struct ***REMOVED***
	kex     string
	hostKey string
	w       directionAlgorithms
	r       directionAlgorithms
***REMOVED***

func findAgreedAlgorithms(clientKexInit, serverKexInit *kexInitMsg) (algs *algorithms, err error) ***REMOVED***
	result := &algorithms***REMOVED******REMOVED***

	result.kex, err = findCommon("key exchange", clientKexInit.KexAlgos, serverKexInit.KexAlgos)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	result.hostKey, err = findCommon("host key", clientKexInit.ServerHostKeyAlgos, serverKexInit.ServerHostKeyAlgos)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	result.w.Cipher, err = findCommon("client to server cipher", clientKexInit.CiphersClientServer, serverKexInit.CiphersClientServer)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	result.r.Cipher, err = findCommon("server to client cipher", clientKexInit.CiphersServerClient, serverKexInit.CiphersServerClient)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	result.w.MAC, err = findCommon("client to server MAC", clientKexInit.MACsClientServer, serverKexInit.MACsClientServer)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	result.r.MAC, err = findCommon("server to client MAC", clientKexInit.MACsServerClient, serverKexInit.MACsServerClient)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	result.w.Compression, err = findCommon("client to server compression", clientKexInit.CompressionClientServer, serverKexInit.CompressionClientServer)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	result.r.Compression, err = findCommon("server to client compression", clientKexInit.CompressionServerClient, serverKexInit.CompressionServerClient)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	return result, nil
***REMOVED***

// If rekeythreshold is too small, we can't make any progress sending
// stuff.
const minRekeyThreshold uint64 = 256

// Config contains configuration data common to both ServerConfig and
// ClientConfig.
type Config struct ***REMOVED***
	// Rand provides the source of entropy for cryptographic
	// primitives. If Rand is nil, the cryptographic random reader
	// in package crypto/rand will be used.
	Rand io.Reader

	// The maximum number of bytes sent or received after which a
	// new key is negotiated. It must be at least 256. If
	// unspecified, a size suitable for the chosen cipher is used.
	RekeyThreshold uint64

	// The allowed key exchanges algorithms. If unspecified then a
	// default set of algorithms is used.
	KeyExchanges []string

	// The allowed cipher algorithms. If unspecified then a sensible
	// default is used.
	Ciphers []string

	// The allowed MAC algorithms. If unspecified then a sensible default
	// is used.
	MACs []string
***REMOVED***

// SetDefaults sets sensible values for unset fields in config. This is
// exported for testing: Configs passed to SSH functions are copied and have
// default values set automatically.
func (c *Config) SetDefaults() ***REMOVED***
	if c.Rand == nil ***REMOVED***
		c.Rand = rand.Reader
	***REMOVED***
	if c.Ciphers == nil ***REMOVED***
		c.Ciphers = preferredCiphers
	***REMOVED***
	var ciphers []string
	for _, c := range c.Ciphers ***REMOVED***
		if cipherModes[c] != nil ***REMOVED***
			// reject the cipher if we have no cipherModes definition
			ciphers = append(ciphers, c)
		***REMOVED***
	***REMOVED***
	c.Ciphers = ciphers

	if c.KeyExchanges == nil ***REMOVED***
		c.KeyExchanges = supportedKexAlgos
	***REMOVED***

	if c.MACs == nil ***REMOVED***
		c.MACs = supportedMACs
	***REMOVED***

	if c.RekeyThreshold == 0 ***REMOVED***
		// cipher specific default
	***REMOVED*** else if c.RekeyThreshold < minRekeyThreshold ***REMOVED***
		c.RekeyThreshold = minRekeyThreshold
	***REMOVED*** else if c.RekeyThreshold >= math.MaxInt64 ***REMOVED***
		// Avoid weirdness if somebody uses -1 as a threshold.
		c.RekeyThreshold = math.MaxInt64
	***REMOVED***
***REMOVED***

// buildDataSignedForAuth returns the data that is signed in order to prove
// possession of a private key. See RFC 4252, section 7.
func buildDataSignedForAuth(sessionID []byte, req userAuthRequestMsg, algo, pubKey []byte) []byte ***REMOVED***
	data := struct ***REMOVED***
		Session []byte
		Type    byte
		User    string
		Service string
		Method  string
		Sign    bool
		Algo    []byte
		PubKey  []byte
	***REMOVED******REMOVED***
		sessionID,
		msgUserAuthRequest,
		req.User,
		req.Service,
		req.Method,
		true,
		algo,
		pubKey,
	***REMOVED***
	return Marshal(data)
***REMOVED***

func appendU16(buf []byte, n uint16) []byte ***REMOVED***
	return append(buf, byte(n>>8), byte(n))
***REMOVED***

func appendU32(buf []byte, n uint32) []byte ***REMOVED***
	return append(buf, byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
***REMOVED***

func appendU64(buf []byte, n uint64) []byte ***REMOVED***
	return append(buf,
		byte(n>>56), byte(n>>48), byte(n>>40), byte(n>>32),
		byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
***REMOVED***

func appendInt(buf []byte, n int) []byte ***REMOVED***
	return appendU32(buf, uint32(n))
***REMOVED***

func appendString(buf []byte, s string) []byte ***REMOVED***
	buf = appendU32(buf, uint32(len(s)))
	buf = append(buf, s...)
	return buf
***REMOVED***

func appendBool(buf []byte, b bool) []byte ***REMOVED***
	if b ***REMOVED***
		return append(buf, 1)
	***REMOVED***
	return append(buf, 0)
***REMOVED***

// newCond is a helper to hide the fact that there is no usable zero
// value for sync.Cond.
func newCond() *sync.Cond ***REMOVED*** return sync.NewCond(new(sync.Mutex)) ***REMOVED***

// window represents the buffer available to clients
// wishing to write to a channel.
type window struct ***REMOVED***
	*sync.Cond
	win          uint32 // RFC 4254 5.2 says the window size can grow to 2^32-1
	writeWaiters int
	closed       bool
***REMOVED***

// add adds win to the amount of window available
// for consumers.
func (w *window) add(win uint32) bool ***REMOVED***
	// a zero sized window adjust is a noop.
	if win == 0 ***REMOVED***
		return true
	***REMOVED***
	w.L.Lock()
	if w.win+win < win ***REMOVED***
		w.L.Unlock()
		return false
	***REMOVED***
	w.win += win
	// It is unusual that multiple goroutines would be attempting to reserve
	// window space, but not guaranteed. Use broadcast to notify all waiters
	// that additional window is available.
	w.Broadcast()
	w.L.Unlock()
	return true
***REMOVED***

// close sets the window to closed, so all reservations fail
// immediately.
func (w *window) close() ***REMOVED***
	w.L.Lock()
	w.closed = true
	w.Broadcast()
	w.L.Unlock()
***REMOVED***

// reserve reserves win from the available window capacity.
// If no capacity remains, reserve will block. reserve may
// return less than requested.
func (w *window) reserve(win uint32) (uint32, error) ***REMOVED***
	var err error
	w.L.Lock()
	w.writeWaiters++
	w.Broadcast()
	for w.win == 0 && !w.closed ***REMOVED***
		w.Wait()
	***REMOVED***
	w.writeWaiters--
	if w.win < win ***REMOVED***
		win = w.win
	***REMOVED***
	w.win -= win
	if w.closed ***REMOVED***
		err = io.EOF
	***REMOVED***
	w.L.Unlock()
	return win, err
***REMOVED***

// waitWriterBlocked waits until some goroutine is blocked for further
// writes. It is used in tests only.
func (w *window) waitWriterBlocked() ***REMOVED***
	w.Cond.L.Lock()
	for w.writeWaiters == 0 ***REMOVED***
		w.Cond.Wait()
	***REMOVED***
	w.Cond.L.Unlock()
***REMOVED***
