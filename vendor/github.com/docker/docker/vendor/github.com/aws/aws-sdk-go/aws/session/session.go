package session

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/corehandlers"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
)

// A Session provides a central location to create service clients from and
// store configurations and request handlers for those services.
//
// Sessions are safe to create service clients concurrently, but it is not safe
// to mutate the Session concurrently.
//
// The Session satisfies the service client's client.ClientConfigProvider.
type Session struct ***REMOVED***
	Config   *aws.Config
	Handlers request.Handlers
***REMOVED***

// New creates a new instance of the handlers merging in the provided configs
// on top of the SDK's default configurations. Once the Session is created it
// can be mutated to modify the Config or Handlers. The Session is safe to be
// read concurrently, but it should not be written to concurrently.
//
// If the AWS_SDK_LOAD_CONFIG environment is set to a truthy value, the New
// method could now encounter an error when loading the configuration. When
// The environment variable is set, and an error occurs, New will return a
// session that will fail all requests reporting the error that occurred while
// loading the session. Use NewSession to get the error when creating the
// session.
//
// If the AWS_SDK_LOAD_CONFIG environment variable is set to a truthy value
// the shared config file (~/.aws/config) will also be loaded, in addition to
// the shared credentials file (~/.aws/credentials). Values set in both the
// shared config, and shared credentials will be taken from the shared
// credentials file.
//
// Deprecated: Use NewSession functions to create sessions instead. NewSession
// has the same functionality as New except an error can be returned when the
// func is called instead of waiting to receive an error until a request is made.
func New(cfgs ...*aws.Config) *Session ***REMOVED***
	// load initial config from environment
	envCfg := loadEnvConfig()

	if envCfg.EnableSharedConfig ***REMOVED***
		s, err := newSession(Options***REMOVED******REMOVED***, envCfg, cfgs...)
		if err != nil ***REMOVED***
			// Old session.New expected all errors to be discovered when
			// a request is made, and would report the errors then. This
			// needs to be replicated if an error occurs while creating
			// the session.
			msg := "failed to create session with AWS_SDK_LOAD_CONFIG enabled. " +
				"Use session.NewSession to handle errors occurring during session creation."

			// Session creation failed, need to report the error and prevent
			// any requests from succeeding.
			s = &Session***REMOVED***Config: defaults.Config()***REMOVED***
			s.Config.MergeIn(cfgs...)
			s.Config.Logger.Log("ERROR:", msg, "Error:", err)
			s.Handlers.Validate.PushBack(func(r *request.Request) ***REMOVED***
				r.Error = err
			***REMOVED***)
		***REMOVED***
		return s
	***REMOVED***

	return deprecatedNewSession(cfgs...)
***REMOVED***

// NewSession returns a new Session created from SDK defaults, config files,
// environment, and user provided config files. Once the Session is created
// it can be mutated to modify the Config or Handlers. The Session is safe to
// be read concurrently, but it should not be written to concurrently.
//
// If the AWS_SDK_LOAD_CONFIG environment variable is set to a truthy value
// the shared config file (~/.aws/config) will also be loaded in addition to
// the shared credentials file (~/.aws/credentials). Values set in both the
// shared config, and shared credentials will be taken from the shared
// credentials file. Enabling the Shared Config will also allow the Session
// to be built with retrieving credentials with AssumeRole set in the config.
//
// See the NewSessionWithOptions func for information on how to override or
// control through code how the Session will be created. Such as specifying the
// config profile, and controlling if shared config is enabled or not.
func NewSession(cfgs ...*aws.Config) (*Session, error) ***REMOVED***
	opts := Options***REMOVED******REMOVED***
	opts.Config.MergeIn(cfgs...)

	return NewSessionWithOptions(opts)
***REMOVED***

// SharedConfigState provides the ability to optionally override the state
// of the session's creation based on the shared config being enabled or
// disabled.
type SharedConfigState int

