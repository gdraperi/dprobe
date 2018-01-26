// Package helpers implements utility functionality common to many
// CFSSL packages.
package helpers

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math/big"

	"strings"
	"time"

	"github.com/cloudflare/cfssl/crypto/pkcs7"
	cferr "github.com/cloudflare/cfssl/errors"
	"github.com/cloudflare/cfssl/helpers/derhelpers"
	"github.com/cloudflare/cfssl/log"
	"golang.org/x/crypto/pkcs12"
)

// OneYear is a time.Duration representing a year's worth of seconds.
const OneYear = 8760 * time.Hour

// OneDay is a time.Duration representing a day's worth of seconds.
const OneDay = 24 * time.Hour

// InclusiveDate returns the time.Time representation of a date - 1
// nanosecond. This allows time.After to be used inclusively.
func InclusiveDate(year int, month time.Month, day int) time.Time ***REMOVED***
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Add(-1 * time.Nanosecond)
***REMOVED***

// Jul2012 is the July 2012 CAB Forum deadline for when CAs must stop
// issuing certificates valid for more than 5 years.
var Jul2012 = InclusiveDate(2012, time.July, 01)

// Apr2015 is the April 2015 CAB Forum deadline for when CAs must stop
// issuing certificates valid for more than 39 months.
var Apr2015 = InclusiveDate(2015, time.April, 01)

// KeyLength returns the bit size of ECDSA or RSA PublicKey
func KeyLength(key interface***REMOVED******REMOVED***) int ***REMOVED***
	if key == nil ***REMOVED***
		return 0
	***REMOVED***
	if ecdsaKey, ok := key.(*ecdsa.PublicKey); ok ***REMOVED***
		return ecdsaKey.Curve.Params().BitSize
	***REMOVED*** else if rsaKey, ok := key.(*rsa.PublicKey); ok ***REMOVED***
		return rsaKey.N.BitLen()
	***REMOVED***

	return 0
***REMOVED***

// ExpiryTime returns the time when the certificate chain is expired.
func ExpiryTime(chain []*x509.Certificate) (notAfter time.Time) ***REMOVED***
	if len(chain) == 0 ***REMOVED***
		return
	***REMOVED***

	notAfter = chain[0].NotAfter
	for _, cert := range chain ***REMOVED***
		if notAfter.After(cert.NotAfter) ***REMOVED***
			notAfter = cert.NotAfter
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// MonthsValid returns the number of months for which a certificate is valid.
func MonthsValid(c *x509.Certificate) int ***REMOVED***
	issued := c.NotBefore
	expiry := c.NotAfter
	years := (expiry.Year() - issued.Year())
	months := years*12 + int(expiry.Month()) - int(issued.Month())

	// Round up if valid for less than a full month
	if expiry.Day() > issued.Day() ***REMOVED***
		months++
	***REMOVED***
	return months
***REMOVED***

// ValidExpiry determines if a certificate is valid for an acceptable
// length of time per the CA/Browser Forum baseline requirements.
// See https://cabforum.org/wp-content/uploads/CAB-Forum-BR-1.3.0.pdf
func ValidExpiry(c *x509.Certificate) bool ***REMOVED***
	issued := c.NotBefore

	var maxMonths int
	switch ***REMOVED***
	case issued.After(Apr2015):
		maxMonths = 39
	case issued.After(Jul2012):
		maxMonths = 60
	case issued.Before(Jul2012):
		maxMonths = 120
	***REMOVED***

	if MonthsValid(c) > maxMonths ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

