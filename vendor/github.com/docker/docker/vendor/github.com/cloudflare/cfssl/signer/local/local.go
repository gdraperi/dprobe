// Package local implements certificate signature functionality for CFSSL.
package local

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/binary"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"math/big"
	"net"
	"net/mail"
	"os"

	"github.com/cloudflare/cfssl/certdb"
	"github.com/cloudflare/cfssl/config"
	cferr "github.com/cloudflare/cfssl/errors"
	"github.com/cloudflare/cfssl/helpers"
	"github.com/cloudflare/cfssl/info"
	"github.com/cloudflare/cfssl/log"
	"github.com/cloudflare/cfssl/signer"
	"github.com/google/certificate-transparency/go"
	"github.com/google/certificate-transparency/go/client"
)

// Signer contains a signer that uses the standard library to
// support both ECDSA and RSA CA keys.
type Signer struct ***REMOVED***
	ca         *x509.Certificate
	priv       crypto.Signer
	policy     *config.Signing
	sigAlgo    x509.SignatureAlgorithm
	dbAccessor certdb.Accessor
***REMOVED***

// NewSigner creates a new Signer directly from a
// private key and certificate, with optional policy.
func NewSigner(priv crypto.Signer, cert *x509.Certificate, sigAlgo x509.SignatureAlgorithm, policy *config.Signing) (*Signer, error) ***REMOVED***
	if policy == nil ***REMOVED***
		policy = &config.Signing***REMOVED***
			Profiles: map[string]*config.SigningProfile***REMOVED******REMOVED***,
			Default:  config.DefaultConfig()***REMOVED***
	***REMOVED***

	if !policy.Valid() ***REMOVED***
		return nil, cferr.New(cferr.PolicyError, cferr.InvalidPolicy)
	***REMOVED***

	return &Signer***REMOVED***
		ca:      cert,
		priv:    priv,
		sigAlgo: sigAlgo,
		policy:  policy,
	***REMOVED***, nil
***REMOVED***

// NewSignerFromFile generates a new local signer from a caFile
// and a caKey file, both PEM encoded.
func NewSignerFromFile(caFile, caKeyFile string, policy *config.Signing) (*Signer, error) ***REMOVED***
	log.Debug("Loading CA: ", caFile)
	ca, err := ioutil.ReadFile(caFile)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	log.Debug("Loading CA key: ", caKeyFile)
	cakey, err := ioutil.ReadFile(caKeyFile)
	if err != nil ***REMOVED***
		return nil, cferr.Wrap(cferr.CertificateError, cferr.ReadFailed, err)
	***REMOVED***

	parsedCa, err := helpers.ParseCertificatePEM(ca)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	strPassword := os.Getenv("CFSSL_CA_PK_PASSWORD")
	password := []byte(strPassword)
	if strPassword == "" ***REMOVED***
		password = nil
	***REMOVED***

	priv, err := helpers.ParsePrivateKeyPEMWithPassword(cakey, password)
	if err != nil ***REMOVED***
		log.Debug("Malformed private key %v", err)
		return nil, err
	***REMOVED***

	return NewSigner(priv, parsedCa, signer.DefaultSigAlgo(priv), policy)
***REMOVED***

