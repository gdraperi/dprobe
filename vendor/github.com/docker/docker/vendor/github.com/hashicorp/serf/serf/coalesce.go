package serf

import (
	"time"
)

// coalescer is a simple interface that must be implemented to be
// used inside of a coalesceLoop
type coalescer interface ***REMOVED***
	// Can the coalescer handle this event, if not it is
	// directly passed through to the destination channel
	Handle(Event) bool

	// Invoked to coalesce the given event
	Coalesce(Event)

	// Invoked to flush the coalesced events
	Flush(outChan chan<- Event)
***REMOVED***

// coalescedEventCh returns an event channel where the events are coalesced
// using the given coalescer.
func coalescedEventCh(outCh chan<- Event, shutdownCh <-chan struct***REMOVED******REMOVED***,
	cPeriod time.Duration, qPeriod time.Duration, c coalescer) chan<- Event ***REMOVED***
	inCh := make(chan Event, 1024)
	go coalesceLoop(inCh, outCh, shutdownCh, cPeriod, qPeriod, c)
	return inCh
***REMOVED***

// coalesceLoop is a simple long-running routine that manages the high-level
// flow of coalescing based on quiescence and a maximum quantum period.
func coalesceLoop(inCh <-chan Event, outCh chan<- Event, shutdownCh <-chan struct***REMOVED******REMOVED***,
	coalescePeriod time.Duration, quiescentPeriod time.Duration, c coalescer) ***REMOVED***
	var quiescent <-chan time.Time
	var quantum <-chan time.Time
	shutdown := false

INGEST:
	// Reset the timers
	quantum = nil
	quiescent = nil

	for ***REMOVED***
		select ***REMOVED***
		case e := <-inCh:
			// Ignore any non handled events
			if !c.Handle(e) ***REMOVED***
				outCh <- e
				continue
			***REMOVED***

			// Start a new quantum if we need to
			// and restart the quiescent timer
			if quantum == nil ***REMOVED***
				quantum = time.After(coalescePeriod)
			***REMOVED***
			quiescent = time.After(quiescentPeriod)

			// Coalesce the event
			c.Coalesce(e)

		case <-quantum:
			goto FLUSH
		case <-quiescent:
			goto FLUSH
		case <-shutdownCh:
			shutdown = true
			goto FLUSH
		***REMOVED***
	***REMOVED***

FLUSH:
	// Flush the coalesced events
	c.Flush(outCh)

	// Restart ingestion if we are not done
	if !shutdown ***REMOVED***
		goto INGEST
	***REMOVED***
***REMOVED***
