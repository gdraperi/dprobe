package libnetwork

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/types"
)

// EndpointInfo provides an interface to retrieve network resources bound to the endpoint.
type EndpointInfo interface ***REMOVED***
	// Iface returns InterfaceInfo, go interface that can be used
	// to get more information on the interface which was assigned to
	// the endpoint by the driver. This can be used after the
	// endpoint has been created.
	Iface() InterfaceInfo

	// Gateway returns the IPv4 gateway assigned by the driver.
	// This will only return a valid value if a container has joined the endpoint.
	Gateway() net.IP

	// GatewayIPv6 returns the IPv6 gateway assigned by the driver.
	// This will only return a valid value if a container has joined the endpoint.
	GatewayIPv6() net.IP

	// StaticRoutes returns the list of static routes configured by the network
	// driver when the container joins a network
	StaticRoutes() []*types.StaticRoute

	// Sandbox returns the attached sandbox if there, nil otherwise.
	Sandbox() Sandbox

	// LoadBalancer returns whether the endpoint is the load balancer endpoint for the network.
	LoadBalancer() bool
***REMOVED***

// InterfaceInfo provides an interface to retrieve interface addresses bound to the endpoint.
type InterfaceInfo interface ***REMOVED***
	// MacAddress returns the MAC address assigned to the endpoint.
	MacAddress() net.HardwareAddr

	// Address returns the IPv4 address assigned to the endpoint.
	Address() *net.IPNet

	// AddressIPv6 returns the IPv6 address assigned to the endpoint.
	AddressIPv6() *net.IPNet

	// LinkLocalAddresses returns the list of link-local (IPv4/IPv6) addresses assigned to the endpoint.
	LinkLocalAddresses() []*net.IPNet
***REMOVED***

type endpointInterface struct ***REMOVED***
	mac       net.HardwareAddr
	addr      *net.IPNet
	addrv6    *net.IPNet
	llAddrs   []*net.IPNet
	srcName   string
	dstPrefix string
	routes    []*net.IPNet
	v4PoolID  string
	v6PoolID  string
***REMOVED***

func (epi *endpointInterface) MarshalJSON() ([]byte, error) ***REMOVED***
	epMap := make(map[string]interface***REMOVED******REMOVED***)
	if epi.mac != nil ***REMOVED***
		epMap["mac"] = epi.mac.String()
	***REMOVED***
	if epi.addr != nil ***REMOVED***
		epMap["addr"] = epi.addr.String()
	***REMOVED***
	if epi.addrv6 != nil ***REMOVED***
		epMap["addrv6"] = epi.addrv6.String()
	***REMOVED***
	if len(epi.llAddrs) != 0 ***REMOVED***
		list := make([]string, 0, len(epi.llAddrs))
		for _, ll := range epi.llAddrs ***REMOVED***
			list = append(list, ll.String())
		***REMOVED***
		epMap["llAddrs"] = list
	***REMOVED***
	epMap["srcName"] = epi.srcName
	epMap["dstPrefix"] = epi.dstPrefix
	var routes []string
	for _, route := range epi.routes ***REMOVED***
		routes = append(routes, route.String())
	***REMOVED***
	epMap["routes"] = routes
	epMap["v4PoolID"] = epi.v4PoolID
	epMap["v6PoolID"] = epi.v6PoolID
	return json.Marshal(epMap)
***REMOVED***

