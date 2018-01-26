package metadata

import (
	"context"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/filters"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/labels"
	"github.com/containerd/containerd/metadata/boltutil"
	"github.com/containerd/containerd/namespaces"
	digest "github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

type imageStore struct ***REMOVED***
	tx *bolt.Tx
***REMOVED***

// NewImageStore returns a store backed by a bolt DB
func NewImageStore(tx *bolt.Tx) images.Store ***REMOVED***
	return &imageStore***REMOVED***tx: tx***REMOVED***
***REMOVED***

func (s *imageStore) Get(ctx context.Context, name string) (images.Image, error) ***REMOVED***
	var image images.Image

	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return images.Image***REMOVED******REMOVED***, err
	***REMOVED***

	bkt := getImagesBucket(s.tx, namespace)
	if bkt == nil ***REMOVED***
		return images.Image***REMOVED******REMOVED***, errors.Wrapf(errdefs.ErrNotFound, "image %q", name)
	***REMOVED***

	ibkt := bkt.Bucket([]byte(name))
	if ibkt == nil ***REMOVED***
		return images.Image***REMOVED******REMOVED***, errors.Wrapf(errdefs.ErrNotFound, "image %q", name)
	***REMOVED***

	image.Name = name
	if err := readImage(&image, ibkt); err != nil ***REMOVED***
		return images.Image***REMOVED******REMOVED***, errors.Wrapf(err, "image %q", name)
	***REMOVED***

	return image, nil
***REMOVED***

func (s *imageStore) List(ctx context.Context, fs ...string) ([]images.Image, error) ***REMOVED***
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	filter, err := filters.ParseAll(fs...)
	if err != nil ***REMOVED***
		return nil, errors.Wrapf(errdefs.ErrInvalidArgument, err.Error())
	***REMOVED***

	bkt := getImagesBucket(s.tx, namespace)
	if bkt == nil ***REMOVED***
		return nil, nil // empty store
	***REMOVED***

	var m []images.Image
	if err := bkt.ForEach(func(k, v []byte) error ***REMOVED***
		var (
			image = images.Image***REMOVED***
				Name: string(k),
			***REMOVED***
			kbkt = bkt.Bucket(k)
		)

		if err := readImage(&image, kbkt); err != nil ***REMOVED***
			return err
		***REMOVED***

		if filter.Match(adaptImage(image)) ***REMOVED***
			m = append(m, image)
		***REMOVED***
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return m, nil
***REMOVED***

func (s *imageStore) Create(ctx context.Context, image images.Image) (images.Image, error) ***REMOVED***
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return images.Image***REMOVED******REMOVED***, err
	***REMOVED***

	if err := validateImage(&image); err != nil ***REMOVED***
		return images.Image***REMOVED******REMOVED***, err
	***REMOVED***

	return image, withImagesBucket(s.tx, namespace, func(bkt *bolt.Bucket) error ***REMOVED***
		ibkt, err := bkt.CreateBucket([]byte(image.Name))
		if err != nil ***REMOVED***
			if err != bolt.ErrBucketExists ***REMOVED***
				return err
			***REMOVED***

			return errors.Wrapf(errdefs.ErrAlreadyExists, "image %q", image.Name)
		***REMOVED***

		image.CreatedAt = time.Now().UTC()
		image.UpdatedAt = image.CreatedAt
		return writeImage(ibkt, &image)
	***REMOVED***)
***REMOVED***

func (s *imageStore) Update(ctx context.Context, image images.Image, fieldpaths ...string) (images.Image, error) ***REMOVED***
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return images.Image***REMOVED******REMOVED***, err
	***REMOVED***

	if image.Name == "" ***REMOVED***
		return images.Image***REMOVED******REMOVED***, errors.Wrapf(errdefs.ErrInvalidArgument, "image name is required for update")
	***REMOVED***

	var updated images.Image
	return updated, withImagesBucket(s.tx, namespace, func(bkt *bolt.Bucket) error ***REMOVED***
		ibkt := bkt.Bucket([]byte(image.Name))
		if ibkt == nil ***REMOVED***
			return errors.Wrapf(errdefs.ErrNotFound, "image %q", image.Name)
		***REMOVED***

		if err := readImage(&updated, ibkt); err != nil ***REMOVED***
			return errors.Wrapf(err, "image %q", image.Name)
		***REMOVED***
		createdat := updated.CreatedAt
		updated.Name = image.Name

		if len(fieldpaths) > 0 ***REMOVED***
			for _, path := range fieldpaths ***REMOVED***
				if strings.HasPrefix(path, "labels.") ***REMOVED***
					if updated.Labels == nil ***REMOVED***
						updated.Labels = map[string]string***REMOVED******REMOVED***
					***REMOVED***

					key := strings.TrimPrefix(path, "labels.")
					updated.Labels[key] = image.Labels[key]
					continue
				***REMOVED***

				switch path ***REMOVED***
				case "labels":
					updated.Labels = image.Labels
				case "target":
					// NOTE(stevvooe): While we allow setting individual labels, we
					// only support replacing the target as a unit, since that is
					// commonly pulled as a unit from other sources. It often doesn't
					// make sense to modify the size or digest without touching the
					// mediatype, as well, for example.
					updated.Target = image.Target
				default:
					return errors.Wrapf(errdefs.ErrInvalidArgument, "cannot update %q field on image %q", path, image.Name)
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			updated = image
		***REMOVED***

		if err := validateImage(&updated); err != nil ***REMOVED***
			return err
		***REMOVED***

		updated.CreatedAt = createdat
		updated.UpdatedAt = time.Now().UTC()
		return writeImage(ibkt, &updated)
	***REMOVED***)
***REMOVED***

func (s *imageStore) Delete(ctx context.Context, name string, opts ...images.DeleteOpt) error ***REMOVED***
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return withImagesBucket(s.tx, namespace, func(bkt *bolt.Bucket) error ***REMOVED***
		err := bkt.DeleteBucket([]byte(name))
		if err == bolt.ErrBucketNotFound ***REMOVED***
			return errors.Wrapf(errdefs.ErrNotFound, "image %q", name)
		***REMOVED***
		return err
	***REMOVED***)
***REMOVED***

func validateImage(image *images.Image) error ***REMOVED***
	if image.Name == "" ***REMOVED***
		return errors.Wrapf(errdefs.ErrInvalidArgument, "image name must not be empty")
	***REMOVED***

	for k, v := range image.Labels ***REMOVED***
		if err := labels.Validate(k, v); err != nil ***REMOVED***
			return errors.Wrapf(err, "image.Labels")
		***REMOVED***
	***REMOVED***

	return validateTarget(&image.Target)
***REMOVED***

func validateTarget(target *ocispec.Descriptor) error ***REMOVED***
	// NOTE(stevvooe): Only validate fields we actually store.

	if err := target.Digest.Validate(); err != nil ***REMOVED***
		return errors.Wrapf(errdefs.ErrInvalidArgument, "Target.Digest %q invalid: %v", target.Digest, err)
	***REMOVED***

	if target.Size <= 0 ***REMOVED***
		return errors.Wrapf(errdefs.ErrInvalidArgument, "Target.Size must be greater than zero")
	***REMOVED***

	if target.MediaType == "" ***REMOVED***
		return errors.Wrapf(errdefs.ErrInvalidArgument, "Target.MediaType must be set")
	***REMOVED***

	return nil
***REMOVED***

func readImage(image *images.Image, bkt *bolt.Bucket) error ***REMOVED***
	if err := boltutil.ReadTimestamps(bkt, &image.CreatedAt, &image.UpdatedAt); err != nil ***REMOVED***
		return err
	***REMOVED***

	labels, err := boltutil.ReadLabels(bkt)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	image.Labels = labels

	tbkt := bkt.Bucket(bucketKeyTarget)
	if tbkt == nil ***REMOVED***
		return errors.New("unable to read target bucket")
	***REMOVED***
	return tbkt.ForEach(func(k, v []byte) error ***REMOVED***
		if v == nil ***REMOVED***
			return nil // skip it? a bkt maybe?
		***REMOVED***

		// TODO(stevvooe): This is why we need to use byte values for
		// keys, rather than full arrays.
		switch string(k) ***REMOVED***
		case string(bucketKeyDigest):
			image.Target.Digest = digest.Digest(v)
		case string(bucketKeyMediaType):
			image.Target.MediaType = string(v)
		case string(bucketKeySize):
			image.Target.Size, _ = binary.Varint(v)
		***REMOVED***

		return nil
	***REMOVED***)
***REMOVED***

func writeImage(bkt *bolt.Bucket, image *images.Image) error ***REMOVED***
	if err := boltutil.WriteTimestamps(bkt, image.CreatedAt, image.UpdatedAt); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := boltutil.WriteLabels(bkt, image.Labels); err != nil ***REMOVED***
		return errors.Wrapf(err, "writing labels for image %v", image.Name)
	***REMOVED***

	// write the target bucket
	tbkt, err := bkt.CreateBucketIfNotExists([]byte(bucketKeyTarget))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	sizeEncoded, err := encodeInt(image.Target.Size)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, v := range [][2][]byte***REMOVED***
		***REMOVED***bucketKeyDigest, []byte(image.Target.Digest)***REMOVED***,
		***REMOVED***bucketKeyMediaType, []byte(image.Target.MediaType)***REMOVED***,
		***REMOVED***bucketKeySize, sizeEncoded***REMOVED***,
	***REMOVED*** ***REMOVED***
		if err := tbkt.Put(v[0], v[1]); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func encodeInt(i int64) ([]byte, error) ***REMOVED***
	var (
		buf      [binary.MaxVarintLen64]byte
		iEncoded = buf[:]
	)
	iEncoded = iEncoded[:binary.PutVarint(iEncoded, i)]

	if len(iEncoded) == 0 ***REMOVED***
		return nil, fmt.Errorf("failed encoding integer = %v", i)
	***REMOVED***
	return iEncoded, nil
***REMOVED***
