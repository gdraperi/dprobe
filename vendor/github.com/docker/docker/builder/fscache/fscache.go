package fscache

import (
	"archive/tar"
	"crypto/sha256"
	"encoding/json"
	"hash"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/docker/docker/builder"
	"github.com/docker/docker/builder/remotecontext"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/directory"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/pkg/tarsum"
	"github.com/moby/buildkit/session/filesync"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tonistiigi/fsutil"
	"golang.org/x/net/context"
	"golang.org/x/sync/singleflight"
)

const dbFile = "fscache.db"
const cacheKey = "cache"
const metaKey = "meta"

// Backend is a backing implementation for FSCache
type Backend interface ***REMOVED***
	Get(id string) (string, error)
	Remove(id string) error
***REMOVED***

// FSCache allows syncing remote resources to cached snapshots
type FSCache struct ***REMOVED***
	opt        Opt
	transports map[string]Transport
	mu         sync.Mutex
	g          singleflight.Group
	store      *fsCacheStore
***REMOVED***

// Opt defines options for initializing FSCache
type Opt struct ***REMOVED***
	Backend  Backend
	Root     string // for storing local metadata
	GCPolicy GCPolicy
***REMOVED***

// GCPolicy defines policy for garbage collection
type GCPolicy struct ***REMOVED***
	MaxSize         uint64
	MaxKeepDuration time.Duration
***REMOVED***

// NewFSCache returns new FSCache object
func NewFSCache(opt Opt) (*FSCache, error) ***REMOVED***
	store, err := newFSCacheStore(opt)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &FSCache***REMOVED***
		store:      store,
		opt:        opt,
		transports: make(map[string]Transport),
	***REMOVED***, nil
***REMOVED***

// Transport defines a method for syncing remote data to FSCache
type Transport interface ***REMOVED***
	Copy(ctx context.Context, id RemoteIdentifier, dest string, cs filesync.CacheUpdater) error
***REMOVED***

// RemoteIdentifier identifies a transfer request
type RemoteIdentifier interface ***REMOVED***
	Key() string
	SharedKey() string
	Transport() string
***REMOVED***

// RegisterTransport registers a new transport method
func (fsc *FSCache) RegisterTransport(id string, transport Transport) error ***REMOVED***
	fsc.mu.Lock()
	defer fsc.mu.Unlock()
	if _, ok := fsc.transports[id]; ok ***REMOVED***
		return errors.Errorf("transport %v already exists", id)
	***REMOVED***
	fsc.transports[id] = transport
	return nil
***REMOVED***

