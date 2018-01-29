// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"bytes"
	"crypto"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"

	"golang.org/x/crypto/ed25519"
)

// These constants represent the algorithm names for key types supported by this
// package.
const (
	KeyAlgoRSA      = "ssh-rsa"
	KeyAlgoDSA      = "ssh-dss"
	KeyAlgoECDSA256 = "ecdsa-sha2-nistp256"
	KeyAlgoECDSA384 = "ecdsa-sha2-nistp384"
	KeyAlgoECDSA521 = "ecdsa-sha2-nistp521"
	KeyAlgoED25519  = "ssh-ed25519"
)

// parsePubKey parses a public key of the given algorithm.
// Use ParsePublicKey for keys with prepended algorithm.
func parsePubKey(in []byte, algo string) (pubKey PublicKey, rest []byte, err error) ***REMOVED***
	switch algo ***REMOVED***
	case KeyAlgoRSA:
		return parseRSA(in)
	case KeyAlgoDSA:
		return parseDSA(in)
	case KeyAlgoECDSA256, KeyAlgoECDSA384, KeyAlgoECDSA521:
		return parseECDSA(in)
	case KeyAlgoED25519:
		return parseED25519(in)
	case CertAlgoRSAv01, CertAlgoDSAv01, CertAlgoECDSA256v01, CertAlgoECDSA384v01, CertAlgoECDSA521v01, CertAlgoED25519v01:
		cert, err := parseCert(in, certToPrivAlgo(algo))
		if err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***
		return cert, nil, nil
	***REMOVED***
	return nil, nil, fmt.Errorf("ssh: unknown key algorithm: %v", algo)
***REMOVED***

// parseAuthorizedKey parses a public key in OpenSSH authorized_keys format
// (see sshd(8) manual page) once the options and key type fields have been
// removed.
func parseAuthorizedKey(in []byte) (out PublicKey, comment string, err error) ***REMOVED***
	in = bytes.TrimSpace(in)

	i := bytes.IndexAny(in, " \t")
	if i == -1 ***REMOVED***
		i = len(in)
	***REMOVED***
	base64Key := in[:i]

	key := make([]byte, base64.StdEncoding.DecodedLen(len(base64Key)))
	n, err := base64.StdEncoding.Decode(key, base64Key)
	if err != nil ***REMOVED***
		return nil, "", err
	***REMOVED***
	key = key[:n]
	out, err = ParsePublicKey(key)
	if err != nil ***REMOVED***
		return nil, "", err
	***REMOVED***
	comment = string(bytes.TrimSpace(in[i:]))
	return out, comment, nil
***REMOVED***

// ParseKnownHosts parses an entry in the format of the known_hosts file.
//
// The known_hosts format is documented in the sshd(8) manual page. This
// function will parse a single entry from in. On successful return, marker
// will contain the optional marker value (i.e. "cert-authority" or "revoked")
// or else be empty, hosts will contain the hosts that this entry matches,
// pubKey will contain the public key and comment will contain any trailing
// comment at the end of the line. See the sshd(8) manual page for the various
// forms that a host string can take.
//
// The unparsed remainder of the input will be returned in rest. This function
// can be called repeatedly to parse multiple entries.
//
// If no entries were found in the input then err will be io.EOF. Otherwise a
// non-nil err value indicates a parse error.
func ParseKnownHosts(in []byte) (marker string, hosts []string, pubKey PublicKey, comment string, rest []byte, err error) ***REMOVED***
	for len(in) > 0 ***REMOVED***
		end := bytes.IndexByte(in, '\n')
		if end != -1 ***REMOVED***
			rest = in[end+1:]
			in = in[:end]
		***REMOVED*** else ***REMOVED***
			rest = nil
		***REMOVED***

		end = bytes.IndexByte(in, '\r')
		if end != -1 ***REMOVED***
			in = in[:end]
		***REMOVED***

		in = bytes.TrimSpace(in)
		if len(in) == 0 || in[0] == '#' ***REMOVED***
			in = rest
			continue
		***REMOVED***

		i := bytes.IndexAny(in, " \t")
		if i == -1 ***REMOVED***
			in = rest
			continue
		***REMOVED***

		// Strip out the beginning of the known_host key.
		// This is either an optional marker or a (set of) hostname(s).
		keyFields := bytes.Fields(in)
		if len(keyFields) < 3 || len(keyFields) > 5 ***REMOVED***
			return "", nil, nil, "", nil, errors.New("ssh: invalid entry in known_hosts data")
		***REMOVED***

		// keyFields[0] is either "@cert-authority", "@revoked" or a comma separated
		// list of hosts
		marker := ""
		if keyFields[0][0] == '@' ***REMOVED***
			marker = string(keyFields[0][1:])
			keyFields = keyFields[1:]
		***REMOVED***

		hosts := string(keyFields[0])
		// keyFields[1] contains the key type (e.g. “ssh-rsa”).
		// However, that information is duplicated inside the
		// base64-encoded key and so is ignored here.

		key := bytes.Join(keyFields[2:], []byte(" "))
		if pubKey, comment, err = parseAuthorizedKey(key); err != nil ***REMOVED***
			return "", nil, nil, "", nil, err
		***REMOVED***

		return marker, strings.Split(hosts, ","), pubKey, comment, rest, nil
	***REMOVED***

	return "", nil, nil, "", nil, io.EOF
