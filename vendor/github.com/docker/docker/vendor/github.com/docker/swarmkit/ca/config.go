package ca

import (
	cryptorand "crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"math/big"
	"math/rand"
	"path/filepath"
	"strings"
	"sync"
	"time"

	cfconfig "github.com/cloudflare/cfssl/config"
	events "github.com/docker/go-events"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/connectionbroker"
	"github.com/docker/swarmkit/identity"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/watch"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/credentials"

	"golang.org/x/net/context"
)

const (
	rootCACertFilename  = "swarm-root-ca.crt"
	rootCAKeyFilename   = "swarm-root-ca.key"
	nodeTLSCertFilename = "swarm-node.crt"
	nodeTLSKeyFilename  = "swarm-node.key"
	nodeCSRFilename     = "swarm-node.csr"

	// DefaultRootCN represents the root CN that we should create roots CAs with by default
	DefaultRootCN = "swarm-ca"
	// ManagerRole represents the Manager node type, and is used for authorization to endpoints
	ManagerRole = "swarm-manager"
	// WorkerRole represents the Worker node type, and is used for authorization to endpoints
	WorkerRole = "swarm-worker"
	// CARole represents the CA node type, and is used for clients attempting to get new certificates issued
	CARole = "swarm-ca"

	generatedSecretEntropyBytes = 16
	joinTokenBase               = 36
	// ceil(log(2^128-1, 36))
	maxGeneratedSecretLength = 25
	// ceil(log(2^256-1, 36))
	base36DigestLen = 50
)

var (
	// GetCertRetryInterval is how long to wait before retrying a node
	// certificate or root certificate request.
	GetCertRetryInterval = 2 * time.Second
)

// SecurityConfig is used to represent a node's security configuration. It includes information about
// the RootCA and ServerTLSCreds/ClientTLSCreds transport authenticators to be used for MTLS
type SecurityConfig struct ***REMOVED***
	// mu protects against concurrent access to fields inside the structure.
	mu sync.Mutex

	// renewalMu makes sure only one certificate renewal attempt happens at
	// a time. It should never be locked after mu is already locked.
	renewalMu sync.Mutex

	rootCA        *RootCA
	keyReadWriter *KeyReadWriter

	certificate *tls.Certificate
	issuerInfo  *IssuerInfo

	ServerTLSCreds *MutableTLSCreds
	ClientTLSCreds *MutableTLSCreds

	// An optional queue for anyone interested in subscribing to SecurityConfig updates
	queue *watch.Queue
***REMOVED***

// CertificateUpdate represents a change in the underlying TLS configuration being returned by
// a certificate renewal event.
type CertificateUpdate struct ***REMOVED***
	Role string
	Err  error
***REMOVED***

func validateRootCAAndTLSCert(rootCA *RootCA, tlsKeyPair *tls.Certificate) error ***REMOVED***
	var (
		leafCert         *x509.Certificate
		intermediatePool *x509.CertPool
	)
	for i, derBytes := range tlsKeyPair.Certificate ***REMOVED***
		parsed, err := x509.ParseCertificate(derBytes)
		if err != nil ***REMOVED***
			return errors.Wrap(err, "could not validate new root certificates due to parse error")
		***REMOVED***
		if i == 0 ***REMOVED***
			leafCert = parsed
		***REMOVED*** else ***REMOVED***
			if intermediatePool == nil ***REMOVED***
				intermediatePool = x509.NewCertPool()
			***REMOVED***
			intermediatePool.AddCert(parsed)
		***REMOVED***
	***REMOVED***
	opts := x509.VerifyOptions***REMOVED***
		Roots:         rootCA.Pool,
		Intermediates: intermediatePool,
	***REMOVED***
	if _, err := leafCert.Verify(opts); err != nil ***REMOVED***
		return errors.Wrap(err, "new root CA does not match existing TLS credentials")
	***REMOVED***
	return nil
***REMOVED***

