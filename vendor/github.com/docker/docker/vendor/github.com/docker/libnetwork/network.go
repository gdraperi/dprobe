package libnetwork

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/libnetwork/common"
	"github.com/docker/libnetwork/config"
	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/etchosts"
	"github.com/docker/libnetwork/ipamapi"
	"github.com/docker/libnetwork/netlabel"
	"github.com/docker/libnetwork/netutils"
	"github.com/docker/libnetwork/networkdb"
	"github.com/docker/libnetwork/options"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

// A Network represents a logical connectivity zone that containers may
// join using the Link method. A Network is managed by a specific driver.
type Network interface ***REMOVED***
	// A user chosen name for this network.
	Name() string

	// A system generated id for this network.
	ID() string

	// The type of network, which corresponds to its managing driver.
	Type() string

	// Create a new endpoint to this network symbolically identified by the
	// specified unique name. The options parameter carries driver specific options.
	CreateEndpoint(name string, options ...EndpointOption) (Endpoint, error)

	// Delete the network.
	Delete() error

	// Endpoints returns the list of Endpoint(s) in this network.
	Endpoints() []Endpoint

	// WalkEndpoints uses the provided function to walk the Endpoints
	WalkEndpoints(walker EndpointWalker)

	// EndpointByName returns the Endpoint which has the passed name. If not found, the error ErrNoSuchEndpoint is returned.
	EndpointByName(name string) (Endpoint, error)

	// EndpointByID returns the Endpoint which has the passed id. If not found, the error ErrNoSuchEndpoint is returned.
	EndpointByID(id string) (Endpoint, error)

	// Return certain operational data belonging to this network
	Info() NetworkInfo
***REMOVED***

// NetworkInfo returns some configuration and operational information about the network
type NetworkInfo interface ***REMOVED***
	IpamConfig() (string, map[string]string, []*IpamConf, []*IpamConf)
	IpamInfo() ([]*IpamInfo, []*IpamInfo)
	DriverOptions() map[string]string
	Scope() string
	IPv6Enabled() bool
	Internal() bool
	Attachable() bool
	Ingress() bool
	ConfigFrom() string
	ConfigOnly() bool
	Labels() map[string]string
	Dynamic() bool
	Created() time.Time
	// Peers returns a slice of PeerInfo structures which has the information about the peer
	// nodes participating in the same overlay network. This is currently the per-network
	// gossip cluster. For non-dynamic overlay networks and bridge networks it returns an
	// empty slice
	Peers() []networkdb.PeerInfo
	//Services returns a map of services keyed by the service name with the details
	//of all the tasks that belong to the service. Applicable only in swarm mode.
	Services() map[string]ServiceInfo
***REMOVED***

// EndpointWalker is a client provided function which will be used to walk the Endpoints.
// When the function returns true, the walk will stop.
type EndpointWalker func(ep Endpoint) bool

// ipInfo is the reverse mapping from IP to service name to serve the PTR query.
// extResolver is set if an externl server resolves a service name to this IP.
// Its an indication to defer PTR queries also to that external server.
type ipInfo struct ***REMOVED***
	name        string
	serviceID   string
	extResolver bool
***REMOVED***

// svcMapEntry is the body of the element into the svcMap
// The ip is a string because the SetMatrix does not accept non hashable values
type svcMapEntry struct ***REMOVED***
	ip        string
	serviceID string
***REMOVED***

type svcInfo struct ***REMOVED***
	svcMap     common.SetMatrix
	svcIPv6Map common.SetMatrix
	ipMap      common.SetMatrix
	service    map[string][]servicePorts
***REMOVED***

// backing container or host's info
type serviceTarget struct ***REMOVED***
	name string
	ip   net.IP
	port uint16
***REMOVED***

type servicePorts struct ***REMOVED***
	portName string
	proto    string
	target   []serviceTarget
***REMOVED***

type networkDBTable struct ***REMOVED***
	name    string
	objType driverapi.ObjectType
***REMOVED***

// IpamConf contains all the ipam related configurations for a network
type IpamConf struct ***REMOVED***
	// The master address pool for containers and network interfaces
	PreferredPool string
	// A subset of the master pool. If specified,
	// this becomes the container pool
	SubPool string
	// Preferred Network Gateway address (optional)
	Gateway string
	// Auxiliary addresses for network driver. Must be within the master pool.
	// libnetwork will reserve them if they fall into the container pool
	AuxAddresses map[string]string
***REMOVED***

// Validate checks whether the configuration is valid
func (c *IpamConf) Validate() error ***REMOVED***
	if c.Gateway != "" && nil == net.ParseIP(c.Gateway) ***REMOVED***
		return types.BadRequestErrorf("invalid gateway address %s in Ipam configuration", c.Gateway)
	***REMOVED***
	return nil
***REMOVED***

// IpamInfo contains all the ipam related operational info for a network
type IpamInfo struct ***REMOVED***
	PoolID string
	Meta   map[string]string
	driverapi.IPAMData
***REMOVED***

// MarshalJSON encodes IpamInfo into json message
func (i *IpamInfo) MarshalJSON() ([]byte, error) ***REMOVED***
	m := map[string]interface***REMOVED******REMOVED******REMOVED***
		"PoolID": i.PoolID,
	***REMOVED***
	v, err := json.Marshal(&i.IPAMData)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	m["IPAMData"] = string(v)

	if i.Meta != nil ***REMOVED***
		m["Meta"] = i.Meta
	***REMOVED***
	return json.Marshal(m)
***REMOVED***

// UnmarshalJSON decodes json message into PoolData
func (i *IpamInfo) UnmarshalJSON(data []byte) error ***REMOVED***
	var (
		m   map[string]interface***REMOVED******REMOVED***
		err error
	)
	if err = json.Unmarshal(data, &m); err != nil ***REMOVED***
		return err
	***REMOVED***
	i.PoolID = m["PoolID"].(string)
	if v, ok := m["Meta"]; ok ***REMOVED***
		b, _ := json.Marshal(v)
		if err = json.Unmarshal(b, &i.Meta); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if v, ok := m["IPAMData"]; ok ***REMOVED***
		if err = json.Unmarshal([]byte(v.(string)), &i.IPAMData); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

type network struct ***REMOVED***
	ctrlr          *controller
	name           string
	networkType    string
	id             string
	created        time.Time
	scope          string // network data scope
	labels         map[string]string
	ipamType       string
	ipamOptions    map[string]string
	addrSpace      string
	ipamV4Config   []*IpamConf
	ipamV6Config   []*IpamConf
	ipamV4Info     []*IpamInfo
	ipamV6Info     []*IpamInfo
	enableIPv6     bool
	postIPv6       bool
	epCnt          *endpointCnt
	generic        options.Generic
	dbIndex        uint64
	dbExists       bool
	persist        bool
	stopWatchCh    chan struct***REMOVED******REMOVED***
	drvOnce        *sync.Once
	resolverOnce   sync.Once
	resolver       []Resolver
	internal       bool
	attachable     bool
	inDelete       bool
	ingress        bool
	driverTables   []networkDBTable
	dynamic        bool
	configOnly     bool
	configFrom     string
	loadBalancerIP net.IP
	sync.Mutex
***REMOVED***

func (n *network) Name() string ***REMOVED***
	n.Lock()
	defer n.Unlock()

	return n.name
***REMOVED***

func (n *network) ID() string ***REMOVED***
	n.Lock()
	defer n.Unlock()

	return n.id
***REMOVED***

func (n *network) Created() time.Time ***REMOVED***
	n.Lock()
	defer n.Unlock()

	return n.created
***REMOVED***

func (n *network) Type() string ***REMOVED***
	n.Lock()
	defer n.Unlock()

	return n.networkType
***REMOVED***

func (n *network) Key() []string ***REMOVED***
	n.Lock()
	defer n.Unlock()
	return []string***REMOVED***datastore.NetworkKeyPrefix, n.id***REMOVED***
***REMOVED***

func (n *network) KeyPrefix() []string ***REMOVED***
	return []string***REMOVED***datastore.NetworkKeyPrefix***REMOVED***
***REMOVED***

func (n *network) Value() []byte ***REMOVED***
	n.Lock()
	defer n.Unlock()
	b, err := json.Marshal(n)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	return b
***REMOVED***

func (n *network) SetValue(value []byte) error ***REMOVED***
	return json.Unmarshal(value, n)
***REMOVED***

func (n *network) Index() uint64 ***REMOVED***
	n.Lock()
	defer n.Unlock()
	return n.dbIndex
***REMOVED***

func (n *network) SetIndex(index uint64) ***REMOVED***
	n.Lock()
	n.dbIndex = index
	n.dbExists = true
	n.Unlock()
***REMOVED***

func (n *network) Exists() bool ***REMOVED***
	n.Lock()
	defer n.Unlock()
	return n.dbExists
***REMOVED***

func (n *network) Skip() bool ***REMOVED***
	n.Lock()
	defer n.Unlock()
	return !n.persist
***REMOVED***

func (n *network) New() datastore.KVObject ***REMOVED***
	n.Lock()
	defer n.Unlock()

	return &network***REMOVED***
		ctrlr:   n.ctrlr,
		drvOnce: &sync.Once***REMOVED******REMOVED***,
		scope:   n.scope,
	***REMOVED***
***REMOVED***

