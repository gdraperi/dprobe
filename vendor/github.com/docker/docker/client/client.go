/*
Package client is a Go client for the Docker Engine API.

For more information about the Engine API, see the documentation:
https://docs.docker.com/engine/reference/api/

Usage

You use the library by creating a client object and calling methods on it. The
client can be created either from environment variables with NewEnvClient, or
configured manually with NewClient.

For example, to list running containers (the equivalent of "docker ps"):

	package main

	import (
		"context"
		"fmt"

		"github.com/docker/docker/api/types"
		"github.com/docker/docker/client"
	)

	func main() ***REMOVED***
		cli, err := client.NewEnvClient()
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***

		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions***REMOVED******REMOVED***)
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***

		for _, container := range containers ***REMOVED***
			fmt.Printf("%s %s\n", container.ID[:10], container.Image)
		***REMOVED***
	***REMOVED***

*/
package client

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/versions"
	"github.com/docker/go-connections/sockets"
	"github.com/docker/go-connections/tlsconfig"
	"golang.org/x/net/context"
)

// ErrRedirect is the error returned by checkRedirect when the request is non-GET.
var ErrRedirect = errors.New("unexpected redirect in response")

// Client is the API client that performs all operations
// against a docker server.
type Client struct ***REMOVED***
	// scheme sets the scheme for the client
	scheme string
	// host holds the server address to connect to
	host string
	// proto holds the client protocol i.e. unix.
	proto string
	// addr holds the client address.
	addr string
	// basePath holds the path to prepend to the requests.
	basePath string
	// client used to send and receive http requests.
	client *http.Client
	// version of the server to talk to.
	version string
	// custom http headers configured by users.
	customHTTPHeaders map[string]string
	// manualOverride is set to true when the version was set by users.
	manualOverride bool
***REMOVED***

// CheckRedirect specifies the policy for dealing with redirect responses:
// If the request is non-GET return `ErrRedirect`. Otherwise use the last response.
//
// Go 1.8 changes behavior for HTTP redirects (specifically 301, 307, and 308) in the client .
// The Docker client (and by extension docker API client) can be made to to send a request
// like POST /containers//start where what would normally be in the name section of the URL is empty.
// This triggers an HTTP 301 from the daemon.
// In go 1.8 this 301 will be converted to a GET request, and ends up getting a 404 from the daemon.
// This behavior change manifests in the client in that before the 301 was not followed and
// the client did not generate an error, but now results in a message like Error response from daemon: page not found.
func CheckRedirect(req *http.Request, via []*http.Request) error ***REMOVED***
	if via[0].Method == http.MethodGet ***REMOVED***
		return http.ErrUseLastResponse
	***REMOVED***
	return ErrRedirect
***REMOVED***

// NewEnvClient initializes a new API client based on environment variables.
// Use DOCKER_HOST to set the url to the docker server.
// Use DOCKER_API_VERSION to set the version of the API to reach, leave empty for latest.
// Use DOCKER_CERT_PATH to load the TLS certificates from.
// Use DOCKER_TLS_VERIFY to enable or disable TLS verification, off by default.
func NewEnvClient() (*Client, error) ***REMOVED***
	var client *http.Client
	if dockerCertPath := os.Getenv("DOCKER_CERT_PATH"); dockerCertPath != "" ***REMOVED***
		options := tlsconfig.Options***REMOVED***
			CAFile:             filepath.Join(dockerCertPath, "ca.pem"),
			CertFile:           filepath.Join(dockerCertPath, "cert.pem"),
			KeyFile:            filepath.Join(dockerCertPath, "key.pem"),
			InsecureSkipVerify: os.Getenv("DOCKER_TLS_VERIFY") == "",
		***REMOVED***
		tlsc, err := tlsconfig.Client(options)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		client = &http.Client***REMOVED***
			Transport: &http.Transport***REMOVED***
				TLSClientConfig: tlsc,
			***REMOVED***,
			CheckRedirect: CheckRedirect,
		***REMOVED***
	***REMOVED***

	host := os.Getenv("DOCKER_HOST")
	if host == "" ***REMOVED***
		host = DefaultDockerHost
	***REMOVED***
	version := os.Getenv("DOCKER_API_VERSION")
	if version == "" ***REMOVED***
		version = api.DefaultVersion
	***REMOVED***

	cli, err := NewClient(host, version, client, nil)
	if err != nil ***REMOVED***
		return cli, err
	***REMOVED***
	if os.Getenv("DOCKER_API_VERSION") != "" ***REMOVED***
		cli.manualOverride = true
	***REMOVED***
	return cli, nil
***REMOVED***

