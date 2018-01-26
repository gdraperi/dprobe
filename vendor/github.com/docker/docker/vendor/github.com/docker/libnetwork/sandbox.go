package libnetwork

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/docker/libnetwork/etchosts"
	"github.com/docker/libnetwork/netlabel"
	"github.com/docker/libnetwork/osl"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

// Sandbox provides the control over the network container entity. It is a one to one mapping with the container.
type Sandbox interface ***REMOVED***
	// ID returns the ID of the sandbox
	ID() string
	// Key returns the sandbox's key
	Key() string
	// ContainerID returns the container id associated to this sandbox
	ContainerID() string
	// Labels returns the sandbox's labels
	Labels() map[string]interface***REMOVED******REMOVED***
	// Statistics retrieves the interfaces' statistics for the sandbox
	Statistics() (map[string]*types.InterfaceStatistics, error)
	// Refresh leaves all the endpoints, resets and re-applies the options,
	// re-joins all the endpoints without destroying the osl sandbox
	Refresh(options ...SandboxOption) error
	// SetKey updates the Sandbox Key
	SetKey(key string) error
	// Rename changes the name of all attached Endpoints
	Rename(name string) error
	// Delete destroys this container after detaching it from all connected endpoints.
	Delete() error
	// Endpoints returns all the endpoints connected to the sandbox
	Endpoints() []Endpoint
	// ResolveService returns all the backend details about the containers or hosts
	// backing a service. Its purpose is to satisfy an SRV query
	ResolveService(name string) ([]*net.SRV, []net.IP)
	// EnableService  makes a managed container's service available by adding the
	// endpoint to the service load balancer and service discovery
	EnableService() error
	// DisableService removes a managed container's endpoints from the load balancer
	// and service discovery
	DisableService() error
***REMOVED***

// SandboxOption is an option setter function type used to pass various options to
// NewNetContainer method. The various setter functions of type SandboxOption are
// provided by libnetwork, they look like ContainerOptionXXXX(...)
type SandboxOption func(sb *sandbox)

func (sb *sandbox) processOptions(options ...SandboxOption) ***REMOVED***
	for _, opt := range options ***REMOVED***
		if opt != nil ***REMOVED***
			opt(sb)
		***REMOVED***
	***REMOVED***
***REMOVED***

type epHeap []*endpoint

type sandbox struct ***REMOVED***
	id                 string
	containerID        string
	config             containerConfig
	extDNS             []extDNSEntry
	osSbox             osl.Sandbox
	controller         *controller
	resolver           Resolver
	resolverOnce       sync.Once
	refCnt             int
	endpoints          epHeap
	epPriority         map[string]int
	populatedEndpoints map[string]struct***REMOVED******REMOVED***
	joinLeaveDone      chan struct***REMOVED******REMOVED***
	dbIndex            uint64
	dbExists           bool
	isStub             bool
	inDelete           bool
	ingress            bool
	ndotsSet           bool
	sync.Mutex
	// This mutex is used to serialize service related operation for an endpoint
	// The lock is here because the endpoint is saved into the store so is not unique
	Service sync.Mutex
***REMOVED***

// These are the container configs used to customize container /etc/hosts file.
type hostsPathConfig struct ***REMOVED***
	hostName        string
	domainName      string
	hostsPath       string
	originHostsPath string
	extraHosts      []extraHost
	parentUpdates   []parentUpdate
***REMOVED***

type parentUpdate struct ***REMOVED***
	cid  string
	name string
	ip   string
***REMOVED***

type extraHost struct ***REMOVED***
	name string
	IP   string
***REMOVED***

// These are the container configs used to customize container /etc/resolv.conf file.
type resolvConfPathConfig struct ***REMOVED***
	resolvConfPath       string
	originResolvConfPath string
	resolvConfHashFile   string
	dnsList              []string
	dnsSearchList        []string
	dnsOptionsList       []string
***REMOVED***

type containerConfig struct ***REMOVED***
	hostsPathConfig
	resolvConfPathConfig
	generic           map[string]interface***REMOVED******REMOVED***
	useDefaultSandBox bool
	useExternalKey    bool
	prio              int // higher the value, more the priority
	exposedPorts      []types.TransportPort
***REMOVED***

const (
	resolverIPSandbox = "127.0.0.11"
)

func (sb *sandbox) ID() string ***REMOVED***
	return sb.id
***REMOVED***

func (sb *sandbox) ContainerID() string ***REMOVED***
	return sb.containerID
***REMOVED***

func (sb *sandbox) Key() string ***REMOVED***
	if sb.config.useDefaultSandBox ***REMOVED***
		return osl.GenerateKey("default")
	***REMOVED***
	return osl.GenerateKey(sb.id)
***REMOVED***

func (sb *sandbox) Labels() map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	sb.Lock()
	defer sb.Unlock()
	opts := make(map[string]interface***REMOVED******REMOVED***, len(sb.config.generic))
	for k, v := range sb.config.generic ***REMOVED***
		opts[k] = v
	***REMOVED***
	return opts
***REMOVED***

func (sb *sandbox) Statistics() (map[string]*types.InterfaceStatistics, error) ***REMOVED***
	m := make(map[string]*types.InterfaceStatistics)

	sb.Lock()
	osb := sb.osSbox
	sb.Unlock()
	if osb == nil ***REMOVED***
		return m, nil
	***REMOVED***

	var err error
	for _, i := range osb.Info().Interfaces() ***REMOVED***
		if m[i.DstName()], err = i.Statistics(); err != nil ***REMOVED***
			return m, err
		***REMOVED***
	***REMOVED***

	return m, nil
