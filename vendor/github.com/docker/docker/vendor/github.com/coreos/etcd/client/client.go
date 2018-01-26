// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/coreos/etcd/version"

	"golang.org/x/net/context"
)

var (
	ErrNoEndpoints           = errors.New("client: no endpoints available")
	ErrTooManyRedirects      = errors.New("client: too many redirects")
	ErrClusterUnavailable    = errors.New("client: etcd cluster is unavailable or misconfigured")
	ErrNoLeaderEndpoint      = errors.New("client: no leader endpoint available")
	errTooManyRedirectChecks = errors.New("client: too many redirect checks")

	// oneShotCtxValue is set on a context using WithValue(&oneShotValue) so
	// that Do() will not retry a request
	oneShotCtxValue interface***REMOVED******REMOVED***
)

var DefaultRequestTimeout = 5 * time.Second

var DefaultTransport CancelableTransport = &http.Transport***REMOVED***
	Proxy: http.ProxyFromEnvironment,
	Dial: (&net.Dialer***REMOVED***
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	***REMOVED***).Dial,
	TLSHandshakeTimeout: 10 * time.Second,
***REMOVED***

type EndpointSelectionMode int

const (
	// EndpointSelectionRandom is the default value of the 'SelectionMode'.
	// As the name implies, the client object will pick a node from the members
	// of the cluster in a random fashion. If the cluster has three members, A, B,
	// and C, the client picks any node from its three members as its request
	// destination.
	EndpointSelectionRandom EndpointSelectionMode = iota

	// If 'SelectionMode' is set to 'EndpointSelectionPrioritizeLeader',
	// requests are sent directly to the cluster leader. This reduces
	// forwarding roundtrips compared to making requests to etcd followers
	// who then forward them to the cluster leader. In the event of a leader
	// failure, however, clients configured this way cannot prioritize among
	// the remaining etcd followers. Therefore, when a client sets 'SelectionMode'
	// to 'EndpointSelectionPrioritizeLeader', it must use 'client.AutoSync()' to
	// maintain its knowledge of current cluster state.
	//
	// This mode should be used with Client.AutoSync().
	EndpointSelectionPrioritizeLeader
)

type Config struct ***REMOVED***
	// Endpoints defines a set of URLs (schemes, hosts and ports only)
	// that can be used to communicate with a logical etcd cluster. For
	// example, a three-node cluster could be provided like so:
	//
	// 	Endpoints: []string***REMOVED***
	//		"http://node1.example.com:2379",
	//		"http://node2.example.com:2379",
	//		"http://node3.example.com:2379",
	//	***REMOVED***
	//
	// If multiple endpoints are provided, the Client will attempt to
	// use them all in the event that one or more of them are unusable.
	//
	// If Client.Sync is ever called, the Client may cache an alternate
	// set of endpoints to continue operation.
	Endpoints []string

	// Transport is used by the Client to drive HTTP requests. If not
	// provided, DefaultTransport will be used.
	Transport CancelableTransport

	// CheckRedirect specifies the policy for handling HTTP redirects.
	// If CheckRedirect is not nil, the Client calls it before
	// following an HTTP redirect. The sole argument is the number of
	// requests that have already been made. If CheckRedirect returns
	// an error, Client.Do will not make any further requests and return
	// the error back it to the caller.
	//
	// If CheckRedirect is nil, the Client uses its default policy,
	// which is to stop after 10 consecutive requests.
	CheckRedirect CheckRedirectFunc

	// Username specifies the user credential to add as an authorization header
	Username string

	// Password is the password for the specified user to add as an authorization header
	// to the request.
	Password string

	// HeaderTimeoutPerRequest specifies the time limit to wait for response
	// header in a single request made by the Client. The timeout includes
	// connection time, any redirects, and header wait time.
	//
	// For non-watch GET request, server returns the response body immediately.
	// For PUT/POST/DELETE request, server will attempt to commit request
	// before responding, which is expected to take `100ms + 2 * RTT`.
	// For watch request, server returns the header immediately to notify Client
	// watch start. But if server is behind some kind of proxy, the response
	// header may be cached at proxy, and Client cannot rely on this behavior.
	//
	// Especially, wait request will ignore this timeout.
	//
	// One API call may send multiple requests to different etcd servers until it
	// succeeds. Use context of the API to specify the overall timeout.
	//
	// A HeaderTimeoutPerRequest of zero means no timeout.
	HeaderTimeoutPerRequest time.Duration

	// SelectionMode is an EndpointSelectionMode enum that specifies the
	// policy for choosing the etcd cluster node to which requests are sent.
	SelectionMode EndpointSelectionMode
