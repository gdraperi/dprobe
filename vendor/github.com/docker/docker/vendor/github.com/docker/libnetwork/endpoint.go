package libnetwork

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/ipamapi"
	"github.com/docker/libnetwork/netlabel"
	"github.com/docker/libnetwork/options"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

// Endpoint represents a logical connection between a network and a sandbox.
type Endpoint interface ***REMOVED***
	// A system generated id for this endpoint.
	ID() string

	// Name returns the name of this endpoint.
	Name() string

	// Network returns the name of the network to which this endpoint is attached.
	Network() string

	// Join joins the sandbox to the endpoint and populates into the sandbox
	// the network resources allocated for the endpoint.
	Join(sandbox Sandbox, options ...EndpointOption) error

	// Leave detaches the network resources populated in the sandbox.
	Leave(sandbox Sandbox, options ...EndpointOption) error

	// Return certain operational data belonging to this endpoint
	Info() EndpointInfo

	// DriverInfo returns a collection of driver operational data related to this endpoint retrieved from the driver
	DriverInfo() (map[string]interface***REMOVED******REMOVED***, error)

	// Delete and detaches this endpoint from the network.
	Delete(force bool) error
***REMOVED***

// EndpointOption is an option setter function type used to pass various options to Network
// and Endpoint interfaces methods. The various setter functions of type EndpointOption are
// provided by libnetwork, they look like <Create|Join|Leave>Option[...](...)
type EndpointOption func(ep *endpoint)

type endpoint struct ***REMOVED***
	name              string
	id                string
	network           *network
	iface             *endpointInterface
	joinInfo          *endpointJoinInfo
	sandboxID         string
	locator           string
	exposedPorts      []types.TransportPort
	anonymous         bool
	disableResolution bool
	generic           map[string]interface***REMOVED******REMOVED***
	joinLeaveDone     chan struct***REMOVED******REMOVED***
	prefAddress       net.IP
	prefAddressV6     net.IP
	ipamOptions       map[string]string
	aliases           map[string]string
	myAliases         []string
	svcID             string
	svcName           string
	virtualIP         net.IP
	svcAliases        []string
	ingressPorts      []*PortConfig
	dbIndex           uint64
	dbExists          bool
	serviceEnabled    bool
	loadBalancer      bool
	sync.Mutex
***REMOVED***

func (ep *endpoint) MarshalJSON() ([]byte, error) ***REMOVED***
	ep.Lock()
	defer ep.Unlock()

	epMap := make(map[string]interface***REMOVED******REMOVED***)
	epMap["name"] = ep.name
	epMap["id"] = ep.id
	epMap["ep_iface"] = ep.iface
	epMap["joinInfo"] = ep.joinInfo
	epMap["exposed_ports"] = ep.exposedPorts
	if ep.generic != nil ***REMOVED***
		epMap["generic"] = ep.generic
	***REMOVED***
	epMap["sandbox"] = ep.sandboxID
	epMap["locator"] = ep.locator
	epMap["anonymous"] = ep.anonymous
	epMap["disableResolution"] = ep.disableResolution
	epMap["myAliases"] = ep.myAliases
	epMap["svcName"] = ep.svcName
	epMap["svcID"] = ep.svcID
	epMap["virtualIP"] = ep.virtualIP.String()
	epMap["ingressPorts"] = ep.ingressPorts
	epMap["svcAliases"] = ep.svcAliases
	epMap["loadBalancer"] = ep.loadBalancer

	return json.Marshal(epMap)
***REMOVED***

