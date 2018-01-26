package metadata

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/labels"
	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/metadata/boltutil"
	"github.com/containerd/containerd/mount"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/snapshots"
	"github.com/pkg/errors"
)

type snapshotter struct ***REMOVED***
	snapshots.Snapshotter
	name string
	db   *DB
	l    sync.RWMutex
***REMOVED***

// newSnapshotter returns a new Snapshotter which namespaces the given snapshot
// using the provided name and database.
func newSnapshotter(db *DB, name string, sn snapshots.Snapshotter) *snapshotter ***REMOVED***
	return &snapshotter***REMOVED***
		Snapshotter: sn,
		name:        name,
		db:          db,
	***REMOVED***
***REMOVED***

func createKey(id uint64, namespace, key string) string ***REMOVED***
	return fmt.Sprintf("%s/%d/%s", namespace, id, key)
***REMOVED***

func getKey(tx *bolt.Tx, ns, name, key string) string ***REMOVED***
	bkt := getSnapshotterBucket(tx, ns, name)
	if bkt == nil ***REMOVED***
		return ""
	***REMOVED***
	bkt = bkt.Bucket([]byte(key))
	if bkt == nil ***REMOVED***
		return ""
	***REMOVED***
	v := bkt.Get(bucketKeyName)
	if len(v) == 0 ***REMOVED***
		return ""
	***REMOVED***
	return string(v)
***REMOVED***

