package dispatcher

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/docker/go-events"
	"github.com/docker/go-metrics"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/equality"
	"github.com/docker/swarmkit/ca"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/drivers"
	"github.com/docker/swarmkit/manager/state/store"
	"github.com/docker/swarmkit/protobuf/ptypes"
	"github.com/docker/swarmkit/remotes"
	"github.com/docker/swarmkit/watch"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/transport"
)

const (
	// DefaultHeartBeatPeriod is used for setting default value in cluster config
	// and in case if cluster config is missing.
	DefaultHeartBeatPeriod       = 5 * time.Second
	defaultHeartBeatEpsilon      = 500 * time.Millisecond
	defaultGracePeriodMultiplier = 3
	defaultRateLimitPeriod       = 8 * time.Second

	// maxBatchItems is the threshold of queued writes that should
	// trigger an actual transaction to commit them to the shared store.
	maxBatchItems = 10000

	// maxBatchInterval needs to strike a balance between keeping
	// latency low, and realizing opportunities to combine many writes
	// into a single transaction. A fraction of a second feels about
	// right.
	maxBatchInterval = 100 * time.Millisecond

	modificationBatchLimit = 100
	batchingWaitTime       = 100 * time.Millisecond

	// defaultNodeDownPeriod specifies the default time period we
	// wait before moving tasks assigned to down nodes to ORPHANED
	// state.
	defaultNodeDownPeriod = 24 * time.Hour
)

var (
	// ErrNodeAlreadyRegistered returned if node with same ID was already
	// registered with this dispatcher.
	ErrNodeAlreadyRegistered = errors.New("node already registered")
	// ErrNodeNotRegistered returned if node with such ID wasn't registered
	// with this dispatcher.
	ErrNodeNotRegistered = errors.New("node not registered")
	// ErrSessionInvalid returned when the session in use is no longer valid.
	// The node should re-register and start a new session.
	ErrSessionInvalid = errors.New("session invalid")
	// ErrNodeNotFound returned when the Node doesn't exist in raft.
	ErrNodeNotFound = errors.New("node not found")

	// Scheduling delay timer.
	schedulingDelayTimer metrics.Timer
)

func init() ***REMOVED***
	ns := metrics.NewNamespace("swarm", "dispatcher", nil)
	schedulingDelayTimer = ns.NewTimer("scheduling_delay",
		"Scheduling delay is the time a task takes to go from NEW to RUNNING state.")
	metrics.Register(ns)
***REMOVED***

// Config is configuration for Dispatcher. For default you should use
// DefaultConfig.
type Config struct ***REMOVED***
	HeartbeatPeriod  time.Duration
	HeartbeatEpsilon time.Duration
	// RateLimitPeriod specifies how often node with same ID can try to register
	// new session.
	RateLimitPeriod       time.Duration
	GracePeriodMultiplier int
***REMOVED***

// DefaultConfig returns default config for Dispatcher.
func DefaultConfig() *Config ***REMOVED***
	return &Config***REMOVED***
		HeartbeatPeriod:       DefaultHeartBeatPeriod,
		HeartbeatEpsilon:      defaultHeartBeatEpsilon,
		RateLimitPeriod:       defaultRateLimitPeriod,
		GracePeriodMultiplier: defaultGracePeriodMultiplier,
	***REMOVED***
***REMOVED***

// Cluster is interface which represent raft cluster. manager/state/raft.Node
// is implements it. This interface needed only for easier unit-testing.
type Cluster interface ***REMOVED***
	GetMemberlist() map[uint64]*api.RaftMember
	SubscribePeers() (chan events.Event, func())
	MemoryStore() *store.MemoryStore
***REMOVED***

// nodeUpdate provides a new status and/or description to apply to a node
// object.
type nodeUpdate struct ***REMOVED***
	status      *api.NodeStatus
	description *api.NodeDescription
***REMOVED***

// clusterUpdate is an object that stores an update to the cluster that should trigger
// a new session message.  These are pointers to indicate the difference between
// "there is no update" and "update this to nil"
type clusterUpdate struct ***REMOVED***
	managerUpdate      *[]*api.WeightedPeer
	bootstrapKeyUpdate *[]*api.EncryptionKey
	rootCAUpdate       *[]byte
***REMOVED***

// Dispatcher is responsible for dispatching tasks and tracking agent health.
type Dispatcher struct ***REMOVED***
	mu                   sync.Mutex
	wg                   sync.WaitGroup
	nodes                *nodeStore
	store                *store.MemoryStore
	lastSeenManagers     []*api.WeightedPeer
	networkBootstrapKeys []*api.EncryptionKey
	lastSeenRootCert     []byte
	config               *Config
	cluster              Cluster
	ctx                  context.Context
	cancel               context.CancelFunc
	clusterUpdateQueue   *watch.Queue
	dp                   *drivers.DriverProvider
	securityConfig       *ca.SecurityConfig

	taskUpdates     map[string]*api.TaskStatus // indexed by task ID
	taskUpdatesLock sync.Mutex

	nodeUpdates     map[string]nodeUpdate // indexed by node ID
	nodeUpdatesLock sync.Mutex

	downNodes *nodeStore

	processUpdatesTrigger chan struct***REMOVED******REMOVED***

	// for waiting for the next task/node batch update
	processUpdatesLock sync.Mutex
	processUpdatesCond *sync.Cond
***REMOVED***

