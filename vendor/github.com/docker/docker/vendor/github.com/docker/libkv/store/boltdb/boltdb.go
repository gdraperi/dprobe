package boltdb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/boltdb/bolt"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
)

var (
	// ErrMultipleEndpointsUnsupported is thrown when multiple endpoints specified for
	// BoltDB. Endpoint has to be a local file path
	ErrMultipleEndpointsUnsupported = errors.New("boltdb supports one endpoint and should be a file path")
	// ErrBoltBucketOptionMissing is thrown when boltBcuket config option is missing
	ErrBoltBucketOptionMissing = errors.New("boltBucket config option missing")
)

const (
	filePerm os.FileMode = 0644
)

//BoltDB type implements the Store interface
type BoltDB struct ***REMOVED***
	client     *bolt.DB
	boltBucket []byte
	dbIndex    uint64
	path       string
	timeout    time.Duration
	// By default libkv opens and closes the bolt DB connection  for every
	// get/put operation. This allows multiple apps to use a Bolt DB at the
	// same time.
	// PersistConnection flag provides an option to override ths behavior.
	// ie: open the connection in New and use it till Close is called.
	PersistConnection bool
	sync.Mutex
***REMOVED***

const (
	libkvmetadatalen = 8
	transientTimeout = time.Duration(10) * time.Second
)

// Register registers boltdb to libkv
func Register() ***REMOVED***
	libkv.AddStore(store.BOLTDB, New)
***REMOVED***

// New opens a new BoltDB connection to the specified path and bucket
func New(endpoints []string, options *store.Config) (store.Store, error) ***REMOVED***
	var (
		db          *bolt.DB
		err         error
		boltOptions *bolt.Options
		timeout     = transientTimeout
	)

	if len(endpoints) > 1 ***REMOVED***
		return nil, ErrMultipleEndpointsUnsupported
	***REMOVED***

	if (options == nil) || (len(options.Bucket) == 0) ***REMOVED***
		return nil, ErrBoltBucketOptionMissing
	***REMOVED***

	dir, _ := filepath.Split(endpoints[0])
	if err = os.MkdirAll(dir, 0750); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if options.PersistConnection ***REMOVED***
		boltOptions = &bolt.Options***REMOVED***Timeout: options.ConnectionTimeout***REMOVED***
		db, err = bolt.Open(endpoints[0], filePerm, boltOptions)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	if options.ConnectionTimeout != 0 ***REMOVED***
		timeout = options.ConnectionTimeout
	***REMOVED***

	b := &BoltDB***REMOVED***
		client:            db,
		path:              endpoints[0],
		boltBucket:        []byte(options.Bucket),
		timeout:           timeout,
		PersistConnection: options.PersistConnection,
	***REMOVED***

	return b, nil
***REMOVED***

func (b *BoltDB) reset() ***REMOVED***
	b.path = ""
	b.boltBucket = []byte***REMOVED******REMOVED***
***REMOVED***

func (b *BoltDB) getDBhandle() (*bolt.DB, error) ***REMOVED***
	var (
		db  *bolt.DB
		err error
	)
	if !b.PersistConnection ***REMOVED***
		boltOptions := &bolt.Options***REMOVED***Timeout: b.timeout***REMOVED***
		if db, err = bolt.Open(b.path, filePerm, boltOptions); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		b.client = db
	***REMOVED***

	return b.client, nil
***REMOVED***

func (b *BoltDB) releaseDBhandle() ***REMOVED***
	if !b.PersistConnection ***REMOVED***
		b.client.Close()
	***REMOVED***
***REMOVED***

