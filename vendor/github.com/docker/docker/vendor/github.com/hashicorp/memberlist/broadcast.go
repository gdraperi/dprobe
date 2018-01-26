package memberlist

/*
The broadcast mechanism works by maintaining a sorted list of messages to be
sent out. When a message is to be broadcast, the retransmit count
is set to zero and appended to the queue. The retransmit count serves
as the "priority", ensuring that newer messages get sent first. Once
a message hits the retransmit limit, it is removed from the queue.

Additionally, older entries can be invalidated by new messages that
are contradictory. For example, if we send "***REMOVED***suspect M1 inc: 1***REMOVED***,
then a following ***REMOVED***alive M1 inc: 2***REMOVED*** will invalidate that message
*/

type memberlistBroadcast struct ***REMOVED***
	node   string
	msg    []byte
	notify chan struct***REMOVED******REMOVED***
***REMOVED***

func (b *memberlistBroadcast) Invalidates(other Broadcast) bool ***REMOVED***
	// Check if that broadcast is a memberlist type
	mb, ok := other.(*memberlistBroadcast)
	if !ok ***REMOVED***
		return false
	***REMOVED***

	// Invalidates any message about the same node
	return b.node == mb.node
***REMOVED***

func (b *memberlistBroadcast) Message() []byte ***REMOVED***
	return b.msg
***REMOVED***

func (b *memberlistBroadcast) Finished() ***REMOVED***
	select ***REMOVED***
	case b.notify <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
	default:
	***REMOVED***
***REMOVED***

// encodeAndBroadcast encodes a message and enqueues it for broadcast. Fails
// silently if there is an encoding error.
func (m *Memberlist) encodeAndBroadcast(node string, msgType messageType, msg interface***REMOVED******REMOVED***) ***REMOVED***
	m.encodeBroadcastNotify(node, msgType, msg, nil)
***REMOVED***

// encodeBroadcastNotify encodes a message and enqueues it for broadcast
// and notifies the given channel when transmission is finished. Fails
// silently if there is an encoding error.
func (m *Memberlist) encodeBroadcastNotify(node string, msgType messageType, msg interface***REMOVED******REMOVED***, notify chan struct***REMOVED******REMOVED***) ***REMOVED***
	buf, err := encode(msgType, msg)
	if err != nil ***REMOVED***
		m.logger.Printf("[ERR] memberlist: Failed to encode message for broadcast: %s", err)
	***REMOVED*** else ***REMOVED***
		m.queueBroadcast(node, buf.Bytes(), notify)
	***REMOVED***
***REMOVED***

// queueBroadcast is used to start dissemination of a message. It will be
// sent up to a configured number of times. The message could potentially
// be invalidated by a future message about the same node
func (m *Memberlist) queueBroadcast(node string, msg []byte, notify chan struct***REMOVED******REMOVED***) ***REMOVED***
	b := &memberlistBroadcast***REMOVED***node, msg, notify***REMOVED***
	m.broadcasts.QueueBroadcast(b)
***REMOVED***

// getBroadcasts is used to return a slice of broadcasts to send up to
// a maximum byte size, while imposing a per-broadcast overhead. This is used
// to fill a UDP packet with piggybacked data
func (m *Memberlist) getBroadcasts(overhead, limit int) [][]byte ***REMOVED***
	// Get memberlist messages first
	toSend := m.broadcasts.GetBroadcasts(overhead, limit)

	// Check if the user has anything to broadcast
	d := m.config.Delegate
	if d != nil ***REMOVED***
		// Determine the bytes used already
		bytesUsed := 0
		for _, msg := range toSend ***REMOVED***
			bytesUsed += len(msg) + overhead
		***REMOVED***

		// Check space remaining for user messages
		avail := limit - bytesUsed
		if avail > overhead+userMsgOverhead ***REMOVED***
			userMsgs := d.GetBroadcasts(overhead+userMsgOverhead, avail)

			// Frame each user message
			for _, msg := range userMsgs ***REMOVED***
				buf := make([]byte, 1, len(msg)+1)
				buf[0] = byte(userMsg)
				buf = append(buf, msg...)
				toSend = append(toSend, buf)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return toSend
***REMOVED***
