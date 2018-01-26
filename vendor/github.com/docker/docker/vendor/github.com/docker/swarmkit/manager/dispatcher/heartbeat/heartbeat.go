package heartbeat

import (
	"sync/atomic"
	"time"
)

// Heartbeat is simple way to track heartbeats.
type Heartbeat struct ***REMOVED***
	timeout int64
	timer   *time.Timer
***REMOVED***

// New creates new Heartbeat with specified duration. timeoutFunc will be called
// if timeout for heartbeat is expired. Note that in case of timeout you need to
// call Beat() to reactivate Heartbeat.
func New(timeout time.Duration, timeoutFunc func()) *Heartbeat ***REMOVED***
	hb := &Heartbeat***REMOVED***
		timeout: int64(timeout),
		timer:   time.AfterFunc(timeout, timeoutFunc),
	***REMOVED***
	return hb
***REMOVED***

// Beat resets internal timer to zero. It also can be used to reactivate
// Heartbeat after timeout.
func (hb *Heartbeat) Beat() ***REMOVED***
	hb.timer.Reset(time.Duration(atomic.LoadInt64(&hb.timeout)))
***REMOVED***

// Update updates internal timeout to d. It does not do Beat.
func (hb *Heartbeat) Update(d time.Duration) ***REMOVED***
	atomic.StoreInt64(&hb.timeout, int64(d))
***REMOVED***

// Stop stops Heartbeat timer.
func (hb *Heartbeat) Stop() ***REMOVED***
	hb.timer.Stop()
***REMOVED***
