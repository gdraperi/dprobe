// Package csr implements certificate requests for CFSSL.
package csr

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"net"
	"net/mail"
	"strings"

	cferr "github.com/cloudflare/cfssl/errors"
	"github.com/cloudflare/cfssl/helpers"
	"github.com/cloudflare/cfssl/log"
)

const (
	curveP256 = 256
	curveP384 = 384
	curveP521 = 521
)

// A Name contains the SubjectInfo fields.
type Name struct ***REMOVED***
	C            string // Country
	ST           string // State
	L            string // Locality
	O            string // OrganisationName
	OU           string // OrganisationalUnitName
	SerialNumber string
***REMOVED***

// A KeyRequest is a generic request for a new key.
type KeyRequest interface ***REMOVED***
	Algo() string
	Size() int
	Generate() (crypto.PrivateKey, error)
	SigAlgo() x509.SignatureAlgorithm
***REMOVED***

// A BasicKeyRequest contains the algorithm and key size for a new private key.
type BasicKeyRequest struct ***REMOVED***
	A string `json:"algo"`
	S int    `json:"size"`
***REMOVED***

// NewBasicKeyRequest returns a default BasicKeyRequest.
func NewBasicKeyRequest() *BasicKeyRequest ***REMOVED***
	return &BasicKeyRequest***REMOVED***"ecdsa", curveP256***REMOVED***
***REMOVED***

// Algo returns the requested key algorithm represented as a string.
func (kr *BasicKeyRequest) Algo() string ***REMOVED***
	return kr.A
***REMOVED***

// Size returns the requested key size.
func (kr *BasicKeyRequest) Size() int ***REMOVED***
	return kr.S
***REMOVED***

// Generate generates a key as specified in the request. Currently,
// only ECDSA and RSA are supported.
func (kr *BasicKeyRequest) Generate() (crypto.PrivateKey, error) ***REMOVED***
	log.Debugf("generate key from request: algo=%s, size=%d", kr.Algo(), kr.Size())
	switch kr.Algo() ***REMOVED***
	case "rsa":
		if kr.Size() < 2048 ***REMOVED***
			return nil, errors.New("RSA key is too weak")
		***REMOVED***
		if kr.Size() > 8192 ***REMOVED***
			return nil, errors.New("RSA key size too large")
		***REMOVED***
		return rsa.GenerateKey(rand.Reader, kr.Size())
	case "ecdsa":
		var curve elliptic.Curve
		switch kr.Size() ***REMOVED***
		case curveP256:
			curve = elliptic.P256()
		case curveP384:
			curve = elliptic.P384()
		case curveP521:
			curve = elliptic.P521()
		default:
			return nil, errors.New("invalid curve")
		***REMOVED***
		return ecdsa.GenerateKey(curve, rand.Reader)
	default:
		return nil, errors.New("invalid algorithm")
	***REMOVED***
***REMOVED***

// SigAlgo returns an appropriate X.509 signature algorithm given the
// key request's type and size.
func (kr *BasicKeyRequest) SigAlgo() x509.SignatureAlgorithm ***REMOVED***
	switch kr.Algo() ***REMOVED***
	case "rsa":
		switch ***REMOVED***
		case kr.Size() >= 4096:
			return x509.SHA512WithRSA
		case kr.Size() >= 3072:
			return x509.SHA384WithRSA
		case kr.Size() >= 2048:
			return x509.SHA256WithRSA
		default:
			return x509.SHA1WithRSA
		***REMOVED***
	case "ecdsa":
		switch kr.Size() ***REMOVED***
		case curveP521:
			return x509.ECDSAWithSHA512
		case curveP384:
			return x509.ECDSAWithSHA384
		case curveP256:
			return x509.ECDSAWithSHA256
		default:
			return x509.ECDSAWithSHA1
		***REMOVED***
	default:
		return x509.UnknownSignatureAlgorithm
	***REMOVED***
***REMOVED***

// CAConfig is a section used in the requests initialising a new CA.
type CAConfig struct ***REMOVED***
	PathLength  int    `json:"pathlen"`
	PathLenZero bool   `json:"pathlenzero"`
	Expiry      string `json:"expiry"`
***REMOVED***

// A CertificateRequest encapsulates the API interface to the
// certificate request functionality.
type CertificateRequest struct ***REMOVED***
	CN           string
	Names        []Name     `json:"names"`
	Hosts        []string   `json:"hosts"`
	KeyRequest   KeyRequest `json:"key,omitempty"`
	CA           *CAConfig  `json:"ca,omitempty"`
	SerialNumber string     `json:"serialnumber,omitempty"`
***REMOVED***

// New returns a new, empty CertificateRequest with a
// BasicKeyRequest.
func New() *CertificateRequest ***REMOVED***
	return &CertificateRequest***REMOVED***
		KeyRequest: NewBasicKeyRequest(),
	***REMOVED***
***REMOVED***

// appendIf appends to a if s is not an empty string.
func appendIf(s string, a *[]string) ***REMOVED***
	if s != "" ***REMOVED***
		*a = append(*a, s)
	***REMOVED***
***REMOVED***

