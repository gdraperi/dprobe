package events

import (
	"container/list"
	"sync"

	"github.com/sirupsen/logrus"
)

// Queue accepts all messages into a queue for asynchronous consumption
// by a sink. It is unbounded and thread safe but the sink must be reliable or
// events will be dropped.
type Queue struct ***REMOVED***
	dst    Sink
	events *list.List
	cond   *sync.Cond
	mu     sync.Mutex
	closed bool
***REMOVED***

// NewQueue returns a queue to the provided Sink dst.
func NewQueue(dst Sink) *Queue ***REMOVED***
	eq := Queue***REMOVED***
		dst:    dst,
		events: list.New(),
	***REMOVED***

	eq.cond = sync.NewCond(&eq.mu)
	go eq.run()
	return &eq
***REMOVED***

// Write accepts the events into the queue, only failing if the queue has
// been closed.
func (eq *Queue) Write(event Event) error ***REMOVED***
	eq.mu.Lock()
	defer eq.mu.Unlock()

	if eq.closed ***REMOVED***
		return ErrSinkClosed
	***REMOVED***

	eq.events.PushBack(event)
	eq.cond.Signal() // signal waiters

	return nil
***REMOVED***

// Close shutsdown the event queue, flushing
func (eq *Queue) Close() error ***REMOVED***
	eq.mu.Lock()
	defer eq.mu.Unlock()

	if eq.closed ***REMOVED***
		return nil
	***REMOVED***

	// set closed flag
	eq.closed = true
	eq.cond.Signal() // signal flushes queue
	eq.cond.Wait()   // wait for signal from last flush
	return eq.dst.Close()
***REMOVED***

// run is the main goroutine to flush events to the target sink.
func (eq *Queue) run() ***REMOVED***
	for ***REMOVED***
		event := eq.next()

		if event == nil ***REMOVED***
			return // nil block means event queue is closed.
		***REMOVED***

		if err := eq.dst.Write(event); err != nil ***REMOVED***
			// TODO(aaronl): Dropping events could be bad depending
			// on the application. We should have a way of
			// communicating this condition. However, logging
			// at a log level above debug may not be appropriate.
			// Eventually, go-events should not use logrus at all,
			// and should bubble up conditions like this through
			// error values.
			logrus.WithFields(logrus.Fields***REMOVED***
				"event": event,
				"sink":  eq.dst,
			***REMOVED***).WithError(err).Debug("eventqueue: dropped event")
		***REMOVED***
	***REMOVED***
***REMOVED***

// next encompasses the critical section of the run loop. When the queue is
// empty, it will block on the condition. If new data arrives, it will wake
// and return a block. When closed, a nil slice will be returned.
func (eq *Queue) next() Event ***REMOVED***
	eq.mu.Lock()
	defer eq.mu.Unlock()

	for eq.events.Len() < 1 ***REMOVED***
		if eq.closed ***REMOVED***
			eq.cond.Broadcast()
			return nil
		***REMOVED***

		eq.cond.Wait()
	***REMOVED***

	front := eq.events.Front()
	block := front.Value.(Event)
	eq.events.Remove(front)

	return block
***REMOVED***