// New returns Dispatcher with cluster interface(usually raft.Node).
func New(cluster Cluster, c *Config, dp *drivers.DriverProvider, securityConfig *ca.SecurityConfig) *Dispatcher ***REMOVED***
	d := &Dispatcher***REMOVED***
		dp:                    dp,
		nodes:                 newNodeStore(c.HeartbeatPeriod, c.HeartbeatEpsilon, c.GracePeriodMultiplier, c.RateLimitPeriod),
		downNodes:             newNodeStore(defaultNodeDownPeriod, 0, 1, 0),
		store:                 cluster.MemoryStore(),
		cluster:               cluster,
		processUpdatesTrigger: make(chan struct***REMOVED******REMOVED***, 1),
		config:                c,
		securityConfig:        securityConfig,
	***REMOVED***

	d.processUpdatesCond = sync.NewCond(&d.processUpdatesLock)

	return d
***REMOVED***

func getWeightedPeers(cluster Cluster) []*api.WeightedPeer ***REMOVED***
	members := cluster.GetMemberlist()
	var mgrs []*api.WeightedPeer
	for _, m := range members ***REMOVED***
		mgrs = append(mgrs, &api.WeightedPeer***REMOVED***
			Peer: &api.Peer***REMOVED***
				NodeID: m.NodeID,
				Addr:   m.Addr,
			***REMOVED***,

			// TODO(stevvooe): Calculate weight of manager selection based on
			// cluster-level observations, such as number of connections and
			// load.
			Weight: remotes.DefaultObservationWeight,
		***REMOVED***)
	***REMOVED***
	return mgrs
***REMOVED***

// Run runs dispatcher tasks which should be run on leader dispatcher.
// Dispatcher can be stopped with cancelling ctx or calling Stop().
func (d *Dispatcher) Run(ctx context.Context) error ***REMOVED***
	d.taskUpdatesLock.Lock()
	d.taskUpdates = make(map[string]*api.TaskStatus)
	d.taskUpdatesLock.Unlock()

	d.nodeUpdatesLock.Lock()
	d.nodeUpdates = make(map[string]nodeUpdate)
	d.nodeUpdatesLock.Unlock()

	d.mu.Lock()
	if d.isRunning() ***REMOVED***
		d.mu.Unlock()
		return errors.New("dispatcher is already running")
	***REMOVED***
	ctx = log.WithModule(ctx, "dispatcher")
	if err := d.markNodesUnknown(ctx); err != nil ***REMOVED***
		log.G(ctx).Errorf(`failed to move all nodes to "unknown" state: %v`, err)
	***REMOVED***
	configWatcher, cancel, err := store.ViewAndWatch(
		d.store,
		func(readTx store.ReadTx) error ***REMOVED***
			clusters, err := store.FindClusters(readTx, store.ByName(store.DefaultClusterName))
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if err == nil && len(clusters) == 1 ***REMOVED***
				heartbeatPeriod, err := gogotypes.DurationFromProto(clusters[0].Spec.Dispatcher.HeartbeatPeriod)
				if err == nil && heartbeatPeriod > 0 ***REMOVED***
					d.config.HeartbeatPeriod = heartbeatPeriod
				***REMOVED***
				if clusters[0].NetworkBootstrapKeys != nil ***REMOVED***
					d.networkBootstrapKeys = clusters[0].NetworkBootstrapKeys
				***REMOVED***
				d.lastSeenRootCert = clusters[0].RootCA.CACert
			***REMOVED***
			return nil
		***REMOVED***,
		api.EventUpdateCluster***REMOVED******REMOVED***,
	)
	if err != nil ***REMOVED***
		d.mu.Unlock()
		return err
	***REMOVED***
	// set queue here to guarantee that Close will close it
	d.clusterUpdateQueue = watch.NewQueue()

	peerWatcher, peerCancel := d.cluster.SubscribePeers()
	defer peerCancel()
	d.lastSeenManagers = getWeightedPeers(d.cluster)

	defer cancel()
	d.ctx, d.cancel = context.WithCancel(ctx)
	ctx = d.ctx
	d.wg.Add(1)
	defer d.wg.Done()
	d.mu.Unlock()

	publishManagers := func(peers []*api.Peer) ***REMOVED***
		var mgrs []*api.WeightedPeer
		for _, p := range peers ***REMOVED***
			mgrs = append(mgrs, &api.WeightedPeer***REMOVED***
				Peer:   p,
				Weight: remotes.DefaultObservationWeight,
			***REMOVED***)
		***REMOVED***
		d.mu.Lock()
		d.lastSeenManagers = mgrs
		d.mu.Unlock()
		d.clusterUpdateQueue.Publish(clusterUpdate***REMOVED***managerUpdate: &mgrs***REMOVED***)
	***REMOVED***

	batchTimer := time.NewTimer(maxBatchInterval)
	defer batchTimer.Stop()

	for ***REMOVED***
		select ***REMOVED***
		case ev := <-peerWatcher:
			publishManagers(ev.([]*api.Peer))
		case <-d.processUpdatesTrigger:
			d.processUpdates(ctx)
			batchTimer.Reset(maxBatchInterval)
		case <-batchTimer.C:
			d.processUpdates(ctx)
			batchTimer.Reset(maxBatchInterval)
		case v := <-configWatcher:
			cluster := v.(api.EventUpdateCluster)
			d.mu.Lock()
			if cluster.Cluster.Spec.Dispatcher.HeartbeatPeriod != nil ***REMOVED***
				// ignore error, since Spec has passed validation before
				heartbeatPeriod, _ := gogotypes.DurationFromProto(cluster.Cluster.Spec.Dispatcher.HeartbeatPeriod)
				if heartbeatPeriod != d.config.HeartbeatPeriod ***REMOVED***
					// only call d.nodes.updatePeriod when heartbeatPeriod changes
					d.config.HeartbeatPeriod = heartbeatPeriod
					d.nodes.updatePeriod(d.config.HeartbeatPeriod, d.config.HeartbeatEpsilon, d.config.GracePeriodMultiplier)
				***REMOVED***
			***REMOVED***
			d.lastSeenRootCert = cluster.Cluster.RootCA.CACert
			d.networkBootstrapKeys = cluster.Cluster.NetworkBootstrapKeys
			d.mu.Unlock()
			d.clusterUpdateQueue.Publish(clusterUpdate***REMOVED***
				bootstrapKeyUpdate: &cluster.Cluster.NetworkBootstrapKeys,
				rootCAUpdate:       &cluster.Cluster.RootCA.CACert,
			***REMOVED***)
		case <-ctx.Done():
			return nil
		***REMOVED***
	***REMOVED***