// Name returns the PKIX name for the request.
func (cr *CertificateRequest) Name() pkix.Name ***REMOVED***
	var name pkix.Name
	name.CommonName = cr.CN

	for _, n := range cr.Names ***REMOVED***
		appendIf(n.C, &name.Country)
		appendIf(n.ST, &name.Province)
		appendIf(n.L, &name.Locality)
		appendIf(n.O, &name.Organization)
		appendIf(n.OU, &name.OrganizationalUnit)
	***REMOVED***
	name.SerialNumber = cr.SerialNumber
	return name
***REMOVED***

// BasicConstraints CSR information RFC 5280, 4.2.1.9
type BasicConstraints struct ***REMOVED***
	IsCA       bool `asn1:"optional"`
	MaxPathLen int  `asn1:"optional,default:-1"`
***REMOVED***

// ParseRequest takes a certificate request and generates a key and
// CSR from it. It does no validation -- caveat emptor. It will,
// however, fail if the key request is not valid (i.e., an unsupported
// curve or RSA key size). The lack of validation was specifically
// chosen to allow the end user to define a policy and validate the
// request appropriately before calling this function.
func ParseRequest(req *CertificateRequest) (csr, key []byte, err error) ***REMOVED***
	log.Info("received CSR")
	if req.KeyRequest == nil ***REMOVED***
		req.KeyRequest = NewBasicKeyRequest()
	***REMOVED***

	log.Infof("generating key: %s-%d", req.KeyRequest.Algo(), req.KeyRequest.Size())
	priv, err := req.KeyRequest.Generate()
	if err != nil ***REMOVED***
		err = cferr.Wrap(cferr.PrivateKeyError, cferr.GenerationFailed, err)
		return
	***REMOVED***

	switch priv := priv.(type) ***REMOVED***
	case *rsa.PrivateKey:
		key = x509.MarshalPKCS1PrivateKey(priv)
		block := pem.Block***REMOVED***
			Type:  "RSA PRIVATE KEY",
			Bytes: key,
		***REMOVED***
		key = pem.EncodeToMemory(&block)
	case *ecdsa.PrivateKey:
		key, err = x509.MarshalECPrivateKey(priv)
		if err != nil ***REMOVED***
			err = cferr.Wrap(cferr.PrivateKeyError, cferr.Unknown, err)
			return
		***REMOVED***
		block := pem.Block***REMOVED***
			Type:  "EC PRIVATE KEY",
			Bytes: key,
		***REMOVED***
		key = pem.EncodeToMemory(&block)
	default:
		panic("Generate should have failed to produce a valid key.")
	***REMOVED***

	csr, err = Generate(priv.(crypto.Signer), req)
	if err != nil ***REMOVED***
		log.Errorf("failed to generate a CSR: %v", err)
		err = cferr.Wrap(cferr.CSRError, cferr.BadRequest, err)
	***REMOVED***
	return
***REMOVED***

// ExtractCertificateRequest extracts a CertificateRequest from
// x509.Certificate. It is aimed to used for generating a new certificate
// from an existing certificate. For a root certificate, the CA expiry
// length is calculated as the duration between cert.NotAfter and cert.NotBefore.
func ExtractCertificateRequest(cert *x509.Certificate) *CertificateRequest ***REMOVED***
	req := New()
	req.CN = cert.Subject.CommonName
	req.Names = getNames(cert.Subject)
	req.Hosts = getHosts(cert)
	req.SerialNumber = cert.Subject.SerialNumber

	if cert.IsCA ***REMOVED***
		req.CA = new(CAConfig)
		// CA expiry length is calculated based on the input cert
		// issue date and expiry date.
		req.CA.Expiry = cert.NotAfter.Sub(cert.NotBefore).String()
		req.CA.PathLength = cert.MaxPathLen
		req.CA.PathLenZero = cert.MaxPathLenZero
	***REMOVED***

	return req
***REMOVED***

func getHosts(cert *x509.Certificate) []string ***REMOVED***
	var hosts []string
	for _, ip := range cert.IPAddresses ***REMOVED***
		hosts = append(hosts, ip.String())
	***REMOVED***
	for _, dns := range cert.DNSNames ***REMOVED***
		hosts = append(hosts, dns)
	***REMOVED***
	for _, email := range cert.EmailAddresses ***REMOVED***
		hosts = append(hosts, email)
	***REMOVED***

	return hosts
***REMOVED***

// getNames returns an array of Names from the certificate
// It onnly cares about Country, Organization, OrganizationalUnit, Locality, Province
func getNames(sub pkix.Name) []Name ***REMOVED***
	// anonymous func for finding the max of a list of interger
	max := func(v1 int, vn ...int) (max int) ***REMOVED***
		max = v1
		for i := 0; i < len(vn); i++ ***REMOVED***
			if vn[i] > max ***REMOVED***
				max = vn[i]
			***REMOVED***
		***REMOVED***
		return max
	***REMOVED***

	nc := len(sub.Country)
	norg := len(sub.Organization)
	nou := len(sub.OrganizationalUnit)
	nl := len(sub.Locality)
	np := len(sub.Province)

	n := max(nc, norg, nou, nl, np)

	names := make([]Name, n)
	for i := range names ***REMOVED***
		if i < nc ***REMOVED***
			names[i].C = sub.Country[i]
		***REMOVED***
		if i < norg ***REMOVED***
			names[i].O = sub.Organization[i]
		***REMOVED***
		if i < nou ***REMOVED***
			names[i].OU = sub.OrganizationalUnit[i]
		***REMOVED***
		if i < nl ***REMOVED***
			names[i].L = sub.Locality[i]
		***REMOVED***
		if i < np ***REMOVED***
			names[i].ST = sub.Province[i]
		***REMOVED***
	***REMOVED***
	return names
