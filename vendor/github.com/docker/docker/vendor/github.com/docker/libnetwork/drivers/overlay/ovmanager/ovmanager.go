package ovmanager

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/discoverapi"
	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/idm"
	"github.com/docker/libnetwork/netlabel"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

const (
	networkType  = "overlay"
	vxlanIDStart = 4096
	vxlanIDEnd   = (1 << 24) - 1
)

type networkTable map[string]*network

type driver struct ***REMOVED***
	config   map[string]interface***REMOVED******REMOVED***
	networks networkTable
	store    datastore.DataStore
	vxlanIdm *idm.Idm
	sync.Mutex
***REMOVED***

type subnet struct ***REMOVED***
	subnetIP *net.IPNet
	gwIP     *net.IPNet
	vni      uint32
***REMOVED***

type network struct ***REMOVED***
	id      string
	driver  *driver
	subnets []*subnet
	sync.Mutex
***REMOVED***

// Init registers a new instance of overlay driver
func Init(dc driverapi.DriverCallback, config map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	var err error
	c := driverapi.Capability***REMOVED***
		DataScope:         datastore.GlobalScope,
		ConnectivityScope: datastore.GlobalScope,
	***REMOVED***

	d := &driver***REMOVED***
		networks: networkTable***REMOVED******REMOVED***,
		config:   config,
	***REMOVED***

	d.vxlanIdm, err = idm.New(nil, "vxlan-id", 0, vxlanIDEnd)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to initialize vxlan id manager: %v", err)
	***REMOVED***

	return dc.RegisterDriver(networkType, d, c)
***REMOVED***

func (d *driver) NetworkAllocate(id string, option map[string]string, ipV4Data, ipV6Data []driverapi.IPAMData) (map[string]string, error) ***REMOVED***
	if id == "" ***REMOVED***
		return nil, fmt.Errorf("invalid network id for overlay network")
	***REMOVED***

	if ipV4Data == nil ***REMOVED***
		return nil, fmt.Errorf("empty ipv4 data passed during overlay network creation")
	***REMOVED***

	n := &network***REMOVED***
		id:      id,
		driver:  d,
		subnets: []*subnet***REMOVED******REMOVED***,
	***REMOVED***

	opts := make(map[string]string)
	vxlanIDList := make([]uint32, 0, len(ipV4Data))
	for key, val := range option ***REMOVED***
		if key == netlabel.OverlayVxlanIDList ***REMOVED***
			logrus.Debugf("overlay network option: %s", val)
			valStrList := strings.Split(val, ",")
			for _, idStr := range valStrList ***REMOVED***
				vni, err := strconv.Atoi(idStr)
				if err != nil ***REMOVED***
					return nil, fmt.Errorf("invalid vxlan id value %q passed", idStr)
				***REMOVED***

				vxlanIDList = append(vxlanIDList, uint32(vni))
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			opts[key] = val
		***REMOVED***
	***REMOVED***

	for i, ipd := range ipV4Data ***REMOVED***
		s := &subnet***REMOVED***
			subnetIP: ipd.Pool,
			gwIP:     ipd.Gateway,
		***REMOVED***

		if len(vxlanIDList) > i ***REMOVED***
			s.vni = vxlanIDList[i]
		***REMOVED***

		if err := n.obtainVxlanID(s); err != nil ***REMOVED***
			n.releaseVxlanID()
			return nil, fmt.Errorf("could not obtain vxlan id for pool %s: %v", s.subnetIP, err)
		***REMOVED***

		n.subnets = append(n.subnets, s)
	***REMOVED***

	val := fmt.Sprintf("%d", n.subnets[0].vni)
	for _, s := range n.subnets[1:] ***REMOVED***
		val = val + fmt.Sprintf(",%d", s.vni)
	***REMOVED***
	opts[netlabel.OverlayVxlanIDList] = val

	d.Lock()
	d.networks[id] = n
	d.Unlock()

	return opts, nil
***REMOVED***

func (d *driver) NetworkFree(id string) error ***REMOVED***
	if id == "" ***REMOVED***
		return fmt.Errorf("invalid network id passed while freeing overlay network")
	***REMOVED***

	d.Lock()
	n, ok := d.networks[id]
	d.Unlock()

	if !ok ***REMOVED***
		return fmt.Errorf("overlay network with id %s not found", id)
	***REMOVED***

	// Release all vxlan IDs in one shot.
	n.releaseVxlanID()

	d.Lock()
	delete(d.networks, id)
	d.Unlock()

	return nil
***REMOVED***

func (n *network) obtainVxlanID(s *subnet) error ***REMOVED***
	var (
		err error
		vni uint64
	)

	n.Lock()
	vni = uint64(s.vni)
	n.Unlock()

	if vni == 0 ***REMOVED***
		vni, err = n.driver.vxlanIdm.GetIDInRange(vxlanIDStart, vxlanIDEnd, true)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		n.Lock()
		s.vni = uint32(vni)
		n.Unlock()
		return nil
	***REMOVED***

	return n.driver.vxlanIdm.GetSpecificID(vni)
***REMOVED***

func (n *network) releaseVxlanID() ***REMOVED***
	n.Lock()
	vnis := make([]uint32, 0, len(n.subnets))
	for _, s := range n.subnets ***REMOVED***
		vnis = append(vnis, s.vni)
		s.vni = 0
	***REMOVED***
	n.Unlock()

	for _, vni := range vnis ***REMOVED***
		n.driver.vxlanIdm.Release(uint64(vni))
	***REMOVED***
***REMOVED***

func (d *driver) CreateNetwork(id string, option map[string]interface***REMOVED******REMOVED***, nInfo driverapi.NetworkInfo, ipV4Data, ipV6Data []driverapi.IPAMData) error ***REMOVED***
	return types.NotImplementedErrorf("not implemented")
***REMOVED***

func (d *driver) EventNotify(etype driverapi.EventType, nid, tableName, key string, value []byte) ***REMOVED***
***REMOVED***

func (d *driver) DecodeTableEntry(tablename string, key string, value []byte) (string, map[string]string) ***REMOVED***
	return "", nil
***REMOVED***

func (d *driver) DeleteNetwork(nid string) error ***REMOVED***
	return types.NotImplementedErrorf("not implemented")
***REMOVED***

func (d *driver) CreateEndpoint(nid, eid string, ifInfo driverapi.InterfaceInfo, epOptions map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	return types.NotImplementedErrorf("not implemented")
***REMOVED***

func (d *driver) DeleteEndpoint(nid, eid string) error ***REMOVED***
	return types.NotImplementedErrorf("not implemented")
***REMOVED***

func (d *driver) EndpointOperInfo(nid, eid string) (map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	return nil, types.NotImplementedErrorf("not implemented")
***REMOVED***

// Join method is invoked when a Sandbox is attached to an endpoint.
func (d *driver) Join(nid, eid string, sboxKey string, jinfo driverapi.JoinInfo, options map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	return types.NotImplementedErrorf("not implemented")
***REMOVED***

// Leave method is invoked when a Sandbox detaches from an endpoint.
func (d *driver) Leave(nid, eid string) error ***REMOVED***
	return types.NotImplementedErrorf("not implemented")
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

func (d *driver) ProgramExternalConnectivity(nid, eid string, options map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	return types.NotImplementedErrorf("not implemented")
***REMOVED***

func (d *driver) RevokeExternalConnectivity(nid, eid string) error ***REMOVED***
	return types.NotImplementedErrorf("not implemented")
***REMOVED***
