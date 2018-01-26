/*
Package libnetwork provides the basic functionality and extension points to
create network namespaces and allocate interfaces for containers to use.

	networkType := "bridge"

	// Create a new controller instance
	driverOptions := options.Generic***REMOVED******REMOVED***
	genericOption := make(map[string]interface***REMOVED******REMOVED***)
	genericOption[netlabel.GenericData] = driverOptions
	controller, err := libnetwork.New(config.OptionDriverConfig(networkType, genericOption))
	if err != nil ***REMOVED***
		return
	***REMOVED***

	// Create a network for containers to join.
	// NewNetwork accepts Variadic optional arguments that libnetwork and Drivers can make use of
	network, err := controller.NewNetwork(networkType, "network1", "")
	if err != nil ***REMOVED***
		return
	***REMOVED***

	// For each new container: allocate IP and interfaces. The returned network
	// settings will be used for container infos (inspect and such), as well as
	// iptables rules for port publishing. This info is contained or accessible
	// from the returned endpoint.
	ep, err := network.CreateEndpoint("Endpoint1")
	if err != nil ***REMOVED***
		return
	***REMOVED***

	// Create the sandbox for the container.
	// NewSandbox accepts Variadic optional arguments which libnetwork can use.
	sbx, err := controller.NewSandbox("container1",
		libnetwork.OptionHostname("test"),
		libnetwork.OptionDomainname("docker.io"))

	// A sandbox can join the endpoint via the join api.
	err = ep.Join(sbx)
	if err != nil ***REMOVED***
		return
	***REMOVED***
*/
package libnetwork