func (epi *endpointInterface) UnmarshalJSON(b []byte) error ***REMOVED***
	var (
		err   error
		epMap map[string]interface***REMOVED******REMOVED***
	)
	if err = json.Unmarshal(b, &epMap); err != nil ***REMOVED***
		return err
	***REMOVED***
	if v, ok := epMap["mac"]; ok ***REMOVED***
		if epi.mac, err = net.ParseMAC(v.(string)); err != nil ***REMOVED***
			return types.InternalErrorf("failed to decode endpoint interface mac address after json unmarshal: %s", v.(string))
		***REMOVED***
	***REMOVED***
	if v, ok := epMap["addr"]; ok ***REMOVED***
		if epi.addr, err = types.ParseCIDR(v.(string)); err != nil ***REMOVED***
			return types.InternalErrorf("failed to decode endpoint interface ipv4 address after json unmarshal: %v", err)
		***REMOVED***
	***REMOVED***
	if v, ok := epMap["addrv6"]; ok ***REMOVED***
		if epi.addrv6, err = types.ParseCIDR(v.(string)); err != nil ***REMOVED***
			return types.InternalErrorf("failed to decode endpoint interface ipv6 address after json unmarshal: %v", err)
		***REMOVED***
	***REMOVED***
	if v, ok := epMap["llAddrs"]; ok ***REMOVED***
		list := v.([]interface***REMOVED******REMOVED***)
		epi.llAddrs = make([]*net.IPNet, 0, len(list))
		for _, llS := range list ***REMOVED***
			ll, err := types.ParseCIDR(llS.(string))
			if err != nil ***REMOVED***
				return types.InternalErrorf("failed to decode endpoint interface link-local address (%v) after json unmarshal: %v", llS, err)
			***REMOVED***
			epi.llAddrs = append(epi.llAddrs, ll)
		***REMOVED***
	***REMOVED***
	epi.srcName = epMap["srcName"].(string)
	epi.dstPrefix = epMap["dstPrefix"].(string)

	rb, _ := json.Marshal(epMap["routes"])
	var routes []string
	json.Unmarshal(rb, &routes)
	epi.routes = make([]*net.IPNet, 0)
	for _, route := range routes ***REMOVED***
		ip, ipr, err := net.ParseCIDR(route)
		if err == nil ***REMOVED***
			ipr.IP = ip
			epi.routes = append(epi.routes, ipr)
		***REMOVED***
	***REMOVED***
	epi.v4PoolID = epMap["v4PoolID"].(string)
	epi.v6PoolID = epMap["v6PoolID"].(string)

	return nil
***REMOVED***

func (epi *endpointInterface) CopyTo(dstEpi *endpointInterface) error ***REMOVED***
	dstEpi.mac = types.GetMacCopy(epi.mac)
	dstEpi.addr = types.GetIPNetCopy(epi.addr)
	dstEpi.addrv6 = types.GetIPNetCopy(epi.addrv6)
	dstEpi.srcName = epi.srcName
	dstEpi.dstPrefix = epi.dstPrefix
	dstEpi.v4PoolID = epi.v4PoolID
	dstEpi.v6PoolID = epi.v6PoolID
	if len(epi.llAddrs) != 0 ***REMOVED***
		dstEpi.llAddrs = make([]*net.IPNet, 0, len(epi.llAddrs))
		dstEpi.llAddrs = append(dstEpi.llAddrs, epi.llAddrs...)
	***REMOVED***

	for _, route := range epi.routes ***REMOVED***
		dstEpi.routes = append(dstEpi.routes, types.GetIPNetCopy(route))
	***REMOVED***

	return nil
***REMOVED***

type endpointJoinInfo struct ***REMOVED***
	gw                    net.IP
	gw6                   net.IP
	StaticRoutes          []*types.StaticRoute
	driverTableEntries    []*tableEntry
	disableGatewayService bool
***REMOVED***

type tableEntry struct ***REMOVED***
	tableName string
	key       string
	value     []byte
***REMOVED***

func (ep *endpoint) Info() EndpointInfo ***REMOVED***
	if ep.sandboxID != "" ***REMOVED***
		return ep
	***REMOVED***
	n, err := ep.getNetworkFromStore()
	if err != nil ***REMOVED***
		return nil
	***REMOVED***

	ep, err = n.getEndpointFromStore(ep.ID())
	if err != nil ***REMOVED***
		return nil
	***REMOVED***

	sb, ok := ep.getSandbox()
	if !ok ***REMOVED***
		// endpoint hasn't joined any sandbox.
		// Just return the endpoint
		return ep
	***REMOVED***

	return sb.getEndpoint(ep.ID())
***REMOVED***

func (ep *endpoint) Iface() InterfaceInfo ***REMOVED***
	ep.Lock()
	defer ep.Unlock()

	if ep.iface != nil ***REMOVED***
		return ep.iface
	***REMOVED***

	return nil
***REMOVED***

func (ep *endpoint) Interface() driverapi.InterfaceInfo ***REMOVED***
	ep.Lock()
	defer ep.Unlock()

	if ep.iface != nil ***REMOVED***
		return ep.iface
	***REMOVED***

	return nil
***REMOVED***

