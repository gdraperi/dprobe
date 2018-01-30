// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webdav

import (
	"container/heap"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	// ErrConfirmationFailed is returned by a LockSystem's Confirm method.
	ErrConfirmationFailed = errors.New("webdav: confirmation failed")
	// ErrForbidden is returned by a LockSystem's Unlock method.
	ErrForbidden = errors.New("webdav: forbidden")
	// ErrLocked is returned by a LockSystem's Create, Refresh and Unlock methods.
	ErrLocked = errors.New("webdav: locked")
	// ErrNoSuchLock is returned by a LockSystem's Refresh and Unlock methods.
	ErrNoSuchLock = errors.New("webdav: no such lock")
)

// Condition can match a WebDAV resource, based on a token or ETag.
// Exactly one of Token and ETag should be non-empty.
type Condition struct ***REMOVED***
	Not   bool
	Token string
	ETag  string
***REMOVED***

// LockSystem manages access to a collection of named resources. The elements
// in a lock name are separated by slash ('/', U+002F) characters, regardless
// of host operating system convention.
type LockSystem interface ***REMOVED***
	// Confirm confirms that the caller can claim all of the locks specified by
	// the given conditions, and that holding the union of all of those locks
	// gives exclusive access to all of the named resources. Up to two resources
	// can be named. Empty names are ignored.
	//
	// Exactly one of release and err will be non-nil. If release is non-nil,
	// all of the requested locks are held until release is called. Calling
	// release does not unlock the lock, in the WebDAV UNLOCK sense, but once
	// Confirm has confirmed that a lock claim is valid, that lock cannot be
	// Confirmed again until it has been released.
	//
	// If Confirm returns ErrConfirmationFailed then the Handler will continue
	// to try any other set of locks presented (a WebDAV HTTP request can
	// present more than one set of locks). If it returns any other non-nil
	// error, the Handler will write a "500 Internal Server Error" HTTP status.
	Confirm(now time.Time, name0, name1 string, conditions ...Condition) (release func(), err error)

	// Create creates a lock with the given depth, duration, owner and root
	// (name). The depth will either be negative (meaning infinite) or zero.
	//
	// If Create returns ErrLocked then the Handler will write a "423 Locked"
	// HTTP status. If it returns any other non-nil error, the Handler will
	// write a "500 Internal Server Error" HTTP status.
	//
	// See http://www.webdav.org/specs/rfc4918.html#rfc.section.9.10.6 for
	// when to use each error.
	//
	// The token returned identifies the created lock. It should be an absolute
	// URI as defined by RFC 3986, Section 4.3. In particular, it should not
	// contain whitespace.
	Create(now time.Time, details LockDetails) (token string, err error)

	// Refresh refreshes the lock with the given token.
	//
	// If Refresh returns ErrLocked then the Handler will write a "423 Locked"
	// HTTP Status. If Refresh returns ErrNoSuchLock then the Handler will write
	// a "412 Precondition Failed" HTTP Status. If it returns any other non-nil
	// error, the Handler will write a "500 Internal Server Error" HTTP status.
	//
	// See http://www.webdav.org/specs/rfc4918.html#rfc.section.9.10.6 for
	// when to use each error.
	Refresh(now time.Time, token string, duration time.Duration) (LockDetails, error)

	// Unlock unlocks the lock with the given token.
	//
	// If Unlock returns ErrForbidden then the Handler will write a "403
	// Forbidden" HTTP Status. If Unlock returns ErrLocked then the Handler
	// will write a "423 Locked" HTTP status. If Unlock returns ErrNoSuchLock
	// then the Handler will write a "409 Conflict" HTTP Status. If it returns
	// any other non-nil error, the Handler will write a "500 Internal Server
	// Error" HTTP status.
	//
	// See http://www.webdav.org/specs/rfc4918.html#rfc.section.9.11.1 for
	// when to use each error.
	Unlock(now time.Time, token string) error
***REMOVED***

