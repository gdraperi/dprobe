// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x509

import (
	"fmt"
	"net"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"
)

type InvalidReason int

const (
	// NotAuthorizedToSign results when a certificate is signed by another
	// which isn't marked as a CA certificate.
	NotAuthorizedToSign InvalidReason = iota
	// Expired results when a certificate has expired, based on the time
	// given in the VerifyOptions.
	Expired
	// CANotAuthorizedForThisName results when an intermediate or root
	// certificate has a name constraint which doesn't include the name
	// being checked.
	CANotAuthorizedForThisName
	// TooManyIntermediates results when a path length constraint is
	// violated.
	TooManyIntermediates
	// IncompatibleUsage results when the certificate's key usage indicates
	// that it may only be used for a different purpose.
	IncompatibleUsage
)

// CertificateInvalidError results when an odd error occurs. Users of this
// library probably want to handle all these errors uniformly.
type CertificateInvalidError struct ***REMOVED***
	Cert   *Certificate
	Reason InvalidReason
***REMOVED***

func (e CertificateInvalidError) Error() string ***REMOVED***
	switch e.Reason ***REMOVED***
	case NotAuthorizedToSign:
		return "x509: certificate is not authorized to sign other certificates"
	case Expired:
		return "x509: certificate has expired or is not yet valid"
	case CANotAuthorizedForThisName:
		return "x509: a root or intermediate certificate is not authorized to sign in this domain"
	case TooManyIntermediates:
		return "x509: too many intermediates for path length constraint"
	case IncompatibleUsage:
		return "x509: certificate specifies an incompatible key usage"
	***REMOVED***
	return "x509: unknown error"
***REMOVED***

// HostnameError results when the set of authorized names doesn't match the
// requested name.
type HostnameError struct ***REMOVED***
	Certificate *Certificate
	Host        string
***REMOVED***

func (h HostnameError) Error() string ***REMOVED***
	c := h.Certificate

	var valid string
	if ip := net.ParseIP(h.Host); ip != nil ***REMOVED***
		// Trying to validate an IP
		if len(c.IPAddresses) == 0 ***REMOVED***
			return "x509: cannot validate certificate for " + h.Host + " because it doesn't contain any IP SANs"
		***REMOVED***
		for _, san := range c.IPAddresses ***REMOVED***
			if len(valid) > 0 ***REMOVED***
				valid += ", "
			***REMOVED***
			valid += san.String()
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if len(c.DNSNames) > 0 ***REMOVED***
			valid = strings.Join(c.DNSNames, ", ")
		***REMOVED*** else ***REMOVED***
			valid = c.Subject.CommonName
		***REMOVED***
	***REMOVED***
	return "x509: certificate is valid for " + valid + ", not " + h.Host
***REMOVED***

// UnknownAuthorityError results when the certificate issuer is unknown
type UnknownAuthorityError struct ***REMOVED***
	cert *Certificate
	// hintErr contains an error that may be helpful in determining why an
	// authority wasn't found.
	hintErr error
	// hintCert contains a possible authority certificate that was rejected
	// because of the error in hintErr.
	hintCert *Certificate
***REMOVED***

func (e UnknownAuthorityError) Error() string ***REMOVED***
	s := "x509: certificate signed by unknown authority"
	if e.hintErr != nil ***REMOVED***
		certName := e.hintCert.Subject.CommonName
		if len(certName) == 0 ***REMOVED***
			if len(e.hintCert.Subject.Organization) > 0 ***REMOVED***
				certName = e.hintCert.Subject.Organization[0]
			***REMOVED***
			certName = "serial:" + e.hintCert.SerialNumber.String()
		***REMOVED***
		s += fmt.Sprintf(" (possibly because of %q while trying to verify candidate authority certificate %q)", e.hintErr, certName)
	***REMOVED***
	return s
***REMOVED***

// SystemRootsError results when we fail to load the system root certificates.
type SystemRootsError struct ***REMOVED***
***REMOVED***

func (e SystemRootsError) Error() string ***REMOVED***
	return "x509: failed to load system roots and no roots provided"
***REMOVED***