***REMOVED***

// ParseAuthorizedKeys parses a public key from an authorized_keys
// file used in OpenSSH according to the sshd(8) manual page.
func ParseAuthorizedKey(in []byte) (out PublicKey, comment string, options []string, rest []byte, err error) ***REMOVED***
	for len(in) > 0 ***REMOVED***
		end := bytes.IndexByte(in, '\n')
		if end != -1 ***REMOVED***
			rest = in[end+1:]
			in = in[:end]
		***REMOVED*** else ***REMOVED***
			rest = nil
		***REMOVED***

		end = bytes.IndexByte(in, '\r')
		if end != -1 ***REMOVED***
			in = in[:end]
		***REMOVED***

		in = bytes.TrimSpace(in)
		if len(in) == 0 || in[0] == '#' ***REMOVED***
			in = rest
			continue
		***REMOVED***

		i := bytes.IndexAny(in, " \t")
		if i == -1 ***REMOVED***
			in = rest
			continue
		***REMOVED***

		if out, comment, err = parseAuthorizedKey(in[i:]); err == nil ***REMOVED***
			return out, comment, options, rest, nil
		***REMOVED***

		// No key type recognised. Maybe there's an options field at
		// the beginning.
		var b byte
		inQuote := false
		var candidateOptions []string
		optionStart := 0
		for i, b = range in ***REMOVED***
			isEnd := !inQuote && (b == ' ' || b == '\t')
			if (b == ',' && !inQuote) || isEnd ***REMOVED***
				if i-optionStart > 0 ***REMOVED***
					candidateOptions = append(candidateOptions, string(in[optionStart:i]))
				***REMOVED***
				optionStart = i + 1
			***REMOVED***
			if isEnd ***REMOVED***
				break
			***REMOVED***
			if b == '"' && (i == 0 || (i > 0 && in[i-1] != '\\')) ***REMOVED***
				inQuote = !inQuote
			***REMOVED***
		***REMOVED***
		for i < len(in) && (in[i] == ' ' || in[i] == '\t') ***REMOVED***
			i++
		***REMOVED***
		if i == len(in) ***REMOVED***
			// Invalid line: unmatched quote
			in = rest
			continue
		***REMOVED***

		in = in[i:]
		i = bytes.IndexAny(in, " \t")
		if i == -1 ***REMOVED***
			in = rest
			continue
		***REMOVED***

		if out, comment, err = parseAuthorizedKey(in[i:]); err == nil ***REMOVED***
			options = candidateOptions
			return out, comment, options, rest, nil
		***REMOVED***

		in = rest
		continue
	***REMOVED***

	return nil, "", nil, nil, errors.New("ssh: no key found")
***REMOVED***

// ParsePublicKey parses an SSH public key formatted for use in
// the SSH wire protocol according to RFC 4253, section 6.6.
func ParsePublicKey(in []byte) (out PublicKey, err error) ***REMOVED***
	algo, in, ok := parseString(in)
	if !ok ***REMOVED***
		return nil, errShortRead
	***REMOVED***
	var rest []byte
	out, rest, err = parsePubKey(in, string(algo))
	if len(rest) > 0 ***REMOVED***
		return nil, errors.New("ssh: trailing junk in public key")
	***REMOVED***

	return out, err
***REMOVED***