***REMOVED***

func (cfg *Config) transport() CancelableTransport ***REMOVED***
	if cfg.Transport == nil ***REMOVED***
		return DefaultTransport
	***REMOVED***
	return cfg.Transport
***REMOVED***

func (cfg *Config) checkRedirect() CheckRedirectFunc ***REMOVED***
	if cfg.CheckRedirect == nil ***REMOVED***
		return DefaultCheckRedirect
	***REMOVED***
	return cfg.CheckRedirect
***REMOVED***

// CancelableTransport mimics net/http.Transport, but requires that
// the object also support request cancellation.
type CancelableTransport interface ***REMOVED***
	http.RoundTripper
	CancelRequest(req *http.Request)
***REMOVED***

type CheckRedirectFunc func(via int) error

// DefaultCheckRedirect follows up to 10 redirects, but no more.
var DefaultCheckRedirect CheckRedirectFunc = func(via int) error ***REMOVED***
	if via > 10 ***REMOVED***
		return ErrTooManyRedirects
	***REMOVED***
	return nil
***REMOVED***

type Client interface ***REMOVED***
	// Sync updates the internal cache of the etcd cluster's membership.
	Sync(context.Context) error

	// AutoSync periodically calls Sync() every given interval.
	// The recommended sync interval is 10 seconds to 1 minute, which does
	// not bring too much overhead to server and makes client catch up the
	// cluster change in time.
	//
	// The example to use it:
	//
	//  for ***REMOVED***
	//      err := client.AutoSync(ctx, 10*time.Second)
	//      if err == context.DeadlineExceeded || err == context.Canceled ***REMOVED***
	//          break
	//  ***REMOVED***
	//      log.Print(err)
	//  ***REMOVED***
	AutoSync(context.Context, time.Duration) error

	// Endpoints returns a copy of the current set of API endpoints used
	// by Client to resolve HTTP requests. If Sync has ever been called,
	// this may differ from the initial Endpoints provided in the Config.
	Endpoints() []string

	// SetEndpoints sets the set of API endpoints used by Client to resolve
	// HTTP requests. If the given endpoints are not valid, an error will be
	// returned
	SetEndpoints(eps []string) error

	// GetVersion retrieves the current etcd server and cluster version
	GetVersion(ctx context.Context) (*version.Versions, error)

	httpClient
***REMOVED***

func New(cfg Config) (Client, error) ***REMOVED***
	c := &httpClusterClient***REMOVED***
		clientFactory: newHTTPClientFactory(cfg.transport(), cfg.checkRedirect(), cfg.HeaderTimeoutPerRequest),
		rand:          rand.New(rand.NewSource(int64(time.Now().Nanosecond()))),
		selectionMode: cfg.SelectionMode,
	***REMOVED***
	if cfg.Username != "" ***REMOVED***
		c.credentials = &credentials***REMOVED***
			username: cfg.Username,
			password: cfg.Password,
		***REMOVED***
	***REMOVED***
	if err := c.SetEndpoints(cfg.Endpoints); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return c, nil
***REMOVED***

type httpClient interface ***REMOVED***
	Do(context.Context, httpAction) (*http.Response, []byte, error)
***REMOVED***

func newHTTPClientFactory(tr CancelableTransport, cr CheckRedirectFunc, headerTimeout time.Duration) httpClientFactory ***REMOVED***
	return func(ep url.URL) httpClient ***REMOVED***
		return &redirectFollowingHTTPClient***REMOVED***
			checkRedirect: cr,
			client: &simpleHTTPClient***REMOVED***
				transport:     tr,
				endpoint:      ep,
				headerTimeout: headerTimeout,
			***REMOVED***,
		***REMOVED***
	***REMOVED***
***REMOVED***

type credentials struct ***REMOVED***
	username string
	password string
***REMOVED***

type httpClientFactory func(url.URL) httpClient

type httpAction interface ***REMOVED***
	HTTPRequest(url.URL) *http.Request
***REMOVED***

type httpClusterClient struct ***REMOVED***
	clientFactory httpClientFactory
	endpoints     []url.URL
	pinned        int
	credentials   *credentials
	sync.RWMutex
	rand          *rand.Rand
	selectionMode EndpointSelectionMode