// VerifyOptions contains parameters for Certificate.Verify. It's a structure
// because other PKIX verification APIs have ended up needing many options.
type VerifyOptions struct ***REMOVED***
	DNSName           string
	Intermediates     *CertPool
	Roots             *CertPool // if nil, the system roots are used
	CurrentTime       time.Time // if zero, the current time is used
	DisableTimeChecks bool
	// KeyUsage specifies which Extended Key Usage values are acceptable.
	// An empty list means ExtKeyUsageServerAuth. Key usage is considered a
	// constraint down the chain which mirrors Windows CryptoAPI behaviour,
	// but not the spec. To accept any key usage, include ExtKeyUsageAny.
	KeyUsages []ExtKeyUsage
***REMOVED***

const (
	leafCertificate = iota
	intermediateCertificate
	rootCertificate
)

// isValid performs validity checks on the c.
func (c *Certificate) isValid(certType int, currentChain []*Certificate, opts *VerifyOptions) error ***REMOVED***
	if !opts.DisableTimeChecks ***REMOVED***
		now := opts.CurrentTime
		if now.IsZero() ***REMOVED***
			now = time.Now()
		***REMOVED***
		if now.Before(c.NotBefore) || now.After(c.NotAfter) ***REMOVED***
			return CertificateInvalidError***REMOVED***c, Expired***REMOVED***
		***REMOVED***
	***REMOVED***

	if len(c.PermittedDNSDomains) > 0 ***REMOVED***
		ok := false
		for _, domain := range c.PermittedDNSDomains ***REMOVED***
			if opts.DNSName == domain ||
				(strings.HasSuffix(opts.DNSName, domain) &&
					len(opts.DNSName) >= 1+len(domain) &&
					opts.DNSName[len(opts.DNSName)-len(domain)-1] == '.') ***REMOVED***
				ok = true
				break
			***REMOVED***
		***REMOVED***

		if !ok ***REMOVED***
			return CertificateInvalidError***REMOVED***c, CANotAuthorizedForThisName***REMOVED***
		***REMOVED***
	***REMOVED***

	// KeyUsage status flags are ignored. From Engineering Security, Peter
	// Gutmann: A European government CA marked its signing certificates as
	// being valid for encryption only, but no-one noticed. Another
	// European CA marked its signature keys as not being valid for
	// signatures. A different CA marked its own trusted root certificate
	// as being invalid for certificate signing.  Another national CA
	// distributed a certificate to be used to encrypt data for the
	// country’s tax authority that was marked as only being usable for
	// digital signatures but not for encryption. Yet another CA reversed
	// the order of the bit flags in the keyUsage due to confusion over
	// encoding endianness, essentially setting a random keyUsage in
	// certificates that it issued. Another CA created a self-invalidating
	// certificate by adding a certificate policy statement stipulating
	// that the certificate had to be used strictly as specified in the
	// keyUsage, and a keyUsage containing a flag indicating that the RSA
	// encryption key could only be used for Diffie-Hellman key agreement.

	if certType == intermediateCertificate && (!c.BasicConstraintsValid || !c.IsCA) ***REMOVED***
		return CertificateInvalidError***REMOVED***c, NotAuthorizedToSign***REMOVED***
	***REMOVED***

	if c.BasicConstraintsValid && c.MaxPathLen >= 0 ***REMOVED***
		numIntermediates := len(currentChain) - 1
		if numIntermediates > c.MaxPathLen ***REMOVED***
			return CertificateInvalidError***REMOVED***c, TooManyIntermediates***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// Verify attempts to verify c by building one or more chains from c to a
