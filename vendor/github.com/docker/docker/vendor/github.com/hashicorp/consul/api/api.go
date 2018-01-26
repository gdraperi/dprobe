package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// QueryOptions are used to parameterize a query
type QueryOptions struct ***REMOVED***
	// Providing a datacenter overwrites the DC provided
	// by the Config
	Datacenter string

	// AllowStale allows any Consul server (non-leader) to service
	// a read. This allows for lower latency and higher throughput
	AllowStale bool

	// RequireConsistent forces the read to be fully consistent.
	// This is more expensive but prevents ever performing a stale
	// read.
	RequireConsistent bool

	// WaitIndex is used to enable a blocking query. Waits
	// until the timeout or the next index is reached
	WaitIndex uint64

	// WaitTime is used to bound the duration of a wait.
	// Defaults to that of the Config, but can be overriden.
	WaitTime time.Duration

	// Token is used to provide a per-request ACL token
	// which overrides the agent's default token.
	Token string
***REMOVED***

// WriteOptions are used to parameterize a write
type WriteOptions struct ***REMOVED***
	// Providing a datacenter overwrites the DC provided
	// by the Config
	Datacenter string

	// Token is used to provide a per-request ACL token
	// which overrides the agent's default token.
	Token string
***REMOVED***

// QueryMeta is used to return meta data about a query
type QueryMeta struct ***REMOVED***
	// LastIndex. This can be used as a WaitIndex to perform
	// a blocking query
	LastIndex uint64

	// Time of last contact from the leader for the
	// server servicing the request
	LastContact time.Duration

	// Is there a known leader
	KnownLeader bool

	// How long did the request take
	RequestTime time.Duration
***REMOVED***

// WriteMeta is used to return meta data about a write
type WriteMeta struct ***REMOVED***
	// How long did the request take
	RequestTime time.Duration
***REMOVED***

// HttpBasicAuth is used to authenticate http client with HTTP Basic Authentication
type HttpBasicAuth struct ***REMOVED***
	// Username to use for HTTP Basic Authentication
	Username string

	// Password to use for HTTP Basic Authentication
	Password string
***REMOVED***

// Config is used to configure the creation of a client
type Config struct ***REMOVED***
	// Address is the address of the Consul server
	Address string

	// Scheme is the URI scheme for the Consul server
	Scheme string

	// Datacenter to use. If not provided, the default agent datacenter is used.
	Datacenter string

	// HttpClient is the client to use. Default will be
	// used if not provided.
	HttpClient *http.Client

	// HttpAuth is the auth info to use for http access.
	HttpAuth *HttpBasicAuth

	// WaitTime limits how long a Watch will block. If not provided,
	// the agent default values will be used.
	WaitTime time.Duration

	// Token is used to provide a per-request ACL token
	// which overrides the agent's default token.
	Token string
***REMOVED***

// DefaultConfig returns a default configuration for the client
func DefaultConfig() *Config ***REMOVED***
	config := &Config***REMOVED***
		Address:    "127.0.0.1:8500",
		Scheme:     "http",
		HttpClient: http.DefaultClient,
	***REMOVED***

	if addr := os.Getenv("CONSUL_HTTP_ADDR"); addr != "" ***REMOVED***
		config.Address = addr
	***REMOVED***

	if token := os.Getenv("CONSUL_HTTP_TOKEN"); token != "" ***REMOVED***
		config.Token = token
	***REMOVED***

	if auth := os.Getenv("CONSUL_HTTP_AUTH"); auth != "" ***REMOVED***
		var username, password string
		if strings.Contains(auth, ":") ***REMOVED***
			split := strings.SplitN(auth, ":", 2)
			username = split[0]
			password = split[1]
		***REMOVED*** else ***REMOVED***
			username = auth
		***REMOVED***

		config.HttpAuth = &HttpBasicAuth***REMOVED***
			Username: username,
			Password: password,
		***REMOVED***
	***REMOVED***

	if ssl := os.Getenv("CONSUL_HTTP_SSL"); ssl != "" ***REMOVED***
		enabled, err := strconv.ParseBool(ssl)
		if err != nil ***REMOVED***
			log.Printf("[WARN] client: could not parse CONSUL_HTTP_SSL: %s", err)
		***REMOVED***

		if enabled ***REMOVED***
			config.Scheme = "https"
		***REMOVED***
	***REMOVED***

	if verify := os.Getenv("CONSUL_HTTP_SSL_VERIFY"); verify != "" ***REMOVED***
		doVerify, err := strconv.ParseBool(verify)
		if err != nil ***REMOVED***
			log.Printf("[WARN] client: could not parse CONSUL_HTTP_SSL_VERIFY: %s", err)
		***REMOVED***

		if !doVerify ***REMOVED***
			config.HttpClient.Transport = &http.Transport***REMOVED***
				TLSClientConfig: &tls.Config***REMOVED***
					InsecureSkipVerify: true,
				***REMOVED***,
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return config
***REMOVED***

