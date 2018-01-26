// +build windows

// Shim for the Host Network Service (HNS) to manage networking for
// Windows Server containers and Hyper-V containers. This module
// is a basic libnetwork driver that passes all the calls to HNS
// It implements the 4 networking modes supported by HNS L2Bridge,
// L2Tunnel, NAT and Transparent(DHCP)
//
// The network are stored in memory and docker daemon ensures discovering
// and loading these networks on startup

package windows

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/Microsoft/hcsshim"
	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/discoverapi"
	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/netlabel"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

// networkConfiguration for network specific configuration
type networkConfiguration struct ***REMOVED***
	ID                 string
	Type               string
	Name               string
	HnsID              string
	RDID               string
	VLAN               uint
	VSID               uint
	DNSServers         string
	MacPools           []hcsshim.MacPool
	DNSSuffix          string
	SourceMac          string
	NetworkAdapterName string
	dbIndex            uint64
	dbExists           bool
	DisableGatewayDNS  bool
***REMOVED***

// endpointConfiguration represents the user specified configuration for the sandbox endpoint
type endpointOption struct ***REMOVED***
	MacAddress  net.HardwareAddr
	QosPolicies []types.QosPolicy
	DNSServers  []string
	DisableDNS  bool
	DisableICC  bool
***REMOVED***

// EndpointConnectivity stores the port bindings and exposed ports that the user has specified in epOptions.
type EndpointConnectivity struct ***REMOVED***
	PortBindings []types.PortBinding
	ExposedPorts []types.TransportPort
***REMOVED***

type hnsEndpoint struct ***REMOVED***
	id        string
	nid       string
	profileID string
	Type      string
	//Note: Currently, the sandboxID is the same as the containerID since windows does
	//not expose the sandboxID.
	//In the future, windows will support a proper sandboxID that is different
	//than the containerID.
	//Therefore, we are using sandboxID now, so that we won't have to change this code
	//when windows properly supports a sandboxID.
	sandboxID      string
	macAddress     net.HardwareAddr
	epOption       *endpointOption       // User specified parameters
	epConnectivity *EndpointConnectivity // User specified parameters
	portMapping    []types.PortBinding   // Operation port bindings
	addr           *net.IPNet
	gateway        net.IP
	dbIndex        uint64
	dbExists       bool
***REMOVED***

type hnsNetwork struct ***REMOVED***
	id        string
	created   bool
	config    *networkConfiguration
	endpoints map[string]*hnsEndpoint // key: endpoint id
	driver    *driver                 // The network's driver
	sync.Mutex
***REMOVED***

type driver struct ***REMOVED***
	name     string
	networks map[string]*hnsNetwork
	store    datastore.DataStore
	sync.Mutex
***REMOVED***

const (
	errNotFound = "HNS failed with error : The object identifier does not represent a valid object. "
)

// IsBuiltinLocalDriver validates if network-type is a builtin local-scoped driver
func IsBuiltinLocalDriver(networkType string) bool ***REMOVED***
	if "l2bridge" == networkType || "l2tunnel" == networkType || "nat" == networkType || "ics" == networkType || "transparent" == networkType ***REMOVED***
		return true
	***REMOVED***

	return false
***REMOVED***

// New constructs a new bridge driver
func newDriver(networkType string) *driver ***REMOVED***
	return &driver***REMOVED***name: networkType, networks: map[string]*hnsNetwork***REMOVED******REMOVED******REMOVED***
***REMOVED***