func (s *Signer) sign(template *x509.Certificate, profile *config.SigningProfile) (cert []byte, err error) ***REMOVED***
	var distPoints = template.CRLDistributionPoints
	err = signer.FillTemplate(template, s.policy.Default, profile)
	if distPoints != nil && len(distPoints) > 0 ***REMOVED***
		template.CRLDistributionPoints = distPoints
	***REMOVED***
	if err != nil ***REMOVED***
		return
	***REMOVED***

	var initRoot bool
	if s.ca == nil ***REMOVED***
		if !template.IsCA ***REMOVED***
			err = cferr.New(cferr.PolicyError, cferr.InvalidRequest)
			return
		***REMOVED***
		template.DNSNames = nil
		template.EmailAddresses = nil
		s.ca = template
		initRoot = true
	***REMOVED***

	derBytes, err := x509.CreateCertificate(rand.Reader, template, s.ca, template.PublicKey, s.priv)
	if err != nil ***REMOVED***
		return nil, cferr.Wrap(cferr.CertificateError, cferr.Unknown, err)
	***REMOVED***
	if initRoot ***REMOVED***
		s.ca, err = x509.ParseCertificate(derBytes)
		if err != nil ***REMOVED***
			return nil, cferr.Wrap(cferr.CertificateError, cferr.ParseFailed, err)
		***REMOVED***
	***REMOVED***

	cert = pem.EncodeToMemory(&pem.Block***REMOVED***Type: "CERTIFICATE", Bytes: derBytes***REMOVED***)
	log.Infof("signed certificate with serial number %d", template.SerialNumber)
	return
***REMOVED***

// replaceSliceIfEmpty replaces the contents of replaced with newContents if
// the slice referenced by replaced is empty
func replaceSliceIfEmpty(replaced, newContents *[]string) ***REMOVED***
	if len(*replaced) == 0 ***REMOVED***
		*replaced = *newContents
	***REMOVED***
***REMOVED***

// PopulateSubjectFromCSR has functionality similar to Name, except
// it fills the fields of the resulting pkix.Name with req's if the
// subject's corresponding fields are empty
func PopulateSubjectFromCSR(s *signer.Subject, req pkix.Name) pkix.Name ***REMOVED***
	// if no subject, use req
	if s == nil ***REMOVED***
		return req
	***REMOVED***

	name := s.Name()

	if name.CommonName == "" ***REMOVED***
		name.CommonName = req.CommonName
	***REMOVED***

	replaceSliceIfEmpty(&name.Country, &req.Country)
	replaceSliceIfEmpty(&name.Province, &req.Province)
	replaceSliceIfEmpty(&name.Locality, &req.Locality)
	replaceSliceIfEmpty(&name.Organization, &req.Organization)
	replaceSliceIfEmpty(&name.OrganizationalUnit, &req.OrganizationalUnit)
	if name.SerialNumber == "" ***REMOVED***
		name.SerialNumber = req.SerialNumber
	***REMOVED***
	return name
***REMOVED***

// OverrideHosts fills template's IPAddresses, EmailAddresses, and DNSNames with the
// content of hosts, if it is not nil.
func OverrideHosts(template *x509.Certificate, hosts []string) ***REMOVED***
	if hosts != nil ***REMOVED***
		template.IPAddresses = []net.IP***REMOVED******REMOVED***
		template.EmailAddresses = []string***REMOVED******REMOVED***
		template.DNSNames = []string***REMOVED******REMOVED***
	***REMOVED***

	for i := range hosts ***REMOVED***
		if ip := net.ParseIP(hosts[i]); ip != nil ***REMOVED***
			template.IPAddresses = append(template.IPAddresses, ip)
		***REMOVED*** else if email, err := mail.ParseAddress(hosts[i]); err == nil && email != nil ***REMOVED***
			template.EmailAddresses = append(template.EmailAddresses, email.Address)
		***REMOVED*** else ***REMOVED***
			template.DNSNames = append(template.DNSNames, hosts[i])
		***REMOVED***
	***REMOVED***

***REMOVED***

