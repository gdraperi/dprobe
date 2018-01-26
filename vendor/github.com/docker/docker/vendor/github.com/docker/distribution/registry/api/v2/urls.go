package v2

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/docker/distribution/reference"
	"github.com/gorilla/mux"
)

// URLBuilder creates registry API urls from a single base endpoint. It can be
// used to create urls for use in a registry client or server.
//
// All urls will be created from the given base, including the api version.
// For example, if a root of "/foo/" is provided, urls generated will be fall
// under "/foo/v2/...". Most application will only provide a schema, host and
// port, such as "https://localhost:5000/".
type URLBuilder struct ***REMOVED***
	root     *url.URL // url root (ie http://localhost/)
	router   *mux.Router
	relative bool
***REMOVED***

// NewURLBuilder creates a URLBuilder with provided root url object.
func NewURLBuilder(root *url.URL, relative bool) *URLBuilder ***REMOVED***
	return &URLBuilder***REMOVED***
		root:     root,
		router:   Router(),
		relative: relative,
	***REMOVED***
***REMOVED***

// NewURLBuilderFromString workes identically to NewURLBuilder except it takes
// a string argument for the root, returning an error if it is not a valid
// url.
func NewURLBuilderFromString(root string, relative bool) (*URLBuilder, error) ***REMOVED***
	u, err := url.Parse(root)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return NewURLBuilder(u, relative), nil
***REMOVED***

// NewURLBuilderFromRequest uses information from an *http.Request to
// construct the root url.
func NewURLBuilderFromRequest(r *http.Request, relative bool) *URLBuilder ***REMOVED***
	var (
		scheme = "http"
		host   = r.Host
	)

	if r.TLS != nil ***REMOVED***
		scheme = "https"
	***REMOVED*** else if len(r.URL.Scheme) > 0 ***REMOVED***
		scheme = r.URL.Scheme
	***REMOVED***

	// Handle fowarded headers
	// Prefer "Forwarded" header as defined by rfc7239 if given
	// see https://tools.ietf.org/html/rfc7239
	if forwarded := r.Header.Get("Forwarded"); len(forwarded) > 0 ***REMOVED***
		forwardedHeader, _, err := parseForwardedHeader(forwarded)
		if err == nil ***REMOVED***
			if fproto := forwardedHeader["proto"]; len(fproto) > 0 ***REMOVED***
				scheme = fproto
			***REMOVED***
			if fhost := forwardedHeader["host"]; len(fhost) > 0 ***REMOVED***
				host = fhost
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if forwardedProto := r.Header.Get("X-Forwarded-Proto"); len(forwardedProto) > 0 ***REMOVED***
			scheme = forwardedProto
		***REMOVED***
		if forwardedHost := r.Header.Get("X-Forwarded-Host"); len(forwardedHost) > 0 ***REMOVED***
			// According to the Apache mod_proxy docs, X-Forwarded-Host can be a
			// comma-separated list of hosts, to which each proxy appends the
			// requested host. We want to grab the first from this comma-separated
			// list.
			hosts := strings.SplitN(forwardedHost, ",", 2)
			host = strings.TrimSpace(hosts[0])
		***REMOVED***
	***REMOVED***

	basePath := routeDescriptorsMap[RouteNameBase].Path

	requestPath := r.URL.Path
	index := strings.Index(requestPath, basePath)

	u := &url.URL***REMOVED***
		Scheme: scheme,
		Host:   host,
	***REMOVED***

	if index > 0 ***REMOVED***
		// N.B. index+1 is important because we want to include the trailing /
		u.Path = requestPath[0 : index+1]
	***REMOVED***

	return NewURLBuilder(u, relative)
***REMOVED***

// BuildBaseURL constructs a base url for the API, typically just "/v2/".
func (ub *URLBuilder) BuildBaseURL() (string, error) ***REMOVED***
	route := ub.cloneRoute(RouteNameBase)

	baseURL, err := route.URL()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return baseURL.String(), nil
***REMOVED***

// BuildCatalogURL constructs a url get a catalog of repositories
func (ub *URLBuilder) BuildCatalogURL(values ...url.Values) (string, error) ***REMOVED***
	route := ub.cloneRoute(RouteNameCatalog)

	catalogURL, err := route.URL()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return appendValuesURL(catalogURL, values...).String(), nil
***REMOVED***

