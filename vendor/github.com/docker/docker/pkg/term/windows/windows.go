// These files implement ANSI-aware input and output streams for use by the Docker Windows client.
// When asked for the set of standard streams (e.g., stdin, stdout, stderr), the code will create
// and return pseudo-streams that convert ANSI sequences to / from Windows Console API calls.

package windowsconsole

import (
	"io/ioutil"
	"os"
	"sync"

	ansiterm "github.com/Azure/go-ansiterm"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger
var initOnce sync.Once

func initLogger() ***REMOVED***
	initOnce.Do(func() ***REMOVED***
		logFile := ioutil.Discard

		if isDebugEnv := os.Getenv(ansiterm.LogEnv); isDebugEnv == "1" ***REMOVED***
			logFile, _ = os.Create("ansiReaderWriter.log")
		***REMOVED***

		logger = &logrus.Logger***REMOVED***
			Out:       logFile,
			Formatter: new(logrus.TextFormatter),
			Level:     logrus.DebugLevel,
		***REMOVED***
	***REMOVED***)
***REMOVED***
