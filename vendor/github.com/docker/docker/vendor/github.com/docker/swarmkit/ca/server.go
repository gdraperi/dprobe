package ca

import (
	"bytes"
	"crypto/subtle"
	"crypto/x509"
	"sync"
	"time"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/equality"
	"github.com/docker/swarmkit/identity"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/state/store"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	defaultReconciliationRetryInterval = 10 * time.Second
	defaultRootReconciliationInterval  = 3 * time.Second
)

// Server is the CA and NodeCA API gRPC server.
// TODO(aaronl): At some point we may want to have separate implementations of
// CA, NodeCA, and other hypothetical future CA services. At the moment,
// breaking it apart doesn't seem worth it.
type Server struct ***REMOVED***
	mu                          sync.Mutex
	wg                          sync.WaitGroup
	ctx                         context.Context
	cancel                      func()
	store                       *store.MemoryStore
	securityConfig              *SecurityConfig
	clusterID                   string
	localRootCA                 *RootCA
	externalCA                  *ExternalCA
	externalCAPool              *x509.CertPool
	joinTokens                  *api.JoinTokens
	reconciliationRetryInterval time.Duration

	// pending is a map of nodes with pending certificates issuance or
	// renewal. They are indexed by node ID.
	pending map[string]*api.Node

	// started is a channel which gets closed once the server is running
	// and able to service RPCs.
	started chan struct***REMOVED******REMOVED***

	// these are cached values to ensure we only update the security config when
	// the cluster root CA and external CAs have changed - the cluster object
	// can change for other reasons, and it would not be necessary to update
	// the security config as a result
	lastSeenClusterRootCA *api.RootCA
	lastSeenExternalCAs   []*api.ExternalCA

	// This mutex protects the components of the CA server used to issue new certificates
	// (and any attributes used to update those components): `lastSeenClusterRootCA` and
	// `lastSeenExternalCA`, which are used to update `externalCA` and the `rootCA` object
	// of the SecurityConfig
	signingMu sync.Mutex

	// lets us monitor and finish root rotations
	rootReconciler                  *rootRotationReconciler
	rootReconciliationRetryInterval time.Duration
***REMOVED***

// DefaultCAConfig returns the default CA Config, with a default expiration.
func DefaultCAConfig() api.CAConfig ***REMOVED***
	return api.CAConfig***REMOVED***
		NodeCertExpiry: gogotypes.DurationProto(DefaultNodeCertExpiration),
	***REMOVED***
***REMOVED***

// NewServer creates a CA API server.
func NewServer(store *store.MemoryStore, securityConfig *SecurityConfig) *Server ***REMOVED***
	return &Server***REMOVED***
		store:                           store,
		securityConfig:                  securityConfig,
		localRootCA:                     securityConfig.RootCA(),
		externalCA:                      NewExternalCA(nil, nil),
		pending:                         make(map[string]*api.Node),
		started:                         make(chan struct***REMOVED******REMOVED***),
		reconciliationRetryInterval:     defaultReconciliationRetryInterval,
		rootReconciliationRetryInterval: defaultRootReconciliationInterval,
		clusterID:                       securityConfig.ClientTLSCreds.Organization(),
	***REMOVED***
***REMOVED***

// ExternalCA returns the current external CA - this is exposed to support unit testing only, and the external CA
// should really be a private field
func (s *Server) ExternalCA() *ExternalCA ***REMOVED***
	s.signingMu.Lock()
	defer s.signingMu.Unlock()
	return s.externalCA
***REMOVED***

// RootCA returns the current local root CA - this is exposed to support unit testing only, and the root CA
// should really be a private field
func (s *Server) RootCA() *RootCA ***REMOVED***
	s.signingMu.Lock()
	defer s.signingMu.Unlock()
	return s.localRootCA
***REMOVED***

// SetReconciliationRetryInterval changes the time interval between
// reconciliation attempts. This function must be called before Run.
func (s *Server) SetReconciliationRetryInterval(reconciliationRetryInterval time.Duration) ***REMOVED***
	s.reconciliationRetryInterval = reconciliationRetryInterval
***REMOVED***

// SetRootReconciliationInterval changes the time interval between root rotation
// reconciliation attempts.  This function must be called before Run.
func (s *Server) SetRootReconciliationInterval(interval time.Duration) ***REMOVED***
	s.rootReconciliationRetryInterval = interval
***REMOVED***

