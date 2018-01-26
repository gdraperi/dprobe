// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x509

import (
	"crypto/rsa"
	// START CT CHANGES
	"github.com/google/certificate-transparency/go/asn1"
	// END CT CHANGES
	"errors"
	"math/big"
)

// pkcs1PrivateKey is a structure which mirrors the PKCS#1 ASN.1 for an RSA private key.
type pkcs1PrivateKey struct ***REMOVED***
	Version int
	N       *big.Int
	E       int
	D       *big.Int
	P       *big.Int
	Q       *big.Int
	// We ignore these values, if present, because rsa will calculate them.
	Dp   *big.Int `asn1:"optional"`
	Dq   *big.Int `asn1:"optional"`
	Qinv *big.Int `asn1:"optional"`

	AdditionalPrimes []pkcs1AdditionalRSAPrime `asn1:"optional,omitempty"`
***REMOVED***

type pkcs1AdditionalRSAPrime struct ***REMOVED***
	Prime *big.Int

	// We ignore these values because rsa will calculate them.
	Exp   *big.Int
	Coeff *big.Int
***REMOVED***

// ParsePKCS1PrivateKey returns an RSA private key from its ASN.1 PKCS#1 DER encoded form.
func ParsePKCS1PrivateKey(der []byte) (key *rsa.PrivateKey, err error) ***REMOVED***
	var priv pkcs1PrivateKey
	rest, err := asn1.Unmarshal(der, &priv)
	if len(rest) > 0 ***REMOVED***
		err = asn1.SyntaxError***REMOVED***Msg: "trailing data"***REMOVED***
		return
	***REMOVED***
	if err != nil ***REMOVED***
		return
	***REMOVED***

	if priv.Version > 1 ***REMOVED***
		return nil, errors.New("x509: unsupported private key version")
	***REMOVED***

	if priv.N.Sign() <= 0 || priv.D.Sign() <= 0 || priv.P.Sign() <= 0 || priv.Q.Sign() <= 0 ***REMOVED***
		return nil, errors.New("x509: private key contains zero or negative value")
	***REMOVED***

	key = new(rsa.PrivateKey)
	key.PublicKey = rsa.PublicKey***REMOVED***
		E: priv.E,
		N: priv.N,
	***REMOVED***

	key.D = priv.D
	key.Primes = make([]*big.Int, 2+len(priv.AdditionalPrimes))
	key.Primes[0] = priv.P
	key.Primes[1] = priv.Q
	for i, a := range priv.AdditionalPrimes ***REMOVED***
		if a.Prime.Sign() <= 0 ***REMOVED***
			return nil, errors.New("x509: private key contains zero or negative prime")
		***REMOVED***
		key.Primes[i+2] = a.Prime
		// We ignore the other two values because rsa will calculate
		// them as needed.
	***REMOVED***

	err = key.Validate()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	key.Precompute()

	return
***REMOVED***

// MarshalPKCS1PrivateKey converts a private key to ASN.1 DER encoded form.
func MarshalPKCS1PrivateKey(key *rsa.PrivateKey) []byte ***REMOVED***
	key.Precompute()

	version := 0
	if len(key.Primes) > 2 ***REMOVED***
		version = 1
	***REMOVED***

	priv := pkcs1PrivateKey***REMOVED***
		Version: version,
		N:       key.N,
		E:       key.PublicKey.E,
		D:       key.D,
		P:       key.Primes[0],
		Q:       key.Primes[1],
		Dp:      key.Precomputed.Dp,
		Dq:      key.Precomputed.Dq,
		Qinv:    key.Precomputed.Qinv,
	***REMOVED***

	priv.AdditionalPrimes = make([]pkcs1AdditionalRSAPrime, len(key.Precomputed.CRTValues))
	for i, values := range key.Precomputed.CRTValues ***REMOVED***
		priv.AdditionalPrimes[i].Prime = key.Primes[2+i]
		priv.AdditionalPrimes[i].Exp = values.Exp
		priv.AdditionalPrimes[i].Coeff = values.Coeff
	***REMOVED***

	b, _ := asn1.Marshal(priv)
	return b
***REMOVED***

// rsaPublicKey reflects the ASN.1 structure of a PKCS#1 public key.
type rsaPublicKey struct ***REMOVED***
	N *big.Int
	E int
***REMOVED***
