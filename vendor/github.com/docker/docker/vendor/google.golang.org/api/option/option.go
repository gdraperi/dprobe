// Package option contains options for Google API clients.
package option

import (
	"net/http"

	"golang.org/x/oauth2"
	"google.golang.org/api/internal"
	"google.golang.org/grpc"
)

// A ClientOption is an option for a Google API client.
type ClientOption interface ***REMOVED***
	Apply(*internal.DialSettings)
***REMOVED***

// WithTokenSource returns a ClientOption that specifies an OAuth2 token
// source to be used as the basis for authentication.
func WithTokenSource(s oauth2.TokenSource) ClientOption ***REMOVED***
	return withTokenSource***REMOVED***s***REMOVED***
***REMOVED***

type withTokenSource struct***REMOVED*** ts oauth2.TokenSource ***REMOVED***

func (w withTokenSource) Apply(o *internal.DialSettings) ***REMOVED***
	o.TokenSource = w.ts
***REMOVED***

// WithServiceAccountFile returns a ClientOption that uses a Google service
// account credentials file to authenticate.
// Use WithTokenSource with a token source created from
// golang.org/x/oauth2/google.JWTConfigFromJSON
// if reading the file from disk is not an option.
func WithServiceAccountFile(filename string) ClientOption ***REMOVED***
	return withServiceAccountFile(filename)
***REMOVED***

type withServiceAccountFile string

func (w withServiceAccountFile) Apply(o *internal.DialSettings) ***REMOVED***
	o.ServiceAccountJSONFilename = string(w)
***REMOVED***

// WithEndpoint returns a ClientOption that overrides the default endpoint
// to be used for a service.
func WithEndpoint(url string) ClientOption ***REMOVED***
	return withEndpoint(url)
***REMOVED***

type withEndpoint string

func (w withEndpoint) Apply(o *internal.DialSettings) ***REMOVED***
	o.Endpoint = string(w)
***REMOVED***

// WithScopes returns a ClientOption that overrides the default OAuth2 scopes
// to be used for a service.
func WithScopes(scope ...string) ClientOption ***REMOVED***
	return withScopes(scope)
***REMOVED***

type withScopes []string

func (w withScopes) Apply(o *internal.DialSettings) ***REMOVED***
	s := make([]string, len(w))
	copy(s, w)
	o.Scopes = s
***REMOVED***

// WithUserAgent returns a ClientOption that sets the User-Agent.
func WithUserAgent(ua string) ClientOption ***REMOVED***
	return withUA(ua)
***REMOVED***

type withUA string

func (w withUA) Apply(o *internal.DialSettings) ***REMOVED*** o.UserAgent = string(w) ***REMOVED***

// WithHTTPClient returns a ClientOption that specifies the HTTP client to use
// as the basis of communications. This option may only be used with services
// that support HTTP as their communication transport. When used, the
// WithHTTPClient option takes precedent over all other supplied options.
func WithHTTPClient(client *http.Client) ClientOption ***REMOVED***
	return withHTTPClient***REMOVED***client***REMOVED***
***REMOVED***

type withHTTPClient struct***REMOVED*** client *http.Client ***REMOVED***

func (w withHTTPClient) Apply(o *internal.DialSettings) ***REMOVED***
	o.HTTPClient = w.client
***REMOVED***

// WithGRPCConn returns a ClientOption that specifies the gRPC client
// connection to use as the basis of communications. This option many only be
// used with services that support gRPC as their communication transport. When
// used, the WithGRPCConn option takes precedent over all other supplied
// options.
func WithGRPCConn(conn *grpc.ClientConn) ClientOption ***REMOVED***
	return withGRPCConn***REMOVED***conn***REMOVED***
***REMOVED***

type withGRPCConn struct***REMOVED*** conn *grpc.ClientConn ***REMOVED***

func (w withGRPCConn) Apply(o *internal.DialSettings) ***REMOVED***
	o.GRPCConn = w.conn
***REMOVED***

// WithGRPCDialOption returns a ClientOption that appends a new grpc.DialOption
// to an underlying gRPC dial. It does not work with WithGRPCConn.
func WithGRPCDialOption(opt grpc.DialOption) ClientOption ***REMOVED***
	return withGRPCDialOption***REMOVED***opt***REMOVED***
***REMOVED***

type withGRPCDialOption struct***REMOVED*** opt grpc.DialOption ***REMOVED***

func (w withGRPCDialOption) Apply(o *internal.DialSettings) ***REMOVED***
	o.GRPCDialOpts = append(o.GRPCDialOpts, w.opt)
***REMOVED***

// WithGRPCConnectionPool returns a ClientOption that creates a pool of gRPC
// connections that requests will be balanced between.
// This is an EXPERIMENTAL API and may be changed or removed in the future.
func WithGRPCConnectionPool(size int) ClientOption ***REMOVED***
	return withGRPCConnectionPool(size)
***REMOVED***

type withGRPCConnectionPool int

func (w withGRPCConnectionPool) Apply(o *internal.DialSettings) ***REMOVED***
	balancer := grpc.RoundRobin(internal.NewPoolResolver(int(w), o))
	o.GRPCDialOpts = append(o.GRPCDialOpts, grpc.WithBalancer(balancer))
***REMOVED***

// WithAPIKey returns a ClientOption that specifies an API key to be used
// as the basis for authentication.
func WithAPIKey(apiKey string) ClientOption ***REMOVED***
	return withAPIKey(apiKey)
***REMOVED***

type withAPIKey string

func (w withAPIKey) Apply(o *internal.DialSettings) ***REMOVED*** o.APIKey = string(w) ***REMOVED***
