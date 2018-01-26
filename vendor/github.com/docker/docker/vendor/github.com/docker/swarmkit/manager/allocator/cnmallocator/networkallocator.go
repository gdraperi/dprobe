package cnmallocator

import (
	"fmt"
	"net"
	"strings"

	"github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/drvregistry"
	"github.com/docker/libnetwork/ipamapi"
	"github.com/docker/libnetwork/netlabel"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/allocator/networkallocator"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const (
	// DefaultDriver defines the name of the driver to be used by
	// default if a network without any driver name specified is
	// created.
	DefaultDriver = "overlay"
)

// cnmNetworkAllocator acts as the controller for all network related operations
// like managing network and IPAM drivers and also creating and
// deleting networks and the associated resources.
type cnmNetworkAllocator struct ***REMOVED***
	// The driver register which manages all internal and external
	// IPAM and network drivers.
	drvRegistry *drvregistry.DrvRegistry

	// The port allocator instance for allocating node ports
	portAllocator *portAllocator

	// Local network state used by cnmNetworkAllocator to do network management.
	networks map[string]*network

	// Allocator state to indicate if allocation has been
	// successfully completed for this service.
	services map[string]struct***REMOVED******REMOVED***

	// Allocator state to indicate if allocation has been
	// successfully completed for this task.
	tasks map[string]struct***REMOVED******REMOVED***

	// Allocator state to indicate if allocation has been
	// successfully completed for this node on this network.
	// outer map key: node id
	// inner map key: network id
	nodes map[string]map[string]struct***REMOVED******REMOVED***
***REMOVED***

// Local in-memory state related to network that need to be tracked by cnmNetworkAllocator
type network struct ***REMOVED***
	// A local cache of the store object.
	nw *api.Network

	// pools is used to save the internal poolIDs needed when
	// releasing the pool.
	pools map[string]string

	// endpoints is a map of endpoint IP to the poolID from which it
	// was allocated.
	endpoints map[string]string

	// isNodeLocal indicates whether the scope of the network's resources
	// is local to the node. If true, it means the resources can only be
	// allocated locally by the node where the network will be deployed.
	// In this the swarm manager will skip the allocations.
	isNodeLocal bool
***REMOVED***

type networkDriver struct ***REMOVED***
	driver     driverapi.Driver
	name       string
	capability *driverapi.Capability
***REMOVED***

type initializer struct ***REMOVED***
	fn    drvregistry.InitFunc
	ntype string
***REMOVED***

// New returns a new NetworkAllocator handle
func New(pg plugingetter.PluginGetter) (networkallocator.NetworkAllocator, error) ***REMOVED***
	na := &cnmNetworkAllocator***REMOVED***
		networks: make(map[string]*network),
		services: make(map[string]struct***REMOVED******REMOVED***),
		tasks:    make(map[string]struct***REMOVED******REMOVED***),
		nodes:    make(map[string]map[string]struct***REMOVED******REMOVED***),
	***REMOVED***

	// There are no driver configurations and notification
	// functions as of now.
	reg, err := drvregistry.New(nil, nil, nil, nil, pg)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := initializeDrivers(reg); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err = initIPAMDrivers(reg); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	pa, err := newPortAllocator()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	na.portAllocator = pa
	na.drvRegistry = reg
	return na, nil
***REMOVED***

// Allocate allocates all the necessary resources both general
// and driver-specific which may be specified in the NetworkSpec
func (na *cnmNetworkAllocator) Allocate(n *api.Network) error ***REMOVED***
	if _, ok := na.networks[n.ID]; ok ***REMOVED***
		return fmt.Errorf("network %s already allocated", n.ID)
	***REMOVED***

	d, err := na.resolveDriver(n)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	nw := &network***REMOVED***
		nw:          n,
		endpoints:   make(map[string]string),
		isNodeLocal: d.capability.DataScope == datastore.LocalScope,
	***REMOVED***

	// No swarm-level allocation can be provided by the network driver for
	// node-local networks. Only thing needed is populating the driver's name
	// in the driver's state.
	if nw.isNodeLocal ***REMOVED***
		n.DriverState = &api.Driver***REMOVED***
			Name: d.name,
		***REMOVED***
		// In order to support backward compatibility with older daemon
		// versions which assumes the network attachment to contains
		// non nil IPAM attribute, passing an empty object
		n.IPAM = &api.IPAMOptions***REMOVED***Driver: &api.Driver***REMOVED******REMOVED******REMOVED***
	***REMOVED*** else ***REMOVED***
		nw.pools, err = na.allocatePools(n)
		if err != nil ***REMOVED***
			return errors.Wrapf(err, "failed allocating pools and gateway IP for network %s", n.ID)
		***REMOVED***

		if err := na.allocateDriverState(n); err != nil ***REMOVED***
			na.freePools(n, nw.pools)
			return errors.Wrapf(err, "failed while allocating driver state for network %s", n.ID)
		***REMOVED***
	***REMOVED***

	na.networks[n.ID] = nw

	return nil
