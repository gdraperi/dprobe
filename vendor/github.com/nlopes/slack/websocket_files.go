package slack

// FileActionEvent represents the File action event
type fileActionEvent struct ***REMOVED***
	Type           string `json:"type"`
	EventTimestamp string `json:"event_ts"`
	File           File   `json:"file"`
	// FileID is used for FileDeletedEvent
	FileID string `json:"file_id,omitempty"`
***REMOVED***

// FileCreatedEvent represents the File created event
type FileCreatedEvent fileActionEvent

// FileSharedEvent represents the File shared event
type FileSharedEvent fileActionEvent

// FilePublicEvent represents the File public event
type FilePublicEvent fileActionEvent

// FileUnsharedEvent represents the File unshared event
type FileUnsharedEvent fileActionEvent

// FileChangeEvent represents the File change event
type FileChangeEvent fileActionEvent

// FileDeletedEvent represents the File deleted event
type FileDeletedEvent fileActionEvent

// FilePrivateEvent represents the File private event
type FilePrivateEvent fileActionEvent

// FileCommentAddedEvent represents the File comment added event
type FileCommentAddedEvent struct ***REMOVED***
	fileActionEvent
	Comment Comment `json:"comment"`
***REMOVED***

// FileCommentEditedEvent represents the File comment edited event
type FileCommentEditedEvent struct ***REMOVED***
	fileActionEvent
	Comment Comment `json:"comment"`
***REMOVED***

// FileCommentDeletedEvent represents the File comment deleted event
type FileCommentDeletedEvent struct ***REMOVED***
	fileActionEvent
	Comment string `json:"comment"`
***REMOVED***
