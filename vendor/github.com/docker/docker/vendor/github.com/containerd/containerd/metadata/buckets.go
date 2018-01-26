package metadata

import (
	"github.com/boltdb/bolt"
	digest "github.com/opencontainers/go-digest"
)

// The layout where a "/" delineates a bucket is desribed in the following
// section. Please try to follow this as closely as possible when adding
// functionality. We can bolster this with helpers and more structure if that
// becomes an issue.
//
// Generically, we try to do the following:
//
// 	<version>/<namespace>/<object>/<key> -> <field>
//
// version: Currently, this is "v1". Additions can be made to v1 in a backwards
// compatible way. If the layout changes, a new version must be made, along
// with a migration.
//
// namespace: the namespace to which this object belongs.
//
// object: defines which object set is stored in the bucket. There are two
// special objects, "labels" and "indexes". The "labels" bucket stores the
// labels for the parent namespace. The "indexes" object is reserved for
// indexing objects, if we require in the future.
//
// key: object-specific key identifying the storage bucket for the objects
// contents.
var (
	bucketKeyVersion          = []byte(schemaVersion)
	bucketKeyDBVersion        = []byte("version")    // stores the version of the schema
	bucketKeyObjectLabels     = []byte("labels")     // stores the labels for a namespace.
	bucketKeyObjectImages     = []byte("images")     // stores image objects
	bucketKeyObjectContainers = []byte("containers") // stores container objects
	bucketKeyObjectSnapshots  = []byte("snapshots")  // stores snapshot references
	bucketKeyObjectContent    = []byte("content")    // stores content references
	bucketKeyObjectBlob       = []byte("blob")       // stores content links
	bucketKeyObjectIngest     = []byte("ingest")     // stores ingest links
	bucketKeyObjectLeases     = []byte("leases")     // stores leases

	bucketKeyDigest      = []byte("digest")
	bucketKeyMediaType   = []byte("mediatype")
	bucketKeySize        = []byte("size")
	bucketKeyImage       = []byte("image")
	bucketKeyRuntime     = []byte("runtime")
	bucketKeyName        = []byte("name")
	bucketKeyParent      = []byte("parent")
	bucketKeyChildren    = []byte("children")
	bucketKeyOptions     = []byte("options")
	bucketKeySpec        = []byte("spec")
	bucketKeySnapshotKey = []byte("snapshotKey")
	bucketKeySnapshotter = []byte("snapshotter")
	bucketKeyTarget      = []byte("target")
	bucketKeyExtensions  = []byte("extensions")
	bucketKeyCreatedAt   = []byte("createdat")
)

func getBucket(tx *bolt.Tx, keys ...[]byte) *bolt.Bucket ***REMOVED***
	bkt := tx.Bucket(keys[0])

	for _, key := range keys[1:] ***REMOVED***
		if bkt == nil ***REMOVED***
			break
		***REMOVED***
		bkt = bkt.Bucket(key)
	***REMOVED***

	return bkt
***REMOVED***

func createBucketIfNotExists(tx *bolt.Tx, keys ...[]byte) (*bolt.Bucket, error) ***REMOVED***
	bkt, err := tx.CreateBucketIfNotExists(keys[0])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	for _, key := range keys[1:] ***REMOVED***
		bkt, err = bkt.CreateBucketIfNotExists(key)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return bkt, nil
***REMOVED***

func namespaceLabelsBucketPath(namespace string) [][]byte ***REMOVED***
	return [][]byte***REMOVED***bucketKeyVersion, []byte(namespace), bucketKeyObjectLabels***REMOVED***
***REMOVED***

func withNamespacesLabelsBucket(tx *bolt.Tx, namespace string, fn func(bkt *bolt.Bucket) error) error ***REMOVED***
	bkt, err := createBucketIfNotExists(tx, namespaceLabelsBucketPath(namespace)...)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return fn(bkt)
***REMOVED***