// Client provides a client to the Consul API
type Client struct ***REMOVED***
	config Config
***REMOVED***

// NewClient returns a new client
func NewClient(config *Config) (*Client, error) ***REMOVED***
	// bootstrap the config
	defConfig := DefaultConfig()

	if len(config.Address) == 0 ***REMOVED***
		config.Address = defConfig.Address
	***REMOVED***

	if len(config.Scheme) == 0 ***REMOVED***
		config.Scheme = defConfig.Scheme
	***REMOVED***

	if config.HttpClient == nil ***REMOVED***
		config.HttpClient = defConfig.HttpClient
	***REMOVED***

	if parts := strings.SplitN(config.Address, "unix://", 2); len(parts) == 2 ***REMOVED***
		config.HttpClient = &http.Client***REMOVED***
			Transport: &http.Transport***REMOVED***
				Dial: func(_, _ string) (net.Conn, error) ***REMOVED***
					return net.Dial("unix", parts[1])
				***REMOVED***,
			***REMOVED***,
		***REMOVED***
		config.Address = parts[1]
	***REMOVED***

	client := &Client***REMOVED***
		config: *config,
	***REMOVED***
	return client, nil
***REMOVED***

// request is used to help build up a request
type request struct ***REMOVED***
	config *Config
	method string
	url    *url.URL
	params url.Values
	body   io.Reader
	obj    interface***REMOVED******REMOVED***
***REMOVED***

// setQueryOptions is used to annotate the request with
// additional query options
func (r *request) setQueryOptions(q *QueryOptions) ***REMOVED***
	if q == nil ***REMOVED***
		return
	***REMOVED***
	if q.Datacenter != "" ***REMOVED***
		r.params.Set("dc", q.Datacenter)
	***REMOVED***
	if q.AllowStale ***REMOVED***
		r.params.Set("stale", "")
	***REMOVED***
	if q.RequireConsistent ***REMOVED***
		r.params.Set("consistent", "")
	***REMOVED***
	if q.WaitIndex != 0 ***REMOVED***
		r.params.Set("index", strconv.FormatUint(q.WaitIndex, 10))
	***REMOVED***
	if q.WaitTime != 0 ***REMOVED***
		r.params.Set("wait", durToMsec(q.WaitTime))
	***REMOVED***
	if q.Token != "" ***REMOVED***
		r.params.Set("token", q.Token)
	***REMOVED***
***REMOVED***

// durToMsec converts a duration to a millisecond specified string
func durToMsec(dur time.Duration) string ***REMOVED***
	return fmt.Sprintf("%dms", dur/time.Millisecond)
***REMOVED***

// setWriteOptions is used to annotate the request with
// additional write options
func (r *request) setWriteOptions(q *WriteOptions) ***REMOVED***
	if q == nil ***REMOVED***
		return
	***REMOVED***
	if q.Datacenter != "" ***REMOVED***
		r.params.Set("dc", q.Datacenter)
	***REMOVED***
	if q.Token != "" ***REMOVED***
		r.params.Set("token", q.Token)
	***REMOVED***
***REMOVED***

// toHTTP converts the request to an HTTP request
func (r *request) toHTTP() (*http.Request, error) ***REMOVED***
	// Encode the query parameters
	r.url.RawQuery = r.params.Encode()

	// Check if we should encode the body
	if r.body == nil && r.obj != nil ***REMOVED***
		if b, err := encodeBody(r.obj); err != nil ***REMOVED***
			return nil, err
		***REMOVED*** else ***REMOVED***
			r.body = b
		***REMOVED***
	***REMOVED***

	// Create the HTTP request
	req, err := http.NewRequest(r.method, r.url.RequestURI(), r.body)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	req.URL.Host = r.url.Host
	req.URL.Scheme = r.url.Scheme
	req.Host = r.url.Host

	// Setup auth
	if r.config.HttpAuth != nil ***REMOVED***
		req.SetBasicAuth(r.config.HttpAuth.Username, r.config.HttpAuth.Password)
	***REMOVED***

	return req, nil
***REMOVED***

