package libtrust

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"time"
)

type certTemplateInfo struct ***REMOVED***
	commonName  string
	domains     []string
	ipAddresses []net.IP
	isCA        bool
	clientAuth  bool
	serverAuth  bool
***REMOVED***

func generateCertTemplate(info *certTemplateInfo) *x509.Certificate ***REMOVED***
	// Generate a certificate template which is valid from the past week to
	// 10 years from now. The usage of the certificate depends on the
	// specified fields in the given certTempInfo object.
	var (
		keyUsage    x509.KeyUsage
		extKeyUsage []x509.ExtKeyUsage
	)

	if info.isCA ***REMOVED***
		keyUsage = x509.KeyUsageCertSign
	***REMOVED***

	if info.clientAuth ***REMOVED***
		extKeyUsage = append(extKeyUsage, x509.ExtKeyUsageClientAuth)
	***REMOVED***

	if info.serverAuth ***REMOVED***
		extKeyUsage = append(extKeyUsage, x509.ExtKeyUsageServerAuth)
	***REMOVED***

	return &x509.Certificate***REMOVED***
		SerialNumber: big.NewInt(0),
		Subject: pkix.Name***REMOVED***
			CommonName: info.commonName,
		***REMOVED***,
		NotBefore:             time.Now().Add(-time.Hour * 24 * 7),
		NotAfter:              time.Now().Add(time.Hour * 24 * 365 * 10),
		DNSNames:              info.domains,
		IPAddresses:           info.ipAddresses,
		IsCA:                  info.isCA,
		KeyUsage:              keyUsage,
		ExtKeyUsage:           extKeyUsage,
		BasicConstraintsValid: info.isCA,
	***REMOVED***
***REMOVED***

func generateCert(pub PublicKey, priv PrivateKey, subInfo, issInfo *certTemplateInfo) (cert *x509.Certificate, err error) ***REMOVED***
	pubCertTemplate := generateCertTemplate(subInfo)
	privCertTemplate := generateCertTemplate(issInfo)

	certDER, err := x509.CreateCertificate(
		rand.Reader, pubCertTemplate, privCertTemplate,
		pub.CryptoPublicKey(), priv.CryptoPrivateKey(),
	)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to create certificate: %s", err)
	***REMOVED***

	cert, err = x509.ParseCertificate(certDER)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to parse certificate: %s", err)
	***REMOVED***

	return
***REMOVED***

// GenerateSelfSignedServerCert creates a self-signed certificate for the
// given key which is to be used for TLS servers with the given domains and
// IP addresses.
func GenerateSelfSignedServerCert(key PrivateKey, domains []string, ipAddresses []net.IP) (*x509.Certificate, error) ***REMOVED***
	info := &certTemplateInfo***REMOVED***
		commonName:  key.KeyID(),
		domains:     domains,
		ipAddresses: ipAddresses,
		serverAuth:  true,
	***REMOVED***

	return generateCert(key.PublicKey(), key, info, info)
***REMOVED***

// GenerateSelfSignedClientCert creates a self-signed certificate for the
// given key which is to be used for TLS clients.
func GenerateSelfSignedClientCert(key PrivateKey) (*x509.Certificate, error) ***REMOVED***
	info := &certTemplateInfo***REMOVED***
		commonName: key.KeyID(),
		clientAuth: true,
	***REMOVED***

	return generateCert(key.PublicKey(), key, info, info)
***REMOVED***

// GenerateCACert creates a certificate which can be used as a trusted
// certificate authority.
func GenerateCACert(signer PrivateKey, trustedKey PublicKey) (*x509.Certificate, error) ***REMOVED***
	subjectInfo := &certTemplateInfo***REMOVED***
		commonName: trustedKey.KeyID(),
		isCA:       true,
	***REMOVED***
	issuerInfo := &certTemplateInfo***REMOVED***
		commonName: signer.KeyID(),
	***REMOVED***

	return generateCert(trustedKey, signer, subjectInfo, issuerInfo)
***REMOVED***

// GenerateCACertPool creates a certificate authority pool to be used for a
// TLS configuration. Any self-signed certificates issued by the specified
// trusted keys will be verified during a TLS handshake
func GenerateCACertPool(signer PrivateKey, trustedKeys []PublicKey) (*x509.CertPool, error) ***REMOVED***
	certPool := x509.NewCertPool()

	for _, trustedKey := range trustedKeys ***REMOVED***
		cert, err := GenerateCACert(signer, trustedKey)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("failed to generate CA certificate: %s", err)
		***REMOVED***

		certPool.AddCert(cert)
	***REMOVED***

	return certPool, nil
***REMOVED***

// LoadCertificateBundle loads certificates from the given file.  The file should be pem encoded
// containing one or more certificates.  The expected pem type is "CERTIFICATE".
func LoadCertificateBundle(filename string) ([]*x509.Certificate, error) ***REMOVED***
	b, err := ioutil.ReadFile(filename)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	certificates := []*x509.Certificate***REMOVED******REMOVED***
	var block *pem.Block
	block, b = pem.Decode(b)
	for ; block != nil; block, b = pem.Decode(b) ***REMOVED***
		if block.Type == "CERTIFICATE" ***REMOVED***
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			certificates = append(certificates, cert)
		***REMOVED*** else ***REMOVED***
			return nil, fmt.Errorf("invalid pem block type: %s", block.Type)
		***REMOVED***
	***REMOVED***

	return certificates, nil
***REMOVED***

// LoadCertificatePool loads a CA pool from the given file.  The file should be pem encoded
// containing one or more certificates. The expected pem type is "CERTIFICATE".
func LoadCertificatePool(filename string) (*x509.CertPool, error) ***REMOVED***
	certs, err := LoadCertificateBundle(filename)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	pool := x509.NewCertPool()
	for _, cert := range certs ***REMOVED***
		pool.AddCert(cert)
	***REMOVED***
	return pool, nil
***REMOVED***
