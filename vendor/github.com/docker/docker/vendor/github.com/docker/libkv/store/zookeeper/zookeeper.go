package zookeeper

import (
	"strings"
	"time"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	zk "github.com/samuel/go-zookeeper/zk"
)

const (
	// SOH control character
	SOH = "\x01"

	defaultTimeout = 10 * time.Second
)

// Zookeeper is the receiver type for
// the Store interface
type Zookeeper struct ***REMOVED***
	timeout time.Duration
	client  *zk.Conn
***REMOVED***

type zookeeperLock struct ***REMOVED***
	client *zk.Conn
	lock   *zk.Lock
	key    string
	value  []byte
***REMOVED***

// Register registers zookeeper to libkv
func Register() ***REMOVED***
	libkv.AddStore(store.ZK, New)
***REMOVED***

// New creates a new Zookeeper client given a
// list of endpoints and an optional tls config
func New(endpoints []string, options *store.Config) (store.Store, error) ***REMOVED***
	s := &Zookeeper***REMOVED******REMOVED***
	s.timeout = defaultTimeout

	// Set options
	if options != nil ***REMOVED***
		if options.ConnectionTimeout != 0 ***REMOVED***
			s.setTimeout(options.ConnectionTimeout)
		***REMOVED***
	***REMOVED***

	// Connect to Zookeeper
	conn, _, err := zk.Connect(endpoints, s.timeout)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	s.client = conn

	return s, nil
***REMOVED***

// setTimeout sets the timeout for connecting to Zookeeper
func (s *Zookeeper) setTimeout(time time.Duration) ***REMOVED***
	s.timeout = time
***REMOVED***

// Get the value at "key", returns the last modified index
// to use in conjunction to Atomic calls
func (s *Zookeeper) Get(key string) (pair *store.KVPair, err error) ***REMOVED***
	resp, meta, err := s.client.Get(s.normalize(key))

	if err != nil ***REMOVED***
		if err == zk.ErrNoNode ***REMOVED***
			return nil, store.ErrKeyNotFound
		***REMOVED***
		return nil, err
	***REMOVED***

	// FIXME handle very rare cases where Get returns the
	// SOH control character instead of the actual value
	if string(resp) == SOH ***REMOVED***
		return s.Get(store.Normalize(key))
	***REMOVED***

	pair = &store.KVPair***REMOVED***
		Key:       key,
		Value:     resp,
		LastIndex: uint64(meta.Version),
	***REMOVED***

	return pair, nil
***REMOVED***