// GetUnlockKey is responsible for returning the current unlock key used for encrypting TLS private keys and
// other at rest data.  Access to this RPC call should only be allowed via mutual TLS from managers.
func (s *Server) GetUnlockKey(ctx context.Context, request *api.GetUnlockKeyRequest) (*api.GetUnlockKeyResponse, error) ***REMOVED***
	// This directly queries the store, rather than storing the unlock key and version on
	// the `Server` object and updating it `updateCluster` is called, because we need this
	// API to return the latest version of the key.  Otherwise, there might be a slight delay
	// between when the cluster gets updated, and when this function returns the latest key.
	// This delay is currently unacceptable because this RPC call is the only way, after a
	// cluster update, to get the actual value of the unlock key, and we don't want to return
	// a cached value.
	resp := api.GetUnlockKeyResponse***REMOVED******REMOVED***
	s.store.View(func(tx store.ReadTx) ***REMOVED***
		cluster := store.GetCluster(tx, s.clusterID)
		resp.Version = cluster.Meta.Version
		if cluster.Spec.EncryptionConfig.AutoLockManagers ***REMOVED***
			for _, encryptionKey := range cluster.UnlockKeys ***REMOVED***
				if encryptionKey.Subsystem == ManagerRole ***REMOVED***
					resp.UnlockKey = encryptionKey.Key
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***)

	return &resp, nil
***REMOVED***

// NodeCertificateStatus returns the current issuance status of an issuance request identified by the nodeID
func (s *Server) NodeCertificateStatus(ctx context.Context, request *api.NodeCertificateStatusRequest) (*api.NodeCertificateStatusResponse, error) ***REMOVED***
	if request.NodeID == "" ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, codes.InvalidArgument.String())
	***REMOVED***

	serverCtx, err := s.isRunningLocked()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var node *api.Node

	event := api.EventUpdateNode***REMOVED***
		Node:   &api.Node***REMOVED***ID: request.NodeID***REMOVED***,
		Checks: []api.NodeCheckFunc***REMOVED***api.NodeCheckID***REMOVED***,
	***REMOVED***

	// Retrieve the current value of the certificate with this token, and create a watcher
	updates, cancel, err := store.ViewAndWatch(
		s.store,
		func(tx store.ReadTx) error ***REMOVED***
			node = store.GetNode(tx, request.NodeID)
			return nil
		***REMOVED***,
		event,
	)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer cancel()

	// This node ID doesn't exist
	if node == nil ***REMOVED***
		return nil, status.Errorf(codes.NotFound, codes.NotFound.String())
	***REMOVED***

	log.G(ctx).WithFields(logrus.Fields***REMOVED***
		"node.id": node.ID,
		"status":  node.Certificate.Status,
		"method":  "NodeCertificateStatus",
	***REMOVED***)

	// If this certificate has a final state, return it immediately (both pending and renew are transition states)
	if isFinalState(node.Certificate.Status) ***REMOVED***
		return &api.NodeCertificateStatusResponse***REMOVED***
			Status:      &node.Certificate.Status,
			Certificate: &node.Certificate,
		***REMOVED***, nil
	***REMOVED***

	log.G(ctx).WithFields(logrus.Fields***REMOVED***
		"node.id": node.ID,
		"status":  node.Certificate.Status,
		"method":  "NodeCertificateStatus",
	***REMOVED***).Debugf("started watching for certificate updates")

	// Certificate is Pending or in an Unknown state, let's wait for changes.
	for ***REMOVED***
		select ***REMOVED***
		case event := <-updates:
			switch v := event.(type) ***REMOVED***
			case api.EventUpdateNode:
				// We got an update on the certificate record. If the status is a final state,
				// return the certificate.
				if isFinalState(v.Node.Certificate.Status) ***REMOVED***
					cert := v.Node.Certificate.Copy()
					return &api.NodeCertificateStatusResponse***REMOVED***
						Status:      &cert.Status,
						Certificate: cert,
					***REMOVED***, nil
				***REMOVED***
			***REMOVED***
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-serverCtx.Done():
			return nil, s.ctx.Err()
		***REMOVED***
	***REMOVED***
***REMOVED***

