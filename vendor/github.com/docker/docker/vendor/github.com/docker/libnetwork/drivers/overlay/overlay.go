package overlay

//go:generate protoc -I.:../../Godeps/_workspace/src/github.com/gogo/protobuf  --gogo_out=import_path=github.com/docker/libnetwork/drivers/overlay,Mgogoproto/gogo.proto=github.com/gogo/protobuf/gogoproto:. overlay.proto

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/discoverapi"
	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/idm"
	"github.com/docker/libnetwork/netlabel"
	"github.com/docker/libnetwork/osl"
	"github.com/docker/libnetwork/types"
	"github.com/hashicorp/serf/serf"
	"github.com/sirupsen/logrus"
)

const (
	networkType  = "overlay"
	vethPrefix   = "veth"
	vethLen      = 7
	vxlanIDStart = 256
	vxlanIDEnd   = (1 << 24) - 1
	vxlanPort    = 4789
	vxlanEncap   = 50
	secureOption = "encrypted"
)

var initVxlanIdm = make(chan (bool), 1)

type driver struct ***REMOVED***
	eventCh          chan serf.Event
	notifyCh         chan ovNotify
	exitCh           chan chan struct***REMOVED******REMOVED***
	bindAddress      string
	advertiseAddress string
	neighIP          string
	config           map[string]interface***REMOVED******REMOVED***
	peerDb           peerNetworkMap
	secMap           *encrMap
	serfInstance     *serf.Serf
	networks         networkTable
	store            datastore.DataStore
	localStore       datastore.DataStore
	vxlanIdm         *idm.Idm
	initOS           sync.Once
	joinOnce         sync.Once
	localJoinOnce    sync.Once
	keys             []*key
	peerOpCh         chan *peerOperation
	peerOpCancel     context.CancelFunc
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
		peerDb: peerNetworkMap***REMOVED***
			mp: map[string]*peerMap***REMOVED******REMOVED***,
		***REMOVED***,
		secMap:   &encrMap***REMOVED***nodes: map[string][]*spi***REMOVED******REMOVED******REMOVED***,
		config:   config,
		peerOpCh: make(chan *peerOperation),
	***REMOVED***

	// Launch the go routine for processing peer operations
	ctx, cancel := context.WithCancel(context.Background())
	d.peerOpCancel = cancel
	go d.peerOpRoutine(ctx, d.peerOpCh)

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

	if err := d.restoreEndpoints(); err != nil ***REMOVED***
		logrus.Warnf("Failure during overlay endpoints restore: %v", err)
	***REMOVED***

	// If an error happened when the network join the sandbox during the endpoints restore
	// we should reset it now along with the once variable, so that subsequent endpoint joins
	// outside of the restore path can potentially fix the network join and succeed.
	for nid, n := range d.networks ***REMOVED***
		if n.initErr != nil ***REMOVED***
			logrus.Infof("resetting init error and once variable for network %s after unsuccessful endpoint restore: %v", nid, n.initErr)
			n.initErr = nil
			n.once = &sync.Once***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***

	return dc.RegisterDriver(networkType, d, c)
***REMOVED***

// Endpoints are stored in the local store. Restore them and reconstruct the overlay sandbox
func (d *driver) restoreEndpoints() error ***REMOVED***
	if d.localStore == nil ***REMOVED***
		logrus.Warn("Cannot restore overlay endpoints because local datastore is missing")
		return nil
	***REMOVED***
	kvol, err := d.localStore.List(datastore.Key(overlayEndpointPrefix), &endpoint***REMOVED******REMOVED***)
	if err != nil && err != datastore.ErrKeyNotFound ***REMOVED***
		return fmt.Errorf("failed to read overlay endpoint from store: %v", err)
	***REMOVED***

	if err == datastore.ErrKeyNotFound ***REMOVED***
		return nil
	***REMOVED***
	for _, kvo := range kvol ***REMOVED***
		ep := kvo.(*endpoint)
		n := d.network(ep.nid)
		if n == nil ***REMOVED***
			logrus.Debugf("Network (%s) not found for restored endpoint (%s)", ep.nid[0:7], ep.id[0:7])
			logrus.Debugf("Deleting stale overlay endpoint (%s) from store", ep.id[0:7])
			if err := d.deleteEndpointFromStore(ep); err != nil ***REMOVED***
				logrus.Debugf("Failed to delete stale overlay endpoint (%s) from store", ep.id[0:7])
			***REMOVED***
			continue
		***REMOVED***
		n.addEndpoint(ep)

		s := n.getSubnetforIP(ep.addr)
		if s == nil ***REMOVED***
			return fmt.Errorf("could not find subnet for endpoint %s", ep.id)
		***REMOVED***

		if err := n.joinSandbox(true); err != nil ***REMOVED***
			return fmt.Errorf("restore network sandbox failed: %v", err)
		***REMOVED***

		if err := n.joinSubnetSandbox(s, true); err != nil ***REMOVED***
			return fmt.Errorf("restore subnet sandbox failed for %q: %v", s.subnetIP.String(), err)
		***REMOVED***

		Ifaces := make(map[string][]osl.IfaceOption)
		vethIfaceOption := make([]osl.IfaceOption, 1)
		vethIfaceOption = append(vethIfaceOption, n.sbox.InterfaceOptions().Master(s.brName))
		Ifaces["veth+veth"] = vethIfaceOption

		err := n.sbox.Restore(Ifaces, nil, nil, nil)
		if err != nil ***REMOVED***
			return fmt.Errorf("failed to restore overlay sandbox: %v", err)
		***REMOVED***

		n.incEndpointCount()
		d.peerAdd(ep.nid, ep.id, ep.addr.IP, ep.addr.Mask, ep.mac, net.ParseIP(d.advertiseAddress), false, false, true)
	***REMOVED***
	return nil