// createFullPath creates the entire path for a directory
// that does not exist
func (s *Zookeeper) createFullPath(path []string, ephemeral bool) error ***REMOVED***
	for i := 1; i <= len(path); i++ ***REMOVED***
		newpath := "/" + strings.Join(path[:i], "/")
		if i == len(path) && ephemeral ***REMOVED***
			_, err := s.client.Create(newpath, []byte***REMOVED******REMOVED***, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
			return err
		***REMOVED***
		_, err := s.client.Create(newpath, []byte***REMOVED******REMOVED***, 0, zk.WorldACL(zk.PermAll))
		if err != nil ***REMOVED***
			// Skip if node already exists
			if err != zk.ErrNodeExists ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Put a value at "key"
func (s *Zookeeper) Put(key string, value []byte, opts *store.WriteOptions) error ***REMOVED***
	fkey := s.normalize(key)

	exists, err := s.Exists(key)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if !exists ***REMOVED***
		if opts != nil && opts.TTL > 0 ***REMOVED***
			s.createFullPath(store.SplitKey(strings.TrimSuffix(key, "/")), true)
		***REMOVED*** else ***REMOVED***
			s.createFullPath(store.SplitKey(strings.TrimSuffix(key, "/")), false)
		***REMOVED***
	***REMOVED***

	_, err = s.client.Set(fkey, value, -1)
	return err
***REMOVED***

// Delete a value at "key"
func (s *Zookeeper) Delete(key string) error ***REMOVED***
	err := s.client.Delete(s.normalize(key), -1)
	if err == zk.ErrNoNode ***REMOVED***
		return store.ErrKeyNotFound
	***REMOVED***
	return err
***REMOVED***

// Exists checks if the key exists inside the store
func (s *Zookeeper) Exists(key string) (bool, error) ***REMOVED***
	exists, _, err := s.client.Exists(s.normalize(key))
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	return exists, nil
***REMOVED***

// Watch for changes on a "key"
// It returns a channel that will receive changes or pass
// on errors. Upon creation, the current value will first
// be sent to the channel. Providing a non-nil stopCh can
// be used to stop watching.
func (s *Zookeeper) Watch(key string, stopCh <-chan struct***REMOVED******REMOVED***) (<-chan *store.KVPair, error) ***REMOVED***
	// Get the key first
	pair, err := s.Get(key)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Catch zk notifications and fire changes into the channel.
	watchCh := make(chan *store.KVPair)
	go func() ***REMOVED***
		defer close(watchCh)

		// Get returns the current value to the channel prior
		// to listening to any event that may occur on that key
		watchCh <- pair
		for ***REMOVED***
			_, _, eventCh, err := s.client.GetW(s.normalize(key))
			if err != nil ***REMOVED***
				return
			***REMOVED***
			select ***REMOVED***
			case e := <-eventCh:
				if e.Type == zk.EventNodeDataChanged ***REMOVED***
					if entry, err := s.Get(key); err == nil ***REMOVED***
						watchCh <- entry
					***REMOVED***
				***REMOVED***
			case <-stopCh:
				// There is no way to stop GetW so just quit
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return watchCh, nil
***REMOVED***

// WatchTree watches for changes on a "directory"
// It returns a channel that will receive changes or pass
// on errors. Upon creating a watch, the current childs values
// will be sent to the channel .Providing a non-nil stopCh can
// be used to stop watching.
func (s *Zookeeper) WatchTree(directory string, stopCh <-chan struct***REMOVED******REMOVED***) (<-chan []*store.KVPair, error) ***REMOVED***
	// List the childrens first
	entries, err := s.List(directory)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Catch zk notifications and fire changes into the channel.
	watchCh := make(chan []*store.KVPair)
	go func() ***REMOVED***
		defer close(watchCh)

		// List returns the children values to the channel
		// prior to listening to any events that may occur
		// on those keys
		watchCh <- entries

		for ***REMOVED***
			_, _, eventCh, err := s.client.ChildrenW(s.normalize(directory))
			if err != nil ***REMOVED***
				return
			***REMOVED***
			select ***REMOVED***
			case e := <-eventCh:
				if e.Type == zk.EventNodeChildrenChanged ***REMOVED***
					if kv, err := s.List(directory); err == nil ***REMOVED***
						watchCh <- kv
					***REMOVED***
				***REMOVED***
			case <-stopCh:
				// There is no way to stop GetW so just quit
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return watchCh, nil
***REMOVED***

// List child nodes of a given directory
func (s *Zookeeper) List(directory string) ([]*store.KVPair, error) ***REMOVED***
	keys, stat, err := s.client.Children(s.normalize(directory))
	if err != nil ***REMOVED***
		if err == zk.ErrNoNode ***REMOVED***
			return nil, store.ErrKeyNotFound
		***REMOVED***
		return nil, err
	***REMOVED***

	kv := []*store.KVPair***REMOVED******REMOVED***

	// FIXME Costly Get request for each child key..
	for _, key := range keys ***REMOVED***
		pair, err := s.Get(strings.TrimSuffix(directory, "/") + s.normalize(key))
		if err != nil ***REMOVED***
			// If node is not found: List is out of date, retry
			if err == store.ErrKeyNotFound ***REMOVED***
				return s.List(directory)
			***REMOVED***
			return nil, err
		***REMOVED***

		kv = append(kv, &store.KVPair***REMOVED***
			Key:       key,
			Value:     []byte(pair.Value),
			LastIndex: uint64(stat.Version),
		***REMOVED***)
	***REMOVED***

	return kv, nil
***REMOVED***

// DeleteTree deletes a range of keys under a given directory
func (s *Zookeeper) DeleteTree(directory string) error ***REMOVED***
	pairs, err := s.List(directory)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var reqs []interface***REMOVED******REMOVED***

	for _, pair := range pairs ***REMOVED***
		reqs = append(reqs, &zk.DeleteRequest***REMOVED***
			Path:    s.normalize(directory + "/" + pair.Key),
			Version: -1,
		***REMOVED***)
	***REMOVED***

	_, err = s.client.Multi(reqs...)
	return err
***REMOVED***

// AtomicPut put a value at "key" if the key has not been
// modified in the meantime, throws an error if this is the case
func (s *Zookeeper) AtomicPut(key string, value []byte, previous *store.KVPair, _ *store.WriteOptions) (bool, *store.KVPair, error) ***REMOVED***
	var lastIndex uint64

	if previous != nil ***REMOVED***
		meta, err := s.client.Set(s.normalize(key), value, int32(previous.LastIndex))
		if err != nil ***REMOVED***
			// Compare Failed
			if err == zk.ErrBadVersion ***REMOVED***
				return false, nil, store.ErrKeyModified
			***REMOVED***
			return false, nil, err
		***REMOVED***
		lastIndex = uint64(meta.Version)
	***REMOVED*** else ***REMOVED***
		// Interpret previous == nil as create operation.
		_, err := s.client.Create(s.normalize(key), value, 0, zk.WorldACL(zk.PermAll))
		if err != nil ***REMOVED***
			// Directory does not exist
			if err == zk.ErrNoNode ***REMOVED***

				// Create the directory
				parts := store.SplitKey(strings.TrimSuffix(key, "/"))
				parts = parts[:len(parts)-1]
				if err = s.createFullPath(parts, false); err != nil ***REMOVED***
					// Failed to create the directory.
					return false, nil, err
				***REMOVED***

				// Create the node
				if _, err := s.client.Create(s.normalize(key), value, 0, zk.WorldACL(zk.PermAll)); err != nil ***REMOVED***
					// Node exist error (when previous nil)
					if err == zk.ErrNodeExists ***REMOVED***
						return false, nil, store.ErrKeyExists
					***REMOVED***
					return false, nil, err
				***REMOVED***

			***REMOVED*** else ***REMOVED***
				// Node Exists error (when previous nil)
				if err == zk.ErrNodeExists ***REMOVED***
					return false, nil, store.ErrKeyExists
				***REMOVED***

				// Unhandled error
				return false, nil, err
			***REMOVED***
		***REMOVED***
		lastIndex = 0 // Newly created nodes have version 0.
	***REMOVED***

	pair := &store.KVPair***REMOVED***
		Key:       key,
		Value:     value,
		LastIndex: lastIndex,
	***REMOVED***

	return true, pair, nil
***REMOVED***

// AtomicDelete deletes a value at "key" if the key
// has not been modified in the meantime, throws an
// error if this is the case
func (s *Zookeeper) AtomicDelete(key string, previous *store.KVPair) (bool, error) ***REMOVED***
	if previous == nil ***REMOVED***
		return false, store.ErrPreviousNotSpecified
	***REMOVED***

	err := s.client.Delete(s.normalize(key), int32(previous.LastIndex))
	if err != nil ***REMOVED***
		// Key not found
		if err == zk.ErrNoNode ***REMOVED***
			return false, store.ErrKeyNotFound
		***REMOVED***
		// Compare failed
		if err == zk.ErrBadVersion ***REMOVED***
			return false, store.ErrKeyModified
		***REMOVED***
		// General store error
		return false, err
	***REMOVED***
	return true, nil
***REMOVED***

// NewLock returns a handle to a lock struct which can
// be used to provide mutual exclusion on a key
func (s *Zookeeper) NewLock(key string, options *store.LockOptions) (lock store.Locker, err error) ***REMOVED***
	value := []byte("")

	// Apply options
	if options != nil ***REMOVED***
		if options.Value != nil ***REMOVED***
			value = options.Value
		***REMOVED***
	***REMOVED***

	lock = &zookeeperLock***REMOVED***
		client: s.client,
		key:    s.normalize(key),
		value:  value,
		lock:   zk.NewLock(s.client, s.normalize(key), zk.WorldACL(zk.PermAll)),
	***REMOVED***

	return lock, err
***REMOVED***

// Lock attempts to acquire the lock and blocks while
// doing so. It returns a channel that is closed if our
// lock is lost or if an error occurs
func (l *zookeeperLock) Lock(stopChan chan struct***REMOVED******REMOVED***) (<-chan struct***REMOVED******REMOVED***, error) ***REMOVED***
	err := l.lock.Lock()

	if err == nil ***REMOVED***
		// We hold the lock, we can set our value
		// FIXME: The value is left behind
		// (problematic for leader election)
		_, err = l.client.Set(l.key, l.value, -1)
	***REMOVED***

	return make(chan struct***REMOVED******REMOVED***), err
***REMOVED***

// Unlock the "key". Calling unlock while
// not holding the lock will throw an error
func (l *zookeeperLock) Unlock() error ***REMOVED***
	return l.lock.Unlock()
***REMOVED***

// Close closes the client connection
func (s *Zookeeper) Close() ***REMOVED***
	s.client.Close()
***REMOVED***

// Normalize the key for usage in Zookeeper
func (s *Zookeeper) normalize(key string) string ***REMOVED***
	key = store.Normalize(key)
	return strings.TrimSuffix(key, "/")
***REMOVED***