// MarshalAuthorizedKey serializes key for inclusion in an OpenSSH
// authorized_keys file. The return value ends with newline.
func MarshalAuthorizedKey(key PublicKey) []byte ***REMOVED***
	b := &bytes.Buffer***REMOVED******REMOVED***
	b.WriteString(key.Type())
	b.WriteByte(' ')
	e := base64.NewEncoder(base64.StdEncoding, b)
	e.Write(key.Marshal())
	e.Close()
	b.WriteByte('\n')
	return b.Bytes()
***REMOVED***

// PublicKey is an abstraction of different types of public keys.
type PublicKey interface ***REMOVED***
	// Type returns the key's type, e.g. "ssh-rsa".
	Type() string

	// Marshal returns the serialized key data in SSH wire format,
	// with the name prefix.
	Marshal() []byte

	// Verify that sig is a signature on the given data using this
	// key. This function will hash the data appropriately first.
	Verify(data []byte, sig *Signature) error
***REMOVED***

// CryptoPublicKey, if implemented by a PublicKey,
// returns the underlying crypto.PublicKey form of the key.
type CryptoPublicKey interface ***REMOVED***
	CryptoPublicKey() crypto.PublicKey
***REMOVED***

// A Signer can create signatures that verify against a public key.
type Signer interface ***REMOVED***
	// PublicKey returns an associated PublicKey instance.
	PublicKey() PublicKey

	// Sign returns raw signature for the given data. This method
	// will apply the hash specified for the keytype to the data.
	Sign(rand io.Reader, data []byte) (*Signature, error)
***REMOVED***

type rsaPublicKey rsa.PublicKey

func (r *rsaPublicKey) Type() string ***REMOVED***
	return "ssh-rsa"
***REMOVED***

// parseRSA parses an RSA key according to RFC 4253, section 6.6.
func parseRSA(in []byte) (out PublicKey, rest []byte, err error) ***REMOVED***
	var w struct ***REMOVED***
		E    *big.Int
		N    *big.Int
		Rest []byte `ssh:"rest"`
	***REMOVED***
	if err := Unmarshal(in, &w); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	if w.E.BitLen() > 24 ***REMOVED***
		return nil, nil, errors.New("ssh: exponent too large")
	***REMOVED***
	e := w.E.Int64()
	if e < 3 || e&1 == 0 ***REMOVED***
		return nil, nil, errors.New("ssh: incorrect exponent")
	***REMOVED***

	var key rsa.PublicKey
	key.E = int(e)
	key.N = w.N
	return (*rsaPublicKey)(&key), w.Rest, nil
***REMOVED***

func (r *rsaPublicKey) Marshal() []byte ***REMOVED***
	e := new(big.Int).SetInt64(int64(r.E))
	// RSA publickey struct layout should match the struct used by
	// parseRSACert in the x/crypto/ssh/agent package.
	wirekey := struct ***REMOVED***
		Name string
		E    *big.Int
		N    *big.Int
	***REMOVED******REMOVED***
		KeyAlgoRSA,
		e,
		r.N,
	***REMOVED***
	return Marshal(&wirekey)
***REMOVED***

func (r *rsaPublicKey) Verify(data []byte, sig *Signature) error ***REMOVED***
	if sig.Format != r.Type() ***REMOVED***
		return fmt.Errorf("ssh: signature type %s for key type %s", sig.Format, r.Type())
	***REMOVED***
	h := crypto.SHA1.New()
	h.Write(data)
	digest := h.Sum(nil)
	return rsa.VerifyPKCS1v15((*rsa.PublicKey)(r), crypto.SHA1, digest, sig.Blob)
***REMOVED***

func (r *rsaPublicKey) CryptoPublicKey() crypto.PublicKey ***REMOVED***
	return (*rsa.PublicKey)(r)
***REMOVED***

type dsaPublicKey dsa.PublicKey

func (k *dsaPublicKey) Type() string ***REMOVED***
	return "ssh-dss"
***REMOVED***

func checkDSAParams(param *dsa.Parameters) error ***REMOVED***
	// SSH specifies FIPS 186-2, which only provided a single size
	// (1024 bits) DSA key. FIPS 186-3 allows for larger key
	// sizes, which would confuse SSH.
	if l := param.P.BitLen(); l != 1024 ***REMOVED***
		return fmt.Errorf("ssh: unsupported DSA key size %d", l)
	***REMOVED***

	return nil
***REMOVED***