func (ep *endpoint) UnmarshalJSON(b []byte) (err error) ***REMOVED***
	ep.Lock()
	defer ep.Unlock()

	var epMap map[string]interface***REMOVED******REMOVED***
	if err := json.Unmarshal(b, &epMap); err != nil ***REMOVED***
		return err
	***REMOVED***
	ep.name = epMap["name"].(string)
	ep.id = epMap["id"].(string)

	ib, _ := json.Marshal(epMap["ep_iface"])
	json.Unmarshal(ib, &ep.iface)

	jb, _ := json.Marshal(epMap["joinInfo"])
	json.Unmarshal(jb, &ep.joinInfo)

	tb, _ := json.Marshal(epMap["exposed_ports"])
	var tPorts []types.TransportPort
	json.Unmarshal(tb, &tPorts)
	ep.exposedPorts = tPorts

	cb, _ := json.Marshal(epMap["sandbox"])
	json.Unmarshal(cb, &ep.sandboxID)

	if v, ok := epMap["generic"]; ok ***REMOVED***
		ep.generic = v.(map[string]interface***REMOVED******REMOVED***)

		if opt, ok := ep.generic[netlabel.PortMap]; ok ***REMOVED***
			pblist := []types.PortBinding***REMOVED******REMOVED***

			for i := 0; i < len(opt.([]interface***REMOVED******REMOVED***)); i++ ***REMOVED***
				pb := types.PortBinding***REMOVED******REMOVED***
				tmp := opt.([]interface***REMOVED******REMOVED***)[i].(map[string]interface***REMOVED******REMOVED***)

				bytes, err := json.Marshal(tmp)
				if err != nil ***REMOVED***
					logrus.Error(err)
					break
				***REMOVED***
				err = json.Unmarshal(bytes, &pb)
				if err != nil ***REMOVED***
					logrus.Error(err)
					break
				***REMOVED***
				pblist = append(pblist, pb)
			***REMOVED***
			ep.generic[netlabel.PortMap] = pblist
		***REMOVED***

		if opt, ok := ep.generic[netlabel.ExposedPorts]; ok ***REMOVED***
			tplist := []types.TransportPort***REMOVED******REMOVED***

			for i := 0; i < len(opt.([]interface***REMOVED******REMOVED***)); i++ ***REMOVED***
				tp := types.TransportPort***REMOVED******REMOVED***
				tmp := opt.([]interface***REMOVED******REMOVED***)[i].(map[string]interface***REMOVED******REMOVED***)

				bytes, err := json.Marshal(tmp)
				if err != nil ***REMOVED***
					logrus.Error(err)
					break
				***REMOVED***
				err = json.Unmarshal(bytes, &tp)
				if err != nil ***REMOVED***
					logrus.Error(err)
					break
				***REMOVED***
				tplist = append(tplist, tp)
			***REMOVED***
			ep.generic[netlabel.ExposedPorts] = tplist

		***REMOVED***
	***REMOVED***

	if v, ok := epMap["anonymous"]; ok ***REMOVED***
		ep.anonymous = v.(bool)
	***REMOVED***
	if v, ok := epMap["disableResolution"]; ok ***REMOVED***
		ep.disableResolution = v.(bool)
	***REMOVED***
	if l, ok := epMap["locator"]; ok ***REMOVED***
		ep.locator = l.(string)
	***REMOVED***

	if sn, ok := epMap["svcName"]; ok ***REMOVED***
		ep.svcName = sn.(string)
	***REMOVED***

	if si, ok := epMap["svcID"]; ok ***REMOVED***
		ep.svcID = si.(string)
	***REMOVED***

	if vip, ok := epMap["virtualIP"]; ok ***REMOVED***
		ep.virtualIP = net.ParseIP(vip.(string))
	***REMOVED***

	if v, ok := epMap["loadBalancer"]; ok ***REMOVED***
		ep.loadBalancer = v.(bool)
	***REMOVED***

	sal, _ := json.Marshal(epMap["svcAliases"])
	var svcAliases []string
	json.Unmarshal(sal, &svcAliases)
	ep.svcAliases = svcAliases

	pc, _ := json.Marshal(epMap["ingressPorts"])
	var ingressPorts []*PortConfig
	json.Unmarshal(pc, &ingressPorts)
	ep.ingressPorts = ingressPorts

	ma, _ := json.Marshal(epMap["myAliases"])
	var myAliases []string
	json.Unmarshal(ma, &myAliases)
	ep.myAliases = myAliases
	return nil
***REMOVED***

func (ep *endpoint) New() datastore.KVObject ***REMOVED***
	return &endpoint***REMOVED***network: ep.getNetwork()***REMOVED***
***REMOVED***

func (ep *endpoint) CopyTo(o datastore.KVObject) error ***REMOVED***
	ep.Lock()
	defer ep.Unlock()

	dstEp := o.(*endpoint)
	dstEp.name = ep.name
	dstEp.id = ep.id
	dstEp.sandboxID = ep.sandboxID
	dstEp.locator = ep.locator
	dstEp.dbIndex = ep.dbIndex
	dstEp.dbExists = ep.dbExists
	dstEp.anonymous = ep.anonymous
	dstEp.disableResolution = ep.disableResolution
	dstEp.svcName = ep.svcName
	dstEp.svcID = ep.svcID
	dstEp.virtualIP = ep.virtualIP
	dstEp.loadBalancer = ep.loadBalancer

	dstEp.svcAliases = make([]string, len(ep.svcAliases))
	copy(dstEp.svcAliases, ep.svcAliases)

	dstEp.ingressPorts = make([]*PortConfig, len(ep.ingressPorts))
	copy(dstEp.ingressPorts, ep.ingressPorts)

	if ep.iface != nil ***REMOVED***
		dstEp.iface = &endpointInterface***REMOVED******REMOVED***
		ep.iface.CopyTo(dstEp.iface)
	***REMOVED***

	if ep.joinInfo != nil ***REMOVED***
		dstEp.joinInfo = &endpointJoinInfo***REMOVED******REMOVED***
		ep.joinInfo.CopyTo(dstEp.joinInfo)
	***REMOVED***

	dstEp.exposedPorts = make([]types.TransportPort, len(ep.exposedPorts))
	copy(dstEp.exposedPorts, ep.exposedPorts)

	dstEp.myAliases = make([]string, len(ep.myAliases))
	copy(dstEp.myAliases, ep.myAliases)

	dstEp.generic = options.Generic***REMOVED******REMOVED***
	for k, v := range ep.generic ***REMOVED***
		dstEp.generic[k] = v
	***REMOVED***

	return nil
***REMOVED***

func (ep *endpoint) ID() string ***REMOVED***
	ep.Lock()
	defer ep.Unlock()

	return ep.id
***REMOVED***

func (ep *endpoint) Name() string ***REMOVED***
	ep.Lock()
	defer ep.Unlock()

	return ep.name
***REMOVED***

func (ep *endpoint) MyAliases() []string ***REMOVED***
	ep.Lock()
	defer ep.Unlock()

	return ep.myAliases
***REMOVED***

func (ep *endpoint) Network() string ***REMOVED***
	if ep.network == nil ***REMOVED***
		return ""
	***REMOVED***

	return ep.network.name
***REMOVED***

func (ep *endpoint) isAnonymous() bool ***REMOVED***
	ep.Lock()
	defer ep.Unlock()
	return ep.anonymous
