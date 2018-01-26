// +build windows

package windows

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
	windowsPrefix         = "windows"
	windowsEndpointPrefix = "windows-endpoint"
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
			return types.InternalErrorf("windows driver failed to initialize data store: %v", err)
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
	kvol, err := d.store.List(datastore.Key(windowsPrefix), &networkConfiguration***REMOVED***Type: d.name***REMOVED***)
	if err != nil && err != datastore.ErrKeyNotFound ***REMOVED***
		return fmt.Errorf("failed to get windows network configurations from store: %v", err)
	***REMOVED***

	// It's normal for network configuration state to be empty. Just return.
	if err == datastore.ErrKeyNotFound ***REMOVED***
		return nil
	***REMOVED***

	for _, kvo := range kvol ***REMOVED***
		ncfg := kvo.(*networkConfiguration)
		if ncfg.Type != d.name ***REMOVED***
			continue
		***REMOVED***
		if err = d.createNetwork(ncfg); err != nil ***REMOVED***
			logrus.Warnf("could not create windows network for id %s hnsid %s while booting up from persistent state: %v", ncfg.ID, ncfg.HnsID, err)
		***REMOVED***
		logrus.Debugf("Network  %v (%s) restored", d.name, ncfg.ID[0:7])
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) populateEndpoints() error ***REMOVED***
	kvol, err := d.store.List(datastore.Key(windowsEndpointPrefix), &hnsEndpoint***REMOVED***Type: d.name***REMOVED***)
	if err != nil && err != datastore.ErrKeyNotFound ***REMOVED***
		return fmt.Errorf("failed to get endpoints from store: %v", err)
	***REMOVED***

	if err == datastore.ErrKeyNotFound ***REMOVED***
		return nil
	***REMOVED***

	for _, kvo := range kvol ***REMOVED***
		ep := kvo.(*hnsEndpoint)
		if ep.Type != d.name ***REMOVED***
			continue
		***REMOVED***
		n, ok := d.networks[ep.nid]
		if !ok ***REMOVED***
			logrus.Debugf("Network (%s) not found for restored endpoint (%s)", ep.nid[0:7], ep.id[0:7])
			logrus.Debugf("Deleting stale endpoint (%s) from store", ep.id[0:7])
			if err := d.storeDelete(ep); err != nil ***REMOVED***
				logrus.Debugf("Failed to delete stale endpoint (%s) from store", ep.id[0:7])
			***REMOVED***
			continue
		***REMOVED***
		n.endpoints[ep.id] = ep
		logrus.Debugf("Endpoint (%s) restored to network (%s)", ep.id[0:7], ep.nid[0:7])
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) storeUpdate(kvObject datastore.KVObject) error ***REMOVED***
	if d.store == nil ***REMOVED***
		logrus.Warnf("store not initialized. kv object %s is not added to the store", datastore.Key(kvObject.Key()...))
		return nil
	***REMOVED***

	if err := d.store.PutObjectAtomic(kvObject); err != nil ***REMOVED***
		return fmt.Errorf("failed to update store for object type %T: %v", kvObject, err)
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) storeDelete(kvObject datastore.KVObject) error ***REMOVED***
	if d.store == nil ***REMOVED***
		logrus.Debugf("store not initialized. kv object %s is not deleted from store", datastore.Key(kvObject.Key()...))
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
	nMap["Type"] = ncfg.Type
	nMap["Name"] = ncfg.Name
	nMap["HnsID"] = ncfg.HnsID
	nMap["VLAN"] = ncfg.VLAN
	nMap["VSID"] = ncfg.VSID
	nMap["DNSServers"] = ncfg.DNSServers
	nMap["DNSSuffix"] = ncfg.DNSSuffix
	nMap["SourceMac"] = ncfg.SourceMac
	nMap["NetworkAdapterName"] = ncfg.NetworkAdapterName

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

	ncfg.ID = nMap["ID"].(string)
	ncfg.Type = nMap["Type"].(string)
	ncfg.Name = nMap["Name"].(string)
	ncfg.HnsID = nMap["HnsID"].(string)
	ncfg.VLAN = uint(nMap["VLAN"].(float64))
	ncfg.VSID = uint(nMap["VSID"].(float64))
	ncfg.DNSServers = nMap["DNSServers"].(string)
	ncfg.DNSSuffix = nMap["DNSSuffix"].(string)
	ncfg.SourceMac = nMap["SourceMac"].(string)
	ncfg.NetworkAdapterName = nMap["NetworkAdapterName"].(string)
	return nil
***REMOVED***

func (ncfg *networkConfiguration) Key() []string ***REMOVED***
	return []string***REMOVED***windowsPrefix + ncfg.Type, ncfg.ID***REMOVED***
***REMOVED***

func (ncfg *networkConfiguration) KeyPrefix() []string ***REMOVED***
	return []string***REMOVED***windowsPrefix + ncfg.Type***REMOVED***
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
	return &networkConfiguration***REMOVED***Type: ncfg.Type***REMOVED***
***REMOVED***