// parseDSA parses an DSA key according to RFC 4253, section 6.6.
func parseDSA(in []byte) (out PublicKey, rest []byte, err error) ***REMOVED***
	var w struct ***REMOVED***
		P, Q, G, Y *big.Int
		Rest       []byte `ssh:"rest"`
	***REMOVED***
	if err := Unmarshal(in, &w); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	param := dsa.Parameters***REMOVED***
		P: w.P,
		Q: w.Q,
		G: w.G,
	***REMOVED***
	if err := checkDSAParams(&param); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	key := &dsaPublicKey***REMOVED***
		Parameters: param,
		Y:          w.Y,
	***REMOVED***
	return key, w.Rest, nil
***REMOVED***

func (k *dsaPublicKey) Marshal() []byte ***REMOVED***
	// DSA publickey struct layout should match the struct used by
	// parseDSACert in the x/crypto/ssh/agent package.
	w := struct ***REMOVED***
		Name       string
		P, Q, G, Y *big.Int
	***REMOVED******REMOVED***
		k.Type(),
		k.P,
		k.Q,
		k.G,
		k.Y,
	***REMOVED***

	return Marshal(&w)
***REMOVED***

func (k *dsaPublicKey) Verify(data []byte, sig *Signature) error ***REMOVED***
	if sig.Format != k.Type() ***REMOVED***
		return fmt.Errorf("ssh: signature type %s for key type %s", sig.Format, k.Type())
	***REMOVED***
	h := crypto.SHA1.New()
	h.Write(data)
	digest := h.Sum(nil)

	// Per RFC 4253, section 6.6,
	// The value for 'dss_signature_blob' is encoded as a string containing
	// r, followed by s (which are 160-bit integers, without lengths or
	// padding, unsigned, and in network byte order).
	// For DSS purposes, sig.Blob should be exactly 40 bytes in length.
	if len(sig.Blob) != 40 ***REMOVED***
		return errors.New("ssh: DSA signature parse error")
	***REMOVED***
	r := new(big.Int).SetBytes(sig.Blob[:20])
	s := new(big.Int).SetBytes(sig.Blob[20:])
	if dsa.Verify((*dsa.PublicKey)(k), digest, r, s) ***REMOVED***
		return nil
	***REMOVED***
	return errors.New("ssh: signature did not verify")
***REMOVED***

func (k *dsaPublicKey) CryptoPublicKey() crypto.PublicKey ***REMOVED***
	return (*dsa.PublicKey)(k)
***REMOVED***

type dsaPrivateKey struct ***REMOVED***
	*dsa.PrivateKey
***REMOVED***

func (k *dsaPrivateKey) PublicKey() PublicKey ***REMOVED***
	return (*dsaPublicKey)(&k.PrivateKey.PublicKey)
***REMOVED***

func (k *dsaPrivateKey) Sign(rand io.Reader, data []byte) (*Signature, error) ***REMOVED***
	h := crypto.SHA1.New()
	h.Write(data)
	digest := h.Sum(nil)
	r, s, err := dsa.Sign(rand, k.PrivateKey, digest)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	sig := make([]byte, 40)
	rb := r.Bytes()
	sb := s.Bytes()

	copy(sig[20-len(rb):20], rb)
	copy(sig[40-len(sb):], sb)

	return &Signature***REMOVED***
		Format: k.PublicKey().Type(),
		Blob:   sig,
	***REMOVED***, nil
***REMOVED***

type ecdsaPublicKey ecdsa.PublicKey

func (k *ecdsaPublicKey) Type() string ***REMOVED***
	return "ecdsa-sha2-" + k.nistID()
***REMOVED***

func (k *ecdsaPublicKey) nistID() string ***REMOVED***
	switch k.Params().BitSize ***REMOVED***
	case 256:
		return "nistp256"
	case 384:
		return "nistp384"
	case 521:
		return "nistp521"
	***REMOVED***
	panic("ssh: unsupported ecdsa key size")
***REMOVED***

type ed25519PublicKey ed25519.PublicKey

func (k ed25519PublicKey) Type() string ***REMOVED***
	return KeyAlgoED25519
***REMOVED***

func parseED25519(in []byte) (out PublicKey, rest []byte, err error) ***REMOVED***
	var w struct ***REMOVED***
		KeyBytes []byte
		Rest     []byte `ssh:"rest"`
	***REMOVED***

	if err := Unmarshal(in, &w); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	key := ed25519.PublicKey(w.KeyBytes)

	return (ed25519PublicKey)(key), w.Rest, nil
***REMOVED***

func (k ed25519PublicKey) Marshal() []byte ***REMOVED***
	w := struct ***REMOVED***
		Name     string
		KeyBytes []byte
	***REMOVED******REMOVED***
		KeyAlgoED25519,
		[]byte(k),
	***REMOVED***
	return Marshal(&w)
