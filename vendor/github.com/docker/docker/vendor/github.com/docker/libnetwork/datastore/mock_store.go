package datastore

import (
	"errors"

	"github.com/docker/libkv/store"
	"github.com/docker/libnetwork/types"
)

var (
	// ErrNotImplmented exported
	ErrNotImplmented = errors.New("Functionality not implemented")
)

// MockData exported
type MockData struct ***REMOVED***
	Data  []byte
	Index uint64
***REMOVED***

// MockStore exported
type MockStore struct ***REMOVED***
	db map[string]*MockData
***REMOVED***

// NewMockStore creates a Map backed Datastore that is useful for mocking
func NewMockStore() *MockStore ***REMOVED***
	db := make(map[string]*MockData)
	return &MockStore***REMOVED***db***REMOVED***
***REMOVED***

// Get the value at "key", returns the last modified index
// to use in conjunction to CAS calls
func (s *MockStore) Get(key string) (*store.KVPair, error) ***REMOVED***
	mData := s.db[key]
	if mData == nil ***REMOVED***
		return nil, nil
	***REMOVED***
	return &store.KVPair***REMOVED***Value: mData.Data, LastIndex: mData.Index***REMOVED***, nil

***REMOVED***

// Put a value at "key"
func (s *MockStore) Put(key string, value []byte, options *store.WriteOptions) error ***REMOVED***
	mData := s.db[key]
	if mData == nil ***REMOVED***
		mData = &MockData***REMOVED***value, 0***REMOVED***
	***REMOVED***
	mData.Index = mData.Index + 1
	s.db[key] = mData
	return nil
***REMOVED***

// Delete a value at "key"
func (s *MockStore) Delete(key string) error ***REMOVED***
	delete(s.db, key)
	return nil
***REMOVED***

// Exists checks that the key exists inside the store
func (s *MockStore) Exists(key string) (bool, error) ***REMOVED***
	_, ok := s.db[key]
	return ok, nil
***REMOVED***

// List gets a range of values at "directory"
func (s *MockStore) List(prefix string) ([]*store.KVPair, error) ***REMOVED***
	return nil, ErrNotImplmented
***REMOVED***

// DeleteTree deletes a range of values at "directory"
func (s *MockStore) DeleteTree(prefix string) error ***REMOVED***
	delete(s.db, prefix)
	return nil
***REMOVED***

// Watch a single key for modifications
func (s *MockStore) Watch(key string, stopCh <-chan struct***REMOVED******REMOVED***) (<-chan *store.KVPair, error) ***REMOVED***
	return nil, ErrNotImplmented
***REMOVED***

// WatchTree triggers a watch on a range of values at "directory"
func (s *MockStore) WatchTree(prefix string, stopCh <-chan struct***REMOVED******REMOVED***) (<-chan []*store.KVPair, error) ***REMOVED***
	return nil, ErrNotImplmented
***REMOVED***

// NewLock exposed
func (s *MockStore) NewLock(key string, options *store.LockOptions) (store.Locker, error) ***REMOVED***
	return nil, ErrNotImplmented
***REMOVED***

// AtomicPut put a value at "key" if the key has not been
// modified in the meantime, throws an error if this is the case
func (s *MockStore) AtomicPut(key string, newValue []byte, previous *store.KVPair, options *store.WriteOptions) (bool, *store.KVPair, error) ***REMOVED***
	mData := s.db[key]

	if previous == nil ***REMOVED***
		if mData != nil ***REMOVED***
			return false, nil, types.BadRequestErrorf("atomic put failed because key exists")
		***REMOVED*** // Else OK.
	***REMOVED*** else ***REMOVED***
		if mData == nil ***REMOVED***
			return false, nil, types.BadRequestErrorf("atomic put failed because key exists")
		***REMOVED***
		if mData != nil && mData.Index != previous.LastIndex ***REMOVED***
			return false, nil, types.BadRequestErrorf("atomic put failed due to mismatched Index")
		***REMOVED*** // Else OK.
	***REMOVED***
	err := s.Put(key, newValue, nil)
	if err != nil ***REMOVED***
		return false, nil, err
	***REMOVED***
	return true, &store.KVPair***REMOVED***Key: key, Value: newValue, LastIndex: s.db[key].Index***REMOVED***, nil
***REMOVED***

// AtomicDelete deletes a value at "key" if the key has not
// been modified in the meantime, throws an error if this is the case
func (s *MockStore) AtomicDelete(key string, previous *store.KVPair) (bool, error) ***REMOVED***
	mData := s.db[key]
	if mData != nil && mData.Index != previous.LastIndex ***REMOVED***
		return false, types.BadRequestErrorf("atomic delete failed due to mismatched Index")
	***REMOVED***
	return true, s.Delete(key)
***REMOVED***

// Close closes the client connection
func (s *MockStore) Close() ***REMOVED***
	return
***REMOVED***
