package registry

import (
	"bytes"
	"crypto/sha256"
	"sync"
	// this is required for some certificates
	_ "crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/api/errcode"
	"github.com/docker/docker/api/types"
	registrytypes "github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/pkg/tarsum"
	"github.com/docker/docker/registry/resumable"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	// ErrRepoNotFound is returned if the repository didn't exist on the
	// remote side
	ErrRepoNotFound notFoundError = "Repository not found"
)

// A Session is used to communicate with a V1 registry
type Session struct ***REMOVED***
	indexEndpoint *V1Endpoint
	client        *http.Client
	// TODO(tiborvass): remove authConfig
	authConfig *types.AuthConfig
	id         string
***REMOVED***

type authTransport struct ***REMOVED***
	http.RoundTripper
	*types.AuthConfig

	alwaysSetBasicAuth bool
	token              []string

	mu     sync.Mutex                      // guards modReq
	modReq map[*http.Request]*http.Request // original -> modified
***REMOVED***

// AuthTransport handles the auth layer when communicating with a v1 registry (private or official)
//
// For private v1 registries, set alwaysSetBasicAuth to true.
//
// For the official v1 registry, if there isn't already an Authorization header in the request,
// but there is an X-Docker-Token header set to true, then Basic Auth will be used to set the Authorization header.
// After sending the request with the provided base http.RoundTripper, if an X-Docker-Token header, representing
// a token, is present in the response, then it gets cached and sent in the Authorization header of all subsequent
// requests.
//
// If the server sends a token without the client having requested it, it is ignored.
//
// This RoundTripper also has a CancelRequest method important for correct timeout handling.
func AuthTransport(base http.RoundTripper, authConfig *types.AuthConfig, alwaysSetBasicAuth bool) http.RoundTripper ***REMOVED***
	if base == nil ***REMOVED***
		base = http.DefaultTransport
	***REMOVED***
	return &authTransport***REMOVED***
		RoundTripper:       base,
		AuthConfig:         authConfig,
		alwaysSetBasicAuth: alwaysSetBasicAuth,
		modReq:             make(map[*http.Request]*http.Request),
	***REMOVED***
***REMOVED***

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request ***REMOVED***
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header ***REMOVED***
		r2.Header[k] = append([]string(nil), s...)
	***REMOVED***

	return r2
***REMOVED***

// RoundTrip changes an HTTP request's headers to add the necessary
// authentication-related headers
func (tr *authTransport) RoundTrip(orig *http.Request) (*http.Response, error) ***REMOVED***
	// Authorization should not be set on 302 redirect for untrusted locations.
	// This logic mirrors the behavior in addRequiredHeadersToRedirectedRequests.
	// As the authorization logic is currently implemented in RoundTrip,
	// a 302 redirect is detected by looking at the Referrer header as go http package adds said header.
	// This is safe as Docker doesn't set Referrer in other scenarios.
	if orig.Header.Get("Referer") != "" && !trustedLocation(orig) ***REMOVED***
		return tr.RoundTripper.RoundTrip(orig)
	***REMOVED***

	req := cloneRequest(orig)
	tr.mu.Lock()
	tr.modReq[orig] = req
	tr.mu.Unlock()

	if tr.alwaysSetBasicAuth ***REMOVED***
		if tr.AuthConfig == nil ***REMOVED***
			return nil, errors.New("unexpected error: empty auth config")
		***REMOVED***
		req.SetBasicAuth(tr.Username, tr.Password)
		return tr.RoundTripper.RoundTrip(req)
	***REMOVED***

	// Don't override
	if req.Header.Get("Authorization") == "" ***REMOVED***
		if req.Header.Get("X-Docker-Token") == "true" && tr.AuthConfig != nil && len(tr.Username) > 0 ***REMOVED***
			req.SetBasicAuth(tr.Username, tr.Password)
		***REMOVED*** else if len(tr.token) > 0 ***REMOVED***
			req.Header.Set("Authorization", "Token "+strings.Join(tr.token, ","))
		***REMOVED***
	***REMOVED***
	resp, err := tr.RoundTripper.RoundTrip(req)
	if err != nil ***REMOVED***
		delete(tr.modReq, orig)
		return nil, err
	***REMOVED***
	if len(resp.Header["X-Docker-Token"]) > 0 ***REMOVED***
		tr.token = resp.Header["X-Docker-Token"]
	***REMOVED***
	resp.Body = &ioutils.OnEOFReader***REMOVED***
		Rc: resp.Body,
		Fn: func() ***REMOVED***
			tr.mu.Lock()
			delete(tr.modReq, orig)
			tr.mu.Unlock()
		***REMOVED***,
	***REMOVED***
	return resp, nil
