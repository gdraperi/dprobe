package bitseq

import (
	"encoding/json"
	"fmt"

	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/types"
)

// Key provides the Key to be used in KV Store
func (h *Handle) Key() []string ***REMOVED***
	h.Lock()
	defer h.Unlock()
	return []string***REMOVED***h.app, h.id***REMOVED***
***REMOVED***

// KeyPrefix returns the immediate parent key that can be used for tree walk
func (h *Handle) KeyPrefix() []string ***REMOVED***
	h.Lock()
	defer h.Unlock()
	return []string***REMOVED***h.app***REMOVED***
***REMOVED***

// Value marshals the data to be stored in the KV store
func (h *Handle) Value() []byte ***REMOVED***
	b, err := json.Marshal(h)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	return b
***REMOVED***

// SetValue unmarshals the data from the KV store
func (h *Handle) SetValue(value []byte) error ***REMOVED***
	return json.Unmarshal(value, h)
***REMOVED***

// Index returns the latest DB Index as seen by this object
func (h *Handle) Index() uint64 ***REMOVED***
	h.Lock()
	defer h.Unlock()
	return h.dbIndex
***REMOVED***

// SetIndex method allows the datastore to store the latest DB Index into this object
func (h *Handle) SetIndex(index uint64) ***REMOVED***
	h.Lock()
	h.dbIndex = index
	h.dbExists = true
	h.Unlock()
***REMOVED***

// Exists method is true if this object has been stored in the DB.
func (h *Handle) Exists() bool ***REMOVED***
	h.Lock()
	defer h.Unlock()
	return h.dbExists
***REMOVED***

// New method returns a handle based on the receiver handle
func (h *Handle) New() datastore.KVObject ***REMOVED***
	h.Lock()
	defer h.Unlock()

	return &Handle***REMOVED***
		app:   h.app,
		store: h.store,
	***REMOVED***
***REMOVED***

// CopyTo deep copies the handle into the passed destination object
func (h *Handle) CopyTo(o datastore.KVObject) error ***REMOVED***
	h.Lock()
	defer h.Unlock()

	dstH := o.(*Handle)
	if h == dstH ***REMOVED***
		return nil
	***REMOVED***
	dstH.Lock()
	dstH.bits = h.bits
	dstH.unselected = h.unselected
	dstH.head = h.head.getCopy()
	dstH.app = h.app
	dstH.id = h.id
	dstH.dbIndex = h.dbIndex
	dstH.dbExists = h.dbExists
	dstH.store = h.store
	dstH.curr = h.curr
	dstH.Unlock()

	return nil
***REMOVED***

// Skip provides a way for a KV Object to avoid persisting it in the KV Store
func (h *Handle) Skip() bool ***REMOVED***
	return false
***REMOVED***

// DataScope method returns the storage scope of the datastore
func (h *Handle) DataScope() string ***REMOVED***
	h.Lock()
	defer h.Unlock()

	return h.store.Scope()
***REMOVED***

func (h *Handle) fromDsValue(value []byte) error ***REMOVED***
	var ba []byte
	if err := json.Unmarshal(value, &ba); err != nil ***REMOVED***
		return fmt.Errorf("failed to decode json: %s", err.Error())
	***REMOVED***
	if err := h.FromByteArray(ba); err != nil ***REMOVED***
		return fmt.Errorf("failed to decode handle: %s", err.Error())
	***REMOVED***
	return nil
***REMOVED***

func (h *Handle) writeToStore() error ***REMOVED***
	h.Lock()
	store := h.store
	h.Unlock()
	if store == nil ***REMOVED***
		return nil
	***REMOVED***
	err := store.PutObjectAtomic(h)
	if err == datastore.ErrKeyModified ***REMOVED***
		return types.RetryErrorf("failed to perform atomic write (%v). Retry might fix the error", err)
	***REMOVED***
	return err
***REMOVED***

func (h *Handle) deleteFromStore() error ***REMOVED***
	h.Lock()
	store := h.store
	h.Unlock()
	if store == nil ***REMOVED***
		return nil
	***REMOVED***
	return store.DeleteObjectAtomic(h)
***REMOVED***
