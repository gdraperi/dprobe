package slack

// ChannelCreatedEvent represents the Channel created event
type ChannelCreatedEvent struct ***REMOVED***
	Type           string             `json:"type"`
	Channel        ChannelCreatedInfo `json:"channel"`
	EventTimestamp string             `json:"event_ts"`
***REMOVED***

// ChannelCreatedInfo represents the information associated with the Channel created event
type ChannelCreatedInfo struct ***REMOVED***
	ID        string `json:"id"`
	IsChannel bool   `json:"is_channel"`
	Name      string `json:"name"`
	Created   int    `json:"created"`
	Creator   string `json:"creator"`
***REMOVED***

// ChannelJoinedEvent represents the Channel joined event
type ChannelJoinedEvent struct ***REMOVED***
	Type    string  `json:"type"`
	Channel Channel `json:"channel"`
***REMOVED***

// ChannelInfoEvent represents the Channel info event
type ChannelInfoEvent struct ***REMOVED***
	// channel_left
	// channel_deleted
	// channel_archive
	// channel_unarchive
	Type      string `json:"type"`
	Channel   string `json:"channel"`
	User      string `json:"user,omitempty"`
	Timestamp string `json:"ts,omitempty"`
***REMOVED***

// ChannelRenameEvent represents the Channel rename event
type ChannelRenameEvent struct ***REMOVED***
	Type      string            `json:"type"`
	Channel   ChannelRenameInfo `json:"channel"`
	Timestamp string            `json:"event_ts"`
***REMOVED***

// ChannelRenameInfo represents the information associated with a Channel rename event
type ChannelRenameInfo struct ***REMOVED***
	ID      string `json:"id"`
	Name    string `json:"name"`
	Created string `json:"created"`
***REMOVED***

// ChannelHistoryChangedEvent represents the Channel history changed event
type ChannelHistoryChangedEvent struct ***REMOVED***
	Type           string `json:"type"`
	Latest         string `json:"latest"`
	Timestamp      string `json:"ts"`
	EventTimestamp string `json:"event_ts"`
***REMOVED***

// ChannelMarkedEvent represents the Channel marked event
type ChannelMarkedEvent ChannelInfoEvent

// ChannelLeftEvent represents the Channel left event
type ChannelLeftEvent ChannelInfoEvent

// ChannelDeletedEvent represents the Channel deleted event
type ChannelDeletedEvent ChannelInfoEvent

// ChannelArchiveEvent represents the Channel archive event
type ChannelArchiveEvent ChannelInfoEvent

// ChannelUnarchiveEvent represents the Channel unarchive event
type ChannelUnarchiveEvent ChannelInfoEvent