// IssueNodeCertificate is responsible for gatekeeping both certificate requests from new nodes in the swarm,
// and authorizing certificate renewals.
// If a node presented a valid certificate, the corresponding certificate is set in a RENEW state.
// If a node failed to present a valid certificate, we check for a valid join token and set the
// role accordingly. A new random node ID is generated, and the corresponding node entry is created.
// IssueNodeCertificate is the only place where new node entries to raft should be created.
func (s *Server) IssueNodeCertificate(ctx context.Context, request *api.IssueNodeCertificateRequest) (*api.IssueNodeCertificateResponse, error) ***REMOVED***
	// First, let's see if the remote node is presenting a non-empty CSR
	if len(request.CSR) == 0 ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, codes.InvalidArgument.String())
	***REMOVED***

	if err := s.isReadyLocked(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var (
		blacklistedCerts map[string]*api.BlacklistedCertificate
		clusters         []*api.Cluster
		err              error
	)

	s.store.View(func(readTx store.ReadTx) ***REMOVED***
		clusters, err = store.FindClusters(readTx, store.ByName(store.DefaultClusterName))
	***REMOVED***)

	// Not having a cluster object yet means we can't check
	// the blacklist.
	if err == nil && len(clusters) == 1 ***REMOVED***
		blacklistedCerts = clusters[0].BlacklistedCertificates
	***REMOVED***

	// Renewing the cert with a local (unix socket) is always valid.
	localNodeInfo := ctx.Value(LocalRequestKey)
	if localNodeInfo != nil ***REMOVED***
		nodeInfo, ok := localNodeInfo.(RemoteNodeInfo)
		if ok && nodeInfo.NodeID != "" ***REMOVED***
			return s.issueRenewCertificate(ctx, nodeInfo.NodeID, request.CSR)
		***REMOVED***
	***REMOVED***

	// If the remote node is a worker (either forwarded by a manager, or calling directly),
	// issue a renew worker certificate entry with the correct ID
	nodeID, err := AuthorizeForwardedRoleAndOrg(ctx, []string***REMOVED***WorkerRole***REMOVED***, []string***REMOVED***ManagerRole***REMOVED***, s.clusterID, blacklistedCerts)
	if err == nil ***REMOVED***
		return s.issueRenewCertificate(ctx, nodeID, request.CSR)
	***REMOVED***

	// If the remote node is a manager (either forwarded by another manager, or calling directly),
	// issue a renew certificate entry with the correct ID
	nodeID, err = AuthorizeForwardedRoleAndOrg(ctx, []string***REMOVED***ManagerRole***REMOVED***, []string***REMOVED***ManagerRole***REMOVED***, s.clusterID, blacklistedCerts)
	if err == nil ***REMOVED***
		return s.issueRenewCertificate(ctx, nodeID, request.CSR)
	***REMOVED***

	// The remote node didn't successfully present a valid MTLS certificate, let's issue a
	// certificate with a new random ID
	role := api.NodeRole(-1)

	s.mu.Lock()
	if subtle.ConstantTimeCompare([]byte(s.joinTokens.Manager), []byte(request.Token)) == 1 ***REMOVED***
		role = api.NodeRoleManager
	***REMOVED*** else if subtle.ConstantTimeCompare([]byte(s.joinTokens.Worker), []byte(request.Token)) == 1 ***REMOVED***
		role = api.NodeRoleWorker
	***REMOVED***
	s.mu.Unlock()

	if role < 0 ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, "A valid join token is necessary to join this cluster")
	***REMOVED***

	// Max number of collisions of ID or CN to tolerate before giving up
	maxRetries := 3
	// Generate a random ID for this new node
	for i := 0; ; i++ ***REMOVED***
		nodeID = identity.NewID()

		// Create a new node
		err := s.store.Update(func(tx store.Tx) error ***REMOVED***
			node := &api.Node***REMOVED***
				Role: role,
				ID:   nodeID,
				Certificate: api.Certificate***REMOVED***
					CSR:  request.CSR,
					CN:   nodeID,
					Role: role,
					Status: api.IssuanceStatus***REMOVED***
						State: api.IssuanceStatePending,
					***REMOVED***,
				***REMOVED***,
				Spec: api.NodeSpec***REMOVED***
					DesiredRole:  role,
					Membership:   api.NodeMembershipAccepted,
					Availability: request.Availability,
				***REMOVED***,
			***REMOVED***

			return store.CreateNode(tx, node)
		***REMOVED***)
		if err == nil ***REMOVED***
			log.G(ctx).WithFields(logrus.Fields***REMOVED***
				"node.id":   nodeID,
				"node.role": role,
				"method":    "IssueNodeCertificate",
			***REMOVED***).Debugf("new certificate entry added")
			break
		***REMOVED***
		if err != store.ErrExist ***REMOVED***
			return nil, err
		***REMOVED***
		if i == maxRetries ***REMOVED***
			return nil, err
		***REMOVED***
		log.G(ctx).WithFields(logrus.Fields***REMOVED***
			"node.id":   nodeID,
			"node.role": role,
			"method":    "IssueNodeCertificate",
		***REMOVED***).Errorf("randomly generated node ID collided with an existing one - retrying")
	***REMOVED***

	return &api.IssueNodeCertificateResponse***REMOVED***
		NodeID:         nodeID,
		NodeMembership: api.NodeMembershipAccepted,
	***REMOVED***, nil
***REMOVED***

