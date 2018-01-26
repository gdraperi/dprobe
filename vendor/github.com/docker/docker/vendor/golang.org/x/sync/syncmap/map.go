// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package syncmap provides a concurrent map implementation.
// It is a prototype for a proposed addition to the sync package
// in the standard library.
// (https://golang.org/issue/18177)
package syncmap

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// Map is a concurrent map with amortized-constant-time loads, stores, and deletes.
// It is safe for multiple goroutines to call a Map's methods concurrently.
//
// The zero Map is valid and empty.
//
// A Map must not be copied after first use.
type Map struct ***REMOVED***
	mu sync.Mutex

	// read contains the portion of the map's contents that are safe for
	// concurrent access (with or without mu held).
	//
	// The read field itself is always safe to load, but must only be stored with
	// mu held.
	//
	// Entries stored in read may be updated concurrently without mu, but updating
	// a previously-expunged entry requires that the entry be copied to the dirty
	// map and unexpunged with mu held.
	read atomic.Value // readOnly

	// dirty contains the portion of the map's contents that require mu to be
	// held. To ensure that the dirty map can be promoted to the read map quickly,
	// it also includes all of the non-expunged entries in the read map.
	//
	// Expunged entries are not stored in the dirty map. An expunged entry in the
	// clean map must be unexpunged and added to the dirty map before a new value
	// can be stored to it.
	//
	// If the dirty map is nil, the next write to the map will initialize it by
	// making a shallow copy of the clean map, omitting stale entries.
	dirty map[interface***REMOVED******REMOVED***]*entry

	// misses counts the number of loads since the read map was last updated that
	// needed to lock mu to determine whether the key was present.
	//
	// Once enough misses have occurred to cover the cost of copying the dirty
	// map, the dirty map will be promoted to the read map (in the unamended
	// state) and the next store to the map will make a new dirty copy.
	misses int
***REMOVED***

// readOnly is an immutable struct stored atomically in the Map.read field.
type readOnly struct ***REMOVED***
	m       map[interface***REMOVED******REMOVED***]*entry
	amended bool // true if the dirty map contains some key not in m.
***REMOVED***

// expunged is an arbitrary pointer that marks entries which have been deleted
// from the dirty map.
var expunged = unsafe.Pointer(new(interface***REMOVED******REMOVED***))

// An entry is a slot in the map corresponding to a particular key.
type entry struct ***REMOVED***
	// p points to the interface***REMOVED******REMOVED*** value stored for the entry.
	//
	// If p == nil, the entry has been deleted and m.dirty == nil.
	//
	// If p == expunged, the entry has been deleted, m.dirty != nil, and the entry
	// is missing from m.dirty.
	//
	// Otherwise, the entry is valid and recorded in m.read.m[key] and, if m.dirty
	// != nil, in m.dirty[key].
	//
	// An entry can be deleted by atomic replacement with nil: when m.dirty is
	// next created, it will atomically replace nil with expunged and leave
	// m.dirty[key] unset.
	//
	// An entry's associated value can be updated by atomic replacement, provided
	// p != expunged. If p == expunged, an entry's associated value can be updated
	// only after first setting m.dirty[key] = e so that lookups using the dirty
	// map find the entry.
	p unsafe.Pointer // *interface***REMOVED******REMOVED***
***REMOVED***

func newEntry(i interface***REMOVED******REMOVED***) *entry ***REMOVED***
	return &entry***REMOVED***p: unsafe.Pointer(&i)***REMOVED***
***REMOVED***

// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (m *Map) Load(key interface***REMOVED******REMOVED***) (value interface***REMOVED******REMOVED***, ok bool) ***REMOVED***
	read, _ := m.read.Load().(readOnly)
	e, ok := read.m[key]
	if !ok && read.amended ***REMOVED***
		m.mu.Lock()
		// Avoid reporting a spurious miss if m.dirty got promoted while we were
		// blocked on m.mu. (If further loads of the same key will not miss, it's
		// not worth copying the dirty map for this key.)
		read, _ = m.read.Load().(readOnly)
		e, ok = read.m[key]
		if !ok && read.amended ***REMOVED***
			e, ok = m.dirty[key]
			// Regardless of whether the entry was present, record a miss: this key
			// will take the slow path until the dirty map is promoted to the read
			// map.
			m.missLocked()
		***REMOVED***
		m.mu.Unlock()
	***REMOVED***
	if !ok ***REMOVED***
		return nil, false
	***REMOVED***
	return e.load()
***REMOVED***