const (
	// SharedConfigStateFromEnv does not override any state of the
	// AWS_SDK_LOAD_CONFIG env var. It is the default value of the
	// SharedConfigState type.
	SharedConfigStateFromEnv SharedConfigState = iota

	// SharedConfigDisable overrides the AWS_SDK_LOAD_CONFIG env var value
	// and disables the shared config functionality.
	SharedConfigDisable

	// SharedConfigEnable overrides the AWS_SDK_LOAD_CONFIG env var value
	// and enables the shared config functionality.
	SharedConfigEnable
)

// Options provides the means to control how a Session is created and what
// configuration values will be loaded.
//
type Options struct ***REMOVED***
	// Provides config values for the SDK to use when creating service clients
	// and making API requests to services. Any value set in with this field
	// will override the associated value provided by the SDK defaults,
	// environment or config files where relevant.
	//
	// If not set, configuration values from from SDK defaults, environment,
	// config will be used.
	Config aws.Config

	// Overrides the config profile the Session should be created from. If not
	// set the value of the environment variable will be loaded (AWS_PROFILE,
	// or AWS_DEFAULT_PROFILE if the Shared Config is enabled).
	//
	// If not set and environment variables are not set the "default"
	// (DefaultSharedConfigProfile) will be used as the profile to load the
	// session config from.
	Profile string

	// Instructs how the Session will be created based on the AWS_SDK_LOAD_CONFIG
	// environment variable. By default a Session will be created using the
	// value provided by the AWS_SDK_LOAD_CONFIG environment variable.
	//
	// Setting this value to SharedConfigEnable or SharedConfigDisable
	// will allow you to override the AWS_SDK_LOAD_CONFIG environment variable
	// and enable or disable the shared config functionality.
	SharedConfigState SharedConfigState

	// Ordered list of files the session will load configuration from.
	// It will override environment variable AWS_SHARED_CREDENTIALS_FILE, AWS_CONFIG_FILE.
	SharedConfigFiles []string

	// When the SDK's shared config is configured to assume a role with MFA
	// this option is required in order to provide the mechanism that will
	// retrieve the MFA token. There is no default value for this field. If
	// it is not set an error will be returned when creating the session.
	//
	// This token provider will be called when ever the assumed role's
	// credentials need to be refreshed. Within the context of service clients
	// all sharing the same session the SDK will ensure calls to the token
	// provider are atomic. When sharing a token provider across multiple
	// sessions additional synchronization logic is needed to ensure the
	// token providers do not introduce race conditions. It is recommend to
	// share the session where possible.
	//
	// stscreds.StdinTokenProvider is a basic implementation that will prompt
	// from stdin for the MFA token code.
	//
	// This field is only used if the shared configuration is enabled, and
	// the config enables assume role wit MFA via the mfa_serial field.
	AssumeRoleTokenProvider func() (string, error)

	// Reader for a custom Credentials Authority (CA) bundle in PEM format that
	// the SDK will use instead of the default system's root CA bundle. Use this
	// only if you want to replace the CA bundle the SDK uses for TLS requests.
	//
	// Enabling this option will attempt to merge the Transport into the SDK's HTTP
	// client. If the client's Transport is not a http.Transport an error will be
	// returned. If the Transport's TLS config is set this option will cause the SDK
	// to overwrite the Transport's TLS config's  RootCAs value. If the CA
	// bundle reader contains multiple certificates all of them will be loaded.
	//
	// The Session option CustomCABundle is also available when creating sessions
	// to also enable this feature. CustomCABundle session option field has priority
	// over the AWS_CA_BUNDLE environment variable, and will be used if both are set.
	CustomCABundle io.Reader
***REMOVED***

