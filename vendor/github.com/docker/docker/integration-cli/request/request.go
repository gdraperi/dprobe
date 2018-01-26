package request

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api"
	dclient "github.com/docker/docker/client"
	"github.com/docker/docker/opts"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/go-connections/sockets"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/pkg/errors"
)

// Method creates a modifier that sets the specified string as the request method
func Method(method string) func(*http.Request) error ***REMOVED***
	return func(req *http.Request) error ***REMOVED***
		req.Method = method
		return nil
	***REMOVED***
***REMOVED***

// RawString sets the specified string as body for the request
func RawString(content string) func(*http.Request) error ***REMOVED***
	return RawContent(ioutil.NopCloser(strings.NewReader(content)))
***REMOVED***

// RawContent sets the specified reader as body for the request
func RawContent(reader io.ReadCloser) func(*http.Request) error ***REMOVED***
	return func(req *http.Request) error ***REMOVED***
		req.Body = reader
		return nil
	***REMOVED***
***REMOVED***

// ContentType sets the specified Content-Type request header
func ContentType(contentType string) func(*http.Request) error ***REMOVED***
	return func(req *http.Request) error ***REMOVED***
		req.Header.Set("Content-Type", contentType)
		return nil
	***REMOVED***
***REMOVED***

// JSON sets the Content-Type request header to json
func JSON(req *http.Request) error ***REMOVED***
	return ContentType("application/json")(req)
***REMOVED***

// JSONBody creates a modifier that encodes the specified data to a JSON string and set it as request body. It also sets
// the Content-Type header of the request.
func JSONBody(data interface***REMOVED******REMOVED***) func(*http.Request) error ***REMOVED***
	return func(req *http.Request) error ***REMOVED***
		jsonData := bytes.NewBuffer(nil)
		if err := json.NewEncoder(jsonData).Encode(data); err != nil ***REMOVED***
			return err
		***REMOVED***
		req.Body = ioutil.NopCloser(jsonData)
		req.Header.Set("Content-Type", "application/json")
		return nil
	***REMOVED***
***REMOVED***

// Post creates and execute a POST request on the specified host and endpoint, with the specified request modifiers
func Post(endpoint string, modifiers ...func(*http.Request) error) (*http.Response, io.ReadCloser, error) ***REMOVED***
	return Do(endpoint, append(modifiers, Method(http.MethodPost))...)
***REMOVED***

// Delete creates and execute a DELETE request on the specified host and endpoint, with the specified request modifiers
func Delete(endpoint string, modifiers ...func(*http.Request) error) (*http.Response, io.ReadCloser, error) ***REMOVED***
	return Do(endpoint, append(modifiers, Method(http.MethodDelete))...)
***REMOVED***

// Get creates and execute a GET request on the specified host and endpoint, with the specified request modifiers
func Get(endpoint string, modifiers ...func(*http.Request) error) (*http.Response, io.ReadCloser, error) ***REMOVED***
	return Do(endpoint, modifiers...)
***REMOVED***

// Do creates and execute a request on the specified endpoint, with the specified request modifiers
func Do(endpoint string, modifiers ...func(*http.Request) error) (*http.Response, io.ReadCloser, error) ***REMOVED***
	return DoOnHost(DaemonHost(), endpoint, modifiers...)
***REMOVED***

// DoOnHost creates and execute a request on the specified host and endpoint, with the specified request modifiers
func DoOnHost(host, endpoint string, modifiers ...func(*http.Request) error) (*http.Response, io.ReadCloser, error) ***REMOVED***
	req, err := New(host, endpoint, modifiers...)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	client, err := NewHTTPClient(host)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	resp, err := client.Do(req)
	var body io.ReadCloser
	if resp != nil ***REMOVED***
		body = ioutils.NewReadCloserWrapper(resp.Body, func() error ***REMOVED***
			defer resp.Body.Close()
			return nil
		***REMOVED***)
	***REMOVED***
	return resp, body, err
***REMOVED***

// New creates a new http Request to the specified host and endpoint, with the specified request modifiers
func New(host, endpoint string, modifiers ...func(*http.Request) error) (*http.Request, error) ***REMOVED***
	_, addr, _, err := dclient.ParseHost(host)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "could not parse url %q", host)
	***REMOVED***
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("could not create new request: %v", err)
	***REMOVED***

	if os.Getenv("DOCKER_TLS_VERIFY") != "" ***REMOVED***
		req.URL.Scheme = "https"
	***REMOVED*** else ***REMOVED***
		req.URL.Scheme = "http"
	***REMOVED***
	req.URL.Host = addr

	for _, config := range modifiers ***REMOVED***
		if err := config(req); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return req, nil