// CopyTo deep copies to the destination IpamConfig
func (c *IpamConf) CopyTo(dstC *IpamConf) error ***REMOVED***
	dstC.PreferredPool = c.PreferredPool
	dstC.SubPool = c.SubPool
	dstC.Gateway = c.Gateway
	if c.AuxAddresses != nil ***REMOVED***
		dstC.AuxAddresses = make(map[string]string, len(c.AuxAddresses))
		for k, v := range c.AuxAddresses ***REMOVED***
			dstC.AuxAddresses[k] = v
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// CopyTo deep copies to the destination IpamInfo
func (i *IpamInfo) CopyTo(dstI *IpamInfo) error ***REMOVED***
	dstI.PoolID = i.PoolID
	if i.Meta != nil ***REMOVED***
		dstI.Meta = make(map[string]string)
		for k, v := range i.Meta ***REMOVED***
			dstI.Meta[k] = v
		***REMOVED***
	***REMOVED***

	dstI.AddressSpace = i.AddressSpace
	dstI.Pool = types.GetIPNetCopy(i.Pool)
	dstI.Gateway = types.GetIPNetCopy(i.Gateway)

	if i.AuxAddresses != nil ***REMOVED***
		dstI.AuxAddresses = make(map[string]*net.IPNet)
		for k, v := range i.AuxAddresses ***REMOVED***
			dstI.AuxAddresses[k] = types.GetIPNetCopy(v)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (n *network) validateConfiguration() error ***REMOVED***
	if n.configOnly ***REMOVED***
		// Only supports network specific configurations.
		// Network operator configurations are not supported.
		if n.ingress || n.internal || n.attachable || n.scope != "" ***REMOVED***
			return types.ForbiddenErrorf("configuration network can only contain network " +
				"specific fields. Network operator fields like " +
				"[ ingress | internal | attachable | scope ] are not supported.")
		***REMOVED***
	***REMOVED***
	if n.configFrom != "" ***REMOVED***
		if n.configOnly ***REMOVED***
			return types.ForbiddenErrorf("a configuration network cannot depend on another configuration network")
		***REMOVED***
		if n.ipamType != "" &&
			n.ipamType != defaultIpamForNetworkType(n.networkType) ||
			n.enableIPv6 ||
			len(n.labels) > 0 || len(n.ipamOptions) > 0 ||
			len(n.ipamV4Config) > 0 || len(n.ipamV6Config) > 0 ***REMOVED***
			return types.ForbiddenErrorf("user specified configurations are not supported if the network depends on a configuration network")
		***REMOVED***
		if len(n.generic) > 0 ***REMOVED***
			if data, ok := n.generic[netlabel.GenericData]; ok ***REMOVED***
				var (
					driverOptions map[string]string
					opts          interface***REMOVED******REMOVED***
				)
				switch data.(type) ***REMOVED***
				case map[string]interface***REMOVED******REMOVED***:
					opts = data.(map[string]interface***REMOVED******REMOVED***)
				case map[string]string:
					opts = data.(map[string]string)
				***REMOVED***
				ba, err := json.Marshal(opts)
				if err != nil ***REMOVED***
					return fmt.Errorf("failed to validate network configuration: %v", err)
				***REMOVED***
				if err := json.Unmarshal(ba, &driverOptions); err != nil ***REMOVED***
					return fmt.Errorf("failed to validate network configuration: %v", err)
				***REMOVED***
				if len(driverOptions) > 0 ***REMOVED***
					return types.ForbiddenErrorf("network driver options are not supported if the network depends on a configuration network")
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Applies network specific configurations
func (n *network) applyConfigurationTo(to *network) error ***REMOVED***
	to.enableIPv6 = n.enableIPv6
	if len(n.labels) > 0 ***REMOVED***
		to.labels = make(map[string]string, len(n.labels))
		for k, v := range n.labels ***REMOVED***
			if _, ok := to.labels[k]; !ok ***REMOVED***
				to.labels[k] = v
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if len(n.ipamType) != 0 ***REMOVED***
		to.ipamType = n.ipamType
	***REMOVED***
	if len(n.ipamOptions) > 0 ***REMOVED***
		to.ipamOptions = make(map[string]string, len(n.ipamOptions))
		for k, v := range n.ipamOptions ***REMOVED***
			if _, ok := to.ipamOptions[k]; !ok ***REMOVED***
				to.ipamOptions[k] = v
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if len(n.ipamV4Config) > 0 ***REMOVED***
		to.ipamV4Config = make([]*IpamConf, 0, len(n.ipamV4Config))
		to.ipamV4Config = append(to.ipamV4Config, n.ipamV4Config...)
	***REMOVED***
	if len(n.ipamV6Config) > 0 ***REMOVED***
		to.ipamV6Config = make([]*IpamConf, 0, len(n.ipamV6Config))
		to.ipamV6Config = append(to.ipamV6Config, n.ipamV6Config...)
	***REMOVED***
	if len(n.generic) > 0 ***REMOVED***
		to.generic = options.Generic***REMOVED******REMOVED***
		for k, v := range n.generic ***REMOVED***
			to.generic[k] = v
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (n *network) CopyTo(o datastore.KVObject) error ***REMOVED***
	n.Lock()
	defer n.Unlock()

	dstN := o.(*network)
	dstN.name = n.name
	dstN.id = n.id
	dstN.created = n.created
	dstN.networkType = n.networkType
	dstN.scope = n.scope
	dstN.dynamic = n.dynamic
	dstN.ipamType = n.ipamType
	dstN.enableIPv6 = n.enableIPv6
	dstN.persist = n.persist
	dstN.postIPv6 = n.postIPv6
	dstN.dbIndex = n.dbIndex
	dstN.dbExists = n.dbExists
	dstN.drvOnce = n.drvOnce
	dstN.internal = n.internal
	dstN.attachable = n.attachable
	dstN.inDelete = n.inDelete
	dstN.ingress = n.ingress
	dstN.configOnly = n.configOnly
	dstN.configFrom = n.configFrom
	dstN.loadBalancerIP = n.loadBalancerIP

	// copy labels
	if dstN.labels == nil ***REMOVED***
		dstN.labels = make(map[string]string, len(n.labels))
	***REMOVED***
	for k, v := range n.labels ***REMOVED***
		dstN.labels[k] = v
	***REMOVED***

	if n.ipamOptions != nil ***REMOVED***
		dstN.ipamOptions = make(map[string]string, len(n.ipamOptions))
		for k, v := range n.ipamOptions ***REMOVED***
			dstN.ipamOptions[k] = v
		***REMOVED***
	***REMOVED***

	for _, v4conf := range n.ipamV4Config ***REMOVED***
		dstV4Conf := &IpamConf***REMOVED******REMOVED***
		v4conf.CopyTo(dstV4Conf)
		dstN.ipamV4Config = append(dstN.ipamV4Config, dstV4Conf)
	***REMOVED***

	for _, v4info := range n.ipamV4Info ***REMOVED***
		dstV4Info := &IpamInfo***REMOVED******REMOVED***
		v4info.CopyTo(dstV4Info)
		dstN.ipamV4Info = append(dstN.ipamV4Info, dstV4Info)
	***REMOVED***

	for _, v6conf := range n.ipamV6Config ***REMOVED***
		dstV6Conf := &IpamConf***REMOVED******REMOVED***
		v6conf.CopyTo(dstV6Conf)
		dstN.ipamV6Config = append(dstN.ipamV6Config, dstV6Conf)
	***REMOVED***

	for _, v6info := range n.ipamV6Info ***REMOVED***
		dstV6Info := &IpamInfo***REMOVED******REMOVED***
		v6info.CopyTo(dstV6Info)
		dstN.ipamV6Info = append(dstN.ipamV6Info, dstV6Info)
	***REMOVED***

	dstN.generic = options.Generic***REMOVED******REMOVED***
	for k, v := range n.generic ***REMOVED***
		dstN.generic[k] = v
	***REMOVED***

	return nil
***REMOVED***

func (n *network) DataScope() string ***REMOVED***
	s := n.Scope()
	// All swarm scope networks have local datascope
	if s == datastore.SwarmScope ***REMOVED***
		s = datastore.LocalScope
	***REMOVED***
	return s
***REMOVED***

func (n *network) getEpCnt() *endpointCnt ***REMOVED***
	n.Lock()
	defer n.Unlock()

	return n.epCnt
***REMOVED***

