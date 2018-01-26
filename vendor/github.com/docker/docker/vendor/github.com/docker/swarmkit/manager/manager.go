package manager

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/cloudflare/cfssl/helpers"
	"github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/go-events"
	gmetrics "github.com/docker/go-metrics"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/ca"
	"github.com/docker/swarmkit/connectionbroker"
	"github.com/docker/swarmkit/identity"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/allocator"
	"github.com/docker/swarmkit/manager/allocator/networkallocator"
	"github.com/docker/swarmkit/manager/controlapi"
	"github.com/docker/swarmkit/manager/dispatcher"
	"github.com/docker/swarmkit/manager/drivers"
	"github.com/docker/swarmkit/manager/health"
	"github.com/docker/swarmkit/manager/keymanager"
	"github.com/docker/swarmkit/manager/logbroker"
	"github.com/docker/swarmkit/manager/metrics"
	"github.com/docker/swarmkit/manager/orchestrator/constraintenforcer"
	"github.com/docker/swarmkit/manager/orchestrator/global"
	"github.com/docker/swarmkit/manager/orchestrator/replicated"
	"github.com/docker/swarmkit/manager/orchestrator/taskreaper"
	"github.com/docker/swarmkit/manager/resourceapi"
	"github.com/docker/swarmkit/manager/scheduler"
	"github.com/docker/swarmkit/manager/state/raft"
	"github.com/docker/swarmkit/manager/state/raft/transport"
	"github.com/docker/swarmkit/manager/state/store"
	"github.com/docker/swarmkit/manager/watchapi"
	"github.com/docker/swarmkit/remotes"
	"github.com/docker/swarmkit/xnet"
	gogotypes "github.com/gogo/protobuf/types"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	// defaultTaskHistoryRetentionLimit is the number of tasks to keep.
	defaultTaskHistoryRetentionLimit = 5
)

// RemoteAddrs provides a listening address and an optional advertise address
// for serving the remote API.
type RemoteAddrs struct ***REMOVED***
	// Address to bind
	ListenAddr string

	// Address to advertise to remote nodes (optional).
	AdvertiseAddr string
***REMOVED***

// Config is used to tune the Manager.
type Config struct ***REMOVED***
	SecurityConfig *ca.SecurityConfig

	// RootCAPaths is the path to which new root certs should be save
	RootCAPaths ca.CertPaths

	// ExternalCAs is a list of initial CAs to which a manager node
	// will make certificate signing requests for node certificates.
	ExternalCAs []*api.ExternalCA

	// ControlAPI is an address for serving the control API.
	ControlAPI string

	// RemoteAPI is a listening address for serving the remote API, and
	// an optional advertise address.
	RemoteAPI *RemoteAddrs

	// JoinRaft is an optional address of a node in an existing raft
	// cluster to join.
	JoinRaft string

	// ForceJoin causes us to invoke raft's Join RPC even if already part
	// of a cluster.
	ForceJoin bool

	// StateDir is the top-level state directory
	StateDir string

	// ForceNewCluster defines if we have to force a new cluster
	// because we are recovering from a backup data directory.
	ForceNewCluster bool

	// ElectionTick defines the amount of ticks needed without
	// leader to trigger a new election
	ElectionTick uint32

	// HeartbeatTick defines the amount of ticks between each
	// heartbeat sent to other members for health-check purposes
	HeartbeatTick uint32

	// AutoLockManagers determines whether or not managers require an unlock key
	// when starting from a stopped state.  This configuration parameter is only
	// applicable when bootstrapping a new cluster for the first time.
	AutoLockManagers bool

	// UnlockKey is the key to unlock a node - used for decrypting manager TLS keys
	// as well as the raft data encryption key (DEK).  It is applicable when
	// bootstrapping a cluster for the first time (it's a cluster-wide setting),
	// and also when loading up any raft data on disk (as a KEK for the raft DEK).
	UnlockKey []byte

	// Availability allows a user to control the current scheduling status of a node
	Availability api.NodeSpec_Availability

	// PluginGetter provides access to docker's plugin inventory.
	PluginGetter plugingetter.PluginGetter
***REMOVED***

// Manager is the cluster manager for Swarm.
// This is the high-level object holding and initializing all the manager
// subsystems.
type Manager struct ***REMOVED***
	config Config

	collector              *metrics.Collector
	caserver               *ca.Server
	dispatcher             *dispatcher.Dispatcher
	logbroker              *logbroker.LogBroker
	watchServer            *watchapi.Server
	replicatedOrchestrator *replicated.Orchestrator
	globalOrchestrator     *global.Orchestrator
	taskReaper             *taskreaper.TaskReaper
	constraintEnforcer     *constraintenforcer.ConstraintEnforcer
	scheduler              *scheduler.Scheduler
	allocator              *allocator.Allocator
	keyManager             *keymanager.KeyManager
	server                 *grpc.Server
	localserver            *grpc.Server
	raftNode               *raft.Node
	dekRotator             *RaftDEKManager
	roleManager            *roleManager

	cancelFunc context.CancelFunc

	// mu is a general mutex used to coordinate starting/stopping and
	// leadership events.
	mu sync.Mutex
	// addrMu is a mutex that protects config.ControlAPI and config.RemoteAPI
	addrMu sync.Mutex

	started chan struct***REMOVED******REMOVED***
	stopped bool

	remoteListener  chan net.Listener
	controlListener chan net.Listener
	errServe        chan error
***REMOVED***

var (
	leaderMetric gmetrics.Gauge
)

func init() ***REMOVED***
	ns := gmetrics.NewNamespace("swarm", "manager", nil)
	leaderMetric = ns.NewGauge("leader", "Indicates if this manager node is a leader", "")
	gmetrics.Register(ns)
***REMOVED***

type closeOnceListener struct ***REMOVED***
	once sync.Once
	net.Listener
***REMOVED***

func (l *closeOnceListener) Close() error ***REMOVED***
	var err error
	l.once.Do(func() ***REMOVED***
		err = l.Listener.Close()
	***REMOVED***)
	return err
***REMOVED***