// NewSecurityConfig initializes and returns a new SecurityConfig.
func NewSecurityConfig(rootCA *RootCA, krw *KeyReadWriter, tlsKeyPair *tls.Certificate, issuerInfo *IssuerInfo) (*SecurityConfig, func() error, error) ***REMOVED***
	// Create the Server TLS Credentials for this node. These will not be used by workers.
	serverTLSCreds, err := rootCA.NewServerTLSCredentials(tlsKeyPair)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	// Create a TLSConfig to be used when this node connects as a client to another remote node.
	// We're using ManagerRole as remote serverName for TLS host verification because both workers
	// and managers always connect to remote managers.
	clientTLSCreds, err := rootCA.NewClientTLSCredentials(tlsKeyPair, ManagerRole)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	q := watch.NewQueue()
	return &SecurityConfig***REMOVED***
		rootCA:        rootCA,
		keyReadWriter: krw,

		certificate: tlsKeyPair,
		issuerInfo:  issuerInfo,
		queue:       q,

		ClientTLSCreds: clientTLSCreds,
		ServerTLSCreds: serverTLSCreds,
	***REMOVED***, q.Close, nil
***REMOVED***

// RootCA returns the root CA.
func (s *SecurityConfig) RootCA() *RootCA ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.rootCA
***REMOVED***

// KeyWriter returns the object that can write keys to disk
func (s *SecurityConfig) KeyWriter() KeyWriter ***REMOVED***
	return s.keyReadWriter
***REMOVED***

// KeyReader returns the object that can read keys from disk
func (s *SecurityConfig) KeyReader() KeyReader ***REMOVED***
	return s.keyReadWriter
***REMOVED***

// UpdateRootCA replaces the root CA with a new root CA
func (s *SecurityConfig) UpdateRootCA(rootCA *RootCA) error ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()

	// refuse to update the root CA if the current TLS credentials do not validate against it
	if err := validateRootCAAndTLSCert(rootCA, s.certificate); err != nil ***REMOVED***
		return err
	***REMOVED***

	s.rootCA = rootCA
	return s.updateTLSCredentials(s.certificate, s.issuerInfo)
***REMOVED***

// Watch allows you to set a watch on the security config, in order to be notified of any changes
func (s *SecurityConfig) Watch() (chan events.Event, func()) ***REMOVED***
	return s.queue.Watch()
***REMOVED***

// IssuerInfo returns the issuer subject and issuer public key
func (s *SecurityConfig) IssuerInfo() *IssuerInfo ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.issuerInfo
***REMOVED***

// This function expects something else to have taken out a lock on the SecurityConfig.
func (s *SecurityConfig) updateTLSCredentials(certificate *tls.Certificate, issuerInfo *IssuerInfo) error ***REMOVED***
	certs := []tls.Certificate***REMOVED****certificate***REMOVED***
	clientConfig, err := NewClientTLSConfig(certs, s.rootCA.Pool, ManagerRole)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed to create a new client config using the new root CA")
	***REMOVED***

	serverConfig, err := NewServerTLSConfig(certs, s.rootCA.Pool)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed to create a new server config using the new root CA")
	***REMOVED***

	if err := s.ClientTLSCreds.loadNewTLSConfig(clientConfig); err != nil ***REMOVED***
		return errors.Wrap(err, "failed to update the client credentials")
	***REMOVED***

	if err := s.ServerTLSCreds.loadNewTLSConfig(serverConfig); err != nil ***REMOVED***
		return errors.Wrap(err, "failed to update the server TLS credentials")
	***REMOVED***

	s.certificate = certificate
	s.issuerInfo = issuerInfo
	if s.queue != nil ***REMOVED***
		s.queue.Publish(&api.NodeTLSInfo***REMOVED***
			TrustRoot:           s.rootCA.Certs,
			CertIssuerPublicKey: s.issuerInfo.PublicKey,
			CertIssuerSubject:   s.issuerInfo.Subject,
		***REMOVED***)
	***REMOVED***
	return nil
***REMOVED***

// UpdateTLSCredentials updates the security config with an updated TLS certificate and issuer info
func (s *SecurityConfig) UpdateTLSCredentials(certificate *tls.Certificate, issuerInfo *IssuerInfo) error ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateTLSCredentials(certificate, issuerInfo)
***REMOVED***