***REMOVED***

func (c *httpClusterClient) getLeaderEndpoint(ctx context.Context, eps []url.URL) (string, error) ***REMOVED***
	ceps := make([]url.URL, len(eps))
	copy(ceps, eps)

	// To perform a lookup on the new endpoint list without using the current
	// client, we'll copy it
	clientCopy := &httpClusterClient***REMOVED***
		clientFactory: c.clientFactory,
		credentials:   c.credentials,
		rand:          c.rand,

		pinned:    0,
		endpoints: ceps,
	***REMOVED***

	mAPI := NewMembersAPI(clientCopy)
	leader, err := mAPI.Leader(ctx)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if len(leader.ClientURLs) == 0 ***REMOVED***
		return "", ErrNoLeaderEndpoint
	***REMOVED***

	return leader.ClientURLs[0], nil // TODO: how to handle multiple client URLs?
***REMOVED***

func (c *httpClusterClient) parseEndpoints(eps []string) ([]url.URL, error) ***REMOVED***
	if len(eps) == 0 ***REMOVED***
		return []url.URL***REMOVED******REMOVED***, ErrNoEndpoints
	***REMOVED***

	neps := make([]url.URL, len(eps))
	for i, ep := range eps ***REMOVED***
		u, err := url.Parse(ep)
		if err != nil ***REMOVED***
			return []url.URL***REMOVED******REMOVED***, err
		***REMOVED***
		neps[i] = *u
	***REMOVED***
	return neps, nil
***REMOVED***

func (c *httpClusterClient) SetEndpoints(eps []string) error ***REMOVED***
	neps, err := c.parseEndpoints(eps)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c.Lock()
	defer c.Unlock()

	c.endpoints = shuffleEndpoints(c.rand, neps)
	// We're not doing anything for PrioritizeLeader here. This is
	// due to not having a context meaning we can't call getLeaderEndpoint
	// However, if you're using PrioritizeLeader, you've already been told
	// to regularly call sync, where we do have a ctx, and can figure the
	// leader. PrioritizeLeader is also quite a loose guarantee, so deal
	// with it
	c.pinned = 0

	return nil
***REMOVED***

func (c *httpClusterClient) Do(ctx context.Context, act httpAction) (*http.Response, []byte, error) ***REMOVED***
	action := act
	c.RLock()
	leps := len(c.endpoints)
	eps := make([]url.URL, leps)
	n := copy(eps, c.endpoints)
	pinned := c.pinned

	if c.credentials != nil ***REMOVED***
		action = &authedAction***REMOVED***
			act:         act,
			credentials: *c.credentials,
		***REMOVED***
	***REMOVED***
	c.RUnlock()

	if leps == 0 ***REMOVED***
		return nil, nil, ErrNoEndpoints
	***REMOVED***

	if leps != n ***REMOVED***
		return nil, nil, errors.New("unable to pick endpoint: copy failed")
	***REMOVED***

	var resp *http.Response
	var body []byte
	var err error
	cerr := &ClusterError***REMOVED******REMOVED***
	isOneShot := ctx.Value(&oneShotCtxValue) != nil

	for i := pinned; i < leps+pinned; i++ ***REMOVED***
		k := i % leps
		hc := c.clientFactory(eps[k])
		resp, body, err = hc.Do(ctx, action)
		if err != nil ***REMOVED***
			cerr.Errors = append(cerr.Errors, err)
			if err == ctx.Err() ***REMOVED***
				return nil, nil, ctx.Err()
			***REMOVED***
			if err == context.Canceled || err == context.DeadlineExceeded ***REMOVED***
				return nil, nil, err
			***REMOVED***
			if isOneShot ***REMOVED***
				return nil, nil, err
			***REMOVED***
			continue
		***REMOVED***
		if resp.StatusCode/100 == 5 ***REMOVED***
			switch resp.StatusCode ***REMOVED***
			case http.StatusInternalServerError, http.StatusServiceUnavailable:
				// TODO: make sure this is a no leader response
				cerr.Errors = append(cerr.Errors, fmt.Errorf("client: etcd member %s has no leader", eps[k].String()))
			default:
				cerr.Errors = append(cerr.Errors, fmt.Errorf("client: etcd member %s returns server error [%s]", eps[k].String(), http.StatusText(resp.StatusCode)))
			***REMOVED***
			if isOneShot ***REMOVED***
				return nil, nil, cerr.Errors[0]
			***REMOVED***
			continue
		***REMOVED***
		if k != pinned ***REMOVED***
			c.Lock()
			c.pinned = k
			c.Unlock()
		***REMOVED***
		return resp, body, nil
	***REMOVED***

	return nil, nil, cerr
