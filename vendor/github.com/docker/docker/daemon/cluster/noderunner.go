package cluster

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	types "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/daemon/cluster/executor/container"
	lncluster "github.com/docker/libnetwork/cluster"
	swarmapi "github.com/docker/swarmkit/api"
	swarmnode "github.com/docker/swarmkit/node"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// nodeRunner implements a manager for continuously running swarmkit node, restarting them with backoff delays if needed.
type nodeRunner struct ***REMOVED***
	nodeState
	mu             sync.RWMutex
	done           chan struct***REMOVED******REMOVED*** // closed when swarmNode exits
	ready          chan struct***REMOVED******REMOVED*** // closed when swarmNode becomes active
	reconnectDelay time.Duration
	config         nodeStartConfig

	repeatedRun     bool
	cancelReconnect func()
	stopping        bool
	cluster         *Cluster // only for accessing config helpers, never call any methods. TODO: change to config struct
***REMOVED***

// nodeStartConfig holds configuration needed to start a new node. Exported
// fields of this structure are saved to disk in json. Unexported fields
// contain data that shouldn't be persisted between daemon reloads.
type nodeStartConfig struct ***REMOVED***
	// LocalAddr is this machine's local IP or hostname, if specified.
	LocalAddr string
	// RemoteAddr is the address that was given to "swarm join". It is used
	// to find LocalAddr if necessary.
	RemoteAddr string
	// ListenAddr is the address we bind to, including a port.
	ListenAddr string
	// AdvertiseAddr is the address other nodes should connect to,
	// including a port.
	AdvertiseAddr string
	// DataPathAddr is the address that has to be used for the data path
	DataPathAddr string
	// JoinInProgress is set to true if a join operation has started, but
	// not completed yet.
	JoinInProgress bool

	joinAddr        string
	forceNewCluster bool
	joinToken       string
	lockKey         []byte
	autolock        bool
	availability    types.NodeAvailability
***REMOVED***

func (n *nodeRunner) Ready() chan error ***REMOVED***
	c := make(chan error, 1)
	n.mu.RLock()
	ready, done := n.ready, n.done
	n.mu.RUnlock()
	go func() ***REMOVED***
		select ***REMOVED***
		case <-ready:
		case <-done:
		***REMOVED***
		select ***REMOVED***
		case <-ready:
		default:
			n.mu.RLock()
			c <- n.err
			n.mu.RUnlock()
		***REMOVED***
		close(c)
	***REMOVED***()
	return c
***REMOVED***

func (n *nodeRunner) Start(conf nodeStartConfig) error ***REMOVED***
	n.mu.Lock()
	defer n.mu.Unlock()

	n.reconnectDelay = initialReconnectDelay

	return n.start(conf)
***REMOVED***

