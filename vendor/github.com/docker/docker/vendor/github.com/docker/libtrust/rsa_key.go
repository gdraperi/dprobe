package libtrust

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
)

/*
 * RSA DSA PUBLIC KEY
 */

// rsaPublicKey implements a JWK Public Key using RSA digital signature algorithms.
type rsaPublicKey struct ***REMOVED***
	*rsa.PublicKey
	extended map[string]interface***REMOVED******REMOVED***
***REMOVED***

func fromRSAPublicKey(cryptoPublicKey *rsa.PublicKey) *rsaPublicKey ***REMOVED***
	return &rsaPublicKey***REMOVED***cryptoPublicKey, map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***
***REMOVED***

// KeyType returns the JWK key type for RSA keys, i.e., "RSA".
func (k *rsaPublicKey) KeyType() string ***REMOVED***
	return "RSA"
***REMOVED***

// KeyID returns a distinct identifier which is unique to this Public Key.
func (k *rsaPublicKey) KeyID() string ***REMOVED***
	return keyIDFromCryptoKey(k)
***REMOVED***

func (k *rsaPublicKey) String() string ***REMOVED***
	return fmt.Sprintf("RSA Public Key <%s>", k.KeyID())
***REMOVED***

// Verify verifyies the signature of the data in the io.Reader using this Public Key.
// The alg parameter should be the name of the JWA digital signature algorithm
// which was used to produce the signature and should be supported by this
// public key. Returns a nil error if the signature is valid.
func (k *rsaPublicKey) Verify(data io.Reader, alg string, signature []byte) error ***REMOVED***
	// Verify the signature of the given date, return non-nil error if valid.
	sigAlg, err := rsaSignatureAlgorithmByName(alg)
	if err != nil ***REMOVED***
		return fmt.Errorf("unable to verify Signature: %s", err)
	***REMOVED***

	hasher := sigAlg.HashID().New()
	_, err = io.Copy(hasher, data)
	if err != nil ***REMOVED***
		return fmt.Errorf("error reading data to sign: %s", err)
	***REMOVED***
	hash := hasher.Sum(nil)

	err = rsa.VerifyPKCS1v15(k.PublicKey, sigAlg.HashID(), hash, signature)
	if err != nil ***REMOVED***
		return fmt.Errorf("invalid %s signature: %s", sigAlg.HeaderParam(), err)
	***REMOVED***

	return nil
***REMOVED***

// CryptoPublicKey returns the internal object which can be used as a
// crypto.PublicKey for use with other standard library operations. The type
// is either *rsa.PublicKey or *ecdsa.PublicKey
func (k *rsaPublicKey) CryptoPublicKey() crypto.PublicKey ***REMOVED***
	return k.PublicKey
***REMOVED***

func (k *rsaPublicKey) toMap() map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	jwk := make(map[string]interface***REMOVED******REMOVED***)
	for k, v := range k.extended ***REMOVED***
		jwk[k] = v
	***REMOVED***
	jwk["kty"] = k.KeyType()
	jwk["kid"] = k.KeyID()
	jwk["n"] = joseBase64UrlEncode(k.N.Bytes())
	jwk["e"] = joseBase64UrlEncode(serializeRSAPublicExponentParam(k.E))

	return jwk
***REMOVED***

// MarshalJSON serializes this Public Key using the JWK JSON serialization format for
// RSA keys.
func (k *rsaPublicKey) MarshalJSON() (data []byte, err error) ***REMOVED***
	return json.Marshal(k.toMap())
***REMOVED***

// PEMBlock serializes this Public Key to DER-encoded PKIX format.
func (k *rsaPublicKey) PEMBlock() (*pem.Block, error) ***REMOVED***
	derBytes, err := x509.MarshalPKIXPublicKey(k.PublicKey)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("unable to serialize RSA PublicKey to DER-encoded PKIX format: %s", err)
	***REMOVED***
	k.extended["kid"] = k.KeyID() // For display purposes.
	return createPemBlock("PUBLIC KEY", derBytes, k.extended)
***REMOVED***

func (k *rsaPublicKey) AddExtendedField(field string, value interface***REMOVED******REMOVED***) ***REMOVED***
	k.extended[field] = value
***REMOVED***

func (k *rsaPublicKey) GetExtendedField(field string) interface***REMOVED******REMOVED*** ***REMOVED***
	v, ok := k.extended[field]
	if !ok ***REMOVED***
		return nil
	***REMOVED***
	return v