// NewClient initializes a new API client for the given host and API version.
// It uses the given http client as transport.
// It also initializes the custom http headers to add to each request.
//
// It won't send any version information if the version number is empty. It is
// highly recommended that you set a version or your client may break if the
// server is upgraded.
func NewClient(host string, version string, client *http.Client, httpHeaders map[string]string) (*Client, error) ***REMOVED***
	hostURL, err := ParseHostURL(host)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if client != nil ***REMOVED***
		if _, ok := client.Transport.(http.RoundTripper); !ok ***REMOVED***
			return nil, fmt.Errorf("unable to verify TLS configuration, invalid transport %v", client.Transport)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		transport := new(http.Transport)
		sockets.ConfigureTransport(transport, hostURL.Scheme, hostURL.Host)
		client = &http.Client***REMOVED***
			Transport:     transport,
			CheckRedirect: CheckRedirect,
		***REMOVED***
	***REMOVED***

	scheme := "http"
	tlsConfig := resolveTLSConfig(client.Transport)
	if tlsConfig != nil ***REMOVED***
		// TODO(stevvooe): This isn't really the right way to write clients in Go.
		// `NewClient` should probably only take an `*http.Client` and work from there.
		// Unfortunately, the model of having a host-ish/url-thingy as the connection
		// string has us confusing protocol and transport layers. We continue doing
		// this to avoid breaking existing clients but this should be addressed.
		scheme = "https"
	***REMOVED***

	// TODO: store URL instead of proto/addr/basePath
	return &Client***REMOVED***
		scheme:            scheme,
		host:              host,
		proto:             hostURL.Scheme,
		addr:              hostURL.Host,
		basePath:          hostURL.Path,
		client:            client,
		version:           version,
		customHTTPHeaders: httpHeaders,
	***REMOVED***, nil
***REMOVED***

// Close the transport used by the client
func (cli *Client) Close() error ***REMOVED***
	if t, ok := cli.client.Transport.(*http.Transport); ok ***REMOVED***
		t.CloseIdleConnections()
	***REMOVED***
	return nil
***REMOVED***

// getAPIPath returns the versioned request path to call the api.
// It appends the query parameters to the path if they are not empty.
func (cli *Client) getAPIPath(p string, query url.Values) string ***REMOVED***
	var apiPath string
	if cli.version != "" ***REMOVED***
		v := strings.TrimPrefix(cli.version, "v")
		apiPath = path.Join(cli.basePath, "/v"+v, p)
	***REMOVED*** else ***REMOVED***
		apiPath = path.Join(cli.basePath, p)
	***REMOVED***
	return (&url.URL***REMOVED***Path: apiPath, RawQuery: query.Encode()***REMOVED***).String()
***REMOVED***

// ClientVersion returns the API version used by this client.
func (cli *Client) ClientVersion() string ***REMOVED***
	return cli.version
***REMOVED***

// NegotiateAPIVersion queries the API and updates the version to match the
// API version. Any errors are silently ignored.
func (cli *Client) NegotiateAPIVersion(ctx context.Context) ***REMOVED***
	ping, _ := cli.Ping(ctx)
	cli.NegotiateAPIVersionPing(ping)
***REMOVED***

// NegotiateAPIVersionPing updates the client version to match the Ping.APIVersion
// if the ping version is less than the default version.
func (cli *Client) NegotiateAPIVersionPing(p types.Ping) ***REMOVED***
	if cli.manualOverride ***REMOVED***
		return
	***REMOVED***

	// try the latest version before versioning headers existed
	if p.APIVersion == "" ***REMOVED***
		p.APIVersion = "1.24"
	***REMOVED***

	// if the client is not initialized with a version, start with the latest supported version
	if cli.version == "" ***REMOVED***
		cli.version = api.DefaultVersion
	***REMOVED***

	// if server version is lower than the client version, downgrade
	if versions.LessThan(p.APIVersion, cli.version) ***REMOVED***
		cli.version = p.APIVersion
	***REMOVED***
***REMOVED***

// DaemonHost returns the host address used by the client
func (cli *Client) DaemonHost() string ***REMOVED***
	return cli.host
***REMOVED***

// ParseHost parses a url string, validates the strings is a host url, and returns
// the parsed host as: protocol, address, and base path
// Deprecated: use ParseHostURL
func ParseHost(host string) (string, string, string, error) ***REMOVED***
	hostURL, err := ParseHostURL(host)
	if err != nil ***REMOVED***
		return "", "", "", err
	***REMOVED***
	return hostURL.Scheme, hostURL.Host, hostURL.Path, nil
***REMOVED***

// ParseHostURL parses a url string, validates the string is a host url, and
// returns the parsed URL
func ParseHostURL(host string) (*url.URL, error) ***REMOVED***
	protoAddrParts := strings.SplitN(host, "://", 2)
	if len(protoAddrParts) == 1 ***REMOVED***
		return nil, fmt.Errorf("unable to parse docker host `%s`", host)
	***REMOVED***

	var basePath string
	proto, addr := protoAddrParts[0], protoAddrParts[1]
	if proto == "tcp" ***REMOVED***
		parsed, err := url.Parse("tcp://" + addr)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		addr = parsed.Host
		basePath = parsed.Path
	***REMOVED***
	return &url.URL***REMOVED***
		Scheme: proto,
		Host:   addr,
		Path:   basePath,
	***REMOVED***, nil
***REMOVED***

// CustomHTTPHeaders returns the custom http headers stored by the client.
func (cli *Client) CustomHTTPHeaders() map[string]string ***REMOVED***
	m := make(map[string]string)
	for k, v := range cli.customHTTPHeaders ***REMOVED***
		m[k] = v
	***REMOVED***
	return m
***REMOVED***

// SetCustomHTTPHeaders that will be set on every HTTP request made by the client.
func (cli *Client) SetCustomHTTPHeaders(headers map[string]string) ***REMOVED***
	cli.customHTTPHeaders = headers
***REMOVED***