***REMOVED***

// NewHTTPClient creates an http client for the specific host
func NewHTTPClient(host string) (*http.Client, error) ***REMOVED***
	// FIXME(vdemeester) 10*time.Second timeout of SockRequest… ?
	proto, addr, _, err := dclient.ParseHost(host)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	transport := new(http.Transport)
	if proto == "tcp" && os.Getenv("DOCKER_TLS_VERIFY") != "" ***REMOVED***
		// Setup the socket TLS configuration.
		tlsConfig, err := getTLSConfig()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		transport = &http.Transport***REMOVED***TLSClientConfig: tlsConfig***REMOVED***
	***REMOVED***
	transport.DisableKeepAlives = true
	err = sockets.ConfigureTransport(transport, proto, addr)
	return &http.Client***REMOVED***
		Transport: transport,
	***REMOVED***, err
***REMOVED***

// NewClient returns a new Docker API client
func NewClient() (dclient.APIClient, error) ***REMOVED***
	return NewClientForHost(DaemonHost())
***REMOVED***

// NewClientForHost returns a Docker API client for the host
func NewClientForHost(host string) (dclient.APIClient, error) ***REMOVED***
	httpClient, err := NewHTTPClient(host)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return dclient.NewClient(host, api.DefaultVersion, httpClient, nil)
***REMOVED***

// FIXME(vdemeester) httputil.ClientConn is deprecated, use http.Client instead (closer to actual client)
// Deprecated: Use New instead of NewRequestClient
// Deprecated: use request.Do (or Get, Delete, Post) instead
func newRequestClient(method, endpoint string, data io.Reader, ct, daemon string, modifiers ...func(*http.Request)) (*http.Request, *httputil.ClientConn, error) ***REMOVED***
	c, err := SockConn(time.Duration(10*time.Second), daemon)
	if err != nil ***REMOVED***
		return nil, nil, fmt.Errorf("could not dial docker daemon: %v", err)
	***REMOVED***

	client := httputil.NewClientConn(c, nil)

	req, err := http.NewRequest(method, endpoint, data)
	if err != nil ***REMOVED***
		client.Close()
		return nil, nil, fmt.Errorf("could not create new request: %v", err)
	***REMOVED***

	for _, opt := range modifiers ***REMOVED***
		opt(req)
	***REMOVED***

	if ct != "" ***REMOVED***
		req.Header.Set("Content-Type", ct)
	***REMOVED***
	return req, client, nil
***REMOVED***

// SockRequest create a request against the specified host (with method, endpoint and other request modifier) and
// returns the status code, and the content as an byte slice
// Deprecated: use request.Do instead
func SockRequest(method, endpoint string, data interface***REMOVED******REMOVED***, daemon string, modifiers ...func(*http.Request)) (int, []byte, error) ***REMOVED***
	jsonData := bytes.NewBuffer(nil)
	if err := json.NewEncoder(jsonData).Encode(data); err != nil ***REMOVED***
		return -1, nil, err
	***REMOVED***

	res, body, err := SockRequestRaw(method, endpoint, jsonData, "application/json", daemon, modifiers...)
	if err != nil ***REMOVED***
		return -1, nil, err
	***REMOVED***
	b, err := ReadBody(body)
	return res.StatusCode, b, err
***REMOVED***

// ReadBody read the specified ReadCloser content and returns it
func ReadBody(b io.ReadCloser) ([]byte, error) ***REMOVED***
	defer b.Close()
	return ioutil.ReadAll(b)
***REMOVED***

// SockRequestRaw create a request against the specified host (with method, endpoint and other request modifier) and
// returns the http response, the output as a io.ReadCloser
// Deprecated: use request.Do (or Get, Delete, Post) instead
func SockRequestRaw(method, endpoint string, data io.Reader, ct, daemon string, modifiers ...func(*http.Request)) (*http.Response, io.ReadCloser, error) ***REMOVED***
	req, client, err := newRequestClient(method, endpoint, data, ct, daemon, modifiers...)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	resp, err := client.Do(req)
	if err != nil ***REMOVED***
		client.Close()
		return resp, nil, err
	***REMOVED***
	body := ioutils.NewReadCloserWrapper(resp.Body, func() error ***REMOVED***
		defer resp.Body.Close()
		return client.Close()
	***REMOVED***)

	return resp, body, err
