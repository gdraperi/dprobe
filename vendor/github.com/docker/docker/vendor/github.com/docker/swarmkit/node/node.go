package node

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/docker/docker/pkg/plugingetter"
	metrics "github.com/docker/go-metrics"
	"github.com/docker/swarmkit/agent"
	"github.com/docker/swarmkit/agent/exec"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/ca"
	"github.com/docker/swarmkit/connectionbroker"
	"github.com/docker/swarmkit/ioutils"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager"
	"github.com/docker/swarmkit/manager/encryption"
	"github.com/docker/swarmkit/remotes"
	"github.com/docker/swarmkit/xnet"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
)

const (
	stateFilename     = "state.json"
	roleChangeTimeout = 16 * time.Second
)

var (
	nodeInfo    metrics.LabeledGauge
	nodeManager metrics.Gauge

	errNodeStarted    = errors.New("node: already started")
	errNodeNotStarted = errors.New("node: not started")
	certDirectory     = "certificates"

	// ErrInvalidUnlockKey is returned when we can't decrypt the TLS certificate
	ErrInvalidUnlockKey = errors.New("node is locked, and needs a valid unlock key")
)

func init() ***REMOVED***
	ns := metrics.NewNamespace("swarm", "node", nil)
	nodeInfo = ns.NewLabeledGauge("info", "Information related to the swarm", "",
		"swarm_id",
		"node_id",
	)
	nodeManager = ns.NewGauge("manager", "Whether this node is a manager or not", "")
	metrics.Register(ns)
***REMOVED***

// Config provides values for a Node.
type Config struct ***REMOVED***
	// Hostname is the name of host for agent instance.
	Hostname string

	// JoinAddr specifies node that should be used for the initial connection to
	// other manager in cluster. This should be only one address and optional,
	// the actual remotes come from the stored state.
	JoinAddr string

	// StateDir specifies the directory the node uses to keep the state of the
	// remote managers and certificates.
	StateDir string

	// JoinToken is the token to be used on the first certificate request.
	JoinToken string

	// ExternalCAs is a list of CAs to which a manager node
	// will make certificate signing requests for node certificates.
	ExternalCAs []*api.ExternalCA

	// ForceNewCluster creates a new cluster from current raft state.
	ForceNewCluster bool

	// ListenControlAPI specifies address the control API should listen on.
	ListenControlAPI string

	// ListenRemoteAPI specifies the address for the remote API that agents
	// and raft members connect to.
	ListenRemoteAPI string

	// AdvertiseRemoteAPI specifies the address that should be advertised
	// for connections to the remote API (including the raft service).
	AdvertiseRemoteAPI string

	// Executor specifies the executor to use for the agent.
	Executor exec.Executor

	// ElectionTick defines the amount of ticks needed without
	// leader to trigger a new election
	ElectionTick uint32

	// HeartbeatTick defines the amount of ticks between each
	// heartbeat sent to other members for health-check purposes
	HeartbeatTick uint32

	// AutoLockManagers determines whether or not an unlock key will be generated
	// when bootstrapping a new cluster for the first time
	AutoLockManagers bool

	// UnlockKey is the key to unlock a node - used for decrypting at rest.  This
	// only applies to nodes that have already joined a cluster.
	UnlockKey []byte

	// Availability allows a user to control the current scheduling status of a node
	Availability api.NodeSpec_Availability

	// PluginGetter provides access to docker's plugin inventory.
	PluginGetter plugingetter.PluginGetter
***REMOVED***

// Node implements the primary node functionality for a member of a swarm
// cluster. Node handles workloads and may also run as a manager.
type Node struct ***REMOVED***
	sync.RWMutex
	config           *Config
	remotes          *persistentRemotes
	connBroker       *connectionbroker.Broker
	role             string
	roleCond         *sync.Cond
	conn             *grpc.ClientConn
	connCond         *sync.Cond
	nodeID           string
	started          chan struct***REMOVED******REMOVED***
	startOnce        sync.Once
	stopped          chan struct***REMOVED******REMOVED***
	stopOnce         sync.Once
	ready            chan struct***REMOVED******REMOVED*** // closed when agent has completed registration and manager(if enabled) is ready to receive control requests
	closed           chan struct***REMOVED******REMOVED***
	err              error
	agent            *agent.Agent
	manager          *manager.Manager
	notifyNodeChange chan *agent.NodeChanges // used by the agent to relay node updates from the dispatcher Session stream to (*Node).run
	unlockKey        []byte
***REMOVED***

type lastSeenRole struct ***REMOVED***
	role api.NodeRole
***REMOVED***

// observe notes the latest value of this node role, and returns true if it
// is the first seen value, or is different from the most recently seen value.
func (l *lastSeenRole) observe(newRole api.NodeRole) bool ***REMOVED***
	changed := l.role != newRole
	l.role = newRole
	return changed