// SyncFrom returns a source based on a remote identifier
func (fsc *FSCache) SyncFrom(ctx context.Context, id RemoteIdentifier) (builder.Source, error) ***REMOVED*** // cacheOpt
	trasportID := id.Transport()
	fsc.mu.Lock()
	transport, ok := fsc.transports[id.Transport()]
	if !ok ***REMOVED***
		fsc.mu.Unlock()
		return nil, errors.Errorf("invalid transport %s", trasportID)
	***REMOVED***

	logrus.Debugf("SyncFrom %s %s", id.Key(), id.SharedKey())
	fsc.mu.Unlock()
	sourceRef, err, _ := fsc.g.Do(id.Key(), func() (interface***REMOVED******REMOVED***, error) ***REMOVED***
		var sourceRef *cachedSourceRef
		sourceRef, err := fsc.store.Get(id.Key())
		if err == nil ***REMOVED***
			return sourceRef, nil
		***REMOVED***

		// check for unused shared cache
		sharedKey := id.SharedKey()
		if sharedKey != "" ***REMOVED***
			r, err := fsc.store.Rebase(sharedKey, id.Key())
			if err == nil ***REMOVED***
				sourceRef = r
			***REMOVED***
		***REMOVED***

		if sourceRef == nil ***REMOVED***
			var err error
			sourceRef, err = fsc.store.New(id.Key(), sharedKey)
			if err != nil ***REMOVED***
				return nil, errors.Wrap(err, "failed to create remote context")
			***REMOVED***
		***REMOVED***

		if err := syncFrom(ctx, sourceRef, transport, id); err != nil ***REMOVED***
			sourceRef.Release()
			return nil, err
		***REMOVED***
		if err := sourceRef.resetSize(-1); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return sourceRef, nil
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ref := sourceRef.(*cachedSourceRef)
	if ref.src == nil ***REMOVED*** // failsafe
		return nil, errors.Errorf("invalid empty pull")
	***REMOVED***
	wc := &wrappedContext***REMOVED***Source: ref.src, closer: func() error ***REMOVED***
		ref.Release()
		return nil
	***REMOVED******REMOVED***
	return wc, nil
***REMOVED***

// DiskUsage reports how much data is allocated by the cache
func (fsc *FSCache) DiskUsage() (int64, error) ***REMOVED***
	return fsc.store.DiskUsage()
***REMOVED***

// Prune allows manually cleaning up the cache
func (fsc *FSCache) Prune(ctx context.Context) (uint64, error) ***REMOVED***
	return fsc.store.Prune(ctx)
***REMOVED***

// Close stops the gc and closes the persistent db
func (fsc *FSCache) Close() error ***REMOVED***
	return fsc.store.Close()
***REMOVED***

func syncFrom(ctx context.Context, cs *cachedSourceRef, transport Transport, id RemoteIdentifier) (retErr error) ***REMOVED***
	src := cs.src
	if src == nil ***REMOVED***
		src = remotecontext.NewCachableSource(cs.Dir())
	***REMOVED***

	if !cs.cached ***REMOVED***
		if err := cs.storage.db.View(func(tx *bolt.Tx) error ***REMOVED***
			b := tx.Bucket([]byte(id.Key()))
			dt := b.Get([]byte(cacheKey))
			if dt != nil ***REMOVED***
				if err := src.UnmarshalBinary(dt); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				return errors.Wrap(src.Scan(), "failed to scan cache records")
			***REMOVED***
			return nil
		***REMOVED***); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	dc := &detectChanges***REMOVED***f: src.HandleChange***REMOVED***

	// todo: probably send a bucket to `Copy` and let it return source
	// but need to make sure that tx is safe
	if err := transport.Copy(ctx, id, cs.Dir(), dc); err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to copy to %s", cs.Dir())
	***REMOVED***

	if !dc.supported ***REMOVED***
		if err := src.Scan(); err != nil ***REMOVED***
			return errors.Wrap(err, "failed to scan cache records after transfer")
		***REMOVED***
	***REMOVED***
	cs.cached = true
	cs.src = src
	return cs.storage.db.Update(func(tx *bolt.Tx) error ***REMOVED***
		dt, err := src.MarshalBinary()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		b := tx.Bucket([]byte(id.Key()))
		return b.Put([]byte(cacheKey), dt)
	***REMOVED***)
***REMOVED***

type fsCacheStore struct ***REMOVED***
	mu       sync.Mutex
	sources  map[string]*cachedSource
	db       *bolt.DB
	fs       Backend
	gcTimer  *time.Timer
	gcPolicy GCPolicy
***REMOVED***

// CachePolicy defines policy for keeping a resource in cache
type CachePolicy struct ***REMOVED***
	Priority int
	LastUsed time.Time
***REMOVED***

func defaultCachePolicy() CachePolicy ***REMOVED***
	return CachePolicy***REMOVED***Priority: 10, LastUsed: time.Now()***REMOVED***
***REMOVED***

func newFSCacheStore(opt Opt) (*fsCacheStore, error) ***REMOVED***
	if err := os.MkdirAll(opt.Root, 0700); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	p := filepath.Join(opt.Root, dbFile)
	db, err := bolt.Open(p, 0600, nil)
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "failed to open database file %s")
	***REMOVED***
	s := &fsCacheStore***REMOVED***db: db, sources: make(map[string]*cachedSource), fs: opt.Backend, gcPolicy: opt.GCPolicy***REMOVED***
	db.View(func(tx *bolt.Tx) error ***REMOVED***
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error ***REMOVED***
			dt := b.Get([]byte(metaKey))
			if dt == nil ***REMOVED***
				return nil
			***REMOVED***
			var sm sourceMeta
			if err := json.Unmarshal(dt, &sm); err != nil ***REMOVED***
				return err
			***REMOVED***
			dir, err := s.fs.Get(sm.BackendID)
			if err != nil ***REMOVED***
				return err // TODO: handle gracefully
			***REMOVED***
			source := &cachedSource***REMOVED***
				refs:       make(map[*cachedSourceRef]struct***REMOVED******REMOVED***),
				id:         string(name),
				dir:        dir,
				sourceMeta: sm,
				storage:    s,
			***REMOVED***
			s.sources[string(name)] = source
			return nil
		***REMOVED***)
	***REMOVED***)

	s.gcTimer = s.startPeriodicGC(5 * time.Minute)
	return s, nil
