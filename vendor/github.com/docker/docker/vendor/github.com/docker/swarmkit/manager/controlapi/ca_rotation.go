package controlapi

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"net/url"
	"time"

	"github.com/cloudflare/cfssl/helpers"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/ca"
	"github.com/docker/swarmkit/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var minRootExpiration = 1 * helpers.OneYear

// determines whether an api.RootCA, api.RootRotation, or api.CAConfig has a signing key (local signer)
func hasSigningKey(a interface***REMOVED******REMOVED***) bool ***REMOVED***
	switch b := a.(type) ***REMOVED***
	case *api.RootCA:
		return len(b.CAKey) > 0
	case *api.RootRotation:
		return b != nil && len(b.CAKey) > 0
	case *api.CAConfig:
		return len(b.SigningCACert) > 0 && len(b.SigningCAKey) > 0
	default:
		panic("needsExternalCAs should be called something of type *api.RootCA, *api.RootRotation, or *api.CAConfig")
	***REMOVED***
***REMOVED***

// Creates a cross-signed intermediate and new api.RootRotation object.
// This function assumes that the root cert and key and the external CAs have already been validated.
func newRootRotationObject(ctx context.Context, securityConfig *ca.SecurityConfig, apiRootCA *api.RootCA, newCARootCA ca.RootCA, extCAs []*api.ExternalCA, version uint64) (*api.RootCA, error) ***REMOVED***
	var (
		rootCert, rootKey, crossSignedCert []byte
		newRootHasSigner                   bool
		err                                error
	)

	rootCert = newCARootCA.Certs
	if s, err := newCARootCA.Signer(); err == nil ***REMOVED***
		rootCert, rootKey = s.Cert, s.Key
		newRootHasSigner = true
	***REMOVED***

	// we have to sign with the original signer, not whatever is in the SecurityConfig's RootCA (which may have an intermediate signer, if
	// a root rotation is already in progress)
	switch ***REMOVED***
	case hasSigningKey(apiRootCA):
		var oldRootCA ca.RootCA
		oldRootCA, err = ca.NewRootCA(apiRootCA.CACert, apiRootCA.CACert, apiRootCA.CAKey, ca.DefaultNodeCertExpiration, nil)
		if err == nil ***REMOVED***
			crossSignedCert, err = oldRootCA.CrossSignCACertificate(rootCert)
		***REMOVED***
	case !newRootHasSigner: // the original CA and the new CA both require external CAs
		return nil, status.Errorf(codes.InvalidArgument, "rotating from one external CA to a different external CA is not supported")
	default:
		// We need the same credentials but to connect to the original URLs (in case we are in the middle of a root rotation already)
		var urls []string
		for _, c := range extCAs ***REMOVED***
			if c.Protocol == api.ExternalCA_CAProtocolCFSSL ***REMOVED***
				urls = append(urls, c.URL)
			***REMOVED***
		***REMOVED***
		if len(urls) == 0 ***REMOVED***
			return nil, status.Errorf(codes.InvalidArgument,
				"must provide an external CA for the current external root CA to generate a cross-signed certificate")
		***REMOVED***
		rootPool := x509.NewCertPool()
		rootPool.AppendCertsFromPEM(apiRootCA.CACert)

		externalCAConfig := ca.NewExternalCATLSConfig(securityConfig.ClientTLSCreds.Config().Certificates, rootPool)
		externalCA := ca.NewExternalCA(nil, externalCAConfig, urls...)
		crossSignedCert, err = externalCA.CrossSignRootCA(ctx, newCARootCA)
	***REMOVED***

	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Error("unable to generate a cross-signed certificate for root rotation")
		return nil, status.Errorf(codes.Internal, "unable to generate a cross-signed certificate for root rotation")
	***REMOVED***

	copied := apiRootCA.Copy()
	copied.RootRotation = &api.RootRotation***REMOVED***
		CACert:            rootCert,
		CAKey:             rootKey,
		CrossSignedCACert: ca.NormalizePEMs(crossSignedCert),
	***REMOVED***
	copied.LastForcedRotation = version
	return copied, nil
***REMOVED***

