package etcd

import (
	"crypto/tls"
	"errors"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/context"

	etcd "github.com/coreos/etcd/client"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
)

var (
	// ErrAbortTryLock is thrown when a user stops trying to seek the lock
	// by sending a signal to the stop chan, this is used to verify if the
	// operation succeeded
	ErrAbortTryLock = errors.New("lock operation aborted")
)

// Etcd is the receiver type for the
// Store interface
type Etcd struct ***REMOVED***
	client etcd.KeysAPI
***REMOVED***

type etcdLock struct ***REMOVED***
	client    etcd.KeysAPI
	stopLock  chan struct***REMOVED******REMOVED***
	stopRenew chan struct***REMOVED******REMOVED***
	key       string
	value     string
	last      *etcd.Response
	ttl       time.Duration
***REMOVED***

const (
	periodicSync      = 5 * time.Minute
	defaultLockTTL    = 20 * time.Second
	defaultUpdateTime = 5 * time.Second
)

// Register registers etcd to libkv
func Register() ***REMOVED***
	libkv.AddStore(store.ETCD, New)
***REMOVED***

// New creates a new Etcd client given a list
// of endpoints and an optional tls config
func New(addrs []string, options *store.Config) (store.Store, error) ***REMOVED***
	s := &Etcd***REMOVED******REMOVED***

	var (
		entries []string
		err     error
	)

	entries = store.CreateEndpoints(addrs, "http")
	cfg := &etcd.Config***REMOVED***
		Endpoints:               entries,
		Transport:               etcd.DefaultTransport,
		HeaderTimeoutPerRequest: 3 * time.Second,
	***REMOVED***

	// Set options
	if options != nil ***REMOVED***
		if options.TLS != nil ***REMOVED***
			setTLS(cfg, options.TLS, addrs)
		***REMOVED***
		if options.ConnectionTimeout != 0 ***REMOVED***
			setTimeout(cfg, options.ConnectionTimeout)
		***REMOVED***
		if options.Username != "" ***REMOVED***
			setCredentials(cfg, options.Username, options.Password)
		***REMOVED***
	***REMOVED***

	c, err := etcd.New(*cfg)
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***

	s.client = etcd.NewKeysAPI(c)

	// Periodic Cluster Sync
	go func() ***REMOVED***
		for ***REMOVED***
			if err := c.AutoSync(context.Background(), periodicSync); err != nil ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return s, nil
***REMOVED***

// SetTLS sets the tls configuration given a tls.Config scheme
func setTLS(cfg *etcd.Config, tls *tls.Config, addrs []string) ***REMOVED***
	entries := store.CreateEndpoints(addrs, "https")
	cfg.Endpoints = entries

	// Set transport
	t := http.Transport***REMOVED***
		Dial: (&net.Dialer***REMOVED***
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		***REMOVED***).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     tls,
	***REMOVED***

	cfg.Transport = &t
***REMOVED***

// setTimeout sets the timeout used for connecting to the store
func setTimeout(cfg *etcd.Config, time time.Duration) ***REMOVED***
	cfg.HeaderTimeoutPerRequest = time
***REMOVED***

// setCredentials sets the username/password credentials for connecting to Etcd
func setCredentials(cfg *etcd.Config, username, password string) ***REMOVED***
	cfg.Username = username
	cfg.Password = password
***REMOVED***

// Normalize the key for usage in Etcd
func (s *Etcd) normalize(key string) string ***REMOVED***
	key = store.Normalize(key)
	return strings.TrimPrefix(key, "/")
***REMOVED***

// keyNotFound checks on the error returned by the KeysAPI
// to verify if the key exists in the store or not
func keyNotFound(err error) bool ***REMOVED***
	if err != nil ***REMOVED***
		if etcdError, ok := err.(etcd.Error); ok ***REMOVED***
			if etcdError.Code == etcd.ErrorCodeKeyNotFound ||
				etcdError.Code == etcd.ErrorCodeNotFile ||
				etcdError.Code == etcd.ErrorCodeNotDir ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// Get the value at "key", returns the last modified
// index to use in conjunction to Atomic calls
func (s *Etcd) Get(key string) (pair *store.KVPair, err error) ***REMOVED***
	getOpts := &etcd.GetOptions***REMOVED***
		Quorum: true,
	***REMOVED***

	result, err := s.client.Get(context.Background(), s.normalize(key), getOpts)
	if err != nil ***REMOVED***
		if keyNotFound(err) ***REMOVED***
			return nil, store.ErrKeyNotFound
		***REMOVED***
		return nil, err
	***REMOVED***

	pair = &store.KVPair***REMOVED***
		Key:       key,
		Value:     []byte(result.Node.Value),
		LastIndex: result.Node.ModifiedIndex,
	***REMOVED***

	return pair, nil
