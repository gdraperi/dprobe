package bridge

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"

	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/discoverapi"
	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/iptables"
	"github.com/docker/libnetwork/netlabel"
	"github.com/docker/libnetwork/netutils"
	"github.com/docker/libnetwork/ns"
	"github.com/docker/libnetwork/options"
	"github.com/docker/libnetwork/osl"
	"github.com/docker/libnetwork/portmapper"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

const (
	networkType                = "bridge"
	vethPrefix                 = "veth"
	vethLen                    = 7
	defaultContainerVethPrefix = "eth"
	maxAllocatePortAttempts    = 10
)

const (
	// DefaultGatewayV4AuxKey represents the default-gateway configured by the user
	DefaultGatewayV4AuxKey = "DefaultGatewayIPv4"
	// DefaultGatewayV6AuxKey represents the ipv6 default-gateway configured by the user
	DefaultGatewayV6AuxKey = "DefaultGatewayIPv6"
)

type defaultBridgeNetworkConflict struct ***REMOVED***
	ID string
***REMOVED***

func (d defaultBridgeNetworkConflict) Error() string ***REMOVED***
	return fmt.Sprintf("Stale default bridge network %s", d.ID)
***REMOVED***

type iptableCleanFunc func() error
type iptablesCleanFuncs []iptableCleanFunc

// configuration info for the "bridge" driver.
type configuration struct ***REMOVED***
	EnableIPForwarding  bool
	EnableIPTables      bool
	EnableUserlandProxy bool
	UserlandProxyPath   string
***REMOVED***

// networkConfiguration for network specific configuration
type networkConfiguration struct ***REMOVED***
	ID                   string
	BridgeName           string
	EnableIPv6           bool
	EnableIPMasquerade   bool
	EnableICC            bool
	Mtu                  int
	DefaultBindingIP     net.IP
	DefaultBridge        bool
	ContainerIfacePrefix string
	// Internal fields set after ipam data parsing
	AddressIPv4        *net.IPNet
	AddressIPv6        *net.IPNet
	DefaultGatewayIPv4 net.IP
	DefaultGatewayIPv6 net.IP
	dbIndex            uint64
	dbExists           bool
	Internal           bool

	BridgeIfaceCreator ifaceCreator
***REMOVED***

// ifaceCreator represents how the bridge interface was created
type ifaceCreator int8

const (
	ifaceCreatorUnknown ifaceCreator = iota
	ifaceCreatedByLibnetwork
	ifaceCreatedByUser
)

// endpointConfiguration represents the user specified configuration for the sandbox endpoint
type endpointConfiguration struct ***REMOVED***
	MacAddress net.HardwareAddr
***REMOVED***

// containerConfiguration represents the user specified configuration for a container
type containerConfiguration struct ***REMOVED***
	ParentEndpoints []string
	ChildEndpoints  []string
***REMOVED***

// cnnectivityConfiguration represents the user specified configuration regarding the external connectivity
type connectivityConfiguration struct ***REMOVED***
	PortBindings []types.PortBinding
	ExposedPorts []types.TransportPort
***REMOVED***

type bridgeEndpoint struct ***REMOVED***
	id              string
	nid             string
	srcName         string
	addr            *net.IPNet
	addrv6          *net.IPNet
	macAddress      net.HardwareAddr
	config          *endpointConfiguration // User specified parameters
	containerConfig *containerConfiguration
	extConnConfig   *connectivityConfiguration
	portMapping     []types.PortBinding // Operation port bindings
	dbIndex         uint64
	dbExists        bool
***REMOVED***

type bridgeNetwork struct ***REMOVED***
	id            string
	bridge        *bridgeInterface // The bridge's L3 interface
	config        *networkConfiguration
	endpoints     map[string]*bridgeEndpoint // key: endpoint id
	portMapper    *portmapper.PortMapper
	driver        *driver // The network's driver
	iptCleanFuncs iptablesCleanFuncs
	sync.Mutex
***REMOVED***

type driver struct ***REMOVED***
	config         *configuration
	network        *bridgeNetwork
	natChain       *iptables.ChainInfo
	filterChain    *iptables.ChainInfo
	isolationChain *iptables.ChainInfo
	networks       map[string]*bridgeNetwork
	store          datastore.DataStore
	nlh            *netlink.Handle
	configNetwork  sync.Mutex
	sync.Mutex
***REMOVED***

// New constructs a new bridge driver
func newDriver() *driver ***REMOVED***
	return &driver***REMOVED***networks: map[string]*bridgeNetwork***REMOVED******REMOVED***, config: &configuration***REMOVED******REMOVED******REMOVED***
***REMOVED***

