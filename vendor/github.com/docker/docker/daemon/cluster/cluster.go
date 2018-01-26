package cluster

//
// ## Swarmkit integration
//
// Cluster - static configurable object for accessing everything swarm related.
// Contains methods for connecting and controlling the cluster. Exists always,
// even if swarm mode is not enabled.
//
// NodeRunner - Manager for starting the swarmkit node. Is present only and
// always if swarm mode is enabled. Implements backoff restart loop in case of
// errors.
//
// NodeState - Information about the current node status including access to
// gRPC clients if a manager is active.
//
// ### Locking
//
// `cluster.controlMutex` - taken for the whole lifecycle of the processes that
// can reconfigure cluster(init/join/leave etc). Protects that one
// reconfiguration action has fully completed before another can start.
//
// `cluster.mu` - taken when the actual changes in cluster configurations
// happen. Different from `controlMutex` because in some cases we need to
// access current cluster state even if the long-running reconfiguration is
// going on. For example network stack may ask for the current cluster state in
// the middle of the shutdown. Any time current cluster state is asked you
// should take the read lock of `cluster.mu`. If you are writing an API
// responder that returns synchronously, hold `cluster.mu.RLock()` for the
// duration of the whole handler function. That ensures that node will not be
// shut down until the handler has finished.
//
// NodeRunner implements its internal locks that should not be used outside of
// the struct. Instead, you should just call `nodeRunner.State()` method to get
// the current state of the cluster(still need `cluster.mu.RLock()` to access
// `cluster.nr` reference itself). Most of the changes in NodeRunner happen
// because of an external event(network problem, unexpected swarmkit error) and
// Docker shouldn't take any locks that delay these changes from happening.
//

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/docker/docker/api/types/network"
	types "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/daemon/cluster/controllers/plugin"
	executorpkg "github.com/docker/docker/daemon/cluster/executor"
	"github.com/docker/docker/pkg/signal"
	lncluster "github.com/docker/libnetwork/cluster"
	swarmapi "github.com/docker/swarmkit/api"
	swarmnode "github.com/docker/swarmkit/node"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const swarmDirName = "swarm"
const controlSocket = "control.sock"
const swarmConnectTimeout = 20 * time.Second
const swarmRequestTimeout = 20 * time.Second
const stateFile = "docker-state.json"
const defaultAddr = "0.0.0.0:2377"

const (
	initialReconnectDelay = 100 * time.Millisecond
	maxReconnectDelay     = 30 * time.Second
	contextPrefix         = "com.docker.swarm"
)

// NetworkSubnetsProvider exposes functions for retrieving the subnets
// of networks managed by Docker, so they can be filtered.
type NetworkSubnetsProvider interface ***REMOVED***
	Subnets() ([]net.IPNet, []net.IPNet)
***REMOVED***

// Config provides values for Cluster.
type Config struct ***REMOVED***
	Root                   string
	Name                   string
	Backend                executorpkg.Backend
	PluginBackend          plugin.Backend
	NetworkSubnetsProvider NetworkSubnetsProvider

	// DefaultAdvertiseAddr is the default host/IP or network interface to use
	// if no AdvertiseAddr value is specified.
	DefaultAdvertiseAddr string

	// path to store runtime state, such as the swarm control socket
	RuntimeRoot string

	// WatchStream is a channel to pass watch API notifications to daemon
	WatchStream chan *swarmapi.WatchMessage
***REMOVED***

// Cluster provides capabilities to participate in a cluster as a worker or a
// manager.
type Cluster struct ***REMOVED***
	mu           sync.RWMutex
	controlMutex sync.RWMutex // protect init/join/leave user operations
	nr           *nodeRunner
	root         string
	runtimeRoot  string
	config       Config
	configEvent  chan lncluster.ConfigEventType // todo: make this array and goroutine safe
	attachers    map[string]*attacher
	watchStream  chan *swarmapi.WatchMessage
***REMOVED***

// attacher manages the in-memory attachment state of a container
// attachment to a global scope network managed by swarm manager. It
// helps in identifying the attachment ID via the taskID and the
// corresponding attachment configuration obtained from the manager.
type attacher struct ***REMOVED***
	taskID           string
	config           *network.NetworkingConfig
	inProgress       bool
	attachWaitCh     chan *network.NetworkingConfig
	attachCompleteCh chan struct***REMOVED******REMOVED***
	detachWaitCh     chan struct***REMOVED******REMOVED***
***REMOVED***

// New creates a new Cluster instance using provided config.
func New(config Config) (*Cluster, error) ***REMOVED***
	root := filepath.Join(config.Root, swarmDirName)
	if err := os.MkdirAll(root, 0700); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if config.RuntimeRoot == "" ***REMOVED***
		config.RuntimeRoot = root
	***REMOVED***
	if err := os.MkdirAll(config.RuntimeRoot, 0700); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	c := &Cluster***REMOVED***
		root:        root,
		config:      config,
		configEvent: make(chan lncluster.ConfigEventType, 10),
		runtimeRoot: config.RuntimeRoot,
		attachers:   make(map[string]*attacher),
		watchStream: config.WatchStream,
	***REMOVED***
	return c, nil
