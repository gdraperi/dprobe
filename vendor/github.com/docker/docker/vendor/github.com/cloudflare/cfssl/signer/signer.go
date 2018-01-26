// Package signer implements certificate signature functionality for CFSSL.
package signer

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/cloudflare/cfssl/certdb"
	"github.com/cloudflare/cfssl/config"
	"github.com/cloudflare/cfssl/csr"
	cferr "github.com/cloudflare/cfssl/errors"
	"github.com/cloudflare/cfssl/helpers"
	"github.com/cloudflare/cfssl/info"
)

// Subject contains the information that should be used to override the
// subject information when signing a certificate.
type Subject struct ***REMOVED***
	CN           string
	Names        []csr.Name `json:"names"`
	SerialNumber string
***REMOVED***

// Extension represents a raw extension to be included in the certificate.  The
// "value" field must be hex encoded.
type Extension struct ***REMOVED***
	ID       config.OID `json:"id"`
	Critical bool       `json:"critical"`
	Value    string     `json:"value"`
***REMOVED***

// SignRequest stores a signature request, which contains the hostname,
// the CSR, optional subject information, and the signature profile.
//
// Extensions provided in the signRequest are copied into the certificate, as
// long as they are in the ExtensionWhitelist for the signer's policy.
// Extensions requested in the CSR are ignored, except for those processed by
// ParseCertificateRequest (mainly subjectAltName).
type SignRequest struct ***REMOVED***
	Hosts       []string    `json:"hosts"`
	Request     string      `json:"certificate_request"`
	Subject     *Subject    `json:"subject,omitempty"`
	Profile     string      `json:"profile"`
	CRLOverride string      `json:"crl_override"`
	Label       string      `json:"label"`
	Serial      *big.Int    `json:"serial,omitempty"`
	Extensions  []Extension `json:"extensions,omitempty"`
***REMOVED***

// appendIf appends to a if s is not an empty string.
func appendIf(s string, a *[]string) ***REMOVED***
	if s != "" ***REMOVED***
		*a = append(*a, s)
	***REMOVED***
***REMOVED***

// Name returns the PKIX name for the subject.
func (s *Subject) Name() pkix.Name ***REMOVED***
	var name pkix.Name
	name.CommonName = s.CN

	for _, n := range s.Names ***REMOVED***
		appendIf(n.C, &name.Country)
		appendIf(n.ST, &name.Province)
		appendIf(n.L, &name.Locality)
		appendIf(n.O, &name.Organization)
		appendIf(n.OU, &name.OrganizationalUnit)
	***REMOVED***
	name.SerialNumber = s.SerialNumber
	return name
***REMOVED***

// SplitHosts takes a comma-spearated list of hosts and returns a slice
// with the hosts split
func SplitHosts(hostList string) []string ***REMOVED***
	if hostList == "" ***REMOVED***
		return nil
	***REMOVED***

	return strings.Split(hostList, ",")
***REMOVED***

// A Signer contains a CA's certificate and private key for signing
// certificates, a Signing policy to refer to and a SignatureAlgorithm.
type Signer interface ***REMOVED***
	Info(info.Req) (*info.Resp, error)
	Policy() *config.Signing
	SetDBAccessor(certdb.Accessor)
	SetPolicy(*config.Signing)
	SigAlgo() x509.SignatureAlgorithm
	Sign(req SignRequest) (cert []byte, err error)
***REMOVED***

// Profile gets the specific profile from the signer
func Profile(s Signer, profile string) (*config.SigningProfile, error) ***REMOVED***
	var p *config.SigningProfile
	policy := s.Policy()
	if policy != nil && policy.Profiles != nil && profile != "" ***REMOVED***
		p = policy.Profiles[profile]
	***REMOVED***

	if p == nil && policy != nil ***REMOVED***
		p = policy.Default
	***REMOVED***

	if p == nil ***REMOVED***
		return nil, cferr.Wrap(cferr.APIClientError, cferr.ClientHTTPError, errors.New("profile must not be nil"))
	***REMOVED***
	return p, nil
***REMOVED***