// SigningPolicy creates a policy used by the signer to ensure that the only fields
// from the remote CSRs we trust are: PublicKey, PublicKeyAlgorithm and SignatureAlgorithm.
// It receives the duration a certificate will be valid for
func SigningPolicy(certExpiry time.Duration) *cfconfig.Signing ***REMOVED***
	// Force the minimum Certificate expiration to be fifteen minutes
	if certExpiry < MinNodeCertExpiration ***REMOVED***
		certExpiry = DefaultNodeCertExpiration
	***REMOVED***

	// Add the backdate
	certExpiry = certExpiry + CertBackdate

	return &cfconfig.Signing***REMOVED***
		Default: &cfconfig.SigningProfile***REMOVED***
			Usage:    []string***REMOVED***"signing", "key encipherment", "server auth", "client auth"***REMOVED***,
			Expiry:   certExpiry,
			Backdate: CertBackdate,
			// Only trust the key components from the CSR. Everything else should
			// come directly from API call params.
			CSRWhitelist: &cfconfig.CSRWhitelist***REMOVED***
				PublicKey:          true,
				PublicKeyAlgorithm: true,
				SignatureAlgorithm: true,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// SecurityConfigPaths is used as a helper to hold all the paths of security relevant files
type SecurityConfigPaths struct ***REMOVED***
	Node, RootCA CertPaths
***REMOVED***

// NewConfigPaths returns the absolute paths to all of the different types of files
func NewConfigPaths(baseCertDir string) *SecurityConfigPaths ***REMOVED***
	return &SecurityConfigPaths***REMOVED***
		Node: CertPaths***REMOVED***
			Cert: filepath.Join(baseCertDir, nodeTLSCertFilename),
			Key:  filepath.Join(baseCertDir, nodeTLSKeyFilename)***REMOVED***,
		RootCA: CertPaths***REMOVED***
			Cert: filepath.Join(baseCertDir, rootCACertFilename),
			Key:  filepath.Join(baseCertDir, rootCAKeyFilename)***REMOVED***,
	***REMOVED***
***REMOVED***

// GenerateJoinToken creates a new join token.
func GenerateJoinToken(rootCA *RootCA) string ***REMOVED***
	var secretBytes [generatedSecretEntropyBytes]byte

	if _, err := cryptorand.Read(secretBytes[:]); err != nil ***REMOVED***
		panic(fmt.Errorf("failed to read random bytes: %v", err))
	***REMOVED***

	var nn, digest big.Int
	nn.SetBytes(secretBytes[:])
	digest.SetString(rootCA.Digest.Hex(), 16)
	return fmt.Sprintf("SWMTKN-1-%0[1]*s-%0[3]*s", base36DigestLen, digest.Text(joinTokenBase), maxGeneratedSecretLength, nn.Text(joinTokenBase))
***REMOVED***

func getCAHashFromToken(token string) (digest.Digest, error) ***REMOVED***
	split := strings.Split(token, "-")
	if len(split) != 4 || split[0] != "SWMTKN" || split[1] != "1" || len(split[2]) != base36DigestLen || len(split[3]) != maxGeneratedSecretLength ***REMOVED***
		return "", errors.New("invalid join token")
	***REMOVED***

	var digestInt big.Int
	digestInt.SetString(split[2], joinTokenBase)

	return digest.Parse(fmt.Sprintf("sha256:%0[1]*s", 64, digestInt.Text(16)))
***REMOVED***

// DownloadRootCA tries to retrieve a remote root CA and matches the digest against the provided token.
func DownloadRootCA(ctx context.Context, paths CertPaths, token string, connBroker *connectionbroker.Broker) (RootCA, error) ***REMOVED***
	var rootCA RootCA
	// Get a digest for the optional CA hash string that we've been provided
	// If we were provided a non-empty string, and it is an invalid hash, return
	// otherwise, allow the invalid digest through.
	var (
		d   digest.Digest
		err error
	)
	if token != "" ***REMOVED***
		d, err = getCAHashFromToken(token)
		if err != nil ***REMOVED***
			return RootCA***REMOVED******REMOVED***, err
		***REMOVED***
	***REMOVED***
	// Get the remote CA certificate, verify integrity with the
	// hash provided. Retry up to 5 times, in case the manager we
	// first try to contact is not responding properly (it may have
	// just been demoted, for example).

	for i := 0; i != 5; i++ ***REMOVED***
		rootCA, err = GetRemoteCA(ctx, d, connBroker)
		if err == nil ***REMOVED***
			break
		***REMOVED***
		log.G(ctx).WithError(err).Errorf("failed to retrieve remote root CA certificate")

		select ***REMOVED***
		case <-time.After(GetCertRetryInterval):
		case <-ctx.Done():
			return RootCA***REMOVED******REMOVED***, ctx.Err()
		***REMOVED***
	***REMOVED***
	if err != nil ***REMOVED***
		return RootCA***REMOVED******REMOVED***, err
	***REMOVED***

	// Save root CA certificate to disk
	if err = SaveRootCA(rootCA, paths); err != nil ***REMOVED***
		return RootCA***REMOVED******REMOVED***, err
	***REMOVED***

	log.G(ctx).Debugf("retrieved remote CA certificate: %s", paths.Cert)
	return rootCA, nil