***REMOVED***

func (sb *sandbox) Delete() error ***REMOVED***
	return sb.delete(false)
***REMOVED***

func (sb *sandbox) delete(force bool) error ***REMOVED***
	sb.Lock()
	if sb.inDelete ***REMOVED***
		sb.Unlock()
		return types.ForbiddenErrorf("another sandbox delete in progress")
	***REMOVED***
	// Set the inDelete flag. This will ensure that we don't
	// update the store until we have completed all the endpoint
	// leaves and deletes. And when endpoint leaves and deletes
	// are completed then we can finally delete the sandbox object
	// altogether from the data store. If the daemon exits
	// ungracefully in the middle of a sandbox delete this way we
	// will have all the references to the endpoints in the
	// sandbox so that we can clean them up when we restart
	sb.inDelete = true
	sb.Unlock()

	c := sb.controller

	// Detach from all endpoints
	retain := false
	for _, ep := range sb.getConnectedEndpoints() ***REMOVED***
		// gw network endpoint detach and removal are automatic
		if ep.endpointInGWNetwork() && !force ***REMOVED***
			continue
		***REMOVED***
		// Retain the sanbdox if we can't obtain the network from store.
		if _, err := c.getNetworkFromStore(ep.getNetwork().ID()); err != nil ***REMOVED***
			if c.isDistributedControl() ***REMOVED***
				retain = true
			***REMOVED***
			logrus.Warnf("Failed getting network for ep %s during sandbox %s delete: %v", ep.ID(), sb.ID(), err)
			continue
		***REMOVED***

		if !force ***REMOVED***
			if err := ep.Leave(sb); err != nil ***REMOVED***
				logrus.Warnf("Failed detaching sandbox %s from endpoint %s: %v\n", sb.ID(), ep.ID(), err)
			***REMOVED***
		***REMOVED***

		if err := ep.Delete(force); err != nil ***REMOVED***
			logrus.Warnf("Failed deleting endpoint %s: %v\n", ep.ID(), err)
		***REMOVED***
	***REMOVED***

	if retain ***REMOVED***
		sb.Lock()
		sb.inDelete = false
		sb.Unlock()
		return fmt.Errorf("could not cleanup all the endpoints in container %s / sandbox %s", sb.containerID, sb.id)
	***REMOVED***
	// Container is going away. Path cache in etchosts is most
	// likely not required any more. Drop it.
	etchosts.Drop(sb.config.hostsPath)

	if sb.resolver != nil ***REMOVED***
		sb.resolver.Stop()
	***REMOVED***

	if sb.osSbox != nil && !sb.config.useDefaultSandBox ***REMOVED***
		sb.osSbox.Destroy()
	***REMOVED***

	if err := sb.storeDelete(); err != nil ***REMOVED***
		logrus.Warnf("Failed to delete sandbox %s from store: %v", sb.ID(), err)
	***REMOVED***

	c.Lock()
	if sb.ingress ***REMOVED***
		c.ingressSandbox = nil
	***REMOVED***
	delete(c.sandboxes, sb.ID())
	c.Unlock()

	return nil
***REMOVED***

func (sb *sandbox) Rename(name string) error ***REMOVED***
	var err error

	for _, ep := range sb.getConnectedEndpoints() ***REMOVED***
		if ep.endpointInGWNetwork() ***REMOVED***
			continue
		***REMOVED***

		oldName := ep.Name()
		lEp := ep
		if err = ep.rename(name); err != nil ***REMOVED***
			break
		***REMOVED***

		defer func() ***REMOVED***
			if err != nil ***REMOVED***
				lEp.rename(oldName)
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	return err
***REMOVED***

func (sb *sandbox) Refresh(options ...SandboxOption) error ***REMOVED***
	// Store connected endpoints
	epList := sb.getConnectedEndpoints()

	// Detach from all endpoints
	for _, ep := range epList ***REMOVED***
		if err := ep.Leave(sb); err != nil ***REMOVED***
			logrus.Warnf("Failed detaching sandbox %s from endpoint %s: %v\n", sb.ID(), ep.ID(), err)
		***REMOVED***
	***REMOVED***

	// Re-apply options
	sb.config = containerConfig***REMOVED******REMOVED***
	sb.processOptions(options...)

	// Setup discovery files
	if err := sb.setupResolutionFiles(); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Re-connect to all endpoints
	for _, ep := range epList ***REMOVED***
		if err := ep.Join(sb); err != nil ***REMOVED***
			logrus.Warnf("Failed attach sandbox %s to endpoint %s: %v\n", sb.ID(), ep.ID(), err)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (sb *sandbox) MarshalJSON() ([]byte, error) ***REMOVED***
	sb.Lock()
	defer sb.Unlock()

	// We are just interested in the container ID. This can be expanded to include all of containerInfo if there is a need
	return json.Marshal(sb.id)
***REMOVED***

func (sb *sandbox) UnmarshalJSON(b []byte) (err error) ***REMOVED***
	sb.Lock()
	defer sb.Unlock()

	var id string
	if err := json.Unmarshal(b, &id); err != nil ***REMOVED***
		return err
	***REMOVED***
	sb.id = id
	return nil
***REMOVED***

