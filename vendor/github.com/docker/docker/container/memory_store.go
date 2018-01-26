package container

import (
	"sync"
)

// memoryStore implements a Store in memory.
type memoryStore struct ***REMOVED***
	s map[string]*Container
	sync.RWMutex
***REMOVED***

// NewMemoryStore initializes a new memory store.
func NewMemoryStore() Store ***REMOVED***
	return &memoryStore***REMOVED***
		s: make(map[string]*Container),
	***REMOVED***
***REMOVED***

// Add appends a new container to the memory store.
// It overrides the id if it existed before.
func (c *memoryStore) Add(id string, cont *Container) ***REMOVED***
	c.Lock()
	c.s[id] = cont
	c.Unlock()
***REMOVED***

// Get returns a container from the store by id.
func (c *memoryStore) Get(id string) *Container ***REMOVED***
	var res *Container
	c.RLock()
	res = c.s[id]
	c.RUnlock()
	return res
***REMOVED***

// Delete removes a container from the store by id.
func (c *memoryStore) Delete(id string) ***REMOVED***
	c.Lock()
	delete(c.s, id)
	c.Unlock()
***REMOVED***

// List returns a sorted list of containers from the store.
// The containers are ordered by creation date.
func (c *memoryStore) List() []*Container ***REMOVED***
	containers := History(c.all())
	containers.sort()
	return containers
***REMOVED***

// Size returns the number of containers in the store.
func (c *memoryStore) Size() int ***REMOVED***
	c.RLock()
	defer c.RUnlock()
	return len(c.s)
***REMOVED***

// First returns the first container found in the store by a given filter.
func (c *memoryStore) First(filter StoreFilter) *Container ***REMOVED***
	for _, cont := range c.all() ***REMOVED***
		if filter(cont) ***REMOVED***
			return cont
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// ApplyAll calls the reducer function with every container in the store.
// This operation is asynchronous in the memory store.
// NOTE: Modifications to the store MUST NOT be done by the StoreReducer.
func (c *memoryStore) ApplyAll(apply StoreReducer) ***REMOVED***
	wg := new(sync.WaitGroup)
	for _, cont := range c.all() ***REMOVED***
		wg.Add(1)
		go func(container *Container) ***REMOVED***
			apply(container)
			wg.Done()
		***REMOVED***(cont)
	***REMOVED***

	wg.Wait()
***REMOVED***

func (c *memoryStore) all() []*Container ***REMOVED***
	c.RLock()
	containers := make([]*Container, 0, len(c.s))
	for _, cont := range c.s ***REMOVED***
		containers = append(containers, cont)
	***REMOVED***
	c.RUnlock()
	return containers
***REMOVED***

var _ Store = &memoryStore***REMOVED******REMOVED***
