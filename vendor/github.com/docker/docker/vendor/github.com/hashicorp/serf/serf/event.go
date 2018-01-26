package serf

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// EventType are all the types of events that may occur and be sent
// along the Serf channel.
type EventType int

const (
	EventMemberJoin EventType = iota
	EventMemberLeave
	EventMemberFailed
	EventMemberUpdate
	EventMemberReap
	EventUser
	EventQuery
)

func (t EventType) String() string ***REMOVED***
	switch t ***REMOVED***
	case EventMemberJoin:
		return "member-join"
	case EventMemberLeave:
		return "member-leave"
	case EventMemberFailed:
		return "member-failed"
	case EventMemberUpdate:
		return "member-update"
	case EventMemberReap:
		return "member-reap"
	case EventUser:
		return "user"
	case EventQuery:
		return "query"
	default:
		panic(fmt.Sprintf("unknown event type: %d", t))
	***REMOVED***
***REMOVED***

// Event is a generic interface for exposing Serf events
// Clients will usually need to use a type switches to get
// to a more useful type
type Event interface ***REMOVED***
	EventType() EventType
	String() string
***REMOVED***

// MemberEvent is the struct used for member related events
// Because Serf coalesces events, an event may contain multiple members.
type MemberEvent struct ***REMOVED***
	Type    EventType
	Members []Member
***REMOVED***

func (m MemberEvent) EventType() EventType ***REMOVED***
	return m.Type
***REMOVED***

func (m MemberEvent) String() string ***REMOVED***
	switch m.Type ***REMOVED***
	case EventMemberJoin:
		return "member-join"
	case EventMemberLeave:
		return "member-leave"
	case EventMemberFailed:
		return "member-failed"
	case EventMemberUpdate:
		return "member-update"
	case EventMemberReap:
		return "member-reap"
	default:
		panic(fmt.Sprintf("unknown event type: %d", m.Type))
	***REMOVED***
***REMOVED***

// UserEvent is the struct used for events that are triggered
// by the user and are not related to members
type UserEvent struct ***REMOVED***
	LTime    LamportTime
	Name     string
	Payload  []byte
	Coalesce bool
***REMOVED***

func (u UserEvent) EventType() EventType ***REMOVED***
	return EventUser
***REMOVED***

func (u UserEvent) String() string ***REMOVED***
	return fmt.Sprintf("user-event: %s", u.Name)
***REMOVED***

// Query is the struct used EventQuery type events
type Query struct ***REMOVED***
	LTime   LamportTime
	Name    string
	Payload []byte

	serf     *Serf
	id       uint32    // ID is not exported, since it may change
	addr     []byte    // Address to respond to
	port     uint16    // Port to respond to
	deadline time.Time // Must respond by this deadline
	respLock sync.Mutex
***REMOVED***

func (q *Query) EventType() EventType ***REMOVED***
	return EventQuery
***REMOVED***

func (q *Query) String() string ***REMOVED***
	return fmt.Sprintf("query: %s", q.Name)
***REMOVED***

// Deadline returns the time by which a response must be sent
func (q *Query) Deadline() time.Time ***REMOVED***
	return q.deadline
***REMOVED***

// Respond is used to send a response to the user query
func (q *Query) Respond(buf []byte) error ***REMOVED***
	q.respLock.Lock()
	defer q.respLock.Unlock()

	// Check if we've already responded
	if q.deadline.IsZero() ***REMOVED***
		return fmt.Errorf("Response already sent")
	***REMOVED***

	// Ensure we aren't past our response deadline
	if time.Now().After(q.deadline) ***REMOVED***
		return fmt.Errorf("Response is past the deadline")
	***REMOVED***

	// Create response
	resp := messageQueryResponse***REMOVED***
		LTime:   q.LTime,
		ID:      q.id,
		From:    q.serf.config.NodeName,
		Payload: buf,
	***REMOVED***

	// Format the response
	raw, err := encodeMessage(messageQueryResponseType, &resp)
	if err != nil ***REMOVED***
		return fmt.Errorf("Failed to format response: %v", err)
	***REMOVED***

	// Check the size limit
	if len(raw) > q.serf.config.QueryResponseSizeLimit ***REMOVED***
		return fmt.Errorf("response exceeds limit of %d bytes", q.serf.config.QueryResponseSizeLimit)
	***REMOVED***

	// Send the response
	addr := net.UDPAddr***REMOVED***IP: q.addr, Port: int(q.port)***REMOVED***
	if err := q.serf.memberlist.SendTo(&addr, raw); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Clera the deadline, response sent
	q.deadline = time.Time***REMOVED******REMOVED***
	return nil
***REMOVED***