// Init registers a new instance of bridge driver
func Init(dc driverapi.DriverCallback, config map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	d := newDriver()
	if err := d.configure(config); err != nil ***REMOVED***
		return err
	***REMOVED***

	c := driverapi.Capability***REMOVED***
		DataScope:         datastore.LocalScope,
		ConnectivityScope: datastore.LocalScope,
	***REMOVED***
	return dc.RegisterDriver(networkType, d, c)
***REMOVED***

// Validate performs a static validation on the network configuration parameters.
// Whatever can be assessed a priori before attempting any programming.
func (c *networkConfiguration) Validate() error ***REMOVED***
	if c.Mtu < 0 ***REMOVED***
		return ErrInvalidMtu(c.Mtu)
	***REMOVED***

	// If bridge v4 subnet is specified
	if c.AddressIPv4 != nil ***REMOVED***
		// If default gw is specified, it must be part of bridge subnet
		if c.DefaultGatewayIPv4 != nil ***REMOVED***
			if !c.AddressIPv4.Contains(c.DefaultGatewayIPv4) ***REMOVED***
				return &ErrInvalidGateway***REMOVED******REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// If default v6 gw is specified, AddressIPv6 must be specified and gw must belong to AddressIPv6 subnet
	if c.EnableIPv6 && c.DefaultGatewayIPv6 != nil ***REMOVED***
		if c.AddressIPv6 == nil || !c.AddressIPv6.Contains(c.DefaultGatewayIPv6) ***REMOVED***
			return &ErrInvalidGateway***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Conflicts check if two NetworkConfiguration objects overlap
func (c *networkConfiguration) Conflicts(o *networkConfiguration) error ***REMOVED***
	if o == nil ***REMOVED***
		return errors.New("same configuration")
	***REMOVED***

	// Also empty, because only one network with empty name is allowed
	if c.BridgeName == o.BridgeName ***REMOVED***
		return errors.New("networks have same bridge name")
	***REMOVED***

	// They must be in different subnets
	if (c.AddressIPv4 != nil && o.AddressIPv4 != nil) &&
		(c.AddressIPv4.Contains(o.AddressIPv4.IP) || o.AddressIPv4.Contains(c.AddressIPv4.IP)) ***REMOVED***
		return errors.New("networks have overlapping IPv4")
	***REMOVED***

	// They must be in different v6 subnets
	if (c.AddressIPv6 != nil && o.AddressIPv6 != nil) &&
		(c.AddressIPv6.Contains(o.AddressIPv6.IP) || o.AddressIPv6.Contains(c.AddressIPv6.IP)) ***REMOVED***
		return errors.New("networks have overlapping IPv6")
	***REMOVED***

	return nil
***REMOVED***

func (c *networkConfiguration) fromLabels(labels map[string]string) error ***REMOVED***
	var err error
	for label, value := range labels ***REMOVED***
		switch label ***REMOVED***
		case BridgeName:
			c.BridgeName = value
		case netlabel.DriverMTU:
			if c.Mtu, err = strconv.Atoi(value); err != nil ***REMOVED***
				return parseErr(label, value, err.Error())
			***REMOVED***
		case netlabel.EnableIPv6:
			if c.EnableIPv6, err = strconv.ParseBool(value); err != nil ***REMOVED***
				return parseErr(label, value, err.Error())
			***REMOVED***
		case EnableIPMasquerade:
			if c.EnableIPMasquerade, err = strconv.ParseBool(value); err != nil ***REMOVED***
				return parseErr(label, value, err.Error())
			***REMOVED***
		case EnableICC:
			if c.EnableICC, err = strconv.ParseBool(value); err != nil ***REMOVED***
				return parseErr(label, value, err.Error())
			***REMOVED***
		case DefaultBridge:
			if c.DefaultBridge, err = strconv.ParseBool(value); err != nil ***REMOVED***
				return parseErr(label, value, err.Error())
			***REMOVED***
		case DefaultBindingIP:
			if c.DefaultBindingIP = net.ParseIP(value); c.DefaultBindingIP == nil ***REMOVED***
				return parseErr(label, value, "nil ip")
			***REMOVED***
		case netlabel.ContainerIfacePrefix:
			c.ContainerIfacePrefix = value
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func parseErr(label, value, errString string) error ***REMOVED***
	return types.BadRequestErrorf("failed to parse %s value: %v (%s)", label, value, errString)
***REMOVED***

func (n *bridgeNetwork) registerIptCleanFunc(clean iptableCleanFunc) ***REMOVED***
	n.iptCleanFuncs = append(n.iptCleanFuncs, clean)
***REMOVED***

func (n *bridgeNetwork) getDriverChains() (*iptables.ChainInfo, *iptables.ChainInfo, *iptables.ChainInfo, error) ***REMOVED***
	n.Lock()
	defer n.Unlock()

	if n.driver == nil ***REMOVED***
		return nil, nil, nil, types.BadRequestErrorf("no driver found")
	***REMOVED***

	return n.driver.natChain, n.driver.filterChain, n.driver.isolationChain, nil
***REMOVED***

func (n *bridgeNetwork) getNetworkBridgeName() string ***REMOVED***
	n.Lock()
	config := n.config
	n.Unlock()

	return config.BridgeName
***REMOVED***

func (n *bridgeNetwork) getEndpoint(eid string) (*bridgeEndpoint, error) ***REMOVED***
	n.Lock()
	defer n.Unlock()

	if eid == "" ***REMOVED***
		return nil, InvalidEndpointIDError(eid)
	***REMOVED***

	if ep, ok := n.endpoints[eid]; ok ***REMOVED***
		return ep, nil
	***REMOVED***

	return nil, nil
***REMOVED***

// Install/Removes the iptables rules needed to isolate this network
// from each of the other networks
func (n *bridgeNetwork) isolateNetwork(others []*bridgeNetwork, enable bool) error ***REMOVED***
	n.Lock()
	thisConfig := n.config
	n.Unlock()

	if thisConfig.Internal ***REMOVED***
		return nil
	***REMOVED***

	// Install the rules to isolate this networks against each of the other networks
	for _, o := range others ***REMOVED***
		o.Lock()
		otherConfig := o.config
		o.Unlock()

		if otherConfig.Internal ***REMOVED***
			continue
		***REMOVED***

		if thisConfig.BridgeName != otherConfig.BridgeName ***REMOVED***
			if err := setINC(thisConfig.BridgeName, otherConfig.BridgeName, enable); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) configure(option map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	var (
		config         *configuration
		err            error
		natChain       *iptables.ChainInfo
		filterChain    *iptables.ChainInfo
		isolationChain *iptables.ChainInfo
	)

	genericData, ok := option[netlabel.GenericData]
	if !ok || genericData == nil ***REMOVED***
		return nil
	***REMOVED***

	switch opt := genericData.(type) ***REMOVED***
	case options.Generic:
		opaqueConfig, err := options.GenerateFromModel(opt, &configuration***REMOVED******REMOVED***)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		config = opaqueConfig.(*configuration)
	case *configuration:
		config = opt
	default:
		return &ErrInvalidDriverConfig***REMOVED******REMOVED***
	***REMOVED***

	if config.EnableIPTables ***REMOVED***
		if _, err := os.Stat("/proc/sys/net/bridge"); err != nil ***REMOVED***
			if out, err := exec.Command("modprobe", "-va", "bridge", "br_netfilter").CombinedOutput(); err != nil ***REMOVED***
				logrus.Warnf("Running modprobe bridge br_netfilter failed with message: %s, error: %v", out, err)
			***REMOVED***
		***REMOVED***
		removeIPChains()
		natChain, filterChain, isolationChain, err = setupIPChains(config)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		// Make sure on firewall reload, first thing being re-played is chains creation
		iptables.OnReloaded(func() ***REMOVED*** logrus.Debugf("Recreating iptables chains on firewall reload"); setupIPChains(config) ***REMOVED***)
	***REMOVED***

	if config.EnableIPForwarding ***REMOVED***
		err = setupIPForwarding(config.EnableIPTables)
		if err != nil ***REMOVED***
			logrus.Warn(err)
			return err
		***REMOVED***
	***REMOVED***

	d.Lock()
	d.natChain = natChain
	d.filterChain = filterChain
	d.isolationChain = isolationChain
	d.config = config
	d.Unlock()

	err = d.initStore(option)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) getNetwork(id string) (*bridgeNetwork, error) ***REMOVED***
	d.Lock()
	defer d.Unlock()

	if id == "" ***REMOVED***
		return nil, types.BadRequestErrorf("invalid network id: %s", id)
	***REMOVED***

	if nw, ok := d.networks[id]; ok ***REMOVED***
		return nw, nil
	***REMOVED***

	return nil, types.NotFoundErrorf("network not found: %s", id)