func (epi *endpointInterface) SetMacAddress(mac net.HardwareAddr) error ***REMOVED***
	if epi.mac != nil ***REMOVED***
		return types.ForbiddenErrorf("endpoint interface MAC address present (%s). Cannot be modified with %s.", epi.mac, mac)
	***REMOVED***
	if mac == nil ***REMOVED***
		return types.BadRequestErrorf("tried to set nil MAC address to endpoint interface")
	***REMOVED***
	epi.mac = types.GetMacCopy(mac)
	return nil
***REMOVED***

func (epi *endpointInterface) SetIPAddress(address *net.IPNet) error ***REMOVED***
	if address.IP == nil ***REMOVED***
		return types.BadRequestErrorf("tried to set nil IP address to endpoint interface")
	***REMOVED***
	if address.IP.To4() == nil ***REMOVED***
		return setAddress(&epi.addrv6, address)
	***REMOVED***
	return setAddress(&epi.addr, address)
***REMOVED***

func setAddress(ifaceAddr **net.IPNet, address *net.IPNet) error ***REMOVED***
	if *ifaceAddr != nil ***REMOVED***
		return types.ForbiddenErrorf("endpoint interface IP present (%s). Cannot be modified with (%s).", *ifaceAddr, address)
	***REMOVED***
	*ifaceAddr = types.GetIPNetCopy(address)
	return nil
***REMOVED***

func (epi *endpointInterface) MacAddress() net.HardwareAddr ***REMOVED***
	return types.GetMacCopy(epi.mac)
***REMOVED***

func (epi *endpointInterface) Address() *net.IPNet ***REMOVED***
	return types.GetIPNetCopy(epi.addr)
***REMOVED***

func (epi *endpointInterface) AddressIPv6() *net.IPNet ***REMOVED***
	return types.GetIPNetCopy(epi.addrv6)
***REMOVED***

func (epi *endpointInterface) LinkLocalAddresses() []*net.IPNet ***REMOVED***
	return epi.llAddrs
***REMOVED***

func (epi *endpointInterface) SetNames(srcName string, dstPrefix string) error ***REMOVED***
	epi.srcName = srcName
	epi.dstPrefix = dstPrefix
	return nil
***REMOVED***

func (ep *endpoint) InterfaceName() driverapi.InterfaceNameInfo ***REMOVED***
	ep.Lock()
	defer ep.Unlock()

	if ep.iface != nil ***REMOVED***
		return ep.iface
	***REMOVED***

	return nil
***REMOVED***

func (ep *endpoint) AddStaticRoute(destination *net.IPNet, routeType int, nextHop net.IP) error ***REMOVED***
	ep.Lock()
	defer ep.Unlock()

	r := types.StaticRoute***REMOVED***Destination: destination, RouteType: routeType, NextHop: nextHop***REMOVED***

	if routeType == types.NEXTHOP ***REMOVED***
		// If the route specifies a next-hop, then it's loosely routed (i.e. not bound to a particular interface).
		ep.joinInfo.StaticRoutes = append(ep.joinInfo.StaticRoutes, &r)
	***REMOVED*** else ***REMOVED***
		// If the route doesn't specify a next-hop, it must be a connected route, bound to an interface.
		ep.iface.routes = append(ep.iface.routes, r.Destination)
	***REMOVED***
	return nil
***REMOVED***

func (ep *endpoint) AddTableEntry(tableName, key string, value []byte) error ***REMOVED***
	ep.Lock()
	defer ep.Unlock()

	ep.joinInfo.driverTableEntries = append(ep.joinInfo.driverTableEntries, &tableEntry***REMOVED***
		tableName: tableName,
		key:       key,
		value:     value,
	***REMOVED***)

	return nil
***REMOVED***

func (ep *endpoint) Sandbox() Sandbox ***REMOVED***
	cnt, ok := ep.getSandbox()
	if !ok ***REMOVED***
		return nil
	***REMOVED***
	return cnt
***REMOVED***

func (ep *endpoint) LoadBalancer() bool ***REMOVED***
	ep.Lock()
	defer ep.Unlock()
	return ep.loadBalancer
***REMOVED***

func (ep *endpoint) StaticRoutes() []*types.StaticRoute ***REMOVED***
	ep.Lock()
	defer ep.Unlock()

	if ep.joinInfo == nil ***REMOVED***
		return nil
	***REMOVED***

	return ep.joinInfo.StaticRoutes