// SignatureString returns the TLS signature string corresponding to
// an X509 signature algorithm.
func SignatureString(alg x509.SignatureAlgorithm) string ***REMOVED***
	switch alg ***REMOVED***
	case x509.MD2WithRSA:
		return "MD2WithRSA"
	case x509.MD5WithRSA:
		return "MD5WithRSA"
	case x509.SHA1WithRSA:
		return "SHA1WithRSA"
	case x509.SHA256WithRSA:
		return "SHA256WithRSA"
	case x509.SHA384WithRSA:
		return "SHA384WithRSA"
	case x509.SHA512WithRSA:
		return "SHA512WithRSA"
	case x509.DSAWithSHA1:
		return "DSAWithSHA1"
	case x509.DSAWithSHA256:
		return "DSAWithSHA256"
	case x509.ECDSAWithSHA1:
		return "ECDSAWithSHA1"
	case x509.ECDSAWithSHA256:
		return "ECDSAWithSHA256"
	case x509.ECDSAWithSHA384:
		return "ECDSAWithSHA384"
	case x509.ECDSAWithSHA512:
		return "ECDSAWithSHA512"
	default:
		return "Unknown Signature"
	***REMOVED***
***REMOVED***

// HashAlgoString returns the hash algorithm name contains in the signature
// method.
func HashAlgoString(alg x509.SignatureAlgorithm) string ***REMOVED***
	switch alg ***REMOVED***
	case x509.MD2WithRSA:
		return "MD2"
	case x509.MD5WithRSA:
		return "MD5"
	case x509.SHA1WithRSA:
		return "SHA1"
	case x509.SHA256WithRSA:
		return "SHA256"
	case x509.SHA384WithRSA:
		return "SHA384"
	case x509.SHA512WithRSA:
		return "SHA512"
	case x509.DSAWithSHA1:
		return "SHA1"
	case x509.DSAWithSHA256:
		return "SHA256"
	case x509.ECDSAWithSHA1:
		return "SHA1"
	case x509.ECDSAWithSHA256:
		return "SHA256"
	case x509.ECDSAWithSHA384:
		return "SHA384"
	case x509.ECDSAWithSHA512:
		return "SHA512"
	default:
		return "Unknown Hash Algorithm"
	***REMOVED***
***REMOVED***

// EncodeCertificatesPEM encodes a number of x509 certficates to PEM
func EncodeCertificatesPEM(certs []*x509.Certificate) []byte ***REMOVED***
	var buffer bytes.Buffer
	for _, cert := range certs ***REMOVED***
		pem.Encode(&buffer, &pem.Block***REMOVED***
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		***REMOVED***)
	***REMOVED***

	return buffer.Bytes()
***REMOVED***

// EncodeCertificatePEM encodes a single x509 certficates to PEM
func EncodeCertificatePEM(cert *x509.Certificate) []byte ***REMOVED***
	return EncodeCertificatesPEM([]*x509.Certificate***REMOVED***cert***REMOVED***)
***REMOVED***

// ParseCertificatesPEM parses a sequence of PEM-encoded certificate and returns them,
// can handle PEM encoded PKCS #7 structures.
func ParseCertificatesPEM(certsPEM []byte) ([]*x509.Certificate, error) ***REMOVED***
	var certs []*x509.Certificate
	var err error
	certsPEM = bytes.TrimSpace(certsPEM)
	for len(certsPEM) > 0 ***REMOVED***
		var cert []*x509.Certificate
		cert, certsPEM, err = ParseOneCertificateFromPEM(certsPEM)
		if err != nil ***REMOVED***

			return nil, cferr.New(cferr.CertificateError, cferr.ParseFailed)
		***REMOVED*** else if cert == nil ***REMOVED***
			break
		***REMOVED***

		certs = append(certs, cert...)
	***REMOVED***
	if len(certsPEM) > 0 ***REMOVED***
		return nil, cferr.New(cferr.CertificateError, cferr.DecodeFailed)
	***REMOVED***
	return certs, nil
***REMOVED***