// newRequest is used to create a new request
func (c *Client) newRequest(method, path string) *request ***REMOVED***
	r := &request***REMOVED***
		config: &c.config,
		method: method,
		url: &url.URL***REMOVED***
			Scheme: c.config.Scheme,
			Host:   c.config.Address,
			Path:   path,
		***REMOVED***,
		params: make(map[string][]string),
	***REMOVED***
	if c.config.Datacenter != "" ***REMOVED***
		r.params.Set("dc", c.config.Datacenter)
	***REMOVED***
	if c.config.WaitTime != 0 ***REMOVED***
		r.params.Set("wait", durToMsec(r.config.WaitTime))
	***REMOVED***
	if c.config.Token != "" ***REMOVED***
		r.params.Set("token", r.config.Token)
	***REMOVED***
	return r
***REMOVED***

// doRequest runs a request with our client
func (c *Client) doRequest(r *request) (time.Duration, *http.Response, error) ***REMOVED***
	req, err := r.toHTTP()
	if err != nil ***REMOVED***
		return 0, nil, err
	***REMOVED***
	start := time.Now()
	resp, err := c.config.HttpClient.Do(req)
	diff := time.Now().Sub(start)
	return diff, resp, err
***REMOVED***

// Query is used to do a GET request against an endpoint
// and deserialize the response into an interface using
// standard Consul conventions.
func (c *Client) query(endpoint string, out interface***REMOVED******REMOVED***, q *QueryOptions) (*QueryMeta, error) ***REMOVED***
	r := c.newRequest("GET", endpoint)
	r.setQueryOptions(q)
	rtt, resp, err := requireOK(c.doRequest(r))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer resp.Body.Close()

	qm := &QueryMeta***REMOVED******REMOVED***
	parseQueryMeta(resp, qm)
	qm.RequestTime = rtt

	if err := decodeBody(resp, out); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return qm, nil
***REMOVED***

// write is used to do a PUT request against an endpoint
// and serialize/deserialized using the standard Consul conventions.
func (c *Client) write(endpoint string, in, out interface***REMOVED******REMOVED***, q *WriteOptions) (*WriteMeta, error) ***REMOVED***
	r := c.newRequest("PUT", endpoint)
	r.setWriteOptions(q)
	r.obj = in
	rtt, resp, err := requireOK(c.doRequest(r))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer resp.Body.Close()

	wm := &WriteMeta***REMOVED***RequestTime: rtt***REMOVED***
	if out != nil ***REMOVED***
		if err := decodeBody(resp, &out); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return wm, nil
***REMOVED***

// parseQueryMeta is used to help parse query meta-data
func parseQueryMeta(resp *http.Response, q *QueryMeta) error ***REMOVED***
	header := resp.Header

	// Parse the X-Consul-Index
	index, err := strconv.ParseUint(header.Get("X-Consul-Index"), 10, 64)
	if err != nil ***REMOVED***
		return fmt.Errorf("Failed to parse X-Consul-Index: %v", err)
	***REMOVED***
	q.LastIndex = index

	// Parse the X-Consul-LastContact
	last, err := strconv.ParseUint(header.Get("X-Consul-LastContact"), 10, 64)
	if err != nil ***REMOVED***
		return fmt.Errorf("Failed to parse X-Consul-LastContact: %v", err)
	***REMOVED***
	q.LastContact = time.Duration(last) * time.Millisecond

	// Parse the X-Consul-KnownLeader
	switch header.Get("X-Consul-KnownLeader") ***REMOVED***
	case "true":
		q.KnownLeader = true
	default:
		q.KnownLeader = false
	***REMOVED***
	return nil
***REMOVED***

// decodeBody is used to JSON decode a body
func decodeBody(resp *http.Response, out interface***REMOVED******REMOVED***) error ***REMOVED***
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(out)
***REMOVED***

// encodeBody is used to encode a request body
func encodeBody(obj interface***REMOVED******REMOVED***) (io.Reader, error) ***REMOVED***
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(obj); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return buf, nil
***REMOVED***

// requireOK is used to wrap doRequest and check for a 200
func requireOK(d time.Duration, resp *http.Response, e error) (time.Duration, *http.Response, error) ***REMOVED***
	if e != nil ***REMOVED***
		if resp != nil ***REMOVED***
			resp.Body.Close()
		***REMOVED***
		return d, nil, e
	***REMOVED***
	if resp.StatusCode != 200 ***REMOVED***
		var buf bytes.Buffer
		io.Copy(&buf, resp.Body)
		resp.Body.Close()
		return d, nil, fmt.Errorf("Unexpected response code: %d (%s)", resp.StatusCode, buf.Bytes())
	***REMOVED***
	return d, resp, nil
***REMOVED***