***REMOVED***

func parseNetworkGenericOptions(data interface***REMOVED******REMOVED***) (*networkConfiguration, error) ***REMOVED***
	var (
		err    error
		config *networkConfiguration
	)

	switch opt := data.(type) ***REMOVED***
	case *networkConfiguration:
		config = opt
	case map[string]string:
		config = &networkConfiguration***REMOVED***
			EnableICC:          true,
			EnableIPMasquerade: true,
		***REMOVED***
		err = config.fromLabels(opt)
	case options.Generic:
		var opaqueConfig interface***REMOVED******REMOVED***
		if opaqueConfig, err = options.GenerateFromModel(opt, config); err == nil ***REMOVED***
			config = opaqueConfig.(*networkConfiguration)
		***REMOVED***
	default:
		err = types.BadRequestErrorf("do not recognize network configuration format: %T", opt)
	***REMOVED***

	return config, err
***REMOVED***

func (c *networkConfiguration) processIPAM(id string, ipamV4Data, ipamV6Data []driverapi.IPAMData) error ***REMOVED***
	if len(ipamV4Data) > 1 || len(ipamV6Data) > 1 ***REMOVED***
		return types.ForbiddenErrorf("bridge driver doesn't support multiple subnets")
	***REMOVED***

	if len(ipamV4Data) == 0 ***REMOVED***
		return types.BadRequestErrorf("bridge network %s requires ipv4 configuration", id)
	***REMOVED***

	if ipamV4Data[0].Gateway != nil ***REMOVED***
		c.AddressIPv4 = types.GetIPNetCopy(ipamV4Data[0].Gateway)
	***REMOVED***

	if gw, ok := ipamV4Data[0].AuxAddresses[DefaultGatewayV4AuxKey]; ok ***REMOVED***
		c.DefaultGatewayIPv4 = gw.IP
	***REMOVED***

	if len(ipamV6Data) > 0 ***REMOVED***
		c.AddressIPv6 = ipamV6Data[0].Pool

		if ipamV6Data[0].Gateway != nil ***REMOVED***
			c.AddressIPv6 = types.GetIPNetCopy(ipamV6Data[0].Gateway)
		***REMOVED***

		if gw, ok := ipamV6Data[0].AuxAddresses[DefaultGatewayV6AuxKey]; ok ***REMOVED***
			c.DefaultGatewayIPv6 = gw.IP
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func parseNetworkOptions(id string, option options.Generic) (*networkConfiguration, error) ***REMOVED***
	var (
		err    error
		config = &networkConfiguration***REMOVED******REMOVED***
	)

	// Parse generic label first, config will be re-assigned
	if genData, ok := option[netlabel.GenericData]; ok && genData != nil ***REMOVED***
		if config, err = parseNetworkGenericOptions(genData); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	// Process well-known labels next
	if val, ok := option[netlabel.EnableIPv6]; ok ***REMOVED***
		config.EnableIPv6 = val.(bool)
	***REMOVED***

	if val, ok := option[netlabel.Internal]; ok ***REMOVED***
		if internal, ok := val.(bool); ok && internal ***REMOVED***
			config.Internal = true
		***REMOVED***
	***REMOVED***

	// Finally validate the configuration
	if err = config.Validate(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if config.BridgeName == "" && config.DefaultBridge == false ***REMOVED***
		config.BridgeName = "br-" + id[:12]
	***REMOVED***

	exists, err := bridgeInterfaceExists(config.BridgeName)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if !exists ***REMOVED***
		config.BridgeIfaceCreator = ifaceCreatedByLibnetwork
	***REMOVED*** else ***REMOVED***
		config.BridgeIfaceCreator = ifaceCreatedByUser
	***REMOVED***

	config.ID = id
	return config, nil
***REMOVED***

// Returns the non link-local IPv6 subnet for the containers attached to this bridge if found, nil otherwise
func getV6Network(config *networkConfiguration, i *bridgeInterface) *net.IPNet ***REMOVED***
	if config.AddressIPv6 != nil ***REMOVED***
		return config.AddressIPv6
	***REMOVED***
	if i.bridgeIPv6 != nil && i.bridgeIPv6.IP != nil && !i.bridgeIPv6.IP.IsLinkLocalUnicast() ***REMOVED***
		return i.bridgeIPv6
	***REMOVED***

	return nil
***REMOVED***

// Return a slice of networks over which caller can iterate safely
func (d *driver) getNetworks() []*bridgeNetwork ***REMOVED***
	d.Lock()
	defer d.Unlock()

	ls := make([]*bridgeNetwork, 0, len(d.networks))
	for _, nw := range d.networks ***REMOVED***
		ls = append(ls, nw)
	***REMOVED***
	return ls
***REMOVED***

func (d *driver) NetworkAllocate(id string, option map[string]string, ipV4Data, ipV6Data []driverapi.IPAMData) (map[string]string, error) ***REMOVED***
	return nil, types.NotImplementedErrorf("not implemented")
***REMOVED***

func (d *driver) NetworkFree(id string) error ***REMOVED***
	return types.NotImplementedErrorf("not implemented")
***REMOVED***

func (d *driver) EventNotify(etype driverapi.EventType, nid, tableName, key string, value []byte) ***REMOVED***
***REMOVED***

func (d *driver) DecodeTableEntry(tablename string, key string, value []byte) (string, map[string]string) ***REMOVED***
	return "", nil
***REMOVED***

// Create a new network using bridge plugin
func (d *driver) CreateNetwork(id string, option map[string]interface***REMOVED******REMOVED***, nInfo driverapi.NetworkInfo, ipV4Data, ipV6Data []driverapi.IPAMData) error ***REMOVED***
	if len(ipV4Data) == 0 || ipV4Data[0].Pool.String() == "0.0.0.0/0" ***REMOVED***
		return types.BadRequestErrorf("ipv4 pool is empty")
	***REMOVED***
	// Sanity checks
	d.Lock()
	if _, ok := d.networks[id]; ok ***REMOVED***
		d.Unlock()
		return types.ForbiddenErrorf("network %s exists", id)
	***REMOVED***
	d.Unlock()

	// Parse and validate the config. It should not be conflict with existing networks' config
	config, err := parseNetworkOptions(id, option)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err = config.processIPAM(id, ipV4Data, ipV6Data); err != nil ***REMOVED***
		return err
	***REMOVED***

	// start the critical section, from this point onward we are dealing with the list of networks
	// so to be consistent we cannot allow that the list changes
	d.configNetwork.Lock()
	defer d.configNetwork.Unlock()

	// check network conflicts
	if err = d.checkConflict(config); err != nil ***REMOVED***
		nerr, ok := err.(defaultBridgeNetworkConflict)
		if !ok ***REMOVED***
			return err
		***REMOVED***
		// Got a conflict with a stale default network, clean that up and continue
		logrus.Warn(nerr)
		d.deleteNetwork(nerr.ID)
	***REMOVED***

	// there is no conflict, now create the network
	if err = d.createNetwork(config); err != nil ***REMOVED***
		return err
	***REMOVED***

	return d.storeUpdate(config)
***REMOVED***

func (d *driver) checkConflict(config *networkConfiguration) error ***REMOVED***
	networkList := d.getNetworks()
	for _, nw := range networkList ***REMOVED***
		nw.Lock()
		nwConfig := nw.config
		nw.Unlock()
		if err := nwConfig.Conflicts(config); err != nil ***REMOVED***
			if config.DefaultBridge ***REMOVED***
				// We encountered and identified a stale default network
				// We must delete it as libnetwork is the source of truth
				// The default network being created must be the only one
				// This can happen only from docker 1.12 on ward
				logrus.Infof("Found stale default bridge network %s (%s)", nwConfig.ID, nwConfig.BridgeName)
				return defaultBridgeNetworkConflict***REMOVED***nwConfig.ID***REMOVED***
			***REMOVED***

			return types.ForbiddenErrorf("cannot create network %s (%s): conflicts with network %s (%s): %s",
				config.ID, config.BridgeName, nwConfig.ID, nwConfig.BridgeName, err.Error())
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (d *driver) createNetwork(config *networkConfiguration) error ***REMOVED***
	var err error

	defer osl.InitOSContext()()

	networkList := d.getNetworks()

	// Initialize handle when needed
	d.Lock()
	if d.nlh == nil ***REMOVED***
		d.nlh = ns.NlHandle()
	***REMOVED***
	d.Unlock()

	// Create or retrieve the bridge L3 interface
	bridgeIface, err := newInterface(d.nlh, config)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Create and set network handler in driver
	network := &bridgeNetwork***REMOVED***
		id:         config.ID,
		endpoints:  make(map[string]*bridgeEndpoint),
		config:     config,
		portMapper: portmapper.New(d.config.UserlandProxyPath),
		bridge:     bridgeIface,
		driver:     d,
	***REMOVED***

	d.Lock()
	d.networks[config.ID] = network
	d.Unlock()

	// On failure make sure to reset driver network handler to nil
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			d.Lock()
			delete(d.networks, config.ID)
			d.Unlock()
		***REMOVED***
	***REMOVED***()

	// Add inter-network communication rules.
	setupNetworkIsolationRules := func(config *networkConfiguration, i *bridgeInterface) error ***REMOVED***
		if err := network.isolateNetwork(networkList, true); err != nil ***REMOVED***
			if err = network.isolateNetwork(networkList, false); err != nil ***REMOVED***
				logrus.Warnf("Failed on removing the inter-network iptables rules on cleanup: %v", err)
			***REMOVED***
			return err
		***REMOVED***
		// register the cleanup function
		network.registerIptCleanFunc(func() error ***REMOVED***
			nwList := d.getNetworks()
			return network.isolateNetwork(nwList, false)
		***REMOVED***)
		return nil
	***REMOVED***

	// Prepare the bridge setup configuration
	bridgeSetup := newBridgeSetup(config, bridgeIface)

	// If the bridge interface doesn't exist, we need to start the setup steps
	// by creating a new device and assigning it an IPv4 address.
	bridgeAlreadyExists := bridgeIface.exists()
	if !bridgeAlreadyExists ***REMOVED***
		bridgeSetup.queueStep(setupDevice)
	***REMOVED***

	// Even if a bridge exists try to setup IPv4.
	bridgeSetup.queueStep(setupBridgeIPv4)

	enableIPv6Forwarding := d.config.EnableIPForwarding && config.AddressIPv6 != nil

	// Conditionally queue setup steps depending on configuration values.
	for _, step := range []struct ***REMOVED***
		Condition bool
		Fn        setupStep
	***REMOVED******REMOVED***
		// Enable IPv6 on the bridge if required. We do this even for a
		// previously  existing bridge, as it may be here from a previous
		// installation where IPv6 wasn't supported yet and needs to be
		// assigned an IPv6 link-local address.
		***REMOVED***config.EnableIPv6, setupBridgeIPv6***REMOVED***,

		// We ensure that the bridge has the expectedIPv4 and IPv6 addresses in
		// the case of a previously existing device.
		***REMOVED***bridgeAlreadyExists, setupVerifyAndReconcile***REMOVED***,

		// Enable IPv6 Forwarding
		***REMOVED***enableIPv6Forwarding, setupIPv6Forwarding***REMOVED***,

		// Setup Loopback Adresses Routing
		***REMOVED***!d.config.EnableUserlandProxy, setupLoopbackAdressesRouting***REMOVED***,

		// Setup IPTables.
		***REMOVED***d.config.EnableIPTables, network.setupIPTables***REMOVED***,

		//We want to track firewalld configuration so that
		//if it is started/reloaded, the rules can be applied correctly
		***REMOVED***d.config.EnableIPTables, network.setupFirewalld***REMOVED***,

		// Setup DefaultGatewayIPv4
		***REMOVED***config.DefaultGatewayIPv4 != nil, setupGatewayIPv4***REMOVED***,

		// Setup DefaultGatewayIPv6
		***REMOVED***config.DefaultGatewayIPv6 != nil, setupGatewayIPv6***REMOVED***,

		// Add inter-network communication rules.
		***REMOVED***d.config.EnableIPTables, setupNetworkIsolationRules***REMOVED***,

		//Configure bridge networking filtering if ICC is off and IP tables are enabled
		***REMOVED***!config.EnableICC && d.config.EnableIPTables, setupBridgeNetFiltering***REMOVED***,
	***REMOVED*** ***REMOVED***
		if step.Condition ***REMOVED***
			bridgeSetup.queueStep(step.Fn)
		***REMOVED***
	***REMOVED***

	// Apply the prepared list of steps, and abort at the first error.
	bridgeSetup.queueStep(setupDeviceUp)
	return bridgeSetup.apply()
***REMOVED***

func (d *driver) DeleteNetwork(nid string) error ***REMOVED***

	d.configNetwork.Lock()
	defer d.configNetwork.Unlock()

	return d.deleteNetwork(nid)
***REMOVED***

func (d *driver) deleteNetwork(nid string) error ***REMOVED***
	var err error

	defer osl.InitOSContext()()
	// Get network handler and remove it from driver
	d.Lock()
	n, ok := d.networks[nid]
	d.Unlock()

	if !ok ***REMOVED***
		return types.InternalMaskableErrorf("network %s does not exist", nid)
	***REMOVED***

	n.Lock()
	config := n.config
	n.Unlock()

	// delele endpoints belong to this network
	for _, ep := range n.endpoints ***REMOVED***
		if err := n.releasePorts(ep); err != nil ***REMOVED***
			logrus.Warn(err)
		***REMOVED***
		if link, err := d.nlh.LinkByName(ep.srcName); err == nil ***REMOVED***
			d.nlh.LinkDel(link)
		***REMOVED***

		if err := d.storeDelete(ep); err != nil ***REMOVED***
			logrus.Warnf("Failed to remove bridge endpoint %s from store: %v", ep.id[0:7], err)
		***REMOVED***
	***REMOVED***

	d.Lock()
	delete(d.networks, nid)
	d.Unlock()

	// On failure set network handler back in driver, but
	// only if is not already taken over by some other thread
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			d.Lock()
			if _, ok := d.networks[nid]; !ok ***REMOVED***
				d.networks[nid] = n
			***REMOVED***
			d.Unlock()
		***REMOVED***
	***REMOVED***()

	switch config.BridgeIfaceCreator ***REMOVED***
	case ifaceCreatedByLibnetwork, ifaceCreatorUnknown:
		// We only delete the bridge if it was created by the bridge driver and
		// it is not the default one (to keep the backward compatible behavior.)
		if !config.DefaultBridge ***REMOVED***
			if err := d.nlh.LinkDel(n.bridge.Link); err != nil ***REMOVED***
				logrus.Warnf("Failed to remove bridge interface %s on network %s delete: %v", config.BridgeName, nid, err)
			***REMOVED***
		***REMOVED***
	case ifaceCreatedByUser:
		// Don't delete the bridge interface if it was not created by libnetwork.
	***REMOVED***

	// clean all relevant iptables rules
	for _, cleanFunc := range n.iptCleanFuncs ***REMOVED***
		if errClean := cleanFunc(); errClean != nil ***REMOVED***
			logrus.Warnf("Failed to clean iptables rules for bridge network: %v", errClean)
		***REMOVED***
	***REMOVED***
	return d.storeDelete(config)
***REMOVED***

func addToBridge(nlh *netlink.Handle, ifaceName, bridgeName string) error ***REMOVED***
	link, err := nlh.LinkByName(ifaceName)
	if err != nil ***REMOVED***
		return fmt.Errorf("could not find interface %s: %v", ifaceName, err)
	***REMOVED***
	if err = nlh.LinkSetMaster(link,
		&netlink.Bridge***REMOVED***LinkAttrs: netlink.LinkAttrs***REMOVED***Name: bridgeName***REMOVED******REMOVED***); err != nil ***REMOVED***
		logrus.Debugf("Failed to add %s to bridge via netlink.Trying ioctl: %v", ifaceName, err)
		iface, err := net.InterfaceByName(ifaceName)
		if err != nil ***REMOVED***
			return fmt.Errorf("could not find network interface %s: %v", ifaceName, err)
		***REMOVED***

		master, err := net.InterfaceByName(bridgeName)
		if err != nil ***REMOVED***
			return fmt.Errorf("could not find bridge %s: %v", bridgeName, err)
		***REMOVED***

		return ioctlAddToBridge(iface, master)
	***REMOVED***
	return nil
***REMOVED***

func setHairpinMode(nlh *netlink.Handle, link netlink.Link, enable bool) error ***REMOVED***
	err := nlh.LinkSetHairpin(link, enable)
	if err != nil && err != syscall.EINVAL ***REMOVED***
		// If error is not EINVAL something else went wrong, bail out right away
		return fmt.Errorf("unable to set hairpin mode on %s via netlink: %v",
			link.Attrs().Name, err)
	***REMOVED***

	// Hairpin mode successfully set up
	if err == nil ***REMOVED***
		return nil
	***REMOVED***

	// The netlink method failed with EINVAL which is probably because of an older
	// kernel. Try one more time via the sysfs method.
	path := filepath.Join("/sys/class/net", link.Attrs().Name, "brport/hairpin_mode")

	var val []byte
	if enable ***REMOVED***
		val = []byte***REMOVED***'1', '\n'***REMOVED***
	***REMOVED*** else ***REMOVED***
		val = []byte***REMOVED***'0', '\n'***REMOVED***
	***REMOVED***

	if err := ioutil.WriteFile(path, val, 0644); err != nil ***REMOVED***
		return fmt.Errorf("unable to set hairpin mode on %s via sysfs: %v", link.Attrs().Name, err)
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) CreateEndpoint(nid, eid string, ifInfo driverapi.InterfaceInfo, epOptions map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	defer osl.InitOSContext()()

	if ifInfo == nil ***REMOVED***
		return errors.New("invalid interface info passed")
	***REMOVED***

	// Get the network handler and make sure it exists
	d.Lock()
	n, ok := d.networks[nid]
	dconfig := d.config
	d.Unlock()

	if !ok ***REMOVED***
		return types.NotFoundErrorf("network %s does not exist", nid)
	***REMOVED***
	if n == nil ***REMOVED***
		return driverapi.ErrNoNetwork(nid)
	***REMOVED***

	// Sanity check
	n.Lock()
	if n.id != nid ***REMOVED***
		n.Unlock()
		return InvalidNetworkIDError(nid)
	***REMOVED***
	n.Unlock()

	// Check if endpoint id is good and retrieve correspondent endpoint
	ep, err := n.getEndpoint(eid)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Endpoint with that id exists either on desired or other sandbox
	if ep != nil ***REMOVED***
		return driverapi.ErrEndpointExists(eid)
	***REMOVED***

	// Try to convert the options to endpoint configuration
	epConfig, err := parseEndpointOptions(epOptions)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Create and add the endpoint
	n.Lock()
	endpoint := &bridgeEndpoint***REMOVED***id: eid, nid: nid, config: epConfig***REMOVED***
	n.endpoints[eid] = endpoint
	n.Unlock()

	// On failure make sure to remove the endpoint
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			n.Lock()
			delete(n.endpoints, eid)
			n.Unlock()
		***REMOVED***
	***REMOVED***()

	// Generate a name for what will be the host side pipe interface
	hostIfName, err := netutils.GenerateIfaceName(d.nlh, vethPrefix, vethLen)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Generate a name for what will be the sandbox side pipe interface
	containerIfName, err := netutils.GenerateIfaceName(d.nlh, vethPrefix, vethLen)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Generate and add the interface pipe host <-> sandbox
	veth := &netlink.Veth***REMOVED***
		LinkAttrs: netlink.LinkAttrs***REMOVED***Name: hostIfName, TxQLen: 0***REMOVED***,
		PeerName:  containerIfName***REMOVED***
	if err = d.nlh.LinkAdd(veth); err != nil ***REMOVED***
		return types.InternalErrorf("failed to add the host (%s) <=> sandbox (%s) pair interfaces: %v", hostIfName, containerIfName, err)
	***REMOVED***

	// Get the host side pipe interface handler
	host, err := d.nlh.LinkByName(hostIfName)
	if err != nil ***REMOVED***
		return types.InternalErrorf("failed to find host side interface %s: %v", hostIfName, err)
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			d.nlh.LinkDel(host)
		***REMOVED***
	***REMOVED***()

	// Get the sandbox side pipe interface handler
	sbox, err := d.nlh.LinkByName(containerIfName)
	if err != nil ***REMOVED***
		return types.InternalErrorf("failed to find sandbox side interface %s: %v", containerIfName, err)
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			d.nlh.LinkDel(sbox)
		***REMOVED***
	***REMOVED***()

	n.Lock()
	config := n.config
	n.Unlock()

	// Add bridge inherited attributes to pipe interfaces
	if config.Mtu != 0 ***REMOVED***
		err = d.nlh.LinkSetMTU(host, config.Mtu)
		if err != nil ***REMOVED***
			return types.InternalErrorf("failed to set MTU on host interface %s: %v", hostIfName, err)
		***REMOVED***
		err = d.nlh.LinkSetMTU(sbox, config.Mtu)
		if err != nil ***REMOVED***
			return types.InternalErrorf("failed to set MTU on sandbox interface %s: %v", containerIfName, err)
		***REMOVED***
	***REMOVED***

	// Attach host side pipe interface into the bridge
	if err = addToBridge(d.nlh, hostIfName, config.BridgeName); err != nil ***REMOVED***
		return fmt.Errorf("adding interface %s to bridge %s failed: %v", hostIfName, config.BridgeName, err)
	***REMOVED***

	if !dconfig.EnableUserlandProxy ***REMOVED***
		err = setHairpinMode(d.nlh, host, true)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Store the sandbox side pipe interface parameters
	endpoint.srcName = containerIfName
	endpoint.macAddress = ifInfo.MacAddress()
	endpoint.addr = ifInfo.Address()
	endpoint.addrv6 = ifInfo.AddressIPv6()

	// Set the sbox's MAC if not provided. If specified, use the one configured by user, otherwise generate one based on IP.
	if endpoint.macAddress == nil ***REMOVED***
		endpoint.macAddress = electMacAddress(epConfig, endpoint.addr.IP)
		if err = ifInfo.SetMacAddress(endpoint.macAddress); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Up the host interface after finishing all netlink configuration
	if err = d.nlh.LinkSetUp(host); err != nil ***REMOVED***
		return fmt.Errorf("could not set link up for host interface %s: %v", hostIfName, err)
	***REMOVED***

	if endpoint.addrv6 == nil && config.EnableIPv6 ***REMOVED***
		var ip6 net.IP
		network := n.bridge.bridgeIPv6
		if config.AddressIPv6 != nil ***REMOVED***
			network = config.AddressIPv6
		***REMOVED***

		ones, _ := network.Mask.Size()
		if ones > 80 ***REMOVED***
			err = types.ForbiddenErrorf("Cannot self generate an IPv6 address on network %v: At least 48 host bits are needed.", network)
			return err
		***REMOVED***

		ip6 = make(net.IP, len(network.IP))
		copy(ip6, network.IP)
		for i, h := range endpoint.macAddress ***REMOVED***
			ip6[i+10] = h
		***REMOVED***

		endpoint.addrv6 = &net.IPNet***REMOVED***IP: ip6, Mask: network.Mask***REMOVED***
		if err = ifInfo.SetIPAddress(endpoint.addrv6); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if err = d.storeUpdate(endpoint); err != nil ***REMOVED***
		return fmt.Errorf("failed to save bridge endpoint %s to store: %v", endpoint.id[0:7], err)
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) DeleteEndpoint(nid, eid string) error ***REMOVED***
	var err error

	defer osl.InitOSContext()()

	// Get the network handler and make sure it exists
	d.Lock()
	n, ok := d.networks[nid]
	d.Unlock()

	if !ok ***REMOVED***
		return types.InternalMaskableErrorf("network %s does not exist", nid)
	***REMOVED***
	if n == nil ***REMOVED***
		return driverapi.ErrNoNetwork(nid)
	***REMOVED***

	// Sanity Check
	n.Lock()
	if n.id != nid ***REMOVED***
		n.Unlock()
		return InvalidNetworkIDError(nid)
	***REMOVED***
	n.Unlock()

	// Check endpoint id and if an endpoint is actually there
	ep, err := n.getEndpoint(eid)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if ep == nil ***REMOVED***
		return EndpointNotFoundError(eid)
	***REMOVED***

	// Remove it
	n.Lock()
	delete(n.endpoints, eid)
	n.Unlock()

	// On failure make sure to set back ep in n.endpoints, but only
	// if it hasn't been taken over already by some other thread.
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			n.Lock()
			if _, ok := n.endpoints[eid]; !ok ***REMOVED***
				n.endpoints[eid] = ep
			***REMOVED***
			n.Unlock()
		***REMOVED***
	***REMOVED***()

	// Try removal of link. Discard error: it is a best effort.
	// Also make sure defer does not see this error either.
	if link, err := d.nlh.LinkByName(ep.srcName); err == nil ***REMOVED***
		d.nlh.LinkDel(link)
	***REMOVED***

	if err := d.storeDelete(ep); err != nil ***REMOVED***
		logrus.Warnf("Failed to remove bridge endpoint %s from store: %v", ep.id[0:7], err)
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) EndpointOperInfo(nid, eid string) (map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	// Get the network handler and make sure it exists
	d.Lock()
	n, ok := d.networks[nid]
	d.Unlock()
	if !ok ***REMOVED***
		return nil, types.NotFoundErrorf("network %s does not exist", nid)
	***REMOVED***
	if n == nil ***REMOVED***
		return nil, driverapi.ErrNoNetwork(nid)
	***REMOVED***

	// Sanity check
	n.Lock()
	if n.id != nid ***REMOVED***
		n.Unlock()
		return nil, InvalidNetworkIDError(nid)
	***REMOVED***
	n.Unlock()

	// Check if endpoint id is good and retrieve correspondent endpoint
	ep, err := n.getEndpoint(eid)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if ep == nil ***REMOVED***
		return nil, driverapi.ErrNoEndpoint(eid)
	***REMOVED***

	m := make(map[string]interface***REMOVED******REMOVED***)

	if ep.extConnConfig != nil && ep.extConnConfig.ExposedPorts != nil ***REMOVED***
		// Return a copy of the config data
		epc := make([]types.TransportPort, 0, len(ep.extConnConfig.ExposedPorts))
		for _, tp := range ep.extConnConfig.ExposedPorts ***REMOVED***
			epc = append(epc, tp.GetCopy())
		***REMOVED***
		m[netlabel.ExposedPorts] = epc
	***REMOVED***

	if ep.portMapping != nil ***REMOVED***
		// Return a copy of the operational data
		pmc := make([]types.PortBinding, 0, len(ep.portMapping))
		for _, pm := range ep.portMapping ***REMOVED***
			pmc = append(pmc, pm.GetCopy())
		***REMOVED***
		m[netlabel.PortMap] = pmc
	***REMOVED***

	if len(ep.macAddress) != 0 ***REMOVED***
		m[netlabel.MacAddress] = ep.macAddress
	***REMOVED***

	return m, nil
***REMOVED***

// Join method is invoked when a Sandbox is attached to an endpoint.
func (d *driver) Join(nid, eid string, sboxKey string, jinfo driverapi.JoinInfo, options map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	defer osl.InitOSContext()()

	network, err := d.getNetwork(nid)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	endpoint, err := network.getEndpoint(eid)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if endpoint == nil ***REMOVED***
		return EndpointNotFoundError(eid)
	***REMOVED***

	endpoint.containerConfig, err = parseContainerOptions(options)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	iNames := jinfo.InterfaceName()
	containerVethPrefix := defaultContainerVethPrefix
	if network.config.ContainerIfacePrefix != "" ***REMOVED***
		containerVethPrefix = network.config.ContainerIfacePrefix
	***REMOVED***
	err = iNames.SetNames(endpoint.srcName, containerVethPrefix)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = jinfo.SetGateway(network.bridge.gatewayIPv4)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = jinfo.SetGatewayIPv6(network.bridge.gatewayIPv6)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

// Leave method is invoked when a Sandbox detaches from an endpoint.
func (d *driver) Leave(nid, eid string) error ***REMOVED***
	defer osl.InitOSContext()()

	network, err := d.getNetwork(nid)
	if err != nil ***REMOVED***
		return types.InternalMaskableErrorf("%s", err)
	***REMOVED***

	endpoint, err := network.getEndpoint(eid)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if endpoint == nil ***REMOVED***
		return EndpointNotFoundError(eid)
	***REMOVED***

	if !network.config.EnableICC ***REMOVED***
		if err = d.link(network, endpoint, false); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) ProgramExternalConnectivity(nid, eid string, options map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	defer osl.InitOSContext()()

	network, err := d.getNetwork(nid)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	endpoint, err := network.getEndpoint(eid)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if endpoint == nil ***REMOVED***
		return EndpointNotFoundError(eid)
	***REMOVED***

	endpoint.extConnConfig, err = parseConnectivityOptions(options)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Program any required port mapping and store them in the endpoint
	endpoint.portMapping, err = network.allocatePorts(endpoint, network.config.DefaultBindingIP, d.config.EnableUserlandProxy)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			if e := network.releasePorts(endpoint); e != nil ***REMOVED***
				logrus.Errorf("Failed to release ports allocated for the bridge endpoint %s on failure %v because of %v",
					eid, err, e)
			***REMOVED***
			endpoint.portMapping = nil
		***REMOVED***
	***REMOVED***()

	if err = d.storeUpdate(endpoint); err != nil ***REMOVED***
		return fmt.Errorf("failed to update bridge endpoint %s to store: %v", endpoint.id[0:7], err)
	***REMOVED***

	if !network.config.EnableICC ***REMOVED***
		return d.link(network, endpoint, true)
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) RevokeExternalConnectivity(nid, eid string) error ***REMOVED***
	defer osl.InitOSContext()()

	network, err := d.getNetwork(nid)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	endpoint, err := network.getEndpoint(eid)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if endpoint == nil ***REMOVED***
		return EndpointNotFoundError(eid)
	***REMOVED***

	err = network.releasePorts(endpoint)
	if err != nil ***REMOVED***
		logrus.Warn(err)
	***REMOVED***

	endpoint.portMapping = nil

	// Clean the connection tracker state of the host for the specific endpoint
	// The host kernel keeps track of the connections (TCP and UDP), so if a new endpoint gets the same IP of
	// this one (that is going down), is possible that some of the packets would not be routed correctly inside
	// the new endpoint
	// Deeper details: https://github.com/docker/docker/issues/8795
	clearEndpointConnections(d.nlh, endpoint)

	if err = d.storeUpdate(endpoint); err != nil ***REMOVED***
		return fmt.Errorf("failed to update bridge endpoint %s to store: %v", endpoint.id[0:7], err)
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) link(network *bridgeNetwork, endpoint *bridgeEndpoint, enable bool) error ***REMOVED***
	var err error

	cc := endpoint.containerConfig
	if cc == nil ***REMOVED***
		return nil
	***REMOVED***
	ec := endpoint.extConnConfig
	if ec == nil ***REMOVED***
		return nil
	***REMOVED***

	if ec.ExposedPorts != nil ***REMOVED***
		for _, p := range cc.ParentEndpoints ***REMOVED***
			var parentEndpoint *bridgeEndpoint
			parentEndpoint, err = network.getEndpoint(p)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if parentEndpoint == nil ***REMOVED***
				err = InvalidEndpointIDError(p)
				return err
			***REMOVED***

			l := newLink(parentEndpoint.addr.IP.String(),
				endpoint.addr.IP.String(),
				ec.ExposedPorts, network.config.BridgeName)
			if enable ***REMOVED***
				err = l.Enable()
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				defer func() ***REMOVED***
					if err != nil ***REMOVED***
						l.Disable()
					***REMOVED***
				***REMOVED***()
			***REMOVED*** else ***REMOVED***
				l.Disable()
			***REMOVED***
		***REMOVED***
	***REMOVED***

	for _, c := range cc.ChildEndpoints ***REMOVED***
		var childEndpoint *bridgeEndpoint
		childEndpoint, err = network.getEndpoint(c)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if childEndpoint == nil ***REMOVED***
			err = InvalidEndpointIDError(c)
			return err
		***REMOVED***
		if childEndpoint.extConnConfig == nil || childEndpoint.extConnConfig.ExposedPorts == nil ***REMOVED***
			continue
		***REMOVED***

		l := newLink(endpoint.addr.IP.String(),
			childEndpoint.addr.IP.String(),
			childEndpoint.extConnConfig.ExposedPorts, network.config.BridgeName)
		if enable ***REMOVED***
			err = l.Enable()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			defer func() ***REMOVED***
				if err != nil ***REMOVED***
					l.Disable()
				***REMOVED***
			***REMOVED***()
		***REMOVED*** else ***REMOVED***
			l.Disable()
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) Type() string ***REMOVED***
	return networkType
