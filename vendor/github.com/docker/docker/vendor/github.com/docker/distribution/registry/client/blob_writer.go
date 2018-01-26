package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
)

type httpBlobUpload struct ***REMOVED***
	statter distribution.BlobStatter
	client  *http.Client

	uuid      string
	startedAt time.Time

	location string // always the last value of the location header.
	offset   int64
	closed   bool
***REMOVED***

func (hbu *httpBlobUpload) Reader() (io.ReadCloser, error) ***REMOVED***
	panic("Not implemented")
***REMOVED***

func (hbu *httpBlobUpload) handleErrorResponse(resp *http.Response) error ***REMOVED***
	if resp.StatusCode == http.StatusNotFound ***REMOVED***
		return distribution.ErrBlobUploadUnknown
	***REMOVED***
	return HandleErrorResponse(resp)
***REMOVED***

func (hbu *httpBlobUpload) ReadFrom(r io.Reader) (n int64, err error) ***REMOVED***
	req, err := http.NewRequest("PATCH", hbu.location, ioutil.NopCloser(r))
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	defer req.Body.Close()

	resp, err := hbu.client.Do(req)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	if !SuccessStatus(resp.StatusCode) ***REMOVED***
		return 0, hbu.handleErrorResponse(resp)
	***REMOVED***

	hbu.uuid = resp.Header.Get("Docker-Upload-UUID")
	hbu.location, err = sanitizeLocation(resp.Header.Get("Location"), hbu.location)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	rng := resp.Header.Get("Range")
	var start, end int64
	if n, err := fmt.Sscanf(rng, "%d-%d", &start, &end); err != nil ***REMOVED***
		return 0, err
	***REMOVED*** else if n != 2 || end < start ***REMOVED***
		return 0, fmt.Errorf("bad range format: %s", rng)
	***REMOVED***

	return (end - start + 1), nil

***REMOVED***

func (hbu *httpBlobUpload) Write(p []byte) (n int, err error) ***REMOVED***
	req, err := http.NewRequest("PATCH", hbu.location, bytes.NewReader(p))
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	req.Header.Set("Content-Range", fmt.Sprintf("%d-%d", hbu.offset, hbu.offset+int64(len(p)-1)))
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(p)))
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := hbu.client.Do(req)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	if !SuccessStatus(resp.StatusCode) ***REMOVED***
		return 0, hbu.handleErrorResponse(resp)
	***REMOVED***

	hbu.uuid = resp.Header.Get("Docker-Upload-UUID")
	hbu.location, err = sanitizeLocation(resp.Header.Get("Location"), hbu.location)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	rng := resp.Header.Get("Range")
	var start, end int
	if n, err := fmt.Sscanf(rng, "%d-%d", &start, &end); err != nil ***REMOVED***
		return 0, err
	***REMOVED*** else if n != 2 || end < start ***REMOVED***
		return 0, fmt.Errorf("bad range format: %s", rng)
	***REMOVED***

	return (end - start + 1), nil

***REMOVED***

func (hbu *httpBlobUpload) Size() int64 ***REMOVED***
	return hbu.offset
***REMOVED***

func (hbu *httpBlobUpload) ID() string ***REMOVED***
	return hbu.uuid
***REMOVED***

func (hbu *httpBlobUpload) StartedAt() time.Time ***REMOVED***
	return hbu.startedAt
***REMOVED***

func (hbu *httpBlobUpload) Commit(ctx context.Context, desc distribution.Descriptor) (distribution.Descriptor, error) ***REMOVED***
	// TODO(dmcgowan): Check if already finished, if so just fetch
	req, err := http.NewRequest("PUT", hbu.location, nil)
	if err != nil ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***

	values := req.URL.Query()
	values.Set("digest", desc.Digest.String())
	req.URL.RawQuery = values.Encode()

	resp, err := hbu.client.Do(req)
	if err != nil ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***
	defer resp.Body.Close()

	if !SuccessStatus(resp.StatusCode) ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, hbu.handleErrorResponse(resp)
	***REMOVED***

	return hbu.statter.Stat(ctx, desc.Digest)
***REMOVED***

func (hbu *httpBlobUpload) Cancel(ctx context.Context) error ***REMOVED***
	req, err := http.NewRequest("DELETE", hbu.location, nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	resp, err := hbu.client.Do(req)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound || SuccessStatus(resp.StatusCode) ***REMOVED***
		return nil
	***REMOVED***
	return hbu.handleErrorResponse(resp)
***REMOVED***

func (hbu *httpBlobUpload) Close() error ***REMOVED***
	hbu.closed = true
	return nil
***REMOVED***
