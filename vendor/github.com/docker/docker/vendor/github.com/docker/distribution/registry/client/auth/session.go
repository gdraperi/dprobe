package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/docker/distribution/registry/client"
	"github.com/docker/distribution/registry/client/auth/challenge"
	"github.com/docker/distribution/registry/client/transport"
	"github.com/sirupsen/logrus"
)

var (
	// ErrNoBasicAuthCredentials is returned if a request can't be authorized with
	// basic auth due to lack of credentials.
	ErrNoBasicAuthCredentials = errors.New("no basic auth credentials")

	// ErrNoToken is returned if a request is successful but the body does not
	// contain an authorization token.
	ErrNoToken = errors.New("authorization server did not include a token in the response")
)

const defaultClientID = "registry-client"

// AuthenticationHandler is an interface for authorizing a request from
// params from a "WWW-Authenicate" header for a single scheme.
type AuthenticationHandler interface ***REMOVED***
	// Scheme returns the scheme as expected from the "WWW-Authenicate" header.
	Scheme() string

	// AuthorizeRequest adds the authorization header to a request (if needed)
	// using the parameters from "WWW-Authenticate" method. The parameters
	// values depend on the scheme.
	AuthorizeRequest(req *http.Request, params map[string]string) error
***REMOVED***

// CredentialStore is an interface for getting credentials for
// a given URL
type CredentialStore interface ***REMOVED***
	// Basic returns basic auth for the given URL
	Basic(*url.URL) (string, string)

	// RefreshToken returns a refresh token for the
	// given URL and service
	RefreshToken(*url.URL, string) string

	// SetRefreshToken sets the refresh token if none
	// is provided for the given url and service
	SetRefreshToken(realm *url.URL, service, token string)
***REMOVED***

// NewAuthorizer creates an authorizer which can handle multiple authentication
// schemes. The handlers are tried in order, the higher priority authentication
// methods should be first. The challengeMap holds a list of challenges for
// a given root API endpoint (for example "https://registry-1.docker.io/v2/").
func NewAuthorizer(manager challenge.Manager, handlers ...AuthenticationHandler) transport.RequestModifier ***REMOVED***
	return &endpointAuthorizer***REMOVED***
		challenges: manager,
		handlers:   handlers,
	***REMOVED***
***REMOVED***

type endpointAuthorizer struct ***REMOVED***
	challenges challenge.Manager
	handlers   []AuthenticationHandler
	transport  http.RoundTripper
***REMOVED***

func (ea *endpointAuthorizer) ModifyRequest(req *http.Request) error ***REMOVED***
	pingPath := req.URL.Path
	if v2Root := strings.Index(req.URL.Path, "/v2/"); v2Root != -1 ***REMOVED***
		pingPath = pingPath[:v2Root+4]
	***REMOVED*** else if v1Root := strings.Index(req.URL.Path, "/v1/"); v1Root != -1 ***REMOVED***
		pingPath = pingPath[:v1Root] + "/v2/"
	***REMOVED*** else ***REMOVED***
		return nil
	***REMOVED***

	ping := url.URL***REMOVED***
		Host:   req.URL.Host,
		Scheme: req.URL.Scheme,
		Path:   pingPath,
	***REMOVED***

	challenges, err := ea.challenges.GetChallenges(ping)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if len(challenges) > 0 ***REMOVED***
		for _, handler := range ea.handlers ***REMOVED***
			for _, c := range challenges ***REMOVED***
				if c.Scheme != handler.Scheme() ***REMOVED***
					continue
				***REMOVED***
				if err := handler.AuthorizeRequest(req, c.Parameters); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// This is the minimum duration a token can last (in seconds).
// A token must not live less than 60 seconds because older versions
// of the Docker client didn't read their expiration from the token
// response and assumed 60 seconds.  So to remain compatible with
// those implementations, a token must live at least this long.
const minimumTokenLifetimeSeconds = 60

// Private interface for time used by this package to enable tests to provide their own implementation.
type clock interface ***REMOVED***
	Now() time.Time
***REMOVED***

type tokenHandler struct ***REMOVED***
	header    http.Header
	creds     CredentialStore
	transport http.RoundTripper
	clock     clock

	offlineAccess bool
	forceOAuth    bool
	clientID      string
	scopes        []Scope

	tokenLock       sync.Mutex
	tokenCache      string
	tokenExpiration time.Time
***REMOVED***

