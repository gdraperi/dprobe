// Package initca contains code to initialise a certificate authority,
// generating a new root key and certificate.
package initca

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"time"

	"github.com/cloudflare/cfssl/config"
	"github.com/cloudflare/cfssl/csr"
	cferr "github.com/cloudflare/cfssl/errors"
	"github.com/cloudflare/cfssl/helpers"
	"github.com/cloudflare/cfssl/log"
	"github.com/cloudflare/cfssl/signer"
	"github.com/cloudflare/cfssl/signer/local"
)

// validator contains the default validation logic for certificate
// authority certificates. The only requirement here is that the
// certificate have a non-empty subject field.
func validator(req *csr.CertificateRequest) error ***REMOVED***
	if req.CN != "" ***REMOVED***
		return nil
	***REMOVED***

	if len(req.Names) == 0 ***REMOVED***
		return cferr.Wrap(cferr.PolicyError, cferr.InvalidRequest, errors.New("missing subject information"))
	***REMOVED***

	for i := range req.Names ***REMOVED***
		if csr.IsNameEmpty(req.Names[i]) ***REMOVED***
			return cferr.Wrap(cferr.PolicyError, cferr.InvalidRequest, errors.New("missing subject information"))
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// New creates a new root certificate from the certificate request.
func New(req *csr.CertificateRequest) (cert, csrPEM, key []byte, err error) ***REMOVED***
	policy := CAPolicy()
	if req.CA != nil ***REMOVED***
		if req.CA.Expiry != "" ***REMOVED***
			policy.Default.ExpiryString = req.CA.Expiry
			policy.Default.Expiry, err = time.ParseDuration(req.CA.Expiry)
			if err != nil ***REMOVED***
				return
			***REMOVED***
		***REMOVED***

		policy.Default.CAConstraint.MaxPathLen = req.CA.PathLength
		if req.CA.PathLength != 0 && req.CA.PathLenZero == true ***REMOVED***
			log.Infof("ignore invalid 'pathlenzero' value")
		***REMOVED*** else ***REMOVED***
			policy.Default.CAConstraint.MaxPathLenZero = req.CA.PathLenZero
		***REMOVED***
	***REMOVED***

	g := &csr.Generator***REMOVED***Validator: validator***REMOVED***
	csrPEM, key, err = g.ProcessRequest(req)
	if err != nil ***REMOVED***
		log.Errorf("failed to process request: %v", err)
		key = nil
		return
	***REMOVED***

	priv, err := helpers.ParsePrivateKeyPEM(key)
	if err != nil ***REMOVED***
		log.Errorf("failed to parse private key: %v", err)
		return
	***REMOVED***

	s, err := local.NewSigner(priv, nil, signer.DefaultSigAlgo(priv), policy)
	if err != nil ***REMOVED***
		log.Errorf("failed to create signer: %v", err)
		return
	***REMOVED***

	signReq := signer.SignRequest***REMOVED***Hosts: req.Hosts, Request: string(csrPEM)***REMOVED***
	cert, err = s.Sign(signReq)

	return

***REMOVED***

// NewFromPEM creates a new root certificate from the key file passed in.
func NewFromPEM(req *csr.CertificateRequest, keyFile string) (cert, csrPEM []byte, err error) ***REMOVED***
	privData, err := ioutil.ReadFile(keyFile)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	priv, err := helpers.ParsePrivateKeyPEM(privData)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	return NewFromSigner(req, priv)
***REMOVED***

// RenewFromPEM re-creates a root certificate from the CA cert and key
// files. The resulting root certificate will have the input CA certificate
// as the template and have the same expiry length. E.g. the exsiting CA
// is valid for a year from Jan 01 2015 to Jan 01 2016, the renewed certificate
// will be valid from now and expire in one year as well.
func RenewFromPEM(caFile, keyFile string) ([]byte, error) ***REMOVED***
	caBytes, err := ioutil.ReadFile(caFile)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ca, err := helpers.ParseCertificatePEM(caBytes)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	keyBytes, err := ioutil.ReadFile(keyFile)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	key, err := helpers.ParsePrivateKeyPEM(keyBytes)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return RenewFromSigner(ca, key)