func (sb *sandbox) Endpoints() []Endpoint ***REMOVED***
	sb.Lock()
	defer sb.Unlock()

	endpoints := make([]Endpoint, len(sb.endpoints))
	for i, ep := range sb.endpoints ***REMOVED***
		endpoints[i] = ep
	***REMOVED***
	return endpoints
***REMOVED***

func (sb *sandbox) getConnectedEndpoints() []*endpoint ***REMOVED***
	sb.Lock()
	defer sb.Unlock()

	eps := make([]*endpoint, len(sb.endpoints))
	for i, ep := range sb.endpoints ***REMOVED***
		eps[i] = ep
	***REMOVED***

	return eps
***REMOVED***

func (sb *sandbox) removeEndpoint(ep *endpoint) ***REMOVED***
	sb.Lock()
	defer sb.Unlock()

	for i, e := range sb.endpoints ***REMOVED***
		if e == ep ***REMOVED***
			heap.Remove(&sb.endpoints, i)
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (sb *sandbox) getEndpoint(id string) *endpoint ***REMOVED***
	sb.Lock()
	defer sb.Unlock()

	for _, ep := range sb.endpoints ***REMOVED***
		if ep.id == id ***REMOVED***
			return ep
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (sb *sandbox) updateGateway(ep *endpoint) error ***REMOVED***
	sb.Lock()
	osSbox := sb.osSbox
	sb.Unlock()
	if osSbox == nil ***REMOVED***
		return nil
	***REMOVED***
	osSbox.UnsetGateway()
	osSbox.UnsetGatewayIPv6()

	if ep == nil ***REMOVED***
		return nil
	***REMOVED***

	ep.Lock()
	joinInfo := ep.joinInfo
	ep.Unlock()

	if err := osSbox.SetGateway(joinInfo.gw); err != nil ***REMOVED***
		return fmt.Errorf("failed to set gateway while updating gateway: %v", err)
	***REMOVED***

	if err := osSbox.SetGatewayIPv6(joinInfo.gw6); err != nil ***REMOVED***
		return fmt.Errorf("failed to set IPv6 gateway while updating gateway: %v", err)
	***REMOVED***

	return nil
***REMOVED***

func (sb *sandbox) HandleQueryResp(name string, ip net.IP) ***REMOVED***
	for _, ep := range sb.getConnectedEndpoints() ***REMOVED***
		n := ep.getNetwork()
		n.HandleQueryResp(name, ip)
	***REMOVED***
***REMOVED***

func (sb *sandbox) ResolveIP(ip string) string ***REMOVED***
	var svc string
	logrus.Debugf("IP To resolve %v", ip)

	for _, ep := range sb.getConnectedEndpoints() ***REMOVED***
		n := ep.getNetwork()
		svc = n.ResolveIP(ip)
		if len(svc) != 0 ***REMOVED***
			return svc
		***REMOVED***
	***REMOVED***

	return svc
***REMOVED***

func (sb *sandbox) ExecFunc(f func()) error ***REMOVED***
	sb.Lock()
	osSbox := sb.osSbox
	sb.Unlock()
	if osSbox != nil ***REMOVED***
		return osSbox.InvokeFunc(f)
	***REMOVED***
	return fmt.Errorf("osl sandbox unavailable in ExecFunc for %v", sb.ContainerID())
***REMOVED***

func (sb *sandbox) ResolveService(name string) ([]*net.SRV, []net.IP) ***REMOVED***
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

	for _, ep := range sb.getConnectedEndpoints() ***REMOVED***
		n := ep.getNetwork()

		srv, ip = n.ResolveService(name)
		if len(srv) > 0 ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return srv, ip
***REMOVED***

func getDynamicNwEndpoints(epList []*endpoint) []*endpoint ***REMOVED***
	eps := []*endpoint***REMOVED******REMOVED***
	for _, ep := range epList ***REMOVED***
		n := ep.getNetwork()
		if n.dynamic && !n.ingress ***REMOVED***
			eps = append(eps, ep)
		***REMOVED***
	***REMOVED***
	return eps
***REMOVED***

func getIngressNwEndpoint(epList []*endpoint) *endpoint ***REMOVED***
	for _, ep := range epList ***REMOVED***
		n := ep.getNetwork()
		if n.ingress ***REMOVED***
			return ep
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func getLocalNwEndpoints(epList []*endpoint) []*endpoint ***REMOVED***
	eps := []*endpoint***REMOVED******REMOVED***
	for _, ep := range epList ***REMOVED***
		n := ep.getNetwork()
		if !n.dynamic && !n.ingress ***REMOVED***
			eps = append(eps, ep)
		***REMOVED***
	***REMOVED***
	return eps
***REMOVED***