***REMOVED***

func (c *httpClusterClient) Endpoints() []string ***REMOVED***
	c.RLock()
	defer c.RUnlock()

	eps := make([]string, len(c.endpoints))
	for i, ep := range c.endpoints ***REMOVED***
		eps[i] = ep.String()
	***REMOVED***

	return eps
***REMOVED***

func (c *httpClusterClient) Sync(ctx context.Context) error ***REMOVED***
	mAPI := NewMembersAPI(c)
	ms, err := mAPI.List(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var eps []string
	for _, m := range ms ***REMOVED***
		eps = append(eps, m.ClientURLs...)
	***REMOVED***

	neps, err := c.parseEndpoints(eps)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	npin := 0

	switch c.selectionMode ***REMOVED***
	case EndpointSelectionRandom:
		c.RLock()
		eq := endpointsEqual(c.endpoints, neps)
		c.RUnlock()

		if eq ***REMOVED***
			return nil
		***REMOVED***
		// When items in the endpoint list changes, we choose a new pin
		neps = shuffleEndpoints(c.rand, neps)
	case EndpointSelectionPrioritizeLeader:
		nle, err := c.getLeaderEndpoint(ctx, neps)
		if err != nil ***REMOVED***
			return ErrNoLeaderEndpoint
		***REMOVED***

		for i, n := range neps ***REMOVED***
			if n.String() == nle ***REMOVED***
				npin = i
				break
			***REMOVED***
		***REMOVED***
	default:
		return fmt.Errorf("invalid endpoint selection mode: %d", c.selectionMode)
	***REMOVED***

	c.Lock()
	defer c.Unlock()
	c.endpoints = neps
	c.pinned = npin

	return nil
***REMOVED***

func (c *httpClusterClient) AutoSync(ctx context.Context, interval time.Duration) error ***REMOVED***
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for ***REMOVED***
		err := c.Sync(ctx)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *httpClusterClient) GetVersion(ctx context.Context) (*version.Versions, error) ***REMOVED***
	act := &getAction***REMOVED***Prefix: "/version"***REMOVED***

	resp, body, err := c.Do(ctx, act)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	switch resp.StatusCode ***REMOVED***
	case http.StatusOK:
		if len(body) == 0 ***REMOVED***
			return nil, ErrEmptyBody
		***REMOVED***
		var vresp version.Versions
		if err := json.Unmarshal(body, &vresp); err != nil ***REMOVED***
			return nil, ErrInvalidJSON
		***REMOVED***
		return &vresp, nil
	default:
		var etcdErr Error
		if err := json.Unmarshal(body, &etcdErr); err != nil ***REMOVED***
			return nil, ErrInvalidJSON
		***REMOVED***
		return nil, etcdErr
	***REMOVED***
***REMOVED***

type roundTripResponse struct ***REMOVED***
	resp *http.Response
	err  error
***REMOVED***

type simpleHTTPClient struct ***REMOVED***
	transport     CancelableTransport
	endpoint      url.URL
	headerTimeout time.Duration
***REMOVED***

