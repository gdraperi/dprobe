package metadata

import (
	"context"
	"encoding/binary"
	"strings"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/gc"
	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/snapshots"
	"github.com/pkg/errors"
)

const (
	// schemaVersion represents the schema version of
	// the database. This schema version represents the
	// structure of the data in the database. The schema
	// can envolve at any time but any backwards
	// incompatible changes or structural changes require
	// bumping the schema version.
	schemaVersion = "v1"

	// dbVersion represents updates to the schema
	// version which are additions and compatible with
	// prior version of the same schema.
	dbVersion = 1
)

// DB represents a metadata database backed by a bolt
// database. The database is fully namespaced and stores
// image, container, namespace, snapshot, and content data
// while proxying data shared across namespaces to backend
// datastores for content and snapshots.
type DB struct ***REMOVED***
	db *bolt.DB
	ss map[string]*snapshotter
	cs *contentStore

	// wlock is used to protect access to the data structures during garbage
	// collection. While the wlock is held no writable transactions can be
	// opened, preventing changes from occurring between the mark and
	// sweep phases without preventing read transactions.
	wlock sync.RWMutex

	// dirty flags and lock keeps track of datastores which have had deletions
	// since the last garbage collection. These datastores will will be garbage
	// collected during the next garbage collection.
	dirtyL  sync.Mutex
	dirtySS map[string]struct***REMOVED******REMOVED***
	dirtyCS bool

	// mutationCallbacks are called after each mutation with the flag
	// set indicating whether any dirty flags are set
	mutationCallbacks []func(bool)
***REMOVED***

// NewDB creates a new metadata database using the provided
// bolt database, content store, and snapshotters.
func NewDB(db *bolt.DB, cs content.Store, ss map[string]snapshots.Snapshotter) *DB ***REMOVED***
	m := &DB***REMOVED***
		db:      db,
		ss:      make(map[string]*snapshotter, len(ss)),
		dirtySS: map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***,
	***REMOVED***

	// Initialize data stores
	m.cs = newContentStore(m, cs)
	for name, sn := range ss ***REMOVED***
		m.ss[name] = newSnapshotter(m, name, sn)
	***REMOVED***

	return m
***REMOVED***

// Init ensures the database is at the correct version
// and performs any needed migrations.
func (m *DB) Init(ctx context.Context) error ***REMOVED***
	// errSkip is used when no migration or version needs to be written
	// to the database and the transaction can be immediately rolled
	// back rather than performing a much slower and unnecessary commit.
	var errSkip = errors.New("skip update")

	err := m.db.Update(func(tx *bolt.Tx) error ***REMOVED***
		var (
			// current schema and version
			schema  = "v0"
			version = 0
		)

		// i represents the index of the first migration
		// which must be run to get the database up to date.
		// The migration's version will be checked in reverse
		// order, decrementing i for each migration which
		// represents a version newer than the current
		// database version
		i := len(migrations)

		for ; i > 0; i-- ***REMOVED***
			migration := migrations[i-1]

			bkt := tx.Bucket([]byte(migration.schema))
			if bkt == nil ***REMOVED***
				// Hasn't encountered another schema, go to next migration
				if schema == "v0" ***REMOVED***
					continue
				***REMOVED***
				break
			***REMOVED***
			if schema == "v0" ***REMOVED***
				schema = migration.schema
				vb := bkt.Get(bucketKeyDBVersion)
				if vb != nil ***REMOVED***
					v, _ := binary.Varint(vb)
					version = int(v)
				***REMOVED***
			***REMOVED***

			if version >= migration.version ***REMOVED***
				break
			***REMOVED***
		***REMOVED***

		// Previous version fo database found
		if schema != "v0" ***REMOVED***
			updates := migrations[i:]

			// No migration updates, return immediately
			if len(updates) == 0 ***REMOVED***
				return errSkip
			***REMOVED***

			for _, m := range updates ***REMOVED***
				t0 := time.Now()
				if err := m.migrate(tx); err != nil ***REMOVED***
					return errors.Wrapf(err, "failed to migrate to %s.%d", m.schema, m.version)
				***REMOVED***
				log.G(ctx).WithField("d", time.Now().Sub(t0)).Debugf("finished database migration to %s.%d", m.schema, m.version)
			***REMOVED***
		***REMOVED***

		bkt, err := tx.CreateBucketIfNotExists(bucketKeyVersion)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		versionEncoded, err := encodeInt(dbVersion)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		return bkt.Put(bucketKeyDBVersion, versionEncoded)
	***REMOVED***)
	if err == errSkip ***REMOVED***
		err = nil
	***REMOVED***
	return err
