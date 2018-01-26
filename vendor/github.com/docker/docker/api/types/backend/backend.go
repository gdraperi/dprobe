// Package backend includes types to send information to server backends.
package backend

import (
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

// ContainerAttachConfig holds the streams to use when connecting to a container to view logs.
type ContainerAttachConfig struct ***REMOVED***
	GetStreams func() (io.ReadCloser, io.Writer, io.Writer, error)
	UseStdin   bool
	UseStdout  bool
	UseStderr  bool
	Logs       bool
	Stream     bool
	DetachKeys string

	// Used to signify that streams are multiplexed and therefore need a StdWriter to encode stdout/stderr messages accordingly.
	// TODO @cpuguy83: This shouldn't be needed. It was only added so that http and websocket endpoints can use the same function, and the websocket function was not using a stdwriter prior to this change...
	// HOWEVER, the websocket endpoint is using a single stream and SHOULD be encoded with stdout/stderr as is done for HTTP since it is still just a single stream.
	// Since such a change is an API change unrelated to the current changeset we'll keep it as is here and change separately.
	MuxStreams bool
***REMOVED***

// LogMessage is datastructure that represents piece of output produced by some
// container.  The Line member is a slice of an array whose contents can be
// changed after a log driver's Log() method returns.
// changes to this struct need to be reflect in the reset method in
// daemon/logger/logger.go
type LogMessage struct ***REMOVED***
	Line      []byte
	Source    string
	Timestamp time.Time
	Attrs     []LogAttr
	Partial   bool

	// Err is an error associated with a message. Completeness of a message
	// with Err is not expected, tho it may be partially complete (fields may
	// be missing, gibberish, or nil)
	Err error
***REMOVED***

// LogAttr is used to hold the extra attributes available in the log message.
type LogAttr struct ***REMOVED***
	Key   string
	Value string
***REMOVED***

// LogSelector is a list of services and tasks that should be returned as part
// of a log stream. It is similar to swarmapi.LogSelector, with the difference
// that the names don't have to be resolved to IDs; this is mostly to avoid
// accidents later where a swarmapi LogSelector might have been incorrectly
// used verbatim (and to avoid the handler having to import swarmapi types)
type LogSelector struct ***REMOVED***
	Services []string
	Tasks    []string
***REMOVED***

// ContainerStatsConfig holds information for configuring the runtime
// behavior of a backend.ContainerStats() call.
type ContainerStatsConfig struct ***REMOVED***
	Stream    bool
	OutStream io.Writer
	Version   string
***REMOVED***

// ExecInspect holds information about a running process started
// with docker exec.
type ExecInspect struct ***REMOVED***
	ID            string
	Running       bool
	ExitCode      *int
	ProcessConfig *ExecProcessConfig
	OpenStdin     bool
	OpenStderr    bool
	OpenStdout    bool
	CanRemove     bool
	ContainerID   string
	DetachKeys    []byte
	Pid           int
***REMOVED***

// ExecProcessConfig holds information about the exec process
// running on the host.
type ExecProcessConfig struct ***REMOVED***
	Tty        bool     `json:"tty"`
	Entrypoint string   `json:"entrypoint"`
	Arguments  []string `json:"arguments"`
	Privileged *bool    `json:"privileged,omitempty"`
	User       string   `json:"user,omitempty"`
***REMOVED***

// ContainerCommitConfig is a wrapper around
// types.ContainerCommitConfig that also
// transports configuration changes for a container.
type ContainerCommitConfig struct ***REMOVED***
	types.ContainerCommitConfig
	Changes []string
	// TODO: ContainerConfig is only used by the dockerfile Builder, so remove it
	// once the Builder has been updated to use a different interface
	ContainerConfig *container.Config
***REMOVED***
