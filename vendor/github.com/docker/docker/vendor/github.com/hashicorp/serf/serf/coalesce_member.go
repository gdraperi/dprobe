package serf

type coalesceEvent struct ***REMOVED***
	Type   EventType
	Member *Member
***REMOVED***

type memberEventCoalescer struct ***REMOVED***
	lastEvents   map[string]EventType
	latestEvents map[string]coalesceEvent
***REMOVED***

func (c *memberEventCoalescer) Handle(e Event) bool ***REMOVED***
	switch e.EventType() ***REMOVED***
	case EventMemberJoin:
		return true
	case EventMemberLeave:
		return true
	case EventMemberFailed:
		return true
	case EventMemberUpdate:
		return true
	case EventMemberReap:
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

func (c *memberEventCoalescer) Coalesce(raw Event) ***REMOVED***
	e := raw.(MemberEvent)
	for _, m := range e.Members ***REMOVED***
		c.latestEvents[m.Name] = coalesceEvent***REMOVED***
			Type:   e.Type,
			Member: &m,
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *memberEventCoalescer) Flush(outCh chan<- Event) ***REMOVED***
	// Coalesce the various events we got into a single set of events.
	events := make(map[EventType]*MemberEvent)
	for name, cevent := range c.latestEvents ***REMOVED***
		previous, ok := c.lastEvents[name]

		// If we sent the same event before, then ignore
		// unless it is a MemberUpdate
		if ok && previous == cevent.Type && cevent.Type != EventMemberUpdate ***REMOVED***
			continue
		***REMOVED***

		// Update our last event
		c.lastEvents[name] = cevent.Type

		// Add it to our event
		newEvent, ok := events[cevent.Type]
		if !ok ***REMOVED***
			newEvent = &MemberEvent***REMOVED***Type: cevent.Type***REMOVED***
			events[cevent.Type] = newEvent
		***REMOVED***
		newEvent.Members = append(newEvent.Members, *cevent.Member)
	***REMOVED***

	// Send out those events
	for _, event := range events ***REMOVED***
		outCh <- *event
	***REMOVED***
***REMOVED***