***REMOVED***

// isServiceEnabled check if service is enabled on the endpoint
func (ep *endpoint) isServiceEnabled() bool ***REMOVED***
	ep.Lock()
	defer ep.Unlock()
	return ep.serviceEnabled
***REMOVED***

// enableService sets service enabled on the endpoint
func (ep *endpoint) enableService() ***REMOVED***
	ep.Lock()
	defer ep.Unlock()
	ep.serviceEnabled = true
***REMOVED***

// disableService disables service on the endpoint
func (ep *endpoint) disableService() ***REMOVED***
	ep.Lock()
	defer ep.Unlock()
	ep.serviceEnabled = false
***REMOVED***

func (ep *endpoint) needResolver() bool ***REMOVED***
	ep.Lock()
	defer ep.Unlock()
	return !ep.disableResolution
***REMOVED***

// endpoint Key structure : endpoint/network-id/endpoint-id
func (ep *endpoint) Key() []string ***REMOVED***
	if ep.network == nil ***REMOVED***
		return nil
	***REMOVED***

	return []string***REMOVED***datastore.EndpointKeyPrefix, ep.network.id, ep.id***REMOVED***
***REMOVED***

func (ep *endpoint) KeyPrefix() []string ***REMOVED***
	if ep.network == nil ***REMOVED***
		return nil
	***REMOVED***

	return []string***REMOVED***datastore.EndpointKeyPrefix, ep.network.id***REMOVED***
***REMOVED***

func (ep *endpoint) networkIDFromKey(key string) (string, error) ***REMOVED***
	// endpoint Key structure : docker/libnetwork/endpoint/$***REMOVED***network-id***REMOVED***/$***REMOVED***endpoint-id***REMOVED***
	// it's an invalid key if the key doesn't have all the 5 key elements above
	keyElements := strings.Split(key, "/")
	if !strings.HasPrefix(key, datastore.Key(datastore.EndpointKeyPrefix)) || len(keyElements) < 5 ***REMOVED***
		return "", fmt.Errorf("invalid endpoint key : %v", key)
	***REMOVED***
	// network-id is placed at index=3. pls refer to endpoint.Key() method
	return strings.Split(key, "/")[3], nil
***REMOVED***

func (ep *endpoint) Value() []byte ***REMOVED***
	b, err := json.Marshal(ep)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	return b
***REMOVED***

func (ep *endpoint) SetValue(value []byte) error ***REMOVED***
	return json.Unmarshal(value, ep)
***REMOVED***

func (ep *endpoint) Index() uint64 ***REMOVED***
	ep.Lock()
	defer ep.Unlock()
	return ep.dbIndex
***REMOVED***

func (ep *endpoint) SetIndex(index uint64) ***REMOVED***
	ep.Lock()
	defer ep.Unlock()
	ep.dbIndex = index
	ep.dbExists = true
***REMOVED***

func (ep *endpoint) Exists() bool ***REMOVED***
	ep.Lock()
	defer ep.Unlock()
	return ep.dbExists
***REMOVED***

func (ep *endpoint) Skip() bool ***REMOVED***
	return ep.getNetwork().Skip()
***REMOVED***

func (ep *endpoint) processOptions(options ...EndpointOption) ***REMOVED***
	ep.Lock()
	defer ep.Unlock()

	for _, opt := range options ***REMOVED***
		if opt != nil ***REMOVED***
			opt(ep)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (ep *endpoint) getNetwork() *network ***REMOVED***
	ep.Lock()
	defer ep.Unlock()

	return ep.network
***REMOVED***

func (ep *endpoint) getNetworkFromStore() (*network, error) ***REMOVED***
	if ep.network == nil ***REMOVED***
		return nil, fmt.Errorf("invalid network object in endpoint %s", ep.Name())
	***REMOVED***

	return ep.network.getController().getNetworkFromStore(ep.network.id)
***REMOVED***

func (ep *endpoint) Join(sbox Sandbox, options ...EndpointOption) error ***REMOVED***
	if sbox == nil ***REMOVED***
		return types.BadRequestErrorf("endpoint cannot be joined by nil container")
	***REMOVED***

	sb, ok := sbox.(*sandbox)
	if !ok ***REMOVED***
		return types.BadRequestErrorf("not a valid Sandbox interface")
	***REMOVED***

	sb.joinLeaveStart()
	defer sb.joinLeaveEnd()

	return ep.sbJoin(sb, options...)
***REMOVED***

