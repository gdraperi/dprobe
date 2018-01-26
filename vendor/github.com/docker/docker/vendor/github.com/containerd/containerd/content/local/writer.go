package local

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/errdefs"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

// writer represents a write transaction against the blob store.
type writer struct ***REMOVED***
	s         *store
	fp        *os.File // opened data file
	path      string   // path to writer dir
	ref       string   // ref key
	offset    int64
	total     int64
	digester  digest.Digester
	startedAt time.Time
	updatedAt time.Time
***REMOVED***

func (w *writer) Status() (content.Status, error) ***REMOVED***
	return content.Status***REMOVED***
		Ref:       w.ref,
		Offset:    w.offset,
		Total:     w.total,
		StartedAt: w.startedAt,
		UpdatedAt: w.updatedAt,
	***REMOVED***, nil
***REMOVED***

// Digest returns the current digest of the content, up to the current write.
//
// Cannot be called concurrently with `Write`.
func (w *writer) Digest() digest.Digest ***REMOVED***
	return w.digester.Digest()
***REMOVED***

// Write p to the transaction.
//
// Note that writes are unbuffered to the backing file. When writing, it is
// recommended to wrap in a bufio.Writer or, preferably, use io.CopyBuffer.
func (w *writer) Write(p []byte) (n int, err error) ***REMOVED***
	n, err = w.fp.Write(p)
	w.digester.Hash().Write(p[:n])
	w.offset += int64(len(p))
	w.updatedAt = time.Now()
	return n, err
***REMOVED***

func (w *writer) Commit(ctx context.Context, size int64, expected digest.Digest, opts ...content.Opt) error ***REMOVED***
	var base content.Info
	for _, opt := range opts ***REMOVED***
		if err := opt(&base); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if w.fp == nil ***REMOVED***
		return errors.Wrap(errdefs.ErrFailedPrecondition, "cannot commit on closed writer")
	***REMOVED***

	if err := w.fp.Sync(); err != nil ***REMOVED***
		return errors.Wrap(err, "sync failed")
	***REMOVED***

	fi, err := w.fp.Stat()
	if err != nil ***REMOVED***
		return errors.Wrap(err, "stat on ingest file failed")
	***REMOVED***

	// change to readonly, more important for read, but provides _some_
	// protection from this point on. We use the existing perms with a mask
	// only allowing reads honoring the umask on creation.
	//
	// This removes write and exec, only allowing read per the creation umask.
	//
	// NOTE: Windows does not support this operation
	if runtime.GOOS != "windows" ***REMOVED***
		if err := w.fp.Chmod((fi.Mode() & os.ModePerm) &^ 0333); err != nil ***REMOVED***
			return errors.Wrap(err, "failed to change ingest file permissions")
		***REMOVED***
	***REMOVED***

	if size > 0 && size != fi.Size() ***REMOVED***
		return errors.Errorf("unexpected commit size %d, expected %d", fi.Size(), size)
	***REMOVED***

	if err := w.fp.Close(); err != nil ***REMOVED***
		return errors.Wrap(err, "failed closing ingest")
	***REMOVED***

	dgst := w.digester.Digest()
	if expected != "" && expected != dgst ***REMOVED***
		return errors.Errorf("unexpected commit digest %s, expected %s", dgst, expected)
	***REMOVED***

	var (
		ingest = filepath.Join(w.path, "data")
		target = w.s.blobPath(dgst)
	)

	// make sure parent directories of blob exist
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil ***REMOVED***
		return err
	***REMOVED***

	// clean up!!
	defer os.RemoveAll(w.path)

	if err := os.Rename(ingest, target); err != nil ***REMOVED***
		if os.IsExist(err) ***REMOVED***
			// collision with the target file!
			return errors.Wrapf(errdefs.ErrAlreadyExists, "content %v", dgst)
		***REMOVED***
		return err
	***REMOVED***
	commitTime := time.Now()
	if err := os.Chtimes(target, commitTime, commitTime); err != nil ***REMOVED***
		return err
	***REMOVED***

	w.fp = nil
	unlock(w.ref)

	if w.s.ls != nil && base.Labels != nil ***REMOVED***
		if err := w.s.ls.Set(dgst, base.Labels); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// Close the writer, flushing any unwritten data and leaving the progress in
// tact.
//
// If one needs to resume the transaction, a new writer can be obtained from
// `Ingester.Writer` using the same key. The write can then be continued
// from it was left off.
//
// To abandon a transaction completely, first call close then `IngestManager.Abort` to
// clean up the associated resources.
func (w *writer) Close() (err error) ***REMOVED***
	if w.fp != nil ***REMOVED***
		w.fp.Sync()
		err = w.fp.Close()
		writeTimestampFile(filepath.Join(w.path, "updatedat"), w.updatedAt)
		w.fp = nil
		unlock(w.ref)
		return
	***REMOVED***

	return nil
***REMOVED***

func (w *writer) Truncate(size int64) error ***REMOVED***
	if size != 0 ***REMOVED***
		return errors.New("Truncate: unsupported size")
	***REMOVED***
	w.offset = 0
	w.digester.Hash().Reset()
	if _, err := w.fp.Seek(0, io.SeekStart); err != nil ***REMOVED***
		return err
	***REMOVED***
	return w.fp.Truncate(0)
***REMOVED***
