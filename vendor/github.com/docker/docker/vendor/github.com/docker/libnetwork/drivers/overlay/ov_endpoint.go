package overlay

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/netutils"
	"github.com/docker/libnetwork/ns"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

type endpointTable map[string]*endpoint

const overlayEndpointPrefix = "overlay/endpoint"

type endpoint struct ***REMOVED***
	id       string
	nid      string
	ifName   string
	mac      net.HardwareAddr
	addr     *net.IPNet
	dbExists bool
	dbIndex  uint64
***REMOVED***

func (n *network) endpoint(eid string) *endpoint ***REMOVED***
	n.Lock()
	defer n.Unlock()

	return n.endpoints[eid]
***REMOVED***

func (n *network) addEndpoint(ep *endpoint) ***REMOVED***
	n.Lock()
	n.endpoints[ep.id] = ep
	n.Unlock()
***REMOVED***

func (n *network) deleteEndpoint(eid string) ***REMOVED***
	n.Lock()
	delete(n.endpoints, eid)
	n.Unlock()
***REMOVED***

func (d *driver) CreateEndpoint(nid, eid string, ifInfo driverapi.InterfaceInfo,
	epOptions map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	var err error

	if err = validateID(nid, eid); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Since we perform lazy configuration make sure we try
	// configuring the driver when we enter CreateEndpoint since
	// CreateNetwork may not be called in every node.
	if err := d.configure(); err != nil ***REMOVED***
		return err
	***REMOVED***

	n := d.network(nid)
	if n == nil ***REMOVED***
		return fmt.Errorf("network id %q not found", nid)
	***REMOVED***

	ep := &endpoint***REMOVED***
		id:   eid,
		nid:  n.id,
		addr: ifInfo.Address(),
		mac:  ifInfo.MacAddress(),
	***REMOVED***
	if ep.addr == nil ***REMOVED***
		return fmt.Errorf("create endpoint was not passed interface IP address")
	***REMOVED***

	if s := n.getSubnetforIP(ep.addr); s == nil ***REMOVED***
		return fmt.Errorf("no matching subnet for IP %q in network %q", ep.addr, nid)
	***REMOVED***

	if ep.mac == nil ***REMOVED***
		ep.mac = netutils.GenerateMACFromIP(ep.addr.IP)
		if err := ifInfo.SetMacAddress(ep.mac); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	n.addEndpoint(ep)

	if err := d.writeEndpointToStore(ep); err != nil ***REMOVED***
		return fmt.Errorf("failed to update overlay endpoint %s to local store: %v", ep.id[0:7], err)
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) DeleteEndpoint(nid, eid string) error ***REMOVED***
	nlh := ns.NlHandle()

	if err := validateID(nid, eid); err != nil ***REMOVED***
		return err
	***REMOVED***

	n := d.network(nid)
	if n == nil ***REMOVED***
		return fmt.Errorf("network id %q not found", nid)
	***REMOVED***

	ep := n.endpoint(eid)
	if ep == nil ***REMOVED***
		return fmt.Errorf("endpoint id %q not found", eid)
	***REMOVED***

	n.deleteEndpoint(eid)

	if err := d.deleteEndpointFromStore(ep); err != nil ***REMOVED***
		logrus.Warnf("Failed to delete overlay endpoint %s from local store: %v", ep.id[0:7], err)
	***REMOVED***

	if ep.ifName == "" ***REMOVED***
		return nil
	***REMOVED***

	link, err := nlh.LinkByName(ep.ifName)
	if err != nil ***REMOVED***
		logrus.Debugf("Failed to retrieve interface (%s)'s link on endpoint (%s) delete: %v", ep.ifName, ep.id, err)
		return nil
	***REMOVED***
	if err := nlh.LinkDel(link); err != nil ***REMOVED***
		logrus.Debugf("Failed to delete interface (%s)'s link on endpoint (%s) delete: %v", ep.ifName, ep.id, err)
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) EndpointOperInfo(nid, eid string) (map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	return make(map[string]interface***REMOVED******REMOVED***, 0), nil
***REMOVED***

func (d *driver) deleteEndpointFromStore(e *endpoint) error ***REMOVED***
	if d.localStore == nil ***REMOVED***
		return fmt.Errorf("overlay local store not initialized, ep not deleted")
	***REMOVED***

	return d.localStore.DeleteObjectAtomic(e)
***REMOVED***

func (d *driver) writeEndpointToStore(e *endpoint) error ***REMOVED***
	if d.localStore == nil ***REMOVED***
		return fmt.Errorf("overlay local store not initialized, ep not added")
	***REMOVED***

	return d.localStore.PutObjectAtomic(e)
***REMOVED***

func (ep *endpoint) DataScope() string ***REMOVED***
	return datastore.LocalScope
***REMOVED***

func (ep *endpoint) New() datastore.KVObject ***REMOVED***
	return &endpoint***REMOVED******REMOVED***
***REMOVED***

func (ep *endpoint) CopyTo(o datastore.KVObject) error ***REMOVED***
	dstep := o.(*endpoint)
	*dstep = *ep
	return nil
***REMOVED***

func (ep *endpoint) Key() []string ***REMOVED***
	return []string***REMOVED***overlayEndpointPrefix, ep.id***REMOVED***
***REMOVED***

func (ep *endpoint) KeyPrefix() []string ***REMOVED***
	return []string***REMOVED***overlayEndpointPrefix***REMOVED***
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

func (ep *endpoint) MarshalJSON() ([]byte, error) ***REMOVED***
	epMap := make(map[string]interface***REMOVED******REMOVED***)

	epMap["id"] = ep.id
	epMap["nid"] = ep.nid
	if ep.ifName != "" ***REMOVED***
		epMap["ifName"] = ep.ifName
	***REMOVED***
	if ep.addr != nil ***REMOVED***
		epMap["addr"] = ep.addr.String()
	***REMOVED***
	if len(ep.mac) != 0 ***REMOVED***
		epMap["mac"] = ep.mac.String()
	***REMOVED***

	return json.Marshal(epMap)
***REMOVED***

func (ep *endpoint) UnmarshalJSON(value []byte) error ***REMOVED***
	var (
		err   error
		epMap map[string]interface***REMOVED******REMOVED***
	)

	json.Unmarshal(value, &epMap)

	ep.id = epMap["id"].(string)
	ep.nid = epMap["nid"].(string)
	if v, ok := epMap["mac"]; ok ***REMOVED***
		if ep.mac, err = net.ParseMAC(v.(string)); err != nil ***REMOVED***
			return types.InternalErrorf("failed to decode endpoint interface mac address after json unmarshal: %s", v.(string))
		***REMOVED***
	***REMOVED***
	if v, ok := epMap["addr"]; ok ***REMOVED***
		if ep.addr, err = types.ParseCIDR(v.(string)); err != nil ***REMOVED***
			return types.InternalErrorf("failed to decode endpoint interface ipv4 address after json unmarshal: %v", err)
		***REMOVED***
	***REMOVED***
	if v, ok := epMap["ifName"]; ok ***REMOVED***
		ep.ifName = v.(string)
	***REMOVED***

	return nil
***REMOVED***