***REMOVED***

// Stop stops dispatcher and closes all grpc streams.
func (d *Dispatcher) Stop() error ***REMOVED***
	d.mu.Lock()
	if !d.isRunning() ***REMOVED***
		d.mu.Unlock()
		return errors.New("dispatcher is already stopped")
	***REMOVED***
	d.cancel()
	d.mu.Unlock()
	d.nodes.Clean()

	d.processUpdatesLock.Lock()
	// In case there are any waiters. There is no chance of any starting
	// after this point, because they check if the context is canceled
	// before waiting.
	d.processUpdatesCond.Broadcast()
	d.processUpdatesLock.Unlock()

	d.clusterUpdateQueue.Close()

	d.wg.Wait()

	return nil
***REMOVED***

func (d *Dispatcher) isRunningLocked() (context.Context, error) ***REMOVED***
	d.mu.Lock()
	if !d.isRunning() ***REMOVED***
		d.mu.Unlock()
		return nil, status.Errorf(codes.Aborted, "dispatcher is stopped")
	***REMOVED***
	ctx := d.ctx
	d.mu.Unlock()
	return ctx, nil
***REMOVED***

func (d *Dispatcher) markNodesUnknown(ctx context.Context) error ***REMOVED***
	log := log.G(ctx).WithField("method", "(*Dispatcher).markNodesUnknown")
	var nodes []*api.Node
	var err error
	d.store.View(func(tx store.ReadTx) ***REMOVED***
		nodes, err = store.FindNodes(tx, store.All)
	***REMOVED***)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed to get list of nodes")
	***REMOVED***
	err = d.store.Batch(func(batch *store.Batch) error ***REMOVED***
		for _, n := range nodes ***REMOVED***
			err := batch.Update(func(tx store.Tx) error ***REMOVED***
				// check if node is still here
				node := store.GetNode(tx, n.ID)
				if node == nil ***REMOVED***
					return nil
				***REMOVED***
				// do not try to resurrect down nodes
				if node.Status.State == api.NodeStatus_DOWN ***REMOVED***
					nodeCopy := node
					expireFunc := func() ***REMOVED***
						if err := d.moveTasksToOrphaned(nodeCopy.ID); err != nil ***REMOVED***
							log.WithError(err).Error(`failed to move all tasks to "ORPHANED" state`)
						***REMOVED***

						d.downNodes.Delete(nodeCopy.ID)
					***REMOVED***

					d.downNodes.Add(nodeCopy, expireFunc)
					return nil
				***REMOVED***

				node.Status.State = api.NodeStatus_UNKNOWN
				node.Status.Message = `Node moved to "unknown" state due to leadership change in cluster`

				nodeID := node.ID

				expireFunc := func() ***REMOVED***
					log := log.WithField("node", nodeID)
					log.Debug("heartbeat expiration for unknown node")
					if err := d.markNodeNotReady(nodeID, api.NodeStatus_DOWN, `heartbeat failure for node in "unknown" state`); err != nil ***REMOVED***
						log.WithError(err).Error(`failed deregistering node after heartbeat expiration for node in "unknown" state`)
					***REMOVED***
				***REMOVED***
				if err := d.nodes.AddUnknown(node, expireFunc); err != nil ***REMOVED***
					return errors.Wrap(err, `adding node in "unknown" state to node store failed`)
				***REMOVED***
				if err := store.UpdateNode(tx, node); err != nil ***REMOVED***
					return errors.Wrap(err, "update failed")
				***REMOVED***
				return nil
			***REMOVED***)
			if err != nil ***REMOVED***
				log.WithField("node", n.ID).WithError(err).Error(`failed to move node to "unknown" state`)
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***)
	return err
***REMOVED***

func (d *Dispatcher) isRunning() bool ***REMOVED***
	if d.ctx == nil ***REMOVED***
		return false
	***REMOVED***
	select ***REMOVED***
	case <-d.ctx.Done():
		return false
	default:
	***REMOVED***
	return true
***REMOVED***

// markNodeReady updates the description of a node, updates its address, and sets status to READY
// this is used during registration when a new node description is provided
// and during node updates when the node description changes
func (d *Dispatcher) markNodeReady(ctx context.Context, nodeID string, description *api.NodeDescription, addr string) error ***REMOVED***
	d.nodeUpdatesLock.Lock()
	d.nodeUpdates[nodeID] = nodeUpdate***REMOVED***
		status: &api.NodeStatus***REMOVED***
			State: api.NodeStatus_READY,
			Addr:  addr,
		***REMOVED***,
		description: description,
	***REMOVED***
	numUpdates := len(d.nodeUpdates)
	d.nodeUpdatesLock.Unlock()

	// Node is marked ready. Remove the node from down nodes if it
	// is there.
	d.downNodes.Delete(nodeID)

	if numUpdates >= maxBatchItems ***REMOVED***
		select ***REMOVED***
		case d.processUpdatesTrigger <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
		case <-ctx.Done():
			return ctx.Err()
		***REMOVED***

	***REMOVED***

	// Wait until the node update batch happens before unblocking register.
	d.processUpdatesLock.Lock()
	defer d.processUpdatesLock.Unlock()

	select ***REMOVED***
	case <-ctx.Done():
		return ctx.Err()
	default:
	***REMOVED***
	d.processUpdatesCond.Wait()

	return nil
***REMOVED***

