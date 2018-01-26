package watch

import (
	"fmt"
	"time"

	events "github.com/docker/go-events"
)

// ErrSinkTimeout is returned from the Write method when a sink times out.
var ErrSinkTimeout = fmt.Errorf("timeout exceeded, tearing down sink")

// timeoutSink is a sink that wraps another sink with a timeout. If the
// embedded sink fails to complete a Write operation within the specified
// timeout, the Write operation of the timeoutSink fails.
type timeoutSink struct ***REMOVED***
	timeout time.Duration
	sink    events.Sink
***REMOVED***

func (s timeoutSink) Write(event events.Event) error ***REMOVED***
	errChan := make(chan error)
	go func(c chan<- error) ***REMOVED***
		c <- s.sink.Write(event)
	***REMOVED***(errChan)

	timer := time.NewTimer(s.timeout)
	select ***REMOVED***
	case err := <-errChan:
		timer.Stop()
		return err
	case <-timer.C:
		s.sink.Close()
		return ErrSinkTimeout
	***REMOVED***
***REMOVED***

func (s timeoutSink) Close() error ***REMOVED***
	return s.sink.Close()
***REMOVED***

// dropErrClosed is a sink that suppresses ErrSinkClosed from Write, to avoid
// debug log messages that may be confusing. It is possible that the queue
// will try to write an event to its destination channel while the queue is
// being removed from the broadcaster. Since the channel is closed before the
// queue, there is a narrow window when this is possible. In some event-based
// dropping events when a sink is removed from a broadcaster is a problem, but
// for the usage in this watch package that's the expected behavior.
type dropErrClosed struct ***REMOVED***
	sink events.Sink
***REMOVED***

func (s dropErrClosed) Write(event events.Event) error ***REMOVED***
	err := s.sink.Write(event)
	if err == events.ErrSinkClosed ***REMOVED***
		return nil
	***REMOVED***
	return err
***REMOVED***

func (s dropErrClosed) Close() error ***REMOVED***
	return s.sink.Close()
***REMOVED***

// dropErrClosedChanGen is a ChannelSinkGenerator for dropErrClosed sinks wrapping
// unbuffered channels.
type dropErrClosedChanGen struct***REMOVED******REMOVED***

func (s *dropErrClosedChanGen) NewChannelSink() (events.Sink, *events.Channel) ***REMOVED***
	ch := events.NewChannel(0)
	return dropErrClosed***REMOVED***sink: ch***REMOVED***, ch
***REMOVED***

// TimeoutDropErrChanGen is a ChannelSinkGenerator that creates a channel,
// wrapped by the dropErrClosed sink and a timeout.
type TimeoutDropErrChanGen struct ***REMOVED***
	timeout time.Duration
***REMOVED***

// NewChannelSink creates a new sink chain of timeoutSink->dropErrClosed->Channel
func (s *TimeoutDropErrChanGen) NewChannelSink() (events.Sink, *events.Channel) ***REMOVED***
	ch := events.NewChannel(0)
	return timeoutSink***REMOVED***
		timeout: s.timeout,
		sink: dropErrClosed***REMOVED***
			sink: ch,
		***REMOVED***,
	***REMOVED***, ch
***REMOVED***

// NewTimeoutDropErrSinkGen returns a generator of timeoutSinks wrapping dropErrClosed
// sinks, wrapping unbuffered channel sinks.
func NewTimeoutDropErrSinkGen(timeout time.Duration) ChannelSinkGenerator ***REMOVED***
	return &TimeoutDropErrChanGen***REMOVED***timeout: timeout***REMOVED***
***REMOVED***
