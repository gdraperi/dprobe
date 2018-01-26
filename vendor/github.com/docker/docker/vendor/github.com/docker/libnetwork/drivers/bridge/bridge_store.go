package bridge

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/discoverapi"
	"github.com/docker/libnetwork/netlabel"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

const (
	// network config prefix was not specific enough.
	// To be backward compatible, need custom endpoint
	// prefix with different root
	bridgePrefix         = "bridge"
	bridgeEndpointPrefix = "bridge-endpoint"
)

func (d *driver) initStore(option map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	if data, ok := option[netlabel.LocalKVClient]; ok ***REMOVED***
		var err error
		dsc, ok := data.(discoverapi.DatastoreConfigData)
		if !ok ***REMOVED***
			return types.InternalErrorf("incorrect data in datastore configuration: %v", data)
		***REMOVED***
		d.store, err = datastore.NewDataStoreFromConfig(dsc)
		if err != nil ***REMOVED***
			return types.InternalErrorf("bridge driver failed to initialize data store: %v", err)
		***REMOVED***

		err = d.populateNetworks()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		err = d.populateEndpoints()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) populateNetworks() error ***REMOVED***
	kvol, err := d.store.List(datastore.Key(bridgePrefix), &networkConfiguration***REMOVED******REMOVED***)
	if err != nil && err != datastore.ErrKeyNotFound ***REMOVED***
		return fmt.Errorf("failed to get bridge network configurations from store: %v", err)
	***REMOVED***

	// It's normal for network configuration state to be empty. Just return.
	if err == datastore.ErrKeyNotFound ***REMOVED***
		return nil
	***REMOVED***

	for _, kvo := range kvol ***REMOVED***
		ncfg := kvo.(*networkConfiguration)
		if err = d.createNetwork(ncfg); err != nil ***REMOVED***
			logrus.Warnf("could not create bridge network for id %s bridge name %s while booting up from persistent state: %v", ncfg.ID, ncfg.BridgeName, err)
		***REMOVED***
		logrus.Debugf("Network (%s) restored", ncfg.ID[0:7])
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) populateEndpoints() error ***REMOVED***
	kvol, err := d.store.List(datastore.Key(bridgeEndpointPrefix), &bridgeEndpoint***REMOVED******REMOVED***)
	if err != nil && err != datastore.ErrKeyNotFound ***REMOVED***
		return fmt.Errorf("failed to get bridge endpoints from store: %v", err)
	***REMOVED***

	if err == datastore.ErrKeyNotFound ***REMOVED***
		return nil
	***REMOVED***

	for _, kvo := range kvol ***REMOVED***
		ep := kvo.(*bridgeEndpoint)
		n, ok := d.networks[ep.nid]
		if !ok ***REMOVED***
			logrus.Debugf("Network (%s) not found for restored bridge endpoint (%s)", ep.nid[0:7], ep.id[0:7])
			logrus.Debugf("Deleting stale bridge endpoint (%s) from store", ep.id[0:7])
			if err := d.storeDelete(ep); err != nil ***REMOVED***
				logrus.Debugf("Failed to delete stale bridge endpoint (%s) from store", ep.id[0:7])
			***REMOVED***
			continue
		***REMOVED***
		n.endpoints[ep.id] = ep
		n.restorePortAllocations(ep)
		logrus.Debugf("Endpoint (%s) restored to network (%s)", ep.id[0:7], ep.nid[0:7])
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) storeUpdate(kvObject datastore.KVObject) error ***REMOVED***
	if d.store == nil ***REMOVED***
		logrus.Warnf("bridge store not initialized. kv object %s is not added to the store", datastore.Key(kvObject.Key()...))
		return nil
	***REMOVED***

	if err := d.store.PutObjectAtomic(kvObject); err != nil ***REMOVED***
		return fmt.Errorf("failed to update bridge store for object type %T: %v", kvObject, err)
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) storeDelete(kvObject datastore.KVObject) error ***REMOVED***
	if d.store == nil ***REMOVED***
		logrus.Debugf("bridge store not initialized. kv object %s is not deleted from store", datastore.Key(kvObject.Key()...))
		return nil
	***REMOVED***