***REMOVED***

func rsaPublicKeyFromMap(jwk map[string]interface***REMOVED******REMOVED***) (*rsaPublicKey, error) ***REMOVED***
	// JWK key type (kty) has already been determined to be "RSA".
	// Need to extract 'n', 'e', and 'kid' and check for
	// consistency.

	// Get the modulus parameter N.
	nB64Url, err := stringFromMap(jwk, "n")
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("JWK RSA Public Key modulus: %s", err)
	***REMOVED***

	n, err := parseRSAModulusParam(nB64Url)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("JWK RSA Public Key modulus: %s", err)
	***REMOVED***

	// Get the public exponent E.
	eB64Url, err := stringFromMap(jwk, "e")
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("JWK RSA Public Key exponent: %s", err)
	***REMOVED***

	e, err := parseRSAPublicExponentParam(eB64Url)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("JWK RSA Public Key exponent: %s", err)
	***REMOVED***

	key := &rsaPublicKey***REMOVED***
		PublicKey: &rsa.PublicKey***REMOVED***N: n, E: e***REMOVED***,
	***REMOVED***

	// Key ID is optional, but if it exists, it should match the key.
	_, ok := jwk["kid"]
	if ok ***REMOVED***
		kid, err := stringFromMap(jwk, "kid")
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("JWK RSA Public Key ID: %s", err)
		***REMOVED***
		if kid != key.KeyID() ***REMOVED***
			return nil, fmt.Errorf("JWK RSA Public Key ID does not match: %s", kid)
		***REMOVED***
	***REMOVED***

	if _, ok := jwk["d"]; ok ***REMOVED***
		return nil, fmt.Errorf("JWK RSA Public Key cannot contain private exponent")
	***REMOVED***

	key.extended = jwk

	return key, nil
***REMOVED***

/*
 * RSA DSA PRIVATE KEY
 */

// rsaPrivateKey implements a JWK Private Key using RSA digital signature algorithms.
type rsaPrivateKey struct ***REMOVED***
	rsaPublicKey
	*rsa.PrivateKey
***REMOVED***

func fromRSAPrivateKey(cryptoPrivateKey *rsa.PrivateKey) *rsaPrivateKey ***REMOVED***
	return &rsaPrivateKey***REMOVED***
		*fromRSAPublicKey(&cryptoPrivateKey.PublicKey),
		cryptoPrivateKey,
	***REMOVED***
***REMOVED***

// PublicKey returns the Public Key data associated with this Private Key.
func (k *rsaPrivateKey) PublicKey() PublicKey ***REMOVED***
	return &k.rsaPublicKey
***REMOVED***

func (k *rsaPrivateKey) String() string ***REMOVED***
	return fmt.Sprintf("RSA Private Key <%s>", k.KeyID())
***REMOVED***

// Sign signs the data read from the io.Reader using a signature algorithm supported
// by the RSA private key. If the specified hashing algorithm is supported by
// this key, that hash function is used to generate the signature otherwise the
// the default hashing algorithm for this key is used. Returns the signature
// and the name of the JWK signature algorithm used, e.g., "RS256", "RS384",
// "RS512".
func (k *rsaPrivateKey) Sign(data io.Reader, hashID crypto.Hash) (signature []byte, alg string, err error) ***REMOVED***
	// Generate a signature of the data using the internal alg.
	sigAlg := rsaPKCS1v15SignatureAlgorithmForHashID(hashID)
	hasher := sigAlg.HashID().New()

	_, err = io.Copy(hasher, data)
	if err != nil ***REMOVED***
		return nil, "", fmt.Errorf("error reading data to sign: %s", err)
	***REMOVED***
	hash := hasher.Sum(nil)

	signature, err = rsa.SignPKCS1v15(rand.Reader, k.PrivateKey, sigAlg.HashID(), hash)
	if err != nil ***REMOVED***
		return nil, "", fmt.Errorf("error producing signature: %s", err)
	***REMOVED***

	alg = sigAlg.HeaderParam()

	return
***REMOVED***

// CryptoPrivateKey returns the internal object which can be used as a
// crypto.PublicKey for use with other standard library operations. The type
// is either *rsa.PublicKey or *ecdsa.PublicKey
func (k *rsaPrivateKey) CryptoPrivateKey() crypto.PrivateKey ***REMOVED***
	return k.PrivateKey
***REMOVED***