func (e *entry) load() (value interface***REMOVED******REMOVED***, ok bool) ***REMOVED***
	p := atomic.LoadPointer(&e.p)
	if p == nil || p == expunged ***REMOVED***
		return nil, false
	***REMOVED***
	return *(*interface***REMOVED******REMOVED***)(p), true
***REMOVED***

// Store sets the value for a key.
func (m *Map) Store(key, value interface***REMOVED******REMOVED***) ***REMOVED***
	read, _ := m.read.Load().(readOnly)
	if e, ok := read.m[key]; ok && e.tryStore(&value) ***REMOVED***
		return
	***REMOVED***

	m.mu.Lock()
	read, _ = m.read.Load().(readOnly)
	if e, ok := read.m[key]; ok ***REMOVED***
		if e.unexpungeLocked() ***REMOVED***
			// The entry was previously expunged, which implies that there is a
			// non-nil dirty map and this entry is not in it.
			m.dirty[key] = e
		***REMOVED***
		e.storeLocked(&value)
	***REMOVED*** else if e, ok := m.dirty[key]; ok ***REMOVED***
		e.storeLocked(&value)
	***REMOVED*** else ***REMOVED***
		if !read.amended ***REMOVED***
			// We're adding the first new key to the dirty map.
			// Make sure it is allocated and mark the read-only map as incomplete.
			m.dirtyLocked()
			m.read.Store(readOnly***REMOVED***m: read.m, amended: true***REMOVED***)
		***REMOVED***
		m.dirty[key] = newEntry(value)
	***REMOVED***
	m.mu.Unlock()
***REMOVED***

// tryStore stores a value if the entry has not been expunged.
//
// If the entry is expunged, tryStore returns false and leaves the entry
// unchanged.
func (e *entry) tryStore(i *interface***REMOVED******REMOVED***) bool ***REMOVED***
	p := atomic.LoadPointer(&e.p)
	if p == expunged ***REMOVED***
		return false
	***REMOVED***
	for ***REMOVED***
		if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(i)) ***REMOVED***
			return true
		***REMOVED***
		p = atomic.LoadPointer(&e.p)
		if p == expunged ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
***REMOVED***

// unexpungeLocked ensures that the entry is not marked as expunged.
//
// If the entry was previously expunged, it must be added to the dirty map
// before m.mu is unlocked.
func (e *entry) unexpungeLocked() (wasExpunged bool) ***REMOVED***
	return atomic.CompareAndSwapPointer(&e.p, expunged, nil)
***REMOVED***

// storeLocked unconditionally stores a value to the entry.
//
// The entry must be known not to be expunged.
func (e *entry) storeLocked(i *interface***REMOVED******REMOVED***) ***REMOVED***
	atomic.StorePointer(&e.p, unsafe.Pointer(i))
***REMOVED***

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *Map) LoadOrStore(key, value interface***REMOVED******REMOVED***) (actual interface***REMOVED******REMOVED***, loaded bool) ***REMOVED***
	// Avoid locking if it's a clean hit.
	read, _ := m.read.Load().(readOnly)
	if e, ok := read.m[key]; ok ***REMOVED***
		actual, loaded, ok := e.tryLoadOrStore(value)
		if ok ***REMOVED***
			return actual, loaded
		***REMOVED***
	***REMOVED***

	m.mu.Lock()
	read, _ = m.read.Load().(readOnly)
	if e, ok := read.m[key]; ok ***REMOVED***
		if e.unexpungeLocked() ***REMOVED***
			m.dirty[key] = e
		***REMOVED***
		actual, loaded, _ = e.tryLoadOrStore(value)
	***REMOVED*** else if e, ok := m.dirty[key]; ok ***REMOVED***
		actual, loaded, _ = e.tryLoadOrStore(value)
		m.missLocked()
	***REMOVED*** else ***REMOVED***
		if !read.amended ***REMOVED***
			// We're adding the first new key to the dirty map.
			// Make sure it is allocated and mark the read-only map as incomplete.
			m.dirtyLocked()
			m.read.Store(readOnly***REMOVED***m: read.m, amended: true***REMOVED***)
		***REMOVED***
		m.dirty[key] = newEntry(value)
		actual, loaded = value, false
	***REMOVED***
	m.mu.Unlock()

	return actual, loaded
***REMOVED***