func (n *nodeRunner) start(conf nodeStartConfig) error ***REMOVED***
	var control string
	if runtime.GOOS == "windows" ***REMOVED***
		control = `\\.\pipe\` + controlSocket
	***REMOVED*** else ***REMOVED***
		control = filepath.Join(n.cluster.runtimeRoot, controlSocket)
	***REMOVED***

	joinAddr := conf.joinAddr
	if joinAddr == "" && conf.JoinInProgress ***REMOVED***
		// We must have been restarted while trying to join a cluster.
		// Continue trying to join instead of forming our own cluster.
		joinAddr = conf.RemoteAddr
	***REMOVED***

	// Hostname is not set here. Instead, it is obtained from
	// the node description that is reported periodically
	swarmnodeConfig := swarmnode.Config***REMOVED***
		ForceNewCluster:    conf.forceNewCluster,
		ListenControlAPI:   control,
		ListenRemoteAPI:    conf.ListenAddr,
		AdvertiseRemoteAPI: conf.AdvertiseAddr,
		JoinAddr:           joinAddr,
		StateDir:           n.cluster.root,
		JoinToken:          conf.joinToken,
		Executor:           container.NewExecutor(n.cluster.config.Backend, n.cluster.config.PluginBackend),
		HeartbeatTick:      1,
		ElectionTick:       3,
		UnlockKey:          conf.lockKey,
		AutoLockManagers:   conf.autolock,
		PluginGetter:       n.cluster.config.Backend.PluginGetter(),
	***REMOVED***
	if conf.availability != "" ***REMOVED***
		avail, ok := swarmapi.NodeSpec_Availability_value[strings.ToUpper(string(conf.availability))]
		if !ok ***REMOVED***
			return fmt.Errorf("invalid Availability: %q", conf.availability)
		***REMOVED***
		swarmnodeConfig.Availability = swarmapi.NodeSpec_Availability(avail)
	***REMOVED***
	node, err := swarmnode.New(&swarmnodeConfig)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := node.Start(context.Background()); err != nil ***REMOVED***
		return err
	***REMOVED***

	n.done = make(chan struct***REMOVED******REMOVED***)
	n.ready = make(chan struct***REMOVED******REMOVED***)
	n.swarmNode = node
	if conf.joinAddr != "" ***REMOVED***
		conf.JoinInProgress = true
	***REMOVED***
	n.config = conf
	savePersistentState(n.cluster.root, conf)

	ctx, cancel := context.WithCancel(context.Background())

	go func() ***REMOVED***
		n.handleNodeExit(node)
		cancel()
	***REMOVED***()

	go n.handleReadyEvent(ctx, node, n.ready)
	go n.handleControlSocketChange(ctx, node)

	return nil
***REMOVED***

func (n *nodeRunner) handleControlSocketChange(ctx context.Context, node *swarmnode.Node) ***REMOVED***
	for conn := range node.ListenControlSocket(ctx) ***REMOVED***
		n.mu.Lock()
		if n.grpcConn != conn ***REMOVED***
			if conn == nil ***REMOVED***
				n.controlClient = nil
				n.logsClient = nil
			***REMOVED*** else ***REMOVED***
				n.controlClient = swarmapi.NewControlClient(conn)
				n.logsClient = swarmapi.NewLogsClient(conn)
				// push store changes to daemon
				go n.watchClusterEvents(ctx, conn)
			***REMOVED***
		***REMOVED***
		n.grpcConn = conn
		n.mu.Unlock()
		n.cluster.SendClusterEvent(lncluster.EventSocketChange)
	***REMOVED***
***REMOVED***

func (n *nodeRunner) watchClusterEvents(ctx context.Context, conn *grpc.ClientConn) ***REMOVED***
	client := swarmapi.NewWatchClient(conn)
	watch, err := client.Watch(ctx, &swarmapi.WatchRequest***REMOVED***
		Entries: []*swarmapi.WatchRequest_WatchEntry***REMOVED***
			***REMOVED***
				Kind:   "node",
				Action: swarmapi.WatchActionKindCreate | swarmapi.WatchActionKindUpdate | swarmapi.WatchActionKindRemove,
			***REMOVED***,
			***REMOVED***
				Kind:   "service",
				Action: swarmapi.WatchActionKindCreate | swarmapi.WatchActionKindUpdate | swarmapi.WatchActionKindRemove,
			***REMOVED***,
			***REMOVED***
				Kind:   "network",
				Action: swarmapi.WatchActionKindCreate | swarmapi.WatchActionKindUpdate | swarmapi.WatchActionKindRemove,
			***REMOVED***,
			***REMOVED***
				Kind:   "secret",
				Action: swarmapi.WatchActionKindCreate | swarmapi.WatchActionKindUpdate | swarmapi.WatchActionKindRemove,
			***REMOVED***,
			***REMOVED***
				Kind:   "config",
				Action: swarmapi.WatchActionKindCreate | swarmapi.WatchActionKindUpdate | swarmapi.WatchActionKindRemove,
			***REMOVED***,
		***REMOVED***,
		IncludeOldObject: true,
	***REMOVED***)
	if err != nil ***REMOVED***
		logrus.WithError(err).Error("failed to watch cluster store")
		return
	***REMOVED***
	for ***REMOVED***
		msg, err := watch.Recv()
		if err != nil ***REMOVED***
			// store watch is broken
			errStatus, ok := status.FromError(err)
			if !ok || errStatus.Code() != codes.Canceled ***REMOVED***
				logrus.WithError(err).Error("failed to receive changes from store watch API")
			***REMOVED***
			return
		***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return
		case n.cluster.watchStream <- msg:
		***REMOVED***
	***REMOVED***
***REMOVED***

func (n *nodeRunner) handleReadyEvent(ctx context.Context, node *swarmnode.Node, ready chan struct***REMOVED******REMOVED***) ***REMOVED***
	select ***REMOVED***
	case <-node.Ready():
		n.mu.Lock()
		n.err = nil
		if n.config.JoinInProgress ***REMOVED***
			n.config.JoinInProgress = false
			savePersistentState(n.cluster.root, n.config)
		***REMOVED***
		n.mu.Unlock()
		close(ready)
	case <-ctx.Done():
	***REMOVED***
	n.cluster.SendClusterEvent(lncluster.EventNodeReady)
***REMOVED***

func (n *nodeRunner) handleNodeExit(node *swarmnode.Node) ***REMOVED***
	err := detectLockedError(node.Err(context.Background()))
	if err != nil ***REMOVED***
		logrus.Errorf("cluster exited with error: %v", err)
	***REMOVED***
	n.mu.Lock()
	n.swarmNode = nil
	n.err = err
	close(n.done)
	select ***REMOVED***
	case <-n.ready:
		n.enableReconnectWatcher()
	default:
		if n.repeatedRun ***REMOVED***
			n.enableReconnectWatcher()
		***REMOVED***
	***REMOVED***
	n.repeatedRun = true
	n.mu.Unlock()
***REMOVED***

// Stop stops the current swarm node if it is running.
func (n *nodeRunner) Stop() error ***REMOVED***
	n.mu.Lock()
	if n.cancelReconnect != nil ***REMOVED*** // between restarts
		n.cancelReconnect()
		n.cancelReconnect = nil
	***REMOVED***
	if n.swarmNode == nil ***REMOVED***
		n.mu.Unlock()
		return nil
	***REMOVED***
	n.stopping = true
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	n.mu.Unlock()
	if err := n.swarmNode.Stop(ctx); err != nil && !strings.Contains(err.Error(), "context canceled") ***REMOVED***
		return err
	***REMOVED***
	n.cluster.SendClusterEvent(lncluster.EventNodeLeave)
	<-n.done
	return nil
***REMOVED***

func (n *nodeRunner) State() nodeState ***REMOVED***
	if n == nil ***REMOVED***
		return nodeState***REMOVED***status: types.LocalNodeStateInactive***REMOVED***
	***REMOVED***
	n.mu.RLock()
	defer n.mu.RUnlock()

	ns := n.nodeState

	if ns.err != nil || n.cancelReconnect != nil ***REMOVED***
		if errors.Cause(ns.err) == errSwarmLocked ***REMOVED***
			ns.status = types.LocalNodeStateLocked
		***REMOVED*** else ***REMOVED***
			ns.status = types.LocalNodeStateError
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		select ***REMOVED***
		case <-n.ready:
			ns.status = types.LocalNodeStateActive
		default:
			ns.status = types.LocalNodeStatePending
		***REMOVED***
	***REMOVED***

	return ns
***REMOVED***

func (n *nodeRunner) enableReconnectWatcher() ***REMOVED***
	if n.stopping ***REMOVED***
		return
	***REMOVED***
	n.reconnectDelay *= 2
	if n.reconnectDelay > maxReconnectDelay ***REMOVED***
		n.reconnectDelay = maxReconnectDelay
	***REMOVED***
	logrus.Warnf("Restarting swarm in %.2f seconds", n.reconnectDelay.Seconds())
	delayCtx, cancel := context.WithTimeout(context.Background(), n.reconnectDelay)
	n.cancelReconnect = cancel

	go func() ***REMOVED***
		<-delayCtx.Done()
		if delayCtx.Err() != context.DeadlineExceeded ***REMOVED***
			return
		***REMOVED***
		n.mu.Lock()
		defer n.mu.Unlock()
		if n.stopping ***REMOVED***
			return
		***REMOVED***

		if err := n.start(n.config); err != nil ***REMOVED***
			n.err = err
		***REMOVED***
	***REMOVED***()
***REMOVED***

// nodeState represents information about the current state of the cluster and
// provides access to the grpc clients.
type nodeState struct ***REMOVED***
	swarmNode       *swarmnode.Node
	grpcConn        *grpc.ClientConn
	controlClient   swarmapi.ControlClient
	logsClient      swarmapi.LogsClient
	status          types.LocalNodeState
	actualLocalAddr string
	err             error
***REMOVED***

// IsActiveManager returns true if node is a manager ready to accept control requests. It is safe to access the client properties if this returns true.
func (ns nodeState) IsActiveManager() bool ***REMOVED***
	return ns.controlClient != nil
***REMOVED***

// IsManager returns true if node is a manager.
func (ns nodeState) IsManager() bool ***REMOVED***
	return ns.swarmNode != nil && ns.swarmNode.Manager() != nil
***REMOVED***

// NodeID returns node's ID or empty string if node is inactive.
func (ns nodeState) NodeID() string ***REMOVED***
	if ns.swarmNode != nil ***REMOVED***
		return ns.swarmNode.NodeID()
	***REMOVED***
	return ""
***REMOVED***
