package api

import (
	"fmt"
)

// AgentCheck represents a check known to the agent
type AgentCheck struct ***REMOVED***
	Node        string
	CheckID     string
	Name        string
	Status      string
	Notes       string
	Output      string
	ServiceID   string
	ServiceName string
***REMOVED***

// AgentService represents a service known to the agent
type AgentService struct ***REMOVED***
	ID      string
	Service string
	Tags    []string
	Port    int
	Address string
***REMOVED***

// AgentMember represents a cluster member known to the agent
type AgentMember struct ***REMOVED***
	Name        string
	Addr        string
	Port        uint16
	Tags        map[string]string
	Status      int
	ProtocolMin uint8
	ProtocolMax uint8
	ProtocolCur uint8
	DelegateMin uint8
	DelegateMax uint8
	DelegateCur uint8
***REMOVED***

// AgentServiceRegistration is used to register a new service
type AgentServiceRegistration struct ***REMOVED***
	ID      string   `json:",omitempty"`
	Name    string   `json:",omitempty"`
	Tags    []string `json:",omitempty"`
	Port    int      `json:",omitempty"`
	Address string   `json:",omitempty"`
	Check   *AgentServiceCheck
	Checks  AgentServiceChecks
***REMOVED***

// AgentCheckRegistration is used to register a new check
type AgentCheckRegistration struct ***REMOVED***
	ID        string `json:",omitempty"`
	Name      string `json:",omitempty"`
	Notes     string `json:",omitempty"`
	ServiceID string `json:",omitempty"`
	AgentServiceCheck
***REMOVED***

// AgentServiceCheck is used to create an associated
// check for a service
type AgentServiceCheck struct ***REMOVED***
	Script   string `json:",omitempty"`
	Interval string `json:",omitempty"`
	Timeout  string `json:",omitempty"`
	TTL      string `json:",omitempty"`
	HTTP     string `json:",omitempty"`
	Status   string `json:",omitempty"`
***REMOVED***
type AgentServiceChecks []*AgentServiceCheck

// Agent can be used to query the Agent endpoints
type Agent struct ***REMOVED***
	c *Client

	// cache the node name
	nodeName string
***REMOVED***

// Agent returns a handle to the agent endpoints
func (c *Client) Agent() *Agent ***REMOVED***
	return &Agent***REMOVED***c: c***REMOVED***
***REMOVED***

