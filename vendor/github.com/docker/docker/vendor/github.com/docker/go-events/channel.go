package events

import (
	"fmt"
	"sync"
)

// Channel provides a sink that can be listened on. The writer and channel
// listener must operate in separate goroutines.
//
// Consumers should listen on Channel.C until Closed is closed.
type Channel struct ***REMOVED***
	C chan Event

	closed chan struct***REMOVED******REMOVED***
	once   sync.Once
***REMOVED***

// NewChannel returns a channel. If buffer is zero, the channel is
// unbuffered.
func NewChannel(buffer int) *Channel ***REMOVED***
	return &Channel***REMOVED***
		C:      make(chan Event, buffer),
		closed: make(chan struct***REMOVED******REMOVED***),
	***REMOVED***
***REMOVED***

// Done returns a channel that will always proceed once the sink is closed.
func (ch *Channel) Done() chan struct***REMOVED******REMOVED*** ***REMOVED***
	return ch.closed
***REMOVED***

// Write the event to the channel. Must be called in a separate goroutine from
// the listener.
func (ch *Channel) Write(event Event) error ***REMOVED***
	select ***REMOVED***
	case ch.C <- event:
		return nil
	case <-ch.closed:
		return ErrSinkClosed
	***REMOVED***
***REMOVED***

// Close the channel sink.
func (ch *Channel) Close() error ***REMOVED***
	ch.once.Do(func() ***REMOVED***
		close(ch.closed)
	***REMOVED***)

	return nil
***REMOVED***

func (ch *Channel) String() string ***REMOVED***
	// Serialize a copy of the Channel that doesn't contain the sync.Once,
	// to avoid a data race.
	ch2 := map[string]interface***REMOVED******REMOVED******REMOVED***
		"C":      ch.C,
		"closed": ch.closed,
	***REMOVED***
	return fmt.Sprint(ch2)
***REMOVED***
