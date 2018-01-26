package registry

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/docker/distribution/registry/client/auth"
	"github.com/docker/distribution/registry/client/auth/challenge"
	"github.com/docker/distribution/registry/client/transport"
	"github.com/docker/docker/api/types"
	registrytypes "github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/errdefs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	// AuthClientID is used the ClientID used for the token server
	AuthClientID = "docker"
)

// loginV1 tries to register/login to the v1 registry server.
func loginV1(authConfig *types.AuthConfig, apiEndpoint APIEndpoint, userAgent string) (string, string, error) ***REMOVED***
	registryEndpoint := apiEndpoint.ToV1Endpoint(userAgent, nil)
	serverAddress := registryEndpoint.String()

	logrus.Debugf("attempting v1 login to registry endpoint %s", serverAddress)

	if serverAddress == "" ***REMOVED***
		return "", "", errdefs.System(errors.New("server Error: Server Address not set"))
	***REMOVED***

	req, err := http.NewRequest("GET", serverAddress+"users/", nil)
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***
	req.SetBasicAuth(authConfig.Username, authConfig.Password)
	resp, err := registryEndpoint.client.Do(req)
	if err != nil ***REMOVED***
		// fallback when request could not be completed
		return "", "", fallbackError***REMOVED***
			err: err,
		***REMOVED***
	***REMOVED***
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil ***REMOVED***
		return "", "", errdefs.System(err)
	***REMOVED***

	switch resp.StatusCode ***REMOVED***
	case http.StatusOK:
		return "Login Succeeded", "", nil
	case http.StatusUnauthorized:
		return "", "", errdefs.Unauthorized(errors.New("Wrong login/password, please try again"))
	case http.StatusForbidden:
		// *TODO: Use registry configuration to determine what this says, if anything?
		return "", "", errdefs.Forbidden(errors.Errorf("Login: Account is not active. Please see the documentation of the registry %s for instructions how to activate it.", serverAddress))
	case http.StatusInternalServerError:
		logrus.Errorf("%s returned status code %d. Response Body :\n%s", req.URL.String(), resp.StatusCode, body)
		return "", "", errdefs.System(errors.New("Internal Server Error"))
	***REMOVED***
	return "", "", errdefs.System(errors.Errorf("Login: %s (Code: %d; Headers: %s)", body,
		resp.StatusCode, resp.Header))
***REMOVED***

type loginCredentialStore struct ***REMOVED***
	authConfig *types.AuthConfig
***REMOVED***

func (lcs loginCredentialStore) Basic(*url.URL) (string, string) ***REMOVED***
	return lcs.authConfig.Username, lcs.authConfig.Password
***REMOVED***

func (lcs loginCredentialStore) RefreshToken(*url.URL, string) string ***REMOVED***
	return lcs.authConfig.IdentityToken
***REMOVED***

func (lcs loginCredentialStore) SetRefreshToken(u *url.URL, service, token string) ***REMOVED***
	lcs.authConfig.IdentityToken = token
***REMOVED***

type staticCredentialStore struct ***REMOVED***
	auth *types.AuthConfig
***REMOVED***

// NewStaticCredentialStore returns a credential store
// which always returns the same credential values.
func NewStaticCredentialStore(auth *types.AuthConfig) auth.CredentialStore ***REMOVED***
	return staticCredentialStore***REMOVED***
		auth: auth,
	***REMOVED***
***REMOVED***

func (scs staticCredentialStore) Basic(*url.URL) (string, string) ***REMOVED***
	if scs.auth == nil ***REMOVED***
		return "", ""
	***REMOVED***
	return scs.auth.Username, scs.auth.Password
***REMOVED***

func (scs staticCredentialStore) RefreshToken(*url.URL, string) string ***REMOVED***
	if scs.auth == nil ***REMOVED***
		return ""
	***REMOVED***
	return scs.auth.IdentityToken
***REMOVED***

func (scs staticCredentialStore) SetRefreshToken(*url.URL, string, string) ***REMOVED***
***REMOVED***

type fallbackError struct ***REMOVED***
	err error
***REMOVED***

func (err fallbackError) Error() string ***REMOVED***
	return err.err.Error()
***REMOVED***

// loginV2 tries to login to the v2 registry server. The given registry
// endpoint will be pinged to get authorization challenges. These challenges
// will be used to authenticate against the registry to validate credentials.
func loginV2(authConfig *types.AuthConfig, endpoint APIEndpoint, userAgent string) (string, string, error) ***REMOVED***
	logrus.Debugf("attempting v2 login to registry endpoint %s", strings.TrimRight(endpoint.URL.String(), "/")+"/v2/")

	modifiers := Headers(userAgent, nil)
	authTransport := transport.NewTransport(NewTransport(endpoint.TLSConfig), modifiers...)

	credentialAuthConfig := *authConfig
	creds := loginCredentialStore***REMOVED***
		authConfig: &credentialAuthConfig,
	***REMOVED***

	loginClient, foundV2, err := v2AuthHTTPClient(endpoint.URL, authTransport, modifiers, creds, nil)
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***

	endpointStr := strings.TrimRight(endpoint.URL.String(), "/") + "/v2/"
	req, err := http.NewRequest("GET", endpointStr, nil)
	if err != nil ***REMOVED***
		if !foundV2 ***REMOVED***
			err = fallbackError***REMOVED***err: err***REMOVED***
		***REMOVED***
		return "", "", err
	***REMOVED***

	resp, err := loginClient.Do(req)
	if err != nil ***REMOVED***
		err = translateV2AuthError(err)
		if !foundV2 ***REMOVED***
			err = fallbackError***REMOVED***err: err***REMOVED***
		***REMOVED***

		return "", "", err
	***REMOVED***
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK ***REMOVED***
		return "Login Succeeded", credentialAuthConfig.IdentityToken, nil
	***REMOVED***

	// TODO(dmcgowan): Attempt to further interpret result, status code and error code string
	err = errors.Errorf("login attempt to %s failed with status: %d %s", endpointStr, resp.StatusCode, http.StatusText(resp.StatusCode))
	if !foundV2 ***REMOVED***
		err = fallbackError***REMOVED***err: err***REMOVED***
	***REMOVED***
	return "", "", err
