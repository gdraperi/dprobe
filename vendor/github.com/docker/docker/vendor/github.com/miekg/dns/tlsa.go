package dns

import (
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"io"
	"net"
	"strconv"
)

// CertificateToDANE converts a certificate to a hex string as used in the TLSA record.
func CertificateToDANE(selector, matchingType uint8, cert *x509.Certificate) (string, error) ***REMOVED***
	switch matchingType ***REMOVED***
	case 0:
		switch selector ***REMOVED***
		case 0:
			return hex.EncodeToString(cert.Raw), nil
		case 1:
			return hex.EncodeToString(cert.RawSubjectPublicKeyInfo), nil
		***REMOVED***
	case 1:
		h := sha256.New()
		switch selector ***REMOVED***
		case 0:
			io.WriteString(h, string(cert.Raw))
			return hex.EncodeToString(h.Sum(nil)), nil
		case 1:
			io.WriteString(h, string(cert.RawSubjectPublicKeyInfo))
			return hex.EncodeToString(h.Sum(nil)), nil
		***REMOVED***
	case 2:
		h := sha512.New()
		switch selector ***REMOVED***
		case 0:
			io.WriteString(h, string(cert.Raw))
			return hex.EncodeToString(h.Sum(nil)), nil
		case 1:
			io.WriteString(h, string(cert.RawSubjectPublicKeyInfo))
			return hex.EncodeToString(h.Sum(nil)), nil
		***REMOVED***
	***REMOVED***
	return "", errors.New("dns: bad TLSA MatchingType or TLSA Selector")
***REMOVED***

// Sign creates a TLSA record from an SSL certificate.
func (r *TLSA) Sign(usage, selector, matchingType int, cert *x509.Certificate) (err error) ***REMOVED***
	r.Hdr.Rrtype = TypeTLSA
	r.Usage = uint8(usage)
	r.Selector = uint8(selector)
	r.MatchingType = uint8(matchingType)

	r.Certificate, err = CertificateToDANE(r.Selector, r.MatchingType, cert)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// Verify verifies a TLSA record against an SSL certificate. If it is OK
// a nil error is returned.
func (r *TLSA) Verify(cert *x509.Certificate) error ***REMOVED***
	c, err := CertificateToDANE(r.Selector, r.MatchingType, cert)
	if err != nil ***REMOVED***
		return err // Not also ErrSig?
	***REMOVED***
	if r.Certificate == c ***REMOVED***
		return nil
	***REMOVED***
	return ErrSig // ErrSig, really?
***REMOVED***

// TLSAName returns the ownername of a TLSA resource record as per the
// rules specified in RFC 6698, Section 3.
func TLSAName(name, service, network string) (string, error) ***REMOVED***
	if !IsFqdn(name) ***REMOVED***
		return "", ErrFqdn
	***REMOVED***
	p, e := net.LookupPort(network, service)
	if e != nil ***REMOVED***
		return "", e
	***REMOVED***
	return "_" + strconv.Itoa(p) + "_" + network + "." + name, nil
***REMOVED***
