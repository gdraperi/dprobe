package client

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client/metadata"
	"github.com/aws/aws-sdk-go/aws/request"
)

// A Config provides configuration to a service client instance.
type Config struct ***REMOVED***
	Config        *aws.Config
	Handlers      request.Handlers
	Endpoint      string
	SigningRegion string
	SigningName   string
***REMOVED***

// ConfigProvider provides a generic way for a service client to receive
// the ClientConfig without circular dependencies.
type ConfigProvider interface ***REMOVED***
	ClientConfig(serviceName string, cfgs ...*aws.Config) Config
***REMOVED***

// ConfigNoResolveEndpointProvider same as ConfigProvider except it will not
// resolve the endpoint automatically. The service client's endpoint must be
// provided via the aws.Config.Endpoint field.
type ConfigNoResolveEndpointProvider interface ***REMOVED***
	ClientConfigNoResolveEndpoint(cfgs ...*aws.Config) Config
***REMOVED***

// A Client implements the base client request and response handling
// used by all service clients.
type Client struct ***REMOVED***
	request.Retryer
	metadata.ClientInfo

	Config   aws.Config
	Handlers request.Handlers
***REMOVED***

// New will return a pointer to a new initialized service client.
func New(cfg aws.Config, info metadata.ClientInfo, handlers request.Handlers, options ...func(*Client)) *Client ***REMOVED***
	svc := &Client***REMOVED***
		Config:     cfg,
		ClientInfo: info,
		Handlers:   handlers.Copy(),
	***REMOVED***

	switch retryer, ok := cfg.Retryer.(request.Retryer); ***REMOVED***
	case ok:
		svc.Retryer = retryer
	case cfg.Retryer != nil && cfg.Logger != nil:
		s := fmt.Sprintf("WARNING: %T does not implement request.Retryer; using DefaultRetryer instead", cfg.Retryer)
		cfg.Logger.Log(s)
		fallthrough
	default:
		maxRetries := aws.IntValue(cfg.MaxRetries)
		if cfg.MaxRetries == nil || maxRetries == aws.UseServiceDefaultRetries ***REMOVED***
			maxRetries = 3
		***REMOVED***
		svc.Retryer = DefaultRetryer***REMOVED***NumMaxRetries: maxRetries***REMOVED***
	***REMOVED***

	svc.AddDebugHandlers()

	for _, option := range options ***REMOVED***
		option(svc)
	***REMOVED***

	return svc
***REMOVED***

// NewRequest returns a new Request pointer for the service API
// operation and parameters.
func (c *Client) NewRequest(operation *request.Operation, params interface***REMOVED******REMOVED***, data interface***REMOVED******REMOVED***) *request.Request ***REMOVED***
	return request.New(c.Config, c.ClientInfo, c.Handlers, c.Retryer, operation, params, data)
***REMOVED***

// AddDebugHandlers injects debug logging handlers into the service to log request
// debug information.
func (c *Client) AddDebugHandlers() ***REMOVED***
	if !c.Config.LogLevel.AtLeast(aws.LogDebug) ***REMOVED***
		return
	***REMOVED***

	c.Handlers.Send.PushFrontNamed(request.NamedHandler***REMOVED***Name: "awssdk.client.LogRequest", Fn: logRequest***REMOVED***)
	c.Handlers.Send.PushBackNamed(request.NamedHandler***REMOVED***Name: "awssdk.client.LogResponse", Fn: logResponse***REMOVED***)
***REMOVED***
