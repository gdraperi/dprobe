package consul

import (
	"crypto/tls"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	api "github.com/hashicorp/consul/api"
)

const (
	// DefaultWatchWaitTime is how long we block for at a
	// time to check if the watched key has changed. This
	// affects the minimum time it takes to cancel a watch.
	DefaultWatchWaitTime = 15 * time.Second

	// RenewSessionRetryMax is the number of time we should try
	// to renew the session before giving up and throwing an error
	RenewSessionRetryMax = 5

	// MaxSessionDestroyAttempts is the maximum times we will try
	// to explicitely destroy the session attached to a lock after
	// the connectivity to the store has been lost
	MaxSessionDestroyAttempts = 5

	// defaultLockTTL is the default ttl for the consul lock
	defaultLockTTL = 20 * time.Second
)

var (
	// ErrMultipleEndpointsUnsupported is thrown when there are
	// multiple endpoints specified for Consul
	ErrMultipleEndpointsUnsupported = errors.New("consul does not support multiple endpoints")

	// ErrSessionRenew is thrown when the session can't be
	// renewed because the Consul version does not support sessions
	ErrSessionRenew = errors.New("cannot set or renew session for ttl, unable to operate on sessions")
)

// Consul is the receiver type for the
// Store interface
type Consul struct ***REMOVED***
	sync.Mutex
	config *api.Config
	client *api.Client
***REMOVED***

type consulLock struct ***REMOVED***
	lock    *api.Lock
	renewCh chan struct***REMOVED******REMOVED***
***REMOVED***

// Register registers consul to libkv
func Register() ***REMOVED***
	libkv.AddStore(store.CONSUL, New)
***REMOVED***

// New creates a new Consul client given a list
// of endpoints and optional tls config
func New(endpoints []string, options *store.Config) (store.Store, error) ***REMOVED***
	if len(endpoints) > 1 ***REMOVED***
		return nil, ErrMultipleEndpointsUnsupported
	***REMOVED***

	s := &Consul***REMOVED******REMOVED***

	// Create Consul client
	config := api.DefaultConfig()
	s.config = config
	config.HttpClient = http.DefaultClient
	config.Address = endpoints[0]
	config.Scheme = "http"

	// Set options
	if options != nil ***REMOVED***
		if options.TLS != nil ***REMOVED***
			s.setTLS(options.TLS)
		***REMOVED***
		if options.ConnectionTimeout != 0 ***REMOVED***
			s.setTimeout(options.ConnectionTimeout)
		***REMOVED***
	***REMOVED***

	// Creates a new client
	client, err := api.NewClient(config)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	s.client = client

	return s, nil
***REMOVED***

// SetTLS sets Consul TLS options
func (s *Consul) setTLS(tls *tls.Config) ***REMOVED***
	s.config.HttpClient.Transport = &http.Transport***REMOVED***
		TLSClientConfig: tls,
	***REMOVED***
	s.config.Scheme = "https"
***REMOVED***

// SetTimeout sets the timeout for connecting to Consul
func (s *Consul) setTimeout(time time.Duration) ***REMOVED***
	s.config.WaitTime = time
***REMOVED***

// Normalize the key for usage in Consul
func (s *Consul) normalize(key string) string ***REMOVED***
	key = store.Normalize(key)
	return strings.TrimPrefix(key, "/")
***REMOVED***

func (s *Consul) renewSession(pair *api.KVPair, ttl time.Duration) error ***REMOVED***
	// Check if there is any previous session with an active TTL
	session, err := s.getActiveSession(pair.Key)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if session == "" ***REMOVED***
		entry := &api.SessionEntry***REMOVED***
			Behavior:  api.SessionBehaviorDelete, // Delete the key when the session expires
			TTL:       (ttl / 2).String(),        // Consul multiplies the TTL by 2x
			LockDelay: 1 * time.Millisecond,      // Virtually disable lock delay
		***REMOVED***

		// Create the key session
		session, _, err = s.client.Session().Create(entry, nil)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		lockOpts := &api.LockOptions***REMOVED***
			Key:     pair.Key,
			Session: session,
		***REMOVED***

		// Lock and ignore if lock is held
		// It's just a placeholder for the
		// ephemeral behavior
		lock, _ := s.client.LockOpts(lockOpts)
		if lock != nil ***REMOVED***
			lock.Lock(nil)
		***REMOVED***
	***REMOVED***

	_, _, err = s.client.Session().Renew(session, nil)
	return err