// certificate in opts.Roots, using certificates in opts.Intermediates if
// needed. If successful, it returns one or more chains where the first
// element of the chain is c and the last element is from opts.Roots.
//
// WARNING: this doesn't do any revocation checking.
func (c *Certificate) Verify(opts VerifyOptions) (chains [][]*Certificate, err error) ***REMOVED***
	// Use Windows's own verification and chain building.
	if opts.Roots == nil && runtime.GOOS == "windows" ***REMOVED***
		return c.systemVerify(&opts)
	***REMOVED***

	if opts.Roots == nil ***REMOVED***
		opts.Roots = systemRootsPool()
		if opts.Roots == nil ***REMOVED***
			return nil, SystemRootsError***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***

	err = c.isValid(leafCertificate, nil, &opts)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	if len(opts.DNSName) > 0 ***REMOVED***
		err = c.VerifyHostname(opts.DNSName)
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	candidateChains, err := c.buildChains(make(map[int][][]*Certificate), []*Certificate***REMOVED***c***REMOVED***, &opts)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	keyUsages := opts.KeyUsages
	if len(keyUsages) == 0 ***REMOVED***
		keyUsages = []ExtKeyUsage***REMOVED***ExtKeyUsageServerAuth***REMOVED***
	***REMOVED***

	// If any key usage is acceptable then we're done.
	for _, usage := range keyUsages ***REMOVED***
		if usage == ExtKeyUsageAny ***REMOVED***
			chains = candidateChains
			return
		***REMOVED***
	***REMOVED***

	for _, candidate := range candidateChains ***REMOVED***
		if checkChainForKeyUsage(candidate, keyUsages) ***REMOVED***
			chains = append(chains, candidate)
		***REMOVED***
	***REMOVED***

	if len(chains) == 0 ***REMOVED***
		err = CertificateInvalidError***REMOVED***c, IncompatibleUsage***REMOVED***
	***REMOVED***

	return
***REMOVED***

func appendToFreshChain(chain []*Certificate, cert *Certificate) []*Certificate ***REMOVED***
	n := make([]*Certificate, len(chain)+1)
	copy(n, chain)
	n[len(chain)] = cert
	return n
***REMOVED***

func (c *Certificate) buildChains(cache map[int][][]*Certificate, currentChain []*Certificate, opts *VerifyOptions) (chains [][]*Certificate, err error) ***REMOVED***
	possibleRoots, failedRoot, rootErr := opts.Roots.findVerifiedParents(c)
	for _, rootNum := range possibleRoots ***REMOVED***
		root := opts.Roots.certs[rootNum]
		err = root.isValid(rootCertificate, currentChain, opts)
		if err != nil ***REMOVED***
			continue
		***REMOVED***
		chains = append(chains, appendToFreshChain(currentChain, root))
	***REMOVED***

	possibleIntermediates, failedIntermediate, intermediateErr := opts.Intermediates.findVerifiedParents(c)
nextIntermediate:
	for _, intermediateNum := range possibleIntermediates ***REMOVED***
		intermediate := opts.Intermediates.certs[intermediateNum]
		for _, cert := range currentChain ***REMOVED***
			if cert == intermediate ***REMOVED***
				continue nextIntermediate
			***REMOVED***
		***REMOVED***
		err = intermediate.isValid(intermediateCertificate, currentChain, opts)
		if err != nil ***REMOVED***
			continue
		***REMOVED***
		var childChains [][]*Certificate
		childChains, ok := cache[intermediateNum]
		if !ok ***REMOVED***
			childChains, err = intermediate.buildChains(cache, appendToFreshChain(currentChain, intermediate), opts)
			cache[intermediateNum] = childChains
		***REMOVED***
		chains = append(chains, childChains...)
	***REMOVED***

	if len(chains) > 0 ***REMOVED***
		err = nil
	***REMOVED***

	if len(chains) == 0 && err == nil ***REMOVED***
		hintErr := rootErr
		hintCert := failedRoot
		if hintErr == nil ***REMOVED***
			hintErr = intermediateErr
			hintCert = failedIntermediate
		***REMOVED***
		err = UnknownAuthorityError***REMOVED***c, hintErr, hintCert***REMOVED***
	***REMOVED***

	return
***REMOVED***

func matchHostnames(pattern, host string) bool ***REMOVED***
	if len(pattern) == 0 || len(host) == 0 ***REMOVED***
		return false
	***REMOVED***

	patternParts := strings.Split(pattern, ".")
	hostParts := strings.Split(host, ".")

	if len(patternParts) != len(hostParts) ***REMOVED***
		return false
	***REMOVED***

	for i, patternPart := range patternParts ***REMOVED***
		if patternPart == "*" ***REMOVED***
			continue
		***REMOVED***
		if patternPart != hostParts[i] ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

