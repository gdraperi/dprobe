package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/api/v2"
	"github.com/docker/distribution/registry/client/transport"
	"github.com/docker/distribution/registry/storage/cache"
	"github.com/docker/distribution/registry/storage/cache/memory"
	"github.com/opencontainers/go-digest"
)

// Registry provides an interface for calling Repositories, which returns a catalog of repositories.
type Registry interface ***REMOVED***
	Repositories(ctx context.Context, repos []string, last string) (n int, err error)
***REMOVED***

// checkHTTPRedirect is a callback that can manipulate redirected HTTP
// requests. It is used to preserve Accept and Range headers.
func checkHTTPRedirect(req *http.Request, via []*http.Request) error ***REMOVED***
	if len(via) >= 10 ***REMOVED***
		return errors.New("stopped after 10 redirects")
	***REMOVED***

	if len(via) > 0 ***REMOVED***
		for headerName, headerVals := range via[0].Header ***REMOVED***
			if headerName != "Accept" && headerName != "Range" ***REMOVED***
				continue
			***REMOVED***
			for _, val := range headerVals ***REMOVED***
				// Don't add to redirected request if redirected
				// request already has a header with the same
				// name and value.
				hasValue := false
				for _, existingVal := range req.Header[headerName] ***REMOVED***
					if existingVal == val ***REMOVED***
						hasValue = true
						break
					***REMOVED***
				***REMOVED***
				if !hasValue ***REMOVED***
					req.Header.Add(headerName, val)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// NewRegistry creates a registry namespace which can be used to get a listing of repositories
func NewRegistry(ctx context.Context, baseURL string, transport http.RoundTripper) (Registry, error) ***REMOVED***
	ub, err := v2.NewURLBuilderFromString(baseURL, false)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	client := &http.Client***REMOVED***
		Transport:     transport,
		Timeout:       1 * time.Minute,
		CheckRedirect: checkHTTPRedirect,
	***REMOVED***

	return &registry***REMOVED***
		client:  client,
		ub:      ub,
		context: ctx,
	***REMOVED***, nil
***REMOVED***

type registry struct ***REMOVED***
	client  *http.Client
	ub      *v2.URLBuilder
	context context.Context
***REMOVED***

// Repositories returns a lexigraphically sorted catalog given a base URL.  The 'entries' slice will be filled up to the size
// of the slice, starting at the value provided in 'last'.  The number of entries will be returned along with io.EOF if there
// are no more entries
func (r *registry) Repositories(ctx context.Context, entries []string, last string) (int, error) ***REMOVED***
	var numFilled int
	var returnErr error

	values := buildCatalogValues(len(entries), last)
	u, err := r.ub.BuildCatalogURL(values)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	resp, err := r.client.Get(u)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	defer resp.Body.Close()

	if SuccessStatus(resp.StatusCode) ***REMOVED***
		var ctlg struct ***REMOVED***
			Repositories []string `json:"repositories"`
		***REMOVED***
		decoder := json.NewDecoder(resp.Body)

		if err := decoder.Decode(&ctlg); err != nil ***REMOVED***
			return 0, err
		***REMOVED***

		for cnt := range ctlg.Repositories ***REMOVED***
			entries[cnt] = ctlg.Repositories[cnt]
		***REMOVED***
		numFilled = len(ctlg.Repositories)

		link := resp.Header.Get("Link")
		if link == "" ***REMOVED***
			returnErr = io.EOF
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		return 0, HandleErrorResponse(resp)
	***REMOVED***

	return numFilled, returnErr
***REMOVED***

// NewRepository creates a new Repository for the given repository name and base URL.
func NewRepository(ctx context.Context, name reference.Named, baseURL string, transport http.RoundTripper) (distribution.Repository, error) ***REMOVED***
	ub, err := v2.NewURLBuilderFromString(baseURL, false)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	client := &http.Client***REMOVED***
		Transport:     transport,
		CheckRedirect: checkHTTPRedirect,
		// TODO(dmcgowan): create cookie jar
	***REMOVED***

	return &repository***REMOVED***
		client:  client,
		ub:      ub,
		name:    name,
		context: ctx,
	***REMOVED***, nil
