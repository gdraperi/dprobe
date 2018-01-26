package overlay

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/Microsoft/hcsshim"
	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/netlabel"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

var (
	hostMode  bool
	networkMu sync.Mutex
)

type networkTable map[string]*network

type subnet struct ***REMOVED***
	vni      uint32
	subnetIP *net.IPNet
	gwIP     *net.IP
***REMOVED***

type subnetJSON struct ***REMOVED***
	SubnetIP string
	GwIP     string
	Vni      uint32
***REMOVED***

type network struct ***REMOVED***
	id              string
	name            string
	hnsID           string
	providerAddress string
	interfaceName   string
	endpoints       endpointTable
	driver          *driver
	initEpoch       int
	initErr         error
	subnets         []*subnet
	secure          bool
	sync.Mutex
***REMOVED***

func (d *driver) NetworkAllocate(id string, option map[string]string, ipV4Data, ipV6Data []driverapi.IPAMData) (map[string]string, error) ***REMOVED***
	return nil, types.NotImplementedErrorf("not implemented")
***REMOVED***

func (d *driver) NetworkFree(id string) error ***REMOVED***
	return types.NotImplementedErrorf("not implemented")
***REMOVED***

func (d *driver) CreateNetwork(id string, option map[string]interface***REMOVED******REMOVED***, nInfo driverapi.NetworkInfo, ipV4Data, ipV6Data []driverapi.IPAMData) error ***REMOVED***
	var (
		networkName   string
		interfaceName string
		staleNetworks []string
	)

	if id == "" ***REMOVED***
		return fmt.Errorf("invalid network id")
	***REMOVED***

	if nInfo == nil ***REMOVED***
		return fmt.Errorf("invalid network info structure")
	***REMOVED***

	if len(ipV4Data) == 0 || ipV4Data[0].Pool.String() == "0.0.0.0/0" ***REMOVED***
		return types.BadRequestErrorf("ipv4 pool is empty")
	***REMOVED***

	staleNetworks = make([]string, 0)
	vnis := make([]uint32, 0, len(ipV4Data))

	existingNetwork := d.network(id)
	if existingNetwork != nil ***REMOVED***
		logrus.Debugf("Network preexists. Deleting %s", id)
		err := d.DeleteNetwork(id)
		if err != nil ***REMOVED***
			logrus.Errorf("Error deleting stale network %s", err.Error())
		***REMOVED***
	***REMOVED***

	n := &network***REMOVED***
		id:        id,
		driver:    d,
		endpoints: endpointTable***REMOVED******REMOVED***,
		subnets:   []*subnet***REMOVED******REMOVED***,
	***REMOVED***

	genData, ok := option[netlabel.GenericData].(map[string]string)

	if !ok ***REMOVED***
		return fmt.Errorf("Unknown generic data option")
	***REMOVED***

	for label, value := range genData ***REMOVED***
		switch label ***REMOVED***
		case "com.docker.network.windowsshim.networkname":
			networkName = value
		case "com.docker.network.windowsshim.interface":
			interfaceName = value
		case "com.docker.network.windowsshim.hnsid":
			n.hnsID = value
		case netlabel.OverlayVxlanIDList:
			vniStrings := strings.Split(value, ",")
			for _, vniStr := range vniStrings ***REMOVED***
				vni, err := strconv.Atoi(vniStr)
				if err != nil ***REMOVED***
					return fmt.Errorf("invalid vxlan id value %q passed", vniStr)
				***REMOVED***

				vnis = append(vnis, uint32(vni))
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// If we are getting vnis from libnetwork, either we get for
	// all subnets or none.
	if len(vnis) < len(ipV4Data) ***REMOVED***
		return fmt.Errorf("insufficient vnis(%d) passed to overlay. Windows driver requires VNIs to be prepopulated", len(vnis))
	***REMOVED***

	for i, ipd := range ipV4Data ***REMOVED***
		s := &subnet***REMOVED***
			subnetIP: ipd.Pool,
			gwIP:     &ipd.Gateway.IP,
		***REMOVED***

		if len(vnis) != 0 ***REMOVED***
			s.vni = vnis[i]
		***REMOVED***

		d.Lock()
		for _, network := range d.networks ***REMOVED***
			found := false
			for _, sub := range network.subnets ***REMOVED***
				if sub.vni == s.vni ***REMOVED***
					staleNetworks = append(staleNetworks, network.id)
					found = true
					break
				***REMOVED***
			***REMOVED***
			if found ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		d.Unlock()

		n.subnets = append(n.subnets, s)
	***REMOVED***

	for _, staleNetwork := range staleNetworks ***REMOVED***
		d.DeleteNetwork(staleNetwork)
	***REMOVED***

	n.name = networkName
	if n.name == "" ***REMOVED***
		n.name = id
	***REMOVED***

	n.interfaceName = interfaceName

	if nInfo != nil ***REMOVED***
		if err := nInfo.TableEventRegister(ovPeerTable, driverapi.EndpointObject); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	d.addNetwork(n)

	err := d.createHnsNetwork(n)

	if err != nil ***REMOVED***
		d.deleteNetwork(id)
	***REMOVED*** else ***REMOVED***
		genData["com.docker.network.windowsshim.hnsid"] = n.hnsID
	***REMOVED***

	return err
***REMOVED***

