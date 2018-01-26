package plugins

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/plugins/transport"
	"github.com/docker/go-connections/sockets"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/sirupsen/logrus"
)

const (
	defaultTimeOut = 30
)

func newTransport(addr string, tlsConfig *tlsconfig.Options) (transport.Transport, error) ***REMOVED***
	tr := &http.Transport***REMOVED******REMOVED***

	if tlsConfig != nil ***REMOVED***
		c, err := tlsconfig.Client(*tlsConfig)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		tr.TLSClientConfig = c
	***REMOVED***

	u, err := url.Parse(addr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	socket := u.Host
	if socket == "" ***REMOVED***
		// valid local socket addresses have the host empty.
		socket = u.Path
	***REMOVED***
	if err := sockets.ConfigureTransport(tr, u.Scheme, socket); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	scheme := httpScheme(u)

	return transport.NewHTTPTransport(tr, scheme, socket), nil
***REMOVED***

// NewClient creates a new plugin client (http).
func NewClient(addr string, tlsConfig *tlsconfig.Options) (*Client, error) ***REMOVED***
	clientTransport, err := newTransport(addr, tlsConfig)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return newClientWithTransport(clientTransport, 0), nil
***REMOVED***

// NewClientWithTimeout creates a new plugin client (http).
func NewClientWithTimeout(addr string, tlsConfig *tlsconfig.Options, timeout time.Duration) (*Client, error) ***REMOVED***
	clientTransport, err := newTransport(addr, tlsConfig)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return newClientWithTransport(clientTransport, timeout), nil
***REMOVED***

// newClientWithTransport creates a new plugin client with a given transport.
func newClientWithTransport(tr transport.Transport, timeout time.Duration) *Client ***REMOVED***
	return &Client***REMOVED***
		http: &http.Client***REMOVED***
			Transport: tr,
			Timeout:   timeout,
		***REMOVED***,
		requestFactory: tr,
	***REMOVED***
***REMOVED***

// Client represents a plugin client.
type Client struct ***REMOVED***
	http           *http.Client // http client to use
	requestFactory transport.RequestFactory
***REMOVED***

// RequestOpts is the set of options that can be passed into a request
type RequestOpts struct ***REMOVED***
	Timeout time.Duration
***REMOVED***

// WithRequestTimeout sets a timeout duration for plugin requests
func WithRequestTimeout(t time.Duration) func(*RequestOpts) ***REMOVED***
	return func(o *RequestOpts) ***REMOVED***
		o.Timeout = t
	***REMOVED***
***REMOVED***

// Call calls the specified method with the specified arguments for the plugin.
// It will retry for 30 seconds if a failure occurs when calling.
func (c *Client) Call(serviceMethod string, args, ret interface***REMOVED******REMOVED***) error ***REMOVED***
	return c.CallWithOptions(serviceMethod, args, ret)
***REMOVED***

// CallWithOptions is just like call except it takes options
func (c *Client) CallWithOptions(serviceMethod string, args interface***REMOVED******REMOVED***, ret interface***REMOVED******REMOVED***, opts ...func(*RequestOpts)) error ***REMOVED***
	var buf bytes.Buffer
	if args != nil ***REMOVED***
		if err := json.NewEncoder(&buf).Encode(args); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	body, err := c.callWithRetry(serviceMethod, &buf, true, opts...)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer body.Close()
	if ret != nil ***REMOVED***
		if err := json.NewDecoder(body).Decode(&ret); err != nil ***REMOVED***
			logrus.Errorf("%s: error reading plugin resp: %v", serviceMethod, err)
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Stream calls the specified method with the specified arguments for the plugin and returns the response body
func (c *Client) Stream(serviceMethod string, args interface***REMOVED******REMOVED***) (io.ReadCloser, error) ***REMOVED***
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(args); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return c.callWithRetry(serviceMethod, &buf, true)
***REMOVED***

// SendFile calls the specified method, and passes through the IO stream
func (c *Client) SendFile(serviceMethod string, data io.Reader, ret interface***REMOVED******REMOVED***) error ***REMOVED***
	body, err := c.callWithRetry(serviceMethod, data, true)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer body.Close()
	if err := json.NewDecoder(body).Decode(&ret); err != nil ***REMOVED***
		logrus.Errorf("%s: error reading plugin resp: %v", serviceMethod, err)
		return err
	***REMOVED***
	return nil
***REMOVED***

func (c *Client) callWithRetry(serviceMethod string, data io.Reader, retry bool, reqOpts ...func(*RequestOpts)) (io.ReadCloser, error) ***REMOVED***
	var retries int
	start := time.Now()

	var opts RequestOpts
	for _, o := range reqOpts ***REMOVED***
		o(&opts)
	***REMOVED***

	for ***REMOVED***
		req, err := c.requestFactory.NewRequest(serviceMethod, data)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		cancelRequest := func() ***REMOVED******REMOVED***
		if opts.Timeout > 0 ***REMOVED***
			var ctx context.Context
			ctx, cancelRequest = context.WithTimeout(req.Context(), opts.Timeout)
			req = req.WithContext(ctx)
		***REMOVED***

		resp, err := c.http.Do(req)
		if err != nil ***REMOVED***
			cancelRequest()
			if !retry ***REMOVED***
				return nil, err
			***REMOVED***

			timeOff := backoff(retries)
			if abort(start, timeOff) ***REMOVED***
				return nil, err
			***REMOVED***
			retries++
			logrus.Warnf("Unable to connect to plugin: %s%s: %v, retrying in %v", req.URL.Host, req.URL.Path, err, timeOff)
			time.Sleep(timeOff)
			continue
		***REMOVED***

		if resp.StatusCode != http.StatusOK ***REMOVED***
			b, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			cancelRequest()
			if err != nil ***REMOVED***
				return nil, &statusError***REMOVED***resp.StatusCode, serviceMethod, err.Error()***REMOVED***
			***REMOVED***

			// Plugins' Response(s) should have an Err field indicating what went
			// wrong. Try to unmarshal into ResponseErr. Otherwise fallback to just
			// return the string(body)
			type responseErr struct ***REMOVED***
				Err string
			***REMOVED***
			remoteErr := responseErr***REMOVED******REMOVED***
			if err := json.Unmarshal(b, &remoteErr); err == nil ***REMOVED***
				if remoteErr.Err != "" ***REMOVED***
					return nil, &statusError***REMOVED***resp.StatusCode, serviceMethod, remoteErr.Err***REMOVED***
				***REMOVED***
			***REMOVED***
			// old way...
			return nil, &statusError***REMOVED***resp.StatusCode, serviceMethod, string(b)***REMOVED***
		***REMOVED***
		return ioutils.NewReadCloserWrapper(resp.Body, func() error ***REMOVED***
			err := resp.Body.Close()
			cancelRequest()
			return err
		***REMOVED***), nil
	***REMOVED***
***REMOVED***

func backoff(retries int) time.Duration ***REMOVED***
	b, max := 1, defaultTimeOut
	for b < max && retries > 0 ***REMOVED***
		b *= 2
		retries--
	***REMOVED***
	if b > max ***REMOVED***
		b = max
	***REMOVED***
	return time.Duration(b) * time.Second
***REMOVED***

func abort(start time.Time, timeOff time.Duration) bool ***REMOVED***
	return timeOff+time.Since(start) >= time.Duration(defaultTimeOut)*time.Second
***REMOVED***

func httpScheme(u *url.URL) string ***REMOVED***
	scheme := u.Scheme
	if scheme != "https" ***REMOVED***
		scheme = "http"
	***REMOVED***
	return scheme
***REMOVED***
