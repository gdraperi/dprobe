package macvlan

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
	macvlanPrefix         = "macvlan"
	macvlanNetworkPrefix  = macvlanPrefix + "/network"
	macvlanEndpointPrefix = macvlanPrefix + "/endpoint"
)

// networkConfiguration for this driver's network specific configuration
type configuration struct ***REMOVED***
	ID               string
	Mtu              int
	dbIndex          uint64
	dbExists         bool
	Internal         bool
	Parent           string
	MacvlanMode      string
	CreatedSlaveLink bool
	Ipv4Subnets      []*ipv4Subnet
	Ipv6Subnets      []*ipv6Subnet
***REMOVED***

type ipv4Subnet struct ***REMOVED***
	SubnetIP string
	GwIP     string
***REMOVED***

type ipv6Subnet struct ***REMOVED***
	SubnetIP string
	GwIP     string
***REMOVED***

// initStore drivers are responsible for caching their own persistent state
func (d *driver) initStore(option map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	if data, ok := option[netlabel.LocalKVClient]; ok ***REMOVED***
		var err error
		dsc, ok := data.(discoverapi.DatastoreConfigData)
		if !ok ***REMOVED***
			return types.InternalErrorf("incorrect data in datastore configuration: %v", data)
		***REMOVED***
		d.store, err = datastore.NewDataStoreFromConfig(dsc)
		if err != nil ***REMOVED***
			return types.InternalErrorf("macvlan driver failed to initialize data store: %v", err)
		***REMOVED***

		return d.populateNetworks()
	***REMOVED***

	return nil
***REMOVED***

// populateNetworks is invoked at driver init to recreate persistently stored networks
func (d *driver) populateNetworks() error ***REMOVED***
	kvol, err := d.store.List(datastore.Key(macvlanPrefix), &configuration***REMOVED******REMOVED***)
	if err != nil && err != datastore.ErrKeyNotFound ***REMOVED***
		return fmt.Errorf("failed to get macvlan network configurations from store: %v", err)
	***REMOVED***
	// If empty it simply means no macvlan networks have been created yet
	if err == datastore.ErrKeyNotFound ***REMOVED***
		return nil
	***REMOVED***
	for _, kvo := range kvol ***REMOVED***
		config := kvo.(*configuration)
		if err = d.createNetwork(config); err != nil ***REMOVED***
			logrus.Warnf("Could not create macvlan network for id %s from persistent state", config.ID)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) populateEndpoints() error ***REMOVED***
	kvol, err := d.store.List(datastore.Key(macvlanEndpointPrefix), &endpoint***REMOVED******REMOVED***)
	if err != nil && err != datastore.ErrKeyNotFound ***REMOVED***
		return fmt.Errorf("failed to get macvlan endpoints from store: %v", err)
	***REMOVED***

	if err == datastore.ErrKeyNotFound ***REMOVED***
		return nil
	***REMOVED***

	for _, kvo := range kvol ***REMOVED***
		ep := kvo.(*endpoint)
		n, ok := d.networks[ep.nid]
		if !ok ***REMOVED***
			logrus.Debugf("Network (%s) not found for restored macvlan endpoint (%s)", ep.nid[0:7], ep.id[0:7])
			logrus.Debugf("Deleting stale macvlan endpoint (%s) from store", ep.id[0:7])
			if err := d.storeDelete(ep); err != nil ***REMOVED***
				logrus.Debugf("Failed to delete stale macvlan endpoint (%s) from store", ep.id[0:7])
			***REMOVED***
			continue
		***REMOVED***
		n.endpoints[ep.id] = ep
		logrus.Debugf("Endpoint (%s) restored to network (%s)", ep.id[0:7], ep.nid[0:7])
	***REMOVED***

	return nil
***REMOVED***

