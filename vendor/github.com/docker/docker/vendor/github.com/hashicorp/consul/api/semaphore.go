package api

import (
	"encoding/json"
	"fmt"
	"path"
	"sync"
	"time"
)

const (
	// DefaultSemaphoreSessionName is the Session Name we assign if none is provided
	DefaultSemaphoreSessionName = "Consul API Semaphore"

	// DefaultSemaphoreSessionTTL is the default session TTL if no Session is provided
	// when creating a new Semaphore. This is used because we do not have another
	// other check to depend upon.
	DefaultSemaphoreSessionTTL = "15s"

	// DefaultSemaphoreWaitTime is how long we block for at a time to check if semaphore
	// acquisition is possible. This affects the minimum time it takes to cancel
	// a Semaphore acquisition.
	DefaultSemaphoreWaitTime = 15 * time.Second

	// DefaultSemaphoreKey is the key used within the prefix to
	// use for coordination between all the contenders.
	DefaultSemaphoreKey = ".lock"

	// SemaphoreFlagValue is a magic flag we set to indicate a key
	// is being used for a semaphore. It is used to detect a potential
	// conflict with a lock.
	SemaphoreFlagValue = 0xe0f69a2baa414de0
)

var (
	// ErrSemaphoreHeld is returned if we attempt to double lock
	ErrSemaphoreHeld = fmt.Errorf("Semaphore already held")

	// ErrSemaphoreNotHeld is returned if we attempt to unlock a semaphore
	// that we do not hold.
	ErrSemaphoreNotHeld = fmt.Errorf("Semaphore not held")

	// ErrSemaphoreInUse is returned if we attempt to destroy a semaphore
	// that is in use.
	ErrSemaphoreInUse = fmt.Errorf("Semaphore in use")

	// ErrSemaphoreConflict is returned if the flags on a key
	// used for a semaphore do not match expectation
	ErrSemaphoreConflict = fmt.Errorf("Existing key does not match semaphore use")
)

// Semaphore is used to implement a distributed semaphore
// using the Consul KV primitives.
type Semaphore struct ***REMOVED***
	c    *Client
	opts *SemaphoreOptions

	isHeld       bool
	sessionRenew chan struct***REMOVED******REMOVED***
	lockSession  string
	l            sync.Mutex
***REMOVED***

// SemaphoreOptions is used to parameterize the Semaphore
type SemaphoreOptions struct ***REMOVED***
	Prefix      string // Must be set and have write permissions
	Limit       int    // Must be set, and be positive
	Value       []byte // Optional, value to associate with the contender entry
	Session     string // OPtional, created if not specified
	SessionName string // Optional, defaults to DefaultLockSessionName
	SessionTTL  string // Optional, defaults to DefaultLockSessionTTL
***REMOVED***

// semaphoreLock is written under the DefaultSemaphoreKey and
// is used to coordinate between all the contenders.
type semaphoreLock struct ***REMOVED***
	// Limit is the integer limit of holders. This is used to
	// verify that all the holders agree on the value.
	Limit int

	// Holders is a list of all the semaphore holders.
	// It maps the session ID to true. It is used as a set effectively.
	Holders map[string]bool
***REMOVED***

// SemaphorePrefix is used to created a Semaphore which will operate
// at the given KV prefix and uses the given limit for the semaphore.
// The prefix must have write privileges, and the limit must be agreed
// upon by all contenders.
func (c *Client) SemaphorePrefix(prefix string, limit int) (*Semaphore, error) ***REMOVED***
	opts := &SemaphoreOptions***REMOVED***
		Prefix: prefix,
		Limit:  limit,
	***REMOVED***
	return c.SemaphoreOpts(opts)
***REMOVED***