func (ep *endpoint) sbJoin(sb *sandbox, options ...EndpointOption) (err error) ***REMOVED***
	n, err := ep.getNetworkFromStore()
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to get network from store during join: %v", err)
	***REMOVED***

	ep, err = n.getEndpointFromStore(ep.ID())
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to get endpoint from store during join: %v", err)
	***REMOVED***

	ep.Lock()
	if ep.sandboxID != "" ***REMOVED***
		ep.Unlock()
		return types.ForbiddenErrorf("another container is attached to the same network endpoint")
	***REMOVED***
	ep.network = n
	ep.sandboxID = sb.ID()
	ep.joinInfo = &endpointJoinInfo***REMOVED******REMOVED***
	epid := ep.id
	ep.Unlock()
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			ep.Lock()
			ep.sandboxID = ""
			ep.Unlock()
		***REMOVED***
	***REMOVED***()

	nid := n.ID()

	ep.processOptions(options...)

	d, err := n.driver(true)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to get driver during join: %v", err)
	***REMOVED***

	err = d.Join(nid, epid, sb.Key(), ep, sb.Labels())
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			if e := d.Leave(nid, epid); e != nil ***REMOVED***
				logrus.Warnf("driver leave failed while rolling back join: %v", e)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	// Watch for service records
	if !n.getController().isAgent() ***REMOVED***
		n.getController().watchSvcRecord(ep)
	***REMOVED***

	if doUpdateHostsFile(n, sb) ***REMOVED***
		address := ""
		if ip := ep.getFirstInterfaceAddress(); ip != nil ***REMOVED***
			address = ip.String()
		***REMOVED***
		if err = sb.updateHostsFile(address); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if err = sb.updateDNS(n.enableIPv6); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Current endpoint providing external connectivity for the sandbox
	extEp := sb.getGatewayEndpoint()

	sb.Lock()
	heap.Push(&sb.endpoints, ep)
	sb.Unlock()
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			sb.removeEndpoint(ep)
		***REMOVED***
	***REMOVED***()

	if err = sb.populateNetworkResources(ep); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err = n.getController().updateToStore(ep); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err = ep.addDriverInfoToCluster(); err != nil ***REMOVED***
		return err
	***REMOVED***

	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			if e := ep.deleteDriverInfoFromCluster(); e != nil ***REMOVED***
				logrus.Errorf("Could not delete endpoint state for endpoint %s from cluster on join failure: %v", ep.Name(), e)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	if sb.needDefaultGW() && sb.getEndpointInGWNetwork() == nil ***REMOVED***
		return sb.setupDefaultGW()
	***REMOVED***

	moveExtConn := sb.getGatewayEndpoint() != extEp

	if moveExtConn ***REMOVED***
		if extEp != nil ***REMOVED***
			logrus.Debugf("Revoking external connectivity on endpoint %s (%s)", extEp.Name(), extEp.ID())
			extN, err := extEp.getNetworkFromStore()
			if err != nil ***REMOVED***
				return fmt.Errorf("failed to get network from store for revoking external connectivity during join: %v", err)
			***REMOVED***
			extD, err := extN.driver(true)
			if err != nil ***REMOVED***
				return fmt.Errorf("failed to get driver for revoking external connectivity during join: %v", err)
			***REMOVED***
			if err = extD.RevokeExternalConnectivity(extEp.network.ID(), extEp.ID()); err != nil ***REMOVED***
				return types.InternalErrorf(
					"driver failed revoking external connectivity on endpoint %s (%s): %v",
					extEp.Name(), extEp.ID(), err)
			***REMOVED***
			defer func() ***REMOVED***
				if err != nil ***REMOVED***
					if e := extD.ProgramExternalConnectivity(extEp.network.ID(), extEp.ID(), sb.Labels()); e != nil ***REMOVED***
						logrus.Warnf("Failed to roll-back external connectivity on endpoint %s (%s): %v",
							extEp.Name(), extEp.ID(), e)
					***REMOVED***
				***REMOVED***
			***REMOVED***()
		***REMOVED***
		if !n.internal ***REMOVED***
			logrus.Debugf("Programming external connectivity on endpoint %s (%s)", ep.Name(), ep.ID())
			if err = d.ProgramExternalConnectivity(n.ID(), ep.ID(), sb.Labels()); err != nil ***REMOVED***
				return types.InternalErrorf(
					"driver failed programming external connectivity on endpoint %s (%s): %v",
					ep.Name(), ep.ID(), err)
			***REMOVED***
		***REMOVED***

	***REMOVED***

	if !sb.needDefaultGW() ***REMOVED***
		if e := sb.clearDefaultGW(); e != nil ***REMOVED***
			logrus.Warnf("Failure while disconnecting sandbox %s (%s) from gateway network: %v",
				sb.ID(), sb.ContainerID(), e)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func doUpdateHostsFile(n *network, sb *sandbox) bool ***REMOVED***
	return !n.ingress && n.Name() != libnGWNetwork
***REMOVED***

func (ep *endpoint) rename(name string) error ***REMOVED***
	var (
		err      error
		netWatch *netWatch
		ok       bool
	)

	n := ep.getNetwork()
	if n == nil ***REMOVED***
		return fmt.Errorf("network not connected for ep %q", ep.name)
	***REMOVED***

	c := n.getController()

	sb, ok := ep.getSandbox()
	if !ok ***REMOVED***
		logrus.Warnf("rename for %s aborted, sandbox %s is not anymore present", ep.ID(), ep.sandboxID)
		return nil
	***REMOVED***

	if c.isAgent() ***REMOVED***
		if err = ep.deleteServiceInfoFromCluster(sb, "rename"); err != nil ***REMOVED***
			return types.InternalErrorf("Could not delete service state for endpoint %s from cluster on rename: %v", ep.Name(), err)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		c.Lock()
		netWatch, ok = c.nmap[n.ID()]
		c.Unlock()
		if !ok ***REMOVED***
			return fmt.Errorf("watch null for network %q", n.Name())
		***REMOVED***
		n.updateSvcRecord(ep, c.getLocalEps(netWatch), false)
	***REMOVED***

	oldName := ep.name
	oldAnonymous := ep.anonymous
	ep.name = name
	ep.anonymous = false

	if c.isAgent() ***REMOVED***
		if err = ep.addServiceInfoToCluster(sb); err != nil ***REMOVED***
			return types.InternalErrorf("Could not add service state for endpoint %s to cluster on rename: %v", ep.Name(), err)
		***REMOVED***
		defer func() ***REMOVED***
			if err != nil ***REMOVED***
				ep.deleteServiceInfoFromCluster(sb, "rename")
				ep.name = oldName
				ep.anonymous = oldAnonymous
				ep.addServiceInfoToCluster(sb)
			***REMOVED***
		***REMOVED***()
	***REMOVED*** else ***REMOVED***
		n.updateSvcRecord(ep, c.getLocalEps(netWatch), true)
		defer func() ***REMOVED***
			if err != nil ***REMOVED***
				n.updateSvcRecord(ep, c.getLocalEps(netWatch), false)
				ep.name = oldName
				ep.anonymous = oldAnonymous
				n.updateSvcRecord(ep, c.getLocalEps(netWatch), true)
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	// Update the store with the updated name
	if err = c.updateToStore(ep); err != nil ***REMOVED***
		return err
	***REMOVED***
	// After the name change do a dummy endpoint count update to
	// trigger the service record update in the peer nodes

	// Ignore the error because updateStore fail for EpCnt is a
	// benign error. Besides there is no meaningful recovery that
	// we can do. When the cluster recovers subsequent EpCnt update
	// will force the peers to get the correct EP name.
	n.getEpCnt().updateStore()

	return err