***REMOVED***

func (d *driver) IsBuiltIn() bool ***REMOVED***
	return true
***REMOVED***

// DiscoverNew is a notification for a new discovery event, such as a new node joining a cluster
func (d *driver) DiscoverNew(dType discoverapi.DiscoveryType, data interface***REMOVED******REMOVED***) error ***REMOVED***
	return nil
***REMOVED***

// DiscoverDelete is a notification for a discovery delete event, such as a node leaving a cluster
func (d *driver) DiscoverDelete(dType discoverapi.DiscoveryType, data interface***REMOVED******REMOVED***) error ***REMOVED***
	return nil
***REMOVED***

func parseEndpointOptions(epOptions map[string]interface***REMOVED******REMOVED***) (*endpointConfiguration, error) ***REMOVED***
	if epOptions == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	ec := &endpointConfiguration***REMOVED******REMOVED***

	if opt, ok := epOptions[netlabel.MacAddress]; ok ***REMOVED***
		if mac, ok := opt.(net.HardwareAddr); ok ***REMOVED***
			ec.MacAddress = mac
		***REMOVED*** else ***REMOVED***
			return nil, &ErrInvalidEndpointConfig***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***

	return ec, nil
***REMOVED***

func parseContainerOptions(cOptions map[string]interface***REMOVED******REMOVED***) (*containerConfiguration, error) ***REMOVED***
	if cOptions == nil ***REMOVED***
		return nil, nil
	***REMOVED***
	genericData := cOptions[netlabel.GenericData]
	if genericData == nil ***REMOVED***
		return nil, nil
	***REMOVED***
	switch opt := genericData.(type) ***REMOVED***
	case options.Generic:
		opaqueConfig, err := options.GenerateFromModel(opt, &containerConfiguration***REMOVED******REMOVED***)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return opaqueConfig.(*containerConfiguration), nil
	case *containerConfiguration:
		return opt, nil
	default:
		return nil, nil
	***REMOVED***