***REMOVED***

// getActiveSession checks if the key already has
// a session attached
func (s *Consul) getActiveSession(key string) (string, error) ***REMOVED***
	pair, _, err := s.client.KV().Get(key, nil)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if pair != nil && pair.Session != "" ***REMOVED***
		return pair.Session, nil
	***REMOVED***
	return "", nil
***REMOVED***

// Get the value at "key", returns the last modified index
// to use in conjunction to CAS calls
func (s *Consul) Get(key string) (*store.KVPair, error) ***REMOVED***
	options := &api.QueryOptions***REMOVED***
		AllowStale:        false,
		RequireConsistent: true,
	***REMOVED***

	pair, meta, err := s.client.KV().Get(s.normalize(key), options)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// If pair is nil then the key does not exist
	if pair == nil ***REMOVED***
		return nil, store.ErrKeyNotFound
	***REMOVED***

	return &store.KVPair***REMOVED***Key: pair.Key, Value: pair.Value, LastIndex: meta.LastIndex***REMOVED***, nil
***REMOVED***

// Put a value at "key"
func (s *Consul) Put(key string, value []byte, opts *store.WriteOptions) error ***REMOVED***
	key = s.normalize(key)

	p := &api.KVPair***REMOVED***
		Key:   key,
		Value: value,
		Flags: api.LockFlagValue,
	***REMOVED***

	if opts != nil && opts.TTL > 0 ***REMOVED***
		// Create or renew a session holding a TTL. Operations on sessions
		// are not deterministic: creating or renewing a session can fail
		for retry := 1; retry <= RenewSessionRetryMax; retry++ ***REMOVED***
			err := s.renewSession(p, opts.TTL)
			if err == nil ***REMOVED***
				break
			***REMOVED***
			if retry == RenewSessionRetryMax ***REMOVED***
				return ErrSessionRenew
			***REMOVED***
		***REMOVED***
	***REMOVED***

	_, err := s.client.KV().Put(p, nil)
	return err
***REMOVED***

// Delete a value at "key"
func (s *Consul) Delete(key string) error ***REMOVED***
	if _, err := s.Get(key); err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err := s.client.KV().Delete(s.normalize(key), nil)
	return err
***REMOVED***

// Exists checks that the key exists inside the store
func (s *Consul) Exists(key string) (bool, error) ***REMOVED***
	_, err := s.Get(key)
	if err != nil ***REMOVED***
		if err == store.ErrKeyNotFound ***REMOVED***
			return false, nil
		***REMOVED***
		return false, err
	***REMOVED***
	return true, nil
***REMOVED***

