package simplelru

import (
	"container/list"
	"errors"
)

// EvictCallback is used to get a callback when a cache entry is evicted
type EvictCallback func(key interface***REMOVED******REMOVED***, value interface***REMOVED******REMOVED***)

// LRU implements a non-thread safe fixed size LRU cache
type LRU struct ***REMOVED***
	size      int
	evictList *list.List
	items     map[interface***REMOVED******REMOVED***]*list.Element
	onEvict   EvictCallback
***REMOVED***

// entry is used to hold a value in the evictList
type entry struct ***REMOVED***
	key   interface***REMOVED******REMOVED***
	value interface***REMOVED******REMOVED***
***REMOVED***

// NewLRU constructs an LRU of the given size
func NewLRU(size int, onEvict EvictCallback) (*LRU, error) ***REMOVED***
	if size <= 0 ***REMOVED***
		return nil, errors.New("Must provide a positive size")
	***REMOVED***
	c := &LRU***REMOVED***
		size:      size,
		evictList: list.New(),
		items:     make(map[interface***REMOVED******REMOVED***]*list.Element),
		onEvict:   onEvict,
	***REMOVED***
	return c, nil
***REMOVED***

// Purge is used to completely clear the cache
func (c *LRU) Purge() ***REMOVED***
	for k, v := range c.items ***REMOVED***
		if c.onEvict != nil ***REMOVED***
			c.onEvict(k, v.Value.(*entry).value)
		***REMOVED***
		delete(c.items, k)
	***REMOVED***
	c.evictList.Init()
***REMOVED***

// Add adds a value to the cache.  Returns true if an eviction occured.
func (c *LRU) Add(key, value interface***REMOVED******REMOVED***) bool ***REMOVED***
	// Check for existing item
	if ent, ok := c.items[key]; ok ***REMOVED***
		c.evictList.MoveToFront(ent)
		ent.Value.(*entry).value = value
		return false
	***REMOVED***

	// Add new item
	ent := &entry***REMOVED***key, value***REMOVED***
	entry := c.evictList.PushFront(ent)
	c.items[key] = entry

	evict := c.evictList.Len() > c.size
	// Verify size not exceeded
	if evict ***REMOVED***
		c.removeOldest()
	***REMOVED***
	return evict
***REMOVED***

// Get looks up a key's value from the cache.
func (c *LRU) Get(key interface***REMOVED******REMOVED***) (value interface***REMOVED******REMOVED***, ok bool) ***REMOVED***
	if ent, ok := c.items[key]; ok ***REMOVED***
		c.evictList.MoveToFront(ent)
		return ent.Value.(*entry).value, true
	***REMOVED***
	return
***REMOVED***

// Check if a key is in the cache, without updating the recent-ness
// or deleting it for being stale.
func (c *LRU) Contains(key interface***REMOVED******REMOVED***) (ok bool) ***REMOVED***
	_, ok = c.items[key]
	return ok
***REMOVED***

// Returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *LRU) Peek(key interface***REMOVED******REMOVED***) (value interface***REMOVED******REMOVED***, ok bool) ***REMOVED***
	if ent, ok := c.items[key]; ok ***REMOVED***
		return ent.Value.(*entry).value, true
	***REMOVED***
	return nil, ok
***REMOVED***

// Remove removes the provided key from the cache, returning if the
// key was contained.
func (c *LRU) Remove(key interface***REMOVED******REMOVED***) bool ***REMOVED***
	if ent, ok := c.items[key]; ok ***REMOVED***
		c.removeElement(ent)
		return true
	***REMOVED***
	return false
***REMOVED***

// RemoveOldest removes the oldest item from the cache.
func (c *LRU) RemoveOldest() (interface***REMOVED******REMOVED***, interface***REMOVED******REMOVED***, bool) ***REMOVED***
	ent := c.evictList.Back()
	if ent != nil ***REMOVED***
		c.removeElement(ent)
		kv := ent.Value.(*entry)
		return kv.key, kv.value, true
	***REMOVED***
	return nil, nil, false
***REMOVED***

// GetOldest returns the oldest entry
func (c *LRU) GetOldest() (interface***REMOVED******REMOVED***, interface***REMOVED******REMOVED***, bool) ***REMOVED***
	ent := c.evictList.Back()
	if ent != nil ***REMOVED***
		kv := ent.Value.(*entry)
		return kv.key, kv.value, true
	***REMOVED***
	return nil, nil, false
***REMOVED***

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (c *LRU) Keys() []interface***REMOVED******REMOVED*** ***REMOVED***
	keys := make([]interface***REMOVED******REMOVED***, len(c.items))
	i := 0
	for ent := c.evictList.Back(); ent != nil; ent = ent.Prev() ***REMOVED***
		keys[i] = ent.Value.(*entry).key
		i++
	***REMOVED***
	return keys
***REMOVED***

// Len returns the number of items in the cache.
func (c *LRU) Len() int ***REMOVED***
	return c.evictList.Len()
***REMOVED***

// removeOldest removes the oldest item from the cache.
func (c *LRU) removeOldest() ***REMOVED***
	ent := c.evictList.Back()
	if ent != nil ***REMOVED***
		c.removeElement(ent)
	***REMOVED***
***REMOVED***

// removeElement is used to remove a given list element from the cache
func (c *LRU) removeElement(e *list.Element) ***REMOVED***
	c.evictList.Remove(e)
	kv := e.Value.(*entry)
	delete(c.items, kv.key)
	if c.onEvict != nil ***REMOVED***
		c.onEvict(kv.key, kv.value)
	***REMOVED***
***REMOVED***