// Scope is a type which is serializable to a string
// using the allow scope grammar.
type Scope interface ***REMOVED***
	String() string
***REMOVED***

// RepositoryScope represents a token scope for access
// to a repository.
type RepositoryScope struct ***REMOVED***
	Repository string
	Class      string
	Actions    []string
***REMOVED***

// String returns the string representation of the repository
// using the scope grammar
func (rs RepositoryScope) String() string ***REMOVED***
	repoType := "repository"
	// Keep existing format for image class to maintain backwards compatibility
	// with authorization servers which do not support the expanded grammar.
	if rs.Class != "" && rs.Class != "image" ***REMOVED***
		repoType = fmt.Sprintf("%s(%s)", repoType, rs.Class)
	***REMOVED***
	return fmt.Sprintf("%s:%s:%s", repoType, rs.Repository, strings.Join(rs.Actions, ","))
***REMOVED***

// RegistryScope represents a token scope for access
// to resources in the registry.
type RegistryScope struct ***REMOVED***
	Name    string
	Actions []string
***REMOVED***

// String returns the string representation of the user
// using the scope grammar
func (rs RegistryScope) String() string ***REMOVED***
	return fmt.Sprintf("registry:%s:%s", rs.Name, strings.Join(rs.Actions, ","))
***REMOVED***

// TokenHandlerOptions is used to configure a new token handler
type TokenHandlerOptions struct ***REMOVED***
	Transport   http.RoundTripper
	Credentials CredentialStore

	OfflineAccess bool
	ForceOAuth    bool
	ClientID      string
	Scopes        []Scope
***REMOVED***

// An implementation of clock for providing real time data.
type realClock struct***REMOVED******REMOVED***

// Now implements clock
func (realClock) Now() time.Time ***REMOVED*** return time.Now() ***REMOVED***

// NewTokenHandler creates a new AuthenicationHandler which supports
// fetching tokens from a remote token server.
func NewTokenHandler(transport http.RoundTripper, creds CredentialStore, scope string, actions ...string) AuthenticationHandler ***REMOVED***
	// Create options...
	return NewTokenHandlerWithOptions(TokenHandlerOptions***REMOVED***
		Transport:   transport,
		Credentials: creds,
		Scopes: []Scope***REMOVED***
			RepositoryScope***REMOVED***
				Repository: scope,
				Actions:    actions,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***)
***REMOVED***

// NewTokenHandlerWithOptions creates a new token handler using the provided
// options structure.
func NewTokenHandlerWithOptions(options TokenHandlerOptions) AuthenticationHandler ***REMOVED***
	handler := &tokenHandler***REMOVED***
		transport:     options.Transport,
		creds:         options.Credentials,
		offlineAccess: options.OfflineAccess,
		forceOAuth:    options.ForceOAuth,
		clientID:      options.ClientID,
		scopes:        options.Scopes,
		clock:         realClock***REMOVED******REMOVED***,
	***REMOVED***

	return handler
***REMOVED***

func (th *tokenHandler) client() *http.Client ***REMOVED***
	return &http.Client***REMOVED***
		Transport: th.transport,
		Timeout:   15 * time.Second,
	***REMOVED***
***REMOVED***

func (th *tokenHandler) Scheme() string ***REMOVED***
	return "bearer"
***REMOVED***

func (th *tokenHandler) AuthorizeRequest(req *http.Request, params map[string]string) error ***REMOVED***
	var additionalScopes []string
	if fromParam := req.URL.Query().Get("from"); fromParam != "" ***REMOVED***
		additionalScopes = append(additionalScopes, RepositoryScope***REMOVED***
			Repository: fromParam,
			Actions:    []string***REMOVED***"pull"***REMOVED***,
		***REMOVED***.String())
	***REMOVED***

	token, err := th.getToken(params, additionalScopes...)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	return nil
***REMOVED***

func (th *tokenHandler) getToken(params map[string]string, additionalScopes ...string) (string, error) ***REMOVED***
	th.tokenLock.Lock()
	defer th.tokenLock.Unlock()
	scopes := make([]string, 0, len(th.scopes)+len(additionalScopes))
	for _, scope := range th.scopes ***REMOVED***
		scopes = append(scopes, scope.String())
	***REMOVED***
	var addedScopes bool
	for _, scope := range additionalScopes ***REMOVED***
		scopes = append(scopes, scope)
		addedScopes = true
	***REMOVED***

	now := th.clock.Now()
	if now.After(th.tokenExpiration) || addedScopes ***REMOVED***
		token, expiration, err := th.fetchToken(params, scopes)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***

		// do not update cache for added scope tokens
		if !addedScopes ***REMOVED***
			th.tokenCache = token
			th.tokenExpiration = expiration
		***REMOVED***

		return token, nil
	***REMOVED***

	return th.tokenCache, nil
