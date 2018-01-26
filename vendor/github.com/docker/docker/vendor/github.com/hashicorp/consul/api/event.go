package api

import (
	"bytes"
	"strconv"
)

// Event can be used to query the Event endpoints
type Event struct ***REMOVED***
	c *Client
***REMOVED***

// UserEvent represents an event that was fired by the user
type UserEvent struct ***REMOVED***
	ID            string
	Name          string
	Payload       []byte
	NodeFilter    string
	ServiceFilter string
	TagFilter     string
	Version       int
	LTime         uint64
***REMOVED***

// Event returns a handle to the event endpoints
func (c *Client) Event() *Event ***REMOVED***
	return &Event***REMOVED***c***REMOVED***
***REMOVED***

// Fire is used to fire a new user event. Only the Name, Payload and Filters
// are respected. This returns the ID or an associated error. Cross DC requests
// are supported.
func (e *Event) Fire(params *UserEvent, q *WriteOptions) (string, *WriteMeta, error) ***REMOVED***
	r := e.c.newRequest("PUT", "/v1/event/fire/"+params.Name)
	r.setWriteOptions(q)
	if params.NodeFilter != "" ***REMOVED***
		r.params.Set("node", params.NodeFilter)
	***REMOVED***
	if params.ServiceFilter != "" ***REMOVED***
		r.params.Set("service", params.ServiceFilter)
	***REMOVED***
	if params.TagFilter != "" ***REMOVED***
		r.params.Set("tag", params.TagFilter)
	***REMOVED***
	if params.Payload != nil ***REMOVED***
		r.body = bytes.NewReader(params.Payload)
	***REMOVED***

	rtt, resp, err := requireOK(e.c.doRequest(r))
	if err != nil ***REMOVED***
		return "", nil, err
	***REMOVED***
	defer resp.Body.Close()

	wm := &WriteMeta***REMOVED***RequestTime: rtt***REMOVED***
	var out UserEvent
	if err := decodeBody(resp, &out); err != nil ***REMOVED***
		return "", nil, err
	***REMOVED***
	return out.ID, wm, nil
***REMOVED***

// List is used to get the most recent events an agent has received.
// This list can be optionally filtered by the name. This endpoint supports
// quasi-blocking queries. The index is not monotonic, nor does it provide provide
// LastContact or KnownLeader.
func (e *Event) List(name string, q *QueryOptions) ([]*UserEvent, *QueryMeta, error) ***REMOVED***
	r := e.c.newRequest("GET", "/v1/event/list")
	r.setQueryOptions(q)
	if name != "" ***REMOVED***
		r.params.Set("name", name)
	***REMOVED***
	rtt, resp, err := requireOK(e.c.doRequest(r))
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	defer resp.Body.Close()

	qm := &QueryMeta***REMOVED******REMOVED***
	parseQueryMeta(resp, qm)
	qm.RequestTime = rtt

	var entries []*UserEvent
	if err := decodeBody(resp, &entries); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return entries, qm, nil
***REMOVED***

// IDToIndex is a bit of a hack. This simulates the index generation to
// convert an event ID into a WaitIndex.
func (e *Event) IDToIndex(uuid string) uint64 ***REMOVED***
	lower := uuid[0:8] + uuid[9:13] + uuid[14:18]
	upper := uuid[19:23] + uuid[24:36]
	lowVal, err := strconv.ParseUint(lower, 16, 64)
	if err != nil ***REMOVED***
		panic("Failed to convert " + lower)
	***REMOVED***
	highVal, err := strconv.ParseUint(upper, 16, 64)
	if err != nil ***REMOVED***
		panic("Failed to convert " + upper)
	***REMOVED***
	return lowVal ^ highVal
***REMOVED***