func (c *simpleHTTPClient) Do(ctx context.Context, act httpAction) (*http.Response, []byte, error) ***REMOVED***
	req := act.HTTPRequest(c.endpoint)

	if err := printcURL(req); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	isWait := false
	if req != nil && req.URL != nil ***REMOVED***
		ws := req.URL.Query().Get("wait")
		if len(ws) != 0 ***REMOVED***
			var err error
			isWait, err = strconv.ParseBool(ws)
			if err != nil ***REMOVED***
				return nil, nil, fmt.Errorf("wrong wait value %s (%v for %+v)", ws, err, req)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	var hctx context.Context
	var hcancel context.CancelFunc
	if !isWait && c.headerTimeout > 0 ***REMOVED***
		hctx, hcancel = context.WithTimeout(ctx, c.headerTimeout)
	***REMOVED*** else ***REMOVED***
		hctx, hcancel = context.WithCancel(ctx)
	***REMOVED***
	defer hcancel()

	reqcancel := requestCanceler(c.transport, req)

	rtchan := make(chan roundTripResponse, 1)
	go func() ***REMOVED***
		resp, err := c.transport.RoundTrip(req)
		rtchan <- roundTripResponse***REMOVED***resp: resp, err: err***REMOVED***
		close(rtchan)
	***REMOVED***()

	var resp *http.Response
	var err error

	select ***REMOVED***
	case rtresp := <-rtchan:
		resp, err = rtresp.resp, rtresp.err
	case <-hctx.Done():
		// cancel and wait for request to actually exit before continuing
		reqcancel()
		rtresp := <-rtchan
		resp = rtresp.resp
		switch ***REMOVED***
		case ctx.Err() != nil:
			err = ctx.Err()
		case hctx.Err() != nil:
			err = fmt.Errorf("client: endpoint %s exceeded header timeout", c.endpoint.String())
		default:
			panic("failed to get error from context")
		***REMOVED***
	***REMOVED***

	// always check for resp nil-ness to deal with possible
	// race conditions between channels above
	defer func() ***REMOVED***
		if resp != nil ***REMOVED***
			resp.Body.Close()
		***REMOVED***
	***REMOVED***()

	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	var body []byte
	done := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		body, err = ioutil.ReadAll(resp.Body)
		done <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***()

	select ***REMOVED***
	case <-ctx.Done():
		resp.Body.Close()
		<-done
		return nil, nil, ctx.Err()
	case <-done:
	***REMOVED***

	return resp, body, err
***REMOVED***

type authedAction struct ***REMOVED***
	act         httpAction
	credentials credentials
***REMOVED***

func (a *authedAction) HTTPRequest(url url.URL) *http.Request ***REMOVED***
	r := a.act.HTTPRequest(url)
	r.SetBasicAuth(a.credentials.username, a.credentials.password)
	return r
***REMOVED***

type redirectFollowingHTTPClient struct ***REMOVED***
	client        httpClient
	checkRedirect CheckRedirectFunc
***REMOVED***

func (r *redirectFollowingHTTPClient) Do(ctx context.Context, act httpAction) (*http.Response, []byte, error) ***REMOVED***
	next := act
	for i := 0; i < 100; i++ ***REMOVED***
		if i > 0 ***REMOVED***
			if err := r.checkRedirect(i); err != nil ***REMOVED***
				return nil, nil, err
			***REMOVED***
		***REMOVED***
		resp, body, err := r.client.Do(ctx, next)
		if err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***
		if resp.StatusCode/100 == 3 ***REMOVED***
			hdr := resp.Header.Get("Location")
			if hdr == "" ***REMOVED***
				return nil, nil, fmt.Errorf("Location header not set")
			***REMOVED***
			loc, err := url.Parse(hdr)
			if err != nil ***REMOVED***
				return nil, nil, fmt.Errorf("Location header not valid URL: %s", hdr)
			***REMOVED***
			next = &redirectedHTTPAction***REMOVED***
				action:   act,
				location: *loc,
			***REMOVED***
			continue
		***REMOVED***
		return resp, body, nil
	***REMOVED***

	return nil, nil, errTooManyRedirectChecks
***REMOVED***

type redirectedHTTPAction struct ***REMOVED***
	action   httpAction
	location url.URL
***REMOVED***

func (r *redirectedHTTPAction) HTTPRequest(ep url.URL) *http.Request ***REMOVED***
	orig := r.action.HTTPRequest(ep)
	orig.URL = &r.location
	return orig
***REMOVED***

func shuffleEndpoints(r *rand.Rand, eps []url.URL) []url.URL ***REMOVED***
	p := r.Perm(len(eps))
	neps := make([]url.URL, len(eps))
	for i, k := range p ***REMOVED***
		neps[i] = eps[k]
	***REMOVED***
	return neps
***REMOVED***

func endpointsEqual(left, right []url.URL) bool ***REMOVED***
	if len(left) != len(right) ***REMOVED***
		return false
	***REMOVED***

	sLeft := make([]string, len(left))
	sRight := make([]string, len(right))
	for i, l := range left ***REMOVED***
		sLeft[i] = l.String()
	***REMOVED***
	for i, r := range right ***REMOVED***
		sRight[i] = r.String()
	***REMOVED***

	sort.Strings(sLeft)
	sort.Strings(sRight)
	for i := range sLeft ***REMOVED***
		if sLeft[i] != sRight[i] ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***
