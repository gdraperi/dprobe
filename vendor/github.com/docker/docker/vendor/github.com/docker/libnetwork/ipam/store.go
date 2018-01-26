package ipam

import (
	"encoding/json"

	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

// Key provides the Key to be used in KV Store
func (aSpace *addrSpace) Key() []string ***REMOVED***
	aSpace.Lock()
	defer aSpace.Unlock()
	return []string***REMOVED***aSpace.id***REMOVED***
***REMOVED***

// KeyPrefix returns the immediate parent key that can be used for tree walk
func (aSpace *addrSpace) KeyPrefix() []string ***REMOVED***
	aSpace.Lock()
	defer aSpace.Unlock()
	return []string***REMOVED***dsConfigKey***REMOVED***
***REMOVED***

// Value marshals the data to be stored in the KV store
func (aSpace *addrSpace) Value() []byte ***REMOVED***
	b, err := json.Marshal(aSpace)
	if err != nil ***REMOVED***
		logrus.Warnf("Failed to marshal ipam configured pools: %v", err)
		return nil
	***REMOVED***
	return b
***REMOVED***

// SetValue unmarshalls the data from the KV store.
func (aSpace *addrSpace) SetValue(value []byte) error ***REMOVED***
	rc := &addrSpace***REMOVED***subnets: make(map[SubnetKey]*PoolData)***REMOVED***
	if err := json.Unmarshal(value, rc); err != nil ***REMOVED***
		return err
	***REMOVED***
	aSpace.subnets = rc.subnets
	return nil
***REMOVED***

// Index returns the latest DB Index as seen by this object
func (aSpace *addrSpace) Index() uint64 ***REMOVED***
	aSpace.Lock()
	defer aSpace.Unlock()
	return aSpace.dbIndex
***REMOVED***

// SetIndex method allows the datastore to store the latest DB Index into this object
func (aSpace *addrSpace) SetIndex(index uint64) ***REMOVED***
	aSpace.Lock()
	aSpace.dbIndex = index
	aSpace.dbExists = true
	aSpace.Unlock()
***REMOVED***

// Exists method is true if this object has been stored in the DB.
func (aSpace *addrSpace) Exists() bool ***REMOVED***
	aSpace.Lock()
	defer aSpace.Unlock()
	return aSpace.dbExists
***REMOVED***

// Skip provides a way for a KV Object to avoid persisting it in the KV Store
func (aSpace *addrSpace) Skip() bool ***REMOVED***
	return false
***REMOVED***

func (a *Allocator) getStore(as string) datastore.DataStore ***REMOVED***
	a.Lock()
	defer a.Unlock()

	if aSpace, ok := a.addrSpaces[as]; ok ***REMOVED***
		return aSpace.ds
	***REMOVED***

	return nil
***REMOVED***

func (a *Allocator) getAddressSpaceFromStore(as string) (*addrSpace, error) ***REMOVED***
	store := a.getStore(as)

	// IPAM may not have a valid store. In such cases it is just in-memory state.
	if store == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	pc := &addrSpace***REMOVED***id: dsConfigKey + "/" + as, ds: store, alloc: a***REMOVED***
	if err := store.GetObject(datastore.Key(pc.Key()...), pc); err != nil ***REMOVED***
		if err == datastore.ErrKeyNotFound ***REMOVED***
			return nil, nil
		***REMOVED***

		return nil, types.InternalErrorf("could not get pools config from store: %v", err)
	***REMOVED***

	return pc, nil
***REMOVED***

func (a *Allocator) writeToStore(aSpace *addrSpace) error ***REMOVED***
	store := aSpace.store()

	// IPAM may not have a valid store. In such cases it is just in-memory state.
	if store == nil ***REMOVED***
		return nil
	***REMOVED***

	err := store.PutObjectAtomic(aSpace)
	if err == datastore.ErrKeyModified ***REMOVED***
		return types.RetryErrorf("failed to perform atomic write (%v). retry might fix the error", err)
	***REMOVED***

	return err
***REMOVED***

func (a *Allocator) deleteFromStore(aSpace *addrSpace) error ***REMOVED***
	store := aSpace.store()

	// IPAM may not have a valid store. In such cases it is just in-memory state.
	if store == nil ***REMOVED***
		return nil
	***REMOVED***

	return store.DeleteObjectAtomic(aSpace)
***REMOVED***

// DataScope method returns the storage scope of the datastore
func (aSpace *addrSpace) DataScope() string ***REMOVED***
	aSpace.Lock()
	defer aSpace.Unlock()

	return aSpace.scope
***REMOVED***