***REMOVED***

func (k ed25519PublicKey) Verify(b []byte, sig *Signature) error ***REMOVED***
	if sig.Format != k.Type() ***REMOVED***
		return fmt.Errorf("ssh: signature type %s for key type %s", sig.Format, k.Type())
	***REMOVED***

	edKey := (ed25519.PublicKey)(k)
	if ok := ed25519.Verify(edKey, b, sig.Blob); !ok ***REMOVED***
		return errors.New("ssh: signature did not verify")
	***REMOVED***

	return nil
***REMOVED***

func (k ed25519PublicKey) CryptoPublicKey() crypto.PublicKey ***REMOVED***
	return ed25519.PublicKey(k)
***REMOVED***

func supportedEllipticCurve(curve elliptic.Curve) bool ***REMOVED***
	return curve == elliptic.P256() || curve == elliptic.P384() || curve == elliptic.P521()
***REMOVED***

// ecHash returns the hash to match the given elliptic curve, see RFC
// 5656, section 6.2.1
func ecHash(curve elliptic.Curve) crypto.Hash ***REMOVED***
	bitSize := curve.Params().BitSize
	switch ***REMOVED***
	case bitSize <= 256:
		return crypto.SHA256
	case bitSize <= 384:
		return crypto.SHA384
	***REMOVED***
	return crypto.SHA512
***REMOVED***

// parseECDSA parses an ECDSA key according to RFC 5656, section 3.1.
func parseECDSA(in []byte) (out PublicKey, rest []byte, err error) ***REMOVED***
	var w struct ***REMOVED***
		Curve    string
		KeyBytes []byte
		Rest     []byte `ssh:"rest"`
	***REMOVED***

	if err := Unmarshal(in, &w); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	key := new(ecdsa.PublicKey)

	switch w.Curve ***REMOVED***
	case "nistp256":
		key.Curve = elliptic.P256()
	case "nistp384":
		key.Curve = elliptic.P384()
	case "nistp521":
		key.Curve = elliptic.P521()
	default:
		return nil, nil, errors.New("ssh: unsupported curve")
	***REMOVED***

	key.X, key.Y = elliptic.Unmarshal(key.Curve, w.KeyBytes)
	if key.X == nil || key.Y == nil ***REMOVED***
		return nil, nil, errors.New("ssh: invalid curve point")
	***REMOVED***
	return (*ecdsaPublicKey)(key), w.Rest, nil
***REMOVED***

func (k *ecdsaPublicKey) Marshal() []byte ***REMOVED***
	// See RFC 5656, section 3.1.
	keyBytes := elliptic.Marshal(k.Curve, k.X, k.Y)
	// ECDSA publickey struct layout should match the struct used by
	// parseECDSACert in the x/crypto/ssh/agent package.
	w := struct ***REMOVED***
		Name string
		ID   string
		Key  []byte
	***REMOVED******REMOVED***
		k.Type(),
		k.nistID(),
		keyBytes,
	***REMOVED***

	return Marshal(&w)
***REMOVED***

func (k *ecdsaPublicKey) Verify(data []byte, sig *Signature) error ***REMOVED***
	if sig.Format != k.Type() ***REMOVED***
		return fmt.Errorf("ssh: signature type %s for key type %s", sig.Format, k.Type())
	***REMOVED***

	h := ecHash(k.Curve).New()
	h.Write(data)
	digest := h.Sum(nil)

	// Per RFC 5656, section 3.1.2,
	// The ecdsa_signature_blob value has the following specific encoding:
	//    mpint    r
	//    mpint    s
	var ecSig struct ***REMOVED***
		R *big.Int
		S *big.Int
	***REMOVED***

	if err := Unmarshal(sig.Blob, &ecSig); err != nil ***REMOVED***
		return err
	***REMOVED***

	if ecdsa.Verify((*ecdsa.PublicKey)(k), digest, ecSig.R, ecSig.S) ***REMOVED***
		return nil
	***REMOVED***
	return errors.New("ssh: signature did not verify")
***REMOVED***

func (k *ecdsaPublicKey) CryptoPublicKey() crypto.PublicKey ***REMOVED***
	return (*ecdsa.PublicKey)(k)
***REMOVED***