// BuildTagsURL constructs a url to list the tags in the named repository.
func (ub *URLBuilder) BuildTagsURL(name reference.Named) (string, error) ***REMOVED***
	route := ub.cloneRoute(RouteNameTags)

	tagsURL, err := route.URL("name", name.Name())
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return tagsURL.String(), nil
***REMOVED***

// BuildManifestURL constructs a url for the manifest identified by name and
// reference. The argument reference may be either a tag or digest.
func (ub *URLBuilder) BuildManifestURL(ref reference.Named) (string, error) ***REMOVED***
	route := ub.cloneRoute(RouteNameManifest)

	tagOrDigest := ""
	switch v := ref.(type) ***REMOVED***
	case reference.Tagged:
		tagOrDigest = v.Tag()
	case reference.Digested:
		tagOrDigest = v.Digest().String()
	default:
		return "", fmt.Errorf("reference must have a tag or digest")
	***REMOVED***

	manifestURL, err := route.URL("name", ref.Name(), "reference", tagOrDigest)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return manifestURL.String(), nil
***REMOVED***

// BuildBlobURL constructs the url for the blob identified by name and dgst.
func (ub *URLBuilder) BuildBlobURL(ref reference.Canonical) (string, error) ***REMOVED***
	route := ub.cloneRoute(RouteNameBlob)

	layerURL, err := route.URL("name", ref.Name(), "digest", ref.Digest().String())
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return layerURL.String(), nil
***REMOVED***

// BuildBlobUploadURL constructs a url to begin a blob upload in the
// repository identified by name.
func (ub *URLBuilder) BuildBlobUploadURL(name reference.Named, values ...url.Values) (string, error) ***REMOVED***
	route := ub.cloneRoute(RouteNameBlobUpload)

	uploadURL, err := route.URL("name", name.Name())
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return appendValuesURL(uploadURL, values...).String(), nil
***REMOVED***

// BuildBlobUploadChunkURL constructs a url for the upload identified by uuid,
// including any url values. This should generally not be used by clients, as
// this url is provided by server implementations during the blob upload
// process.
func (ub *URLBuilder) BuildBlobUploadChunkURL(name reference.Named, uuid string, values ...url.Values) (string, error) ***REMOVED***
	route := ub.cloneRoute(RouteNameBlobUploadChunk)

	uploadURL, err := route.URL("name", name.Name(), "uuid", uuid)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return appendValuesURL(uploadURL, values...).String(), nil
***REMOVED***

// clondedRoute returns a clone of the named route from the router. Routes
// must be cloned to avoid modifying them during url generation.
func (ub *URLBuilder) cloneRoute(name string) clonedRoute ***REMOVED***
	route := new(mux.Route)
	root := new(url.URL)

	*route = *ub.router.GetRoute(name) // clone the route
	*root = *ub.root

	return clonedRoute***REMOVED***Route: route, root: root, relative: ub.relative***REMOVED***
***REMOVED***

type clonedRoute struct ***REMOVED***
	*mux.Route
	root     *url.URL
	relative bool
***REMOVED***

func (cr clonedRoute) URL(pairs ...string) (*url.URL, error) ***REMOVED***
	routeURL, err := cr.Route.URL(pairs...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if cr.relative ***REMOVED***
		return routeURL, nil
	***REMOVED***

	if routeURL.Scheme == "" && routeURL.User == nil && routeURL.Host == "" ***REMOVED***
		routeURL.Path = routeURL.Path[1:]
	***REMOVED***

	url := cr.root.ResolveReference(routeURL)
	url.Scheme = cr.root.Scheme
	return url, nil
***REMOVED***

// appendValuesURL appends the parameters to the url.
func appendValuesURL(u *url.URL, values ...url.Values) *url.URL ***REMOVED***
	merged := u.Query()

	for _, v := range values ***REMOVED***
		for k, vv := range v ***REMOVED***
			merged[k] = append(merged[k], vv...)
		***REMOVED***
	***REMOVED***

	u.RawQuery = merged.Encode()
	return u
***REMOVED***

// appendValues appends the parameters to the url. Panics if the string is not
// a url.
func appendValues(u string, values ...url.Values) string ***REMOVED***
	up, err := url.Parse(u)

	if err != nil ***REMOVED***
		panic(err) // should never happen
	***REMOVED***

	return appendValuesURL(up, values...).String()
***REMOVED***