***REMOVED***

func (ep *endpoint) hasInterface(iName string) bool ***REMOVED***
	ep.Lock()
	defer ep.Unlock()

	return ep.iface != nil && ep.iface.srcName == iName
***REMOVED***

func (ep *endpoint) Leave(sbox Sandbox, options ...EndpointOption) error ***REMOVED***
	if sbox == nil || sbox.ID() == "" || sbox.Key() == "" ***REMOVED***
		return types.BadRequestErrorf("invalid Sandbox passed to endpoint leave: %v", sbox)
	***REMOVED***

	sb, ok := sbox.(*sandbox)
	if !ok ***REMOVED***
		return types.BadRequestErrorf("not a valid Sandbox interface")
	***REMOVED***

	sb.joinLeaveStart()
	defer sb.joinLeaveEnd()

	return ep.sbLeave(sb, false, options...)
***REMOVED***

func (ep *endpoint) sbLeave(sb *sandbox, force bool, options ...EndpointOption) error ***REMOVED***
	n, err := ep.getNetworkFromStore()
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to get network from store during leave: %v", err)
	***REMOVED***

	ep, err = n.getEndpointFromStore(ep.ID())
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to get endpoint from store during leave: %v", err)
	***REMOVED***

	ep.Lock()
	sid := ep.sandboxID
	ep.Unlock()

	if sid == "" ***REMOVED***
		return types.ForbiddenErrorf("cannot leave endpoint with no attached sandbox")
	***REMOVED***
	if sid != sb.ID() ***REMOVED***
		return types.ForbiddenErrorf("unexpected sandbox ID in leave request. Expected %s. Got %s", ep.sandboxID, sb.ID())
	***REMOVED***

	ep.processOptions(options...)

	d, err := n.driver(!force)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to get driver during endpoint leave: %v", err)
	***REMOVED***

	ep.Lock()
	ep.sandboxID = ""
	ep.network = n
	ep.Unlock()

	// Current endpoint providing external connectivity to the sandbox
	extEp := sb.getGatewayEndpoint()
	moveExtConn := extEp != nil && (extEp.ID() == ep.ID())

	if d != nil ***REMOVED***
		if moveExtConn ***REMOVED***
			logrus.Debugf("Revoking external connectivity on endpoint %s (%s)", ep.Name(), ep.ID())
			if err := d.RevokeExternalConnectivity(n.id, ep.id); err != nil ***REMOVED***
				logrus.Warnf("driver failed revoking external connectivity on endpoint %s (%s): %v",
					ep.Name(), ep.ID(), err)
			***REMOVED***
		***REMOVED***

		if err := d.Leave(n.id, ep.id); err != nil ***REMOVED***
			if _, ok := err.(types.MaskableError); !ok ***REMOVED***
				logrus.Warnf("driver error disconnecting container %s : %v", ep.name, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if err := sb.clearNetworkResources(ep); err != nil ***REMOVED***
		logrus.Warnf("Could not cleanup network resources on container %s disconnect: %v", ep.name, err)
	***REMOVED***

	// Update the store about the sandbox detach only after we
	// have completed sb.clearNetworkresources above to avoid
	// spurious logs when cleaning up the sandbox when the daemon
	// ungracefully exits and restarts before completing sandbox
	// detach but after store has been updated.
	if err := n.getController().updateToStore(ep); err != nil ***REMOVED***
		return err
	***REMOVED***

	if e := ep.deleteDriverInfoFromCluster(); e != nil ***REMOVED***
		logrus.Errorf("Could not delete endpoint state for endpoint %s from cluster: %v", ep.Name(), e)
	***REMOVED***

	sb.deleteHostsEntries(n.getSvcRecords(ep))
	if !sb.inDelete && sb.needDefaultGW() && sb.getEndpointInGWNetwork() == nil ***REMOVED***
		return sb.setupDefaultGW()
	***REMOVED***

	// New endpoint providing external connectivity for the sandbox
	extEp = sb.getGatewayEndpoint()
	if moveExtConn && extEp != nil ***REMOVED***
		logrus.Debugf("Programming external connectivity on endpoint %s (%s)", extEp.Name(), extEp.ID())
		extN, err := extEp.getNetworkFromStore()
		if err != nil ***REMOVED***
			return fmt.Errorf("failed to get network from store for programming external connectivity during leave: %v", err)
		***REMOVED***
		extD, err := extN.driver(true)
		if err != nil ***REMOVED***
			return fmt.Errorf("failed to get driver for programming external connectivity during leave: %v", err)
		***REMOVED***
		if err := extD.ProgramExternalConnectivity(extEp.network.ID(), extEp.ID(), sb.Labels()); err != nil ***REMOVED***
			logrus.Warnf("driver failed programming external connectivity on endpoint %s: (%s) %v",
				extEp.Name(), extEp.ID(), err)
		***REMOVED***
	***REMOVED***

	if !sb.needDefaultGW() ***REMOVED***
		if err := sb.clearDefaultGW(); err != nil ***REMOVED***
			logrus.Warnf("Failure while disconnecting sandbox %s (%s) from gateway network: %v",
				sb.ID(), sb.ContainerID(), err)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (ep *endpoint) Delete(force bool) error ***REMOVED***
	var err error
	n, err := ep.getNetworkFromStore()
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to get network during Delete: %v", err)
	***REMOVED***

	ep, err = n.getEndpointFromStore(ep.ID())
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to get endpoint from store during Delete: %v", err)
	***REMOVED***

	ep.Lock()
	epid := ep.id
	name := ep.name
	sbid := ep.sandboxID
	ep.Unlock()

	sb, _ := n.getController().SandboxByID(sbid)
	if sb != nil && !force ***REMOVED***
		return &ActiveContainerError***REMOVED***name: name, id: epid***REMOVED***
	***REMOVED***

	if sb != nil ***REMOVED***
		if e := ep.sbLeave(sb.(*sandbox), force); e != nil ***REMOVED***
			logrus.Warnf("failed to leave sandbox for endpoint %s : %v", name, e)
		***REMOVED***
	***REMOVED***

	if err = n.getController().deleteFromStore(ep); err != nil ***REMOVED***
		return err
	***REMOVED***

	defer func() ***REMOVED***
		if err != nil && !force ***REMOVED***
			ep.dbExists = false
			if e := n.getController().updateToStore(ep); e != nil ***REMOVED***
				logrus.Warnf("failed to recreate endpoint in store %s : %v", name, e)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	// unwatch for service records
	n.getController().unWatchSvcRecord(ep)

	if err = ep.deleteEndpoint(force); err != nil && !force ***REMOVED***
		return err
	***REMOVED***

	ep.releaseAddress()

	if err := n.getEpCnt().DecEndpointCnt(); err != nil ***REMOVED***
		logrus.Warnf("failed to decrement endpoint count for ep %s: %v", ep.ID(), err)
	***REMOVED***

	return nil
***REMOVED***

func (ep *endpoint) deleteEndpoint(force bool) error ***REMOVED***
	ep.Lock()
	n := ep.network
	name := ep.name
	epid := ep.id
	ep.Unlock()

	driver, err := n.driver(!force)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to delete endpoint: %v", err)
	***REMOVED***

	if driver == nil ***REMOVED***
		return nil
	***REMOVED***

	if err := driver.DeleteEndpoint(n.id, epid); err != nil ***REMOVED***
		if _, ok := err.(types.ForbiddenError); ok ***REMOVED***
			return err
		***REMOVED***

		if _, ok := err.(types.MaskableError); !ok ***REMOVED***
			logrus.Warnf("driver error deleting endpoint %s : %v", name, err)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (ep *endpoint) getSandbox() (*sandbox, bool) ***REMOVED***
	c := ep.network.getController()
	ep.Lock()
	sid := ep.sandboxID
	ep.Unlock()

	c.Lock()
	ps, ok := c.sandboxes[sid]
	c.Unlock()

	return ps, ok
***REMOVED***

func (ep *endpoint) getFirstInterfaceAddress() net.IP ***REMOVED***
	ep.Lock()
	defer ep.Unlock()

	if ep.iface.addr != nil ***REMOVED***
		return ep.iface.addr.IP
	***REMOVED***

	return nil
***REMOVED***

// EndpointOptionGeneric function returns an option setter for a Generic option defined
// in a Dictionary of Key-Value pair
func EndpointOptionGeneric(generic map[string]interface***REMOVED******REMOVED***) EndpointOption ***REMOVED***
	return func(ep *endpoint) ***REMOVED***
		for k, v := range generic ***REMOVED***
			ep.generic[k] = v
		***REMOVED***
	***REMOVED***
***REMOVED***

var (
	linkLocalMask     = net.CIDRMask(16, 32)
	linkLocalMaskIPv6 = net.CIDRMask(64, 128)
)

// CreateOptionIpam function returns an option setter for the ipam configuration for this endpoint
func CreateOptionIpam(ipV4, ipV6 net.IP, llIPs []net.IP, ipamOptions map[string]string) EndpointOption ***REMOVED***
	return func(ep *endpoint) ***REMOVED***
		ep.prefAddress = ipV4
		ep.prefAddressV6 = ipV6
		if len(llIPs) != 0 ***REMOVED***
			for _, ip := range llIPs ***REMOVED***
				nw := &net.IPNet***REMOVED***IP: ip, Mask: linkLocalMask***REMOVED***
				if ip.To4() == nil ***REMOVED***
					nw.Mask = linkLocalMaskIPv6
				***REMOVED***
				ep.iface.llAddrs = append(ep.iface.llAddrs, nw)
			***REMOVED***
		***REMOVED***
		ep.ipamOptions = ipamOptions
	***REMOVED***
***REMOVED***

// CreateOptionExposedPorts function returns an option setter for the container exposed
// ports option to be passed to network.CreateEndpoint() method.
func CreateOptionExposedPorts(exposedPorts []types.TransportPort) EndpointOption ***REMOVED***
	return func(ep *endpoint) ***REMOVED***
		// Defensive copy
		eps := make([]types.TransportPort, len(exposedPorts))
		copy(eps, exposedPorts)
		// Store endpoint label and in generic because driver needs it
		ep.exposedPorts = eps
		ep.generic[netlabel.ExposedPorts] = eps
	***REMOVED***
***REMOVED***

// CreateOptionPortMapping function returns an option setter for the mapping
// ports option to be passed to network.CreateEndpoint() method.
func CreateOptionPortMapping(portBindings []types.PortBinding) EndpointOption ***REMOVED***
	return func(ep *endpoint) ***REMOVED***
		// Store a copy of the bindings as generic data to pass to the driver
		pbs := make([]types.PortBinding, len(portBindings))
		copy(pbs, portBindings)
		ep.generic[netlabel.PortMap] = pbs
	***REMOVED***
***REMOVED***

// CreateOptionDNS function returns an option setter for dns entry option to
// be passed to container Create method.
func CreateOptionDNS(dns []string) EndpointOption ***REMOVED***
	return func(ep *endpoint) ***REMOVED***
		ep.generic[netlabel.DNSServers] = dns
	***REMOVED***
***REMOVED***

// CreateOptionAnonymous function returns an option setter for setting
// this endpoint as anonymous
func CreateOptionAnonymous() EndpointOption ***REMOVED***
	return func(ep *endpoint) ***REMOVED***
		ep.anonymous = true
	***REMOVED***
***REMOVED***

// CreateOptionDisableResolution function returns an option setter to indicate
// this endpoint doesn't want embedded DNS server functionality
func CreateOptionDisableResolution() EndpointOption ***REMOVED***
	return func(ep *endpoint) ***REMOVED***
		ep.disableResolution = true
	***REMOVED***
***REMOVED***

// CreateOptionAlias function returns an option setter for setting endpoint alias
func CreateOptionAlias(name string, alias string) EndpointOption ***REMOVED***
	return func(ep *endpoint) ***REMOVED***
		if ep.aliases == nil ***REMOVED***
			ep.aliases = make(map[string]string)
		***REMOVED***
		ep.aliases[alias] = name
	***REMOVED***
***REMOVED***

// CreateOptionService function returns an option setter for setting service binding configuration
func CreateOptionService(name, id string, vip net.IP, ingressPorts []*PortConfig, aliases []string) EndpointOption ***REMOVED***
	return func(ep *endpoint) ***REMOVED***
		ep.svcName = name
		ep.svcID = id
		ep.virtualIP = vip
		ep.ingressPorts = ingressPorts
		ep.svcAliases = aliases
	***REMOVED***
***REMOVED***

// CreateOptionMyAlias function returns an option setter for setting endpoint's self alias
func CreateOptionMyAlias(alias string) EndpointOption ***REMOVED***
	return func(ep *endpoint) ***REMOVED***
		ep.myAliases = append(ep.myAliases, alias)
	***REMOVED***
***REMOVED***

// CreateOptionLoadBalancer function returns an option setter for denoting the endpoint is a load balancer for a network
func CreateOptionLoadBalancer() EndpointOption ***REMOVED***
	return func(ep *endpoint) ***REMOVED***
		ep.loadBalancer = true
	***REMOVED***
***REMOVED***

// JoinOptionPriority function returns an option setter for priority option to
// be passed to the endpoint.Join() method.
func JoinOptionPriority(ep Endpoint, prio int) EndpointOption ***REMOVED***
	return func(ep *endpoint) ***REMOVED***
		// ep lock already acquired
		c := ep.network.getController()
		c.Lock()
		sb, ok := c.sandboxes[ep.sandboxID]
		c.Unlock()
		if !ok ***REMOVED***
			logrus.Errorf("Could not set endpoint priority value during Join to endpoint %s: No sandbox id present in endpoint", ep.id)
			return
		***REMOVED***
		sb.epPriority[ep.id] = prio
	***REMOVED***
***REMOVED***

func (ep *endpoint) DataScope() string ***REMOVED***
	return ep.getNetwork().DataScope()
***REMOVED***

func (ep *endpoint) assignAddress(ipam ipamapi.Ipam, assignIPv4, assignIPv6 bool) error ***REMOVED***
	var err error

	n := ep.getNetwork()
	if n.hasSpecialDriver() ***REMOVED***
		return nil
	***REMOVED***

	logrus.Debugf("Assigning addresses for endpoint %s's interface on network %s", ep.Name(), n.Name())

	if assignIPv4 ***REMOVED***
		if err = ep.assignAddressVersion(4, ipam); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if assignIPv6 ***REMOVED***
		err = ep.assignAddressVersion(6, ipam)
	***REMOVED***

	return err
***REMOVED***

func (ep *endpoint) assignAddressVersion(ipVer int, ipam ipamapi.Ipam) error ***REMOVED***
	var (
		poolID  *string
		address **net.IPNet
		prefAdd net.IP
		progAdd net.IP
	)

	n := ep.getNetwork()
	switch ipVer ***REMOVED***
	case 4:
		poolID = &ep.iface.v4PoolID
		address = &ep.iface.addr
		prefAdd = ep.prefAddress
	case 6:
		poolID = &ep.iface.v6PoolID
		address = &ep.iface.addrv6
		prefAdd = ep.prefAddressV6
	default:
		return types.InternalErrorf("incorrect ip version number passed: %d", ipVer)
	***REMOVED***

	ipInfo := n.getIPInfo(ipVer)

	// ipv6 address is not mandatory
	if len(ipInfo) == 0 && ipVer == 6 ***REMOVED***
		return nil
	***REMOVED***

	// The address to program may be chosen by the user or by the network driver in one specific
	// case to support backward compatibility with `docker daemon --fixed-cidrv6` use case
	if prefAdd != nil ***REMOVED***
		progAdd = prefAdd
	***REMOVED*** else if *address != nil ***REMOVED***
		progAdd = (*address).IP
	***REMOVED***

	for _, d := range ipInfo ***REMOVED***
		if progAdd != nil && !d.Pool.Contains(progAdd) ***REMOVED***
			continue
		***REMOVED***
		addr, _, err := ipam.RequestAddress(d.PoolID, progAdd, ep.ipamOptions)
		if err == nil ***REMOVED***
			ep.Lock()
			*address = addr
			*poolID = d.PoolID
			ep.Unlock()
			return nil
		***REMOVED***
		if err != ipamapi.ErrNoAvailableIPs || progAdd != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if progAdd != nil ***REMOVED***
		return types.BadRequestErrorf("Invalid address %s: It does not belong to any of this network's subnets", prefAdd)
	***REMOVED***
	return fmt.Errorf("no available IPv%d addresses on this network's address pools: %s (%s)", ipVer, n.Name(), n.ID())
***REMOVED***

func (ep *endpoint) releaseAddress() ***REMOVED***
	n := ep.getNetwork()
	if n.hasSpecialDriver() ***REMOVED***
		return
	***REMOVED***

	logrus.Debugf("Releasing addresses for endpoint %s's interface on network %s", ep.Name(), n.Name())

	ipam, _, err := n.getController().getIPAMDriver(n.ipamType)
	if err != nil ***REMOVED***
		logrus.Warnf("Failed to retrieve ipam driver to release interface address on delete of endpoint %s (%s): %v", ep.Name(), ep.ID(), err)
		return
	***REMOVED***

	if ep.iface.addr != nil ***REMOVED***
		if err := ipam.ReleaseAddress(ep.iface.v4PoolID, ep.iface.addr.IP); err != nil ***REMOVED***
			logrus.Warnf("Failed to release ip address %s on delete of endpoint %s (%s): %v", ep.iface.addr.IP, ep.Name(), ep.ID(), err)
		***REMOVED***
	***REMOVED***

	if ep.iface.addrv6 != nil && ep.iface.addrv6.IP.IsGlobalUnicast() ***REMOVED***
		if err := ipam.ReleaseAddress(ep.iface.v6PoolID, ep.iface.addrv6.IP); err != nil ***REMOVED***
			logrus.Warnf("Failed to release ip address %s on delete of endpoint %s (%s): %v", ep.iface.addrv6.IP, ep.Name(), ep.ID(), err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *controller) cleanupLocalEndpoints() ***REMOVED***
	// Get used endpoints
	eps := make(map[string]interface***REMOVED******REMOVED***)
	for _, sb := range c.sandboxes ***REMOVED***
		for _, ep := range sb.endpoints ***REMOVED***
			eps[ep.id] = true
		***REMOVED***
	***REMOVED***
	nl, err := c.getNetworksForScope(datastore.LocalScope)
	if err != nil ***REMOVED***
		logrus.Warnf("Could not get list of networks during endpoint cleanup: %v", err)
		return
	***REMOVED***

	for _, n := range nl ***REMOVED***
		if n.ConfigOnly() ***REMOVED***
			continue
		***REMOVED***
		epl, err := n.getEndpointsFromStore()
		if err != nil ***REMOVED***
			logrus.Warnf("Could not get list of endpoints in network %s during endpoint cleanup: %v", n.name, err)
			continue
		***REMOVED***

		for _, ep := range epl ***REMOVED***
			if _, ok := eps[ep.id]; ok ***REMOVED***
				continue
			***REMOVED***
			logrus.Infof("Removing stale endpoint %s (%s)", ep.name, ep.id)
			if err := ep.Delete(true); err != nil ***REMOVED***
				logrus.Warnf("Could not delete local endpoint %s during endpoint cleanup: %v", ep.name, err)
			***REMOVED***
		***REMOVED***

		epl, err = n.getEndpointsFromStore()
		if err != nil ***REMOVED***
			logrus.Warnf("Could not get list of endpoints in network %s for count update: %v", n.name, err)
			continue
		***REMOVED***

		epCnt := n.getEpCnt().EndpointCnt()
		if epCnt != uint64(len(epl)) ***REMOVED***
			logrus.Infof("Fixing inconsistent endpoint_cnt for network %s. Expected=%d, Actual=%d", n.name, len(epl), epCnt)
			n.getEpCnt().setCnt(uint64(len(epl)))
		***REMOVED***
	***REMOVED***
***REMOVED***
