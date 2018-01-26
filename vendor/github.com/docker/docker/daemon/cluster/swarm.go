package cluster

import (
	"fmt"
	"net"
	"strings"
	"time"

	apitypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	types "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/daemon/cluster/convert"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/opts"
	"github.com/docker/docker/pkg/signal"
	swarmapi "github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/manager/encryption"
	swarmnode "github.com/docker/swarmkit/node"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// Init initializes new cluster from user provided request.
func (c *Cluster) Init(req types.InitRequest) (string, error) ***REMOVED***
	c.controlMutex.Lock()
	defer c.controlMutex.Unlock()
	if c.nr != nil ***REMOVED***
		if req.ForceNewCluster ***REMOVED***

			// Take c.mu temporarily to wait for presently running
			// API handlers to finish before shutting down the node.
			c.mu.Lock()
			if !c.nr.nodeState.IsManager() ***REMOVED***
				return "", errSwarmNotManager
			***REMOVED***
			c.mu.Unlock()

			if err := c.nr.Stop(); err != nil ***REMOVED***
				return "", err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			return "", errSwarmExists
		***REMOVED***
	***REMOVED***

	if err := validateAndSanitizeInitRequest(&req); err != nil ***REMOVED***
		return "", errdefs.InvalidParameter(err)
	***REMOVED***

	listenHost, listenPort, err := resolveListenAddr(req.ListenAddr)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	advertiseHost, advertisePort, err := c.resolveAdvertiseAddr(req.AdvertiseAddr, listenPort)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	dataPathAddr, err := resolveDataPathAddr(req.DataPathAddr)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	localAddr := listenHost

	// If the local address is undetermined, the advertise address
	// will be used as local address, if it belongs to this system.
	// If the advertise address is not local, then we try to find
	// a system address to use as local address. If this fails,
	// we give up and ask the user to pass the listen address.
	if net.ParseIP(localAddr).IsUnspecified() ***REMOVED***
		advertiseIP := net.ParseIP(advertiseHost)

		found := false
		for _, systemIP := range listSystemIPs() ***REMOVED***
			if systemIP.Equal(advertiseIP) ***REMOVED***
				localAddr = advertiseIP.String()
				found = true
				break
			***REMOVED***
		***REMOVED***

		if !found ***REMOVED***
			ip, err := c.resolveSystemAddr()
			if err != nil ***REMOVED***
				logrus.Warnf("Could not find a local address: %v", err)
				return "", errMustSpecifyListenAddr
			***REMOVED***
			localAddr = ip.String()
		***REMOVED***
	***REMOVED***

	nr, err := c.newNodeRunner(nodeStartConfig***REMOVED***
		forceNewCluster: req.ForceNewCluster,
		autolock:        req.AutoLockManagers,
		LocalAddr:       localAddr,
		ListenAddr:      net.JoinHostPort(listenHost, listenPort),
		AdvertiseAddr:   net.JoinHostPort(advertiseHost, advertisePort),
		DataPathAddr:    dataPathAddr,
		availability:    req.Availability,
	***REMOVED***)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	c.mu.Lock()
	c.nr = nr
	c.mu.Unlock()

	if err := <-nr.Ready(); err != nil ***REMOVED***
		c.mu.Lock()
		c.nr = nil
		c.mu.Unlock()
		if !req.ForceNewCluster ***REMOVED*** // if failure on first attempt don't keep state
			if err := clearPersistentState(c.root); err != nil ***REMOVED***
				return "", err
			***REMOVED***
		***REMOVED***
		return "", err
	***REMOVED***
	state := nr.State()
	if state.swarmNode == nil ***REMOVED*** // should never happen but protect from panic
		return "", errors.New("invalid cluster state for spec initialization")
	***REMOVED***
	if err := initClusterSpec(state.swarmNode, req.Spec); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return state.NodeID(), nil
***REMOVED***

// Join makes current Cluster part of an existing swarm cluster.
func (c *Cluster) Join(req types.JoinRequest) error ***REMOVED***
	c.controlMutex.Lock()
	defer c.controlMutex.Unlock()
	c.mu.Lock()
	if c.nr != nil ***REMOVED***
		c.mu.Unlock()
		return errors.WithStack(errSwarmExists)
	***REMOVED***
	c.mu.Unlock()

	if err := validateAndSanitizeJoinRequest(&req); err != nil ***REMOVED***
		return errdefs.InvalidParameter(err)
	***REMOVED***

	listenHost, listenPort, err := resolveListenAddr(req.ListenAddr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var advertiseAddr string
	if req.AdvertiseAddr != "" ***REMOVED***
		advertiseHost, advertisePort, err := c.resolveAdvertiseAddr(req.AdvertiseAddr, listenPort)
		// For joining, we don't need to provide an advertise address,
		// since the remote side can detect it.
		if err == nil ***REMOVED***
			advertiseAddr = net.JoinHostPort(advertiseHost, advertisePort)
		***REMOVED***
	***REMOVED***

	dataPathAddr, err := resolveDataPathAddr(req.DataPathAddr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	nr, err := c.newNodeRunner(nodeStartConfig***REMOVED***
		RemoteAddr:    req.RemoteAddrs[0],
		ListenAddr:    net.JoinHostPort(listenHost, listenPort),
		AdvertiseAddr: advertiseAddr,
		DataPathAddr:  dataPathAddr,
		joinAddr:      req.RemoteAddrs[0],
		joinToken:     req.JoinToken,
		availability:  req.Availability,
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c.mu.Lock()
	c.nr = nr
	c.mu.Unlock()

	select ***REMOVED***
	case <-time.After(swarmConnectTimeout):
		return errSwarmJoinTimeoutReached
	case err := <-nr.Ready():
		if err != nil ***REMOVED***
			c.mu.Lock()
			c.nr = nil
			c.mu.Unlock()
			if err := clearPersistentState(c.root); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		return err
	***REMOVED***
***REMOVED***

// Inspect retrieves the configuration properties of a managed swarm cluster.
func (c *Cluster) Inspect() (types.Swarm, error) ***REMOVED***
	var swarm types.Swarm
	if err := c.lockedManagerAction(func(ctx context.Context, state nodeState) error ***REMOVED***
		s, err := c.inspect(ctx, state)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		swarm = s
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return types.Swarm***REMOVED******REMOVED***, err
	***REMOVED***
	return swarm, nil
***REMOVED***

func (c *Cluster) inspect(ctx context.Context, state nodeState) (types.Swarm, error) ***REMOVED***
	s, err := getSwarm(ctx, state.controlClient)
	if err != nil ***REMOVED***
		return types.Swarm***REMOVED******REMOVED***, err
	***REMOVED***
	return convert.SwarmFromGRPC(*s), nil
***REMOVED***

// Update updates configuration of a managed swarm cluster.
func (c *Cluster) Update(version uint64, spec types.Spec, flags types.UpdateFlags) error ***REMOVED***
	return c.lockedManagerAction(func(ctx context.Context, state nodeState) error ***REMOVED***
		swarm, err := getSwarm(ctx, state.controlClient)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Validate spec name.
		if spec.Annotations.Name == "" ***REMOVED***
			spec.Annotations.Name = "default"
		***REMOVED*** else if spec.Annotations.Name != "default" ***REMOVED***
			return errdefs.InvalidParameter(errors.New(`swarm spec must be named "default"`))
		***REMOVED***

		// In update, client should provide the complete spec of the swarm, including
		// Name and Labels. If a field is specified with 0 or nil, then the default value
		// will be used to swarmkit.
		clusterSpec, err := convert.SwarmSpecToGRPC(spec)
		if err != nil ***REMOVED***
			return errdefs.InvalidParameter(err)
		***REMOVED***

		_, err = state.controlClient.UpdateCluster(
			ctx,
			&swarmapi.UpdateClusterRequest***REMOVED***
				ClusterID: swarm.ID,
				Spec:      &clusterSpec,
				ClusterVersion: &swarmapi.Version***REMOVED***
					Index: version,
				***REMOVED***,
				Rotation: swarmapi.KeyRotation***REMOVED***
					WorkerJoinToken:  flags.RotateWorkerToken,
					ManagerJoinToken: flags.RotateManagerToken,
					ManagerUnlockKey: flags.RotateManagerUnlockKey,
				***REMOVED***,
			***REMOVED***,
		)
		return err
	***REMOVED***)
***REMOVED***

// GetUnlockKey returns the unlock key for the swarm.
func (c *Cluster) GetUnlockKey() (string, error) ***REMOVED***
	var resp *swarmapi.GetUnlockKeyResponse
	if err := c.lockedManagerAction(func(ctx context.Context, state nodeState) error ***REMOVED***
		client := swarmapi.NewCAClient(state.grpcConn)

		r, err := client.GetUnlockKey(ctx, &swarmapi.GetUnlockKeyRequest***REMOVED******REMOVED***)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		resp = r
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if len(resp.UnlockKey) == 0 ***REMOVED***
		// no key
		return "", nil
	***REMOVED***
	return encryption.HumanReadableKey(resp.UnlockKey), nil
***REMOVED***

// UnlockSwarm provides a key to decrypt data that is encrypted at rest.
func (c *Cluster) UnlockSwarm(req types.UnlockRequest) error ***REMOVED***
	c.controlMutex.Lock()
	defer c.controlMutex.Unlock()

	c.mu.RLock()
	state := c.currentNodeState()

	if !state.IsActiveManager() ***REMOVED***
		// when manager is not active,
		// unless it is locked, otherwise return error.
		if err := c.errNoManager(state); err != errSwarmLocked ***REMOVED***
			c.mu.RUnlock()
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// when manager is active, return an error of "not locked"
		c.mu.RUnlock()
		return notLockedError***REMOVED******REMOVED***
	***REMOVED***

	// only when swarm is locked, code running reaches here
	nr := c.nr
	c.mu.RUnlock()

	key, err := encryption.ParseHumanReadableKey(req.UnlockKey)
	if err != nil ***REMOVED***
		return errdefs.InvalidParameter(err)
	***REMOVED***

	config := nr.config
	config.lockKey = key
	if err := nr.Stop(); err != nil ***REMOVED***
		return err
	***REMOVED***
	nr, err = c.newNodeRunner(config)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c.mu.Lock()
	c.nr = nr
	c.mu.Unlock()

	if err := <-nr.Ready(); err != nil ***REMOVED***
		if errors.Cause(err) == errSwarmLocked ***REMOVED***
			return invalidUnlockKey***REMOVED******REMOVED***
		***REMOVED***
		return errors.Errorf("swarm component could not be started: %v", err)
	***REMOVED***
	return nil
***REMOVED***

// Leave shuts down Cluster and removes current state.
func (c *Cluster) Leave(force bool) error ***REMOVED***
	c.controlMutex.Lock()
	defer c.controlMutex.Unlock()

	c.mu.Lock()
	nr := c.nr
	if nr == nil ***REMOVED***
		c.mu.Unlock()
		return errors.WithStack(errNoSwarm)
	***REMOVED***

	state := c.currentNodeState()

	c.mu.Unlock()

	if errors.Cause(state.err) == errSwarmLocked && !force ***REMOVED***
		// leave a locked swarm without --force is not allowed
		return errors.WithStack(notAvailableError("Swarm is encrypted and locked. Please unlock it first or use `--force` to ignore this message."))
	***REMOVED***

	if state.IsManager() && !force ***REMOVED***
		msg := "You are attempting to leave the swarm on a node that is participating as a manager. "
		if state.IsActiveManager() ***REMOVED***
			active, reachable, unreachable, err := managerStats(state.controlClient, state.NodeID())
			if err == nil ***REMOVED***
				if active && removingManagerCausesLossOfQuorum(reachable, unreachable) ***REMOVED***
					if isLastManager(reachable, unreachable) ***REMOVED***
						msg += "Removing the last manager erases all current state of the swarm. Use `--force` to ignore this message. "
						return errors.WithStack(notAvailableError(msg))
					***REMOVED***
					msg += fmt.Sprintf("Removing this node leaves %v managers out of %v. Without a Raft quorum your swarm will be inaccessible. ", reachable-1, reachable+unreachable)
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			msg += "Doing so may lose the consensus of your cluster. "
		***REMOVED***

		msg += "The only way to restore a swarm that has lost consensus is to reinitialize it with `--force-new-cluster`. Use `--force` to suppress this message."
		return errors.WithStack(notAvailableError(msg))
	***REMOVED***
	// release readers in here
	if err := nr.Stop(); err != nil ***REMOVED***
		logrus.Errorf("failed to shut down cluster node: %v", err)
		signal.DumpStacks("")
		return err
	***REMOVED***

	c.mu.Lock()
	c.nr = nil
	c.mu.Unlock()

	if nodeID := state.NodeID(); nodeID != "" ***REMOVED***
		nodeContainers, err := c.listContainerForNode(nodeID)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		for _, id := range nodeContainers ***REMOVED***
			if err := c.config.Backend.ContainerRm(id, &apitypes.ContainerRmConfig***REMOVED***ForceRemove: true***REMOVED***); err != nil ***REMOVED***
				logrus.Errorf("error removing %v: %v", id, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// todo: cleanup optional?
	if err := clearPersistentState(c.root); err != nil ***REMOVED***
		return err
	***REMOVED***
	c.config.Backend.DaemonLeavesCluster()
	return nil
***REMOVED***

// Info returns information about the current cluster state.
func (c *Cluster) Info() types.Info ***REMOVED***
	info := types.Info***REMOVED***
		NodeAddr: c.GetAdvertiseAddress(),
	***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := c.currentNodeState()
	info.LocalNodeState = state.status
	if state.err != nil ***REMOVED***
		info.Error = state.err.Error()
	***REMOVED***

	ctx, cancel := c.getRequestContext()
	defer cancel()

	if state.IsActiveManager() ***REMOVED***
		info.ControlAvailable = true
		swarm, err := c.inspect(ctx, state)
		if err != nil ***REMOVED***
			info.Error = err.Error()
		***REMOVED***

		info.Cluster = &swarm.ClusterInfo

		if r, err := state.controlClient.ListNodes(ctx, &swarmapi.ListNodesRequest***REMOVED******REMOVED***); err != nil ***REMOVED***
			info.Error = err.Error()
		***REMOVED*** else ***REMOVED***
			info.Nodes = len(r.Nodes)
			for _, n := range r.Nodes ***REMOVED***
				if n.ManagerStatus != nil ***REMOVED***
					info.Managers = info.Managers + 1
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if state.swarmNode != nil ***REMOVED***
		for _, r := range state.swarmNode.Remotes() ***REMOVED***
			info.RemoteManagers = append(info.RemoteManagers, types.Peer***REMOVED***NodeID: r.NodeID, Addr: r.Addr***REMOVED***)
		***REMOVED***
		info.NodeID = state.swarmNode.NodeID()
	***REMOVED***

	return info
***REMOVED***

func validateAndSanitizeInitRequest(req *types.InitRequest) error ***REMOVED***
	var err error
	req.ListenAddr, err = validateAddr(req.ListenAddr)
	if err != nil ***REMOVED***
		return fmt.Errorf("invalid ListenAddr %q: %v", req.ListenAddr, err)
	***REMOVED***

	if req.Spec.Annotations.Name == "" ***REMOVED***
		req.Spec.Annotations.Name = "default"
	***REMOVED*** else if req.Spec.Annotations.Name != "default" ***REMOVED***
		return errors.New(`swarm spec must be named "default"`)
	***REMOVED***

	return nil
***REMOVED***

func validateAndSanitizeJoinRequest(req *types.JoinRequest) error ***REMOVED***
	var err error
	req.ListenAddr, err = validateAddr(req.ListenAddr)
	if err != nil ***REMOVED***
		return fmt.Errorf("invalid ListenAddr %q: %v", req.ListenAddr, err)
	***REMOVED***
	if len(req.RemoteAddrs) == 0 ***REMOVED***
		return errors.New("at least 1 RemoteAddr is required to join")
	***REMOVED***
	for i := range req.RemoteAddrs ***REMOVED***
		req.RemoteAddrs[i], err = validateAddr(req.RemoteAddrs[i])
		if err != nil ***REMOVED***
			return fmt.Errorf("invalid remoteAddr %q: %v", req.RemoteAddrs[i], err)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func validateAddr(addr string) (string, error) ***REMOVED***
	if addr == "" ***REMOVED***
		return addr, errors.New("invalid empty address")
	***REMOVED***
	newaddr, err := opts.ParseTCPAddr(addr, defaultAddr)
	if err != nil ***REMOVED***
		return addr, nil
	***REMOVED***
	return strings.TrimPrefix(newaddr, "tcp://"), nil
***REMOVED***

func initClusterSpec(node *swarmnode.Node, spec types.Spec) error ***REMOVED***
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	for conn := range node.ListenControlSocket(ctx) ***REMOVED***
		if ctx.Err() != nil ***REMOVED***
			return ctx.Err()
		***REMOVED***
		if conn != nil ***REMOVED***
			client := swarmapi.NewControlClient(conn)
			var cluster *swarmapi.Cluster
			for i := 0; ; i++ ***REMOVED***
				lcr, err := client.ListClusters(ctx, &swarmapi.ListClustersRequest***REMOVED******REMOVED***)
				if err != nil ***REMOVED***
					return fmt.Errorf("error on listing clusters: %v", err)
				***REMOVED***
				if len(lcr.Clusters) == 0 ***REMOVED***
					if i < 10 ***REMOVED***
						time.Sleep(200 * time.Millisecond)
						continue
					***REMOVED***
					return errors.New("empty list of clusters was returned")
				***REMOVED***
				cluster = lcr.Clusters[0]
				break
			***REMOVED***
			// In init, we take the initial default values from swarmkit, and merge
			// any non nil or 0 value from spec to GRPC spec. This will leave the
			// default value alone.
			// Note that this is different from Update(), as in Update() we expect
			// user to specify the complete spec of the cluster (as they already know
			// the existing one and knows which field to update)
			clusterSpec, err := convert.MergeSwarmSpecToGRPC(spec, cluster.Spec)
			if err != nil ***REMOVED***
				return fmt.Errorf("error updating cluster settings: %v", err)
			***REMOVED***
			_, err = client.UpdateCluster(ctx, &swarmapi.UpdateClusterRequest***REMOVED***
				ClusterID:      cluster.ID,
				ClusterVersion: &cluster.Meta.Version,
				Spec:           &clusterSpec,
			***REMOVED***)
			if err != nil ***REMOVED***
				return fmt.Errorf("error updating cluster settings: %v", err)
			***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	return ctx.Err()
***REMOVED***

func (c *Cluster) listContainerForNode(nodeID string) ([]string, error) ***REMOVED***
	var ids []string
	filters := filters.NewArgs()
	filters.Add("label", fmt.Sprintf("com.docker.swarm.node.id=%s", nodeID))
	containers, err := c.config.Backend.Containers(&apitypes.ContainerListOptions***REMOVED***
		Filters: filters,
	***REMOVED***)
	if err != nil ***REMOVED***
		return []string***REMOVED******REMOVED***, err
	***REMOVED***
	for _, c := range containers ***REMOVED***
		ids = append(ids, c.ID)
	***REMOVED***
	return ids, nil
***REMOVED***