***REMOVED***

// Put a value at "key"
func (s *Etcd) Put(key string, value []byte, opts *store.WriteOptions) error ***REMOVED***
	setOpts := &etcd.SetOptions***REMOVED******REMOVED***

	// Set options
	if opts != nil ***REMOVED***
		setOpts.Dir = opts.IsDir
		setOpts.TTL = opts.TTL
	***REMOVED***

	_, err := s.client.Set(context.Background(), s.normalize(key), string(value), setOpts)
	return err
***REMOVED***

// Delete a value at "key"
func (s *Etcd) Delete(key string) error ***REMOVED***
	opts := &etcd.DeleteOptions***REMOVED***
		Recursive: false,
	***REMOVED***

	_, err := s.client.Delete(context.Background(), s.normalize(key), opts)
	if keyNotFound(err) ***REMOVED***
		return store.ErrKeyNotFound
	***REMOVED***
	return err
***REMOVED***

// Exists checks if the key exists inside the store
func (s *Etcd) Exists(key string) (bool, error) ***REMOVED***
	_, err := s.Get(key)
	if err != nil ***REMOVED***
		if err == store.ErrKeyNotFound ***REMOVED***
			return false, nil
		***REMOVED***
		return false, err
	***REMOVED***
	return true, nil
***REMOVED***

// Watch for changes on a "key"
// It returns a channel that will receive changes or pass
// on errors. Upon creation, the current value will first
// be sent to the channel. Providing a non-nil stopCh can
// be used to stop watching.
func (s *Etcd) Watch(key string, stopCh <-chan struct***REMOVED******REMOVED***) (<-chan *store.KVPair, error) ***REMOVED***
	opts := &etcd.WatcherOptions***REMOVED***Recursive: false***REMOVED***
	watcher := s.client.Watcher(s.normalize(key), opts)

	// watchCh is sending back events to the caller
	watchCh := make(chan *store.KVPair)

	go func() ***REMOVED***
		defer close(watchCh)

		// Get the current value
		pair, err := s.Get(key)
		if err != nil ***REMOVED***
			return
		***REMOVED***

		// Push the current value through the channel.
		watchCh <- pair

		for ***REMOVED***
			// Check if the watch was stopped by the caller
			select ***REMOVED***
			case <-stopCh:
				return
			default:
			***REMOVED***

			result, err := watcher.Next(context.Background())

			if err != nil ***REMOVED***
				return
			***REMOVED***

			watchCh <- &store.KVPair***REMOVED***
				Key:       key,
				Value:     []byte(result.Node.Value),
				LastIndex: result.Node.ModifiedIndex,
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return watchCh, nil
***REMOVED***

// WatchTree watches for changes on a "directory"
// It returns a channel that will receive changes or pass
// on errors. Upon creating a watch, the current childs values
// will be sent to the channel. Providing a non-nil stopCh can
// be used to stop watching.
func (s *Etcd) WatchTree(directory string, stopCh <-chan struct***REMOVED******REMOVED***) (<-chan []*store.KVPair, error) ***REMOVED***
	watchOpts := &etcd.WatcherOptions***REMOVED***Recursive: true***REMOVED***
	watcher := s.client.Watcher(s.normalize(directory), watchOpts)

	// watchCh is sending back events to the caller
	watchCh := make(chan []*store.KVPair)

	go func() ***REMOVED***
		defer close(watchCh)

		// Get child values
		list, err := s.List(directory)
		if err != nil ***REMOVED***
			return
		***REMOVED***

		// Push the current value through the channel.
		watchCh <- list

		for ***REMOVED***
			// Check if the watch was stopped by the caller
			select ***REMOVED***
			case <-stopCh:
				return
			default:
			***REMOVED***

			_, err := watcher.Next(context.Background())

			if err != nil ***REMOVED***
				return
			***REMOVED***

			list, err = s.List(directory)
			if err != nil ***REMOVED***
				return
			***REMOVED***

			watchCh <- list
		***REMOVED***
	***REMOVED***()

	return watchCh, nil
***REMOVED***