// ParseCertificatesDER parses a DER encoding of a certificate object and possibly private key,
// either PKCS #7, PKCS #12, or raw x509.
func ParseCertificatesDER(certsDER []byte, password string) (certs []*x509.Certificate, key crypto.Signer, err error) ***REMOVED***
	certsDER = bytes.TrimSpace(certsDER)
	pkcs7data, err := pkcs7.ParsePKCS7(certsDER)
	if err != nil ***REMOVED***
		var pkcs12data interface***REMOVED******REMOVED***
		certs = make([]*x509.Certificate, 1)
		pkcs12data, certs[0], err = pkcs12.Decode(certsDER, password)
		if err != nil ***REMOVED***
			certs, err = x509.ParseCertificates(certsDER)
			if err != nil ***REMOVED***
				return nil, nil, cferr.New(cferr.CertificateError, cferr.DecodeFailed)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			key = pkcs12data.(crypto.Signer)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if pkcs7data.ContentInfo != "SignedData" ***REMOVED***
			return nil, nil, cferr.Wrap(cferr.CertificateError, cferr.DecodeFailed, errors.New("can only extract certificates from signed data content info"))
		***REMOVED***
		certs = pkcs7data.Content.SignedData.Certificates
	***REMOVED***
	if certs == nil ***REMOVED***
		return nil, key, cferr.New(cferr.CertificateError, cferr.DecodeFailed)
	***REMOVED***
	return certs, key, nil
***REMOVED***

// ParseSelfSignedCertificatePEM parses a PEM-encoded certificate and check if it is self-signed.
func ParseSelfSignedCertificatePEM(certPEM []byte) (*x509.Certificate, error) ***REMOVED***
	cert, err := ParseCertificatePEM(certPEM)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := cert.CheckSignature(cert.SignatureAlgorithm, cert.RawTBSCertificate, cert.Signature); err != nil ***REMOVED***
		return nil, cferr.Wrap(cferr.CertificateError, cferr.VerifyFailed, err)
	***REMOVED***
	return cert, nil
***REMOVED***

// ParseCertificatePEM parses and returns a PEM-encoded certificate,
// can handle PEM encoded PKCS #7 structures.
func ParseCertificatePEM(certPEM []byte) (*x509.Certificate, error) ***REMOVED***
	certPEM = bytes.TrimSpace(certPEM)
	cert, rest, err := ParseOneCertificateFromPEM(certPEM)
	if err != nil ***REMOVED***
		// Log the actual parsing error but throw a default parse error message.
		log.Debugf("Certificate parsing error: %v", err)
		return nil, cferr.New(cferr.CertificateError, cferr.ParseFailed)
	***REMOVED*** else if cert == nil ***REMOVED***
		return nil, cferr.New(cferr.CertificateError, cferr.DecodeFailed)
	***REMOVED*** else if len(rest) > 0 ***REMOVED***
		return nil, cferr.Wrap(cferr.CertificateError, cferr.ParseFailed, errors.New("the PEM file should contain only one object"))
	***REMOVED*** else if len(cert) > 1 ***REMOVED***
		return nil, cferr.Wrap(cferr.CertificateError, cferr.ParseFailed, errors.New("the PKCS7 object in the PEM file should contain only one certificate"))
	***REMOVED***
	return cert[0], nil
***REMOVED***

// ParseOneCertificateFromPEM attempts to parse one PEM encoded certificate object,
// either a raw x509 certificate or a PKCS #7 structure possibly containing
// multiple certificates, from the top of certsPEM, which itself may
// contain multiple PEM encoded certificate objects.
func ParseOneCertificateFromPEM(certsPEM []byte) ([]*x509.Certificate, []byte, error) ***REMOVED***

	block, rest := pem.Decode(certsPEM)
	if block == nil ***REMOVED***
		return nil, rest, nil
	***REMOVED***

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil ***REMOVED***
		pkcs7data, err := pkcs7.ParsePKCS7(block.Bytes)
		if err != nil ***REMOVED***
			return nil, rest, err
		***REMOVED***
		if pkcs7data.ContentInfo != "SignedData" ***REMOVED***
			return nil, rest, errors.New("only PKCS #7 Signed Data Content Info supported for certificate parsing")
		***REMOVED***
		certs := pkcs7data.Content.SignedData.Certificates
		if certs == nil ***REMOVED***
			return nil, rest, errors.New("PKCS #7 structure contains no certificates")
		***REMOVED***
		return certs, rest, nil
	***REMOVED***
	var certs = []*x509.Certificate***REMOVED***cert***REMOVED***
	return certs, rest, nil
***REMOVED***