// SemaphoreOpts is used to create a Semaphore with the given options.
// The prefix must have write privileges, and the limit must be agreed
// upon by all contenders. If a Session is not provided, one will be created.
func (c *Client) SemaphoreOpts(opts *SemaphoreOptions) (*Semaphore, error) ***REMOVED***
	if opts.Prefix == "" ***REMOVED***
		return nil, fmt.Errorf("missing prefix")
	***REMOVED***
	if opts.Limit <= 0 ***REMOVED***
		return nil, fmt.Errorf("semaphore limit must be positive")
	***REMOVED***
	if opts.SessionName == "" ***REMOVED***
		opts.SessionName = DefaultSemaphoreSessionName
	***REMOVED***
	if opts.SessionTTL == "" ***REMOVED***
		opts.SessionTTL = DefaultSemaphoreSessionTTL
	***REMOVED*** else ***REMOVED***
		if _, err := time.ParseDuration(opts.SessionTTL); err != nil ***REMOVED***
			return nil, fmt.Errorf("invalid SessionTTL: %v", err)
		***REMOVED***
	***REMOVED***
	s := &Semaphore***REMOVED***
		c:    c,
		opts: opts,
	***REMOVED***
	return s, nil
***REMOVED***

// Acquire attempts to reserve a slot in the semaphore, blocking until
// success, interrupted via the stopCh or an error is encounted.
// Providing a non-nil stopCh can be used to abort the attempt.
// On success, a channel is returned that represents our slot.
// This channel could be closed at any time due to session invalidation,
// communication errors, operator intervention, etc. It is NOT safe to
// assume that the slot is held until Release() unless the Session is specifically
// created without any associated health checks. By default Consul sessions
// prefer liveness over safety and an application must be able to handle
// the session being lost.
func (s *Semaphore) Acquire(stopCh <-chan struct***REMOVED******REMOVED***) (<-chan struct***REMOVED******REMOVED***, error) ***REMOVED***
	// Hold the lock as we try to acquire
	s.l.Lock()
	defer s.l.Unlock()

	// Check if we already hold the semaphore
	if s.isHeld ***REMOVED***
		return nil, ErrSemaphoreHeld
	***REMOVED***

	// Check if we need to create a session first
	s.lockSession = s.opts.Session
	if s.lockSession == "" ***REMOVED***
		if sess, err := s.createSession(); err != nil ***REMOVED***
			return nil, fmt.Errorf("failed to create session: %v", err)
		***REMOVED*** else ***REMOVED***
			s.sessionRenew = make(chan struct***REMOVED******REMOVED***)
			s.lockSession = sess
			session := s.c.Session()
			go session.RenewPeriodic(s.opts.SessionTTL, sess, nil, s.sessionRenew)

			// If we fail to acquire the lock, cleanup the session
			defer func() ***REMOVED***
				if !s.isHeld ***REMOVED***
					close(s.sessionRenew)
					s.sessionRenew = nil
				***REMOVED***
			***REMOVED***()
		***REMOVED***
	***REMOVED***

	// Create the contender entry
	kv := s.c.KV()
	made, _, err := kv.Acquire(s.contenderEntry(s.lockSession), nil)
	if err != nil || !made ***REMOVED***
		return nil, fmt.Errorf("failed to make contender entry: %v", err)
	***REMOVED***

	// Setup the query options
	qOpts := &QueryOptions***REMOVED***
		WaitTime: DefaultSemaphoreWaitTime,
	***REMOVED***

WAIT:
	// Check if we should quit
	select ***REMOVED***
	case <-stopCh:
		return nil, nil
	default:
	***REMOVED***

	// Read the prefix
	pairs, meta, err := kv.List(s.opts.Prefix, qOpts)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to read prefix: %v", err)
	***REMOVED***

	// Decode the lock
	lockPair := s.findLock(pairs)
	if lockPair.Flags != SemaphoreFlagValue ***REMOVED***
		return nil, ErrSemaphoreConflict
	***REMOVED***
	lock, err := s.decodeLock(lockPair)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Verify we agree with the limit
	if lock.Limit != s.opts.Limit ***REMOVED***
		return nil, fmt.Errorf("semaphore limit conflict (lock: %d, local: %d)",
			lock.Limit, s.opts.Limit)
	***REMOVED***

	// Prune the dead holders
	s.pruneDeadHolders(lock, pairs)

	// Check if the lock is held
	if len(lock.Holders) >= lock.Limit ***REMOVED***
		qOpts.WaitIndex = meta.LastIndex
		goto WAIT
	***REMOVED***

	// Create a new lock with us as a holder
	lock.Holders[s.lockSession] = true
	newLock, err := s.encodeLock(lock, lockPair.ModifyIndex)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Attempt the acquisition
	didSet, _, err := kv.CAS(newLock, nil)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to update lock: %v", err)
	***REMOVED***
	if !didSet ***REMOVED***
		// Update failed, could have been a race with another contender,
		// retry the operation
		goto WAIT
	***REMOVED***

	// Watch to ensure we maintain ownership of the slot
	lockCh := make(chan struct***REMOVED******REMOVED***)
	go s.monitorLock(s.lockSession, lockCh)

	// Set that we own the lock
	s.isHeld = true

	// Acquired! All done
	return lockCh, nil
