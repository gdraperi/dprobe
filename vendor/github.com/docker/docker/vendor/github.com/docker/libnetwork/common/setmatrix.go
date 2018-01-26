package common

import (
	"sync"

	mapset "github.com/deckarep/golang-set"
)

// SetMatrix is a map of Sets
type SetMatrix interface ***REMOVED***
	// Get returns the members of the set for a specific key as a slice.
	Get(key string) ([]interface***REMOVED******REMOVED***, bool)
	// Contains is used to verify if an element is in a set for a specific key
	// returns true if the element is in the set
	// returns true if there is a set for the key
	Contains(key string, value interface***REMOVED******REMOVED***) (bool, bool)
	// Insert inserts the value in the set of a key
	// returns true if the value is inserted (was not already in the set), false otherwise
	// returns also the length of the set for the key
	Insert(key string, value interface***REMOVED******REMOVED***) (bool, int)
	// Remove removes the value in the set for a specific key
	// returns true if the value is deleted, false otherwise
	// returns also the length of the set for the key
	Remove(key string, value interface***REMOVED******REMOVED***) (bool, int)
	// Cardinality returns the number of elements in the set for a key
	// returns false if the set is not present
	Cardinality(key string) (int, bool)
	// String returns the string version of the set, empty otherwise
	// returns false if the set is not present
	String(key string) (string, bool)
	// Returns all the keys in the map
	Keys() []string
***REMOVED***

type setMatrix struct ***REMOVED***
	matrix map[string]mapset.Set

	sync.Mutex
***REMOVED***

// NewSetMatrix creates a new set matrix object
func NewSetMatrix() SetMatrix ***REMOVED***
	s := &setMatrix***REMOVED******REMOVED***
	s.init()
	return s
***REMOVED***

func (s *setMatrix) init() ***REMOVED***
	s.matrix = make(map[string]mapset.Set)
***REMOVED***

func (s *setMatrix) Get(key string) ([]interface***REMOVED******REMOVED***, bool) ***REMOVED***
	s.Lock()
	defer s.Unlock()
	set, ok := s.matrix[key]
	if !ok ***REMOVED***
		return nil, ok
	***REMOVED***
	return set.ToSlice(), ok
***REMOVED***

func (s *setMatrix) Contains(key string, value interface***REMOVED******REMOVED***) (bool, bool) ***REMOVED***
	s.Lock()
	defer s.Unlock()
	set, ok := s.matrix[key]
	if !ok ***REMOVED***
		return false, ok
	***REMOVED***
	return set.Contains(value), ok
***REMOVED***

func (s *setMatrix) Insert(key string, value interface***REMOVED******REMOVED***) (bool, int) ***REMOVED***
	s.Lock()
	defer s.Unlock()
	set, ok := s.matrix[key]
	if !ok ***REMOVED***
		s.matrix[key] = mapset.NewSet()
		s.matrix[key].Add(value)
		return true, 1
	***REMOVED***

	return set.Add(value), set.Cardinality()
***REMOVED***

func (s *setMatrix) Remove(key string, value interface***REMOVED******REMOVED***) (bool, int) ***REMOVED***
	s.Lock()
	defer s.Unlock()
	set, ok := s.matrix[key]
	if !ok ***REMOVED***
		return false, 0
	***REMOVED***

	var removed bool
	if set.Contains(value) ***REMOVED***
		set.Remove(value)
		removed = true
		// If the set is empty remove it from the matrix
		if set.Cardinality() == 0 ***REMOVED***
			delete(s.matrix, key)
		***REMOVED***
	***REMOVED***

	return removed, set.Cardinality()
***REMOVED***

func (s *setMatrix) Cardinality(key string) (int, bool) ***REMOVED***
	s.Lock()
	defer s.Unlock()
	set, ok := s.matrix[key]
	if !ok ***REMOVED***
		return 0, ok
	***REMOVED***

	return set.Cardinality(), ok
***REMOVED***

func (s *setMatrix) String(key string) (string, bool) ***REMOVED***
	s.Lock()
	defer s.Unlock()
	set, ok := s.matrix[key]
	if !ok ***REMOVED***
		return "", ok
	***REMOVED***
	return set.String(), ok
***REMOVED***

func (s *setMatrix) Keys() []string ***REMOVED***
	s.Lock()
	defer s.Unlock()
	keys := make([]string, 0, len(s.matrix))
	for k := range s.matrix ***REMOVED***
		keys = append(keys, k)
	***REMOVED***
	return keys
***REMOVED***
