package metadata

import (
	"context"
	"encoding/binary"
	"strings"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/filters"
	"github.com/containerd/containerd/labels"
	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/metadata/boltutil"
	"github.com/containerd/containerd/namespaces"
	digest "github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

type contentStore struct ***REMOVED***
	content.Store
	db *DB
	l  sync.RWMutex
***REMOVED***

// newContentStore returns a namespaced content store using an existing
// content store interface.
func newContentStore(db *DB, cs content.Store) *contentStore ***REMOVED***
	return &contentStore***REMOVED***
		Store: cs,
		db:    db,
	***REMOVED***
***REMOVED***

func (cs *contentStore) Info(ctx context.Context, dgst digest.Digest) (content.Info, error) ***REMOVED***
	ns, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return content.Info***REMOVED******REMOVED***, err
	***REMOVED***

	var info content.Info
	if err := view(ctx, cs.db, func(tx *bolt.Tx) error ***REMOVED***
		bkt := getBlobBucket(tx, ns, dgst)
		if bkt == nil ***REMOVED***
			return errors.Wrapf(errdefs.ErrNotFound, "content digest %v", dgst)
		***REMOVED***

		info.Digest = dgst
		return readInfo(&info, bkt)
	***REMOVED***); err != nil ***REMOVED***
		return content.Info***REMOVED******REMOVED***, err
	***REMOVED***

	return info, nil
***REMOVED***

func (cs *contentStore) Update(ctx context.Context, info content.Info, fieldpaths ...string) (content.Info, error) ***REMOVED***
	ns, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return content.Info***REMOVED******REMOVED***, err
	***REMOVED***

	cs.l.RLock()
	defer cs.l.RUnlock()

	updated := content.Info***REMOVED***
		Digest: info.Digest,
	***REMOVED***
	if err := update(ctx, cs.db, func(tx *bolt.Tx) error ***REMOVED***
		bkt := getBlobBucket(tx, ns, info.Digest)
		if bkt == nil ***REMOVED***
			return errors.Wrapf(errdefs.ErrNotFound, "content digest %v", info.Digest)
		***REMOVED***

		if err := readInfo(&updated, bkt); err != nil ***REMOVED***
			return errors.Wrapf(err, "info %q", info.Digest)
		***REMOVED***

		if len(fieldpaths) > 0 ***REMOVED***
			for _, path := range fieldpaths ***REMOVED***
				if strings.HasPrefix(path, "labels.") ***REMOVED***
					if updated.Labels == nil ***REMOVED***
						updated.Labels = map[string]string***REMOVED******REMOVED***
					***REMOVED***

					key := strings.TrimPrefix(path, "labels.")
					updated.Labels[key] = info.Labels[key]
					continue
				***REMOVED***

				switch path ***REMOVED***
				case "labels":
					updated.Labels = info.Labels
				default:
					return errors.Wrapf(errdefs.ErrInvalidArgument, "cannot update %q field on content info %q", path, info.Digest)
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// Set mutable fields
			updated.Labels = info.Labels
		***REMOVED***
		if err := validateInfo(&updated); err != nil ***REMOVED***
			return err
		***REMOVED***

		updated.UpdatedAt = time.Now().UTC()
		return writeInfo(&updated, bkt)
	***REMOVED***); err != nil ***REMOVED***
		return content.Info***REMOVED******REMOVED***, err
	***REMOVED***
	return updated, nil
***REMOVED***

func (cs *contentStore) Walk(ctx context.Context, fn content.WalkFunc, fs ...string) error ***REMOVED***
	ns, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	filter, err := filters.ParseAll(fs...)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// TODO: Batch results to keep from reading all info into memory
	var infos []content.Info
	if err := view(ctx, cs.db, func(tx *bolt.Tx) error ***REMOVED***
		bkt := getBlobsBucket(tx, ns)
		if bkt == nil ***REMOVED***
			return nil
		***REMOVED***

		return bkt.ForEach(func(k, v []byte) error ***REMOVED***
			dgst, err := digest.Parse(string(k))
			if err != nil ***REMOVED***
				// Not a digest, skip
				return nil
			***REMOVED***
			bbkt := bkt.Bucket(k)
			if bbkt == nil ***REMOVED***
				return nil
			***REMOVED***
			info := content.Info***REMOVED***
				Digest: dgst,
			***REMOVED***
			if err := readInfo(&info, bkt.Bucket(k)); err != nil ***REMOVED***
				return err
			***REMOVED***
			if filter.Match(adaptContentInfo(info)) ***REMOVED***
				infos = append(infos, info)
			***REMOVED***
			return nil
		***REMOVED***)
	***REMOVED***); err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, info := range infos ***REMOVED***
		if err := fn(info); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (cs *contentStore) Delete(ctx context.Context, dgst digest.Digest) error ***REMOVED***
	ns, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	cs.l.RLock()
	defer cs.l.RUnlock()

	return update(ctx, cs.db, func(tx *bolt.Tx) error ***REMOVED***
		bkt := getBlobBucket(tx, ns, dgst)
		if bkt == nil ***REMOVED***
			return errors.Wrapf(errdefs.ErrNotFound, "content digest %v", dgst)
		***REMOVED***

		if err := getBlobsBucket(tx, ns).DeleteBucket([]byte(dgst.String())); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := removeContentLease(ctx, tx, dgst); err != nil ***REMOVED***
			return err
		***REMOVED***

		// Mark content store as dirty for triggering garbage collection
		cs.db.dirtyL.Lock()
		cs.db.dirtyCS = true
		cs.db.dirtyL.Unlock()

		return nil
	***REMOVED***)
***REMOVED***

func (cs *contentStore) ListStatuses(ctx context.Context, fs ...string) ([]content.Status, error) ***REMOVED***
	ns, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	filter, err := filters.ParseAll(fs...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	brefs := map[string]string***REMOVED******REMOVED***
	if err := view(ctx, cs.db, func(tx *bolt.Tx) error ***REMOVED***
		bkt := getIngestBucket(tx, ns)
		if bkt == nil ***REMOVED***
			return nil
		***REMOVED***

		return bkt.ForEach(func(k, v []byte) error ***REMOVED***
			// TODO(dmcgowan): match name and potentially labels here
			brefs[string(k)] = string(v)
			return nil
		***REMOVED***)
	***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	statuses := make([]content.Status, 0, len(brefs))
	for k, bref := range brefs ***REMOVED***
		status, err := cs.Store.Status(ctx, bref)
		if err != nil ***REMOVED***
			if errdefs.IsNotFound(err) ***REMOVED***
				continue
			***REMOVED***
			return nil, err
		***REMOVED***
		status.Ref = k

		if filter.Match(adaptContentStatus(status)) ***REMOVED***
			statuses = append(statuses, status)
		***REMOVED***
	***REMOVED***

	return statuses, nil

***REMOVED***

func getRef(tx *bolt.Tx, ns, ref string) string ***REMOVED***
	bkt := getIngestBucket(tx, ns)
	if bkt == nil ***REMOVED***
		return ""
	***REMOVED***
	v := bkt.Get([]byte(ref))
	if len(v) == 0 ***REMOVED***
		return ""
	***REMOVED***
	return string(v)
***REMOVED***

func (cs *contentStore) Status(ctx context.Context, ref string) (content.Status, error) ***REMOVED***
	ns, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return content.Status***REMOVED******REMOVED***, err
	***REMOVED***

	var bref string
	if err := view(ctx, cs.db, func(tx *bolt.Tx) error ***REMOVED***
		bref = getRef(tx, ns, ref)
		if bref == "" ***REMOVED***
			return errors.Wrapf(errdefs.ErrNotFound, "reference %v", ref)
		***REMOVED***

		return nil
	***REMOVED***); err != nil ***REMOVED***
		return content.Status***REMOVED******REMOVED***, err
	***REMOVED***

	st, err := cs.Store.Status(ctx, bref)
	if err != nil ***REMOVED***
		return content.Status***REMOVED******REMOVED***, err
	***REMOVED***
	st.Ref = ref
	return st, nil
***REMOVED***

func (cs *contentStore) Abort(ctx context.Context, ref string) error ***REMOVED***
	ns, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	cs.l.RLock()
	defer cs.l.RUnlock()

	return update(ctx, cs.db, func(tx *bolt.Tx) error ***REMOVED***
		bkt := getIngestBucket(tx, ns)
		if bkt == nil ***REMOVED***
			return errors.Wrapf(errdefs.ErrNotFound, "reference %v", ref)
		***REMOVED***
		bref := string(bkt.Get([]byte(ref)))
		if bref == "" ***REMOVED***
			return errors.Wrapf(errdefs.ErrNotFound, "reference %v", ref)
		***REMOVED***
		if err := bkt.Delete([]byte(ref)); err != nil ***REMOVED***
			return err
		***REMOVED***

		return cs.Store.Abort(ctx, bref)
	***REMOVED***)

***REMOVED***

func (cs *contentStore) Writer(ctx context.Context, ref string, size int64, expected digest.Digest) (content.Writer, error) ***REMOVED***
	ns, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	cs.l.RLock()
	defer cs.l.RUnlock()

	var w content.Writer
	if err := update(ctx, cs.db, func(tx *bolt.Tx) error ***REMOVED***
		if expected != "" ***REMOVED***
			cbkt := getBlobBucket(tx, ns, expected)
			if cbkt != nil ***REMOVED***
				return errors.Wrapf(errdefs.ErrAlreadyExists, "content %v", expected)
			***REMOVED***
		***REMOVED***

		bkt, err := createIngestBucket(tx, ns)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		var (
			bref  string
			brefb = bkt.Get([]byte(ref))
		)

		if brefb == nil ***REMOVED***
			sid, err := bkt.NextSequence()
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			bref = createKey(sid, ns, ref)
			if err := bkt.Put([]byte(ref), []byte(bref)); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			bref = string(brefb)
		***REMOVED***

		// Do not use the passed in expected value here since it was
		// already checked against the user metadata. If the content
		// store has the content, it must still be written before
		// linked into the given namespace. It is possible in the future
		// to allow content which exists in content store but not
		// namespace to be linked here and returned an exist error, but
		// this would require more configuration to make secure.
		w, err = cs.Store.Writer(ctx, bref, size, "")
		return err
	***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// TODO: keep the expected in the writer to use on commit
	// when no expected is provided there.
	return &namespacedWriter***REMOVED***
		Writer:    w,
		ref:       ref,
		namespace: ns,
		db:        cs.db,
		l:         &cs.l,
	***REMOVED***, nil
***REMOVED***

type namespacedWriter struct ***REMOVED***
	content.Writer
	ref       string
	namespace string
	db        transactor
	l         *sync.RWMutex
***REMOVED***

func (nw *namespacedWriter) Commit(ctx context.Context, size int64, expected digest.Digest, opts ...content.Opt) error ***REMOVED***
	nw.l.RLock()
	defer nw.l.RUnlock()

	return update(ctx, nw.db, func(tx *bolt.Tx) error ***REMOVED***
		bkt := getIngestBucket(tx, nw.namespace)
		if bkt != nil ***REMOVED***
			if err := bkt.Delete([]byte(nw.ref)); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		dgst, err := nw.commit(ctx, tx, size, expected, opts...)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		return addContentLease(ctx, tx, dgst)
	***REMOVED***)
***REMOVED***

func (nw *namespacedWriter) commit(ctx context.Context, tx *bolt.Tx, size int64, expected digest.Digest, opts ...content.Opt) (digest.Digest, error) ***REMOVED***
	var base content.Info
	for _, opt := range opts ***REMOVED***
		if err := opt(&base); err != nil ***REMOVED***
			return "", err
		***REMOVED***
	***REMOVED***
	if err := validateInfo(&base); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	status, err := nw.Writer.Status()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if size != 0 && size != status.Offset ***REMOVED***
		return "", errors.Errorf("%q failed size validation: %v != %v", nw.ref, status.Offset, size)
	***REMOVED***
	size = status.Offset

	actual := nw.Writer.Digest()

	if err := nw.Writer.Commit(ctx, size, expected); err != nil ***REMOVED***
		if !errdefs.IsAlreadyExists(err) ***REMOVED***
			return "", err
		***REMOVED***
		if getBlobBucket(tx, nw.namespace, actual) != nil ***REMOVED***
			return "", errors.Wrapf(errdefs.ErrAlreadyExists, "content %v", actual)
		***REMOVED***
	***REMOVED***

	bkt, err := createBlobBucket(tx, nw.namespace, actual)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	commitTime := time.Now().UTC()

	sizeEncoded, err := encodeInt(size)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	if err := boltutil.WriteTimestamps(bkt, commitTime, commitTime); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if err := boltutil.WriteLabels(bkt, base.Labels); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return actual, bkt.Put(bucketKeySize, sizeEncoded)
***REMOVED***

func (nw *namespacedWriter) Status() (content.Status, error) ***REMOVED***
	st, err := nw.Writer.Status()
	if err == nil ***REMOVED***
		st.Ref = nw.ref
	***REMOVED***
	return st, err
***REMOVED***

func (cs *contentStore) ReaderAt(ctx context.Context, dgst digest.Digest) (content.ReaderAt, error) ***REMOVED***
	if err := cs.checkAccess(ctx, dgst); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return cs.Store.ReaderAt(ctx, dgst)
***REMOVED***

func (cs *contentStore) checkAccess(ctx context.Context, dgst digest.Digest) error ***REMOVED***
	ns, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return view(ctx, cs.db, func(tx *bolt.Tx) error ***REMOVED***
		bkt := getBlobBucket(tx, ns, dgst)
		if bkt == nil ***REMOVED***
			return errors.Wrapf(errdefs.ErrNotFound, "content digest %v", dgst)
		***REMOVED***
		return nil
	***REMOVED***)
***REMOVED***

func validateInfo(info *content.Info) error ***REMOVED***
	for k, v := range info.Labels ***REMOVED***
		if err := labels.Validate(k, v); err == nil ***REMOVED***
			return errors.Wrapf(err, "info.Labels")
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func readInfo(info *content.Info, bkt *bolt.Bucket) error ***REMOVED***
	if err := boltutil.ReadTimestamps(bkt, &info.CreatedAt, &info.UpdatedAt); err != nil ***REMOVED***
		return err
	***REMOVED***

	labels, err := boltutil.ReadLabels(bkt)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	info.Labels = labels

	if v := bkt.Get(bucketKeySize); len(v) > 0 ***REMOVED***
		info.Size, _ = binary.Varint(v)
	***REMOVED***

	return nil
***REMOVED***

func writeInfo(info *content.Info, bkt *bolt.Bucket) error ***REMOVED***
	if err := boltutil.WriteTimestamps(bkt, info.CreatedAt, info.UpdatedAt); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := boltutil.WriteLabels(bkt, info.Labels); err != nil ***REMOVED***
		return errors.Wrapf(err, "writing labels for info %v", info.Digest)
	***REMOVED***

	// Write size
	sizeEncoded, err := encodeInt(info.Size)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return bkt.Put(bucketKeySize, sizeEncoded)
***REMOVED***

func (cs *contentStore) garbageCollect(ctx context.Context) (d time.Duration, err error) ***REMOVED***
	cs.l.Lock()
	t1 := time.Now()
	defer func() ***REMOVED***
		if err == nil ***REMOVED***
			d = time.Now().Sub(t1)
		***REMOVED***
		cs.l.Unlock()
	***REMOVED***()

	seen := map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	if err := cs.db.View(func(tx *bolt.Tx) error ***REMOVED***
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

			cbkt := v1bkt.Bucket(k).Bucket(bucketKeyObjectContent)
			if cbkt == nil ***REMOVED***
				continue
			***REMOVED***
			bbkt := cbkt.Bucket(bucketKeyObjectBlob)
			if err := bbkt.ForEach(func(ck, cv []byte) error ***REMOVED***
				if cv == nil ***REMOVED***
					seen[string(ck)] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
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

	err = cs.Store.Walk(ctx, func(info content.Info) error ***REMOVED***
		if _, ok := seen[info.Digest.String()]; !ok ***REMOVED***
			if err := cs.Store.Delete(ctx, info.Digest); err != nil ***REMOVED***
				return err
			***REMOVED***
			log.G(ctx).WithField("digest", info.Digest).Debug("removed content")
		***REMOVED***
		return nil
	***REMOVED***)
	return
***REMOVED***