***REMOVED***

// RemoteAPIAddr returns address on which remote manager api listens.
// Returns nil if node is not manager.
func (n *Node) RemoteAPIAddr() (string, error) ***REMOVED***
	n.RLock()
	defer n.RUnlock()
	if n.manager == nil ***REMOVED***
		return "", errors.New("manager is not running")
	***REMOVED***
	addr := n.manager.Addr()
	if addr == "" ***REMOVED***
		return "", errors.New("manager addr is not set")
	***REMOVED***
	return addr, nil
***REMOVED***

// New returns new Node instance.
func New(c *Config) (*Node, error) ***REMOVED***
	if err := os.MkdirAll(c.StateDir, 0700); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	stateFile := filepath.Join(c.StateDir, stateFilename)
	dt, err := ioutil.ReadFile(stateFile)
	var p []api.Peer
	if err != nil && !os.IsNotExist(err) ***REMOVED***
		return nil, err
	***REMOVED***
	if err == nil ***REMOVED***
		if err := json.Unmarshal(dt, &p); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	n := &Node***REMOVED***
		remotes:          newPersistentRemotes(stateFile, p...),
		role:             ca.WorkerRole,
		config:           c,
		started:          make(chan struct***REMOVED******REMOVED***),
		stopped:          make(chan struct***REMOVED******REMOVED***),
		closed:           make(chan struct***REMOVED******REMOVED***),
		ready:            make(chan struct***REMOVED******REMOVED***),
		notifyNodeChange: make(chan *agent.NodeChanges, 1),
		unlockKey:        c.UnlockKey,
	***REMOVED***

	if n.config.JoinAddr != "" || n.config.ForceNewCluster ***REMOVED***
		n.remotes = newPersistentRemotes(filepath.Join(n.config.StateDir, stateFilename))
		if n.config.JoinAddr != "" ***REMOVED***
			n.remotes.Observe(api.Peer***REMOVED***Addr: n.config.JoinAddr***REMOVED***, remotes.DefaultObservationWeight)
		***REMOVED***
	***REMOVED***

	n.connBroker = connectionbroker.New(n.remotes)

	n.roleCond = sync.NewCond(n.RLocker())
	n.connCond = sync.NewCond(n.RLocker())
	return n, nil
***REMOVED***

// BindRemote starts a listener that exposes the remote API.
func (n *Node) BindRemote(ctx context.Context, listenAddr string, advertiseAddr string) error ***REMOVED***
	n.RLock()
	defer n.RUnlock()

	if n.manager == nil ***REMOVED***
		return errors.New("manager is not running")
	***REMOVED***

	return n.manager.BindRemote(ctx, manager.RemoteAddrs***REMOVED***
		ListenAddr:    listenAddr,
		AdvertiseAddr: advertiseAddr,
	***REMOVED***)
***REMOVED***

// Start starts a node instance.
func (n *Node) Start(ctx context.Context) error ***REMOVED***
	err := errNodeStarted

	n.startOnce.Do(func() ***REMOVED***
		close(n.started)
		go n.run(ctx)
		err = nil // clear error above, only once.
	***REMOVED***)

	return err
***REMOVED***

func (n *Node) currentRole() api.NodeRole ***REMOVED***
	n.Lock()
	currentRole := api.NodeRoleWorker
	if n.role == ca.ManagerRole ***REMOVED***
		currentRole = api.NodeRoleManager
	***REMOVED***
	n.Unlock()
	return currentRole
***REMOVED***

