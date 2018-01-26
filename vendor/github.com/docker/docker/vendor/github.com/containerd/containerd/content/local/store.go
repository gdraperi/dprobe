package local

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/filters"
	"github.com/containerd/containerd/log"
	digest "github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

var bufPool = sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED***
		buffer := make([]byte, 1<<20)
		return &buffer
	***REMOVED***,
***REMOVED***

// LabelStore is used to store mutable labels for digests
type LabelStore interface ***REMOVED***
	// Get returns all the labels for the given digest
	Get(digest.Digest) (map[string]string, error)

	// Set sets all the labels for a given digest
	Set(digest.Digest, map[string]string) error

	// Update replaces the given labels for a digest,
	// a key with an empty value removes a label.
	Update(digest.Digest, map[string]string) (map[string]string, error)
***REMOVED***

// Store is digest-keyed store for content. All data written into the store is
// stored under a verifiable digest.
//
// Store can generally support multi-reader, single-writer ingest of data,
// including resumable ingest.
type store struct ***REMOVED***
	root string
	ls   LabelStore
***REMOVED***

// NewStore returns a local content store
func NewStore(root string) (content.Store, error) ***REMOVED***
	return NewLabeledStore(root, nil)
***REMOVED***

// NewLabeledStore returns a new content store using the provided label store
//
// Note: content stores which are used underneath a metadata store may not
// require labels and should use `NewStore`. `NewLabeledStore` is primarily
// useful for tests or standalone implementations.
func NewLabeledStore(root string, ls LabelStore) (content.Store, error) ***REMOVED***
	if err := os.MkdirAll(filepath.Join(root, "ingest"), 0777); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &store***REMOVED***
		root: root,
		ls:   ls,
	***REMOVED***, nil
***REMOVED***

func (s *store) Info(ctx context.Context, dgst digest.Digest) (content.Info, error) ***REMOVED***
	p := s.blobPath(dgst)
	fi, err := os.Stat(p)
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			err = errors.Wrapf(errdefs.ErrNotFound, "content %v", dgst)
		***REMOVED***

		return content.Info***REMOVED******REMOVED***, err
	***REMOVED***
	var labels map[string]string
	if s.ls != nil ***REMOVED***
		labels, err = s.ls.Get(dgst)
		if err != nil ***REMOVED***
			return content.Info***REMOVED******REMOVED***, err
		***REMOVED***
	***REMOVED***
	return s.info(dgst, fi, labels), nil
***REMOVED***

func (s *store) info(dgst digest.Digest, fi os.FileInfo, labels map[string]string) content.Info ***REMOVED***
	return content.Info***REMOVED***
		Digest:    dgst,
		Size:      fi.Size(),
		CreatedAt: fi.ModTime(),
		UpdatedAt: getATime(fi),
		Labels:    labels,
	***REMOVED***
***REMOVED***

// ReaderAt returns an io.ReaderAt for the blob.
func (s *store) ReaderAt(ctx context.Context, dgst digest.Digest) (content.ReaderAt, error) ***REMOVED***
	p := s.blobPath(dgst)
	fi, err := os.Stat(p)
	if err != nil ***REMOVED***
		if !os.IsNotExist(err) ***REMOVED***
			return nil, err
		***REMOVED***

		return nil, errors.Wrapf(errdefs.ErrNotFound, "blob %s expected at %s", dgst, p)
	***REMOVED***

	fp, err := os.Open(p)
	if err != nil ***REMOVED***
		if !os.IsNotExist(err) ***REMOVED***
			return nil, err
		***REMOVED***

		return nil, errors.Wrapf(errdefs.ErrNotFound, "blob %s expected at %s", dgst, p)
	***REMOVED***

	return sizeReaderAt***REMOVED***size: fi.Size(), fp: fp***REMOVED***, nil
***REMOVED***

// Delete removes a blob by its digest.
//
// While this is safe to do concurrently, safe exist-removal logic must hold
// some global lock on the store.
func (s *store) Delete(ctx context.Context, dgst digest.Digest) error ***REMOVED***
	if err := os.RemoveAll(s.blobPath(dgst)); err != nil ***REMOVED***
		if !os.IsNotExist(err) ***REMOVED***
			return err
		***REMOVED***

		return errors.Wrapf(errdefs.ErrNotFound, "content %v", dgst)
	***REMOVED***

	return nil
***REMOVED***

