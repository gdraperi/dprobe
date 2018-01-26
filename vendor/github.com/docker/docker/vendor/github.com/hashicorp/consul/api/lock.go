package api

import (
	"fmt"
	"sync"
	"time"
)

const (
	// DefaultLockSessionName is the Session Name we assign if none is provided
	DefaultLockSessionName = "Consul API Lock"

	// DefaultLockSessionTTL is the default session TTL if no Session is provided
	// when creating a new Lock. This is used because we do not have another
	// other check to depend upon.
	DefaultLockSessionTTL = "15s"

	// DefaultLockWaitTime is how long we block for at a time to check if lock
	// acquisition is possible. This affects the minimum time it takes to cancel
	// a Lock acquisition.
	DefaultLockWaitTime = 15 * time.Second

	// DefaultLockRetryTime is how long we wait after a failed lock acquisition
	// before attempting to do the lock again. This is so that once a lock-delay
	// is in affect, we do not hot loop retrying the acquisition.
	DefaultLockRetryTime = 5 * time.Second

	// LockFlagValue is a magic flag we set to indicate a key
	// is being used for a lock. It is used to detect a potential
	// conflict with a semaphore.
	LockFlagValue = 0x2ddccbc058a50c18
)

var (
	// ErrLockHeld is returned if we attempt to double lock
	ErrLockHeld = fmt.Errorf("Lock already held")

	// ErrLockNotHeld is returned if we attempt to unlock a lock
	// that we do not hold.
	ErrLockNotHeld = fmt.Errorf("Lock not held")

	// ErrLockInUse is returned if we attempt to destroy a lock
	// that is in use.
	ErrLockInUse = fmt.Errorf("Lock in use")

	// ErrLockConflict is returned if the flags on a key
	// used for a lock do not match expectation
	ErrLockConflict = fmt.Errorf("Existing key does not match lock use")
)

// Lock is used to implement client-side leader election. It is follows the
// algorithm as described here: https://consul.io/docs/guides/leader-election.html.
type Lock struct ***REMOVED***
	c    *Client
	opts *LockOptions

	isHeld       bool
	sessionRenew chan struct***REMOVED******REMOVED***
	lockSession  string
	l            sync.Mutex
***REMOVED***

// LockOptions is used to parameterize the Lock behavior.
type LockOptions struct ***REMOVED***
	Key         string // Must be set and have write permissions
	Value       []byte // Optional, value to associate with the lock
	Session     string // Optional, created if not specified
	SessionName string // Optional, defaults to DefaultLockSessionName
	SessionTTL  string // Optional, defaults to DefaultLockSessionTTL
***REMOVED***

// LockKey returns a handle to a lock struct which can be used
// to acquire and release the mutex. The key used must have
// write permissions.
func (c *Client) LockKey(key string) (*Lock, error) ***REMOVED***
	opts := &LockOptions***REMOVED***
		Key: key,
	***REMOVED***
	return c.LockOpts(opts)
***REMOVED***

// LockOpts returns a handle to a lock struct which can be used
// to acquire and release the mutex. The key used must have
// write permissions.
func (c *Client) LockOpts(opts *LockOptions) (*Lock, error) ***REMOVED***
	if opts.Key == "" ***REMOVED***
		return nil, fmt.Errorf("missing key")
	***REMOVED***
	if opts.SessionName == "" ***REMOVED***
		opts.SessionName = DefaultLockSessionName
	***REMOVED***
	if opts.SessionTTL == "" ***REMOVED***
		opts.SessionTTL = DefaultLockSessionTTL
	***REMOVED*** else ***REMOVED***
		if _, err := time.ParseDuration(opts.SessionTTL); err != nil ***REMOVED***
			return nil, fmt.Errorf("invalid SessionTTL: %v", err)
		***REMOVED***
	***REMOVED***
	l := &Lock***REMOVED***
		c:    c,
		opts: opts,
	***REMOVED***
	return l, nil
***REMOVED***

// Lock attempts to acquire the lock and blocks while doing so.
// Providing a non-nil stopCh can be used to abort the lock attempt.
// Returns a channel that is closed if our lock is lost or an error.
// This channel could be closed at any time due to session invalidation,
// communication errors, operator intervention, etc. It is NOT safe to
// assume that the lock is held until Unlock() unless the Session is specifically
// created without any associated health checks. By default Consul sessions
// prefer liveness over safety and an application must be able to handle
// the lock being lost.
func (l *Lock) Lock(stopCh <-chan struct***REMOVED******REMOVED***) (<-chan struct***REMOVED******REMOVED***, error) ***REMOVED***
	// Hold the lock as we try to acquire
	l.l.Lock()
	defer l.l.Unlock()

	// Check if we already hold the lock
	if l.isHeld ***REMOVED***
		return nil, ErrLockHeld
	***REMOVED***

	// Check if we need to create a session first
	l.lockSession = l.opts.Session
	if l.lockSession == "" ***REMOVED***
		if s, err := l.createSession(); err != nil ***REMOVED***
			return nil, fmt.Errorf("failed to create session: %v", err)
		***REMOVED*** else ***REMOVED***
			l.sessionRenew = make(chan struct***REMOVED******REMOVED***)
			l.lockSession = s
			session := l.c.Session()
			go session.RenewPeriodic(l.opts.SessionTTL, s, nil, l.sessionRenew)

			// If we fail to acquire the lock, cleanup the session
			defer func() ***REMOVED***
				if !l.isHeld ***REMOVED***
					close(l.sessionRenew)
					l.sessionRenew = nil
				***REMOVED***
			***REMOVED***()
		***REMOVED***
	***REMOVED***

	// Setup the query options
	kv := l.c.KV()
	qOpts := &QueryOptions***REMOVED***
		WaitTime: DefaultLockWaitTime,
	***REMOVED***