// List child nodes of a given directory
func (s *Consul) List(directory string) ([]*store.KVPair, error) ***REMOVED***
	pairs, _, err := s.client.KV().List(s.normalize(directory), nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(pairs) == 0 ***REMOVED***
		return nil, store.ErrKeyNotFound
	***REMOVED***

	kv := []*store.KVPair***REMOVED******REMOVED***

	for _, pair := range pairs ***REMOVED***
		if pair.Key == directory ***REMOVED***
			continue
		***REMOVED***
		kv = append(kv, &store.KVPair***REMOVED***
			Key:       pair.Key,
			Value:     pair.Value,
			LastIndex: pair.ModifyIndex,
		***REMOVED***)
	***REMOVED***

	return kv, nil
***REMOVED***

// DeleteTree deletes a range of keys under a given directory
func (s *Consul) DeleteTree(directory string) error ***REMOVED***
	if _, err := s.List(directory); err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err := s.client.KV().DeleteTree(s.normalize(directory), nil)
	return err
***REMOVED***

// Watch for changes on a "key"
// It returns a channel that will receive changes or pass
// on errors. Upon creation, the current value will first
// be sent to the channel. Providing a non-nil stopCh can
// be used to stop watching.
func (s *Consul) Watch(key string, stopCh <-chan struct***REMOVED******REMOVED***) (<-chan *store.KVPair, error) ***REMOVED***
	kv := s.client.KV()
	watchCh := make(chan *store.KVPair)

	go func() ***REMOVED***
		defer close(watchCh)

		// Use a wait time in order to check if we should quit
		// from time to time.
		opts := &api.QueryOptions***REMOVED***WaitTime: DefaultWatchWaitTime***REMOVED***

		for ***REMOVED***
			// Check if we should quit
			select ***REMOVED***
			case <-stopCh:
				return
			default:
			***REMOVED***

			// Get the key
			pair, meta, err := kv.Get(key, opts)
			if err != nil ***REMOVED***
				return
			***REMOVED***

			// If LastIndex didn't change then it means `Get` returned
			// because of the WaitTime and the key didn't changed.
			if opts.WaitIndex == meta.LastIndex ***REMOVED***
				continue
			***REMOVED***
			opts.WaitIndex = meta.LastIndex

			// Return the value to the channel
			// FIXME: What happens when a key is deleted?
			if pair != nil ***REMOVED***
				watchCh <- &store.KVPair***REMOVED***
					Key:       pair.Key,
					Value:     pair.Value,
					LastIndex: pair.ModifyIndex,
				***REMOVED***
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
func (s *Consul) WatchTree(directory string, stopCh <-chan struct***REMOVED******REMOVED***) (<-chan []*store.KVPair, error) ***REMOVED***
	kv := s.client.KV()
	watchCh := make(chan []*store.KVPair)

	go func() ***REMOVED***
		defer close(watchCh)

		// Use a wait time in order to check if we should quit
		// from time to time.
		opts := &api.QueryOptions***REMOVED***WaitTime: DefaultWatchWaitTime***REMOVED***
		for ***REMOVED***
			// Check if we should quit
			select ***REMOVED***
			case <-stopCh:
				return
			default:
			***REMOVED***

			// Get all the childrens
			pairs, meta, err := kv.List(directory, opts)
			if err != nil ***REMOVED***
				return
			***REMOVED***

			// If LastIndex didn't change then it means `Get` returned
			// because of the WaitTime and the child keys didn't change.
			if opts.WaitIndex == meta.LastIndex ***REMOVED***
				continue
			***REMOVED***
			opts.WaitIndex = meta.LastIndex

			// Return children KV pairs to the channel
			kvpairs := []*store.KVPair***REMOVED******REMOVED***
			for _, pair := range pairs ***REMOVED***
				if pair.Key == directory ***REMOVED***
					continue
				***REMOVED***
				kvpairs = append(kvpairs, &store.KVPair***REMOVED***
					Key:       pair.Key,
					Value:     pair.Value,
					LastIndex: pair.ModifyIndex,
				***REMOVED***)
			***REMOVED***
			watchCh <- kvpairs
		***REMOVED***
	***REMOVED***()

	return watchCh, nil
***REMOVED***

// NewLock returns a handle to a lock struct which can
// be used to provide mutual exclusion on a key
func (s *Consul) NewLock(key string, options *store.LockOptions) (store.Locker, error) ***REMOVED***
	lockOpts := &api.LockOptions***REMOVED***
		Key: s.normalize(key),
	***REMOVED***

	lock := &consulLock***REMOVED******REMOVED***

	ttl := defaultLockTTL

	if options != nil ***REMOVED***
		// Set optional TTL on Lock
		if options.TTL != 0 ***REMOVED***
			ttl = options.TTL
		***REMOVED***
		// Set optional value on Lock
		if options.Value != nil ***REMOVED***
			lockOpts.Value = options.Value
		***REMOVED***
	***REMOVED***

	entry := &api.SessionEntry***REMOVED***
		Behavior:  api.SessionBehaviorRelease, // Release the lock when the session expires
		TTL:       (ttl / 2).String(),         // Consul multiplies the TTL by 2x
		LockDelay: 1 * time.Millisecond,       // Virtually disable lock delay
	***REMOVED***

	// Create the key session
	session, _, err := s.client.Session().Create(entry, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Place the session and renew chan on lock
	lockOpts.Session = session
	lock.renewCh = options.RenewLock

	l, err := s.client.LockOpts(lockOpts)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Renew the session ttl lock periodically
	s.renewLockSession(entry.TTL, session, options.RenewLock)

	lock.lock = l
	return lock, nil
***REMOVED***

// renewLockSession is used to renew a session Lock, it takes
// a stopRenew chan which is used to explicitely stop the session
// renew process. The renew routine never stops until a signal is
// sent to this channel. If deleting the session fails because the
// connection to the store is lost, it keeps trying to delete the
// session periodically until it can contact the store, this ensures
// that the lock is not maintained indefinitely which ensures liveness
// over safety for the lock when the store becomes unavailable.
func (s *Consul) renewLockSession(initialTTL string, id string, stopRenew chan struct***REMOVED******REMOVED***) ***REMOVED***
	sessionDestroyAttempts := 0
	ttl, err := time.ParseDuration(initialTTL)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	go func() ***REMOVED***
		for ***REMOVED***
			select ***REMOVED***
			case <-time.After(ttl / 2):
				entry, _, err := s.client.Session().Renew(id, nil)
				if err != nil ***REMOVED***
					// If an error occurs, continue until the
					// session gets destroyed explicitely or
					// the session ttl times out
					continue
				***REMOVED***
				if entry == nil ***REMOVED***
					return
				***REMOVED***

				// Handle the server updating the TTL
				ttl, _ = time.ParseDuration(entry.TTL)

			case <-stopRenew:
				// Attempt a session destroy
				_, err := s.client.Session().Destroy(id, nil)
				if err == nil ***REMOVED***
					return
				***REMOVED***

				if sessionDestroyAttempts >= MaxSessionDestroyAttempts ***REMOVED***
					return
				***REMOVED***

				// We can't destroy the session because the store
				// is unavailable, wait for the session renew period
				sessionDestroyAttempts++
				time.Sleep(ttl / 2)
			***REMOVED***
		***REMOVED***
	***REMOVED***()
***REMOVED***

// Lock attempts to acquire the lock and blocks while
// doing so. It returns a channel that is closed if our
// lock is lost or if an error occurs
func (l *consulLock) Lock(stopChan chan struct***REMOVED******REMOVED***) (<-chan struct***REMOVED******REMOVED***, error) ***REMOVED***
	return l.lock.Lock(stopChan)
***REMOVED***

// Unlock the "key". Calling unlock while
// not holding the lock will throw an error
func (l *consulLock) Unlock() error ***REMOVED***
	if l.renewCh != nil ***REMOVED***
		close(l.renewCh)
	***REMOVED***
	return l.lock.Unlock()
***REMOVED***

// AtomicPut put a value at "key" if the key has not been
// modified in the meantime, throws an error if this is the case
func (s *Consul) AtomicPut(key string, value []byte, previous *store.KVPair, options *store.WriteOptions) (bool, *store.KVPair, error) ***REMOVED***

	p := &api.KVPair***REMOVED***Key: s.normalize(key), Value: value, Flags: api.LockFlagValue***REMOVED***

	if previous == nil ***REMOVED***
		// Consul interprets ModifyIndex = 0 as new key.
		p.ModifyIndex = 0
	***REMOVED*** else ***REMOVED***
		p.ModifyIndex = previous.LastIndex
	***REMOVED***

	ok, _, err := s.client.KV().CAS(p, nil)
	if err != nil ***REMOVED***
		return false, nil, err
	***REMOVED***
	if !ok ***REMOVED***
		if previous == nil ***REMOVED***
			return false, nil, store.ErrKeyExists
		***REMOVED***
		return false, nil, store.ErrKeyModified
	***REMOVED***

	pair, err := s.Get(key)
	if err != nil ***REMOVED***
		return false, nil, err
	***REMOVED***

	return true, pair, nil
***REMOVED***

// AtomicDelete deletes a value at "key" if the key has not
// been modified in the meantime, throws an error if this is the case
func (s *Consul) AtomicDelete(key string, previous *store.KVPair) (bool, error) ***REMOVED***
	if previous == nil ***REMOVED***
		return false, store.ErrPreviousNotSpecified
	***REMOVED***

	p := &api.KVPair***REMOVED***Key: s.normalize(key), ModifyIndex: previous.LastIndex, Flags: api.LockFlagValue***REMOVED***

	// Extra Get operation to check on the key
	_, err := s.Get(key)
	if err != nil && err == store.ErrKeyNotFound ***REMOVED***
		return false, err
	***REMOVED***

	if work, _, err := s.client.KV().DeleteCAS(p, nil); err != nil ***REMOVED***
		return false, err
	***REMOVED*** else if !work ***REMOVED***
		return false, store.ErrKeyModified
	***REMOVED***

	return true, nil
***REMOVED***

// Close closes the client connection
func (s *Consul) Close() ***REMOVED***
	return
***REMOVED***