// Sign signs a new certificate based on the PEM-encoded client
// certificate or certificate request with the signing profile,
// specified by profileName.
func (s *Signer) Sign(req signer.SignRequest) (cert []byte, err error) ***REMOVED***
	profile, err := signer.Profile(s, req.Profile)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	block, _ := pem.Decode([]byte(req.Request))
	if block == nil ***REMOVED***
		return nil, cferr.New(cferr.CSRError, cferr.DecodeFailed)
	***REMOVED***

	if block.Type != "NEW CERTIFICATE REQUEST" && block.Type != "CERTIFICATE REQUEST" ***REMOVED***
		return nil, cferr.Wrap(cferr.CSRError,
			cferr.BadRequest, errors.New("not a certificate or csr"))
	***REMOVED***

	csrTemplate, err := signer.ParseCertificateRequest(s, block.Bytes)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Copy out only the fields from the CSR authorized by policy.
	safeTemplate := x509.Certificate***REMOVED******REMOVED***
	// If the profile contains no explicit whitelist, assume that all fields
	// should be copied from the CSR.
	if profile.CSRWhitelist == nil ***REMOVED***
		safeTemplate = *csrTemplate
	***REMOVED*** else ***REMOVED***
		if profile.CSRWhitelist.Subject ***REMOVED***
			safeTemplate.Subject = csrTemplate.Subject
		***REMOVED***
		if profile.CSRWhitelist.PublicKeyAlgorithm ***REMOVED***
			safeTemplate.PublicKeyAlgorithm = csrTemplate.PublicKeyAlgorithm
		***REMOVED***
		if profile.CSRWhitelist.PublicKey ***REMOVED***
			safeTemplate.PublicKey = csrTemplate.PublicKey
		***REMOVED***
		if profile.CSRWhitelist.SignatureAlgorithm ***REMOVED***
			safeTemplate.SignatureAlgorithm = csrTemplate.SignatureAlgorithm
		***REMOVED***
		if profile.CSRWhitelist.DNSNames ***REMOVED***
			safeTemplate.DNSNames = csrTemplate.DNSNames
		***REMOVED***
		if profile.CSRWhitelist.IPAddresses ***REMOVED***
			safeTemplate.IPAddresses = csrTemplate.IPAddresses
		***REMOVED***
		if profile.CSRWhitelist.EmailAddresses ***REMOVED***
			safeTemplate.EmailAddresses = csrTemplate.EmailAddresses
		***REMOVED***
	***REMOVED***

	if req.CRLOverride != "" ***REMOVED***
		safeTemplate.CRLDistributionPoints = []string***REMOVED***req.CRLOverride***REMOVED***
	***REMOVED***

	if safeTemplate.IsCA ***REMOVED***
		if !profile.CAConstraint.IsCA ***REMOVED***
			log.Error("local signer policy disallows issuing CA certificate")
			return nil, cferr.New(cferr.PolicyError, cferr.InvalidRequest)
		***REMOVED***

		if s.ca != nil && s.ca.MaxPathLen > 0 ***REMOVED***
			if safeTemplate.MaxPathLen >= s.ca.MaxPathLen ***REMOVED***
				log.Error("local signer certificate disallows CA MaxPathLen extending")
				// do not sign a cert with pathlen > current
				return nil, cferr.New(cferr.PolicyError, cferr.InvalidRequest)
			***REMOVED***
		***REMOVED*** else if s.ca != nil && s.ca.MaxPathLen == 0 && s.ca.MaxPathLenZero ***REMOVED***
			log.Error("local signer certificate disallows issuing CA certificate")
			// signer has pathlen of 0, do not sign more intermediate CAs
			return nil, cferr.New(cferr.PolicyError, cferr.InvalidRequest)
		***REMOVED***
	***REMOVED***

	OverrideHosts(&safeTemplate, req.Hosts)
	safeTemplate.Subject = PopulateSubjectFromCSR(req.Subject, safeTemplate.Subject)

	// If there is a whitelist, ensure that both the Common Name and SAN DNSNames match
	if profile.NameWhitelist != nil ***REMOVED***
		if safeTemplate.Subject.CommonName != "" ***REMOVED***
			if profile.NameWhitelist.Find([]byte(safeTemplate.Subject.CommonName)) == nil ***REMOVED***
				return nil, cferr.New(cferr.PolicyError, cferr.UnmatchedWhitelist)
			***REMOVED***
		***REMOVED***
		for _, name := range safeTemplate.DNSNames ***REMOVED***
			if profile.NameWhitelist.Find([]byte(name)) == nil ***REMOVED***
				return nil, cferr.New(cferr.PolicyError, cferr.UnmatchedWhitelist)
			***REMOVED***
		***REMOVED***
		for _, name := range safeTemplate.EmailAddresses ***REMOVED***
			if profile.NameWhitelist.Find([]byte(name)) == nil ***REMOVED***
				return nil, cferr.New(cferr.PolicyError, cferr.UnmatchedWhitelist)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if profile.ClientProvidesSerialNumbers ***REMOVED***
		if req.Serial == nil ***REMOVED***
			return nil, cferr.New(cferr.CertificateError, cferr.MissingSerial)
		***REMOVED***
		safeTemplate.SerialNumber = req.Serial
	***REMOVED*** else ***REMOVED***
		// RFC 5280 4.1.2.2:
		// Certificate users MUST be able to handle serialNumber
		// values up to 20 octets.  Conforming CAs MUST NOT use
		// serialNumber values longer than 20 octets.
		//
		// If CFSSL is providing the serial numbers, it makes
		// sense to use the max supported size.
		serialNumber := make([]byte, 20)
		_, err = io.ReadFull(rand.Reader, serialNumber)
		if err != nil ***REMOVED***
			return nil, cferr.Wrap(cferr.CertificateError, cferr.Unknown, err)
		***REMOVED***

		// SetBytes interprets buf as the bytes of a big-endian
		// unsigned integer. The leading byte should be masked
		// off to ensure it isn't negative.
		serialNumber[0] &= 0x7F

		safeTemplate.SerialNumber = new(big.Int).SetBytes(serialNumber)
	***REMOVED***

	if len(req.Extensions) > 0 ***REMOVED***
		for _, ext := range req.Extensions ***REMOVED***
			oid := asn1.ObjectIdentifier(ext.ID)
			if !profile.ExtensionWhitelist[oid.String()] ***REMOVED***
				return nil, cferr.New(cferr.CertificateError, cferr.InvalidRequest)
			***REMOVED***

			rawValue, err := hex.DecodeString(ext.Value)
			if err != nil ***REMOVED***
				return nil, cferr.Wrap(cferr.CertificateError, cferr.InvalidRequest, err)
			***REMOVED***

			safeTemplate.ExtraExtensions = append(safeTemplate.ExtraExtensions, pkix.Extension***REMOVED***
				Id:       oid,
				Critical: ext.Critical,
				Value:    rawValue,
			***REMOVED***)
		***REMOVED***
	***REMOVED***

	var certTBS = safeTemplate

	if len(profile.CTLogServers) > 0 ***REMOVED***
		// Add a poison extension which prevents validation
		var poisonExtension = pkix.Extension***REMOVED***Id: signer.CTPoisonOID, Critical: true, Value: []byte***REMOVED***0x05, 0x00***REMOVED******REMOVED***
		var poisonedPreCert = certTBS
		poisonedPreCert.ExtraExtensions = append(safeTemplate.ExtraExtensions, poisonExtension)
		cert, err = s.sign(&poisonedPreCert, profile)
		if err != nil ***REMOVED***
			return
		***REMOVED***

		derCert, _ := pem.Decode(cert)
		prechain := []ct.ASN1Cert***REMOVED***derCert.Bytes, s.ca.Raw***REMOVED***
		var sctList []ct.SignedCertificateTimestamp

		for _, server := range profile.CTLogServers ***REMOVED***
			log.Infof("submitting poisoned precertificate to %s", server)
			var ctclient = client.New(server, nil)
			var resp *ct.SignedCertificateTimestamp
			resp, err = ctclient.AddPreChain(prechain)
			if err != nil ***REMOVED***
				return nil, cferr.Wrap(cferr.CTError, cferr.PrecertSubmissionFailed, err)
			***REMOVED***
			sctList = append(sctList, *resp)
		***REMOVED***

		var serializedSCTList []byte
		serializedSCTList, err = serializeSCTList(sctList)
		if err != nil ***REMOVED***
			return nil, cferr.Wrap(cferr.CTError, cferr.Unknown, err)
		***REMOVED***

		// Serialize again as an octet string before embedding
		serializedSCTList, err = asn1.Marshal(serializedSCTList)
		if err != nil ***REMOVED***
			return nil, cferr.Wrap(cferr.CTError, cferr.Unknown, err)
		***REMOVED***

		var SCTListExtension = pkix.Extension***REMOVED***Id: signer.SCTListOID, Critical: false, Value: serializedSCTList***REMOVED***
		certTBS.ExtraExtensions = append(certTBS.ExtraExtensions, SCTListExtension)
	***REMOVED***
	var signedCert []byte
	signedCert, err = s.sign(&certTBS, profile)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if s.dbAccessor != nil ***REMOVED***
		var certRecord = certdb.CertificateRecord***REMOVED***
			Serial: certTBS.SerialNumber.String(),
			// this relies on the specific behavior of x509.CreateCertificate
			// which updates certTBS AuthorityKeyId from the signer's SubjectKeyId
			AKI:     hex.EncodeToString(certTBS.AuthorityKeyId),
			CALabel: req.Label,
			Status:  "good",
			Expiry:  certTBS.NotAfter,
			PEM:     string(signedCert),
		***REMOVED***

		err = s.dbAccessor.InsertCertificate(certRecord)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		log.Debug("saved certificate with serial number ", certTBS.SerialNumber)
	***REMOVED***

	return signedCert, nil