retry:
	if err := d.store.DeleteObjectAtomic(kvObject); err != nil ***REMOVED***
		if err == datastore.ErrKeyModified ***REMOVED***
			if err := d.store.GetObject(datastore.Key(kvObject.Key()...), kvObject); err != nil ***REMOVED***
				return fmt.Errorf("could not update the kvobject to latest when trying to delete: %v", err)
			***REMOVED***
			goto retry
		***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (ncfg *networkConfiguration) MarshalJSON() ([]byte, error) ***REMOVED***
	nMap := make(map[string]interface***REMOVED******REMOVED***)
	nMap["ID"] = ncfg.ID
	nMap["BridgeName"] = ncfg.BridgeName
	nMap["EnableIPv6"] = ncfg.EnableIPv6
	nMap["EnableIPMasquerade"] = ncfg.EnableIPMasquerade
	nMap["EnableICC"] = ncfg.EnableICC
	nMap["Mtu"] = ncfg.Mtu
	nMap["Internal"] = ncfg.Internal
	nMap["DefaultBridge"] = ncfg.DefaultBridge
	nMap["DefaultBindingIP"] = ncfg.DefaultBindingIP.String()
	nMap["DefaultGatewayIPv4"] = ncfg.DefaultGatewayIPv4.String()
	nMap["DefaultGatewayIPv6"] = ncfg.DefaultGatewayIPv6.String()
	nMap["ContainerIfacePrefix"] = ncfg.ContainerIfacePrefix
	nMap["BridgeIfaceCreator"] = ncfg.BridgeIfaceCreator

	if ncfg.AddressIPv4 != nil ***REMOVED***
		nMap["AddressIPv4"] = ncfg.AddressIPv4.String()
	***REMOVED***

	if ncfg.AddressIPv6 != nil ***REMOVED***
		nMap["AddressIPv6"] = ncfg.AddressIPv6.String()
	***REMOVED***

	return json.Marshal(nMap)
***REMOVED***

func (ncfg *networkConfiguration) UnmarshalJSON(b []byte) error ***REMOVED***
	var (
		err  error
		nMap map[string]interface***REMOVED******REMOVED***
	)

	if err = json.Unmarshal(b, &nMap); err != nil ***REMOVED***
		return err
	***REMOVED***

	if v, ok := nMap["AddressIPv4"]; ok ***REMOVED***
		if ncfg.AddressIPv4, err = types.ParseCIDR(v.(string)); err != nil ***REMOVED***
			return types.InternalErrorf("failed to decode bridge network address IPv4 after json unmarshal: %s", v.(string))
		***REMOVED***
	***REMOVED***

	if v, ok := nMap["AddressIPv6"]; ok ***REMOVED***
		if ncfg.AddressIPv6, err = types.ParseCIDR(v.(string)); err != nil ***REMOVED***
			return types.InternalErrorf("failed to decode bridge network address IPv6 after json unmarshal: %s", v.(string))
		***REMOVED***
	***REMOVED***

	if v, ok := nMap["ContainerIfacePrefix"]; ok ***REMOVED***
		ncfg.ContainerIfacePrefix = v.(string)
	***REMOVED***

	ncfg.DefaultBridge = nMap["DefaultBridge"].(bool)
	ncfg.DefaultBindingIP = net.ParseIP(nMap["DefaultBindingIP"].(string))
	ncfg.DefaultGatewayIPv4 = net.ParseIP(nMap["DefaultGatewayIPv4"].(string))
	ncfg.DefaultGatewayIPv6 = net.ParseIP(nMap["DefaultGatewayIPv6"].(string))
	ncfg.ID = nMap["ID"].(string)
	ncfg.BridgeName = nMap["BridgeName"].(string)
	ncfg.EnableIPv6 = nMap["EnableIPv6"].(bool)
	ncfg.EnableIPMasquerade = nMap["EnableIPMasquerade"].(bool)
	ncfg.EnableICC = nMap["EnableICC"].(bool)
	ncfg.Mtu = int(nMap["Mtu"].(float64))
	if v, ok := nMap["Internal"]; ok ***REMOVED***
		ncfg.Internal = v.(bool)
	***REMOVED***

	if v, ok := nMap["BridgeIfaceCreator"]; ok ***REMOVED***
		ncfg.BridgeIfaceCreator = ifaceCreator(v.(float64))
	***REMOVED***

	return nil