// New creates a Manager which has not started to accept requests yet.
func New(config *Config) (*Manager, error) ***REMOVED***
	err := os.MkdirAll(config.StateDir, 0700)
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "failed to create state directory")
	***REMOVED***

	raftStateDir := filepath.Join(config.StateDir, "raft")
	err = os.MkdirAll(raftStateDir, 0700)
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "failed to create raft state directory")
	***REMOVED***

	raftCfg := raft.DefaultNodeConfig()

	if config.ElectionTick > 0 ***REMOVED***
		raftCfg.ElectionTick = int(config.ElectionTick)
	***REMOVED***
	if config.HeartbeatTick > 0 ***REMOVED***
		raftCfg.HeartbeatTick = int(config.HeartbeatTick)
	***REMOVED***

	dekRotator, err := NewRaftDEKManager(config.SecurityConfig.KeyWriter())
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	newNodeOpts := raft.NodeOptions***REMOVED***
		ID:              config.SecurityConfig.ClientTLSCreds.NodeID(),
		JoinAddr:        config.JoinRaft,
		ForceJoin:       config.ForceJoin,
		Config:          raftCfg,
		StateDir:        raftStateDir,
		ForceNewCluster: config.ForceNewCluster,
		TLSCredentials:  config.SecurityConfig.ClientTLSCreds,
		KeyRotator:      dekRotator,
	***REMOVED***
	raftNode := raft.NewNode(newNodeOpts)

	opts := []grpc.ServerOption***REMOVED***
		grpc.Creds(config.SecurityConfig.ServerTLSCreds),
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
		grpc.MaxMsgSize(transport.GRPCMaxMsgSize),
	***REMOVED***

	m := &Manager***REMOVED***
		config:          *config,
		caserver:        ca.NewServer(raftNode.MemoryStore(), config.SecurityConfig),
		dispatcher:      dispatcher.New(raftNode, dispatcher.DefaultConfig(), drivers.New(config.PluginGetter), config.SecurityConfig),
		logbroker:       logbroker.New(raftNode.MemoryStore()),
		watchServer:     watchapi.NewServer(raftNode.MemoryStore()),
		server:          grpc.NewServer(opts...),
		localserver:     grpc.NewServer(opts...),
		raftNode:        raftNode,
		started:         make(chan struct***REMOVED******REMOVED***),
		dekRotator:      dekRotator,
		remoteListener:  make(chan net.Listener, 1),
		controlListener: make(chan net.Listener, 1),
		errServe:        make(chan error, 2),
	***REMOVED***

	if config.ControlAPI != "" ***REMOVED***
		m.config.ControlAPI = ""
		if err := m.BindControl(config.ControlAPI); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	if config.RemoteAPI != nil ***REMOVED***
		m.config.RemoteAPI = nil
		// The context isn't used in this case (before (*Manager).Run).
		if err := m.BindRemote(context.Background(), *config.RemoteAPI); err != nil ***REMOVED***
			if config.ControlAPI != "" ***REMOVED***
				l := <-m.controlListener
				l.Close()
			***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return m, nil
***REMOVED***

// BindControl binds a local socket for the control API.
func (m *Manager) BindControl(addr string) error ***REMOVED***
	m.addrMu.Lock()
	defer m.addrMu.Unlock()

	if m.config.ControlAPI != "" ***REMOVED***
		return errors.New("manager already has a control API address")
	***REMOVED***

	// don't create a socket directory if we're on windows. we used named pipe
	if runtime.GOOS != "windows" ***REMOVED***
		err := os.MkdirAll(filepath.Dir(addr), 0700)
		if err != nil ***REMOVED***
			return errors.Wrap(err, "failed to create socket directory")
		***REMOVED***
	***REMOVED***

	l, err := xnet.ListenLocal(addr)

	// A unix socket may fail to bind if the file already
	// exists. Try replacing the file.
	if runtime.GOOS != "windows" ***REMOVED***
		unwrappedErr := err
		if op, ok := unwrappedErr.(*net.OpError); ok ***REMOVED***
			unwrappedErr = op.Err
		***REMOVED***
		if sys, ok := unwrappedErr.(*os.SyscallError); ok ***REMOVED***
			unwrappedErr = sys.Err
		***REMOVED***
		if unwrappedErr == syscall.EADDRINUSE ***REMOVED***
			os.Remove(addr)
			l, err = xnet.ListenLocal(addr)
		***REMOVED***
	***REMOVED***
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed to listen on control API address")
	***REMOVED***

	m.config.ControlAPI = addr
	m.controlListener <- l
	return nil
***REMOVED***

// BindRemote binds a port for the remote API.
func (m *Manager) BindRemote(ctx context.Context, addrs RemoteAddrs) error ***REMOVED***
	m.addrMu.Lock()
	defer m.addrMu.Unlock()

	if m.config.RemoteAPI != nil ***REMOVED***
		return errors.New("manager already has remote API address")
	***REMOVED***

	// If an AdvertiseAddr was specified, we use that as our
	// externally-reachable address.
	advertiseAddr := addrs.AdvertiseAddr

	var advertiseAddrPort string
	if advertiseAddr == "" ***REMOVED***
		// Otherwise, we know we are joining an existing swarm. Use a
		// wildcard address to trigger remote autodetection of our
		// address.
		var err error
		_, advertiseAddrPort, err = net.SplitHostPort(addrs.ListenAddr)
		if err != nil ***REMOVED***
			return fmt.Errorf("missing or invalid listen address %s", addrs.ListenAddr)
		***REMOVED***

		// Even with an IPv6 listening address, it's okay to use
		// 0.0.0.0 here. Any "unspecified" (wildcard) IP will
		// be substituted with the actual source address.
		advertiseAddr = net.JoinHostPort("0.0.0.0", advertiseAddrPort)
	***REMOVED***

	l, err := net.Listen("tcp", addrs.ListenAddr)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed to listen on remote API address")
	***REMOVED***
	if advertiseAddrPort == "0" ***REMOVED***
		advertiseAddr = l.Addr().String()
		addrs.ListenAddr = advertiseAddr
	***REMOVED***

	m.config.RemoteAPI = &addrs

	m.raftNode.SetAddr(ctx, advertiseAddr)
	m.remoteListener <- l

	return nil
***REMOVED***

