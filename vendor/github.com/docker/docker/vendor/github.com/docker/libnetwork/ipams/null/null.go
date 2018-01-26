// Package null implements the null ipam driver. Null ipam driver satisfies ipamapi contract,
// but does not effectively reserve/allocate any address pool or address
package null

import (
	"fmt"
	"net"

	"github.com/docker/libnetwork/discoverapi"
	"github.com/docker/libnetwork/ipamapi"
	"github.com/docker/libnetwork/types"
)

var (
	defaultAS      = "null"
	defaultPool, _ = types.ParseCIDR("0.0.0.0/0")
	defaultPoolID  = fmt.Sprintf("%s/%s", defaultAS, defaultPool.String())
)

type allocator struct***REMOVED******REMOVED***

func (a *allocator) GetDefaultAddressSpaces() (string, string, error) ***REMOVED***
	return defaultAS, defaultAS, nil
***REMOVED***

func (a *allocator) RequestPool(addressSpace, pool, subPool string, options map[string]string, v6 bool) (string, *net.IPNet, map[string]string, error) ***REMOVED***
	if addressSpace != defaultAS ***REMOVED***
		return "", nil, nil, types.BadRequestErrorf("unknown address space: %s", addressSpace)
	***REMOVED***
	if pool != "" ***REMOVED***
		return "", nil, nil, types.BadRequestErrorf("null ipam driver does not handle specific address pool requests")
	***REMOVED***
	if subPool != "" ***REMOVED***
		return "", nil, nil, types.BadRequestErrorf("null ipam driver does not handle specific address subpool requests")
	***REMOVED***
	if v6 ***REMOVED***
		return "", nil, nil, types.BadRequestErrorf("null ipam driver does not handle IPv6 address pool pool requests")
	***REMOVED***
	return defaultPoolID, defaultPool, nil, nil
***REMOVED***

func (a *allocator) ReleasePool(poolID string) error ***REMOVED***
	return nil
***REMOVED***

func (a *allocator) RequestAddress(poolID string, ip net.IP, opts map[string]string) (*net.IPNet, map[string]string, error) ***REMOVED***
	if poolID != defaultPoolID ***REMOVED***
		return nil, nil, types.BadRequestErrorf("unknown pool id: %s", poolID)
	***REMOVED***
	return nil, nil, nil
***REMOVED***

func (a *allocator) ReleaseAddress(poolID string, ip net.IP) error ***REMOVED***
	if poolID != defaultPoolID ***REMOVED***
		return types.BadRequestErrorf("unknown pool id: %s", poolID)
	***REMOVED***
	return nil
***REMOVED***

func (a *allocator) DiscoverNew(dType discoverapi.DiscoveryType, data interface***REMOVED******REMOVED***) error ***REMOVED***
	return nil
***REMOVED***

func (a *allocator) DiscoverDelete(dType discoverapi.DiscoveryType, data interface***REMOVED******REMOVED***) error ***REMOVED***
	return nil
***REMOVED***

func (a *allocator) IsBuiltIn() bool ***REMOVED***
	return true
***REMOVED***

// Init registers a remote ipam when its plugin is activated
func Init(ic ipamapi.Callback, l, g interface***REMOVED******REMOVED***) error ***REMOVED***
	return ic.RegisterIpamDriver(ipamapi.NullIPAM, &allocator***REMOVED******REMOVED***)
***REMOVED***