func (n *Node) run(ctx context.Context) (err error) ***REMOVED***
	defer func() ***REMOVED***
		n.err = err
		close(n.closed)
	***REMOVED***()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ctx = log.WithModule(ctx, "node")

	go func(ctx context.Context) ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
		case <-n.stopped:
			cancel()
		***REMOVED***
	***REMOVED***(ctx)

	paths := ca.NewConfigPaths(filepath.Join(n.config.StateDir, certDirectory))
	securityConfig, secConfigCancel, err := n.loadSecurityConfig(ctx, paths)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer secConfigCancel()

	renewer := ca.NewTLSRenewer(securityConfig, n.connBroker, paths.RootCA)

	ctx = log.WithLogger(ctx, log.G(ctx).WithField("node.id", n.NodeID()))

	taskDBPath := filepath.Join(n.config.StateDir, "worker", "tasks.db")
	if err := os.MkdirAll(filepath.Dir(taskDBPath), 0777); err != nil ***REMOVED***
		return err
	***REMOVED***

	db, err := bolt.Open(taskDBPath, 0666, nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer db.Close()

	agentDone := make(chan struct***REMOVED******REMOVED***)

	go func() ***REMOVED***
		// lastNodeDesiredRole is the last-seen value of Node.Spec.DesiredRole,
		// used to make role changes "edge triggered" and avoid renewal loops.
		lastNodeDesiredRole := lastSeenRole***REMOVED***role: n.currentRole()***REMOVED***

		for ***REMOVED***
			select ***REMOVED***
			case <-agentDone:
				return
			case nodeChanges := <-n.notifyNodeChange:
				if nodeChanges.Node != nil ***REMOVED***
					// This is a bit complex to be backward compatible with older CAs that
					// don't support the Node.Role field. They only use what's presently
					// called DesiredRole.
					// 1) If DesiredRole changes, kick off a certificate renewal. The renewal
					//    is delayed slightly to give Role time to change as well if this is
					//    a newer CA. If the certificate we get back doesn't have the expected
					//    role, we continue renewing with exponential backoff.
					// 2) If the server is sending us IssuanceStateRotate, renew the cert as
					//    requested by the CA.
					desiredRoleChanged := lastNodeDesiredRole.observe(nodeChanges.Node.Spec.DesiredRole)
					if desiredRoleChanged ***REMOVED***
						switch nodeChanges.Node.Spec.DesiredRole ***REMOVED***
						case api.NodeRoleManager:
							renewer.SetExpectedRole(ca.ManagerRole)
						case api.NodeRoleWorker:
							renewer.SetExpectedRole(ca.WorkerRole)
						***REMOVED***
					***REMOVED***
					if desiredRoleChanged || nodeChanges.Node.Certificate.Status.State == api.IssuanceStateRotate ***REMOVED***
						renewer.Renew()
					***REMOVED***
				***REMOVED***

				if nodeChanges.RootCert != nil ***REMOVED***
					if bytes.Equal(nodeChanges.RootCert, securityConfig.RootCA().Certs) ***REMOVED***
						continue
					***REMOVED***
					newRootCA, err := ca.NewRootCA(nodeChanges.RootCert, nil, nil, ca.DefaultNodeCertExpiration, nil)
					if err != nil ***REMOVED***
						log.G(ctx).WithError(err).Error("invalid new root certificate from the dispatcher")
						continue
					***REMOVED***
					if err := securityConfig.UpdateRootCA(&newRootCA); err != nil ***REMOVED***
						log.G(ctx).WithError(err).Error("could not use new root CA from dispatcher")
						continue
					***REMOVED***
					if err := ca.SaveRootCA(newRootCA, paths.RootCA); err != nil ***REMOVED***
						log.G(ctx).WithError(err).Error("could not save new root certificate from the dispatcher")
						continue
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	var wg sync.WaitGroup
	wg.Add(3)

	nodeInfo.WithValues(
		securityConfig.ClientTLSCreds.Organization(),
		securityConfig.ClientTLSCreds.NodeID(),
	).Set(1)

	if n.currentRole() == api.NodeRoleManager ***REMOVED***
		nodeManager.Set(1)
	***REMOVED*** else ***REMOVED***
		nodeManager.Set(0)
	***REMOVED***

	updates := renewer.Start(ctx)
	go func() ***REMOVED***
		for certUpdate := range updates ***REMOVED***
			if certUpdate.Err != nil ***REMOVED***
				logrus.Warnf("error renewing TLS certificate: %v", certUpdate.Err)
				continue
			***REMOVED***
			n.Lock()
			n.role = certUpdate.Role
			n.roleCond.Broadcast()
			n.Unlock()

			// Export the new role.
			if n.currentRole() == api.NodeRoleManager ***REMOVED***
				nodeManager.Set(1)
			***REMOVED*** else ***REMOVED***
				nodeManager.Set(0)
			***REMOVED***
		***REMOVED***

		wg.Done()
	***REMOVED***()

	role := n.role

	managerReady := make(chan struct***REMOVED******REMOVED***)
	agentReady := make(chan struct***REMOVED******REMOVED***)
	var managerErr error
	var agentErr error
	go func() ***REMOVED***
		managerErr = n.superviseManager(ctx, securityConfig, paths.RootCA, managerReady, renewer) // store err and loop
		wg.Done()
		cancel()
	***REMOVED***()
	go func() ***REMOVED***
		agentErr = n.runAgent(ctx, db, securityConfig, agentReady)
		wg.Done()
		cancel()
		close(agentDone)
	***REMOVED***()

	go func() ***REMOVED***
		<-agentReady
		if role == ca.ManagerRole ***REMOVED***
			workerRole := make(chan struct***REMOVED******REMOVED***)
			waitRoleCtx, waitRoleCancel := context.WithCancel(ctx)
			go func() ***REMOVED***
				if n.waitRole(waitRoleCtx, ca.WorkerRole) == nil ***REMOVED***
					close(workerRole)
				***REMOVED***
			***REMOVED***()
			select ***REMOVED***
			case <-managerReady:
			case <-workerRole:
			***REMOVED***
			waitRoleCancel()
		***REMOVED***
		close(n.ready)
	***REMOVED***()

	wg.Wait()
	if managerErr != nil && errors.Cause(managerErr) != context.Canceled ***REMOVED***
		return managerErr
	***REMOVED***
	if agentErr != nil && errors.Cause(agentErr) != context.Canceled ***REMOVED***
		return agentErr
	***REMOVED***
	return err
***REMOVED***

// Stop stops node execution
func (n *Node) Stop(ctx context.Context) error ***REMOVED***
	select ***REMOVED***
	case <-n.started:
	default:
		return errNodeNotStarted
	***REMOVED***
	// ask agent to clean up assignments
	n.Lock()
	if n.agent != nil ***REMOVED***
		if err := n.agent.Leave(ctx); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("agent failed to clean up assignments")
		***REMOVED***
	***REMOVED***
	n.Unlock()

	n.stopOnce.Do(func() ***REMOVED***
		close(n.stopped)
	***REMOVED***)

	select ***REMOVED***
	case <-n.closed:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	***REMOVED***
***REMOVED***

// Err returns the error that caused the node to shutdown or nil. Err blocks
// until the node has fully shut down.
func (n *Node) Err(ctx context.Context) error ***REMOVED***
	select ***REMOVED***
	case <-n.closed:
		return n.err
	case <-ctx.Done():
		return ctx.Err()
	***REMOVED***
***REMOVED***

func (n *Node) runAgent(ctx context.Context, db *bolt.DB, securityConfig *ca.SecurityConfig, ready chan<- struct***REMOVED******REMOVED***) error ***REMOVED***
	waitCtx, waitCancel := context.WithCancel(ctx)
	remotesCh := n.remotes.WaitSelect(ctx)
	controlCh := n.ListenControlSocket(waitCtx)

waitPeer:
	for ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			break waitPeer
		case <-remotesCh:
			break waitPeer
		case conn := <-controlCh:
			if conn != nil ***REMOVED***
				break waitPeer
			***REMOVED***
		***REMOVED***
	***REMOVED***

	waitCancel()

	select ***REMOVED***
	case <-ctx.Done():
		return ctx.Err()
	default:
	***REMOVED***

	secChangesCh, secChangesCancel := securityConfig.Watch()
	defer secChangesCancel()

	rootCA := securityConfig.RootCA()
	issuer := securityConfig.IssuerInfo()

	agentConfig := &agent.Config***REMOVED***
		Hostname:         n.config.Hostname,
		ConnBroker:       n.connBroker,
		Executor:         n.config.Executor,
		DB:               db,
		NotifyNodeChange: n.notifyNodeChange,
		NotifyTLSChange:  secChangesCh,
		Credentials:      securityConfig.ClientTLSCreds,
		NodeTLSInfo: &api.NodeTLSInfo***REMOVED***
			TrustRoot:           rootCA.Certs,
			CertIssuerPublicKey: issuer.PublicKey,
			CertIssuerSubject:   issuer.Subject,
		***REMOVED***,
	***REMOVED***
	// if a join address has been specified, then if the agent fails to connect due to a TLS error, fail fast - don't
	// keep re-trying to join
	if n.config.JoinAddr != "" ***REMOVED***
		agentConfig.SessionTracker = &firstSessionErrorTracker***REMOVED******REMOVED***
	***REMOVED***

	a, err := agent.New(agentConfig)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := a.Start(ctx); err != nil ***REMOVED***
		return err
	***REMOVED***

	n.Lock()
	n.agent = a
	n.Unlock()

	defer func() ***REMOVED***
		n.Lock()
		n.agent = nil
		n.Unlock()
	***REMOVED***()

	go func() ***REMOVED***
		<-a.Ready()
		close(ready)
	***REMOVED***()

	// todo: manually call stop on context cancellation?

	return a.Err(context.Background())
***REMOVED***

// Ready returns a channel that is closed after node's initialization has
// completes for the first time.
func (n *Node) Ready() <-chan struct***REMOVED******REMOVED*** ***REMOVED***
	return n.ready
***REMOVED***

func (n *Node) setControlSocket(conn *grpc.ClientConn) ***REMOVED***
	n.Lock()
	if n.conn != nil ***REMOVED***
		n.conn.Close()
	***REMOVED***
	n.conn = conn
	n.connBroker.SetLocalConn(conn)
	n.connCond.Broadcast()
	n.Unlock()
***REMOVED***

// ListenControlSocket listens changes of a connection for managing the
// cluster control api
func (n *Node) ListenControlSocket(ctx context.Context) <-chan *grpc.ClientConn ***REMOVED***
	c := make(chan *grpc.ClientConn, 1)
	n.RLock()
	conn := n.conn
	c <- conn
	done := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			n.connCond.Broadcast()
		case <-done:
		***REMOVED***
	***REMOVED***()
	go func() ***REMOVED***
		defer close(c)
		defer close(done)
		defer n.RUnlock()
		for ***REMOVED***
			select ***REMOVED***
			case <-ctx.Done():
				return
			default:
			***REMOVED***
			if conn == n.conn ***REMOVED***
				n.connCond.Wait()
				continue
			***REMOVED***
			conn = n.conn
			select ***REMOVED***
			case c <- conn:
			case <-ctx.Done():
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	return c
***REMOVED***

// NodeID returns current node's ID. May be empty if not set.
func (n *Node) NodeID() string ***REMOVED***
	n.RLock()
	defer n.RUnlock()
	return n.nodeID
***REMOVED***

// Manager returns manager instance started by node. May be nil.
func (n *Node) Manager() *manager.Manager ***REMOVED***
	n.RLock()
	defer n.RUnlock()
	return n.manager
***REMOVED***

// Agent returns agent instance started by node. May be nil.
func (n *Node) Agent() *agent.Agent ***REMOVED***
	n.RLock()
	defer n.RUnlock()
	return n.agent
***REMOVED***

// IsStateDirty returns true if any objects have been added to raft which make
// the state "dirty". Currently, the existence of any object other than the
// default cluster or the local node implies a dirty state.
func (n *Node) IsStateDirty() (bool, error) ***REMOVED***
	n.RLock()
	defer n.RUnlock()

	if n.manager == nil ***REMOVED***
		return false, errors.New("node is not a manager")
	***REMOVED***

	return n.manager.IsStateDirty()
***REMOVED***

// Remotes returns a list of known peers known to node.
func (n *Node) Remotes() []api.Peer ***REMOVED***
	weights := n.remotes.Weights()
	remotes := make([]api.Peer, 0, len(weights))
	for p := range weights ***REMOVED***
		remotes = append(remotes, p)
	***REMOVED***
	return remotes
***REMOVED***

func (n *Node) loadSecurityConfig(ctx context.Context, paths *ca.SecurityConfigPaths) (*ca.SecurityConfig, func() error, error) ***REMOVED***
	var (
		securityConfig *ca.SecurityConfig
		cancel         func() error
	)

	krw := ca.NewKeyReadWriter(paths.Node, n.unlockKey, &manager.RaftDEKData***REMOVED******REMOVED***)
	if err := krw.Migrate(); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	// Check if we already have a valid certificates on disk.
	rootCA, err := ca.GetLocalRootCA(paths.RootCA)
	if err != nil && err != ca.ErrNoLocalRootCA ***REMOVED***
		return nil, nil, err
	***REMOVED***
	if err == nil ***REMOVED***
		// if forcing a new cluster, we allow the certificates to be expired - a new set will be generated
		securityConfig, cancel, err = ca.LoadSecurityConfig(ctx, rootCA, krw, n.config.ForceNewCluster)
		if err != nil ***REMOVED***
			_, isInvalidKEK := errors.Cause(err).(ca.ErrInvalidKEK)
			if isInvalidKEK ***REMOVED***
				return nil, nil, ErrInvalidUnlockKey
			***REMOVED*** else if !os.IsNotExist(err) ***REMOVED***
				return nil, nil, errors.Wrapf(err, "error while loading TLS certificate in %s", paths.Node.Cert)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if securityConfig == nil ***REMOVED***
		if n.config.JoinAddr == "" ***REMOVED***
			// if we're not joining a cluster, bootstrap a new one - and we have to set the unlock key
			n.unlockKey = nil
			if n.config.AutoLockManagers ***REMOVED***
				n.unlockKey = encryption.GenerateSecretKey()
			***REMOVED***
			krw = ca.NewKeyReadWriter(paths.Node, n.unlockKey, &manager.RaftDEKData***REMOVED******REMOVED***)
			rootCA, err = ca.CreateRootCA(ca.DefaultRootCN)
			if err != nil ***REMOVED***
				return nil, nil, err
			***REMOVED***
			if err := ca.SaveRootCA(rootCA, paths.RootCA); err != nil ***REMOVED***
				return nil, nil, err
			***REMOVED***
			log.G(ctx).Debug("generated CA key and certificate")
		***REMOVED*** else if err == ca.ErrNoLocalRootCA ***REMOVED*** // from previous error loading the root CA from disk
			rootCA, err = ca.DownloadRootCA(ctx, paths.RootCA, n.config.JoinToken, n.connBroker)
			if err != nil ***REMOVED***
				return nil, nil, err
			***REMOVED***
			log.G(ctx).Debug("downloaded CA certificate")
		***REMOVED***

		// Obtain new certs and setup TLS certificates renewal for this node:
		// - If certificates weren't present on disk, we call CreateSecurityConfig, which blocks
		//   until a valid certificate has been issued.
		// - We wait for CreateSecurityConfig to finish since we need a certificate to operate.

		// Attempt to load certificate from disk
		securityConfig, cancel, err = ca.LoadSecurityConfig(ctx, rootCA, krw, n.config.ForceNewCluster)
		if err == nil ***REMOVED***
			log.G(ctx).WithFields(logrus.Fields***REMOVED***
				"node.id": securityConfig.ClientTLSCreds.NodeID(),
			***REMOVED***).Debugf("loaded TLS certificate")
		***REMOVED*** else ***REMOVED***
			if _, ok := errors.Cause(err).(ca.ErrInvalidKEK); ok ***REMOVED***
				return nil, nil, ErrInvalidUnlockKey
			***REMOVED***
			log.G(ctx).WithError(err).Debugf("no node credentials found in: %s", krw.Target())

			securityConfig, cancel, err = rootCA.CreateSecurityConfig(ctx, krw, ca.CertificateRequestConfig***REMOVED***
				Token:        n.config.JoinToken,
				Availability: n.config.Availability,
				ConnBroker:   n.connBroker,
			***REMOVED***)

			if err != nil ***REMOVED***
				return nil, nil, err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	n.Lock()
	n.role = securityConfig.ClientTLSCreds.Role()
	n.nodeID = securityConfig.ClientTLSCreds.NodeID()
	n.roleCond.Broadcast()
	n.Unlock()

	return securityConfig, cancel, nil
***REMOVED***

func (n *Node) initManagerConnection(ctx context.Context, ready chan<- struct***REMOVED******REMOVED***) error ***REMOVED***
	opts := []grpc.DialOption***REMOVED***
		grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
		grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor),
	***REMOVED***
	insecureCreds := credentials.NewTLS(&tls.Config***REMOVED***InsecureSkipVerify: true***REMOVED***)
	opts = append(opts, grpc.WithTransportCredentials(insecureCreds))
	addr := n.config.ListenControlAPI
	opts = append(opts, grpc.WithDialer(
		func(addr string, timeout time.Duration) (net.Conn, error) ***REMOVED***
			return xnet.DialTimeoutLocal(addr, timeout)
		***REMOVED***))
	conn, err := grpc.Dial(addr, opts...)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	client := api.NewHealthClient(conn)
	for ***REMOVED***
		resp, err := client.Check(ctx, &api.HealthCheckRequest***REMOVED***Service: "ControlAPI"***REMOVED***)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if resp.Status == api.HealthCheckResponse_SERVING ***REMOVED***
			break
		***REMOVED***
		time.Sleep(500 * time.Millisecond)
	***REMOVED***
	n.setControlSocket(conn)
	if ready != nil ***REMOVED***
		close(ready)
	***REMOVED***
	return nil