// NewSignerFromKey takes an *rsa.PrivateKey, *dsa.PrivateKey,
// *ecdsa.PrivateKey or any other crypto.Signer and returns a
// corresponding Signer instance. ECDSA keys must use P-256, P-384 or
// P-521. DSA keys must use parameter size L1024N160.
func NewSignerFromKey(key interface***REMOVED******REMOVED***) (Signer, error) ***REMOVED***
	switch key := key.(type) ***REMOVED***
	case crypto.Signer:
		return NewSignerFromSigner(key)
	case *dsa.PrivateKey:
		return newDSAPrivateKey(key)
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %T", key)
	***REMOVED***
***REMOVED***

func newDSAPrivateKey(key *dsa.PrivateKey) (Signer, error) ***REMOVED***
	if err := checkDSAParams(&key.PublicKey.Parameters); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &dsaPrivateKey***REMOVED***key***REMOVED***, nil
***REMOVED***

type wrappedSigner struct ***REMOVED***
	signer crypto.Signer
	pubKey PublicKey
***REMOVED***

// NewSignerFromSigner takes any crypto.Signer implementation and
// returns a corresponding Signer interface. This can be used, for
// example, with keys kept in hardware modules.
func NewSignerFromSigner(signer crypto.Signer) (Signer, error) ***REMOVED***
	pubKey, err := NewPublicKey(signer.Public())
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &wrappedSigner***REMOVED***signer, pubKey***REMOVED***, nil
***REMOVED***

func (s *wrappedSigner) PublicKey() PublicKey ***REMOVED***
	return s.pubKey
***REMOVED***

func (s *wrappedSigner) Sign(rand io.Reader, data []byte) (*Signature, error) ***REMOVED***
	var hashFunc crypto.Hash

	switch key := s.pubKey.(type) ***REMOVED***
	case *rsaPublicKey, *dsaPublicKey:
		hashFunc = crypto.SHA1
	case *ecdsaPublicKey:
		hashFunc = ecHash(key.Curve)
	case ed25519PublicKey:
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %T", key)
	***REMOVED***

	var digest []byte
	if hashFunc != 0 ***REMOVED***
		h := hashFunc.New()
		h.Write(data)
		digest = h.Sum(nil)
	***REMOVED*** else ***REMOVED***
		digest = data
	***REMOVED***

	signature, err := s.signer.Sign(rand, digest, hashFunc)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// crypto.Signer.Sign is expected to return an ASN.1-encoded signature
	// for ECDSA and DSA, but that's not the encoding expected by SSH, so
	// re-encode.
	switch s.pubKey.(type) ***REMOVED***
	case *ecdsaPublicKey, *dsaPublicKey:
		type asn1Signature struct ***REMOVED***
			R, S *big.Int
		***REMOVED***
		asn1Sig := new(asn1Signature)
		_, err := asn1.Unmarshal(signature, asn1Sig)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		switch s.pubKey.(type) ***REMOVED***
		case *ecdsaPublicKey:
			signature = Marshal(asn1Sig)

		case *dsaPublicKey:
			signature = make([]byte, 40)
			r := asn1Sig.R.Bytes()
			s := asn1Sig.S.Bytes()
			copy(signature[20-len(r):20], r)
			copy(signature[40-len(s):40], s)
		***REMOVED***
	***REMOVED***

	return &Signature***REMOVED***
		Format: s.pubKey.Type(),
		Blob:   signature,
	***REMOVED***, nil
***REMOVED***

// NewPublicKey takes an *rsa.PublicKey, *dsa.PublicKey, *ecdsa.PublicKey,
// or ed25519.PublicKey returns a corresponding PublicKey instance.
// ECDSA keys must use P-256, P-384 or P-521.
func NewPublicKey(key interface***REMOVED******REMOVED***) (PublicKey, error) ***REMOVED***
	switch key := key.(type) ***REMOVED***
	case *rsa.PublicKey:
		return (*rsaPublicKey)(key), nil
	case *ecdsa.PublicKey:
		if !supportedEllipticCurve(key.Curve) ***REMOVED***
			return nil, errors.New("ssh: only P-256, P-384 and P-521 EC keys are supported")
		***REMOVED***
		return (*ecdsaPublicKey)(key), nil
	case *dsa.PublicKey:
		return (*dsaPublicKey)(key), nil
	case ed25519.PublicKey:
		return (ed25519PublicKey)(key), nil
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %T", key)
	***REMOVED***
***REMOVED***

