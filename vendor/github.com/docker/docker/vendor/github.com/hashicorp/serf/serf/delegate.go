package serf

import (
	"fmt"

	"github.com/armon/go-metrics"
)

// delegate is the memberlist.Delegate implementation that Serf uses.
type delegate struct ***REMOVED***
	serf *Serf
***REMOVED***

func (d *delegate) NodeMeta(limit int) []byte ***REMOVED***
	roleBytes := d.serf.encodeTags(d.serf.config.Tags)
	if len(roleBytes) > limit ***REMOVED***
		panic(fmt.Errorf("Node tags '%v' exceeds length limit of %d bytes", d.serf.config.Tags, limit))
	***REMOVED***

	return roleBytes
***REMOVED***

func (d *delegate) NotifyMsg(buf []byte) ***REMOVED***
	// If we didn't actually receive any data, then ignore it.
	if len(buf) == 0 ***REMOVED***
		return
	***REMOVED***
	metrics.AddSample([]string***REMOVED***"serf", "msgs", "received"***REMOVED***, float32(len(buf)))

	rebroadcast := false
	rebroadcastQueue := d.serf.broadcasts
	t := messageType(buf[0])
	switch t ***REMOVED***
	case messageLeaveType:
		var leave messageLeave
		if err := decodeMessage(buf[1:], &leave); err != nil ***REMOVED***
			d.serf.logger.Printf("[ERR] serf: Error decoding leave message: %s", err)
			break
		***REMOVED***

		d.serf.logger.Printf("[DEBUG] serf: messageLeaveType: %s", leave.Node)
		rebroadcast = d.serf.handleNodeLeaveIntent(&leave)

	case messageJoinType:
		var join messageJoin
		if err := decodeMessage(buf[1:], &join); err != nil ***REMOVED***
			d.serf.logger.Printf("[ERR] serf: Error decoding join message: %s", err)
			break
		***REMOVED***

		d.serf.logger.Printf("[DEBUG] serf: messageJoinType: %s", join.Node)
		rebroadcast = d.serf.handleNodeJoinIntent(&join)

	case messageUserEventType:
		var event messageUserEvent
		if err := decodeMessage(buf[1:], &event); err != nil ***REMOVED***
			d.serf.logger.Printf("[ERR] serf: Error decoding user event message: %s", err)
			break
		***REMOVED***

		d.serf.logger.Printf("[DEBUG] serf: messageUserEventType: %s", event.Name)
		rebroadcast = d.serf.handleUserEvent(&event)
		rebroadcastQueue = d.serf.eventBroadcasts

	case messageQueryType:
		var query messageQuery
		if err := decodeMessage(buf[1:], &query); err != nil ***REMOVED***
			d.serf.logger.Printf("[ERR] serf: Error decoding query message: %s", err)
			break
		***REMOVED***

		d.serf.logger.Printf("[DEBUG] serf: messageQueryType: %s", query.Name)
		rebroadcast = d.serf.handleQuery(&query)
		rebroadcastQueue = d.serf.queryBroadcasts

	case messageQueryResponseType:
		var resp messageQueryResponse
		if err := decodeMessage(buf[1:], &resp); err != nil ***REMOVED***
			d.serf.logger.Printf("[ERR] serf: Error decoding query response message: %s", err)
			break
		***REMOVED***

		d.serf.logger.Printf("[DEBUG] serf: messageQueryResponseType: %v", resp.From)
		d.serf.handleQueryResponse(&resp)

	default:
		d.serf.logger.Printf("[WARN] serf: Received message of unknown type: %d", t)
	***REMOVED***

	if rebroadcast ***REMOVED***
		// Copy the buffer since it we cannot rely on the slice not changing
		newBuf := make([]byte, len(buf))
		copy(newBuf, buf)

		rebroadcastQueue.QueueBroadcast(&broadcast***REMOVED***
			msg:    newBuf,
			notify: nil,
		***REMOVED***)
	***REMOVED***
***REMOVED***

func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte ***REMOVED***
	msgs := d.serf.broadcasts.GetBroadcasts(overhead, limit)

	// Determine the bytes used already
	bytesUsed := 0
	for _, msg := range msgs ***REMOVED***
		lm := len(msg)
		bytesUsed += lm + overhead
		metrics.AddSample([]string***REMOVED***"serf", "msgs", "sent"***REMOVED***, float32(lm))
	***REMOVED***

	// Get any additional query broadcasts
	queryMsgs := d.serf.queryBroadcasts.GetBroadcasts(overhead, limit-bytesUsed)
	if queryMsgs != nil ***REMOVED***
		for _, m := range queryMsgs ***REMOVED***
			lm := len(m)
			bytesUsed += lm + overhead
			metrics.AddSample([]string***REMOVED***"serf", "msgs", "sent"***REMOVED***, float32(lm))
		***REMOVED***
		msgs = append(msgs, queryMsgs...)
	***REMOVED***

	// Get any additional event broadcasts
	eventMsgs := d.serf.eventBroadcasts.GetBroadcasts(overhead, limit-bytesUsed)
	if eventMsgs != nil ***REMOVED***
		for _, m := range eventMsgs ***REMOVED***
			lm := len(m)
			bytesUsed += lm + overhead
			metrics.AddSample([]string***REMOVED***"serf", "msgs", "sent"***REMOVED***, float32(lm))
		***REMOVED***
		msgs = append(msgs, eventMsgs...)
	***REMOVED***

	return msgs
