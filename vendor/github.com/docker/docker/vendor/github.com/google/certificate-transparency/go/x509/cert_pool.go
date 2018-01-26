// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x509

import (
	"encoding/pem"
)

// CertPool is a set of certificates.
type CertPool struct ***REMOVED***
	bySubjectKeyId map[string][]int
	byName         map[string][]int
	certs          []*Certificate
***REMOVED***

// NewCertPool returns a new, empty CertPool.
func NewCertPool() *CertPool ***REMOVED***
	return &CertPool***REMOVED***
		make(map[string][]int),
		make(map[string][]int),
		nil,
	***REMOVED***
***REMOVED***

// findVerifiedParents attempts to find certificates in s which have signed the
// given certificate. If any candidates were rejected then errCert will be set
// to one of them, arbitrarily, and err will contain the reason that it was
// rejected.
func (s *CertPool) findVerifiedParents(cert *Certificate) (parents []int, errCert *Certificate, err error) ***REMOVED***
	if s == nil ***REMOVED***
		return
	***REMOVED***
	var candidates []int

	if len(cert.AuthorityKeyId) > 0 ***REMOVED***
		candidates = s.bySubjectKeyId[string(cert.AuthorityKeyId)]
	***REMOVED***
	if len(candidates) == 0 ***REMOVED***
		candidates = s.byName[string(cert.RawIssuer)]
	***REMOVED***

	for _, c := range candidates ***REMOVED***
		if err = cert.CheckSignatureFrom(s.certs[c]); err == nil ***REMOVED***
			parents = append(parents, c)
		***REMOVED*** else ***REMOVED***
			errCert = s.certs[c]
		***REMOVED***
	***REMOVED***

	return
***REMOVED***

// AddCert adds a certificate to a pool.
func (s *CertPool) AddCert(cert *Certificate) ***REMOVED***
	if cert == nil ***REMOVED***
		panic("adding nil Certificate to CertPool")
	***REMOVED***

	// Check that the certificate isn't being added twice.
	for _, c := range s.certs ***REMOVED***
		if c.Equal(cert) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	n := len(s.certs)
	s.certs = append(s.certs, cert)

	if len(cert.SubjectKeyId) > 0 ***REMOVED***
		keyId := string(cert.SubjectKeyId)
		s.bySubjectKeyId[keyId] = append(s.bySubjectKeyId[keyId], n)
	***REMOVED***
	name := string(cert.RawSubject)
	s.byName[name] = append(s.byName[name], n)
***REMOVED***

// AppendCertsFromPEM attempts to parse a series of PEM encoded certificates.
// It appends any certificates found to s and returns true if any certificates
// were successfully parsed.
//
// On many Linux systems, /etc/ssl/cert.pem will contain the system wide set
// of root CAs in a format suitable for this function.
func (s *CertPool) AppendCertsFromPEM(pemCerts []byte) (ok bool) ***REMOVED***
	for len(pemCerts) > 0 ***REMOVED***
		var block *pem.Block
		block, pemCerts = pem.Decode(pemCerts)
		if block == nil ***REMOVED***
			break
		***REMOVED***
		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 ***REMOVED***
			continue
		***REMOVED***

		cert, err := ParseCertificate(block.Bytes)
		if err != nil ***REMOVED***
			continue
		***REMOVED***

		s.AddCert(cert)
		ok = true
	***REMOVED***

	return
***REMOVED***

// Subjects returns a list of the DER-encoded subjects of
// all of the certificates in the pool.
func (s *CertPool) Subjects() (res [][]byte) ***REMOVED***
	res = make([][]byte, len(s.certs))
	for i, c := range s.certs ***REMOVED***
		res[i] = c.RawSubject
	***REMOVED***
	return
***REMOVED***