// LockDetails are a lock's metadata.
type LockDetails struct ***REMOVED***
	// Root is the root resource name being locked. For a zero-depth lock, the
	// root is the only resource being locked.
	Root string
	// Duration is the lock timeout. A negative duration means infinite.
	Duration time.Duration
	// OwnerXML is the verbatim <owner> XML given in a LOCK HTTP request.
	//
	// TODO: does the "verbatim" nature play well with XML namespaces?
	// Does the OwnerXML field need to have more structure? See
	// https://codereview.appspot.com/175140043/#msg2
	OwnerXML string
	// ZeroDepth is whether the lock has zero depth. If it does not have zero
	// depth, it has infinite depth.
	ZeroDepth bool
***REMOVED***

// NewMemLS returns a new in-memory LockSystem.
func NewMemLS() LockSystem ***REMOVED***
	return &memLS***REMOVED***
		byName:  make(map[string]*memLSNode),
		byToken: make(map[string]*memLSNode),
		gen:     uint64(time.Now().Unix()),
	***REMOVED***
***REMOVED***

type memLS struct ***REMOVED***
	mu      sync.Mutex
	byName  map[string]*memLSNode
	byToken map[string]*memLSNode
	gen     uint64
	// byExpiry only contains those nodes whose LockDetails have a finite
	// Duration and are yet to expire.
	byExpiry byExpiry
***REMOVED***

func (m *memLS) nextToken() string ***REMOVED***
	m.gen++
	return strconv.FormatUint(m.gen, 10)
***REMOVED***

func (m *memLS) collectExpiredNodes(now time.Time) ***REMOVED***
	for len(m.byExpiry) > 0 ***REMOVED***
		if now.Before(m.byExpiry[0].expiry) ***REMOVED***
			break
		***REMOVED***
		m.remove(m.byExpiry[0])
	***REMOVED***
***REMOVED***

func (m *memLS) Confirm(now time.Time, name0, name1 string, conditions ...Condition) (func(), error) ***REMOVED***
	m.mu.Lock()
	defer m.mu.Unlock()
	m.collectExpiredNodes(now)

	var n0, n1 *memLSNode
	if name0 != "" ***REMOVED***
		if n0 = m.lookup(slashClean(name0), conditions...); n0 == nil ***REMOVED***
			return nil, ErrConfirmationFailed
		***REMOVED***
	***REMOVED***
	if name1 != "" ***REMOVED***
		if n1 = m.lookup(slashClean(name1), conditions...); n1 == nil ***REMOVED***
			return nil, ErrConfirmationFailed
		***REMOVED***
	***REMOVED***

	// Don't hold the same node twice.
	if n1 == n0 ***REMOVED***
		n1 = nil
	***REMOVED***

	if n0 != nil ***REMOVED***
		m.hold(n0)
	***REMOVED***
	if n1 != nil ***REMOVED***
		m.hold(n1)
	***REMOVED***
	return func() ***REMOVED***
		m.mu.Lock()
		defer m.mu.Unlock()
		if n1 != nil ***REMOVED***
			m.unhold(n1)
		***REMOVED***
		if n0 != nil ***REMOVED***
			m.unhold(n0)
		***REMOVED***
	***REMOVED***, nil
***REMOVED***

// lookup returns the node n that locks the named resource, provided that n
// matches at least one of the given conditions and that lock isn't held by
// another party. Otherwise, it returns nil.
//
// n may be a parent of the named resource, if n is an infinite depth lock.
func (m *memLS) lookup(name string, conditions ...Condition) (n *memLSNode) ***REMOVED***
	// TODO: support Condition.Not and Condition.ETag.
	for _, c := range conditions ***REMOVED***
		n = m.byToken[c.Token]
		if n == nil || n.held ***REMOVED***
			continue
		***REMOVED***
		if name == n.details.Root ***REMOVED***
			return n
		***REMOVED***
		if n.details.ZeroDepth ***REMOVED***
			continue
		***REMOVED***
		if n.details.Root == "/" || strings.HasPrefix(name, n.details.Root+"/") ***REMOVED***
			return n
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (m *memLS) hold(n *memLSNode) ***REMOVED***
	if n.held ***REMOVED***
		panic("webdav: memLS inconsistent held state")
	***REMOVED***
	n.held = true
	if n.details.Duration >= 0 && n.byExpiryIndex >= 0 ***REMOVED***
		heap.Remove(&m.byExpiry, n.byExpiryIndex)
	***REMOVED***
