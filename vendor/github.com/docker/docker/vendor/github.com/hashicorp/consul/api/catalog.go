package api

type Node struct ***REMOVED***
	Node    string
	Address string
***REMOVED***

type CatalogService struct ***REMOVED***
	Node           string
	Address        string
	ServiceID      string
	ServiceName    string
	ServiceAddress string
	ServiceTags    []string
	ServicePort    int
***REMOVED***

type CatalogNode struct ***REMOVED***
	Node     *Node
	Services map[string]*AgentService
***REMOVED***

type CatalogRegistration struct ***REMOVED***
	Node       string
	Address    string
	Datacenter string
	Service    *AgentService
	Check      *AgentCheck
***REMOVED***

type CatalogDeregistration struct ***REMOVED***
	Node       string
	Address    string
	Datacenter string
	ServiceID  string
	CheckID    string
***REMOVED***

// Catalog can be used to query the Catalog endpoints
type Catalog struct ***REMOVED***
	c *Client
***REMOVED***

// Catalog returns a handle to the catalog endpoints
func (c *Client) Catalog() *Catalog ***REMOVED***
	return &Catalog***REMOVED***c***REMOVED***
***REMOVED***

func (c *Catalog) Register(reg *CatalogRegistration, q *WriteOptions) (*WriteMeta, error) ***REMOVED***
	r := c.c.newRequest("PUT", "/v1/catalog/register")
	r.setWriteOptions(q)
	r.obj = reg
	rtt, resp, err := requireOK(c.c.doRequest(r))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	resp.Body.Close()

	wm := &WriteMeta***REMOVED******REMOVED***
	wm.RequestTime = rtt

	return wm, nil
***REMOVED***

func (c *Catalog) Deregister(dereg *CatalogDeregistration, q *WriteOptions) (*WriteMeta, error) ***REMOVED***
	r := c.c.newRequest("PUT", "/v1/catalog/deregister")
	r.setWriteOptions(q)
	r.obj = dereg
	rtt, resp, err := requireOK(c.c.doRequest(r))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	resp.Body.Close()

	wm := &WriteMeta***REMOVED******REMOVED***
	wm.RequestTime = rtt

	return wm, nil
***REMOVED***

// Datacenters is used to query for all the known datacenters
func (c *Catalog) Datacenters() ([]string, error) ***REMOVED***
	r := c.c.newRequest("GET", "/v1/catalog/datacenters")
	_, resp, err := requireOK(c.c.doRequest(r))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer resp.Body.Close()

	var out []string
	if err := decodeBody(resp, &out); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return out, nil
***REMOVED***

// Nodes is used to query all the known nodes
func (c *Catalog) Nodes(q *QueryOptions) ([]*Node, *QueryMeta, error) ***REMOVED***
	r := c.c.newRequest("GET", "/v1/catalog/nodes")
	r.setQueryOptions(q)
	rtt, resp, err := requireOK(c.c.doRequest(r))
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	defer resp.Body.Close()

	qm := &QueryMeta***REMOVED******REMOVED***
	parseQueryMeta(resp, qm)
	qm.RequestTime = rtt

	var out []*Node
	if err := decodeBody(resp, &out); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return out, qm, nil
***REMOVED***

// Services is used to query for all known services
func (c *Catalog) Services(q *QueryOptions) (map[string][]string, *QueryMeta, error) ***REMOVED***
	r := c.c.newRequest("GET", "/v1/catalog/services")
	r.setQueryOptions(q)
	rtt, resp, err := requireOK(c.c.doRequest(r))
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	defer resp.Body.Close()

	qm := &QueryMeta***REMOVED******REMOVED***
	parseQueryMeta(resp, qm)
	qm.RequestTime = rtt

	var out map[string][]string
	if err := decodeBody(resp, &out); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return out, qm, nil
***REMOVED***

// Service is used to query catalog entries for a given service
func (c *Catalog) Service(service, tag string, q *QueryOptions) ([]*CatalogService, *QueryMeta, error) ***REMOVED***
	r := c.c.newRequest("GET", "/v1/catalog/service/"+service)
	r.setQueryOptions(q)
	if tag != "" ***REMOVED***
		r.params.Set("tag", tag)
	***REMOVED***
	rtt, resp, err := requireOK(c.c.doRequest(r))
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	defer resp.Body.Close()

	qm := &QueryMeta***REMOVED******REMOVED***
	parseQueryMeta(resp, qm)
	qm.RequestTime = rtt

	var out []*CatalogService
	if err := decodeBody(resp, &out); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return out, qm, nil
***REMOVED***

// Node is used to query for service information about a single node
func (c *Catalog) Node(node string, q *QueryOptions) (*CatalogNode, *QueryMeta, error) ***REMOVED***
	r := c.c.newRequest("GET", "/v1/catalog/node/"+node)
	r.setQueryOptions(q)
	rtt, resp, err := requireOK(c.c.doRequest(r))
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	defer resp.Body.Close()

	qm := &QueryMeta***REMOVED******REMOVED***
	parseQueryMeta(resp, qm)
	qm.RequestTime = rtt

	var out *CatalogNode
	if err := decodeBody(resp, &out); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return out, qm, nil
***REMOVED***
