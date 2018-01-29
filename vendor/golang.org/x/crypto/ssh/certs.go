// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"sort"
	"time"
)

// These constants from [PROTOCOL.certkeys] represent the algorithm names
// for certificate types supported by this package.
const (
	CertAlgoRSAv01      = "ssh-rsa-cert-v01@openssh.com"
	CertAlgoDSAv01      = "ssh-dss-cert-v01@openssh.com"
	CertAlgoECDSA256v01 = "ecdsa-sha2-nistp256-cert-v01@openssh.com"
	CertAlgoECDSA384v01 = "ecdsa-sha2-nistp384-cert-v01@openssh.com"
	CertAlgoECDSA521v01 = "ecdsa-sha2-nistp521-cert-v01@openssh.com"
	CertAlgoED25519v01  = "ssh-ed25519-cert-v01@openssh.com"
)

// Certificate types distinguish between host and user
// certificates. The values can be set in the CertType field of
// Certificate.
const (
	UserCert = 1
	HostCert = 2
)

// Signature represents a cryptographic signature.
type Signature struct ***REMOVED***
	Format string
	Blob   []byte
***REMOVED***

// CertTimeInfinity can be used for OpenSSHCertV01.ValidBefore to indicate that
// a certificate does not expire.
const CertTimeInfinity = 1<<64 - 1

// An Certificate represents an OpenSSH certificate as defined in
// [PROTOCOL.certkeys]?rev=1.8.
type Certificate struct ***REMOVED***
	Nonce           []byte
	Key             PublicKey
	Serial          uint64
	CertType        uint32
	KeyId           string
	ValidPrincipals []string
	ValidAfter      uint64
	ValidBefore     uint64
	Permissions
	Reserved     []byte
	SignatureKey PublicKey
	Signature    *Signature
***REMOVED***

// genericCertData holds the key-independent part of the certificate data.
// Overall, certificates contain an nonce, public key fields and
// key-independent fields.
type genericCertData struct ***REMOVED***
	Serial          uint64
	CertType        uint32
	KeyId           string
	ValidPrincipals []byte
	ValidAfter      uint64
	ValidBefore     uint64
	CriticalOptions []byte
	Extensions      []byte
	Reserved        []byte
	SignatureKey    []byte
	Signature       []byte
***REMOVED***

func marshalStringList(namelist []string) []byte ***REMOVED***
	var to []byte
	for _, name := range namelist ***REMOVED***
		s := struct***REMOVED*** N string ***REMOVED******REMOVED***name***REMOVED***
		to = append(to, Marshal(&s)...)
	***REMOVED***
	return to
***REMOVED***

type optionsTuple struct ***REMOVED***
	Key   string
	Value []byte
***REMOVED***

type optionsTupleValue struct ***REMOVED***
	Value string
***REMOVED***

// serialize a map of critical options or extensions
// issue #10569 - per [PROTOCOL.certkeys] and SSH implementation,
// we need two length prefixes for a non-empty string value
func marshalTuples(tups map[string]string) []byte ***REMOVED***
	keys := make([]string, 0, len(tups))
	for key := range tups ***REMOVED***
		keys = append(keys, key)
	***REMOVED***
	sort.Strings(keys)

	var ret []byte
	for _, key := range keys ***REMOVED***
		s := optionsTuple***REMOVED***Key: key***REMOVED***
		if value := tups[key]; len(value) > 0 ***REMOVED***
			s.Value = Marshal(&optionsTupleValue***REMOVED***value***REMOVED***)
		***REMOVED***
		ret = append(ret, Marshal(&s)...)
	***REMOVED***
	return ret
***REMOVED***

// issue #10569 - per [PROTOCOL.certkeys] and SSH implementation,
// we need two length prefixes for a non-empty option value
func parseTuples(in []byte) (map[string]string, error) ***REMOVED***
	tups := map[string]string***REMOVED******REMOVED***
	var lastKey string
	var haveLastKey bool

	for len(in) > 0 ***REMOVED***
		var key, val, extra []byte
		var ok bool

		if key, in, ok = parseString(in); !ok ***REMOVED***
			return nil, errShortRead
		***REMOVED***
		keyStr := string(key)
		// according to [PROTOCOL.certkeys], the names must be in
		// lexical order.
		if haveLastKey && keyStr <= lastKey ***REMOVED***
			return nil, fmt.Errorf("ssh: certificate options are not in lexical order")
		***REMOVED***
		lastKey, haveLastKey = keyStr, true
		// the next field is a data field, which if non-empty has a string embedded
		if val, in, ok = parseString(in); !ok ***REMOVED***
			return nil, errShortRead
		***REMOVED***
		if len(val) > 0 ***REMOVED***
			val, extra, ok = parseString(val)
			if !ok ***REMOVED***
				return nil, errShortRead
			***REMOVED***
			if len(extra) > 0 ***REMOVED***
				return nil, fmt.Errorf("ssh: unexpected trailing data after certificate option value")
			***REMOVED***
			tups[keyStr] = string(val)
		***REMOVED*** else ***REMOVED***
			tups[keyStr] = ""
		***REMOVED***
	***REMOVED***
	return tups, nil