// RemovedFromRaft returns a channel that's closed if the manager is removed
// from the raft cluster. This should be used to trigger a manager shutdown.
func (m *Manager) RemovedFromRaft() <-chan struct***REMOVED******REMOVED*** ***REMOVED***
	return m.raftNode.RemovedFromRaft
***REMOVED***

// Addr returns tcp address on which remote api listens.
func (m *Manager) Addr() string ***REMOVED***
	m.addrMu.Lock()
	defer m.addrMu.Unlock()

	if m.config.RemoteAPI == nil ***REMOVED***
		return ""
	***REMOVED***
	return m.config.RemoteAPI.ListenAddr
***REMOVED***

// Run starts all manager sub-systems and the gRPC server at the configured
// address.
// The call never returns unless an error occurs or `Stop()` is called.
func (m *Manager) Run(parent context.Context) error ***REMOVED***
	ctx, ctxCancel := context.WithCancel(parent)
	defer ctxCancel()

	m.cancelFunc = ctxCancel

	leadershipCh, cancel := m.raftNode.SubscribeLeadership()
	defer cancel()

	go m.handleLeadershipEvents(ctx, leadershipCh)

	authorize := func(ctx context.Context, roles []string) error ***REMOVED***
		var (
			blacklistedCerts map[string]*api.BlacklistedCertificate
			clusters         []*api.Cluster
			err              error
		)

		m.raftNode.MemoryStore().View(func(readTx store.ReadTx) ***REMOVED***
			clusters, err = store.FindClusters(readTx, store.ByName(store.DefaultClusterName))

		***REMOVED***)

		// Not having a cluster object yet means we can't check
		// the blacklist.
		if err == nil && len(clusters) == 1 ***REMOVED***
			blacklistedCerts = clusters[0].BlacklistedCertificates
		***REMOVED***

		// Authorize the remote roles, ensure they can only be forwarded by managers
		_, err = ca.AuthorizeForwardedRoleAndOrg(ctx, roles, []string***REMOVED***ca.ManagerRole***REMOVED***, m.config.SecurityConfig.ClientTLSCreds.Organization(), blacklistedCerts)
		return err
	***REMOVED***

	baseControlAPI := controlapi.NewServer(m.raftNode.MemoryStore(), m.raftNode, m.config.SecurityConfig, m.config.PluginGetter, drivers.New(m.config.PluginGetter))
	baseResourceAPI := resourceapi.New(m.raftNode.MemoryStore())
	healthServer := health.NewHealthServer()
	localHealthServer := health.NewHealthServer()

	authenticatedControlAPI := api.NewAuthenticatedWrapperControlServer(baseControlAPI, authorize)
	authenticatedWatchAPI := api.NewAuthenticatedWrapperWatchServer(m.watchServer, authorize)
	authenticatedResourceAPI := api.NewAuthenticatedWrapperResourceAllocatorServer(baseResourceAPI, authorize)
	authenticatedLogsServerAPI := api.NewAuthenticatedWrapperLogsServer(m.logbroker, authorize)
	authenticatedLogBrokerAPI := api.NewAuthenticatedWrapperLogBrokerServer(m.logbroker, authorize)
	authenticatedDispatcherAPI := api.NewAuthenticatedWrapperDispatcherServer(m.dispatcher, authorize)
	authenticatedCAAPI := api.NewAuthenticatedWrapperCAServer(m.caserver, authorize)
	authenticatedNodeCAAPI := api.NewAuthenticatedWrapperNodeCAServer(m.caserver, authorize)
	authenticatedRaftAPI := api.NewAuthenticatedWrapperRaftServer(m.raftNode, authorize)
	authenticatedHealthAPI := api.NewAuthenticatedWrapperHealthServer(healthServer, authorize)
	authenticatedRaftMembershipAPI := api.NewAuthenticatedWrapperRaftMembershipServer(m.raftNode, authorize)

	proxyDispatcherAPI := api.NewRaftProxyDispatcherServer(authenticatedDispatcherAPI, m.raftNode, nil, ca.WithMetadataForwardTLSInfo)
	proxyCAAPI := api.NewRaftProxyCAServer(authenticatedCAAPI, m.raftNode, nil, ca.WithMetadataForwardTLSInfo)
	proxyNodeCAAPI := api.NewRaftProxyNodeCAServer(authenticatedNodeCAAPI, m.raftNode, nil, ca.WithMetadataForwardTLSInfo)
	proxyRaftMembershipAPI := api.NewRaftProxyRaftMembershipServer(authenticatedRaftMembershipAPI, m.raftNode, nil, ca.WithMetadataForwardTLSInfo)
	proxyResourceAPI := api.NewRaftProxyResourceAllocatorServer(authenticatedResourceAPI, m.raftNode, nil, ca.WithMetadataForwardTLSInfo)
	proxyLogBrokerAPI := api.NewRaftProxyLogBrokerServer(authenticatedLogBrokerAPI, m.raftNode, nil, ca.WithMetadataForwardTLSInfo)

	// The following local proxies are only wired up to receive requests
	// from a trusted local socket, and these requests don't use TLS,
	// therefore the requests they handle locally should bypass
	// authorization. When requests are proxied from these servers, they
	// are sent as requests from this manager rather than forwarded
	// requests (it has no TLS information to put in the metadata map).
	forwardAsOwnRequest := func(ctx context.Context) (context.Context, error) ***REMOVED*** return ctx, nil ***REMOVED***
	handleRequestLocally := func(ctx context.Context) (context.Context, error) ***REMOVED***
		remoteAddr := "127.0.0.1:0"

		m.addrMu.Lock()
		if m.config.RemoteAPI != nil ***REMOVED***
			if m.config.RemoteAPI.AdvertiseAddr != "" ***REMOVED***
				remoteAddr = m.config.RemoteAPI.AdvertiseAddr
			***REMOVED*** else ***REMOVED***
				remoteAddr = m.config.RemoteAPI.ListenAddr
			***REMOVED***
		***REMOVED***
		m.addrMu.Unlock()

		creds := m.config.SecurityConfig.ClientTLSCreds

		nodeInfo := ca.RemoteNodeInfo***REMOVED***
			Roles:        []string***REMOVED***creds.Role()***REMOVED***,
			Organization: creds.Organization(),
			NodeID:       creds.NodeID(),
			RemoteAddr:   remoteAddr,
		***REMOVED***

		return context.WithValue(ctx, ca.LocalRequestKey, nodeInfo), nil
	***REMOVED***
	localProxyControlAPI := api.NewRaftProxyControlServer(baseControlAPI, m.raftNode, handleRequestLocally, forwardAsOwnRequest)
	localProxyLogsAPI := api.NewRaftProxyLogsServer(m.logbroker, m.raftNode, handleRequestLocally, forwardAsOwnRequest)
	localProxyDispatcherAPI := api.NewRaftProxyDispatcherServer(m.dispatcher, m.raftNode, handleRequestLocally, forwardAsOwnRequest)
	localProxyCAAPI := api.NewRaftProxyCAServer(m.caserver, m.raftNode, handleRequestLocally, forwardAsOwnRequest)
	localProxyNodeCAAPI := api.NewRaftProxyNodeCAServer(m.caserver, m.raftNode, handleRequestLocally, forwardAsOwnRequest)
	localProxyResourceAPI := api.NewRaftProxyResourceAllocatorServer(baseResourceAPI, m.raftNode, handleRequestLocally, forwardAsOwnRequest)
	localProxyLogBrokerAPI := api.NewRaftProxyLogBrokerServer(m.logbroker, m.raftNode, handleRequestLocally, forwardAsOwnRequest)

	// Everything registered on m.server should be an authenticated
	// wrapper, or a proxy wrapping an authenticated wrapper!
	api.RegisterCAServer(m.server, proxyCAAPI)
	api.RegisterNodeCAServer(m.server, proxyNodeCAAPI)
	api.RegisterRaftServer(m.server, authenticatedRaftAPI)
	api.RegisterHealthServer(m.server, authenticatedHealthAPI)
	api.RegisterRaftMembershipServer(m.server, proxyRaftMembershipAPI)
	api.RegisterControlServer(m.server, authenticatedControlAPI)
	api.RegisterWatchServer(m.server, authenticatedWatchAPI)
	api.RegisterLogsServer(m.server, authenticatedLogsServerAPI)
	api.RegisterLogBrokerServer(m.server, proxyLogBrokerAPI)
	api.RegisterResourceAllocatorServer(m.server, proxyResourceAPI)
	api.RegisterDispatcherServer(m.server, proxyDispatcherAPI)
	grpc_prometheus.Register(m.server)

	api.RegisterControlServer(m.localserver, localProxyControlAPI)
	api.RegisterWatchServer(m.localserver, m.watchServer)
	api.RegisterLogsServer(m.localserver, localProxyLogsAPI)
	api.RegisterHealthServer(m.localserver, localHealthServer)
	api.RegisterDispatcherServer(m.localserver, localProxyDispatcherAPI)
	api.RegisterCAServer(m.localserver, localProxyCAAPI)
	api.RegisterNodeCAServer(m.localserver, localProxyNodeCAAPI)
	api.RegisterResourceAllocatorServer(m.localserver, localProxyResourceAPI)
	api.RegisterLogBrokerServer(m.localserver, localProxyLogBrokerAPI)
	grpc_prometheus.Register(m.localserver)

	healthServer.SetServingStatus("Raft", api.HealthCheckResponse_NOT_SERVING)
	localHealthServer.SetServingStatus("ControlAPI", api.HealthCheckResponse_NOT_SERVING)

	if err := m.watchServer.Start(ctx); err != nil ***REMOVED***
		log.G(ctx).WithError(err).Error("watch server failed to start")
	***REMOVED***

	go m.serveListener(ctx, m.remoteListener)
	go m.serveListener(ctx, m.controlListener)

	defer func() ***REMOVED***
		m.server.Stop()
		m.localserver.Stop()
	***REMOVED***()

	// Set the raft server as serving for the health server
	healthServer.SetServingStatus("Raft", api.HealthCheckResponse_SERVING)

	if err := m.raftNode.JoinAndStart(ctx); err != nil ***REMOVED***
		// Don't block future calls to Stop.
		close(m.started)
		return errors.Wrap(err, "can't initialize raft node")
	***REMOVED***

	localHealthServer.SetServingStatus("ControlAPI", api.HealthCheckResponse_SERVING)

	// Start metrics collection.

	m.collector = metrics.NewCollector(m.raftNode.MemoryStore())
	go func(collector *metrics.Collector) ***REMOVED***
		if err := collector.Run(ctx); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("collector failed with an error")
		***REMOVED***
	***REMOVED***(m.collector)

	close(m.started)

	go func() ***REMOVED***
		err := m.raftNode.Run(ctx)
		if err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("raft node stopped")
			m.Stop(ctx, false)
		***REMOVED***
	***REMOVED***()

	if err := raft.WaitForLeader(ctx, m.raftNode); err != nil ***REMOVED***
		return err
	***REMOVED***

	c, err := raft.WaitForCluster(ctx, m.raftNode)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	raftConfig := c.Spec.Raft

	if err := m.watchForClusterChanges(ctx); err != nil ***REMOVED***
		return err
	***REMOVED***

	if int(raftConfig.ElectionTick) != m.raftNode.Config.ElectionTick ***REMOVED***
		log.G(ctx).Warningf("election tick value (%ds) is different from the one defined in the cluster config (%vs), the cluster may be unstable", m.raftNode.Config.ElectionTick, raftConfig.ElectionTick)
	***REMOVED***
	if int(raftConfig.HeartbeatTick) != m.raftNode.Config.HeartbeatTick ***REMOVED***
		log.G(ctx).Warningf("heartbeat tick value (%ds) is different from the one defined in the cluster config (%vs), the cluster may be unstable", m.raftNode.Config.HeartbeatTick, raftConfig.HeartbeatTick)
	***REMOVED***

	// wait for an error in serving.
	err = <-m.errServe
	m.mu.Lock()
	if m.stopped ***REMOVED***
		m.mu.Unlock()
		return nil
	***REMOVED***
	m.mu.Unlock()
	m.Stop(ctx, false)

	return err