***REMOVED***

type repository struct ***REMOVED***
	client  *http.Client
	ub      *v2.URLBuilder
	context context.Context
	name    reference.Named
***REMOVED***

func (r *repository) Named() reference.Named ***REMOVED***
	return r.name
***REMOVED***

func (r *repository) Blobs(ctx context.Context) distribution.BlobStore ***REMOVED***
	statter := &blobStatter***REMOVED***
		name:   r.name,
		ub:     r.ub,
		client: r.client,
	***REMOVED***
	return &blobs***REMOVED***
		name:    r.name,
		ub:      r.ub,
		client:  r.client,
		statter: cache.NewCachedBlobStatter(memory.NewInMemoryBlobDescriptorCacheProvider(), statter),
	***REMOVED***
***REMOVED***

func (r *repository) Manifests(ctx context.Context, options ...distribution.ManifestServiceOption) (distribution.ManifestService, error) ***REMOVED***
	// todo(richardscothern): options should be sent over the wire
	return &manifests***REMOVED***
		name:   r.name,
		ub:     r.ub,
		client: r.client,
		etags:  make(map[string]string),
	***REMOVED***, nil
***REMOVED***

func (r *repository) Tags(ctx context.Context) distribution.TagService ***REMOVED***
	return &tags***REMOVED***
		client:  r.client,
		ub:      r.ub,
		context: r.context,
		name:    r.Named(),
	***REMOVED***
***REMOVED***

// tags implements remote tagging operations.
type tags struct ***REMOVED***
	client  *http.Client
	ub      *v2.URLBuilder
	context context.Context
	name    reference.Named
***REMOVED***

// All returns all tags
func (t *tags) All(ctx context.Context) ([]string, error) ***REMOVED***
	var tags []string

	u, err := t.ub.BuildTagsURL(t.name)
	if err != nil ***REMOVED***
		return tags, err
	***REMOVED***

	for ***REMOVED***
		resp, err := t.client.Get(u)
		if err != nil ***REMOVED***
			return tags, err
		***REMOVED***
		defer resp.Body.Close()

		if SuccessStatus(resp.StatusCode) ***REMOVED***
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil ***REMOVED***
				return tags, err
			***REMOVED***

			tagsResponse := struct ***REMOVED***
				Tags []string `json:"tags"`
			***REMOVED******REMOVED******REMOVED***
			if err := json.Unmarshal(b, &tagsResponse); err != nil ***REMOVED***
				return tags, err
			***REMOVED***
			tags = append(tags, tagsResponse.Tags...)
			if link := resp.Header.Get("Link"); link != "" ***REMOVED***
				u = strings.Trim(strings.Split(link, ";")[0], "<>")
			***REMOVED*** else ***REMOVED***
				return tags, nil
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			return tags, HandleErrorResponse(resp)
		***REMOVED***
	***REMOVED***
***REMOVED***

func descriptorFromResponse(response *http.Response) (distribution.Descriptor, error) ***REMOVED***
	desc := distribution.Descriptor***REMOVED******REMOVED***
	headers := response.Header

	ctHeader := headers.Get("Content-Type")
	if ctHeader == "" ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, errors.New("missing or empty Content-Type header")
	***REMOVED***
	desc.MediaType = ctHeader

	digestHeader := headers.Get("Docker-Content-Digest")
	if digestHeader == "" ***REMOVED***
		bytes, err := ioutil.ReadAll(response.Body)
		if err != nil ***REMOVED***
			return distribution.Descriptor***REMOVED******REMOVED***, err
		***REMOVED***
		_, desc, err := distribution.UnmarshalManifest(ctHeader, bytes)
		if err != nil ***REMOVED***
			return distribution.Descriptor***REMOVED******REMOVED***, err
		***REMOVED***
		return desc, nil
	***REMOVED***

	dgst, err := digest.Parse(digestHeader)
	if err != nil ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***
	desc.Digest = dgst

	lengthHeader := headers.Get("Content-Length")
	if lengthHeader == "" ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, errors.New("missing or empty Content-Length header")
	***REMOVED***
	length, err := strconv.ParseInt(lengthHeader, 10, 64)
	if err != nil ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***
	desc.Size = length

	return desc, nil

