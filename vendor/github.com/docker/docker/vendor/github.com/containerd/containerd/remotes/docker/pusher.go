package docker

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/remotes"
	digest "github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

type dockerPusher struct ***REMOVED***
	*dockerBase
	tag string

	// TODO: namespace tracker
	tracker StatusTracker
***REMOVED***

func (p dockerPusher) Push(ctx context.Context, desc ocispec.Descriptor) (content.Writer, error) ***REMOVED***
	ctx, err := contextWithRepositoryScope(ctx, p.refspec, true)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ref := remotes.MakeRefKey(ctx, desc)
	status, err := p.tracker.GetStatus(ref)
	if err == nil ***REMOVED***
		if status.Offset == status.Total ***REMOVED***
			return nil, errors.Wrapf(errdefs.ErrAlreadyExists, "ref %v", ref)
		***REMOVED***
		// TODO: Handle incomplete status
	***REMOVED*** else if !errdefs.IsNotFound(err) ***REMOVED***
		return nil, errors.Wrap(err, "failed to get status")
	***REMOVED***

	var (
		isManifest bool
		existCheck string
	)

	switch desc.MediaType ***REMOVED***
	case images.MediaTypeDockerSchema2Manifest, images.MediaTypeDockerSchema2ManifestList,
		ocispec.MediaTypeImageManifest, ocispec.MediaTypeImageIndex:
		isManifest = true
		if p.tag == "" ***REMOVED***
			existCheck = path.Join("manifests", desc.Digest.String())
		***REMOVED*** else ***REMOVED***
			existCheck = path.Join("manifests", p.tag)
		***REMOVED***
	default:
		existCheck = path.Join("blobs", desc.Digest.String())
	***REMOVED***

	req, err := http.NewRequest(http.MethodHead, p.url(existCheck), nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	req.Header.Set("Accept", strings.Join([]string***REMOVED***desc.MediaType, `*`***REMOVED***, ", "))
	resp, err := p.doRequestWithRetries(ctx, req, nil)
	if err != nil ***REMOVED***
		if errors.Cause(err) != ErrInvalidAuthorization ***REMOVED***
			return nil, err
		***REMOVED***
		log.G(ctx).WithError(err).Debugf("Unable to check existence, continuing with push")
	***REMOVED*** else ***REMOVED***
		if resp.StatusCode == http.StatusOK ***REMOVED***
			var exists bool
			if isManifest && p.tag != "" ***REMOVED***
				dgstHeader := digest.Digest(resp.Header.Get("Docker-Content-Digest"))
				if dgstHeader == desc.Digest ***REMOVED***
					exists = true
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				exists = true
			***REMOVED***

			if exists ***REMOVED***
				p.tracker.SetStatus(ref, Status***REMOVED***
					Status: content.Status***REMOVED***
						Ref: ref,
						// TODO: Set updated time?
					***REMOVED***,
				***REMOVED***)
				return nil, errors.Wrapf(errdefs.ErrAlreadyExists, "content %v on remote", desc.Digest)
			***REMOVED***
		***REMOVED*** else if resp.StatusCode != http.StatusNotFound ***REMOVED***
			// TODO: log error
			return nil, errors.Errorf("unexpected response: %s", resp.Status)
		***REMOVED***
	***REMOVED***

	// TODO: Lookup related objects for cross repository push

	if isManifest ***REMOVED***
		var putPath string
		if p.tag != "" ***REMOVED***
			putPath = path.Join("manifests", p.tag)
		***REMOVED*** else ***REMOVED***
			putPath = path.Join("manifests", desc.Digest.String())
		***REMOVED***

		req, err = http.NewRequest(http.MethodPut, p.url(putPath), nil)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		req.Header.Add("Content-Type", desc.MediaType)
	***REMOVED*** else ***REMOVED***
		// TODO: Do monolithic upload if size is small

		// Start upload request
		req, err = http.NewRequest(http.MethodPost, p.url("blobs", "uploads")+"/", nil)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		resp, err := p.doRequestWithRetries(ctx, req, nil)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		switch resp.StatusCode ***REMOVED***
		case http.StatusOK, http.StatusAccepted, http.StatusNoContent:
		default:
			// TODO: log error
			return nil, errors.Errorf("unexpected response: %s", resp.Status)
		***REMOVED***

		location := resp.Header.Get("Location")
		// Support paths without host in location
		if strings.HasPrefix(location, "/") ***REMOVED***
			u := p.base
			u.Path = location
			location = u.String()
		***REMOVED***

		req, err = http.NewRequest(http.MethodPut, location, nil)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		q := req.URL.Query()
		q.Add("digest", desc.Digest.String())
		req.URL.RawQuery = q.Encode()

	***REMOVED***
	p.tracker.SetStatus(ref, Status***REMOVED***
		Status: content.Status***REMOVED***
			Ref:       ref,
			Total:     desc.Size,
			Expected:  desc.Digest,
			StartedAt: time.Now(),
		***REMOVED***,
	***REMOVED***)

	// TODO: Support chunked upload

	pr, pw := io.Pipe()
	respC := make(chan *http.Response, 1)

	req.Body = ioutil.NopCloser(pr)
	req.ContentLength = desc.Size

	go func() ***REMOVED***
		defer close(respC)
		resp, err = p.doRequest(ctx, req)
		if err != nil ***REMOVED***
			pr.CloseWithError(err)
			return
		***REMOVED***

		switch resp.StatusCode ***REMOVED***
		case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		default:
			// TODO: log error
			pr.CloseWithError(errors.Errorf("unexpected response: %s", resp.Status))
		***REMOVED***
		respC <- resp
	***REMOVED***()

	return &pushWriter***REMOVED***
		base:       p.dockerBase,
		ref:        ref,
		pipe:       pw,
		responseC:  respC,
		isManifest: isManifest,
		expected:   desc.Digest,
		tracker:    p.tracker,
	***REMOVED***, nil
***REMOVED***

type pushWriter struct ***REMOVED***
	base *dockerBase
	ref  string

	pipe       *io.PipeWriter
	responseC  <-chan *http.Response
	isManifest bool

	expected digest.Digest
	tracker  StatusTracker
***REMOVED***

func (pw *pushWriter) Write(p []byte) (n int, err error) ***REMOVED***
	status, err := pw.tracker.GetStatus(pw.ref)
	if err != nil ***REMOVED***
		return n, err
	***REMOVED***
	n, err = pw.pipe.Write(p)
	status.Offset += int64(n)
	status.UpdatedAt = time.Now()
	pw.tracker.SetStatus(pw.ref, status)
	return
***REMOVED***

func (pw *pushWriter) Close() error ***REMOVED***
	return pw.pipe.Close()
***REMOVED***

func (pw *pushWriter) Status() (content.Status, error) ***REMOVED***
	status, err := pw.tracker.GetStatus(pw.ref)
	if err != nil ***REMOVED***
		return content.Status***REMOVED******REMOVED***, err
	***REMOVED***
	return status.Status, nil

***REMOVED***

func (pw *pushWriter) Digest() digest.Digest ***REMOVED***
	// TODO: Get rid of this function?
	return pw.expected
***REMOVED***

func (pw *pushWriter) Commit(ctx context.Context, size int64, expected digest.Digest, opts ...content.Opt) error ***REMOVED***
	// Check whether read has already thrown an error
	if _, err := pw.pipe.Write([]byte***REMOVED******REMOVED***); err != nil && err != io.ErrClosedPipe ***REMOVED***
		return errors.Wrap(err, "pipe error before commit")
	***REMOVED***

	if err := pw.pipe.Close(); err != nil ***REMOVED***
		return err
	***REMOVED***
	// TODO: Update status to determine committing

	// TODO: timeout waiting for response
	resp := <-pw.responseC
	if resp == nil ***REMOVED***
		return errors.New("no response")
	***REMOVED***

	// 201 is specified return status, some registries return
	// 200 or 204.
	switch resp.StatusCode ***REMOVED***
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
	default:
		return errors.Errorf("unexpected status: %s", resp.Status)
	***REMOVED***

	status, err := pw.tracker.GetStatus(pw.ref)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed to get status")
	***REMOVED***

	if size > 0 && size != status.Offset ***REMOVED***
		return errors.Errorf("unxpected size %d, expected %d", status.Offset, size)
	***REMOVED***

	if expected == "" ***REMOVED***
		expected = status.Expected
	***REMOVED***

	actual, err := digest.Parse(resp.Header.Get("Docker-Content-Digest"))
	if err != nil ***REMOVED***
		return errors.Wrap(err, "invalid content digest in response")
	***REMOVED***

	if actual != expected ***REMOVED***
		return errors.Errorf("got digest %s, expected %s", actual, expected)
	***REMOVED***

	return nil
***REMOVED***

func (pw *pushWriter) Truncate(size int64) error ***REMOVED***
	// TODO: if blob close request and start new request at offset
	// TODO: always error on manifest
	return errors.New("cannot truncate remote upload")
***REMOVED***