***REMOVED***

const stopTimeout = 8 * time.Second

// Stop stops the manager. It immediately closes all open connections and
// active RPCs as well as stopping the manager's subsystems. If clearData is
// set, the raft logs, snapshots, and keys will be erased.
func (m *Manager) Stop(ctx context.Context, clearData bool) ***REMOVED***
	log.G(ctx).Info("Stopping manager")
	// It's not safe to start shutting down while the manager is still
	// starting up.
	<-m.started

	// the mutex stops us from trying to stop while we're already stopping, or
	// from returning before we've finished stopping.
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.stopped ***REMOVED***
		return
	***REMOVED***
	m.stopped = true

	srvDone, localSrvDone := make(chan struct***REMOVED******REMOVED***), make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		m.server.GracefulStop()
		close(srvDone)
	***REMOVED***()
	go func() ***REMOVED***
		m.localserver.GracefulStop()
		close(localSrvDone)
	***REMOVED***()

	m.raftNode.Cancel()

	if m.collector != nil ***REMOVED***
		m.collector.Stop()
	***REMOVED***

	m.dispatcher.Stop()
	m.logbroker.Stop()
	m.watchServer.Stop()
	m.caserver.Stop()

	if m.allocator != nil ***REMOVED***
		m.allocator.Stop()
	***REMOVED***
	if m.replicatedOrchestrator != nil ***REMOVED***
		m.replicatedOrchestrator.Stop()
	***REMOVED***
	if m.globalOrchestrator != nil ***REMOVED***
		m.globalOrchestrator.Stop()
	***REMOVED***
	if m.taskReaper != nil ***REMOVED***
		m.taskReaper.Stop()
	***REMOVED***
	if m.constraintEnforcer != nil ***REMOVED***
		m.constraintEnforcer.Stop()
	***REMOVED***
	if m.scheduler != nil ***REMOVED***
		m.scheduler.Stop()
	***REMOVED***
	if m.roleManager != nil ***REMOVED***
		m.roleManager.Stop()
	***REMOVED***
	if m.keyManager != nil ***REMOVED***
		m.keyManager.Stop()
	***REMOVED***

	if clearData ***REMOVED***
		m.raftNode.ClearData()
	***REMOVED***
	m.cancelFunc()
	<-m.raftNode.Done()

	timer := time.AfterFunc(stopTimeout, func() ***REMOVED***
		m.server.Stop()
		m.localserver.Stop()
	***REMOVED***)
	defer timer.Stop()
	// TODO: we're not waiting on ctx because it very well could be passed from Run,
	// which is already cancelled here. We need to refactor that.
	select ***REMOVED***
	case <-srvDone:
		<-localSrvDone
	case <-localSrvDone:
		<-srvDone
	***REMOVED***

	log.G(ctx).Info("Manager shut down")
	// mutex is released and Run can return now