// TODO : Can be made much more generic with the help of reflection (but has some golang limitations)
func (n *network) MarshalJSON() ([]byte, error) ***REMOVED***
	netMap := make(map[string]interface***REMOVED******REMOVED***)
	netMap["name"] = n.name
	netMap["id"] = n.id
	netMap["created"] = n.created
	netMap["networkType"] = n.networkType
	netMap["scope"] = n.scope
	netMap["labels"] = n.labels
	netMap["ipamType"] = n.ipamType
	netMap["ipamOptions"] = n.ipamOptions
	netMap["addrSpace"] = n.addrSpace
	netMap["enableIPv6"] = n.enableIPv6
	if n.generic != nil ***REMOVED***
		netMap["generic"] = n.generic
	***REMOVED***
	netMap["persist"] = n.persist
	netMap["postIPv6"] = n.postIPv6
	if len(n.ipamV4Config) > 0 ***REMOVED***
		ics, err := json.Marshal(n.ipamV4Config)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		netMap["ipamV4Config"] = string(ics)
	***REMOVED***
	if len(n.ipamV4Info) > 0 ***REMOVED***
		iis, err := json.Marshal(n.ipamV4Info)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		netMap["ipamV4Info"] = string(iis)
	***REMOVED***
	if len(n.ipamV6Config) > 0 ***REMOVED***
		ics, err := json.Marshal(n.ipamV6Config)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		netMap["ipamV6Config"] = string(ics)
	***REMOVED***
	if len(n.ipamV6Info) > 0 ***REMOVED***
		iis, err := json.Marshal(n.ipamV6Info)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		netMap["ipamV6Info"] = string(iis)
	***REMOVED***
	netMap["internal"] = n.internal
	netMap["attachable"] = n.attachable
	netMap["inDelete"] = n.inDelete
	netMap["ingress"] = n.ingress
	netMap["configOnly"] = n.configOnly
	netMap["configFrom"] = n.configFrom
	netMap["loadBalancerIP"] = n.loadBalancerIP
	return json.Marshal(netMap)
***REMOVED***

// TODO : Can be made much more generic with the help of reflection (but has some golang limitations)
func (n *network) UnmarshalJSON(b []byte) (err error) ***REMOVED***
	var netMap map[string]interface***REMOVED******REMOVED***
	if err := json.Unmarshal(b, &netMap); err != nil ***REMOVED***
		return err
	***REMOVED***
	n.name = netMap["name"].(string)
	n.id = netMap["id"].(string)
	// "created" is not available in older versions
	if v, ok := netMap["created"]; ok ***REMOVED***
		// n.created is time.Time but marshalled as string
		if err = n.created.UnmarshalText([]byte(v.(string))); err != nil ***REMOVED***
			logrus.Warnf("failed to unmarshal creation time %v: %v", v, err)
			n.created = time.Time***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
	n.networkType = netMap["networkType"].(string)
	n.enableIPv6 = netMap["enableIPv6"].(bool)

	// if we weren't unmarshaling to netMap we could simply set n.labels
	// unfortunately, we can't because map[string]interface***REMOVED******REMOVED*** != map[string]string
	if labels, ok := netMap["labels"].(map[string]interface***REMOVED******REMOVED***); ok ***REMOVED***
		n.labels = make(map[string]string, len(labels))
		for label, value := range labels ***REMOVED***
			n.labels[label] = value.(string)
		***REMOVED***
	***REMOVED***

	if v, ok := netMap["ipamOptions"]; ok ***REMOVED***
		if iOpts, ok := v.(map[string]interface***REMOVED******REMOVED***); ok ***REMOVED***
			n.ipamOptions = make(map[string]string, len(iOpts))
			for k, v := range iOpts ***REMOVED***
				n.ipamOptions[k] = v.(string)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if v, ok := netMap["generic"]; ok ***REMOVED***
		n.generic = v.(map[string]interface***REMOVED******REMOVED***)
		// Restore opts in their map[string]string form
		if v, ok := n.generic[netlabel.GenericData]; ok ***REMOVED***
			var lmap map[string]string
			ba, err := json.Marshal(v)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if err := json.Unmarshal(ba, &lmap); err != nil ***REMOVED***
				return err
			***REMOVED***
			n.generic[netlabel.GenericData] = lmap
		***REMOVED***
	***REMOVED***
	if v, ok := netMap["persist"]; ok ***REMOVED***
		n.persist = v.(bool)
	***REMOVED***
	if v, ok := netMap["postIPv6"]; ok ***REMOVED***
		n.postIPv6 = v.(bool)
	***REMOVED***
	if v, ok := netMap["ipamType"]; ok ***REMOVED***
		n.ipamType = v.(string)
	***REMOVED*** else ***REMOVED***
		n.ipamType = ipamapi.DefaultIPAM
	***REMOVED***
	if v, ok := netMap["addrSpace"]; ok ***REMOVED***
		n.addrSpace = v.(string)
	***REMOVED***
	if v, ok := netMap["ipamV4Config"]; ok ***REMOVED***
		if err := json.Unmarshal([]byte(v.(string)), &n.ipamV4Config); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if v, ok := netMap["ipamV4Info"]; ok ***REMOVED***
		if err := json.Unmarshal([]byte(v.(string)), &n.ipamV4Info); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if v, ok := netMap["ipamV6Config"]; ok ***REMOVED***
		if err := json.Unmarshal([]byte(v.(string)), &n.ipamV6Config); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if v, ok := netMap["ipamV6Info"]; ok ***REMOVED***
		if err := json.Unmarshal([]byte(v.(string)), &n.ipamV6Info); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if v, ok := netMap["internal"]; ok ***REMOVED***
		n.internal = v.(bool)
	***REMOVED***
	if v, ok := netMap["attachable"]; ok ***REMOVED***
		n.attachable = v.(bool)
	***REMOVED***
	if s, ok := netMap["scope"]; ok ***REMOVED***
		n.scope = s.(string)
	***REMOVED***
	if v, ok := netMap["inDelete"]; ok ***REMOVED***
		n.inDelete = v.(bool)
	***REMOVED***
	if v, ok := netMap["ingress"]; ok ***REMOVED***
		n.ingress = v.(bool)
	***REMOVED***
	if v, ok := netMap["configOnly"]; ok ***REMOVED***
		n.configOnly = v.(bool)
	***REMOVED***
	if v, ok := netMap["configFrom"]; ok ***REMOVED***
		n.configFrom = v.(string)
	***REMOVED***
	if v, ok := netMap["loadBalancerIP"]; ok ***REMOVED***
		n.loadBalancerIP = net.ParseIP(v.(string))
	***REMOVED***
	// Reconcile old networks with the recently added `--ipv6` flag
	if !n.enableIPv6 ***REMOVED***
		n.enableIPv6 = len(n.ipamV6Info) > 0
	***REMOVED***
	return nil
***REMOVED***

// NetworkOption is an option setter function type used to pass various options to
// NewNetwork method. The various setter functions of type NetworkOption are
// provided by libnetwork, they look like NetworkOptionXXXX(...)
type NetworkOption func(n *network)

// NetworkOptionGeneric function returns an option setter for a Generic option defined
// in a Dictionary of Key-Value pair
func NetworkOptionGeneric(generic map[string]interface***REMOVED******REMOVED***) NetworkOption ***REMOVED***
	return func(n *network) ***REMOVED***
		if n.generic == nil ***REMOVED***
			n.generic = make(map[string]interface***REMOVED******REMOVED***)
		***REMOVED***
		if val, ok := generic[netlabel.EnableIPv6]; ok ***REMOVED***
			n.enableIPv6 = val.(bool)
		***REMOVED***
		if val, ok := generic[netlabel.Internal]; ok ***REMOVED***
			n.internal = val.(bool)
		***REMOVED***
		for k, v := range generic ***REMOVED***
			n.generic[k] = v
		***REMOVED***
	***REMOVED***
***REMOVED***

// NetworkOptionIngress returns an option setter to indicate if a network is
// an ingress network.
func NetworkOptionIngress(ingress bool) NetworkOption ***REMOVED***
	return func(n *network) ***REMOVED***
		n.ingress = ingress
	***REMOVED***
***REMOVED***

// NetworkOptionPersist returns an option setter to set persistence policy for a network
func NetworkOptionPersist(persist bool) NetworkOption ***REMOVED***
	return func(n *network) ***REMOVED***
		n.persist = persist
	***REMOVED***
***REMOVED***

// NetworkOptionEnableIPv6 returns an option setter to explicitly configure IPv6
func NetworkOptionEnableIPv6(enableIPv6 bool) NetworkOption ***REMOVED***
	return func(n *network) ***REMOVED***
		if n.generic == nil ***REMOVED***
			n.generic = make(map[string]interface***REMOVED******REMOVED***)
		***REMOVED***
		n.enableIPv6 = enableIPv6
		n.generic[netlabel.EnableIPv6] = enableIPv6
	***REMOVED***
***REMOVED***

// NetworkOptionInternalNetwork returns an option setter to config the network
// to be internal which disables default gateway service
func NetworkOptionInternalNetwork() NetworkOption ***REMOVED***
	return func(n *network) ***REMOVED***
		if n.generic == nil ***REMOVED***
			n.generic = make(map[string]interface***REMOVED******REMOVED***)
		***REMOVED***
		n.internal = true
		n.generic[netlabel.Internal] = true
	***REMOVED***
***REMOVED***

// NetworkOptionAttachable returns an option setter to set attachable for a network
func NetworkOptionAttachable(attachable bool) NetworkOption ***REMOVED***
	return func(n *network) ***REMOVED***
		n.attachable = attachable
	***REMOVED***
***REMOVED***

// NetworkOptionScope returns an option setter to overwrite the network's scope.
// By default the network's scope is set to the network driver's datascope.
func NetworkOptionScope(scope string) NetworkOption ***REMOVED***
	return func(n *network) ***REMOVED***
		n.scope = scope
	***REMOVED***
***REMOVED***