***REMOVED***

// A Generator is responsible for validating certificate requests.
type Generator struct ***REMOVED***
	Validator func(*CertificateRequest) error
***REMOVED***

// ProcessRequest validates and processes the incoming request. It is
// a wrapper around a validator and the ParseRequest function.
func (g *Generator) ProcessRequest(req *CertificateRequest) (csr, key []byte, err error) ***REMOVED***

	log.Info("generate received request")
	err = g.Validator(req)
	if err != nil ***REMOVED***
		log.Warningf("invalid request: %v", err)
		return
	***REMOVED***

	csr, key, err = ParseRequest(req)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return
***REMOVED***

// IsNameEmpty returns true if the name has no identifying information in it.
func IsNameEmpty(n Name) bool ***REMOVED***
	empty := func(s string) bool ***REMOVED*** return strings.TrimSpace(s) == "" ***REMOVED***

	if empty(n.C) && empty(n.ST) && empty(n.L) && empty(n.O) && empty(n.OU) ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

// Regenerate uses the provided CSR as a template for signing a new
// CSR using priv.
func Regenerate(priv crypto.Signer, csr []byte) ([]byte, error) ***REMOVED***
	req, extra, err := helpers.ParseCSR(csr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED*** else if len(extra) > 0 ***REMOVED***
		return nil, errors.New("csr: trailing data in certificate request")
	***REMOVED***

	return x509.CreateCertificateRequest(rand.Reader, req, priv)
***REMOVED***

// Generate creates a new CSR from a CertificateRequest structure and
// an existing key. The KeyRequest field is ignored.
func Generate(priv crypto.Signer, req *CertificateRequest) (csr []byte, err error) ***REMOVED***
	sigAlgo := helpers.SignerAlgo(priv)
	if sigAlgo == x509.UnknownSignatureAlgorithm ***REMOVED***
		return nil, cferr.New(cferr.PrivateKeyError, cferr.Unavailable)
	***REMOVED***

	var tpl = x509.CertificateRequest***REMOVED***
		Subject:            req.Name(),
		SignatureAlgorithm: sigAlgo,
	***REMOVED***

	for i := range req.Hosts ***REMOVED***
		if ip := net.ParseIP(req.Hosts[i]); ip != nil ***REMOVED***
			tpl.IPAddresses = append(tpl.IPAddresses, ip)
		***REMOVED*** else if email, err := mail.ParseAddress(req.Hosts[i]); err == nil && email != nil ***REMOVED***
			tpl.EmailAddresses = append(tpl.EmailAddresses, email.Address)
		***REMOVED*** else ***REMOVED***
			tpl.DNSNames = append(tpl.DNSNames, req.Hosts[i])
		***REMOVED***
	***REMOVED***

	if req.CA != nil ***REMOVED***
		err = appendCAInfoToCSR(req.CA, &tpl)
		if err != nil ***REMOVED***
			err = cferr.Wrap(cferr.CSRError, cferr.GenerationFailed, err)
			return
		***REMOVED***
	***REMOVED***

	csr, err = x509.CreateCertificateRequest(rand.Reader, &tpl, priv)
	if err != nil ***REMOVED***
		log.Errorf("failed to generate a CSR: %v", err)
		err = cferr.Wrap(cferr.CSRError, cferr.BadRequest, err)
		return
	***REMOVED***
	block := pem.Block***REMOVED***
		Type:  "CERTIFICATE REQUEST",
		Bytes: csr,
	***REMOVED***

	log.Info("encoded CSR")
	csr = pem.EncodeToMemory(&block)
	return
***REMOVED***

// appendCAInfoToCSR appends CAConfig BasicConstraint extension to a CSR
func appendCAInfoToCSR(reqConf *CAConfig, csr *x509.CertificateRequest) error ***REMOVED***
	pathlen := reqConf.PathLength
	if pathlen == 0 && !reqConf.PathLenZero ***REMOVED***
		pathlen = -1
	***REMOVED***
	val, err := asn1.Marshal(BasicConstraints***REMOVED***true, pathlen***REMOVED***)

	if err != nil ***REMOVED***
		return err
	***REMOVED***

	csr.ExtraExtensions = []pkix.Extension***REMOVED***
		***REMOVED***
			Id:       asn1.ObjectIdentifier***REMOVED***2, 5, 29, 19***REMOVED***,
			Value:    val,
			Critical: true,
		***REMOVED***,
	***REMOVED***

	return nil
***REMOVED***
