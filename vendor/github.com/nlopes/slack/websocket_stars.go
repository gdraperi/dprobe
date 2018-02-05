package slack

type starEvent struct ***REMOVED***
	Type           string      `json:"type"`
	User           string      `json:"user"`
	Item           StarredItem `json:"item"`
	EventTimestamp string      `json:"event_ts"`
***REMOVED***

// StarAddedEvent represents the Star added event
type StarAddedEvent starEvent

// StarRemovedEvent represents the Star removed event
type StarRemovedEvent starEvent