func (s *store) Update(ctx context.Context, info content.Info, fieldpaths ...string) (content.Info, error) ***REMOVED***
	if s.ls == nil ***REMOVED***
		return content.Info***REMOVED******REMOVED***, errors.Wrapf(errdefs.ErrFailedPrecondition, "update not supported on immutable content store")
	***REMOVED***

	p := s.blobPath(info.Digest)
	fi, err := os.Stat(p)
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			err = errors.Wrapf(errdefs.ErrNotFound, "content %v", info.Digest)
		***REMOVED***

		return content.Info***REMOVED******REMOVED***, err
	***REMOVED***

	var (
		all    bool
		labels map[string]string
	)
	if len(fieldpaths) > 0 ***REMOVED***
		for _, path := range fieldpaths ***REMOVED***
			if strings.HasPrefix(path, "labels.") ***REMOVED***
				if labels == nil ***REMOVED***
					labels = map[string]string***REMOVED******REMOVED***
				***REMOVED***

				key := strings.TrimPrefix(path, "labels.")
				labels[key] = info.Labels[key]
				continue
			***REMOVED***

			switch path ***REMOVED***
			case "labels":
				all = true
				labels = info.Labels
			default:
				return content.Info***REMOVED******REMOVED***, errors.Wrapf(errdefs.ErrInvalidArgument, "cannot update %q field on content info %q", path, info.Digest)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		all = true
		labels = info.Labels
	***REMOVED***

	if all ***REMOVED***
		err = s.ls.Set(info.Digest, labels)
	***REMOVED*** else ***REMOVED***
		labels, err = s.ls.Update(info.Digest, labels)
	***REMOVED***
	if err != nil ***REMOVED***
		return content.Info***REMOVED******REMOVED***, err
	***REMOVED***

	info = s.info(info.Digest, fi, labels)
	info.UpdatedAt = time.Now()

	if err := os.Chtimes(p, info.UpdatedAt, info.CreatedAt); err != nil ***REMOVED***
		log.G(ctx).WithError(err).Warnf("could not change access time for %s", info.Digest)
	***REMOVED***

	return info, nil
***REMOVED***