***REMOVED***

func (s *fsCacheStore) startPeriodicGC(interval time.Duration) *time.Timer ***REMOVED***
	var t *time.Timer
	t = time.AfterFunc(interval, func() ***REMOVED***
		if err := s.GC(); err != nil ***REMOVED***
			logrus.Errorf("build gc error: %v", err)
		***REMOVED***
		t.Reset(interval)
	***REMOVED***)
	return t
***REMOVED***

func (s *fsCacheStore) Close() error ***REMOVED***
	s.gcTimer.Stop()
	return s.db.Close()
***REMOVED***

func (s *fsCacheStore) New(id, sharedKey string) (*cachedSourceRef, error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	var ret *cachedSource
	if err := s.db.Update(func(tx *bolt.Tx) error ***REMOVED***
		b, err := tx.CreateBucket([]byte(id))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		backendID := stringid.GenerateRandomID()
		dir, err := s.fs.Get(backendID)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		source := &cachedSource***REMOVED***
			refs: make(map[*cachedSourceRef]struct***REMOVED******REMOVED***),
			id:   id,
			dir:  dir,
			sourceMeta: sourceMeta***REMOVED***
				BackendID:   backendID,
				SharedKey:   sharedKey,
				CachePolicy: defaultCachePolicy(),
			***REMOVED***,
			storage: s,
		***REMOVED***
		dt, err := json.Marshal(source.sourceMeta)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := b.Put([]byte(metaKey), dt); err != nil ***REMOVED***
			return err
		***REMOVED***
		s.sources[id] = source
		ret = source
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return ret.getRef(), nil
***REMOVED***

func (s *fsCacheStore) Rebase(sharedKey, newid string) (*cachedSourceRef, error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	var ret *cachedSource
	for id, snap := range s.sources ***REMOVED***
		if snap.SharedKey == sharedKey && len(snap.refs) == 0 ***REMOVED***
			if err := s.db.Update(func(tx *bolt.Tx) error ***REMOVED***
				if err := tx.DeleteBucket([]byte(id)); err != nil ***REMOVED***
					return err
				***REMOVED***
				b, err := tx.CreateBucket([]byte(newid))
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				snap.id = newid
				snap.CachePolicy = defaultCachePolicy()
				dt, err := json.Marshal(snap.sourceMeta)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				if err := b.Put([]byte(metaKey), dt); err != nil ***REMOVED***
					return err
				***REMOVED***
				delete(s.sources, id)
				s.sources[newid] = snap
				return nil
			***REMOVED***); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			ret = snap
			break
		***REMOVED***
	***REMOVED***
	if ret == nil ***REMOVED***
		return nil, errors.Errorf("no candidate for rebase")
	***REMOVED***
	return ret.getRef(), nil
***REMOVED***

func (s *fsCacheStore) Get(id string) (*cachedSourceRef, error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	src, ok := s.sources[id]
	if !ok ***REMOVED***
		return nil, errors.Errorf("not found")
	***REMOVED***
	return src.getRef(), nil
***REMOVED***

// DiskUsage reports how much data is allocated by the cache
func (s *fsCacheStore) DiskUsage() (int64, error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	var size int64

	for _, snap := range s.sources ***REMOVED***
		if len(snap.refs) == 0 ***REMOVED***
			ss, err := snap.getSize()
			if err != nil ***REMOVED***
				return 0, err
			***REMOVED***
			size += ss
		***REMOVED***
	***REMOVED***
	return size, nil
***REMOVED***

// Prune allows manually cleaning up the cache
func (s *fsCacheStore) Prune(ctx context.Context) (uint64, error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	var size uint64

	for id, snap := range s.sources ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			logrus.Debugf("Cache prune operation cancelled, pruned size: %d", size)
			// when the context is cancelled, only return current size and nil
			return size, nil
		default:
		***REMOVED***
		if len(snap.refs) == 0 ***REMOVED***
			ss, err := snap.getSize()
			if err != nil ***REMOVED***
				return size, err
			***REMOVED***
			if err := s.delete(id); err != nil ***REMOVED***
				return size, errors.Wrapf(err, "failed to delete %s", id)
			***REMOVED***
			size += uint64(ss)
		***REMOVED***
	***REMOVED***
	return size, nil
***REMOVED***

// GC runs a garbage collector on FSCache
func (s *fsCacheStore) GC() error ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	var size uint64

	cutoff := time.Now().Add(-s.gcPolicy.MaxKeepDuration)
	var blacklist []*cachedSource

	for id, snap := range s.sources ***REMOVED***
		if len(snap.refs) == 0 ***REMOVED***
			if cutoff.After(snap.CachePolicy.LastUsed) ***REMOVED***
				if err := s.delete(id); err != nil ***REMOVED***
					return errors.Wrapf(err, "failed to delete %s", id)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				ss, err := snap.getSize()
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				size += uint64(ss)
				blacklist = append(blacklist, snap)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	sort.Sort(sortableCacheSources(blacklist))
	for _, snap := range blacklist ***REMOVED***
		if size <= s.gcPolicy.MaxSize ***REMOVED***
			break
		***REMOVED***
		ss, err := snap.getSize()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := s.delete(snap.id); err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to delete %s", snap.id)
		***REMOVED***
		size -= uint64(ss)
	***REMOVED***
	return nil
***REMOVED***

// keep mu while calling this
func (s *fsCacheStore) delete(id string) error ***REMOVED***
	src, ok := s.sources[id]
	if !ok ***REMOVED***
		return nil
	***REMOVED***
	if len(src.refs) > 0 ***REMOVED***
		return errors.Errorf("can't delete %s because it has active references", id)
	***REMOVED***
	delete(s.sources, id)
	if err := s.db.Update(func(tx *bolt.Tx) error ***REMOVED***
		return tx.DeleteBucket([]byte(id))
	***REMOVED***); err != nil ***REMOVED***
		return err
	***REMOVED***
	return s.fs.Remove(src.BackendID)
***REMOVED***

type sourceMeta struct ***REMOVED***
	SharedKey   string
	BackendID   string
	CachePolicy CachePolicy
	Size        int64
***REMOVED***

type cachedSource struct ***REMOVED***
	sourceMeta
	refs    map[*cachedSourceRef]struct***REMOVED******REMOVED***
	id      string
	dir     string
	src     *remotecontext.CachableSource
	storage *fsCacheStore
	cached  bool // keep track if cache is up to date
***REMOVED***

type cachedSourceRef struct ***REMOVED***
	*cachedSource
***REMOVED***

func (cs *cachedSource) Dir() string ***REMOVED***
	return cs.dir
***REMOVED***

// hold storage lock before calling
func (cs *cachedSource) getRef() *cachedSourceRef ***REMOVED***
	ref := &cachedSourceRef***REMOVED***cachedSource: cs***REMOVED***
	cs.refs[ref] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	return ref
***REMOVED***

// hold storage lock before calling
func (cs *cachedSource) getSize() (int64, error) ***REMOVED***
	if cs.sourceMeta.Size < 0 ***REMOVED***
		ss, err := directory.Size(cs.dir)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		if err := cs.resetSize(ss); err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		return ss, nil
	***REMOVED***
	return cs.sourceMeta.Size, nil
***REMOVED***

func (cs *cachedSource) resetSize(val int64) error ***REMOVED***
	cs.sourceMeta.Size = val
	return cs.saveMeta()
***REMOVED***
func (cs *cachedSource) saveMeta() error ***REMOVED***
	return cs.storage.db.Update(func(tx *bolt.Tx) error ***REMOVED***
		b := tx.Bucket([]byte(cs.id))
		dt, err := json.Marshal(cs.sourceMeta)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		return b.Put([]byte(metaKey), dt)
	***REMOVED***)
***REMOVED***

func (csr *cachedSourceRef) Release() error ***REMOVED***
	csr.cachedSource.storage.mu.Lock()
	defer csr.cachedSource.storage.mu.Unlock()
	delete(csr.cachedSource.refs, csr)
	if len(csr.cachedSource.refs) == 0 ***REMOVED***
		go csr.cachedSource.storage.GC()
	***REMOVED***
	return nil
***REMOVED***

type detectChanges struct ***REMOVED***
	f         fsutil.ChangeFunc
	supported bool
***REMOVED***

func (dc *detectChanges) HandleChange(kind fsutil.ChangeKind, path string, fi os.FileInfo, err error) error ***REMOVED***
	if dc == nil ***REMOVED***
		return nil
	***REMOVED***
	return dc.f(kind, path, fi, err)
***REMOVED***

func (dc *detectChanges) MarkSupported(v bool) ***REMOVED***
	if dc == nil ***REMOVED***
		return
	***REMOVED***
	dc.supported = v
***REMOVED***

func (dc *detectChanges) ContentHasher() fsutil.ContentHasher ***REMOVED***
	return newTarsumHash
***REMOVED***

type wrappedContext struct ***REMOVED***
	builder.Source
	closer func() error
***REMOVED***

func (wc *wrappedContext) Close() error ***REMOVED***
	if err := wc.Source.Close(); err != nil ***REMOVED***
		return err
	***REMOVED***
	return wc.closer()
***REMOVED***

type sortableCacheSources []*cachedSource

// Len is the number of elements in the collection.
func (s sortableCacheSources) Len() int ***REMOVED***
	return len(s)
***REMOVED***

// Less reports whether the element with
// index i should sort before the element with index j.
func (s sortableCacheSources) Less(i, j int) bool ***REMOVED***
	return s[i].CachePolicy.LastUsed.Before(s[j].CachePolicy.LastUsed)
***REMOVED***

// Swap swaps the elements with indexes i and j.
func (s sortableCacheSources) Swap(i, j int) ***REMOVED***
	s[i], s[j] = s[j], s[i]
***REMOVED***

func newTarsumHash(stat *fsutil.Stat) (hash.Hash, error) ***REMOVED***
	fi := &fsutil.StatInfo***REMOVED***stat***REMOVED***
	p := stat.Path
	if fi.IsDir() ***REMOVED***
		p += string(os.PathSeparator)
	***REMOVED***
	h, err := archive.FileInfoHeader(p, fi, stat.Linkname)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	h.Name = p
	h.Uid = int(stat.Uid)
	h.Gid = int(stat.Gid)
	h.Linkname = stat.Linkname
	if stat.Xattrs != nil ***REMOVED***
		h.Xattrs = make(map[string]string)
		for k, v := range stat.Xattrs ***REMOVED***
			h.Xattrs[k] = string(v)
		***REMOVED***
	***REMOVED***

	tsh := &tarsumHash***REMOVED***h: h, Hash: sha256.New()***REMOVED***
	tsh.Reset()
	return tsh, nil
***REMOVED***

// Reset resets the Hash to its initial state.
func (tsh *tarsumHash) Reset() ***REMOVED***
	tsh.Hash.Reset()
	tarsum.WriteV1Header(tsh.h, tsh.Hash)
***REMOVED***

type tarsumHash struct ***REMOVED***
	hash.Hash
	h *tar.Header
***REMOVED***