// Checks that a CA URL is connectable using the credentials we have and that its server certificate is signed by the
// root CA that we expect.  This uses a TCP dialer rather than an HTTP client; because we have custom TLS configuration,
// if we wanted to use an HTTP client we'd have to create a new transport for every connection.  The docs specify that
// Transports cache connections for future re-use, which could cause many open connections.
func validateExternalCAURL(dialer *net.Dialer, tlsOpts *tls.Config, caURL string) error ***REMOVED***
	parsed, err := url.Parse(caURL)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if parsed.Scheme != "https" ***REMOVED***
		return errors.New("invalid HTTP scheme")
	***REMOVED***
	host, port, err := net.SplitHostPort(parsed.Host)
	if err != nil ***REMOVED***
		// It either has no port or is otherwise invalid (e.g. too many colons).  If it's otherwise invalid the dialer
		// will error later, so just assume it's no port and set the port to the default HTTPS port.
		host = parsed.Host
		port = "443"
	***REMOVED***

	conn, err := tls.DialWithDialer(dialer, "tcp", net.JoinHostPort(host, port), tlsOpts)
	if conn != nil ***REMOVED***
		conn.Close()
	***REMOVED***
	return err
***REMOVED***

// Validates that there is at least 1 reachable, valid external CA for the given CA certificate.  Returns true if there is, false otherwise.
// Requires that the wanted cert is already normalized.
func validateHasAtLeastOneExternalCA(ctx context.Context, externalCAs map[string][]*api.ExternalCA, securityConfig *ca.SecurityConfig,
	wantedCert []byte, desc string) ([]*api.ExternalCA, error) ***REMOVED***
	specific, ok := externalCAs[string(wantedCert)]
	if ok ***REMOVED***
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(wantedCert)
		dialer := net.Dialer***REMOVED***Timeout: 5 * time.Second***REMOVED***
		opts := tls.Config***REMOVED***
			RootCAs:      pool,
			Certificates: securityConfig.ClientTLSCreds.Config().Certificates,
		***REMOVED***
		for i, ca := range specific ***REMOVED***
			if ca.Protocol == api.ExternalCA_CAProtocolCFSSL ***REMOVED***
				if err := validateExternalCAURL(&dialer, &opts, ca.URL); err != nil ***REMOVED***
					log.G(ctx).WithError(err).Warnf("external CA # %d is unreachable or invalid", i+1)
				***REMOVED*** else ***REMOVED***
					return specific, nil
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil, status.Errorf(codes.InvalidArgument, "there must be at least one valid, reachable external CA corresponding to the %s CA certificate", desc)
***REMOVED***

// validates that the list of external CAs have valid certs associated with them, and produce a mapping of subject/pubkey:external
// for later validation of required external CAs
func getNormalizedExtCAs(caConfig *api.CAConfig, normalizedCurrentRootCACert []byte) (map[string][]*api.ExternalCA, error) ***REMOVED***
	extCAs := make(map[string][]*api.ExternalCA)

	for _, extCA := range caConfig.ExternalCAs ***REMOVED***
		associatedCert := normalizedCurrentRootCACert
		// if no associated cert is provided, assume it's the current root cert
		if len(extCA.CACert) > 0 ***REMOVED***
			associatedCert = ca.NormalizePEMs(extCA.CACert)
		***REMOVED***
		certKey := string(associatedCert)
		extCAs[certKey] = append(extCAs[certKey], extCA)
	***REMOVED***

	return extCAs, nil
***REMOVED***