// NetworkOptionIpam function returns an option setter for the ipam configuration for this network
func NetworkOptionIpam(ipamDriver string, addrSpace string, ipV4 []*IpamConf, ipV6 []*IpamConf, opts map[string]string) NetworkOption ***REMOVED***
	return func(n *network) ***REMOVED***
		if ipamDriver != "" ***REMOVED***
			n.ipamType = ipamDriver
			if ipamDriver == ipamapi.DefaultIPAM ***REMOVED***
				n.ipamType = defaultIpamForNetworkType(n.Type())
			***REMOVED***
		***REMOVED***
		n.ipamOptions = opts
		n.addrSpace = addrSpace
		n.ipamV4Config = ipV4
		n.ipamV6Config = ipV6
	***REMOVED***
***REMOVED***

// NetworkOptionLBEndpoint function returns an option setter for the configuration of the load balancer endpoint for this network
func NetworkOptionLBEndpoint(ip net.IP) NetworkOption ***REMOVED***
	return func(n *network) ***REMOVED***
		n.loadBalancerIP = ip
	***REMOVED***
***REMOVED***

// NetworkOptionDriverOpts function returns an option setter for any driver parameter described by a map
func NetworkOptionDriverOpts(opts map[string]string) NetworkOption ***REMOVED***
	return func(n *network) ***REMOVED***
		if n.generic == nil ***REMOVED***
			n.generic = make(map[string]interface***REMOVED******REMOVED***)
		***REMOVED***
		if opts == nil ***REMOVED***
			opts = make(map[string]string)
		***REMOVED***
		// Store the options
		n.generic[netlabel.GenericData] = opts
	***REMOVED***
***REMOVED***

// NetworkOptionLabels function returns an option setter for labels specific to a network
func NetworkOptionLabels(labels map[string]string) NetworkOption ***REMOVED***
	return func(n *network) ***REMOVED***
		n.labels = labels
	***REMOVED***
***REMOVED***

// NetworkOptionDynamic function returns an option setter for dynamic option for a network
func NetworkOptionDynamic() NetworkOption ***REMOVED***
	return func(n *network) ***REMOVED***
		n.dynamic = true
	***REMOVED***
***REMOVED***

// NetworkOptionDeferIPv6Alloc instructs the network to defer the IPV6 address allocation until after the endpoint has been created
// It is being provided to support the specific docker daemon flags where user can deterministically assign an IPv6 address
// to a container as combination of fixed-cidr-v6 + mac-address
// TODO: Remove this option setter once we support endpoint ipam options
func NetworkOptionDeferIPv6Alloc(enable bool) NetworkOption ***REMOVED***
	return func(n *network) ***REMOVED***
		n.postIPv6 = enable
	***REMOVED***
***REMOVED***

// NetworkOptionConfigOnly tells controller this network is
// a configuration only network. It serves as a configuration
// for other networks.
func NetworkOptionConfigOnly() NetworkOption ***REMOVED***
	return func(n *network) ***REMOVED***
		n.configOnly = true
	***REMOVED***
***REMOVED***

// NetworkOptionConfigFrom tells controller to pick the
// network configuration from a configuration only network
func NetworkOptionConfigFrom(name string) NetworkOption ***REMOVED***
	return func(n *network) ***REMOVED***
		n.configFrom = name
	***REMOVED***
***REMOVED***

func (n *network) processOptions(options ...NetworkOption) ***REMOVED***
	for _, opt := range options ***REMOVED***
		if opt != nil ***REMOVED***
			opt(n)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (n *network) resolveDriver(name string, load bool) (driverapi.Driver, *driverapi.Capability, error) ***REMOVED***
	c := n.getController()

	// Check if a driver for the specified network type is available
	d, cap := c.drvRegistry.Driver(name)
	if d == nil ***REMOVED***
		if load ***REMOVED***
			err := c.loadDriver(name)
			if err != nil ***REMOVED***
				return nil, nil, err
			***REMOVED***

			d, cap = c.drvRegistry.Driver(name)
			if d == nil ***REMOVED***
				return nil, nil, fmt.Errorf("could not resolve driver %s in registry", name)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// don't fail if driver loading is not required
			return nil, nil, nil
		***REMOVED***
	***REMOVED***

	return d, cap, nil
***REMOVED***

func (n *network) driverScope() string ***REMOVED***
	_, cap, err := n.resolveDriver(n.networkType, true)
	if err != nil ***REMOVED***
		// If driver could not be resolved simply return an empty string
		return ""
	***REMOVED***

	return cap.DataScope
***REMOVED***

func (n *network) driverIsMultihost() bool ***REMOVED***
	_, cap, err := n.resolveDriver(n.networkType, true)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	return cap.ConnectivityScope == datastore.GlobalScope
***REMOVED***

func (n *network) driver(load bool) (driverapi.Driver, error) ***REMOVED***
	d, cap, err := n.resolveDriver(n.networkType, load)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	n.Lock()
	// If load is not required, driver, cap and err may all be nil
	if n.scope == "" && cap != nil ***REMOVED***
		n.scope = cap.DataScope
	***REMOVED***
	if n.dynamic ***REMOVED***
		// If the network is dynamic, then it is swarm
		// scoped regardless of the backing driver.
		n.scope = datastore.SwarmScope
	***REMOVED***
	n.Unlock()
	return d, nil
***REMOVED***

func (n *network) Delete() error ***REMOVED***
	return n.delete(false)
***REMOVED***

func (n *network) delete(force bool) error ***REMOVED***
	n.Lock()
	c := n.ctrlr
	name := n.name
	id := n.id
	n.Unlock()

	c.networkLocker.Lock(id)
	defer c.networkLocker.Unlock(id)

	n, err := c.getNetworkFromStore(id)
	if err != nil ***REMOVED***
		return &UnknownNetworkError***REMOVED***name: name, id: id***REMOVED***
	***REMOVED***

	if len(n.loadBalancerIP) != 0 ***REMOVED***
		endpoints := n.Endpoints()
		if force || len(endpoints) == 1 ***REMOVED***
			n.deleteLoadBalancerSandbox()
		***REMOVED***
		//Reload the network from the store to update the epcnt.
		n, err = c.getNetworkFromStore(id)
		if err != nil ***REMOVED***
			return &UnknownNetworkError***REMOVED***name: name, id: id***REMOVED***
		***REMOVED***
	***REMOVED***

	if !force && n.getEpCnt().EndpointCnt() != 0 ***REMOVED***
		if n.configOnly ***REMOVED***
			return types.ForbiddenErrorf("configuration network %q is in use", n.Name())
		***REMOVED***
		return &ActiveEndpointsError***REMOVED***name: n.name, id: n.id***REMOVED***
	***REMOVED***

	// Mark the network for deletion
	n.inDelete = true
	if err = c.updateToStore(n); err != nil ***REMOVED***
		return fmt.Errorf("error marking network %s (%s) for deletion: %v", n.Name(), n.ID(), err)
	***REMOVED***

	if n.ConfigFrom() != "" ***REMOVED***
		if t, err := c.getConfigNetwork(n.ConfigFrom()); err == nil ***REMOVED***
			if err := t.getEpCnt().DecEndpointCnt(); err != nil ***REMOVED***
				logrus.Warnf("Failed to update reference count for configuration network %q on removal of network %q: %v",
					t.Name(), n.Name(), err)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			logrus.Warnf("Could not find configuration network %q during removal of network %q", n.configOnly, n.Name())
		***REMOVED***
	***REMOVED***

	if n.configOnly ***REMOVED***
		goto removeFromStore
	***REMOVED***

	if err = n.deleteNetwork(); err != nil ***REMOVED***
		if !force ***REMOVED***
			return err
		***REMOVED***
		logrus.Debugf("driver failed to delete stale network %s (%s): %v", n.Name(), n.ID(), err)
	***REMOVED***

	n.ipamRelease()
	if err = c.updateToStore(n); err != nil ***REMOVED***
		logrus.Warnf("Failed to update store after ipam release for network %s (%s): %v", n.Name(), n.ID(), err)
	***REMOVED***

	// We are about to delete the network. Leave the gossip
	// cluster for the network to stop all incoming network
	// specific gossip updates before cleaning up all the service
	// bindings for the network. But cleanup service binding
	// before deleting the network from the store since service
	// bindings cleanup requires the network in the store.
	n.cancelDriverWatches()
	if err = n.leaveCluster(); err != nil ***REMOVED***
		logrus.Errorf("Failed leaving network %s from the agent cluster: %v", n.Name(), err)
	***REMOVED***

	// Cleanup the service discovery for this network
	c.cleanupServiceDiscovery(n.ID())

	// Cleanup the load balancer
	c.cleanupServiceBindings(n.ID())

removeFromStore:
	// deleteFromStore performs an atomic delete operation and the
	// network.epCnt will help prevent any possible
	// race between endpoint join and network delete
	if err = c.deleteFromStore(n.getEpCnt()); err != nil ***REMOVED***
		if !force ***REMOVED***
			return fmt.Errorf("error deleting network endpoint count from store: %v", err)
		***REMOVED***
		logrus.Debugf("Error deleting endpoint count from store for stale network %s (%s) for deletion: %v", n.Name(), n.ID(), err)
	***REMOVED***

	if err = c.deleteFromStore(n); err != nil ***REMOVED***
		return fmt.Errorf("error deleting network from store: %v", err)
	***REMOVED***

	return nil
***REMOVED***

