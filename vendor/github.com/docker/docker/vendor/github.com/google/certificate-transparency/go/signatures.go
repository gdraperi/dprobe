package ct

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/big"
)

var allowVerificationWithNonCompliantKeys = flag.Bool("allow_verification_with_non_compliant_keys", false,
	"Allow a SignatureVerifier to use keys which are technically non-compliant with RFC6962.")

// PublicKeyFromPEM parses a PEM formatted block and returns the public key contained within and any remaining unread bytes, or an error.
func PublicKeyFromPEM(b []byte) (crypto.PublicKey, SHA256Hash, []byte, error) ***REMOVED***
	p, rest := pem.Decode(b)
	if p == nil ***REMOVED***
		return nil, [sha256.Size]byte***REMOVED******REMOVED***, rest, fmt.Errorf("no PEM block found in %s", string(b))
	***REMOVED***
	k, err := x509.ParsePKIXPublicKey(p.Bytes)
	return k, sha256.Sum256(p.Bytes), rest, err
***REMOVED***

// SignatureVerifier can verify signatures on SCTs and STHs
type SignatureVerifier struct ***REMOVED***
	pubKey crypto.PublicKey
***REMOVED***

// NewSignatureVerifier creates a new SignatureVerifier using the passed in PublicKey.
func NewSignatureVerifier(pk crypto.PublicKey) (*SignatureVerifier, error) ***REMOVED***
	switch pkType := pk.(type) ***REMOVED***
	case *rsa.PublicKey:
		if pkType.N.BitLen() < 2048 ***REMOVED***
			e := fmt.Errorf("public key is RSA with < 2048 bits (size:%d)", pkType.N.BitLen())
			if !(*allowVerificationWithNonCompliantKeys) ***REMOVED***
				return nil, e
			***REMOVED***
			log.Printf("WARNING: %v", e)
		***REMOVED***
	case *ecdsa.PublicKey:
		params := *(pkType.Params())
		if params != *elliptic.P256().Params() ***REMOVED***
			e := fmt.Errorf("public is ECDSA, but not on the P256 curve")
			if !(*allowVerificationWithNonCompliantKeys) ***REMOVED***
				return nil, e
			***REMOVED***
			log.Printf("WARNING: %v", e)

		***REMOVED***
	default:
		return nil, fmt.Errorf("Unsupported public key type %v", pkType)
	***REMOVED***

	return &SignatureVerifier***REMOVED***
		pubKey: pk,
	***REMOVED***, nil
***REMOVED***

// verifySignature verifies that the passed in signature over data was created by our PublicKey.
// Currently, only SHA256 is supported as a HashAlgorithm, and only ECDSA and RSA signatures are supported.
func (s SignatureVerifier) verifySignature(data []byte, sig DigitallySigned) error ***REMOVED***
	if sig.HashAlgorithm != SHA256 ***REMOVED***
		return fmt.Errorf("unsupported HashAlgorithm in signature: %v", sig.HashAlgorithm)
	***REMOVED***

	hasherType := crypto.SHA256
	hasher := hasherType.New()
	if _, err := hasher.Write(data); err != nil ***REMOVED***
		return fmt.Errorf("failed to write to hasher: %v", err)
	***REMOVED***
	hash := hasher.Sum([]byte***REMOVED******REMOVED***)

	switch sig.SignatureAlgorithm ***REMOVED***
	case RSA:
		rsaKey, ok := s.pubKey.(*rsa.PublicKey)
		if !ok ***REMOVED***
			return fmt.Errorf("cannot verify RSA signature with %T key", s.pubKey)
		***REMOVED***
		if err := rsa.VerifyPKCS1v15(rsaKey, hasherType, hash, sig.Signature); err != nil ***REMOVED***
			return fmt.Errorf("failed to verify rsa signature: %v", err)
		***REMOVED***
	case ECDSA:
		ecdsaKey, ok := s.pubKey.(*ecdsa.PublicKey)
		if !ok ***REMOVED***
			return fmt.Errorf("cannot verify ECDSA signature with %T key", s.pubKey)
		***REMOVED***
		var ecdsaSig struct ***REMOVED***
			R, S *big.Int
		***REMOVED***
		rest, err := asn1.Unmarshal(sig.Signature, &ecdsaSig)
		if err != nil ***REMOVED***
			return fmt.Errorf("failed to unmarshal ECDSA signature: %v", err)
		***REMOVED***
		if len(rest) != 0 ***REMOVED***
			log.Printf("Garbage following signature %v", rest)
		***REMOVED***

		if !ecdsa.Verify(ecdsaKey, hash, ecdsaSig.R, ecdsaSig.S) ***REMOVED***
			return errors.New("failed to verify ecdsa signature")
		***REMOVED***
	default:
		return fmt.Errorf("unsupported signature type %v", sig.SignatureAlgorithm)
	***REMOVED***
	return nil
***REMOVED***

// VerifySCTSignature verifies that the SCT's signature is valid for the given LogEntry
func (s SignatureVerifier) VerifySCTSignature(sct SignedCertificateTimestamp, entry LogEntry) error ***REMOVED***
	sctData, err := SerializeSCTSignatureInput(sct, entry)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return s.verifySignature(sctData, sct.Signature)
***REMOVED***

// VerifySTHSignature verifies that the STH's signature is valid.
func (s SignatureVerifier) VerifySTHSignature(sth SignedTreeHead) error ***REMOVED***
	sthData, err := SerializeSTHSignatureInput(sth)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return s.verifySignature(sthData, sth.TreeHeadSignature)
***REMOVED***
