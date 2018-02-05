package slack

import (
	"fmt"
	"sync"
)

// SetLogger let's library users supply a logger, so that api debugging
// can be logged along with the application's debugging info.
func SetLogger(l logProvider) ***REMOVED***
	loggerMutex.Lock()
	logger = ilogger***REMOVED***logProvider: l***REMOVED***
	loggerMutex.Unlock()
***REMOVED***

var (
	loggerMutex = new(sync.Mutex)
	logger      logInternal // A logger that can be set by consumers
)

// logProvider is a logger interface compatible with both stdlib and some
// 3rd party loggers such as logrus.
type logProvider interface ***REMOVED***
	Output(int, string) error
***REMOVED***

// logInternal represents the internal logging api we use.
type logInternal interface ***REMOVED***
	Print(...interface***REMOVED******REMOVED***)
	Printf(string, ...interface***REMOVED******REMOVED***)
	Println(...interface***REMOVED******REMOVED***)
	Output(int, string) error
***REMOVED***

// ilogger implements the additional methods used by our internal logging.
type ilogger struct ***REMOVED***
	logProvider
***REMOVED***

// Println replicates the behaviour of the standard logger.
func (t ilogger) Println(v ...interface***REMOVED******REMOVED***) ***REMOVED***
	t.Output(2, fmt.Sprintln(v...))
***REMOVED***

// Printf replicates the behaviour of the standard logger.
func (t ilogger) Printf(format string, v ...interface***REMOVED******REMOVED***) ***REMOVED***
	t.Output(2, fmt.Sprintf(format, v...))
***REMOVED***

// Print replicates the behaviour of the standard logger.
func (t ilogger) Print(v ...interface***REMOVED******REMOVED***) ***REMOVED***
	t.Output(2, fmt.Sprint(v...))
***REMOVED***
