package serf

import (
	"bytes"
	"github.com/hashicorp/go-msgpack/codec"
	"time"
)

// messageType are the types of gossip messages Serf will send along
// memberlist.
type messageType uint8

const (
	messageLeaveType messageType = iota
	messageJoinType
	messagePushPullType
	messageUserEventType
	messageQueryType
	messageQueryResponseType
	messageConflictResponseType
	messageKeyRequestType
	messageKeyResponseType
)

const (
	// Ack flag is used to force receiver to send an ack back
	queryFlagAck uint32 = 1 << iota

	// NoBroadcast is used to prevent re-broadcast of a query.
	// this can be used to selectively send queries to individual members
	queryFlagNoBroadcast
)

// filterType is used with a queryFilter to specify the type of
// filter we are sending
type filterType uint8

const (
	filterNodeType filterType = iota
	filterTagType
)

// messageJoin is the message broadcasted after we join to
// associated the node with a lamport clock
type messageJoin struct ***REMOVED***
	LTime LamportTime
	Node  string
***REMOVED***

// messageLeave is the message broadcasted to signal the intentional to
// leave.
type messageLeave struct ***REMOVED***
	LTime LamportTime
	Node  string
***REMOVED***

// messagePushPullType is used when doing a state exchange. This
// is a relatively large message, but is sent infrequently
type messagePushPull struct ***REMOVED***
	LTime        LamportTime            // Current node lamport time
	StatusLTimes map[string]LamportTime // Maps the node to its status time
	LeftMembers  []string               // List of left nodes
	EventLTime   LamportTime            // Lamport time for event clock
	Events       []*userEvents          // Recent events
	QueryLTime   LamportTime            // Lamport time for query clock
***REMOVED***

// messageUserEvent is used for user-generated events
type messageUserEvent struct ***REMOVED***
	LTime   LamportTime
	Name    string
	Payload []byte
	CC      bool // "Can Coalesce". Zero value is compatible with Serf 0.1
***REMOVED***

// messageQuery is used for query events
type messageQuery struct ***REMOVED***
	LTime   LamportTime   // Event lamport time
	ID      uint32        // Query ID, randomly generated
	Addr    []byte        // Source address, used for a direct reply
	Port    uint16        // Source port, used for a direct reply
	Filters [][]byte      // Potential query filters
	Flags   uint32        // Used to provide various flags
	Timeout time.Duration // Maximum time between delivery and response
	Name    string        // Query name
	Payload []byte        // Query payload
***REMOVED***

// Ack checks if the ack flag is set
func (m *messageQuery) Ack() bool ***REMOVED***
	return (m.Flags & queryFlagAck) != 0
***REMOVED***

// NoBroadcast checks if the no broadcast flag is set
func (m *messageQuery) NoBroadcast() bool ***REMOVED***
	return (m.Flags & queryFlagNoBroadcast) != 0
***REMOVED***

// filterNode is used with the filterNodeType, and is a list
// of node names
type filterNode []string

// filterTag is used with the filterTagType and is a regular
// expression to apply to a tag
type filterTag struct ***REMOVED***
	Tag  string
	Expr string
***REMOVED***

// messageQueryResponse is used to respond to a query
type messageQueryResponse struct ***REMOVED***
	LTime   LamportTime // Event lamport time
	ID      uint32      // Query ID
	From    string      // Node name
	Flags   uint32      // Used to provide various flags
	Payload []byte      // Optional response payload
***REMOVED***

// Ack checks if the ack flag is set
func (m *messageQueryResponse) Ack() bool ***REMOVED***
	return (m.Flags & queryFlagAck) != 0
***REMOVED***

func decodeMessage(buf []byte, out interface***REMOVED******REMOVED***) error ***REMOVED***
	var handle codec.MsgpackHandle
	return codec.NewDecoder(bytes.NewReader(buf), &handle).Decode(out)
***REMOVED***

func encodeMessage(t messageType, msg interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	buf := bytes.NewBuffer(nil)
	buf.WriteByte(uint8(t))

	handle := codec.MsgpackHandle***REMOVED******REMOVED***
	encoder := codec.NewEncoder(buf, &handle)
	err := encoder.Encode(msg)
	return buf.Bytes(), err
***REMOVED***

func encodeFilter(f filterType, filt interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	buf := bytes.NewBuffer(nil)
	buf.WriteByte(uint8(f))

	handle := codec.MsgpackHandle***REMOVED******REMOVED***
	encoder := codec.NewEncoder(buf, &handle)
	err := encoder.Encode(filt)
	return buf.Bytes(), err
***REMOVED***
