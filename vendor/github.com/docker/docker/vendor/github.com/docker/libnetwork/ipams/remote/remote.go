package remote

import (
	"fmt"
	"net"

	"github.com/docker/docker/pkg/plugins"
	"github.com/docker/libnetwork/discoverapi"
	"github.com/docker/libnetwork/ipamapi"
	"github.com/docker/libnetwork/ipams/remote/api"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

type allocator struct ***REMOVED***
	endpoint *plugins.Client
	name     string
***REMOVED***

// PluginResponse is the interface for the plugin request responses
type PluginResponse interface ***REMOVED***
	IsSuccess() bool
	GetError() string
***REMOVED***

func newAllocator(name string, client *plugins.Client) ipamapi.Ipam ***REMOVED***
	a := &allocator***REMOVED***name: name, endpoint: client***REMOVED***
	return a
***REMOVED***

// Init registers a remote ipam when its plugin is activated
func Init(cb ipamapi.Callback, l, g interface***REMOVED******REMOVED***) error ***REMOVED***

	newPluginHandler := func(name string, client *plugins.Client) ***REMOVED***
		a := newAllocator(name, client)
		if cps, err := a.(*allocator).getCapabilities(); err == nil ***REMOVED***
			if err := cb.RegisterIpamDriverWithCapabilities(name, a, cps); err != nil ***REMOVED***
				logrus.Errorf("error registering remote ipam driver %s due to %v", name, err)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			logrus.Infof("remote ipam driver %s does not support capabilities", name)
			logrus.Debug(err)
			if err := cb.RegisterIpamDriver(name, a); err != nil ***REMOVED***
				logrus.Errorf("error registering remote ipam driver %s due to %v", name, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Unit test code is unaware of a true PluginStore. So we fall back to v1 plugins.
	handleFunc := plugins.Handle
	if pg := cb.GetPluginGetter(); pg != nil ***REMOVED***
		handleFunc = pg.Handle
		activePlugins := pg.GetAllManagedPluginsByCap(ipamapi.PluginEndpointType)
		for _, ap := range activePlugins ***REMOVED***
			newPluginHandler(ap.Name(), ap.Client())
		***REMOVED***
	***REMOVED***
	handleFunc(ipamapi.PluginEndpointType, newPluginHandler)
	return nil
***REMOVED***

func (a *allocator) call(methodName string, arg interface***REMOVED******REMOVED***, retVal PluginResponse) error ***REMOVED***
	method := ipamapi.PluginEndpointType + "." + methodName
	err := a.endpoint.Call(method, arg, retVal)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if !retVal.IsSuccess() ***REMOVED***
		return fmt.Errorf("remote: %s", retVal.GetError())
	***REMOVED***
	return nil
***REMOVED***

func (a *allocator) getCapabilities() (*ipamapi.Capability, error) ***REMOVED***
	var res api.GetCapabilityResponse
	if err := a.call("GetCapabilities", nil, &res); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return res.ToCapability(), nil
***REMOVED***

// GetDefaultAddressSpaces returns the local and global default address spaces
func (a *allocator) GetDefaultAddressSpaces() (string, string, error) ***REMOVED***
	res := &api.GetAddressSpacesResponse***REMOVED******REMOVED***
	if err := a.call("GetDefaultAddressSpaces", nil, res); err != nil ***REMOVED***
		return "", "", err
	***REMOVED***
	return res.LocalDefaultAddressSpace, res.GlobalDefaultAddressSpace, nil
***REMOVED***

// RequestPool requests an address pool in the specified address space
func (a *allocator) RequestPool(addressSpace, pool, subPool string, options map[string]string, v6 bool) (string, *net.IPNet, map[string]string, error) ***REMOVED***
	req := &api.RequestPoolRequest***REMOVED***AddressSpace: addressSpace, Pool: pool, SubPool: subPool, Options: options, V6: v6***REMOVED***
	res := &api.RequestPoolResponse***REMOVED******REMOVED***
	if err := a.call("RequestPool", req, res); err != nil ***REMOVED***
		return "", nil, nil, err
	***REMOVED***
	retPool, err := types.ParseCIDR(res.Pool)
	return res.PoolID, retPool, res.Data, err
***REMOVED***

// ReleasePool removes an address pool from the specified address space
func (a *allocator) ReleasePool(poolID string) error ***REMOVED***
	req := &api.ReleasePoolRequest***REMOVED***PoolID: poolID***REMOVED***
	res := &api.ReleasePoolResponse***REMOVED******REMOVED***
	return a.call("ReleasePool", req, res)
***REMOVED***

// RequestAddress requests an address from the address pool
func (a *allocator) RequestAddress(poolID string, address net.IP, options map[string]string) (*net.IPNet, map[string]string, error) ***REMOVED***
	var (
		prefAddress string
		retAddress  *net.IPNet
		err         error
	)
	if address != nil ***REMOVED***
		prefAddress = address.String()
	***REMOVED***
	req := &api.RequestAddressRequest***REMOVED***PoolID: poolID, Address: prefAddress, Options: options***REMOVED***
	res := &api.RequestAddressResponse***REMOVED******REMOVED***
	if err := a.call("RequestAddress", req, res); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	if res.Address != "" ***REMOVED***
		retAddress, err = types.ParseCIDR(res.Address)
	***REMOVED*** else ***REMOVED***
		return nil, nil, ipamapi.ErrNoIPReturned
	***REMOVED***
	return retAddress, res.Data, err
***REMOVED***

// ReleaseAddress releases the address from the specified address pool
func (a *allocator) ReleaseAddress(poolID string, address net.IP) error ***REMOVED***
	var relAddress string
	if address != nil ***REMOVED***
		relAddress = address.String()
	***REMOVED***
	req := &api.ReleaseAddressRequest***REMOVED***PoolID: poolID, Address: relAddress***REMOVED***
	res := &api.ReleaseAddressResponse***REMOVED******REMOVED***
	return a.call("ReleaseAddress", req, res)
***REMOVED***

// DiscoverNew is a notification for a new discovery event, such as a new global datastore
func (a *allocator) DiscoverNew(dType discoverapi.DiscoveryType, data interface***REMOVED******REMOVED***) error ***REMOVED***
	return nil
***REMOVED***

// DiscoverDelete is a notification for a discovery delete event, such as a node leaving a cluster
func (a *allocator) DiscoverDelete(dType discoverapi.DiscoveryType, data interface***REMOVED******REMOVED***) error ***REMOVED***
	return nil
***REMOVED***

func (a *allocator) IsBuiltIn() bool ***REMOVED***
	return false
***REMOVED***
