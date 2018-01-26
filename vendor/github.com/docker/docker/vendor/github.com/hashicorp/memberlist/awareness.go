package memberlist

import (
	"sync"
	"time"

	"github.com/armon/go-metrics"
)

// awareness manages a simple metric for tracking the estimated health of the
// local node. Health is primary the node's ability to respond in the soft
// real-time manner required for correct health checking of other nodes in the
// cluster.
type awareness struct ***REMOVED***
	sync.RWMutex

	// max is the upper threshold for the timeout scale (the score will be
	// constrained to be from 0 <= score < max).
	max int

	// score is the current awareness score. Lower values are healthier and
	// zero is the minimum value.
	score int
***REMOVED***

// newAwareness returns a new awareness object.
func newAwareness(max int) *awareness ***REMOVED***
	return &awareness***REMOVED***
		max:   max,
		score: 0,
	***REMOVED***
***REMOVED***

// ApplyDelta takes the given delta and applies it to the score in a thread-safe
// manner. It also enforces a floor of zero and a max of max, so deltas may not
// change the overall score if it's railed at one of the extremes.
func (a *awareness) ApplyDelta(delta int) ***REMOVED***
	a.Lock()
	initial := a.score
	a.score += delta
	if a.score < 0 ***REMOVED***
		a.score = 0
	***REMOVED*** else if a.score > (a.max - 1) ***REMOVED***
		a.score = (a.max - 1)
	***REMOVED***
	final := a.score
	a.Unlock()

	if initial != final ***REMOVED***
		metrics.SetGauge([]string***REMOVED***"memberlist", "health", "score"***REMOVED***, float32(final))
	***REMOVED***
***REMOVED***

// GetHealthScore returns the raw health score.
func (a *awareness) GetHealthScore() int ***REMOVED***
	a.RLock()
	score := a.score
	a.RUnlock()
	return score
***REMOVED***

// ScaleTimeout takes the given duration and scales it based on the current
// score. Less healthyness will lead to longer timeouts.
func (a *awareness) ScaleTimeout(timeout time.Duration) time.Duration ***REMOVED***
	a.RLock()
	score := a.score
	a.RUnlock()
	return timeout * (time.Duration(score) + 1)
***REMOVED***