***REMOVED***

func (m *Manager) updateKEK(ctx context.Context, cluster *api.Cluster) error ***REMOVED***
	securityConfig := m.config.SecurityConfig
	nodeID := m.config.SecurityConfig.ClientTLSCreds.NodeID()
	logger := log.G(ctx).WithFields(logrus.Fields***REMOVED***
		"node.id":   nodeID,
		"node.role": ca.ManagerRole,
	***REMOVED***)

	kekData := ca.KEKData***REMOVED***Version: cluster.Meta.Version.Index***REMOVED***
	for _, encryptionKey := range cluster.UnlockKeys ***REMOVED***
		if encryptionKey.Subsystem == ca.ManagerRole ***REMOVED***
			kekData.KEK = encryptionKey.Key
			break
		***REMOVED***
	***REMOVED***
	updated, unlockedToLocked, err := m.dekRotator.MaybeUpdateKEK(kekData)
	if err != nil ***REMOVED***
		logger.WithError(err).Errorf("failed to re-encrypt TLS key with a new KEK")
		return err
	***REMOVED***
	if updated ***REMOVED***
		logger.Debug("successfully rotated KEK")
	***REMOVED***
	if unlockedToLocked ***REMOVED***
		// a best effort attempt to update the TLS certificate - if it fails, it'll be updated the next time it renews;
		// don't wait because it might take a bit
		go func() ***REMOVED***
			insecureCreds := credentials.NewTLS(&tls.Config***REMOVED***InsecureSkipVerify: true***REMOVED***)

			conn, err := grpc.Dial(
				m.config.ControlAPI,
				grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
				grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor),
				grpc.WithTransportCredentials(insecureCreds),
				grpc.WithDialer(
					func(addr string, timeout time.Duration) (net.Conn, error) ***REMOVED***
						return xnet.DialTimeoutLocal(addr, timeout)
					***REMOVED***),
			)
			if err != nil ***REMOVED***
				logger.WithError(err).Error("failed to connect to local manager socket after locking the cluster")
				return
			***REMOVED***

			defer conn.Close()

			connBroker := connectionbroker.New(remotes.NewRemotes())
			connBroker.SetLocalConn(conn)
			if err := ca.RenewTLSConfigNow(ctx, securityConfig, connBroker, m.config.RootCAPaths); err != nil ***REMOVED***
				logger.WithError(err).Error("failed to download new TLS certificate after locking the cluster")
			***REMOVED***
		***REMOVED***()
	***REMOVED***
	return nil
***REMOVED***

func (m *Manager) watchForClusterChanges(ctx context.Context) error ***REMOVED***
	clusterID := m.config.SecurityConfig.ClientTLSCreds.Organization()
	var cluster *api.Cluster
	clusterWatch, clusterWatchCancel, err := store.ViewAndWatch(m.raftNode.MemoryStore(),
		func(tx store.ReadTx) error ***REMOVED***
			cluster = store.GetCluster(tx, clusterID)
			if cluster == nil ***REMOVED***
				return fmt.Errorf("unable to get current cluster")
			***REMOVED***
			return nil
		***REMOVED***,
		api.EventUpdateCluster***REMOVED***
			Cluster: &api.Cluster***REMOVED***ID: clusterID***REMOVED***,
			Checks:  []api.ClusterCheckFunc***REMOVED***api.ClusterCheckID***REMOVED***,
		***REMOVED***,
	)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := m.updateKEK(ctx, cluster); err != nil ***REMOVED***
		return err
	***REMOVED***

	go func() ***REMOVED***
		for ***REMOVED***
			select ***REMOVED***
			case event := <-clusterWatch:
				clusterEvent := event.(api.EventUpdateCluster)
				m.updateKEK(ctx, clusterEvent.Cluster)
			case <-ctx.Done():
				clusterWatchCancel()
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	return nil
***REMOVED***