***REMOVED***

func (ncfg *networkConfiguration) Key() []string ***REMOVED***
	return []string***REMOVED***bridgePrefix, ncfg.ID***REMOVED***
***REMOVED***

func (ncfg *networkConfiguration) KeyPrefix() []string ***REMOVED***
	return []string***REMOVED***bridgePrefix***REMOVED***
***REMOVED***

func (ncfg *networkConfiguration) Value() []byte ***REMOVED***
	b, err := json.Marshal(ncfg)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	return b
***REMOVED***

func (ncfg *networkConfiguration) SetValue(value []byte) error ***REMOVED***
	return json.Unmarshal(value, ncfg)
***REMOVED***

func (ncfg *networkConfiguration) Index() uint64 ***REMOVED***
	return ncfg.dbIndex
***REMOVED***

func (ncfg *networkConfiguration) SetIndex(index uint64) ***REMOVED***
	ncfg.dbIndex = index
	ncfg.dbExists = true
***REMOVED***

func (ncfg *networkConfiguration) Exists() bool ***REMOVED***
	return ncfg.dbExists
***REMOVED***

func (ncfg *networkConfiguration) Skip() bool ***REMOVED***
	return false
***REMOVED***

func (ncfg *networkConfiguration) New() datastore.KVObject ***REMOVED***
	return &networkConfiguration***REMOVED******REMOVED***
***REMOVED***

func (ncfg *networkConfiguration) CopyTo(o datastore.KVObject) error ***REMOVED***
	dstNcfg := o.(*networkConfiguration)
	*dstNcfg = *ncfg
	return nil
***REMOVED***

func (ncfg *networkConfiguration) DataScope() string ***REMOVED***
	return datastore.LocalScope
***REMOVED***

func (ep *bridgeEndpoint) MarshalJSON() ([]byte, error) ***REMOVED***
	epMap := make(map[string]interface***REMOVED******REMOVED***)
	epMap["id"] = ep.id
	epMap["nid"] = ep.nid
	epMap["SrcName"] = ep.srcName
	epMap["MacAddress"] = ep.macAddress.String()
	epMap["Addr"] = ep.addr.String()
	if ep.addrv6 != nil ***REMOVED***
		epMap["Addrv6"] = ep.addrv6.String()
	***REMOVED***
	epMap["Config"] = ep.config
	epMap["ContainerConfig"] = ep.containerConfig
	epMap["ExternalConnConfig"] = ep.extConnConfig
	epMap["PortMapping"] = ep.portMapping

	return json.Marshal(epMap)
***REMOVED***

