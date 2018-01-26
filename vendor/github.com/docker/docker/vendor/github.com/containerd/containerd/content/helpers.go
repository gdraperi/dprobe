package content

import (
	"context"
	"io"
	"sync"

	"github.com/containerd/containerd/errdefs"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

var bufPool = sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED***
		buffer := make([]byte, 1<<20)
		return &buffer
	***REMOVED***,
***REMOVED***

// NewReader returns a io.Reader from a ReaderAt
func NewReader(ra ReaderAt) io.Reader ***REMOVED***
	rd := io.NewSectionReader(ra, 0, ra.Size())
	return rd
***REMOVED***

// ReadBlob retrieves the entire contents of the blob from the provider.
//
// Avoid using this for large blobs, such as layers.
func ReadBlob(ctx context.Context, provider Provider, dgst digest.Digest) ([]byte, error) ***REMOVED***
	ra, err := provider.ReaderAt(ctx, dgst)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer ra.Close()

	p := make([]byte, ra.Size())

	_, err = ra.ReadAt(p, 0)
	return p, err
***REMOVED***

// WriteBlob writes data with the expected digest into the content store. If
// expected already exists, the method returns immediately and the reader will
// not be consumed.
//
// This is useful when the digest and size are known beforehand.
//
// Copy is buffered, so no need to wrap reader in buffered io.
func WriteBlob(ctx context.Context, cs Ingester, ref string, r io.Reader, size int64, expected digest.Digest, opts ...Opt) error ***REMOVED***
	cw, err := cs.Writer(ctx, ref, size, expected)
	if err != nil ***REMOVED***
		if !errdefs.IsAlreadyExists(err) ***REMOVED***
			return err
		***REMOVED***

		return nil // all ready present
	***REMOVED***
	defer cw.Close()

	return Copy(ctx, cw, r, size, expected, opts...)
***REMOVED***

// Copy copies data with the expected digest from the reader into the
// provided content store writer.
//
// This is useful when the digest and size are known beforehand. When
// the size or digest is unknown, these values may be empty.
//
// Copy is buffered, so no need to wrap reader in buffered io.
func Copy(ctx context.Context, cw Writer, r io.Reader, size int64, expected digest.Digest, opts ...Opt) error ***REMOVED***
	ws, err := cw.Status()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if ws.Offset > 0 ***REMOVED***
		r, err = seekReader(r, ws.Offset, size)
		if err != nil ***REMOVED***
			if !isUnseekable(err) ***REMOVED***
				return errors.Wrapf(err, "unable to resume write to %v", ws.Ref)
			***REMOVED***

			// reader is unseekable, try to move the writer back to the start.
			if err := cw.Truncate(0); err != nil ***REMOVED***
				return errors.Wrapf(err, "content writer truncate failed")
			***REMOVED***
		***REMOVED***
	***REMOVED***

	buf := bufPool.Get().(*[]byte)
	defer bufPool.Put(buf)

	if _, err := io.CopyBuffer(cw, r, *buf); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := cw.Commit(ctx, size, expected, opts...); err != nil ***REMOVED***
		if !errdefs.IsAlreadyExists(err) ***REMOVED***
			return errors.Wrapf(err, "failed commit on ref %q", ws.Ref)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

var errUnseekable = errors.New("seek not supported")

func isUnseekable(err error) bool ***REMOVED***
	return errors.Cause(err) == errUnseekable
***REMOVED***

// seekReader attempts to seek the reader to the given offset, either by
// resolving `io.Seeker` or by detecting `io.ReaderAt`.
func seekReader(r io.Reader, offset, size int64) (io.Reader, error) ***REMOVED***
	// attempt to resolve r as a seeker and setup the offset.
	seeker, ok := r.(io.Seeker)
	if ok ***REMOVED***
		nn, err := seeker.Seek(offset, io.SeekStart)
		if nn != offset ***REMOVED***
			return nil, errors.Wrapf(err, "failed to seek to offset %v", offset)
		***REMOVED***

		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		return r, nil
	***REMOVED***

	// ok, let's try io.ReaderAt!
	readerAt, ok := r.(io.ReaderAt)
	if ok && size > offset ***REMOVED***
		sr := io.NewSectionReader(readerAt, offset, size)
		return sr, nil
	***REMOVED***

	return r, errors.Wrapf(errUnseekable, "seek to offset %v failed", offset)
***REMOVED***
