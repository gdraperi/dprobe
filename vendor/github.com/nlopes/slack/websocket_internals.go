package slack

import (
	"fmt"
	"time"
)

/**
 * Internal events, created by this lib and not mapped to Slack APIs.
 */

// ConnectedEvent is used for when we connect to Slack
type ConnectedEvent struct ***REMOVED***
	ConnectionCount int // 1 = first time, 2 = second time
	Info            *Info
***REMOVED***

// ConnectionErrorEvent contains information about a connection error
type ConnectionErrorEvent struct ***REMOVED***
	Attempt  int
	ErrorObj error
***REMOVED***

func (c *ConnectionErrorEvent) Error() string ***REMOVED***
	return c.ErrorObj.Error()
***REMOVED***

// ConnectingEvent contains information about our connection attempt
type ConnectingEvent struct ***REMOVED***
	Attempt         int // 1 = first attempt, 2 = second attempt
	ConnectionCount int
***REMOVED***

// DisconnectedEvent contains information about how we disconnected
type DisconnectedEvent struct ***REMOVED***
	Intentional bool
***REMOVED***

// LatencyReport contains information about connection latency
type LatencyReport struct ***REMOVED***
	Value time.Duration
***REMOVED***

// InvalidAuthEvent is used in case we can't even authenticate with the API
type InvalidAuthEvent struct***REMOVED******REMOVED***

// UnmarshallingErrorEvent is used when there are issues deconstructing a response
type UnmarshallingErrorEvent struct ***REMOVED***
	ErrorObj error
***REMOVED***

func (u UnmarshallingErrorEvent) Error() string ***REMOVED***
	return u.ErrorObj.Error()
***REMOVED***

// MessageTooLongEvent is used when sending a message that is too long
type MessageTooLongEvent struct ***REMOVED***
	Message   OutgoingMessage
	MaxLength int
***REMOVED***

func (m *MessageTooLongEvent) Error() string ***REMOVED***
	return fmt.Sprintf("Message too long (max %d characters)", m.MaxLength)
***REMOVED***

// RateLimitEvent is used when Slack warns that rate-limits are being hit.
type RateLimitEvent struct***REMOVED******REMOVED***

func (e *RateLimitEvent) Error() string ***REMOVED***
	return "Messages are being sent too fast."
***REMOVED***

// OutgoingErrorEvent contains information in case there were errors sending messages
type OutgoingErrorEvent struct ***REMOVED***
	Message  OutgoingMessage
	ErrorObj error
***REMOVED***

func (o OutgoingErrorEvent) Error() string ***REMOVED***
	return o.ErrorObj.Error()
***REMOVED***

// IncomingEventError contains information about an unexpected error receiving a websocket event
type IncomingEventError struct ***REMOVED***
	ErrorObj error
***REMOVED***

func (i *IncomingEventError) Error() string ***REMOVED***
	return i.ErrorObj.Error()
***REMOVED***

// AckErrorEvent i
type AckErrorEvent struct ***REMOVED***
	ErrorObj error
***REMOVED***

func (a *AckErrorEvent) Error() string ***REMOVED***
	return a.ErrorObj.Error()
***REMOVED***
