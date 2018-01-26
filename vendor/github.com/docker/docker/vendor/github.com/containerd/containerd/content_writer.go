package containerd

import (
	"context"
	"io"

	contentapi "github.com/containerd/containerd/api/services/content/v1"
	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/errdefs"
	digest "github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

type remoteWriter struct ***REMOVED***
	ref    string
	client contentapi.Content_WriteClient
	offset int64
	digest digest.Digest
***REMOVED***

// send performs a synchronous req-resp cycle on the client.
func (rw *remoteWriter) send(req *contentapi.WriteContentRequest) (*contentapi.WriteContentResponse, error) ***REMOVED***
	if err := rw.client.Send(req); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	resp, err := rw.client.Recv()

	if err == nil ***REMOVED***
		// try to keep these in sync
		if resp.Digest != "" ***REMOVED***
			rw.digest = resp.Digest
		***REMOVED***
	***REMOVED***

	return resp, err
***REMOVED***

func (rw *remoteWriter) Status() (content.Status, error) ***REMOVED***
	resp, err := rw.send(&contentapi.WriteContentRequest***REMOVED***
		Action: contentapi.WriteActionStat,
	***REMOVED***)
	if err != nil ***REMOVED***
		return content.Status***REMOVED******REMOVED***, errors.Wrap(err, "error getting writer status")
	***REMOVED***

	return content.Status***REMOVED***
		Ref:       rw.ref,
		Offset:    resp.Offset,
		Total:     resp.Total,
		StartedAt: resp.StartedAt,
		UpdatedAt: resp.UpdatedAt,
	***REMOVED***, nil
***REMOVED***

func (rw *remoteWriter) Digest() digest.Digest ***REMOVED***
	return rw.digest
***REMOVED***

func (rw *remoteWriter) Write(p []byte) (n int, err error) ***REMOVED***
	offset := rw.offset

	resp, err := rw.send(&contentapi.WriteContentRequest***REMOVED***
		Action: contentapi.WriteActionWrite,
		Offset: offset,
		Data:   p,
	***REMOVED***)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	n = int(resp.Offset - offset)
	if n < len(p) ***REMOVED***
		err = io.ErrShortWrite
	***REMOVED***

	rw.offset += int64(n)
	if resp.Digest != "" ***REMOVED***
		rw.digest = resp.Digest
	***REMOVED***
	return
***REMOVED***

func (rw *remoteWriter) Commit(ctx context.Context, size int64, expected digest.Digest, opts ...content.Opt) error ***REMOVED***
	var base content.Info
	for _, opt := range opts ***REMOVED***
		if err := opt(&base); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	resp, err := rw.send(&contentapi.WriteContentRequest***REMOVED***
		Action:   contentapi.WriteActionCommit,
		Total:    size,
		Offset:   rw.offset,
		Expected: expected,
		Labels:   base.Labels,
	***REMOVED***)
	if err != nil ***REMOVED***
		return errdefs.FromGRPC(err)
	***REMOVED***

	if size != 0 && resp.Offset != size ***REMOVED***
		return errors.Errorf("unexpected size: %v != %v", resp.Offset, size)
	***REMOVED***

	if expected != "" && resp.Digest != expected ***REMOVED***
		return errors.Errorf("unexpected digest: %v != %v", resp.Digest, expected)
	***REMOVED***

	rw.digest = resp.Digest
	rw.offset = resp.Offset
	return nil
***REMOVED***

func (rw *remoteWriter) Truncate(size int64) error ***REMOVED***
	// This truncation won't actually be validated until a write is issued.
	rw.offset = size
	return nil
***REMOVED***

func (rw *remoteWriter) Close() error ***REMOVED***
	return rw.client.CloseSend()
***REMOVED***