***REMOVED***

// CancelRequest cancels an in-flight request by closing its connection.
func (tr *authTransport) CancelRequest(req *http.Request) ***REMOVED***
	type canceler interface ***REMOVED***
		CancelRequest(*http.Request)
	***REMOVED***
	if cr, ok := tr.RoundTripper.(canceler); ok ***REMOVED***
		tr.mu.Lock()
		modReq := tr.modReq[req]
		delete(tr.modReq, req)
		tr.mu.Unlock()
		cr.CancelRequest(modReq)
	***REMOVED***
***REMOVED***

func authorizeClient(client *http.Client, authConfig *types.AuthConfig, endpoint *V1Endpoint) error ***REMOVED***
	var alwaysSetBasicAuth bool

	// If we're working with a standalone private registry over HTTPS, send Basic Auth headers
	// alongside all our requests.
	if endpoint.String() != IndexServer && endpoint.URL.Scheme == "https" ***REMOVED***
		info, err := endpoint.Ping()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if info.Standalone && authConfig != nil ***REMOVED***
			logrus.Debugf("Endpoint %s is eligible for private registry. Enabling decorator.", endpoint.String())
			alwaysSetBasicAuth = true
		***REMOVED***
	***REMOVED***

	// Annotate the transport unconditionally so that v2 can
	// properly fallback on v1 when an image is not found.
	client.Transport = AuthTransport(client.Transport, authConfig, alwaysSetBasicAuth)

	jar, err := cookiejar.New(nil)
	if err != nil ***REMOVED***
		return errors.New("cookiejar.New is not supposed to return an error")
	***REMOVED***
	client.Jar = jar

	return nil
***REMOVED***

func newSession(client *http.Client, authConfig *types.AuthConfig, endpoint *V1Endpoint) *Session ***REMOVED***
	return &Session***REMOVED***
		authConfig:    authConfig,
		client:        client,
		indexEndpoint: endpoint,
		id:            stringid.GenerateRandomID(),
	***REMOVED***
***REMOVED***

// NewSession creates a new session
// TODO(tiborvass): remove authConfig param once registry client v2 is vendored
func NewSession(client *http.Client, authConfig *types.AuthConfig, endpoint *V1Endpoint) (*Session, error) ***REMOVED***
	if err := authorizeClient(client, authConfig, endpoint); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return newSession(client, authConfig, endpoint), nil
***REMOVED***

// ID returns this registry session's ID.
func (r *Session) ID() string ***REMOVED***
	return r.id
***REMOVED***

// GetRemoteHistory retrieves the history of a given image from the registry.
// It returns a list of the parent's JSON files (including the requested image).
func (r *Session) GetRemoteHistory(imgID, registry string) ([]string, error) ***REMOVED***
	res, err := r.client.Get(registry + "images/" + imgID + "/ancestry")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer res.Body.Close()
	if res.StatusCode != 200 ***REMOVED***
		if res.StatusCode == 401 ***REMOVED***
			return nil, errcode.ErrorCodeUnauthorized.WithArgs()
		***REMOVED***
		return nil, newJSONError(fmt.Sprintf("Server error: %d trying to fetch remote history for %s", res.StatusCode, imgID), res)
	***REMOVED***

	var history []string
	if err := json.NewDecoder(res.Body).Decode(&history); err != nil ***REMOVED***
		return nil, fmt.Errorf("Error while reading the http response: %v", err)
	***REMOVED***

	logrus.Debugf("Ancestry: %v", history)
	return history, nil
***REMOVED***

// LookupRemoteImage checks if an image exists in the registry
func (r *Session) LookupRemoteImage(imgID, registry string) error ***REMOVED***
	res, err := r.client.Get(registry + "images/" + imgID + "/json")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	res.Body.Close()
	if res.StatusCode != 200 ***REMOVED***
		return newJSONError(fmt.Sprintf("HTTP code %d", res.StatusCode), res)
	***REMOVED***
	return nil
