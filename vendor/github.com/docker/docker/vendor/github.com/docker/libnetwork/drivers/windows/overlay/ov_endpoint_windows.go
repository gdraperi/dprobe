package overlay

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/Microsoft/hcsshim"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/drivers/windows"
	"github.com/docker/libnetwork/netlabel"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

type endpointTable map[string]*endpoint

const overlayEndpointPrefix = "overlay/endpoint"

type endpoint struct ***REMOVED***
	id             string
	nid            string
	profileID      string
	remote         bool
	mac            net.HardwareAddr
	addr           *net.IPNet
	disablegateway bool
	portMapping    []types.PortBinding // Operation port bindings
***REMOVED***

func validateID(nid, eid string) error ***REMOVED***
	if nid == "" ***REMOVED***
		return fmt.Errorf("invalid network id")
	***REMOVED***

	if eid == "" ***REMOVED***
		return fmt.Errorf("invalid endpoint id")
	***REMOVED***

	return nil
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

func (n *network) removeEndpointWithAddress(addr *net.IPNet) ***REMOVED***
	var networkEndpoint *endpoint
	n.Lock()
	for _, ep := range n.endpoints ***REMOVED***
		if ep.addr.IP.Equal(addr.IP) ***REMOVED***
			networkEndpoint = ep
			break
		***REMOVED***
	***REMOVED***

	if networkEndpoint != nil ***REMOVED***
		delete(n.endpoints, networkEndpoint.id)
	***REMOVED***
	n.Unlock()

	if networkEndpoint != nil ***REMOVED***
		logrus.Debugf("Removing stale endpoint from HNS")
		_, err := hcsshim.HNSEndpointRequest("DELETE", networkEndpoint.profileID, "")

		if err != nil ***REMOVED***
			logrus.Debugf("Failed to delete stale overlay endpoint (%s) from hns", networkEndpoint.id[0:7])
		***REMOVED***
	***REMOVED***
***REMOVED***

func (d *driver) CreateEndpoint(nid, eid string, ifInfo driverapi.InterfaceInfo,
	epOptions map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	var err error
	if err = validateID(nid, eid); err != nil ***REMOVED***
		return err
	***REMOVED***

	n := d.network(nid)
	if n == nil ***REMOVED***
		return fmt.Errorf("network id %q not found", nid)
	***REMOVED***

	ep := n.endpoint(eid)
	if ep != nil ***REMOVED***
		logrus.Debugf("Deleting stale endpoint %s", eid)
		n.deleteEndpoint(eid)

		_, err := hcsshim.HNSEndpointRequest("DELETE", ep.profileID, "")
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	ep = &endpoint***REMOVED***
		id:   eid,
		nid:  n.id,
		addr: ifInfo.Address(),
		mac:  ifInfo.MacAddress(),
	***REMOVED***

	if ep.addr == nil ***REMOVED***
		return fmt.Errorf("create endpoint was not passed interface IP address")
	***REMOVED***

	s := n.getSubnetforIP(ep.addr)
	if s == nil ***REMOVED***
		return fmt.Errorf("no matching subnet for IP %q in network %q", ep.addr, nid)
	***REMOVED***

	// Todo: Add port bindings and qos policies here

	hnsEndpoint := &hcsshim.HNSEndpoint***REMOVED***
		Name:              eid,
		VirtualNetwork:    n.hnsID,
		IPAddress:         ep.addr.IP,
		EnableInternalDNS: true,
		GatewayAddress:    s.gwIP.String(),
	***REMOVED***

	if ep.mac != nil ***REMOVED***
		hnsEndpoint.MacAddress = ep.mac.String()
	***REMOVED***

	paPolicy, err := json.Marshal(hcsshim.PaPolicy***REMOVED***
		Type: "PA",
		PA:   n.providerAddress,
	***REMOVED***)

	if err != nil ***REMOVED***
		return err
	***REMOVED***

	hnsEndpoint.Policies = append(hnsEndpoint.Policies, paPolicy)

	if system.GetOSVersion().Build > 16236 ***REMOVED***
		natPolicy, err := json.Marshal(hcsshim.PaPolicy***REMOVED***
			Type: "OutBoundNAT",
		***REMOVED***)

		if err != nil ***REMOVED***
			return err
		***REMOVED***

		hnsEndpoint.Policies = append(hnsEndpoint.Policies, natPolicy)

		epConnectivity, err := windows.ParseEndpointConnectivity(epOptions)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		pbPolicy, err := windows.ConvertPortBindings(epConnectivity.PortBindings)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		hnsEndpoint.Policies = append(hnsEndpoint.Policies, pbPolicy...)

		ep.disablegateway = true
	***REMOVED***

	configurationb, err := json.Marshal(hnsEndpoint)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	hnsresponse, err := hcsshim.HNSEndpointRequest("POST", "", string(configurationb))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	ep.profileID = hnsresponse.Id

	if ep.mac == nil ***REMOVED***
		ep.mac, err = net.ParseMAC(hnsresponse.MacAddress)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := ifInfo.SetMacAddress(ep.mac); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	ep.portMapping, err = windows.ParsePortBindingPolicies(hnsresponse.Policies)
	if err != nil ***REMOVED***
		hcsshim.HNSEndpointRequest("DELETE", hnsresponse.Id, "")
		return err
	***REMOVED***

	n.addEndpoint(ep)

	return nil
***REMOVED***

func (d *driver) DeleteEndpoint(nid, eid string) error ***REMOVED***
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

	_, err := hcsshim.HNSEndpointRequest("DELETE", ep.profileID, "")
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) EndpointOperInfo(nid, eid string) (map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	if err := validateID(nid, eid); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	n := d.network(nid)
	if n == nil ***REMOVED***
		return nil, fmt.Errorf("network id %q not found", nid)
	***REMOVED***

	ep := n.endpoint(eid)
	if ep == nil ***REMOVED***
		return nil, fmt.Errorf("endpoint id %q not found", eid)
	***REMOVED***

	data := make(map[string]interface***REMOVED******REMOVED***, 1)
	data["hnsid"] = ep.profileID
	data["AllowUnqualifiedDNSQuery"] = true

	if ep.portMapping != nil ***REMOVED***
		// Return a copy of the operational data
		pmc := make([]types.PortBinding, 0, len(ep.portMapping))
		for _, pm := range ep.portMapping ***REMOVED***
			pmc = append(pmc, pm.GetCopy())
		***REMOVED***
		data[netlabel.PortMap] = pmc
	***REMOVED***

	return data, nil
***REMOVED***