func (n *network) deleteNetwork() error ***REMOVED***
	d, err := n.driver(true)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed deleting network: %v", err)
	***REMOVED***

	if err := d.DeleteNetwork(n.ID()); err != nil ***REMOVED***
		// Forbidden Errors should be honored
		if _, ok := err.(types.ForbiddenError); ok ***REMOVED***
			return err
		***REMOVED***

		if _, ok := err.(types.MaskableError); !ok ***REMOVED***
			logrus.Warnf("driver error deleting network %s : %v", n.name, err)
		***REMOVED***
	***REMOVED***

	for _, resolver := range n.resolver ***REMOVED***
		resolver.Stop()
	***REMOVED***
	return nil
***REMOVED***

func (n *network) addEndpoint(ep *endpoint) error ***REMOVED***
	d, err := n.driver(true)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to add endpoint: %v", err)
	***REMOVED***

	err = d.CreateEndpoint(n.id, ep.id, ep.Interface(), ep.generic)
	if err != nil ***REMOVED***
		return types.InternalErrorf("failed to create endpoint %s on network %s: %v",
			ep.Name(), n.Name(), err)
	***REMOVED***

	return nil
***REMOVED***

func (n *network) CreateEndpoint(name string, options ...EndpointOption) (Endpoint, error) ***REMOVED***
	var err error
	if !config.IsValidName(name) ***REMOVED***
		return nil, ErrInvalidName(name)
	***REMOVED***

	if n.ConfigOnly() ***REMOVED***
		return nil, types.ForbiddenErrorf("cannot create endpoint on configuration-only network")
	***REMOVED***

	if _, err = n.EndpointByName(name); err == nil ***REMOVED***
		return nil, types.ForbiddenErrorf("endpoint with name %s already exists in network %s", name, n.Name())
	***REMOVED***

	n.ctrlr.networkLocker.Lock(n.id)
	defer n.ctrlr.networkLocker.Unlock(n.id)

	return n.createEndpoint(name, options...)

***REMOVED***