***REMOVED***

func (na *cnmNetworkAllocator) getNetwork(id string) *network ***REMOVED***
	return na.networks[id]
***REMOVED***

// Deallocate frees all the general and driver specific resources
// which were assigned to the passed network.
func (na *cnmNetworkAllocator) Deallocate(n *api.Network) error ***REMOVED***
	localNet := na.getNetwork(n.ID)
	if localNet == nil ***REMOVED***
		return fmt.Errorf("could not get networker state for network %s", n.ID)
	***REMOVED***

	// No swarm-level resource deallocation needed for node-local networks
	if localNet.isNodeLocal ***REMOVED***
		delete(na.networks, n.ID)
		return nil
	***REMOVED***

	if err := na.freeDriverState(n); err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to free driver state for network %s", n.ID)
	***REMOVED***

	delete(na.networks, n.ID)

	return na.freePools(n, localNet.pools)
***REMOVED***

// AllocateService allocates all the network resources such as virtual
// IP and ports needed by the service.
func (na *cnmNetworkAllocator) AllocateService(s *api.Service) (err error) ***REMOVED***
	if err = na.portAllocator.serviceAllocatePorts(s); err != nil ***REMOVED***
		return err
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			na.DeallocateService(s)
		***REMOVED***
	***REMOVED***()

	if s.Endpoint == nil ***REMOVED***
		s.Endpoint = &api.Endpoint***REMOVED******REMOVED***
	***REMOVED***
	s.Endpoint.Spec = s.Spec.Endpoint.Copy()

	// If ResolutionMode is DNSRR do not try allocating VIPs, but
	// free any VIP from previous state.
	if s.Spec.Endpoint != nil && s.Spec.Endpoint.Mode == api.ResolutionModeDNSRoundRobin ***REMOVED***
		for _, vip := range s.Endpoint.VirtualIPs ***REMOVED***
			if err := na.deallocateVIP(vip); err != nil ***REMOVED***
				// don't bail here, deallocate as many as possible.
				log.L.WithError(err).
					WithField("vip.network", vip.NetworkID).
					WithField("vip.addr", vip.Addr).Error("error deallocating vip")
			***REMOVED***
		***REMOVED***

		s.Endpoint.VirtualIPs = nil

		delete(na.services, s.ID)
		return nil
	***REMOVED***

	specNetworks := serviceNetworks(s)

	// Allocate VIPs for all the pre-populated endpoint attachments
	eVIPs := s.Endpoint.VirtualIPs[:0]

vipLoop:
	for _, eAttach := range s.Endpoint.VirtualIPs ***REMOVED***
		if na.IsVIPOnIngressNetwork(eAttach) && networkallocator.IsIngressNetworkNeeded(s) ***REMOVED***
			if err = na.allocateVIP(eAttach); err != nil ***REMOVED***
				return err
			***REMOVED***
			eVIPs = append(eVIPs, eAttach)
			continue vipLoop

		***REMOVED***
		for _, nAttach := range specNetworks ***REMOVED***
			if nAttach.Target == eAttach.NetworkID ***REMOVED***
				if err = na.allocateVIP(eAttach); err != nil ***REMOVED***
					return err
				***REMOVED***
				eVIPs = append(eVIPs, eAttach)
				continue vipLoop
			***REMOVED***
		***REMOVED***
		// If the network of the VIP is not part of the service spec,
		// deallocate the vip
		na.deallocateVIP(eAttach)
	***REMOVED***