***REMOVED***

func (n *Node) waitRole(ctx context.Context, role string) error ***REMOVED***
	n.roleCond.L.Lock()
	if role == n.role ***REMOVED***
		n.roleCond.L.Unlock()
		return nil
	***REMOVED***
	finishCh := make(chan struct***REMOVED******REMOVED***)
	defer close(finishCh)
	go func() ***REMOVED***
		select ***REMOVED***
		case <-finishCh:
		case <-ctx.Done():
			// call broadcast to shutdown this function
			n.roleCond.Broadcast()
		***REMOVED***
	***REMOVED***()
	defer n.roleCond.L.Unlock()
	for role != n.role ***REMOVED***
		n.roleCond.Wait()
		select ***REMOVED***
		case <-ctx.Done():
			return ctx.Err()
		default:
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (n *Node) runManager(ctx context.Context, securityConfig *ca.SecurityConfig, rootPaths ca.CertPaths, ready chan struct***REMOVED******REMOVED***, workerRole <-chan struct***REMOVED******REMOVED***) (bool, error) ***REMOVED***
	var remoteAPI *manager.RemoteAddrs
	if n.config.ListenRemoteAPI != "" ***REMOVED***
		remoteAPI = &manager.RemoteAddrs***REMOVED***
			ListenAddr:    n.config.ListenRemoteAPI,
			AdvertiseAddr: n.config.AdvertiseRemoteAPI,
		***REMOVED***
	***REMOVED***

	joinAddr := n.config.JoinAddr
	if joinAddr == "" ***REMOVED***
		remoteAddr, err := n.remotes.Select(n.NodeID())
		if err == nil ***REMOVED***
			joinAddr = remoteAddr.Addr
		***REMOVED***
	***REMOVED***

	m, err := manager.New(&manager.Config***REMOVED***
		ForceNewCluster:  n.config.ForceNewCluster,
		RemoteAPI:        remoteAPI,
		ControlAPI:       n.config.ListenControlAPI,
		SecurityConfig:   securityConfig,
		ExternalCAs:      n.config.ExternalCAs,
		JoinRaft:         joinAddr,
		ForceJoin:        n.config.JoinAddr != "",
		StateDir:         n.config.StateDir,
		HeartbeatTick:    n.config.HeartbeatTick,
		ElectionTick:     n.config.ElectionTick,
		AutoLockManagers: n.config.AutoLockManagers,
		UnlockKey:        n.unlockKey,
		Availability:     n.config.Availability,
		PluginGetter:     n.config.PluginGetter,
		RootCAPaths:      rootPaths,
	***REMOVED***)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	done := make(chan struct***REMOVED******REMOVED***)
	var runErr error
	go func(logger *logrus.Entry) ***REMOVED***
		if err := m.Run(log.WithLogger(context.Background(), logger)); err != nil ***REMOVED***
			runErr = err
		***REMOVED***
		close(done)
	***REMOVED***(log.G(ctx))

	var clearData bool
	defer func() ***REMOVED***
		n.Lock()
		n.manager = nil
		n.Unlock()
		m.Stop(ctx, clearData)
		<-done
		n.setControlSocket(nil)
	***REMOVED***()

	n.Lock()
	n.manager = m
	n.Unlock()

	connCtx, connCancel := context.WithCancel(ctx)
	defer connCancel()

	go n.initManagerConnection(connCtx, ready)

	// wait for manager stop or for role change
	select ***REMOVED***
	case <-done:
		return false, runErr
	case <-workerRole:
		log.G(ctx).Info("role changed to worker, stopping manager")
		clearData = true
	case <-m.RemovedFromRaft():
		log.G(ctx).Info("manager removed from raft cluster, stopping manager")
		clearData = true
	case <-ctx.Done():
		return false, ctx.Err()
	***REMOVED***
	return clearData, nil