***REMOVED***

// ContentStore returns a namespaced content store
// proxied to a content store.
func (m *DB) ContentStore() content.Store ***REMOVED***
	if m.cs == nil ***REMOVED***
		return nil
	***REMOVED***
	return m.cs
***REMOVED***

// Snapshotter returns a namespaced content store for
// the requested snapshotter name proxied to a snapshotter.
func (m *DB) Snapshotter(name string) snapshots.Snapshotter ***REMOVED***
	sn, ok := m.ss[name]
	if !ok ***REMOVED***
		return nil
	***REMOVED***
	return sn
***REMOVED***

// View runs a readonly transaction on the metadata store.
func (m *DB) View(fn func(*bolt.Tx) error) error ***REMOVED***
	return m.db.View(fn)
***REMOVED***

// Update runs a writable transaction on the metadata store.
func (m *DB) Update(fn func(*bolt.Tx) error) error ***REMOVED***
	m.wlock.RLock()
	defer m.wlock.RUnlock()
	err := m.db.Update(fn)
	if err == nil ***REMOVED***
		m.dirtyL.Lock()
		dirty := m.dirtyCS || len(m.dirtySS) > 0
		for _, fn := range m.mutationCallbacks ***REMOVED***
			fn(dirty)
		***REMOVED***
		m.dirtyL.Unlock()
	***REMOVED***

	return err
***REMOVED***

// RegisterMutationCallback registers a function to be called after a metadata
// mutations has been performed.
//
// The callback function in an argument for whether a deletion has occurred
// since the last garbage collection.
func (m *DB) RegisterMutationCallback(fn func(bool)) ***REMOVED***
	m.dirtyL.Lock()
	m.mutationCallbacks = append(m.mutationCallbacks, fn)
	m.dirtyL.Unlock()
***REMOVED***

// GCStats holds the duration for the different phases of the garbage collector
type GCStats struct ***REMOVED***
	MetaD     time.Duration
	ContentD  time.Duration
	SnapshotD map[string]time.Duration
***REMOVED***