func (sb *sandbox) ResolveName(name string, ipType int) ([]net.IP, bool) ***REMOVED***
	// Embedded server owns the docker network domain. Resolution should work
	// for both container_name and container_name.network_name
	// We allow '.' in service name and network name. For a name a.b.c.d the
	// following have to tried;
	// ***REMOVED***a.b.c.d in the networks container is connected to***REMOVED***
	// ***REMOVED***a.b.c in network d***REMOVED***,
	// ***REMOVED***a.b in network c.d***REMOVED***,
	// ***REMOVED***a in network b.c.d***REMOVED***,

	logrus.Debugf("Name To resolve: %v", name)
	name = strings.TrimSuffix(name, ".")
	reqName := []string***REMOVED***name***REMOVED***
	networkName := []string***REMOVED***""***REMOVED***

	if strings.Contains(name, ".") ***REMOVED***
		var i int
		dup := name
		for ***REMOVED***
			if i = strings.LastIndex(dup, "."); i == -1 ***REMOVED***
				break
			***REMOVED***
			networkName = append(networkName, name[i+1:])
			reqName = append(reqName, name[:i])

			dup = dup[:i]
		***REMOVED***
	***REMOVED***

	epList := sb.getConnectedEndpoints()

	// In swarm mode services with exposed ports are connected to user overlay
	// network, ingress network and docker_gwbridge network. Name resolution
	// should prioritize returning the VIP/IPs on user overlay network.
	newList := []*endpoint***REMOVED******REMOVED***
	if !sb.controller.isDistributedControl() ***REMOVED***
		newList = append(newList, getDynamicNwEndpoints(epList)...)
		ingressEP := getIngressNwEndpoint(epList)
		if ingressEP != nil ***REMOVED***
			newList = append(newList, ingressEP)
		***REMOVED***
		newList = append(newList, getLocalNwEndpoints(epList)...)
		epList = newList
	***REMOVED***

	for i := 0; i < len(reqName); i++ ***REMOVED***

		// First check for local container alias
		ip, ipv6Miss := sb.resolveName(reqName[i], networkName[i], epList, true, ipType)
		if ip != nil ***REMOVED***
			return ip, false
		***REMOVED***
		if ipv6Miss ***REMOVED***
			return ip, ipv6Miss
		***REMOVED***

		// Resolve the actual container name
		ip, ipv6Miss = sb.resolveName(reqName[i], networkName[i], epList, false, ipType)
		if ip != nil ***REMOVED***
			return ip, false
		***REMOVED***
		if ipv6Miss ***REMOVED***
			return ip, ipv6Miss
		***REMOVED***
	***REMOVED***
	return nil, false
***REMOVED***

func (sb *sandbox) resolveName(req string, networkName string, epList []*endpoint, alias bool, ipType int) ([]net.IP, bool) ***REMOVED***
	var ipv6Miss bool

	for _, ep := range epList ***REMOVED***
		name := req
		n := ep.getNetwork()

		if networkName != "" && networkName != n.Name() ***REMOVED***
			continue
		***REMOVED***

		if alias ***REMOVED***
			if ep.aliases == nil ***REMOVED***
				continue
			***REMOVED***

			var ok bool
			ep.Lock()
			name, ok = ep.aliases[req]
			ep.Unlock()
			if !ok ***REMOVED***
				continue
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// If it is a regular lookup and if the requested name is an alias
			// don't perform a svc lookup for this endpoint.
			ep.Lock()
			if _, ok := ep.aliases[req]; ok ***REMOVED***
				ep.Unlock()
				continue
			***REMOVED***
			ep.Unlock()
		***REMOVED***

		ip, miss := n.ResolveName(name, ipType)

		if ip != nil ***REMOVED***
			return ip, false
		***REMOVED***

		if miss ***REMOVED***
			ipv6Miss = miss
		***REMOVED***
	***REMOVED***
	return nil, ipv6Miss
***REMOVED***

func (sb *sandbox) SetKey(basePath string) error ***REMOVED***
	start := time.Now()
	defer func() ***REMOVED***
		logrus.Debugf("sandbox set key processing took %s for container %s", time.Since(start), sb.ContainerID())
	***REMOVED***()

	if basePath == "" ***REMOVED***
		return types.BadRequestErrorf("invalid sandbox key")
	***REMOVED***

	sb.Lock()
	if sb.inDelete ***REMOVED***
		sb.Unlock()
		return types.ForbiddenErrorf("failed to SetKey: sandbox %q delete in progress", sb.id)
	***REMOVED***
	oldosSbox := sb.osSbox
	sb.Unlock()

	if oldosSbox != nil ***REMOVED***
		// If we already have an OS sandbox, release the network resources from that
		// and destroy the OS snab. We are moving into a new home further down. Note that none
		// of the network resources gets destroyed during the move.
		sb.releaseOSSbox()
	***REMOVED***

	osSbox, err := osl.GetSandboxForExternalKey(basePath, sb.Key())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	sb.Lock()
	sb.osSbox = osSbox
	sb.Unlock()

	// If the resolver was setup before stop it and set it up in the
	// new osl sandbox.
	if oldosSbox != nil && sb.resolver != nil ***REMOVED***
		sb.resolver.Stop()

		if err := sb.osSbox.InvokeFunc(sb.resolver.SetupFunc(0)); err == nil ***REMOVED***
			if err := sb.resolver.Start(); err != nil ***REMOVED***
				logrus.Errorf("Resolver Start failed for container %s, %q", sb.ContainerID(), err)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			logrus.Errorf("Resolver Setup Function failed for container %s, %q", sb.ContainerID(), err)
		***REMOVED***
	***REMOVED***

	for _, ep := range sb.getConnectedEndpoints() ***REMOVED***
		if err = sb.populateNetworkResources(ep); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (sb *sandbox) EnableService() (err error) ***REMOVED***
	logrus.Debugf("EnableService %s START", sb.containerID)
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			sb.DisableService()
		***REMOVED***
	***REMOVED***()
	for _, ep := range sb.getConnectedEndpoints() ***REMOVED***
		if !ep.isServiceEnabled() ***REMOVED***
			if err := ep.addServiceInfoToCluster(sb); err != nil ***REMOVED***
				return fmt.Errorf("could not update state for endpoint %s into cluster: %v", ep.Name(), err)
			***REMOVED***
			ep.enableService()
		***REMOVED***
	***REMOVED***
	logrus.Debugf("EnableService %s DONE", sb.containerID)
	return nil
