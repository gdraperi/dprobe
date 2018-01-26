package memberlist

import (
	"math"
	"sync/atomic"
	"time"
)

// suspicion manages the suspect timer for a node and provides an interface
// to accelerate the timeout as we get more independent confirmations that
// a node is suspect.
type suspicion struct ***REMOVED***
	// n is the number of independent confirmations we've seen. This must
	// be updated using atomic instructions to prevent contention with the
	// timer callback.
	n int32

	// k is the number of independent confirmations we'd like to see in
	// order to drive the timer to its minimum value.
	k int32

	// min is the minimum timer value.
	min time.Duration

	// max is the maximum timer value.
	max time.Duration

	// start captures the timestamp when we began the timer. This is used
	// so we can calculate durations to feed the timer during updates in
	// a way the achieves the overall time we'd like.
	start time.Time

	// timer is the underlying timer that implements the timeout.
	timer *time.Timer

	// f is the function to call when the timer expires. We hold on to this
	// because there are cases where we call it directly.
	timeoutFn func()

	// confirmations is a map of "from" nodes that have confirmed a given
	// node is suspect. This prevents double counting.
	confirmations map[string]struct***REMOVED******REMOVED***
***REMOVED***

// newSuspicion returns a timer started with the max time, and that will drive
// to the min time after seeing k or more confirmations. The from node will be
// excluded from confirmations since we might get our own suspicion message
// gossiped back to us. The minimum time will be used if no confirmations are
// called for (k <= 0).
func newSuspicion(from string, k int, min time.Duration, max time.Duration, fn func(int)) *suspicion ***REMOVED***
	s := &suspicion***REMOVED***
		k:             int32(k),
		min:           min,
		max:           max,
		confirmations: make(map[string]struct***REMOVED******REMOVED***),
	***REMOVED***

	// Exclude the from node from any confirmations.
	s.confirmations[from] = struct***REMOVED******REMOVED******REMOVED******REMOVED***

	// Pass the number of confirmations into the timeout function for
	// easy telemetry.
	s.timeoutFn = func() ***REMOVED***
		fn(int(atomic.LoadInt32(&s.n)))
	***REMOVED***

	// If there aren't any confirmations to be made then take the min
	// time from the start.
	timeout := max
	if k < 1 ***REMOVED***
		timeout = min
	***REMOVED***
	s.timer = time.AfterFunc(timeout, s.timeoutFn)

	// Capture the start time right after starting the timer above so
	// we should always err on the side of a little longer timeout if
	// there's any preemption that separates this and the step above.
	s.start = time.Now()
	return s
***REMOVED***

// remainingSuspicionTime takes the state variables of the suspicion timer and
// calculates the remaining time to wait before considering a node dead. The
// return value can be negative, so be prepared to fire the timer immediately in
// that case.
func remainingSuspicionTime(n, k int32, elapsed time.Duration, min, max time.Duration) time.Duration ***REMOVED***
	frac := math.Log(float64(n)+1.0) / math.Log(float64(k)+1.0)
	raw := max.Seconds() - frac*(max.Seconds()-min.Seconds())
	timeout := time.Duration(math.Floor(1000.0*raw)) * time.Millisecond
	if timeout < min ***REMOVED***
		timeout = min
	***REMOVED***

	// We have to take into account the amount of time that has passed so
	// far, so we get the right overall timeout.
	return timeout - elapsed
***REMOVED***

// Confirm registers that a possibly new peer has also determined the given
// node is suspect. This returns true if this was new information, and false
// if it was a duplicate confirmation, or if we've got enough confirmations to
// hit the minimum.
func (s *suspicion) Confirm(from string) bool ***REMOVED***
	// If we've got enough confirmations then stop accepting them.
	if atomic.LoadInt32(&s.n) >= s.k ***REMOVED***
		return false
	***REMOVED***

	// Only allow one confirmation from each possible peer.
	if _, ok := s.confirmations[from]; ok ***REMOVED***
		return false
	***REMOVED***
	s.confirmations[from] = struct***REMOVED******REMOVED******REMOVED******REMOVED***

	// Compute the new timeout given the current number of confirmations and
	// adjust the timer. If the timeout becomes negative *and* we can cleanly
	// stop the timer then we will call the timeout function directly from
	// here.
	n := atomic.AddInt32(&s.n, 1)
	elapsed := time.Now().Sub(s.start)
	remaining := remainingSuspicionTime(n, s.k, elapsed, s.min, s.max)
	if s.timer.Stop() ***REMOVED***
		if remaining > 0 ***REMOVED***
			s.timer.Reset(remaining)
		***REMOVED*** else ***REMOVED***
			go s.timeoutFn()
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***
