package clock

import "time"

type Ticker interface ***REMOVED***
	C() <-chan time.Time
	Stop()
***REMOVED***

type realTicker struct ***REMOVED***
	t *time.Ticker
***REMOVED***

func (t *realTicker) C() <-chan time.Time ***REMOVED***
	return t.t.C
***REMOVED***

func (t *realTicker) Stop() ***REMOVED***
	t.t.Stop()
***REMOVED***