// DefaultSigAlgo returns an appropriate X.509 signature algorithm given
// the CA's private key.
func DefaultSigAlgo(priv crypto.Signer) x509.SignatureAlgorithm ***REMOVED***
	pub := priv.Public()
	switch pub := pub.(type) ***REMOVED***
	case *rsa.PublicKey:
		keySize := pub.N.BitLen()
		switch ***REMOVED***
		case keySize >= 4096:
			return x509.SHA512WithRSA
		case keySize >= 3072:
			return x509.SHA384WithRSA
		case keySize >= 2048:
			return x509.SHA256WithRSA
		default:
			return x509.SHA1WithRSA
		***REMOVED***
	case *ecdsa.PublicKey:
		switch pub.Curve ***REMOVED***
		case elliptic.P256():
			return x509.ECDSAWithSHA256
		case elliptic.P384():
			return x509.ECDSAWithSHA384
		case elliptic.P521():
			return x509.ECDSAWithSHA512
		default:
			return x509.ECDSAWithSHA1
		***REMOVED***
	default:
		return x509.UnknownSignatureAlgorithm
	***REMOVED***
***REMOVED***

// ParseCertificateRequest takes an incoming certificate request and
// builds a certificate template from it.
func ParseCertificateRequest(s Signer, csrBytes []byte) (template *x509.Certificate, err error) ***REMOVED***
	csrv, err := x509.ParseCertificateRequest(csrBytes)
	if err != nil ***REMOVED***
		err = cferr.Wrap(cferr.CSRError, cferr.ParseFailed, err)
		return
	***REMOVED***

	err = helpers.CheckSignature(csrv, csrv.SignatureAlgorithm, csrv.RawTBSCertificateRequest, csrv.Signature)
	if err != nil ***REMOVED***
		err = cferr.Wrap(cferr.CSRError, cferr.KeyMismatch, err)
		return
	***REMOVED***

	template = &x509.Certificate***REMOVED***
		Subject:            csrv.Subject,
		PublicKeyAlgorithm: csrv.PublicKeyAlgorithm,
		PublicKey:          csrv.PublicKey,
		SignatureAlgorithm: s.SigAlgo(),
		DNSNames:           csrv.DNSNames,
		IPAddresses:        csrv.IPAddresses,
		EmailAddresses:     csrv.EmailAddresses,
	***REMOVED***

	for _, val := range csrv.Extensions ***REMOVED***
		// Check the CSR for the X.509 BasicConstraints (RFC 5280, 4.2.1.9)
		// extension and append to template if necessary
		if val.Id.Equal(asn1.ObjectIdentifier***REMOVED***2, 5, 29, 19***REMOVED***) ***REMOVED***
			var constraints csr.BasicConstraints
			var rest []byte

			if rest, err = asn1.Unmarshal(val.Value, &constraints); err != nil ***REMOVED***
				return nil, cferr.Wrap(cferr.CSRError, cferr.ParseFailed, err)
			***REMOVED*** else if len(rest) != 0 ***REMOVED***
				return nil, cferr.Wrap(cferr.CSRError, cferr.ParseFailed, errors.New("x509: trailing data after X.509 BasicConstraints"))
			***REMOVED***

			template.BasicConstraintsValid = true
			template.IsCA = constraints.IsCA
			template.MaxPathLen = constraints.MaxPathLen
			template.MaxPathLenZero = template.MaxPathLen == 0
		***REMOVED***
	***REMOVED***

	return
***REMOVED***

type subjectPublicKeyInfo struct ***REMOVED***
	Algorithm        pkix.AlgorithmIdentifier
	SubjectPublicKey asn1.BitString
***REMOVED***

// ComputeSKI derives an SKI from the certificate's public key in a
// standard manner. This is done by computing the SHA-1 digest of the
// SubjectPublicKeyInfo component of the certificate.
func ComputeSKI(template *x509.Certificate) ([]byte, error) ***REMOVED***
	pub := template.PublicKey
	encodedPub, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var subPKI subjectPublicKeyInfo
	_, err = asn1.Unmarshal(encodedPub, &subPKI)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	pubHash := sha1.Sum(subPKI.SubjectPublicKey.Bytes)
	return pubHash[:], nil
***REMOVED***

