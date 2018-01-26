package docker

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/log"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type dockerFetcher struct ***REMOVED***
	*dockerBase
***REMOVED***

func (r dockerFetcher) Fetch(ctx context.Context, desc ocispec.Descriptor) (io.ReadCloser, error) ***REMOVED***
	ctx = log.WithLogger(ctx, log.G(ctx).WithFields(
		logrus.Fields***REMOVED***
			"base":   r.base.String(),
			"digest": desc.Digest,
		***REMOVED***,
	))

	urls, err := r.getV2URLPaths(ctx, desc)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ctx, err = contextWithRepositoryScope(ctx, r.refspec, false)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return newHTTPReadSeeker(desc.Size, func(offset int64) (io.ReadCloser, error) ***REMOVED***
		for _, u := range urls ***REMOVED***
			rc, err := r.open(ctx, u, desc.MediaType, offset)
			if err != nil ***REMOVED***
				if errdefs.IsNotFound(err) ***REMOVED***
					continue // try one of the other urls.
				***REMOVED***

				return nil, err
			***REMOVED***

			return rc, nil
		***REMOVED***

		return nil, errors.Wrapf(errdefs.ErrNotFound,
			"could not fetch content descriptor %v (%v) from remote",
			desc.Digest, desc.MediaType)

	***REMOVED***)
***REMOVED***

func (r dockerFetcher) open(ctx context.Context, u, mediatype string, offset int64) (io.ReadCloser, error) ***REMOVED***
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	req.Header.Set("Accept", strings.Join([]string***REMOVED***mediatype, `*`***REMOVED***, ", "))

	if offset > 0 ***REMOVED***
		// TODO(stevvooe): Only set this header in response to the
		// "Accept-Ranges: bytes" header.
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", offset))
	***REMOVED***

	resp, err := r.doRequestWithRetries(ctx, req, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if resp.StatusCode > 299 ***REMOVED***
		// TODO(stevvooe): When doing a offset specific request, we should
		// really distinguish between a 206 and a 200. In the case of 200, we
		// can discard the bytes, hiding the seek behavior from the
		// implementation.

		resp.Body.Close()
		if resp.StatusCode == http.StatusNotFound ***REMOVED***
			return nil, errors.Wrapf(errdefs.ErrNotFound, "content at %v not found", u)
		***REMOVED***
		return nil, errors.Errorf("unexpected status code %v: %v", u, resp.Status)
	***REMOVED***

	return resp.Body, nil
***REMOVED***

// getV2URLPaths generates the candidate urls paths for the object based on the
// set of hints and the provided object id. URLs are returned in the order of
// most to least likely succeed.
func (r *dockerFetcher) getV2URLPaths(ctx context.Context, desc ocispec.Descriptor) ([]string, error) ***REMOVED***
	var urls []string

	if len(desc.URLs) > 0 ***REMOVED***
		// handle fetch via external urls.
		for _, u := range desc.URLs ***REMOVED***
			log.G(ctx).WithField("url", u).Debug("adding alternative url")
			urls = append(urls, u)
		***REMOVED***
	***REMOVED***

	switch desc.MediaType ***REMOVED***
	case images.MediaTypeDockerSchema2Manifest, images.MediaTypeDockerSchema2ManifestList,
		images.MediaTypeDockerSchema1Manifest,
		ocispec.MediaTypeImageManifest, ocispec.MediaTypeImageIndex:
		urls = append(urls, r.url(path.Join("manifests", desc.Digest.String())))
	***REMOVED***

	// always fallback to attempting to get the object out of the blobs store.
	urls = append(urls, r.url(path.Join("blobs", desc.Digest.String())))

	return urls, nil
***REMOVED***
