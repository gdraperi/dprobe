package graphdriver

import "sync"

type minfo struct ***REMOVED***
	check bool
	count int
***REMOVED***

// RefCounter is a generic counter for use by graphdriver Get/Put calls
type RefCounter struct ***REMOVED***
	counts  map[string]*minfo
	mu      sync.Mutex
	checker Checker
***REMOVED***

// NewRefCounter returns a new RefCounter
func NewRefCounter(c Checker) *RefCounter ***REMOVED***
	return &RefCounter***REMOVED***
		checker: c,
		counts:  make(map[string]*minfo),
	***REMOVED***
***REMOVED***

// Increment increases the ref count for the given id and returns the current count
func (c *RefCounter) Increment(path string) int ***REMOVED***
	return c.incdec(path, func(minfo *minfo) ***REMOVED***
		minfo.count++
	***REMOVED***)
***REMOVED***

// Decrement decreases the ref count for the given id and returns the current count
func (c *RefCounter) Decrement(path string) int ***REMOVED***
	return c.incdec(path, func(minfo *minfo) ***REMOVED***
		minfo.count--
	***REMOVED***)
***REMOVED***

func (c *RefCounter) incdec(path string, infoOp func(minfo *minfo)) int ***REMOVED***
	c.mu.Lock()
	m := c.counts[path]
	if m == nil ***REMOVED***
		m = &minfo***REMOVED******REMOVED***
		c.counts[path] = m
	***REMOVED***
	// if we are checking this path for the first time check to make sure
	// if it was already mounted on the system and make sure we have a correct ref
	// count if it is mounted as it is in use.
	if !m.check ***REMOVED***
		m.check = true
		if c.checker.IsMounted(path) ***REMOVED***
			m.count++
		***REMOVED***
	***REMOVED***
	infoOp(m)
	count := m.count
	c.mu.Unlock()
	return count
***REMOVED***