// LoadPEMCertPool loads a pool of PEM certificates from file.
func LoadPEMCertPool(certsFile string) (*x509.CertPool, error) ***REMOVED***
	if certsFile == "" ***REMOVED***
		return nil, nil
	***REMOVED***
	pemCerts, err := ioutil.ReadFile(certsFile)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return PEMToCertPool(pemCerts)
***REMOVED***

// PEMToCertPool concerts PEM certificates to a CertPool.
func PEMToCertPool(pemCerts []byte) (*x509.CertPool, error) ***REMOVED***
	if len(pemCerts) == 0 ***REMOVED***
		return nil, nil
	***REMOVED***

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemCerts) ***REMOVED***
		return nil, errors.New("failed to load cert pool")
	***REMOVED***

	return certPool, nil
***REMOVED***

// ParsePrivateKeyPEM parses and returns a PEM-encoded private
// key. The private key may be either an unencrypted PKCS#8, PKCS#1,
// or elliptic private key.
func ParsePrivateKeyPEM(keyPEM []byte) (key crypto.Signer, err error) ***REMOVED***
	return ParsePrivateKeyPEMWithPassword(keyPEM, nil)
***REMOVED***

// ParsePrivateKeyPEMWithPassword parses and returns a PEM-encoded private
// key. The private key may be a potentially encrypted PKCS#8, PKCS#1,
// or elliptic private key.
func ParsePrivateKeyPEMWithPassword(keyPEM []byte, password []byte) (key crypto.Signer, err error) ***REMOVED***
	keyDER, err := GetKeyDERFromPEM(keyPEM, password)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return derhelpers.ParsePrivateKeyDER(keyDER)
***REMOVED***

// GetKeyDERFromPEM parses a PEM-encoded private key and returns DER-format key bytes.
func GetKeyDERFromPEM(in []byte, password []byte) ([]byte, error) ***REMOVED***
	keyDER, _ := pem.Decode(in)
	if keyDER != nil ***REMOVED***
		if procType, ok := keyDER.Headers["Proc-Type"]; ok ***REMOVED***
			if strings.Contains(procType, "ENCRYPTED") ***REMOVED***
				if password != nil ***REMOVED***
					return x509.DecryptPEMBlock(keyDER, password)
				***REMOVED***
				return nil, cferr.New(cferr.PrivateKeyError, cferr.Encrypted)
			***REMOVED***
		***REMOVED***
		return keyDER.Bytes, nil
	***REMOVED***

	return nil, cferr.New(cferr.PrivateKeyError, cferr.DecodeFailed)
***REMOVED***

// CheckSignature verifies a signature made by the key on a CSR, such
// as on the CSR itself.
func CheckSignature(csr *x509.CertificateRequest, algo x509.SignatureAlgorithm, signed, signature []byte) error ***REMOVED***
	var hashType crypto.Hash

	switch algo ***REMOVED***
	case x509.SHA1WithRSA, x509.ECDSAWithSHA1:
		hashType = crypto.SHA1
	case x509.SHA256WithRSA, x509.ECDSAWithSHA256:
		hashType = crypto.SHA256
	case x509.SHA384WithRSA, x509.ECDSAWithSHA384:
		hashType = crypto.SHA384
	case x509.SHA512WithRSA, x509.ECDSAWithSHA512:
		hashType = crypto.SHA512
	default:
		return x509.ErrUnsupportedAlgorithm
	***REMOVED***

	if !hashType.Available() ***REMOVED***
		return x509.ErrUnsupportedAlgorithm
	***REMOVED***
	h := hashType.New()

	h.Write(signed)
	digest := h.Sum(nil)

	switch pub := csr.PublicKey.(type) ***REMOVED***
	case *rsa.PublicKey:
		return rsa.VerifyPKCS1v15(pub, hashType, digest, signature)
	case *ecdsa.PublicKey:
		ecdsaSig := new(struct***REMOVED*** R, S *big.Int ***REMOVED***)
		if _, err := asn1.Unmarshal(signature, ecdsaSig); err != nil ***REMOVED***
			return err
		***REMOVED***
		if ecdsaSig.R.Sign() <= 0 || ecdsaSig.S.Sign() <= 0 ***REMOVED***
			return errors.New("x509: ECDSA signature contained zero or negative values")
		***REMOVED***
		if !ecdsa.Verify(pub, digest, ecdsaSig.R, ecdsaSig.S) ***REMOVED***
			return errors.New("x509: ECDSA verification failure")
		***REMOVED***
		return nil
	***REMOVED***
	return x509.ErrUnsupportedAlgorithm