// GarbageCollect starts garbage collection
func (m *DB) GarbageCollect(ctx context.Context) (stats GCStats, err error) ***REMOVED***
	m.wlock.Lock()
	t1 := time.Now()

	marked, err := m.getMarked(ctx)
	if err != nil ***REMOVED***
		m.wlock.Unlock()
		return GCStats***REMOVED******REMOVED***, err
	***REMOVED***

	m.dirtyL.Lock()

	if err := m.db.Update(func(tx *bolt.Tx) error ***REMOVED***
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		rm := func(ctx context.Context, n gc.Node) error ***REMOVED***
			if _, ok := marked[n]; ok ***REMOVED***
				return nil
			***REMOVED***

			if n.Type == ResourceSnapshot ***REMOVED***
				if idx := strings.IndexRune(n.Key, '/'); idx > 0 ***REMOVED***
					m.dirtySS[n.Key[:idx]] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
				***REMOVED***
			***REMOVED*** else if n.Type == ResourceContent ***REMOVED***
				m.dirtyCS = true
			***REMOVED***
			return remove(ctx, tx, n)
		***REMOVED***

		if err := scanAll(ctx, tx, rm); err != nil ***REMOVED***
			return errors.Wrap(err, "failed to scan and remove")
		***REMOVED***

		return nil
	***REMOVED***); err != nil ***REMOVED***
		m.dirtyL.Unlock()
		m.wlock.Unlock()
		return GCStats***REMOVED******REMOVED***, err
	***REMOVED***

	var wg sync.WaitGroup

	if len(m.dirtySS) > 0 ***REMOVED***
		var sl sync.Mutex
		stats.SnapshotD = map[string]time.Duration***REMOVED******REMOVED***
		wg.Add(len(m.dirtySS))
		for snapshotterName := range m.dirtySS ***REMOVED***
			log.G(ctx).WithField("snapshotter", snapshotterName).Debug("schedule snapshotter cleanup")
			go func(snapshotterName string) ***REMOVED***
				st1 := time.Now()
				m.cleanupSnapshotter(snapshotterName)

				sl.Lock()
				stats.SnapshotD[snapshotterName] = time.Now().Sub(st1)
				sl.Unlock()

				wg.Done()
			***REMOVED***(snapshotterName)
		***REMOVED***
		m.dirtySS = map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	if m.dirtyCS ***REMOVED***
		wg.Add(1)
		log.G(ctx).Debug("schedule content cleanup")
		go func() ***REMOVED***
			ct1 := time.Now()
			m.cleanupContent()
			stats.ContentD = time.Now().Sub(ct1)
			wg.Done()
		***REMOVED***()
		m.dirtyCS = false
	***REMOVED***

	m.dirtyL.Unlock()

	stats.MetaD = time.Now().Sub(t1)
	m.wlock.Unlock()

	wg.Wait()

	return
***REMOVED***

func (m *DB) getMarked(ctx context.Context) (map[gc.Node]struct***REMOVED******REMOVED***, error) ***REMOVED***
	var marked map[gc.Node]struct***REMOVED******REMOVED***
	if err := m.db.View(func(tx *bolt.Tx) error ***REMOVED***
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		var (
			nodes []gc.Node
			wg    sync.WaitGroup
			roots = make(chan gc.Node)
		)
		wg.Add(1)
		go func() ***REMOVED***
			defer wg.Done()
			for n := range roots ***REMOVED***
				nodes = append(nodes, n)
			***REMOVED***
		***REMOVED***()
		// Call roots
		if err := scanRoots(ctx, tx, roots); err != nil ***REMOVED***
			cancel()
			return err
		***REMOVED***
		close(roots)
		wg.Wait()

		refs := func(n gc.Node) ([]gc.Node, error) ***REMOVED***
			var sn []gc.Node
			if err := references(ctx, tx, n, func(nn gc.Node) ***REMOVED***
				sn = append(sn, nn)
			***REMOVED***); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return sn, nil
		***REMOVED***

		reachable, err := gc.Tricolor(nodes, refs)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		marked = reachable
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return marked, nil
***REMOVED***

func (m *DB) cleanupSnapshotter(name string) (time.Duration, error) ***REMOVED***
	ctx := context.Background()
	sn, ok := m.ss[name]
	if !ok ***REMOVED***
		return 0, nil
	***REMOVED***

	d, err := sn.garbageCollect(ctx)
	logger := log.G(ctx).WithField("snapshotter", name)
	if err != nil ***REMOVED***
		logger.WithError(err).Warn("snapshot garbage collection failed")
	***REMOVED*** else ***REMOVED***
		logger.WithField("d", d).Debugf("snapshot garbage collected")
	***REMOVED***
	return d, err
***REMOVED***

func (m *DB) cleanupContent() (time.Duration, error) ***REMOVED***
	ctx := context.Background()
	if m.cs == nil ***REMOVED***
		return 0, nil
	***REMOVED***

	d, err := m.cs.garbageCollect(ctx)
	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Warn("content garbage collection failed")
	***REMOVED*** else ***REMOVED***
		log.G(ctx).WithField("d", d).Debugf("content garbage collected")
	***REMOVED***

	return d, err
***REMOVED***