networkLoop:
	for _, nAttach := range specNetworks ***REMOVED***
		for _, vip := range s.Endpoint.VirtualIPs ***REMOVED***
			if vip.NetworkID == nAttach.Target ***REMOVED***
				continue networkLoop
			***REMOVED***
		***REMOVED***

		vip := &api.Endpoint_VirtualIP***REMOVED***NetworkID: nAttach.Target***REMOVED***
		if err = na.allocateVIP(vip); err != nil ***REMOVED***
			return err
		***REMOVED***

		eVIPs = append(eVIPs, vip)
	***REMOVED***

	if len(eVIPs) > 0 ***REMOVED***
		na.services[s.ID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED*** else ***REMOVED***
		delete(na.services, s.ID)
	***REMOVED***

	s.Endpoint.VirtualIPs = eVIPs
	return nil
***REMOVED***

// DeallocateService de-allocates all the network resources such as
// virtual IP and ports associated with the service.
func (na *cnmNetworkAllocator) DeallocateService(s *api.Service) error ***REMOVED***
	if s.Endpoint == nil ***REMOVED***
		return nil
	***REMOVED***

	for _, vip := range s.Endpoint.VirtualIPs ***REMOVED***
		if err := na.deallocateVIP(vip); err != nil ***REMOVED***
			// don't bail here, deallocate as many as possible.
			log.L.WithError(err).
				WithField("vip.network", vip.NetworkID).
				WithField("vip.addr", vip.Addr).Error("error deallocating vip")
		***REMOVED***
	***REMOVED***
	s.Endpoint.VirtualIPs = nil

	na.portAllocator.serviceDeallocatePorts(s)
	delete(na.services, s.ID)

	return nil
***REMOVED***

// IsAllocated returns if the passed network has been allocated or not.
func (na *cnmNetworkAllocator) IsAllocated(n *api.Network) bool ***REMOVED***
	_, ok := na.networks[n.ID]
	return ok
***REMOVED***

// IsTaskAllocated returns if the passed task has its network resources allocated or not.
func (na *cnmNetworkAllocator) IsTaskAllocated(t *api.Task) bool ***REMOVED***
	// If the task is not found in the allocated set, then it is
	// not allocated.
	if _, ok := na.tasks[t.ID]; !ok ***REMOVED***
		return false
	***REMOVED***

	// If Networks is empty there is no way this Task is allocated.
	if len(t.Networks) == 0 ***REMOVED***
		return false
	***REMOVED***

	// To determine whether the task has its resources allocated,
	// we just need to look at one global scope network (in case of
	// multi-network attachment).  This is because we make sure we
	// allocate for every network or we allocate for none.

	// Find the first global scope network
	for _, nAttach := range t.Networks ***REMOVED***
		// If the network is not allocated, the task cannot be allocated.
		localNet, ok := na.networks[nAttach.Network.ID]
		if !ok ***REMOVED***
			return false
		***REMOVED***

		// Nothing else to check for local scope network
		if localNet.isNodeLocal ***REMOVED***
			continue
		***REMOVED***

		// Addresses empty. Task is not allocated.
		if len(nAttach.Addresses) == 0 ***REMOVED***
			return false
		***REMOVED***

		// The allocated IP address not found in local endpoint state. Not allocated.
		if _, ok := localNet.endpoints[nAttach.Addresses[0]]; !ok ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

// HostPublishPortsNeedUpdate returns true if the passed service needs
// allocations for its published ports in host (non ingress) mode
func (na *cnmNetworkAllocator) HostPublishPortsNeedUpdate(s *api.Service) bool ***REMOVED***
	return na.portAllocator.hostPublishPortsNeedUpdate(s)
***REMOVED***

// IsServiceAllocated returns false if the passed service needs to have network resources allocated/updated.
func (na *cnmNetworkAllocator) IsServiceAllocated(s *api.Service, flags ...func(*networkallocator.ServiceAllocationOpts)) bool ***REMOVED***
	var options networkallocator.ServiceAllocationOpts
	for _, flag := range flags ***REMOVED***
		flag(&options)
	***REMOVED***

	specNetworks := serviceNetworks(s)

	// If endpoint mode is VIP and allocator does not have the
	// service in VIP allocated set then it needs to be allocated.
	if len(specNetworks) != 0 &&
		(s.Spec.Endpoint == nil ||
			s.Spec.Endpoint.Mode == api.ResolutionModeVirtualIP) ***REMOVED***

		if _, ok := na.services[s.ID]; !ok ***REMOVED***
			return false
		***REMOVED***

		if s.Endpoint == nil || len(s.Endpoint.VirtualIPs) == 0 ***REMOVED***
			return false
		***REMOVED***

		// If the spec has networks which don't have a corresponding VIP,
		// the service needs to be allocated.
	networkLoop:
		for _, net := range specNetworks ***REMOVED***
			for _, vip := range s.Endpoint.VirtualIPs ***REMOVED***
				if vip.NetworkID == net.Target ***REMOVED***
					continue networkLoop
				***REMOVED***
			***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	// If the spec no longer has networks attached and has a vip allocated
	// from previous spec the service needs to allocated.
	if s.Endpoint != nil ***REMOVED***
	vipLoop:
		for _, vip := range s.Endpoint.VirtualIPs ***REMOVED***
			if na.IsVIPOnIngressNetwork(vip) && networkallocator.IsIngressNetworkNeeded(s) ***REMOVED***
				// This checks the condition when ingress network is needed
				// but allocation has not been done.
				if _, ok := na.services[s.ID]; !ok ***REMOVED***
					return false
				***REMOVED***
				continue vipLoop
			***REMOVED***
			for _, net := range specNetworks ***REMOVED***
				if vip.NetworkID == net.Target ***REMOVED***
					continue vipLoop
				***REMOVED***
			***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	// If the endpoint mode is DNSRR and allocator has the service
	// in VIP allocated set then we return to be allocated to make
	// sure the allocator triggers networkallocator to free up the
	// resources if any.
	if s.Spec.Endpoint != nil && s.Spec.Endpoint.Mode == api.ResolutionModeDNSRoundRobin ***REMOVED***
		if _, ok := na.services[s.ID]; ok ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	if (s.Spec.Endpoint != nil && len(s.Spec.Endpoint.Ports) != 0) ||
		(s.Endpoint != nil && len(s.Endpoint.Ports) != 0) ***REMOVED***
		return na.portAllocator.isPortsAllocatedOnInit(s, options.OnInit)
	***REMOVED***
	return true
***REMOVED***

// AllocateTask allocates all the endpoint resources for all the
// networks that a task is attached to.
func (na *cnmNetworkAllocator) AllocateTask(t *api.Task) error ***REMOVED***
	for i, nAttach := range t.Networks ***REMOVED***
		if localNet := na.getNetwork(nAttach.Network.ID); localNet != nil && localNet.isNodeLocal ***REMOVED***
			continue
		***REMOVED***
		if err := na.allocateNetworkIPs(nAttach); err != nil ***REMOVED***
			if err := na.releaseEndpoints(t.Networks[:i]); err != nil ***REMOVED***
				log.G(context.TODO()).WithError(err).Errorf("failed to release IP addresses while rolling back allocation for task %s network %s", t.ID, nAttach.Network.ID)
			***REMOVED***
			return errors.Wrapf(err, "failed to allocate network IP for task %s network %s", t.ID, nAttach.Network.ID)
		***REMOVED***
	***REMOVED***

	na.tasks[t.ID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***

	return nil
***REMOVED***

// DeallocateTask releases all the endpoint resources for all the
// networks that a task is attached to.
func (na *cnmNetworkAllocator) DeallocateTask(t *api.Task) error ***REMOVED***
	delete(na.tasks, t.ID)
	return na.releaseEndpoints(t.Networks)
***REMOVED***

// IsAttachmentAllocated returns if the passed node and network has resources allocated or not.
func (na *cnmNetworkAllocator) IsAttachmentAllocated(node *api.Node, networkAttachment *api.NetworkAttachment) bool ***REMOVED***
	if node == nil ***REMOVED***
		return false
	***REMOVED***

	if networkAttachment == nil ***REMOVED***
		return false
	***REMOVED***

	// If the node is not found in the allocated set, then it is
	// not allocated.
	if _, ok := na.nodes[node.ID]; !ok ***REMOVED***
		return false
	***REMOVED***

	// If the nework is not found in the allocated set, then it is
	// not allocated.
	if _, ok := na.nodes[node.ID][networkAttachment.Network.ID]; !ok ***REMOVED***
		return false
	***REMOVED***

	// If the network is not allocated, the node cannot be allocated.
	localNet, ok := na.networks[networkAttachment.Network.ID]
	if !ok ***REMOVED***
		return false
	***REMOVED***

	// Addresses empty, not allocated.
	if len(networkAttachment.Addresses) == 0 ***REMOVED***
		return false
	***REMOVED***

	// The allocated IP address not found in local endpoint state. Not allocated.
	if _, ok := localNet.endpoints[networkAttachment.Addresses[0]]; !ok ***REMOVED***
		return false
	***REMOVED***

	return true
***REMOVED***

// AllocateAttachment allocates the IP addresses for a LB in a network
// on a given node
func (na *cnmNetworkAllocator) AllocateAttachment(node *api.Node, networkAttachment *api.NetworkAttachment) error ***REMOVED***

	if err := na.allocateNetworkIPs(networkAttachment); err != nil ***REMOVED***
		return err
	***REMOVED***

	if na.nodes[node.ID] == nil ***REMOVED***
		na.nodes[node.ID] = make(map[string]struct***REMOVED******REMOVED***)
	***REMOVED***
	na.nodes[node.ID][networkAttachment.Network.ID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***

	return nil
***REMOVED***

// DeallocateAttachment deallocates the IP addresses for a LB in a network to
// which the node is attached.
func (na *cnmNetworkAllocator) DeallocateAttachment(node *api.Node, networkAttachment *api.NetworkAttachment) error ***REMOVED***

	delete(na.nodes[node.ID], networkAttachment.Network.ID)
	if len(na.nodes[node.ID]) == 0 ***REMOVED***
		delete(na.nodes, node.ID)
	***REMOVED***

	return na.releaseEndpoints([]*api.NetworkAttachment***REMOVED***networkAttachment***REMOVED***)
***REMOVED***

func (na *cnmNetworkAllocator) releaseEndpoints(networks []*api.NetworkAttachment) error ***REMOVED***
	for _, nAttach := range networks ***REMOVED***
		localNet := na.getNetwork(nAttach.Network.ID)
		if localNet == nil ***REMOVED***
			return fmt.Errorf("could not find network allocator state for network %s", nAttach.Network.ID)
		***REMOVED***

		if localNet.isNodeLocal ***REMOVED***
			continue
		***REMOVED***

		ipam, _, _, err := na.resolveIPAM(nAttach.Network)
		if err != nil ***REMOVED***
			return errors.Wrap(err, "failed to resolve IPAM while releasing")
		***REMOVED***

		// Do not fail and bail out if we fail to release IP
		// address here. Keep going and try releasing as many
		// addresses as possible.
		for _, addr := range nAttach.Addresses ***REMOVED***
			// Retrieve the poolID and immediately nuke
			// out the mapping.
			poolID := localNet.endpoints[addr]
			delete(localNet.endpoints, addr)

			ip, _, err := net.ParseCIDR(addr)
			if err != nil ***REMOVED***
				log.G(context.TODO()).Errorf("Could not parse IP address %s while releasing", addr)
				continue
			***REMOVED***

			if err := ipam.ReleaseAddress(poolID, ip); err != nil ***REMOVED***
				log.G(context.TODO()).WithError(err).Errorf("IPAM failure while releasing IP address %s", addr)
			***REMOVED***
		***REMOVED***

		// Clear out the address list when we are done with
		// this network.
		nAttach.Addresses = nil
	***REMOVED***

	return nil
***REMOVED***

// allocate virtual IP for a single endpoint attachment of the service.
func (na *cnmNetworkAllocator) allocateVIP(vip *api.Endpoint_VirtualIP) error ***REMOVED***
	var opts map[string]string
	localNet := na.getNetwork(vip.NetworkID)
	if localNet == nil ***REMOVED***
		return errors.New("networkallocator: could not find local network state")
	***REMOVED***

	if localNet.isNodeLocal ***REMOVED***
		return nil
	***REMOVED***

	// If this IP is already allocated in memory we don't need to
	// do anything.
	if _, ok := localNet.endpoints[vip.Addr]; ok ***REMOVED***
		return nil
	***REMOVED***

	ipam, _, _, err := na.resolveIPAM(localNet.nw)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed to resolve IPAM while allocating")
	***REMOVED***

	var addr net.IP
	if vip.Addr != "" ***REMOVED***
		var err error

		addr, _, err = net.ParseCIDR(vip.Addr)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if localNet.nw.IPAM != nil && localNet.nw.IPAM.Driver != nil ***REMOVED***
		// set ipam allocation method to serial
		opts = setIPAMSerialAlloc(localNet.nw.IPAM.Driver.Options)
	***REMOVED***

	for _, poolID := range localNet.pools ***REMOVED***
		ip, _, err := ipam.RequestAddress(poolID, addr, opts)
		if err != nil && err != ipamapi.ErrNoAvailableIPs && err != ipamapi.ErrIPOutOfRange ***REMOVED***
			return errors.Wrap(err, "could not allocate VIP from IPAM")
		***REMOVED***

		// If we got an address then we are done.
		if err == nil ***REMOVED***
			ipStr := ip.String()
			localNet.endpoints[ipStr] = poolID
			vip.Addr = ipStr
			return nil
		***REMOVED***
	***REMOVED***

	return errors.New("could not find an available IP while allocating VIP")
***REMOVED***

func (na *cnmNetworkAllocator) deallocateVIP(vip *api.Endpoint_VirtualIP) error ***REMOVED***
	localNet := na.getNetwork(vip.NetworkID)
	if localNet == nil ***REMOVED***
		return errors.New("networkallocator: could not find local network state")
	***REMOVED***
	if localNet.isNodeLocal ***REMOVED***
		return nil
	***REMOVED***
	ipam, _, _, err := na.resolveIPAM(localNet.nw)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed to resolve IPAM while allocating")
	***REMOVED***

	// Retrieve the poolID and immediately nuke
	// out the mapping.
	poolID := localNet.endpoints[vip.Addr]
	delete(localNet.endpoints, vip.Addr)

	ip, _, err := net.ParseCIDR(vip.Addr)
	if err != nil ***REMOVED***
		log.G(context.TODO()).Errorf("Could not parse VIP address %s while releasing", vip.Addr)
		return err
	***REMOVED***

	if err := ipam.ReleaseAddress(poolID, ip); err != nil ***REMOVED***
		log.G(context.TODO()).WithError(err).Errorf("IPAM failure while releasing VIP address %s", vip.Addr)
		return err
	***REMOVED***

	return nil
***REMOVED***

// allocate the IP addresses for a single network attachment of the task.
func (na *cnmNetworkAllocator) allocateNetworkIPs(nAttach *api.NetworkAttachment) error ***REMOVED***
	var ip *net.IPNet
	var opts map[string]string

	ipam, _, _, err := na.resolveIPAM(nAttach.Network)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed to resolve IPAM while allocating")
	***REMOVED***

	localNet := na.getNetwork(nAttach.Network.ID)
	if localNet == nil ***REMOVED***
		return fmt.Errorf("could not find network allocator state for network %s", nAttach.Network.ID)
	***REMOVED***

	addresses := nAttach.Addresses
	if len(addresses) == 0 ***REMOVED***
		addresses = []string***REMOVED***""***REMOVED***
	***REMOVED***

	for i, rawAddr := range addresses ***REMOVED***
		var addr net.IP
		if rawAddr != "" ***REMOVED***
			var err error
			addr, _, err = net.ParseCIDR(rawAddr)
			if err != nil ***REMOVED***
				addr = net.ParseIP(rawAddr)

				if addr == nil ***REMOVED***
					return errors.Wrapf(err, "could not parse address string %s", rawAddr)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		// Set the ipam options if the network has an ipam driver.
		if localNet.nw.IPAM != nil && localNet.nw.IPAM.Driver != nil ***REMOVED***
			// set ipam allocation method to serial
			opts = setIPAMSerialAlloc(localNet.nw.IPAM.Driver.Options)
		***REMOVED***

		for _, poolID := range localNet.pools ***REMOVED***
			var err error

			ip, _, err = ipam.RequestAddress(poolID, addr, opts)
			if err != nil && err != ipamapi.ErrNoAvailableIPs && err != ipamapi.ErrIPOutOfRange ***REMOVED***
				return errors.Wrap(err, "could not allocate IP from IPAM")
			***REMOVED***

			// If we got an address then we are done.
			if err == nil ***REMOVED***
				ipStr := ip.String()
				localNet.endpoints[ipStr] = poolID
				addresses[i] = ipStr
				nAttach.Addresses = addresses
				return nil
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return errors.New("could not find an available IP")
***REMOVED***

func (na *cnmNetworkAllocator) freeDriverState(n *api.Network) error ***REMOVED***
	d, err := na.resolveDriver(n)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return d.driver.NetworkFree(n.ID)
***REMOVED***

func (na *cnmNetworkAllocator) allocateDriverState(n *api.Network) error ***REMOVED***
	d, err := na.resolveDriver(n)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	options := make(map[string]string)
	// reconcile the driver specific options from the network spec
	// and from the operational state retrieved from the store
	if n.Spec.DriverConfig != nil ***REMOVED***
		for k, v := range n.Spec.DriverConfig.Options ***REMOVED***
			options[k] = v
		***REMOVED***
	***REMOVED***
	if n.DriverState != nil ***REMOVED***
		for k, v := range n.DriverState.Options ***REMOVED***
			options[k] = v
		***REMOVED***
	***REMOVED***

	// Construct IPAM data for driver consumption.
	ipv4Data := make([]driverapi.IPAMData, 0, len(n.IPAM.Configs))
	for _, ic := range n.IPAM.Configs ***REMOVED***
		if ic.Family == api.IPAMConfig_IPV6 ***REMOVED***
			continue
		***REMOVED***

		_, subnet, err := net.ParseCIDR(ic.Subnet)
		if err != nil ***REMOVED***
			return errors.Wrapf(err, "error parsing subnet %s while allocating driver state", ic.Subnet)
		***REMOVED***

		gwIP := net.ParseIP(ic.Gateway)
		gwNet := &net.IPNet***REMOVED***
			IP:   gwIP,
			Mask: subnet.Mask,
		***REMOVED***

		data := driverapi.IPAMData***REMOVED***
			Pool:    subnet,
			Gateway: gwNet,
		***REMOVED***

		ipv4Data = append(ipv4Data, data)
	***REMOVED***

	ds, err := d.driver.NetworkAllocate(n.ID, options, ipv4Data, nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Update network object with the obtained driver state.
	n.DriverState = &api.Driver***REMOVED***
		Name:    d.name,
		Options: ds,
	***REMOVED***

	return nil
***REMOVED***

// Resolve network driver
func (na *cnmNetworkAllocator) resolveDriver(n *api.Network) (*networkDriver, error) ***REMOVED***
	dName := DefaultDriver
	if n.Spec.DriverConfig != nil && n.Spec.DriverConfig.Name != "" ***REMOVED***
		dName = n.Spec.DriverConfig.Name
	***REMOVED***

	d, drvcap := na.drvRegistry.Driver(dName)
	if d == nil ***REMOVED***
		var err error
		err = na.loadDriver(dName)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		d, drvcap = na.drvRegistry.Driver(dName)
		if d == nil ***REMOVED***
			return nil, fmt.Errorf("could not resolve network driver %s", dName)
		***REMOVED***
	***REMOVED***

	return &networkDriver***REMOVED***driver: d, capability: drvcap, name: dName***REMOVED***, nil
***REMOVED***

func (na *cnmNetworkAllocator) loadDriver(name string) error ***REMOVED***
	pg := na.drvRegistry.GetPluginGetter()
	if pg == nil ***REMOVED***
		return errors.New("plugin store is uninitialized")
	***REMOVED***
	_, err := pg.Get(name, driverapi.NetworkPluginEndpointType, plugingetter.Lookup)
	return err
***REMOVED***

// Resolve the IPAM driver
func (na *cnmNetworkAllocator) resolveIPAM(n *api.Network) (ipamapi.Ipam, string, map[string]string, error) ***REMOVED***
	dName := ipamapi.DefaultIPAM
	if n.Spec.IPAM != nil && n.Spec.IPAM.Driver != nil && n.Spec.IPAM.Driver.Name != "" ***REMOVED***
		dName = n.Spec.IPAM.Driver.Name
	***REMOVED***

	var dOptions map[string]string
	if n.Spec.IPAM != nil && n.Spec.IPAM.Driver != nil && len(n.Spec.IPAM.Driver.Options) != 0 ***REMOVED***
		dOptions = n.Spec.IPAM.Driver.Options
	***REMOVED***

	ipam, _ := na.drvRegistry.IPAM(dName)
	if ipam == nil ***REMOVED***
		return nil, "", nil, fmt.Errorf("could not resolve IPAM driver %s", dName)
	***REMOVED***

	return ipam, dName, dOptions, nil
***REMOVED***

func (na *cnmNetworkAllocator) freePools(n *api.Network, pools map[string]string) error ***REMOVED***
	ipam, _, _, err := na.resolveIPAM(n)
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to resolve IPAM while freeing pools for network %s", n.ID)
	***REMOVED***

	releasePools(ipam, n.IPAM.Configs, pools)
	return nil
***REMOVED***

func releasePools(ipam ipamapi.Ipam, icList []*api.IPAMConfig, pools map[string]string) ***REMOVED***
	for _, ic := range icList ***REMOVED***
		if err := ipam.ReleaseAddress(pools[ic.Subnet], net.ParseIP(ic.Gateway)); err != nil ***REMOVED***
			log.G(context.TODO()).WithError(err).Errorf("Failed to release address %s", ic.Subnet)
		***REMOVED***
	***REMOVED***

	for k, p := range pools ***REMOVED***
		if err := ipam.ReleasePool(p); err != nil ***REMOVED***
			log.G(context.TODO()).WithError(err).Errorf("Failed to release pool %s", k)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (na *cnmNetworkAllocator) allocatePools(n *api.Network) (map[string]string, error) ***REMOVED***
	ipam, dName, dOptions, err := na.resolveIPAM(n)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// We don't support user defined address spaces yet so just
	// retrieve default address space names for the driver.
	_, asName, err := na.drvRegistry.IPAMDefaultAddressSpaces(dName)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	pools := make(map[string]string)

	var ipamConfigs []*api.IPAMConfig

	// If there is non-nil IPAM state always prefer those subnet
	// configs over Spec configs.
	if n.IPAM != nil ***REMOVED***
		ipamConfigs = n.IPAM.Configs
	***REMOVED*** else if n.Spec.IPAM != nil ***REMOVED***
		ipamConfigs = make([]*api.IPAMConfig, len(n.Spec.IPAM.Configs))
		copy(ipamConfigs, n.Spec.IPAM.Configs)
	***REMOVED***

	// Append an empty slot for subnet allocation if there are no
	// IPAM configs from either spec or state.
	if len(ipamConfigs) == 0 ***REMOVED***
		ipamConfigs = append(ipamConfigs, &api.IPAMConfig***REMOVED***Family: api.IPAMConfig_IPV4***REMOVED***)
	***REMOVED***

	// Update the runtime IPAM configurations with initial state
	n.IPAM = &api.IPAMOptions***REMOVED***
		Driver:  &api.Driver***REMOVED***Name: dName, Options: dOptions***REMOVED***,
		Configs: ipamConfigs,
	***REMOVED***

	for i, ic := range ipamConfigs ***REMOVED***
		poolID, poolIP, meta, err := ipam.RequestPool(asName, ic.Subnet, ic.Range, dOptions, false)
		if err != nil ***REMOVED***
			// Rollback by releasing all the resources allocated so far.
			releasePools(ipam, ipamConfigs[:i], pools)
			return nil, err
		***REMOVED***
		pools[poolIP.String()] = poolID

		// The IPAM contract allows the IPAM driver to autonomously
		// provide a network gateway in response to the pool request.
		// But if the network spec contains a gateway, we will allocate
		// it irrespective of whether the ipam driver returned one already.
		// If none of the above is true, we need to allocate one now, and
		// let the driver know this request is for the network gateway.
		var (
			gwIP *net.IPNet
			ip   net.IP
		)
		if gws, ok := meta[netlabel.Gateway]; ok ***REMOVED***
			if ip, gwIP, err = net.ParseCIDR(gws); err != nil ***REMOVED***
				return nil, fmt.Errorf("failed to parse gateway address (%v) returned by ipam driver: %v", gws, err)
			***REMOVED***
			gwIP.IP = ip
		***REMOVED***
		if dOptions == nil ***REMOVED***
			dOptions = make(map[string]string)
		***REMOVED***
		dOptions[ipamapi.RequestAddressType] = netlabel.Gateway
		// set ipam allocation method to serial
		dOptions = setIPAMSerialAlloc(dOptions)
		defer delete(dOptions, ipamapi.RequestAddressType)

		if ic.Gateway != "" || gwIP == nil ***REMOVED***
			gwIP, _, err = ipam.RequestAddress(poolID, net.ParseIP(ic.Gateway), dOptions)
			if err != nil ***REMOVED***
				// Rollback by releasing all the resources allocated so far.
				releasePools(ipam, ipamConfigs[:i], pools)
				return nil, err
			***REMOVED***
		***REMOVED***

		if ic.Subnet == "" ***REMOVED***
			ic.Subnet = poolIP.String()
		***REMOVED***

		if ic.Gateway == "" ***REMOVED***
			ic.Gateway = gwIP.IP.String()
		***REMOVED***

	***REMOVED***

	return pools, nil
***REMOVED***

func initializeDrivers(reg *drvregistry.DrvRegistry) error ***REMOVED***
	for _, i := range initializers ***REMOVED***
		if err := reg.AddDriver(i.ntype, i.fn, nil); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func serviceNetworks(s *api.Service) []*api.NetworkAttachmentConfig ***REMOVED***
	// Always prefer NetworkAttachmentConfig in the TaskSpec
	if len(s.Spec.Task.Networks) == 0 && len(s.Spec.Networks) != 0 ***REMOVED***
		return s.Spec.Networks
	***REMOVED***
	return s.Spec.Task.Networks
***REMOVED***

// IsVIPOnIngressNetwork check if the vip is in ingress network
func (na *cnmNetworkAllocator) IsVIPOnIngressNetwork(vip *api.Endpoint_VirtualIP) bool ***REMOVED***
	if vip == nil ***REMOVED***
		return false
	***REMOVED***

	localNet := na.getNetwork(vip.NetworkID)
	if localNet != nil && localNet.nw != nil ***REMOVED***
		return networkallocator.IsIngressNetwork(localNet.nw)
	***REMOVED***
	return false
***REMOVED***

// IsBuiltInDriver returns whether the passed driver is an internal network driver
func IsBuiltInDriver(name string) bool ***REMOVED***
	n := strings.ToLower(name)
	for _, d := range initializers ***REMOVED***
		if n == d.ntype ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// setIPAMSerialAlloc sets the ipam allocation method to serial
func setIPAMSerialAlloc(opts map[string]string) map[string]string ***REMOVED***
	if opts == nil ***REMOVED***
		opts = make(map[string]string)
	***REMOVED***
	if _, ok := opts[ipamapi.AllocSerialPrefix]; !ok ***REMOVED***
		opts[ipamapi.AllocSerialPrefix] = "true"
	***REMOVED***
	return opts
***REMOVED***