***REMOVED***

func (m *memLS) unhold(n *memLSNode) ***REMOVED***
	if !n.held ***REMOVED***
		panic("webdav: memLS inconsistent held state")
	***REMOVED***
	n.held = false
	if n.details.Duration >= 0 ***REMOVED***
		heap.Push(&m.byExpiry, n)
	***REMOVED***
***REMOVED***

func (m *memLS) Create(now time.Time, details LockDetails) (string, error) ***REMOVED***
	m.mu.Lock()
	defer m.mu.Unlock()
	m.collectExpiredNodes(now)
	details.Root = slashClean(details.Root)

	if !m.canCreate(details.Root, details.ZeroDepth) ***REMOVED***
		return "", ErrLocked
	***REMOVED***
	n := m.create(details.Root)
	n.token = m.nextToken()
	m.byToken[n.token] = n
	n.details = details
	if n.details.Duration >= 0 ***REMOVED***
		n.expiry = now.Add(n.details.Duration)
		heap.Push(&m.byExpiry, n)
	***REMOVED***
	return n.token, nil
***REMOVED***

func (m *memLS) Refresh(now time.Time, token string, duration time.Duration) (LockDetails, error) ***REMOVED***
	m.mu.Lock()
	defer m.mu.Unlock()
	m.collectExpiredNodes(now)

	n := m.byToken[token]
	if n == nil ***REMOVED***
		return LockDetails***REMOVED******REMOVED***, ErrNoSuchLock
	***REMOVED***
	if n.held ***REMOVED***
		return LockDetails***REMOVED******REMOVED***, ErrLocked
	***REMOVED***
	if n.byExpiryIndex >= 0 ***REMOVED***
		heap.Remove(&m.byExpiry, n.byExpiryIndex)
	***REMOVED***
	n.details.Duration = duration
	if n.details.Duration >= 0 ***REMOVED***
		n.expiry = now.Add(n.details.Duration)
		heap.Push(&m.byExpiry, n)
	***REMOVED***
	return n.details, nil
***REMOVED***

func (m *memLS) Unlock(now time.Time, token string) error ***REMOVED***
	m.mu.Lock()
	defer m.mu.Unlock()
	m.collectExpiredNodes(now)

	n := m.byToken[token]
	if n == nil ***REMOVED***
		return ErrNoSuchLock
	***REMOVED***
	if n.held ***REMOVED***
		return ErrLocked
	***REMOVED***
	m.remove(n)
	return nil
***REMOVED***

func (m *memLS) canCreate(name string, zeroDepth bool) bool ***REMOVED***
	return walkToRoot(name, func(name0 string, first bool) bool ***REMOVED***
		n := m.byName[name0]
		if n == nil ***REMOVED***
			return true
		***REMOVED***
		if first ***REMOVED***
			if n.token != "" ***REMOVED***
				// The target node is already locked.
				return false
			***REMOVED***
			if !zeroDepth ***REMOVED***
				// The requested lock depth is infinite, and the fact that n exists
				// (n != nil) means that a descendent of the target node is locked.
				return false
			***REMOVED***
		***REMOVED*** else if n.token != "" && !n.details.ZeroDepth ***REMOVED***
			// An ancestor of the target node is locked with infinite depth.
			return false
		***REMOVED***
		return true
	***REMOVED***)
***REMOVED***

func (m *memLS) create(name string) (ret *memLSNode) ***REMOVED***
	walkToRoot(name, func(name0 string, first bool) bool ***REMOVED***
		n := m.byName[name0]
		if n == nil ***REMOVED***
			n = &memLSNode***REMOVED***
				details: LockDetails***REMOVED***
					Root: name0,
				***REMOVED***,
				byExpiryIndex: -1,
			***REMOVED***
			m.byName[name0] = n
		***REMOVED***
		n.refCount++
		if first ***REMOVED***
			ret = n
		***REMOVED***
		return true
	***REMOVED***)
	return ret
