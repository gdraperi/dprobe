package hcsshim

import (
	"encoding/json"
	"net"

	"github.com/sirupsen/logrus"
)

// Subnet is assoicated with a network and represents a list
// of subnets available to the network
type Subnet struct ***REMOVED***
	AddressPrefix  string            `json:",omitempty"`
	GatewayAddress string            `json:",omitempty"`
	Policies       []json.RawMessage `json:",omitempty"`
***REMOVED***

// MacPool is assoicated with a network and represents a list
// of macaddresses available to the network
type MacPool struct ***REMOVED***
	StartMacAddress string `json:",omitempty"`
	EndMacAddress   string `json:",omitempty"`
***REMOVED***

// HNSNetwork represents a network in HNS
type HNSNetwork struct ***REMOVED***
	Id                   string            `json:"ID,omitempty"`
	Name                 string            `json:",omitempty"`
	Type                 string            `json:",omitempty"`
	NetworkAdapterName   string            `json:",omitempty"`
	SourceMac            string            `json:",omitempty"`
	Policies             []json.RawMessage `json:",omitempty"`
	MacPools             []MacPool         `json:",omitempty"`
	Subnets              []Subnet          `json:",omitempty"`
	DNSSuffix            string            `json:",omitempty"`
	DNSServerList        string            `json:",omitempty"`
	DNSServerCompartment uint32            `json:",omitempty"`
	ManagementIP         string            `json:",omitempty"`
	AutomaticDNS         bool              `json:",omitempty"`
***REMOVED***

type hnsNetworkResponse struct ***REMOVED***
	Success bool
	Error   string
	Output  HNSNetwork
***REMOVED***

type hnsResponse struct ***REMOVED***
	Success bool
	Error   string
	Output  json.RawMessage
***REMOVED***

// HNSNetworkRequest makes a call into HNS to update/query a single network
func HNSNetworkRequest(method, path, request string) (*HNSNetwork, error) ***REMOVED***
	var network HNSNetwork
	err := hnsCall(method, "/networks/"+path, request, &network)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &network, nil
***REMOVED***

// HNSListNetworkRequest makes a HNS call to query the list of available networks
func HNSListNetworkRequest(method, path, request string) ([]HNSNetwork, error) ***REMOVED***
	var network []HNSNetwork
	err := hnsCall(method, "/networks/"+path, request, &network)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return network, nil
***REMOVED***

// GetHNSNetworkByID
func GetHNSNetworkByID(networkID string) (*HNSNetwork, error) ***REMOVED***
	return HNSNetworkRequest("GET", networkID, "")
***REMOVED***

// GetHNSNetworkName filtered by Name
func GetHNSNetworkByName(networkName string) (*HNSNetwork, error) ***REMOVED***
	hsnnetworks, err := HNSListNetworkRequest("GET", "", "")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	for _, hnsnetwork := range hsnnetworks ***REMOVED***
		if hnsnetwork.Name == networkName ***REMOVED***
			return &hnsnetwork, nil
		***REMOVED***
	***REMOVED***
	return nil, NetworkNotFoundError***REMOVED***NetworkName: networkName***REMOVED***
***REMOVED***

// Create Network by sending NetworkRequest to HNS.
func (network *HNSNetwork) Create() (*HNSNetwork, error) ***REMOVED***
	operation := "Create"
	title := "HCSShim::HNSNetwork::" + operation
	logrus.Debugf(title+" id=%s", network.Id)

	jsonString, err := json.Marshal(network)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return HNSNetworkRequest("POST", "", string(jsonString))
***REMOVED***

// Delete Network by sending NetworkRequest to HNS
func (network *HNSNetwork) Delete() (*HNSNetwork, error) ***REMOVED***
	operation := "Delete"
	title := "HCSShim::HNSNetwork::" + operation
	logrus.Debugf(title+" id=%s", network.Id)

	return HNSNetworkRequest("DELETE", network.Id, "")
***REMOVED***

// Creates an endpoint on the Network.
func (network *HNSNetwork) NewEndpoint(ipAddress net.IP, macAddress net.HardwareAddr) *HNSEndpoint ***REMOVED***
	return &HNSEndpoint***REMOVED***
		VirtualNetwork: network.Id,
		IPAddress:      ipAddress,
		MacAddress:     string(macAddress),
	***REMOVED***
***REMOVED***

func (network *HNSNetwork) CreateEndpoint(endpoint *HNSEndpoint) (*HNSEndpoint, error) ***REMOVED***
	operation := "CreateEndpoint"
	title := "HCSShim::HNSNetwork::" + operation
	logrus.Debugf(title+" id=%s, endpointId=%s", network.Id, endpoint.Id)

	endpoint.VirtualNetwork = network.Id
	return endpoint.Create()
***REMOVED***

func (network *HNSNetwork) CreateRemoteEndpoint(endpoint *HNSEndpoint) (*HNSEndpoint, error) ***REMOVED***
	operation := "CreateRemoteEndpoint"
	title := "HCSShim::HNSNetwork::" + operation
	logrus.Debugf(title+" id=%s", network.Id)
	endpoint.IsRemoteEndpoint = true
	return network.CreateEndpoint(endpoint)
***REMOVED***
