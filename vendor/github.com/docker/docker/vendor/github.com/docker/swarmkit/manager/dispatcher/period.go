package dispatcher

import (
	"math/rand"
	"time"
)

type periodChooser struct ***REMOVED***
	period  time.Duration
	epsilon time.Duration
	rand    *rand.Rand
***REMOVED***

func newPeriodChooser(period, eps time.Duration) *periodChooser ***REMOVED***
	return &periodChooser***REMOVED***
		period:  period,
		epsilon: eps,
		rand:    rand.New(rand.NewSource(time.Now().UnixNano())),
	***REMOVED***
***REMOVED***

func (pc *periodChooser) Choose() time.Duration ***REMOVED***
	var adj int64
	if pc.epsilon > 0 ***REMOVED***
		adj = rand.Int63n(int64(2*pc.epsilon)) - int64(pc.epsilon)
	***REMOVED***
	return pc.period + time.Duration(adj)
***REMOVED***