***REMOVED***

// Fini cleans up the driver resources
func Fini(drv driverapi.Driver) ***REMOVED***
	d := drv.(*driver)

	// Notify the peer go routine to return
	if d.peerOpCancel != nil ***REMOVED***
		d.peerOpCancel()
	***REMOVED***

	if d.exitCh != nil ***REMOVED***
		waitCh := make(chan struct***REMOVED******REMOVED***)

		d.exitCh <- waitCh

		<-waitCh
	***REMOVED***
***REMOVED***

func (d *driver) configure() error ***REMOVED***

	// Apply OS specific kernel configs if needed
	d.initOS.Do(applyOStweaks)

	if d.store == nil ***REMOVED***
		return nil
	***REMOVED***

	if d.vxlanIdm == nil ***REMOVED***
		return d.initializeVxlanIdm()
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) initializeVxlanIdm() error ***REMOVED***
	var err error

	initVxlanIdm <- true
	defer func() ***REMOVED*** <-initVxlanIdm ***REMOVED***()

	if d.vxlanIdm != nil ***REMOVED***
		return nil
	***REMOVED***

	d.vxlanIdm, err = idm.New(d.store, "vxlan-id", vxlanIDStart, vxlanIDEnd)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to initialize vxlan id manager: %v", err)
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) Type() string ***REMOVED***
	return networkType
***REMOVED***

func (d *driver) IsBuiltIn() bool ***REMOVED***
	return true
***REMOVED***

func validateSelf(node string) error ***REMOVED***
	advIP := net.ParseIP(node)
	if advIP == nil ***REMOVED***
		return fmt.Errorf("invalid self address (%s)", node)
	***REMOVED***

	addrs, err := net.InterfaceAddrs()
	if err != nil ***REMOVED***
		return fmt.Errorf("Unable to get interface addresses %v", err)
	***REMOVED***
	for _, addr := range addrs ***REMOVED***
		ip, _, err := net.ParseCIDR(addr.String())
		if err == nil && ip.Equal(advIP) ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	return fmt.Errorf("Multi-Host overlay networking requires cluster-advertise(%s) to be configured with a local ip-address that is reachable within the cluster", advIP.String())
***REMOVED***