// FillTemplate is a utility function that tries to load as much of
// the certificate template as possible from the profiles and current
// template. It fills in the key uses, expiration, revocation URLs
// and SKI.
func FillTemplate(template *x509.Certificate, defaultProfile, profile *config.SigningProfile) error ***REMOVED***
	ski, err := ComputeSKI(template)

	var (
		eku             []x509.ExtKeyUsage
		ku              x509.KeyUsage
		backdate        time.Duration
		expiry          time.Duration
		notBefore       time.Time
		notAfter        time.Time
		crlURL, ocspURL string
		issuerURL       = profile.IssuerURL
	)

	// The third value returned from Usages is a list of unknown key usages.
	// This should be used when validating the profile at load, and isn't used
	// here.
	ku, eku, _ = profile.Usages()
	if profile.IssuerURL == nil ***REMOVED***
		issuerURL = defaultProfile.IssuerURL
	***REMOVED***

	if ku == 0 && len(eku) == 0 ***REMOVED***
		return cferr.New(cferr.PolicyError, cferr.NoKeyUsages)
	***REMOVED***

	if expiry = profile.Expiry; expiry == 0 ***REMOVED***
		expiry = defaultProfile.Expiry
	***REMOVED***

	if crlURL = profile.CRL; crlURL == "" ***REMOVED***
		crlURL = defaultProfile.CRL
	***REMOVED***
	if ocspURL = profile.OCSP; ocspURL == "" ***REMOVED***
		ocspURL = defaultProfile.OCSP
	***REMOVED***
	if backdate = profile.Backdate; backdate == 0 ***REMOVED***
		backdate = -5 * time.Minute
	***REMOVED*** else ***REMOVED***
		backdate = -1 * profile.Backdate
	***REMOVED***

	if !profile.NotBefore.IsZero() ***REMOVED***
		notBefore = profile.NotBefore.UTC()
	***REMOVED*** else ***REMOVED***
		notBefore = time.Now().Round(time.Minute).Add(backdate).UTC()
	***REMOVED***

	if !profile.NotAfter.IsZero() ***REMOVED***
		notAfter = profile.NotAfter.UTC()
	***REMOVED*** else ***REMOVED***
		notAfter = notBefore.Add(expiry).UTC()
	***REMOVED***

	template.NotBefore = notBefore
	template.NotAfter = notAfter
	template.KeyUsage = ku
	template.ExtKeyUsage = eku
	template.BasicConstraintsValid = true
	template.IsCA = profile.CAConstraint.IsCA
	if template.IsCA ***REMOVED***
		template.MaxPathLen = profile.CAConstraint.MaxPathLen
		if template.MaxPathLen == 0 ***REMOVED***
			template.MaxPathLenZero = profile.CAConstraint.MaxPathLenZero
		***REMOVED***
		template.DNSNames = nil
		template.EmailAddresses = nil
	***REMOVED***
	template.SubjectKeyId = ski

	if ocspURL != "" ***REMOVED***
		template.OCSPServer = []string***REMOVED***ocspURL***REMOVED***
	***REMOVED***
	if crlURL != "" ***REMOVED***
		template.CRLDistributionPoints = []string***REMOVED***crlURL***REMOVED***
	***REMOVED***

	if len(issuerURL) != 0 ***REMOVED***
		template.IssuingCertificateURL = issuerURL
	***REMOVED***
	if len(profile.Policies) != 0 ***REMOVED***
		err = addPolicies(template, profile.Policies)
		if err != nil ***REMOVED***
			return cferr.Wrap(cferr.PolicyError, cferr.InvalidPolicy, err)
		***REMOVED***
	***REMOVED***
	if profile.OCSPNoCheck ***REMOVED***
		ocspNoCheckExtension := pkix.Extension***REMOVED***
			Id:       asn1.ObjectIdentifier***REMOVED***1, 3, 6, 1, 5, 5, 7, 48, 1, 5***REMOVED***,
			Critical: false,
			Value:    []byte***REMOVED***0x05, 0x00***REMOVED***,
		***REMOVED***
		template.ExtraExtensions = append(template.ExtraExtensions, ocspNoCheckExtension)
	***REMOVED***

	return nil
***REMOVED***

