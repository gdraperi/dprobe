package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/reference"
	"github.com/containerd/containerd/remotes"
	digest "github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context/ctxhttp"
)

var (
	// ErrNoToken is returned if a request is successful but the body does not
	// contain an authorization token.
	ErrNoToken = errors.New("authorization server did not include a token in the response")

	// ErrInvalidAuthorization is used when credentials are passed to a server but
	// those credentials are rejected.
	ErrInvalidAuthorization = errors.New("authorization failed")
)

type dockerResolver struct ***REMOVED***
	credentials func(string) (string, string, error)
	plainHTTP   bool
	client      *http.Client
	tracker     StatusTracker
***REMOVED***

// ResolverOptions are used to configured a new Docker register resolver
type ResolverOptions struct ***REMOVED***
	// Credentials provides username and secret given a host.
	// If username is empty but a secret is given, that secret
	// is interpretted as a long lived token.
	Credentials func(string) (string, string, error)

	// PlainHTTP specifies to use plain http and not https
	PlainHTTP bool

	// Client is the http client to used when making registry requests
	Client *http.Client

	// Tracker is used to track uploads to the registry. This is used
	// since the registry does not have upload tracking and the existing
	// mechanism for getting blob upload status is expensive.
	Tracker StatusTracker
***REMOVED***

// NewResolver returns a new resolver to a Docker registry
func NewResolver(options ResolverOptions) remotes.Resolver ***REMOVED***
	tracker := options.Tracker
	if tracker == nil ***REMOVED***
		tracker = NewInMemoryTracker()
	***REMOVED***
	return &dockerResolver***REMOVED***
		credentials: options.Credentials,
		plainHTTP:   options.PlainHTTP,
		client:      options.Client,
		tracker:     tracker,
	***REMOVED***
***REMOVED***

var _ remotes.Resolver = &dockerResolver***REMOVED******REMOVED***

func (r *dockerResolver) Resolve(ctx context.Context, ref string) (string, ocispec.Descriptor, error) ***REMOVED***
	refspec, err := reference.Parse(ref)
	if err != nil ***REMOVED***
		return "", ocispec.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***

	if refspec.Object == "" ***REMOVED***
		return "", ocispec.Descriptor***REMOVED******REMOVED***, reference.ErrObjectRequired
	***REMOVED***

	base, err := r.base(refspec)
	if err != nil ***REMOVED***
		return "", ocispec.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***

	fetcher := dockerFetcher***REMOVED***
		dockerBase: base,
	***REMOVED***

	var (
		urls []string
		dgst = refspec.Digest()
	)

	if dgst != "" ***REMOVED***
		if err := dgst.Validate(); err != nil ***REMOVED***
			// need to fail here, since we can't actually resolve the invalid
			// digest.
			return "", ocispec.Descriptor***REMOVED******REMOVED***, err
		***REMOVED***

		// turns out, we have a valid digest, make a url.
		urls = append(urls, fetcher.url("manifests", dgst.String()))

		// fallback to blobs on not found.
		urls = append(urls, fetcher.url("blobs", dgst.String()))
	***REMOVED*** else ***REMOVED***
		urls = append(urls, fetcher.url("manifests", refspec.Object))
	***REMOVED***

	ctx, err = contextWithRepositoryScope(ctx, refspec, false)
	if err != nil ***REMOVED***
		return "", ocispec.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***
	for _, u := range urls ***REMOVED***
		req, err := http.NewRequest(http.MethodHead, u, nil)
		if err != nil ***REMOVED***
			return "", ocispec.Descriptor***REMOVED******REMOVED***, err
		***REMOVED***

		// set headers for all the types we support for resolution.
		req.Header.Set("Accept", strings.Join([]string***REMOVED***
			images.MediaTypeDockerSchema2Manifest,
			images.MediaTypeDockerSchema2ManifestList,
			ocispec.MediaTypeImageManifest,
			ocispec.MediaTypeImageIndex, "*"***REMOVED***, ", "))

		log.G(ctx).Debug("resolving")
		resp, err := fetcher.doRequestWithRetries(ctx, req, nil)
		if err != nil ***REMOVED***
			return "", ocispec.Descriptor***REMOVED******REMOVED***, err
		***REMOVED***
		resp.Body.Close() // don't care about body contents.

		if resp.StatusCode > 299 ***REMOVED***
			if resp.StatusCode == http.StatusNotFound ***REMOVED***
				continue
			***REMOVED***
			return "", ocispec.Descriptor***REMOVED******REMOVED***, errors.Errorf("unexpected status code %v: %v", u, resp.Status)
		***REMOVED***

		// this is the only point at which we trust the registry. we use the
		// content headers to assemble a descriptor for the name. when this becomes
		// more robust, we mostly get this information from a secure trust store.
		dgstHeader := digest.Digest(resp.Header.Get("Docker-Content-Digest"))

		if dgstHeader != "" ***REMOVED***
			if err := dgstHeader.Validate(); err != nil ***REMOVED***
				return "", ocispec.Descriptor***REMOVED******REMOVED***, errors.Wrapf(err, "%q in header not a valid digest", dgstHeader)
			***REMOVED***
			dgst = dgstHeader
		***REMOVED***

		if dgst == "" ***REMOVED***
			return "", ocispec.Descriptor***REMOVED******REMOVED***, errors.Errorf("could not resolve digest for %v", ref)
		***REMOVED***

		var (
			size       int64
			sizeHeader = resp.Header.Get("Content-Length")
		)

		size, err = strconv.ParseInt(sizeHeader, 10, 64)
		if err != nil ***REMOVED***

			return "", ocispec.Descriptor***REMOVED******REMOVED***, errors.Wrapf(err, "invalid size header: %q", sizeHeader)
		***REMOVED***
		if size < 0 ***REMOVED***
			return "", ocispec.Descriptor***REMOVED******REMOVED***, errors.Errorf("%q in header not a valid size", sizeHeader)
		***REMOVED***

		desc := ocispec.Descriptor***REMOVED***
			Digest:    dgst,
			MediaType: resp.Header.Get("Content-Type"), // need to strip disposition?
			Size:      size,
		***REMOVED***

		log.G(ctx).WithField("desc.digest", desc.Digest).Debug("resolved")
		return ref, desc, nil
	***REMOVED***

	return "", ocispec.Descriptor***REMOVED******REMOVED***, errors.Errorf("%v not found", ref)