***REMOVED***

func (ep *endpoint) Gateway() net.IP ***REMOVED***
	ep.Lock()
	defer ep.Unlock()

	if ep.joinInfo == nil ***REMOVED***
		return net.IP***REMOVED******REMOVED***
	***REMOVED***

	return types.GetIPCopy(ep.joinInfo.gw)
***REMOVED***

func (ep *endpoint) GatewayIPv6() net.IP ***REMOVED***
	ep.Lock()
	defer ep.Unlock()

	if ep.joinInfo == nil ***REMOVED***
		return net.IP***REMOVED******REMOVED***
	***REMOVED***

	return types.GetIPCopy(ep.joinInfo.gw6)
***REMOVED***

func (ep *endpoint) SetGateway(gw net.IP) error ***REMOVED***
	ep.Lock()
	defer ep.Unlock()

	ep.joinInfo.gw = types.GetIPCopy(gw)
	return nil
***REMOVED***

func (ep *endpoint) SetGatewayIPv6(gw6 net.IP) error ***REMOVED***
	ep.Lock()
	defer ep.Unlock()

	ep.joinInfo.gw6 = types.GetIPCopy(gw6)
	return nil
***REMOVED***

func (ep *endpoint) retrieveFromStore() (*endpoint, error) ***REMOVED***
	n, err := ep.getNetworkFromStore()
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("could not find network in store to get latest endpoint %s: %v", ep.Name(), err)
	***REMOVED***
	return n.getEndpointFromStore(ep.ID())
***REMOVED***

func (ep *endpoint) DisableGatewayService() ***REMOVED***
	ep.Lock()
	defer ep.Unlock()

	ep.joinInfo.disableGatewayService = true
***REMOVED***

func (epj *endpointJoinInfo) MarshalJSON() ([]byte, error) ***REMOVED***
	epMap := make(map[string]interface***REMOVED******REMOVED***)
	if epj.gw != nil ***REMOVED***
		epMap["gw"] = epj.gw.String()
	***REMOVED***
	if epj.gw6 != nil ***REMOVED***
		epMap["gw6"] = epj.gw6.String()
	***REMOVED***
	epMap["disableGatewayService"] = epj.disableGatewayService
	epMap["StaticRoutes"] = epj.StaticRoutes
	return json.Marshal(epMap)
***REMOVED***

func (epj *endpointJoinInfo) UnmarshalJSON(b []byte) error ***REMOVED***
	var (
		err   error
		epMap map[string]interface***REMOVED******REMOVED***
	)
	if err = json.Unmarshal(b, &epMap); err != nil ***REMOVED***
		return err
	***REMOVED***
	if v, ok := epMap["gw"]; ok ***REMOVED***
		epj.gw = net.ParseIP(v.(string))
	***REMOVED***
	if v, ok := epMap["gw6"]; ok ***REMOVED***
		epj.gw6 = net.ParseIP(v.(string))
	***REMOVED***
	epj.disableGatewayService = epMap["disableGatewayService"].(bool)

	var tStaticRoute []types.StaticRoute
	if v, ok := epMap["StaticRoutes"]; ok ***REMOVED***
		tb, _ := json.Marshal(v)
		var tStaticRoute []types.StaticRoute
		json.Unmarshal(tb, &tStaticRoute)
	***REMOVED***
	var StaticRoutes []*types.StaticRoute
	for _, r := range tStaticRoute ***REMOVED***
		StaticRoutes = append(StaticRoutes, &r)
	***REMOVED***
	epj.StaticRoutes = StaticRoutes

	return nil
***REMOVED***

func (epj *endpointJoinInfo) CopyTo(dstEpj *endpointJoinInfo) error ***REMOVED***
	dstEpj.disableGatewayService = epj.disableGatewayService
	dstEpj.StaticRoutes = make([]*types.StaticRoute, len(epj.StaticRoutes))
	copy(dstEpj.StaticRoutes, epj.StaticRoutes)
	dstEpj.driverTableEntries = make([]*tableEntry, len(epj.driverTableEntries))
	copy(dstEpj.driverTableEntries, epj.driverTableEntries)
	dstEpj.gw = types.GetIPCopy(epj.gw)
	dstEpj.gw6 = types.GetIPCopy(epj.gw6)
	return nil
***REMOVED***