// NewSessionWithOptions returns a new Session created from SDK defaults, config files,
// environment, and user provided config files. This func uses the Options
// values to configure how the Session is created.
//
// If the AWS_SDK_LOAD_CONFIG environment variable is set to a truthy value
// the shared config file (~/.aws/config) will also be loaded in addition to
// the shared credentials file (~/.aws/credentials). Values set in both the
// shared config, and shared credentials will be taken from the shared
// credentials file. Enabling the Shared Config will also allow the Session
// to be built with retrieving credentials with AssumeRole set in the config.
//
//     // Equivalent to session.New
//     sess := session.Must(session.NewSessionWithOptions(session.Options***REMOVED******REMOVED***))
//
//     // Specify profile to load for the session's config
//     sess := session.Must(session.NewSessionWithOptions(session.Options***REMOVED***
//          Profile: "profile_name",
// ***REMOVED***))
//
//     // Specify profile for config and region for requests
//     sess := session.Must(session.NewSessionWithOptions(session.Options***REMOVED***
//          Config: aws.Config***REMOVED***Region: aws.String("us-east-1")***REMOVED***,
//          Profile: "profile_name",
// ***REMOVED***))
//
//     // Force enable Shared Config support
//     sess := session.Must(session.NewSessionWithOptions(session.Options***REMOVED***
//         SharedConfigState: session.SharedConfigEnable,
// ***REMOVED***))
func NewSessionWithOptions(opts Options) (*Session, error) ***REMOVED***
	var envCfg envConfig
	if opts.SharedConfigState == SharedConfigEnable ***REMOVED***
		envCfg = loadSharedEnvConfig()
	***REMOVED*** else ***REMOVED***
		envCfg = loadEnvConfig()
	***REMOVED***

	if len(opts.Profile) > 0 ***REMOVED***
		envCfg.Profile = opts.Profile
	***REMOVED***

	switch opts.SharedConfigState ***REMOVED***
	case SharedConfigDisable:
		envCfg.EnableSharedConfig = false
	case SharedConfigEnable:
		envCfg.EnableSharedConfig = true
	***REMOVED***

	if len(envCfg.SharedCredentialsFile) == 0 ***REMOVED***
		envCfg.SharedCredentialsFile = defaults.SharedCredentialsFilename()
	***REMOVED***
	if len(envCfg.SharedConfigFile) == 0 ***REMOVED***
		envCfg.SharedConfigFile = defaults.SharedConfigFilename()
	***REMOVED***

	// Only use AWS_CA_BUNDLE if session option is not provided.
	if len(envCfg.CustomCABundle) != 0 && opts.CustomCABundle == nil ***REMOVED***
		f, err := os.Open(envCfg.CustomCABundle)
		if err != nil ***REMOVED***
			return nil, awserr.New("LoadCustomCABundleError",
				"failed to open custom CA bundle PEM file", err)
		***REMOVED***
		defer f.Close()
		opts.CustomCABundle = f
	***REMOVED***

	return newSession(opts, envCfg, &opts.Config)
***REMOVED***

// Must is a helper function to ensure the Session is valid and there was no
// error when calling a NewSession function.
//
// This helper is intended to be used in variable initialization to load the
// Session and configuration at startup. Such as:
//
//     var sess = session.Must(session.NewSession())
func Must(sess *Session, err error) *Session ***REMOVED***
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	return sess
***REMOVED***

func deprecatedNewSession(cfgs ...*aws.Config) *Session ***REMOVED***
	cfg := defaults.Config()
	handlers := defaults.Handlers()

	// Apply the passed in configs so the configuration can be applied to the
	// default credential chain
	cfg.MergeIn(cfgs...)
	if cfg.EndpointResolver == nil ***REMOVED***
		// An endpoint resolver is required for a session to be able to provide
		// endpoints for service client configurations.
		cfg.EndpointResolver = endpoints.DefaultResolver()
	***REMOVED***
	cfg.Credentials = defaults.CredChain(cfg, handlers)

	// Reapply any passed in configs to override credentials if set
	cfg.MergeIn(cfgs...)

	s := &Session***REMOVED***
		Config:   cfg,
		Handlers: handlers,
	***REMOVED***

	initHandlers(s)

	return s
***REMOVED***