// toLowerCaseASCII returns a lower-case version of in. See RFC 6125 6.4.1. We use
// an explicitly ASCII function to avoid any sharp corners resulting from
// performing Unicode operations on DNS labels.
func toLowerCaseASCII(in string) string ***REMOVED***
	// If the string is already lower-case then there's nothing to do.
	isAlreadyLowerCase := true
	for _, c := range in ***REMOVED***
		if c == utf8.RuneError ***REMOVED***
			// If we get a UTF-8 error then there might be
			// upper-case ASCII bytes in the invalid sequence.
			isAlreadyLowerCase = false
			break
		***REMOVED***
		if 'A' <= c && c <= 'Z' ***REMOVED***
			isAlreadyLowerCase = false
			break
		***REMOVED***
	***REMOVED***

	if isAlreadyLowerCase ***REMOVED***
		return in
	***REMOVED***

	out := []byte(in)
	for i, c := range out ***REMOVED***
		if 'A' <= c && c <= 'Z' ***REMOVED***
			out[i] += 'a' - 'A'
		***REMOVED***
	***REMOVED***
	return string(out)
***REMOVED***

// VerifyHostname returns nil if c is a valid certificate for the named host.
// Otherwise it returns an error describing the mismatch.
func (c *Certificate) VerifyHostname(h string) error ***REMOVED***
	// IP addresses may be written in [ ].
	candidateIP := h
	if len(h) >= 3 && h[0] == '[' && h[len(h)-1] == ']' ***REMOVED***
		candidateIP = h[1 : len(h)-1]
	***REMOVED***
	if ip := net.ParseIP(candidateIP); ip != nil ***REMOVED***
		// We only match IP addresses against IP SANs.
		// https://tools.ietf.org/html/rfc6125#appendix-B.2
		for _, candidate := range c.IPAddresses ***REMOVED***
			if ip.Equal(candidate) ***REMOVED***
				return nil
			***REMOVED***
		***REMOVED***
		return HostnameError***REMOVED***c, candidateIP***REMOVED***
	***REMOVED***

	lowered := toLowerCaseASCII(h)

	if len(c.DNSNames) > 0 ***REMOVED***
		for _, match := range c.DNSNames ***REMOVED***
			if matchHostnames(toLowerCaseASCII(match), lowered) ***REMOVED***
				return nil
			***REMOVED***
		***REMOVED***
		// If Subject Alt Name is given, we ignore the common name.
	***REMOVED*** else if matchHostnames(toLowerCaseASCII(c.Subject.CommonName), lowered) ***REMOVED***
		return nil
	***REMOVED***

	return HostnameError***REMOVED***c, h***REMOVED***
***REMOVED***

func checkChainForKeyUsage(chain []*Certificate, keyUsages []ExtKeyUsage) bool ***REMOVED***
	usages := make([]ExtKeyUsage, len(keyUsages))
	copy(usages, keyUsages)

	if len(chain) == 0 ***REMOVED***
		return false
	***REMOVED***

	usagesRemaining := len(usages)

	// We walk down the list and cross out any usages that aren't supported
	// by each certificate. If we cross out all the usages, then the chain
	// is unacceptable.

	for i := len(chain) - 1; i >= 0; i-- ***REMOVED***
		cert := chain[i]
		if len(cert.ExtKeyUsage) == 0 && len(cert.UnknownExtKeyUsage) == 0 ***REMOVED***
			// The certificate doesn't have any extended key usage specified.
			continue
		***REMOVED***

		for _, usage := range cert.ExtKeyUsage ***REMOVED***
			if usage == ExtKeyUsageAny ***REMOVED***
				// The certificate is explicitly good for any usage.
				continue
			***REMOVED***
		***REMOVED***

		const invalidUsage ExtKeyUsage = -1

	NextRequestedUsage:
		for i, requestedUsage := range usages ***REMOVED***
			if requestedUsage == invalidUsage ***REMOVED***
				continue
			***REMOVED***

			for _, usage := range cert.ExtKeyUsage ***REMOVED***
				if requestedUsage == usage ***REMOVED***
					continue NextRequestedUsage
				***REMOVED*** else if requestedUsage == ExtKeyUsageServerAuth &&
					(usage == ExtKeyUsageNetscapeServerGatedCrypto ||
						usage == ExtKeyUsageMicrosoftServerGatedCrypto) ***REMOVED***
					// In order to support COMODO
					// certificate chains, we have to
					// accept Netscape or Microsoft SGC
					// usages as equal to ServerAuth.
					continue NextRequestedUsage
				***REMOVED***
			***REMOVED***

			usages[i] = invalidUsage
			usagesRemaining--
			if usagesRemaining == 0 ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***
