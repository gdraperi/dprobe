package metadata

import (
	"context"

	"github.com/boltdb/bolt"
	"github.com/containerd/containerd/errdefs"
	l "github.com/containerd/containerd/labels"
	"github.com/containerd/containerd/namespaces"
	"github.com/pkg/errors"
)

type namespaceStore struct ***REMOVED***
	tx *bolt.Tx
***REMOVED***

// NewNamespaceStore returns a store backed by a bolt DB
func NewNamespaceStore(tx *bolt.Tx) namespaces.Store ***REMOVED***
	return &namespaceStore***REMOVED***tx: tx***REMOVED***
***REMOVED***

func (s *namespaceStore) Create(ctx context.Context, namespace string, labels map[string]string) error ***REMOVED***
	topbkt, err := createBucketIfNotExists(s.tx, bucketKeyVersion)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := namespaces.Validate(namespace); err != nil ***REMOVED***
		return err
	***REMOVED***

	for k, v := range labels ***REMOVED***
		if err := l.Validate(k, v); err != nil ***REMOVED***
			return errors.Wrapf(err, "namespace.Labels")
		***REMOVED***
	***REMOVED***

	// provides the already exists error.
	bkt, err := topbkt.CreateBucket([]byte(namespace))
	if err != nil ***REMOVED***
		if err == bolt.ErrBucketExists ***REMOVED***
			return errors.Wrapf(errdefs.ErrAlreadyExists, "namespace %q", namespace)
		***REMOVED***

		return err
	***REMOVED***

	lbkt, err := bkt.CreateBucketIfNotExists(bucketKeyObjectLabels)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for k, v := range labels ***REMOVED***
		if err := lbkt.Put([]byte(k), []byte(v)); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (s *namespaceStore) Labels(ctx context.Context, namespace string) (map[string]string, error) ***REMOVED***
	labels := map[string]string***REMOVED******REMOVED***

	bkt := getNamespaceLabelsBucket(s.tx, namespace)
	if bkt == nil ***REMOVED***
		return labels, nil
	***REMOVED***

	if err := bkt.ForEach(func(k, v []byte) error ***REMOVED***
		labels[string(k)] = string(v)
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return labels, nil
***REMOVED***

func (s *namespaceStore) SetLabel(ctx context.Context, namespace, key, value string) error ***REMOVED***
	if err := l.Validate(key, value); err != nil ***REMOVED***
		return errors.Wrapf(err, "namespace.Labels")
	***REMOVED***

	return withNamespacesLabelsBucket(s.tx, namespace, func(bkt *bolt.Bucket) error ***REMOVED***
		if value == "" ***REMOVED***
			return bkt.Delete([]byte(key))
		***REMOVED***

		return bkt.Put([]byte(key), []byte(value))
	***REMOVED***)

***REMOVED***

func (s *namespaceStore) List(ctx context.Context) ([]string, error) ***REMOVED***
	bkt := getBucket(s.tx, bucketKeyVersion)
	if bkt == nil ***REMOVED***
		return nil, nil // no namespaces!
	***REMOVED***

	var namespaces []string
	if err := bkt.ForEach(func(k, v []byte) error ***REMOVED***
		if v != nil ***REMOVED***
			return nil // not a bucket
		***REMOVED***

		namespaces = append(namespaces, string(k))
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return namespaces, nil
***REMOVED***

func (s *namespaceStore) Delete(ctx context.Context, namespace string) error ***REMOVED***
	bkt := getBucket(s.tx, bucketKeyVersion)
	if empty, err := s.namespaceEmpty(ctx, namespace); err != nil ***REMOVED***
		return err
	***REMOVED*** else if !empty ***REMOVED***
		return errors.Wrapf(errdefs.ErrFailedPrecondition, "namespace %q must be empty", namespace)
	***REMOVED***

	if err := bkt.DeleteBucket([]byte(namespace)); err != nil ***REMOVED***
		if err == bolt.ErrBucketNotFound ***REMOVED***
			return errors.Wrapf(errdefs.ErrNotFound, "namespace %q", namespace)
		***REMOVED***

		return err
	***REMOVED***

	return nil
***REMOVED***

func (s *namespaceStore) namespaceEmpty(ctx context.Context, namespace string) (bool, error) ***REMOVED***
	ctx = namespaces.WithNamespace(ctx, namespace)

	// need to check the various object stores.

	imageStore := NewImageStore(s.tx)
	images, err := imageStore.List(ctx)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	if len(images) > 0 ***REMOVED***
		return false, nil
	***REMOVED***

	containerStore := NewContainerStore(s.tx)
	containers, err := containerStore.List(ctx)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	if len(containers) > 0 ***REMOVED***
		return false, nil
	***REMOVED***

	// TODO(stevvooe): Need to add check for content store, as well. Still need
	// to make content store namespace aware.

	return true, nil
***REMOVED***