***REMOVED***

func (sb *sandbox) DisableService() (err error) ***REMOVED***
	logrus.Debugf("DisableService %s START", sb.containerID)
	failedEps := []string***REMOVED******REMOVED***
	defer func() ***REMOVED***
		if len(failedEps) > 0 ***REMOVED***
			err = fmt.Errorf("failed to disable service on sandbox:%s, for endpoints %s", sb.ID(), strings.Join(failedEps, ","))
		***REMOVED***
	***REMOVED***()
	for _, ep := range sb.getConnectedEndpoints() ***REMOVED***
		if ep.isServiceEnabled() ***REMOVED***
			if err := ep.deleteServiceInfoFromCluster(sb, "DisableService"); err != nil ***REMOVED***
				failedEps = append(failedEps, ep.Name())
				logrus.Warnf("failed update state for endpoint %s into cluster: %v", ep.Name(), err)
			***REMOVED***
			ep.disableService()
		***REMOVED***
	***REMOVED***
	logrus.Debugf("DisableService %s DONE", sb.containerID)
	return nil
***REMOVED***

func releaseOSSboxResources(osSbox osl.Sandbox, ep *endpoint) ***REMOVED***
	for _, i := range osSbox.Info().Interfaces() ***REMOVED***
		// Only remove the interfaces owned by this endpoint from the sandbox.
		if ep.hasInterface(i.SrcName()) ***REMOVED***
			if err := i.Remove(); err != nil ***REMOVED***
				logrus.Debugf("Remove interface %s failed: %v", i.SrcName(), err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	ep.Lock()
	joinInfo := ep.joinInfo
	vip := ep.virtualIP
	ep.Unlock()

	if len(vip) != 0 ***REMOVED***
		if err := osSbox.RemoveLoopbackAliasIP(&net.IPNet***REMOVED***IP: vip, Mask: net.CIDRMask(32, 32)***REMOVED***); err != nil ***REMOVED***
			logrus.Warnf("Remove virtual IP %v failed: %v", vip, err)
		***REMOVED***
	***REMOVED***

	if joinInfo == nil ***REMOVED***
		return
	***REMOVED***

	// Remove non-interface routes.
	for _, r := range joinInfo.StaticRoutes ***REMOVED***
		if err := osSbox.RemoveStaticRoute(r); err != nil ***REMOVED***
			logrus.Debugf("Remove route failed: %v", err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (sb *sandbox) releaseOSSbox() ***REMOVED***
	sb.Lock()
	osSbox := sb.osSbox
	sb.osSbox = nil
	sb.Unlock()

	if osSbox == nil ***REMOVED***
		return
	***REMOVED***

	for _, ep := range sb.getConnectedEndpoints() ***REMOVED***
		releaseOSSboxResources(osSbox, ep)
	***REMOVED***

	osSbox.Destroy()
***REMOVED***

func (sb *sandbox) restoreOslSandbox() error ***REMOVED***
	var routes []*types.StaticRoute

	// restore osl sandbox
	Ifaces := make(map[string][]osl.IfaceOption)
	for _, ep := range sb.endpoints ***REMOVED***
		var ifaceOptions []osl.IfaceOption
		ep.Lock()
		joinInfo := ep.joinInfo
		i := ep.iface
		ep.Unlock()

		if i == nil ***REMOVED***
			logrus.Errorf("error restoring endpoint %s for container %s", ep.Name(), sb.ContainerID())
			continue
		***REMOVED***

		ifaceOptions = append(ifaceOptions, sb.osSbox.InterfaceOptions().Address(i.addr), sb.osSbox.InterfaceOptions().Routes(i.routes))
		if i.addrv6 != nil && i.addrv6.IP.To16() != nil ***REMOVED***
			ifaceOptions = append(ifaceOptions, sb.osSbox.InterfaceOptions().AddressIPv6(i.addrv6))
		***REMOVED***
		if i.mac != nil ***REMOVED***
			ifaceOptions = append(ifaceOptions, sb.osSbox.InterfaceOptions().MacAddress(i.mac))
		***REMOVED***
		if len(i.llAddrs) != 0 ***REMOVED***
			ifaceOptions = append(ifaceOptions, sb.osSbox.InterfaceOptions().LinkLocalAddresses(i.llAddrs))
		***REMOVED***
		Ifaces[fmt.Sprintf("%s+%s", i.srcName, i.dstPrefix)] = ifaceOptions
		if joinInfo != nil ***REMOVED***
			routes = append(routes, joinInfo.StaticRoutes...)
		***REMOVED***
		if ep.needResolver() ***REMOVED***
			sb.startResolver(true)
		***REMOVED***
	***REMOVED***

	gwep := sb.getGatewayEndpoint()
	if gwep == nil ***REMOVED***
		return nil
	***REMOVED***

	// restore osl sandbox
	err := sb.osSbox.Restore(Ifaces, routes, gwep.joinInfo.gw, gwep.joinInfo.gw6)
	return err
***REMOVED***

func (sb *sandbox) populateNetworkResources(ep *endpoint) error ***REMOVED***
	sb.Lock()
	if sb.osSbox == nil ***REMOVED***
		sb.Unlock()
		return nil
	***REMOVED***
	inDelete := sb.inDelete
	sb.Unlock()

	ep.Lock()
	joinInfo := ep.joinInfo
	i := ep.iface
	ep.Unlock()

	if ep.needResolver() ***REMOVED***
		sb.startResolver(false)
	***REMOVED***

	if i != nil && i.srcName != "" ***REMOVED***
		var ifaceOptions []osl.IfaceOption

		ifaceOptions = append(ifaceOptions, sb.osSbox.InterfaceOptions().Address(i.addr), sb.osSbox.InterfaceOptions().Routes(i.routes))
		if i.addrv6 != nil && i.addrv6.IP.To16() != nil ***REMOVED***
			ifaceOptions = append(ifaceOptions, sb.osSbox.InterfaceOptions().AddressIPv6(i.addrv6))
		***REMOVED***
		if len(i.llAddrs) != 0 ***REMOVED***
			ifaceOptions = append(ifaceOptions, sb.osSbox.InterfaceOptions().LinkLocalAddresses(i.llAddrs))
		***REMOVED***
		if i.mac != nil ***REMOVED***
			ifaceOptions = append(ifaceOptions, sb.osSbox.InterfaceOptions().MacAddress(i.mac))
		***REMOVED***

		if err := sb.osSbox.AddInterface(i.srcName, i.dstPrefix, ifaceOptions...); err != nil ***REMOVED***
			return fmt.Errorf("failed to add interface %s to sandbox: %v", i.srcName, err)
		***REMOVED***
	***REMOVED***

	if len(ep.virtualIP) != 0 ***REMOVED***
		err := sb.osSbox.AddLoopbackAliasIP(&net.IPNet***REMOVED***IP: ep.virtualIP, Mask: net.CIDRMask(32, 32)***REMOVED***)
		if err != nil ***REMOVED***
			return fmt.Errorf("failed to add virtual IP %v: %v", ep.virtualIP, err)
		***REMOVED***
	***REMOVED***

	if joinInfo != nil ***REMOVED***
		// Set up non-interface routes.
		for _, r := range joinInfo.StaticRoutes ***REMOVED***
			if err := sb.osSbox.AddStaticRoute(r); err != nil ***REMOVED***
				return fmt.Errorf("failed to add static route %s: %v", r.Destination.String(), err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if ep == sb.getGatewayEndpoint() ***REMOVED***
		if err := sb.updateGateway(ep); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Make sure to add the endpoint to the populated endpoint set
	// before populating loadbalancers.
	sb.Lock()
	sb.populatedEndpoints[ep.ID()] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	sb.Unlock()

	// Populate load balancer only after updating all the other
	// information including gateway and other routes so that
	// loadbalancers are populated all the network state is in
	// place in the sandbox.
	sb.populateLoadbalancers(ep)

	// Only update the store if we did not come here as part of
	// sandbox delete. If we came here as part of delete then do
	// not bother updating the store. The sandbox object will be
	// deleted anyway
	if !inDelete ***REMOVED***
		return sb.storeUpdate()
	***REMOVED***

	return nil
***REMOVED***

func (sb *sandbox) clearNetworkResources(origEp *endpoint) error ***REMOVED***
	ep := sb.getEndpoint(origEp.id)
	if ep == nil ***REMOVED***
		return fmt.Errorf("could not find the sandbox endpoint data for endpoint %s",
			origEp.id)
	***REMOVED***

	sb.Lock()
	osSbox := sb.osSbox
	inDelete := sb.inDelete
	sb.Unlock()
	if osSbox != nil ***REMOVED***
		releaseOSSboxResources(osSbox, ep)
	***REMOVED***

	sb.Lock()
	delete(sb.populatedEndpoints, ep.ID())

	if len(sb.endpoints) == 0 ***REMOVED***
		// sb.endpoints should never be empty and this is unexpected error condition
		// We log an error message to note this down for debugging purposes.
		logrus.Errorf("No endpoints in sandbox while trying to remove endpoint %s", ep.Name())
		sb.Unlock()
		return nil
	***REMOVED***

	var (
		gwepBefore, gwepAfter *endpoint
		index                 = -1
	)
	for i, e := range sb.endpoints ***REMOVED***
		if e == ep ***REMOVED***
			index = i
		***REMOVED***
		if len(e.Gateway()) > 0 && gwepBefore == nil ***REMOVED***
			gwepBefore = e
		***REMOVED***
		if index != -1 && gwepBefore != nil ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	if index == -1 ***REMOVED***
		logrus.Warnf("Endpoint %s has already been deleted", ep.Name())
		sb.Unlock()
		return nil
	***REMOVED***

	heap.Remove(&sb.endpoints, index)
	for _, e := range sb.endpoints ***REMOVED***
		if len(e.Gateway()) > 0 ***REMOVED***
			gwepAfter = e
			break
		***REMOVED***
	***REMOVED***
	delete(sb.epPriority, ep.ID())
	sb.Unlock()

	if gwepAfter != nil && gwepBefore != gwepAfter ***REMOVED***
		sb.updateGateway(gwepAfter)
	***REMOVED***

	// Only update the store if we did not come here as part of
	// sandbox delete. If we came here as part of delete then do
	// not bother updating the store. The sandbox object will be
	// deleted anyway
	if !inDelete ***REMOVED***
		return sb.storeUpdate()
	***REMOVED***

	return nil
***REMOVED***

func (sb *sandbox) isEndpointPopulated(ep *endpoint) bool ***REMOVED***
	sb.Lock()
	_, ok := sb.populatedEndpoints[ep.ID()]
	sb.Unlock()
	return ok
***REMOVED***

// joinLeaveStart waits to ensure there are no joins or leaves in progress and
// marks this join/leave in progress without race
func (sb *sandbox) joinLeaveStart() ***REMOVED***
	sb.Lock()
	defer sb.Unlock()

	for sb.joinLeaveDone != nil ***REMOVED***
		joinLeaveDone := sb.joinLeaveDone
		sb.Unlock()

		<-joinLeaveDone

		sb.Lock()
	***REMOVED***

	sb.joinLeaveDone = make(chan struct***REMOVED******REMOVED***)
***REMOVED***

// joinLeaveEnd marks the end of this join/leave operation and
// signals the same without race to other join and leave waiters
func (sb *sandbox) joinLeaveEnd() ***REMOVED***
	sb.Lock()
	defer sb.Unlock()

	if sb.joinLeaveDone != nil ***REMOVED***
		close(sb.joinLeaveDone)
		sb.joinLeaveDone = nil
	***REMOVED***
***REMOVED***

func (sb *sandbox) hasPortConfigs() bool ***REMOVED***
	opts := sb.Labels()
	_, hasExpPorts := opts[netlabel.ExposedPorts]
	_, hasPortMaps := opts[netlabel.PortMap]
	return hasExpPorts || hasPortMaps
***REMOVED***

// OptionHostname function returns an option setter for hostname option to
// be passed to NewSandbox method.
func OptionHostname(name string) SandboxOption ***REMOVED***
	return func(sb *sandbox) ***REMOVED***
		sb.config.hostName = name
	***REMOVED***
***REMOVED***

// OptionDomainname function returns an option setter for domainname option to
// be passed to NewSandbox method.
func OptionDomainname(name string) SandboxOption ***REMOVED***
	return func(sb *sandbox) ***REMOVED***
		sb.config.domainName = name
	***REMOVED***
***REMOVED***

// OptionHostsPath function returns an option setter for hostspath option to
// be passed to NewSandbox method.
func OptionHostsPath(path string) SandboxOption ***REMOVED***
	return func(sb *sandbox) ***REMOVED***
		sb.config.hostsPath = path
	***REMOVED***
***REMOVED***

// OptionOriginHostsPath function returns an option setter for origin hosts file path
// to be passed to NewSandbox method.
func OptionOriginHostsPath(path string) SandboxOption ***REMOVED***
	return func(sb *sandbox) ***REMOVED***
		sb.config.originHostsPath = path
	***REMOVED***
***REMOVED***

// OptionExtraHost function returns an option setter for extra /etc/hosts options
// which is a name and IP as strings.
func OptionExtraHost(name string, IP string) SandboxOption ***REMOVED***
	return func(sb *sandbox) ***REMOVED***
		sb.config.extraHosts = append(sb.config.extraHosts, extraHost***REMOVED***name: name, IP: IP***REMOVED***)
	***REMOVED***
***REMOVED***

// OptionParentUpdate function returns an option setter for parent container
// which needs to update the IP address for the linked container.
func OptionParentUpdate(cid string, name, ip string) SandboxOption ***REMOVED***
	return func(sb *sandbox) ***REMOVED***
		sb.config.parentUpdates = append(sb.config.parentUpdates, parentUpdate***REMOVED***cid: cid, name: name, ip: ip***REMOVED***)
	***REMOVED***
***REMOVED***

// OptionResolvConfPath function returns an option setter for resolvconfpath option to
// be passed to net container methods.
func OptionResolvConfPath(path string) SandboxOption ***REMOVED***
	return func(sb *sandbox) ***REMOVED***
		sb.config.resolvConfPath = path
	***REMOVED***
***REMOVED***

// OptionOriginResolvConfPath function returns an option setter to set the path to the
// origin resolv.conf file to be passed to net container methods.
func OptionOriginResolvConfPath(path string) SandboxOption ***REMOVED***
	return func(sb *sandbox) ***REMOVED***
		sb.config.originResolvConfPath = path
	***REMOVED***
***REMOVED***

// OptionDNS function returns an option setter for dns entry option to
// be passed to container Create method.
func OptionDNS(dns string) SandboxOption ***REMOVED***
	return func(sb *sandbox) ***REMOVED***
		sb.config.dnsList = append(sb.config.dnsList, dns)
	***REMOVED***
***REMOVED***

// OptionDNSSearch function returns an option setter for dns search entry option to
// be passed to container Create method.
func OptionDNSSearch(search string) SandboxOption ***REMOVED***
	return func(sb *sandbox) ***REMOVED***
		sb.config.dnsSearchList = append(sb.config.dnsSearchList, search)
	***REMOVED***
***REMOVED***

// OptionDNSOptions function returns an option setter for dns options entry option to
// be passed to container Create method.
func OptionDNSOptions(options string) SandboxOption ***REMOVED***
	return func(sb *sandbox) ***REMOVED***
		sb.config.dnsOptionsList = append(sb.config.dnsOptionsList, options)
	***REMOVED***
***REMOVED***

// OptionUseDefaultSandbox function returns an option setter for using default sandbox to
// be passed to container Create method.
func OptionUseDefaultSandbox() SandboxOption ***REMOVED***
	return func(sb *sandbox) ***REMOVED***
		sb.config.useDefaultSandBox = true
	***REMOVED***
***REMOVED***

// OptionUseExternalKey function returns an option setter for using provided namespace
// instead of creating one.
func OptionUseExternalKey() SandboxOption ***REMOVED***
	return func(sb *sandbox) ***REMOVED***
		sb.config.useExternalKey = true
	***REMOVED***
***REMOVED***

// OptionGeneric function returns an option setter for Generic configuration
// that is not managed by libNetwork but can be used by the Drivers during the call to
// net container creation method. Container Labels are a good example.
func OptionGeneric(generic map[string]interface***REMOVED******REMOVED***) SandboxOption ***REMOVED***
	return func(sb *sandbox) ***REMOVED***
		if sb.config.generic == nil ***REMOVED***
			sb.config.generic = make(map[string]interface***REMOVED******REMOVED***, len(generic))
		***REMOVED***
		for k, v := range generic ***REMOVED***
			sb.config.generic[k] = v
		***REMOVED***
	***REMOVED***
***REMOVED***

// OptionExposedPorts function returns an option setter for the container exposed
// ports option to be passed to container Create method.
func OptionExposedPorts(exposedPorts []types.TransportPort) SandboxOption ***REMOVED***
	return func(sb *sandbox) ***REMOVED***
		if sb.config.generic == nil ***REMOVED***
			sb.config.generic = make(map[string]interface***REMOVED******REMOVED***)
		***REMOVED***
		// Defensive copy
		eps := make([]types.TransportPort, len(exposedPorts))
		copy(eps, exposedPorts)
		// Store endpoint label and in generic because driver needs it
		sb.config.exposedPorts = eps
		sb.config.generic[netlabel.ExposedPorts] = eps
	***REMOVED***
***REMOVED***

// OptionPortMapping function returns an option setter for the mapping
// ports option to be passed to container Create method.
func OptionPortMapping(portBindings []types.PortBinding) SandboxOption ***REMOVED***
	return func(sb *sandbox) ***REMOVED***
		if sb.config.generic == nil ***REMOVED***
			sb.config.generic = make(map[string]interface***REMOVED******REMOVED***)
		***REMOVED***
		// Store a copy of the bindings as generic data to pass to the driver
		pbs := make([]types.PortBinding, len(portBindings))
		copy(pbs, portBindings)
		sb.config.generic[netlabel.PortMap] = pbs
	***REMOVED***
***REMOVED***

// OptionIngress function returns an option setter for marking a
// sandbox as the controller's ingress sandbox.
func OptionIngress() SandboxOption ***REMOVED***
	return func(sb *sandbox) ***REMOVED***
		sb.ingress = true
	***REMOVED***
***REMOVED***

func (eh epHeap) Len() int ***REMOVED*** return len(eh) ***REMOVED***

func (eh epHeap) Less(i, j int) bool ***REMOVED***
	var (
		cip, cjp int
		ok       bool
	)

	ci, _ := eh[i].getSandbox()
	cj, _ := eh[j].getSandbox()

	epi := eh[i]
	epj := eh[j]

	if epi.endpointInGWNetwork() ***REMOVED***
		return false
	***REMOVED***

	if epj.endpointInGWNetwork() ***REMOVED***
		return true
	***REMOVED***

	if epi.getNetwork().Internal() ***REMOVED***
		return false
	***REMOVED***

	if epj.getNetwork().Internal() ***REMOVED***
		return true
	***REMOVED***

	if epi.joinInfo != nil && epj.joinInfo != nil ***REMOVED***
		if (epi.joinInfo.gw != nil && epi.joinInfo.gw6 != nil) &&
			(epj.joinInfo.gw == nil || epj.joinInfo.gw6 == nil) ***REMOVED***
			return true
		***REMOVED***
		if (epj.joinInfo.gw != nil && epj.joinInfo.gw6 != nil) &&
			(epi.joinInfo.gw == nil || epi.joinInfo.gw6 == nil) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	if ci != nil ***REMOVED***
		cip, ok = ci.epPriority[eh[i].ID()]
		if !ok ***REMOVED***
			cip = 0
		***REMOVED***
	***REMOVED***

	if cj != nil ***REMOVED***
		cjp, ok = cj.epPriority[eh[j].ID()]
		if !ok ***REMOVED***
			cjp = 0
		***REMOVED***
	***REMOVED***

	if cip == cjp ***REMOVED***
		return eh[i].network.Name() < eh[j].network.Name()
	***REMOVED***

	return cip > cjp
***REMOVED***

func (eh epHeap) Swap(i, j int) ***REMOVED*** eh[i], eh[j] = eh[j], eh[i] ***REMOVED***

func (eh *epHeap) Push(x interface***REMOVED******REMOVED***) ***REMOVED***
	*eh = append(*eh, x.(*endpoint))
***REMOVED***

func (eh *epHeap) Pop() interface***REMOVED******REMOVED*** ***REMOVED***
	old := *eh
	n := len(old)
	x := old[n-1]
	*eh = old[0 : n-1]
	return x
***REMOVED***

func (sb *sandbox) NdotsSet() bool ***REMOVED***
	return sb.ndotsSet
***REMOVED***