***REMOVED***

func (n *Node) superviseManager(ctx context.Context, securityConfig *ca.SecurityConfig, rootPaths ca.CertPaths, ready chan struct***REMOVED******REMOVED***, renewer *ca.TLSRenewer) error ***REMOVED***
	for ***REMOVED***
		if err := n.waitRole(ctx, ca.ManagerRole); err != nil ***REMOVED***
			return err
		***REMOVED***

		workerRole := make(chan struct***REMOVED******REMOVED***)
		waitRoleCtx, waitRoleCancel := context.WithCancel(ctx)
		go func() ***REMOVED***
			if n.waitRole(waitRoleCtx, ca.WorkerRole) == nil ***REMOVED***
				close(workerRole)
			***REMOVED***
		***REMOVED***()

		wasRemoved, err := n.runManager(ctx, securityConfig, rootPaths, ready, workerRole)
		if err != nil ***REMOVED***
			waitRoleCancel()
			return errors.Wrap(err, "manager stopped")
		***REMOVED***

		// If the manager stopped running and our role is still
		// "manager", it's possible that the manager was demoted and
		// the agent hasn't realized this yet. We should wait for the
		// role to change instead of restarting the manager immediately.
		err = func() error ***REMOVED***
			timer := time.NewTimer(roleChangeTimeout)
			defer timer.Stop()
			defer waitRoleCancel()

			select ***REMOVED***
			case <-timer.C:
			case <-workerRole:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			***REMOVED***

			if !wasRemoved ***REMOVED***
				log.G(ctx).Warn("failed to get worker role after manager stop, restarting manager")
				return nil
			***REMOVED***
			// We need to be extra careful about restarting the
			// manager. It may cause the node to wrongly join under
			// a new Raft ID. Since we didn't see a role change
			// yet, force a certificate renewal. If the certificate
			// comes back with a worker role, we know we shouldn't
			// restart the manager. However, if we don't see
			// workerRole get closed, it means we didn't switch to
			// a worker certificate, either because we couldn't
			// contact a working CA, or because we've been
			// re-promoted. In this case, we must assume we were
			// re-promoted, and restart the manager.
			log.G(ctx).Warn("failed to get worker role after manager stop, forcing certificate renewal")
			timer.Reset(roleChangeTimeout)

			renewer.Renew()

			// Now that the renewal request has been sent to the
			// renewal goroutine, wait for a change in role.
			select ***REMOVED***
			case <-timer.C:
				log.G(ctx).Warn("failed to get worker role after manager stop, restarting manager")
			case <-workerRole:
			case <-ctx.Done():
				return ctx.Err()
			***REMOVED***
			return nil
		***REMOVED***()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		ready = nil
	***REMOVED***
