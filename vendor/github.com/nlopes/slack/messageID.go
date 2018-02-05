package slack

import "sync"

// IDGenerator provides an interface for generating integer ID values.
type IDGenerator interface ***REMOVED***
	Next() int
***REMOVED***

// NewSafeID returns a new instance of an IDGenerator which is safe for
// concurrent use by multiple goroutines.
func NewSafeID(startID int) IDGenerator ***REMOVED***
	return &safeID***REMOVED***
		nextID: startID,
		mutex:  &sync.Mutex***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

type safeID struct ***REMOVED***
	nextID int
	mutex  *sync.Mutex
***REMOVED***

func (s *safeID) Next() int ***REMOVED***
	s.mutex.Lock()
	defer s.mutex.Unlock()
	id := s.nextID
	s.nextID++
	return id
***REMOVED***