***REMOVED***

func (r *dockerResolver) Fetcher(ctx context.Context, ref string) (remotes.Fetcher, error) ***REMOVED***
	refspec, err := reference.Parse(ref)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	base, err := r.base(refspec)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return dockerFetcher***REMOVED***
		dockerBase: base,
	***REMOVED***, nil
***REMOVED***

func (r *dockerResolver) Pusher(ctx context.Context, ref string) (remotes.Pusher, error) ***REMOVED***
	refspec, err := reference.Parse(ref)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Manifests can be pushed by digest like any other object, but the passed in
	// reference cannot take a digest without the associated content. A tag is allowed
	// and will be used to tag pushed manifests.
	if refspec.Object != "" && strings.Contains(refspec.Object, "@") ***REMOVED***
		return nil, errors.New("cannot use digest reference for push locator")
	***REMOVED***

	base, err := r.base(refspec)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return dockerPusher***REMOVED***
		dockerBase: base,
		tag:        refspec.Object,
		tracker:    r.tracker,
	***REMOVED***, nil
***REMOVED***

type dockerBase struct ***REMOVED***
	refspec reference.Spec
	base    url.URL
	token   string

	client   *http.Client
	useBasic bool
	username string
	secret   string
***REMOVED***

func (r *dockerResolver) base(refspec reference.Spec) (*dockerBase, error) ***REMOVED***
	var (
		err              error
		base             url.URL
		username, secret string
	)

	host := refspec.Hostname()
	base.Scheme = "https"

	if host == "docker.io" ***REMOVED***
		base.Host = "registry-1.docker.io"
	***REMOVED*** else ***REMOVED***
		base.Host = host

		if r.plainHTTP || strings.HasPrefix(host, "localhost:") ***REMOVED***
			base.Scheme = "http"
		***REMOVED***
	***REMOVED***

	if r.credentials != nil ***REMOVED***
		username, secret, err = r.credentials(base.Host)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	prefix := strings.TrimPrefix(refspec.Locator, host+"/")
	base.Path = path.Join("/v2", prefix)

	return &dockerBase***REMOVED***
		refspec:  refspec,
		base:     base,
		client:   r.client,
		username: username,
		secret:   secret,
	***REMOVED***, nil
***REMOVED***

func (r *dockerBase) url(ps ...string) string ***REMOVED***
	url := r.base
	url.Path = path.Join(url.Path, path.Join(ps...))
	return url.String()
***REMOVED***

func (r *dockerBase) authorize(req *http.Request) ***REMOVED***
	if r.useBasic ***REMOVED***
		req.SetBasicAuth(r.username, r.secret)
	***REMOVED*** else if r.token != "" ***REMOVED***
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.token))
	***REMOVED***