***REMOVED***

type persistentRemotes struct ***REMOVED***
	sync.RWMutex
	c *sync.Cond
	remotes.Remotes
	storePath      string
	lastSavedState []api.Peer
***REMOVED***

func newPersistentRemotes(f string, peers ...api.Peer) *persistentRemotes ***REMOVED***
	pr := &persistentRemotes***REMOVED***
		storePath: f,
		Remotes:   remotes.NewRemotes(peers...),
	***REMOVED***
	pr.c = sync.NewCond(pr.RLocker())
	return pr
***REMOVED***

func (s *persistentRemotes) Observe(peer api.Peer, weight int) ***REMOVED***
	s.Lock()
	defer s.Unlock()
	s.Remotes.Observe(peer, weight)
	s.c.Broadcast()
	if err := s.save(); err != nil ***REMOVED***
		logrus.Errorf("error writing cluster state file: %v", err)
		return
	***REMOVED***
	return
***REMOVED***
func (s *persistentRemotes) Remove(peers ...api.Peer) ***REMOVED***
	s.Lock()
	defer s.Unlock()
	s.Remotes.Remove(peers...)
	if err := s.save(); err != nil ***REMOVED***
		logrus.Errorf("error writing cluster state file: %v", err)
		return
	***REMOVED***
	return