func (k *rsaPrivateKey) toMap() map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	k.Precompute() // Make sure the precomputed values are stored.
	jwk := k.rsaPublicKey.toMap()

	jwk["d"] = joseBase64UrlEncode(k.D.Bytes())
	jwk["p"] = joseBase64UrlEncode(k.Primes[0].Bytes())
	jwk["q"] = joseBase64UrlEncode(k.Primes[1].Bytes())
	jwk["dp"] = joseBase64UrlEncode(k.Precomputed.Dp.Bytes())
	jwk["dq"] = joseBase64UrlEncode(k.Precomputed.Dq.Bytes())
	jwk["qi"] = joseBase64UrlEncode(k.Precomputed.Qinv.Bytes())

	otherPrimes := k.Primes[2:]

	if len(otherPrimes) > 0 ***REMOVED***
		otherPrimesInfo := make([]interface***REMOVED******REMOVED***, len(otherPrimes))
		for i, r := range otherPrimes ***REMOVED***
			otherPrimeInfo := make(map[string]string, 3)
			otherPrimeInfo["r"] = joseBase64UrlEncode(r.Bytes())
			crtVal := k.Precomputed.CRTValues[i]
			otherPrimeInfo["d"] = joseBase64UrlEncode(crtVal.Exp.Bytes())
			otherPrimeInfo["t"] = joseBase64UrlEncode(crtVal.Coeff.Bytes())
			otherPrimesInfo[i] = otherPrimeInfo
		***REMOVED***
		jwk["oth"] = otherPrimesInfo
	***REMOVED***

	return jwk
***REMOVED***

// MarshalJSON serializes this Private Key using the JWK JSON serialization format for
// RSA keys.
func (k *rsaPrivateKey) MarshalJSON() (data []byte, err error) ***REMOVED***
	return json.Marshal(k.toMap())
***REMOVED***

// PEMBlock serializes this Private Key to DER-encoded PKIX format.
func (k *rsaPrivateKey) PEMBlock() (*pem.Block, error) ***REMOVED***
	derBytes := x509.MarshalPKCS1PrivateKey(k.PrivateKey)
	k.extended["keyID"] = k.KeyID() // For display purposes.
	return createPemBlock("RSA PRIVATE KEY", derBytes, k.extended)
***REMOVED***

