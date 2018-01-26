package registry

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"golang.org/x/net/context"

	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/client/auth"
	"github.com/docker/docker/api/types"
	registrytypes "github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/errdefs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	// DefaultSearchLimit is the default value for maximum number of returned search results.
	DefaultSearchLimit = 25
)

// Service is the interface defining what a registry service should implement.
type Service interface ***REMOVED***
	Auth(ctx context.Context, authConfig *types.AuthConfig, userAgent string) (status, token string, err error)
	LookupPullEndpoints(hostname string) (endpoints []APIEndpoint, err error)
	LookupPushEndpoints(hostname string) (endpoints []APIEndpoint, err error)
	ResolveRepository(name reference.Named) (*RepositoryInfo, error)
	Search(ctx context.Context, term string, limit int, authConfig *types.AuthConfig, userAgent string, headers map[string][]string) (*registrytypes.SearchResults, error)
	ServiceConfig() *registrytypes.ServiceConfig
	TLSConfig(hostname string) (*tls.Config, error)
	LoadAllowNondistributableArtifacts([]string) error
	LoadMirrors([]string) error
	LoadInsecureRegistries([]string) error
***REMOVED***

// DefaultService is a registry service. It tracks configuration data such as a list
// of mirrors.
type DefaultService struct ***REMOVED***
	config *serviceConfig
	mu     sync.Mutex
***REMOVED***

// NewService returns a new instance of DefaultService ready to be
// installed into an engine.
func NewService(options ServiceOptions) (*DefaultService, error) ***REMOVED***
	config, err := newServiceConfig(options)

	return &DefaultService***REMOVED***config: config***REMOVED***, err
***REMOVED***

// ServiceConfig returns the public registry service configuration.
func (s *DefaultService) ServiceConfig() *registrytypes.ServiceConfig ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()

	servConfig := registrytypes.ServiceConfig***REMOVED***
		AllowNondistributableArtifactsCIDRs:     make([]*(registrytypes.NetIPNet), 0),
		AllowNondistributableArtifactsHostnames: make([]string, 0),
		InsecureRegistryCIDRs:                   make([]*(registrytypes.NetIPNet), 0),
		IndexConfigs:                            make(map[string]*(registrytypes.IndexInfo)),
		Mirrors:                                 make([]string, 0),
	***REMOVED***

	// construct a new ServiceConfig which will not retrieve s.Config directly,
	// and look up items in s.config with mu locked
	servConfig.AllowNondistributableArtifactsCIDRs = append(servConfig.AllowNondistributableArtifactsCIDRs, s.config.ServiceConfig.AllowNondistributableArtifactsCIDRs...)
	servConfig.AllowNondistributableArtifactsHostnames = append(servConfig.AllowNondistributableArtifactsHostnames, s.config.ServiceConfig.AllowNondistributableArtifactsHostnames...)
	servConfig.InsecureRegistryCIDRs = append(servConfig.InsecureRegistryCIDRs, s.config.ServiceConfig.InsecureRegistryCIDRs...)

	for key, value := range s.config.ServiceConfig.IndexConfigs ***REMOVED***
		servConfig.IndexConfigs[key] = value
	***REMOVED***

	servConfig.Mirrors = append(servConfig.Mirrors, s.config.ServiceConfig.Mirrors...)

	return &servConfig
***REMOVED***

// LoadAllowNondistributableArtifacts loads allow-nondistributable-artifacts registries for Service.
func (s *DefaultService) LoadAllowNondistributableArtifacts(registries []string) error ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.config.LoadAllowNondistributableArtifacts(registries)
***REMOVED***

// LoadMirrors loads registry mirrors for Service
func (s *DefaultService) LoadMirrors(mirrors []string) error ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.config.LoadMirrors(mirrors)
***REMOVED***

// LoadInsecureRegistries loads insecure registries for Service
func (s *DefaultService) LoadInsecureRegistries(registries []string) error ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.config.LoadInsecureRegistries(registries)
***REMOVED***

