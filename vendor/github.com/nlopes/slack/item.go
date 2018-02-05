package slack

const (
	TYPE_MESSAGE      = "message"
	TYPE_FILE         = "file"
	TYPE_FILE_COMMENT = "file_comment"
	TYPE_CHANNEL      = "channel"
	TYPE_IM           = "im"
	TYPE_GROUP        = "group"
)

// Item is any type of slack message - message, file, or file comment.
type Item struct ***REMOVED***
	Type      string   `json:"type"`
	Channel   string   `json:"channel,omitempty"`
	Message   *Message `json:"message,omitempty"`
	File      *File    `json:"file,omitempty"`
	Comment   *Comment `json:"comment,omitempty"`
	Timestamp string   `json:"ts,omitempty"`
***REMOVED***

// NewMessageItem turns a message on a channel into a typed message struct.
func NewMessageItem(ch string, m *Message) Item ***REMOVED***
	return Item***REMOVED***Type: TYPE_MESSAGE, Channel: ch, Message: m***REMOVED***
***REMOVED***

// NewFileItem turns a file into a typed file struct.
func NewFileItem(f *File) Item ***REMOVED***
	return Item***REMOVED***Type: TYPE_FILE, File: f***REMOVED***
***REMOVED***

// NewFileCommentItem turns a file and comment into a typed file_comment struct.
func NewFileCommentItem(f *File, c *Comment) Item ***REMOVED***
	return Item***REMOVED***Type: TYPE_FILE_COMMENT, File: f, Comment: c***REMOVED***
***REMOVED***

// NewChannelItem turns a channel id into a typed channel struct.
func NewChannelItem(ch string) Item ***REMOVED***
	return Item***REMOVED***Type: TYPE_CHANNEL, Channel: ch***REMOVED***
***REMOVED***

// NewIMItem turns a channel id into a typed im struct.
func NewIMItem(ch string) Item ***REMOVED***
	return Item***REMOVED***Type: TYPE_IM, Channel: ch***REMOVED***
***REMOVED***

// NewGroupItem turns a channel id into a typed group struct.
func NewGroupItem(ch string) Item ***REMOVED***
	return Item***REMOVED***Type: TYPE_GROUP, Channel: ch***REMOVED***
***REMOVED***

// ItemRef is a reference to a message of any type. One of FileID,
// CommentId, or the combination of ChannelId and Timestamp must be
// specified.
type ItemRef struct ***REMOVED***
	Channel   string `json:"channel"`
	Timestamp string `json:"timestamp"`
	File      string `json:"file"`
	Comment   string `json:"file_comment"`
***REMOVED***

// NewRefToMessage initializes a reference to to a message.
func NewRefToMessage(channel, timestamp string) ItemRef ***REMOVED***
	return ItemRef***REMOVED***Channel: channel, Timestamp: timestamp***REMOVED***
***REMOVED***

// NewRefToFile initializes a reference to a file.
func NewRefToFile(file string) ItemRef ***REMOVED***
	return ItemRef***REMOVED***File: file***REMOVED***
***REMOVED***

// NewRefToComment initializes a reference to a file comment.
func NewRefToComment(comment string) ItemRef ***REMOVED***
	return ItemRef***REMOVED***Comment: comment***REMOVED***
***REMOVED***