// storeUpdate used to update persistent macvlan network records as they are created
func (d *driver) storeUpdate(kvObject datastore.KVObject) error ***REMOVED***
	if d.store == nil ***REMOVED***
		logrus.Warnf("macvlan store not initialized. kv object %s is not added to the store", datastore.Key(kvObject.Key()...))
		return nil
	***REMOVED***
	if err := d.store.PutObjectAtomic(kvObject); err != nil ***REMOVED***
		return fmt.Errorf("failed to update macvlan store for object type %T: %v", kvObject, err)
	***REMOVED***

	return nil
***REMOVED***

// storeDelete used to delete macvlan records from persistent cache as they are deleted
func (d *driver) storeDelete(kvObject datastore.KVObject) error ***REMOVED***
	if d.store == nil ***REMOVED***
		logrus.Debugf("macvlan store not initialized. kv object %s is not deleted from store", datastore.Key(kvObject.Key()...))
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

func (config *configuration) MarshalJSON() ([]byte, error) ***REMOVED***
	nMap := make(map[string]interface***REMOVED******REMOVED***)
	nMap["ID"] = config.ID
	nMap["Mtu"] = config.Mtu
	nMap["Parent"] = config.Parent
	nMap["MacvlanMode"] = config.MacvlanMode
	nMap["Internal"] = config.Internal
	nMap["CreatedSubIface"] = config.CreatedSlaveLink
	if len(config.Ipv4Subnets) > 0 ***REMOVED***
		iis, err := json.Marshal(config.Ipv4Subnets)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		nMap["Ipv4Subnets"] = string(iis)
	***REMOVED***
	if len(config.Ipv6Subnets) > 0 ***REMOVED***
		iis, err := json.Marshal(config.Ipv6Subnets)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		nMap["Ipv6Subnets"] = string(iis)
	***REMOVED***

	return json.Marshal(nMap)
***REMOVED***

func (config *configuration) UnmarshalJSON(b []byte) error ***REMOVED***
	var (
		err  error
		nMap map[string]interface***REMOVED******REMOVED***
	)

	if err = json.Unmarshal(b, &nMap); err != nil ***REMOVED***
		return err
	***REMOVED***
	config.ID = nMap["ID"].(string)
	config.Mtu = int(nMap["Mtu"].(float64))
	config.Parent = nMap["Parent"].(string)
	config.MacvlanMode = nMap["MacvlanMode"].(string)
	config.Internal = nMap["Internal"].(bool)
	config.CreatedSlaveLink = nMap["CreatedSubIface"].(bool)
	if v, ok := nMap["Ipv4Subnets"]; ok ***REMOVED***
		if err := json.Unmarshal([]byte(v.(string)), &config.Ipv4Subnets); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if v, ok := nMap["Ipv6Subnets"]; ok ***REMOVED***
		if err := json.Unmarshal([]byte(v.(string)), &config.Ipv6Subnets); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (config *configuration) Key() []string ***REMOVED***
	return []string***REMOVED***macvlanNetworkPrefix, config.ID***REMOVED***
***REMOVED***

func (config *configuration) KeyPrefix() []string ***REMOVED***
	return []string***REMOVED***macvlanNetworkPrefix***REMOVED***
***REMOVED***

func (config *configuration) Value() []byte ***REMOVED***
	b, err := json.Marshal(config)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***

	return b
***REMOVED***

func (config *configuration) SetValue(value []byte) error ***REMOVED***
	return json.Unmarshal(value, config)
***REMOVED***

func (config *configuration) Index() uint64 ***REMOVED***
	return config.dbIndex
***REMOVED***

func (config *configuration) SetIndex(index uint64) ***REMOVED***
	config.dbIndex = index
	config.dbExists = true
***REMOVED***

func (config *configuration) Exists() bool ***REMOVED***
	return config.dbExists
***REMOVED***

func (config *configuration) Skip() bool ***REMOVED***
	return false
***REMOVED***

func (config *configuration) New() datastore.KVObject ***REMOVED***
	return &configuration***REMOVED******REMOVED***
***REMOVED***

func (config *configuration) CopyTo(o datastore.KVObject) error ***REMOVED***
	dstNcfg := o.(*configuration)
	*dstNcfg = *config

	return nil