// Auth contacts the public registry with the provided credentials,
// and returns OK if authentication was successful.
// It can be used to verify the validity of a client's credentials.
func (s *DefaultService) Auth(ctx context.Context, authConfig *types.AuthConfig, userAgent string) (status, token string, err error) ***REMOVED***
	// TODO Use ctx when searching for repositories
	serverAddress := authConfig.ServerAddress
	if serverAddress == "" ***REMOVED***
		serverAddress = IndexServer
	***REMOVED***
	if !strings.HasPrefix(serverAddress, "https://") && !strings.HasPrefix(serverAddress, "http://") ***REMOVED***
		serverAddress = "https://" + serverAddress
	***REMOVED***
	u, err := url.Parse(serverAddress)
	if err != nil ***REMOVED***
		return "", "", errdefs.InvalidParameter(errors.Errorf("unable to parse server address: %v", err))
	***REMOVED***

	endpoints, err := s.LookupPushEndpoints(u.Host)
	if err != nil ***REMOVED***
		return "", "", errdefs.InvalidParameter(err)
	***REMOVED***

	for _, endpoint := range endpoints ***REMOVED***
		login := loginV2
		if endpoint.Version == APIVersion1 ***REMOVED***
			login = loginV1
		***REMOVED***

		status, token, err = login(authConfig, endpoint, userAgent)
		if err == nil ***REMOVED***
			return
		***REMOVED***
		if fErr, ok := err.(fallbackError); ok ***REMOVED***
			err = fErr.err
			logrus.Infof("Error logging in to %s endpoint, trying next endpoint: %v", endpoint.Version, err)
			continue
		***REMOVED***

		return "", "", err
	***REMOVED***

	return "", "", err
***REMOVED***

// splitReposSearchTerm breaks a search term into an index name and remote name
func splitReposSearchTerm(reposName string) (string, string) ***REMOVED***
	nameParts := strings.SplitN(reposName, "/", 2)
	var indexName, remoteName string
	if len(nameParts) == 1 || (!strings.Contains(nameParts[0], ".") &&
		!strings.Contains(nameParts[0], ":") && nameParts[0] != "localhost") ***REMOVED***
		// This is a Docker Index repos (ex: samalba/hipache or ubuntu)
		// 'docker.io'
		indexName = IndexName
		remoteName = reposName
	***REMOVED*** else ***REMOVED***
		indexName = nameParts[0]
		remoteName = nameParts[1]
	***REMOVED***
	return indexName, remoteName
***REMOVED***

// Search queries the public registry for images matching the specified
// search terms, and returns the results.
func (s *DefaultService) Search(ctx context.Context, term string, limit int, authConfig *types.AuthConfig, userAgent string, headers map[string][]string) (*registrytypes.SearchResults, error) ***REMOVED***
	// TODO Use ctx when searching for repositories
	if err := validateNoScheme(term); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	indexName, remoteName := splitReposSearchTerm(term)

	// Search is a long-running operation, just lock s.config to avoid block others.
	s.mu.Lock()
	index, err := newIndexInfo(s.config, indexName)
	s.mu.Unlock()

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// *TODO: Search multiple indexes.
	endpoint, err := NewV1Endpoint(index, userAgent, http.Header(headers))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var client *http.Client
	if authConfig != nil && authConfig.IdentityToken != "" && authConfig.Username != "" ***REMOVED***
		creds := NewStaticCredentialStore(authConfig)
		scopes := []auth.Scope***REMOVED***
			auth.RegistryScope***REMOVED***
				Name:    "catalog",
				Actions: []string***REMOVED***"search"***REMOVED***,
			***REMOVED***,
		***REMOVED***

		modifiers := Headers(userAgent, nil)
		v2Client, foundV2, err := v2AuthHTTPClient(endpoint.URL, endpoint.client.Transport, modifiers, creds, scopes)
		if err != nil ***REMOVED***
			if fErr, ok := err.(fallbackError); ok ***REMOVED***
				logrus.Errorf("Cannot use identity token for search, v2 auth not supported: %v", fErr.err)
			***REMOVED*** else ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED*** else if foundV2 ***REMOVED***
			// Copy non transport http client features
			v2Client.Timeout = endpoint.client.Timeout
			v2Client.CheckRedirect = endpoint.client.CheckRedirect
			v2Client.Jar = endpoint.client.Jar

			logrus.Debugf("using v2 client for search to %s", endpoint.URL)
			client = v2Client
		***REMOVED***
	***REMOVED***

	if client == nil ***REMOVED***
		client = endpoint.client
		if err := authorizeClient(client, authConfig, endpoint); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	r := newSession(client, authConfig, endpoint)

	if index.Official ***REMOVED***
		localName := remoteName
		if strings.HasPrefix(localName, "library/") ***REMOVED***
			// If pull "library/foo", it's stored locally under "foo"
			localName = strings.SplitN(localName, "/", 2)[1]
		***REMOVED***

		return r.SearchRepositories(localName, limit)
	***REMOVED***
	return r.SearchRepositories(remoteName, limit)
