package overlay

//go:generate protoc -I.:../../Godeps/_workspace/src/github.com/gogo/protobuf  --gogo_out=import_path=github.com/docker/libnetwork/drivers/overlay,Mgogoproto/gogo.proto=github.com/gogo/protobuf/gogoproto:. overlay.proto

import (
	"encoding/json"
	"net"
	"sync"

	"github.com/Microsoft/hcsshim"
	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/discoverapi"
	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/netlabel"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

const (
	networkType  = "overlay"
	vethPrefix   = "veth"
	vethLen      = 7
	secureOption = "encrypted"
)

type driver struct ***REMOVED***
	config     map[string]interface***REMOVED******REMOVED***
	networks   networkTable
	store      datastore.DataStore
	localStore datastore.DataStore
	once       sync.Once
	joinOnce   sync.Once
	sync.Mutex
***REMOVED***

// Init registers a new instance of overlay driver
func Init(dc driverapi.DriverCallback, config map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	c := driverapi.Capability***REMOVED***
		DataScope:         datastore.GlobalScope,
		ConnectivityScope: datastore.GlobalScope,
	***REMOVED***

	d := &driver***REMOVED***
		networks: networkTable***REMOVED******REMOVED***,
		config:   config,
	***REMOVED***

	if data, ok := config[netlabel.GlobalKVClient]; ok ***REMOVED***
		var err error
		dsc, ok := data.(discoverapi.DatastoreConfigData)
		if !ok ***REMOVED***
			return types.InternalErrorf("incorrect data in datastore configuration: %v", data)
		***REMOVED***
		d.store, err = datastore.NewDataStoreFromConfig(dsc)
		if err != nil ***REMOVED***
			return types.InternalErrorf("failed to initialize data store: %v", err)
		***REMOVED***
	***REMOVED***

	if data, ok := config[netlabel.LocalKVClient]; ok ***REMOVED***
		var err error
		dsc, ok := data.(discoverapi.DatastoreConfigData)
		if !ok ***REMOVED***
			return types.InternalErrorf("incorrect data in datastore configuration: %v", data)
		***REMOVED***
		d.localStore, err = datastore.NewDataStoreFromConfig(dsc)
		if err != nil ***REMOVED***
			return types.InternalErrorf("failed to initialize local data store: %v", err)
		***REMOVED***
	***REMOVED***

	d.restoreHNSNetworks()

	return dc.RegisterDriver(networkType, d, c)
***REMOVED***

func (d *driver) restoreHNSNetworks() error ***REMOVED***
	logrus.Infof("Restoring existing overlay networks from HNS into docker")

	hnsresponse, err := hcsshim.HNSListNetworkRequest("GET", "", "")
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, v := range hnsresponse ***REMOVED***
		if v.Type != networkType ***REMOVED***
			continue
		***REMOVED***

		logrus.Infof("Restoring overlay network: %s", v.Name)
		n := d.convertToOverlayNetwork(&v)
		d.addNetwork(n)

		//
		// We assume that any network will be recreated on daemon restart
		// and therefore don't restore hns endpoints for now
		//
		//n.restoreNetworkEndpoints()
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) convertToOverlayNetwork(v *hcsshim.HNSNetwork) *network ***REMOVED***
	n := &network***REMOVED***
		id:              v.Name,
		hnsID:           v.Id,
		driver:          d,
		endpoints:       endpointTable***REMOVED******REMOVED***,
		subnets:         []*subnet***REMOVED******REMOVED***,
		providerAddress: v.ManagementIP,
	***REMOVED***

	for _, hnsSubnet := range v.Subnets ***REMOVED***
		vsidPolicy := &hcsshim.VsidPolicy***REMOVED******REMOVED***
		for _, policy := range hnsSubnet.Policies ***REMOVED***
			if err := json.Unmarshal([]byte(policy), &vsidPolicy); err == nil && vsidPolicy.Type == "VSID" ***REMOVED***
				break
			***REMOVED***
		***REMOVED***

		gwIP := net.ParseIP(hnsSubnet.GatewayAddress)
		localsubnet := &subnet***REMOVED***
			vni:  uint32(vsidPolicy.VSID),
			gwIP: &gwIP,
		***REMOVED***

		_, subnetIP, err := net.ParseCIDR(hnsSubnet.AddressPrefix)

		if err != nil ***REMOVED***
			logrus.Errorf("Error parsing subnet address %s ", hnsSubnet.AddressPrefix)
			continue
		***REMOVED***

		localsubnet.subnetIP = subnetIP

		n.subnets = append(n.subnets, localsubnet)
	***REMOVED***

	return n
***REMOVED***

func (d *driver) Type() string ***REMOVED***
	return networkType
***REMOVED***

func (d *driver) IsBuiltIn() bool ***REMOVED***
	return true
***REMOVED***

// DiscoverNew is a notification for a new discovery event, such as a new node joining a cluster
func (d *driver) DiscoverNew(dType discoverapi.DiscoveryType, data interface***REMOVED******REMOVED***) error ***REMOVED***
	return types.NotImplementedErrorf("not implemented")
***REMOVED***

// DiscoverDelete is a notification for a discovery delete event, such as a node leaving a cluster
func (d *driver) DiscoverDelete(dType discoverapi.DiscoveryType, data interface***REMOVED******REMOVED***) error ***REMOVED***
	return types.NotImplementedErrorf("not implemented")
***REMOVED***