func (s *snapshotter) resolveKey(ctx context.Context, key string) (string, error) ***REMOVED***
	ns, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	var id string
	if err := view(ctx, s.db, func(tx *bolt.Tx) error ***REMOVED***
		id = getKey(tx, ns, s.name, key)
		if id == "" ***REMOVED***
			return errors.Wrapf(errdefs.ErrNotFound, "snapshot %v does not exist", key)
		***REMOVED***
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return id, nil
***REMOVED***

func (s *snapshotter) Stat(ctx context.Context, key string) (snapshots.Info, error) ***REMOVED***
	ns, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return snapshots.Info***REMOVED******REMOVED***, err
	***REMOVED***

	var (
		bkey  string
		local = snapshots.Info***REMOVED***
			Name: key,
		***REMOVED***
	)
	if err := view(ctx, s.db, func(tx *bolt.Tx) error ***REMOVED***
		bkt := getSnapshotterBucket(tx, ns, s.name)
		if bkt == nil ***REMOVED***
			return errors.Wrapf(errdefs.ErrNotFound, "snapshot %v does not exist", key)
		***REMOVED***
		sbkt := bkt.Bucket([]byte(key))
		if sbkt == nil ***REMOVED***
			return errors.Wrapf(errdefs.ErrNotFound, "snapshot %v does not exist", key)
		***REMOVED***
		local.Labels, err = boltutil.ReadLabels(sbkt)
		if err != nil ***REMOVED***
			return errors.Wrap(err, "failed to read labels")
		***REMOVED***
		if err := boltutil.ReadTimestamps(sbkt, &local.Created, &local.Updated); err != nil ***REMOVED***
			return errors.Wrap(err, "failed to read timestamps")
		***REMOVED***
		bkey = string(sbkt.Get(bucketKeyName))
		local.Parent = string(sbkt.Get(bucketKeyParent))

		return nil
	***REMOVED***); err != nil ***REMOVED***
		return snapshots.Info***REMOVED******REMOVED***, err
	***REMOVED***

	info, err := s.Snapshotter.Stat(ctx, bkey)
	if err != nil ***REMOVED***
		return snapshots.Info***REMOVED******REMOVED***, err
	***REMOVED***

	return overlayInfo(info, local), nil
***REMOVED***

func (s *snapshotter) Update(ctx context.Context, info snapshots.Info, fieldpaths ...string) (snapshots.Info, error) ***REMOVED***
	s.l.RLock()
	defer s.l.RUnlock()

	ns, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return snapshots.Info***REMOVED******REMOVED***, err
	***REMOVED***

	if info.Name == "" ***REMOVED***
		return snapshots.Info***REMOVED******REMOVED***, errors.Wrap(errdefs.ErrInvalidArgument, "")
	***REMOVED***

	var (
		bkey  string
		local = snapshots.Info***REMOVED***
			Name: info.Name,
		***REMOVED***
	)
	if err := update(ctx, s.db, func(tx *bolt.Tx) error ***REMOVED***
		bkt := getSnapshotterBucket(tx, ns, s.name)
		if bkt == nil ***REMOVED***
			return errors.Wrapf(errdefs.ErrNotFound, "snapshot %v does not exist", info.Name)
		***REMOVED***
		sbkt := bkt.Bucket([]byte(info.Name))
		if sbkt == nil ***REMOVED***
			return errors.Wrapf(errdefs.ErrNotFound, "snapshot %v does not exist", info.Name)
		***REMOVED***

		local.Labels, err = boltutil.ReadLabels(sbkt)
		if err != nil ***REMOVED***
			return errors.Wrap(err, "failed to read labels")
		***REMOVED***
		if err := boltutil.ReadTimestamps(sbkt, &local.Created, &local.Updated); err != nil ***REMOVED***
			return errors.Wrap(err, "failed to read timestamps")
		***REMOVED***

		// Handle field updates
		if len(fieldpaths) > 0 ***REMOVED***
			for _, path := range fieldpaths ***REMOVED***
				if strings.HasPrefix(path, "labels.") ***REMOVED***
					if local.Labels == nil ***REMOVED***
						local.Labels = map[string]string***REMOVED******REMOVED***
					***REMOVED***

					key := strings.TrimPrefix(path, "labels.")
					local.Labels[key] = info.Labels[key]
					continue
				***REMOVED***

				switch path ***REMOVED***
				case "labels":
					local.Labels = info.Labels
				default:
					return errors.Wrapf(errdefs.ErrInvalidArgument, "cannot update %q field on snapshot %q", path, info.Name)
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			local.Labels = info.Labels
		***REMOVED***
		if err := validateSnapshot(&local); err != nil ***REMOVED***
			return err
		***REMOVED***
		local.Updated = time.Now().UTC()

		if err := boltutil.WriteTimestamps(sbkt, local.Created, local.Updated); err != nil ***REMOVED***
			return errors.Wrap(err, "failed to read timestamps")
		***REMOVED***
		if err := boltutil.WriteLabels(sbkt, local.Labels); err != nil ***REMOVED***
			return errors.Wrap(err, "failed to read labels")
		***REMOVED***
		bkey = string(sbkt.Get(bucketKeyName))
		local.Parent = string(sbkt.Get(bucketKeyParent))

		return nil
	***REMOVED***); err != nil ***REMOVED***
		return snapshots.Info***REMOVED******REMOVED***, err
	***REMOVED***

	info, err = s.Snapshotter.Stat(ctx, bkey)
	if err != nil ***REMOVED***
		return snapshots.Info***REMOVED******REMOVED***, err
	***REMOVED***

	return overlayInfo(info, local), nil
***REMOVED***

func overlayInfo(info, overlay snapshots.Info) snapshots.Info ***REMOVED***
	// Merge info
	info.Name = overlay.Name
	info.Created = overlay.Created
	info.Updated = overlay.Updated
	info.Parent = overlay.Parent
	if info.Labels == nil ***REMOVED***
		info.Labels = overlay.Labels
	***REMOVED*** else ***REMOVED***
		for k, v := range overlay.Labels ***REMOVED***
			overlay.Labels[k] = v
		***REMOVED***
	***REMOVED***
	return info
***REMOVED***

func (s *snapshotter) Usage(ctx context.Context, key string) (snapshots.Usage, error) ***REMOVED***
	bkey, err := s.resolveKey(ctx, key)
	if err != nil ***REMOVED***
		return snapshots.Usage***REMOVED******REMOVED***, err
	***REMOVED***
	return s.Snapshotter.Usage(ctx, bkey)
***REMOVED***

func (s *snapshotter) Mounts(ctx context.Context, key string) ([]mount.Mount, error) ***REMOVED***
	bkey, err := s.resolveKey(ctx, key)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return s.Snapshotter.Mounts(ctx, bkey)
***REMOVED***

func (s *snapshotter) Prepare(ctx context.Context, key, parent string, opts ...snapshots.Opt) ([]mount.Mount, error) ***REMOVED***
	return s.createSnapshot(ctx, key, parent, false, opts)
***REMOVED***

func (s *snapshotter) View(ctx context.Context, key, parent string, opts ...snapshots.Opt) ([]mount.Mount, error) ***REMOVED***
	return s.createSnapshot(ctx, key, parent, true, opts)
***REMOVED***

func (s *snapshotter) createSnapshot(ctx context.Context, key, parent string, readonly bool, opts []snapshots.Opt) ([]mount.Mount, error) ***REMOVED***
	s.l.RLock()
	defer s.l.RUnlock()

	ns, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var base snapshots.Info
	for _, opt := range opts ***REMOVED***
		if err := opt(&base); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	if err := validateSnapshot(&base); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var m []mount.Mount
	if err := update(ctx, s.db, func(tx *bolt.Tx) error ***REMOVED***
		bkt, err := createSnapshotterBucket(tx, ns, s.name)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		bbkt, err := bkt.CreateBucket([]byte(key))
		if err != nil ***REMOVED***
			if err == bolt.ErrBucketExists ***REMOVED***
				err = errors.Wrapf(errdefs.ErrAlreadyExists, "snapshot %q", key)
			***REMOVED***
			return err
		***REMOVED***
		if err := addSnapshotLease(ctx, tx, s.name, key); err != nil ***REMOVED***
			return err
		***REMOVED***

		var bparent string
		if parent != "" ***REMOVED***
			pbkt := bkt.Bucket([]byte(parent))
			if pbkt == nil ***REMOVED***
				return errors.Wrapf(errdefs.ErrNotFound, "parent snapshot %v does not exist", parent)
			***REMOVED***
			bparent = string(pbkt.Get(bucketKeyName))

			cbkt, err := pbkt.CreateBucketIfNotExists(bucketKeyChildren)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if err := cbkt.Put([]byte(key), nil); err != nil ***REMOVED***
				return err
			***REMOVED***

			if err := bbkt.Put(bucketKeyParent, []byte(parent)); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		sid, err := bkt.NextSequence()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		bkey := createKey(sid, ns, key)
		if err := bbkt.Put(bucketKeyName, []byte(bkey)); err != nil ***REMOVED***
			return err
		***REMOVED***

		ts := time.Now().UTC()
		if err := boltutil.WriteTimestamps(bbkt, ts, ts); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := boltutil.WriteLabels(bbkt, base.Labels); err != nil ***REMOVED***
			return err
		***REMOVED***

		// TODO: Consider doing this outside of transaction to lessen
		// metadata lock time
		if readonly ***REMOVED***
			m, err = s.Snapshotter.View(ctx, bkey, bparent)
		***REMOVED*** else ***REMOVED***
			m, err = s.Snapshotter.Prepare(ctx, bkey, bparent)
		***REMOVED***
		return err
	***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return m, nil
***REMOVED***

func (s *snapshotter) Commit(ctx context.Context, name, key string, opts ...snapshots.Opt) error ***REMOVED***
	s.l.RLock()
	defer s.l.RUnlock()

	ns, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var base snapshots.Info
	for _, opt := range opts ***REMOVED***
		if err := opt(&base); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if err := validateSnapshot(&base); err != nil ***REMOVED***
		return err
	***REMOVED***

	return update(ctx, s.db, func(tx *bolt.Tx) error ***REMOVED***
		bkt := getSnapshotterBucket(tx, ns, s.name)
		if bkt == nil ***REMOVED***
			return errors.Wrapf(errdefs.ErrNotFound,
				"can not find snapshotter %q", s.name)
		***REMOVED***

		bbkt, err := bkt.CreateBucket([]byte(name))
		if err != nil ***REMOVED***
			if err == bolt.ErrBucketExists ***REMOVED***
				err = errors.Wrapf(errdefs.ErrAlreadyExists, "snapshot %q", name)
			***REMOVED***
			return err
		***REMOVED***
		if err := addSnapshotLease(ctx, tx, s.name, name); err != nil ***REMOVED***
			return err
		***REMOVED***

		obkt := bkt.Bucket([]byte(key))
		if obkt == nil ***REMOVED***
			return errors.Wrapf(errdefs.ErrNotFound, "snapshot %v does not exist", key)
		***REMOVED***

		bkey := string(obkt.Get(bucketKeyName))

		sid, err := bkt.NextSequence()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		nameKey := createKey(sid, ns, name)

		if err := bbkt.Put(bucketKeyName, []byte(nameKey)); err != nil ***REMOVED***
			return err
		***REMOVED***

		parent := obkt.Get(bucketKeyParent)
		if len(parent) > 0 ***REMOVED***
			pbkt := bkt.Bucket(parent)
			if pbkt == nil ***REMOVED***
				return errors.Wrapf(errdefs.ErrNotFound, "parent snapshot %v does not exist", string(parent))
			***REMOVED***

			cbkt, err := pbkt.CreateBucketIfNotExists(bucketKeyChildren)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if err := cbkt.Delete([]byte(key)); err != nil ***REMOVED***
				return err
			***REMOVED***
			if err := cbkt.Put([]byte(name), nil); err != nil ***REMOVED***
				return err
			***REMOVED***

			if err := bbkt.Put(bucketKeyParent, parent); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		ts := time.Now().UTC()
		if err := boltutil.WriteTimestamps(bbkt, ts, ts); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := boltutil.WriteLabels(bbkt, base.Labels); err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := bkt.DeleteBucket([]byte(key)); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := removeSnapshotLease(ctx, tx, s.name, key); err != nil ***REMOVED***
			return err
		***REMOVED***

		// TODO: Consider doing this outside of transaction to lessen
		// metadata lock time
		return s.Snapshotter.Commit(ctx, nameKey, bkey)
	***REMOVED***)

***REMOVED***

func (s *snapshotter) Remove(ctx context.Context, key string) error ***REMOVED***
	s.l.RLock()
	defer s.l.RUnlock()

	ns, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return update(ctx, s.db, func(tx *bolt.Tx) error ***REMOVED***
		var sbkt *bolt.Bucket
		bkt := getSnapshotterBucket(tx, ns, s.name)
		if bkt != nil ***REMOVED***
			sbkt = bkt.Bucket([]byte(key))
		***REMOVED***
		if sbkt == nil ***REMOVED***
			return errors.Wrapf(errdefs.ErrNotFound, "snapshot %v does not exist", key)
		***REMOVED***

		cbkt := sbkt.Bucket(bucketKeyChildren)
		if cbkt != nil ***REMOVED***
			if child, _ := cbkt.Cursor().First(); child != nil ***REMOVED***
				return errors.Wrap(errdefs.ErrFailedPrecondition, "cannot remove snapshot with child")
			***REMOVED***
		***REMOVED***

		parent := sbkt.Get(bucketKeyParent)
		if len(parent) > 0 ***REMOVED***
			pbkt := bkt.Bucket(parent)
			if pbkt == nil ***REMOVED***
				return errors.Wrapf(errdefs.ErrNotFound, "parent snapshot %v does not exist", string(parent))
			***REMOVED***
			cbkt := pbkt.Bucket(bucketKeyChildren)
			if cbkt != nil ***REMOVED***
				if err := cbkt.Delete([]byte(key)); err != nil ***REMOVED***
					return errors.Wrap(err, "failed to remove child link")
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if err := bkt.DeleteBucket([]byte(key)); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := removeSnapshotLease(ctx, tx, s.name, key); err != nil ***REMOVED***
			return err
		***REMOVED***

		// Mark snapshotter as dirty for triggering garbage collection
		s.db.dirtyL.Lock()
		s.db.dirtySS[s.name] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		s.db.dirtyL.Unlock()

		return nil
	***REMOVED***)
***REMOVED***

type infoPair struct ***REMOVED***
	bkey string
	info snapshots.Info
***REMOVED***

func (s *snapshotter) Walk(ctx context.Context, fn func(context.Context, snapshots.Info) error) error ***REMOVED***
	ns, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var (
		batchSize = 100
		pairs     = []infoPair***REMOVED******REMOVED***
		lastKey   string
	)

	for ***REMOVED***
		if err := view(ctx, s.db, func(tx *bolt.Tx) error ***REMOVED***
			bkt := getSnapshotterBucket(tx, ns, s.name)
			if bkt == nil ***REMOVED***
				return nil
			***REMOVED***

			c := bkt.Cursor()

			var k, v []byte
			if lastKey == "" ***REMOVED***
				k, v = c.First()
			***REMOVED*** else ***REMOVED***
				k, v = c.Seek([]byte(lastKey))
			***REMOVED***

			for k != nil ***REMOVED***
				if v == nil ***REMOVED***
					if len(pairs) >= batchSize ***REMOVED***
						break
					***REMOVED***
					sbkt := bkt.Bucket(k)

					pair := infoPair***REMOVED***
						bkey: string(sbkt.Get(bucketKeyName)),
						info: snapshots.Info***REMOVED***
							Name:   string(k),
							Parent: string(sbkt.Get(bucketKeyParent)),
						***REMOVED***,
					***REMOVED***

					err := boltutil.ReadTimestamps(sbkt, &pair.info.Created, &pair.info.Updated)
					if err != nil ***REMOVED***
						return err
					***REMOVED***
					pair.info.Labels, err = boltutil.ReadLabels(sbkt)
					if err != nil ***REMOVED***
						return err
					***REMOVED***

					pairs = append(pairs, pair)
				***REMOVED***

				k, v = c.Next()
			***REMOVED***

			lastKey = string(k)

			return nil
		***REMOVED***); err != nil ***REMOVED***
			return err
		***REMOVED***

		for _, pair := range pairs ***REMOVED***
			info, err := s.Snapshotter.Stat(ctx, pair.bkey)
			if err != nil ***REMOVED***
				if errdefs.IsNotFound(err) ***REMOVED***
					continue
				***REMOVED***
				return err
			***REMOVED***

			if err := fn(ctx, overlayInfo(info, pair.info)); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		if lastKey == "" ***REMOVED***
			break
		***REMOVED***

		pairs = pairs[:0]

	***REMOVED***

	return nil
***REMOVED***

func validateSnapshot(info *snapshots.Info) error ***REMOVED***
	for k, v := range info.Labels ***REMOVED***
		if err := labels.Validate(k, v); err != nil ***REMOVED***
			return errors.Wrapf(err, "info.Labels")
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (s *snapshotter) garbageCollect(ctx context.Context) (d time.Duration, err error) ***REMOVED***
	s.l.Lock()
	t1 := time.Now()
	defer func() ***REMOVED***
		if err == nil ***REMOVED***
			d = time.Now().Sub(t1)
		***REMOVED***
		s.l.Unlock()
	***REMOVED***()

	seen := map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	if err := s.db.View(func(tx *bolt.Tx) error ***REMOVED***
		v1bkt := tx.Bucket(bucketKeyVersion)
		if v1bkt == nil ***REMOVED***
			return nil
		***REMOVED***

		// iterate through each namespace
		v1c := v1bkt.Cursor()

		for k, v := v1c.First(); k != nil; k, v = v1c.Next() ***REMOVED***
			if v != nil ***REMOVED***
				continue
			***REMOVED***

			sbkt := v1bkt.Bucket(k).Bucket(bucketKeyObjectSnapshots)
			if sbkt == nil ***REMOVED***
				continue
			***REMOVED***

			// Load specific snapshotter
			ssbkt := sbkt.Bucket([]byte(s.name))
			if ssbkt == nil ***REMOVED***
				continue
			***REMOVED***

			if err := ssbkt.ForEach(func(sk, sv []byte) error ***REMOVED***
				if sv == nil ***REMOVED***
					bkey := ssbkt.Bucket(sk).Get(bucketKeyName)
					if len(bkey) > 0 ***REMOVED***
						seen[string(bkey)] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
					***REMOVED***
				***REMOVED***
				return nil
			***REMOVED***); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		return nil
	***REMOVED***); err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	roots, err := s.walkTree(ctx, seen)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	// TODO: Unlock before removal (once nodes are fully unavailable).
	// This could be achieved through doing prune inside the lock
	// and having a cleanup method which actually performs the
	// deletions on the snapshotters which support it.

	for _, node := range roots ***REMOVED***
		if err := s.pruneBranch(ctx, node); err != nil ***REMOVED***
			return 0, err
		***REMOVED***
	***REMOVED***

	return
***REMOVED***

type treeNode struct ***REMOVED***
	info     snapshots.Info
	remove   bool
	children []*treeNode
***REMOVED***

func (s *snapshotter) walkTree(ctx context.Context, seen map[string]struct***REMOVED******REMOVED***) ([]*treeNode, error) ***REMOVED***
	roots := []*treeNode***REMOVED******REMOVED***
	nodes := map[string]*treeNode***REMOVED******REMOVED***

	if err := s.Snapshotter.Walk(ctx, func(ctx context.Context, info snapshots.Info) error ***REMOVED***
		_, isSeen := seen[info.Name]
		node, ok := nodes[info.Name]
		if !ok ***REMOVED***
			node = &treeNode***REMOVED******REMOVED***
			nodes[info.Name] = node
		***REMOVED***

		node.remove = !isSeen
		node.info = info

		if info.Parent == "" ***REMOVED***
			roots = append(roots, node)
		***REMOVED*** else ***REMOVED***
			parent, ok := nodes[info.Parent]
			if !ok ***REMOVED***
				parent = &treeNode***REMOVED******REMOVED***
				nodes[info.Parent] = parent
			***REMOVED***
			parent.children = append(parent.children, node)
		***REMOVED***

		return nil
	***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return roots, nil
***REMOVED***

func (s *snapshotter) pruneBranch(ctx context.Context, node *treeNode) error ***REMOVED***
	for _, child := range node.children ***REMOVED***
		if err := s.pruneBranch(ctx, child); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if node.remove ***REMOVED***
		logger := log.G(ctx).WithField("snapshotter", s.name)
		if err := s.Snapshotter.Remove(ctx, node.info.Name); err != nil ***REMOVED***
			if !errdefs.IsFailedPrecondition(err) ***REMOVED***
				return err
			***REMOVED***
			logger.WithError(err).WithField("key", node.info.Name).Warnf("failed to remove snapshot")
		***REMOVED*** else ***REMOVED***
			logger.WithField("key", node.info.Name).Debug("removed snapshot")
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// Close closes s.Snapshotter but not db
func (s *snapshotter) Close() error ***REMOVED***
	return s.Snapshotter.Close()
***REMOVED***
