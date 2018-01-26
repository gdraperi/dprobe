package api

const (
	// ACLCLientType is the client type token
	ACLClientType = "client"

	// ACLManagementType is the management type token
	ACLManagementType = "management"
)

// ACLEntry is used to represent an ACL entry
type ACLEntry struct ***REMOVED***
	CreateIndex uint64
	ModifyIndex uint64
	ID          string
	Name        string
	Type        string
	Rules       string
***REMOVED***

// ACL can be used to query the ACL endpoints
type ACL struct ***REMOVED***
	c *Client
***REMOVED***

// ACL returns a handle to the ACL endpoints
func (c *Client) ACL() *ACL ***REMOVED***
	return &ACL***REMOVED***c***REMOVED***
***REMOVED***

// Create is used to generate a new token with the given parameters
func (a *ACL) Create(acl *ACLEntry, q *WriteOptions) (string, *WriteMeta, error) ***REMOVED***
	r := a.c.newRequest("PUT", "/v1/acl/create")
	r.setWriteOptions(q)
	r.obj = acl
	rtt, resp, err := requireOK(a.c.doRequest(r))
	if err != nil ***REMOVED***
		return "", nil, err
	***REMOVED***
	defer resp.Body.Close()

	wm := &WriteMeta***REMOVED***RequestTime: rtt***REMOVED***
	var out struct***REMOVED*** ID string ***REMOVED***
	if err := decodeBody(resp, &out); err != nil ***REMOVED***
		return "", nil, err
	***REMOVED***
	return out.ID, wm, nil
***REMOVED***

// Update is used to update the rules of an existing token
func (a *ACL) Update(acl *ACLEntry, q *WriteOptions) (*WriteMeta, error) ***REMOVED***
	r := a.c.newRequest("PUT", "/v1/acl/update")
	r.setWriteOptions(q)
	r.obj = acl
	rtt, resp, err := requireOK(a.c.doRequest(r))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer resp.Body.Close()

	wm := &WriteMeta***REMOVED***RequestTime: rtt***REMOVED***
	return wm, nil
***REMOVED***

// Destroy is used to destroy a given ACL token ID
func (a *ACL) Destroy(id string, q *WriteOptions) (*WriteMeta, error) ***REMOVED***
	r := a.c.newRequest("PUT", "/v1/acl/destroy/"+id)
	r.setWriteOptions(q)
	rtt, resp, err := requireOK(a.c.doRequest(r))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	resp.Body.Close()

	wm := &WriteMeta***REMOVED***RequestTime: rtt***REMOVED***
	return wm, nil
***REMOVED***

// Clone is used to return a new token cloned from an existing one
func (a *ACL) Clone(id string, q *WriteOptions) (string, *WriteMeta, error) ***REMOVED***
	r := a.c.newRequest("PUT", "/v1/acl/clone/"+id)
	r.setWriteOptions(q)
	rtt, resp, err := requireOK(a.c.doRequest(r))
	if err != nil ***REMOVED***
		return "", nil, err
	***REMOVED***
	defer resp.Body.Close()

	wm := &WriteMeta***REMOVED***RequestTime: rtt***REMOVED***
	var out struct***REMOVED*** ID string ***REMOVED***
	if err := decodeBody(resp, &out); err != nil ***REMOVED***
		return "", nil, err
	***REMOVED***
	return out.ID, wm, nil
***REMOVED***

// Info is used to query for information about an ACL token
func (a *ACL) Info(id string, q *QueryOptions) (*ACLEntry, *QueryMeta, error) ***REMOVED***
	r := a.c.newRequest("GET", "/v1/acl/info/"+id)
	r.setQueryOptions(q)
	rtt, resp, err := requireOK(a.c.doRequest(r))
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	defer resp.Body.Close()

	qm := &QueryMeta***REMOVED******REMOVED***
	parseQueryMeta(resp, qm)
	qm.RequestTime = rtt

	var entries []*ACLEntry
	if err := decodeBody(resp, &entries); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	if len(entries) > 0 ***REMOVED***
		return entries[0], qm, nil
	***REMOVED***
	return nil, qm, nil
***REMOVED***

// List is used to get all the ACL tokens
func (a *ACL) List(q *QueryOptions) ([]*ACLEntry, *QueryMeta, error) ***REMOVED***
	r := a.c.newRequest("GET", "/v1/acl/list")
	r.setQueryOptions(q)
	rtt, resp, err := requireOK(a.c.doRequest(r))
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	defer resp.Body.Close()

	qm := &QueryMeta***REMOVED******REMOVED***
	parseQueryMeta(resp, qm)
	qm.RequestTime = rtt

	var entries []*ACLEntry
	if err := decodeBody(resp, &entries); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return entries, qm, nil
***REMOVED***