type policyInformation struct ***REMOVED***
	PolicyIdentifier asn1.ObjectIdentifier
	Qualifiers       []interface***REMOVED******REMOVED*** `asn1:"tag:optional,omitempty"`
***REMOVED***

type cpsPolicyQualifier struct ***REMOVED***
	PolicyQualifierID asn1.ObjectIdentifier
	Qualifier         string `asn1:"tag:optional,ia5"`
***REMOVED***

type userNotice struct ***REMOVED***
	ExplicitText string `asn1:"tag:optional,utf8"`
***REMOVED***
type userNoticePolicyQualifier struct ***REMOVED***
	PolicyQualifierID asn1.ObjectIdentifier
	Qualifier         userNotice
***REMOVED***

var (
	// Per https://tools.ietf.org/html/rfc3280.html#page-106, this represents:
	// iso(1) identified-organization(3) dod(6) internet(1) security(5)
	//   mechanisms(5) pkix(7) id-qt(2) id-qt-cps(1)
	iDQTCertificationPracticeStatement = asn1.ObjectIdentifier***REMOVED***1, 3, 6, 1, 5, 5, 7, 2, 1***REMOVED***
	// iso(1) identified-organization(3) dod(6) internet(1) security(5)
	//   mechanisms(5) pkix(7) id-qt(2) id-qt-unotice(2)
	iDQTUserNotice = asn1.ObjectIdentifier***REMOVED***1, 3, 6, 1, 5, 5, 7, 2, 2***REMOVED***

	// CTPoisonOID is the object ID of the critical poison extension for precertificates
	// https://tools.ietf.org/html/rfc6962#page-9
	CTPoisonOID = asn1.ObjectIdentifier***REMOVED***1, 3, 6, 1, 4, 1, 11129, 2, 4, 3***REMOVED***

	// SCTListOID is the object ID for the Signed Certificate Timestamp certificate extension
	// https://tools.ietf.org/html/rfc6962#page-14
	SCTListOID = asn1.ObjectIdentifier***REMOVED***1, 3, 6, 1, 4, 1, 11129, 2, 4, 2***REMOVED***
)

// addPolicies adds Certificate Policies and optional Policy Qualifiers to a
// certificate, based on the input config. Go's x509 library allows setting
// Certificate Policies easily, but does not support nested Policy Qualifiers
// under those policies. So we need to construct the ASN.1 structure ourselves.
func addPolicies(template *x509.Certificate, policies []config.CertificatePolicy) error ***REMOVED***
	asn1PolicyList := []policyInformation***REMOVED******REMOVED***

	for _, policy := range policies ***REMOVED***
		pi := policyInformation***REMOVED***
			// The PolicyIdentifier is an OID assigned to a given issuer.
			PolicyIdentifier: asn1.ObjectIdentifier(policy.ID),
		***REMOVED***
		for _, qualifier := range policy.Qualifiers ***REMOVED***
			switch qualifier.Type ***REMOVED***
			case "id-qt-unotice":
				pi.Qualifiers = append(pi.Qualifiers,
					userNoticePolicyQualifier***REMOVED***
						PolicyQualifierID: iDQTUserNotice,
						Qualifier: userNotice***REMOVED***
							ExplicitText: qualifier.Value,
						***REMOVED***,
					***REMOVED***)
			case "id-qt-cps":
				pi.Qualifiers = append(pi.Qualifiers,
					cpsPolicyQualifier***REMOVED***
						PolicyQualifierID: iDQTCertificationPracticeStatement,
						Qualifier:         qualifier.Value,
					***REMOVED***)
			default:
				return errors.New("Invalid qualifier type in Policies " + qualifier.Type)
			***REMOVED***
		***REMOVED***
		asn1PolicyList = append(asn1PolicyList, pi)
	***REMOVED***

	asn1Bytes, err := asn1.Marshal(asn1PolicyList)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	template.ExtraExtensions = append(template.ExtraExtensions, pkix.Extension***REMOVED***
		Id:       asn1.ObjectIdentifier***REMOVED***2, 5, 29, 32***REMOVED***,
		Critical: false,
		Value:    asn1Bytes,
	***REMOVED***)
	return nil
***REMOVED***