// Self is used to query the agent we are speaking to for
// information about itself
func (a *Agent) Self() (map[string]map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	r := a.c.newRequest("GET", "/v1/agent/self")
	_, resp, err := requireOK(a.c.doRequest(r))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer resp.Body.Close()

	var out map[string]map[string]interface***REMOVED******REMOVED***
	if err := decodeBody(resp, &out); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return out, nil
***REMOVED***

// NodeName is used to get the node name of the agent
func (a *Agent) NodeName() (string, error) ***REMOVED***
	if a.nodeName != "" ***REMOVED***
		return a.nodeName, nil
	***REMOVED***
	info, err := a.Self()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	name := info["Config"]["NodeName"].(string)
	a.nodeName = name
	return name, nil
***REMOVED***

// Checks returns the locally registered checks
func (a *Agent) Checks() (map[string]*AgentCheck, error) ***REMOVED***
	r := a.c.newRequest("GET", "/v1/agent/checks")
	_, resp, err := requireOK(a.c.doRequest(r))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer resp.Body.Close()

	var out map[string]*AgentCheck
	if err := decodeBody(resp, &out); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return out, nil
***REMOVED***

// Services returns the locally registered services
func (a *Agent) Services() (map[string]*AgentService, error) ***REMOVED***
	r := a.c.newRequest("GET", "/v1/agent/services")
	_, resp, err := requireOK(a.c.doRequest(r))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer resp.Body.Close()

	var out map[string]*AgentService
	if err := decodeBody(resp, &out); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return out, nil
***REMOVED***

// Members returns the known gossip members. The WAN
// flag can be used to query a server for WAN members.
func (a *Agent) Members(wan bool) ([]*AgentMember, error) ***REMOVED***
	r := a.c.newRequest("GET", "/v1/agent/members")
	if wan ***REMOVED***
		r.params.Set("wan", "1")
	***REMOVED***
	_, resp, err := requireOK(a.c.doRequest(r))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer resp.Body.Close()

	var out []*AgentMember
	if err := decodeBody(resp, &out); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return out, nil
***REMOVED***

// ServiceRegister is used to register a new service with
// the local agent
func (a *Agent) ServiceRegister(service *AgentServiceRegistration) error ***REMOVED***
	r := a.c.newRequest("PUT", "/v1/agent/service/register")
	r.obj = service
	_, resp, err := requireOK(a.c.doRequest(r))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	resp.Body.Close()
	return nil
***REMOVED***

// ServiceDeregister is used to deregister a service with
// the local agent
func (a *Agent) ServiceDeregister(serviceID string) error ***REMOVED***
	r := a.c.newRequest("PUT", "/v1/agent/service/deregister/"+serviceID)
	_, resp, err := requireOK(a.c.doRequest(r))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	resp.Body.Close()
	return nil
***REMOVED***

// PassTTL is used to set a TTL check to the passing state
func (a *Agent) PassTTL(checkID, note string) error ***REMOVED***
	return a.UpdateTTL(checkID, note, "pass")
***REMOVED***

// WarnTTL is used to set a TTL check to the warning state
func (a *Agent) WarnTTL(checkID, note string) error ***REMOVED***
	return a.UpdateTTL(checkID, note, "warn")
***REMOVED***

// FailTTL is used to set a TTL check to the failing state
func (a *Agent) FailTTL(checkID, note string) error ***REMOVED***
	return a.UpdateTTL(checkID, note, "fail")
***REMOVED***

// UpdateTTL is used to update the TTL of a check
func (a *Agent) UpdateTTL(checkID, note, status string) error ***REMOVED***
	switch status ***REMOVED***
	case "pass":
	case "warn":
	case "fail":
	default:
		return fmt.Errorf("Invalid status: %s", status)
	***REMOVED***
	endpoint := fmt.Sprintf("/v1/agent/check/%s/%s", status, checkID)
	r := a.c.newRequest("PUT", endpoint)
	r.params.Set("note", note)
	_, resp, err := requireOK(a.c.doRequest(r))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	resp.Body.Close()
	return nil
***REMOVED***

// CheckRegister is used to register a new check with
// the local agent
func (a *Agent) CheckRegister(check *AgentCheckRegistration) error ***REMOVED***
	r := a.c.newRequest("PUT", "/v1/agent/check/register")
	r.obj = check
	_, resp, err := requireOK(a.c.doRequest(r))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	resp.Body.Close()
	return nil
***REMOVED***

// CheckDeregister is used to deregister a check with
// the local agent
func (a *Agent) CheckDeregister(checkID string) error ***REMOVED***
	r := a.c.newRequest("PUT", "/v1/agent/check/deregister/"+checkID)
	_, resp, err := requireOK(a.c.doRequest(r))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	resp.Body.Close()
	return nil
***REMOVED***

// Join is used to instruct the agent to attempt a join to
// another cluster member
func (a *Agent) Join(addr string, wan bool) error ***REMOVED***
	r := a.c.newRequest("PUT", "/v1/agent/join/"+addr)
	if wan ***REMOVED***
		r.params.Set("wan", "1")
	***REMOVED***
	_, resp, err := requireOK(a.c.doRequest(r))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	resp.Body.Close()
	return nil
***REMOVED***

// ForceLeave is used to have the agent eject a failed node
func (a *Agent) ForceLeave(node string) error ***REMOVED***
	r := a.c.newRequest("PUT", "/v1/agent/force-leave/"+node)
	_, resp, err := requireOK(a.c.doRequest(r))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	resp.Body.Close()
	return nil
***REMOVED***

// EnableServiceMaintenance toggles service maintenance mode on
// for the given service ID.
func (a *Agent) EnableServiceMaintenance(serviceID, reason string) error ***REMOVED***
	r := a.c.newRequest("PUT", "/v1/agent/service/maintenance/"+serviceID)
	r.params.Set("enable", "true")
	r.params.Set("reason", reason)
	_, resp, err := requireOK(a.c.doRequest(r))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	resp.Body.Close()
	return nil
***REMOVED***

// DisableServiceMaintenance toggles service maintenance mode off
// for the given service ID.
func (a *Agent) DisableServiceMaintenance(serviceID string) error ***REMOVED***
	r := a.c.newRequest("PUT", "/v1/agent/service/maintenance/"+serviceID)
	r.params.Set("enable", "false")
	_, resp, err := requireOK(a.c.doRequest(r))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	resp.Body.Close()
	return nil
***REMOVED***

// EnableNodeMaintenance toggles node maintenance mode on for the
// agent we are connected to.
func (a *Agent) EnableNodeMaintenance(reason string) error ***REMOVED***
	r := a.c.newRequest("PUT", "/v1/agent/maintenance")
	r.params.Set("enable", "true")
	r.params.Set("reason", reason)
	_, resp, err := requireOK(a.c.doRequest(r))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	resp.Body.Close()
	return nil
***REMOVED***

// DisableNodeMaintenance toggles node maintenance mode off for the
// agent we are connected to.
func (a *Agent) DisableNodeMaintenance() error ***REMOVED***
	r := a.c.newRequest("PUT", "/v1/agent/maintenance")
	r.params.Set("enable", "false")
	_, resp, err := requireOK(a.c.doRequest(r))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	resp.Body.Close()
	return nil
***REMOVED***
