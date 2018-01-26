package local

import (
	"sync"

	"github.com/containerd/containerd/errdefs"
	"github.com/pkg/errors"
)

// Handles locking references

var (
	// locks lets us lock in process
	locks   = map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	locksMu sync.Mutex
)

func tryLock(ref string) error ***REMOVED***
	locksMu.Lock()
	defer locksMu.Unlock()

	if _, ok := locks[ref]; ok ***REMOVED***
		return errors.Wrapf(errdefs.ErrUnavailable, "ref %s locked", ref)
	***REMOVED***

	locks[ref] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	return nil
***REMOVED***

func unlock(ref string) ***REMOVED***
	locksMu.Lock()
	defer locksMu.Unlock()

	if _, ok := locks[ref]; ok ***REMOVED***
		delete(locks, ref)
	***REMOVED***
***REMOVED***