WAIT:
	// Check if we should quit
	select ***REMOVED***
	case <-stopCh:
		return nil, nil
	default:
	***REMOVED***

	// Look for an existing lock, blocking until not taken
	pair, meta, err := kv.Get(l.opts.Key, qOpts)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to read lock: %v", err)
	***REMOVED***
	if pair != nil && pair.Flags != LockFlagValue ***REMOVED***
		return nil, ErrLockConflict
	***REMOVED***
	locked := false
	if pair != nil && pair.Session == l.lockSession ***REMOVED***
		goto HELD
	***REMOVED***
	if pair != nil && pair.Session != "" ***REMOVED***
		qOpts.WaitIndex = meta.LastIndex
		goto WAIT
	***REMOVED***

	// Try to acquire the lock
	pair = l.lockEntry(l.lockSession)
	locked, _, err = kv.Acquire(pair, nil)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to acquire lock: %v", err)
	***REMOVED***

	// Handle the case of not getting the lock
	if !locked ***REMOVED***
		select ***REMOVED***
		case <-time.After(DefaultLockRetryTime):
			goto WAIT
		case <-stopCh:
			return nil, nil
		***REMOVED***
	***REMOVED***

HELD:
	// Watch to ensure we maintain leadership
	leaderCh := make(chan struct***REMOVED******REMOVED***)
	go l.monitorLock(l.lockSession, leaderCh)

	// Set that we own the lock
	l.isHeld = true

	// Locked! All done
	return leaderCh, nil
***REMOVED***

// Unlock released the lock. It is an error to call this
// if the lock is not currently held.
func (l *Lock) Unlock() error ***REMOVED***
	// Hold the lock as we try to release
	l.l.Lock()
	defer l.l.Unlock()

	// Ensure the lock is actually held
	if !l.isHeld ***REMOVED***
		return ErrLockNotHeld
	***REMOVED***

	// Set that we no longer own the lock
	l.isHeld = false

	// Stop the session renew
	if l.sessionRenew != nil ***REMOVED***
		defer func() ***REMOVED***
			close(l.sessionRenew)
			l.sessionRenew = nil
		***REMOVED***()
	***REMOVED***

	// Get the lock entry, and clear the lock session
	lockEnt := l.lockEntry(l.lockSession)
	l.lockSession = ""

	// Release the lock explicitly
	kv := l.c.KV()
	_, _, err := kv.Release(lockEnt, nil)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to release lock: %v", err)
	***REMOVED***
	return nil
***REMOVED***

// Destroy is used to cleanup the lock entry. It is not necessary
// to invoke. It will fail if the lock is in use.
func (l *Lock) Destroy() error ***REMOVED***
	// Hold the lock as we try to release
	l.l.Lock()
	defer l.l.Unlock()

	// Check if we already hold the lock
	if l.isHeld ***REMOVED***
		return ErrLockHeld
	***REMOVED***

	// Look for an existing lock
	kv := l.c.KV()
	pair, _, err := kv.Get(l.opts.Key, nil)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to read lock: %v", err)
	***REMOVED***

	// Nothing to do if the lock does not exist
	if pair == nil ***REMOVED***
		return nil
	***REMOVED***

	// Check for possible flag conflict
	if pair.Flags != LockFlagValue ***REMOVED***
		return ErrLockConflict
	***REMOVED***

	// Check if it is in use
	if pair.Session != "" ***REMOVED***
		return ErrLockInUse
	***REMOVED***

	// Attempt the delete
	didRemove, _, err := kv.DeleteCAS(pair, nil)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to remove lock: %v", err)
	***REMOVED***
	if !didRemove ***REMOVED***
		return ErrLockInUse
	***REMOVED***
	return nil
***REMOVED***

// createSession is used to create a new managed session
func (l *Lock) createSession() (string, error) ***REMOVED***
	session := l.c.Session()
	se := &SessionEntry***REMOVED***
		Name: l.opts.SessionName,
		TTL:  l.opts.SessionTTL,
	***REMOVED***
	id, _, err := session.Create(se, nil)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return id, nil
***REMOVED***

// lockEntry returns a formatted KVPair for the lock
func (l *Lock) lockEntry(session string) *KVPair ***REMOVED***
	return &KVPair***REMOVED***
		Key:     l.opts.Key,
		Value:   l.opts.Value,
		Session: session,
		Flags:   LockFlagValue,
	***REMOVED***
***REMOVED***

// monitorLock is a long running routine to monitor a lock ownership
// It closes the stopCh if we lose our leadership.
func (l *Lock) monitorLock(session string, stopCh chan struct***REMOVED******REMOVED***) ***REMOVED***
	defer close(stopCh)
	kv := l.c.KV()
	opts := &QueryOptions***REMOVED***RequireConsistent: true***REMOVED***
WAIT:
	pair, meta, err := kv.Get(l.opts.Key, opts)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if pair != nil && pair.Session == session ***REMOVED***
		opts.WaitIndex = meta.LastIndex
		goto WAIT
	***REMOVED***
***REMOVED***