***REMOVED***

// LoadSecurityConfig loads TLS credentials from disk, or returns an error if
// these credentials do not exist or are unusable.
func LoadSecurityConfig(ctx context.Context, rootCA RootCA, krw *KeyReadWriter, allowExpired bool) (*SecurityConfig, func() error, error) ***REMOVED***
	ctx = log.WithModule(ctx, "tls")

	// At this point we've successfully loaded the CA details from disk, or
	// successfully downloaded them remotely. The next step is to try to
	// load our certificates.

	// Read both the Cert and Key from disk
	cert, key, err := krw.Read()
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	// Check to see if this certificate was signed by our CA, and isn't expired
	_, chains, err := ValidateCertChain(rootCA.Pool, cert, allowExpired)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	// ValidateChain, if successful, will always return at least 1 chain containing
	// at least 2 certificates:  the leaf and the root.
	issuer := chains[0][1]

	// Now that we know this certificate is valid, create a TLS Certificate for our
	// credentials
	keyPair, err := tls.X509KeyPair(cert, key)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	secConfig, cleanup, err := NewSecurityConfig(&rootCA, krw, &keyPair, &IssuerInfo***REMOVED***
		Subject:   issuer.RawSubject,
		PublicKey: issuer.RawSubjectPublicKeyInfo,
	***REMOVED***)
	if err == nil ***REMOVED***
		log.G(ctx).WithFields(logrus.Fields***REMOVED***
			"node.id":   secConfig.ClientTLSCreds.NodeID(),
			"node.role": secConfig.ClientTLSCreds.Role(),
		***REMOVED***).Debug("loaded node credentials")
	***REMOVED***
	return secConfig, cleanup, err
***REMOVED***

// CertificateRequestConfig contains the information needed to request a
// certificate from a remote CA.
type CertificateRequestConfig struct ***REMOVED***
	// Token is the join token that authenticates us with the CA.
	Token string
	// Availability allows a user to control the current scheduling status of a node
	Availability api.NodeSpec_Availability
	// ConnBroker provides connections to CAs.
	ConnBroker *connectionbroker.Broker
	// Credentials provides transport credentials for communicating with the
	// remote server.
	Credentials credentials.TransportCredentials
	// ForceRemote specifies that only a remote (TCP) connection should
	// be used to request the certificate. This may be necessary in cases
	// where the local node is running a manager, but is in the process of
	// being demoted.
	ForceRemote bool
	// NodeCertificateStatusRequestTimeout determines how long to wait for a node
	// status RPC result.  If not provided (zero value), will default to 5 seconds.
	NodeCertificateStatusRequestTimeout time.Duration
	// RetryInterval specifies how long to delay between retries, if non-zero.
	RetryInterval time.Duration
***REMOVED***

// CreateSecurityConfig creates a new key and cert for this node, either locally
// or via a remote CA.
func (rootCA RootCA) CreateSecurityConfig(ctx context.Context, krw *KeyReadWriter, config CertificateRequestConfig) (*SecurityConfig, func() error, error) ***REMOVED***
	ctx = log.WithModule(ctx, "tls")

	// Create a new random ID for this certificate
	cn := identity.NewID()
	org := identity.NewID()

	proposedRole := ManagerRole
	tlsKeyPair, issuerInfo, err := rootCA.IssueAndSaveNewCertificates(krw, cn, proposedRole, org)
	switch errors.Cause(err) ***REMOVED***
	case ErrNoValidSigner:
		config.RetryInterval = GetCertRetryInterval
		// Request certificate issuance from a remote CA.
		// Last argument is nil because at this point we don't have any valid TLS creds
		tlsKeyPair, issuerInfo, err = rootCA.RequestAndSaveNewCertificates(ctx, krw, config)
		if err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("failed to request and save new certificate")
			return nil, nil, err
		***REMOVED***
	case nil:
		log.G(ctx).WithFields(logrus.Fields***REMOVED***
			"node.id":   cn,
			"node.role": proposedRole,
		***REMOVED***).Debug("issued new TLS certificate")
	default:
		log.G(ctx).WithFields(logrus.Fields***REMOVED***
			"node.id":   cn,
			"node.role": proposedRole,
		***REMOVED***).WithError(err).Errorf("failed to issue and save new certificate")
		return nil, nil, err
	***REMOVED***

	secConfig, cleanup, err := NewSecurityConfig(&rootCA, krw, tlsKeyPair, issuerInfo)
	if err == nil ***REMOVED***
		log.G(ctx).WithFields(logrus.Fields***REMOVED***
			"node.id":   secConfig.ClientTLSCreds.NodeID(),
			"node.role": secConfig.ClientTLSCreds.Role(),
		***REMOVED***).Debugf("new node credentials generated: %s", krw.Target())
	***REMOVED***
	return secConfig, cleanup, err
