package windowsipam

import (
	"net"

	"github.com/docker/libnetwork/discoverapi"
	"github.com/docker/libnetwork/ipamapi"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

const (
	localAddressSpace  = "LocalDefault"
	globalAddressSpace = "GlobalDefault"
)

// DefaultIPAM defines the default ipam-driver for local-scoped windows networks
const DefaultIPAM = "windows"

var (
	defaultPool, _ = types.ParseCIDR("0.0.0.0/0")
)

type allocator struct ***REMOVED***
***REMOVED***

// GetInit registers the built-in ipam service with libnetwork
func GetInit(ipamName string) func(ic ipamapi.Callback, l, g interface***REMOVED******REMOVED***) error ***REMOVED***
	return func(ic ipamapi.Callback, l, g interface***REMOVED******REMOVED***) error ***REMOVED***
		return ic.RegisterIpamDriver(ipamName, &allocator***REMOVED******REMOVED***)
	***REMOVED***
***REMOVED***

func (a *allocator) GetDefaultAddressSpaces() (string, string, error) ***REMOVED***
	return localAddressSpace, globalAddressSpace, nil
***REMOVED***

// RequestPool returns an address pool along with its unique id. This is a null ipam driver. It allocates the
// subnet user asked and does not validate anything. Doesn't support subpool allocation
func (a *allocator) RequestPool(addressSpace, pool, subPool string, options map[string]string, v6 bool) (string, *net.IPNet, map[string]string, error) ***REMOVED***
	logrus.Debugf("RequestPool(%s, %s, %s, %v, %t)", addressSpace, pool, subPool, options, v6)
	if subPool != "" || v6 ***REMOVED***
		return "", nil, nil, types.InternalErrorf("This request is not supported by null ipam driver")
	***REMOVED***

	var ipNet *net.IPNet
	var err error

	if pool != "" ***REMOVED***
		_, ipNet, err = net.ParseCIDR(pool)
		if err != nil ***REMOVED***
			return "", nil, nil, err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		ipNet = defaultPool
	***REMOVED***

	return ipNet.String(), ipNet, nil, nil
***REMOVED***

// ReleasePool releases the address pool - always succeeds
func (a *allocator) ReleasePool(poolID string) error ***REMOVED***
	logrus.Debugf("ReleasePool(%s)", poolID)
	return nil
***REMOVED***

// RequestAddress returns an address from the specified pool ID.
// Always allocate the 0.0.0.0/32 ip if no preferred address was specified
func (a *allocator) RequestAddress(poolID string, prefAddress net.IP, opts map[string]string) (*net.IPNet, map[string]string, error) ***REMOVED***
	logrus.Debugf("RequestAddress(%s, %v, %v)", poolID, prefAddress, opts)
	_, ipNet, err := net.ParseCIDR(poolID)

	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	if prefAddress != nil ***REMOVED***
		return &net.IPNet***REMOVED***IP: prefAddress, Mask: ipNet.Mask***REMOVED***, nil, nil
	***REMOVED***

	return nil, nil, nil
***REMOVED***

// ReleaseAddress releases the address - always succeeds
func (a *allocator) ReleaseAddress(poolID string, address net.IP) error ***REMOVED***
	logrus.Debugf("ReleaseAddress(%s, %v)", poolID, address)
	return nil
***REMOVED***

// DiscoverNew informs the allocator about a new global scope datastore
func (a *allocator) DiscoverNew(dType discoverapi.DiscoveryType, data interface***REMOVED******REMOVED***) error ***REMOVED***
	return nil
***REMOVED***

// DiscoverDelete is a notification of no interest for the allocator
func (a *allocator) DiscoverDelete(dType discoverapi.DiscoveryType, data interface***REMOVED******REMOVED***) error ***REMOVED***
	return nil
***REMOVED***

func (a *allocator) IsBuiltIn() bool ***REMOVED***
	return true
***REMOVED***