***REMOVED***

func (s *persistentRemotes) save() error ***REMOVED***
	weights := s.Weights()
	remotes := make([]api.Peer, 0, len(weights))
	for r := range weights ***REMOVED***
		remotes = append(remotes, r)
	***REMOVED***
	sort.Sort(sortablePeers(remotes))
	if reflect.DeepEqual(remotes, s.lastSavedState) ***REMOVED***
		return nil
	***REMOVED***
	dt, err := json.Marshal(remotes)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	s.lastSavedState = remotes
	return ioutils.AtomicWriteFile(s.storePath, dt, 0600)
***REMOVED***

// WaitSelect waits until at least one remote becomes available and then selects one.
func (s *persistentRemotes) WaitSelect(ctx context.Context) <-chan api.Peer ***REMOVED***
	c := make(chan api.Peer, 1)
	s.RLock()
	done := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			s.c.Broadcast()
		case <-done:
		***REMOVED***
	***REMOVED***()
	go func() ***REMOVED***
		defer s.RUnlock()
		defer close(c)
		defer close(done)
		for ***REMOVED***
			if ctx.Err() != nil ***REMOVED***
				return
			***REMOVED***
			p, err := s.Select()
			if err == nil ***REMOVED***
				c <- p
				return
			***REMOVED***
			s.c.Wait()
		***REMOVED***
	***REMOVED***()
	return c