// validateAndUpdateCA validates a cluster's desired CA configuration spec, and returns a RootCA value on success representing
// current RootCA as it should be.  Validation logic and return values are as follows:
// 1. Validates that the contents are complete - e.g. a signing key is not provided without a signing cert, and that external
//    CAs are not removed if they are needed.  Otherwise, returns an error.
// 2. If no desired signing cert or key are provided, then either:
//    - we are happy with the current CA configuration (force rotation value has not changed), and we return the current RootCA
//      object as is
//    - we want to generate a new internal CA cert and key (force rotation value has changed), and we return the updated RootCA
//      object
// 3. Signing cert and key have been provided: validate that these match (the cert and key match). Otherwise, return an error.
// 4. Return the updated RootCA object according to the following criteria:
//    - If the desired cert is the same as the current CA cert then abort any outstanding rotations. The current signing key
//      is replaced with the desired signing key (this could lets us switch between external->internal or internal->external
//      without an actual CA rotation, which is not needed because any leaf cert issued with one CA cert can be validated using
//       the second CA certificate).
//    - If the desired cert is the same as the current to-be-rotated-to CA cert then a new root rotation is not needed. The
//      current to-be-rotated-to signing key is replaced with the desired signing key (this could lets us switch between
//      external->internal or internal->external without an actual CA rotation, which is not needed because any leaf cert
//      issued with one CA cert can be validated using the second CA certificate).
//    - Otherwise, start a new root rotation using the desired signing cert and desired signing key as the root rotation
//      signing cert and key.  If a root rotation is already in progress, just replace it and start over.
func validateCAConfig(ctx context.Context, securityConfig *ca.SecurityConfig, cluster *api.Cluster) (*api.RootCA, error) ***REMOVED***
	newConfig := cluster.Spec.CAConfig.Copy()
	newConfig.SigningCACert = ca.NormalizePEMs(newConfig.SigningCACert) // ensure this is normalized before we use it

	if len(newConfig.SigningCAKey) > 0 && len(newConfig.SigningCACert) == 0 ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, "if a signing CA key is provided, the signing CA cert must also be provided")
	***REMOVED***

	normalizedRootCA := ca.NormalizePEMs(cluster.RootCA.CACert)
	extCAs, err := getNormalizedExtCAs(newConfig, normalizedRootCA) // validate that the list of external CAs is not malformed
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var oldCertExtCAs []*api.ExternalCA
	if !hasSigningKey(&cluster.RootCA) ***REMOVED***
		oldCertExtCAs, err = validateHasAtLeastOneExternalCA(ctx, extCAs, securityConfig, normalizedRootCA, "current")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	// if the desired CA cert and key are not set, then we are happy with the current root CA configuration, unless
	// the ForceRotate version has changed
	if len(newConfig.SigningCACert) == 0 ***REMOVED***
		if cluster.RootCA.LastForcedRotation != newConfig.ForceRotate ***REMOVED***
			newRootCA, err := ca.CreateRootCA(ca.DefaultRootCN)
			if err != nil ***REMOVED***
				return nil, status.Errorf(codes.Internal, err.Error())
			***REMOVED***
			return newRootRotationObject(ctx, securityConfig, &cluster.RootCA, newRootCA, oldCertExtCAs, newConfig.ForceRotate)
		***REMOVED***

		// we also need to make sure that if the current root rotation requires an external CA, those external CAs are
		// still valid
		if cluster.RootCA.RootRotation != nil && !hasSigningKey(cluster.RootCA.RootRotation) ***REMOVED***
			_, err := validateHasAtLeastOneExternalCA(ctx, extCAs, securityConfig, ca.NormalizePEMs(cluster.RootCA.RootRotation.CACert), "next")
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***

		return &cluster.RootCA, nil // no change, return as is
	***REMOVED***

	// A desired cert and maybe key were provided - we need to make sure the cert and key (if provided) match.
	var signingCert []byte
	if hasSigningKey(newConfig) ***REMOVED***
		signingCert = newConfig.SigningCACert
	***REMOVED***
	newRootCA, err := ca.NewRootCA(newConfig.SigningCACert, signingCert, newConfig.SigningCAKey, ca.DefaultNodeCertExpiration, nil)
	if err != nil ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	***REMOVED***

	if len(newRootCA.Pool.Subjects()) != 1 ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, "the desired CA certificate cannot contain multiple certificates")
	***REMOVED***

	parsedCert, err := helpers.ParseCertificatePEM(newConfig.SigningCACert)
	if err != nil ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, "could not parse the desired CA certificate")
	***REMOVED***

	// The new certificate's expiry must be at least one year away
	if parsedCert.NotAfter.Before(time.Now().Add(minRootExpiration)) ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, "CA certificate expires too soon")
	***REMOVED***

	if !hasSigningKey(newConfig) ***REMOVED***
		if _, err := validateHasAtLeastOneExternalCA(ctx, extCAs, securityConfig, newConfig.SigningCACert, "desired"); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	// check if we can abort any existing root rotations
	if bytes.Equal(normalizedRootCA, newConfig.SigningCACert) ***REMOVED***
		copied := cluster.RootCA.Copy()
		copied.CAKey = newConfig.SigningCAKey
		copied.RootRotation = nil
		copied.LastForcedRotation = newConfig.ForceRotate
		return copied, nil
	***REMOVED***

	// check if this is the same desired cert as an existing root rotation
	if r := cluster.RootCA.RootRotation; r != nil && bytes.Equal(ca.NormalizePEMs(r.CACert), newConfig.SigningCACert) ***REMOVED***
		copied := cluster.RootCA.Copy()
		copied.RootRotation.CAKey = newConfig.SigningCAKey
		copied.LastForcedRotation = newConfig.ForceRotate
		return copied, nil
	***REMOVED***

	// ok, everything's different; we have to begin a new root rotation which means generating a new cross-signed cert
	return newRootRotationObject(ctx, securityConfig, &cluster.RootCA, newRootCA, oldCertExtCAs, newConfig.ForceRotate)
***REMOVED***