***REMOVED***

func parseConnectivityOptions(cOptions map[string]interface***REMOVED******REMOVED***) (*connectivityConfiguration, error) ***REMOVED***
	if cOptions == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	cc := &connectivityConfiguration***REMOVED******REMOVED***

	if opt, ok := cOptions[netlabel.PortMap]; ok ***REMOVED***
		if pb, ok := opt.([]types.PortBinding); ok ***REMOVED***
			cc.PortBindings = pb
		***REMOVED*** else ***REMOVED***
			return nil, types.BadRequestErrorf("Invalid port mapping data in connectivity configuration: %v", opt)
		***REMOVED***
	***REMOVED***

	if opt, ok := cOptions[netlabel.ExposedPorts]; ok ***REMOVED***
		if ports, ok := opt.([]types.TransportPort); ok ***REMOVED***
			cc.ExposedPorts = ports
		***REMOVED*** else ***REMOVED***
			return nil, types.BadRequestErrorf("Invalid exposed ports data in connectivity configuration: %v", opt)
		***REMOVED***
	***REMOVED***

	return cc, nil
***REMOVED***

func electMacAddress(epConfig *endpointConfiguration, ip net.IP) net.HardwareAddr ***REMOVED***
	if epConfig != nil && epConfig.MacAddress != nil ***REMOVED***
		return epConfig.MacAddress
	***REMOVED***
	return netutils.GenerateMACFromIP(ip)
***REMOVED***