***REMOVED***

// GetRemoteImageJSON retrieves an image's JSON metadata from the registry.
func (r *Session) GetRemoteImageJSON(imgID, registry string) ([]byte, int64, error) ***REMOVED***
	res, err := r.client.Get(registry + "images/" + imgID + "/json")
	if err != nil ***REMOVED***
		return nil, -1, fmt.Errorf("Failed to download json: %s", err)
	***REMOVED***
	defer res.Body.Close()
	if res.StatusCode != 200 ***REMOVED***
		return nil, -1, newJSONError(fmt.Sprintf("HTTP code %d", res.StatusCode), res)
	***REMOVED***
	// if the size header is not present, then set it to '-1'
	imageSize := int64(-1)
	if hdr := res.Header.Get("X-Docker-Size"); hdr != "" ***REMOVED***
		imageSize, err = strconv.ParseInt(hdr, 10, 64)
		if err != nil ***REMOVED***
			return nil, -1, err
		***REMOVED***
	***REMOVED***

	jsonString, err := ioutil.ReadAll(res.Body)
	if err != nil ***REMOVED***
		return nil, -1, fmt.Errorf("Failed to parse downloaded json: %v (%s)", err, jsonString)
	***REMOVED***
	return jsonString, imageSize, nil
***REMOVED***

// GetRemoteImageLayer retrieves an image layer from the registry
func (r *Session) GetRemoteImageLayer(imgID, registry string, imgSize int64) (io.ReadCloser, error) ***REMOVED***
	var (
		statusCode = 0
		res        *http.Response
		err        error
		imageURL   = fmt.Sprintf("%simages/%s/layer", registry, imgID)
	)

	req, err := http.NewRequest("GET", imageURL, nil)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("Error while getting from the server: %v", err)
	***REMOVED***

	res, err = r.client.Do(req)
	if err != nil ***REMOVED***
		logrus.Debugf("Error contacting registry %s: %v", registry, err)
		// the only case err != nil && res != nil is https://golang.org/src/net/http/client.go#L515
		if res != nil ***REMOVED***
			if res.Body != nil ***REMOVED***
				res.Body.Close()
			***REMOVED***
			statusCode = res.StatusCode
		***REMOVED***
		return nil, fmt.Errorf("Server error: Status %d while fetching image layer (%s)",
			statusCode, imgID)
	***REMOVED***

	if res.StatusCode != 200 ***REMOVED***
		res.Body.Close()
		return nil, fmt.Errorf("Server error: Status %d while fetching image layer (%s)",
			res.StatusCode, imgID)
	***REMOVED***

	if res.Header.Get("Accept-Ranges") == "bytes" && imgSize > 0 ***REMOVED***
		logrus.Debug("server supports resume")
		return resumable.NewRequestReaderWithInitialResponse(r.client, req, 5, imgSize, res), nil
	***REMOVED***
	logrus.Debug("server doesn't support resume")
	return res.Body, nil
***REMOVED***

// GetRemoteTag retrieves the tag named in the askedTag argument from the given
// repository. It queries each of the registries supplied in the registries
// argument, and returns data from the first one that answers the query
// successfully.
func (r *Session) GetRemoteTag(registries []string, repositoryRef reference.Named, askedTag string) (string, error) ***REMOVED***
	repository := reference.Path(repositoryRef)

	if strings.Count(repository, "/") == 0 ***REMOVED***
		// This will be removed once the registry supports auto-resolution on
		// the "library" namespace
		repository = "library/" + repository
	***REMOVED***
	for _, host := range registries ***REMOVED***
		endpoint := fmt.Sprintf("%srepositories/%s/tags/%s", host, repository, askedTag)
		res, err := r.client.Get(endpoint)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***

		logrus.Debugf("Got status code %d from %s", res.StatusCode, endpoint)
		defer res.Body.Close()

		if res.StatusCode == 404 ***REMOVED***
			return "", ErrRepoNotFound
		***REMOVED***
		if res.StatusCode != 200 ***REMOVED***
			continue
		***REMOVED***

		var tagID string
		if err := json.NewDecoder(res.Body).Decode(&tagID); err != nil ***REMOVED***
			return "", err
		***REMOVED***
		return tagID, nil
	***REMOVED***
	return "", fmt.Errorf("Could not reach any registry endpoint")
***REMOVED***