// tryLoadOrStore atomically loads or stores a value if the entry is not
// expunged.
//
// If the entry is expunged, tryLoadOrStore leaves the entry unchanged and
// returns with ok==false.
func (e *entry) tryLoadOrStore(i interface***REMOVED******REMOVED***) (actual interface***REMOVED******REMOVED***, loaded, ok bool) ***REMOVED***
	p := atomic.LoadPointer(&e.p)
	if p == expunged ***REMOVED***
		return nil, false, false
	***REMOVED***
	if p != nil ***REMOVED***
		return *(*interface***REMOVED******REMOVED***)(p), true, true
	***REMOVED***

	// Copy the interface after the first load to make this method more amenable
	// to escape analysis: if we hit the "load" path or the entry is expunged, we
	// shouldn't bother heap-allocating.
	ic := i
	for ***REMOVED***
		if atomic.CompareAndSwapPointer(&e.p, nil, unsafe.Pointer(&ic)) ***REMOVED***
			return i, false, true
		***REMOVED***
		p = atomic.LoadPointer(&e.p)
		if p == expunged ***REMOVED***
			return nil, false, false
		***REMOVED***
		if p != nil ***REMOVED***
			return *(*interface***REMOVED******REMOVED***)(p), true, true
		***REMOVED***
	***REMOVED***
***REMOVED***

// Delete deletes the value for a key.
func (m *Map) Delete(key interface***REMOVED******REMOVED***) ***REMOVED***
	read, _ := m.read.Load().(readOnly)
	e, ok := read.m[key]
	if !ok && read.amended ***REMOVED***
		m.mu.Lock()
		read, _ = m.read.Load().(readOnly)
		e, ok = read.m[key]
		if !ok && read.amended ***REMOVED***
			delete(m.dirty, key)
		***REMOVED***
		m.mu.Unlock()
	***REMOVED***
	if ok ***REMOVED***
		e.delete()
	***REMOVED***
***REMOVED***

func (e *entry) delete() (hadValue bool) ***REMOVED***
	for ***REMOVED***
		p := atomic.LoadPointer(&e.p)
		if p == nil || p == expunged ***REMOVED***
			return false
		***REMOVED***
		if atomic.CompareAndSwapPointer(&e.p, p, nil) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
***REMOVED***

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the Map's
// contents: no key will be visited more than once, but if the value for any key
// is stored or deleted concurrently, Range may reflect any mapping for that key
// from any point during the Range call.
//
// Range may be O(N) with the number of elements in the map even if f returns
// false after a constant number of calls.
func (m *Map) Range(f func(key, value interface***REMOVED******REMOVED***) bool) ***REMOVED***
	// We need to be able to iterate over all of the keys that were already
	// present at the start of the call to Range.
	// If read.amended is false, then read.m satisfies that property without
	// requiring us to hold m.mu for a long time.
	read, _ := m.read.Load().(readOnly)
	if read.amended ***REMOVED***
		// m.dirty contains keys not in read.m. Fortunately, Range is already O(N)
		// (assuming the caller does not break out early), so a call to Range
		// amortizes an entire copy of the map: we can promote the dirty copy
		// immediately!
		m.mu.Lock()
		read, _ = m.read.Load().(readOnly)
		if read.amended ***REMOVED***
			read = readOnly***REMOVED***m: m.dirty***REMOVED***
			m.read.Store(read)
			m.dirty = nil
			m.misses = 0
		***REMOVED***
		m.mu.Unlock()
	***REMOVED***

	for k, e := range read.m ***REMOVED***
		v, ok := e.load()
		if !ok ***REMOVED***
			continue
		***REMOVED***
		if !f(k, v) ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

func (m *Map) missLocked() ***REMOVED***
	m.misses++
	if m.misses < len(m.dirty) ***REMOVED***
		return
	***REMOVED***
	m.read.Store(readOnly***REMOVED***m: m.dirty***REMOVED***)
	m.dirty = nil
	m.misses = 0
***REMOVED***

func (m *Map) dirtyLocked() ***REMOVED***
	if m.dirty != nil ***REMOVED***
		return
	***REMOVED***

	read, _ := m.read.Load().(readOnly)
	m.dirty = make(map[interface***REMOVED******REMOVED***]*entry, len(read.m))
	for k, e := range read.m ***REMOVED***
		if !e.tryExpungeLocked() ***REMOVED***
			m.dirty[k] = e
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *entry) tryExpungeLocked() (isExpunged bool) ***REMOVED***
	p := atomic.LoadPointer(&e.p)
	for p == nil ***REMOVED***
		if atomic.CompareAndSwapPointer(&e.p, nil, expunged) ***REMOVED***
			return true
		***REMOVED***
		p = atomic.LoadPointer(&e.p)
	***REMOVED***
	return p == expunged
***REMOVED***
