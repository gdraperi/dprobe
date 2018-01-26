package events

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

// Broadcaster sends events to multiple, reliable Sinks. The goal of this
// component is to dispatch events to configured endpoints. Reliability can be
// provided by wrapping incoming sinks.
type Broadcaster struct ***REMOVED***
	sinks   []Sink
	events  chan Event
	adds    chan configureRequest
	removes chan configureRequest

	shutdown chan struct***REMOVED******REMOVED***
	closed   chan struct***REMOVED******REMOVED***
	once     sync.Once
***REMOVED***

// NewBroadcaster appends one or more sinks to the list of sinks. The
// broadcaster behavior will be affected by the properties of the sink.
// Generally, the sink should accept all messages and deal with reliability on
// its own. Use of EventQueue and RetryingSink should be used here.
func NewBroadcaster(sinks ...Sink) *Broadcaster ***REMOVED***
	b := Broadcaster***REMOVED***
		sinks:    sinks,
		events:   make(chan Event),
		adds:     make(chan configureRequest),
		removes:  make(chan configureRequest),
		shutdown: make(chan struct***REMOVED******REMOVED***),
		closed:   make(chan struct***REMOVED******REMOVED***),
	***REMOVED***

	// Start the broadcaster
	go b.run()

	return &b
***REMOVED***

// Write accepts an event to be dispatched to all sinks. This method will never
// fail and should never block (hopefully!). The caller cedes the memory to the
// broadcaster and should not modify it after calling write.
func (b *Broadcaster) Write(event Event) error ***REMOVED***
	select ***REMOVED***
	case b.events <- event:
	case <-b.closed:
		return ErrSinkClosed
	***REMOVED***
	return nil
***REMOVED***

// Add the sink to the broadcaster.
//
// The provided sink must be comparable with equality. Typically, this just
// works with a regular pointer type.
func (b *Broadcaster) Add(sink Sink) error ***REMOVED***
	return b.configure(b.adds, sink)
***REMOVED***

// Remove the provided sink.
func (b *Broadcaster) Remove(sink Sink) error ***REMOVED***
	return b.configure(b.removes, sink)
***REMOVED***

type configureRequest struct ***REMOVED***
	sink     Sink
	response chan error
***REMOVED***

func (b *Broadcaster) configure(ch chan configureRequest, sink Sink) error ***REMOVED***
	response := make(chan error, 1)

	for ***REMOVED***
		select ***REMOVED***
		case ch <- configureRequest***REMOVED***
			sink:     sink,
			response: response***REMOVED***:
			ch = nil
		case err := <-response:
			return err
		case <-b.closed:
			return ErrSinkClosed
		***REMOVED***
	***REMOVED***
***REMOVED***

// Close the broadcaster, ensuring that all messages are flushed to the
// underlying sink before returning.
func (b *Broadcaster) Close() error ***REMOVED***
	b.once.Do(func() ***REMOVED***
		close(b.shutdown)
	***REMOVED***)

	<-b.closed
	return nil
***REMOVED***

// run is the main broadcast loop, started when the broadcaster is created.
// Under normal conditions, it waits for events on the event channel. After
// Close is called, this goroutine will exit.
func (b *Broadcaster) run() ***REMOVED***
	defer close(b.closed)
	remove := func(target Sink) ***REMOVED***
		for i, sink := range b.sinks ***REMOVED***
			if sink == target ***REMOVED***
				b.sinks = append(b.sinks[:i], b.sinks[i+1:]...)
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	for ***REMOVED***
		select ***REMOVED***
		case event := <-b.events:
			for _, sink := range b.sinks ***REMOVED***
				if err := sink.Write(event); err != nil ***REMOVED***
					if err == ErrSinkClosed ***REMOVED***
						// remove closed sinks
						remove(sink)
						continue
					***REMOVED***
					logrus.WithField("event", event).WithField("events.sink", sink).WithError(err).
						Errorf("broadcaster: dropping event")
				***REMOVED***
			***REMOVED***
		case request := <-b.adds:
			// while we have to iterate for add/remove, common iteration for
			// send is faster against slice.

			var found bool
			for _, sink := range b.sinks ***REMOVED***
				if request.sink == sink ***REMOVED***
					found = true
					break
				***REMOVED***
			***REMOVED***

			if !found ***REMOVED***
				b.sinks = append(b.sinks, request.sink)
			***REMOVED***
			// b.sinks[request.sink] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			request.response <- nil
		case request := <-b.removes:
			remove(request.sink)
			request.response <- nil
		case <-b.shutdown:
			// close all the underlying sinks
			for _, sink := range b.sinks ***REMOVED***
				if err := sink.Close(); err != nil && err != ErrSinkClosed ***REMOVED***
					logrus.WithField("events.sink", sink).WithError(err).
						Errorf("broadcaster: closing sink failed")
				***REMOVED***
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (b *Broadcaster) String() string ***REMOVED***
	// Serialize copy of this broadcaster without the sync.Once, to avoid
	// a data race.

	b2 := map[string]interface***REMOVED******REMOVED******REMOVED***
		"sinks":   b.sinks,
		"events":  b.events,
		"adds":    b.adds,
		"removes": b.removes,

		"shutdown": b.shutdown,
		"closed":   b.closed,
	***REMOVED***

	return fmt.Sprint(b2)
***REMOVED***