// ParsePrivateKey returns a Signer from a PEM encoded private key. It supports
// the same keys as ParseRawPrivateKey.
func ParsePrivateKey(pemBytes []byte) (Signer, error) ***REMOVED***
	key, err := ParseRawPrivateKey(pemBytes)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return NewSignerFromKey(key)
***REMOVED***

// ParsePrivateKeyWithPassphrase returns a Signer from a PEM encoded private
// key and passphrase. It supports the same keys as
// ParseRawPrivateKeyWithPassphrase.
func ParsePrivateKeyWithPassphrase(pemBytes, passPhrase []byte) (Signer, error) ***REMOVED***
	key, err := ParseRawPrivateKeyWithPassphrase(pemBytes, passPhrase)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return NewSignerFromKey(key)
***REMOVED***

// encryptedBlock tells whether a private key is
// encrypted by examining its Proc-Type header
// for a mention of ENCRYPTED
// according to RFC 1421 Section 4.6.1.1.
func encryptedBlock(block *pem.Block) bool ***REMOVED***
	return strings.Contains(block.Headers["Proc-Type"], "ENCRYPTED")
***REMOVED***

// ParseRawPrivateKey returns a private key from a PEM encoded private key. It
// supports RSA (PKCS#1), DSA (OpenSSL), and ECDSA private keys.
func ParseRawPrivateKey(pemBytes []byte) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	block, _ := pem.Decode(pemBytes)
	if block == nil ***REMOVED***
		return nil, errors.New("ssh: no key found")
	***REMOVED***

	if encryptedBlock(block) ***REMOVED***
		return nil, errors.New("ssh: cannot decode encrypted private keys")
	***REMOVED***

	switch block.Type ***REMOVED***
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(block.Bytes)
	case "DSA PRIVATE KEY":
		return ParseDSAPrivateKey(block.Bytes)
	case "OPENSSH PRIVATE KEY":
		return parseOpenSSHPrivateKey(block.Bytes)
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	***REMOVED***
***REMOVED***