***REMOVED***

// Get issues a HEAD request for a Manifest against its named endpoint in order
// to construct a descriptor for the tag.  If the registry doesn't support HEADing
// a manifest, fallback to GET.
func (t *tags) Get(ctx context.Context, tag string) (distribution.Descriptor, error) ***REMOVED***
	ref, err := reference.WithTag(t.name, tag)
	if err != nil ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***
	u, err := t.ub.BuildManifestURL(ref)
	if err != nil ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***

	newRequest := func(method string) (*http.Response, error) ***REMOVED***
		req, err := http.NewRequest(method, u, nil)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		for _, t := range distribution.ManifestMediaTypes() ***REMOVED***
			req.Header.Add("Accept", t)
		***REMOVED***
		resp, err := t.client.Do(req)
		return resp, err
	***REMOVED***

	resp, err := newRequest("HEAD")
	if err != nil ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***
	defer resp.Body.Close()

	switch ***REMOVED***
	case resp.StatusCode >= 200 && resp.StatusCode < 400:
		return descriptorFromResponse(resp)
	default:
		// if the response is an error - there will be no body to decode.
		// Issue a GET request:
		//   - for data from a server that does not handle HEAD
		//   - to get error details in case of a failure
		resp, err = newRequest("GET")
		if err != nil ***REMOVED***
			return distribution.Descriptor***REMOVED******REMOVED***, err
		***REMOVED***
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 400 ***REMOVED***
			return descriptorFromResponse(resp)
		***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, HandleErrorResponse(resp)
	***REMOVED***
***REMOVED***

func (t *tags) Lookup(ctx context.Context, digest distribution.Descriptor) ([]string, error) ***REMOVED***
	panic("not implemented")
***REMOVED***

func (t *tags) Tag(ctx context.Context, tag string, desc distribution.Descriptor) error ***REMOVED***
	panic("not implemented")
***REMOVED***

func (t *tags) Untag(ctx context.Context, tag string) error ***REMOVED***
	panic("not implemented")
***REMOVED***

type manifests struct ***REMOVED***
	name   reference.Named
	ub     *v2.URLBuilder
	client *http.Client
	etags  map[string]string
***REMOVED***

func (ms *manifests) Exists(ctx context.Context, dgst digest.Digest) (bool, error) ***REMOVED***
	ref, err := reference.WithDigest(ms.name, dgst)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	u, err := ms.ub.BuildManifestURL(ref)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	resp, err := ms.client.Head(u)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	if SuccessStatus(resp.StatusCode) ***REMOVED***
		return true, nil
	***REMOVED*** else if resp.StatusCode == http.StatusNotFound ***REMOVED***
		return false, nil
	***REMOVED***
	return false, HandleErrorResponse(resp)
***REMOVED***

// AddEtagToTag allows a client to supply an eTag to Get which will be
// used for a conditional HTTP request.  If the eTag matches, a nil manifest
// and ErrManifestNotModified error will be returned. etag is automatically
// quoted when added to this map.
func AddEtagToTag(tag, etag string) distribution.ManifestServiceOption ***REMOVED***
	return etagOption***REMOVED***tag, etag***REMOVED***
***REMOVED***

type etagOption struct***REMOVED*** tag, etag string ***REMOVED***

func (o etagOption) Apply(ms distribution.ManifestService) error ***REMOVED***
	if ms, ok := ms.(*manifests); ok ***REMOVED***
		ms.etags[o.tag] = fmt.Sprintf(`"%s"`, o.etag)
		return nil
	***REMOVED***
	return fmt.Errorf("etag options is a client-only option")
***REMOVED***

// ReturnContentDigest allows a client to set a the content digest on
// a successful request from the 'Docker-Content-Digest' header. This
// returned digest is represents the digest which the registry uses
// to refer to the content and can be used to delete the content.
func ReturnContentDigest(dgst *digest.Digest) distribution.ManifestServiceOption ***REMOVED***
	return contentDigestOption***REMOVED***dgst***REMOVED***
***REMOVED***

type contentDigestOption struct***REMOVED*** digest *digest.Digest ***REMOVED***