func (n *network) createEndpoint(name string, options ...EndpointOption) (Endpoint, error) ***REMOVED***
	var err error

	ep := &endpoint***REMOVED***name: name, generic: make(map[string]interface***REMOVED******REMOVED***), iface: &endpointInterface***REMOVED******REMOVED******REMOVED***
	ep.id = stringid.GenerateRandomID()

	// Initialize ep.network with a possibly stale copy of n. We need this to get network from
	// store. But once we get it from store we will have the most uptodate copy possibly.
	ep.network = n
	ep.locator = n.getController().clusterHostID()
	ep.network, err = ep.getNetworkFromStore()
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to get network during CreateEndpoint: %v", err)
	***REMOVED***
	n = ep.network

	ep.processOptions(options...)

	for _, llIPNet := range ep.Iface().LinkLocalAddresses() ***REMOVED***
		if !llIPNet.IP.IsLinkLocalUnicast() ***REMOVED***
			return nil, types.BadRequestErrorf("invalid link local IP address: %v", llIPNet.IP)
		***REMOVED***
	***REMOVED***

	if opt, ok := ep.generic[netlabel.MacAddress]; ok ***REMOVED***
		if mac, ok := opt.(net.HardwareAddr); ok ***REMOVED***
			ep.iface.mac = mac
		***REMOVED***
	***REMOVED***

	ipam, cap, err := n.getController().getIPAMDriver(n.ipamType)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if cap.RequiresMACAddress ***REMOVED***
		if ep.iface.mac == nil ***REMOVED***
			ep.iface.mac = netutils.GenerateRandomMAC()
		***REMOVED***
		if ep.ipamOptions == nil ***REMOVED***
			ep.ipamOptions = make(map[string]string)
		***REMOVED***
		ep.ipamOptions[netlabel.MacAddress] = ep.iface.mac.String()
	***REMOVED***

	if err = ep.assignAddress(ipam, true, n.enableIPv6 && !n.postIPv6); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			ep.releaseAddress()
		***REMOVED***
	***REMOVED***()
	// Moving updateToSTore before calling addEndpoint so that we shall clean up VETH interfaces in case
	// DockerD get killed between addEndpoint and updateSTore call
	if err = n.getController().updateToStore(ep); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			if e := n.getController().deleteFromStore(ep); e != nil ***REMOVED***
				logrus.Warnf("error rolling back endpoint %s from store: %v", name, e)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	if err = n.addEndpoint(ep); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			if e := ep.deleteEndpoint(false); e != nil ***REMOVED***
				logrus.Warnf("cleaning up endpoint failed %s : %v", name, e)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	if err = ep.assignAddress(ipam, false, n.enableIPv6 && n.postIPv6); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Watch for service records
	n.getController().watchSvcRecord(ep)
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			n.getController().unWatchSvcRecord(ep)
		***REMOVED***
	***REMOVED***()

	// Increment endpoint count to indicate completion of endpoint addition
	if err = n.getEpCnt().IncEndpointCnt(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return ep, nil
***REMOVED***

func (n *network) Endpoints() []Endpoint ***REMOVED***
	var list []Endpoint

	endpoints, err := n.getEndpointsFromStore()
	if err != nil ***REMOVED***
		logrus.Error(err)
	***REMOVED***

	for _, ep := range endpoints ***REMOVED***
		list = append(list, ep)
	***REMOVED***

	return list
***REMOVED***

func (n *network) WalkEndpoints(walker EndpointWalker) ***REMOVED***
	for _, e := range n.Endpoints() ***REMOVED***
		if walker(e) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (n *network) EndpointByName(name string) (Endpoint, error) ***REMOVED***
	if name == "" ***REMOVED***
		return nil, ErrInvalidName(name)
	***REMOVED***
	var e Endpoint

	s := func(current Endpoint) bool ***REMOVED***
		if current.Name() == name ***REMOVED***
			e = current
			return true
		***REMOVED***
		return false
	***REMOVED***

	n.WalkEndpoints(s)

	if e == nil ***REMOVED***
		return nil, ErrNoSuchEndpoint(name)
	***REMOVED***

	return e, nil
***REMOVED***

func (n *network) EndpointByID(id string) (Endpoint, error) ***REMOVED***
	if id == "" ***REMOVED***
		return nil, ErrInvalidID(id)
	***REMOVED***

	ep, err := n.getEndpointFromStore(id)
	if err != nil ***REMOVED***
		return nil, ErrNoSuchEndpoint(id)
	***REMOVED***

	return ep, nil
***REMOVED***

func (n *network) updateSvcRecord(ep *endpoint, localEps []*endpoint, isAdd bool) ***REMOVED***
	var ipv6 net.IP
	epName := ep.Name()
	if iface := ep.Iface(); iface.Address() != nil ***REMOVED***
		myAliases := ep.MyAliases()
		if iface.AddressIPv6() != nil ***REMOVED***
			ipv6 = iface.AddressIPv6().IP
		***REMOVED***

		serviceID := ep.svcID
		if serviceID == "" ***REMOVED***
			serviceID = ep.ID()
		***REMOVED***
		if isAdd ***REMOVED***
			// If anonymous endpoint has an alias use the first alias
			// for ip->name mapping. Not having the reverse mapping
			// breaks some apps
			if ep.isAnonymous() ***REMOVED***
				if len(myAliases) > 0 ***REMOVED***
					n.addSvcRecords(ep.ID(), myAliases[0], serviceID, iface.Address().IP, ipv6, true, "updateSvcRecord")
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				n.addSvcRecords(ep.ID(), epName, serviceID, iface.Address().IP, ipv6, true, "updateSvcRecord")
			***REMOVED***
			for _, alias := range myAliases ***REMOVED***
				n.addSvcRecords(ep.ID(), alias, serviceID, iface.Address().IP, ipv6, false, "updateSvcRecord")
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if ep.isAnonymous() ***REMOVED***
				if len(myAliases) > 0 ***REMOVED***
					n.deleteSvcRecords(ep.ID(), myAliases[0], serviceID, iface.Address().IP, ipv6, true, "updateSvcRecord")
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				n.deleteSvcRecords(ep.ID(), epName, serviceID, iface.Address().IP, ipv6, true, "updateSvcRecord")
			***REMOVED***
			for _, alias := range myAliases ***REMOVED***
				n.deleteSvcRecords(ep.ID(), alias, serviceID, iface.Address().IP, ipv6, false, "updateSvcRecord")
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func addIPToName(ipMap common.SetMatrix, name, serviceID string, ip net.IP) ***REMOVED***
	reverseIP := netutils.ReverseIP(ip.String())
	ipMap.Insert(reverseIP, ipInfo***REMOVED***
		name:      name,
		serviceID: serviceID,
	***REMOVED***)
***REMOVED***

func delIPToName(ipMap common.SetMatrix, name, serviceID string, ip net.IP) ***REMOVED***
	reverseIP := netutils.ReverseIP(ip.String())
	ipMap.Remove(reverseIP, ipInfo***REMOVED***
		name:      name,
		serviceID: serviceID,
	***REMOVED***)
***REMOVED***

func addNameToIP(svcMap common.SetMatrix, name, serviceID string, epIP net.IP) ***REMOVED***
	svcMap.Insert(name, svcMapEntry***REMOVED***
		ip:        epIP.String(),
		serviceID: serviceID,
	***REMOVED***)
***REMOVED***

func delNameToIP(svcMap common.SetMatrix, name, serviceID string, epIP net.IP) ***REMOVED***
	svcMap.Remove(name, svcMapEntry***REMOVED***
		ip:        epIP.String(),
		serviceID: serviceID,
	***REMOVED***)
***REMOVED***

func (n *network) addSvcRecords(eID, name, serviceID string, epIP, epIPv6 net.IP, ipMapUpdate bool, method string) ***REMOVED***
	// Do not add service names for ingress network as this is a
	// routing only network
	if n.ingress ***REMOVED***
		return
	***REMOVED***

	logrus.Debugf("%s (%s).addSvcRecords(%s, %s, %s, %t) %s sid:%s", eID, n.ID()[0:7], name, epIP, epIPv6, ipMapUpdate, method, serviceID)

	c := n.getController()
	c.Lock()
	defer c.Unlock()

	sr, ok := c.svcRecords[n.ID()]
	if !ok ***REMOVED***
		sr = svcInfo***REMOVED***
			svcMap:     common.NewSetMatrix(),
			svcIPv6Map: common.NewSetMatrix(),
			ipMap:      common.NewSetMatrix(),
		***REMOVED***
		c.svcRecords[n.ID()] = sr
	***REMOVED***

	if ipMapUpdate ***REMOVED***
		addIPToName(sr.ipMap, name, serviceID, epIP)
		if epIPv6 != nil ***REMOVED***
			addIPToName(sr.ipMap, name, serviceID, epIPv6)
		***REMOVED***
	***REMOVED***

	addNameToIP(sr.svcMap, name, serviceID, epIP)
	if epIPv6 != nil ***REMOVED***
		addNameToIP(sr.svcIPv6Map, name, serviceID, epIPv6)
	***REMOVED***
***REMOVED***

func (n *network) deleteSvcRecords(eID, name, serviceID string, epIP net.IP, epIPv6 net.IP, ipMapUpdate bool, method string) ***REMOVED***
	// Do not delete service names from ingress network as this is a
	// routing only network
	if n.ingress ***REMOVED***
		return
	***REMOVED***

	logrus.Debugf("%s (%s).deleteSvcRecords(%s, %s, %s, %t) %s sid:%s ", eID, n.ID()[0:7], name, epIP, epIPv6, ipMapUpdate, method, serviceID)

	c := n.getController()
	c.Lock()
	defer c.Unlock()

	sr, ok := c.svcRecords[n.ID()]
	if !ok ***REMOVED***
		return
	***REMOVED***

	if ipMapUpdate ***REMOVED***
		delIPToName(sr.ipMap, name, serviceID, epIP)

		if epIPv6 != nil ***REMOVED***
			delIPToName(sr.ipMap, name, serviceID, epIPv6)
		***REMOVED***
	***REMOVED***

	delNameToIP(sr.svcMap, name, serviceID, epIP)

	if epIPv6 != nil ***REMOVED***
		delNameToIP(sr.svcIPv6Map, name, serviceID, epIPv6)
	***REMOVED***
***REMOVED***

func (n *network) getSvcRecords(ep *endpoint) []etchosts.Record ***REMOVED***
	n.Lock()
	defer n.Unlock()

	if ep == nil ***REMOVED***
		return nil
	***REMOVED***

	var recs []etchosts.Record

	epName := ep.Name()

	n.ctrlr.Lock()
	defer n.ctrlr.Unlock()
	sr, ok := n.ctrlr.svcRecords[n.id]
	if !ok || sr.svcMap == nil ***REMOVED***
		return nil
	***REMOVED***

	svcMapKeys := sr.svcMap.Keys()
	// Loop on service names on this network
	for _, k := range svcMapKeys ***REMOVED***
		if strings.Split(k, ".")[0] == epName ***REMOVED***
			continue
		***REMOVED***
		// Get all the IPs associated to this service
		mapEntryList, ok := sr.svcMap.Get(k)
		if !ok ***REMOVED***
			// The key got deleted
			continue
		***REMOVED***
		if len(mapEntryList) == 0 ***REMOVED***
			logrus.Warnf("Found empty list of IP addresses for service %s on network %s (%s)", k, n.name, n.id)
			continue
		***REMOVED***

		recs = append(recs, etchosts.Record***REMOVED***
			Hosts: k,
			IP:    mapEntryList[0].(svcMapEntry).ip,
		***REMOVED***)
	***REMOVED***

	return recs
***REMOVED***

func (n *network) getController() *controller ***REMOVED***
	n.Lock()
	defer n.Unlock()
	return n.ctrlr
***REMOVED***

func (n *network) ipamAllocate() error ***REMOVED***
	if n.hasSpecialDriver() ***REMOVED***
		return nil
	***REMOVED***

	ipam, _, err := n.getController().getIPAMDriver(n.ipamType)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if n.addrSpace == "" ***REMOVED***
		if n.addrSpace, err = n.deriveAddressSpace(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	err = n.ipamAllocateVersion(4, ipam)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			n.ipamReleaseVersion(4, ipam)
		***REMOVED***
	***REMOVED***()

	if !n.enableIPv6 ***REMOVED***
		return nil
	***REMOVED***

	err = n.ipamAllocateVersion(6, ipam)
	return err
***REMOVED***

func (n *network) requestPoolHelper(ipam ipamapi.Ipam, addressSpace, preferredPool, subPool string, options map[string]string, v6 bool) (string, *net.IPNet, map[string]string, error) ***REMOVED***
	for ***REMOVED***
		poolID, pool, meta, err := ipam.RequestPool(addressSpace, preferredPool, subPool, options, v6)
		if err != nil ***REMOVED***
			return "", nil, nil, err
		***REMOVED***

		// If the network belongs to global scope or the pool was
		// explicitly chosen or it is invalid, do not perform the overlap check.
		if n.Scope() == datastore.GlobalScope || preferredPool != "" || !types.IsIPNetValid(pool) ***REMOVED***
			return poolID, pool, meta, nil
		***REMOVED***

		// Check for overlap and if none found, we have found the right pool.
		if _, err := netutils.FindAvailableNetwork([]*net.IPNet***REMOVED***pool***REMOVED***); err == nil ***REMOVED***
			return poolID, pool, meta, nil
		***REMOVED***

		// Pool obtained in this iteration is
		// overlapping. Hold onto the pool and don't release
		// it yet, because we don't want ipam to give us back
		// the same pool over again. But make sure we still do
		// a deferred release when we have either obtained a
		// non-overlapping pool or ran out of pre-defined
		// pools.
		defer func() ***REMOVED***
			if err := ipam.ReleasePool(poolID); err != nil ***REMOVED***
				logrus.Warnf("Failed to release overlapping pool %s while returning from pool request helper for network %s", pool, n.Name())
			***REMOVED***
		***REMOVED***()

		// If this is a preferred pool request and the network
		// is local scope and there is an overlap, we fail the
		// network creation right here. The pool will be
		// released in the defer.
		if preferredPool != "" ***REMOVED***
			return "", nil, nil, fmt.Errorf("requested subnet %s overlaps in the host", preferredPool)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (n *network) ipamAllocateVersion(ipVer int, ipam ipamapi.Ipam) error ***REMOVED***
	var (
		cfgList  *[]*IpamConf
		infoList *[]*IpamInfo
		err      error
	)

	switch ipVer ***REMOVED***
	case 4:
		cfgList = &n.ipamV4Config
		infoList = &n.ipamV4Info
	case 6:
		cfgList = &n.ipamV6Config
		infoList = &n.ipamV6Info
	default:
		return types.InternalErrorf("incorrect ip version passed to ipam allocate: %d", ipVer)
	***REMOVED***

	if len(*cfgList) == 0 ***REMOVED***
		*cfgList = []*IpamConf***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	*infoList = make([]*IpamInfo, len(*cfgList))

	logrus.Debugf("Allocating IPv%d pools for network %s (%s)", ipVer, n.Name(), n.ID())

	for i, cfg := range *cfgList ***REMOVED***
		if err = cfg.Validate(); err != nil ***REMOVED***
			return err
		***REMOVED***
		d := &IpamInfo***REMOVED******REMOVED***
		(*infoList)[i] = d

		d.AddressSpace = n.addrSpace
		d.PoolID, d.Pool, d.Meta, err = n.requestPoolHelper(ipam, n.addrSpace, cfg.PreferredPool, cfg.SubPool, n.ipamOptions, ipVer == 6)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		defer func() ***REMOVED***
			if err != nil ***REMOVED***
				if err := ipam.ReleasePool(d.PoolID); err != nil ***REMOVED***
					logrus.Warnf("Failed to release address pool %s after failure to create network %s (%s)", d.PoolID, n.Name(), n.ID())
				***REMOVED***
			***REMOVED***
		***REMOVED***()

		if gws, ok := d.Meta[netlabel.Gateway]; ok ***REMOVED***
			if d.Gateway, err = types.ParseCIDR(gws); err != nil ***REMOVED***
				return types.BadRequestErrorf("failed to parse gateway address (%v) returned by ipam driver: %v", gws, err)
			***REMOVED***
		***REMOVED***

		// If user requested a specific gateway, libnetwork will allocate it
		// irrespective of whether ipam driver returned a gateway already.
		// If none of the above is true, libnetwork will allocate one.
		if cfg.Gateway != "" || d.Gateway == nil ***REMOVED***
			var gatewayOpts = map[string]string***REMOVED***
				ipamapi.RequestAddressType: netlabel.Gateway,
			***REMOVED***
			if d.Gateway, _, err = ipam.RequestAddress(d.PoolID, net.ParseIP(cfg.Gateway), gatewayOpts); err != nil ***REMOVED***
				return types.InternalErrorf("failed to allocate gateway (%v): %v", cfg.Gateway, err)
			***REMOVED***
		***REMOVED***

		// Auxiliary addresses must be part of the master address pool
		// If they fall into the container addressable pool, libnetwork will reserve them
		if cfg.AuxAddresses != nil ***REMOVED***
			var ip net.IP
			d.IPAMData.AuxAddresses = make(map[string]*net.IPNet, len(cfg.AuxAddresses))
			for k, v := range cfg.AuxAddresses ***REMOVED***
				if ip = net.ParseIP(v); ip == nil ***REMOVED***
					return types.BadRequestErrorf("non parsable secondary ip address (%s:%s) passed for network %s", k, v, n.Name())
				***REMOVED***
				if !d.Pool.Contains(ip) ***REMOVED***
					return types.ForbiddenErrorf("auxilairy address: (%s:%s) must belong to the master pool: %s", k, v, d.Pool)
				***REMOVED***
				// Attempt reservation in the container addressable pool, silent the error if address does not belong to that pool
				if d.IPAMData.AuxAddresses[k], _, err = ipam.RequestAddress(d.PoolID, ip, nil); err != nil && err != ipamapi.ErrIPOutOfRange ***REMOVED***
					return types.InternalErrorf("failed to allocate secondary ip address (%s:%s): %v", k, v, err)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (n *network) ipamRelease() ***REMOVED***
	if n.hasSpecialDriver() ***REMOVED***
		return
	***REMOVED***
	ipam, _, err := n.getController().getIPAMDriver(n.ipamType)
	if err != nil ***REMOVED***
		logrus.Warnf("Failed to retrieve ipam driver to release address pool(s) on delete of network %s (%s): %v", n.Name(), n.ID(), err)
		return
	***REMOVED***
	n.ipamReleaseVersion(4, ipam)
	n.ipamReleaseVersion(6, ipam)
***REMOVED***

func (n *network) ipamReleaseVersion(ipVer int, ipam ipamapi.Ipam) ***REMOVED***
	var infoList *[]*IpamInfo

	switch ipVer ***REMOVED***
	case 4:
		infoList = &n.ipamV4Info
	case 6:
		infoList = &n.ipamV6Info
	default:
		logrus.Warnf("incorrect ip version passed to ipam release: %d", ipVer)
		return
	***REMOVED***

	if len(*infoList) == 0 ***REMOVED***
		return
	***REMOVED***

	logrus.Debugf("releasing IPv%d pools from network %s (%s)", ipVer, n.Name(), n.ID())

	for _, d := range *infoList ***REMOVED***
		if d.Gateway != nil ***REMOVED***
			if err := ipam.ReleaseAddress(d.PoolID, d.Gateway.IP); err != nil ***REMOVED***
				logrus.Warnf("Failed to release gateway ip address %s on delete of network %s (%s): %v", d.Gateway.IP, n.Name(), n.ID(), err)
			***REMOVED***
		***REMOVED***
		if d.IPAMData.AuxAddresses != nil ***REMOVED***
			for k, nw := range d.IPAMData.AuxAddresses ***REMOVED***
				if d.Pool.Contains(nw.IP) ***REMOVED***
					if err := ipam.ReleaseAddress(d.PoolID, nw.IP); err != nil && err != ipamapi.ErrIPOutOfRange ***REMOVED***
						logrus.Warnf("Failed to release secondary ip address %s (%v) on delete of network %s (%s): %v", k, nw.IP, n.Name(), n.ID(), err)
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if err := ipam.ReleasePool(d.PoolID); err != nil ***REMOVED***
			logrus.Warnf("Failed to release address pool %s on delete of network %s (%s): %v", d.PoolID, n.Name(), n.ID(), err)
		***REMOVED***
	***REMOVED***

	*infoList = nil
***REMOVED***

func (n *network) getIPInfo(ipVer int) []*IpamInfo ***REMOVED***
	var info []*IpamInfo
	switch ipVer ***REMOVED***
	case 4:
		info = n.ipamV4Info
	case 6:
		info = n.ipamV6Info
	default:
		return nil
	***REMOVED***
	l := make([]*IpamInfo, 0, len(info))
	n.Lock()
	l = append(l, info...)
	n.Unlock()
	return l
***REMOVED***

func (n *network) getIPData(ipVer int) []driverapi.IPAMData ***REMOVED***
	var info []*IpamInfo
	switch ipVer ***REMOVED***
	case 4:
		info = n.ipamV4Info
	case 6:
		info = n.ipamV6Info
	default:
		return nil
	***REMOVED***
	l := make([]driverapi.IPAMData, 0, len(info))
	n.Lock()
	for _, d := range info ***REMOVED***
		l = append(l, d.IPAMData)
	***REMOVED***
	n.Unlock()
	return l
***REMOVED***

func (n *network) deriveAddressSpace() (string, error) ***REMOVED***
	local, global, err := n.getController().drvRegistry.IPAMDefaultAddressSpaces(n.ipamType)
	if err != nil ***REMOVED***
		return "", types.NotFoundErrorf("failed to get default address space: %v", err)
	***REMOVED***
	if n.DataScope() == datastore.GlobalScope ***REMOVED***
		return global, nil
	***REMOVED***
	return local, nil
***REMOVED***

func (n *network) Info() NetworkInfo ***REMOVED***
	return n
***REMOVED***

func (n *network) Peers() []networkdb.PeerInfo ***REMOVED***
	if !n.Dynamic() ***REMOVED***
		return []networkdb.PeerInfo***REMOVED******REMOVED***
	***REMOVED***

	agent := n.getController().getAgent()
	if agent == nil ***REMOVED***
		return []networkdb.PeerInfo***REMOVED******REMOVED***
	***REMOVED***

	return agent.networkDB.Peers(n.ID())
***REMOVED***

func (n *network) DriverOptions() map[string]string ***REMOVED***
	n.Lock()
	defer n.Unlock()
	if n.generic != nil ***REMOVED***
		if m, ok := n.generic[netlabel.GenericData]; ok ***REMOVED***
			return m.(map[string]string)
		***REMOVED***
	***REMOVED***
	return map[string]string***REMOVED******REMOVED***
***REMOVED***

func (n *network) Scope() string ***REMOVED***
	n.Lock()
	defer n.Unlock()
	return n.scope
***REMOVED***

func (n *network) IpamConfig() (string, map[string]string, []*IpamConf, []*IpamConf) ***REMOVED***
	n.Lock()
	defer n.Unlock()

	v4L := make([]*IpamConf, len(n.ipamV4Config))
	v6L := make([]*IpamConf, len(n.ipamV6Config))

	for i, c := range n.ipamV4Config ***REMOVED***
		cc := &IpamConf***REMOVED******REMOVED***
		c.CopyTo(cc)
		v4L[i] = cc
	***REMOVED***

	for i, c := range n.ipamV6Config ***REMOVED***
		cc := &IpamConf***REMOVED******REMOVED***
		c.CopyTo(cc)
		v6L[i] = cc
	***REMOVED***

	return n.ipamType, n.ipamOptions, v4L, v6L
***REMOVED***

func (n *network) IpamInfo() ([]*IpamInfo, []*IpamInfo) ***REMOVED***
	n.Lock()
	defer n.Unlock()

	v4Info := make([]*IpamInfo, len(n.ipamV4Info))
	v6Info := make([]*IpamInfo, len(n.ipamV6Info))

	for i, info := range n.ipamV4Info ***REMOVED***
		ic := &IpamInfo***REMOVED******REMOVED***
		info.CopyTo(ic)
		v4Info[i] = ic
	***REMOVED***

	for i, info := range n.ipamV6Info ***REMOVED***
		ic := &IpamInfo***REMOVED******REMOVED***
		info.CopyTo(ic)
		v6Info[i] = ic
	***REMOVED***

	return v4Info, v6Info
***REMOVED***

func (n *network) Internal() bool ***REMOVED***
	n.Lock()
	defer n.Unlock()

	return n.internal
***REMOVED***

func (n *network) Attachable() bool ***REMOVED***
	n.Lock()
	defer n.Unlock()

	return n.attachable
***REMOVED***

func (n *network) Ingress() bool ***REMOVED***
	n.Lock()
	defer n.Unlock()

	return n.ingress
***REMOVED***

func (n *network) Dynamic() bool ***REMOVED***
	n.Lock()
	defer n.Unlock()

	return n.dynamic
***REMOVED***

func (n *network) IPv6Enabled() bool ***REMOVED***
	n.Lock()
	defer n.Unlock()

	return n.enableIPv6
***REMOVED***

func (n *network) ConfigFrom() string ***REMOVED***
	n.Lock()
	defer n.Unlock()

	return n.configFrom
***REMOVED***

func (n *network) ConfigOnly() bool ***REMOVED***
	n.Lock()
	defer n.Unlock()

	return n.configOnly
***REMOVED***

func (n *network) Labels() map[string]string ***REMOVED***
	n.Lock()
	defer n.Unlock()

	var lbls = make(map[string]string, len(n.labels))
	for k, v := range n.labels ***REMOVED***
		lbls[k] = v
	***REMOVED***

	return lbls
***REMOVED***

func (n *network) TableEventRegister(tableName string, objType driverapi.ObjectType) error ***REMOVED***
	if !driverapi.IsValidType(objType) ***REMOVED***
		return fmt.Errorf("invalid object type %v in registering table, %s", objType, tableName)
	***REMOVED***

	t := networkDBTable***REMOVED***
		name:    tableName,
		objType: objType,
	***REMOVED***
	n.Lock()
	defer n.Unlock()
	n.driverTables = append(n.driverTables, t)
	return nil
***REMOVED***

// Special drivers are ones which do not need to perform any network plumbing
func (n *network) hasSpecialDriver() bool ***REMOVED***
	return n.Type() == "host" || n.Type() == "null"
***REMOVED***

func (n *network) ResolveName(req string, ipType int) ([]net.IP, bool) ***REMOVED***
	var ipv6Miss bool

	c := n.getController()
	c.Lock()
	defer c.Unlock()
	sr, ok := c.svcRecords[n.ID()]

	if !ok ***REMOVED***
		return nil, false
	***REMOVED***

	req = strings.TrimSuffix(req, ".")
	ipSet, ok := sr.svcMap.Get(req)

	if ipType == types.IPv6 ***REMOVED***
		// If the name resolved to v4 address then its a valid name in
		// the docker network domain. If the network is not v6 enabled
		// set ipv6Miss to filter the DNS query from going to external
		// resolvers.
		if ok && !n.enableIPv6 ***REMOVED***
			ipv6Miss = true
		***REMOVED***
		ipSet, ok = sr.svcIPv6Map.Get(req)
	***REMOVED***

	if ok && len(ipSet) > 0 ***REMOVED***
		// this map is to avoid IP duplicates, this can happen during a transition period where 2 services are using the same IP
		noDup := make(map[string]bool)
		var ipLocal []net.IP
		for _, ip := range ipSet ***REMOVED***
			if _, dup := noDup[ip.(svcMapEntry).ip]; !dup ***REMOVED***
				noDup[ip.(svcMapEntry).ip] = true
				ipLocal = append(ipLocal, net.ParseIP(ip.(svcMapEntry).ip))
			***REMOVED***
		***REMOVED***
		return ipLocal, ok
	***REMOVED***

	return nil, ipv6Miss
***REMOVED***

func (n *network) HandleQueryResp(name string, ip net.IP) ***REMOVED***
	c := n.getController()
	c.Lock()
	defer c.Unlock()
	sr, ok := c.svcRecords[n.ID()]

	if !ok ***REMOVED***
		return
	***REMOVED***

	ipStr := netutils.ReverseIP(ip.String())
	// If an object with extResolver == true is already in the set this call will fail
	// but anyway it means that has already been inserted before
	if ok, _ := sr.ipMap.Contains(ipStr, ipInfo***REMOVED***name: name***REMOVED***); ok ***REMOVED***
		sr.ipMap.Remove(ipStr, ipInfo***REMOVED***name: name***REMOVED***)
		sr.ipMap.Insert(ipStr, ipInfo***REMOVED***name: name, extResolver: true***REMOVED***)
	***REMOVED***
***REMOVED***

func (n *network) ResolveIP(ip string) string ***REMOVED***
	c := n.getController()
	c.Lock()
	defer c.Unlock()
	sr, ok := c.svcRecords[n.ID()]

	if !ok ***REMOVED***
		return ""
	***REMOVED***

	nwName := n.Name()

	elemSet, ok := sr.ipMap.Get(ip)
	if !ok || len(elemSet) == 0 ***REMOVED***
		return ""
	***REMOVED***
	// NOTE it is possible to have more than one element in the Set, this will happen
	// because of interleave of different events from different sources (local container create vs
	// network db notifications)
	// In such cases the resolution will be based on the first element of the set, and can vary
	// during the system stabilitation
	elem, ok := elemSet[0].(ipInfo)
	if !ok ***REMOVED***
		setStr, b := sr.ipMap.String(ip)
		logrus.Errorf("expected set of ipInfo type for key %s set:%t %s", ip, b, setStr)
		return ""
	***REMOVED***

	if elem.extResolver ***REMOVED***
		return ""
	***REMOVED***

	return elem.name + "." + nwName
***REMOVED***

func (n *network) ResolveService(name string) ([]*net.SRV, []net.IP) ***REMOVED***
	c := n.getController()

	srv := []*net.SRV***REMOVED******REMOVED***
	ip := []net.IP***REMOVED******REMOVED***

	logrus.Debugf("Service name To resolve: %v", name)

	// There are DNS implementaions that allow SRV queries for names not in
	// the format defined by RFC 2782. Hence specific validations checks are
	// not done
	parts := strings.Split(name, ".")
	if len(parts) < 3 ***REMOVED***
		return nil, nil
	***REMOVED***

	portName := parts[0]
	proto := parts[1]
	svcName := strings.Join(parts[2:], ".")

	c.Lock()
	defer c.Unlock()
	sr, ok := c.svcRecords[n.ID()]

	if !ok ***REMOVED***
		return nil, nil
	***REMOVED***

	svcs, ok := sr.service[svcName]
	if !ok ***REMOVED***
		return nil, nil
	***REMOVED***

	for _, svc := range svcs ***REMOVED***
		if svc.portName != portName ***REMOVED***
			continue
		***REMOVED***
		if svc.proto != proto ***REMOVED***
			continue
		***REMOVED***
		for _, t := range svc.target ***REMOVED***
			srv = append(srv,
				&net.SRV***REMOVED***
					Target: t.name,
					Port:   t.port,
				***REMOVED***)

			ip = append(ip, t.ip)
		***REMOVED***
	***REMOVED***

	return srv, ip
***REMOVED***

func (n *network) ExecFunc(f func()) error ***REMOVED***
	return types.NotImplementedErrorf("ExecFunc not supported by network")
***REMOVED***

func (n *network) NdotsSet() bool ***REMOVED***
	return false
***REMOVED***

// config-only network is looked up by name
func (c *controller) getConfigNetwork(name string) (*network, error) ***REMOVED***
	var n Network

	s := func(current Network) bool ***REMOVED***
		if current.Info().ConfigOnly() && current.Name() == name ***REMOVED***
			n = current
			return true
		***REMOVED***
		return false
	***REMOVED***

	c.WalkNetworks(s)

	if n == nil ***REMOVED***
		return nil, types.NotFoundErrorf("configuration network %q not found", name)
	***REMOVED***

	return n.(*network), nil
***REMOVED***

func (n *network) createLoadBalancerSandbox() error ***REMOVED***
	sandboxName := n.name + "-sbox"
	sbOptions := []SandboxOption***REMOVED******REMOVED***
	if n.ingress ***REMOVED***
		sbOptions = append(sbOptions, OptionIngress())
	***REMOVED***
	sb, err := n.ctrlr.NewSandbox(sandboxName, sbOptions...)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			if e := n.ctrlr.SandboxDestroy(sandboxName); e != nil ***REMOVED***
				logrus.Warnf("could not delete sandbox %s on failure on failure (%v): %v", sandboxName, err, e)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	endpointName := n.name + "-endpoint"
	epOptions := []EndpointOption***REMOVED***
		CreateOptionIpam(n.loadBalancerIP, nil, nil, nil),
		CreateOptionLoadBalancer(),
	***REMOVED***
	ep, err := n.createEndpoint(endpointName, epOptions...)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			if e := ep.Delete(true); e != nil ***REMOVED***
				logrus.Warnf("could not delete endpoint %s on failure on failure (%v): %v", endpointName, err, e)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	if err := ep.Join(sb, nil); err != nil ***REMOVED***
		return err
	***REMOVED***
	return sb.EnableService()
***REMOVED***

func (n *network) deleteLoadBalancerSandbox() ***REMOVED***
	n.Lock()
	c := n.ctrlr
	name := n.name
	n.Unlock()

	endpointName := name + "-endpoint"
	sandboxName := name + "-sbox"

	endpoint, err := n.EndpointByName(endpointName)
	if err != nil ***REMOVED***
		logrus.Warnf("Failed to find load balancer endpoint %s on network %s: %v", endpointName, name, err)
	***REMOVED*** else ***REMOVED***

		info := endpoint.Info()
		if info != nil ***REMOVED***
			sb := info.Sandbox()
			if sb != nil ***REMOVED***
				if err := sb.DisableService(); err != nil ***REMOVED***
					logrus.Warnf("Failed to disable service on sandbox %s: %v", sandboxName, err)
					//Ignore error and attempt to delete the load balancer endpoint
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if err := endpoint.Delete(true); err != nil ***REMOVED***
			logrus.Warnf("Failed to delete endpoint %s (%s) in %s: %v", endpoint.Name(), endpoint.ID(), sandboxName, err)
			//Ignore error and attempt to delete the sandbox.
		***REMOVED***
	***REMOVED***

	if err := c.SandboxDestroy(sandboxName); err != nil ***REMOVED***
		logrus.Warnf("Failed to delete %s sandbox: %v", sandboxName, err)
	***REMOVED***
***REMOVED***