func newSession(opts Options, envCfg envConfig, cfgs ...*aws.Config) (*Session, error) ***REMOVED***
	cfg := defaults.Config()
	handlers := defaults.Handlers()

	// Get a merged version of the user provided config to determine if
	// credentials were.
	userCfg := &aws.Config***REMOVED******REMOVED***
	userCfg.MergeIn(cfgs...)

	// Ordered config files will be loaded in with later files overwriting
	// previous config file values.
	var cfgFiles []string
	if opts.SharedConfigFiles != nil ***REMOVED***
		cfgFiles = opts.SharedConfigFiles
	***REMOVED*** else ***REMOVED***
		cfgFiles = []string***REMOVED***envCfg.SharedConfigFile, envCfg.SharedCredentialsFile***REMOVED***
		if !envCfg.EnableSharedConfig ***REMOVED***
			// The shared config file (~/.aws/config) is only loaded if instructed
			// to load via the envConfig.EnableSharedConfig (AWS_SDK_LOAD_CONFIG).
			cfgFiles = cfgFiles[1:]
		***REMOVED***
	***REMOVED***

	// Load additional config from file(s)
	sharedCfg, err := loadSharedConfig(envCfg.Profile, cfgFiles)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := mergeConfigSrcs(cfg, userCfg, envCfg, sharedCfg, handlers, opts); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	s := &Session***REMOVED***
		Config:   cfg,
		Handlers: handlers,
	***REMOVED***

	initHandlers(s)

	// Setup HTTP client with custom cert bundle if enabled
	if opts.CustomCABundle != nil ***REMOVED***
		if err := loadCustomCABundle(s, opts.CustomCABundle); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return s, nil
***REMOVED***

func loadCustomCABundle(s *Session, bundle io.Reader) error ***REMOVED***
	var t *http.Transport
	switch v := s.Config.HTTPClient.Transport.(type) ***REMOVED***
	case *http.Transport:
		t = v
	default:
		if s.Config.HTTPClient.Transport != nil ***REMOVED***
			return awserr.New("LoadCustomCABundleError",
				"unable to load custom CA bundle, HTTPClient's transport unsupported type", nil)
		***REMOVED***
	***REMOVED***
	if t == nil ***REMOVED***
		t = &http.Transport***REMOVED******REMOVED***
	***REMOVED***

	p, err := loadCertPool(bundle)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if t.TLSClientConfig == nil ***REMOVED***
		t.TLSClientConfig = &tls.Config***REMOVED******REMOVED***
	***REMOVED***
	t.TLSClientConfig.RootCAs = p

	s.Config.HTTPClient.Transport = t

	return nil
***REMOVED***

func loadCertPool(r io.Reader) (*x509.CertPool, error) ***REMOVED***
	b, err := ioutil.ReadAll(r)
	if err != nil ***REMOVED***
		return nil, awserr.New("LoadCustomCABundleError",
			"failed to read custom CA bundle PEM file", err)
	***REMOVED***

	p := x509.NewCertPool()
	if !p.AppendCertsFromPEM(b) ***REMOVED***
		return nil, awserr.New("LoadCustomCABundleError",
			"failed to load custom CA bundle PEM file", err)
	***REMOVED***

	return p, nil
***REMOVED***