func (o contentDigestOption) Apply(ms distribution.ManifestService) error ***REMOVED***
	return nil
***REMOVED***

func (ms *manifests) Get(ctx context.Context, dgst digest.Digest, options ...distribution.ManifestServiceOption) (distribution.Manifest, error) ***REMOVED***
	var (
		digestOrTag string
		ref         reference.Named
		err         error
		contentDgst *digest.Digest
	)

	for _, option := range options ***REMOVED***
		if opt, ok := option.(distribution.WithTagOption); ok ***REMOVED***
			digestOrTag = opt.Tag
			ref, err = reference.WithTag(ms.name, opt.Tag)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED*** else if opt, ok := option.(contentDigestOption); ok ***REMOVED***
			contentDgst = opt.digest
		***REMOVED*** else ***REMOVED***
			err := option.Apply(ms)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if digestOrTag == "" ***REMOVED***
		digestOrTag = dgst.String()
		ref, err = reference.WithDigest(ms.name, dgst)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	u, err := ms.ub.BuildManifestURL(ref)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	req, err := http.NewRequest("GET", u, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	for _, t := range distribution.ManifestMediaTypes() ***REMOVED***
		req.Header.Add("Accept", t)
	***REMOVED***

	if _, ok := ms.etags[digestOrTag]; ok ***REMOVED***
		req.Header.Set("If-None-Match", ms.etags[digestOrTag])
	***REMOVED***

	resp, err := ms.client.Do(req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotModified ***REMOVED***
		return nil, distribution.ErrManifestNotModified
	***REMOVED*** else if SuccessStatus(resp.StatusCode) ***REMOVED***
		if contentDgst != nil ***REMOVED***
			dgst, err := digest.Parse(resp.Header.Get("Docker-Content-Digest"))
			if err == nil ***REMOVED***
				*contentDgst = dgst
			***REMOVED***
		***REMOVED***
		mt := resp.Header.Get("Content-Type")
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		m, _, err := distribution.UnmarshalManifest(mt, body)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return m, nil
	***REMOVED***
	return nil, HandleErrorResponse(resp)
***REMOVED***

// Put puts a manifest.  A tag can be specified using an options parameter which uses some shared state to hold the
// tag name in order to build the correct upload URL.
func (ms *manifests) Put(ctx context.Context, m distribution.Manifest, options ...distribution.ManifestServiceOption) (digest.Digest, error) ***REMOVED***
	ref := ms.name
	var tagged bool

	for _, option := range options ***REMOVED***
		if opt, ok := option.(distribution.WithTagOption); ok ***REMOVED***
			var err error
			ref, err = reference.WithTag(ref, opt.Tag)
			if err != nil ***REMOVED***
				return "", err
			***REMOVED***
			tagged = true
		***REMOVED*** else ***REMOVED***
			err := option.Apply(ms)
			if err != nil ***REMOVED***
				return "", err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	mediaType, p, err := m.Payload()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	if !tagged ***REMOVED***
		// generate a canonical digest and Put by digest
		_, d, err := distribution.UnmarshalManifest(mediaType, p)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		ref, err = reference.WithDigest(ref, d.Digest)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
	***REMOVED***

	manifestURL, err := ms.ub.BuildManifestURL(ref)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	putRequest, err := http.NewRequest("PUT", manifestURL, bytes.NewReader(p))
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	putRequest.Header.Set("Content-Type", mediaType)

	resp, err := ms.client.Do(putRequest)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	defer resp.Body.Close()

	if SuccessStatus(resp.StatusCode) ***REMOVED***
		dgstHeader := resp.Header.Get("Docker-Content-Digest")
		dgst, err := digest.Parse(dgstHeader)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***

		return dgst, nil
	***REMOVED***

	return "", HandleErrorResponse(resp)
***REMOVED***

func (ms *manifests) Delete(ctx context.Context, dgst digest.Digest) error ***REMOVED***
	ref, err := reference.WithDigest(ms.name, dgst)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	u, err := ms.ub.BuildManifestURL(ref)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	req, err := http.NewRequest("DELETE", u, nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	resp, err := ms.client.Do(req)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer resp.Body.Close()

	if SuccessStatus(resp.StatusCode) ***REMOVED***
		return nil
	***REMOVED***
	return HandleErrorResponse(resp)
***REMOVED***

// todo(richardscothern): Restore interface and implementation with merge of #1050
/*func (ms *manifests) Enumerate(ctx context.Context, manifests []distribution.Manifest, last distribution.Manifest) (n int, err error) ***REMOVED***
	panic("not supported")
***REMOVED****/

type blobs struct ***REMOVED***
	name   reference.Named
	ub     *v2.URLBuilder
	client *http.Client

	statter distribution.BlobDescriptorService
	distribution.BlobDeleter
***REMOVED***

func sanitizeLocation(location, base string) (string, error) ***REMOVED***
	baseURL, err := url.Parse(base)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	locationURL, err := url.Parse(location)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return baseURL.ResolveReference(locationURL).String(), nil
***REMOVED***

func (bs *blobs) Stat(ctx context.Context, dgst digest.Digest) (distribution.Descriptor, error) ***REMOVED***
	return bs.statter.Stat(ctx, dgst)

***REMOVED***

func (bs *blobs) Get(ctx context.Context, dgst digest.Digest) ([]byte, error) ***REMOVED***
	reader, err := bs.Open(ctx, dgst)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer reader.Close()

	return ioutil.ReadAll(reader)
***REMOVED***

func (bs *blobs) Open(ctx context.Context, dgst digest.Digest) (distribution.ReadSeekCloser, error) ***REMOVED***
	ref, err := reference.WithDigest(bs.name, dgst)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	blobURL, err := bs.ub.BuildBlobURL(ref)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return transport.NewHTTPReadSeeker(bs.client, blobURL,
		func(resp *http.Response) error ***REMOVED***
			if resp.StatusCode == http.StatusNotFound ***REMOVED***
				return distribution.ErrBlobUnknown
			***REMOVED***
			return HandleErrorResponse(resp)
		***REMOVED***), nil
***REMOVED***

func (bs *blobs) ServeBlob(ctx context.Context, w http.ResponseWriter, r *http.Request, dgst digest.Digest) error ***REMOVED***
	panic("not implemented")
***REMOVED***

func (bs *blobs) Put(ctx context.Context, mediaType string, p []byte) (distribution.Descriptor, error) ***REMOVED***
	writer, err := bs.Create(ctx)
	if err != nil ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***
	dgstr := digest.Canonical.Digester()
	n, err := io.Copy(writer, io.TeeReader(bytes.NewReader(p), dgstr.Hash()))
	if err != nil ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***
	if n < int64(len(p)) ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, fmt.Errorf("short copy: wrote %d of %d", n, len(p))
	***REMOVED***

	desc := distribution.Descriptor***REMOVED***
		MediaType: mediaType,
		Size:      int64(len(p)),
		Digest:    dgstr.Digest(),
	***REMOVED***

	return writer.Commit(ctx, desc)
***REMOVED***

type optionFunc func(interface***REMOVED******REMOVED***) error

func (f optionFunc) Apply(v interface***REMOVED******REMOVED***) error ***REMOVED***
	return f(v)
***REMOVED***

// WithMountFrom returns a BlobCreateOption which designates that the blob should be
// mounted from the given canonical reference.
func WithMountFrom(ref reference.Canonical) distribution.BlobCreateOption ***REMOVED***
	return optionFunc(func(v interface***REMOVED******REMOVED***) error ***REMOVED***
		opts, ok := v.(*distribution.CreateOptions)
		if !ok ***REMOVED***
			return fmt.Errorf("unexpected options type: %T", v)
		***REMOVED***

		opts.Mount.ShouldMount = true
		opts.Mount.From = ref

		return nil
	***REMOVED***)
***REMOVED***

func (bs *blobs) Create(ctx context.Context, options ...distribution.BlobCreateOption) (distribution.BlobWriter, error) ***REMOVED***
	var opts distribution.CreateOptions

	for _, option := range options ***REMOVED***
		err := option.Apply(&opts)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	var values []url.Values

	if opts.Mount.ShouldMount ***REMOVED***
		values = append(values, url.Values***REMOVED***"from": ***REMOVED***opts.Mount.From.Name()***REMOVED***, "mount": ***REMOVED***opts.Mount.From.Digest().String()***REMOVED******REMOVED***)
	***REMOVED***

	u, err := bs.ub.BuildBlobUploadURL(bs.name, values...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	resp, err := bs.client.Post(u, "", nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer resp.Body.Close()

	switch resp.StatusCode ***REMOVED***
	case http.StatusCreated:
		desc, err := bs.statter.Stat(ctx, opts.Mount.From.Digest())
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return nil, distribution.ErrBlobMounted***REMOVED***From: opts.Mount.From, Descriptor: desc***REMOVED***
	case http.StatusAccepted:
		// TODO(dmcgowan): Check for invalid UUID
		uuid := resp.Header.Get("Docker-Upload-UUID")
		location, err := sanitizeLocation(resp.Header.Get("Location"), u)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		return &httpBlobUpload***REMOVED***
			statter:   bs.statter,
			client:    bs.client,
			uuid:      uuid,
			startedAt: time.Now(),
			location:  location,
		***REMOVED***, nil
	default:
		return nil, HandleErrorResponse(resp)
	***REMOVED***
***REMOVED***

func (bs *blobs) Resume(ctx context.Context, id string) (distribution.BlobWriter, error) ***REMOVED***
	panic("not implemented")
***REMOVED***

func (bs *blobs) Delete(ctx context.Context, dgst digest.Digest) error ***REMOVED***
	return bs.statter.Clear(ctx, dgst)
***REMOVED***

type blobStatter struct ***REMOVED***
	name   reference.Named
	ub     *v2.URLBuilder
	client *http.Client
***REMOVED***

func (bs *blobStatter) Stat(ctx context.Context, dgst digest.Digest) (distribution.Descriptor, error) ***REMOVED***
	ref, err := reference.WithDigest(bs.name, dgst)
	if err != nil ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***
	u, err := bs.ub.BuildBlobURL(ref)
	if err != nil ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***

	resp, err := bs.client.Head(u)
	if err != nil ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***
	defer resp.Body.Close()

	if SuccessStatus(resp.StatusCode) ***REMOVED***
		lengthHeader := resp.Header.Get("Content-Length")
		if lengthHeader == "" ***REMOVED***
			return distribution.Descriptor***REMOVED******REMOVED***, fmt.Errorf("missing content-length header for request: %s", u)
		***REMOVED***

		length, err := strconv.ParseInt(lengthHeader, 10, 64)
		if err != nil ***REMOVED***
			return distribution.Descriptor***REMOVED******REMOVED***, fmt.Errorf("error parsing content-length: %v", err)
		***REMOVED***

		return distribution.Descriptor***REMOVED***
			MediaType: resp.Header.Get("Content-Type"),
			Size:      length,
			Digest:    dgst,
		***REMOVED***, nil
	***REMOVED*** else if resp.StatusCode == http.StatusNotFound ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, distribution.ErrBlobUnknown
	***REMOVED***
	return distribution.Descriptor***REMOVED******REMOVED***, HandleErrorResponse(resp)
***REMOVED***

func buildCatalogValues(maxEntries int, last string) url.Values ***REMOVED***
	values := url.Values***REMOVED******REMOVED***

	if maxEntries > 0 ***REMOVED***
		values.Add("n", strconv.Itoa(maxEntries))
	***REMOVED***

	if last != "" ***REMOVED***
		values.Add("last", last)
	***REMOVED***

	return values
***REMOVED***

func (bs *blobStatter) Clear(ctx context.Context, dgst digest.Digest) error ***REMOVED***
	ref, err := reference.WithDigest(bs.name, dgst)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	blobURL, err := bs.ub.BuildBlobURL(ref)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	req, err := http.NewRequest("DELETE", blobURL, nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	resp, err := bs.client.Do(req)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer resp.Body.Close()

	if SuccessStatus(resp.StatusCode) ***REMOVED***
		return nil
	***REMOVED***
	return HandleErrorResponse(resp)
***REMOVED***

func (bs *blobStatter) SetDescriptor(ctx context.Context, dgst digest.Digest, desc distribution.Descriptor) error ***REMOVED***
	return nil
***REMOVED***