// gets the node IP from the context of a grpc call
func nodeIPFromContext(ctx context.Context) (string, error) ***REMOVED***
	nodeInfo, err := ca.RemoteNode(ctx)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	addr, _, err := net.SplitHostPort(nodeInfo.RemoteAddr)
	if err != nil ***REMOVED***
		return "", errors.Wrap(err, "unable to get ip from addr:port")
	***REMOVED***
	return addr, nil
***REMOVED***

// register is used for registration of node with particular dispatcher.
func (d *Dispatcher) register(ctx context.Context, nodeID string, description *api.NodeDescription) (string, error) ***REMOVED***
	// prevent register until we're ready to accept it
	dctx, err := d.isRunningLocked()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	if err := d.nodes.CheckRateLimit(nodeID); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	// TODO(stevvooe): Validate node specification.
	var node *api.Node
	d.store.View(func(tx store.ReadTx) ***REMOVED***
		node = store.GetNode(tx, nodeID)
	***REMOVED***)
	if node == nil ***REMOVED***
		return "", ErrNodeNotFound
	***REMOVED***

	addr, err := nodeIPFromContext(ctx)
	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Debug("failed to get remote node IP")
	***REMOVED***

	if err := d.markNodeReady(dctx, nodeID, description, addr); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	expireFunc := func() ***REMOVED***
		log.G(ctx).Debug("heartbeat expiration")
		if err := d.markNodeNotReady(nodeID, api.NodeStatus_DOWN, "heartbeat failure"); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("failed deregistering node after heartbeat expiration")
		***REMOVED***
	***REMOVED***

	rn := d.nodes.Add(node, expireFunc)

	// NOTE(stevvooe): We need be a little careful with re-registration. The
	// current implementation just matches the node id and then gives away the
	// sessionID. If we ever want to use sessionID as a secret, which we may
	// want to, this is giving away the keys to the kitchen.
	//
	// The right behavior is going to be informed by identity. Basically, each
	// time a node registers, we invalidate the session and issue a new
	// session, once identity is proven. This will cause misbehaved agents to
	// be kicked when multiple connections are made.
	return rn.SessionID, nil
***REMOVED***