***REMOVED***

func v2AuthHTTPClient(endpoint *url.URL, authTransport http.RoundTripper, modifiers []transport.RequestModifier, creds auth.CredentialStore, scopes []auth.Scope) (*http.Client, bool, error) ***REMOVED***
	challengeManager, foundV2, err := PingV2Registry(endpoint, authTransport)
	if err != nil ***REMOVED***
		if !foundV2 ***REMOVED***
			err = fallbackError***REMOVED***err: err***REMOVED***
		***REMOVED***
		return nil, foundV2, err
	***REMOVED***

	tokenHandlerOptions := auth.TokenHandlerOptions***REMOVED***
		Transport:     authTransport,
		Credentials:   creds,
		OfflineAccess: true,
		ClientID:      AuthClientID,
		Scopes:        scopes,
	***REMOVED***
	tokenHandler := auth.NewTokenHandlerWithOptions(tokenHandlerOptions)
	basicHandler := auth.NewBasicHandler(creds)
	modifiers = append(modifiers, auth.NewAuthorizer(challengeManager, tokenHandler, basicHandler))
	tr := transport.NewTransport(authTransport, modifiers...)

	return &http.Client***REMOVED***
		Transport: tr,
		Timeout:   15 * time.Second,
	***REMOVED***, foundV2, nil

***REMOVED***

// ConvertToHostname converts a registry url which has http|https prepended
// to just an hostname.
func ConvertToHostname(url string) string ***REMOVED***
	stripped := url
	if strings.HasPrefix(url, "http://") ***REMOVED***
		stripped = strings.TrimPrefix(url, "http://")
	***REMOVED*** else if strings.HasPrefix(url, "https://") ***REMOVED***
		stripped = strings.TrimPrefix(url, "https://")
	***REMOVED***

	nameParts := strings.SplitN(stripped, "/", 2)

	return nameParts[0]
***REMOVED***

// ResolveAuthConfig matches an auth configuration to a server address or a URL
func ResolveAuthConfig(authConfigs map[string]types.AuthConfig, index *registrytypes.IndexInfo) types.AuthConfig ***REMOVED***
	configKey := GetAuthConfigKey(index)
	// First try the happy case
	if c, found := authConfigs[configKey]; found || index.Official ***REMOVED***
		return c
	***REMOVED***

	// Maybe they have a legacy config file, we will iterate the keys converting
	// them to the new format and testing
	for registry, ac := range authConfigs ***REMOVED***
		if configKey == ConvertToHostname(registry) ***REMOVED***
			return ac
		***REMOVED***
	***REMOVED***

	// When all else fails, return an empty auth config
	return types.AuthConfig***REMOVED******REMOVED***
***REMOVED***

// PingResponseError is used when the response from a ping
// was received but invalid.
type PingResponseError struct ***REMOVED***
	Err error
***REMOVED***

func (err PingResponseError) Error() string ***REMOVED***
	return err.Err.Error()
***REMOVED***

// PingV2Registry attempts to ping a v2 registry and on success return a
// challenge manager for the supported authentication types and
// whether v2 was confirmed by the response. If a response is received but
// cannot be interpreted a PingResponseError will be returned.
// nolint: interfacer
func PingV2Registry(endpoint *url.URL, transport http.RoundTripper) (challenge.Manager, bool, error) ***REMOVED***
	var (
		foundV2   = false
		v2Version = auth.APIVersion***REMOVED***
			Type:    "registry",
			Version: "2.0",
		***REMOVED***
	)

	pingClient := &http.Client***REMOVED***
		Transport: transport,
		Timeout:   15 * time.Second,
	***REMOVED***
	endpointStr := strings.TrimRight(endpoint.String(), "/") + "/v2/"
	req, err := http.NewRequest("GET", endpointStr, nil)
	if err != nil ***REMOVED***
		return nil, false, err
	***REMOVED***
	resp, err := pingClient.Do(req)
	if err != nil ***REMOVED***
		return nil, false, err
	***REMOVED***
	defer resp.Body.Close()

	versions := auth.APIVersions(resp, DefaultRegistryVersionHeader)
	for _, pingVersion := range versions ***REMOVED***
		if pingVersion == v2Version ***REMOVED***
			// The version header indicates we're definitely
			// talking to a v2 registry. So don't allow future
			// fallbacks to the v1 protocol.

			foundV2 = true
			break
		***REMOVED***
	***REMOVED***

	challengeManager := challenge.NewSimpleManager()
	if err := challengeManager.AddResponse(resp); err != nil ***REMOVED***
		return nil, foundV2, PingResponseError***REMOVED***
			Err: err,
		***REMOVED***
	***REMOVED***

	return challengeManager, foundV2, nil
***REMOVED***