func rsaPrivateKeyFromMap(jwk map[string]interface***REMOVED******REMOVED***) (*rsaPrivateKey, error) ***REMOVED***
	// The JWA spec for RSA Private Keys (draft rfc section 5.3.2) states that
	// only the private key exponent 'd' is REQUIRED, the others are just for
	// signature/decryption optimizations and SHOULD be included when the JWK
	// is produced. We MAY choose to accept a JWK which only includes 'd', but
	// we're going to go ahead and not choose to accept it without the extra
	// fields. Only the 'oth' field will be optional (for multi-prime keys).
	privateExponent, err := parseRSAPrivateKeyParamFromMap(jwk, "d")
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("JWK RSA Private Key exponent: %s", err)
	***REMOVED***
	firstPrimeFactor, err := parseRSAPrivateKeyParamFromMap(jwk, "p")
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("JWK RSA Private Key prime factor: %s", err)
	***REMOVED***
	secondPrimeFactor, err := parseRSAPrivateKeyParamFromMap(jwk, "q")
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("JWK RSA Private Key prime factor: %s", err)
	***REMOVED***
	firstFactorCRT, err := parseRSAPrivateKeyParamFromMap(jwk, "dp")
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("JWK RSA Private Key CRT exponent: %s", err)
	***REMOVED***
	secondFactorCRT, err := parseRSAPrivateKeyParamFromMap(jwk, "dq")
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("JWK RSA Private Key CRT exponent: %s", err)
	***REMOVED***
	crtCoeff, err := parseRSAPrivateKeyParamFromMap(jwk, "qi")
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("JWK RSA Private Key CRT coefficient: %s", err)
	***REMOVED***

	var oth interface***REMOVED******REMOVED***
	if _, ok := jwk["oth"]; ok ***REMOVED***
		oth = jwk["oth"]
		delete(jwk, "oth")
	***REMOVED***

	// JWK key type (kty) has already been determined to be "RSA".
	// Need to extract the public key information, then extract the private
	// key values.
	publicKey, err := rsaPublicKeyFromMap(jwk)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	privateKey := &rsa.PrivateKey***REMOVED***
		PublicKey: *publicKey.PublicKey,
		D:         privateExponent,
		Primes:    []*big.Int***REMOVED***firstPrimeFactor, secondPrimeFactor***REMOVED***,
		Precomputed: rsa.PrecomputedValues***REMOVED***
			Dp:   firstFactorCRT,
			Dq:   secondFactorCRT,
			Qinv: crtCoeff,
		***REMOVED***,
	***REMOVED***

	if oth != nil ***REMOVED***
		// Should be an array of more JSON objects.
		otherPrimesInfo, ok := oth.([]interface***REMOVED******REMOVED***)
		if !ok ***REMOVED***
			return nil, errors.New("JWK RSA Private Key: Invalid other primes info: must be an array")
		***REMOVED***
		numOtherPrimeFactors := len(otherPrimesInfo)
		if numOtherPrimeFactors == 0 ***REMOVED***
			return nil, errors.New("JWK RSA Privake Key: Invalid other primes info: must be absent or non-empty")
		***REMOVED***
		otherPrimeFactors := make([]*big.Int, numOtherPrimeFactors)
		productOfPrimes := new(big.Int).Mul(firstPrimeFactor, secondPrimeFactor)
		crtValues := make([]rsa.CRTValue, numOtherPrimeFactors)

		for i, val := range otherPrimesInfo ***REMOVED***
			otherPrimeinfo, ok := val.(map[string]interface***REMOVED******REMOVED***)
			if !ok ***REMOVED***
				return nil, errors.New("JWK RSA Private Key: Invalid other prime info: must be a JSON object")
			***REMOVED***

			otherPrimeFactor, err := parseRSAPrivateKeyParamFromMap(otherPrimeinfo, "r")
			if err != nil ***REMOVED***
				return nil, fmt.Errorf("JWK RSA Private Key prime factor: %s", err)
			***REMOVED***
			otherFactorCRT, err := parseRSAPrivateKeyParamFromMap(otherPrimeinfo, "d")
			if err != nil ***REMOVED***
				return nil, fmt.Errorf("JWK RSA Private Key CRT exponent: %s", err)
			***REMOVED***
			otherCrtCoeff, err := parseRSAPrivateKeyParamFromMap(otherPrimeinfo, "t")
			if err != nil ***REMOVED***
				return nil, fmt.Errorf("JWK RSA Private Key CRT coefficient: %s", err)
			***REMOVED***

			crtValue := crtValues[i]
			crtValue.Exp = otherFactorCRT
			crtValue.Coeff = otherCrtCoeff
			crtValue.R = productOfPrimes
			otherPrimeFactors[i] = otherPrimeFactor
			productOfPrimes = new(big.Int).Mul(productOfPrimes, otherPrimeFactor)
		***REMOVED***

		privateKey.Primes = append(privateKey.Primes, otherPrimeFactors...)
		privateKey.Precomputed.CRTValues = crtValues
	***REMOVED***

	key := &rsaPrivateKey***REMOVED***
		rsaPublicKey: *publicKey,
		PrivateKey:   privateKey,
	***REMOVED***

	return key, nil
***REMOVED***

/*
 *	Key Generation Functions.
 */

func generateRSAPrivateKey(bits int) (k *rsaPrivateKey, err error) ***REMOVED***
	k = new(rsaPrivateKey)
	k.PrivateKey, err = rsa.GenerateKey(rand.Reader, bits)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	k.rsaPublicKey.PublicKey = &k.PrivateKey.PublicKey
	k.extended = make(map[string]interface***REMOVED******REMOVED***)

	return
***REMOVED***

// GenerateRSA2048PrivateKey generates a key pair using 2048-bit RSA.
func GenerateRSA2048PrivateKey() (PrivateKey, error) ***REMOVED***
	k, err := generateRSAPrivateKey(2048)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("error generating RSA 2048-bit key: %s", err)
	***REMOVED***

	return k, nil
***REMOVED***

// GenerateRSA3072PrivateKey generates a key pair using 3072-bit RSA.
func GenerateRSA3072PrivateKey() (PrivateKey, error) ***REMOVED***
	k, err := generateRSAPrivateKey(3072)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("error generating RSA 3072-bit key: %s", err)
	***REMOVED***

	return k, nil
***REMOVED***

// GenerateRSA4096PrivateKey generates a key pair using 4096-bit RSA.
func GenerateRSA4096PrivateKey() (PrivateKey, error) ***REMOVED***
	k, err := generateRSAPrivateKey(4096)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("error generating RSA 4096-bit key: %s", err)
	***REMOVED***

	return k, nil
***REMOVED***
