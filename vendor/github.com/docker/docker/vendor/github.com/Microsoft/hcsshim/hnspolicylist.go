package hcsshim

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

// RoutePolicy is a structure defining schema for Route based Policy
type RoutePolicy struct ***REMOVED***
	Policy
	DestinationPrefix string `json:"DestinationPrefix,omitempty"`
	NextHop           string `json:"NextHop,omitempty"`
	EncapEnabled      bool   `json:"NeedEncap,omitempty"`
***REMOVED***

// ELBPolicy is a structure defining schema for ELB LoadBalancing based Policy
type ELBPolicy struct ***REMOVED***
	LBPolicy
	SourceVIP string   `json:"SourceVIP,omitempty"`
	VIPs      []string `json:"VIPs,omitempty"`
	ILB       bool     `json:"ILB,omitempty"`
***REMOVED***

// LBPolicy is a structure defining schema for LoadBalancing based Policy
type LBPolicy struct ***REMOVED***
	Policy
	Protocol     uint16 `json:"Protocol,omitempty"`
	InternalPort uint16
	ExternalPort uint16
***REMOVED***

// PolicyList is a structure defining schema for Policy list request
type PolicyList struct ***REMOVED***
	ID                 string            `json:"ID,omitempty"`
	EndpointReferences []string          `json:"References,omitempty"`
	Policies           []json.RawMessage `json:"Policies,omitempty"`
***REMOVED***

// HNSPolicyListRequest makes a call into HNS to update/query a single network
func HNSPolicyListRequest(method, path, request string) (*PolicyList, error) ***REMOVED***
	var policy PolicyList
	err := hnsCall(method, "/policylists/"+path, request, &policy)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &policy, nil
***REMOVED***

// HNSListPolicyListRequest gets all the policy list
func HNSListPolicyListRequest() ([]PolicyList, error) ***REMOVED***
	var plist []PolicyList
	err := hnsCall("GET", "/policylists/", "", &plist)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return plist, nil
***REMOVED***

// PolicyListRequest makes a HNS call to modify/query a network policy list
func PolicyListRequest(method, path, request string) (*PolicyList, error) ***REMOVED***
	policylist := &PolicyList***REMOVED******REMOVED***
	err := hnsCall(method, "/policylists/"+path, request, &policylist)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return policylist, nil
***REMOVED***

// GetPolicyListByID get the policy list by ID
func GetPolicyListByID(policyListID string) (*PolicyList, error) ***REMOVED***
	return PolicyListRequest("GET", policyListID, "")
***REMOVED***

// Create PolicyList by sending PolicyListRequest to HNS.
func (policylist *PolicyList) Create() (*PolicyList, error) ***REMOVED***
	operation := "Create"
	title := "HCSShim::PolicyList::" + operation
	logrus.Debugf(title+" id=%s", policylist.ID)
	jsonString, err := json.Marshal(policylist)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return PolicyListRequest("POST", "", string(jsonString))
***REMOVED***

// Delete deletes PolicyList
func (policylist *PolicyList) Delete() (*PolicyList, error) ***REMOVED***
	operation := "Delete"
	title := "HCSShim::PolicyList::" + operation
	logrus.Debugf(title+" id=%s", policylist.ID)

	return PolicyListRequest("DELETE", policylist.ID, "")
***REMOVED***

// AddEndpoint add an endpoint to a Policy List
func (policylist *PolicyList) AddEndpoint(endpoint *HNSEndpoint) (*PolicyList, error) ***REMOVED***
	operation := "AddEndpoint"
	title := "HCSShim::PolicyList::" + operation
	logrus.Debugf(title+" id=%s, endpointId:%s", policylist.ID, endpoint.Id)

	_, err := policylist.Delete()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Add Endpoint to the Existing List
	policylist.EndpointReferences = append(policylist.EndpointReferences, "/endpoints/"+endpoint.Id)

	return policylist.Create()
***REMOVED***

// RemoveEndpoint removes an endpoint from the Policy List
func (policylist *PolicyList) RemoveEndpoint(endpoint *HNSEndpoint) (*PolicyList, error) ***REMOVED***
	operation := "RemoveEndpoint"
	title := "HCSShim::PolicyList::" + operation
	logrus.Debugf(title+" id=%s, endpointId:%s", policylist.ID, endpoint.Id)

	_, err := policylist.Delete()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	elementToRemove := "/endpoints/" + endpoint.Id

	var references []string

	for _, endpointReference := range policylist.EndpointReferences ***REMOVED***
		if endpointReference == elementToRemove ***REMOVED***
			continue
		***REMOVED***
		references = append(references, endpointReference)
	***REMOVED***
	policylist.EndpointReferences = references
	return policylist.Create()
***REMOVED***

// AddLoadBalancer policy list for the specified endpoints
func AddLoadBalancer(endpoints []HNSEndpoint, isILB bool, sourceVIP, vip string, protocol uint16, internalPort uint16, externalPort uint16) (*PolicyList, error) ***REMOVED***
	operation := "AddLoadBalancer"
	title := "HCSShim::PolicyList::" + operation
	logrus.Debugf(title+" endpointId=%v, isILB=%v, sourceVIP=%s, vip=%s, protocol=%v, internalPort=%v, externalPort=%v", endpoints, isILB, sourceVIP, vip, protocol, internalPort, externalPort)

	policylist := &PolicyList***REMOVED******REMOVED***

	elbPolicy := &ELBPolicy***REMOVED***
		SourceVIP: sourceVIP,
		ILB:       isILB,
	***REMOVED***

	if len(vip) > 0 ***REMOVED***
		elbPolicy.VIPs = []string***REMOVED***vip***REMOVED***
	***REMOVED***
	elbPolicy.Type = ExternalLoadBalancer
	elbPolicy.Protocol = protocol
	elbPolicy.InternalPort = internalPort
	elbPolicy.ExternalPort = externalPort

	for _, endpoint := range endpoints ***REMOVED***
		policylist.EndpointReferences = append(policylist.EndpointReferences, "/endpoints/"+endpoint.Id)
	***REMOVED***

	jsonString, err := json.Marshal(elbPolicy)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	policylist.Policies = append(policylist.Policies, jsonString)
	return policylist.Create()
***REMOVED***

// AddRoute adds route policy list for the specified endpoints
func AddRoute(endpoints []HNSEndpoint, destinationPrefix string, nextHop string, encapEnabled bool) (*PolicyList, error) ***REMOVED***
	operation := "AddRoute"
	title := "HCSShim::PolicyList::" + operation
	logrus.Debugf(title+" destinationPrefix:%s", destinationPrefix)

	policylist := &PolicyList***REMOVED******REMOVED***

	rPolicy := &RoutePolicy***REMOVED***
		DestinationPrefix: destinationPrefix,
		NextHop:           nextHop,
		EncapEnabled:      encapEnabled,
	***REMOVED***
	rPolicy.Type = Route

	for _, endpoint := range endpoints ***REMOVED***
		policylist.EndpointReferences = append(policylist.EndpointReferences, "/endpoints/"+endpoint.Id)
	***REMOVED***

	jsonString, err := json.Marshal(rPolicy)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	policylist.Policies = append(policylist.Policies, jsonString)
	return policylist.Create()
***REMOVED***