***REMOVED***

// NewFromSigner creates a new root certificate from a crypto.Signer.
func NewFromSigner(req *csr.CertificateRequest, priv crypto.Signer) (cert, csrPEM []byte, err error) ***REMOVED***
	policy := CAPolicy()
	if req.CA != nil ***REMOVED***
		if req.CA.Expiry != "" ***REMOVED***
			policy.Default.ExpiryString = req.CA.Expiry
			policy.Default.Expiry, err = time.ParseDuration(req.CA.Expiry)
			if err != nil ***REMOVED***
				return nil, nil, err
			***REMOVED***
		***REMOVED***

		policy.Default.CAConstraint.MaxPathLen = req.CA.PathLength
		if req.CA.PathLength != 0 && req.CA.PathLenZero == true ***REMOVED***
			log.Infof("ignore invalid 'pathlenzero' value")
		***REMOVED*** else ***REMOVED***
			policy.Default.CAConstraint.MaxPathLenZero = req.CA.PathLenZero
		***REMOVED***
	***REMOVED***

	csrPEM, err = csr.Generate(priv, req)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	s, err := local.NewSigner(priv, nil, signer.DefaultSigAlgo(priv), policy)
	if err != nil ***REMOVED***
		log.Errorf("failed to create signer: %v", err)
		return
	***REMOVED***

	signReq := signer.SignRequest***REMOVED***Request: string(csrPEM)***REMOVED***
	cert, err = s.Sign(signReq)
	return
***REMOVED***

// RenewFromSigner re-creates a root certificate from the CA cert and crypto.Signer.
// The resulting root certificate will have ca certificate
// as the template and have the same expiry length. E.g. the exsiting CA
// is valid for a year from Jan 01 2015 to Jan 01 2016, the renewed certificate
// will be valid from now and expire in one year as well.
func RenewFromSigner(ca *x509.Certificate, priv crypto.Signer) ([]byte, error) ***REMOVED***
	if !ca.IsCA ***REMOVED***
		return nil, errors.New("input certificate is not a CA cert")
	***REMOVED***

	// matching certificate public key vs private key
	switch ***REMOVED***
	case ca.PublicKeyAlgorithm == x509.RSA:

		var rsaPublicKey *rsa.PublicKey
		var ok bool
		if rsaPublicKey, ok = priv.Public().(*rsa.PublicKey); !ok ***REMOVED***
			return nil, cferr.New(cferr.PrivateKeyError, cferr.KeyMismatch)
		***REMOVED***
		if ca.PublicKey.(*rsa.PublicKey).N.Cmp(rsaPublicKey.N) != 0 ***REMOVED***
			return nil, cferr.New(cferr.PrivateKeyError, cferr.KeyMismatch)
		***REMOVED***
	case ca.PublicKeyAlgorithm == x509.ECDSA:
		var ecdsaPublicKey *ecdsa.PublicKey
		var ok bool
		if ecdsaPublicKey, ok = priv.Public().(*ecdsa.PublicKey); !ok ***REMOVED***
			return nil, cferr.New(cferr.PrivateKeyError, cferr.KeyMismatch)
		***REMOVED***
		if ca.PublicKey.(*ecdsa.PublicKey).X.Cmp(ecdsaPublicKey.X) != 0 ***REMOVED***
			return nil, cferr.New(cferr.PrivateKeyError, cferr.KeyMismatch)
		***REMOVED***
	default:
		return nil, cferr.New(cferr.PrivateKeyError, cferr.NotRSAOrECC)
	***REMOVED***

	req := csr.ExtractCertificateRequest(ca)

	cert, _, err := NewFromSigner(req, priv)
	return cert, err

***REMOVED***

// CAPolicy contains the CA issuing policy as default policy.
var CAPolicy = func() *config.Signing ***REMOVED***
	return &config.Signing***REMOVED***
		Default: &config.SigningProfile***REMOVED***
			Usage:        []string***REMOVED***"cert sign", "crl sign"***REMOVED***,
			ExpiryString: "43800h",
			Expiry:       5 * helpers.OneYear,
			CAConstraint: config.CAConstraint***REMOVED***IsCA: true***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***