func (s *store) Walk(ctx context.Context, fn content.WalkFunc, filters ...string) error ***REMOVED***
	// TODO: Support filters
	root := filepath.Join(s.root, "blobs")
	var alg digest.Algorithm
	return filepath.Walk(root, func(path string, fi os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if !fi.IsDir() && !alg.Available() ***REMOVED***
			return nil
		***REMOVED***

		// TODO(stevvooe): There are few more cases with subdirs that should be
		// handled in case the layout gets corrupted. This isn't strict enough
		// and may spew bad data.

		if path == root ***REMOVED***
			return nil
		***REMOVED***
		if filepath.Dir(path) == root ***REMOVED***
			alg = digest.Algorithm(filepath.Base(path))

			if !alg.Available() ***REMOVED***
				alg = ""
				return filepath.SkipDir
			***REMOVED***

			// descending into a hash directory
			return nil
		***REMOVED***

		dgst := digest.NewDigestFromHex(alg.String(), filepath.Base(path))
		if err := dgst.Validate(); err != nil ***REMOVED***
			// log error but don't report
			log.L.WithError(err).WithField("path", path).Error("invalid digest for blob path")
			// if we see this, it could mean some sort of corruption of the
			// store or extra paths not expected previously.
		***REMOVED***

		var labels map[string]string
		if s.ls != nil ***REMOVED***
			labels, err = s.ls.Get(dgst)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		return fn(s.info(dgst, fi, labels))
	***REMOVED***)
***REMOVED***

func (s *store) Status(ctx context.Context, ref string) (content.Status, error) ***REMOVED***
	return s.status(s.ingestRoot(ref))
***REMOVED***

func (s *store) ListStatuses(ctx context.Context, fs ...string) ([]content.Status, error) ***REMOVED***
	fp, err := os.Open(filepath.Join(s.root, "ingest"))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	defer fp.Close()

	fis, err := fp.Readdir(-1)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	filter, err := filters.ParseAll(fs...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var active []content.Status
	for _, fi := range fis ***REMOVED***
		p := filepath.Join(s.root, "ingest", fi.Name())
		stat, err := s.status(p)
		if err != nil ***REMOVED***
			if !os.IsNotExist(err) ***REMOVED***
				return nil, err
			***REMOVED***

			// TODO(stevvooe): This is a common error if uploads are being
			// completed while making this listing. Need to consider taking a
			// lock on the whole store to coordinate this aspect.
			//
			// Another option is to cleanup downloads asynchronously and
			// coordinate this method with the cleanup process.
			//
			// For now, we just skip them, as they really don't exist.
			continue
		***REMOVED***

		if filter.Match(adaptStatus(stat)) ***REMOVED***
			active = append(active, stat)
		***REMOVED***
	***REMOVED***

	return active, nil
***REMOVED***

// status works like stat above except uses the path to the ingest.
func (s *store) status(ingestPath string) (content.Status, error) ***REMOVED***
	dp := filepath.Join(ingestPath, "data")
	fi, err := os.Stat(dp)
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			err = errors.Wrap(errdefs.ErrNotFound, err.Error())
		***REMOVED***
		return content.Status***REMOVED******REMOVED***, err
	***REMOVED***

	ref, err := readFileString(filepath.Join(ingestPath, "ref"))
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			err = errors.Wrap(errdefs.ErrNotFound, err.Error())
		***REMOVED***
		return content.Status***REMOVED******REMOVED***, err
	***REMOVED***

	startedAt, err := readFileTimestamp(filepath.Join(ingestPath, "startedat"))
	if err != nil ***REMOVED***
		return content.Status***REMOVED******REMOVED***, errors.Wrapf(err, "could not read startedat")
	***REMOVED***

	updatedAt, err := readFileTimestamp(filepath.Join(ingestPath, "updatedat"))
	if err != nil ***REMOVED***
		return content.Status***REMOVED******REMOVED***, errors.Wrapf(err, "could not read updatedat")
	***REMOVED***

	// because we don't write updatedat on every write, the mod time may
	// actually be more up to date.
	if fi.ModTime().After(updatedAt) ***REMOVED***
		updatedAt = fi.ModTime()
	***REMOVED***

	return content.Status***REMOVED***
		Ref:       ref,
		Offset:    fi.Size(),
		Total:     s.total(ingestPath),
		UpdatedAt: updatedAt,
		StartedAt: startedAt,
	***REMOVED***, nil
***REMOVED***

func adaptStatus(status content.Status) filters.Adaptor ***REMOVED***
	return filters.AdapterFunc(func(fieldpath []string) (string, bool) ***REMOVED***
		if len(fieldpath) == 0 ***REMOVED***
			return "", false
		***REMOVED***
		switch fieldpath[0] ***REMOVED***
		case "ref":
			return status.Ref, true
		***REMOVED***

		return "", false
	***REMOVED***)
***REMOVED***

// total attempts to resolve the total expected size for the write.
func (s *store) total(ingestPath string) int64 ***REMOVED***
	totalS, err := readFileString(filepath.Join(ingestPath, "total"))
	if err != nil ***REMOVED***
		return 0
	***REMOVED***

	total, err := strconv.ParseInt(totalS, 10, 64)
	if err != nil ***REMOVED***
		// represents a corrupted file, should probably remove.
		return 0
	***REMOVED***

	return total
***REMOVED***

// Writer begins or resumes the active writer identified by ref. If the writer
// is already in use, an error is returned. Only one writer may be in use per
// ref at a time.
//
// The argument `ref` is used to uniquely identify a long-lived writer transaction.
func (s *store) Writer(ctx context.Context, ref string, total int64, expected digest.Digest) (content.Writer, error) ***REMOVED***
	var lockErr error
	for count := uint64(0); count < 10; count++ ***REMOVED***
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1<<count)))
		if err := tryLock(ref); err != nil ***REMOVED***
			if !errdefs.IsUnavailable(err) ***REMOVED***
				return nil, err
			***REMOVED***

			lockErr = err
		***REMOVED*** else ***REMOVED***
			lockErr = nil
			break
		***REMOVED***
	***REMOVED***

	if lockErr != nil ***REMOVED***
		return nil, lockErr
	***REMOVED***

	w, err := s.writer(ctx, ref, total, expected)
	if err != nil ***REMOVED***
		unlock(ref)
		return nil, err
	***REMOVED***

	return w, nil // lock is now held by w.
***REMOVED***