func (d *driver) nodeJoin(advertiseAddress, bindAddress string, self bool) ***REMOVED***
	if self && !d.isSerfAlive() ***REMOVED***
		d.Lock()
		d.advertiseAddress = advertiseAddress
		d.bindAddress = bindAddress
		d.Unlock()

		// If containers are already running on this network update the
		// advertise address in the peerDB
		d.localJoinOnce.Do(func() ***REMOVED***
			d.peerDBUpdateSelf()
		***REMOVED***)

		// If there is no cluster store there is no need to start serf.
		if d.store != nil ***REMOVED***
			if err := validateSelf(advertiseAddress); err != nil ***REMOVED***
				logrus.Warn(err.Error())
			***REMOVED***
			err := d.serfInit()
			if err != nil ***REMOVED***
				logrus.Errorf("initializing serf instance failed: %v", err)
				d.Lock()
				d.advertiseAddress = ""
				d.bindAddress = ""
				d.Unlock()
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***

	d.Lock()
	if !self ***REMOVED***
		d.neighIP = advertiseAddress
	***REMOVED***
	neighIP := d.neighIP
	d.Unlock()

	if d.serfInstance != nil && neighIP != "" ***REMOVED***
		var err error
		d.joinOnce.Do(func() ***REMOVED***
			err = d.serfJoin(neighIP)
			if err == nil ***REMOVED***
				d.pushLocalDb()
			***REMOVED***
		***REMOVED***)
		if err != nil ***REMOVED***
			logrus.Errorf("joining serf neighbor %s failed: %v", advertiseAddress, err)
			d.Lock()
			d.joinOnce = sync.Once***REMOVED******REMOVED***
			d.Unlock()
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (d *driver) pushLocalEndpointEvent(action, nid, eid string) ***REMOVED***
	n := d.network(nid)
	if n == nil ***REMOVED***
		logrus.Debugf("Error pushing local endpoint event for network %s", nid)
		return
	***REMOVED***
	ep := n.endpoint(eid)
	if ep == nil ***REMOVED***
		logrus.Debugf("Error pushing local endpoint event for ep %s / %s", nid, eid)
		return
	***REMOVED***

	if !d.isSerfAlive() ***REMOVED***
		return
	***REMOVED***
	d.notifyCh <- ovNotify***REMOVED***
		action: "join",
		nw:     n,
		ep:     ep,
	***REMOVED***
***REMOVED***

// DiscoverNew is a notification for a new discovery event, such as a new node joining a cluster
func (d *driver) DiscoverNew(dType discoverapi.DiscoveryType, data interface***REMOVED******REMOVED***) error ***REMOVED***
	var err error
	switch dType ***REMOVED***
	case discoverapi.NodeDiscovery:
		nodeData, ok := data.(discoverapi.NodeDiscoveryData)
		if !ok || nodeData.Address == "" ***REMOVED***
			return fmt.Errorf("invalid discovery data")
		***REMOVED***
		d.nodeJoin(nodeData.Address, nodeData.BindAddress, nodeData.Self)
	case discoverapi.DatastoreConfig:
		if d.store != nil ***REMOVED***
			return types.ForbiddenErrorf("cannot accept datastore configuration: Overlay driver has a datastore configured already")
		***REMOVED***
		dsc, ok := data.(discoverapi.DatastoreConfigData)
		if !ok ***REMOVED***
			return types.InternalErrorf("incorrect data in datastore configuration: %v", data)
		***REMOVED***
		d.store, err = datastore.NewDataStoreFromConfig(dsc)
		if err != nil ***REMOVED***
			return types.InternalErrorf("failed to initialize data store: %v", err)
		***REMOVED***
	case discoverapi.EncryptionKeysConfig:
		encrData, ok := data.(discoverapi.DriverEncryptionConfig)
		if !ok ***REMOVED***
			return fmt.Errorf("invalid encryption key notification data")
		***REMOVED***
		keys := make([]*key, 0, len(encrData.Keys))
		for i := 0; i < len(encrData.Keys); i++ ***REMOVED***
			k := &key***REMOVED***
				value: encrData.Keys[i],
				tag:   uint32(encrData.Tags[i]),
			***REMOVED***
			keys = append(keys, k)
		***REMOVED***
		if err := d.setKeys(keys); err != nil ***REMOVED***
			logrus.Warn(err)
		***REMOVED***
	case discoverapi.EncryptionKeysUpdate:
		var newKey, delKey, priKey *key
		encrData, ok := data.(discoverapi.DriverEncryptionUpdate)
		if !ok ***REMOVED***
			return fmt.Errorf("invalid encryption key notification data")
		***REMOVED***
		if encrData.Key != nil ***REMOVED***
			newKey = &key***REMOVED***
				value: encrData.Key,
				tag:   uint32(encrData.Tag),
			***REMOVED***
		***REMOVED***
		if encrData.Primary != nil ***REMOVED***
			priKey = &key***REMOVED***
				value: encrData.Primary,
				tag:   uint32(encrData.PrimaryTag),
			***REMOVED***
		***REMOVED***
		if encrData.Prune != nil ***REMOVED***
			delKey = &key***REMOVED***
				value: encrData.Prune,
				tag:   uint32(encrData.PruneTag),
			***REMOVED***
		***REMOVED***
		if err := d.updateKeys(newKey, priKey, delKey); err != nil ***REMOVED***
			logrus.Warn(err)
		***REMOVED***
	default:
	***REMOVED***
	return nil
***REMOVED***

// DiscoverDelete is a notification for a discovery delete event, such as a node leaving a cluster
func (d *driver) DiscoverDelete(dType discoverapi.DiscoveryType, data interface***REMOVED******REMOVED***) error ***REMOVED***
	return nil
***REMOVED***
