package utils

import (
	"io"
	"net"
	"os"
	"syscall"

	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/streamformatter"
	"github.com/sirupsen/logrus"
)

// WriteDistributionProgress is a helper for writing progress from chan to JSON
// stream with an optional cancel function.
func WriteDistributionProgress(cancelFunc func(), outStream io.Writer, progressChan <-chan progress.Progress) ***REMOVED***
	progressOutput := streamformatter.NewJSONProgressOutput(outStream, false)
	operationCancelled := false

	for prog := range progressChan ***REMOVED***
		if err := progressOutput.WriteProgress(prog); err != nil && !operationCancelled ***REMOVED***
			// don't log broken pipe errors as this is the normal case when a client aborts
			if isBrokenPipe(err) ***REMOVED***
				logrus.Info("Pull session cancelled")
			***REMOVED*** else ***REMOVED***
				logrus.Errorf("error writing progress to client: %v", err)
			***REMOVED***
			cancelFunc()
			operationCancelled = true
			// Don't return, because we need to continue draining
			// progressChan until it's closed to avoid a deadlock.
		***REMOVED***
	***REMOVED***
***REMOVED***

func isBrokenPipe(e error) bool ***REMOVED***
	if netErr, ok := e.(*net.OpError); ok ***REMOVED***
		e = netErr.Err
		if sysErr, ok := netErr.Err.(*os.SyscallError); ok ***REMOVED***
			e = sysErr.Err
		***REMOVED***
	***REMOVED***
	return e == syscall.EPIPE
***REMOVED***
