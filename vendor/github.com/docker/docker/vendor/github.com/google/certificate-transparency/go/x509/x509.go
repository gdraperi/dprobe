// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package x509 parses X.509-encoded keys and certificates.
//
// START CT CHANGES
// This is a fork of the go library crypto/x509 package, it's more relaxed
// about certificates that it'll accept, and exports the TBSCertificate
// structure.
// END CT CHANGES
package x509

import (
	"bytes"
	"crypto"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/sha1"
	// START CT CHANGES
	"github.com/google/certificate-transparency/go/asn1"
	"github.com/google/certificate-transparency/go/x509/pkix"
	// END CT CHANGES
	"encoding/pem"
	"errors"
	// START CT CHANGES
	"fmt"
	// END CT CHANGES
	"io"
	"math/big"
	"net"
	"time"
)

// pkixPublicKey reflects a PKIX public key structure. See SubjectPublicKeyInfo
// in RFC 3280.
type pkixPublicKey struct ***REMOVED***
	Algo      pkix.AlgorithmIdentifier
	BitString asn1.BitString
***REMOVED***

// ParsePKIXPublicKey parses a DER encoded public key. These values are
// typically found in PEM blocks with "BEGIN PUBLIC KEY".
func ParsePKIXPublicKey(derBytes []byte) (pub interface***REMOVED******REMOVED***, err error) ***REMOVED***
	var pki publicKeyInfo
	if _, err = asn1.Unmarshal(derBytes, &pki); err != nil ***REMOVED***
		return
	***REMOVED***
	algo := getPublicKeyAlgorithmFromOID(pki.Algorithm.Algorithm)
	if algo == UnknownPublicKeyAlgorithm ***REMOVED***
		return nil, errors.New("x509: unknown public key algorithm")
	***REMOVED***
	return parsePublicKey(algo, &pki)
***REMOVED***

func marshalPublicKey(pub interface***REMOVED******REMOVED***) (publicKeyBytes []byte, publicKeyAlgorithm pkix.AlgorithmIdentifier, err error) ***REMOVED***
	switch pub := pub.(type) ***REMOVED***
	case *rsa.PublicKey:
		publicKeyBytes, err = asn1.Marshal(rsaPublicKey***REMOVED***
			N: pub.N,
			E: pub.E,
		***REMOVED***)
		publicKeyAlgorithm.Algorithm = oidPublicKeyRSA
		// This is a NULL parameters value which is technically
		// superfluous, but most other code includes it and, by
		// doing this, we match their public key hashes.
		publicKeyAlgorithm.Parameters = asn1.RawValue***REMOVED***
			Tag: 5,
		***REMOVED***
	case *ecdsa.PublicKey:
		publicKeyBytes = elliptic.Marshal(pub.Curve, pub.X, pub.Y)
		oid, ok := oidFromNamedCurve(pub.Curve)
		if !ok ***REMOVED***
			return nil, pkix.AlgorithmIdentifier***REMOVED******REMOVED***, errors.New("x509: unsupported elliptic curve")
		***REMOVED***
		publicKeyAlgorithm.Algorithm = oidPublicKeyECDSA
		var paramBytes []byte
		paramBytes, err = asn1.Marshal(oid)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		publicKeyAlgorithm.Parameters.FullBytes = paramBytes
	default:
		return nil, pkix.AlgorithmIdentifier***REMOVED******REMOVED***, errors.New("x509: only RSA and ECDSA public keys supported")
	***REMOVED***

	return publicKeyBytes, publicKeyAlgorithm, nil
***REMOVED***

