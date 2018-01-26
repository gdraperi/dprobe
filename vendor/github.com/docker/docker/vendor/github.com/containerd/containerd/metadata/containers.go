package metadata

import (
	"context"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/filters"
	"github.com/containerd/containerd/identifiers"
	"github.com/containerd/containerd/labels"
	"github.com/containerd/containerd/metadata/boltutil"
	"github.com/containerd/containerd/namespaces"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
)

type containerStore struct ***REMOVED***
	tx *bolt.Tx
***REMOVED***

// NewContainerStore returns a Store backed by an underlying bolt DB
func NewContainerStore(tx *bolt.Tx) containers.Store ***REMOVED***
	return &containerStore***REMOVED***
		tx: tx,
	***REMOVED***
***REMOVED***

func (s *containerStore) Get(ctx context.Context, id string) (containers.Container, error) ***REMOVED***
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return containers.Container***REMOVED******REMOVED***, err
	***REMOVED***

	bkt := getContainerBucket(s.tx, namespace, id)
	if bkt == nil ***REMOVED***
		return containers.Container***REMOVED******REMOVED***, errors.Wrapf(errdefs.ErrNotFound, "container %q in namespace %q", id, namespace)
	***REMOVED***

	container := containers.Container***REMOVED***ID: id***REMOVED***
	if err := readContainer(&container, bkt); err != nil ***REMOVED***
		return containers.Container***REMOVED******REMOVED***, errors.Wrapf(err, "failed to read container %q", id)
	***REMOVED***

	return container, nil
***REMOVED***

