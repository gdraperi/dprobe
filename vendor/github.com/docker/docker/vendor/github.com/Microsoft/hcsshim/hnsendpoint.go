package hcsshim

import (
	"encoding/json"
	"net"

	"github.com/sirupsen/logrus"
)

// HNSEndpoint represents a network endpoint in HNS
type HNSEndpoint struct ***REMOVED***
	Id                 string            `json:"ID,omitempty"`
	Name               string            `json:",omitempty"`
	VirtualNetwork     string            `json:",omitempty"`
	VirtualNetworkName string            `json:",omitempty"`
	Policies           []json.RawMessage `json:",omitempty"`
	MacAddress         string            `json:",omitempty"`
	IPAddress          net.IP            `json:",omitempty"`
	DNSSuffix          string            `json:",omitempty"`
	DNSServerList      string            `json:",omitempty"`
	GatewayAddress     string            `json:",omitempty"`
	EnableInternalDNS  bool              `json:",omitempty"`
	DisableICC         bool              `json:",omitempty"`
	PrefixLength       uint8             `json:",omitempty"`
	IsRemoteEndpoint   bool              `json:",omitempty"`
***REMOVED***

//SystemType represents the type of the system on which actions are done
type SystemType string

// SystemType const
const (
	ContainerType      SystemType = "Container"
	VirtualMachineType SystemType = "VirtualMachine"
	HostType           SystemType = "Host"
)

// EndpointAttachDetachRequest is the structure used to send request to the container to modify the system
// Supported resource types are Network and Request Types are Add/Remove
type EndpointAttachDetachRequest struct ***REMOVED***
	ContainerID    string     `json:"ContainerId,omitempty"`
	SystemType     SystemType `json:"SystemType"`
	CompartmentID  uint16     `json:"CompartmentId,omitempty"`
	VirtualNICName string     `json:"VirtualNicName,omitempty"`
***REMOVED***

// EndpointResquestResponse is object to get the endpoint request response
type EndpointResquestResponse struct ***REMOVED***
	Success bool
	Error   string
***REMOVED***

// HNSEndpointRequest makes a HNS call to modify/query a network endpoint
func HNSEndpointRequest(method, path, request string) (*HNSEndpoint, error) ***REMOVED***
	endpoint := &HNSEndpoint***REMOVED******REMOVED***
	err := hnsCall(method, "/endpoints/"+path, request, &endpoint)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return endpoint, nil
***REMOVED***

// HNSListEndpointRequest makes a HNS call to query the list of available endpoints
func HNSListEndpointRequest() ([]HNSEndpoint, error) ***REMOVED***
	var endpoint []HNSEndpoint
	err := hnsCall("GET", "/endpoints/", "", &endpoint)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return endpoint, nil
***REMOVED***

// HotAttachEndpoint makes a HCS Call to attach the endpoint to the container
func HotAttachEndpoint(containerID string, endpointID string) error ***REMOVED***
	return modifyNetworkEndpoint(containerID, endpointID, Add)
***REMOVED***

// HotDetachEndpoint makes a HCS Call to detach the endpoint from the container
func HotDetachEndpoint(containerID string, endpointID string) error ***REMOVED***
	return modifyNetworkEndpoint(containerID, endpointID, Remove)
***REMOVED***

// ModifyContainer corresponding to the container id, by sending a request
func modifyContainer(id string, request *ResourceModificationRequestResponse) error ***REMOVED***
	container, err := OpenContainer(id)
	if err != nil ***REMOVED***
		if IsNotExist(err) ***REMOVED***
			return ErrComputeSystemDoesNotExist
		***REMOVED***
		return getInnerError(err)
	***REMOVED***
	defer container.Close()
	err = container.Modify(request)
	if err != nil ***REMOVED***
		if IsNotSupported(err) ***REMOVED***
			return ErrPlatformNotSupported
		***REMOVED***
		return getInnerError(err)
	***REMOVED***

	return nil
***REMOVED***

func modifyNetworkEndpoint(containerID string, endpointID string, request RequestType) error ***REMOVED***
	requestMessage := &ResourceModificationRequestResponse***REMOVED***
		Resource: Network,
		Request:  request,
		Data:     endpointID,
	***REMOVED***
	err := modifyContainer(containerID, requestMessage)

	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

