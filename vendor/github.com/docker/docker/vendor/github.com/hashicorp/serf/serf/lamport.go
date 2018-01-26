package serf

import (
	"sync/atomic"
)

// LamportClock is a thread safe implementation of a lamport clock. It
// uses efficient atomic operations for all of its functions, falling back
// to a heavy lock only if there are enough CAS failures.
type LamportClock struct ***REMOVED***
	counter uint64
***REMOVED***

// LamportTime is the value of a LamportClock.
type LamportTime uint64

// Time is used to return the current value of the lamport clock
func (l *LamportClock) Time() LamportTime ***REMOVED***
	return LamportTime(atomic.LoadUint64(&l.counter))
***REMOVED***

// Increment is used to increment and return the value of the lamport clock
func (l *LamportClock) Increment() LamportTime ***REMOVED***
	return LamportTime(atomic.AddUint64(&l.counter, 1))
***REMOVED***

// Witness is called to update our local clock if necessary after
// witnessing a clock value received from another process
func (l *LamportClock) Witness(v LamportTime) ***REMOVED***
WITNESS:
	// If the other value is old, we do not need to do anything
	cur := atomic.LoadUint64(&l.counter)
	other := uint64(v)
	if other < cur ***REMOVED***
		return
	***REMOVED***

	// Ensure that our local clock is at least one ahead.
	if !atomic.CompareAndSwapUint64(&l.counter, cur, other+1) ***REMOVED***
		// The CAS failed, so we just retry. Eventually our CAS should
		// succeed or a future witness will pass us by and our witness
		// will end.
		goto WITNESS
	***REMOVED***
***REMOVED***