// writer provides the main implementation of the Writer method. The caller
// must hold the lock correctly and release on error if there is a problem.
func (s *store) writer(ctx context.Context, ref string, total int64, expected digest.Digest) (content.Writer, error) ***REMOVED***
	// TODO(stevvooe): Need to actually store expected here. We have
	// code in the service that shouldn't be dealing with this.
	if expected != "" ***REMOVED***
		p := s.blobPath(expected)
		if _, err := os.Stat(p); err == nil ***REMOVED***
			return nil, errors.Wrapf(errdefs.ErrAlreadyExists, "content %v", expected)
		***REMOVED***
	***REMOVED***

	path, refp, data := s.ingestPaths(ref)

	var (
		digester  = digest.Canonical.Digester()
		offset    int64
		startedAt time.Time
		updatedAt time.Time
	)

	// ensure that the ingest path has been created.
	if err := os.Mkdir(path, 0755); err != nil ***REMOVED***
		if !os.IsExist(err) ***REMOVED***
			return nil, err
		***REMOVED***

		status, err := s.status(path)
		if err != nil ***REMOVED***
			return nil, errors.Wrap(err, "failed reading status of resume write")
		***REMOVED***

		if ref != status.Ref ***REMOVED***
			// NOTE(stevvooe): This is fairly catastrophic. Either we have some
			// layout corruption or a hash collision for the ref key.
			return nil, errors.Wrapf(err, "ref key does not match: %v != %v", ref, status.Ref)
		***REMOVED***

		if total > 0 && status.Total > 0 && total != status.Total ***REMOVED***
			return nil, errors.Errorf("provided total differs from status: %v != %v", total, status.Total)
		***REMOVED***

		// TODO(stevvooe): slow slow slow!!, send to goroutine or use resumable hashes
		fp, err := os.Open(data)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		defer fp.Close()

		p := bufPool.Get().(*[]byte)
		defer bufPool.Put(p)

		offset, err = io.CopyBuffer(digester.Hash(), fp, *p)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		updatedAt = status.UpdatedAt
		startedAt = status.StartedAt
		total = status.Total
	***REMOVED*** else ***REMOVED***
		startedAt = time.Now()
		updatedAt = startedAt

		// the ingest is new, we need to setup the target location.
		// write the ref to a file for later use
		if err := ioutil.WriteFile(refp, []byte(ref), 0666); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		if writeTimestampFile(filepath.Join(path, "startedat"), startedAt); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		if writeTimestampFile(filepath.Join(path, "updatedat"), startedAt); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		if total > 0 ***REMOVED***
			if err := ioutil.WriteFile(filepath.Join(path, "total"), []byte(fmt.Sprint(total)), 0666); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	fp, err := os.OpenFile(data, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "failed to open data file")
	***REMOVED***

	if _, err := fp.Seek(offset, io.SeekStart); err != nil ***REMOVED***
		return nil, errors.Wrap(err, "could not seek to current write offset")
	***REMOVED***

	return &writer***REMOVED***
		s:         s,
		fp:        fp,
		ref:       ref,
		path:      path,
		offset:    offset,
		total:     total,
		digester:  digester,
		startedAt: startedAt,
		updatedAt: updatedAt,
	***REMOVED***, nil
***REMOVED***

// Abort an active transaction keyed by ref. If the ingest is active, it will
// be cancelled. Any resources associated with the ingest will be cleaned.
func (s *store) Abort(ctx context.Context, ref string) error ***REMOVED***
	root := s.ingestRoot(ref)
	if err := os.RemoveAll(root); err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return errors.Wrapf(errdefs.ErrNotFound, "ingest ref %q", ref)
		***REMOVED***

		return err
	***REMOVED***

	return nil
***REMOVED***

func (s *store) blobPath(dgst digest.Digest) string ***REMOVED***
	return filepath.Join(s.root, "blobs", dgst.Algorithm().String(), dgst.Hex())
***REMOVED***

func (s *store) ingestRoot(ref string) string ***REMOVED***
	dgst := digest.FromString(ref)
	return filepath.Join(s.root, "ingest", dgst.Hex())
***REMOVED***

// ingestPaths are returned. The paths are the following:
//
// - root: entire ingest directory
// - ref: name of the starting ref, must be unique
// - data: file where data is written
//
func (s *store) ingestPaths(ref string) (string, string, string) ***REMOVED***
	var (
		fp = s.ingestRoot(ref)
		rp = filepath.Join(fp, "ref")
		dp = filepath.Join(fp, "data")
	)

	return fp, rp, dp
***REMOVED***

func readFileString(path string) (string, error) ***REMOVED***
	p, err := ioutil.ReadFile(path)
	return string(p), err
***REMOVED***

// readFileTimestamp reads a file with just a timestamp present.
func readFileTimestamp(p string) (time.Time, error) ***REMOVED***
	b, err := ioutil.ReadFile(p)
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			err = errors.Wrap(errdefs.ErrNotFound, err.Error())
		***REMOVED***
		return time.Time***REMOVED******REMOVED***, err
	***REMOVED***

	var t time.Time
	if err := t.UnmarshalText(b); err != nil ***REMOVED***
		return time.Time***REMOVED******REMOVED***, errors.Wrapf(err, "could not parse timestamp file %v", p)
	***REMOVED***

	return t, nil
***REMOVED***

func writeTimestampFile(p string, t time.Time) error ***REMOVED***
	b, err := t.MarshalText()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return ioutil.WriteFile(p, b, 0666)
***REMOVED***
