package store

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var volumeBucketName = []byte("volumes")

type volumeMetadata struct ***REMOVED***
	Name    string
	Driver  string
	Labels  map[string]string
	Options map[string]string
***REMOVED***

func (s *VolumeStore) setMeta(name string, meta volumeMetadata) error ***REMOVED***
	return s.db.Update(func(tx *bolt.Tx) error ***REMOVED***
		return setMeta(tx, name, meta)
	***REMOVED***)
***REMOVED***

func setMeta(tx *bolt.Tx, name string, meta volumeMetadata) error ***REMOVED***
	metaJSON, err := json.Marshal(meta)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	b := tx.Bucket(volumeBucketName)
	return errors.Wrap(b.Put([]byte(name), metaJSON), "error setting volume metadata")
***REMOVED***

func (s *VolumeStore) getMeta(name string) (volumeMetadata, error) ***REMOVED***
	var meta volumeMetadata
	err := s.db.View(func(tx *bolt.Tx) error ***REMOVED***
		return getMeta(tx, name, &meta)
	***REMOVED***)
	return meta, err
***REMOVED***

func getMeta(tx *bolt.Tx, name string, meta *volumeMetadata) error ***REMOVED***
	b := tx.Bucket(volumeBucketName)
	val := b.Get([]byte(name))
	if string(val) == "" ***REMOVED***
		return nil
	***REMOVED***
	if err := json.Unmarshal(val, meta); err != nil ***REMOVED***
		return errors.Wrap(err, "error unmarshaling volume metadata")
	***REMOVED***
	return nil
***REMOVED***

func (s *VolumeStore) removeMeta(name string) error ***REMOVED***
	return s.db.Update(func(tx *bolt.Tx) error ***REMOVED***
		return removeMeta(tx, name)
	***REMOVED***)
***REMOVED***

func removeMeta(tx *bolt.Tx, name string) error ***REMOVED***
	b := tx.Bucket(volumeBucketName)
	return errors.Wrap(b.Delete([]byte(name)), "error removing volume metadata")
***REMOVED***

// listMeta is used during restore to get the list of volume metadata
// from the on-disk database.
// Any errors that occur are only logged.
func listMeta(tx *bolt.Tx) []volumeMetadata ***REMOVED***
	var ls []volumeMetadata
	b := tx.Bucket(volumeBucketName)
	b.ForEach(func(k, v []byte) error ***REMOVED***
		if len(v) == 0 ***REMOVED***
			// don't try to unmarshal an empty value
			return nil
		***REMOVED***

		var m volumeMetadata
		if err := json.Unmarshal(v, &m); err != nil ***REMOVED***
			// Just log the error
			logrus.Errorf("Error while reading volume metadata for volume %q: %v", string(k), err)
			return nil
		***REMOVED***
		ls = append(ls, m)
		return nil
	***REMOVED***)
	return ls
***REMOVED***