***REMOVED***

// ParseCSR parses a PEM- or DER-encoded PKCS #10 certificate signing request.
func ParseCSR(in []byte) (csr *x509.CertificateRequest, rest []byte, err error) ***REMOVED***
	in = bytes.TrimSpace(in)
	p, rest := pem.Decode(in)
	if p != nil ***REMOVED***
		if p.Type != "NEW CERTIFICATE REQUEST" && p.Type != "CERTIFICATE REQUEST" ***REMOVED***
			return nil, rest, cferr.New(cferr.CSRError, cferr.BadRequest)
		***REMOVED***

		csr, err = x509.ParseCertificateRequest(p.Bytes)
	***REMOVED*** else ***REMOVED***
		csr, err = x509.ParseCertificateRequest(in)
	***REMOVED***

	if err != nil ***REMOVED***
		return nil, rest, err
	***REMOVED***

	err = CheckSignature(csr, csr.SignatureAlgorithm, csr.RawTBSCertificateRequest, csr.Signature)
	if err != nil ***REMOVED***
		return nil, rest, err
	***REMOVED***

	return csr, rest, nil
***REMOVED***

// ParseCSRPEM parses a PEM-encoded certificiate signing request.
// It does not check the signature. This is useful for dumping data from a CSR
// locally.
func ParseCSRPEM(csrPEM []byte) (*x509.CertificateRequest, error) ***REMOVED***
	block, _ := pem.Decode([]byte(csrPEM))
	der := block.Bytes
	csrObject, err := x509.ParseCertificateRequest(der)

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return csrObject, nil
***REMOVED***

// SignerAlgo returns an X.509 signature algorithm from a crypto.Signer.
func SignerAlgo(priv crypto.Signer) x509.SignatureAlgorithm ***REMOVED***
	switch pub := priv.Public().(type) ***REMOVED***
	case *rsa.PublicKey:
		bitLength := pub.N.BitLen()
		switch ***REMOVED***
		case bitLength >= 4096:
			return x509.SHA512WithRSA
		case bitLength >= 3072:
			return x509.SHA384WithRSA
		case bitLength >= 2048:
			return x509.SHA256WithRSA
		default:
			return x509.SHA1WithRSA
		***REMOVED***
	case *ecdsa.PublicKey:
		switch pub.Curve ***REMOVED***
		case elliptic.P521():
			return x509.ECDSAWithSHA512
		case elliptic.P384():
			return x509.ECDSAWithSHA384
		case elliptic.P256():
			return x509.ECDSAWithSHA256
		default:
			return x509.ECDSAWithSHA1
		***REMOVED***
	default:
		return x509.UnknownSignatureAlgorithm
	***REMOVED***
***REMOVED***

// LoadClientCertificate load key/certificate from pem files
func LoadClientCertificate(certFile string, keyFile string) (*tls.Certificate, error) ***REMOVED***
	if certFile != "" && keyFile != "" ***REMOVED***
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil ***REMOVED***
			log.Critical("Unable to read client certificate from file: %s or key from file: %s", certFile, keyFile)
			return nil, err
		***REMOVED***
		log.Debug("Client certificate loaded ")
		return &cert, nil
	***REMOVED***
	return nil, nil
***REMOVED***

// CreateTLSConfig creates a tls.Config object from certs and roots
func CreateTLSConfig(remoteCAs *x509.CertPool, cert *tls.Certificate) *tls.Config ***REMOVED***
	var certs []tls.Certificate
	if cert != nil ***REMOVED***
		certs = []tls.Certificate***REMOVED****cert***REMOVED***
	***REMOVED***
	return &tls.Config***REMOVED***
		Certificates: certs,
		RootCAs:      remoteCAs,
	***REMOVED***
***REMOVED***
