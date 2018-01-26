package host

import (
	"sync"

	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/discoverapi"
	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/types"
)

const networkType = "host"

type driver struct ***REMOVED***
	network string
	sync.Mutex
***REMOVED***

// Init registers a new instance of host driver
func Init(dc driverapi.DriverCallback, config map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	c := driverapi.Capability***REMOVED***
		DataScope:         datastore.LocalScope,
		ConnectivityScope: datastore.LocalScope,
	***REMOVED***
	return dc.RegisterDriver(networkType, &driver***REMOVED******REMOVED***, c)
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

func (d *driver) CreateNetwork(id string, option map[string]interface***REMOVED******REMOVED***, nInfo driverapi.NetworkInfo, ipV4Data, ipV6Data []driverapi.IPAMData) error ***REMOVED***
	d.Lock()
	defer d.Unlock()

	if d.network != "" ***REMOVED***
		return types.ForbiddenErrorf("only one instance of \"%s\" network is allowed", networkType)
	***REMOVED***

	d.network = id

	return nil
***REMOVED***

func (d *driver) DeleteNetwork(nid string) error ***REMOVED***
	return types.ForbiddenErrorf("network of type \"%s\" cannot be deleted", networkType)
***REMOVED***

func (d *driver) CreateEndpoint(nid, eid string, ifInfo driverapi.InterfaceInfo, epOptions map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	return nil
***REMOVED***

func (d *driver) DeleteEndpoint(nid, eid string) error ***REMOVED***
	return nil
***REMOVED***

func (d *driver) EndpointOperInfo(nid, eid string) (map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	return make(map[string]interface***REMOVED******REMOVED***, 0), nil
***REMOVED***

// Join method is invoked when a Sandbox is attached to an endpoint.
func (d *driver) Join(nid, eid string, sboxKey string, jinfo driverapi.JoinInfo, options map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	return nil
***REMOVED***

// Leave method is invoked when a Sandbox detaches from an endpoint.
func (d *driver) Leave(nid, eid string) error ***REMOVED***
	return nil
***REMOVED***

func (d *driver) ProgramExternalConnectivity(nid, eid string, options map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	return nil
***REMOVED***

func (d *driver) RevokeExternalConnectivity(nid, eid string) error ***REMOVED***
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