func (ep *bridgeEndpoint) UnmarshalJSON(b []byte) error ***REMOVED***
	var (
		err   error
		epMap map[string]interface***REMOVED******REMOVED***
	)

	if err = json.Unmarshal(b, &epMap); err != nil ***REMOVED***
		return fmt.Errorf("Failed to unmarshal to bridge endpoint: %v", err)
	***REMOVED***

	if v, ok := epMap["MacAddress"]; ok ***REMOVED***
		if ep.macAddress, err = net.ParseMAC(v.(string)); err != nil ***REMOVED***
			return types.InternalErrorf("failed to decode bridge endpoint MAC address (%s) after json unmarshal: %v", v.(string), err)
		***REMOVED***
	***REMOVED***
	if v, ok := epMap["Addr"]; ok ***REMOVED***
		if ep.addr, err = types.ParseCIDR(v.(string)); err != nil ***REMOVED***
			return types.InternalErrorf("failed to decode bridge endpoint IPv4 address (%s) after json unmarshal: %v", v.(string), err)
		***REMOVED***
	***REMOVED***
	if v, ok := epMap["Addrv6"]; ok ***REMOVED***
		if ep.addrv6, err = types.ParseCIDR(v.(string)); err != nil ***REMOVED***
			return types.InternalErrorf("failed to decode bridge endpoint IPv6 address (%s) after json unmarshal: %v", v.(string), err)
		***REMOVED***
	***REMOVED***
	ep.id = epMap["id"].(string)
	ep.nid = epMap["nid"].(string)
	ep.srcName = epMap["SrcName"].(string)
	d, _ := json.Marshal(epMap["Config"])
	if err := json.Unmarshal(d, &ep.config); err != nil ***REMOVED***
		logrus.Warnf("Failed to decode endpoint config %v", err)
	***REMOVED***
	d, _ = json.Marshal(epMap["ContainerConfig"])
	if err := json.Unmarshal(d, &ep.containerConfig); err != nil ***REMOVED***
		logrus.Warnf("Failed to decode endpoint container config %v", err)
	***REMOVED***
	d, _ = json.Marshal(epMap["ExternalConnConfig"])
	if err := json.Unmarshal(d, &ep.extConnConfig); err != nil ***REMOVED***
		logrus.Warnf("Failed to decode endpoint external connectivity configuration %v", err)
	***REMOVED***
	d, _ = json.Marshal(epMap["PortMapping"])
	if err := json.Unmarshal(d, &ep.portMapping); err != nil ***REMOVED***
		logrus.Warnf("Failed to decode endpoint port mapping %v", err)
	***REMOVED***

	return nil
***REMOVED***

func (ep *bridgeEndpoint) Key() []string ***REMOVED***
	return []string***REMOVED***bridgeEndpointPrefix, ep.id***REMOVED***
***REMOVED***

func (ep *bridgeEndpoint) KeyPrefix() []string ***REMOVED***
	return []string***REMOVED***bridgeEndpointPrefix***REMOVED***
***REMOVED***

func (ep *bridgeEndpoint) Value() []byte ***REMOVED***
	b, err := json.Marshal(ep)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	return b
***REMOVED***

func (ep *bridgeEndpoint) SetValue(value []byte) error ***REMOVED***
	return json.Unmarshal(value, ep)
***REMOVED***

func (ep *bridgeEndpoint) Index() uint64 ***REMOVED***
	return ep.dbIndex
***REMOVED***

func (ep *bridgeEndpoint) SetIndex(index uint64) ***REMOVED***
	ep.dbIndex = index
	ep.dbExists = true
***REMOVED***

func (ep *bridgeEndpoint) Exists() bool ***REMOVED***
	return ep.dbExists
***REMOVED***

func (ep *bridgeEndpoint) Skip() bool ***REMOVED***
	return false
***REMOVED***

func (ep *bridgeEndpoint) New() datastore.KVObject ***REMOVED***
	return &bridgeEndpoint***REMOVED******REMOVED***
***REMOVED***

func (ep *bridgeEndpoint) CopyTo(o datastore.KVObject) error ***REMOVED***
	dstEp := o.(*bridgeEndpoint)
	*dstEp = *ep
	return nil
***REMOVED***

func (ep *bridgeEndpoint) DataScope() string ***REMOVED***
	return datastore.LocalScope
***REMOVED***

func (n *bridgeNetwork) restorePortAllocations(ep *bridgeEndpoint) ***REMOVED***
	if ep.extConnConfig == nil ||
		ep.extConnConfig.ExposedPorts == nil ||
		ep.extConnConfig.PortBindings == nil ***REMOVED***
		return
	***REMOVED***
	tmp := ep.extConnConfig.PortBindings
	ep.extConnConfig.PortBindings = ep.portMapping
	_, err := n.allocatePorts(ep, n.config.DefaultBindingIP, n.driver.config.EnableUserlandProxy)
	if err != nil ***REMOVED***
		logrus.Warnf("Failed to reserve existing port mapping for endpoint %s:%v", ep.id[0:7], err)
	***REMOVED***
	ep.extConnConfig.PortBindings = tmp
***REMOVED***