// AtomicPut puts a value at "key" if the key has not been
// modified in the meantime, throws an error if this is the case
func (s *Etcd) AtomicPut(key string, value []byte, previous *store.KVPair, opts *store.WriteOptions) (bool, *store.KVPair, error) ***REMOVED***
	var (
		meta *etcd.Response
		err  error
	)

	setOpts := &etcd.SetOptions***REMOVED******REMOVED***

	if previous != nil ***REMOVED***
		setOpts.PrevExist = etcd.PrevExist
		setOpts.PrevIndex = previous.LastIndex
		if previous.Value != nil ***REMOVED***
			setOpts.PrevValue = string(previous.Value)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		setOpts.PrevExist = etcd.PrevNoExist
	***REMOVED***

	if opts != nil ***REMOVED***
		if opts.TTL > 0 ***REMOVED***
			setOpts.TTL = opts.TTL
		***REMOVED***
	***REMOVED***

	meta, err = s.client.Set(context.Background(), s.normalize(key), string(value), setOpts)
	if err != nil ***REMOVED***
		if etcdError, ok := err.(etcd.Error); ok ***REMOVED***
			// Compare failed
			if etcdError.Code == etcd.ErrorCodeTestFailed ***REMOVED***
				return false, nil, store.ErrKeyModified
			***REMOVED***
			// Node exists error (when PrevNoExist)
			if etcdError.Code == etcd.ErrorCodeNodeExist ***REMOVED***
				return false, nil, store.ErrKeyExists
			***REMOVED***
		***REMOVED***
		return false, nil, err
	***REMOVED***

	updated := &store.KVPair***REMOVED***
		Key:       key,
		Value:     value,
		LastIndex: meta.Node.ModifiedIndex,
	***REMOVED***

	return true, updated, nil
***REMOVED***

// AtomicDelete deletes a value at "key" if the key
// has not been modified in the meantime, throws an
// error if this is the case
func (s *Etcd) AtomicDelete(key string, previous *store.KVPair) (bool, error) ***REMOVED***
	if previous == nil ***REMOVED***
		return false, store.ErrPreviousNotSpecified
	***REMOVED***

	delOpts := &etcd.DeleteOptions***REMOVED******REMOVED***

	if previous != nil ***REMOVED***
		delOpts.PrevIndex = previous.LastIndex
		if previous.Value != nil ***REMOVED***
			delOpts.PrevValue = string(previous.Value)
		***REMOVED***
	***REMOVED***

	_, err := s.client.Delete(context.Background(), s.normalize(key), delOpts)
	if err != nil ***REMOVED***
		if etcdError, ok := err.(etcd.Error); ok ***REMOVED***
			// Key Not Found
			if etcdError.Code == etcd.ErrorCodeKeyNotFound ***REMOVED***
				return false, store.ErrKeyNotFound
			***REMOVED***
			// Compare failed
			if etcdError.Code == etcd.ErrorCodeTestFailed ***REMOVED***
				return false, store.ErrKeyModified
			***REMOVED***
		***REMOVED***
		return false, err
	***REMOVED***

	return true, nil
***REMOVED***

// List child nodes of a given directory
func (s *Etcd) List(directory string) ([]*store.KVPair, error) ***REMOVED***
	getOpts := &etcd.GetOptions***REMOVED***
		Quorum:    true,
		Recursive: true,
		Sort:      true,
	***REMOVED***

	resp, err := s.client.Get(context.Background(), s.normalize(directory), getOpts)
	if err != nil ***REMOVED***
		if keyNotFound(err) ***REMOVED***
			return nil, store.ErrKeyNotFound
		***REMOVED***
		return nil, err
	***REMOVED***

	kv := []*store.KVPair***REMOVED******REMOVED***
	for _, n := range resp.Node.Nodes ***REMOVED***
		kv = append(kv, &store.KVPair***REMOVED***
			Key:       n.Key,
			Value:     []byte(n.Value),
			LastIndex: n.ModifiedIndex,
		***REMOVED***)
	***REMOVED***
	return kv, nil
***REMOVED***

// DeleteTree deletes a range of keys under a given directory
func (s *Etcd) DeleteTree(directory string) error ***REMOVED***
	delOpts := &etcd.DeleteOptions***REMOVED***
		Recursive: true,
	***REMOVED***

	_, err := s.client.Delete(context.Background(), s.normalize(directory), delOpts)
	if keyNotFound(err) ***REMOVED***
		return store.ErrKeyNotFound
	***REMOVED***
	return err
***REMOVED***

// NewLock returns a handle to a lock struct which can
// be used to provide mutual exclusion on a key
func (s *Etcd) NewLock(key string, options *store.LockOptions) (lock store.Locker, err error) ***REMOVED***
	var value string
	ttl := defaultLockTTL
	renewCh := make(chan struct***REMOVED******REMOVED***)

	// Apply options on Lock
	if options != nil ***REMOVED***
		if options.Value != nil ***REMOVED***
			value = string(options.Value)
		***REMOVED***
		if options.TTL != 0 ***REMOVED***
			ttl = options.TTL
		***REMOVED***
		if options.RenewLock != nil ***REMOVED***
			renewCh = options.RenewLock
		***REMOVED***
	***REMOVED***

	// Create lock object
	lock = &etcdLock***REMOVED***
		client:    s.client,
		stopRenew: renewCh,
		key:       s.normalize(key),
		value:     value,
		ttl:       ttl,
	***REMOVED***

	return lock, nil
