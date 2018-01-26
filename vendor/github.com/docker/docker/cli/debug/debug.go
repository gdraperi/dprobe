package debug

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Enable sets the DEBUG env var to true
// and makes the logger to log at debug level.
func Enable() ***REMOVED***
	os.Setenv("DEBUG", "1")
	logrus.SetLevel(logrus.DebugLevel)
***REMOVED***

// Disable sets the DEBUG env var to false
// and makes the logger to log at info level.
func Disable() ***REMOVED***
	os.Setenv("DEBUG", "")
	logrus.SetLevel(logrus.InfoLevel)
***REMOVED***

// IsEnabled checks whether the debug flag is set or not.
func IsEnabled() bool ***REMOVED***
	return os.Getenv("DEBUG") != ""
***REMOVED***