***REMOVED***

func serializeSCTList(sctList []ct.SignedCertificateTimestamp) ([]byte, error) ***REMOVED***
	var buf bytes.Buffer
	for _, sct := range sctList ***REMOVED***
		sct, err := ct.SerializeSCT(sct)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		binary.Write(&buf, binary.BigEndian, uint16(len(sct)))
		buf.Write(sct)
	***REMOVED***

	var sctListLengthField = make([]byte, 2)
	binary.BigEndian.PutUint16(sctListLengthField, uint16(buf.Len()))
	return bytes.Join([][]byte***REMOVED***sctListLengthField, buf.Bytes()***REMOVED***, nil), nil
***REMOVED***

// Info return a populated info.Resp struct or an error.
func (s *Signer) Info(req info.Req) (resp *info.Resp, err error) ***REMOVED***
	cert, err := s.Certificate(req.Label, req.Profile)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	profile, err := signer.Profile(s, req.Profile)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	resp = new(info.Resp)
	if cert.Raw != nil ***REMOVED***
		resp.Certificate = string(bytes.TrimSpace(pem.EncodeToMemory(&pem.Block***REMOVED***Type: "CERTIFICATE", Bytes: cert.Raw***REMOVED***)))
	***REMOVED***
	resp.Usage = profile.Usage
	resp.ExpiryString = profile.ExpiryString

	return
***REMOVED***

// SigAlgo returns the RSA signer's signature algorithm.
func (s *Signer) SigAlgo() x509.SignatureAlgorithm ***REMOVED***
	return s.sigAlgo
***REMOVED***

// Certificate returns the signer's certificate.
func (s *Signer) Certificate(label, profile string) (*x509.Certificate, error) ***REMOVED***
	cert := *s.ca
	return &cert, nil
***REMOVED***

// SetPolicy sets the signer's signature policy.
func (s *Signer) SetPolicy(policy *config.Signing) ***REMOVED***
	s.policy = policy
***REMOVED***

// SetDBAccessor sets the signers' cert db accessor
func (s *Signer) SetDBAccessor(dba certdb.Accessor) ***REMOVED***
	s.dbAccessor = dba
***REMOVED***

// Policy returns the signer's policy.
func (s *Signer) Policy() *config.Signing ***REMOVED***
	return s.policy
***REMOVED***