func (s *containerStore) List(ctx context.Context, fs ...string) ([]containers.Container, error) ***REMOVED***
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	filter, err := filters.ParseAll(fs...)
	if err != nil ***REMOVED***
		return nil, errors.Wrapf(errdefs.ErrInvalidArgument, err.Error())
	***REMOVED***

	bkt := getContainersBucket(s.tx, namespace)
	if bkt == nil ***REMOVED***
		return nil, nil // empty store
	***REMOVED***

	var m []containers.Container
	if err := bkt.ForEach(func(k, v []byte) error ***REMOVED***
		cbkt := bkt.Bucket(k)
		if cbkt == nil ***REMOVED***
			return nil
		***REMOVED***
		container := containers.Container***REMOVED***ID: string(k)***REMOVED***

		if err := readContainer(&container, cbkt); err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to read container %q", string(k))
		***REMOVED***

		if filter.Match(adaptContainer(container)) ***REMOVED***
			m = append(m, container)
		***REMOVED***
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return m, nil
***REMOVED***

func (s *containerStore) Create(ctx context.Context, container containers.Container) (containers.Container, error) ***REMOVED***
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return containers.Container***REMOVED******REMOVED***, err
	***REMOVED***

	if err := validateContainer(&container); err != nil ***REMOVED***
		return containers.Container***REMOVED******REMOVED***, errors.Wrap(err, "create container failed validation")
	***REMOVED***

	bkt, err := createContainersBucket(s.tx, namespace)
	if err != nil ***REMOVED***
		return containers.Container***REMOVED******REMOVED***, err
	***REMOVED***

	cbkt, err := bkt.CreateBucket([]byte(container.ID))
	if err != nil ***REMOVED***
		if err == bolt.ErrBucketExists ***REMOVED***
			err = errors.Wrapf(errdefs.ErrAlreadyExists, "container %q", container.ID)
		***REMOVED***
		return containers.Container***REMOVED******REMOVED***, err
	***REMOVED***

	container.CreatedAt = time.Now().UTC()
	container.UpdatedAt = container.CreatedAt
	if err := writeContainer(cbkt, &container); err != nil ***REMOVED***
		return containers.Container***REMOVED******REMOVED***, errors.Wrapf(err, "failed to write container %q", container.ID)
	***REMOVED***

	return container, nil
***REMOVED***

func (s *containerStore) Update(ctx context.Context, container containers.Container, fieldpaths ...string) (containers.Container, error) ***REMOVED***
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return containers.Container***REMOVED******REMOVED***, err
	***REMOVED***

	if container.ID == "" ***REMOVED***
		return containers.Container***REMOVED******REMOVED***, errors.Wrapf(errdefs.ErrInvalidArgument, "must specify a container id")
	***REMOVED***

	bkt := getContainersBucket(s.tx, namespace)
	if bkt == nil ***REMOVED***
		return containers.Container***REMOVED******REMOVED***, errors.Wrapf(errdefs.ErrNotFound, "cannot update container %q in namespace %q", container.ID, namespace)
	***REMOVED***

	cbkt := bkt.Bucket([]byte(container.ID))
	if cbkt == nil ***REMOVED***
		return containers.Container***REMOVED******REMOVED***, errors.Wrapf(errdefs.ErrNotFound, "container %q", container.ID)
	***REMOVED***

	var updated containers.Container
	if err := readContainer(&updated, cbkt); err != nil ***REMOVED***
		return updated, errors.Wrapf(err, "failed to read container %q", container.ID)
	***REMOVED***
	createdat := updated.CreatedAt
	updated.ID = container.ID

	if len(fieldpaths) == 0 ***REMOVED***
		// only allow updates to these field on full replace.
		fieldpaths = []string***REMOVED***"labels", "spec", "extensions"***REMOVED***

		// Fields that are immutable must cause an error when no field paths
		// are provided. This allows these fields to become mutable in the
		// future.
		if updated.Image != container.Image ***REMOVED***
			return containers.Container***REMOVED******REMOVED***, errors.Wrapf(errdefs.ErrInvalidArgument, "container.Image field is immutable")
		***REMOVED***

		if updated.SnapshotKey != container.SnapshotKey ***REMOVED***
			return containers.Container***REMOVED******REMOVED***, errors.Wrapf(errdefs.ErrInvalidArgument, "container.SnapshotKey field is immutable")
		***REMOVED***

		if updated.Snapshotter != container.Snapshotter ***REMOVED***
			return containers.Container***REMOVED******REMOVED***, errors.Wrapf(errdefs.ErrInvalidArgument, "container.Snapshotter field is immutable")
		***REMOVED***

		if updated.Runtime.Name != container.Runtime.Name ***REMOVED***
			return containers.Container***REMOVED******REMOVED***, errors.Wrapf(errdefs.ErrInvalidArgument, "container.Runtime.Name field is immutable")
		***REMOVED***
	***REMOVED***

	// apply the field mask. If you update this code, you better follow the
	// field mask rules in field_mask.proto. If you don't know what this
	// is, do not update this code.
	for _, path := range fieldpaths ***REMOVED***
		if strings.HasPrefix(path, "labels.") ***REMOVED***
			if updated.Labels == nil ***REMOVED***
				updated.Labels = map[string]string***REMOVED******REMOVED***
			***REMOVED***
			key := strings.TrimPrefix(path, "labels.")
			updated.Labels[key] = container.Labels[key]
			continue
		***REMOVED***

		if strings.HasPrefix(path, "extensions.") ***REMOVED***
			if updated.Extensions == nil ***REMOVED***
				updated.Extensions = map[string]types.Any***REMOVED******REMOVED***
			***REMOVED***
			key := strings.TrimPrefix(path, "extensions.")
			updated.Extensions[key] = container.Extensions[key]
			continue
		***REMOVED***

		switch path ***REMOVED***
		case "labels":
			updated.Labels = container.Labels
		case "spec":
			updated.Spec = container.Spec
		case "extensions":
			updated.Extensions = container.Extensions
		default:
			return containers.Container***REMOVED******REMOVED***, errors.Wrapf(errdefs.ErrInvalidArgument, "cannot update %q field on %q", path, container.ID)
		***REMOVED***
	***REMOVED***

	if err := validateContainer(&updated); err != nil ***REMOVED***
		return containers.Container***REMOVED******REMOVED***, errors.Wrap(err, "update failed validation")
	***REMOVED***

	updated.CreatedAt = createdat
	updated.UpdatedAt = time.Now().UTC()
	if err := writeContainer(cbkt, &updated); err != nil ***REMOVED***
		return containers.Container***REMOVED******REMOVED***, errors.Wrapf(err, "failed to write container %q", container.ID)
	***REMOVED***

	return updated, nil
***REMOVED***

func (s *containerStore) Delete(ctx context.Context, id string) error ***REMOVED***
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	bkt := getContainersBucket(s.tx, namespace)
	if bkt == nil ***REMOVED***
		return errors.Wrapf(errdefs.ErrNotFound, "cannot delete container %q in namespace %q", id, namespace)
	***REMOVED***

	if err := bkt.DeleteBucket([]byte(id)); err == bolt.ErrBucketNotFound ***REMOVED***
		return errors.Wrapf(errdefs.ErrNotFound, "container %v", id)
	***REMOVED***
	return err
***REMOVED***

func validateContainer(container *containers.Container) error ***REMOVED***
	if err := identifiers.Validate(container.ID); err != nil ***REMOVED***
		return errors.Wrap(err, "container.ID")
	***REMOVED***

	for k := range container.Extensions ***REMOVED***
		if k == "" ***REMOVED***
			return errors.Wrapf(errdefs.ErrInvalidArgument, "container.Extension keys must not be zero-length")
		***REMOVED***
	***REMOVED***

	// image has no validation
	for k, v := range container.Labels ***REMOVED***
		if err := labels.Validate(k, v); err == nil ***REMOVED***
			return errors.Wrapf(err, "containers.Labels")
		***REMOVED***
	***REMOVED***

	if container.Runtime.Name == "" ***REMOVED***
		return errors.Wrapf(errdefs.ErrInvalidArgument, "container.Runtime.Name must be set")
	***REMOVED***

	if container.Spec == nil ***REMOVED***
		return errors.Wrapf(errdefs.ErrInvalidArgument, "container.Spec must be set")
	***REMOVED***

	if container.SnapshotKey != "" && container.Snapshotter == "" ***REMOVED***
		return errors.Wrapf(errdefs.ErrInvalidArgument, "container.Snapshotter must be set if container.SnapshotKey is set")
	***REMOVED***

	return nil
***REMOVED***

func readContainer(container *containers.Container, bkt *bolt.Bucket) error ***REMOVED***
	labels, err := boltutil.ReadLabels(bkt)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	container.Labels = labels

	if err := boltutil.ReadTimestamps(bkt, &container.CreatedAt, &container.UpdatedAt); err != nil ***REMOVED***
		return err
	***REMOVED***

	return bkt.ForEach(func(k, v []byte) error ***REMOVED***
		switch string(k) ***REMOVED***
		case string(bucketKeyImage):
			container.Image = string(v)
		case string(bucketKeyRuntime):
			rbkt := bkt.Bucket(bucketKeyRuntime)
			if rbkt == nil ***REMOVED***
				return nil // skip runtime. should be an error?
			***REMOVED***

			n := rbkt.Get(bucketKeyName)
			if n != nil ***REMOVED***
				container.Runtime.Name = string(n)
			***REMOVED***

			obkt := rbkt.Get(bucketKeyOptions)
			if obkt == nil ***REMOVED***
				return nil
			***REMOVED***

			var any types.Any
			if err := proto.Unmarshal(obkt, &any); err != nil ***REMOVED***
				return err
			***REMOVED***
			container.Runtime.Options = &any
		case string(bucketKeySpec):
			var any types.Any
			if err := proto.Unmarshal(v, &any); err != nil ***REMOVED***
				return err
			***REMOVED***
			container.Spec = &any
		case string(bucketKeySnapshotKey):
			container.SnapshotKey = string(v)
		case string(bucketKeySnapshotter):
			container.Snapshotter = string(v)
		case string(bucketKeyExtensions):
			ebkt := bkt.Bucket(bucketKeyExtensions)
			if ebkt == nil ***REMOVED***
				return nil
			***REMOVED***

			extensions := make(map[string]types.Any)
			if err := ebkt.ForEach(func(k, v []byte) error ***REMOVED***
				var a types.Any
				if err := proto.Unmarshal(v, &a); err != nil ***REMOVED***
					return err
				***REMOVED***

				extensions[string(k)] = a
				return nil
			***REMOVED***); err != nil ***REMOVED***

				return err
			***REMOVED***

			container.Extensions = extensions
		***REMOVED***

		return nil
	***REMOVED***)
***REMOVED***

func writeContainer(bkt *bolt.Bucket, container *containers.Container) error ***REMOVED***
	if err := boltutil.WriteTimestamps(bkt, container.CreatedAt, container.UpdatedAt); err != nil ***REMOVED***
		return err
	***REMOVED***

	if container.Spec != nil ***REMOVED***
		spec, err := container.Spec.Marshal()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := bkt.Put(bucketKeySpec, spec); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	for _, v := range [][2][]byte***REMOVED***
		***REMOVED***bucketKeyImage, []byte(container.Image)***REMOVED***,
		***REMOVED***bucketKeySnapshotter, []byte(container.Snapshotter)***REMOVED***,
		***REMOVED***bucketKeySnapshotKey, []byte(container.SnapshotKey)***REMOVED***,
	***REMOVED*** ***REMOVED***
		if err := bkt.Put(v[0], v[1]); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if rbkt := bkt.Bucket(bucketKeyRuntime); rbkt != nil ***REMOVED***
		if err := bkt.DeleteBucket(bucketKeyRuntime); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	rbkt, err := bkt.CreateBucket(bucketKeyRuntime)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := rbkt.Put(bucketKeyName, []byte(container.Runtime.Name)); err != nil ***REMOVED***
		return err
	***REMOVED***

	if len(container.Extensions) > 0 ***REMOVED***
		ebkt, err := bkt.CreateBucketIfNotExists(bucketKeyExtensions)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		for name, ext := range container.Extensions ***REMOVED***
			p, err := proto.Marshal(&ext)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			if err := ebkt.Put([]byte(name), p); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if container.Runtime.Options != nil ***REMOVED***
		data, err := proto.Marshal(container.Runtime.Options)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := rbkt.Put(bucketKeyOptions, data); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return boltutil.WriteLabels(bkt, container.Labels)
***REMOVED***