// GetInit returns an initializer for the given network type
func GetInit(networkType string) func(dc driverapi.DriverCallback, config map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	return func(dc driverapi.DriverCallback, config map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
		if !IsBuiltinLocalDriver(networkType) ***REMOVED***
			return types.BadRequestErrorf("Network type not supported: %s", networkType)
		***REMOVED***

		d := newDriver(networkType)

		err := d.initStore(config)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		return dc.RegisterDriver(networkType, d, driverapi.Capability***REMOVED***
			DataScope:         datastore.LocalScope,
			ConnectivityScope: datastore.LocalScope,
		***REMOVED***)
	***REMOVED***
***REMOVED***

func (d *driver) getNetwork(id string) (*hnsNetwork, error) ***REMOVED***
	d.Lock()
	defer d.Unlock()

	if nw, ok := d.networks[id]; ok ***REMOVED***
		return nw, nil
	***REMOVED***

	return nil, types.NotFoundErrorf("network not found: %s", id)
***REMOVED***

func (n *hnsNetwork) getEndpoint(eid string) (*hnsEndpoint, error) ***REMOVED***
	n.Lock()
	defer n.Unlock()

	if ep, ok := n.endpoints[eid]; ok ***REMOVED***
		return ep, nil
	***REMOVED***

	return nil, types.NotFoundErrorf("Endpoint not found: %s", eid)
***REMOVED***

func (d *driver) parseNetworkOptions(id string, genericOptions map[string]string) (*networkConfiguration, error) ***REMOVED***
	config := &networkConfiguration***REMOVED***Type: d.name***REMOVED***

	for label, value := range genericOptions ***REMOVED***
		switch label ***REMOVED***
		case NetworkName:
			config.Name = value
		case HNSID:
			config.HnsID = value
		case RoutingDomain:
			config.RDID = value
		case Interface:
			config.NetworkAdapterName = value
		case DNSSuffix:
			config.DNSSuffix = value
		case DNSServers:
			config.DNSServers = value
		case DisableGatewayDNS:
			b, err := strconv.ParseBool(value)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			config.DisableGatewayDNS = b
		case MacPool:
			config.MacPools = make([]hcsshim.MacPool, 0)
			s := strings.Split(value, ",")
			if len(s)%2 != 0 ***REMOVED***
				return nil, types.BadRequestErrorf("Invalid mac pool. You must specify both a start range and an end range")
			***REMOVED***
			for i := 0; i < len(s)-1; i += 2 ***REMOVED***
				config.MacPools = append(config.MacPools, hcsshim.MacPool***REMOVED***
					StartMacAddress: s[i],
					EndMacAddress:   s[i+1],
				***REMOVED***)
			***REMOVED***
		case VLAN:
			vlan, err := strconv.ParseUint(value, 10, 32)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			config.VLAN = uint(vlan)
		case VSID:
			vsid, err := strconv.ParseUint(value, 10, 32)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			config.VSID = uint(vsid)
		***REMOVED***
	***REMOVED***

	config.ID = id
	config.Type = d.name
	return config, nil
***REMOVED***

func (c *networkConfiguration) processIPAM(id string, ipamV4Data, ipamV6Data []driverapi.IPAMData) error ***REMOVED***
	if len(ipamV6Data) > 0 ***REMOVED***
		return types.ForbiddenErrorf("windowsshim driver doesn't support v6 subnets")
	***REMOVED***

	if len(ipamV4Data) == 0 ***REMOVED***
		return types.BadRequestErrorf("network %s requires ipv4 configuration", id)
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) EventNotify(etype driverapi.EventType, nid, tableName, key string, value []byte) ***REMOVED***
***REMOVED***

func (d *driver) DecodeTableEntry(tablename string, key string, value []byte) (string, map[string]string) ***REMOVED***
	return "", nil
***REMOVED***

func (d *driver) createNetwork(config *networkConfiguration) error ***REMOVED***
	network := &hnsNetwork***REMOVED***
		id:        config.ID,
		endpoints: make(map[string]*hnsEndpoint),
		config:    config,
		driver:    d,
	***REMOVED***

	d.Lock()
	d.networks[config.ID] = network
	d.Unlock()

	return nil
***REMOVED***

// Create a new network
func (d *driver) CreateNetwork(id string, option map[string]interface***REMOVED******REMOVED***, nInfo driverapi.NetworkInfo, ipV4Data, ipV6Data []driverapi.IPAMData) error ***REMOVED***
	if _, err := d.getNetwork(id); err == nil ***REMOVED***
		return types.ForbiddenErrorf("network %s exists", id)
	***REMOVED***

	genData, ok := option[netlabel.GenericData].(map[string]string)
	if !ok ***REMOVED***
		return fmt.Errorf("Unknown generic data option")
	***REMOVED***

	// Parse and validate the config. It should not conflict with existing networks' config
	config, err := d.parseNetworkOptions(id, genData)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = config.processIPAM(id, ipV4Data, ipV6Data)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = d.createNetwork(config)

	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// A non blank hnsid indicates that the network was discovered
	// from HNS. No need to call HNS if this network was discovered
	// from HNS
	if config.HnsID == "" ***REMOVED***
		subnets := []hcsshim.Subnet***REMOVED******REMOVED***

		for _, ipData := range ipV4Data ***REMOVED***
			subnet := hcsshim.Subnet***REMOVED***
				AddressPrefix: ipData.Pool.String(),
			***REMOVED***

			if ipData.Gateway != nil ***REMOVED***
				subnet.GatewayAddress = ipData.Gateway.IP.String()
			***REMOVED***

			subnets = append(subnets, subnet)
		***REMOVED***

		network := &hcsshim.HNSNetwork***REMOVED***
			Name:               config.Name,
			Type:               d.name,
			Subnets:            subnets,
			DNSServerList:      config.DNSServers,
			DNSSuffix:          config.DNSSuffix,
			MacPools:           config.MacPools,
			SourceMac:          config.SourceMac,
			NetworkAdapterName: config.NetworkAdapterName,
		***REMOVED***

		if config.VLAN != 0 ***REMOVED***
			vlanPolicy, err := json.Marshal(hcsshim.VlanPolicy***REMOVED***
				Type: "VLAN",
				VLAN: config.VLAN,
			***REMOVED***)

			if err != nil ***REMOVED***
				return err
			***REMOVED***
			network.Policies = append(network.Policies, vlanPolicy)
		***REMOVED***

		if config.VSID != 0 ***REMOVED***
			vsidPolicy, err := json.Marshal(hcsshim.VsidPolicy***REMOVED***
				Type: "VSID",
				VSID: config.VSID,
			***REMOVED***)

			if err != nil ***REMOVED***
				return err
			***REMOVED***
			network.Policies = append(network.Policies, vsidPolicy)
		***REMOVED***

		if network.Name == "" ***REMOVED***
			network.Name = id
		***REMOVED***

		configurationb, err := json.Marshal(network)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		configuration := string(configurationb)
		logrus.Debugf("HNSNetwork Request =%v Address Space=%v", configuration, subnets)

		hnsresponse, err := hcsshim.HNSNetworkRequest("POST", "", configuration)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		config.HnsID = hnsresponse.Id
		genData[HNSID] = config.HnsID
	***REMOVED***

	n, err := d.getNetwork(id)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	n.created = true
	return d.storeUpdate(config)
***REMOVED***

func (d *driver) DeleteNetwork(nid string) error ***REMOVED***
	n, err := d.getNetwork(nid)
	if err != nil ***REMOVED***
		return types.InternalMaskableErrorf("%s", err)
	***REMOVED***

	n.Lock()
	config := n.config
	n.Unlock()

	if n.created ***REMOVED***
		_, err = hcsshim.HNSNetworkRequest("DELETE", config.HnsID, "")
		if err != nil && err.Error() != errNotFound ***REMOVED***
			return types.ForbiddenErrorf(err.Error())
		***REMOVED***
	***REMOVED***

	d.Lock()
	delete(d.networks, nid)
	d.Unlock()

	// delele endpoints belong to this network
	for _, ep := range n.endpoints ***REMOVED***
		if err := d.storeDelete(ep); err != nil ***REMOVED***
			logrus.Warnf("Failed to remove bridge endpoint %s from store: %v", ep.id[0:7], err)
		***REMOVED***
	***REMOVED***

	return d.storeDelete(config)
***REMOVED***

func convertQosPolicies(qosPolicies []types.QosPolicy) ([]json.RawMessage, error) ***REMOVED***
	var qps []json.RawMessage

	// Enumerate through the qos policies specified by the user and convert
	// them into the internal structure matching the JSON blob that can be
	// understood by the HCS.
	for _, elem := range qosPolicies ***REMOVED***
		encodedPolicy, err := json.Marshal(hcsshim.QosPolicy***REMOVED***
			Type: "QOS",
			MaximumOutgoingBandwidthInBytes: elem.MaxEgressBandwidth,
		***REMOVED***)

		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		qps = append(qps, encodedPolicy)
	***REMOVED***
	return qps, nil
***REMOVED***

// ConvertPortBindings converts PortBindings to JSON for HNS request
func ConvertPortBindings(portBindings []types.PortBinding) ([]json.RawMessage, error) ***REMOVED***
	var pbs []json.RawMessage

	// Enumerate through the port bindings specified by the user and convert
	// them into the internal structure matching the JSON blob that can be
	// understood by the HCS.
	for _, elem := range portBindings ***REMOVED***
		proto := strings.ToUpper(elem.Proto.String())
		if proto != "TCP" && proto != "UDP" ***REMOVED***
			return nil, fmt.Errorf("invalid protocol %s", elem.Proto.String())
		***REMOVED***

		if elem.HostPort != elem.HostPortEnd ***REMOVED***
			return nil, fmt.Errorf("Windows does not support more than one host port in NAT settings")
		***REMOVED***

		if len(elem.HostIP) != 0 ***REMOVED***
			return nil, fmt.Errorf("Windows does not support host IP addresses in NAT settings")
		***REMOVED***

		encodedPolicy, err := json.Marshal(hcsshim.NatPolicy***REMOVED***
			Type:         "NAT",
			ExternalPort: elem.HostPort,
			InternalPort: elem.Port,
			Protocol:     elem.Proto.String(),
		***REMOVED***)

		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		pbs = append(pbs, encodedPolicy)
	***REMOVED***
	return pbs, nil
***REMOVED***

// ParsePortBindingPolicies parses HNS endpoint response message to PortBindings
func ParsePortBindingPolicies(policies []json.RawMessage) ([]types.PortBinding, error) ***REMOVED***
	var bindings []types.PortBinding
	hcsPolicy := &hcsshim.NatPolicy***REMOVED******REMOVED***

	for _, elem := range policies ***REMOVED***

		if err := json.Unmarshal([]byte(elem), &hcsPolicy); err != nil || hcsPolicy.Type != "NAT" ***REMOVED***
			continue
		***REMOVED***

		binding := types.PortBinding***REMOVED***
			HostPort:    hcsPolicy.ExternalPort,
			HostPortEnd: hcsPolicy.ExternalPort,
			Port:        hcsPolicy.InternalPort,
			Proto:       types.ParseProtocol(hcsPolicy.Protocol),
			HostIP:      net.IPv4(0, 0, 0, 0),
		***REMOVED***

		bindings = append(bindings, binding)
	***REMOVED***

	return bindings, nil
***REMOVED***

func parseEndpointOptions(epOptions map[string]interface***REMOVED******REMOVED***) (*endpointOption, error) ***REMOVED***
	if epOptions == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	ec := &endpointOption***REMOVED******REMOVED***

	if opt, ok := epOptions[netlabel.MacAddress]; ok ***REMOVED***
		if mac, ok := opt.(net.HardwareAddr); ok ***REMOVED***
			ec.MacAddress = mac
		***REMOVED*** else ***REMOVED***
			return nil, fmt.Errorf("Invalid endpoint configuration")
		***REMOVED***
	***REMOVED***

	if opt, ok := epOptions[QosPolicies]; ok ***REMOVED***
		if policies, ok := opt.([]types.QosPolicy); ok ***REMOVED***
			ec.QosPolicies = policies
		***REMOVED*** else ***REMOVED***
			return nil, fmt.Errorf("Invalid endpoint configuration")
		***REMOVED***
	***REMOVED***

	if opt, ok := epOptions[netlabel.DNSServers]; ok ***REMOVED***
		if dns, ok := opt.([]string); ok ***REMOVED***
			ec.DNSServers = dns
		***REMOVED*** else ***REMOVED***
			return nil, fmt.Errorf("Invalid endpoint configuration")
		***REMOVED***
	***REMOVED***

	if opt, ok := epOptions[DisableICC]; ok ***REMOVED***
		if disableICC, ok := opt.(bool); ok ***REMOVED***
			ec.DisableICC = disableICC
		***REMOVED*** else ***REMOVED***
			return nil, fmt.Errorf("Invalid endpoint configuration")
		***REMOVED***
	***REMOVED***

	if opt, ok := epOptions[DisableDNS]; ok ***REMOVED***
		if disableDNS, ok := opt.(bool); ok ***REMOVED***
			ec.DisableDNS = disableDNS
		***REMOVED*** else ***REMOVED***
			return nil, fmt.Errorf("Invalid endpoint configuration")
		***REMOVED***
	***REMOVED***

	return ec, nil
***REMOVED***

// ParseEndpointConnectivity parses options passed to CreateEndpoint, specifically port bindings, and store in a endpointConnectivity object.
func ParseEndpointConnectivity(epOptions map[string]interface***REMOVED******REMOVED***) (*EndpointConnectivity, error) ***REMOVED***
	if epOptions == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	ec := &EndpointConnectivity***REMOVED******REMOVED***

	if opt, ok := epOptions[netlabel.PortMap]; ok ***REMOVED***
		if bs, ok := opt.([]types.PortBinding); ok ***REMOVED***
			ec.PortBindings = bs
		***REMOVED*** else ***REMOVED***
			return nil, fmt.Errorf("Invalid endpoint configuration")
		***REMOVED***
	***REMOVED***

	if opt, ok := epOptions[netlabel.ExposedPorts]; ok ***REMOVED***
		if ports, ok := opt.([]types.TransportPort); ok ***REMOVED***
			ec.ExposedPorts = ports
		***REMOVED*** else ***REMOVED***
			return nil, fmt.Errorf("Invalid endpoint configuration")
		***REMOVED***
	***REMOVED***
	return ec, nil
***REMOVED***

func (d *driver) CreateEndpoint(nid, eid string, ifInfo driverapi.InterfaceInfo, epOptions map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	n, err := d.getNetwork(nid)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Check if endpoint id is good and retrieve corresponding endpoint
	ep, err := n.getEndpoint(eid)
	if err == nil && ep != nil ***REMOVED***
		return driverapi.ErrEndpointExists(eid)
	***REMOVED***

	endpointStruct := &hcsshim.HNSEndpoint***REMOVED***
		VirtualNetwork: n.config.HnsID,
	***REMOVED***

	epOption, err := parseEndpointOptions(epOptions)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	epConnectivity, err := ParseEndpointConnectivity(epOptions)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	macAddress := ifInfo.MacAddress()
	// Use the macaddress if it was provided
	if macAddress != nil ***REMOVED***
		endpointStruct.MacAddress = strings.Replace(macAddress.String(), ":", "-", -1)
	***REMOVED***

	endpointStruct.Policies, err = ConvertPortBindings(epConnectivity.PortBindings)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	qosPolicies, err := convertQosPolicies(epOption.QosPolicies)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	endpointStruct.Policies = append(endpointStruct.Policies, qosPolicies...)

	if ifInfo.Address() != nil ***REMOVED***
		endpointStruct.IPAddress = ifInfo.Address().IP
	***REMOVED***

	endpointStruct.DNSServerList = strings.Join(epOption.DNSServers, ",")

	// overwrite the ep DisableDNS option if DisableGatewayDNS was set to true during the network creation option
	if n.config.DisableGatewayDNS ***REMOVED***
		logrus.Debugf("n.config.DisableGatewayDNS[%v] overwrites epOption.DisableDNS[%v]", n.config.DisableGatewayDNS, epOption.DisableDNS)
		epOption.DisableDNS = n.config.DisableGatewayDNS
	***REMOVED***

	if n.driver.name == "nat" && !epOption.DisableDNS ***REMOVED***
		logrus.Debugf("endpointStruct.EnableInternalDNS =[%v]", endpointStruct.EnableInternalDNS)
		endpointStruct.EnableInternalDNS = true
	***REMOVED***

	endpointStruct.DisableICC = epOption.DisableICC

	configurationb, err := json.Marshal(endpointStruct)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	hnsresponse, err := hcsshim.HNSEndpointRequest("POST", "", string(configurationb))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	mac, err := net.ParseMAC(hnsresponse.MacAddress)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// TODO For now the ip mask is not in the info generated by HNS
	endpoint := &hnsEndpoint***REMOVED***
		id:         eid,
		nid:        n.id,
		Type:       d.name,
		addr:       &net.IPNet***REMOVED***IP: hnsresponse.IPAddress, Mask: hnsresponse.IPAddress.DefaultMask()***REMOVED***,
		macAddress: mac,
	***REMOVED***

	if hnsresponse.GatewayAddress != "" ***REMOVED***
		endpoint.gateway = net.ParseIP(hnsresponse.GatewayAddress)
	***REMOVED***

	endpoint.profileID = hnsresponse.Id
	endpoint.epConnectivity = epConnectivity
	endpoint.epOption = epOption
	endpoint.portMapping, err = ParsePortBindingPolicies(hnsresponse.Policies)

	if err != nil ***REMOVED***
		hcsshim.HNSEndpointRequest("DELETE", hnsresponse.Id, "")
		return err
	***REMOVED***

	n.Lock()
	n.endpoints[eid] = endpoint
	n.Unlock()

	if ifInfo.Address() == nil ***REMOVED***
		ifInfo.SetIPAddress(endpoint.addr)
	***REMOVED***

	if macAddress == nil ***REMOVED***
		ifInfo.SetMacAddress(endpoint.macAddress)
	***REMOVED***

	if err = d.storeUpdate(endpoint); err != nil ***REMOVED***
		logrus.Errorf("Failed to save endpoint %s to store: %v", endpoint.id[0:7], err)
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) DeleteEndpoint(nid, eid string) error ***REMOVED***
	n, err := d.getNetwork(nid)
	if err != nil ***REMOVED***
		return types.InternalMaskableErrorf("%s", err)
	***REMOVED***

	ep, err := n.getEndpoint(eid)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	n.Lock()
	delete(n.endpoints, eid)
	n.Unlock()

	_, err = hcsshim.HNSEndpointRequest("DELETE", ep.profileID, "")
	if err != nil && err.Error() != errNotFound ***REMOVED***
		return err
	***REMOVED***

	if err := d.storeDelete(ep); err != nil ***REMOVED***
		logrus.Warnf("Failed to remove bridge endpoint %s from store: %v", ep.id[0:7], err)
	***REMOVED***
	return nil
***REMOVED***

func (d *driver) EndpointOperInfo(nid, eid string) (map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	network, err := d.getNetwork(nid)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ep, err := network.getEndpoint(eid)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	data := make(map[string]interface***REMOVED******REMOVED***, 1)
	if network.driver.name == "nat" ***REMOVED***
		data["AllowUnqualifiedDNSQuery"] = true
	***REMOVED***

	data["hnsid"] = ep.profileID
	if ep.epConnectivity.ExposedPorts != nil ***REMOVED***
		// Return a copy of the config data
		epc := make([]types.TransportPort, 0, len(ep.epConnectivity.ExposedPorts))
		for _, tp := range ep.epConnectivity.ExposedPorts ***REMOVED***
			epc = append(epc, tp.GetCopy())
		***REMOVED***
		data[netlabel.ExposedPorts] = epc
	***REMOVED***

	if ep.portMapping != nil ***REMOVED***
		// Return a copy of the operational data
		pmc := make([]types.PortBinding, 0, len(ep.portMapping))
		for _, pm := range ep.portMapping ***REMOVED***
			pmc = append(pmc, pm.GetCopy())
		***REMOVED***
		data[netlabel.PortMap] = pmc
	***REMOVED***

	if len(ep.macAddress) != 0 ***REMOVED***
		data[netlabel.MacAddress] = ep.macAddress
	***REMOVED***
	return data, nil
***REMOVED***

// Join method is invoked when a Sandbox is attached to an endpoint.
func (d *driver) Join(nid, eid string, sboxKey string, jinfo driverapi.JoinInfo, options map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	network, err := d.getNetwork(nid)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Ensure that the endpoint exists
	endpoint, err := network.getEndpoint(eid)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = jinfo.SetGateway(endpoint.gateway)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	endpoint.sandboxID = sboxKey

	err = hcsshim.HotAttachEndpoint(endpoint.sandboxID, endpoint.profileID)
	if err != nil ***REMOVED***
		// If container doesn't exists in hcs, do not throw error for hot add/remove
		if err != hcsshim.ErrComputeSystemDoesNotExist ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	jinfo.DisableGatewayService()
	return nil
***REMOVED***

// Leave method is invoked when a Sandbox detaches from an endpoint.
func (d *driver) Leave(nid, eid string) error ***REMOVED***
	network, err := d.getNetwork(nid)
	if err != nil ***REMOVED***
		return types.InternalMaskableErrorf("%s", err)
	***REMOVED***

	// Ensure that the endpoint exists
	endpoint, err := network.getEndpoint(eid)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = hcsshim.HotDetachEndpoint(endpoint.sandboxID, endpoint.profileID)
	if err != nil ***REMOVED***
		// If container doesn't exists in hcs, do not throw error for hot add/remove
		if err != hcsshim.ErrComputeSystemDoesNotExist ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (d *driver) ProgramExternalConnectivity(nid, eid string, options map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	return nil
***REMOVED***

func (d *driver) RevokeExternalConnectivity(nid, eid string) error ***REMOVED***
	return nil
***REMOVED***

func (d *driver) NetworkAllocate(id string, option map[string]string, ipV4Data, ipV6Data []driverapi.IPAMData) (map[string]string, error) ***REMOVED***
	return nil, types.NotImplementedErrorf("not implemented")
***REMOVED***

func (d *driver) NetworkFree(id string) error ***REMOVED***
	return types.NotImplementedErrorf("not implemented")
***REMOVED***

func (d *driver) Type() string ***REMOVED***
	return d.name
***REMOVED***

func (d *driver) IsBuiltIn() bool ***REMOVED***
	return true
***REMOVED***

// DiscoverNew is a notification for a new discovery event, such as a new node joining a cluster
func (d *driver) DiscoverNew(dType discoverapi.DiscoveryType, data interface***REMOVED******REMOVED***) error ***REMOVED***
	return nil
***REMOVED***

// DiscoverDelete is a notification for a discovery delete event, such as a node leaving a cluster
func (d *driver) DiscoverDelete(dType discoverapi.DiscoveryType, data interface***REMOVED******REMOVED***) error ***REMOVED***
	return nil
***REMOVED***