***REMOVED***

// sortablePeers is a sort wrapper for []api.Peer
type sortablePeers []api.Peer

func (sp sortablePeers) Less(i, j int) bool ***REMOVED*** return sp[i].NodeID < sp[j].NodeID ***REMOVED***

func (sp sortablePeers) Len() int ***REMOVED*** return len(sp) ***REMOVED***

func (sp sortablePeers) Swap(i, j int) ***REMOVED*** sp[i], sp[j] = sp[j], sp[i] ***REMOVED***

// firstSessionErrorTracker is a utility that helps determine whether the agent should exit after
// a TLS failure on establishing the first session.  This should only happen if a join address
// is specified.  If establishing the first session succeeds, but later on some session fails
// because of a TLS error, we don't want to exit the agent because a previously successful
// session indicates that the TLS error may be a transient issue.
type firstSessionErrorTracker struct ***REMOVED***
	mu               sync.Mutex
	pastFirstSession bool
	err              error
***REMOVED***

func (fs *firstSessionErrorTracker) SessionEstablished() ***REMOVED***
	fs.mu.Lock()
	fs.pastFirstSession = true
	fs.mu.Unlock()
***REMOVED***

func (fs *firstSessionErrorTracker) SessionError(err error) ***REMOVED***
	fs.mu.Lock()
	fs.err = err
	fs.mu.Unlock()
***REMOVED***

func (fs *firstSessionErrorTracker) SessionClosed() error ***REMOVED***
	fs.mu.Lock()
	defer fs.mu.Unlock()
	// unfortunately grpc connection errors are type grpc.rpcError, which are not exposed, and we can't get at the underlying error type
	if !fs.pastFirstSession && grpc.Code(fs.err) == codes.Internal &&
		strings.HasPrefix(grpc.ErrorDesc(fs.err), "connection error") && strings.Contains(grpc.ErrorDesc(fs.err), "transport: x509:") ***REMOVED***
		return fs.err
	***REMOVED***
	return nil
***REMOVED***