// rotateRootCAKEK will attempt to rotate the key-encryption-key for root CA key-material in raft.
// If there is no passphrase set in ENV, it returns.
// If there is plain-text root key-material, and a passphrase set, it encrypts it.
// If there is encrypted root key-material and it is using the current passphrase, it returns.
// If there is encrypted root key-material, and it is using the previous passphrase, it
// re-encrypts it with the current passphrase.
func (m *Manager) rotateRootCAKEK(ctx context.Context, clusterID string) error ***REMOVED***
	// If we don't have a KEK, we won't ever be rotating anything
	strPassphrase := os.Getenv(ca.PassphraseENVVar)
	strPassphrasePrev := os.Getenv(ca.PassphraseENVVarPrev)
	if strPassphrase == "" && strPassphrasePrev == "" ***REMOVED***
		return nil
	***REMOVED***
	if strPassphrase != "" ***REMOVED***
		log.G(ctx).Warn("Encrypting the root CA key in swarm using environment variables is deprecated. " +
			"Support for decrypting or rotating the key will be removed in the future.")
	***REMOVED***

	passphrase := []byte(strPassphrase)
	passphrasePrev := []byte(strPassphrasePrev)

	s := m.raftNode.MemoryStore()
	var (
		cluster  *api.Cluster
		err      error
		finalKey []byte
	)
	// Retrieve the cluster identified by ClusterID
	return s.Update(func(tx store.Tx) error ***REMOVED***
		cluster = store.GetCluster(tx, clusterID)
		if cluster == nil ***REMOVED***
			return fmt.Errorf("cluster not found: %s", clusterID)
		***REMOVED***

		// Try to get the private key from the cluster
		privKeyPEM := cluster.RootCA.CAKey
		if len(privKeyPEM) == 0 ***REMOVED***
			// We have no PEM root private key in this cluster.
			log.G(ctx).Warnf("cluster %s does not have private key material", clusterID)
			return nil
		***REMOVED***

		// Decode the PEM private key
		keyBlock, _ := pem.Decode(privKeyPEM)
		if keyBlock == nil ***REMOVED***
			return fmt.Errorf("invalid PEM-encoded private key inside of cluster %s", clusterID)
		***REMOVED***

		if x509.IsEncryptedPEMBlock(keyBlock) ***REMOVED***
			// PEM encryption does not have a digest, so sometimes decryption doesn't
			// error even with the wrong passphrase.  So actually try to parse it into a valid key.
			_, err := helpers.ParsePrivateKeyPEMWithPassword(privKeyPEM, []byte(passphrase))
			if err == nil ***REMOVED***
				// This key is already correctly encrypted with the correct KEK, nothing to do here
				return nil
			***REMOVED***

			// This key is already encrypted, but failed with current main passphrase.
			// Let's try to decrypt with the previous passphrase, and parse into a valid key, for the
			// same reason as above.
			_, err = helpers.ParsePrivateKeyPEMWithPassword(privKeyPEM, []byte(passphrasePrev))
			if err != nil ***REMOVED***
				// We were not able to decrypt either with the main or backup passphrase, error
				return err
			***REMOVED***
			// ok the above passphrase is correct, so decrypt the PEM block so we can re-encrypt -
			// since the key was successfully decrypted above, there will be no error doing PEM
			// decryption
			unencryptedDER, _ := x509.DecryptPEMBlock(keyBlock, []byte(passphrasePrev))
			unencryptedKeyBlock := &pem.Block***REMOVED***
				Type:  keyBlock.Type,
				Bytes: unencryptedDER,
			***REMOVED***

			// we were able to decrypt the key with the previous passphrase - if the current passphrase is empty,
			// the we store the decrypted key in raft
			finalKey = pem.EncodeToMemory(unencryptedKeyBlock)

			// the current passphrase is not empty, so let's encrypt with the new one and store it in raft
			if strPassphrase != "" ***REMOVED***
				finalKey, err = ca.EncryptECPrivateKey(finalKey, strPassphrase)
				if err != nil ***REMOVED***
					log.G(ctx).WithError(err).Debugf("failed to rotate the key-encrypting-key for the root key material of cluster %s", clusterID)
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED*** else if strPassphrase != "" ***REMOVED***
			// If this key is not encrypted, and the passphrase is not nil, then we have to encrypt it
			finalKey, err = ca.EncryptECPrivateKey(privKeyPEM, strPassphrase)
			if err != nil ***REMOVED***
				log.G(ctx).WithError(err).Debugf("failed to rotate the key-encrypting-key for the root key material of cluster %s", clusterID)
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			return nil // don't update if it's not encrypted and we don't want it encrypted
		***REMOVED***

		log.G(ctx).Infof("Updating the encryption on the root key material of cluster %s", clusterID)
		cluster.RootCA.CAKey = finalKey
		return store.UpdateCluster(tx, cluster)
	***REMOVED***)
***REMOVED***

