package serf

type latestUserEvents struct ***REMOVED***
	LTime  LamportTime
	Events []Event
***REMOVED***

type userEventCoalescer struct ***REMOVED***
	// Maps an event name into the latest versions
	events map[string]*latestUserEvents
***REMOVED***

func (c *userEventCoalescer) Handle(e Event) bool ***REMOVED***
	// Only handle EventUser messages
	if e.EventType() != EventUser ***REMOVED***
		return false
	***REMOVED***

	// Check if coalescing is enabled
	user := e.(UserEvent)
	return user.Coalesce
***REMOVED***

func (c *userEventCoalescer) Coalesce(e Event) ***REMOVED***
	user := e.(UserEvent)
	latest, ok := c.events[user.Name]

	// Create a new entry if there are none, or
	// if this message has the newest LTime
	if !ok || latest.LTime < user.LTime ***REMOVED***
		latest = &latestUserEvents***REMOVED***
			LTime:  user.LTime,
			Events: []Event***REMOVED***e***REMOVED***,
		***REMOVED***
		c.events[user.Name] = latest
		return
	***REMOVED***

	// If the the same age, save it
	if latest.LTime == user.LTime ***REMOVED***
		latest.Events = append(latest.Events, e)
	***REMOVED***
***REMOVED***

func (c *userEventCoalescer) Flush(outChan chan<- Event) ***REMOVED***
	for _, latest := range c.events ***REMOVED***
		for _, e := range latest.Events ***REMOVED***
			outChan <- e
		***REMOVED***
	***REMOVED***
	c.events = make(map[string]*latestUserEvents)
***REMOVED***
