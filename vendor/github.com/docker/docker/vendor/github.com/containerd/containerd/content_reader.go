package containerd

import (
	"context"

	contentapi "github.com/containerd/containerd/api/services/content/v1"
	digest "github.com/opencontainers/go-digest"
)

type remoteReaderAt struct ***REMOVED***
	ctx    context.Context
	digest digest.Digest
	size   int64
	client contentapi.ContentClient
***REMOVED***

func (ra *remoteReaderAt) Size() int64 ***REMOVED***
	return ra.size
***REMOVED***

func (ra *remoteReaderAt) ReadAt(p []byte, off int64) (n int, err error) ***REMOVED***
	rr := &contentapi.ReadContentRequest***REMOVED***
		Digest: ra.digest,
		Offset: off,
		Size_:  int64(len(p)),
	***REMOVED***
	rc, err := ra.client.Read(ra.ctx, rr)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	for len(p) > 0 ***REMOVED***
		var resp *contentapi.ReadContentResponse
		// fill our buffer up until we can fill p.
		resp, err = rc.Recv()
		if err != nil ***REMOVED***
			return n, err
		***REMOVED***

		copied := copy(p, resp.Data)
		n += copied
		p = p[copied:]
	***REMOVED***
	return n, nil
***REMOVED***

func (ra *remoteReaderAt) Close() error ***REMOVED***
	return nil
***REMOVED***
