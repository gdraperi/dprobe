// Package idm manages reservation/release of numerical ids from a configured set of contiguous ids
package idm

import (
	"errors"
	"fmt"

	"github.com/docker/libnetwork/bitseq"
	"github.com/docker/libnetwork/datastore"
)

// Idm manages the reservation/release of numerical ids from a contiguous set
type Idm struct ***REMOVED***
	start  uint64
	end    uint64
	handle *bitseq.Handle
***REMOVED***

// New returns an instance of id manager for a [start,end] set of numerical ids
func New(ds datastore.DataStore, id string, start, end uint64) (*Idm, error) ***REMOVED***
	if id == "" ***REMOVED***
		return nil, errors.New("Invalid id")
	***REMOVED***
	if end <= start ***REMOVED***
		return nil, fmt.Errorf("Invalid set range: [%d, %d]", start, end)
	***REMOVED***

	h, err := bitseq.NewHandle("idm", ds, id, 1+end-start)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to initialize bit sequence handler: %s", err.Error())
	***REMOVED***

	return &Idm***REMOVED***start: start, end: end, handle: h***REMOVED***, nil
***REMOVED***

// GetID returns the first available id in the set
func (i *Idm) GetID(serial bool) (uint64, error) ***REMOVED***
	if i.handle == nil ***REMOVED***
		return 0, errors.New("ID set is not initialized")
	***REMOVED***
	ordinal, err := i.handle.SetAny(serial)
	return i.start + ordinal, err
***REMOVED***

// GetSpecificID tries to reserve the specified id
func (i *Idm) GetSpecificID(id uint64) error ***REMOVED***
	if i.handle == nil ***REMOVED***
		return errors.New("ID set is not initialized")
	***REMOVED***

	if id < i.start || id > i.end ***REMOVED***
		return errors.New("Requested id does not belong to the set")
	***REMOVED***

	return i.handle.Set(id - i.start)
***REMOVED***

// GetIDInRange returns the first available id in the set within a [start,end] range
func (i *Idm) GetIDInRange(start, end uint64, serial bool) (uint64, error) ***REMOVED***
	if i.handle == nil ***REMOVED***
		return 0, errors.New("ID set is not initialized")
	***REMOVED***

	if start < i.start || end > i.end ***REMOVED***
		return 0, errors.New("Requested range does not belong to the set")
	***REMOVED***

	ordinal, err := i.handle.SetAnyInRange(start-i.start, end-i.start, serial)

	return i.start + ordinal, err
***REMOVED***

// Release releases the specified id
func (i *Idm) Release(id uint64) ***REMOVED***
	i.handle.Unset(id - i.start)
***REMOVED***
