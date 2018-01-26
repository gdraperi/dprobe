// Package defaults is a collection of helpers to retrieve the SDK's default
// configuration and handlers.
//
// Generally this package shouldn't be used directly, but session.Session
// instead. This package is useful when you need to reset the defaults
// of a session or service client to the SDK defaults before setting
// additional parameters.
package defaults

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/corehandlers"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/credentials/endpointcreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
)

// A Defaults provides a collection of default values for SDK clients.
type Defaults struct ***REMOVED***
	Config   *aws.Config
	Handlers request.Handlers
***REMOVED***

// Get returns the SDK's default values with Config and handlers pre-configured.
func Get() Defaults ***REMOVED***
	cfg := Config()
	handlers := Handlers()
	cfg.Credentials = CredChain(cfg, handlers)

	return Defaults***REMOVED***
		Config:   cfg,
		Handlers: handlers,
	***REMOVED***
***REMOVED***

// Config returns the default configuration without credentials.
// To retrieve a config with credentials also included use
// `defaults.Get().Config` instead.
//
// Generally you shouldn't need to use this method directly, but
// is available if you need to reset the configuration of an
// existing service client or session.
func Config() *aws.Config ***REMOVED***
	return aws.NewConfig().
		WithCredentials(credentials.AnonymousCredentials).
		WithRegion(os.Getenv("AWS_REGION")).
		WithHTTPClient(http.DefaultClient).
		WithMaxRetries(aws.UseServiceDefaultRetries).
		WithLogger(aws.NewDefaultLogger()).
		WithLogLevel(aws.LogOff).
		WithEndpointResolver(endpoints.DefaultResolver())
***REMOVED***

// Handlers returns the default request handlers.
//
// Generally you shouldn't need to use this method directly, but
// is available if you need to reset the request handlers of an
// existing service client or session.
func Handlers() request.Handlers ***REMOVED***
	var handlers request.Handlers

	handlers.Validate.PushBackNamed(corehandlers.ValidateEndpointHandler)
	handlers.Validate.AfterEachFn = request.HandlerListStopOnError
	handlers.Build.PushBackNamed(corehandlers.SDKVersionUserAgentHandler)
	handlers.Build.AfterEachFn = request.HandlerListStopOnError
	handlers.Sign.PushBackNamed(corehandlers.BuildContentLengthHandler)
	handlers.Send.PushBackNamed(corehandlers.ValidateReqSigHandler)
	handlers.Send.PushBackNamed(corehandlers.SendHandler)
	handlers.AfterRetry.PushBackNamed(corehandlers.AfterRetryHandler)
	handlers.ValidateResponse.PushBackNamed(corehandlers.ValidateResponseHandler)

	return handlers
***REMOVED***

// CredChain returns the default credential chain.
//
// Generally you shouldn't need to use this method directly, but
// is available if you need to reset the credentials of an
// existing service client or session's Config.
func CredChain(cfg *aws.Config, handlers request.Handlers) *credentials.Credentials ***REMOVED***
	return credentials.NewCredentials(&credentials.ChainProvider***REMOVED***
		VerboseErrors: aws.BoolValue(cfg.CredentialsChainVerboseErrors),
		Providers: []credentials.Provider***REMOVED***
			&credentials.EnvProvider***REMOVED******REMOVED***,
			&credentials.SharedCredentialsProvider***REMOVED***Filename: "", Profile: ""***REMOVED***,
			RemoteCredProvider(*cfg, handlers),
		***REMOVED***,
	***REMOVED***)
***REMOVED***

const (
	httpProviderEnvVar     = "AWS_CONTAINER_CREDENTIALS_FULL_URI"
	ecsCredsProviderEnvVar = "AWS_CONTAINER_CREDENTIALS_RELATIVE_URI"
)

// RemoteCredProvider returns a credentials provider for the default remote
// endpoints such as EC2 or ECS Roles.
func RemoteCredProvider(cfg aws.Config, handlers request.Handlers) credentials.Provider ***REMOVED***
	if u := os.Getenv(httpProviderEnvVar); len(u) > 0 ***REMOVED***
		return localHTTPCredProvider(cfg, handlers, u)
	***REMOVED***

	if uri := os.Getenv(ecsCredsProviderEnvVar); len(uri) > 0 ***REMOVED***
		u := fmt.Sprintf("http://169.254.170.2%s", uri)
		return httpCredProvider(cfg, handlers, u)
	***REMOVED***

	return ec2RoleProvider(cfg, handlers)
***REMOVED***

var lookupHostFn = net.LookupHost

func isLoopbackHost(host string) (bool, error) ***REMOVED***
	ip := net.ParseIP(host)
	if ip != nil ***REMOVED***
		return ip.IsLoopback(), nil
	***REMOVED***

	// Host is not an ip, perform lookup
	addrs, err := lookupHostFn(host)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	for _, addr := range addrs ***REMOVED***
		if !net.ParseIP(addr).IsLoopback() ***REMOVED***
			return false, nil
		***REMOVED***
	***REMOVED***

	return true, nil
***REMOVED***

func localHTTPCredProvider(cfg aws.Config, handlers request.Handlers, u string) credentials.Provider ***REMOVED***
	var errMsg string

	parsed, err := url.Parse(u)
	if err != nil ***REMOVED***
		errMsg = fmt.Sprintf("invalid URL, %v", err)
	***REMOVED*** else ***REMOVED***
		host := aws.URLHostname(parsed)
		if len(host) == 0 ***REMOVED***
			errMsg = "unable to parse host from local HTTP cred provider URL"
		***REMOVED*** else if isLoopback, loopbackErr := isLoopbackHost(host); loopbackErr != nil ***REMOVED***
			errMsg = fmt.Sprintf("failed to resolve host %q, %v", host, loopbackErr)
		***REMOVED*** else if !isLoopback ***REMOVED***
			errMsg = fmt.Sprintf("invalid endpoint host, %q, only loopback hosts are allowed.", host)
		***REMOVED***
	***REMOVED***

	if len(errMsg) > 0 ***REMOVED***
		if cfg.Logger != nil ***REMOVED***
			cfg.Logger.Log("Ignoring, HTTP credential provider", errMsg, err)
		***REMOVED***
		return credentials.ErrorProvider***REMOVED***
			Err:          awserr.New("CredentialsEndpointError", errMsg, err),
			ProviderName: endpointcreds.ProviderName,
		***REMOVED***
	***REMOVED***

	return httpCredProvider(cfg, handlers, u)
***REMOVED***

func httpCredProvider(cfg aws.Config, handlers request.Handlers, u string) credentials.Provider ***REMOVED***
	return endpointcreds.NewProviderClient(cfg, handlers, u,
		func(p *endpointcreds.Provider) ***REMOVED***
			p.ExpiryWindow = 5 * time.Minute
		***REMOVED***,
	)
***REMOVED***

func ec2RoleProvider(cfg aws.Config, handlers request.Handlers) credentials.Provider ***REMOVED***
	resolver := cfg.EndpointResolver
	if resolver == nil ***REMOVED***
		resolver = endpoints.DefaultResolver()
	***REMOVED***

	e, _ := resolver.EndpointFor(endpoints.Ec2metadataServiceID, "")
	return &ec2rolecreds.EC2RoleProvider***REMOVED***
		Client:       ec2metadata.NewClient(cfg, handlers, e.URL, e.SigningRegion),
		ExpiryWindow: 5 * time.Minute,
	***REMOVED***
***REMOVED***