***REMOVED***

func (m *memLS) remove(n *memLSNode) ***REMOVED***
	delete(m.byToken, n.token)
	n.token = ""
	walkToRoot(n.details.Root, func(name0 string, first bool) bool ***REMOVED***
		x := m.byName[name0]
		x.refCount--
		if x.refCount == 0 ***REMOVED***
			delete(m.byName, name0)
		***REMOVED***
		return true
	***REMOVED***)
	if n.byExpiryIndex >= 0 ***REMOVED***
		heap.Remove(&m.byExpiry, n.byExpiryIndex)
	***REMOVED***
***REMOVED***

func walkToRoot(name string, f func(name0 string, first bool) bool) bool ***REMOVED***
	for first := true; ; first = false ***REMOVED***
		if !f(name, first) ***REMOVED***
			return false
		***REMOVED***
		if name == "/" ***REMOVED***
			break
		***REMOVED***
		name = name[:strings.LastIndex(name, "/")]
		if name == "" ***REMOVED***
			name = "/"
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

type memLSNode struct ***REMOVED***
	// details are the lock metadata. Even if this node's name is not explicitly locked,
	// details.Root will still equal the node's name.
	details LockDetails
	// token is the unique identifier for this node's lock. An empty token means that
	// this node is not explicitly locked.
	token string
	// refCount is the number of self-or-descendent nodes that are explicitly locked.
	refCount int
	// expiry is when this node's lock expires.
	expiry time.Time
	// byExpiryIndex is the index of this node in memLS.byExpiry. It is -1
	// if this node does not expire, or has expired.
	byExpiryIndex int
	// held is whether this node's lock is actively held by a Confirm call.
	held bool
***REMOVED***

type byExpiry []*memLSNode

func (b *byExpiry) Len() int ***REMOVED***
	return len(*b)
***REMOVED***

func (b *byExpiry) Less(i, j int) bool ***REMOVED***
	return (*b)[i].expiry.Before((*b)[j].expiry)
***REMOVED***

func (b *byExpiry) Swap(i, j int) ***REMOVED***
	(*b)[i], (*b)[j] = (*b)[j], (*b)[i]
	(*b)[i].byExpiryIndex = i
	(*b)[j].byExpiryIndex = j
***REMOVED***

func (b *byExpiry) Push(x interface***REMOVED******REMOVED***) ***REMOVED***
	n := x.(*memLSNode)
	n.byExpiryIndex = len(*b)
	*b = append(*b, n)
***REMOVED***

func (b *byExpiry) Pop() interface***REMOVED******REMOVED*** ***REMOVED***
	i := len(*b) - 1
	n := (*b)[i]
	(*b)[i] = nil
	n.byExpiryIndex = -1
	*b = (*b)[:i]
	return n
***REMOVED***

const infiniteTimeout = -1

// parseTimeout parses the Timeout HTTP header, as per section 10.7. If s is
// empty, an infiniteTimeout is returned.
func parseTimeout(s string) (time.Duration, error) ***REMOVED***
	if s == "" ***REMOVED***
		return infiniteTimeout, nil
	***REMOVED***
	if i := strings.IndexByte(s, ','); i >= 0 ***REMOVED***
		s = s[:i]
	***REMOVED***
	s = strings.TrimSpace(s)
	if s == "Infinite" ***REMOVED***
		return infiniteTimeout, nil
	***REMOVED***
	const pre = "Second-"
	if !strings.HasPrefix(s, pre) ***REMOVED***
		return 0, errInvalidTimeout
	***REMOVED***
	s = s[len(pre):]
	if s == "" || s[0] < '0' || '9' < s[0] ***REMOVED***
		return 0, errInvalidTimeout
	***REMOVED***
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil || 1<<32-1 < n ***REMOVED***
		return 0, errInvalidTimeout
	***REMOVED***
	return time.Duration(n) * time.Second, nil
***REMOVED***
