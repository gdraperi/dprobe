package boltutil

import (
	"time"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

var (
	bucketKeyLabels    = []byte("labels")
	bucketKeyCreatedAt = []byte("createdat")
	bucketKeyUpdatedAt = []byte("updatedat")
)

// ReadLabels reads the labels key from the bucket
// Uses the key "labels"
func ReadLabels(bkt *bolt.Bucket) (map[string]string, error) ***REMOVED***
	lbkt := bkt.Bucket(bucketKeyLabels)
	if lbkt == nil ***REMOVED***
		return nil, nil
	***REMOVED***
	labels := map[string]string***REMOVED******REMOVED***
	if err := lbkt.ForEach(func(k, v []byte) error ***REMOVED***
		labels[string(k)] = string(v)
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return labels, nil
***REMOVED***

// WriteLabels will write a new labels bucket to the provided bucket at key
// bucketKeyLabels, replacing the contents of the bucket with the provided map.
//
// The provide map labels will be modified to have the final contents of the
// bucket. Typically, this removes zero-value entries.
// Uses the key "labels"
func WriteLabels(bkt *bolt.Bucket, labels map[string]string) error ***REMOVED***
	// Remove existing labels to keep from merging
	if lbkt := bkt.Bucket(bucketKeyLabels); lbkt != nil ***REMOVED***
		if err := bkt.DeleteBucket(bucketKeyLabels); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if len(labels) == 0 ***REMOVED***
		return nil
	***REMOVED***

	lbkt, err := bkt.CreateBucket(bucketKeyLabels)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for k, v := range labels ***REMOVED***
		if v == "" ***REMOVED***
			delete(labels, k) // remove since we don't actually set it
			continue
		***REMOVED***

		if err := lbkt.Put([]byte(k), []byte(v)); err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to set label %q=%q", k, v)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// ReadTimestamps reads created and updated timestamps from a bucket.
// Uses keys "createdat" and "updatedat"
func ReadTimestamps(bkt *bolt.Bucket, created, updated *time.Time) error ***REMOVED***
	for _, f := range []struct ***REMOVED***
		b []byte
		t *time.Time
	***REMOVED******REMOVED***
		***REMOVED***bucketKeyCreatedAt, created***REMOVED***,
		***REMOVED***bucketKeyUpdatedAt, updated***REMOVED***,
	***REMOVED*** ***REMOVED***
		v := bkt.Get(f.b)
		if v != nil ***REMOVED***
			if err := f.t.UnmarshalBinary(v); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// WriteTimestamps writes created and updated timestamps to a bucket.
// Uses keys "createdat" and "updatedat"
func WriteTimestamps(bkt *bolt.Bucket, created, updated time.Time) error ***REMOVED***
	createdAt, err := created.MarshalBinary()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	updatedAt, err := updated.MarshalBinary()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, v := range [][2][]byte***REMOVED***
		***REMOVED***bucketKeyCreatedAt, createdAt***REMOVED***,
		***REMOVED***bucketKeyUpdatedAt, updatedAt***REMOVED***,
	***REMOVED*** ***REMOVED***
		if err := bkt.Put(v[0], v[1]); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