// UpdateTaskStatus updates status of task. Node should send such updates
// on every status change of its tasks.
func (d *Dispatcher) UpdateTaskStatus(ctx context.Context, r *api.UpdateTaskStatusRequest) (*api.UpdateTaskStatusResponse, error) ***REMOVED***
	nodeInfo, err := ca.RemoteNode(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	nodeID := nodeInfo.NodeID
	fields := logrus.Fields***REMOVED***
		"node.id":      nodeID,
		"node.session": r.SessionID,
		"method":       "(*Dispatcher).UpdateTaskStatus",
	***REMOVED***
	if nodeInfo.ForwardedBy != nil ***REMOVED***
		fields["forwarder.id"] = nodeInfo.ForwardedBy.NodeID
	***REMOVED***
	log := log.G(ctx).WithFields(fields)

	dctx, err := d.isRunningLocked()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if _, err := d.nodes.GetWithSession(nodeID, r.SessionID); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	validTaskUpdates := make([]*api.UpdateTaskStatusRequest_TaskStatusUpdate, 0, len(r.Updates))

	// Validate task updates
	for _, u := range r.Updates ***REMOVED***
		if u.Status == nil ***REMOVED***
			log.WithField("task.id", u.TaskID).Warn("task report has nil status")
			continue
		***REMOVED***

		var t *api.Task
		d.store.View(func(tx store.ReadTx) ***REMOVED***
			t = store.GetTask(tx, u.TaskID)
		***REMOVED***)
		if t == nil ***REMOVED***
			// Task may have been deleted
			log.WithField("task.id", u.TaskID).Debug("cannot find target task in store")
			continue
		***REMOVED***

		if t.NodeID != nodeID ***REMOVED***
			err := status.Errorf(codes.PermissionDenied, "cannot update a task not assigned this node")
			log.WithField("task.id", u.TaskID).Error(err)
			return nil, err
		***REMOVED***

		validTaskUpdates = append(validTaskUpdates, u)
	***REMOVED***

	d.taskUpdatesLock.Lock()
	// Enqueue task updates
	for _, u := range validTaskUpdates ***REMOVED***
		d.taskUpdates[u.TaskID] = u.Status
	***REMOVED***

	numUpdates := len(d.taskUpdates)
	d.taskUpdatesLock.Unlock()

	if numUpdates >= maxBatchItems ***REMOVED***
		select ***REMOVED***
		case d.processUpdatesTrigger <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
		case <-dctx.Done():
		***REMOVED***
	***REMOVED***
	return nil, nil
***REMOVED***

func (d *Dispatcher) processUpdates(ctx context.Context) ***REMOVED***
	var (
		taskUpdates map[string]*api.TaskStatus
		nodeUpdates map[string]nodeUpdate
	)
	d.taskUpdatesLock.Lock()
	if len(d.taskUpdates) != 0 ***REMOVED***
		taskUpdates = d.taskUpdates
		d.taskUpdates = make(map[string]*api.TaskStatus)
	***REMOVED***
	d.taskUpdatesLock.Unlock()

	d.nodeUpdatesLock.Lock()
	if len(d.nodeUpdates) != 0 ***REMOVED***
		nodeUpdates = d.nodeUpdates
		d.nodeUpdates = make(map[string]nodeUpdate)
	***REMOVED***
	d.nodeUpdatesLock.Unlock()

	if len(taskUpdates) == 0 && len(nodeUpdates) == 0 ***REMOVED***
		return
	***REMOVED***

	log := log.G(ctx).WithFields(logrus.Fields***REMOVED***
		"method": "(*Dispatcher).processUpdates",
	***REMOVED***)

	err := d.store.Batch(func(batch *store.Batch) error ***REMOVED***
		for taskID, status := range taskUpdates ***REMOVED***
			err := batch.Update(func(tx store.Tx) error ***REMOVED***
				logger := log.WithField("task.id", taskID)
				task := store.GetTask(tx, taskID)
				if task == nil ***REMOVED***
					// Task may have been deleted
					logger.Debug("cannot find target task in store")
					return nil
				***REMOVED***

				logger = logger.WithField("state.transition", fmt.Sprintf("%v->%v", task.Status.State, status.State))

				if task.Status == *status ***REMOVED***
					logger.Debug("task status identical, ignoring")
					return nil
				***REMOVED***

				if task.Status.State > status.State ***REMOVED***
					logger.Debug("task status invalid transition")
					return nil
				***REMOVED***

				// Update scheduling delay metric for running tasks.
				// We use the status update time on the leader to calculate the scheduling delay.
				// Because of this, the recorded scheduling delay will be an overestimate and include
				// the network delay between the worker and the leader.
				// This is not ideal, but its a known overestimation, rather than using the status update time
				// from the worker node, which may cause unknown incorrect results due to possible clock skew.
				if status.State == api.TaskStateRunning ***REMOVED***
					start := time.Unix(status.AppliedAt.GetSeconds(), int64(status.AppliedAt.GetNanos()))
					schedulingDelayTimer.UpdateSince(start)
				***REMOVED***

				task.Status = *status
				task.Status.AppliedBy = d.securityConfig.ClientTLSCreds.NodeID()
				task.Status.AppliedAt = ptypes.MustTimestampProto(time.Now())
				if err := store.UpdateTask(tx, task); err != nil ***REMOVED***
					logger.WithError(err).Error("failed to update task status")
					return nil
				***REMOVED***
				logger.Debug("dispatcher committed status update to store")
				return nil
			***REMOVED***)
			if err != nil ***REMOVED***
				log.WithError(err).Error("dispatcher task update transaction failed")
			***REMOVED***
		***REMOVED***

		for nodeID, nodeUpdate := range nodeUpdates ***REMOVED***
			err := batch.Update(func(tx store.Tx) error ***REMOVED***
				logger := log.WithField("node.id", nodeID)
				node := store.GetNode(tx, nodeID)
				if node == nil ***REMOVED***
					logger.Errorf("node unavailable")
					return nil
				***REMOVED***

				if nodeUpdate.status != nil ***REMOVED***
					node.Status.State = nodeUpdate.status.State
					node.Status.Message = nodeUpdate.status.Message
					if nodeUpdate.status.Addr != "" ***REMOVED***
						node.Status.Addr = nodeUpdate.status.Addr
					***REMOVED***
				***REMOVED***
				if nodeUpdate.description != nil ***REMOVED***
					node.Description = nodeUpdate.description
				***REMOVED***

				if err := store.UpdateNode(tx, node); err != nil ***REMOVED***
					logger.WithError(err).Error("failed to update node status")
					return nil
				***REMOVED***
				logger.Debug("node status updated")
				return nil
			***REMOVED***)
			if err != nil ***REMOVED***
				log.WithError(err).Error("dispatcher node update transaction failed")
			***REMOVED***
		***REMOVED***

		return nil
	***REMOVED***)
	if err != nil ***REMOVED***
		log.WithError(err).Error("dispatcher batch failed")
	***REMOVED***

	d.processUpdatesCond.Broadcast()
***REMOVED***

// Tasks is a stream of tasks state for node. Each message contains full list
// of tasks which should be run on node, if task is not present in that list,
// it should be terminated.
func (d *Dispatcher) Tasks(r *api.TasksRequest, stream api.Dispatcher_TasksServer) error ***REMOVED***
	nodeInfo, err := ca.RemoteNode(stream.Context())
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	nodeID := nodeInfo.NodeID

	dctx, err := d.isRunningLocked()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	fields := logrus.Fields***REMOVED***
		"node.id":      nodeID,
		"node.session": r.SessionID,
		"method":       "(*Dispatcher).Tasks",
	***REMOVED***
	if nodeInfo.ForwardedBy != nil ***REMOVED***
		fields["forwarder.id"] = nodeInfo.ForwardedBy.NodeID
	***REMOVED***
	log.G(stream.Context()).WithFields(fields).Debug("")

	if _, err = d.nodes.GetWithSession(nodeID, r.SessionID); err != nil ***REMOVED***
		return err
	***REMOVED***

	tasksMap := make(map[string]*api.Task)
	nodeTasks, cancel, err := store.ViewAndWatch(
		d.store,
		func(readTx store.ReadTx) error ***REMOVED***
			tasks, err := store.FindTasks(readTx, store.ByNodeID(nodeID))
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			for _, t := range tasks ***REMOVED***
				tasksMap[t.ID] = t
			***REMOVED***
			return nil
		***REMOVED***,
		api.EventCreateTask***REMOVED***Task: &api.Task***REMOVED***NodeID: nodeID***REMOVED***,
			Checks: []api.TaskCheckFunc***REMOVED***api.TaskCheckNodeID***REMOVED******REMOVED***,
		api.EventUpdateTask***REMOVED***Task: &api.Task***REMOVED***NodeID: nodeID***REMOVED***,
			Checks: []api.TaskCheckFunc***REMOVED***api.TaskCheckNodeID***REMOVED******REMOVED***,
		api.EventDeleteTask***REMOVED***Task: &api.Task***REMOVED***NodeID: nodeID***REMOVED***,
			Checks: []api.TaskCheckFunc***REMOVED***api.TaskCheckNodeID***REMOVED******REMOVED***,
	)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer cancel()

	for ***REMOVED***
		if _, err := d.nodes.GetWithSession(nodeID, r.SessionID); err != nil ***REMOVED***
			return err
		***REMOVED***

		var tasks []*api.Task
		for _, t := range tasksMap ***REMOVED***
			// dispatcher only sends tasks that have been assigned to a node
			if t != nil && t.Status.State >= api.TaskStateAssigned ***REMOVED***
				tasks = append(tasks, t)
			***REMOVED***
		***REMOVED***

		if err := stream.Send(&api.TasksMessage***REMOVED***Tasks: tasks***REMOVED***); err != nil ***REMOVED***
			return err
		***REMOVED***

		// bursty events should be processed in batches and sent out snapshot
		var (
			modificationCnt int
			batchingTimer   *time.Timer
			batchingTimeout <-chan time.Time
		)

	batchingLoop:
		for modificationCnt < modificationBatchLimit ***REMOVED***
			select ***REMOVED***
			case event := <-nodeTasks:
				switch v := event.(type) ***REMOVED***
				case api.EventCreateTask:
					tasksMap[v.Task.ID] = v.Task
					modificationCnt++
				case api.EventUpdateTask:
					if oldTask, exists := tasksMap[v.Task.ID]; exists ***REMOVED***
						// States ASSIGNED and below are set by the orchestrator/scheduler,
						// not the agent, so tasks in these states need to be sent to the
						// agent even if nothing else has changed.
						if equality.TasksEqualStable(oldTask, v.Task) && v.Task.Status.State > api.TaskStateAssigned ***REMOVED***
							// this update should not trigger action at agent
							tasksMap[v.Task.ID] = v.Task
							continue
						***REMOVED***
					***REMOVED***
					tasksMap[v.Task.ID] = v.Task
					modificationCnt++
				case api.EventDeleteTask:
					delete(tasksMap, v.Task.ID)
					modificationCnt++
				***REMOVED***
				if batchingTimer != nil ***REMOVED***
					batchingTimer.Reset(batchingWaitTime)
				***REMOVED*** else ***REMOVED***
					batchingTimer = time.NewTimer(batchingWaitTime)
					batchingTimeout = batchingTimer.C
				***REMOVED***
			case <-batchingTimeout:
				break batchingLoop
			case <-stream.Context().Done():
				return stream.Context().Err()
			case <-dctx.Done():
				return dctx.Err()
			***REMOVED***
		***REMOVED***

		if batchingTimer != nil ***REMOVED***
			batchingTimer.Stop()
		***REMOVED***
	***REMOVED***
***REMOVED***

// Assignments is a stream of assignments for a node. Each message contains
// either full list of tasks and secrets for the node, or an incremental update.
func (d *Dispatcher) Assignments(r *api.AssignmentsRequest, stream api.Dispatcher_AssignmentsServer) error ***REMOVED***
	nodeInfo, err := ca.RemoteNode(stream.Context())
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	nodeID := nodeInfo.NodeID

	dctx, err := d.isRunningLocked()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	fields := logrus.Fields***REMOVED***
		"node.id":      nodeID,
		"node.session": r.SessionID,
		"method":       "(*Dispatcher).Assignments",
	***REMOVED***
	if nodeInfo.ForwardedBy != nil ***REMOVED***
		fields["forwarder.id"] = nodeInfo.ForwardedBy.NodeID
	***REMOVED***
	log := log.G(stream.Context()).WithFields(fields)
	log.Debug("")

	if _, err = d.nodes.GetWithSession(nodeID, r.SessionID); err != nil ***REMOVED***
		return err
	***REMOVED***

	var (
		sequence    int64
		appliesTo   string
		assignments = newAssignmentSet(log, d.dp)
	)

	sendMessage := func(msg api.AssignmentsMessage, assignmentType api.AssignmentsMessage_Type) error ***REMOVED***
		sequence++
		msg.AppliesTo = appliesTo
		msg.ResultsIn = strconv.FormatInt(sequence, 10)
		appliesTo = msg.ResultsIn
		msg.Type = assignmentType

		return stream.Send(&msg)
	***REMOVED***

	// TODO(aaronl): Also send node secrets that should be exposed to
	// this node.
	nodeTasks, cancel, err := store.ViewAndWatch(
		d.store,
		func(readTx store.ReadTx) error ***REMOVED***
			tasks, err := store.FindTasks(readTx, store.ByNodeID(nodeID))
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			for _, t := range tasks ***REMOVED***
				assignments.addOrUpdateTask(readTx, t)
			***REMOVED***

			return nil
		***REMOVED***,
		api.EventUpdateTask***REMOVED***Task: &api.Task***REMOVED***NodeID: nodeID***REMOVED***,
			Checks: []api.TaskCheckFunc***REMOVED***api.TaskCheckNodeID***REMOVED******REMOVED***,
		api.EventDeleteTask***REMOVED***Task: &api.Task***REMOVED***NodeID: nodeID***REMOVED***,
			Checks: []api.TaskCheckFunc***REMOVED***api.TaskCheckNodeID***REMOVED******REMOVED***,
	)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer cancel()

	if err := sendMessage(assignments.message(), api.AssignmentsMessage_COMPLETE); err != nil ***REMOVED***
		return err
	***REMOVED***

	for ***REMOVED***
		// Check for session expiration
		if _, err := d.nodes.GetWithSession(nodeID, r.SessionID); err != nil ***REMOVED***
			return err
		***REMOVED***

		// bursty events should be processed in batches and sent out together
		var (
			modificationCnt int
			batchingTimer   *time.Timer
			batchingTimeout <-chan time.Time
		)

		oneModification := func() ***REMOVED***
			modificationCnt++

			if batchingTimer != nil ***REMOVED***
				batchingTimer.Reset(batchingWaitTime)
			***REMOVED*** else ***REMOVED***
				batchingTimer = time.NewTimer(batchingWaitTime)
				batchingTimeout = batchingTimer.C
			***REMOVED***
		***REMOVED***

		// The batching loop waits for 50 ms after the most recent
		// change, or until modificationBatchLimit is reached. The
		// worst case latency is modificationBatchLimit * batchingWaitTime,
		// which is 10 seconds.
	batchingLoop:
		for modificationCnt < modificationBatchLimit ***REMOVED***
			select ***REMOVED***
			case event := <-nodeTasks:
				switch v := event.(type) ***REMOVED***
				// We don't monitor EventCreateTask because tasks are
				// never created in the ASSIGNED state. First tasks are
				// created by the orchestrator, then the scheduler moves
				// them to ASSIGNED. If this ever changes, we will need
				// to monitor task creations as well.
				case api.EventUpdateTask:
					d.store.View(func(readTx store.ReadTx) ***REMOVED***
						if assignments.addOrUpdateTask(readTx, v.Task) ***REMOVED***
							oneModification()
						***REMOVED***
					***REMOVED***)
				case api.EventDeleteTask:
					if assignments.removeTask(v.Task) ***REMOVED***
						oneModification()
					***REMOVED***
					// TODO(aaronl): For node secrets, we'll need to handle
					// EventCreateSecret.
				***REMOVED***
			case <-batchingTimeout:
				break batchingLoop
			case <-stream.Context().Done():
				return stream.Context().Err()
			case <-dctx.Done():
				return dctx.Err()
			***REMOVED***
		***REMOVED***

		if batchingTimer != nil ***REMOVED***
			batchingTimer.Stop()
		***REMOVED***

		if modificationCnt > 0 ***REMOVED***
			if err := sendMessage(assignments.message(), api.AssignmentsMessage_INCREMENTAL); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (d *Dispatcher) moveTasksToOrphaned(nodeID string) error ***REMOVED***
	err := d.store.Batch(func(batch *store.Batch) error ***REMOVED***
		var (
			tasks []*api.Task
			err   error
		)

		d.store.View(func(tx store.ReadTx) ***REMOVED***
			tasks, err = store.FindTasks(tx, store.ByNodeID(nodeID))
		***REMOVED***)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		for _, task := range tasks ***REMOVED***
			// Tasks running on an unreachable node need to be marked as
			// orphaned since we have no idea whether the task is still running
			// or not.
			//
			// This only applies for tasks that could have made progress since
			// the agent became unreachable (assigned<->running)
			//
			// Tasks in a final state (e.g. rejected) *cannot* have made
			// progress, therefore there's no point in marking them as orphaned
			if task.Status.State >= api.TaskStateAssigned && task.Status.State <= api.TaskStateRunning ***REMOVED***
				task.Status.State = api.TaskStateOrphaned
			***REMOVED***

			if err := batch.Update(func(tx store.Tx) error ***REMOVED***
				err := store.UpdateTask(tx, task)
				if err != nil ***REMOVED***
					return err
				***REMOVED***

				return nil
			***REMOVED***); err != nil ***REMOVED***
				return err
			***REMOVED***

		***REMOVED***

		return nil
	***REMOVED***)

	return err
***REMOVED***

// markNodeNotReady sets the node state to some state other than READY
func (d *Dispatcher) markNodeNotReady(id string, state api.NodeStatus_State, message string) error ***REMOVED***
	dctx, err := d.isRunningLocked()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Node is down. Add it to down nodes so that we can keep
	// track of tasks assigned to the node.
	var node *api.Node
	d.store.View(func(readTx store.ReadTx) ***REMOVED***
		node = store.GetNode(readTx, id)
		if node == nil ***REMOVED***
			err = fmt.Errorf("could not find node %s while trying to add to down nodes store", id)
		***REMOVED***
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	expireFunc := func() ***REMOVED***
		if err := d.moveTasksToOrphaned(id); err != nil ***REMOVED***
			log.G(dctx).WithError(err).Error(`failed to move all tasks to "ORPHANED" state`)
		***REMOVED***

		d.downNodes.Delete(id)
	***REMOVED***

	d.downNodes.Add(node, expireFunc)

	status := &api.NodeStatus***REMOVED***
		State:   state,
		Message: message,
	***REMOVED***

	d.nodeUpdatesLock.Lock()
	// pluck the description out of nodeUpdates. this protects against a case
	// where a node is marked ready and a description is added, but then the
	// node is immediately marked not ready. this preserves that description
	d.nodeUpdates[id] = nodeUpdate***REMOVED***status: status, description: d.nodeUpdates[id].description***REMOVED***
	numUpdates := len(d.nodeUpdates)
	d.nodeUpdatesLock.Unlock()

	if numUpdates >= maxBatchItems ***REMOVED***
		select ***REMOVED***
		case d.processUpdatesTrigger <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
		case <-dctx.Done():
		***REMOVED***
	***REMOVED***

	if rn := d.nodes.Delete(id); rn == nil ***REMOVED***
		return errors.Errorf("node %s is not found in local storage", id)
	***REMOVED***

	return nil
***REMOVED***

// Heartbeat is heartbeat method for nodes. It returns new TTL in response.
// Node should send new heartbeat earlier than now + TTL, otherwise it will
// be deregistered from dispatcher and its status will be updated to NodeStatus_DOWN
func (d *Dispatcher) Heartbeat(ctx context.Context, r *api.HeartbeatRequest) (*api.HeartbeatResponse, error) ***REMOVED***
	nodeInfo, err := ca.RemoteNode(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	period, err := d.nodes.Heartbeat(nodeInfo.NodeID, r.SessionID)
	return &api.HeartbeatResponse***REMOVED***Period: period***REMOVED***, err
***REMOVED***

func (d *Dispatcher) getManagers() []*api.WeightedPeer ***REMOVED***
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.lastSeenManagers
***REMOVED***

func (d *Dispatcher) getNetworkBootstrapKeys() []*api.EncryptionKey ***REMOVED***
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.networkBootstrapKeys
***REMOVED***

func (d *Dispatcher) getRootCACert() []byte ***REMOVED***
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.lastSeenRootCert
***REMOVED***

// Session is a stream which controls agent connection.
// Each message contains list of backup Managers with weights. Also there is
// a special boolean field Disconnect which if true indicates that node should
// reconnect to another Manager immediately.
func (d *Dispatcher) Session(r *api.SessionRequest, stream api.Dispatcher_SessionServer) error ***REMOVED***
	ctx := stream.Context()
	nodeInfo, err := ca.RemoteNode(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	nodeID := nodeInfo.NodeID

	dctx, err := d.isRunningLocked()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var sessionID string
	if _, err := d.nodes.GetWithSession(nodeID, r.SessionID); err != nil ***REMOVED***
		// register the node.
		sessionID, err = d.register(ctx, nodeID, r.Description)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		sessionID = r.SessionID
		// get the node IP addr
		addr, err := nodeIPFromContext(stream.Context())
		if err != nil ***REMOVED***
			log.G(ctx).WithError(err).Debug("failed to get remote node IP")
		***REMOVED***
		// update the node description
		if err := d.markNodeReady(dctx, nodeID, r.Description, addr); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	fields := logrus.Fields***REMOVED***
		"node.id":      nodeID,
		"node.session": sessionID,
		"method":       "(*Dispatcher).Session",
	***REMOVED***
	if nodeInfo.ForwardedBy != nil ***REMOVED***
		fields["forwarder.id"] = nodeInfo.ForwardedBy.NodeID
	***REMOVED***
	log := log.G(ctx).WithFields(fields)

	var nodeObj *api.Node
	nodeUpdates, cancel, err := store.ViewAndWatch(d.store, func(readTx store.ReadTx) error ***REMOVED***
		nodeObj = store.GetNode(readTx, nodeID)
		return nil
	***REMOVED***, api.EventUpdateNode***REMOVED***Node: &api.Node***REMOVED***ID: nodeID***REMOVED***,
		Checks: []api.NodeCheckFunc***REMOVED***api.NodeCheckID***REMOVED******REMOVED***,
	)
	if cancel != nil ***REMOVED***
		defer cancel()
	***REMOVED***

	if err != nil ***REMOVED***
		log.WithError(err).Error("ViewAndWatch Node failed")
	***REMOVED***

	if _, err = d.nodes.GetWithSession(nodeID, sessionID); err != nil ***REMOVED***
		return err
	***REMOVED***

	clusterUpdatesCh, clusterCancel := d.clusterUpdateQueue.Watch()
	defer clusterCancel()

	if err := stream.Send(&api.SessionMessage***REMOVED***
		SessionID:            sessionID,
		Node:                 nodeObj,
		Managers:             d.getManagers(),
		NetworkBootstrapKeys: d.getNetworkBootstrapKeys(),
		RootCA:               d.getRootCACert(),
	***REMOVED***); err != nil ***REMOVED***
		return err
	***REMOVED***

	// disconnectNode is a helper forcibly shutdown connection
	disconnectNode := func() error ***REMOVED***
		// force disconnect by shutting down the stream.
		transportStream, ok := transport.StreamFromContext(stream.Context())
		if ok ***REMOVED***
			// if we have the transport stream, we can signal a disconnect
			// in the client.
			if err := transportStream.ServerTransport().Close(); err != nil ***REMOVED***
				log.WithError(err).Error("session end")
			***REMOVED***
		***REMOVED***

		if err := d.markNodeNotReady(nodeID, api.NodeStatus_DISCONNECTED, "node is currently trying to find new manager"); err != nil ***REMOVED***
			log.WithError(err).Error("failed to remove node")
		***REMOVED***
		// still return an abort if the transport closure was ineffective.
		return status.Errorf(codes.Aborted, "node must disconnect")
	***REMOVED***

	for ***REMOVED***
		// After each message send, we need to check the nodes sessionID hasn't
		// changed. If it has, we will shut down the stream and make the node
		// re-register.
		node, err := d.nodes.GetWithSession(nodeID, sessionID)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		var (
			disconnect bool
			mgrs       []*api.WeightedPeer
			netKeys    []*api.EncryptionKey
			rootCert   []byte
		)

		select ***REMOVED***
		case ev := <-clusterUpdatesCh:
			update := ev.(clusterUpdate)
			if update.managerUpdate != nil ***REMOVED***
				mgrs = *update.managerUpdate
			***REMOVED***
			if update.bootstrapKeyUpdate != nil ***REMOVED***
				netKeys = *update.bootstrapKeyUpdate
			***REMOVED***
			if update.rootCAUpdate != nil ***REMOVED***
				rootCert = *update.rootCAUpdate
			***REMOVED***
		case ev := <-nodeUpdates:
			nodeObj = ev.(api.EventUpdateNode).Node
		case <-stream.Context().Done():
			return stream.Context().Err()
		case <-node.Disconnect:
			disconnect = true
		case <-dctx.Done():
			disconnect = true
		***REMOVED***
		if mgrs == nil ***REMOVED***
			mgrs = d.getManagers()
		***REMOVED***
		if netKeys == nil ***REMOVED***
			netKeys = d.getNetworkBootstrapKeys()
		***REMOVED***
		if rootCert == nil ***REMOVED***
			rootCert = d.getRootCACert()
		***REMOVED***

		if err := stream.Send(&api.SessionMessage***REMOVED***
			SessionID:            sessionID,
			Node:                 nodeObj,
			Managers:             mgrs,
			NetworkBootstrapKeys: netKeys,
			RootCA:               rootCert,
		***REMOVED***); err != nil ***REMOVED***
			return err
		***REMOVED***
		if disconnect ***REMOVED***
			return disconnectNode()
		***REMOVED***
	***REMOVED***
***REMOVED***