***REMOVED***

// SockRequestHijack creates a connection to specified host (with method, contenttype, …) and returns a hijacked connection
// and the output as a `bufio.Reader`
func SockRequestHijack(method, endpoint string, data io.Reader, ct string, daemon string, modifiers ...func(*http.Request)) (net.Conn, *bufio.Reader, error) ***REMOVED***
	req, client, err := newRequestClient(method, endpoint, data, ct, daemon, modifiers...)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	client.Do(req)
	conn, br := client.Hijack()
	return conn, br, nil
***REMOVED***

// SockConn opens a connection on the specified socket
func SockConn(timeout time.Duration, daemon string) (net.Conn, error) ***REMOVED***
	daemonURL, err := url.Parse(daemon)
	if err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "could not parse url %q", daemon)
	***REMOVED***

	var c net.Conn
	switch daemonURL.Scheme ***REMOVED***
	case "npipe":
		return npipeDial(daemonURL.Path, timeout)
	case "unix":
		return net.DialTimeout(daemonURL.Scheme, daemonURL.Path, timeout)
	case "tcp":
		if os.Getenv("DOCKER_TLS_VERIFY") != "" ***REMOVED***
			// Setup the socket TLS configuration.
			tlsConfig, err := getTLSConfig()
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			dialer := &net.Dialer***REMOVED***Timeout: timeout***REMOVED***
			return tls.DialWithDialer(dialer, daemonURL.Scheme, daemonURL.Host, tlsConfig)
		***REMOVED***
		return net.DialTimeout(daemonURL.Scheme, daemonURL.Host, timeout)
	default:
		return c, errors.Errorf("unknown scheme %v (%s)", daemonURL.Scheme, daemon)
	***REMOVED***
***REMOVED***

func getTLSConfig() (*tls.Config, error) ***REMOVED***
	dockerCertPath := os.Getenv("DOCKER_CERT_PATH")

	if dockerCertPath == "" ***REMOVED***
		return nil, errors.New("DOCKER_TLS_VERIFY specified, but no DOCKER_CERT_PATH environment variable")
	***REMOVED***

	option := &tlsconfig.Options***REMOVED***
		CAFile:   filepath.Join(dockerCertPath, "ca.pem"),
		CertFile: filepath.Join(dockerCertPath, "cert.pem"),
		KeyFile:  filepath.Join(dockerCertPath, "key.pem"),
	***REMOVED***
	tlsConfig, err := tlsconfig.Client(*option)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return tlsConfig, nil
***REMOVED***

// DaemonHost return the daemon host string for this test execution
func DaemonHost() string ***REMOVED***
	daemonURLStr := "unix://" + opts.DefaultUnixSocket
	if daemonHostVar := os.Getenv("DOCKER_HOST"); daemonHostVar != "" ***REMOVED***
		daemonURLStr = daemonHostVar
	***REMOVED***
	return daemonURLStr
***REMOVED***

// NewEnvClientWithVersion returns a docker client with a specified version.
// See: github.com/docker/docker/client `NewEnvClient()`
func NewEnvClientWithVersion(version string) (*dclient.Client, error) ***REMOVED***
	if version == "" ***REMOVED***
		return nil, errors.New("version not specified")
	***REMOVED***

	var httpClient *http.Client
	if os.Getenv("DOCKER_CERT_PATH") != "" ***REMOVED***
		tlsConfig, err := getTLSConfig()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		httpClient = &http.Client***REMOVED***
			Transport: &http.Transport***REMOVED***
				TLSClientConfig: tlsConfig,
			***REMOVED***,
		***REMOVED***
	***REMOVED***

	host := os.Getenv("DOCKER_HOST")
	if host == "" ***REMOVED***
		host = dclient.DefaultDockerHost
	***REMOVED***

	cli, err := dclient.NewClient(host, version, httpClient, nil)
	if err != nil ***REMOVED***
		return cli, err
	***REMOVED***
	return cli, nil
***REMOVED***