***REMOVED***

func (config *configuration) DataScope() string ***REMOVED***
	return datastore.LocalScope
***REMOVED***

func (ep *endpoint) MarshalJSON() ([]byte, error) ***REMOVED***
	epMap := make(map[string]interface***REMOVED******REMOVED***)
	epMap["id"] = ep.id
	epMap["nid"] = ep.nid
	epMap["SrcName"] = ep.srcName
	if len(ep.mac) != 0 ***REMOVED***
		epMap["MacAddress"] = ep.mac.String()
	***REMOVED***
	if ep.addr != nil ***REMOVED***
		epMap["Addr"] = ep.addr.String()
	***REMOVED***
	if ep.addrv6 != nil ***REMOVED***
		epMap["Addrv6"] = ep.addrv6.String()
	***REMOVED***
	return json.Marshal(epMap)
***REMOVED***

func (ep *endpoint) UnmarshalJSON(b []byte) error ***REMOVED***
	var (
		err   error
		epMap map[string]interface***REMOVED******REMOVED***
	)

	if err = json.Unmarshal(b, &epMap); err != nil ***REMOVED***
		return fmt.Errorf("Failed to unmarshal to macvlan endpoint: %v", err)
	***REMOVED***

	if v, ok := epMap["MacAddress"]; ok ***REMOVED***
		if ep.mac, err = net.ParseMAC(v.(string)); err != nil ***REMOVED***
			return types.InternalErrorf("failed to decode macvlan endpoint MAC address (%s) after json unmarshal: %v", v.(string), err)
		***REMOVED***
	***REMOVED***
	if v, ok := epMap["Addr"]; ok ***REMOVED***
		if ep.addr, err = types.ParseCIDR(v.(string)); err != nil ***REMOVED***
			return types.InternalErrorf("failed to decode macvlan endpoint IPv4 address (%s) after json unmarshal: %v", v.(string), err)
		***REMOVED***
	***REMOVED***
	if v, ok := epMap["Addrv6"]; ok ***REMOVED***
		if ep.addrv6, err = types.ParseCIDR(v.(string)); err != nil ***REMOVED***
			return types.InternalErrorf("failed to decode macvlan endpoint IPv6 address (%s) after json unmarshal: %v", v.(string), err)
		***REMOVED***
	***REMOVED***
	ep.id = epMap["id"].(string)
	ep.nid = epMap["nid"].(string)
	ep.srcName = epMap["SrcName"].(string)

	return nil
***REMOVED***

func (ep *endpoint) Key() []string ***REMOVED***
	return []string***REMOVED***macvlanEndpointPrefix, ep.id***REMOVED***
***REMOVED***

func (ep *endpoint) KeyPrefix() []string ***REMOVED***
	return []string***REMOVED***macvlanEndpointPrefix***REMOVED***
***REMOVED***

func (ep *endpoint) Value() []byte ***REMOVED***
	b, err := json.Marshal(ep)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	return b
***REMOVED***

func (ep *endpoint) SetValue(value []byte) error ***REMOVED***
	return json.Unmarshal(value, ep)
***REMOVED***

func (ep *endpoint) Index() uint64 ***REMOVED***
	return ep.dbIndex
***REMOVED***

func (ep *endpoint) SetIndex(index uint64) ***REMOVED***
	ep.dbIndex = index
	ep.dbExists = true
***REMOVED***

func (ep *endpoint) Exists() bool ***REMOVED***
	return ep.dbExists
***REMOVED***

func (ep *endpoint) Skip() bool ***REMOVED***
	return false
***REMOVED***

func (ep *endpoint) New() datastore.KVObject ***REMOVED***
	return &endpoint***REMOVED******REMOVED***
***REMOVED***

func (ep *endpoint) CopyTo(o datastore.KVObject) error ***REMOVED***
	dstEp := o.(*endpoint)
	*dstEp = *ep
	return nil
***REMOVED***

func (ep *endpoint) DataScope() string ***REMOVED***
	return datastore.LocalScope
***REMOVED***