import (
	"container/heap"
	"fmt"
	"net"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/pkg/discovery"
	"github.com/docker/docker/pkg/locker"
	"github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/docker/pkg/plugins"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/libnetwork/cluster"
	"github.com/docker/libnetwork/config"
	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/diagnose"
	"github.com/docker/libnetwork/discoverapi"
	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/drvregistry"
	"github.com/docker/libnetwork/hostdiscovery"
	"github.com/docker/libnetwork/ipamapi"
	"github.com/docker/libnetwork/netlabel"
	"github.com/docker/libnetwork/osl"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

// NetworkController provides the interface for controller instance which manages
// networks.
type NetworkController interface ***REMOVED***
	// ID provides a unique identity for the controller
	ID() string

	// BuiltinDrivers returns list of builtin drivers
	BuiltinDrivers() []string

	// BuiltinIPAMDrivers returns list of builtin ipam drivers
	BuiltinIPAMDrivers() []string

	// Config method returns the bootup configuration for the controller
	Config() config.Config

	// Create a new network. The options parameter carries network specific options.
	NewNetwork(networkType, name string, id string, options ...NetworkOption) (Network, error)

	// Networks returns the list of Network(s) managed by this controller.
	Networks() []Network

	// WalkNetworks uses the provided function to walk the Network(s) managed by this controller.
	WalkNetworks(walker NetworkWalker)

	// NetworkByName returns the Network which has the passed name. If not found, the error ErrNoSuchNetwork is returned.
	NetworkByName(name string) (Network, error)

	// NetworkByID returns the Network which has the passed id. If not found, the error ErrNoSuchNetwork is returned.
	NetworkByID(id string) (Network, error)

	// NewSandbox creates a new network sandbox for the passed container id
	NewSandbox(containerID string, options ...SandboxOption) (Sandbox, error)

	// Sandboxes returns the list of Sandbox(s) managed by this controller.
	Sandboxes() []Sandbox

	// WalkSandboxes uses the provided function to walk the Sandbox(s) managed by this controller.
	WalkSandboxes(walker SandboxWalker)

	// SandboxByID returns the Sandbox which has the passed id. If not found, a types.NotFoundError is returned.
	SandboxByID(id string) (Sandbox, error)

	// SandboxDestroy destroys a sandbox given a container ID
	SandboxDestroy(id string) error

	// Stop network controller
	Stop()

	// ReloadCondfiguration updates the controller configuration
	ReloadConfiguration(cfgOptions ...config.Option) error

	// SetClusterProvider sets cluster provider
	SetClusterProvider(provider cluster.Provider)

	// Wait for agent initialization complete in libnetwork controller
	AgentInitWait()

	// Wait for agent to stop if running
	AgentStopWait()

	// SetKeys configures the encryption key for gossip and overlay data path
	SetKeys(keys []*types.EncryptionKey) error

	// StartDiagnose start the network diagnose mode
	StartDiagnose(port int)
	// StopDiagnose start the network diagnose mode
	StopDiagnose()
	// IsDiagnoseEnabled returns true if the diagnose is enabled
	IsDiagnoseEnabled() bool
***REMOVED***

// NetworkWalker is a client provided function which will be used to walk the Networks.
// When the function returns true, the walk will stop.
type NetworkWalker func(nw Network) bool

// SandboxWalker is a client provided function which will be used to walk the Sandboxes.
// When the function returns true, the walk will stop.
type SandboxWalker func(sb Sandbox) bool

type sandboxTable map[string]*sandbox

type controller struct ***REMOVED***
	id                     string
	drvRegistry            *drvregistry.DrvRegistry
	sandboxes              sandboxTable
	cfg                    *config.Config
	stores                 []datastore.DataStore
	discovery              hostdiscovery.HostDiscovery
	extKeyListener         net.Listener
	watchCh                chan *endpoint
	unWatchCh              chan *endpoint
	svcRecords             map[string]svcInfo
	nmap                   map[string]*netWatch
	serviceBindings        map[serviceKey]*service
	defOsSbox              osl.Sandbox
	ingressSandbox         *sandbox
	sboxOnce               sync.Once
	agent                  *agent
	networkLocker          *locker.Locker
	agentInitDone          chan struct***REMOVED******REMOVED***
	agentStopDone          chan struct***REMOVED******REMOVED***
	keys                   []*types.EncryptionKey
	clusterConfigAvailable bool
	DiagnoseServer         *diagnose.Server
	sync.Mutex
***REMOVED***

type initializer struct ***REMOVED***
	fn    drvregistry.InitFunc
	ntype string
***REMOVED***

// New creates a new instance of network controller.
func New(cfgOptions ...config.Option) (NetworkController, error) ***REMOVED***
	c := &controller***REMOVED***
		id:              stringid.GenerateRandomID(),
		cfg:             config.ParseConfigOptions(cfgOptions...),
		sandboxes:       sandboxTable***REMOVED******REMOVED***,
		svcRecords:      make(map[string]svcInfo),
		serviceBindings: make(map[serviceKey]*service),
		agentInitDone:   make(chan struct***REMOVED******REMOVED***),
		networkLocker:   locker.New(),
		DiagnoseServer:  diagnose.New(),
	***REMOVED***
	c.DiagnoseServer.Init()

	if err := c.initStores(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	drvRegistry, err := drvregistry.New(c.getStore(datastore.LocalScope), c.getStore(datastore.GlobalScope), c.RegisterDriver, nil, c.cfg.PluginGetter)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	for _, i := range getInitializers(c.cfg.Daemon.Experimental) ***REMOVED***
		var dcfg map[string]interface***REMOVED******REMOVED***

		// External plugins don't need config passed through daemon. They can
		// bootstrap themselves
		if i.ntype != "remote" ***REMOVED***
			dcfg = c.makeDriverConfig(i.ntype)
		***REMOVED***

		if err := drvRegistry.AddDriver(i.ntype, i.fn, dcfg); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	if err = initIPAMDrivers(drvRegistry, nil, c.getStore(datastore.GlobalScope)); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	c.drvRegistry = drvRegistry

	if c.cfg != nil && c.cfg.Cluster.Watcher != nil ***REMOVED***
		if err := c.initDiscovery(c.cfg.Cluster.Watcher); err != nil ***REMOVED***
			// Failing to initialize discovery is a bad situation to be in.
			// But it cannot fail creating the Controller
			logrus.Errorf("Failed to Initialize Discovery : %v", err)
		***REMOVED***
	***REMOVED***

	c.WalkNetworks(populateSpecial)

	// Reserve pools first before doing cleanup. Otherwise the
	// cleanups of endpoint/network and sandbox below will
	// generate many unnecessary warnings
	c.reservePools()

	// Cleanup resources
	c.sandboxCleanup(c.cfg.ActiveSandboxes)
	c.cleanupLocalEndpoints()
	c.networkCleanup()

	if err := c.startExternalKeyListener(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return c, nil
***REMOVED***

func (c *controller) SetClusterProvider(provider cluster.Provider) ***REMOVED***
	var sameProvider bool
	c.Lock()
	// Avoids to spawn multiple goroutine for the same cluster provider
	if c.cfg.Daemon.ClusterProvider == provider ***REMOVED***
		// If the cluster provider is already set, there is already a go routine spawned
		// that is listening for events, so nothing to do here
		sameProvider = true
	***REMOVED*** else ***REMOVED***
		c.cfg.Daemon.ClusterProvider = provider
	***REMOVED***
	c.Unlock()

	if provider == nil || sameProvider ***REMOVED***
		return
	***REMOVED***
	// We don't want to spawn a new go routine if the previous one did not exit yet
	c.AgentStopWait()
	go c.clusterAgentInit()
***REMOVED***

func isValidClusteringIP(addr string) bool ***REMOVED***
	return addr != "" && !net.ParseIP(addr).IsLoopback() && !net.ParseIP(addr).IsUnspecified()
***REMOVED***

// libnetwork side of agent depends on the keys. On the first receipt of
// keys setup the agent. For subsequent key set handle the key change
func (c *controller) SetKeys(keys []*types.EncryptionKey) error ***REMOVED***
	subsysKeys := make(map[string]int)
	for _, key := range keys ***REMOVED***
		if key.Subsystem != subsysGossip &&
			key.Subsystem != subsysIPSec ***REMOVED***
			return fmt.Errorf("key received for unrecognized subsystem")
		***REMOVED***
		subsysKeys[key.Subsystem]++
	***REMOVED***
	for s, count := range subsysKeys ***REMOVED***
		if count != keyringSize ***REMOVED***
			return fmt.Errorf("incorrect number of keys for subsystem %v", s)
		***REMOVED***
	***REMOVED***

	agent := c.getAgent()

	if agent == nil ***REMOVED***
		c.Lock()
		c.keys = keys
		c.Unlock()
		return nil
	***REMOVED***
	return c.handleKeyChange(keys)
***REMOVED***

func (c *controller) getAgent() *agent ***REMOVED***
	c.Lock()
	defer c.Unlock()
	return c.agent
***REMOVED***

func (c *controller) clusterAgentInit() ***REMOVED***
	clusterProvider := c.cfg.Daemon.ClusterProvider
	var keysAvailable bool
	for ***REMOVED***
		eventType := <-clusterProvider.ListenClusterEvents()
		// The events: EventSocketChange, EventNodeReady and EventNetworkKeysAvailable are not ordered
		// when all the condition for the agent initialization are met then proceed with it
		switch eventType ***REMOVED***
		case cluster.EventNetworkKeysAvailable:
			// Validates that the keys are actually available before starting the initialization
			// This will handle old spurious messages left on the channel
			c.Lock()
			keysAvailable = c.keys != nil
			c.Unlock()
			fallthrough
		case cluster.EventSocketChange, cluster.EventNodeReady:
			if keysAvailable && !c.isDistributedControl() ***REMOVED***
				c.agentOperationStart()
				if err := c.agentSetup(clusterProvider); err != nil ***REMOVED***
					c.agentStopComplete()
				***REMOVED*** else ***REMOVED***
					c.agentInitComplete()
				***REMOVED***
			***REMOVED***
		case cluster.EventNodeLeave:
			keysAvailable = false
			c.agentOperationStart()
			c.Lock()
			c.keys = nil
			c.Unlock()

			// We are leaving the cluster. Make sure we
			// close the gossip so that we stop all
			// incoming gossip updates before cleaning up
			// any remaining service bindings. But before
			// deleting the networks since the networks
			// should still be present when cleaning up
			// service bindings
			c.agentClose()
			c.cleanupServiceDiscovery("")
			c.cleanupServiceBindings("")

			c.agentStopComplete()

			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// AgentInitWait waits for agent initialization to be completed in the controller.
func (c *controller) AgentInitWait() ***REMOVED***
	c.Lock()
	agentInitDone := c.agentInitDone
	c.Unlock()

	if agentInitDone != nil ***REMOVED***
		<-agentInitDone
	***REMOVED***
***REMOVED***

// AgentStopWait waits for the Agent stop to be completed in the controller
func (c *controller) AgentStopWait() ***REMOVED***
	c.Lock()
	agentStopDone := c.agentStopDone
	c.Unlock()
	if agentStopDone != nil ***REMOVED***
		<-agentStopDone
	***REMOVED***
***REMOVED***

// agentOperationStart marks the start of an Agent Init or Agent Stop
func (c *controller) agentOperationStart() ***REMOVED***
	c.Lock()
	if c.agentInitDone == nil ***REMOVED***
		c.agentInitDone = make(chan struct***REMOVED******REMOVED***)
	***REMOVED***
	if c.agentStopDone == nil ***REMOVED***
		c.agentStopDone = make(chan struct***REMOVED******REMOVED***)
	***REMOVED***
	c.Unlock()
***REMOVED***

// agentInitComplete notifies the successful completion of the Agent initialization
func (c *controller) agentInitComplete() ***REMOVED***
	c.Lock()
	if c.agentInitDone != nil ***REMOVED***
		close(c.agentInitDone)
		c.agentInitDone = nil
	***REMOVED***
	c.Unlock()
***REMOVED***

// agentStopComplete notifies the successful completion of the Agent stop
func (c *controller) agentStopComplete() ***REMOVED***
	c.Lock()
	if c.agentStopDone != nil ***REMOVED***
		close(c.agentStopDone)
		c.agentStopDone = nil
	***REMOVED***
	c.Unlock()
***REMOVED***

func (c *controller) makeDriverConfig(ntype string) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	if c.cfg == nil ***REMOVED***
		return nil
	***REMOVED***

	config := make(map[string]interface***REMOVED******REMOVED***)

	for _, label := range c.cfg.Daemon.Labels ***REMOVED***
		if !strings.HasPrefix(netlabel.Key(label), netlabel.DriverPrefix+"."+ntype) ***REMOVED***
			continue
		***REMOVED***

		config[netlabel.Key(label)] = netlabel.Value(label)
	***REMOVED***

	drvCfg, ok := c.cfg.Daemon.DriverCfg[ntype]
	if ok ***REMOVED***
		for k, v := range drvCfg.(map[string]interface***REMOVED******REMOVED***) ***REMOVED***
			config[k] = v
		***REMOVED***
	***REMOVED***

	for k, v := range c.cfg.Scopes ***REMOVED***
		if !v.IsValid() ***REMOVED***
			continue
		***REMOVED***
		config[netlabel.MakeKVClient(k)] = discoverapi.DatastoreConfigData***REMOVED***
			Scope:    k,
			Provider: v.Client.Provider,
			Address:  v.Client.Address,
			Config:   v.Client.Config,
		***REMOVED***
	***REMOVED***

	return config
***REMOVED***

var procReloadConfig = make(chan (bool), 1)

func (c *controller) ReloadConfiguration(cfgOptions ...config.Option) error ***REMOVED***
	procReloadConfig <- true
	defer func() ***REMOVED*** <-procReloadConfig ***REMOVED***()

	// For now we accept the configuration reload only as a mean to provide a global store config after boot.
	// Refuse the configuration if it alters an existing datastore client configuration.
	update := false
	cfg := config.ParseConfigOptions(cfgOptions...)

	for s := range c.cfg.Scopes ***REMOVED***
		if _, ok := cfg.Scopes[s]; !ok ***REMOVED***
			return types.ForbiddenErrorf("cannot accept new configuration because it removes an existing datastore client")
		***REMOVED***
	***REMOVED***
	for s, nSCfg := range cfg.Scopes ***REMOVED***
		if eSCfg, ok := c.cfg.Scopes[s]; ok ***REMOVED***
			if eSCfg.Client.Provider != nSCfg.Client.Provider ||
				eSCfg.Client.Address != nSCfg.Client.Address ***REMOVED***
				return types.ForbiddenErrorf("cannot accept new configuration because it modifies an existing datastore client")
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if err := c.initScopedStore(s, nSCfg); err != nil ***REMOVED***
				return err
			***REMOVED***
			update = true
		***REMOVED***
	***REMOVED***
	if !update ***REMOVED***
		return nil
	***REMOVED***

	c.Lock()
	c.cfg = cfg
	c.Unlock()

	var dsConfig *discoverapi.DatastoreConfigData
	for scope, sCfg := range cfg.Scopes ***REMOVED***
		if scope == datastore.LocalScope || !sCfg.IsValid() ***REMOVED***
			continue
		***REMOVED***
		dsConfig = &discoverapi.DatastoreConfigData***REMOVED***
			Scope:    scope,
			Provider: sCfg.Client.Provider,
			Address:  sCfg.Client.Address,
			Config:   sCfg.Client.Config,
		***REMOVED***
		break
	***REMOVED***
	if dsConfig == nil ***REMOVED***
		return nil
	***REMOVED***

	c.drvRegistry.WalkIPAMs(func(name string, driver ipamapi.Ipam, cap *ipamapi.Capability) bool ***REMOVED***
		err := driver.DiscoverNew(discoverapi.DatastoreConfig, *dsConfig)
		if err != nil ***REMOVED***
			logrus.Errorf("Failed to set datastore in driver %s: %v", name, err)
		***REMOVED***
		return false
	***REMOVED***)

	c.drvRegistry.WalkDrivers(func(name string, driver driverapi.Driver, capability driverapi.Capability) bool ***REMOVED***
		err := driver.DiscoverNew(discoverapi.DatastoreConfig, *dsConfig)
		if err != nil ***REMOVED***
			logrus.Errorf("Failed to set datastore in driver %s: %v", name, err)
		***REMOVED***
		return false
	***REMOVED***)

	if c.discovery == nil && c.cfg.Cluster.Watcher != nil ***REMOVED***
		if err := c.initDiscovery(c.cfg.Cluster.Watcher); err != nil ***REMOVED***
			logrus.Errorf("Failed to Initialize Discovery after configuration update: %v", err)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (c *controller) ID() string ***REMOVED***
	return c.id
***REMOVED***

func (c *controller) BuiltinDrivers() []string ***REMOVED***
	drivers := []string***REMOVED******REMOVED***
	c.drvRegistry.WalkDrivers(func(name string, driver driverapi.Driver, capability driverapi.Capability) bool ***REMOVED***
		if driver.IsBuiltIn() ***REMOVED***
			drivers = append(drivers, name)
		***REMOVED***
		return false
	***REMOVED***)
	return drivers
***REMOVED***

func (c *controller) BuiltinIPAMDrivers() []string ***REMOVED***
	drivers := []string***REMOVED******REMOVED***
	c.drvRegistry.WalkIPAMs(func(name string, driver ipamapi.Ipam, cap *ipamapi.Capability) bool ***REMOVED***
		if driver.IsBuiltIn() ***REMOVED***
			drivers = append(drivers, name)
		***REMOVED***
		return false
	***REMOVED***)
	return drivers
***REMOVED***

func (c *controller) validateHostDiscoveryConfig() bool ***REMOVED***
	if c.cfg == nil || c.cfg.Cluster.Discovery == "" || c.cfg.Cluster.Address == "" ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

func (c *controller) clusterHostID() string ***REMOVED***
	c.Lock()
	defer c.Unlock()
	if c.cfg == nil || c.cfg.Cluster.Address == "" ***REMOVED***
		return ""
	***REMOVED***
	addr := strings.Split(c.cfg.Cluster.Address, ":")
	return addr[0]
***REMOVED***

func (c *controller) isNodeAlive(node string) bool ***REMOVED***
	if c.discovery == nil ***REMOVED***
		return false
	***REMOVED***

	nodes := c.discovery.Fetch()
	for _, n := range nodes ***REMOVED***
		if n.String() == node ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

func (c *controller) initDiscovery(watcher discovery.Watcher) error ***REMOVED***
	if c.cfg == nil ***REMOVED***
		return fmt.Errorf("discovery initialization requires a valid configuration")
	***REMOVED***

	c.discovery = hostdiscovery.NewHostDiscovery(watcher)
	return c.discovery.Watch(c.activeCallback, c.hostJoinCallback, c.hostLeaveCallback)
***REMOVED***

func (c *controller) activeCallback() ***REMOVED***
	ds := c.getStore(datastore.GlobalScope)
	if ds != nil && !ds.Active() ***REMOVED***
		ds.RestartWatch()
	***REMOVED***
***REMOVED***

func (c *controller) hostJoinCallback(nodes []net.IP) ***REMOVED***
	c.processNodeDiscovery(nodes, true)
***REMOVED***

func (c *controller) hostLeaveCallback(nodes []net.IP) ***REMOVED***
	c.processNodeDiscovery(nodes, false)
***REMOVED***

func (c *controller) processNodeDiscovery(nodes []net.IP, add bool) ***REMOVED***
	c.drvRegistry.WalkDrivers(func(name string, driver driverapi.Driver, capability driverapi.Capability) bool ***REMOVED***
		c.pushNodeDiscovery(driver, capability, nodes, add)
		return false
	***REMOVED***)
***REMOVED***

func (c *controller) pushNodeDiscovery(d driverapi.Driver, cap driverapi.Capability, nodes []net.IP, add bool) ***REMOVED***
	var self net.IP
	if c.cfg != nil ***REMOVED***
		addr := strings.Split(c.cfg.Cluster.Address, ":")
		self = net.ParseIP(addr[0])
		// if external kvstore is not configured, try swarm-mode config
		if self == nil ***REMOVED***
			if agent := c.getAgent(); agent != nil ***REMOVED***
				self = net.ParseIP(agent.advertiseAddr)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if d == nil || cap.ConnectivityScope != datastore.GlobalScope || nodes == nil ***REMOVED***
		return
	***REMOVED***

	for _, node := range nodes ***REMOVED***
		nodeData := discoverapi.NodeDiscoveryData***REMOVED***Address: node.String(), Self: node.Equal(self)***REMOVED***
		var err error
		if add ***REMOVED***
			err = d.DiscoverNew(discoverapi.NodeDiscovery, nodeData)
		***REMOVED*** else ***REMOVED***
			err = d.DiscoverDelete(discoverapi.NodeDiscovery, nodeData)
		***REMOVED***
		if err != nil ***REMOVED***
			logrus.Debugf("discovery notification error: %v", err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *controller) Config() config.Config ***REMOVED***
	c.Lock()
	defer c.Unlock()
	if c.cfg == nil ***REMOVED***
		return config.Config***REMOVED******REMOVED***
	***REMOVED***
	return *c.cfg
***REMOVED***

func (c *controller) isManager() bool ***REMOVED***
	c.Lock()
	defer c.Unlock()
	if c.cfg == nil || c.cfg.Daemon.ClusterProvider == nil ***REMOVED***
		return false
	***REMOVED***
	return c.cfg.Daemon.ClusterProvider.IsManager()
***REMOVED***

func (c *controller) isAgent() bool ***REMOVED***
	c.Lock()
	defer c.Unlock()
	if c.cfg == nil || c.cfg.Daemon.ClusterProvider == nil ***REMOVED***
		return false
	***REMOVED***
	return c.cfg.Daemon.ClusterProvider.IsAgent()
***REMOVED***

func (c *controller) isDistributedControl() bool ***REMOVED***
	return !c.isManager() && !c.isAgent()
***REMOVED***

func (c *controller) GetPluginGetter() plugingetter.PluginGetter ***REMOVED***
	return c.drvRegistry.GetPluginGetter()
***REMOVED***

func (c *controller) RegisterDriver(networkType string, driver driverapi.Driver, capability driverapi.Capability) error ***REMOVED***
	c.Lock()
	hd := c.discovery
	c.Unlock()

	if hd != nil ***REMOVED***
		c.pushNodeDiscovery(driver, capability, hd.Fetch(), true)
	***REMOVED***

	c.agentDriverNotify(driver)
	return nil
***REMOVED***

// NewNetwork creates a new network of the specified network type. The options
// are network specific and modeled in a generic way.
func (c *controller) NewNetwork(networkType, name string, id string, options ...NetworkOption) (Network, error) ***REMOVED***
	if id != "" ***REMOVED***
		c.networkLocker.Lock(id)
		defer c.networkLocker.Unlock(id)

		if _, err := c.NetworkByID(id); err == nil ***REMOVED***
			return nil, NetworkNameError(id)
		***REMOVED***
	***REMOVED***

	if !config.IsValidName(name) ***REMOVED***
		return nil, ErrInvalidName(name)
	***REMOVED***

	if id == "" ***REMOVED***
		id = stringid.GenerateRandomID()
	***REMOVED***

	defaultIpam := defaultIpamForNetworkType(networkType)
	// Construct the network object
	network := &network***REMOVED***
		name:        name,
		networkType: networkType,
		generic:     map[string]interface***REMOVED******REMOVED******REMOVED***netlabel.GenericData: make(map[string]string)***REMOVED***,
		ipamType:    defaultIpam,
		id:          id,
		created:     time.Now(),
		ctrlr:       c,
		persist:     true,
		drvOnce:     &sync.Once***REMOVED******REMOVED***,
	***REMOVED***

	network.processOptions(options...)
	if err := network.validateConfiguration(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var (
		cap *driverapi.Capability
		err error
	)

	// Reset network types, force local scope and skip allocation and
	// plumbing for configuration networks. Reset of the config-only
	// network drivers is needed so that this special network is not
	// usable by old engine versions.
	if network.configOnly ***REMOVED***
		network.scope = datastore.LocalScope
		network.networkType = "null"
		goto addToStore
	***REMOVED***

	_, cap, err = network.resolveDriver(network.networkType, true)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if network.scope == datastore.LocalScope && cap.DataScope == datastore.GlobalScope ***REMOVED***
		return nil, types.ForbiddenErrorf("cannot downgrade network scope for %s networks", networkType)

	***REMOVED***
	if network.ingress && cap.DataScope != datastore.GlobalScope ***REMOVED***
		return nil, types.ForbiddenErrorf("Ingress network can only be global scope network")
	***REMOVED***

	// At this point the network scope is still unknown if not set by user
	if (cap.DataScope == datastore.GlobalScope || network.scope == datastore.SwarmScope) &&
		!c.isDistributedControl() && !network.dynamic ***REMOVED***
		if c.isManager() ***REMOVED***
			// For non-distributed controlled environment, globalscoped non-dynamic networks are redirected to Manager
			return nil, ManagerRedirectError(name)
		***REMOVED***
		return nil, types.ForbiddenErrorf("Cannot create a multi-host network from a worker node. Please create the network from a manager node.")
	***REMOVED***

	if network.scope == datastore.SwarmScope && c.isDistributedControl() ***REMOVED***
		return nil, types.ForbiddenErrorf("cannot create a swarm scoped network when swarm is not active")
	***REMOVED***

	// Make sure we have a driver available for this network type
	// before we allocate anything.
	if _, err := network.driver(true); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// From this point on, we need the network specific configuration,
	// which may come from a configuration-only network
	if network.configFrom != "" ***REMOVED***
		t, err := c.getConfigNetwork(network.configFrom)
		if err != nil ***REMOVED***
			return nil, types.NotFoundErrorf("configuration network %q does not exist", network.configFrom)
		***REMOVED***
		if err := t.applyConfigurationTo(network); err != nil ***REMOVED***
			return nil, types.InternalErrorf("Failed to apply configuration: %v", err)
		***REMOVED***
		defer func() ***REMOVED***
			if err == nil ***REMOVED***
				if err := t.getEpCnt().IncEndpointCnt(); err != nil ***REMOVED***
					logrus.Warnf("Failed to update reference count for configuration network %q on creation of network %q: %v",
						t.Name(), network.Name(), err)
				***REMOVED***
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	err = network.ipamAllocate()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			network.ipamRelease()
		***REMOVED***
	***REMOVED***()

	err = c.addNetwork(network)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			if e := network.deleteNetwork(); e != nil ***REMOVED***
				logrus.Warnf("couldn't roll back driver network on network %s creation failure: %v", network.name, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

addToStore:
	// First store the endpoint count, then the network. To avoid to
	// end up with a datastore containing a network and not an epCnt,
	// in case of an ungraceful shutdown during this function call.
	epCnt := &endpointCnt***REMOVED***n: network***REMOVED***
	if err = c.updateToStore(epCnt); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			if e := c.deleteFromStore(epCnt); e != nil ***REMOVED***
				logrus.Warnf("could not rollback from store, epCnt %v on failure (%v): %v", epCnt, err, e)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	network.epCnt = epCnt
	if err = c.updateToStore(network); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			if e := c.deleteFromStore(network); e != nil ***REMOVED***
				logrus.Warnf("could not rollback from store, network %v on failure (%v): %v", network, err, e)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	if network.configOnly ***REMOVED***
		return network, nil
	***REMOVED***

	joinCluster(network)
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			network.cancelDriverWatches()
			if e := network.leaveCluster(); e != nil ***REMOVED***
				logrus.Warnf("Failed to leave agent cluster on network %s on failure (%v): %v", network.name, err, e)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	if len(network.loadBalancerIP) != 0 ***REMOVED***
		if err = network.createLoadBalancerSandbox(); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	if !c.isDistributedControl() ***REMOVED***
		c.Lock()
		arrangeIngressFilterRule()
		c.Unlock()
	***REMOVED***

	c.arrangeUserFilterRule()

	return network, nil
***REMOVED***

var joinCluster NetworkWalker = func(nw Network) bool ***REMOVED***
	n := nw.(*network)
	if n.configOnly ***REMOVED***
		return false
	***REMOVED***
	if err := n.joinCluster(); err != nil ***REMOVED***
		logrus.Errorf("Failed to join network %s (%s) into agent cluster: %v", n.Name(), n.ID(), err)
	***REMOVED***
	n.addDriverWatches()
	return false
***REMOVED***

func (c *controller) reservePools() ***REMOVED***
	networks, err := c.getNetworksForScope(datastore.LocalScope)
	if err != nil ***REMOVED***
		logrus.Warnf("Could not retrieve networks from local store during ipam allocation for existing networks: %v", err)
		return
	***REMOVED***

	for _, n := range networks ***REMOVED***
		if n.configOnly ***REMOVED***
			continue
		***REMOVED***
		if !doReplayPoolReserve(n) ***REMOVED***
			continue
		***REMOVED***
		// Construct pseudo configs for the auto IP case
		autoIPv4 := (len(n.ipamV4Config) == 0 || (len(n.ipamV4Config) == 1 && n.ipamV4Config[0].PreferredPool == "")) && len(n.ipamV4Info) > 0
		autoIPv6 := (len(n.ipamV6Config) == 0 || (len(n.ipamV6Config) == 1 && n.ipamV6Config[0].PreferredPool == "")) && len(n.ipamV6Info) > 0
		if autoIPv4 ***REMOVED***
			n.ipamV4Config = []*IpamConf***REMOVED******REMOVED***PreferredPool: n.ipamV4Info[0].Pool.String()***REMOVED******REMOVED***
		***REMOVED***
		if n.enableIPv6 && autoIPv6 ***REMOVED***
			n.ipamV6Config = []*IpamConf***REMOVED******REMOVED***PreferredPool: n.ipamV6Info[0].Pool.String()***REMOVED******REMOVED***
		***REMOVED***
		// Account current network gateways
		for i, c := range n.ipamV4Config ***REMOVED***
			if c.Gateway == "" && n.ipamV4Info[i].Gateway != nil ***REMOVED***
				c.Gateway = n.ipamV4Info[i].Gateway.IP.String()
			***REMOVED***
		***REMOVED***
		if n.enableIPv6 ***REMOVED***
			for i, c := range n.ipamV6Config ***REMOVED***
				if c.Gateway == "" && n.ipamV6Info[i].Gateway != nil ***REMOVED***
					c.Gateway = n.ipamV6Info[i].Gateway.IP.String()
				***REMOVED***
			***REMOVED***
		***REMOVED***
		// Reserve pools
		if err := n.ipamAllocate(); err != nil ***REMOVED***
			logrus.Warnf("Failed to allocate ipam pool(s) for network %q (%s): %v", n.Name(), n.ID(), err)
		***REMOVED***
		// Reserve existing endpoints' addresses
		ipam, _, err := n.getController().getIPAMDriver(n.ipamType)
		if err != nil ***REMOVED***
			logrus.Warnf("Failed to retrieve ipam driver for network %q (%s) during address reservation", n.Name(), n.ID())
			continue
		***REMOVED***
		epl, err := n.getEndpointsFromStore()
		if err != nil ***REMOVED***
			logrus.Warnf("Failed to retrieve list of current endpoints on network %q (%s)", n.Name(), n.ID())
			continue
		***REMOVED***
		for _, ep := range epl ***REMOVED***
			if err := ep.assignAddress(ipam, true, ep.Iface().AddressIPv6() != nil); err != nil ***REMOVED***
				logrus.Warnf("Failed to reserve current address for endpoint %q (%s) on network %q (%s)",
					ep.Name(), ep.ID(), n.Name(), n.ID())
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func doReplayPoolReserve(n *network) bool ***REMOVED***
	_, caps, err := n.getController().getIPAMDriver(n.ipamType)
	if err != nil ***REMOVED***
		logrus.Warnf("Failed to retrieve ipam driver for network %q (%s): %v", n.Name(), n.ID(), err)
		return false
	***REMOVED***
	return caps.RequiresRequestReplay
***REMOVED***

func (c *controller) addNetwork(n *network) error ***REMOVED***
	d, err := n.driver(true)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Create the network
	if err := d.CreateNetwork(n.id, n.generic, n, n.getIPData(4), n.getIPData(6)); err != nil ***REMOVED***
		return err
	***REMOVED***

	n.startResolver()

	return nil
***REMOVED***

func (c *controller) Networks() []Network ***REMOVED***
	var list []Network

	networks, err := c.getNetworksFromStore()
	if err != nil ***REMOVED***
		logrus.Error(err)
	***REMOVED***

	for _, n := range networks ***REMOVED***
		if n.inDelete ***REMOVED***
			continue
		***REMOVED***
		list = append(list, n)
	***REMOVED***

	return list
***REMOVED***

func (c *controller) WalkNetworks(walker NetworkWalker) ***REMOVED***
	for _, n := range c.Networks() ***REMOVED***
		if walker(n) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *controller) NetworkByName(name string) (Network, error) ***REMOVED***
	if name == "" ***REMOVED***
		return nil, ErrInvalidName(name)
	***REMOVED***
	var n Network

	s := func(current Network) bool ***REMOVED***
		if current.Name() == name ***REMOVED***
			n = current
			return true
		***REMOVED***
		return false
	***REMOVED***

	c.WalkNetworks(s)

	if n == nil ***REMOVED***
		return nil, ErrNoSuchNetwork(name)
	***REMOVED***

	return n, nil
***REMOVED***

func (c *controller) NetworkByID(id string) (Network, error) ***REMOVED***
	if id == "" ***REMOVED***
		return nil, ErrInvalidID(id)
	***REMOVED***

	n, err := c.getNetworkFromStore(id)
	if err != nil ***REMOVED***
		return nil, ErrNoSuchNetwork(id)
	***REMOVED***

	return n, nil
***REMOVED***

// NewSandbox creates a new sandbox for the passed container id
func (c *controller) NewSandbox(containerID string, options ...SandboxOption) (Sandbox, error) ***REMOVED***
	if containerID == "" ***REMOVED***
		return nil, types.BadRequestErrorf("invalid container ID")
	***REMOVED***

	var sb *sandbox
	c.Lock()
	for _, s := range c.sandboxes ***REMOVED***
		if s.containerID == containerID ***REMOVED***
			// If not a stub, then we already have a complete sandbox.
			if !s.isStub ***REMOVED***
				sbID := s.ID()
				c.Unlock()
				return nil, types.ForbiddenErrorf("container %s is already present in sandbox %s", containerID, sbID)
			***REMOVED***

			// We already have a stub sandbox from the
			// store. Make use of it so that we don't lose
			// the endpoints from store but reset the
			// isStub flag.
			sb = s
			sb.isStub = false
			break
		***REMOVED***
	***REMOVED***
	c.Unlock()

	// Create sandbox and process options first. Key generation depends on an option
	if sb == nil ***REMOVED***
		sb = &sandbox***REMOVED***
			id:                 stringid.GenerateRandomID(),
			containerID:        containerID,
			endpoints:          epHeap***REMOVED******REMOVED***,
			epPriority:         map[string]int***REMOVED******REMOVED***,
			populatedEndpoints: map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			config:             containerConfig***REMOVED******REMOVED***,
			controller:         c,
			extDNS:             []extDNSEntry***REMOVED******REMOVED***,
		***REMOVED***
	***REMOVED***

	heap.Init(&sb.endpoints)

	sb.processOptions(options...)

	c.Lock()
	if sb.ingress && c.ingressSandbox != nil ***REMOVED***
		c.Unlock()
		return nil, types.ForbiddenErrorf("ingress sandbox already present")
	***REMOVED***

	if sb.ingress ***REMOVED***
		c.ingressSandbox = sb
		sb.config.hostsPath = filepath.Join(c.cfg.Daemon.DataDir, "/network/files/hosts")
		sb.config.resolvConfPath = filepath.Join(c.cfg.Daemon.DataDir, "/network/files/resolv.conf")
		sb.id = "ingress_sbox"
	***REMOVED***
	c.Unlock()

	var err error
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			c.Lock()
			if sb.ingress ***REMOVED***
				c.ingressSandbox = nil
			***REMOVED***
			c.Unlock()
		***REMOVED***
	***REMOVED***()

	if err = sb.setupResolutionFiles(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if sb.config.useDefaultSandBox ***REMOVED***
		c.sboxOnce.Do(func() ***REMOVED***
			c.defOsSbox, err = osl.NewSandbox(sb.Key(), false, false)
		***REMOVED***)

		if err != nil ***REMOVED***
			c.sboxOnce = sync.Once***REMOVED******REMOVED***
			return nil, fmt.Errorf("failed to create default sandbox: %v", err)
		***REMOVED***

		sb.osSbox = c.defOsSbox
	***REMOVED***

	if sb.osSbox == nil && !sb.config.useExternalKey ***REMOVED***
		if sb.osSbox, err = osl.NewSandbox(sb.Key(), !sb.config.useDefaultSandBox, false); err != nil ***REMOVED***
			return nil, fmt.Errorf("failed to create new osl sandbox: %v", err)
		***REMOVED***
	***REMOVED***

	c.Lock()
	c.sandboxes[sb.id] = sb
	c.Unlock()
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			c.Lock()
			delete(c.sandboxes, sb.id)
			c.Unlock()
		***REMOVED***
	***REMOVED***()

	err = sb.storeUpdate()
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to update the store state of sandbox: %v", err)
	***REMOVED***

	return sb, nil
***REMOVED***

func (c *controller) Sandboxes() []Sandbox ***REMOVED***
	c.Lock()
	defer c.Unlock()

	list := make([]Sandbox, 0, len(c.sandboxes))
	for _, s := range c.sandboxes ***REMOVED***
		// Hide stub sandboxes from libnetwork users
		if s.isStub ***REMOVED***
			continue
		***REMOVED***

		list = append(list, s)
	***REMOVED***

	return list
***REMOVED***

func (c *controller) WalkSandboxes(walker SandboxWalker) ***REMOVED***
	for _, sb := range c.Sandboxes() ***REMOVED***
		if walker(sb) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *controller) SandboxByID(id string) (Sandbox, error) ***REMOVED***
	if id == "" ***REMOVED***
		return nil, ErrInvalidID(id)
	***REMOVED***
	c.Lock()
	s, ok := c.sandboxes[id]
	c.Unlock()
	if !ok ***REMOVED***
		return nil, types.NotFoundErrorf("sandbox %s not found", id)
	***REMOVED***
	return s, nil
***REMOVED***

// SandboxDestroy destroys a sandbox given a container ID
func (c *controller) SandboxDestroy(id string) error ***REMOVED***
	var sb *sandbox
	c.Lock()
	for _, s := range c.sandboxes ***REMOVED***
		if s.containerID == id ***REMOVED***
			sb = s
			break
		***REMOVED***
	***REMOVED***
	c.Unlock()

	// It is not an error if sandbox is not available
	if sb == nil ***REMOVED***
		return nil
	***REMOVED***

	return sb.Delete()
***REMOVED***

// SandboxContainerWalker returns a Sandbox Walker function which looks for an existing Sandbox with the passed containerID
func SandboxContainerWalker(out *Sandbox, containerID string) SandboxWalker ***REMOVED***
	return func(sb Sandbox) bool ***REMOVED***
		if sb.ContainerID() == containerID ***REMOVED***
			*out = sb
			return true
		***REMOVED***
		return false
	***REMOVED***
***REMOVED***

// SandboxKeyWalker returns a Sandbox Walker function which looks for an existing Sandbox with the passed key
func SandboxKeyWalker(out *Sandbox, key string) SandboxWalker ***REMOVED***
	return func(sb Sandbox) bool ***REMOVED***
		if sb.Key() == key ***REMOVED***
			*out = sb
			return true
		***REMOVED***
		return false
	***REMOVED***
***REMOVED***

func (c *controller) loadDriver(networkType string) error ***REMOVED***
	var err error

	if pg := c.GetPluginGetter(); pg != nil ***REMOVED***
		_, err = pg.Get(networkType, driverapi.NetworkPluginEndpointType, plugingetter.Lookup)
	***REMOVED*** else ***REMOVED***
		_, err = plugins.Get(networkType, driverapi.NetworkPluginEndpointType)
	***REMOVED***

	if err != nil ***REMOVED***
		if err == plugins.ErrNotFound ***REMOVED***
			return types.NotFoundErrorf(err.Error())
		***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (c *controller) loadIPAMDriver(name string) error ***REMOVED***
	var err error

	if pg := c.GetPluginGetter(); pg != nil ***REMOVED***
		_, err = pg.Get(name, ipamapi.PluginEndpointType, plugingetter.Lookup)
	***REMOVED*** else ***REMOVED***
		_, err = plugins.Get(name, ipamapi.PluginEndpointType)
	***REMOVED***

	if err != nil ***REMOVED***
		if err == plugins.ErrNotFound ***REMOVED***
			return types.NotFoundErrorf(err.Error())
		***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (c *controller) getIPAMDriver(name string) (ipamapi.Ipam, *ipamapi.Capability, error) ***REMOVED***
	id, cap := c.drvRegistry.IPAM(name)
	if id == nil ***REMOVED***
		// Might be a plugin name. Try loading it
		if err := c.loadIPAMDriver(name); err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***

		// Now that we resolved the plugin, try again looking up the registry
		id, cap = c.drvRegistry.IPAM(name)
		if id == nil ***REMOVED***
			return nil, nil, types.BadRequestErrorf("invalid ipam driver: %q", name)
		***REMOVED***
	***REMOVED***

	return id, cap, nil
***REMOVED***

func (c *controller) Stop() ***REMOVED***
	c.closeStores()
	c.stopExternalKeyListener()
	osl.GC()
***REMOVED***

// StartDiagnose start the network diagnose mode
func (c *controller) StartDiagnose(port int) ***REMOVED***
	c.Lock()
	if !c.DiagnoseServer.IsDebugEnable() ***REMOVED***
		c.DiagnoseServer.EnableDebug("127.0.0.1", port)
	***REMOVED***
	c.Unlock()
***REMOVED***

// StopDiagnose start the network diagnose mode
func (c *controller) StopDiagnose() ***REMOVED***
	c.Lock()
	if c.DiagnoseServer.IsDebugEnable() ***REMOVED***
		c.DiagnoseServer.DisableDebug()
	***REMOVED***
	c.Unlock()
***REMOVED***

// IsDiagnoseEnabled returns true if the diagnose is enabled
func (c *controller) IsDiagnoseEnabled() bool ***REMOVED***
	c.Lock()
	defer c.Unlock()
	return c.DiagnoseServer.IsDebugEnable()
***REMOVED***