func (d *driver) DeleteNetwork(nid string) error ***REMOVED***
	if nid == "" ***REMOVED***
		return fmt.Errorf("invalid network id")
	***REMOVED***

	n := d.network(nid)
	if n == nil ***REMOVED***
		return types.ForbiddenErrorf("could not find network with id %s", nid)
	***REMOVED***

	_, err := hcsshim.HNSNetworkRequest("DELETE", n.hnsID, "")
	if err != nil ***REMOVED***
		return types.ForbiddenErrorf(err.Error())
	***REMOVED***

	d.deleteNetwork(nid)

	return nil
***REMOVED***

func (d *driver) ProgramExternalConnectivity(nid, eid string, options map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	return nil
***REMOVED***

func (d *driver) RevokeExternalConnectivity(nid, eid string) error ***REMOVED***
	return nil
***REMOVED***

func (d *driver) addNetwork(n *network) ***REMOVED***
	d.Lock()
	d.networks[n.id] = n
	d.Unlock()
***REMOVED***

func (d *driver) deleteNetwork(nid string) ***REMOVED***
	d.Lock()
	delete(d.networks, nid)
	d.Unlock()
***REMOVED***

func (d *driver) network(nid string) *network ***REMOVED***
	d.Lock()
	defer d.Unlock()
	return d.networks[nid]
***REMOVED***

// func (n *network) restoreNetworkEndpoints() error ***REMOVED***
// 	logrus.Infof("Restoring endpoints for overlay network: %s", n.id)

// 	hnsresponse, err := hcsshim.HNSListEndpointRequest("GET", "", "")
// 	if err != nil ***REMOVED***
// 		return err
// 	***REMOVED***

// 	for _, endpoint := range hnsresponse ***REMOVED***
// 		if endpoint.VirtualNetwork != n.hnsID ***REMOVED***
// 			continue
// 		***REMOVED***

// 		ep := n.convertToOverlayEndpoint(&endpoint)

// 		if ep != nil ***REMOVED***
// 			logrus.Debugf("Restored endpoint:%s Remote:%t", ep.id, ep.remote)
// 			n.addEndpoint(ep)
// 		***REMOVED***
// 	***REMOVED***

// 	return nil
// ***REMOVED***

func (n *network) convertToOverlayEndpoint(v *hcsshim.HNSEndpoint) *endpoint ***REMOVED***
	ep := &endpoint***REMOVED***
		id:        v.Name,
		profileID: v.Id,
		nid:       n.id,
		remote:    v.IsRemoteEndpoint,
	***REMOVED***

	mac, err := net.ParseMAC(v.MacAddress)

	if err != nil ***REMOVED***
		return nil
	***REMOVED***

	ep.mac = mac
	ep.addr = &net.IPNet***REMOVED***
		IP:   v.IPAddress,
		Mask: net.CIDRMask(32, 32),
	***REMOVED***

	return ep
***REMOVED***

func (d *driver) createHnsNetwork(n *network) error ***REMOVED***

	subnets := []hcsshim.Subnet***REMOVED******REMOVED***

	for _, s := range n.subnets ***REMOVED***
		subnet := hcsshim.Subnet***REMOVED***
			AddressPrefix: s.subnetIP.String(),
		***REMOVED***

		if s.gwIP != nil ***REMOVED***
			subnet.GatewayAddress = s.gwIP.String()
		***REMOVED***

		vsidPolicy, err := json.Marshal(hcsshim.VsidPolicy***REMOVED***
			Type: "VSID",
			VSID: uint(s.vni),
		***REMOVED***)

		if err != nil ***REMOVED***
			return err
		***REMOVED***

		subnet.Policies = append(subnet.Policies, vsidPolicy)
		subnets = append(subnets, subnet)
	***REMOVED***

	network := &hcsshim.HNSNetwork***REMOVED***
		Name:               n.name,
		Type:               d.Type(),
		Subnets:            subnets,
		NetworkAdapterName: n.interfaceName,
		AutomaticDNS:       true,
	***REMOVED***

	configurationb, err := json.Marshal(network)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	configuration := string(configurationb)
	logrus.Infof("HNSNetwork Request =%v", configuration)

	hnsresponse, err := hcsshim.HNSNetworkRequest("POST", "", configuration)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	n.hnsID = hnsresponse.Id
	n.providerAddress = hnsresponse.ManagementIP

	return nil
***REMOVED***

// contains return true if the passed ip belongs to one the network's
// subnets
func (n *network) contains(ip net.IP) bool ***REMOVED***
	for _, s := range n.subnets ***REMOVED***
		if s.subnetIP.Contains(ip) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

// getSubnetforIP returns the subnet to which the given IP belongs
func (n *network) getSubnetforIP(ip *net.IPNet) *subnet ***REMOVED***
	for _, s := range n.subnets ***REMOVED***
		// first check if the mask lengths are the same
		i, _ := s.subnetIP.Mask.Size()
		j, _ := ip.Mask.Size()
		if i != j ***REMOVED***
			continue
		***REMOVED***
		if s.subnetIP.Contains(ip.IP) ***REMOVED***
			return s
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// getMatchingSubnet return the network's subnet that matches the input
func (n *network) getMatchingSubnet(ip *net.IPNet) *subnet ***REMOVED***
	if ip == nil ***REMOVED***
		return nil
	***REMOVED***
	for _, s := range n.subnets ***REMOVED***
		// first check if the mask lengths are the same
		i, _ := s.subnetIP.Mask.Size()
		j, _ := ip.Mask.Size()
		if i != j ***REMOVED***
			continue
		***REMOVED***
		if s.subnetIP.IP.Equal(ip.IP) ***REMOVED***
			return s
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