func (ncfg *networkConfiguration) CopyTo(o datastore.KVObject) error ***REMOVED***
	dstNcfg := o.(*networkConfiguration)
	*dstNcfg = *ncfg
	return nil
***REMOVED***

func (ncfg *networkConfiguration) DataScope() string ***REMOVED***
	return datastore.LocalScope
***REMOVED***

func (ep *hnsEndpoint) MarshalJSON() ([]byte, error) ***REMOVED***
	epMap := make(map[string]interface***REMOVED******REMOVED***)
	epMap["id"] = ep.id
	epMap["nid"] = ep.nid
	epMap["Type"] = ep.Type
	epMap["profileID"] = ep.profileID
	epMap["MacAddress"] = ep.macAddress.String()
	if ep.addr.IP != nil ***REMOVED***
		epMap["Addr"] = ep.addr.String()
	***REMOVED***
	if ep.gateway != nil ***REMOVED***
		epMap["gateway"] = ep.gateway.String()
	***REMOVED***
	epMap["epOption"] = ep.epOption
	epMap["epConnectivity"] = ep.epConnectivity
	epMap["PortMapping"] = ep.portMapping

	return json.Marshal(epMap)
***REMOVED***

func (ep *hnsEndpoint) UnmarshalJSON(b []byte) error ***REMOVED***
	var (
		err   error
		epMap map[string]interface***REMOVED******REMOVED***
	)

	if err = json.Unmarshal(b, &epMap); err != nil ***REMOVED***
		return fmt.Errorf("Failed to unmarshal to endpoint: %v", err)
	***REMOVED***
	if v, ok := epMap["MacAddress"]; ok ***REMOVED***
		if ep.macAddress, err = net.ParseMAC(v.(string)); err != nil ***REMOVED***
			return types.InternalErrorf("failed to decode endpoint MAC address (%s) after json unmarshal: %v", v.(string), err)
		***REMOVED***
	***REMOVED***
	if v, ok := epMap["Addr"]; ok ***REMOVED***
		if ep.addr, err = types.ParseCIDR(v.(string)); err != nil ***REMOVED***
			return types.InternalErrorf("failed to decode endpoint IPv4 address (%s) after json unmarshal: %v", v.(string), err)
		***REMOVED***
	***REMOVED***
	if v, ok := epMap["gateway"]; ok ***REMOVED***
		ep.gateway = net.ParseIP(v.(string))
	***REMOVED***
	ep.id = epMap["id"].(string)
	ep.Type = epMap["Type"].(string)
	ep.nid = epMap["nid"].(string)
	ep.profileID = epMap["profileID"].(string)
	d, _ := json.Marshal(epMap["epOption"])
	if err := json.Unmarshal(d, &ep.epOption); err != nil ***REMOVED***
		logrus.Warnf("Failed to decode endpoint container config %v", err)
	***REMOVED***
	d, _ = json.Marshal(epMap["epConnectivity"])
	if err := json.Unmarshal(d, &ep.epConnectivity); err != nil ***REMOVED***
		logrus.Warnf("Failed to decode endpoint external connectivity configuration %v", err)
	***REMOVED***
	d, _ = json.Marshal(epMap["PortMapping"])
	if err := json.Unmarshal(d, &ep.portMapping); err != nil ***REMOVED***
		logrus.Warnf("Failed to decode endpoint port mapping %v", err)
	***REMOVED***

	return nil
***REMOVED***

func (ep *hnsEndpoint) Key() []string ***REMOVED***
	return []string***REMOVED***windowsEndpointPrefix + ep.Type, ep.id***REMOVED***
***REMOVED***

func (ep *hnsEndpoint) KeyPrefix() []string ***REMOVED***
	return []string***REMOVED***windowsEndpointPrefix + ep.Type***REMOVED***
***REMOVED***

func (ep *hnsEndpoint) Value() []byte ***REMOVED***
	b, err := json.Marshal(ep)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	return b
***REMOVED***

func (ep *hnsEndpoint) SetValue(value []byte) error ***REMOVED***
	return json.Unmarshal(value, ep)
***REMOVED***

func (ep *hnsEndpoint) Index() uint64 ***REMOVED***
	return ep.dbIndex
***REMOVED***

func (ep *hnsEndpoint) SetIndex(index uint64) ***REMOVED***
	ep.dbIndex = index
	ep.dbExists = true
***REMOVED***

func (ep *hnsEndpoint) Exists() bool ***REMOVED***
	return ep.dbExists
***REMOVED***

func (ep *hnsEndpoint) Skip() bool ***REMOVED***
	return false
***REMOVED***

func (ep *hnsEndpoint) New() datastore.KVObject ***REMOVED***
	return &hnsEndpoint***REMOVED***Type: ep.Type***REMOVED***
***REMOVED***

func (ep *hnsEndpoint) CopyTo(o datastore.KVObject) error ***REMOVED***
	dstEp := o.(*hnsEndpoint)
	*dstEp = *ep
	return nil
***REMOVED***

func (ep *hnsEndpoint) DataScope() string ***REMOVED***
	return datastore.LocalScope
***REMOVED***