***REMOVED***

// Lock attempts to acquire the lock and blocks while
// doing so. It returns a channel that is closed if our
// lock is lost or if an error occurs
func (l *etcdLock) Lock(stopChan chan struct***REMOVED******REMOVED***) (<-chan struct***REMOVED******REMOVED***, error) ***REMOVED***

	// Lock holder channel
	lockHeld := make(chan struct***REMOVED******REMOVED***)
	stopLocking := l.stopRenew

	setOpts := &etcd.SetOptions***REMOVED***
		TTL: l.ttl,
	***REMOVED***

	for ***REMOVED***
		setOpts.PrevExist = etcd.PrevNoExist
		resp, err := l.client.Set(context.Background(), l.key, l.value, setOpts)
		if err != nil ***REMOVED***
			if etcdError, ok := err.(etcd.Error); ok ***REMOVED***
				if etcdError.Code != etcd.ErrorCodeNodeExist ***REMOVED***
					return nil, err
				***REMOVED***
				setOpts.PrevIndex = ^uint64(0)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			setOpts.PrevIndex = resp.Node.ModifiedIndex
		***REMOVED***

		setOpts.PrevExist = etcd.PrevExist
		l.last, err = l.client.Set(context.Background(), l.key, l.value, setOpts)

		if err == nil ***REMOVED***
			// Leader section
			l.stopLock = stopLocking
			go l.holdLock(l.key, lockHeld, stopLocking)
			break
		***REMOVED*** else ***REMOVED***
			// If this is a legitimate error, return
			if etcdError, ok := err.(etcd.Error); ok ***REMOVED***
				if etcdError.Code != etcd.ErrorCodeTestFailed ***REMOVED***
					return nil, err
				***REMOVED***
			***REMOVED***

			// Seeker section
			errorCh := make(chan error)
			chWStop := make(chan bool)
			free := make(chan bool)

			go l.waitLock(l.key, errorCh, chWStop, free)

			// Wait for the key to be available or for
			// a signal to stop trying to lock the key
			select ***REMOVED***
			case <-free:
				break
			case err := <-errorCh:
				return nil, err
			case <-stopChan:
				return nil, ErrAbortTryLock
			***REMOVED***

			// Delete or Expire event occurred
			// Retry
		***REMOVED***
	***REMOVED***

	return lockHeld, nil
***REMOVED***

// Hold the lock as long as we can
// Updates the key ttl periodically until we receive
// an explicit stop signal from the Unlock method
func (l *etcdLock) holdLock(key string, lockHeld chan struct***REMOVED******REMOVED***, stopLocking <-chan struct***REMOVED******REMOVED***) ***REMOVED***
	defer close(lockHeld)

	update := time.NewTicker(l.ttl / 3)
	defer update.Stop()

	var err error
	setOpts := &etcd.SetOptions***REMOVED***TTL: l.ttl***REMOVED***

	for ***REMOVED***
		select ***REMOVED***
		case <-update.C:
			setOpts.PrevIndex = l.last.Node.ModifiedIndex
			l.last, err = l.client.Set(context.Background(), key, l.value, setOpts)
			if err != nil ***REMOVED***
				return
			***REMOVED***

		case <-stopLocking:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// WaitLock simply waits for the key to be available for creation
func (l *etcdLock) waitLock(key string, errorCh chan error, stopWatchCh chan bool, free chan<- bool) ***REMOVED***
	opts := &etcd.WatcherOptions***REMOVED***Recursive: false***REMOVED***
	watcher := l.client.Watcher(key, opts)

	for ***REMOVED***
		event, err := watcher.Next(context.Background())
		if err != nil ***REMOVED***
			errorCh <- err
			return
		***REMOVED***
		if event.Action == "delete" || event.Action == "expire" ***REMOVED***
			free <- true
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// Unlock the "key". Calling unlock while
// not holding the lock will throw an error
func (l *etcdLock) Unlock() error ***REMOVED***
	if l.stopLock != nil ***REMOVED***
		l.stopLock <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
	if l.last != nil ***REMOVED***
		delOpts := &etcd.DeleteOptions***REMOVED***
			PrevIndex: l.last.Node.ModifiedIndex,
		***REMOVED***
		_, err := l.client.Delete(context.Background(), l.key, delOpts)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Close closes the client connection
func (s *Etcd) Close() ***REMOVED***
	return
***REMOVED***
