package metadata

import "github.com/boltdb/bolt"

type migration struct ***REMOVED***
	schema  string
	version int
	migrate func(*bolt.Tx) error
***REMOVED***

// migrations stores the list of database migrations
// for each update to the database schema. The migrations
// array MUST be ordered by version from least to greatest.
// The last entry in the array should correspond to the
// schemaVersion and dbVersion constants.
// A migration test MUST be added for each migration in
// the array.
// The migrate function can safely assume the version
// of the data it is migrating from is the previous version
// of the database.
var migrations = []migration***REMOVED***
	***REMOVED***
		schema:  "v1",
		version: 1,
		migrate: addChildLinks,
	***REMOVED***,
***REMOVED***

// addChildLinks Adds children key to the snapshotters to enforce snapshot
// entries cannot be removed which have children
func addChildLinks(tx *bolt.Tx) error ***REMOVED***
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
		nbkt := v1bkt.Bucket(k)

		sbkt := nbkt.Bucket(bucketKeyObjectSnapshots)
		if sbkt != nil ***REMOVED***
			// Iterate through each snapshotter
			if err := sbkt.ForEach(func(sk, sv []byte) error ***REMOVED***
				if sv != nil ***REMOVED***
					return nil
				***REMOVED***
				snbkt := sbkt.Bucket(sk)

				// Iterate through each snapshot
				return snbkt.ForEach(func(k, v []byte) error ***REMOVED***
					if v != nil ***REMOVED***
						return nil
					***REMOVED***
					parent := snbkt.Bucket(k).Get(bucketKeyParent)
					if len(parent) > 0 ***REMOVED***
						pbkt := snbkt.Bucket(parent)
						if pbkt == nil ***REMOVED***
							// Not enforcing consistency during migration, skip
							return nil
						***REMOVED***
						cbkt, err := pbkt.CreateBucketIfNotExists(bucketKeyChildren)
						if err != nil ***REMOVED***
							return err
						***REMOVED***
						if err := cbkt.Put(k, nil); err != nil ***REMOVED***
							return err
						***REMOVED***
					***REMOVED***

					return nil
				***REMOVED***)
			***REMOVED***); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