// MarshalPKIXPublicKey serialises a public key to DER-encoded PKIX format.
func MarshalPKIXPublicKey(pub interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	var publicKeyBytes []byte
	var publicKeyAlgorithm pkix.AlgorithmIdentifier
	var err error

	if publicKeyBytes, publicKeyAlgorithm, err = marshalPublicKey(pub); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	pkix := pkixPublicKey***REMOVED***
		Algo: publicKeyAlgorithm,
		BitString: asn1.BitString***REMOVED***
			Bytes:     publicKeyBytes,
			BitLength: 8 * len(publicKeyBytes),
		***REMOVED***,
	***REMOVED***

	ret, _ := asn1.Marshal(pkix)
	return ret, nil
***REMOVED***

// These structures reflect the ASN.1 structure of X.509 certificates.:

type certificate struct ***REMOVED***
	Raw                asn1.RawContent
	TBSCertificate     tbsCertificate
	SignatureAlgorithm pkix.AlgorithmIdentifier
	SignatureValue     asn1.BitString
***REMOVED***

type tbsCertificate struct ***REMOVED***
	Raw                asn1.RawContent
	Version            int `asn1:"optional,explicit,default:1,tag:0"`
	SerialNumber       *big.Int
	SignatureAlgorithm pkix.AlgorithmIdentifier
	Issuer             asn1.RawValue
	Validity           validity
	Subject            asn1.RawValue
	PublicKey          publicKeyInfo
	UniqueId           asn1.BitString   `asn1:"optional,tag:1"`
	SubjectUniqueId    asn1.BitString   `asn1:"optional,tag:2"`
	Extensions         []pkix.Extension `asn1:"optional,explicit,tag:3"`
***REMOVED***

type dsaAlgorithmParameters struct ***REMOVED***
	P, Q, G *big.Int
***REMOVED***

type dsaSignature struct ***REMOVED***
	R, S *big.Int
***REMOVED***

type ecdsaSignature dsaSignature

type validity struct ***REMOVED***
	NotBefore, NotAfter time.Time
***REMOVED***

type publicKeyInfo struct ***REMOVED***
	Raw       asn1.RawContent
	Algorithm pkix.AlgorithmIdentifier
	PublicKey asn1.BitString
***REMOVED***

// RFC 5280,  4.2.1.1
type authKeyId struct ***REMOVED***
	Id []byte `asn1:"optional,tag:0"`
***REMOVED***

type SignatureAlgorithm int

const (
	UnknownSignatureAlgorithm SignatureAlgorithm = iota
	MD2WithRSA
	MD5WithRSA
	SHA1WithRSA
	SHA256WithRSA
	SHA384WithRSA
	SHA512WithRSA
	DSAWithSHA1
	DSAWithSHA256
	ECDSAWithSHA1
	ECDSAWithSHA256
	ECDSAWithSHA384
	ECDSAWithSHA512
)

type PublicKeyAlgorithm int

const (
	UnknownPublicKeyAlgorithm PublicKeyAlgorithm = iota
	RSA
	DSA
	ECDSA
)

// OIDs for signature algorithms
//
// pkcs-1 OBJECT IDENTIFIER ::= ***REMOVED***
//    iso(1) member-body(2) us(840) rsadsi(113549) pkcs(1) 1 ***REMOVED***
//
//
// RFC 3279 2.2.1 RSA Signature Algorithms
//
// md2WithRSAEncryption OBJECT IDENTIFIER ::= ***REMOVED*** pkcs-1 2 ***REMOVED***
//
// md5WithRSAEncryption OBJECT IDENTIFIER ::= ***REMOVED*** pkcs-1 4 ***REMOVED***
//
// sha-1WithRSAEncryption OBJECT IDENTIFIER ::= ***REMOVED*** pkcs-1 5 ***REMOVED***
//
// dsaWithSha1 OBJECT IDENTIFIER ::= ***REMOVED***
//    iso(1) member-body(2) us(840) x9-57(10040) x9cm(4) 3 ***REMOVED***
//
// RFC 3279 2.2.3 ECDSA Signature Algorithm
//
// ecdsa-with-SHA1 OBJECT IDENTIFIER ::= ***REMOVED***
// 	  iso(1) member-body(2) us(840) ansi-x962(10045)
//    signatures(4) ecdsa-with-SHA1(1)***REMOVED***
//
//
// RFC 4055 5 PKCS #1 Version 1.5
//
// sha256WithRSAEncryption OBJECT IDENTIFIER ::= ***REMOVED*** pkcs-1 11 ***REMOVED***
//
// sha384WithRSAEncryption OBJECT IDENTIFIER ::= ***REMOVED*** pkcs-1 12 ***REMOVED***
//
// sha512WithRSAEncryption OBJECT IDENTIFIER ::= ***REMOVED*** pkcs-1 13 ***REMOVED***
//
//
// RFC 5758 3.1 DSA Signature Algorithms
//
// dsaWithSha256 OBJECT IDENTIFIER ::= ***REMOVED***
//    joint-iso-ccitt(2) country(16) us(840) organization(1) gov(101)
//    csor(3) algorithms(4) id-dsa-with-sha2(3) 2***REMOVED***
//
// RFC 5758 3.2 ECDSA Signature Algorithm
//
// ecdsa-with-SHA256 OBJECT IDENTIFIER ::= ***REMOVED*** iso(1) member-body(2)
//    us(840) ansi-X9-62(10045) signatures(4) ecdsa-with-SHA2(3) 2 ***REMOVED***
//
// ecdsa-with-SHA384 OBJECT IDENTIFIER ::= ***REMOVED*** iso(1) member-body(2)
//    us(840) ansi-X9-62(10045) signatures(4) ecdsa-with-SHA2(3) 3 ***REMOVED***
//
// ecdsa-with-SHA512 OBJECT IDENTIFIER ::= ***REMOVED*** iso(1) member-body(2)
//    us(840) ansi-X9-62(10045) signatures(4) ecdsa-with-SHA2(3) 4 ***REMOVED***

var (
	oidSignatureMD2WithRSA      = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 113549, 1, 1, 2***REMOVED***
	oidSignatureMD5WithRSA      = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 113549, 1, 1, 4***REMOVED***
	oidSignatureSHA1WithRSA     = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 113549, 1, 1, 5***REMOVED***
	oidSignatureSHA256WithRSA   = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 113549, 1, 1, 11***REMOVED***
	oidSignatureSHA384WithRSA   = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 113549, 1, 1, 12***REMOVED***
	oidSignatureSHA512WithRSA   = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 113549, 1, 1, 13***REMOVED***
	oidSignatureDSAWithSHA1     = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 10040, 4, 3***REMOVED***
	oidSignatureDSAWithSHA256   = asn1.ObjectIdentifier***REMOVED***2, 16, 840, 1, 101, 4, 3, 2***REMOVED***
	oidSignatureECDSAWithSHA1   = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 10045, 4, 1***REMOVED***
	oidSignatureECDSAWithSHA256 = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 10045, 4, 3, 2***REMOVED***
	oidSignatureECDSAWithSHA384 = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 10045, 4, 3, 3***REMOVED***
	oidSignatureECDSAWithSHA512 = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 10045, 4, 3, 4***REMOVED***
)

func getSignatureAlgorithmFromOID(oid asn1.ObjectIdentifier) SignatureAlgorithm ***REMOVED***
	switch ***REMOVED***
	case oid.Equal(oidSignatureMD2WithRSA):
		return MD2WithRSA
	case oid.Equal(oidSignatureMD5WithRSA):
		return MD5WithRSA
	case oid.Equal(oidSignatureSHA1WithRSA):
		return SHA1WithRSA
	case oid.Equal(oidSignatureSHA256WithRSA):
		return SHA256WithRSA
	case oid.Equal(oidSignatureSHA384WithRSA):
		return SHA384WithRSA
	case oid.Equal(oidSignatureSHA512WithRSA):
		return SHA512WithRSA
	case oid.Equal(oidSignatureDSAWithSHA1):
		return DSAWithSHA1
	case oid.Equal(oidSignatureDSAWithSHA256):
		return DSAWithSHA256
	case oid.Equal(oidSignatureECDSAWithSHA1):
		return ECDSAWithSHA1
	case oid.Equal(oidSignatureECDSAWithSHA256):
		return ECDSAWithSHA256
	case oid.Equal(oidSignatureECDSAWithSHA384):
		return ECDSAWithSHA384
	case oid.Equal(oidSignatureECDSAWithSHA512):
		return ECDSAWithSHA512
	***REMOVED***
	return UnknownSignatureAlgorithm
***REMOVED***

// RFC 3279, 2.3 Public Key Algorithms
//
// pkcs-1 OBJECT IDENTIFIER ::== ***REMOVED*** iso(1) member-body(2) us(840)
//    rsadsi(113549) pkcs(1) 1 ***REMOVED***
//
// rsaEncryption OBJECT IDENTIFIER ::== ***REMOVED*** pkcs1-1 1 ***REMOVED***
//
// id-dsa OBJECT IDENTIFIER ::== ***REMOVED*** iso(1) member-body(2) us(840)
//    x9-57(10040) x9cm(4) 1 ***REMOVED***
//
// RFC 5480, 2.1.1 Unrestricted Algorithm Identifier and Parameters
//
// id-ecPublicKey OBJECT IDENTIFIER ::= ***REMOVED***
//       iso(1) member-body(2) us(840) ansi-X9-62(10045) keyType(2) 1 ***REMOVED***
var (
	oidPublicKeyRSA   = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 113549, 1, 1, 1***REMOVED***
	oidPublicKeyDSA   = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 10040, 4, 1***REMOVED***
	oidPublicKeyECDSA = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 10045, 2, 1***REMOVED***
)

func getPublicKeyAlgorithmFromOID(oid asn1.ObjectIdentifier) PublicKeyAlgorithm ***REMOVED***
	switch ***REMOVED***
	case oid.Equal(oidPublicKeyRSA):
		return RSA
	case oid.Equal(oidPublicKeyDSA):
		return DSA
	case oid.Equal(oidPublicKeyECDSA):
		return ECDSA
	***REMOVED***
	return UnknownPublicKeyAlgorithm
***REMOVED***

// RFC 5480, 2.1.1.1. Named Curve
//
// secp224r1 OBJECT IDENTIFIER ::= ***REMOVED***
//   iso(1) identified-organization(3) certicom(132) curve(0) 33 ***REMOVED***
//
// secp256r1 OBJECT IDENTIFIER ::= ***REMOVED***
//   iso(1) member-body(2) us(840) ansi-X9-62(10045) curves(3)
//   prime(1) 7 ***REMOVED***
//
// secp384r1 OBJECT IDENTIFIER ::= ***REMOVED***
//   iso(1) identified-organization(3) certicom(132) curve(0) 34 ***REMOVED***
//
// secp521r1 OBJECT IDENTIFIER ::= ***REMOVED***
//   iso(1) identified-organization(3) certicom(132) curve(0) 35 ***REMOVED***
//
// NB: secp256r1 is equivalent to prime256v1
var (
	oidNamedCurveP224 = asn1.ObjectIdentifier***REMOVED***1, 3, 132, 0, 33***REMOVED***
	oidNamedCurveP256 = asn1.ObjectIdentifier***REMOVED***1, 2, 840, 10045, 3, 1, 7***REMOVED***
	oidNamedCurveP384 = asn1.ObjectIdentifier***REMOVED***1, 3, 132, 0, 34***REMOVED***
	oidNamedCurveP521 = asn1.ObjectIdentifier***REMOVED***1, 3, 132, 0, 35***REMOVED***
)

func namedCurveFromOID(oid asn1.ObjectIdentifier) elliptic.Curve ***REMOVED***
	switch ***REMOVED***
	case oid.Equal(oidNamedCurveP224):
		return elliptic.P224()
	case oid.Equal(oidNamedCurveP256):
		return elliptic.P256()
	case oid.Equal(oidNamedCurveP384):
		return elliptic.P384()
	case oid.Equal(oidNamedCurveP521):
		return elliptic.P521()
	***REMOVED***
	return nil
***REMOVED***

func oidFromNamedCurve(curve elliptic.Curve) (asn1.ObjectIdentifier, bool) ***REMOVED***
	switch curve ***REMOVED***
	case elliptic.P224():
		return oidNamedCurveP224, true
	case elliptic.P256():
		return oidNamedCurveP256, true
	case elliptic.P384():
		return oidNamedCurveP384, true
	case elliptic.P521():
		return oidNamedCurveP521, true
	***REMOVED***

	return nil, false
***REMOVED***

// KeyUsage represents the set of actions that are valid for a given key. It's
// a bitmap of the KeyUsage* constants.
type KeyUsage int

const (
	KeyUsageDigitalSignature KeyUsage = 1 << iota
	KeyUsageContentCommitment
	KeyUsageKeyEncipherment
	KeyUsageDataEncipherment
	KeyUsageKeyAgreement
	KeyUsageCertSign
	KeyUsageCRLSign
	KeyUsageEncipherOnly
	KeyUsageDecipherOnly
)

// RFC 5280, 4.2.1.12  Extended Key Usage
//
// anyExtendedKeyUsage OBJECT IDENTIFIER ::= ***REMOVED*** id-ce-extKeyUsage 0 ***REMOVED***
//
// id-kp OBJECT IDENTIFIER ::= ***REMOVED*** id-pkix 3 ***REMOVED***
//
// id-kp-serverAuth             OBJECT IDENTIFIER ::= ***REMOVED*** id-kp 1 ***REMOVED***
// id-kp-clientAuth             OBJECT IDENTIFIER ::= ***REMOVED*** id-kp 2 ***REMOVED***
// id-kp-codeSigning            OBJECT IDENTIFIER ::= ***REMOVED*** id-kp 3 ***REMOVED***
// id-kp-emailProtection        OBJECT IDENTIFIER ::= ***REMOVED*** id-kp 4 ***REMOVED***
// id-kp-timeStamping           OBJECT IDENTIFIER ::= ***REMOVED*** id-kp 8 ***REMOVED***
// id-kp-OCSPSigning            OBJECT IDENTIFIER ::= ***REMOVED*** id-kp 9 ***REMOVED***
var (
	oidExtKeyUsageAny                        = asn1.ObjectIdentifier***REMOVED***2, 5, 29, 37, 0***REMOVED***
	oidExtKeyUsageServerAuth                 = asn1.ObjectIdentifier***REMOVED***1, 3, 6, 1, 5, 5, 7, 3, 1***REMOVED***
	oidExtKeyUsageClientAuth                 = asn1.ObjectIdentifier***REMOVED***1, 3, 6, 1, 5, 5, 7, 3, 2***REMOVED***
	oidExtKeyUsageCodeSigning                = asn1.ObjectIdentifier***REMOVED***1, 3, 6, 1, 5, 5, 7, 3, 3***REMOVED***
	oidExtKeyUsageEmailProtection            = asn1.ObjectIdentifier***REMOVED***1, 3, 6, 1, 5, 5, 7, 3, 4***REMOVED***
	oidExtKeyUsageIPSECEndSystem             = asn1.ObjectIdentifier***REMOVED***1, 3, 6, 1, 5, 5, 7, 3, 5***REMOVED***
	oidExtKeyUsageIPSECTunnel                = asn1.ObjectIdentifier***REMOVED***1, 3, 6, 1, 5, 5, 7, 3, 6***REMOVED***
	oidExtKeyUsageIPSECUser                  = asn1.ObjectIdentifier***REMOVED***1, 3, 6, 1, 5, 5, 7, 3, 7***REMOVED***
	oidExtKeyUsageTimeStamping               = asn1.ObjectIdentifier***REMOVED***1, 3, 6, 1, 5, 5, 7, 3, 8***REMOVED***
	oidExtKeyUsageOCSPSigning                = asn1.ObjectIdentifier***REMOVED***1, 3, 6, 1, 5, 5, 7, 3, 9***REMOVED***
	oidExtKeyUsageMicrosoftServerGatedCrypto = asn1.ObjectIdentifier***REMOVED***1, 3, 6, 1, 4, 1, 311, 10, 3, 3***REMOVED***
	oidExtKeyUsageNetscapeServerGatedCrypto  = asn1.ObjectIdentifier***REMOVED***2, 16, 840, 1, 113730, 4, 1***REMOVED***
)

// ExtKeyUsage represents an extended set of actions that are valid for a given key.
// Each of the ExtKeyUsage* constants define a unique action.
type ExtKeyUsage int

const (
	ExtKeyUsageAny ExtKeyUsage = iota
	ExtKeyUsageServerAuth
	ExtKeyUsageClientAuth
	ExtKeyUsageCodeSigning
	ExtKeyUsageEmailProtection
	ExtKeyUsageIPSECEndSystem
	ExtKeyUsageIPSECTunnel
	ExtKeyUsageIPSECUser
	ExtKeyUsageTimeStamping
	ExtKeyUsageOCSPSigning
	ExtKeyUsageMicrosoftServerGatedCrypto
	ExtKeyUsageNetscapeServerGatedCrypto
)

// extKeyUsageOIDs contains the mapping between an ExtKeyUsage and its OID.
var extKeyUsageOIDs = []struct ***REMOVED***
	extKeyUsage ExtKeyUsage
	oid         asn1.ObjectIdentifier
***REMOVED******REMOVED***
	***REMOVED***ExtKeyUsageAny, oidExtKeyUsageAny***REMOVED***,
	***REMOVED***ExtKeyUsageServerAuth, oidExtKeyUsageServerAuth***REMOVED***,
	***REMOVED***ExtKeyUsageClientAuth, oidExtKeyUsageClientAuth***REMOVED***,
	***REMOVED***ExtKeyUsageCodeSigning, oidExtKeyUsageCodeSigning***REMOVED***,
	***REMOVED***ExtKeyUsageEmailProtection, oidExtKeyUsageEmailProtection***REMOVED***,
	***REMOVED***ExtKeyUsageIPSECEndSystem, oidExtKeyUsageIPSECEndSystem***REMOVED***,
	***REMOVED***ExtKeyUsageIPSECTunnel, oidExtKeyUsageIPSECTunnel***REMOVED***,
	***REMOVED***ExtKeyUsageIPSECUser, oidExtKeyUsageIPSECUser***REMOVED***,
	***REMOVED***ExtKeyUsageTimeStamping, oidExtKeyUsageTimeStamping***REMOVED***,
	***REMOVED***ExtKeyUsageOCSPSigning, oidExtKeyUsageOCSPSigning***REMOVED***,
	***REMOVED***ExtKeyUsageMicrosoftServerGatedCrypto, oidExtKeyUsageMicrosoftServerGatedCrypto***REMOVED***,
	***REMOVED***ExtKeyUsageNetscapeServerGatedCrypto, oidExtKeyUsageNetscapeServerGatedCrypto***REMOVED***,
***REMOVED***

func extKeyUsageFromOID(oid asn1.ObjectIdentifier) (eku ExtKeyUsage, ok bool) ***REMOVED***
	for _, pair := range extKeyUsageOIDs ***REMOVED***
		if oid.Equal(pair.oid) ***REMOVED***
			return pair.extKeyUsage, true
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func oidFromExtKeyUsage(eku ExtKeyUsage) (oid asn1.ObjectIdentifier, ok bool) ***REMOVED***
	for _, pair := range extKeyUsageOIDs ***REMOVED***
		if eku == pair.extKeyUsage ***REMOVED***
			return pair.oid, true
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// A Certificate represents an X.509 certificate.
type Certificate struct ***REMOVED***
	Raw                     []byte // Complete ASN.1 DER content (certificate, signature algorithm and signature).
	RawTBSCertificate       []byte // Certificate part of raw ASN.1 DER content.
	RawSubjectPublicKeyInfo []byte // DER encoded SubjectPublicKeyInfo.
	RawSubject              []byte // DER encoded Subject
	RawIssuer               []byte // DER encoded Issuer

	Signature          []byte
	SignatureAlgorithm SignatureAlgorithm

	PublicKeyAlgorithm PublicKeyAlgorithm
	PublicKey          interface***REMOVED******REMOVED***

	Version             int
	SerialNumber        *big.Int
	Issuer              pkix.Name
	Subject             pkix.Name
	NotBefore, NotAfter time.Time // Validity bounds.
	KeyUsage            KeyUsage

	// Extensions contains raw X.509 extensions. When parsing certificates,
	// this can be used to extract non-critical extensions that are not
	// parsed by this package. When marshaling certificates, the Extensions
	// field is ignored, see ExtraExtensions.
	Extensions []pkix.Extension

	// ExtraExtensions contains extensions to be copied, raw, into any
	// marshaled certificates. Values override any extensions that would
	// otherwise be produced based on the other fields. The ExtraExtensions
	// field is not populated when parsing certificates, see Extensions.
	ExtraExtensions []pkix.Extension

	ExtKeyUsage        []ExtKeyUsage           // Sequence of extended key usages.
	UnknownExtKeyUsage []asn1.ObjectIdentifier // Encountered extended key usages unknown to this package.

	BasicConstraintsValid bool // if true then the next two fields are valid.
	IsCA                  bool
	MaxPathLen            int

	SubjectKeyId   []byte
	AuthorityKeyId []byte

	// RFC 5280, 4.2.2.1 (Authority Information Access)
	OCSPServer            []string
	IssuingCertificateURL []string

	// Subject Alternate Name values
	DNSNames       []string
	EmailAddresses []string
	IPAddresses    []net.IP

	// Name constraints
	PermittedDNSDomainsCritical bool // if true then the name constraints are marked critical.
	PermittedDNSDomains         []string

	// CRL Distribution Points
	CRLDistributionPoints []string

	PolicyIdentifiers []asn1.ObjectIdentifier
***REMOVED***

// ErrUnsupportedAlgorithm results from attempting to perform an operation that
// involves algorithms that are not currently implemented.
var ErrUnsupportedAlgorithm = errors.New("x509: cannot verify signature: algorithm unimplemented")

// ConstraintViolationError results when a requested usage is not permitted by
// a certificate. For example: checking a signature when the public key isn't a
// certificate signing key.
type ConstraintViolationError struct***REMOVED******REMOVED***

func (ConstraintViolationError) Error() string ***REMOVED***
	return "x509: invalid signature: parent certificate cannot sign this kind of certificate"
***REMOVED***

func (c *Certificate) Equal(other *Certificate) bool ***REMOVED***
	return bytes.Equal(c.Raw, other.Raw)
***REMOVED***

// Entrust have a broken root certificate (CN=Entrust.net Certification
// Authority (2048)) which isn't marked as a CA certificate and is thus invalid
// according to PKIX.
// We recognise this certificate by its SubjectPublicKeyInfo and exempt it
// from the Basic Constraints requirement.
// See http://www.entrust.net/knowledge-base/technote.cfm?tn=7869
//
// TODO(agl): remove this hack once their reissued root is sufficiently
// widespread.
var entrustBrokenSPKI = []byte***REMOVED***
	0x30, 0x82, 0x01, 0x22, 0x30, 0x0d, 0x06, 0x09,
	0x2a, 0x86, 0x48, 0x86, 0xf7, 0x0d, 0x01, 0x01,
	0x01, 0x05, 0x00, 0x03, 0x82, 0x01, 0x0f, 0x00,
	0x30, 0x82, 0x01, 0x0a, 0x02, 0x82, 0x01, 0x01,
	0x00, 0x97, 0xa3, 0x2d, 0x3c, 0x9e, 0xde, 0x05,
	0xda, 0x13, 0xc2, 0x11, 0x8d, 0x9d, 0x8e, 0xe3,
	0x7f, 0xc7, 0x4b, 0x7e, 0x5a, 0x9f, 0xb3, 0xff,
	0x62, 0xab, 0x73, 0xc8, 0x28, 0x6b, 0xba, 0x10,
	0x64, 0x82, 0x87, 0x13, 0xcd, 0x57, 0x18, 0xff,
	0x28, 0xce, 0xc0, 0xe6, 0x0e, 0x06, 0x91, 0x50,
	0x29, 0x83, 0xd1, 0xf2, 0xc3, 0x2a, 0xdb, 0xd8,
	0xdb, 0x4e, 0x04, 0xcc, 0x00, 0xeb, 0x8b, 0xb6,
	0x96, 0xdc, 0xbc, 0xaa, 0xfa, 0x52, 0x77, 0x04,
	0xc1, 0xdb, 0x19, 0xe4, 0xae, 0x9c, 0xfd, 0x3c,
	0x8b, 0x03, 0xef, 0x4d, 0xbc, 0x1a, 0x03, 0x65,
	0xf9, 0xc1, 0xb1, 0x3f, 0x72, 0x86, 0xf2, 0x38,
	0xaa, 0x19, 0xae, 0x10, 0x88, 0x78, 0x28, 0xda,
	0x75, 0xc3, 0x3d, 0x02, 0x82, 0x02, 0x9c, 0xb9,
	0xc1, 0x65, 0x77, 0x76, 0x24, 0x4c, 0x98, 0xf7,
	0x6d, 0x31, 0x38, 0xfb, 0xdb, 0xfe, 0xdb, 0x37,
	0x02, 0x76, 0xa1, 0x18, 0x97, 0xa6, 0xcc, 0xde,
	0x20, 0x09, 0x49, 0x36, 0x24, 0x69, 0x42, 0xf6,
	0xe4, 0x37, 0x62, 0xf1, 0x59, 0x6d, 0xa9, 0x3c,
	0xed, 0x34, 0x9c, 0xa3, 0x8e, 0xdb, 0xdc, 0x3a,
	0xd7, 0xf7, 0x0a, 0x6f, 0xef, 0x2e, 0xd8, 0xd5,
	0x93, 0x5a, 0x7a, 0xed, 0x08, 0x49, 0x68, 0xe2,
	0x41, 0xe3, 0x5a, 0x90, 0xc1, 0x86, 0x55, 0xfc,
	0x51, 0x43, 0x9d, 0xe0, 0xb2, 0xc4, 0x67, 0xb4,
	0xcb, 0x32, 0x31, 0x25, 0xf0, 0x54, 0x9f, 0x4b,
	0xd1, 0x6f, 0xdb, 0xd4, 0xdd, 0xfc, 0xaf, 0x5e,
	0x6c, 0x78, 0x90, 0x95, 0xde, 0xca, 0x3a, 0x48,
	0xb9, 0x79, 0x3c, 0x9b, 0x19, 0xd6, 0x75, 0x05,
	0xa0, 0xf9, 0x88, 0xd7, 0xc1, 0xe8, 0xa5, 0x09,
	0xe4, 0x1a, 0x15, 0xdc, 0x87, 0x23, 0xaa, 0xb2,
	0x75, 0x8c, 0x63, 0x25, 0x87, 0xd8, 0xf8, 0x3d,
	0xa6, 0xc2, 0xcc, 0x66, 0xff, 0xa5, 0x66, 0x68,
	0x55, 0x02, 0x03, 0x01, 0x00, 0x01,
***REMOVED***

// CheckSignatureFrom verifies that the signature on c is a valid signature
// from parent.
func (c *Certificate) CheckSignatureFrom(parent *Certificate) (err error) ***REMOVED***
	// RFC 5280, 4.2.1.9:
	// "If the basic constraints extension is not present in a version 3
	// certificate, or the extension is present but the cA boolean is not
	// asserted, then the certified public key MUST NOT be used to verify
	// certificate signatures."
	// (except for Entrust, see comment above entrustBrokenSPKI)
	if (parent.Version == 3 && !parent.BasicConstraintsValid ||
		parent.BasicConstraintsValid && !parent.IsCA) &&
		!bytes.Equal(c.RawSubjectPublicKeyInfo, entrustBrokenSPKI) ***REMOVED***
		return ConstraintViolationError***REMOVED******REMOVED***
	***REMOVED***

	if parent.KeyUsage != 0 && parent.KeyUsage&KeyUsageCertSign == 0 ***REMOVED***
		return ConstraintViolationError***REMOVED******REMOVED***
	***REMOVED***

	if parent.PublicKeyAlgorithm == UnknownPublicKeyAlgorithm ***REMOVED***
		return ErrUnsupportedAlgorithm
	***REMOVED***

	// TODO(agl): don't ignore the path length constraint.

	return parent.CheckSignature(c.SignatureAlgorithm, c.RawTBSCertificate, c.Signature)
***REMOVED***

// CheckSignature verifies that signature is a valid signature over signed from
// c's public key.
func (c *Certificate) CheckSignature(algo SignatureAlgorithm, signed, signature []byte) (err error) ***REMOVED***
	var hashType crypto.Hash

	switch algo ***REMOVED***
	case SHA1WithRSA, DSAWithSHA1, ECDSAWithSHA1:
		hashType = crypto.SHA1
	case SHA256WithRSA, DSAWithSHA256, ECDSAWithSHA256:
		hashType = crypto.SHA256
	case SHA384WithRSA, ECDSAWithSHA384:
		hashType = crypto.SHA384
	case SHA512WithRSA, ECDSAWithSHA512:
		hashType = crypto.SHA512
	default:
		return ErrUnsupportedAlgorithm
	***REMOVED***

	if !hashType.Available() ***REMOVED***
		return ErrUnsupportedAlgorithm
	***REMOVED***
	h := hashType.New()

	h.Write(signed)
	digest := h.Sum(nil)

	switch pub := c.PublicKey.(type) ***REMOVED***
	case *rsa.PublicKey:
		return rsa.VerifyPKCS1v15(pub, hashType, digest, signature)
	case *dsa.PublicKey:
		dsaSig := new(dsaSignature)
		if _, err := asn1.Unmarshal(signature, dsaSig); err != nil ***REMOVED***
			return err
		***REMOVED***
		if dsaSig.R.Sign() <= 0 || dsaSig.S.Sign() <= 0 ***REMOVED***
			return errors.New("x509: DSA signature contained zero or negative values")
		***REMOVED***
		if !dsa.Verify(pub, digest, dsaSig.R, dsaSig.S) ***REMOVED***
			return errors.New("x509: DSA verification failure")
		***REMOVED***
		return
	case *ecdsa.PublicKey:
		ecdsaSig := new(ecdsaSignature)
		if _, err := asn1.Unmarshal(signature, ecdsaSig); err != nil ***REMOVED***
			return err
		***REMOVED***
		if ecdsaSig.R.Sign() <= 0 || ecdsaSig.S.Sign() <= 0 ***REMOVED***
			return errors.New("x509: ECDSA signature contained zero or negative values")
		***REMOVED***
		if !ecdsa.Verify(pub, digest, ecdsaSig.R, ecdsaSig.S) ***REMOVED***
			return errors.New("x509: ECDSA verification failure")
		***REMOVED***
		return
	***REMOVED***
	return ErrUnsupportedAlgorithm
***REMOVED***

// CheckCRLSignature checks that the signature in crl is from c.
func (c *Certificate) CheckCRLSignature(crl *pkix.CertificateList) (err error) ***REMOVED***
	algo := getSignatureAlgorithmFromOID(crl.SignatureAlgorithm.Algorithm)
	return c.CheckSignature(algo, crl.TBSCertList.Raw, crl.SignatureValue.RightAlign())
***REMOVED***

// START CT CHANGES
type UnhandledCriticalExtension struct ***REMOVED***
	ID asn1.ObjectIdentifier
***REMOVED***

func (h UnhandledCriticalExtension) Error() string ***REMOVED***
	return fmt.Sprintf("x509: unhandled critical extension (%v)", h.ID)
***REMOVED***

// END CT CHANGES

type basicConstraints struct ***REMOVED***
	IsCA       bool `asn1:"optional"`
	MaxPathLen int  `asn1:"optional,default:-1"`
***REMOVED***

// RFC 5280 4.2.1.4
type policyInformation struct ***REMOVED***
	Policy asn1.ObjectIdentifier
	// policyQualifiers omitted
***REMOVED***

// RFC 5280, 4.2.1.10
type nameConstraints struct ***REMOVED***
	Permitted []generalSubtree `asn1:"optional,tag:0"`
	Excluded  []generalSubtree `asn1:"optional,tag:1"`
***REMOVED***

type generalSubtree struct ***REMOVED***
	Name string `asn1:"tag:2,optional,ia5"`
***REMOVED***

// RFC 5280, 4.2.2.1
type authorityInfoAccess struct ***REMOVED***
	Method   asn1.ObjectIdentifier
	Location asn1.RawValue
***REMOVED***

// RFC 5280, 4.2.1.14
type distributionPoint struct ***REMOVED***
	DistributionPoint distributionPointName `asn1:"optional,tag:0"`
	Reason            asn1.BitString        `asn1:"optional,tag:1"`
	CRLIssuer         asn1.RawValue         `asn1:"optional,tag:2"`
***REMOVED***

type distributionPointName struct ***REMOVED***
	FullName     asn1.RawValue    `asn1:"optional,tag:0"`
	RelativeName pkix.RDNSequence `asn1:"optional,tag:1"`
***REMOVED***

func parsePublicKey(algo PublicKeyAlgorithm, keyData *publicKeyInfo) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	asn1Data := keyData.PublicKey.RightAlign()
	switch algo ***REMOVED***
	case RSA:
		p := new(rsaPublicKey)
		_, err := asn1.Unmarshal(asn1Data, p)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		if p.N.Sign() <= 0 ***REMOVED***
			return nil, errors.New("x509: RSA modulus is not a positive number")
		***REMOVED***
		if p.E <= 0 ***REMOVED***
			return nil, errors.New("x509: RSA public exponent is not a positive number")
		***REMOVED***

		pub := &rsa.PublicKey***REMOVED***
			E: p.E,
			N: p.N,
		***REMOVED***
		return pub, nil
	case DSA:
		var p *big.Int
		_, err := asn1.Unmarshal(asn1Data, &p)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		paramsData := keyData.Algorithm.Parameters.FullBytes
		params := new(dsaAlgorithmParameters)
		_, err = asn1.Unmarshal(paramsData, params)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if p.Sign() <= 0 || params.P.Sign() <= 0 || params.Q.Sign() <= 0 || params.G.Sign() <= 0 ***REMOVED***
			return nil, errors.New("x509: zero or negative DSA parameter")
		***REMOVED***
		pub := &dsa.PublicKey***REMOVED***
			Parameters: dsa.Parameters***REMOVED***
				P: params.P,
				Q: params.Q,
				G: params.G,
			***REMOVED***,
			Y: p,
		***REMOVED***
		return pub, nil
	case ECDSA:
		paramsData := keyData.Algorithm.Parameters.FullBytes
		namedCurveOID := new(asn1.ObjectIdentifier)
		_, err := asn1.Unmarshal(paramsData, namedCurveOID)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		namedCurve := namedCurveFromOID(*namedCurveOID)
		if namedCurve == nil ***REMOVED***
			return nil, errors.New("x509: unsupported elliptic curve")
		***REMOVED***
		x, y := elliptic.Unmarshal(namedCurve, asn1Data)
		if x == nil ***REMOVED***
			return nil, errors.New("x509: failed to unmarshal elliptic curve point")
		***REMOVED***
		pub := &ecdsa.PublicKey***REMOVED***
			Curve: namedCurve,
			X:     x,
			Y:     y,
		***REMOVED***
		return pub, nil
	default:
		return nil, nil
	***REMOVED***
***REMOVED***

// START CT CHANGES

// NonFatalErrors is an error type which can hold a number of other errors.
// It's used to collect a range of non-fatal errors which occur while parsing
// a certificate, that way we can still match on certs which technically are
// invalid.
type NonFatalErrors struct ***REMOVED***
	Errors []error
***REMOVED***

// Adds an error to the list of errors contained by NonFatalErrors.
func (e *NonFatalErrors) AddError(err error) ***REMOVED***
	e.Errors = append(e.Errors, err)
***REMOVED***

// Returns a string consisting of the values of Error() from all of the errors
// contained in |e|
func (e NonFatalErrors) Error() string ***REMOVED***
	r := "NonFatalErrors: "
	for _, err := range e.Errors ***REMOVED***
		r += err.Error() + "; "
	***REMOVED***
	return r
***REMOVED***

// Returns true if |e| contains at least one error
func (e *NonFatalErrors) HasError() bool ***REMOVED***
	return len(e.Errors) > 0
***REMOVED***

// END CT CHANGES

func parseCertificate(in *certificate) (*Certificate, error) ***REMOVED***
	// START CT CHANGES
	var nfe NonFatalErrors
	// END CT CHANGES

	out := new(Certificate)
	out.Raw = in.Raw
	out.RawTBSCertificate = in.TBSCertificate.Raw
	out.RawSubjectPublicKeyInfo = in.TBSCertificate.PublicKey.Raw
	out.RawSubject = in.TBSCertificate.Subject.FullBytes
	out.RawIssuer = in.TBSCertificate.Issuer.FullBytes

	out.Signature = in.SignatureValue.RightAlign()
	out.SignatureAlgorithm =
		getSignatureAlgorithmFromOID(in.TBSCertificate.SignatureAlgorithm.Algorithm)

	out.PublicKeyAlgorithm =
		getPublicKeyAlgorithmFromOID(in.TBSCertificate.PublicKey.Algorithm.Algorithm)
	var err error
	out.PublicKey, err = parsePublicKey(out.PublicKeyAlgorithm, &in.TBSCertificate.PublicKey)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if in.TBSCertificate.SerialNumber.Sign() < 0 ***REMOVED***
		// START CT CHANGES
		nfe.AddError(errors.New("x509: negative serial number"))
		// END CT CHANGES
	***REMOVED***

	out.Version = in.TBSCertificate.Version + 1
	out.SerialNumber = in.TBSCertificate.SerialNumber

	var issuer, subject pkix.RDNSequence
	if _, err := asn1.Unmarshal(in.TBSCertificate.Subject.FullBytes, &subject); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if _, err := asn1.Unmarshal(in.TBSCertificate.Issuer.FullBytes, &issuer); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	out.Issuer.FillFromRDNSequence(&issuer)
	out.Subject.FillFromRDNSequence(&subject)

	out.NotBefore = in.TBSCertificate.Validity.NotBefore
	out.NotAfter = in.TBSCertificate.Validity.NotAfter

	for _, e := range in.TBSCertificate.Extensions ***REMOVED***
		out.Extensions = append(out.Extensions, e)

		if len(e.Id) == 4 && e.Id[0] == 2 && e.Id[1] == 5 && e.Id[2] == 29 ***REMOVED***
			switch e.Id[3] ***REMOVED***
			case 15:
				// RFC 5280, 4.2.1.3
				var usageBits asn1.BitString
				_, err := asn1.Unmarshal(e.Value, &usageBits)

				if err == nil ***REMOVED***
					var usage int
					for i := 0; i < 9; i++ ***REMOVED***
						if usageBits.At(i) != 0 ***REMOVED***
							usage |= 1 << uint(i)
						***REMOVED***
					***REMOVED***
					out.KeyUsage = KeyUsage(usage)
					continue
				***REMOVED***
			case 19:
				// RFC 5280, 4.2.1.9
				var constraints basicConstraints
				_, err := asn1.Unmarshal(e.Value, &constraints)

				if err == nil ***REMOVED***
					out.BasicConstraintsValid = true
					out.IsCA = constraints.IsCA
					out.MaxPathLen = constraints.MaxPathLen
					continue
				***REMOVED***
			case 17:
				// RFC 5280, 4.2.1.6

				// SubjectAltName ::= GeneralNames
				//
				// GeneralNames ::= SEQUENCE SIZE (1..MAX) OF GeneralName
				//
				// GeneralName ::= CHOICE ***REMOVED***
				//      otherName                       [0]     OtherName,
				//      rfc822Name                      [1]     IA5String,
				//      dNSName                         [2]     IA5String,
				//      x400Address                     [3]     ORAddress,
				//      directoryName                   [4]     Name,
				//      ediPartyName                    [5]     EDIPartyName,
				//      uniformResourceIdentifier       [6]     IA5String,
				//      iPAddress                       [7]     OCTET STRING,
				//      registeredID                    [8]     OBJECT IDENTIFIER ***REMOVED***
				var seq asn1.RawValue
				_, err := asn1.Unmarshal(e.Value, &seq)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				if !seq.IsCompound || seq.Tag != 16 || seq.Class != 0 ***REMOVED***
					return nil, asn1.StructuralError***REMOVED***Msg: "bad SAN sequence"***REMOVED***
				***REMOVED***

				parsedName := false

				rest := seq.Bytes
				for len(rest) > 0 ***REMOVED***
					var v asn1.RawValue
					rest, err = asn1.Unmarshal(rest, &v)
					if err != nil ***REMOVED***
						return nil, err
					***REMOVED***
					switch v.Tag ***REMOVED***
					case 1:
						out.EmailAddresses = append(out.EmailAddresses, string(v.Bytes))
						parsedName = true
					case 2:
						out.DNSNames = append(out.DNSNames, string(v.Bytes))
						parsedName = true
					case 7:
						switch len(v.Bytes) ***REMOVED***
						case net.IPv4len, net.IPv6len:
							out.IPAddresses = append(out.IPAddresses, v.Bytes)
						default:
							// START CT CHANGES
							nfe.AddError(fmt.Errorf("x509: certificate contained IP address of length %d : %v", len(v.Bytes), v.Bytes))
							// END CT CHANGES
						***REMOVED***
					***REMOVED***
				***REMOVED***

				if parsedName ***REMOVED***
					continue
				***REMOVED***
				// If we didn't parse any of the names then we
				// fall through to the critical check below.

			case 30:
				// RFC 5280, 4.2.1.10

				// NameConstraints ::= SEQUENCE ***REMOVED***
				//      permittedSubtrees       [0]     GeneralSubtrees OPTIONAL,
				//      excludedSubtrees        [1]     GeneralSubtrees OPTIONAL ***REMOVED***
				//
				// GeneralSubtrees ::= SEQUENCE SIZE (1..MAX) OF GeneralSubtree
				//
				// GeneralSubtree ::= SEQUENCE ***REMOVED***
				//      base                    GeneralName,
				//      minimum         [0]     BaseDistance DEFAULT 0,
				//      maximum         [1]     BaseDistance OPTIONAL ***REMOVED***
				//
				// BaseDistance ::= INTEGER (0..MAX)

				var constraints nameConstraints
				_, err := asn1.Unmarshal(e.Value, &constraints)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***

				if len(constraints.Excluded) > 0 && e.Critical ***REMOVED***
					// START CT CHANGES
					nfe.AddError(UnhandledCriticalExtension***REMOVED***e.Id***REMOVED***)
					// END CT CHANGES
				***REMOVED***

				for _, subtree := range constraints.Permitted ***REMOVED***
					if len(subtree.Name) == 0 ***REMOVED***
						if e.Critical ***REMOVED***
							// START CT CHANGES
							nfe.AddError(UnhandledCriticalExtension***REMOVED***e.Id***REMOVED***)
							// END CT CHANGES
						***REMOVED***
						continue
					***REMOVED***
					out.PermittedDNSDomains = append(out.PermittedDNSDomains, subtree.Name)
				***REMOVED***
				continue

			case 31:
				// RFC 5280, 4.2.1.14

				// CRLDistributionPoints ::= SEQUENCE SIZE (1..MAX) OF DistributionPoint
				//
				// DistributionPoint ::= SEQUENCE ***REMOVED***
				//     distributionPoint       [0]     DistributionPointName OPTIONAL,
				//     reasons                 [1]     ReasonFlags OPTIONAL,
				//     cRLIssuer               [2]     GeneralNames OPTIONAL ***REMOVED***
				//
				// DistributionPointName ::= CHOICE ***REMOVED***
				//     fullName                [0]     GeneralNames,
				//     nameRelativeToCRLIssuer [1]     RelativeDistinguishedName ***REMOVED***

				var cdp []distributionPoint
				_, err := asn1.Unmarshal(e.Value, &cdp)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***

				for _, dp := range cdp ***REMOVED***
					var n asn1.RawValue
					_, err = asn1.Unmarshal(dp.DistributionPoint.FullName.Bytes, &n)
					if err != nil ***REMOVED***
						return nil, err
					***REMOVED***

					if n.Tag == 6 ***REMOVED***
						out.CRLDistributionPoints = append(out.CRLDistributionPoints, string(n.Bytes))
					***REMOVED***
				***REMOVED***
				continue

			case 35:
				// RFC 5280, 4.2.1.1
				var a authKeyId
				_, err = asn1.Unmarshal(e.Value, &a)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				out.AuthorityKeyId = a.Id
				continue

			case 37:
				// RFC 5280, 4.2.1.12.  Extended Key Usage

				// id-ce-extKeyUsage OBJECT IDENTIFIER ::= ***REMOVED*** id-ce 37 ***REMOVED***
				//
				// ExtKeyUsageSyntax ::= SEQUENCE SIZE (1..MAX) OF KeyPurposeId
				//
				// KeyPurposeId ::= OBJECT IDENTIFIER

				var keyUsage []asn1.ObjectIdentifier
				_, err = asn1.Unmarshal(e.Value, &keyUsage)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***

				for _, u := range keyUsage ***REMOVED***
					if extKeyUsage, ok := extKeyUsageFromOID(u); ok ***REMOVED***
						out.ExtKeyUsage = append(out.ExtKeyUsage, extKeyUsage)
					***REMOVED*** else ***REMOVED***
						out.UnknownExtKeyUsage = append(out.UnknownExtKeyUsage, u)
					***REMOVED***
				***REMOVED***

				continue

			case 14:
				// RFC 5280, 4.2.1.2
				var keyid []byte
				_, err = asn1.Unmarshal(e.Value, &keyid)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				out.SubjectKeyId = keyid
				continue

			case 32:
				// RFC 5280 4.2.1.4: Certificate Policies
				var policies []policyInformation
				if _, err = asn1.Unmarshal(e.Value, &policies); err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				out.PolicyIdentifiers = make([]asn1.ObjectIdentifier, len(policies))
				for i, policy := range policies ***REMOVED***
					out.PolicyIdentifiers[i] = policy.Policy
				***REMOVED***
			***REMOVED***
		***REMOVED*** else if e.Id.Equal(oidExtensionAuthorityInfoAccess) ***REMOVED***
			// RFC 5280 4.2.2.1: Authority Information Access
			var aia []authorityInfoAccess
			if _, err = asn1.Unmarshal(e.Value, &aia); err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			for _, v := range aia ***REMOVED***
				// GeneralName: uniformResourceIdentifier [6] IA5String
				if v.Location.Tag != 6 ***REMOVED***
					continue
				***REMOVED***
				if v.Method.Equal(oidAuthorityInfoAccessOcsp) ***REMOVED***
					out.OCSPServer = append(out.OCSPServer, string(v.Location.Bytes))
				***REMOVED*** else if v.Method.Equal(oidAuthorityInfoAccessIssuers) ***REMOVED***
					out.IssuingCertificateURL = append(out.IssuingCertificateURL, string(v.Location.Bytes))
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if e.Critical ***REMOVED***
			// START CT CHANGES
			nfe.AddError(UnhandledCriticalExtension***REMOVED***e.Id***REMOVED***)
			// END CT CHANGES
		***REMOVED***
	***REMOVED***
	// START CT CHANGES
	if nfe.HasError() ***REMOVED***
		return out, nfe
	***REMOVED***
	// END CT CHANGES
	return out, nil
***REMOVED***

// START CT CHANGES

// ParseTBSCertificate parses a single TBSCertificate from the given ASN.1 DER data.
// The parsed data is returned in a Certificate struct for ease of access.
func ParseTBSCertificate(asn1Data []byte) (*Certificate, error) ***REMOVED***
	var tbsCert tbsCertificate
	rest, err := asn1.Unmarshal(asn1Data, &tbsCert)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(rest) > 0 ***REMOVED***
		return nil, asn1.SyntaxError***REMOVED***Msg: "trailing data"***REMOVED***
	***REMOVED***
	return parseCertificate(&certificate***REMOVED***
		Raw:            tbsCert.Raw,
		TBSCertificate: tbsCert***REMOVED***)
***REMOVED***

// END CT CHANGES

// ParseCertificate parses a single certificate from the given ASN.1 DER data.
func ParseCertificate(asn1Data []byte) (*Certificate, error) ***REMOVED***
	var cert certificate
	rest, err := asn1.Unmarshal(asn1Data, &cert)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(rest) > 0 ***REMOVED***
		return nil, asn1.SyntaxError***REMOVED***Msg: "trailing data"***REMOVED***
	***REMOVED***

	return parseCertificate(&cert)
***REMOVED***

// ParseCertificates parses one or more certificates from the given ASN.1 DER
// data. The certificates must be concatenated with no intermediate padding.
func ParseCertificates(asn1Data []byte) ([]*Certificate, error) ***REMOVED***
	var v []*certificate

	for len(asn1Data) > 0 ***REMOVED***
		cert := new(certificate)
		var err error
		asn1Data, err = asn1.Unmarshal(asn1Data, cert)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		v = append(v, cert)
	***REMOVED***

	ret := make([]*Certificate, len(v))
	for i, ci := range v ***REMOVED***
		cert, err := parseCertificate(ci)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		ret[i] = cert
	***REMOVED***

	return ret, nil
***REMOVED***

func reverseBitsInAByte(in byte) byte ***REMOVED***
	b1 := in>>4 | in<<4
	b2 := b1>>2&0x33 | b1<<2&0xcc
	b3 := b2>>1&0x55 | b2<<1&0xaa
	return b3
***REMOVED***

var (
	oidExtensionSubjectKeyId          = []int***REMOVED***2, 5, 29, 14***REMOVED***
	oidExtensionKeyUsage              = []int***REMOVED***2, 5, 29, 15***REMOVED***
	oidExtensionExtendedKeyUsage      = []int***REMOVED***2, 5, 29, 37***REMOVED***
	oidExtensionAuthorityKeyId        = []int***REMOVED***2, 5, 29, 35***REMOVED***
	oidExtensionBasicConstraints      = []int***REMOVED***2, 5, 29, 19***REMOVED***
	oidExtensionSubjectAltName        = []int***REMOVED***2, 5, 29, 17***REMOVED***
	oidExtensionCertificatePolicies   = []int***REMOVED***2, 5, 29, 32***REMOVED***
	oidExtensionNameConstraints       = []int***REMOVED***2, 5, 29, 30***REMOVED***
	oidExtensionCRLDistributionPoints = []int***REMOVED***2, 5, 29, 31***REMOVED***
	oidExtensionAuthorityInfoAccess   = []int***REMOVED***1, 3, 6, 1, 5, 5, 7, 1, 1***REMOVED***
)

var (
	oidAuthorityInfoAccessOcsp    = asn1.ObjectIdentifier***REMOVED***1, 3, 6, 1, 5, 5, 7, 48, 1***REMOVED***
	oidAuthorityInfoAccessIssuers = asn1.ObjectIdentifier***REMOVED***1, 3, 6, 1, 5, 5, 7, 48, 2***REMOVED***
)

// oidNotInExtensions returns whether an extension with the given oid exists in
// extensions.
func oidInExtensions(oid asn1.ObjectIdentifier, extensions []pkix.Extension) bool ***REMOVED***
	for _, e := range extensions ***REMOVED***
		if e.Id.Equal(oid) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func buildExtensions(template *Certificate) (ret []pkix.Extension, err error) ***REMOVED***
	ret = make([]pkix.Extension, 10 /* maximum number of elements. */)
	n := 0

	if template.KeyUsage != 0 &&
		!oidInExtensions(oidExtensionKeyUsage, template.ExtraExtensions) ***REMOVED***
		ret[n].Id = oidExtensionKeyUsage
		ret[n].Critical = true

		var a [2]byte
		a[0] = reverseBitsInAByte(byte(template.KeyUsage))
		a[1] = reverseBitsInAByte(byte(template.KeyUsage >> 8))

		l := 1
		if a[1] != 0 ***REMOVED***
			l = 2
		***REMOVED***

		ret[n].Value, err = asn1.Marshal(asn1.BitString***REMOVED***Bytes: a[0:l], BitLength: l * 8***REMOVED***)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		n++
	***REMOVED***

	if (len(template.ExtKeyUsage) > 0 || len(template.UnknownExtKeyUsage) > 0) &&
		!oidInExtensions(oidExtensionExtendedKeyUsage, template.ExtraExtensions) ***REMOVED***
		ret[n].Id = oidExtensionExtendedKeyUsage

		var oids []asn1.ObjectIdentifier
		for _, u := range template.ExtKeyUsage ***REMOVED***
			if oid, ok := oidFromExtKeyUsage(u); ok ***REMOVED***
				oids = append(oids, oid)
			***REMOVED*** else ***REMOVED***
				panic("internal error")
			***REMOVED***
		***REMOVED***

		oids = append(oids, template.UnknownExtKeyUsage...)

		ret[n].Value, err = asn1.Marshal(oids)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		n++
	***REMOVED***

	if template.BasicConstraintsValid && !oidInExtensions(oidExtensionBasicConstraints, template.ExtraExtensions) ***REMOVED***
		ret[n].Id = oidExtensionBasicConstraints
		ret[n].Value, err = asn1.Marshal(basicConstraints***REMOVED***template.IsCA, template.MaxPathLen***REMOVED***)
		ret[n].Critical = true
		if err != nil ***REMOVED***
			return
		***REMOVED***
		n++
	***REMOVED***

	if len(template.SubjectKeyId) > 0 && !oidInExtensions(oidExtensionSubjectKeyId, template.ExtraExtensions) ***REMOVED***
		ret[n].Id = oidExtensionSubjectKeyId
		ret[n].Value, err = asn1.Marshal(template.SubjectKeyId)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		n++
	***REMOVED***

	if len(template.AuthorityKeyId) > 0 && !oidInExtensions(oidExtensionAuthorityKeyId, template.ExtraExtensions) ***REMOVED***
		ret[n].Id = oidExtensionAuthorityKeyId
		ret[n].Value, err = asn1.Marshal(authKeyId***REMOVED***template.AuthorityKeyId***REMOVED***)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		n++
	***REMOVED***

	if (len(template.OCSPServer) > 0 || len(template.IssuingCertificateURL) > 0) &&
		!oidInExtensions(oidExtensionAuthorityInfoAccess, template.ExtraExtensions) ***REMOVED***
		ret[n].Id = oidExtensionAuthorityInfoAccess
		var aiaValues []authorityInfoAccess
		for _, name := range template.OCSPServer ***REMOVED***
			aiaValues = append(aiaValues, authorityInfoAccess***REMOVED***
				Method:   oidAuthorityInfoAccessOcsp,
				Location: asn1.RawValue***REMOVED***Tag: 6, Class: 2, Bytes: []byte(name)***REMOVED***,
			***REMOVED***)
		***REMOVED***
		for _, name := range template.IssuingCertificateURL ***REMOVED***
			aiaValues = append(aiaValues, authorityInfoAccess***REMOVED***
				Method:   oidAuthorityInfoAccessIssuers,
				Location: asn1.RawValue***REMOVED***Tag: 6, Class: 2, Bytes: []byte(name)***REMOVED***,
			***REMOVED***)
		***REMOVED***
		ret[n].Value, err = asn1.Marshal(aiaValues)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		n++
	***REMOVED***

	if (len(template.DNSNames) > 0 || len(template.EmailAddresses) > 0 || len(template.IPAddresses) > 0) &&
		!oidInExtensions(oidExtensionSubjectAltName, template.ExtraExtensions) ***REMOVED***
		ret[n].Id = oidExtensionSubjectAltName
		var rawValues []asn1.RawValue
		for _, name := range template.DNSNames ***REMOVED***
			rawValues = append(rawValues, asn1.RawValue***REMOVED***Tag: 2, Class: 2, Bytes: []byte(name)***REMOVED***)
		***REMOVED***
		for _, email := range template.EmailAddresses ***REMOVED***
			rawValues = append(rawValues, asn1.RawValue***REMOVED***Tag: 1, Class: 2, Bytes: []byte(email)***REMOVED***)
		***REMOVED***
		for _, rawIP := range template.IPAddresses ***REMOVED***
			// If possible, we always want to encode IPv4 addresses in 4 bytes.
			ip := rawIP.To4()
			if ip == nil ***REMOVED***
				ip = rawIP
			***REMOVED***
			rawValues = append(rawValues, asn1.RawValue***REMOVED***Tag: 7, Class: 2, Bytes: ip***REMOVED***)
		***REMOVED***
		ret[n].Value, err = asn1.Marshal(rawValues)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		n++
	***REMOVED***

	if len(template.PolicyIdentifiers) > 0 &&
		!oidInExtensions(oidExtensionCertificatePolicies, template.ExtraExtensions) ***REMOVED***
		ret[n].Id = oidExtensionCertificatePolicies
		policies := make([]policyInformation, len(template.PolicyIdentifiers))
		for i, policy := range template.PolicyIdentifiers ***REMOVED***
			policies[i].Policy = policy
		***REMOVED***
		ret[n].Value, err = asn1.Marshal(policies)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		n++
	***REMOVED***

	if len(template.PermittedDNSDomains) > 0 &&
		!oidInExtensions(oidExtensionNameConstraints, template.ExtraExtensions) ***REMOVED***
		ret[n].Id = oidExtensionNameConstraints
		ret[n].Critical = template.PermittedDNSDomainsCritical

		var out nameConstraints
		out.Permitted = make([]generalSubtree, len(template.PermittedDNSDomains))
		for i, permitted := range template.PermittedDNSDomains ***REMOVED***
			out.Permitted[i] = generalSubtree***REMOVED***Name: permitted***REMOVED***
		***REMOVED***
		ret[n].Value, err = asn1.Marshal(out)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		n++
	***REMOVED***

	if len(template.CRLDistributionPoints) > 0 &&
		!oidInExtensions(oidExtensionCRLDistributionPoints, template.ExtraExtensions) ***REMOVED***
		ret[n].Id = oidExtensionCRLDistributionPoints

		var crlDp []distributionPoint
		for _, name := range template.CRLDistributionPoints ***REMOVED***
			rawFullName, _ := asn1.Marshal(asn1.RawValue***REMOVED***Tag: 6, Class: 2, Bytes: []byte(name)***REMOVED***)

			dp := distributionPoint***REMOVED***
				DistributionPoint: distributionPointName***REMOVED***
					FullName: asn1.RawValue***REMOVED***Tag: 0, Class: 2, Bytes: rawFullName***REMOVED***,
				***REMOVED***,
			***REMOVED***
			crlDp = append(crlDp, dp)
		***REMOVED***

		ret[n].Value, err = asn1.Marshal(crlDp)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		n++
	***REMOVED***

	// Adding another extension here? Remember to update the maximum number
	// of elements in the make() at the top of the function.

	return append(ret[:n], template.ExtraExtensions...), nil
***REMOVED***

func subjectBytes(cert *Certificate) ([]byte, error) ***REMOVED***
	if len(cert.RawSubject) > 0 ***REMOVED***
		return cert.RawSubject, nil
	***REMOVED***

	return asn1.Marshal(cert.Subject.ToRDNSequence())
***REMOVED***

// CreateCertificate creates a new certificate based on a template. The
// following members of template are used: SerialNumber, Subject, NotBefore,
// NotAfter, KeyUsage, ExtKeyUsage, UnknownExtKeyUsage, BasicConstraintsValid,
// IsCA, MaxPathLen, SubjectKeyId, DNSNames, PermittedDNSDomainsCritical,
// PermittedDNSDomains.
//
// The certificate is signed by parent. If parent is equal to template then the
// certificate is self-signed. The parameter pub is the public key of the
// signee and priv is the private key of the signer.
//
// The returned slice is the certificate in DER encoding.
//
// The only supported key types are RSA and ECDSA (*rsa.PublicKey or
// *ecdsa.PublicKey for pub, *rsa.PrivateKey or *ecdsa.PublicKey for priv).
func CreateCertificate(rand io.Reader, template, parent *Certificate, pub interface***REMOVED******REMOVED***, priv interface***REMOVED******REMOVED***) (cert []byte, err error) ***REMOVED***
	var publicKeyBytes []byte
	var publicKeyAlgorithm pkix.AlgorithmIdentifier

	if publicKeyBytes, publicKeyAlgorithm, err = marshalPublicKey(pub); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var signatureAlgorithm pkix.AlgorithmIdentifier
	var hashFunc crypto.Hash

	switch priv := priv.(type) ***REMOVED***
	case *rsa.PrivateKey:
		signatureAlgorithm.Algorithm = oidSignatureSHA1WithRSA
		hashFunc = crypto.SHA1
	case *ecdsa.PrivateKey:
		switch priv.Curve ***REMOVED***
		case elliptic.P224(), elliptic.P256():
			hashFunc = crypto.SHA256
			signatureAlgorithm.Algorithm = oidSignatureECDSAWithSHA256
		case elliptic.P384():
			hashFunc = crypto.SHA384
			signatureAlgorithm.Algorithm = oidSignatureECDSAWithSHA384
		case elliptic.P521():
			hashFunc = crypto.SHA512
			signatureAlgorithm.Algorithm = oidSignatureECDSAWithSHA512
		default:
			return nil, errors.New("x509: unknown elliptic curve")
		***REMOVED***
	default:
		return nil, errors.New("x509: only RSA and ECDSA private keys supported")
	***REMOVED***

	if err != nil ***REMOVED***
		return
	***REMOVED***

	if len(parent.SubjectKeyId) > 0 ***REMOVED***
		template.AuthorityKeyId = parent.SubjectKeyId
	***REMOVED***

	extensions, err := buildExtensions(template)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	asn1Issuer, err := subjectBytes(parent)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	asn1Subject, err := subjectBytes(template)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	encodedPublicKey := asn1.BitString***REMOVED***BitLength: len(publicKeyBytes) * 8, Bytes: publicKeyBytes***REMOVED***
	c := tbsCertificate***REMOVED***
		Version:            2,
		SerialNumber:       template.SerialNumber,
		SignatureAlgorithm: signatureAlgorithm,
		Issuer:             asn1.RawValue***REMOVED***FullBytes: asn1Issuer***REMOVED***,
		Validity:           validity***REMOVED***template.NotBefore.UTC(), template.NotAfter.UTC()***REMOVED***,
		Subject:            asn1.RawValue***REMOVED***FullBytes: asn1Subject***REMOVED***,
		PublicKey:          publicKeyInfo***REMOVED***nil, publicKeyAlgorithm, encodedPublicKey***REMOVED***,
		Extensions:         extensions,
	***REMOVED***

	tbsCertContents, err := asn1.Marshal(c)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	c.Raw = tbsCertContents

	h := hashFunc.New()
	h.Write(tbsCertContents)
	digest := h.Sum(nil)

	var signature []byte

	switch priv := priv.(type) ***REMOVED***
	case *rsa.PrivateKey:
		signature, err = rsa.SignPKCS1v15(rand, priv, hashFunc, digest)
	case *ecdsa.PrivateKey:
		var r, s *big.Int
		if r, s, err = ecdsa.Sign(rand, priv, digest); err == nil ***REMOVED***
			signature, err = asn1.Marshal(ecdsaSignature***REMOVED***r, s***REMOVED***)
		***REMOVED***
	default:
		panic("internal error")
	***REMOVED***

	if err != nil ***REMOVED***
		return
	***REMOVED***

	cert, err = asn1.Marshal(certificate***REMOVED***
		nil,
		c,
		signatureAlgorithm,
		asn1.BitString***REMOVED***Bytes: signature, BitLength: len(signature) * 8***REMOVED***,
	***REMOVED***)
	return
***REMOVED***

// pemCRLPrefix is the magic string that indicates that we have a PEM encoded
// CRL.
var pemCRLPrefix = []byte("-----BEGIN X509 CRL")

// pemType is the type of a PEM encoded CRL.
var pemType = "X509 CRL"

// ParseCRL parses a CRL from the given bytes. It's often the case that PEM
// encoded CRLs will appear where they should be DER encoded, so this function
// will transparently handle PEM encoding as long as there isn't any leading
// garbage.
func ParseCRL(crlBytes []byte) (certList *pkix.CertificateList, err error) ***REMOVED***
	if bytes.HasPrefix(crlBytes, pemCRLPrefix) ***REMOVED***
		block, _ := pem.Decode(crlBytes)
		if block != nil && block.Type == pemType ***REMOVED***
			crlBytes = block.Bytes
		***REMOVED***
	***REMOVED***
	return ParseDERCRL(crlBytes)
***REMOVED***

// ParseDERCRL parses a DER encoded CRL from the given bytes.
func ParseDERCRL(derBytes []byte) (certList *pkix.CertificateList, err error) ***REMOVED***
	certList = new(pkix.CertificateList)
	_, err = asn1.Unmarshal(derBytes, certList)
	if err != nil ***REMOVED***
		certList = nil
	***REMOVED***
	return
***REMOVED***

// CreateCRL returns a DER encoded CRL, signed by this Certificate, that
// contains the given list of revoked certificates.
//
// The only supported key type is RSA (*rsa.PrivateKey for priv).
func (c *Certificate) CreateCRL(rand io.Reader, priv interface***REMOVED******REMOVED***, revokedCerts []pkix.RevokedCertificate, now, expiry time.Time) (crlBytes []byte, err error) ***REMOVED***
	rsaPriv, ok := priv.(*rsa.PrivateKey)
	if !ok ***REMOVED***
		return nil, errors.New("x509: non-RSA private keys not supported")
	***REMOVED***
	tbsCertList := pkix.TBSCertificateList***REMOVED***
		Version: 2,
		Signature: pkix.AlgorithmIdentifier***REMOVED***
			Algorithm: oidSignatureSHA1WithRSA,
		***REMOVED***,
		Issuer:              c.Subject.ToRDNSequence(),
		ThisUpdate:          now.UTC(),
		NextUpdate:          expiry.UTC(),
		RevokedCertificates: revokedCerts,
	***REMOVED***

	tbsCertListContents, err := asn1.Marshal(tbsCertList)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	h := sha1.New()
	h.Write(tbsCertListContents)
	digest := h.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand, rsaPriv, crypto.SHA1, digest)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	return asn1.Marshal(pkix.CertificateList***REMOVED***
		TBSCertList: tbsCertList,
		SignatureAlgorithm: pkix.AlgorithmIdentifier***REMOVED***
			Algorithm: oidSignatureSHA1WithRSA,
		***REMOVED***,
		SignatureValue: asn1.BitString***REMOVED***Bytes: signature, BitLength: len(signature) * 8***REMOVED***,
	***REMOVED***)
***REMOVED***