***REMOVED***

func (d *delegate) LocalState(join bool) []byte ***REMOVED***
	d.serf.memberLock.RLock()
	defer d.serf.memberLock.RUnlock()
	d.serf.eventLock.RLock()
	defer d.serf.eventLock.RUnlock()

	// Create the message to send
	pp := messagePushPull***REMOVED***
		LTime:        d.serf.clock.Time(),
		StatusLTimes: make(map[string]LamportTime, len(d.serf.members)),
		LeftMembers:  make([]string, 0, len(d.serf.leftMembers)),
		EventLTime:   d.serf.eventClock.Time(),
		Events:       d.serf.eventBuffer,
		QueryLTime:   d.serf.queryClock.Time(),
	***REMOVED***

	// Add all the join LTimes
	for name, member := range d.serf.members ***REMOVED***
		pp.StatusLTimes[name] = member.statusLTime
	***REMOVED***

	// Add all the left nodes
	for _, member := range d.serf.leftMembers ***REMOVED***
		pp.LeftMembers = append(pp.LeftMembers, member.Name)
	***REMOVED***

	// Encode the push pull state
	buf, err := encodeMessage(messagePushPullType, &pp)
	if err != nil ***REMOVED***
		d.serf.logger.Printf("[ERR] serf: Failed to encode local state: %v", err)
		return nil
	***REMOVED***
	return buf
***REMOVED***

func (d *delegate) MergeRemoteState(buf []byte, isJoin bool) ***REMOVED***
	// Ensure we have a message
	if len(buf) == 0 ***REMOVED***
		d.serf.logger.Printf("[ERR] serf: Remote state is zero bytes")
		return
	***REMOVED***

	// Check the message type
	if messageType(buf[0]) != messagePushPullType ***REMOVED***
		d.serf.logger.Printf("[ERR] serf: Remote state has bad type prefix: %v", buf[0])
		return
	***REMOVED***

	// Attempt a decode
	pp := messagePushPull***REMOVED******REMOVED***
	if err := decodeMessage(buf[1:], &pp); err != nil ***REMOVED***
		d.serf.logger.Printf("[ERR] serf: Failed to decode remote state: %v", err)
		return
	***REMOVED***

	// Witness the Lamport clocks first.
	// We subtract 1 since no message with that clock has been sent yet
	if pp.LTime > 0 ***REMOVED***
		d.serf.clock.Witness(pp.LTime - 1)
	***REMOVED***
	if pp.EventLTime > 0 ***REMOVED***
		d.serf.eventClock.Witness(pp.EventLTime - 1)
	***REMOVED***
	if pp.QueryLTime > 0 ***REMOVED***
		d.serf.queryClock.Witness(pp.QueryLTime - 1)
	***REMOVED***

	// Process the left nodes first to avoid the LTimes from being increment
	// in the wrong order
	leftMap := make(map[string]struct***REMOVED******REMOVED***, len(pp.LeftMembers))
	leave := messageLeave***REMOVED******REMOVED***
	for _, name := range pp.LeftMembers ***REMOVED***
		leftMap[name] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		leave.LTime = pp.StatusLTimes[name]
		leave.Node = name
		d.serf.handleNodeLeaveIntent(&leave)
	***REMOVED***

	// Update any other LTimes
	join := messageJoin***REMOVED******REMOVED***
	for name, statusLTime := range pp.StatusLTimes ***REMOVED***
		// Skip the left nodes
		if _, ok := leftMap[name]; ok ***REMOVED***
			continue
		***REMOVED***

		// Create an artificial join message
		join.LTime = statusLTime
		join.Node = name
		d.serf.handleNodeJoinIntent(&join)
	***REMOVED***

	// If we are doing a join, and eventJoinIgnore is set
	// then we set the eventMinTime to the EventLTime. This
	// prevents any of the incoming events from being processed
	if isJoin && d.serf.eventJoinIgnore ***REMOVED***
		d.serf.eventLock.Lock()
		if pp.EventLTime > d.serf.eventMinTime ***REMOVED***
			d.serf.eventMinTime = pp.EventLTime
		***REMOVED***
		d.serf.eventLock.Unlock()
	***REMOVED***

	// Process all the events
	userEvent := messageUserEvent***REMOVED******REMOVED***
	for _, events := range pp.Events ***REMOVED***
		if events == nil ***REMOVED***
			continue
		***REMOVED***
		userEvent.LTime = events.LTime
		for _, e := range events.Events ***REMOVED***
			userEvent.Name = e.Name
			userEvent.Payload = e.Payload
			d.serf.handleUserEvent(&userEvent)
		***REMOVED***
	***REMOVED***
***REMOVED***