// GetHNSEndpointByID get the Endpoint by ID
func GetHNSEndpointByID(endpointID string) (*HNSEndpoint, error) ***REMOVED***
	return HNSEndpointRequest("GET", endpointID, "")
***REMOVED***

// GetHNSEndpointByName gets the endpoint filtered by Name
func GetHNSEndpointByName(endpointName string) (*HNSEndpoint, error) ***REMOVED***
	hnsResponse, err := HNSListEndpointRequest()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	for _, hnsEndpoint := range hnsResponse ***REMOVED***
		if hnsEndpoint.Name == endpointName ***REMOVED***
			return &hnsEndpoint, nil
		***REMOVED***
	***REMOVED***
	return nil, EndpointNotFoundError***REMOVED***EndpointName: endpointName***REMOVED***
***REMOVED***

// Create Endpoint by sending EndpointRequest to HNS. TODO: Create a separate HNS interface to place all these methods
func (endpoint *HNSEndpoint) Create() (*HNSEndpoint, error) ***REMOVED***
	operation := "Create"
	title := "HCSShim::HNSEndpoint::" + operation
	logrus.Debugf(title+" id=%s", endpoint.Id)

	jsonString, err := json.Marshal(endpoint)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return HNSEndpointRequest("POST", "", string(jsonString))
***REMOVED***

// Delete Endpoint by sending EndpointRequest to HNS
func (endpoint *HNSEndpoint) Delete() (*HNSEndpoint, error) ***REMOVED***
	operation := "Delete"
	title := "HCSShim::HNSEndpoint::" + operation
	logrus.Debugf(title+" id=%s", endpoint.Id)

	return HNSEndpointRequest("DELETE", endpoint.Id, "")
***REMOVED***

// Update Endpoint
func (endpoint *HNSEndpoint) Update() (*HNSEndpoint, error) ***REMOVED***
	operation := "Update"
	title := "HCSShim::HNSEndpoint::" + operation
	logrus.Debugf(title+" id=%s", endpoint.Id)
	jsonString, err := json.Marshal(endpoint)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	err = hnsCall("POST", "/endpoints/"+endpoint.Id, string(jsonString), &endpoint)

	return endpoint, err
***REMOVED***

// ContainerHotAttach attaches an endpoint to a running container
func (endpoint *HNSEndpoint) ContainerHotAttach(containerID string) error ***REMOVED***
	operation := "ContainerHotAttach"
	title := "HCSShim::HNSEndpoint::" + operation
	logrus.Debugf(title+" id=%s, containerId=%s", endpoint.Id, containerID)

	return modifyNetworkEndpoint(containerID, endpoint.Id, Add)
***REMOVED***

// ContainerHotDetach detaches an endpoint from a running container
func (endpoint *HNSEndpoint) ContainerHotDetach(containerID string) error ***REMOVED***
	operation := "ContainerHotDetach"
	title := "HCSShim::HNSEndpoint::" + operation
	logrus.Debugf(title+" id=%s, containerId=%s", endpoint.Id, containerID)

	return modifyNetworkEndpoint(containerID, endpoint.Id, Remove)
***REMOVED***

// ApplyACLPolicy applies a set of ACL Policies on the Endpoint
func (endpoint *HNSEndpoint) ApplyACLPolicy(policies ...*ACLPolicy) error ***REMOVED***
	operation := "ApplyACLPolicy"
	title := "HCSShim::HNSEndpoint::" + operation
	logrus.Debugf(title+" id=%s", endpoint.Id)

	for _, policy := range policies ***REMOVED***
		if policy == nil ***REMOVED***
			continue
		***REMOVED***
		jsonString, err := json.Marshal(policy)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		endpoint.Policies = append(endpoint.Policies, jsonString)
	***REMOVED***

	_, err := endpoint.Update()
	return err
***REMOVED***

