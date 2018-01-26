package distribution

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/client"
	"github.com/docker/distribution/registry/client/auth"
	"github.com/docker/distribution/registry/client/transport"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/dockerversion"
	"github.com/docker/docker/registry"
	"github.com/docker/go-connections/sockets"
	"golang.org/x/net/context"
)

// ImageTypes represents the schema2 config types for images
var ImageTypes = []string***REMOVED***
	schema2.MediaTypeImageConfig,
	// Handle unexpected values from https://github.com/docker/distribution/issues/1621
	// (see also https://github.com/docker/docker/issues/22378,
	// https://github.com/docker/docker/issues/30083)
	"application/octet-stream",
	"application/json",
	"text/html",
	// Treat defaulted values as images, newer types cannot be implied
	"",
***REMOVED***

// PluginTypes represents the schema2 config types for plugins
var PluginTypes = []string***REMOVED***
	schema2.MediaTypePluginConfig,
***REMOVED***

var mediaTypeClasses map[string]string

func init() ***REMOVED***
	// initialize media type classes with all know types for
	// plugin
	mediaTypeClasses = map[string]string***REMOVED******REMOVED***
	for _, t := range ImageTypes ***REMOVED***
		mediaTypeClasses[t] = "image"
	***REMOVED***
	for _, t := range PluginTypes ***REMOVED***
		mediaTypeClasses[t] = "plugin"
	***REMOVED***
***REMOVED***

// NewV2Repository returns a repository (v2 only). It creates an HTTP transport
// providing timeout settings and authentication support, and also verifies the
// remote API version.
func NewV2Repository(ctx context.Context, repoInfo *registry.RepositoryInfo, endpoint registry.APIEndpoint, metaHeaders http.Header, authConfig *types.AuthConfig, actions ...string) (repo distribution.Repository, foundVersion bool, err error) ***REMOVED***
	repoName := repoInfo.Name.Name()
	// If endpoint does not support CanonicalName, use the RemoteName instead
	if endpoint.TrimHostname ***REMOVED***
		repoName = reference.Path(repoInfo.Name)
	***REMOVED***

	direct := &net.Dialer***REMOVED***
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	***REMOVED***

	// TODO(dmcgowan): Call close idle connections when complete, use keep alive
	base := &http.Transport***REMOVED***
		Proxy:               http.ProxyFromEnvironment,
		Dial:                direct.Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     endpoint.TLSConfig,
		// TODO(dmcgowan): Call close idle connections when complete and use keep alive
		DisableKeepAlives: true,
	***REMOVED***

	proxyDialer, err := sockets.DialerFromEnvironment(direct)
	if err == nil ***REMOVED***
		base.Dial = proxyDialer.Dial
	***REMOVED***

	modifiers := registry.Headers(dockerversion.DockerUserAgent(ctx), metaHeaders)
	authTransport := transport.NewTransport(base, modifiers...)

	challengeManager, foundVersion, err := registry.PingV2Registry(endpoint.URL, authTransport)
	if err != nil ***REMOVED***
		transportOK := false
		if responseErr, ok := err.(registry.PingResponseError); ok ***REMOVED***
			transportOK = true
			err = responseErr.Err
		***REMOVED***
		return nil, foundVersion, fallbackError***REMOVED***
			err:         err,
			confirmedV2: foundVersion,
			transportOK: transportOK,
		***REMOVED***
	***REMOVED***

	if authConfig.RegistryToken != "" ***REMOVED***
		passThruTokenHandler := &existingTokenHandler***REMOVED***token: authConfig.RegistryToken***REMOVED***
		modifiers = append(modifiers, auth.NewAuthorizer(challengeManager, passThruTokenHandler))
	***REMOVED*** else ***REMOVED***
		scope := auth.RepositoryScope***REMOVED***
			Repository: repoName,
			Actions:    actions,
			Class:      repoInfo.Class,
		***REMOVED***

		creds := registry.NewStaticCredentialStore(authConfig)
		tokenHandlerOptions := auth.TokenHandlerOptions***REMOVED***
			Transport:   authTransport,
			Credentials: creds,
			Scopes:      []auth.Scope***REMOVED***scope***REMOVED***,
			ClientID:    registry.AuthClientID,
		***REMOVED***
		tokenHandler := auth.NewTokenHandlerWithOptions(tokenHandlerOptions)
		basicHandler := auth.NewBasicHandler(creds)
		modifiers = append(modifiers, auth.NewAuthorizer(challengeManager, tokenHandler, basicHandler))
	***REMOVED***
	tr := transport.NewTransport(base, modifiers...)

	repoNameRef, err := reference.WithName(repoName)
	if err != nil ***REMOVED***
		return nil, foundVersion, fallbackError***REMOVED***
			err:         err,
			confirmedV2: foundVersion,
			transportOK: true,
		***REMOVED***
	***REMOVED***

	repo, err = client.NewRepository(ctx, repoNameRef, endpoint.URL.String(), tr)
	if err != nil ***REMOVED***
		err = fallbackError***REMOVED***
			err:         err,
			confirmedV2: foundVersion,
			transportOK: true,
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

type existingTokenHandler struct ***REMOVED***
	token string
***REMOVED***

func (th *existingTokenHandler) Scheme() string ***REMOVED***
	return "bearer"
***REMOVED***

func (th *existingTokenHandler) AuthorizeRequest(req *http.Request, params map[string]string) error ***REMOVED***
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", th.token))
	return nil
***REMOVED***
