package api

// Raw can be used to do raw queries against custom endpoints
type Raw struct ***REMOVED***
	c *Client
***REMOVED***

// Raw returns a handle to query endpoints
func (c *Client) Raw() *Raw ***REMOVED***
	return &Raw***REMOVED***c***REMOVED***
***REMOVED***

// Query is used to do a GET request against an endpoint
// and deserialize the response into an interface using
// standard Consul conventions.
func (raw *Raw) Query(endpoint string, out interface***REMOVED******REMOVED***, q *QueryOptions) (*QueryMeta, error) ***REMOVED***
	return raw.c.query(endpoint, out, q)
***REMOVED***

// Write is used to do a PUT request against an endpoint
// and serialize/deserialized using the standard Consul conventions.
func (raw *Raw) Write(endpoint string, in, out interface***REMOVED******REMOVED***, q *WriteOptions) (*WriteMeta, error) ***REMOVED***
	return raw.c.write(endpoint, in, out, q)
***REMOVED***
