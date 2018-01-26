package api

import (
	"fmt"
)

// HealthCheck is used to represent a single check
type HealthCheck struct ***REMOVED***
	Node        string
	CheckID     string
	Name        string
	Status      string
	Notes       string
	Output      string
	ServiceID   string
	ServiceName string
***REMOVED***

// ServiceEntry is used for the health service endpoint
type ServiceEntry struct ***REMOVED***
	Node    *Node
	Service *AgentService
	Checks  []*HealthCheck
***REMOVED***

// Health can be used to query the Health endpoints
type Health struct ***REMOVED***
	c *Client
***REMOVED***

// Health returns a handle to the health endpoints
func (c *Client) Health() *Health ***REMOVED***
	return &Health***REMOVED***c***REMOVED***
***REMOVED***

// Node is used to query for checks belonging to a given node
func (h *Health) Node(node string, q *QueryOptions) ([]*HealthCheck, *QueryMeta, error) ***REMOVED***
	r := h.c.newRequest("GET", "/v1/health/node/"+node)
	r.setQueryOptions(q)
	rtt, resp, err := requireOK(h.c.doRequest(r))
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	defer resp.Body.Close()

	qm := &QueryMeta***REMOVED******REMOVED***
	parseQueryMeta(resp, qm)
	qm.RequestTime = rtt

	var out []*HealthCheck
	if err := decodeBody(resp, &out); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return out, qm, nil
***REMOVED***

// Checks is used to return the checks associated with a service
func (h *Health) Checks(service string, q *QueryOptions) ([]*HealthCheck, *QueryMeta, error) ***REMOVED***
	r := h.c.newRequest("GET", "/v1/health/checks/"+service)
	r.setQueryOptions(q)
	rtt, resp, err := requireOK(h.c.doRequest(r))
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	defer resp.Body.Close()

	qm := &QueryMeta***REMOVED******REMOVED***
	parseQueryMeta(resp, qm)
	qm.RequestTime = rtt

	var out []*HealthCheck
	if err := decodeBody(resp, &out); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return out, qm, nil
***REMOVED***

// Service is used to query health information along with service info
// for a given service. It can optionally do server-side filtering on a tag
// or nodes with passing health checks only.
func (h *Health) Service(service, tag string, passingOnly bool, q *QueryOptions) ([]*ServiceEntry, *QueryMeta, error) ***REMOVED***
	r := h.c.newRequest("GET", "/v1/health/service/"+service)
	r.setQueryOptions(q)
	if tag != "" ***REMOVED***
		r.params.Set("tag", tag)
	***REMOVED***
	if passingOnly ***REMOVED***
		r.params.Set("passing", "1")
	***REMOVED***
	rtt, resp, err := requireOK(h.c.doRequest(r))
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	defer resp.Body.Close()

	qm := &QueryMeta***REMOVED******REMOVED***
	parseQueryMeta(resp, qm)
	qm.RequestTime = rtt

	var out []*ServiceEntry
	if err := decodeBody(resp, &out); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return out, qm, nil
***REMOVED***

// State is used to retreive all the checks in a given state.
// The wildcard "any" state can also be used for all checks.
func (h *Health) State(state string, q *QueryOptions) ([]*HealthCheck, *QueryMeta, error) ***REMOVED***
	switch state ***REMOVED***
	case "any":
	case "warning":
	case "critical":
	case "passing":
	case "unknown":
	default:
		return nil, nil, fmt.Errorf("Unsupported state: %v", state)
	***REMOVED***
	r := h.c.newRequest("GET", "/v1/health/state/"+state)
	r.setQueryOptions(q)
	rtt, resp, err := requireOK(h.c.doRequest(r))
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	defer resp.Body.Close()

	qm := &QueryMeta***REMOVED******REMOVED***
	parseQueryMeta(resp, qm)
	qm.RequestTime = rtt

	var out []*HealthCheck
	if err := decodeBody(resp, &out); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return out, qm, nil
***REMOVED***