***REMOVED***

func parseCert(in []byte, privAlgo string) (*Certificate, error) ***REMOVED***
	nonce, rest, ok := parseString(in)
	if !ok ***REMOVED***
		return nil, errShortRead
	***REMOVED***

	key, rest, err := parsePubKey(rest, privAlgo)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var g genericCertData
	if err := Unmarshal(rest, &g); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	c := &Certificate***REMOVED***
		Nonce:       nonce,
		Key:         key,
		Serial:      g.Serial,
		CertType:    g.CertType,
		KeyId:       g.KeyId,
		ValidAfter:  g.ValidAfter,
		ValidBefore: g.ValidBefore,
	***REMOVED***

	for principals := g.ValidPrincipals; len(principals) > 0; ***REMOVED***
		principal, rest, ok := parseString(principals)
		if !ok ***REMOVED***
			return nil, errShortRead
		***REMOVED***
		c.ValidPrincipals = append(c.ValidPrincipals, string(principal))
		principals = rest
	***REMOVED***

	c.CriticalOptions, err = parseTuples(g.CriticalOptions)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	c.Extensions, err = parseTuples(g.Extensions)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	c.Reserved = g.Reserved
	k, err := ParsePublicKey(g.SignatureKey)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	c.SignatureKey = k
	c.Signature, rest, ok = parseSignatureBody(g.Signature)
	if !ok || len(rest) > 0 ***REMOVED***
		return nil, errors.New("ssh: signature parse error")
	***REMOVED***

	return c, nil
***REMOVED***

type openSSHCertSigner struct ***REMOVED***
	pub    *Certificate
	signer Signer
***REMOVED***

// NewCertSigner returns a Signer that signs with the given Certificate, whose
// private key is held by signer. It returns an error if the public key in cert
// doesn't match the key used by signer.
func NewCertSigner(cert *Certificate, signer Signer) (Signer, error) ***REMOVED***
	if bytes.Compare(cert.Key.Marshal(), signer.PublicKey().Marshal()) != 0 ***REMOVED***
		return nil, errors.New("ssh: signer and cert have different public key")
	***REMOVED***

	return &openSSHCertSigner***REMOVED***cert, signer***REMOVED***, nil
***REMOVED***

func (s *openSSHCertSigner) Sign(rand io.Reader, data []byte) (*Signature, error) ***REMOVED***
	return s.signer.Sign(rand, data)
***REMOVED***

func (s *openSSHCertSigner) PublicKey() PublicKey ***REMOVED***
	return s.pub
***REMOVED***

const sourceAddressCriticalOption = "source-address"

// CertChecker does the work of verifying a certificate. Its methods
// can be plugged into ClientConfig.HostKeyCallback and
// ServerConfig.PublicKeyCallback. For the CertChecker to work,
// minimally, the IsAuthority callback should be set.
type CertChecker struct ***REMOVED***
	// SupportedCriticalOptions lists the CriticalOptions that the
	// server application layer understands. These are only used
	// for user certificates.
	SupportedCriticalOptions []string

	// IsUserAuthority should return true if the key is recognized as an
	// authority for the given user certificate. This allows for
	// certificates to be signed by other certificates. This must be set
	// if this CertChecker will be checking user certificates.
	IsUserAuthority func(auth PublicKey) bool

	// IsHostAuthority should report whether the key is recognized as
	// an authority for this host. This allows for certificates to be
	// signed by other keys, and for those other keys to only be valid
	// signers for particular hostnames. This must be set if this
	// CertChecker will be checking host certificates.
	IsHostAuthority func(auth PublicKey, address string) bool

	// Clock is used for verifying time stamps. If nil, time.Now
	// is used.
	Clock func() time.Time

	// UserKeyFallback is called when CertChecker.Authenticate encounters a
	// public key that is not a certificate. It must implement validation
	// of user keys or else, if nil, all such keys are rejected.
	UserKeyFallback func(conn ConnMetadata, key PublicKey) (*Permissions, error)

	// HostKeyFallback is called when CertChecker.CheckHostKey encounters a
	// public key that is not a certificate. It must implement host key
	// validation or else, if nil, all such keys are rejected.
	HostKeyFallback HostKeyCallback

	// IsRevoked is called for each certificate so that revocation checking
	// can be implemented. It should return true if the given certificate
	// is revoked and false otherwise. If nil, no certificates are
	// considered to have been revoked.
	IsRevoked func(cert *Certificate) bool
