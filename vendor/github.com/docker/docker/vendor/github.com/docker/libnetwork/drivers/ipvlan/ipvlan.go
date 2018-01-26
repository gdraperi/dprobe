package ipvlan

import (
	"net"
	"sync"

	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/discoverapi"
	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/osl"
	"github.com/docker/libnetwork/types"
)

const (
	vethLen             = 7
	containerVethPrefix = "eth"
	vethPrefix          = "veth"
	ipvlanType          = "ipvlan" // driver type name
	modeL2              = "l2"     // ipvlan mode l2 is the default
	modeL3              = "l3"     // ipvlan L3 mode
	parentOpt           = "parent" // parent interface -o parent
	modeOpt             = "_mode"  // ipvlan mode ux opt suffix
)

var driverModeOpt = ipvlanType + modeOpt // mode -o ipvlan_mode

type endpointTable map[string]*endpoint

type networkTable map[string]*network

type driver struct ***REMOVED***
	networks networkTable
	sync.Once
	sync.Mutex
	store datastore.DataStore
***REMOVED***

type endpoint struct ***REMOVED***
	id       string
	nid      string
	mac      net.HardwareAddr
	addr     *net.IPNet
	addrv6   *net.IPNet
	srcName  string
	dbIndex  uint64
	dbExists bool
***REMOVED***

type network struct ***REMOVED***
	id        string
	sbox      osl.Sandbox
	endpoints endpointTable
	driver    *driver
	config    *configuration
	sync.Mutex
***REMOVED***

// Init initializes and registers the libnetwork ipvlan driver
func Init(dc driverapi.DriverCallback, config map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	c := driverapi.Capability***REMOVED***
		DataScope:         datastore.LocalScope,
		ConnectivityScope: datastore.GlobalScope,
	***REMOVED***
	d := &driver***REMOVED***
		networks: networkTable***REMOVED******REMOVED***,
	***REMOVED***
	d.initStore(config)

	return dc.RegisterDriver(ipvlanType, d, c)
***REMOVED***

func (d *driver) NetworkAllocate(id string, option map[string]string, ipV4Data, ipV6Data []driverapi.IPAMData) (map[string]string, error) ***REMOVED***
	return nil, types.NotImplementedErrorf("not implemented")
***REMOVED***

func (d *driver) NetworkFree(id string) error ***REMOVED***
	return types.NotImplementedErrorf("not implemented")
***REMOVED***

func (d *driver) EndpointOperInfo(nid, eid string) (map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	return make(map[string]interface***REMOVED******REMOVED***, 0), nil
***REMOVED***

func (d *driver) Type() string ***REMOVED***
	return ipvlanType
***REMOVED***

func (d *driver) IsBuiltIn() bool ***REMOVED***
	return true
***REMOVED***

func (d *driver) ProgramExternalConnectivity(nid, eid string, options map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	return nil
***REMOVED***

func (d *driver) RevokeExternalConnectivity(nid, eid string) error ***REMOVED***
	return nil
***REMOVED***

// DiscoverNew is a notification for a new discovery event.
func (d *driver) DiscoverNew(dType discoverapi.DiscoveryType, data interface***REMOVED******REMOVED***) error ***REMOVED***
	return nil
***REMOVED***

// DiscoverDelete is a notification for a discovery delete event.
func (d *driver) DiscoverDelete(dType discoverapi.DiscoveryType, data interface***REMOVED******REMOVED***) error ***REMOVED***
	return nil
***REMOVED***

func (d *driver) EventNotify(etype driverapi.EventType, nid, tableName, key string, value []byte) ***REMOVED***
***REMOVED***

func (d *driver) DecodeTableEntry(tablename string, key string, value []byte) (string, map[string]string) ***REMOVED***
	return "", nil
***REMOVED***