// Get the value at "key". BoltDB doesn't provide an inbuilt last modified index with every kv pair. Its implemented by
// by a atomic counter maintained by the libkv and appened to the value passed by the client.
func (b *BoltDB) Get(key string) (*store.KVPair, error) ***REMOVED***
	var (
		val []byte
		db  *bolt.DB
		err error
	)
	b.Lock()
	defer b.Unlock()

	if db, err = b.getDBhandle(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer b.releaseDBhandle()

	err = db.View(func(tx *bolt.Tx) error ***REMOVED***
		bucket := tx.Bucket(b.boltBucket)
		if bucket == nil ***REMOVED***
			return store.ErrKeyNotFound
		***REMOVED***

		v := bucket.Get([]byte(key))
		val = make([]byte, len(v))
		copy(val, v)

		return nil
	***REMOVED***)

	if len(val) == 0 ***REMOVED***
		return nil, store.ErrKeyNotFound
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	dbIndex := binary.LittleEndian.Uint64(val[:libkvmetadatalen])
	val = val[libkvmetadatalen:]

	return &store.KVPair***REMOVED***Key: key, Value: val, LastIndex: (dbIndex)***REMOVED***, nil
***REMOVED***

//Put the key, value pair. index number metadata is prepended to the value
func (b *BoltDB) Put(key string, value []byte, opts *store.WriteOptions) error ***REMOVED***
	var (
		dbIndex uint64
		db      *bolt.DB
		err     error
	)
	b.Lock()
	defer b.Unlock()

	dbval := make([]byte, libkvmetadatalen)

	if db, err = b.getDBhandle(); err != nil ***REMOVED***
		return err
	***REMOVED***
	defer b.releaseDBhandle()

	err = db.Update(func(tx *bolt.Tx) error ***REMOVED***
		bucket, err := tx.CreateBucketIfNotExists(b.boltBucket)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		dbIndex = atomic.AddUint64(&b.dbIndex, 1)
		binary.LittleEndian.PutUint64(dbval, dbIndex)
		dbval = append(dbval, value...)

		err = bucket.Put([]byte(key), dbval)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		return nil
	***REMOVED***)
	return err
***REMOVED***

//Delete the value for the given key.
func (b *BoltDB) Delete(key string) error ***REMOVED***
	var (
		db  *bolt.DB
		err error
	)
	b.Lock()
	defer b.Unlock()

	if db, err = b.getDBhandle(); err != nil ***REMOVED***
		return err
	***REMOVED***
	defer b.releaseDBhandle()

	err = db.Update(func(tx *bolt.Tx) error ***REMOVED***
		bucket := tx.Bucket(b.boltBucket)
		if bucket == nil ***REMOVED***
			return store.ErrKeyNotFound
		***REMOVED***
		err := bucket.Delete([]byte(key))
		return err
	***REMOVED***)
	return err
***REMOVED***

// Exists checks if the key exists inside the store
func (b *BoltDB) Exists(key string) (bool, error) ***REMOVED***
	var (
		val []byte
		db  *bolt.DB
		err error
	)
	b.Lock()
	defer b.Unlock()

	if db, err = b.getDBhandle(); err != nil ***REMOVED***
		return false, err
	***REMOVED***
	defer b.releaseDBhandle()

	err = db.View(func(tx *bolt.Tx) error ***REMOVED***
		bucket := tx.Bucket(b.boltBucket)
		if bucket == nil ***REMOVED***
			return store.ErrKeyNotFound
		***REMOVED***

		val = bucket.Get([]byte(key))

		return nil
	***REMOVED***)

	if len(val) == 0 ***REMOVED***
		return false, err
	***REMOVED***
	return true, err
***REMOVED***

// List returns the range of keys starting with the passed in prefix
func (b *BoltDB) List(keyPrefix string) ([]*store.KVPair, error) ***REMOVED***
	var (
		db  *bolt.DB
		err error
	)
	b.Lock()
	defer b.Unlock()

	kv := []*store.KVPair***REMOVED******REMOVED***

	if db, err = b.getDBhandle(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer b.releaseDBhandle()

	err = db.View(func(tx *bolt.Tx) error ***REMOVED***
		bucket := tx.Bucket(b.boltBucket)
		if bucket == nil ***REMOVED***
			return store.ErrKeyNotFound
		***REMOVED***

		cursor := bucket.Cursor()
		prefix := []byte(keyPrefix)

		for key, v := cursor.Seek(prefix); bytes.HasPrefix(key, prefix); key, v = cursor.Next() ***REMOVED***

			dbIndex := binary.LittleEndian.Uint64(v[:libkvmetadatalen])
			v = v[libkvmetadatalen:]
			val := make([]byte, len(v))
			copy(val, v)

			kv = append(kv, &store.KVPair***REMOVED***
				Key:       string(key),
				Value:     val,
				LastIndex: dbIndex,
			***REMOVED***)
		***REMOVED***
		return nil
	***REMOVED***)
	if len(kv) == 0 ***REMOVED***
		return nil, store.ErrKeyNotFound
	***REMOVED***
	return kv, err
***REMOVED***

// AtomicDelete deletes a value at "key" if the key
// has not been modified in the meantime, throws an
// error if this is the case
func (b *BoltDB) AtomicDelete(key string, previous *store.KVPair) (bool, error) ***REMOVED***
	var (
		val []byte
		db  *bolt.DB
		err error
	)
	b.Lock()
	defer b.Unlock()

	if previous == nil ***REMOVED***
		return false, store.ErrPreviousNotSpecified
	***REMOVED***
	if db, err = b.getDBhandle(); err != nil ***REMOVED***
		return false, err
	***REMOVED***
	defer b.releaseDBhandle()

	err = db.Update(func(tx *bolt.Tx) error ***REMOVED***
		bucket := tx.Bucket(b.boltBucket)
		if bucket == nil ***REMOVED***
			return store.ErrKeyNotFound
		***REMOVED***

		val = bucket.Get([]byte(key))
		if val == nil ***REMOVED***
			return store.ErrKeyNotFound
		***REMOVED***
		dbIndex := binary.LittleEndian.Uint64(val[:libkvmetadatalen])
		if dbIndex != previous.LastIndex ***REMOVED***
			return store.ErrKeyModified
		***REMOVED***
		err := bucket.Delete([]byte(key))
		return err
	***REMOVED***)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	return true, err
***REMOVED***

// AtomicPut puts a value at "key" if the key has not been
// modified since the last Put, throws an error if this is the case
func (b *BoltDB) AtomicPut(key string, value []byte, previous *store.KVPair, options *store.WriteOptions) (bool, *store.KVPair, error) ***REMOVED***
	var (
		val     []byte
		dbIndex uint64
		db      *bolt.DB
		err     error
	)
	b.Lock()
	defer b.Unlock()

	dbval := make([]byte, libkvmetadatalen)

	if db, err = b.getDBhandle(); err != nil ***REMOVED***
		return false, nil, err
	***REMOVED***
	defer b.releaseDBhandle()

	err = db.Update(func(tx *bolt.Tx) error ***REMOVED***
		var err error
		bucket := tx.Bucket(b.boltBucket)
		if bucket == nil ***REMOVED***
			if previous != nil ***REMOVED***
				return store.ErrKeyNotFound
			***REMOVED***
			bucket, err = tx.CreateBucket(b.boltBucket)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		// AtomicPut is equivalent to Put if previous is nil and the Ky
		// doesn't exist in the DB.
		val = bucket.Get([]byte(key))
		if previous == nil && len(val) != 0 ***REMOVED***
			return store.ErrKeyExists
		***REMOVED***
		if previous != nil ***REMOVED***
			if len(val) == 0 ***REMOVED***
				return store.ErrKeyNotFound
			***REMOVED***
			dbIndex = binary.LittleEndian.Uint64(val[:libkvmetadatalen])
			if dbIndex != previous.LastIndex ***REMOVED***
				return store.ErrKeyModified
			***REMOVED***
		***REMOVED***
		dbIndex = atomic.AddUint64(&b.dbIndex, 1)
		binary.LittleEndian.PutUint64(dbval, b.dbIndex)
		dbval = append(dbval, value...)
		return (bucket.Put([]byte(key), dbval))
	***REMOVED***)
	if err != nil ***REMOVED***
		return false, nil, err
	***REMOVED***

	updated := &store.KVPair***REMOVED***
		Key:       key,
		Value:     value,
		LastIndex: dbIndex,
	***REMOVED***

	return true, updated, nil
***REMOVED***

// Close the db connection to the BoltDB
func (b *BoltDB) Close() ***REMOVED***
	b.Lock()
	defer b.Unlock()

	if !b.PersistConnection ***REMOVED***
		b.reset()
	***REMOVED*** else ***REMOVED***
		b.client.Close()
	***REMOVED***
	return
***REMOVED***

// DeleteTree deletes a range of keys with a given prefix
func (b *BoltDB) DeleteTree(keyPrefix string) error ***REMOVED***
	var (
		db  *bolt.DB
		err error
	)
	b.Lock()
	defer b.Unlock()

	if db, err = b.getDBhandle(); err != nil ***REMOVED***
		return err
	***REMOVED***
	defer b.releaseDBhandle()

	err = db.Update(func(tx *bolt.Tx) error ***REMOVED***
		bucket := tx.Bucket(b.boltBucket)
		if bucket == nil ***REMOVED***
			return store.ErrKeyNotFound
		***REMOVED***

		cursor := bucket.Cursor()
		prefix := []byte(keyPrefix)

		for key, _ := cursor.Seek(prefix); bytes.HasPrefix(key, prefix); key, _ = cursor.Next() ***REMOVED***
			_ = bucket.Delete([]byte(key))
		***REMOVED***
		return nil
	***REMOVED***)

	return err
***REMOVED***

// NewLock has to implemented at the library level since its not supported by BoltDB
func (b *BoltDB) NewLock(key string, options *store.LockOptions) (store.Locker, error) ***REMOVED***
	return nil, store.ErrCallNotSupported
***REMOVED***

// Watch has to implemented at the library level since its not supported by BoltDB
func (b *BoltDB) Watch(key string, stopCh <-chan struct***REMOVED******REMOVED***) (<-chan *store.KVPair, error) ***REMOVED***
	return nil, store.ErrCallNotSupported
***REMOVED***

// WatchTree has to implemented at the library level since its not supported by BoltDB
func (b *BoltDB) WatchTree(directory string, stopCh <-chan struct***REMOVED******REMOVED***) (<-chan []*store.KVPair, error) ***REMOVED***
	return nil, store.ErrCallNotSupported
***REMOVED***