***REMOVED***

// TODO(cyli): currently we have to only update if it's a worker role - if we have a single root CA update path for
// both managers and workers, we won't need to check any more.
func updateRootThenUpdateCert(ctx context.Context, s *SecurityConfig, connBroker *connectionbroker.Broker, rootPaths CertPaths, failedCert *x509.Certificate) (*tls.Certificate, *IssuerInfo, error) ***REMOVED***
	if len(failedCert.Subject.OrganizationalUnit) == 0 || failedCert.Subject.OrganizationalUnit[0] != WorkerRole ***REMOVED***
		return nil, nil, errors.New("cannot update root CA since this is not a worker")
	***REMOVED***
	// try downloading a new root CA if it's an unknown authority issue, in case there was a root rotation completion
	// and we just didn't get the new root
	rootCA, err := GetRemoteCA(ctx, "", connBroker)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	// validate against the existing security config creds
	if err := s.UpdateRootCA(&rootCA); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	if err := SaveRootCA(rootCA, rootPaths); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return rootCA.RequestAndSaveNewCertificates(ctx, s.KeyWriter(),
		CertificateRequestConfig***REMOVED***
			ConnBroker:  connBroker,
			Credentials: s.ClientTLSCreds,
		***REMOVED***)
***REMOVED***

// RenewTLSConfigNow gets a new TLS cert and key, and updates the security config if provided.  This is similar to
// RenewTLSConfig, except while that monitors for expiry, and periodically renews, this renews once and is blocking
func RenewTLSConfigNow(ctx context.Context, s *SecurityConfig, connBroker *connectionbroker.Broker, rootPaths CertPaths) error ***REMOVED***
	s.renewalMu.Lock()
	defer s.renewalMu.Unlock()

	ctx = log.WithModule(ctx, "tls")
	log := log.G(ctx).WithFields(logrus.Fields***REMOVED***
		"node.id":   s.ClientTLSCreds.NodeID(),
		"node.role": s.ClientTLSCreds.Role(),
	***REMOVED***)

	// Let's request new certs. Renewals don't require a token.
	rootCA := s.RootCA()
	tlsKeyPair, issuerInfo, err := rootCA.RequestAndSaveNewCertificates(ctx,
		s.KeyWriter(),
		CertificateRequestConfig***REMOVED***
			ConnBroker:  connBroker,
			Credentials: s.ClientTLSCreds,
		***REMOVED***)
	if wrappedError, ok := err.(x509UnknownAuthError); ok ***REMOVED***
		var newErr error
		tlsKeyPair, issuerInfo, newErr = updateRootThenUpdateCert(ctx, s, connBroker, rootPaths, wrappedError.failedLeafCert)
		if newErr != nil ***REMOVED***
			err = wrappedError.error
		***REMOVED*** else ***REMOVED***
			err = nil
		***REMOVED***
	***REMOVED***
	if err != nil ***REMOVED***
		log.WithError(err).Errorf("failed to renew the certificate")
		return err
	***REMOVED***

	return s.UpdateTLSCredentials(tlsKeyPair, issuerInfo)
***REMOVED***