***REMOVED***

// ResolveRepository splits a repository name into its components
// and configuration of the associated registry.
func (s *DefaultService) ResolveRepository(name reference.Named) (*RepositoryInfo, error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	return newRepositoryInfo(s.config, name)
***REMOVED***

// APIEndpoint represents a remote API endpoint
type APIEndpoint struct ***REMOVED***
	Mirror                         bool
	URL                            *url.URL
	Version                        APIVersion
	AllowNondistributableArtifacts bool
	Official                       bool
	TrimHostname                   bool
	TLSConfig                      *tls.Config
***REMOVED***

// ToV1Endpoint returns a V1 API endpoint based on the APIEndpoint
func (e APIEndpoint) ToV1Endpoint(userAgent string, metaHeaders http.Header) *V1Endpoint ***REMOVED***
	return newV1Endpoint(*e.URL, e.TLSConfig, userAgent, metaHeaders)
***REMOVED***

// TLSConfig constructs a client TLS configuration based on server defaults
func (s *DefaultService) TLSConfig(hostname string) (*tls.Config, error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()

	return newTLSConfig(hostname, isSecureIndex(s.config, hostname))
***REMOVED***

// tlsConfig constructs a client TLS configuration based on server defaults
func (s *DefaultService) tlsConfig(hostname string) (*tls.Config, error) ***REMOVED***
	return newTLSConfig(hostname, isSecureIndex(s.config, hostname))
***REMOVED***

func (s *DefaultService) tlsConfigForMirror(mirrorURL *url.URL) (*tls.Config, error) ***REMOVED***
	return s.tlsConfig(mirrorURL.Host)
***REMOVED***

// LookupPullEndpoints creates a list of endpoints to try to pull from, in order of preference.
// It gives preference to v2 endpoints over v1, mirrors over the actual
// registry, and HTTPS over plain HTTP.
func (s *DefaultService) LookupPullEndpoints(hostname string) (endpoints []APIEndpoint, err error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.lookupEndpoints(hostname)
***REMOVED***

// LookupPushEndpoints creates a list of endpoints to try to push to, in order of preference.
// It gives preference to v2 endpoints over v1, and HTTPS over plain HTTP.
// Mirrors are not included.
func (s *DefaultService) LookupPushEndpoints(hostname string) (endpoints []APIEndpoint, err error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()

	allEndpoints, err := s.lookupEndpoints(hostname)
	if err == nil ***REMOVED***
		for _, endpoint := range allEndpoints ***REMOVED***
			if !endpoint.Mirror ***REMOVED***
				endpoints = append(endpoints, endpoint)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return endpoints, err
***REMOVED***

func (s *DefaultService) lookupEndpoints(hostname string) (endpoints []APIEndpoint, err error) ***REMOVED***
	endpoints, err = s.lookupV2Endpoints(hostname)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if s.config.V2Only ***REMOVED***
		return endpoints, nil
	***REMOVED***

	legacyEndpoints, err := s.lookupV1Endpoints(hostname)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	endpoints = append(endpoints, legacyEndpoints...)

	return endpoints, nil
***REMOVED***
