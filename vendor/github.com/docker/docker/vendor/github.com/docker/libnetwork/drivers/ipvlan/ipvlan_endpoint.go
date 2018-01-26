package ipvlan

import (
	"fmt"

	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/netlabel"
	"github.com/docker/libnetwork/ns"
	"github.com/docker/libnetwork/osl"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

// CreateEndpoint assigns the mac, ip and endpoint id for the new container
func (d *driver) CreateEndpoint(nid, eid string, ifInfo driverapi.InterfaceInfo,
	epOptions map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	defer osl.InitOSContext()()

	if err := validateID(nid, eid); err != nil ***REMOVED***
		return err
	***REMOVED***
	n, err := d.getNetwork(nid)
	if err != nil ***REMOVED***
		return fmt.Errorf("network id %q not found", nid)
	***REMOVED***
	if ifInfo.MacAddress() != nil ***REMOVED***
		return fmt.Errorf("%s interfaces do not support custom mac address assigment", ipvlanType)
	***REMOVED***
	ep := &endpoint***REMOVED***
		id:     eid,
		nid:    nid,
		addr:   ifInfo.Address(),
		addrv6: ifInfo.AddressIPv6(),
	***REMOVED***
	if ep.addr == nil ***REMOVED***
		return fmt.Errorf("create endpoint was not passed an IP address")
	***REMOVED***
	// disallow port mapping -p
	if opt, ok := epOptions[netlabel.PortMap]; ok ***REMOVED***
		if _, ok := opt.([]types.PortBinding); ok ***REMOVED***
			if len(opt.([]types.PortBinding)) > 0 ***REMOVED***
				logrus.Warnf("%s driver does not support port mappings", ipvlanType)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// disallow port exposure --expose
	if opt, ok := epOptions[netlabel.ExposedPorts]; ok ***REMOVED***
		if _, ok := opt.([]types.TransportPort); ok ***REMOVED***
			if len(opt.([]types.TransportPort)) > 0 ***REMOVED***
				logrus.Warnf("%s driver does not support port exposures", ipvlanType)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if err := d.storeUpdate(ep); err != nil ***REMOVED***
		return fmt.Errorf("failed to save ipvlan endpoint %s to store: %v", ep.id[0:7], err)
	***REMOVED***

	n.addEndpoint(ep)

	return nil
***REMOVED***

// DeleteEndpoint remove the endpoint and associated netlink interface
func (d *driver) DeleteEndpoint(nid, eid string) error ***REMOVED***
	defer osl.InitOSContext()()
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
	if link, err := ns.NlHandle().LinkByName(ep.srcName); err == nil ***REMOVED***
		ns.NlHandle().LinkDel(link)
	***REMOVED***

	if err := d.storeDelete(ep); err != nil ***REMOVED***
		logrus.Warnf("Failed to remove ipvlan endpoint %s from store: %v", ep.id[0:7], err)
	***REMOVED***
	n.deleteEndpoint(ep.id)
	return nil
***REMOVED***