***REMOVED***

// Release is used to voluntarily give up our semaphore slot. It is
// an error to call this if the semaphore has not been acquired.
func (s *Semaphore) Release() error ***REMOVED***
	// Hold the lock as we try to release
	s.l.Lock()
	defer s.l.Unlock()

	// Ensure the lock is actually held
	if !s.isHeld ***REMOVED***
		return ErrSemaphoreNotHeld
	***REMOVED***

	// Set that we no longer own the lock
	s.isHeld = false

	// Stop the session renew
	if s.sessionRenew != nil ***REMOVED***
		defer func() ***REMOVED***
			close(s.sessionRenew)
			s.sessionRenew = nil
		***REMOVED***()
	***REMOVED***

	// Get and clear the lock session
	lockSession := s.lockSession
	s.lockSession = ""

	// Remove ourselves as a lock holder
	kv := s.c.KV()
	key := path.Join(s.opts.Prefix, DefaultSemaphoreKey)
READ:
	pair, _, err := kv.Get(key, nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if pair == nil ***REMOVED***
		pair = &KVPair***REMOVED******REMOVED***
	***REMOVED***
	lock, err := s.decodeLock(pair)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Create a new lock without us as a holder
	if _, ok := lock.Holders[lockSession]; ok ***REMOVED***
		delete(lock.Holders, lockSession)
		newLock, err := s.encodeLock(lock, pair.ModifyIndex)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Swap the locks
		didSet, _, err := kv.CAS(newLock, nil)
		if err != nil ***REMOVED***
			return fmt.Errorf("failed to update lock: %v", err)
		***REMOVED***
		if !didSet ***REMOVED***
			goto READ
		***REMOVED***
	***REMOVED***

	// Destroy the contender entry
	contenderKey := path.Join(s.opts.Prefix, lockSession)
	if _, err := kv.Delete(contenderKey, nil); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// Destroy is used to cleanup the semaphore entry. It is not necessary
// to invoke. It will fail if the semaphore is in use.
func (s *Semaphore) Destroy() error ***REMOVED***
	// Hold the lock as we try to acquire
	s.l.Lock()
	defer s.l.Unlock()

	// Check if we already hold the semaphore
	if s.isHeld ***REMOVED***
		return ErrSemaphoreHeld
	***REMOVED***

	// List for the semaphore
	kv := s.c.KV()
	pairs, _, err := kv.List(s.opts.Prefix, nil)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to read prefix: %v", err)
	***REMOVED***

	// Find the lock pair, bail if it doesn't exist
	lockPair := s.findLock(pairs)
	if lockPair.ModifyIndex == 0 ***REMOVED***
		return nil
	***REMOVED***
	if lockPair.Flags != SemaphoreFlagValue ***REMOVED***
		return ErrSemaphoreConflict
	***REMOVED***

	// Decode the lock
	lock, err := s.decodeLock(lockPair)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Prune the dead holders
	s.pruneDeadHolders(lock, pairs)

	// Check if there are any holders
	if len(lock.Holders) > 0 ***REMOVED***
		return ErrSemaphoreInUse
	***REMOVED***

	// Attempt the delete
	didRemove, _, err := kv.DeleteCAS(lockPair, nil)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to remove semaphore: %v", err)
	***REMOVED***
	if !didRemove ***REMOVED***
		return ErrSemaphoreInUse
	***REMOVED***
	return nil
***REMOVED***

// createSession is used to create a new managed session
func (s *Semaphore) createSession() (string, error) ***REMOVED***
	session := s.c.Session()
	se := &SessionEntry***REMOVED***
		Name:     s.opts.SessionName,
		TTL:      s.opts.SessionTTL,
		Behavior: SessionBehaviorDelete,
	***REMOVED***
	id, _, err := session.Create(se, nil)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return id, nil
***REMOVED***

// contenderEntry returns a formatted KVPair for the contender
func (s *Semaphore) contenderEntry(session string) *KVPair ***REMOVED***
	return &KVPair***REMOVED***
		Key:     path.Join(s.opts.Prefix, session),
		Value:   s.opts.Value,
		Session: session,
		Flags:   SemaphoreFlagValue,
	***REMOVED***
***REMOVED***

// findLock is used to find the KV Pair which is used for coordination
func (s *Semaphore) findLock(pairs KVPairs) *KVPair ***REMOVED***
	key := path.Join(s.opts.Prefix, DefaultSemaphoreKey)
	for _, pair := range pairs ***REMOVED***
		if pair.Key == key ***REMOVED***
			return pair
		***REMOVED***
	***REMOVED***
	return &KVPair***REMOVED***Flags: SemaphoreFlagValue***REMOVED***
***REMOVED***

// decodeLock is used to decode a semaphoreLock from an
// entry in Consul
func (s *Semaphore) decodeLock(pair *KVPair) (*semaphoreLock, error) ***REMOVED***
	// Handle if there is no lock
	if pair == nil || pair.Value == nil ***REMOVED***
		return &semaphoreLock***REMOVED***
			Limit:   s.opts.Limit,
			Holders: make(map[string]bool),
		***REMOVED***, nil
	***REMOVED***

	l := &semaphoreLock***REMOVED******REMOVED***
	if err := json.Unmarshal(pair.Value, l); err != nil ***REMOVED***
		return nil, fmt.Errorf("lock decoding failed: %v", err)
	***REMOVED***
	return l, nil
***REMOVED***

// encodeLock is used to encode a semaphoreLock into a KVPair
// that can be PUT
func (s *Semaphore) encodeLock(l *semaphoreLock, oldIndex uint64) (*KVPair, error) ***REMOVED***
	enc, err := json.Marshal(l)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("lock encoding failed: %v", err)
	***REMOVED***
	pair := &KVPair***REMOVED***
		Key:         path.Join(s.opts.Prefix, DefaultSemaphoreKey),
		Value:       enc,
		Flags:       SemaphoreFlagValue,
		ModifyIndex: oldIndex,
	***REMOVED***
	return pair, nil
***REMOVED***

// pruneDeadHolders is used to remove all the dead lock holders
func (s *Semaphore) pruneDeadHolders(lock *semaphoreLock, pairs KVPairs) ***REMOVED***
	// Gather all the live holders
	alive := make(map[string]struct***REMOVED******REMOVED***, len(pairs))
	for _, pair := range pairs ***REMOVED***
		if pair.Session != "" ***REMOVED***
			alive[pair.Session] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***

	// Remove any holders that are dead
	for holder := range lock.Holders ***REMOVED***
		if _, ok := alive[holder]; !ok ***REMOVED***
			delete(lock.Holders, holder)
		***REMOVED***
	***REMOVED***
***REMOVED***

// monitorLock is a long running routine to monitor a semaphore ownership
// It closes the stopCh if we lose our slot.
func (s *Semaphore) monitorLock(session string, stopCh chan struct***REMOVED******REMOVED***) ***REMOVED***
	defer close(stopCh)
	kv := s.c.KV()
	opts := &QueryOptions***REMOVED***RequireConsistent: true***REMOVED***
WAIT:
	pairs, meta, err := kv.List(s.opts.Prefix, opts)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	lockPair := s.findLock(pairs)
	lock, err := s.decodeLock(lockPair)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	s.pruneDeadHolders(lock, pairs)
	if _, ok := lock.Holders[session]; ok ***REMOVED***
		opts.WaitIndex = meta.LastIndex
		goto WAIT
	***REMOVED***
***REMOVED***