***REMOVED***

func (r *dockerBase) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) ***REMOVED***
	ctx = log.WithLogger(ctx, log.G(ctx).WithField("url", req.URL.String()))
	log.G(ctx).WithField("request.headers", req.Header).WithField("request.method", req.Method).Debug("do request")
	r.authorize(req)
	resp, err := ctxhttp.Do(ctx, r.client, req)
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "failed to do request")
	***REMOVED***
	log.G(ctx).WithFields(logrus.Fields***REMOVED***
		"status":           resp.Status,
		"response.headers": resp.Header,
	***REMOVED***).Debug("fetch response received")
	return resp, nil
***REMOVED***

func (r *dockerBase) doRequestWithRetries(ctx context.Context, req *http.Request, responses []*http.Response) (*http.Response, error) ***REMOVED***
	resp, err := r.doRequest(ctx, req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	responses = append(responses, resp)
	req, err = r.retryRequest(ctx, req, responses)
	if err != nil ***REMOVED***
		resp.Body.Close()
		return nil, err
	***REMOVED***
	if req != nil ***REMOVED***
		resp.Body.Close()
		return r.doRequestWithRetries(ctx, req, responses)
	***REMOVED***
	return resp, err
***REMOVED***

func (r *dockerBase) retryRequest(ctx context.Context, req *http.Request, responses []*http.Response) (*http.Request, error) ***REMOVED***
	if len(responses) > 5 ***REMOVED***
		return nil, nil
	***REMOVED***
	last := responses[len(responses)-1]
	if last.StatusCode == http.StatusUnauthorized ***REMOVED***
		log.G(ctx).WithField("header", last.Header.Get("WWW-Authenticate")).Debug("Unauthorized")
		for _, c := range parseAuthHeader(last.Header) ***REMOVED***
			if c.scheme == bearerAuth ***REMOVED***
				if err := invalidAuthorization(c, responses); err != nil ***REMOVED***
					r.token = ""
					return nil, err
				***REMOVED***
				if err := r.setTokenAuth(ctx, c.parameters); err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				return copyRequest(req)
			***REMOVED*** else if c.scheme == basicAuth ***REMOVED***
				if r.username != "" && r.secret != "" ***REMOVED***
					r.useBasic = true
				***REMOVED***
				return copyRequest(req)
			***REMOVED***
		***REMOVED***
		return nil, nil
	***REMOVED*** else if last.StatusCode == http.StatusMethodNotAllowed && req.Method == http.MethodHead ***REMOVED***
		// Support registries which have not properly implemented the HEAD method for
		// manifests endpoint
		if strings.Contains(req.URL.Path, "/manifests/") ***REMOVED***
			// TODO: copy request?
			req.Method = http.MethodGet
			return copyRequest(req)
		***REMOVED***
	***REMOVED***

	// TODO: Handle 50x errors accounting for attempt history
	return nil, nil
***REMOVED***

func invalidAuthorization(c challenge, responses []*http.Response) error ***REMOVED***
	errStr := c.parameters["error"]
	if errStr == "" ***REMOVED***
		return nil
	***REMOVED***

	n := len(responses)
	if n == 1 || (n > 1 && !sameRequest(responses[n-2].Request, responses[n-1].Request)) ***REMOVED***
		return nil
	***REMOVED***

	return errors.Wrapf(ErrInvalidAuthorization, "server message: %s", errStr)
***REMOVED***

func sameRequest(r1, r2 *http.Request) bool ***REMOVED***
	if r1.Method != r2.Method ***REMOVED***
		return false
	***REMOVED***
	if *r1.URL != *r2.URL ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

func copyRequest(req *http.Request) (*http.Request, error) ***REMOVED***
	ireq := *req
	if ireq.GetBody != nil ***REMOVED***
		var err error
		ireq.Body, err = ireq.GetBody()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return &ireq, nil
***REMOVED***

func (r *dockerBase) setTokenAuth(ctx context.Context, params map[string]string) error ***REMOVED***
	realm, ok := params["realm"]
	if !ok ***REMOVED***
		return errors.New("no realm specified for token auth challenge")
	***REMOVED***

	realmURL, err := url.Parse(realm)
	if err != nil ***REMOVED***
		return fmt.Errorf("invalid token auth challenge realm: %s", err)
	***REMOVED***

	to := tokenOptions***REMOVED***
		realm:   realmURL.String(),
		service: params["service"],
	***REMOVED***

	to.scopes = getTokenScopes(ctx, params)
	if len(to.scopes) == 0 ***REMOVED***
		return errors.Errorf("no scope specified for token auth challenge")
	***REMOVED***
	if r.secret != "" ***REMOVED***
		// Credential information is provided, use oauth POST endpoint
		r.token, err = r.fetchTokenWithOAuth(ctx, to)
		if err != nil ***REMOVED***
			return errors.Wrap(err, "failed to fetch oauth token")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// Do request anonymously
		r.token, err = r.getToken(ctx, to)
		if err != nil ***REMOVED***
			return errors.Wrap(err, "failed to fetch anonymous token")
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

type tokenOptions struct ***REMOVED***
	realm   string
	service string
	scopes  []string
***REMOVED***

type postTokenResponse struct ***REMOVED***
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"`
	IssuedAt     time.Time `json:"issued_at"`
	Scope        string    `json:"scope"`
***REMOVED***

func (r *dockerBase) fetchTokenWithOAuth(ctx context.Context, to tokenOptions) (string, error) ***REMOVED***
	form := url.Values***REMOVED******REMOVED***
	form.Set("scope", strings.Join(to.scopes, " "))
	form.Set("service", to.service)
	// TODO: Allow setting client_id
	form.Set("client_id", "containerd-dist-tool")

	if r.username == "" ***REMOVED***
		form.Set("grant_type", "refresh_token")
		form.Set("refresh_token", r.secret)
	***REMOVED*** else ***REMOVED***
		form.Set("grant_type", "password")
		form.Set("username", r.username)
		form.Set("password", r.secret)
	***REMOVED***

	resp, err := ctxhttp.PostForm(ctx, r.client, to.realm, form)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	defer resp.Body.Close()

	// Registries without support for POST may return 404 for POST /v2/token.
	// As of September 2017, GCR is known to return 404.
	if (resp.StatusCode == 405 && r.username != "") || resp.StatusCode == 404 ***REMOVED***
		return r.getToken(ctx, to)
	***REMOVED*** else if resp.StatusCode < 200 || resp.StatusCode >= 400 ***REMOVED***
		b, _ := ioutil.ReadAll(io.LimitReader(resp.Body, 64000)) // 64KB
		log.G(ctx).WithFields(logrus.Fields***REMOVED***
			"status": resp.Status,
			"body":   string(b),
		***REMOVED***).Debugf("token request failed")
		// TODO: handle error body and write debug output
		return "", errors.Errorf("unexpected status: %s", resp.Status)
	***REMOVED***

	decoder := json.NewDecoder(resp.Body)

	var tr postTokenResponse
	if err = decoder.Decode(&tr); err != nil ***REMOVED***
		return "", fmt.Errorf("unable to decode token response: %s", err)
	***REMOVED***

	return tr.AccessToken, nil
***REMOVED***

type getTokenResponse struct ***REMOVED***
	Token        string    `json:"token"`
	AccessToken  string    `json:"access_token"`
	ExpiresIn    int       `json:"expires_in"`
	IssuedAt     time.Time `json:"issued_at"`
	RefreshToken string    `json:"refresh_token"`
***REMOVED***

// getToken fetches a token using a GET request
func (r *dockerBase) getToken(ctx context.Context, to tokenOptions) (string, error) ***REMOVED***
	req, err := http.NewRequest("GET", to.realm, nil)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	reqParams := req.URL.Query()

	if to.service != "" ***REMOVED***
		reqParams.Add("service", to.service)
	***REMOVED***

	for _, scope := range to.scopes ***REMOVED***
		reqParams.Add("scope", scope)
	***REMOVED***

	if r.secret != "" ***REMOVED***
		req.SetBasicAuth(r.username, r.secret)
	***REMOVED***

	req.URL.RawQuery = reqParams.Encode()

	resp, err := ctxhttp.Do(ctx, r.client, req)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 ***REMOVED***
		// TODO: handle error body and write debug output
		return "", errors.Errorf("unexpected status: %s", resp.Status)
	***REMOVED***

	decoder := json.NewDecoder(resp.Body)

	var tr getTokenResponse
	if err = decoder.Decode(&tr); err != nil ***REMOVED***
		return "", fmt.Errorf("unable to decode token response: %s", err)
	***REMOVED***

	// `access_token` is equivalent to `token` and if both are specified
	// the choice is undefined.  Canonicalize `access_token` by sticking
	// things in `token`.
	if tr.AccessToken != "" ***REMOVED***
		tr.Token = tr.AccessToken
	***REMOVED***

	if tr.Token == "" ***REMOVED***
		return "", ErrNoToken
	***REMOVED***

	return tr.Token, nil
***REMOVED***