// handleLeadershipEvents handles the is leader event or is follower event.
func (m *Manager) handleLeadershipEvents(ctx context.Context, leadershipCh chan events.Event) ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case leadershipEvent := <-leadershipCh:
			m.mu.Lock()
			if m.stopped ***REMOVED***
				m.mu.Unlock()
				return
			***REMOVED***
			newState := leadershipEvent.(raft.LeadershipState)

			if newState == raft.IsLeader ***REMOVED***
				m.becomeLeader(ctx)
				leaderMetric.Set(1)
			***REMOVED*** else if newState == raft.IsFollower ***REMOVED***
				m.becomeFollower()
				leaderMetric.Set(0)
			***REMOVED***
			m.mu.Unlock()
		case <-ctx.Done():
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// serveListener serves a listener for local and non local connections.
func (m *Manager) serveListener(ctx context.Context, lCh <-chan net.Listener) ***REMOVED***
	var l net.Listener
	select ***REMOVED***
	case l = <-lCh:
	case <-ctx.Done():
		return
	***REMOVED***
	ctx = log.WithLogger(ctx, log.G(ctx).WithFields(
		logrus.Fields***REMOVED***
			"proto": l.Addr().Network(),
			"addr":  l.Addr().String(),
		***REMOVED***))
	if _, ok := l.(*net.TCPListener); !ok ***REMOVED***
		log.G(ctx).Info("Listening for local connections")
		// we need to disallow double closes because UnixListener.Close
		// can delete unix-socket file of newer listener. grpc calls
		// Close twice indeed: in Serve and in Stop.
		m.errServe <- m.localserver.Serve(&closeOnceListener***REMOVED***Listener: l***REMOVED***)
	***REMOVED*** else ***REMOVED***
		log.G(ctx).Info("Listening for connections")
		m.errServe <- m.server.Serve(l)
	***REMOVED***
***REMOVED***

// becomeLeader starts the subsystems that are run on the leader.
func (m *Manager) becomeLeader(ctx context.Context) ***REMOVED***
	s := m.raftNode.MemoryStore()

	rootCA := m.config.SecurityConfig.RootCA()
	nodeID := m.config.SecurityConfig.ClientTLSCreds.NodeID()

	raftCfg := raft.DefaultRaftConfig()
	raftCfg.ElectionTick = uint32(m.raftNode.Config.ElectionTick)
	raftCfg.HeartbeatTick = uint32(m.raftNode.Config.HeartbeatTick)

	clusterID := m.config.SecurityConfig.ClientTLSCreds.Organization()

	initialCAConfig := ca.DefaultCAConfig()
	initialCAConfig.ExternalCAs = m.config.ExternalCAs

	var unlockKeys []*api.EncryptionKey
	if m.config.AutoLockManagers ***REMOVED***
		unlockKeys = []*api.EncryptionKey***REMOVED******REMOVED***
			Subsystem: ca.ManagerRole,
			Key:       m.config.UnlockKey,
		***REMOVED******REMOVED***
	***REMOVED***

	s.Update(func(tx store.Tx) error ***REMOVED***
		// Add a default cluster object to the
		// store. Don't check the error because
		// we expect this to fail unless this
		// is a brand new cluster.
		err := store.CreateCluster(tx, defaultClusterObject(
			clusterID,
			initialCAConfig,
			raftCfg,
			api.EncryptionConfig***REMOVED***AutoLockManagers: m.config.AutoLockManagers***REMOVED***,
			unlockKeys,
			rootCA))

		if err != nil && err != store.ErrExist ***REMOVED***
			log.G(ctx).WithError(err).Errorf("error creating cluster object")
		***REMOVED***

		// Add Node entry for ourself, if one
		// doesn't exist already.
		freshCluster := nil == store.CreateNode(tx, managerNode(nodeID, m.config.Availability))

		if freshCluster ***REMOVED***
			// This is a fresh swarm cluster. Add to store now any initial
			// cluster resource, like the default ingress network which
			// provides the routing mesh for this cluster.
			log.G(ctx).Info("Creating default ingress network")
			if err := store.CreateNetwork(tx, newIngressNetwork()); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Error("failed to create default ingress network")
			***REMOVED***
		***REMOVED***
		// Create now the static predefined if the store does not contain predefined
		//networks like bridge/host node-local networks which
		// are known to be present in each cluster node. This is needed
		// in order to allow running services on the predefined docker
		// networks like `bridge` and `host`.
		for _, p := range allocator.PredefinedNetworks() ***REMOVED***
			if store.GetNetwork(tx, p.Name) == nil ***REMOVED***
				if err := store.CreateNetwork(tx, newPredefinedNetwork(p.Name, p.Driver)); err != nil ***REMOVED***
					log.G(ctx).WithError(err).Error("failed to create predefined network " + p.Name)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***)

	// Attempt to rotate the key-encrypting-key of the root CA key-material
	err := m.rotateRootCAKEK(ctx, clusterID)
	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Error("root key-encrypting-key rotation failed")
	***REMOVED***

	m.replicatedOrchestrator = replicated.NewReplicatedOrchestrator(s)
	m.constraintEnforcer = constraintenforcer.New(s)
	m.globalOrchestrator = global.NewGlobalOrchestrator(s)
	m.taskReaper = taskreaper.New(s)
	m.scheduler = scheduler.New(s)
	m.keyManager = keymanager.New(s, keymanager.DefaultConfig())
	m.roleManager = newRoleManager(s, m.raftNode)

	// TODO(stevvooe): Allocate a context that can be used to
	// shutdown underlying manager processes when leadership is
	// lost.

	m.allocator, err = allocator.New(s, m.config.PluginGetter)
	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Error("failed to create allocator")
		// TODO(stevvooe): It doesn't seem correct here to fail
		// creating the allocator but then use it anyway.
	***REMOVED***

	if m.keyManager != nil ***REMOVED***
		go func(keyManager *keymanager.KeyManager) ***REMOVED***
			if err := keyManager.Run(ctx); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Error("keymanager failed with an error")
			***REMOVED***
		***REMOVED***(m.keyManager)
	***REMOVED***

	go func(d *dispatcher.Dispatcher) ***REMOVED***
		if err := d.Run(ctx); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("Dispatcher exited with an error")
		***REMOVED***
	***REMOVED***(m.dispatcher)

	if err := m.logbroker.Start(ctx); err != nil ***REMOVED***
		log.G(ctx).WithError(err).Error("LogBroker failed to start")
	***REMOVED***

	go func(server *ca.Server) ***REMOVED***
		if err := server.Run(ctx); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("CA signer exited with an error")
		***REMOVED***
	***REMOVED***(m.caserver)

	// Start all sub-components in separate goroutines.
	// TODO(aluzzardi): This should have some kind of error handling so that
	// any component that goes down would bring the entire manager down.
	if m.allocator != nil ***REMOVED***
		go func(allocator *allocator.Allocator) ***REMOVED***
			if err := allocator.Run(ctx); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Error("allocator exited with an error")
			***REMOVED***
		***REMOVED***(m.allocator)
	***REMOVED***

	go func(scheduler *scheduler.Scheduler) ***REMOVED***
		if err := scheduler.Run(ctx); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("scheduler exited with an error")
		***REMOVED***
	***REMOVED***(m.scheduler)

	go func(constraintEnforcer *constraintenforcer.ConstraintEnforcer) ***REMOVED***
		constraintEnforcer.Run()
	***REMOVED***(m.constraintEnforcer)

	go func(taskReaper *taskreaper.TaskReaper) ***REMOVED***
		taskReaper.Run(ctx)
	***REMOVED***(m.taskReaper)

	go func(orchestrator *replicated.Orchestrator) ***REMOVED***
		if err := orchestrator.Run(ctx); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("replicated orchestrator exited with an error")
		***REMOVED***
	***REMOVED***(m.replicatedOrchestrator)

	go func(globalOrchestrator *global.Orchestrator) ***REMOVED***
		if err := globalOrchestrator.Run(ctx); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("global orchestrator exited with an error")
		***REMOVED***
	***REMOVED***(m.globalOrchestrator)

	go func(roleManager *roleManager) ***REMOVED***
		roleManager.Run(ctx)
	***REMOVED***(m.roleManager)
