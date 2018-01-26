package clock

import "time"

type Clock interface ***REMOVED***
	Now() time.Time
	Sleep(d time.Duration)
	Since(t time.Time) time.Duration

	NewTimer(d time.Duration) Timer
	NewTicker(d time.Duration) Ticker
***REMOVED***

type realClock struct***REMOVED******REMOVED***

func NewClock() Clock ***REMOVED***
	return &realClock***REMOVED******REMOVED***
***REMOVED***

func (clock *realClock) Now() time.Time ***REMOVED***
	return time.Now()
***REMOVED***

func (clock *realClock) Since(t time.Time) time.Duration ***REMOVED***
	return time.Now().Sub(t)
***REMOVED***

func (clock *realClock) Sleep(d time.Duration) ***REMOVED***
	<-clock.NewTimer(d).C()
***REMOVED***

func (clock *realClock) NewTimer(d time.Duration) Timer ***REMOVED***
	return &realTimer***REMOVED***
		t: time.NewTimer(d),
	***REMOVED***
***REMOVED***

func (clock *realClock) NewTicker(d time.Duration) Ticker ***REMOVED***
	return &realTicker***REMOVED***
		t: time.NewTicker(d),
	***REMOVED***
***REMOVED***
