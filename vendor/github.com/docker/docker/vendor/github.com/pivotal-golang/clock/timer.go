package clock

import "time"

type Timer interface ***REMOVED***
	C() <-chan time.Time
	Reset(d time.Duration) bool
	Stop() bool
***REMOVED***

type realTimer struct ***REMOVED***
	t *time.Timer
***REMOVED***

func (t *realTimer) C() <-chan time.Time ***REMOVED***
	return t.t.C
***REMOVED***

func (t *realTimer) Reset(d time.Duration) bool ***REMOVED***
	return t.t.Reset(d)
***REMOVED***

func (t *realTimer) Stop() bool ***REMOVED***
	return t.t.Stop()
***REMOVED***