***REMOVED***

// becomeFollower shuts down the subsystems that are only run by the leader.
func (m *Manager) becomeFollower() ***REMOVED***
	m.dispatcher.Stop()
	m.logbroker.Stop()
	m.caserver.Stop()

	if m.allocator != nil ***REMOVED***
		m.allocator.Stop()
		m.allocator = nil
	***REMOVED***

	m.constraintEnforcer.Stop()
	m.constraintEnforcer = nil

	m.replicatedOrchestrator.Stop()
	m.replicatedOrchestrator = nil

	m.globalOrchestrator.Stop()
	m.globalOrchestrator = nil

	m.taskReaper.Stop()
	m.taskReaper = nil

	m.scheduler.Stop()
	m.scheduler = nil

	m.roleManager.Stop()
	m.roleManager = nil

	if m.keyManager != nil ***REMOVED***
		m.keyManager.Stop()
		m.keyManager = nil
	***REMOVED***
***REMOVED***

// defaultClusterObject creates a default cluster.
func defaultClusterObject(
	clusterID string,
	initialCAConfig api.CAConfig,
	raftCfg api.RaftConfig,
	encryptionConfig api.EncryptionConfig,
	initialUnlockKeys []*api.EncryptionKey,
	rootCA *ca.RootCA) *api.Cluster ***REMOVED***
	var caKey []byte
	if rcaSigner, err := rootCA.Signer(); err == nil ***REMOVED***
		caKey = rcaSigner.Key
	***REMOVED***

	return &api.Cluster***REMOVED***
		ID: clusterID,
		Spec: api.ClusterSpec***REMOVED***
			Annotations: api.Annotations***REMOVED***
				Name: store.DefaultClusterName,
			***REMOVED***,
			Orchestration: api.OrchestrationConfig***REMOVED***
				TaskHistoryRetentionLimit: defaultTaskHistoryRetentionLimit,
			***REMOVED***,
			Dispatcher: api.DispatcherConfig***REMOVED***
				HeartbeatPeriod: gogotypes.DurationProto(dispatcher.DefaultHeartBeatPeriod),
			***REMOVED***,
			Raft:             raftCfg,
			CAConfig:         initialCAConfig,
			EncryptionConfig: encryptionConfig,
		***REMOVED***,
		RootCA: api.RootCA***REMOVED***
			CAKey:      caKey,
			CACert:     rootCA.Certs,
			CACertHash: rootCA.Digest.String(),
			JoinTokens: api.JoinTokens***REMOVED***
				Worker:  ca.GenerateJoinToken(rootCA),
				Manager: ca.GenerateJoinToken(rootCA),
			***REMOVED***,
		***REMOVED***,
		UnlockKeys: initialUnlockKeys,
	***REMOVED***
***REMOVED***

// managerNode creates a new node with NodeRoleManager role.
func managerNode(nodeID string, availability api.NodeSpec_Availability) *api.Node ***REMOVED***
	return &api.Node***REMOVED***
		ID: nodeID,
		Certificate: api.Certificate***REMOVED***
			CN:   nodeID,
			Role: api.NodeRoleManager,
			Status: api.IssuanceStatus***REMOVED***
				State: api.IssuanceStateIssued,
			***REMOVED***,
		***REMOVED***,
		Spec: api.NodeSpec***REMOVED***
			DesiredRole:  api.NodeRoleManager,
			Membership:   api.NodeMembershipAccepted,
			Availability: availability,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// newIngressNetwork returns the network object for the default ingress
// network, the network which provides the routing mesh. Caller will save to
// store this object once, at fresh cluster creation. It is expected to
// call this function inside a store update transaction.
func newIngressNetwork() *api.Network ***REMOVED***
	return &api.Network***REMOVED***
		ID: identity.NewID(),
		Spec: api.NetworkSpec***REMOVED***
			Ingress: true,
			Annotations: api.Annotations***REMOVED***
				Name: "ingress",
			***REMOVED***,
			DriverConfig: &api.Driver***REMOVED******REMOVED***,
			IPAM: &api.IPAMOptions***REMOVED***
				Driver: &api.Driver***REMOVED******REMOVED***,
				Configs: []*api.IPAMConfig***REMOVED***
					***REMOVED***
						Subnet: "10.255.0.0/16",
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// Creates a network object representing one of the predefined networks
// known to be statically created on the cluster nodes. These objects
// are populated in the store at cluster creation solely in order to
// support running services on the nodes' predefined networks.
// External clients can filter these predefined networks by looking
// at the predefined label.
func newPredefinedNetwork(name, driver string) *api.Network ***REMOVED***
	return &api.Network***REMOVED***
		ID: identity.NewID(),
		Spec: api.NetworkSpec***REMOVED***
			Annotations: api.Annotations***REMOVED***
				Name: name,
				Labels: map[string]string***REMOVED***
					networkallocator.PredefinedLabel: "true",
				***REMOVED***,
			***REMOVED***,
			DriverConfig: &api.Driver***REMOVED***Name: driver***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***