func getNamespaceLabelsBucket(tx *bolt.Tx, namespace string) *bolt.Bucket ***REMOVED***
	return getBucket(tx, namespaceLabelsBucketPath(namespace)...)
***REMOVED***

func imagesBucketPath(namespace string) [][]byte ***REMOVED***
	return [][]byte***REMOVED***bucketKeyVersion, []byte(namespace), bucketKeyObjectImages***REMOVED***
***REMOVED***

func withImagesBucket(tx *bolt.Tx, namespace string, fn func(bkt *bolt.Bucket) error) error ***REMOVED***
	bkt, err := createBucketIfNotExists(tx, imagesBucketPath(namespace)...)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return fn(bkt)
***REMOVED***

func getImagesBucket(tx *bolt.Tx, namespace string) *bolt.Bucket ***REMOVED***
	return getBucket(tx, imagesBucketPath(namespace)...)
***REMOVED***

func createContainersBucket(tx *bolt.Tx, namespace string) (*bolt.Bucket, error) ***REMOVED***
	bkt, err := createBucketIfNotExists(tx, bucketKeyVersion, []byte(namespace), bucketKeyObjectContainers)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return bkt, nil
***REMOVED***

func getContainersBucket(tx *bolt.Tx, namespace string) *bolt.Bucket ***REMOVED***
	return getBucket(tx, bucketKeyVersion, []byte(namespace), bucketKeyObjectContainers)
***REMOVED***

func getContainerBucket(tx *bolt.Tx, namespace, id string) *bolt.Bucket ***REMOVED***
	return getBucket(tx, bucketKeyVersion, []byte(namespace), bucketKeyObjectContainers, []byte(id))
***REMOVED***

func createSnapshotterBucket(tx *bolt.Tx, namespace, snapshotter string) (*bolt.Bucket, error) ***REMOVED***
	bkt, err := createBucketIfNotExists(tx, bucketKeyVersion, []byte(namespace), bucketKeyObjectSnapshots, []byte(snapshotter))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return bkt, nil
***REMOVED***

func getSnapshotterBucket(tx *bolt.Tx, namespace, snapshotter string) *bolt.Bucket ***REMOVED***
	return getBucket(tx, bucketKeyVersion, []byte(namespace), bucketKeyObjectSnapshots, []byte(snapshotter))
***REMOVED***

func createBlobBucket(tx *bolt.Tx, namespace string, dgst digest.Digest) (*bolt.Bucket, error) ***REMOVED***
	bkt, err := createBucketIfNotExists(tx, bucketKeyVersion, []byte(namespace), bucketKeyObjectContent, bucketKeyObjectBlob, []byte(dgst.String()))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return bkt, nil
***REMOVED***

func getBlobsBucket(tx *bolt.Tx, namespace string) *bolt.Bucket ***REMOVED***
	return getBucket(tx, bucketKeyVersion, []byte(namespace), bucketKeyObjectContent, bucketKeyObjectBlob)
***REMOVED***

func getBlobBucket(tx *bolt.Tx, namespace string, dgst digest.Digest) *bolt.Bucket ***REMOVED***
	return getBucket(tx, bucketKeyVersion, []byte(namespace), bucketKeyObjectContent, bucketKeyObjectBlob, []byte(dgst.String()))
***REMOVED***

func createIngestBucket(tx *bolt.Tx, namespace string) (*bolt.Bucket, error) ***REMOVED***
	bkt, err := createBucketIfNotExists(tx, bucketKeyVersion, []byte(namespace), bucketKeyObjectContent, bucketKeyObjectIngest)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return bkt, nil
***REMOVED***

func getIngestBucket(tx *bolt.Tx, namespace string) *bolt.Bucket ***REMOVED***
	return getBucket(tx, bucketKeyVersion, []byte(namespace), bucketKeyObjectContent, bucketKeyObjectIngest)
***REMOVED***