***REMOVED***

type postTokenResponse struct ***REMOVED***
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"`
	IssuedAt     time.Time `json:"issued_at"`
	Scope        string    `json:"scope"`
***REMOVED***

func (th *tokenHandler) fetchTokenWithOAuth(realm *url.URL, refreshToken, service string, scopes []string) (token string, expiration time.Time, err error) ***REMOVED***
	form := url.Values***REMOVED******REMOVED***
	form.Set("scope", strings.Join(scopes, " "))
	form.Set("service", service)

	clientID := th.clientID
	if clientID == "" ***REMOVED***
		// Use default client, this is a required field
		clientID = defaultClientID
	***REMOVED***
	form.Set("client_id", clientID)

	if refreshToken != "" ***REMOVED***
		form.Set("grant_type", "refresh_token")
		form.Set("refresh_token", refreshToken)
	***REMOVED*** else if th.creds != nil ***REMOVED***
		form.Set("grant_type", "password")
		username, password := th.creds.Basic(realm)
		form.Set("username", username)
		form.Set("password", password)

		// attempt to get a refresh token
		form.Set("access_type", "offline")
	***REMOVED*** else ***REMOVED***
		// refuse to do oauth without a grant type
		return "", time.Time***REMOVED******REMOVED***, fmt.Errorf("no supported grant type")
	***REMOVED***

	resp, err := th.client().PostForm(realm.String(), form)
	if err != nil ***REMOVED***
		return "", time.Time***REMOVED******REMOVED***, err
	***REMOVED***
	defer resp.Body.Close()

	if !client.SuccessStatus(resp.StatusCode) ***REMOVED***
		err := client.HandleErrorResponse(resp)
		return "", time.Time***REMOVED******REMOVED***, err
	***REMOVED***

	decoder := json.NewDecoder(resp.Body)

	var tr postTokenResponse
	if err = decoder.Decode(&tr); err != nil ***REMOVED***
		return "", time.Time***REMOVED******REMOVED***, fmt.Errorf("unable to decode token response: %s", err)
	***REMOVED***

	if tr.RefreshToken != "" && tr.RefreshToken != refreshToken ***REMOVED***
		th.creds.SetRefreshToken(realm, service, tr.RefreshToken)
	***REMOVED***

	if tr.ExpiresIn < minimumTokenLifetimeSeconds ***REMOVED***
		// The default/minimum lifetime.
		tr.ExpiresIn = minimumTokenLifetimeSeconds
		logrus.Debugf("Increasing token expiration to: %d seconds", tr.ExpiresIn)
	***REMOVED***

	if tr.IssuedAt.IsZero() ***REMOVED***
		// issued_at is optional in the token response.
		tr.IssuedAt = th.clock.Now().UTC()
	***REMOVED***

	return tr.AccessToken, tr.IssuedAt.Add(time.Duration(tr.ExpiresIn) * time.Second), nil
***REMOVED***

type getTokenResponse struct ***REMOVED***
	Token        string    `json:"token"`
	AccessToken  string    `json:"access_token"`
	ExpiresIn    int       `json:"expires_in"`
	IssuedAt     time.Time `json:"issued_at"`
	RefreshToken string    `json:"refresh_token"`
***REMOVED***

