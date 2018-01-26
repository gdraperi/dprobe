// Package logger defines interfaces that logger drivers implement to
// log messages.
//
// The other half of a logger driver is the implementation of the
// factory, which holds the contextual instance information that
// allows multiple loggers of the same type to perform different
// actions, such as logging to different locations.
package logger

import (
	"sync"
	"time"

	"github.com/docker/docker/api/types/backend"
)

// ErrReadLogsNotSupported is returned when the underlying log driver does not support reading
type ErrReadLogsNotSupported struct***REMOVED******REMOVED***

func (ErrReadLogsNotSupported) Error() string ***REMOVED***
	return "configured logging driver does not support reading"
***REMOVED***

// NotImplemented makes this error implement the `NotImplemented` interface from api/errdefs
func (ErrReadLogsNotSupported) NotImplemented() ***REMOVED******REMOVED***

const (
	logWatcherBufferSize = 4096
)

var messagePool = &sync.Pool***REMOVED***New: func() interface***REMOVED******REMOVED*** ***REMOVED*** return &Message***REMOVED***Line: make([]byte, 0, 256)***REMOVED*** ***REMOVED******REMOVED***

// NewMessage returns a new message from the message sync.Pool
func NewMessage() *Message ***REMOVED***
	return messagePool.Get().(*Message)
***REMOVED***

// PutMessage puts the specified message back n the message pool.
// The message fields are reset before putting into the pool.
func PutMessage(msg *Message) ***REMOVED***
	msg.reset()
	messagePool.Put(msg)
***REMOVED***

// Message is datastructure that represents piece of output produced by some
// container.  The Line member is a slice of an array whose contents can be
// changed after a log driver's Log() method returns.
//
// Message is subtyped from backend.LogMessage because there is a lot of
// internal complexity around the Message type that should not be exposed
// to any package not explicitly importing the logger type.
//
// Any changes made to this struct must also be updated in the `reset` function
type Message backend.LogMessage

// reset sets the message back to default values
// This is used when putting a message back into the message pool.
// Any changes to the `Message` struct should be reflected here.
func (m *Message) reset() ***REMOVED***
	m.Line = m.Line[:0]
	m.Source = ""
	m.Attrs = nil
	m.Partial = false

	m.Err = nil
***REMOVED***

// AsLogMessage returns a pointer to the message as a pointer to
// backend.LogMessage, which is an identical type with a different purpose
func (m *Message) AsLogMessage() *backend.LogMessage ***REMOVED***
	return (*backend.LogMessage)(m)
***REMOVED***

// Logger is the interface for docker logging drivers.
type Logger interface ***REMOVED***
	Log(*Message) error
	Name() string
	Close() error
***REMOVED***

// SizedLogger is the interface for logging drivers that can control
// the size of buffer used for their messages.
type SizedLogger interface ***REMOVED***
	Logger
	BufSize() int
***REMOVED***

// ReadConfig is the configuration passed into ReadLogs.
type ReadConfig struct ***REMOVED***
	Since  time.Time
	Until  time.Time
	Tail   int
	Follow bool
***REMOVED***

// LogReader is the interface for reading log messages for loggers that support reading.
type LogReader interface ***REMOVED***
	// Read logs from underlying logging backend
	ReadLogs(ReadConfig) *LogWatcher
***REMOVED***

// LogWatcher is used when consuming logs read from the LogReader interface.
type LogWatcher struct ***REMOVED***
	// For sending log messages to a reader.
	Msg chan *Message
	// For sending error messages that occur while while reading logs.
	Err           chan error
	closeOnce     sync.Once
	closeNotifier chan struct***REMOVED******REMOVED***
***REMOVED***

// NewLogWatcher returns a new LogWatcher.
func NewLogWatcher() *LogWatcher ***REMOVED***
	return &LogWatcher***REMOVED***
		Msg:           make(chan *Message, logWatcherBufferSize),
		Err:           make(chan error, 1),
		closeNotifier: make(chan struct***REMOVED******REMOVED***),
	***REMOVED***
***REMOVED***

// Close notifies the underlying log reader to stop.
func (w *LogWatcher) Close() ***REMOVED***
	// only close if not already closed
	w.closeOnce.Do(func() ***REMOVED***
		close(w.closeNotifier)
	***REMOVED***)
***REMOVED***

// WatchClose returns a channel receiver that receives notification
// when the watcher has been closed. This should only be called from
// one goroutine.
func (w *LogWatcher) WatchClose() <-chan struct***REMOVED******REMOVED*** ***REMOVED***
	return w.closeNotifier
***REMOVED***

// Capability defines the list of capabilties that a driver can implement
// These capabilities are not required to be a logging driver, however do
// determine how a logging driver can be used
type Capability struct ***REMOVED***
	// Determines if a log driver can read back logs
	ReadLogs bool
***REMOVED***

// MarshalFunc is a func that marshals a message into an arbitrary format
type MarshalFunc func(*Message) ([]byte, error)