// issueRenewCertificate receives a nodeID and a CSR and modifies the node's certificate entry with the new CSR
// and changes the state to RENEW, so it can be picked up and signed by the signing reconciliation loop
func (s *Server) issueRenewCertificate(ctx context.Context, nodeID string, csr []byte) (*api.IssueNodeCertificateResponse, error) ***REMOVED***
	var (
		cert api.Certificate
		node *api.Node
	)
	err := s.store.Update(func(tx store.Tx) error ***REMOVED***
		// Attempt to retrieve the node with nodeID
		node = store.GetNode(tx, nodeID)
		if node == nil ***REMOVED***
			log.G(ctx).WithFields(logrus.Fields***REMOVED***
				"node.id": nodeID,
				"method":  "issueRenewCertificate",
			***REMOVED***).Warnf("node does not exist")
			// If this node doesn't exist, we shouldn't be renewing a certificate for it
			return status.Errorf(codes.NotFound, "node %s not found when attempting to renew certificate", nodeID)
		***REMOVED***

		// Create a new Certificate entry for this node with the new CSR and a RENEW state
		cert = api.Certificate***REMOVED***
			CSR:  csr,
			CN:   node.ID,
			Role: node.Role,
			Status: api.IssuanceStatus***REMOVED***
				State: api.IssuanceStateRenew,
			***REMOVED***,
		***REMOVED***

		node.Certificate = cert
		return store.UpdateNode(tx, node)
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	log.G(ctx).WithFields(logrus.Fields***REMOVED***
		"cert.cn":   cert.CN,
		"cert.role": cert.Role,
		"method":    "issueRenewCertificate",
	***REMOVED***).Debugf("node certificate updated")

	return &api.IssueNodeCertificateResponse***REMOVED***
		NodeID:         nodeID,
		NodeMembership: node.Spec.Membership,
	***REMOVED***, nil
***REMOVED***

// GetRootCACertificate returns the certificate of the Root CA. It is used as a convenience for distributing
// the root of trust for the swarm. Clients should be using the CA hash to verify if they weren't target to
// a MiTM. If they fail to do so, node bootstrap works with TOFU semantics.
func (s *Server) GetRootCACertificate(ctx context.Context, request *api.GetRootCACertificateRequest) (*api.GetRootCACertificateResponse, error) ***REMOVED***
	log.G(ctx).WithFields(logrus.Fields***REMOVED***
		"method": "GetRootCACertificate",
	***REMOVED***)

	s.signingMu.Lock()
	defer s.signingMu.Unlock()

	return &api.GetRootCACertificateResponse***REMOVED***
		Certificate: s.localRootCA.Certs,
	***REMOVED***, nil
***REMOVED***

// Run runs the CA signer main loop.
// The CA signer can be stopped with cancelling ctx or calling Stop().
func (s *Server) Run(ctx context.Context) error ***REMOVED***
	s.mu.Lock()
	if s.isRunning() ***REMOVED***
		s.mu.Unlock()
		return errors.New("CA signer is already running")
	***REMOVED***
	s.wg.Add(1)
	s.ctx, s.cancel = context.WithCancel(log.WithModule(ctx, "ca"))
	ctx = s.ctx
	s.mu.Unlock()
	defer s.wg.Done()
	defer func() ***REMOVED***
		s.mu.Lock()
		s.mu.Unlock()
	***REMOVED***()

	// Retrieve the channels to keep track of changes in the cluster
	// Retrieve all the currently registered nodes
	var (
		nodes   []*api.Node
		cluster *api.Cluster
		err     error
	)
	updates, cancel, err := store.ViewAndWatch(
		s.store,
		func(readTx store.ReadTx) error ***REMOVED***
			cluster = store.GetCluster(readTx, s.clusterID)
			if cluster == nil ***REMOVED***
				return errors.New("could not find cluster object")
			***REMOVED***
			nodes, err = store.FindNodes(readTx, store.All)
			return err
		***REMOVED***,
		api.EventCreateNode***REMOVED******REMOVED***,
		api.EventUpdateNode***REMOVED******REMOVED***,
		api.EventDeleteNode***REMOVED******REMOVED***,
		api.EventUpdateCluster***REMOVED***
			Cluster: &api.Cluster***REMOVED***ID: s.clusterID***REMOVED***,
			Checks:  []api.ClusterCheckFunc***REMOVED***api.ClusterCheckID***REMOVED***,
		***REMOVED***,
	)

	// call once to ensure that the join tokens and local/external CA signer are always set
	rootReconciler := &rootRotationReconciler***REMOVED***
		ctx:                 log.WithField(ctx, "method", "(*Server).rootRotationReconciler"),
		clusterID:           s.clusterID,
		store:               s.store,
		batchUpdateInterval: s.rootReconciliationRetryInterval,
	***REMOVED***

	s.UpdateRootCA(ctx, cluster, rootReconciler)

	// Do this after updateCluster has been called, so Ready() and isRunning never returns true without
	// the join tokens and external CA/security config's root CA being set correctly
	s.mu.Lock()
	close(s.started)
	s.mu.Unlock()

	if err != nil ***REMOVED***
		log.G(ctx).WithFields(logrus.Fields***REMOVED***
			"method": "(*Server).Run",
		***REMOVED***).WithError(err).Errorf("snapshot store view failed")
		return err
	***REMOVED***
	defer cancel()

	// We might have missed some updates if there was a leader election,
	// so let's pick up the slack.
	if err := s.reconcileNodeCertificates(ctx, nodes); err != nil ***REMOVED***
		// We don't return here because that means the Run loop would
		// never run. Log an error instead.
		log.G(ctx).WithFields(logrus.Fields***REMOVED***
			"method": "(*Server).Run",
		***REMOVED***).WithError(err).Errorf("error attempting to reconcile certificates")
	***REMOVED***

	ticker := time.NewTicker(s.reconciliationRetryInterval)
	defer ticker.Stop()

	externalTLSCredsChange, externalTLSWatchCancel := s.securityConfig.Watch()
	defer externalTLSWatchCancel()

	// Watch for new nodes being created, new nodes being updated, and changes
	// to the cluster
	for ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return nil
		default:
		***REMOVED***

		select ***REMOVED***
		case event := <-updates:
			switch v := event.(type) ***REMOVED***
			case api.EventCreateNode:
				s.evaluateAndSignNodeCert(ctx, v.Node)
				rootReconciler.UpdateNode(v.Node)
			case api.EventUpdateNode:
				// If this certificate is already at a final state
				// no need to evaluate and sign it.
				if !isFinalState(v.Node.Certificate.Status) ***REMOVED***
					s.evaluateAndSignNodeCert(ctx, v.Node)
				***REMOVED***
				rootReconciler.UpdateNode(v.Node)
			case api.EventDeleteNode:
				rootReconciler.DeleteNode(v.Node)
			case api.EventUpdateCluster:
				if v.Cluster.ID == s.clusterID ***REMOVED***
					s.UpdateRootCA(ctx, v.Cluster, rootReconciler)
				***REMOVED***
			***REMOVED***
		case <-externalTLSCredsChange:
			// The TLS certificates can rotate independently of the root CA (and hence which roots the
			// external CA trusts) and external CA URLs.  It's possible that the root CA update is received
			// before the external TLS cred change notification.  During that period, it is possible that
			// the TLS creds will expire or otherwise fail to authorize against external CAs.  However, in
			// that case signing will just fail with a recoverable connectivity error - the state of the
			// certificate issuance is left as pending, and on the next tick, the server will try to sign
			// all nodes with pending certs again (by which time the TLS cred change will have been
			// received).

			// Note that if the external CA changes, the new external CA *MUST* trust the current server's
			// certificate issuer, and this server's certificates should not be extremely close to expiry,
			// otherwise this server would not be able to get new TLS certificates and will no longer be
			// able to function.
			s.signingMu.Lock()
			s.externalCA.UpdateTLSConfig(NewExternalCATLSConfig(
				s.securityConfig.ClientTLSCreds.Config().Certificates, s.externalCAPool))
			s.signingMu.Unlock()
		case <-ticker.C:
			for _, node := range s.pending ***REMOVED***
				if err := s.evaluateAndSignNodeCert(ctx, node); err != nil ***REMOVED***
					// If this sign operation did not succeed, the rest are
					// unlikely to. Yield so that we don't hammer an external CA.
					// Since the map iteration order is randomized, there is no
					// risk of getting stuck on a problematic CSR.
					break
				***REMOVED***
			***REMOVED***
		case <-ctx.Done():
			return nil
		***REMOVED***
	***REMOVED***
***REMOVED***

// Stop stops the CA and closes all grpc streams.
func (s *Server) Stop() error ***REMOVED***
	s.mu.Lock()

	if !s.isRunning() ***REMOVED***
		s.mu.Unlock()
		return errors.New("CA signer is already stopped")
	***REMOVED***
	s.cancel()
	s.started = make(chan struct***REMOVED******REMOVED***)
	s.joinTokens = nil
	s.mu.Unlock()

	// Wait for Run to complete
	s.wg.Wait()

	return nil
***REMOVED***

// Ready waits on the ready channel and returns when the server is ready to serve.
func (s *Server) Ready() <-chan struct***REMOVED******REMOVED*** ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.started
***REMOVED***

func (s *Server) isRunningLocked() (context.Context, error) ***REMOVED***
	s.mu.Lock()
	if !s.isRunning() ***REMOVED***
		s.mu.Unlock()
		return nil, status.Errorf(codes.Aborted, "CA signer is stopped")
	***REMOVED***
	ctx := s.ctx
	s.mu.Unlock()
	return ctx, nil
***REMOVED***

func (s *Server) isReadyLocked() error ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.isRunning() ***REMOVED***
		return status.Errorf(codes.Aborted, "CA signer is stopped")
	***REMOVED***
	if s.joinTokens == nil ***REMOVED***
		return status.Errorf(codes.Aborted, "CA signer is still starting")
	***REMOVED***
	return nil
***REMOVED***

func (s *Server) isRunning() bool ***REMOVED***
	if s.ctx == nil ***REMOVED***
		return false
	***REMOVED***
	select ***REMOVED***
	case <-s.ctx.Done():
		return false
	default:
	***REMOVED***
	return true
***REMOVED***

// filterExternalCAURLS returns a list of external CA urls filtered by the desired cert.
func filterExternalCAURLS(ctx context.Context, desiredCert, defaultCert []byte, apiExternalCAs []*api.ExternalCA) (urls []string) ***REMOVED***
	desiredCert = NormalizePEMs(desiredCert)

	// TODO(aaronl): In the future, this will be abstracted with an ExternalCA interface that has different
	// implementations for different CA types. At the moment, only CFSSL is supported.
	for i, extCA := range apiExternalCAs ***REMOVED***
		// We want to support old external CA specifications which did not have a CA cert.  If there is no cert specified,
		// we assume it's the old cert
		certForExtCA := extCA.CACert
		if len(certForExtCA) == 0 ***REMOVED***
			certForExtCA = defaultCert
		***REMOVED***
		certForExtCA = NormalizePEMs(certForExtCA)
		if extCA.Protocol != api.ExternalCA_CAProtocolCFSSL ***REMOVED***
			log.G(ctx).Debugf("skipping external CA %d (url: %s) due to unknown protocol type", i, extCA.URL)
			continue
		***REMOVED***
		if !bytes.Equal(certForExtCA, desiredCert) ***REMOVED***
			log.G(ctx).Debugf("skipping external CA %d (url: %s) because it has the wrong CA cert", i, extCA.URL)
			continue
		***REMOVED***
		urls = append(urls, extCA.URL)
	***REMOVED***
	return
***REMOVED***

// UpdateRootCA is called when there are cluster changes, and it ensures that the local RootCA is
// always aware of changes in clusterExpiry and the Root CA key material - this can be called by
// anything to update the root CA material
func (s *Server) UpdateRootCA(ctx context.Context, cluster *api.Cluster, reconciler *rootRotationReconciler) error ***REMOVED***
	s.mu.Lock()
	s.joinTokens = cluster.RootCA.JoinTokens.Copy()
	s.mu.Unlock()
	rCA := cluster.RootCA.Copy()
	if reconciler != nil ***REMOVED***
		reconciler.UpdateRootCA(rCA)
	***REMOVED***

	s.signingMu.Lock()
	defer s.signingMu.Unlock()
	firstSeenCluster := s.lastSeenClusterRootCA == nil && s.lastSeenExternalCAs == nil
	rootCAChanged := len(rCA.CACert) != 0 && !equality.RootCAEqualStable(s.lastSeenClusterRootCA, rCA)
	externalCAChanged := !equality.ExternalCAsEqualStable(s.lastSeenExternalCAs, cluster.Spec.CAConfig.ExternalCAs)
	ctx = log.WithLogger(ctx, log.G(ctx).WithFields(logrus.Fields***REMOVED***
		"cluster.id": cluster.ID,
		"method":     "(*Server).UpdateRootCA",
	***REMOVED***))

	if rootCAChanged ***REMOVED***
		setOrUpdate := "set"
		if !firstSeenCluster ***REMOVED***
			log.G(ctx).Debug("Updating signing root CA and external CA due to change in cluster Root CA")
			setOrUpdate = "updated"
		***REMOVED***
		expiry := DefaultNodeCertExpiration
		if cluster.Spec.CAConfig.NodeCertExpiry != nil ***REMOVED***
			// NodeCertExpiry exists, let's try to parse the duration out of it
			clusterExpiry, err := gogotypes.DurationFromProto(cluster.Spec.CAConfig.NodeCertExpiry)
			if err != nil ***REMOVED***
				log.G(ctx).WithError(err).Warn("failed to parse certificate expiration, using default")
			***REMOVED*** else ***REMOVED***
				// We were able to successfully parse the expiration out of the cluster.
				expiry = clusterExpiry
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// NodeCertExpiry seems to be nil
			log.G(ctx).Warn("no certificate expiration specified, using default")
		***REMOVED***
		// Attempt to update our local RootCA with the new parameters
		updatedRootCA, err := RootCAFromAPI(ctx, rCA, expiry)
		if err != nil ***REMOVED***
			return errors.Wrap(err, "invalid Root CA object in cluster")
		***REMOVED***

		s.localRootCA = &updatedRootCA
		s.externalCAPool = updatedRootCA.Pool
		externalCACert := rCA.CACert
		if rCA.RootRotation != nil ***REMOVED***
			externalCACert = rCA.RootRotation.CACert
			// the external CA has to trust the new CA cert
			s.externalCAPool = x509.NewCertPool()
			s.externalCAPool.AppendCertsFromPEM(rCA.CACert)
			s.externalCAPool.AppendCertsFromPEM(rCA.RootRotation.CACert)
		***REMOVED***
		s.lastSeenExternalCAs = cluster.Spec.CAConfig.Copy().ExternalCAs
		urls := filterExternalCAURLS(ctx, externalCACert, rCA.CACert, s.lastSeenExternalCAs)
		// Replace the external CA with the relevant intermediates, URLS, and TLS config
		s.externalCA = NewExternalCA(updatedRootCA.Intermediates,
			NewExternalCATLSConfig(s.securityConfig.ClientTLSCreds.Config().Certificates, s.externalCAPool), urls...)

		// only update the server cache if we've successfully updated the root CA
		log.G(ctx).Debugf("Root CA %s successfully", setOrUpdate)
		s.lastSeenClusterRootCA = rCA
	***REMOVED*** else if externalCAChanged ***REMOVED***
		// we want to update only if the external CA URLS have changed, since if the root CA has changed we already
		// run similar logic
		if !firstSeenCluster ***REMOVED***
			log.G(ctx).Debug("Updating security config external CA URLs due to change in cluster spec's list of external CAs")
		***REMOVED***
		wantedExternalCACert := rCA.CACert // we want to only add external CA URLs that use this cert
		if rCA.RootRotation != nil ***REMOVED***
			// we're rotating to a new root, so we only want external CAs with the new root cert
			wantedExternalCACert = rCA.RootRotation.CACert
		***REMOVED***
		// Update our external CA with the list of External CA URLs from the new cluster state
		s.lastSeenExternalCAs = cluster.Spec.CAConfig.Copy().ExternalCAs
		urls := filterExternalCAURLS(ctx, wantedExternalCACert, rCA.CACert, s.lastSeenExternalCAs)
		s.externalCA.UpdateURLs(urls...)
	***REMOVED***
	return nil
***REMOVED***

// evaluateAndSignNodeCert implements the logic of which certificates to sign
func (s *Server) evaluateAndSignNodeCert(ctx context.Context, node *api.Node) error ***REMOVED***
	// If the desired membership and actual state are in sync, there's
	// nothing to do.
	certState := node.Certificate.Status.State
	if node.Spec.Membership == api.NodeMembershipAccepted &&
		(certState == api.IssuanceStateIssued || certState == api.IssuanceStateRotate) ***REMOVED***
		return nil
	***REMOVED***

	// If the certificate state is renew, then it is a server-sided accepted cert (cert renewals)
	if certState == api.IssuanceStateRenew ***REMOVED***
		return s.signNodeCert(ctx, node)
	***REMOVED***

	// Sign this certificate if a user explicitly changed it to Accepted, and
	// the certificate is in pending state
	if node.Spec.Membership == api.NodeMembershipAccepted && certState == api.IssuanceStatePending ***REMOVED***
		return s.signNodeCert(ctx, node)
	***REMOVED***

	return nil
***REMOVED***

// signNodeCert does the bulk of the work for signing a certificate
func (s *Server) signNodeCert(ctx context.Context, node *api.Node) error ***REMOVED***
	s.signingMu.Lock()
	rootCA := s.localRootCA
	externalCA := s.externalCA
	s.signingMu.Unlock()

	node = node.Copy()
	nodeID := node.ID
	// Convert the role from proto format
	role, err := ParseRole(node.Certificate.Role)
	if err != nil ***REMOVED***
		log.G(ctx).WithFields(logrus.Fields***REMOVED***
			"node.id": node.ID,
			"method":  "(*Server).signNodeCert",
		***REMOVED***).WithError(err).Errorf("failed to parse role")
		return errors.New("failed to parse role")
	***REMOVED***

	s.pending[node.ID] = node

	// Attempt to sign the CSR
	var (
		rawCSR = node.Certificate.CSR
		cn     = node.Certificate.CN
		ou     = role
		org    = s.clusterID
	)

	// Try using the external CA first.
	cert, err := externalCA.Sign(ctx, PrepareCSR(rawCSR, cn, ou, org))
	if err == ErrNoExternalCAURLs ***REMOVED***
		// No external CA servers configured. Try using the local CA.
		cert, err = rootCA.ParseValidateAndSignCSR(rawCSR, cn, ou, org)
	***REMOVED***

	if err != nil ***REMOVED***
		log.G(ctx).WithFields(logrus.Fields***REMOVED***
			"node.id": node.ID,
			"method":  "(*Server).signNodeCert",
		***REMOVED***).WithError(err).Errorf("failed to sign CSR")

		// If the current state is already Failed, no need to change it
		if node.Certificate.Status.State == api.IssuanceStateFailed ***REMOVED***
			delete(s.pending, node.ID)
			return errors.New("failed to sign CSR")
		***REMOVED***

		if _, ok := err.(recoverableErr); ok ***REMOVED***
			// Return without changing the state of the certificate. We may
			// retry signing it in the future.
			return errors.New("failed to sign CSR")
		***REMOVED***

		// We failed to sign this CSR, change the state to FAILED
		err = s.store.Update(func(tx store.Tx) error ***REMOVED***
			node := store.GetNode(tx, nodeID)
			if node == nil ***REMOVED***
				return errors.Errorf("node %s not found", nodeID)
			***REMOVED***

			node.Certificate.Status = api.IssuanceStatus***REMOVED***
				State: api.IssuanceStateFailed,
				Err:   err.Error(),
			***REMOVED***

			return store.UpdateNode(tx, node)
		***REMOVED***)
		if err != nil ***REMOVED***
			log.G(ctx).WithFields(logrus.Fields***REMOVED***
				"node.id": nodeID,
				"method":  "(*Server).signNodeCert",
			***REMOVED***).WithError(err).Errorf("transaction failed when setting state to FAILED")
		***REMOVED***

		delete(s.pending, node.ID)
		return errors.New("failed to sign CSR")
	***REMOVED***

	// We were able to successfully sign the new CSR. Let's try to update the nodeStore
	for ***REMOVED***
		err = s.store.Update(func(tx store.Tx) error ***REMOVED***
			node.Certificate.Certificate = cert
			node.Certificate.Status = api.IssuanceStatus***REMOVED***
				State: api.IssuanceStateIssued,
			***REMOVED***

			err := store.UpdateNode(tx, node)
			if err != nil ***REMOVED***
				node = store.GetNode(tx, nodeID)
				if node == nil ***REMOVED***
					err = errors.Errorf("node %s does not exist", nodeID)
				***REMOVED***
			***REMOVED***
			return err
		***REMOVED***)
		if err == nil ***REMOVED***
			log.G(ctx).WithFields(logrus.Fields***REMOVED***
				"node.id":   node.ID,
				"node.role": node.Certificate.Role,
				"method":    "(*Server).signNodeCert",
			***REMOVED***).Debugf("certificate issued")
			delete(s.pending, node.ID)
			break
		***REMOVED***
		if err == store.ErrSequenceConflict ***REMOVED***
			continue
		***REMOVED***

		log.G(ctx).WithFields(logrus.Fields***REMOVED***
			"node.id": nodeID,
			"method":  "(*Server).signNodeCert",
		***REMOVED***).WithError(err).Errorf("transaction failed")
		return errors.New("transaction failed")
	***REMOVED***
	return nil
***REMOVED***

// reconcileNodeCertificates is a helper method that calls evaluateAndSignNodeCert on all the
// nodes.
func (s *Server) reconcileNodeCertificates(ctx context.Context, nodes []*api.Node) error ***REMOVED***
	for _, node := range nodes ***REMOVED***
		s.evaluateAndSignNodeCert(ctx, node)
	***REMOVED***

	return nil
***REMOVED***

// A successfully issued certificate and a failed certificate are our current final states
func isFinalState(status api.IssuanceStatus) bool ***REMOVED***
	if status.State == api.IssuanceStateIssued || status.State == api.IssuanceStateFailed ||
		status.State == api.IssuanceStateRotate ***REMOVED***
		return true
	***REMOVED***

	return false
***REMOVED***

// RootCAFromAPI creates a RootCA object from an api.RootCA object
func RootCAFromAPI(ctx context.Context, apiRootCA *api.RootCA, expiry time.Duration) (RootCA, error) ***REMOVED***
	var intermediates []byte
	signingCert := apiRootCA.CACert
	signingKey := apiRootCA.CAKey
	if apiRootCA.RootRotation != nil ***REMOVED***
		signingCert = apiRootCA.RootRotation.CrossSignedCACert
		signingKey = apiRootCA.RootRotation.CAKey
		intermediates = apiRootCA.RootRotation.CrossSignedCACert
	***REMOVED***
	if signingKey == nil ***REMOVED***
		signingCert = nil
	***REMOVED***
	return NewRootCA(apiRootCA.CACert, signingCert, signingKey, expiry, intermediates)
***REMOVED***