// ContainerAttach attaches an endpoint to container
func (endpoint *HNSEndpoint) ContainerAttach(containerID string, compartmentID uint16) error ***REMOVED***
	operation := "ContainerAttach"
	title := "HCSShim::HNSEndpoint::" + operation
	logrus.Debugf(title+" id=%s", endpoint.Id)

	requestMessage := &EndpointAttachDetachRequest***REMOVED***
		ContainerID:   containerID,
		CompartmentID: compartmentID,
		SystemType:    ContainerType,
	***REMOVED***
	response := &EndpointResquestResponse***REMOVED******REMOVED***
	jsonString, err := json.Marshal(requestMessage)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return hnsCall("POST", "/endpoints/"+endpoint.Id+"/attach", string(jsonString), &response)
***REMOVED***

// ContainerDetach detaches an endpoint from container
func (endpoint *HNSEndpoint) ContainerDetach(containerID string) error ***REMOVED***
	operation := "ContainerDetach"
	title := "HCSShim::HNSEndpoint::" + operation
	logrus.Debugf(title+" id=%s", endpoint.Id)

	requestMessage := &EndpointAttachDetachRequest***REMOVED***
		ContainerID: containerID,
		SystemType:  ContainerType,
	***REMOVED***
	response := &EndpointResquestResponse***REMOVED******REMOVED***

	jsonString, err := json.Marshal(requestMessage)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return hnsCall("POST", "/endpoints/"+endpoint.Id+"/detach", string(jsonString), &response)
***REMOVED***

// HostAttach attaches a nic on the host
func (endpoint *HNSEndpoint) HostAttach(compartmentID uint16) error ***REMOVED***
	operation := "HostAttach"
	title := "HCSShim::HNSEndpoint::" + operation
	logrus.Debugf(title+" id=%s", endpoint.Id)
	requestMessage := &EndpointAttachDetachRequest***REMOVED***
		CompartmentID: compartmentID,
		SystemType:    HostType,
	***REMOVED***
	response := &EndpointResquestResponse***REMOVED******REMOVED***

	jsonString, err := json.Marshal(requestMessage)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return hnsCall("POST", "/endpoints/"+endpoint.Id+"/attach", string(jsonString), &response)

***REMOVED***

// HostDetach detaches a nic on the host
func (endpoint *HNSEndpoint) HostDetach() error ***REMOVED***
	operation := "HostDetach"
	title := "HCSShim::HNSEndpoint::" + operation
	logrus.Debugf(title+" id=%s", endpoint.Id)
	requestMessage := &EndpointAttachDetachRequest***REMOVED***
		SystemType: HostType,
	***REMOVED***
	response := &EndpointResquestResponse***REMOVED******REMOVED***

	jsonString, err := json.Marshal(requestMessage)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return hnsCall("POST", "/endpoints/"+endpoint.Id+"/detach", string(jsonString), &response)
***REMOVED***

// VirtualMachineNICAttach attaches a endpoint to a virtual machine
func (endpoint *HNSEndpoint) VirtualMachineNICAttach(virtualMachineNICName string) error ***REMOVED***
	operation := "VirtualMachineNicAttach"
	title := "HCSShim::HNSEndpoint::" + operation
	logrus.Debugf(title+" id=%s", endpoint.Id)
	requestMessage := &EndpointAttachDetachRequest***REMOVED***
		VirtualNICName: virtualMachineNICName,
		SystemType:     VirtualMachineType,
	***REMOVED***
	response := &EndpointResquestResponse***REMOVED******REMOVED***

	jsonString, err := json.Marshal(requestMessage)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return hnsCall("POST", "/endpoints/"+endpoint.Id+"/attach", string(jsonString), &response)
***REMOVED***

// VirtualMachineNICDetach detaches a endpoint  from a virtual machine
func (endpoint *HNSEndpoint) VirtualMachineNICDetach() error ***REMOVED***
	operation := "VirtualMachineNicDetach"
	title := "HCSShim::HNSEndpoint::" + operation
	logrus.Debugf(title+" id=%s", endpoint.Id)

	requestMessage := &EndpointAttachDetachRequest***REMOVED***
		SystemType: VirtualMachineType,
	***REMOVED***
	response := &EndpointResquestResponse***REMOVED******REMOVED***

	jsonString, err := json.Marshal(requestMessage)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return hnsCall("POST", "/endpoints/"+endpoint.Id+"/detach", string(jsonString), &response)
***REMOVED***