func mergeConfigSrcs(cfg, userCfg *aws.Config, envCfg envConfig, sharedCfg sharedConfig, handlers request.Handlers, sessOpts Options) error ***REMOVED***
	// Merge in user provided configuration
	cfg.MergeIn(userCfg)

	// Region if not already set by user
	if len(aws.StringValue(cfg.Region)) == 0 ***REMOVED***
		if len(envCfg.Region) > 0 ***REMOVED***
			cfg.WithRegion(envCfg.Region)
		***REMOVED*** else if envCfg.EnableSharedConfig && len(sharedCfg.Region) > 0 ***REMOVED***
			cfg.WithRegion(sharedCfg.Region)
		***REMOVED***
	***REMOVED***

	// Configure credentials if not already set
	if cfg.Credentials == credentials.AnonymousCredentials && userCfg.Credentials == nil ***REMOVED***
		if len(envCfg.Creds.AccessKeyID) > 0 ***REMOVED***
			cfg.Credentials = credentials.NewStaticCredentialsFromCreds(
				envCfg.Creds,
			)
		***REMOVED*** else if envCfg.EnableSharedConfig && len(sharedCfg.AssumeRole.RoleARN) > 0 && sharedCfg.AssumeRoleSource != nil ***REMOVED***
			cfgCp := *cfg
			cfgCp.Credentials = credentials.NewStaticCredentialsFromCreds(
				sharedCfg.AssumeRoleSource.Creds,
			)
			if len(sharedCfg.AssumeRole.MFASerial) > 0 && sessOpts.AssumeRoleTokenProvider == nil ***REMOVED***
				// AssumeRole Token provider is required if doing Assume Role
				// with MFA.
				return AssumeRoleTokenProviderNotSetError***REMOVED******REMOVED***
			***REMOVED***
			cfg.Credentials = stscreds.NewCredentials(
				&Session***REMOVED***
					Config:   &cfgCp,
					Handlers: handlers.Copy(),
				***REMOVED***,
				sharedCfg.AssumeRole.RoleARN,
				func(opt *stscreds.AssumeRoleProvider) ***REMOVED***
					opt.RoleSessionName = sharedCfg.AssumeRole.RoleSessionName

					// Assume role with external ID
					if len(sharedCfg.AssumeRole.ExternalID) > 0 ***REMOVED***
						opt.ExternalID = aws.String(sharedCfg.AssumeRole.ExternalID)
					***REMOVED***

					// Assume role with MFA
					if len(sharedCfg.AssumeRole.MFASerial) > 0 ***REMOVED***
						opt.SerialNumber = aws.String(sharedCfg.AssumeRole.MFASerial)
						opt.TokenProvider = sessOpts.AssumeRoleTokenProvider
					***REMOVED***
				***REMOVED***,
			)
		***REMOVED*** else if len(sharedCfg.Creds.AccessKeyID) > 0 ***REMOVED***
			cfg.Credentials = credentials.NewStaticCredentialsFromCreds(
				sharedCfg.Creds,
			)
		***REMOVED*** else ***REMOVED***
			// Fallback to default credentials provider, include mock errors
			// for the credential chain so user can identify why credentials
			// failed to be retrieved.
			cfg.Credentials = credentials.NewCredentials(&credentials.ChainProvider***REMOVED***
				VerboseErrors: aws.BoolValue(cfg.CredentialsChainVerboseErrors),
				Providers: []credentials.Provider***REMOVED***
					&credProviderError***REMOVED***Err: awserr.New("EnvAccessKeyNotFound", "failed to find credentials in the environment.", nil)***REMOVED***,
					&credProviderError***REMOVED***Err: awserr.New("SharedCredsLoad", fmt.Sprintf("failed to load profile, %s.", envCfg.Profile), nil)***REMOVED***,
					defaults.RemoteCredProvider(*cfg, handlers),
				***REMOVED***,
			***REMOVED***)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// AssumeRoleTokenProviderNotSetError is an error returned when creating a session when the
// MFAToken option is not set when shared config is configured load assume a
// role with an MFA token.
type AssumeRoleTokenProviderNotSetError struct***REMOVED******REMOVED***

// Code is the short id of the error.
func (e AssumeRoleTokenProviderNotSetError) Code() string ***REMOVED***
	return "AssumeRoleTokenProviderNotSetError"
***REMOVED***

// Message is the description of the error
func (e AssumeRoleTokenProviderNotSetError) Message() string ***REMOVED***
	return fmt.Sprintf("assume role with MFA enabled, but AssumeRoleTokenProvider session option not set.")
***REMOVED***

// OrigErr is the underlying error that caused the failure.
func (e AssumeRoleTokenProviderNotSetError) OrigErr() error ***REMOVED***
	return nil
***REMOVED***

// Error satisfies the error interface.
func (e AssumeRoleTokenProviderNotSetError) Error() string ***REMOVED***
	return awserr.SprintError(e.Code(), e.Message(), "", nil)
***REMOVED***

type credProviderError struct ***REMOVED***
	Err error
***REMOVED***

var emptyCreds = credentials.Value***REMOVED******REMOVED***

func (c credProviderError) Retrieve() (credentials.Value, error) ***REMOVED***
	return credentials.Value***REMOVED******REMOVED***, c.Err
***REMOVED***
func (c credProviderError) IsExpired() bool ***REMOVED***
	return true