***REMOVED***

// CheckHostKey checks a host key certificate. This method can be
// plugged into ClientConfig.HostKeyCallback.
func (c *CertChecker) CheckHostKey(addr string, remote net.Addr, key PublicKey) error ***REMOVED***
	cert, ok := key.(*Certificate)
	if !ok ***REMOVED***
		if c.HostKeyFallback != nil ***REMOVED***
			return c.HostKeyFallback(addr, remote, key)
		***REMOVED***
		return errors.New("ssh: non-certificate host key")
	***REMOVED***
	if cert.CertType != HostCert ***REMOVED***
		return fmt.Errorf("ssh: certificate presented as a host key has type %d", cert.CertType)
	***REMOVED***
	if !c.IsHostAuthority(cert.SignatureKey, addr) ***REMOVED***
		return fmt.Errorf("ssh: no authorities for hostname: %v", addr)
	***REMOVED***

	hostname, _, err := net.SplitHostPort(addr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Pass hostname only as principal for host certificates (consistent with OpenSSH)
	return c.CheckCert(hostname, cert)
***REMOVED***

// Authenticate checks a user certificate. Authenticate can be used as
// a value for ServerConfig.PublicKeyCallback.
func (c *CertChecker) Authenticate(conn ConnMetadata, pubKey PublicKey) (*Permissions, error) ***REMOVED***
	cert, ok := pubKey.(*Certificate)
	if !ok ***REMOVED***
		if c.UserKeyFallback != nil ***REMOVED***
			return c.UserKeyFallback(conn, pubKey)
		***REMOVED***
		return nil, errors.New("ssh: normal key pairs not accepted")
	***REMOVED***

	if cert.CertType != UserCert ***REMOVED***
		return nil, fmt.Errorf("ssh: cert has type %d", cert.CertType)
	***REMOVED***
	if !c.IsUserAuthority(cert.SignatureKey) ***REMOVED***
		return nil, fmt.Errorf("ssh: certificate signed by unrecognized authority")
	***REMOVED***

	if err := c.CheckCert(conn.User(), cert); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &cert.Permissions, nil
***REMOVED***

// CheckCert checks CriticalOptions, ValidPrincipals, revocation, timestamp and
// the signature of the certificate.
func (c *CertChecker) CheckCert(principal string, cert *Certificate) error ***REMOVED***
	if c.IsRevoked != nil && c.IsRevoked(cert) ***REMOVED***
		return fmt.Errorf("ssh: certificate serial %d revoked", cert.Serial)
	***REMOVED***

	for opt := range cert.CriticalOptions ***REMOVED***
		// sourceAddressCriticalOption will be enforced by
		// serverAuthenticate
		if opt == sourceAddressCriticalOption ***REMOVED***
			continue
		***REMOVED***

		found := false
		for _, supp := range c.SupportedCriticalOptions ***REMOVED***
			if supp == opt ***REMOVED***
				found = true
				break
			***REMOVED***
		***REMOVED***
		if !found ***REMOVED***
			return fmt.Errorf("ssh: unsupported critical option %q in certificate", opt)
		***REMOVED***
	***REMOVED***

	if len(cert.ValidPrincipals) > 0 ***REMOVED***
		// By default, certs are valid for all users/hosts.
		found := false
		for _, p := range cert.ValidPrincipals ***REMOVED***
			if p == principal ***REMOVED***
				found = true
				break
			***REMOVED***
		***REMOVED***
		if !found ***REMOVED***
			return fmt.Errorf("ssh: principal %q not in the set of valid principals for given certificate: %q", principal, cert.ValidPrincipals)
		***REMOVED***
	***REMOVED***

	clock := c.Clock
	if clock == nil ***REMOVED***
		clock = time.Now
	***REMOVED***

	unixNow := clock().Unix()
	if after := int64(cert.ValidAfter); after < 0 || unixNow < int64(cert.ValidAfter) ***REMOVED***
		return fmt.Errorf("ssh: cert is not yet valid")
	***REMOVED***
	if before := int64(cert.ValidBefore); cert.ValidBefore != uint64(CertTimeInfinity) && (unixNow >= before || before < 0) ***REMOVED***
		return fmt.Errorf("ssh: cert has expired")
	***REMOVED***
	if err := cert.SignatureKey.Verify(cert.bytesForSigning(), cert.Signature); err != nil ***REMOVED***
		return fmt.Errorf("ssh: certificate signature does not verify")
	***REMOVED***

	return nil
***REMOVED***

// SignCert sets c.SignatureKey to the authority's public key and stores a
// Signature, by authority, in the certificate.
func (c *Certificate) SignCert(rand io.Reader, authority Signer) error ***REMOVED***
	c.Nonce = make([]byte, 32)
	if _, err := io.ReadFull(rand, c.Nonce); err != nil ***REMOVED***
		return err
	***REMOVED***
	c.SignatureKey = authority.PublicKey()

	sig, err := authority.Sign(rand, c.bytesForSigning())
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	c.Signature = sig
	return nil
***REMOVED***

var certAlgoNames = map[string]string***REMOVED***
	KeyAlgoRSA:      CertAlgoRSAv01,
	KeyAlgoDSA:      CertAlgoDSAv01,
	KeyAlgoECDSA256: CertAlgoECDSA256v01,
	KeyAlgoECDSA384: CertAlgoECDSA384v01,
	KeyAlgoECDSA521: CertAlgoECDSA521v01,
	KeyAlgoED25519:  CertAlgoED25519v01,
***REMOVED***

// certToPrivAlgo returns the underlying algorithm for a certificate algorithm.
// Panics if a non-certificate algorithm is passed.
func certToPrivAlgo(algo string) string ***REMOVED***
	for privAlgo, pubAlgo := range certAlgoNames ***REMOVED***
		if pubAlgo == algo ***REMOVED***
			return privAlgo
		***REMOVED***
	***REMOVED***
	panic("unknown cert algorithm")
***REMOVED***

func (cert *Certificate) bytesForSigning() []byte ***REMOVED***
	c2 := *cert
	c2.Signature = nil
	out := c2.Marshal()
	// Drop trailing signature length.
	return out[:len(out)-4]
***REMOVED***

// Marshal serializes c into OpenSSH's wire format. It is part of the
// PublicKey interface.
func (c *Certificate) Marshal() []byte ***REMOVED***
	generic := genericCertData***REMOVED***
		Serial:          c.Serial,
		CertType:        c.CertType,
		KeyId:           c.KeyId,
		ValidPrincipals: marshalStringList(c.ValidPrincipals),
		ValidAfter:      uint64(c.ValidAfter),
		ValidBefore:     uint64(c.ValidBefore),
		CriticalOptions: marshalTuples(c.CriticalOptions),
		Extensions:      marshalTuples(c.Extensions),
		Reserved:        c.Reserved,
		SignatureKey:    c.SignatureKey.Marshal(),
	***REMOVED***
	if c.Signature != nil ***REMOVED***
		generic.Signature = Marshal(c.Signature)
	***REMOVED***
	genericBytes := Marshal(&generic)
	keyBytes := c.Key.Marshal()
	_, keyBytes, _ = parseString(keyBytes)
	prefix := Marshal(&struct ***REMOVED***
		Name  string
		Nonce []byte
		Key   []byte `ssh:"rest"`
	***REMOVED******REMOVED***c.Type(), c.Nonce, keyBytes***REMOVED***)

	result := make([]byte, 0, len(prefix)+len(genericBytes))
	result = append(result, prefix...)
	result = append(result, genericBytes...)
	return result
***REMOVED***

// Type returns the key name. It is part of the PublicKey interface.
func (c *Certificate) Type() string ***REMOVED***
	algo, ok := certAlgoNames[c.Key.Type()]
	if !ok ***REMOVED***
		panic("unknown cert key type " + c.Key.Type())
	***REMOVED***
	return algo
***REMOVED***

// Verify verifies a signature against the certificate's public
// key. It is part of the PublicKey interface.
func (c *Certificate) Verify(data []byte, sig *Signature) error ***REMOVED***
	return c.Key.Verify(data, sig)
***REMOVED***

func parseSignatureBody(in []byte) (out *Signature, rest []byte, ok bool) ***REMOVED***
	format, in, ok := parseString(in)
	if !ok ***REMOVED***
		return
	***REMOVED***

	out = &Signature***REMOVED***
		Format: string(format),
	***REMOVED***

	if out.Blob, in, ok = parseString(in); !ok ***REMOVED***
		return
	***REMOVED***

	return out, in, ok
***REMOVED***

func parseSignature(in []byte) (out *Signature, rest []byte, ok bool) ***REMOVED***
	sigBytes, rest, ok := parseString(in)
	if !ok ***REMOVED***
		return
	***REMOVED***

	out, trailing, ok := parseSignatureBody(sigBytes)
	if !ok || len(trailing) > 0 ***REMOVED***
		return nil, nil, false
	***REMOVED***
	return
***REMOVED***
