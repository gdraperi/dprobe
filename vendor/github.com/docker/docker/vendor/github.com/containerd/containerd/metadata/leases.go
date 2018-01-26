package metadata

import (
	"context"
	"time"

	"github.com/boltdb/bolt"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/leases"
	"github.com/containerd/containerd/metadata/boltutil"
	"github.com/containerd/containerd/namespaces"
	digest "github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

// Lease retains resources to prevent garbage collection before
// the resources can be fully referenced.
type Lease struct ***REMOVED***
	ID        string
	CreatedAt time.Time
	Labels    map[string]string

	Content   []string
	Snapshots map[string][]string
***REMOVED***

// LeaseManager manages the create/delete lifecyle of leases
// and also returns existing leases
type LeaseManager struct ***REMOVED***
	tx *bolt.Tx
***REMOVED***

// NewLeaseManager creates a new lease manager for managing leases using
// the provided database transaction.
func NewLeaseManager(tx *bolt.Tx) *LeaseManager ***REMOVED***
	return &LeaseManager***REMOVED***
		tx: tx,
	***REMOVED***
***REMOVED***

// Create creates a new lease using the provided lease
func (lm *LeaseManager) Create(ctx context.Context, lid string, labels map[string]string) (Lease, error) ***REMOVED***
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return Lease***REMOVED******REMOVED***, err
	***REMOVED***

	topbkt, err := createBucketIfNotExists(lm.tx, bucketKeyVersion, []byte(namespace), bucketKeyObjectLeases)
	if err != nil ***REMOVED***
		return Lease***REMOVED******REMOVED***, err
	***REMOVED***

	txbkt, err := topbkt.CreateBucket([]byte(lid))
	if err != nil ***REMOVED***
		if err == bolt.ErrBucketExists ***REMOVED***
			err = errdefs.ErrAlreadyExists
		***REMOVED***
		return Lease***REMOVED******REMOVED***, errors.Wrapf(err, "lease %q", lid)
	***REMOVED***

	t := time.Now().UTC()
	createdAt, err := t.MarshalBinary()
	if err != nil ***REMOVED***
		return Lease***REMOVED******REMOVED***, err
	***REMOVED***
	if err := txbkt.Put(bucketKeyCreatedAt, createdAt); err != nil ***REMOVED***
		return Lease***REMOVED******REMOVED***, err
	***REMOVED***

	if labels != nil ***REMOVED***
		if err := boltutil.WriteLabels(txbkt, labels); err != nil ***REMOVED***
			return Lease***REMOVED******REMOVED***, err
		***REMOVED***
	***REMOVED***

	return Lease***REMOVED***
		ID:        lid,
		CreatedAt: t,
		Labels:    labels,
	***REMOVED***, nil
***REMOVED***

// Delete delets the lease with the provided lease ID
func (lm *LeaseManager) Delete(ctx context.Context, lid string) error ***REMOVED***
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	topbkt := getBucket(lm.tx, bucketKeyVersion, []byte(namespace), bucketKeyObjectLeases)
	if topbkt == nil ***REMOVED***
		return nil
	***REMOVED***
	if err := topbkt.DeleteBucket([]byte(lid)); err != nil && err != bolt.ErrBucketNotFound ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// List lists all active leases
func (lm *LeaseManager) List(ctx context.Context, includeResources bool, filter ...string) ([]Lease, error) ***REMOVED***
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var leases []Lease

	topbkt := getBucket(lm.tx, bucketKeyVersion, []byte(namespace), bucketKeyObjectLeases)
	if topbkt == nil ***REMOVED***
		return leases, nil
	***REMOVED***

	if err := topbkt.ForEach(func(k, v []byte) error ***REMOVED***
		if v != nil ***REMOVED***
			return nil
		***REMOVED***
		txbkt := topbkt.Bucket(k)

		l := Lease***REMOVED***
			ID: string(k),
		***REMOVED***

		if v := txbkt.Get(bucketKeyCreatedAt); v != nil ***REMOVED***
			t := &l.CreatedAt
			if err := t.UnmarshalBinary(v); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		labels, err := boltutil.ReadLabels(txbkt)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		l.Labels = labels

		// TODO: Read Snapshots
		// TODO: Read Content

		leases = append(leases, l)

		return nil
	***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return leases, nil
***REMOVED***

func addSnapshotLease(ctx context.Context, tx *bolt.Tx, snapshotter, key string) error ***REMOVED***
	lid, ok := leases.Lease(ctx)
	if !ok ***REMOVED***
		return nil
	***REMOVED***

	namespace, ok := namespaces.Namespace(ctx)
	if !ok ***REMOVED***
		panic("namespace must already be checked")
	***REMOVED***

	bkt := getBucket(tx, bucketKeyVersion, []byte(namespace), bucketKeyObjectLeases, []byte(lid))
	if bkt == nil ***REMOVED***
		return errors.Wrap(errdefs.ErrNotFound, "lease does not exist")
	***REMOVED***

	bkt, err := bkt.CreateBucketIfNotExists(bucketKeyObjectSnapshots)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	bkt, err = bkt.CreateBucketIfNotExists([]byte(snapshotter))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return bkt.Put([]byte(key), nil)
***REMOVED***

func removeSnapshotLease(ctx context.Context, tx *bolt.Tx, snapshotter, key string) error ***REMOVED***
	lid, ok := leases.Lease(ctx)
	if !ok ***REMOVED***
		return nil
	***REMOVED***

	namespace, ok := namespaces.Namespace(ctx)
	if !ok ***REMOVED***
		panic("namespace must already be checked")
	***REMOVED***

	bkt := getBucket(tx, bucketKeyVersion, []byte(namespace), bucketKeyObjectLeases, []byte(lid), bucketKeyObjectSnapshots, []byte(snapshotter))
	if bkt == nil ***REMOVED***
		// Key does not exist so we return nil
		return nil
	***REMOVED***

	return bkt.Delete([]byte(key))
***REMOVED***

func addContentLease(ctx context.Context, tx *bolt.Tx, dgst digest.Digest) error ***REMOVED***
	lid, ok := leases.Lease(ctx)
	if !ok ***REMOVED***
		return nil
	***REMOVED***

	namespace, ok := namespaces.Namespace(ctx)
	if !ok ***REMOVED***
		panic("namespace must already be required")
	***REMOVED***

	bkt := getBucket(tx, bucketKeyVersion, []byte(namespace), bucketKeyObjectLeases, []byte(lid))
	if bkt == nil ***REMOVED***
		return errors.Wrap(errdefs.ErrNotFound, "lease does not exist")
	***REMOVED***

	bkt, err := bkt.CreateBucketIfNotExists(bucketKeyObjectContent)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return bkt.Put([]byte(dgst.String()), nil)
***REMOVED***

func removeContentLease(ctx context.Context, tx *bolt.Tx, dgst digest.Digest) error ***REMOVED***
	lid, ok := leases.Lease(ctx)
	if !ok ***REMOVED***
		return nil
	***REMOVED***

	namespace, ok := namespaces.Namespace(ctx)
	if !ok ***REMOVED***
		panic("namespace must already be checked")
	***REMOVED***

	bkt := getBucket(tx, bucketKeyVersion, []byte(namespace), bucketKeyObjectLeases, []byte(lid), bucketKeyObjectContent)
	if bkt == nil ***REMOVED***
		// Key does not exist so we return nil
		return nil
	***REMOVED***

	return bkt.Delete([]byte(dgst.String()))
***REMOVED***