***REMOVED***

func initHandlers(s *Session) ***REMOVED***
	// Add the Validate parameter handler if it is not disabled.
	s.Handlers.Validate.Remove(corehandlers.ValidateParametersHandler)
	if !aws.BoolValue(s.Config.DisableParamValidation) ***REMOVED***
		s.Handlers.Validate.PushBackNamed(corehandlers.ValidateParametersHandler)
	***REMOVED***
***REMOVED***

// Copy creates and returns a copy of the current Session, coping the config
// and handlers. If any additional configs are provided they will be merged
// on top of the Session's copied config.
//
//     // Create a copy of the current Session, configured for the us-west-2 region.
//     sess.Copy(&aws.Config***REMOVED***Region: aws.String("us-west-2")***REMOVED***)
func (s *Session) Copy(cfgs ...*aws.Config) *Session ***REMOVED***
	newSession := &Session***REMOVED***
		Config:   s.Config.Copy(cfgs...),
		Handlers: s.Handlers.Copy(),
	***REMOVED***

	initHandlers(newSession)

	return newSession
***REMOVED***

// ClientConfig satisfies the client.ConfigProvider interface and is used to
// configure the service client instances. Passing the Session to the service
// client's constructor (New) will use this method to configure the client.
func (s *Session) ClientConfig(serviceName string, cfgs ...*aws.Config) client.Config ***REMOVED***
	// Backwards compatibility, the error will be eaten if user calls ClientConfig
	// directly. All SDK services will use ClientconfigWithError.
	cfg, _ := s.clientConfigWithErr(serviceName, cfgs...)

	return cfg
***REMOVED***

func (s *Session) clientConfigWithErr(serviceName string, cfgs ...*aws.Config) (client.Config, error) ***REMOVED***
	s = s.Copy(cfgs...)

	var resolved endpoints.ResolvedEndpoint
	var err error

	region := aws.StringValue(s.Config.Region)

	if endpoint := aws.StringValue(s.Config.Endpoint); len(endpoint) != 0 ***REMOVED***
		resolved.URL = endpoints.AddScheme(endpoint, aws.BoolValue(s.Config.DisableSSL))
		resolved.SigningRegion = region
	***REMOVED*** else ***REMOVED***
		resolved, err = s.Config.EndpointResolver.EndpointFor(
			serviceName, region,
			func(opt *endpoints.Options) ***REMOVED***
				opt.DisableSSL = aws.BoolValue(s.Config.DisableSSL)
				opt.UseDualStack = aws.BoolValue(s.Config.UseDualStack)

				// Support the condition where the service is modeled but its
				// endpoint metadata is not available.
				opt.ResolveUnknownService = true
			***REMOVED***,
		)
	***REMOVED***

	return client.Config***REMOVED***
		Config:        s.Config,
		Handlers:      s.Handlers,
		Endpoint:      resolved.URL,
		SigningRegion: resolved.SigningRegion,
		SigningName:   resolved.SigningName,
	***REMOVED***, err
***REMOVED***

// ClientConfigNoResolveEndpoint is the same as ClientConfig with the exception
// that the EndpointResolver will not be used to resolve the endpoint. The only
// endpoint set must come from the aws.Config.Endpoint field.
func (s *Session) ClientConfigNoResolveEndpoint(cfgs ...*aws.Config) client.Config ***REMOVED***
	s = s.Copy(cfgs...)

	var resolved endpoints.ResolvedEndpoint

	region := aws.StringValue(s.Config.Region)

	if ep := aws.StringValue(s.Config.Endpoint); len(ep) > 0 ***REMOVED***
		resolved.URL = endpoints.AddScheme(ep, aws.BoolValue(s.Config.DisableSSL))
		resolved.SigningRegion = region
	***REMOVED***

	return client.Config***REMOVED***
		Config:        s.Config,
		Handlers:      s.Handlers,
		Endpoint:      resolved.URL,
		SigningRegion: resolved.SigningRegion,
		SigningName:   resolved.SigningName,
	***REMOVED***
***REMOVED***
