package libnetwork

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/docker/libnetwork/datastore"
)

type endpointCnt struct ***REMOVED***
	n        *network
	Count    uint64
	dbIndex  uint64
	dbExists bool
	sync.Mutex
***REMOVED***

const epCntKeyPrefix = "endpoint_count"

func (ec *endpointCnt) Key() []string ***REMOVED***
	ec.Lock()
	defer ec.Unlock()

	return []string***REMOVED***epCntKeyPrefix, ec.n.id***REMOVED***
***REMOVED***

func (ec *endpointCnt) KeyPrefix() []string ***REMOVED***
	ec.Lock()
	defer ec.Unlock()

	return []string***REMOVED***epCntKeyPrefix, ec.n.id***REMOVED***
***REMOVED***

func (ec *endpointCnt) Value() []byte ***REMOVED***
	ec.Lock()
	defer ec.Unlock()

	b, err := json.Marshal(ec)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	return b
***REMOVED***

func (ec *endpointCnt) SetValue(value []byte) error ***REMOVED***
	ec.Lock()
	defer ec.Unlock()

	return json.Unmarshal(value, &ec)
***REMOVED***

func (ec *endpointCnt) Index() uint64 ***REMOVED***
	ec.Lock()
	defer ec.Unlock()
	return ec.dbIndex
***REMOVED***

func (ec *endpointCnt) SetIndex(index uint64) ***REMOVED***
	ec.Lock()
	ec.dbIndex = index
	ec.dbExists = true
	ec.Unlock()
***REMOVED***

func (ec *endpointCnt) Exists() bool ***REMOVED***
	ec.Lock()
	defer ec.Unlock()
	return ec.dbExists
***REMOVED***

func (ec *endpointCnt) Skip() bool ***REMOVED***
	ec.Lock()
	defer ec.Unlock()
	return !ec.n.persist
***REMOVED***

func (ec *endpointCnt) New() datastore.KVObject ***REMOVED***
	ec.Lock()
	defer ec.Unlock()

	return &endpointCnt***REMOVED***
		n: ec.n,
	***REMOVED***
***REMOVED***

func (ec *endpointCnt) CopyTo(o datastore.KVObject) error ***REMOVED***
	ec.Lock()
	defer ec.Unlock()

	dstEc := o.(*endpointCnt)
	dstEc.n = ec.n
	dstEc.Count = ec.Count
	dstEc.dbExists = ec.dbExists
	dstEc.dbIndex = ec.dbIndex

	return nil
***REMOVED***

func (ec *endpointCnt) DataScope() string ***REMOVED***
	return ec.n.DataScope()
***REMOVED***

func (ec *endpointCnt) EndpointCnt() uint64 ***REMOVED***
	ec.Lock()
	defer ec.Unlock()

	return ec.Count
***REMOVED***

func (ec *endpointCnt) updateStore() error ***REMOVED***
	store := ec.n.getController().getStore(ec.DataScope())
	if store == nil ***REMOVED***
		return fmt.Errorf("store not found for scope %s on endpoint count update", ec.DataScope())
	***REMOVED***
	// make a copy of count and n to avoid being overwritten by store.GetObject
	count := ec.EndpointCnt()
	n := ec.n
	for ***REMOVED***
		if err := ec.n.getController().updateToStore(ec); err == nil || err != datastore.ErrKeyModified ***REMOVED***
			return err
		***REMOVED***
		if err := store.GetObject(datastore.Key(ec.Key()...), ec); err != nil ***REMOVED***
			return fmt.Errorf("could not update the kvobject to latest on endpoint count update: %v", err)
		***REMOVED***
		ec.Lock()
		ec.Count = count
		ec.n = n
		ec.Unlock()
	***REMOVED***
***REMOVED***

func (ec *endpointCnt) setCnt(cnt uint64) error ***REMOVED***
	ec.Lock()
	ec.Count = cnt
	ec.Unlock()
	return ec.updateStore()
***REMOVED***

func (ec *endpointCnt) atomicIncDecEpCnt(inc bool) error ***REMOVED***
	store := ec.n.getController().getStore(ec.DataScope())
	if store == nil ***REMOVED***
		return fmt.Errorf("store not found for scope %s", ec.DataScope())
	***REMOVED***

	tmp := &endpointCnt***REMOVED***n: ec.n***REMOVED***
	if err := store.GetObject(datastore.Key(ec.Key()...), tmp); err != nil ***REMOVED***
		return err
	***REMOVED***
retry:
	ec.Lock()
	if inc ***REMOVED***
		ec.Count++
	***REMOVED*** else ***REMOVED***
		if ec.Count > 0 ***REMOVED***
			ec.Count--
		***REMOVED***
	***REMOVED***
	ec.Unlock()

	if err := ec.n.getController().updateToStore(ec); err != nil ***REMOVED***
		if err == datastore.ErrKeyModified ***REMOVED***
			if err := store.GetObject(datastore.Key(ec.Key()...), ec); err != nil ***REMOVED***
				return fmt.Errorf("could not update the kvobject to latest when trying to atomic add endpoint count: %v", err)
			***REMOVED***

			goto retry
		***REMOVED***

		return err
	***REMOVED***

	return nil
***REMOVED***

func (ec *endpointCnt) IncEndpointCnt() error ***REMOVED***
	return ec.atomicIncDecEpCnt(true)
***REMOVED***

func (ec *endpointCnt) DecEndpointCnt() error ***REMOVED***
	return ec.atomicIncDecEpCnt(false)
***REMOVED***
