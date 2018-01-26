package remote

import (
	"errors"
	"fmt"
	"net"

	"github.com/docker/docker/pkg/plugins"
	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/discoverapi"
	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/drivers/remote/api"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

type driver struct ***REMOVED***
	endpoint    *plugins.Client
	networkType string
***REMOVED***

type maybeError interface ***REMOVED***
	GetError() string
***REMOVED***

func newDriver(name string, client *plugins.Client) driverapi.Driver ***REMOVED***
	return &driver***REMOVED***networkType: name, endpoint: client***REMOVED***
***REMOVED***

// Init makes sure a remote driver is registered when a network driver
// plugin is activated.
func Init(dc driverapi.DriverCallback, config map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	newPluginHandler := func(name string, client *plugins.Client) ***REMOVED***
		// negotiate driver capability with client
		d := newDriver(name, client)
		c, err := d.(*driver).getCapabilities()
		if err != nil ***REMOVED***
			logrus.Errorf("error getting capability for %s due to %v", name, err)
			return
		***REMOVED***
		if err = dc.RegisterDriver(name, d, *c); err != nil ***REMOVED***
			logrus.Errorf("error registering driver for %s due to %v", name, err)
		***REMOVED***
	***REMOVED***

	// Unit test code is unaware of a true PluginStore. So we fall back to v1 plugins.
	handleFunc := plugins.Handle
	if pg := dc.GetPluginGetter(); pg != nil ***REMOVED***
		handleFunc = pg.Handle
		activePlugins := pg.GetAllManagedPluginsByCap(driverapi.NetworkPluginEndpointType)
		for _, ap := range activePlugins ***REMOVED***
			newPluginHandler(ap.Name(), ap.Client())
		***REMOVED***
	***REMOVED***
	handleFunc(driverapi.NetworkPluginEndpointType, newPluginHandler)

	return nil
***REMOVED***

// Get capability from client
func (d *driver) getCapabilities() (*driverapi.Capability, error) ***REMOVED***
	var capResp api.GetCapabilityResponse
	if err := d.call("GetCapabilities", nil, &capResp); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	c := &driverapi.Capability***REMOVED******REMOVED***
	switch capResp.Scope ***REMOVED***
	case "global":
		c.DataScope = datastore.GlobalScope
	case "local":
		c.DataScope = datastore.LocalScope
	default:
		return nil, fmt.Errorf("invalid capability: expecting 'local' or 'global', got %s", capResp.Scope)
	***REMOVED***

	switch capResp.ConnectivityScope ***REMOVED***
	case "global":
		c.ConnectivityScope = datastore.GlobalScope
	case "local":
		c.ConnectivityScope = datastore.LocalScope
	case "":
		c.ConnectivityScope = c.DataScope
	default:
		return nil, fmt.Errorf("invalid capability: expecting 'local' or 'global', got %s", capResp.Scope)
	***REMOVED***

	return c, nil
***REMOVED***

// Config is not implemented for remote drivers, since it is assumed
// to be supplied to the remote process out-of-band (e.g., as command
// line arguments).
func (d *driver) Config(option map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	return &driverapi.ErrNotImplemented***REMOVED******REMOVED***
***REMOVED***

func (d *driver) call(methodName string, arg interface***REMOVED******REMOVED***, retVal maybeError) error ***REMOVED***
	method := driverapi.NetworkPluginEndpointType + "." + methodName
	err := d.endpoint.Call(method, arg, retVal)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if e := retVal.GetError(); e != "" ***REMOVED***
		return fmt.Errorf("remote: %s", e)
	***REMOVED***
	return nil
***REMOVED***

func (d *driver) NetworkAllocate(id string, options map[string]string, ipV4Data, ipV6Data []driverapi.IPAMData) (map[string]string, error) ***REMOVED***
	create := &api.AllocateNetworkRequest***REMOVED***
		NetworkID: id,
		Options:   options,
		IPv4Data:  ipV4Data,
		IPv6Data:  ipV6Data,
	***REMOVED***
	retVal := api.AllocateNetworkResponse***REMOVED******REMOVED***
	err := d.call("AllocateNetwork", create, &retVal)
	return retVal.Options, err
***REMOVED***

func (d *driver) NetworkFree(id string) error ***REMOVED***
	fr := &api.FreeNetworkRequest***REMOVED***NetworkID: id***REMOVED***
	return d.call("FreeNetwork", fr, &api.FreeNetworkResponse***REMOVED******REMOVED***)
***REMOVED***

func (d *driver) EventNotify(etype driverapi.EventType, nid, tableName, key string, value []byte) ***REMOVED***
***REMOVED***

func (d *driver) DecodeTableEntry(tablename string, key string, value []byte) (string, map[string]string) ***REMOVED***
	return "", nil
***REMOVED***

func (d *driver) CreateNetwork(id string, options map[string]interface***REMOVED******REMOVED***, nInfo driverapi.NetworkInfo, ipV4Data, ipV6Data []driverapi.IPAMData) error ***REMOVED***
	create := &api.CreateNetworkRequest***REMOVED***
		NetworkID: id,
		Options:   options,
		IPv4Data:  ipV4Data,
		IPv6Data:  ipV6Data,
	***REMOVED***
	return d.call("CreateNetwork", create, &api.CreateNetworkResponse***REMOVED******REMOVED***)
***REMOVED***

func (d *driver) DeleteNetwork(nid string) error ***REMOVED***
	delete := &api.DeleteNetworkRequest***REMOVED***NetworkID: nid***REMOVED***
	return d.call("DeleteNetwork", delete, &api.DeleteNetworkResponse***REMOVED******REMOVED***)
***REMOVED***

func (d *driver) CreateEndpoint(nid, eid string, ifInfo driverapi.InterfaceInfo, epOptions map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	if ifInfo == nil ***REMOVED***
		return errors.New("must not be called with nil InterfaceInfo")
	***REMOVED***

	reqIface := &api.EndpointInterface***REMOVED******REMOVED***
	if ifInfo.Address() != nil ***REMOVED***
		reqIface.Address = ifInfo.Address().String()
	***REMOVED***
	if ifInfo.AddressIPv6() != nil ***REMOVED***
		reqIface.AddressIPv6 = ifInfo.AddressIPv6().String()
	***REMOVED***
	if ifInfo.MacAddress() != nil ***REMOVED***
		reqIface.MacAddress = ifInfo.MacAddress().String()
	***REMOVED***

	create := &api.CreateEndpointRequest***REMOVED***
		NetworkID:  nid,
		EndpointID: eid,
		Interface:  reqIface,
		Options:    epOptions,
	***REMOVED***
	var res api.CreateEndpointResponse
	if err := d.call("CreateEndpoint", create, &res); err != nil ***REMOVED***
		return err
	***REMOVED***

	inIface, err := parseInterface(res)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if inIface == nil ***REMOVED***
		// Remote driver did not set any field
		return nil
	***REMOVED***

	if inIface.MacAddress != nil ***REMOVED***
		if err := ifInfo.SetMacAddress(inIface.MacAddress); err != nil ***REMOVED***
			return errorWithRollback(fmt.Sprintf("driver modified interface MAC address: %v", err), d.DeleteEndpoint(nid, eid))
		***REMOVED***
	***REMOVED***
	if inIface.Address != nil ***REMOVED***
		if err := ifInfo.SetIPAddress(inIface.Address); err != nil ***REMOVED***
			return errorWithRollback(fmt.Sprintf("driver modified interface address: %v", err), d.DeleteEndpoint(nid, eid))
		***REMOVED***
	***REMOVED***
	if inIface.AddressIPv6 != nil ***REMOVED***
		if err := ifInfo.SetIPAddress(inIface.AddressIPv6); err != nil ***REMOVED***
			return errorWithRollback(fmt.Sprintf("driver modified interface address: %v", err), d.DeleteEndpoint(nid, eid))
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func errorWithRollback(msg string, err error) error ***REMOVED***
	rollback := "rolled back"
	if err != nil ***REMOVED***
		rollback = "failed to roll back: " + err.Error()
	***REMOVED***
	return fmt.Errorf("%s; %s", msg, rollback)
***REMOVED***

func (d *driver) DeleteEndpoint(nid, eid string) error ***REMOVED***
	delete := &api.DeleteEndpointRequest***REMOVED***
		NetworkID:  nid,
		EndpointID: eid,
	***REMOVED***
	return d.call("DeleteEndpoint", delete, &api.DeleteEndpointResponse***REMOVED******REMOVED***)
***REMOVED***

func (d *driver) EndpointOperInfo(nid, eid string) (map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	info := &api.EndpointInfoRequest***REMOVED***
		NetworkID:  nid,
		EndpointID: eid,
	***REMOVED***
	var res api.EndpointInfoResponse
	if err := d.call("EndpointOperInfo", info, &res); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return res.Value, nil
***REMOVED***

// Join method is invoked when a Sandbox is attached to an endpoint.
func (d *driver) Join(nid, eid string, sboxKey string, jinfo driverapi.JoinInfo, options map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	join := &api.JoinRequest***REMOVED***
		NetworkID:  nid,
		EndpointID: eid,
		SandboxKey: sboxKey,
		Options:    options,
	***REMOVED***
	var (
		res api.JoinResponse
		err error
	)
	if err = d.call("Join", join, &res); err != nil ***REMOVED***
		return err
	***REMOVED***

	ifaceName := res.InterfaceName
	if iface := jinfo.InterfaceName(); iface != nil && ifaceName != nil ***REMOVED***
		if err := iface.SetNames(ifaceName.SrcName, ifaceName.DstPrefix); err != nil ***REMOVED***
			return errorWithRollback(fmt.Sprintf("failed to set interface name: %s", err), d.Leave(nid, eid))
		***REMOVED***
	***REMOVED***

	var addr net.IP
	if res.Gateway != "" ***REMOVED***
		if addr = net.ParseIP(res.Gateway); addr == nil ***REMOVED***
			return fmt.Errorf(`unable to parse Gateway "%s"`, res.Gateway)
		***REMOVED***
		if jinfo.SetGateway(addr) != nil ***REMOVED***
			return errorWithRollback(fmt.Sprintf("failed to set gateway: %v", addr), d.Leave(nid, eid))
		***REMOVED***
	***REMOVED***
	if res.GatewayIPv6 != "" ***REMOVED***
		if addr = net.ParseIP(res.GatewayIPv6); addr == nil ***REMOVED***
			return fmt.Errorf(`unable to parse GatewayIPv6 "%s"`, res.GatewayIPv6)
		***REMOVED***
		if jinfo.SetGatewayIPv6(addr) != nil ***REMOVED***
			return errorWithRollback(fmt.Sprintf("failed to set gateway IPv6: %v", addr), d.Leave(nid, eid))
		***REMOVED***
	***REMOVED***
	if len(res.StaticRoutes) > 0 ***REMOVED***
		routes, err := parseStaticRoutes(res)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		for _, route := range routes ***REMOVED***
			if jinfo.AddStaticRoute(route.Destination, route.RouteType, route.NextHop) != nil ***REMOVED***
				return errorWithRollback(fmt.Sprintf("failed to set static route: %v", route), d.Leave(nid, eid))
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if res.DisableGatewayService ***REMOVED***
		jinfo.DisableGatewayService()
	***REMOVED***
	return nil
***REMOVED***

// Leave method is invoked when a Sandbox detaches from an endpoint.
func (d *driver) Leave(nid, eid string) error ***REMOVED***
	leave := &api.LeaveRequest***REMOVED***
		NetworkID:  nid,
		EndpointID: eid,
	***REMOVED***
	return d.call("Leave", leave, &api.LeaveResponse***REMOVED******REMOVED***)
***REMOVED***

// ProgramExternalConnectivity is invoked to program the rules to allow external connectivity for the endpoint.
func (d *driver) ProgramExternalConnectivity(nid, eid string, options map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	data := &api.ProgramExternalConnectivityRequest***REMOVED***
		NetworkID:  nid,
		EndpointID: eid,
		Options:    options,
	***REMOVED***
	err := d.call("ProgramExternalConnectivity", data, &api.ProgramExternalConnectivityResponse***REMOVED******REMOVED***)
	if err != nil && plugins.IsNotFound(err) ***REMOVED***
		// It is not mandatory yet to support this method
		return nil
	***REMOVED***
	return err
***REMOVED***

// RevokeExternalConnectivity method is invoked to remove any external connectivity programming related to the endpoint.
func (d *driver) RevokeExternalConnectivity(nid, eid string) error ***REMOVED***
	data := &api.RevokeExternalConnectivityRequest***REMOVED***
		NetworkID:  nid,
		EndpointID: eid,
	***REMOVED***
	err := d.call("RevokeExternalConnectivity", data, &api.RevokeExternalConnectivityResponse***REMOVED******REMOVED***)
	if err != nil && plugins.IsNotFound(err) ***REMOVED***
		// It is not mandatory yet to support this method
		return nil
	***REMOVED***
	return err
***REMOVED***

func (d *driver) Type() string ***REMOVED***
	return d.networkType
***REMOVED***

func (d *driver) IsBuiltIn() bool ***REMOVED***
	return false
***REMOVED***

// DiscoverNew is a notification for a new discovery event, such as a new node joining a cluster
func (d *driver) DiscoverNew(dType discoverapi.DiscoveryType, data interface***REMOVED******REMOVED***) error ***REMOVED***
	if dType != discoverapi.NodeDiscovery ***REMOVED***
		return nil
	***REMOVED***
	notif := &api.DiscoveryNotification***REMOVED***
		DiscoveryType: dType,
		DiscoveryData: data,
	***REMOVED***
	return d.call("DiscoverNew", notif, &api.DiscoveryResponse***REMOVED******REMOVED***)
***REMOVED***

// DiscoverDelete is a notification for a discovery delete event, such as a node leaving a cluster
func (d *driver) DiscoverDelete(dType discoverapi.DiscoveryType, data interface***REMOVED******REMOVED***) error ***REMOVED***
	if dType != discoverapi.NodeDiscovery ***REMOVED***
		return nil
	***REMOVED***
	notif := &api.DiscoveryNotification***REMOVED***
		DiscoveryType: dType,
		DiscoveryData: data,
	***REMOVED***
	return d.call("DiscoverDelete", notif, &api.DiscoveryResponse***REMOVED******REMOVED***)
***REMOVED***

func parseStaticRoutes(r api.JoinResponse) ([]*types.StaticRoute, error) ***REMOVED***
	var routes = make([]*types.StaticRoute, len(r.StaticRoutes))
	for i, inRoute := range r.StaticRoutes ***REMOVED***
		var err error
		outRoute := &types.StaticRoute***REMOVED***RouteType: inRoute.RouteType***REMOVED***

		if inRoute.Destination != "" ***REMOVED***
			if outRoute.Destination, err = types.ParseCIDR(inRoute.Destination); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***

		if inRoute.NextHop != "" ***REMOVED***
			outRoute.NextHop = net.ParseIP(inRoute.NextHop)
			if outRoute.NextHop == nil ***REMOVED***
				return nil, fmt.Errorf("failed to parse nexthop IP %s", inRoute.NextHop)
			***REMOVED***
		***REMOVED***

		routes[i] = outRoute
	***REMOVED***
	return routes, nil
***REMOVED***

// parseInterfaces validates all the parameters of an Interface and returns them.
func parseInterface(r api.CreateEndpointResponse) (*api.Interface, error) ***REMOVED***
	var outIf *api.Interface

	inIf := r.Interface
	if inIf != nil ***REMOVED***
		var err error
		outIf = &api.Interface***REMOVED******REMOVED***
		if inIf.Address != "" ***REMOVED***
			if outIf.Address, err = types.ParseCIDR(inIf.Address); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
		if inIf.AddressIPv6 != "" ***REMOVED***
			if outIf.AddressIPv6, err = types.ParseCIDR(inIf.AddressIPv6); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
		if inIf.MacAddress != "" ***REMOVED***
			if outIf.MacAddress, err = net.ParseMAC(inIf.MacAddress); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return outIf, nil
***REMOVED***