// ParseRawPrivateKeyWithPassphrase returns a private key decrypted with
// passphrase from a PEM encoded private key. If wrong passphrase, return
// x509.IncorrectPasswordError.
func ParseRawPrivateKeyWithPassphrase(pemBytes, passPhrase []byte) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	block, _ := pem.Decode(pemBytes)
	if block == nil ***REMOVED***
		return nil, errors.New("ssh: no key found")
	***REMOVED***
	buf := block.Bytes

	if encryptedBlock(block) ***REMOVED***
		if x509.IsEncryptedPEMBlock(block) ***REMOVED***
			var err error
			buf, err = x509.DecryptPEMBlock(block, passPhrase)
			if err != nil ***REMOVED***
				if err == x509.IncorrectPasswordError ***REMOVED***
					return nil, err
				***REMOVED***
				return nil, fmt.Errorf("ssh: cannot decode encrypted private keys: %v", err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	switch block.Type ***REMOVED***
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(buf)
	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(buf)
	case "DSA PRIVATE KEY":
		return ParseDSAPrivateKey(buf)
	case "OPENSSH PRIVATE KEY":
		return parseOpenSSHPrivateKey(buf)
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	***REMOVED***
***REMOVED***

// ParseDSAPrivateKey returns a DSA private key from its ASN.1 DER encoding, as
// specified by the OpenSSL DSA man page.
func ParseDSAPrivateKey(der []byte) (*dsa.PrivateKey, error) ***REMOVED***
	var k struct ***REMOVED***
		Version int
		P       *big.Int
		Q       *big.Int
		G       *big.Int
		Pub     *big.Int
		Priv    *big.Int
	***REMOVED***
	rest, err := asn1.Unmarshal(der, &k)
	if err != nil ***REMOVED***
		return nil, errors.New("ssh: failed to parse DSA key: " + err.Error())
	***REMOVED***
	if len(rest) > 0 ***REMOVED***
		return nil, errors.New("ssh: garbage after DSA key")
	***REMOVED***

	return &dsa.PrivateKey***REMOVED***
		PublicKey: dsa.PublicKey***REMOVED***
			Parameters: dsa.Parameters***REMOVED***
				P: k.P,
				Q: k.Q,
				G: k.G,
			***REMOVED***,
			Y: k.Pub,
		***REMOVED***,
		X: k.Priv,
	***REMOVED***, nil
***REMOVED***

// Implemented based on the documentation at
// https://github.com/openssh/openssh-portable/blob/master/PROTOCOL.key
func parseOpenSSHPrivateKey(key []byte) (crypto.PrivateKey, error) ***REMOVED***
	magic := append([]byte("openssh-key-v1"), 0)
	if !bytes.Equal(magic, key[0:len(magic)]) ***REMOVED***
		return nil, errors.New("ssh: invalid openssh private key format")
	***REMOVED***
	remaining := key[len(magic):]

	var w struct ***REMOVED***
		CipherName   string
		KdfName      string
		KdfOpts      string
		NumKeys      uint32
		PubKey       []byte
		PrivKeyBlock []byte
	***REMOVED***

	if err := Unmarshal(remaining, &w); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if w.KdfName != "none" || w.CipherName != "none" ***REMOVED***
		return nil, errors.New("ssh: cannot decode encrypted private keys")
	***REMOVED***

	pk1 := struct ***REMOVED***
		Check1  uint32
		Check2  uint32
		Keytype string
		Rest    []byte `ssh:"rest"`
	***REMOVED******REMOVED******REMOVED***

	if err := Unmarshal(w.PrivKeyBlock, &pk1); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if pk1.Check1 != pk1.Check2 ***REMOVED***
		return nil, errors.New("ssh: checkint mismatch")
	***REMOVED***

	// we only handle ed25519 and rsa keys currently
	switch pk1.Keytype ***REMOVED***
	case KeyAlgoRSA:
		// https://github.com/openssh/openssh-portable/blob/master/sshkey.c#L2760-L2773
		key := struct ***REMOVED***
			N       *big.Int
			E       *big.Int
			D       *big.Int
			Iqmp    *big.Int
			P       *big.Int
			Q       *big.Int
			Comment string
			Pad     []byte `ssh:"rest"`
		***REMOVED******REMOVED******REMOVED***

		if err := Unmarshal(pk1.Rest, &key); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		for i, b := range key.Pad ***REMOVED***
			if int(b) != i+1 ***REMOVED***
				return nil, errors.New("ssh: padding not as expected")
			***REMOVED***
		***REMOVED***

		pk := &rsa.PrivateKey***REMOVED***
			PublicKey: rsa.PublicKey***REMOVED***
				N: key.N,
				E: int(key.E.Int64()),
			***REMOVED***,
			D:      key.D,
			Primes: []*big.Int***REMOVED***key.P, key.Q***REMOVED***,
		***REMOVED***

		if err := pk.Validate(); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		pk.Precompute()

		return pk, nil
	case KeyAlgoED25519:
		key := struct ***REMOVED***
			Pub     []byte
			Priv    []byte
			Comment string
			Pad     []byte `ssh:"rest"`
		***REMOVED******REMOVED******REMOVED***

		if err := Unmarshal(pk1.Rest, &key); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		if len(key.Priv) != ed25519.PrivateKeySize ***REMOVED***
			return nil, errors.New("ssh: private key unexpected length")
		***REMOVED***

		for i, b := range key.Pad ***REMOVED***
			if int(b) != i+1 ***REMOVED***
				return nil, errors.New("ssh: padding not as expected")
			***REMOVED***
		***REMOVED***

		pk := ed25519.PrivateKey(make([]byte, ed25519.PrivateKeySize))
		copy(pk, key.Priv)
		return &pk, nil
	default:
		return nil, errors.New("ssh: unhandled key type")
	***REMOVED***
***REMOVED***

// FingerprintLegacyMD5 returns the user presentation of the key's
// fingerprint as described by RFC 4716 section 4.
func FingerprintLegacyMD5(pubKey PublicKey) string ***REMOVED***
	md5sum := md5.Sum(pubKey.Marshal())
	hexarray := make([]string, len(md5sum))
	for i, c := range md5sum ***REMOVED***
		hexarray[i] = hex.EncodeToString([]byte***REMOVED***c***REMOVED***)
	***REMOVED***
	return strings.Join(hexarray, ":")
***REMOVED***

// FingerprintSHA256 returns the user presentation of the key's
// fingerprint as unpadded base64 encoded sha256 hash.
// This format was introduced from OpenSSH 6.8.
// https://www.openssh.com/txt/release-6.8
// https://tools.ietf.org/html/rfc4648#section-3.2 (unpadded base64 encoding)
func FingerprintSHA256(pubKey PublicKey) string ***REMOVED***
	sha256sum := sha256.Sum256(pubKey.Marshal())
	hash := base64.RawStdEncoding.EncodeToString(sha256sum[:])
	return "SHA256:" + hash
***REMOVED***
