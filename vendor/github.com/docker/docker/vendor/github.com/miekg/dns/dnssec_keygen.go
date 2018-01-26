package dns

import (
	"crypto"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"math/big"
)

// Generate generates a DNSKEY of the given bit size.
// The public part is put inside the DNSKEY record.
// The Algorithm in the key must be set as this will define
// what kind of DNSKEY will be generated.
// The ECDSA algorithms imply a fixed keysize, in that case
// bits should be set to the size of the algorithm.
func (k *DNSKEY) Generate(bits int) (crypto.PrivateKey, error) ***REMOVED***
	switch k.Algorithm ***REMOVED***
	case DSA, DSANSEC3SHA1:
		if bits != 1024 ***REMOVED***
			return nil, ErrKeySize
		***REMOVED***
	case RSAMD5, RSASHA1, RSASHA256, RSASHA1NSEC3SHA1:
		if bits < 512 || bits > 4096 ***REMOVED***
			return nil, ErrKeySize
		***REMOVED***
	case RSASHA512:
		if bits < 1024 || bits > 4096 ***REMOVED***
			return nil, ErrKeySize
		***REMOVED***
	case ECDSAP256SHA256:
		if bits != 256 ***REMOVED***
			return nil, ErrKeySize
		***REMOVED***
	case ECDSAP384SHA384:
		if bits != 384 ***REMOVED***
			return nil, ErrKeySize
		***REMOVED***
	***REMOVED***

	switch k.Algorithm ***REMOVED***
	case DSA, DSANSEC3SHA1:
		params := new(dsa.Parameters)
		if err := dsa.GenerateParameters(params, rand.Reader, dsa.L1024N160); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		priv := new(dsa.PrivateKey)
		priv.PublicKey.Parameters = *params
		err := dsa.GenerateKey(priv, rand.Reader)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		k.setPublicKeyDSA(params.Q, params.P, params.G, priv.PublicKey.Y)
		return priv, nil
	case RSAMD5, RSASHA1, RSASHA256, RSASHA512, RSASHA1NSEC3SHA1:
		priv, err := rsa.GenerateKey(rand.Reader, bits)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		k.setPublicKeyRSA(priv.PublicKey.E, priv.PublicKey.N)
		return priv, nil
	case ECDSAP256SHA256, ECDSAP384SHA384:
		var c elliptic.Curve
		switch k.Algorithm ***REMOVED***
		case ECDSAP256SHA256:
			c = elliptic.P256()
		case ECDSAP384SHA384:
			c = elliptic.P384()
		***REMOVED***
		priv, err := ecdsa.GenerateKey(c, rand.Reader)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		k.setPublicKeyECDSA(priv.PublicKey.X, priv.PublicKey.Y)
		return priv, nil
	default:
		return nil, ErrAlg
	***REMOVED***
***REMOVED***

// Set the public key (the value E and N)
func (k *DNSKEY) setPublicKeyRSA(_E int, _N *big.Int) bool ***REMOVED***
	if _E == 0 || _N == nil ***REMOVED***
		return false
	***REMOVED***
	buf := exponentToBuf(_E)
	buf = append(buf, _N.Bytes()...)
	k.PublicKey = toBase64(buf)
	return true
***REMOVED***

// Set the public key for Elliptic Curves
func (k *DNSKEY) setPublicKeyECDSA(_X, _Y *big.Int) bool ***REMOVED***
	if _X == nil || _Y == nil ***REMOVED***
		return false
	***REMOVED***
	var intlen int
	switch k.Algorithm ***REMOVED***
	case ECDSAP256SHA256:
		intlen = 32
	case ECDSAP384SHA384:
		intlen = 48
	***REMOVED***
	k.PublicKey = toBase64(curveToBuf(_X, _Y, intlen))
	return true
***REMOVED***

// Set the public key for DSA
func (k *DNSKEY) setPublicKeyDSA(_Q, _P, _G, _Y *big.Int) bool ***REMOVED***
	if _Q == nil || _P == nil || _G == nil || _Y == nil ***REMOVED***
		return false
	***REMOVED***
	buf := dsaToBuf(_Q, _P, _G, _Y)
	k.PublicKey = toBase64(buf)
	return true
***REMOVED***

// Set the public key (the values E and N) for RSA
// RFC 3110: Section 2. RSA Public KEY Resource Records
func exponentToBuf(_E int) []byte ***REMOVED***
	var buf []byte
	i := big.NewInt(int64(_E))
	if len(i.Bytes()) < 256 ***REMOVED***
		buf = make([]byte, 1)
		buf[0] = uint8(len(i.Bytes()))
	***REMOVED*** else ***REMOVED***
		buf = make([]byte, 3)
		buf[0] = 0
		buf[1] = uint8(len(i.Bytes()) >> 8)
		buf[2] = uint8(len(i.Bytes()))
	***REMOVED***
	buf = append(buf, i.Bytes()...)
	return buf
***REMOVED***

// Set the public key for X and Y for Curve. The two
// values are just concatenated.
func curveToBuf(_X, _Y *big.Int, intlen int) []byte ***REMOVED***
	buf := intToBytes(_X, intlen)
	buf = append(buf, intToBytes(_Y, intlen)...)
	return buf
***REMOVED***

// Set the public key for X and Y for Curve. The two
// values are just concatenated.
func dsaToBuf(_Q, _P, _G, _Y *big.Int) []byte ***REMOVED***
	t := divRoundUp(divRoundUp(_G.BitLen(), 8)-64, 8)
	buf := []byte***REMOVED***byte(t)***REMOVED***
	buf = append(buf, intToBytes(_Q, 20)...)
	buf = append(buf, intToBytes(_P, 64+t*8)...)
	buf = append(buf, intToBytes(_G, 64+t*8)...)
	buf = append(buf, intToBytes(_Y, 64+t*8)...)
	return buf
***REMOVED***