// calculateRandomExpiry returns a random duration between 50% and 80% of the
// original validity period
func calculateRandomExpiry(validFrom, validUntil time.Time) time.Duration ***REMOVED***
	duration := validUntil.Sub(validFrom)

	var randomExpiry int
	// Our lower bound of renewal will be half of the total expiration time
	minValidity := int(duration.Minutes() * CertLowerRotationRange)
	// Our upper bound of renewal will be 80% of the total expiration time
	maxValidity := int(duration.Minutes() * CertUpperRotationRange)
	// Let's select a random number of minutes between min and max, and set our retry for that
	// Using randomly selected rotation allows us to avoid certificate thundering herds.
	if maxValidity-minValidity < 1 ***REMOVED***
		randomExpiry = minValidity
	***REMOVED*** else ***REMOVED***
		randomExpiry = rand.Intn(maxValidity-minValidity) + int(minValidity)
	***REMOVED***

	expiry := validFrom.Add(time.Duration(randomExpiry) * time.Minute).Sub(time.Now())
	if expiry < 0 ***REMOVED***
		return 0
	***REMOVED***
	return expiry
***REMOVED***

// NewServerTLSConfig returns a tls.Config configured for a TLS Server, given a tls.Certificate
// and the PEM-encoded root CA Certificate
func NewServerTLSConfig(certs []tls.Certificate, rootCAPool *x509.CertPool) (*tls.Config, error) ***REMOVED***
	if rootCAPool == nil ***REMOVED***
		return nil, errors.New("valid root CA pool required")
	***REMOVED***

	return &tls.Config***REMOVED***
		Certificates: certs,
		// Since we're using the same CA server to issue Certificates to new nodes, we can't
		// use tls.RequireAndVerifyClientCert
		ClientAuth:               tls.VerifyClientCertIfGiven,
		RootCAs:                  rootCAPool,
		ClientCAs:                rootCAPool,
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
	***REMOVED***, nil
***REMOVED***

// NewClientTLSConfig returns a tls.Config configured for a TLS Client, given a tls.Certificate
// the PEM-encoded root CA Certificate, and the name of the remote server the client wants to connect to.
func NewClientTLSConfig(certs []tls.Certificate, rootCAPool *x509.CertPool, serverName string) (*tls.Config, error) ***REMOVED***
	if rootCAPool == nil ***REMOVED***
		return nil, errors.New("valid root CA pool required")
	***REMOVED***

	return &tls.Config***REMOVED***
		ServerName:   serverName,
		Certificates: certs,
		RootCAs:      rootCAPool,
		MinVersion:   tls.VersionTLS12,
	***REMOVED***, nil
***REMOVED***

// NewClientTLSCredentials returns GRPC credentials for a TLS GRPC client, given a tls.Certificate
// a PEM-Encoded root CA Certificate, and the name of the remote server the client wants to connect to.
func (rootCA *RootCA) NewClientTLSCredentials(cert *tls.Certificate, serverName string) (*MutableTLSCreds, error) ***REMOVED***
	tlsConfig, err := NewClientTLSConfig([]tls.Certificate***REMOVED****cert***REMOVED***, rootCA.Pool, serverName)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	mtls, err := NewMutableTLS(tlsConfig)

	return mtls, err
***REMOVED***

// NewServerTLSCredentials returns GRPC credentials for a TLS GRPC client, given a tls.Certificate
// a PEM-Encoded root CA Certificate, and the name of the remote server the client wants to connect to.
func (rootCA *RootCA) NewServerTLSCredentials(cert *tls.Certificate) (*MutableTLSCreds, error) ***REMOVED***
	tlsConfig, err := NewServerTLSConfig([]tls.Certificate***REMOVED****cert***REMOVED***, rootCA.Pool)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	mtls, err := NewMutableTLS(tlsConfig)

	return mtls, err
***REMOVED***

// ParseRole parses an apiRole into an internal role string
func ParseRole(apiRole api.NodeRole) (string, error) ***REMOVED***
	switch apiRole ***REMOVED***
	case api.NodeRoleManager:
		return ManagerRole, nil
	case api.NodeRoleWorker:
		return WorkerRole, nil
	default:
		return "", errors.Errorf("failed to parse api role: %v", apiRole)
	***REMOVED***
***REMOVED***

// FormatRole parses an internal role string into an apiRole
func FormatRole(role string) (api.NodeRole, error) ***REMOVED***
	switch strings.ToLower(role) ***REMOVED***
	case strings.ToLower(ManagerRole):
		return api.NodeRoleManager, nil
	case strings.ToLower(WorkerRole):
		return api.NodeRoleWorker, nil
	default:
		return 0, errors.Errorf("failed to parse role: %s", role)
	***REMOVED***
***REMOVED***