// GetRemoteTags retrieves all tags from the given repository. It queries each
// of the registries supplied in the registries argument, and returns data from
// the first one that answers the query successfully. It returns a map with
// tag names as the keys and image IDs as the values.
func (r *Session) GetRemoteTags(registries []string, repositoryRef reference.Named) (map[string]string, error) ***REMOVED***
	repository := reference.Path(repositoryRef)

	if strings.Count(repository, "/") == 0 ***REMOVED***
		// This will be removed once the registry supports auto-resolution on
		// the "library" namespace
		repository = "library/" + repository
	***REMOVED***
	for _, host := range registries ***REMOVED***
		endpoint := fmt.Sprintf("%srepositories/%s/tags", host, repository)
		res, err := r.client.Get(endpoint)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		logrus.Debugf("Got status code %d from %s", res.StatusCode, endpoint)
		defer res.Body.Close()

		if res.StatusCode == 404 ***REMOVED***
			return nil, ErrRepoNotFound
		***REMOVED***
		if res.StatusCode != 200 ***REMOVED***
			continue
		***REMOVED***

		result := make(map[string]string)
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return result, nil
	***REMOVED***
	return nil, fmt.Errorf("Could not reach any registry endpoint")
***REMOVED***

func buildEndpointsList(headers []string, indexEp string) ([]string, error) ***REMOVED***
	var endpoints []string
	parsedURL, err := url.Parse(indexEp)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var urlScheme = parsedURL.Scheme
	// The registry's URL scheme has to match the Index'
	for _, ep := range headers ***REMOVED***
		epList := strings.Split(ep, ",")
		for _, epListElement := range epList ***REMOVED***
			endpoints = append(
				endpoints,
				fmt.Sprintf("%s://%s/v1/", urlScheme, strings.TrimSpace(epListElement)))
		***REMOVED***
	***REMOVED***
	return endpoints, nil
***REMOVED***

