package slack

// TeamJoinEvent represents the Team join event
type TeamJoinEvent struct ***REMOVED***
	Type string `json:"type"`
	User User   `json:"user"`
***REMOVED***

// TeamRenameEvent represents the Team rename event
type TeamRenameEvent struct ***REMOVED***
	Type           string `json:"type"`
	Name           string `json:"name,omitempty"`
	EventTimestamp string `json:"event_ts,omitempty"`
***REMOVED***

// TeamPrefChangeEvent represents the Team preference change event
type TeamPrefChangeEvent struct ***REMOVED***
	Type  string   `json:"type"`
	Name  string   `json:"name,omitempty"`
	Value []string `json:"value,omitempty"`
***REMOVED***

// TeamDomainChangeEvent represents the Team domain change event
type TeamDomainChangeEvent struct ***REMOVED***
	Type   string `json:"type"`
	URL    string `json:"url"`
	Domain string `json:"domain"`
***REMOVED***

// TeamMigrationStartedEvent represents the Team migration started event
type TeamMigrationStartedEvent struct ***REMOVED***
	Type string `json:"type"`
***REMOVED***