func (th *tokenHandler) fetchTokenWithBasicAuth(realm *url.URL, service string, scopes []string) (token string, expiration time.Time, err error) ***REMOVED***

	req, err := http.NewRequest("GET", realm.String(), nil)
	if err != nil ***REMOVED***
		return "", time.Time***REMOVED******REMOVED***, err
	***REMOVED***

	reqParams := req.URL.Query()

	if service != "" ***REMOVED***
		reqParams.Add("service", service)
	***REMOVED***

	for _, scope := range scopes ***REMOVED***
		reqParams.Add("scope", scope)
	***REMOVED***

	if th.offlineAccess ***REMOVED***
		reqParams.Add("offline_token", "true")
		clientID := th.clientID
		if clientID == "" ***REMOVED***
			clientID = defaultClientID
		***REMOVED***
		reqParams.Add("client_id", clientID)
	***REMOVED***

	if th.creds != nil ***REMOVED***
		username, password := th.creds.Basic(realm)
		if username != "" && password != "" ***REMOVED***
			reqParams.Add("account", username)
			req.SetBasicAuth(username, password)
		***REMOVED***
	***REMOVED***

	req.URL.RawQuery = reqParams.Encode()

	resp, err := th.client().Do(req)
	if err != nil ***REMOVED***
		return "", time.Time***REMOVED******REMOVED***, err
	***REMOVED***
	defer resp.Body.Close()

	if !client.SuccessStatus(resp.StatusCode) ***REMOVED***
		err := client.HandleErrorResponse(resp)
		return "", time.Time***REMOVED******REMOVED***, err
	***REMOVED***

	decoder := json.NewDecoder(resp.Body)

	var tr getTokenResponse
	if err = decoder.Decode(&tr); err != nil ***REMOVED***
		return "", time.Time***REMOVED******REMOVED***, fmt.Errorf("unable to decode token response: %s", err)
	***REMOVED***

	if tr.RefreshToken != "" && th.creds != nil ***REMOVED***
		th.creds.SetRefreshToken(realm, service, tr.RefreshToken)
	***REMOVED***

	// `access_token` is equivalent to `token` and if both are specified
	// the choice is undefined.  Canonicalize `access_token` by sticking
	// things in `token`.
	if tr.AccessToken != "" ***REMOVED***
		tr.Token = tr.AccessToken
	***REMOVED***

	if tr.Token == "" ***REMOVED***
		return "", time.Time***REMOVED******REMOVED***, ErrNoToken
	***REMOVED***

	if tr.ExpiresIn < minimumTokenLifetimeSeconds ***REMOVED***
		// The default/minimum lifetime.
		tr.ExpiresIn = minimumTokenLifetimeSeconds
		logrus.Debugf("Increasing token expiration to: %d seconds", tr.ExpiresIn)
	***REMOVED***

	if tr.IssuedAt.IsZero() ***REMOVED***
		// issued_at is optional in the token response.
		tr.IssuedAt = th.clock.Now().UTC()
	***REMOVED***

	return tr.Token, tr.IssuedAt.Add(time.Duration(tr.ExpiresIn) * time.Second), nil
***REMOVED***

func (th *tokenHandler) fetchToken(params map[string]string, scopes []string) (token string, expiration time.Time, err error) ***REMOVED***
	realm, ok := params["realm"]
	if !ok ***REMOVED***
		return "", time.Time***REMOVED******REMOVED***, errors.New("no realm specified for token auth challenge")
	***REMOVED***

	// TODO(dmcgowan): Handle empty scheme and relative realm
	realmURL, err := url.Parse(realm)
	if err != nil ***REMOVED***
		return "", time.Time***REMOVED******REMOVED***, fmt.Errorf("invalid token auth challenge realm: %s", err)
	***REMOVED***

	service := params["service"]

	var refreshToken string

	if th.creds != nil ***REMOVED***
		refreshToken = th.creds.RefreshToken(realmURL, service)
	***REMOVED***

	if refreshToken != "" || th.forceOAuth ***REMOVED***
		return th.fetchTokenWithOAuth(realmURL, refreshToken, service, scopes)
	***REMOVED***

	return th.fetchTokenWithBasicAuth(realmURL, service, scopes)
***REMOVED***

type basicHandler struct ***REMOVED***
	creds CredentialStore
***REMOVED***

// NewBasicHandler creaters a new authentiation handler which adds
// basic authentication credentials to a request.
func NewBasicHandler(creds CredentialStore) AuthenticationHandler ***REMOVED***
	return &basicHandler***REMOVED***
		creds: creds,
	***REMOVED***
***REMOVED***

func (*basicHandler) Scheme() string ***REMOVED***
	return "basic"
***REMOVED***

func (bh *basicHandler) AuthorizeRequest(req *http.Request, params map[string]string) error ***REMOVED***
	if bh.creds != nil ***REMOVED***
		username, password := bh.creds.Basic(req.URL)
		if username != "" && password != "" ***REMOVED***
			req.SetBasicAuth(username, password)
			return nil
		***REMOVED***
	***REMOVED***
	return ErrNoBasicAuthCredentials
***REMOVED***