// GetRepositoryData returns lists of images and endpoints for the repository
func (r *Session) GetRepositoryData(name reference.Named) (*RepositoryData, error) ***REMOVED***
	repositoryTarget := fmt.Sprintf("%srepositories/%s/images", r.indexEndpoint.String(), reference.Path(name))

	logrus.Debugf("[registry] Calling GET %s", repositoryTarget)

	req, err := http.NewRequest("GET", repositoryTarget, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// this will set basic auth in r.client.Transport and send cached X-Docker-Token headers for all subsequent requests
	req.Header.Set("X-Docker-Token", "true")
	res, err := r.client.Do(req)
	if err != nil ***REMOVED***
		// check if the error is because of i/o timeout
		// and return a non-obtuse error message for users
		// "Get https://index.docker.io/v1/repositories/library/busybox/images: i/o timeout"
		// was a top search on the docker user forum
		if isTimeout(err) ***REMOVED***
			return nil, fmt.Errorf("network timed out while trying to connect to %s. You may want to check your internet connection or if you are behind a proxy", repositoryTarget)
		***REMOVED***
		return nil, fmt.Errorf("Error while pulling image: %v", err)
	***REMOVED***
	defer res.Body.Close()
	if res.StatusCode == 401 ***REMOVED***
		return nil, errcode.ErrorCodeUnauthorized.WithArgs()
	***REMOVED***
	// TODO: Right now we're ignoring checksums in the response body.
	// In the future, we need to use them to check image validity.
	if res.StatusCode == 404 ***REMOVED***
		return nil, newJSONError(fmt.Sprintf("HTTP code: %d", res.StatusCode), res)
	***REMOVED*** else if res.StatusCode != 200 ***REMOVED***
		errBody, err := ioutil.ReadAll(res.Body)
		if err != nil ***REMOVED***
			logrus.Debugf("Error reading response body: %s", err)
		***REMOVED***
		return nil, newJSONError(fmt.Sprintf("Error: Status %d trying to pull repository %s: %q", res.StatusCode, reference.Path(name), errBody), res)
	***REMOVED***

	var endpoints []string
	if res.Header.Get("X-Docker-Endpoints") != "" ***REMOVED***
		endpoints, err = buildEndpointsList(res.Header["X-Docker-Endpoints"], r.indexEndpoint.String())
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// Assume the endpoint is on the same host
		endpoints = append(endpoints, fmt.Sprintf("%s://%s/v1/", r.indexEndpoint.URL.Scheme, req.URL.Host))
	***REMOVED***

	remoteChecksums := []*ImgData***REMOVED******REMOVED***
	if err := json.NewDecoder(res.Body).Decode(&remoteChecksums); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Forge a better object from the retrieved data
	imgsData := make(map[string]*ImgData, len(remoteChecksums))
	for _, elem := range remoteChecksums ***REMOVED***
		imgsData[elem.ID] = elem
	***REMOVED***

	return &RepositoryData***REMOVED***
		ImgList:   imgsData,
		Endpoints: endpoints,
	***REMOVED***, nil
***REMOVED***

// PushImageChecksumRegistry uploads checksums for an image
func (r *Session) PushImageChecksumRegistry(imgData *ImgData, registry string) error ***REMOVED***
	u := registry + "images/" + imgData.ID + "/checksum"

	logrus.Debugf("[registry] Calling PUT %s", u)

	req, err := http.NewRequest("PUT", u, nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	req.Header.Set("X-Docker-Checksum", imgData.Checksum)
	req.Header.Set("X-Docker-Checksum-Payload", imgData.ChecksumPayload)

	res, err := r.client.Do(req)
	if err != nil ***REMOVED***
		return fmt.Errorf("Failed to upload metadata: %v", err)
	***REMOVED***
	defer res.Body.Close()
	if len(res.Cookies()) > 0 ***REMOVED***
		r.client.Jar.SetCookies(req.URL, res.Cookies())
	***REMOVED***
	if res.StatusCode != 200 ***REMOVED***
		errBody, err := ioutil.ReadAll(res.Body)
		if err != nil ***REMOVED***
			return fmt.Errorf("HTTP code %d while uploading metadata and error when trying to parse response body: %s", res.StatusCode, err)
		***REMOVED***
		var jsonBody map[string]string
		if err := json.Unmarshal(errBody, &jsonBody); err != nil ***REMOVED***
			errBody = []byte(err.Error())
		***REMOVED*** else if jsonBody["error"] == "Image already exists" ***REMOVED***
			return ErrAlreadyExists
		***REMOVED***
		return fmt.Errorf("HTTP code %d while uploading metadata: %q", res.StatusCode, errBody)
	***REMOVED***
	return nil
***REMOVED***

// PushImageJSONRegistry pushes JSON metadata for a local image to the registry
func (r *Session) PushImageJSONRegistry(imgData *ImgData, jsonRaw []byte, registry string) error ***REMOVED***

	u := registry + "images/" + imgData.ID + "/json"

	logrus.Debugf("[registry] Calling PUT %s", u)

	req, err := http.NewRequest("PUT", u, bytes.NewReader(jsonRaw))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	req.Header.Add("Content-type", "application/json")

	res, err := r.client.Do(req)
	if err != nil ***REMOVED***
		return fmt.Errorf("Failed to upload metadata: %s", err)
	***REMOVED***
	defer res.Body.Close()
	if res.StatusCode == 401 && strings.HasPrefix(registry, "http://") ***REMOVED***
		return newJSONError("HTTP code 401, Docker will not send auth headers over HTTP.", res)
	***REMOVED***
	if res.StatusCode != 200 ***REMOVED***
		errBody, err := ioutil.ReadAll(res.Body)
		if err != nil ***REMOVED***
			return newJSONError(fmt.Sprintf("HTTP code %d while uploading metadata and error when trying to parse response body: %s", res.StatusCode, err), res)
		***REMOVED***
		var jsonBody map[string]string
		if err := json.Unmarshal(errBody, &jsonBody); err != nil ***REMOVED***
			errBody = []byte(err.Error())
		***REMOVED*** else if jsonBody["error"] == "Image already exists" ***REMOVED***
			return ErrAlreadyExists
		***REMOVED***
		return newJSONError(fmt.Sprintf("HTTP code %d while uploading metadata: %q", res.StatusCode, errBody), res)
	***REMOVED***
	return nil
***REMOVED***

// PushImageLayerRegistry sends the checksum of an image layer to the registry
func (r *Session) PushImageLayerRegistry(imgID string, layer io.Reader, registry string, jsonRaw []byte) (checksum string, checksumPayload string, err error) ***REMOVED***
	u := registry + "images/" + imgID + "/layer"

	logrus.Debugf("[registry] Calling PUT %s", u)

	tarsumLayer, err := tarsum.NewTarSum(layer, false, tarsum.Version0)
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***
	h := sha256.New()
	h.Write(jsonRaw)
	h.Write([]byte***REMOVED***'\n'***REMOVED***)
	checksumLayer := io.TeeReader(tarsumLayer, h)

	req, err := http.NewRequest("PUT", u, checksumLayer)
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***
	req.Header.Add("Content-Type", "application/octet-stream")
	req.ContentLength = -1
	req.TransferEncoding = []string***REMOVED***"chunked"***REMOVED***
	res, err := r.client.Do(req)
	if err != nil ***REMOVED***
		return "", "", fmt.Errorf("Failed to upload layer: %v", err)
	***REMOVED***
	if rc, ok := layer.(io.Closer); ok ***REMOVED***
		if err := rc.Close(); err != nil ***REMOVED***
			return "", "", err
		***REMOVED***
	***REMOVED***
	defer res.Body.Close()

	if res.StatusCode != 200 ***REMOVED***
		errBody, err := ioutil.ReadAll(res.Body)
		if err != nil ***REMOVED***
			return "", "", newJSONError(fmt.Sprintf("HTTP code %d while uploading metadata and error when trying to parse response body: %s", res.StatusCode, err), res)
		***REMOVED***
		return "", "", newJSONError(fmt.Sprintf("Received HTTP code %d while uploading layer: %q", res.StatusCode, errBody), res)
	***REMOVED***

	checksumPayload = "sha256:" + hex.EncodeToString(h.Sum(nil))
	return tarsumLayer.Sum(jsonRaw), checksumPayload, nil
***REMOVED***

// PushRegistryTag pushes a tag on the registry.
// Remote has the format '<user>/<repo>
func (r *Session) PushRegistryTag(remote reference.Named, revision, tag, registry string) error ***REMOVED***
	// "jsonify" the string
	revision = "\"" + revision + "\""
	path := fmt.Sprintf("repositories/%s/tags/%s", reference.Path(remote), tag)

	req, err := http.NewRequest("PUT", registry+path, strings.NewReader(revision))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	req.Header.Add("Content-type", "application/json")
	req.ContentLength = int64(len(revision))
	res, err := r.client.Do(req)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	res.Body.Close()
	if res.StatusCode != 200 && res.StatusCode != 201 ***REMOVED***
		return newJSONError(fmt.Sprintf("Internal server error: %d trying to push tag %s on %s", res.StatusCode, tag, reference.Path(remote)), res)
	***REMOVED***
	return nil
***REMOVED***

// PushImageJSONIndex uploads an image list to the repository
func (r *Session) PushImageJSONIndex(remote reference.Named, imgList []*ImgData, validate bool, regs []string) (*RepositoryData, error) ***REMOVED***
	cleanImgList := []*ImgData***REMOVED******REMOVED***
	if validate ***REMOVED***
		for _, elem := range imgList ***REMOVED***
			if elem.Checksum != "" ***REMOVED***
				cleanImgList = append(cleanImgList, elem)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		cleanImgList = imgList
	***REMOVED***

	imgListJSON, err := json.Marshal(cleanImgList)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var suffix string
	if validate ***REMOVED***
		suffix = "images"
	***REMOVED***
	u := fmt.Sprintf("%srepositories/%s/%s", r.indexEndpoint.String(), reference.Path(remote), suffix)
	logrus.Debugf("[registry] PUT %s", u)
	logrus.Debugf("Image list pushed to index:\n%s", imgListJSON)
	headers := map[string][]string***REMOVED***
		"Content-type": ***REMOVED***"application/json"***REMOVED***,
		// this will set basic auth in r.client.Transport and send cached X-Docker-Token headers for all subsequent requests
		"X-Docker-Token": ***REMOVED***"true"***REMOVED***,
	***REMOVED***
	if validate ***REMOVED***
		headers["X-Docker-Endpoints"] = regs
	***REMOVED***

	// Redirect if necessary
	var res *http.Response
	for ***REMOVED***
		if res, err = r.putImageRequest(u, headers, imgListJSON); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if !shouldRedirect(res) ***REMOVED***
			break
		***REMOVED***
		res.Body.Close()
		u = res.Header.Get("Location")
		logrus.Debugf("Redirected to %s", u)
	***REMOVED***
	defer res.Body.Close()

	if res.StatusCode == 401 ***REMOVED***
		return nil, errcode.ErrorCodeUnauthorized.WithArgs()
	***REMOVED***

	var tokens, endpoints []string
	if !validate ***REMOVED***
		if res.StatusCode != 200 && res.StatusCode != 201 ***REMOVED***
			errBody, err := ioutil.ReadAll(res.Body)
			if err != nil ***REMOVED***
				logrus.Debugf("Error reading response body: %s", err)
			***REMOVED***
			return nil, newJSONError(fmt.Sprintf("Error: Status %d trying to push repository %s: %q", res.StatusCode, reference.Path(remote), errBody), res)
		***REMOVED***
		tokens = res.Header["X-Docker-Token"]
		logrus.Debugf("Auth token: %v", tokens)

		if res.Header.Get("X-Docker-Endpoints") == "" ***REMOVED***
			return nil, fmt.Errorf("Index response didn't contain any endpoints")
		***REMOVED***
		endpoints, err = buildEndpointsList(res.Header["X-Docker-Endpoints"], r.indexEndpoint.String())
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if res.StatusCode != 204 ***REMOVED***
			errBody, err := ioutil.ReadAll(res.Body)
			if err != nil ***REMOVED***
				logrus.Debugf("Error reading response body: %s", err)
			***REMOVED***
			return nil, newJSONError(fmt.Sprintf("Error: Status %d trying to push checksums %s: %q", res.StatusCode, reference.Path(remote), errBody), res)
		***REMOVED***
	***REMOVED***

	return &RepositoryData***REMOVED***
		Endpoints: endpoints,
	***REMOVED***, nil
***REMOVED***

func (r *Session) putImageRequest(u string, headers map[string][]string, body []byte) (*http.Response, error) ***REMOVED***
	req, err := http.NewRequest("PUT", u, bytes.NewReader(body))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	req.ContentLength = int64(len(body))
	for k, v := range headers ***REMOVED***
		req.Header[k] = v
	***REMOVED***
	response, err := r.client.Do(req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return response, nil
***REMOVED***

func shouldRedirect(response *http.Response) bool ***REMOVED***
	return response.StatusCode >= 300 && response.StatusCode < 400
***REMOVED***

// SearchRepositories performs a search against the remote repository
func (r *Session) SearchRepositories(term string, limit int) (*registrytypes.SearchResults, error) ***REMOVED***
	if limit < 1 || limit > 100 ***REMOVED***
		return nil, errdefs.InvalidParameter(errors.Errorf("Limit %d is outside the range of [1, 100]", limit))
	***REMOVED***
	logrus.Debugf("Index server: %s", r.indexEndpoint)
	u := r.indexEndpoint.String() + "search?q=" + url.QueryEscape(term) + "&n=" + url.QueryEscape(fmt.Sprintf("%d", limit))

	req, err := http.NewRequest("GET", u, nil)
	if err != nil ***REMOVED***
		return nil, errors.Wrap(errdefs.InvalidParameter(err), "Error building request")
	***REMOVED***
	// Have the AuthTransport send authentication, when logged in.
	req.Header.Set("X-Docker-Token", "true")
	res, err := r.client.Do(req)
	if err != nil ***REMOVED***
		return nil, errdefs.System(err)
	***REMOVED***
	defer res.Body.Close()
	if res.StatusCode != 200 ***REMOVED***
		return nil, newJSONError(fmt.Sprintf("Unexpected status code %d", res.StatusCode), res)
	***REMOVED***
	result := new(registrytypes.SearchResults)
	return result, errors.Wrap(json.NewDecoder(res.Body).Decode(result), "error decoding registry search results")
***REMOVED***

func isTimeout(err error) bool ***REMOVED***
	type timeout interface ***REMOVED***
		Timeout() bool
	***REMOVED***
	e := err
	switch urlErr := err.(type) ***REMOVED***
	case *url.Error:
		e = urlErr.Err
	***REMOVED***
	t, ok := e.(timeout)
	return ok && t.Timeout()
***REMOVED***

func newJSONError(msg string, res *http.Response) error ***REMOVED***
	return &jsonmessage.JSONError***REMOVED***
		Message: msg,
		Code:    res.StatusCode,
	***REMOVED***
***REMOVED***
