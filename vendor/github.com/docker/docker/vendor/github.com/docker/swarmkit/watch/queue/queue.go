package queue

import (
	"container/list"
	"fmt"
	"sync"

	"github.com/docker/go-events"
	"github.com/sirupsen/logrus"
)

// ErrQueueFull is returned by a Write operation when that Write causes the
// queue to reach its size limit.
var ErrQueueFull = fmt.Errorf("queue closed due to size limit")

// LimitQueue accepts all messages into a queue for asynchronous consumption by
// a sink until an upper limit of messages is reached. When that limit is
// reached, the entire Queue is Closed. It is thread safe but the
// sink must be reliable or events will be dropped.
// If a size of 0 is provided, the LimitQueue is considered limitless.
type LimitQueue struct ***REMOVED***
	dst        events.Sink
	events     *list.List
	limit      uint64
	cond       *sync.Cond
	mu         sync.Mutex
	closed     bool
	full       chan struct***REMOVED******REMOVED***
	fullClosed bool
***REMOVED***

// NewLimitQueue returns a queue to the provided Sink dst.
func NewLimitQueue(dst events.Sink, limit uint64) *LimitQueue ***REMOVED***
	eq := LimitQueue***REMOVED***
		dst:    dst,
		events: list.New(),
		limit:  limit,
		full:   make(chan struct***REMOVED******REMOVED***),
	***REMOVED***

	eq.cond = sync.NewCond(&eq.mu)
	go eq.run()
	return &eq
***REMOVED***

// Write accepts the events into the queue, only failing if the queue has
// been closed or has reached its size limit.
func (eq *LimitQueue) Write(event events.Event) error ***REMOVED***
	eq.mu.Lock()
	defer eq.mu.Unlock()

	if eq.closed ***REMOVED***
		return events.ErrSinkClosed
	***REMOVED***

	if eq.limit > 0 && uint64(eq.events.Len()) >= eq.limit ***REMOVED***
		// If the limit has been reached, don't write the event to the queue,
		// and close the Full channel. This notifies listeners that the queue
		// is now full, but the sink is still permitted to consume events. It's
		// the responsibility of the listener to decide whether they want to
		// live with dropped events or whether they want to Close() the
		// LimitQueue
		if !eq.fullClosed ***REMOVED***
			eq.fullClosed = true
			close(eq.full)
		***REMOVED***
		return ErrQueueFull
	***REMOVED***

	eq.events.PushBack(event)
	eq.cond.Signal() // signal waiters

	return nil
***REMOVED***

// Full returns a channel that is closed when the queue becomes full for the
// first time.
func (eq *LimitQueue) Full() chan struct***REMOVED******REMOVED*** ***REMOVED***
	return eq.full
***REMOVED***

// Close shuts down the event queue, flushing all events
func (eq *LimitQueue) Close() error ***REMOVED***
	eq.mu.Lock()
	defer eq.mu.Unlock()

	if eq.closed ***REMOVED***
		return nil
	***REMOVED***

	// set the closed flag
	eq.closed = true
	eq.cond.Signal() // signal flushes queue
	eq.cond.Wait()   // wait for signal from last flush
	return eq.dst.Close()
***REMOVED***

// run is the main goroutine to flush events to the target sink.
func (eq *LimitQueue) run() ***REMOVED***
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

// Len returns the number of items that are currently stored in the queue and
// not consumed by its sink.
func (eq *LimitQueue) Len() int ***REMOVED***
	eq.mu.Lock()
	defer eq.mu.Unlock()
	return eq.events.Len()
***REMOVED***

func (eq *LimitQueue) String() string ***REMOVED***
	eq.mu.Lock()
	defer eq.mu.Unlock()
	return fmt.Sprintf("%v", eq.events)
***REMOVED***

// next encompasses the critical section of the run loop. When the queue is
// empty, it will block on the condition. If new data arrives, it will wake
// and return a block. When closed, a nil slice will be returned.
func (eq *LimitQueue) next() events.Event ***REMOVED***
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
	block := front.Value.(events.Event)
	eq.events.Remove(front)

	return block
***REMOVED***