***REMOVED***

// Start the Cluster instance
// TODO The split between New and Start can be join again when the SendClusterEvent
// method is no longer required
func (c *Cluster) Start() error ***REMOVED***
	root := filepath.Join(c.config.Root, swarmDirName)

	nodeConfig, err := loadPersistentState(root)
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return nil
		***REMOVED***
		return err
	***REMOVED***

	nr, err := c.newNodeRunner(*nodeConfig)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	c.nr = nr

	select ***REMOVED***
	case <-time.After(swarmConnectTimeout):
		logrus.Error("swarm component could not be started before timeout was reached")
	case err := <-nr.Ready():
		if err != nil ***REMOVED***
			logrus.WithError(err).Error("swarm component could not be started")
			return nil
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (c *Cluster) newNodeRunner(conf nodeStartConfig) (*nodeRunner, error) ***REMOVED***
	if err := c.config.Backend.IsSwarmCompatible(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	actualLocalAddr := conf.LocalAddr
	if actualLocalAddr == "" ***REMOVED***
		// If localAddr was not specified, resolve it automatically
		// based on the route to joinAddr. localAddr can only be left
		// empty on "join".
		listenHost, _, err := net.SplitHostPort(conf.ListenAddr)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("could not parse listen address: %v", err)
		***REMOVED***

		listenAddrIP := net.ParseIP(listenHost)
		if listenAddrIP == nil || !listenAddrIP.IsUnspecified() ***REMOVED***
			actualLocalAddr = listenHost
		***REMOVED*** else ***REMOVED***
			if conf.RemoteAddr == "" ***REMOVED***
				// Should never happen except using swarms created by
				// old versions that didn't save remoteAddr.
				conf.RemoteAddr = "8.8.8.8:53"
			***REMOVED***
			conn, err := net.Dial("udp", conf.RemoteAddr)
			if err != nil ***REMOVED***
				return nil, fmt.Errorf("could not find local IP address: %v", err)
			***REMOVED***
			localHostPort := conn.LocalAddr().String()
			actualLocalAddr, _, _ = net.SplitHostPort(localHostPort)
			conn.Close()
		***REMOVED***
	***REMOVED***

	nr := &nodeRunner***REMOVED***cluster: c***REMOVED***
	nr.actualLocalAddr = actualLocalAddr

	if err := nr.Start(conf); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	c.config.Backend.DaemonJoinsCluster(c)

	return nr, nil
***REMOVED***

func (c *Cluster) getRequestContext() (context.Context, func()) ***REMOVED*** // TODO: not needed when requests don't block on qourum lost
	return context.WithTimeout(context.Background(), swarmRequestTimeout)
***REMOVED***

// IsManager returns true if Cluster is participating as a manager.
func (c *Cluster) IsManager() bool ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentNodeState().IsActiveManager()
***REMOVED***

// IsAgent returns true if Cluster is participating as a worker/agent.
func (c *Cluster) IsAgent() bool ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentNodeState().status == types.LocalNodeStateActive
***REMOVED***

// GetLocalAddress returns the local address.
func (c *Cluster) GetLocalAddress() string ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentNodeState().actualLocalAddr
***REMOVED***

// GetListenAddress returns the listen address.
func (c *Cluster) GetListenAddress() string ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.nr != nil ***REMOVED***
		return c.nr.config.ListenAddr
	***REMOVED***
	return ""
***REMOVED***

// GetAdvertiseAddress returns the remotely reachable address of this node.
func (c *Cluster) GetAdvertiseAddress() string ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.nr != nil && c.nr.config.AdvertiseAddr != "" ***REMOVED***
		advertiseHost, _, _ := net.SplitHostPort(c.nr.config.AdvertiseAddr)
		return advertiseHost
	***REMOVED***
	return c.currentNodeState().actualLocalAddr
***REMOVED***

// GetDataPathAddress returns the address to be used for the data path traffic, if specified.
func (c *Cluster) GetDataPathAddress() string ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.nr != nil ***REMOVED***
		return c.nr.config.DataPathAddr
	***REMOVED***
	return ""
***REMOVED***

// GetRemoteAddressList returns the advertise address for each of the remote managers if
// available.
func (c *Cluster) GetRemoteAddressList() []string ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.getRemoteAddressList()
***REMOVED***

func (c *Cluster) getRemoteAddressList() []string ***REMOVED***
	state := c.currentNodeState()
	if state.swarmNode == nil ***REMOVED***
		return []string***REMOVED******REMOVED***
	***REMOVED***

	nodeID := state.swarmNode.NodeID()
	remotes := state.swarmNode.Remotes()
	addressList := make([]string, 0, len(remotes))
	for _, r := range remotes ***REMOVED***
		if r.NodeID != nodeID ***REMOVED***
			addressList = append(addressList, r.Addr)
		***REMOVED***
	***REMOVED***
	return addressList
***REMOVED***

// ListenClusterEvents returns a channel that receives messages on cluster
// participation changes.
// todo: make cancelable and accessible to multiple callers
func (c *Cluster) ListenClusterEvents() <-chan lncluster.ConfigEventType ***REMOVED***
	return c.configEvent
***REMOVED***

// currentNodeState should not be called without a read lock
func (c *Cluster) currentNodeState() nodeState ***REMOVED***
	return c.nr.State()
***REMOVED***

// errNoManager returns error describing why manager commands can't be used.
// Call with read lock.
func (c *Cluster) errNoManager(st nodeState) error ***REMOVED***
	if st.swarmNode == nil ***REMOVED***
		if errors.Cause(st.err) == errSwarmLocked ***REMOVED***
			return errSwarmLocked
		***REMOVED***
		if st.err == errSwarmCertificatesExpired ***REMOVED***
			return errSwarmCertificatesExpired
		***REMOVED***
		return errors.WithStack(notAvailableError("This node is not a swarm manager. Use \"docker swarm init\" or \"docker swarm join\" to connect this node to swarm and try again."))
	***REMOVED***
	if st.swarmNode.Manager() != nil ***REMOVED***
		return errors.WithStack(notAvailableError("This node is not a swarm manager. Manager is being prepared or has trouble connecting to the cluster."))
	***REMOVED***
	return errors.WithStack(notAvailableError("This node is not a swarm manager. Worker nodes can't be used to view or modify cluster state. Please run this command on a manager node or promote the current node to a manager."))
***REMOVED***

// Cleanup stops active swarm node. This is run before daemon shutdown.
func (c *Cluster) Cleanup() ***REMOVED***
	c.controlMutex.Lock()
	defer c.controlMutex.Unlock()

	c.mu.Lock()
	node := c.nr
	if node == nil ***REMOVED***
		c.mu.Unlock()
		return
	***REMOVED***
	state := c.currentNodeState()
	c.mu.Unlock()

	if state.IsActiveManager() ***REMOVED***
		active, reachable, unreachable, err := managerStats(state.controlClient, state.NodeID())
		if err == nil ***REMOVED***
			singlenode := active && isLastManager(reachable, unreachable)
			if active && !singlenode && removingManagerCausesLossOfQuorum(reachable, unreachable) ***REMOVED***
				logrus.Errorf("Leaving cluster with %v managers left out of %v. Raft quorum will be lost.", reachable-1, reachable+unreachable)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if err := node.Stop(); err != nil ***REMOVED***
		logrus.Errorf("failed to shut down cluster node: %v", err)
		signal.DumpStacks("")
	***REMOVED***

	c.mu.Lock()
	c.nr = nil
	c.mu.Unlock()
***REMOVED***

func managerStats(client swarmapi.ControlClient, currentNodeID string) (current bool, reachable int, unreachable int, err error) ***REMOVED***
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	nodes, err := client.ListNodes(ctx, &swarmapi.ListNodesRequest***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		return false, 0, 0, err
	***REMOVED***
	for _, n := range nodes.Nodes ***REMOVED***
		if n.ManagerStatus != nil ***REMOVED***
			if n.ManagerStatus.Reachability == swarmapi.RaftMemberStatus_REACHABLE ***REMOVED***
				reachable++
				if n.ID == currentNodeID ***REMOVED***
					current = true
				***REMOVED***
			***REMOVED***
			if n.ManagerStatus.Reachability == swarmapi.RaftMemberStatus_UNREACHABLE ***REMOVED***
				unreachable++
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func detectLockedError(err error) error ***REMOVED***
	if err == swarmnode.ErrInvalidUnlockKey ***REMOVED***
		return errors.WithStack(errSwarmLocked)
	***REMOVED***
	return err
***REMOVED***

func (c *Cluster) lockedManagerAction(fn func(ctx context.Context, state nodeState) error) error ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := c.currentNodeState()
	if !state.IsActiveManager() ***REMOVED***
		return c.errNoManager(state)
	***REMOVED***

	ctx, cancel := c.getRequestContext()
	defer cancel()

	return fn(ctx, state)
***REMOVED***

// SendClusterEvent allows to send cluster events on the configEvent channel
// TODO This method should not be exposed.
// Currently it is used to notify the network controller that the keys are
// available
func (c *Cluster) SendClusterEvent(event lncluster.ConfigEventType) ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.configEvent <- event
***REMOVED***
