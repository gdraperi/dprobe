package icmd

import (
	"os/exec"
	"syscall"

	"github.com/pkg/errors"
)

// getExitCode returns the ExitStatus of a process from the error returned by
// exec.Run(). If the exit status could not be parsed an error is returned.
func getExitCode(err error) (int, error) ***REMOVED***
	if exiterr, ok := err.(*exec.ExitError); ok ***REMOVED***
		if procExit, ok := exiterr.Sys().(syscall.WaitStatus); ok ***REMOVED***
			return procExit.ExitStatus(), nil
		***REMOVED***
	***REMOVED***
	return 0, errors.Wrap(err, "failed to get exit code")
***REMOVED***

func processExitCode(err error) (exitCode int) ***REMOVED***
	if err == nil ***REMOVED***
		return 0
	***REMOVED***
	exitCode, exiterr := getExitCode(err)
	if exiterr != nil ***REMOVED***
		// TODO: Fix this so we check the error's text.
		// we've failed to retrieve exit code, so we set it to 127
		return 127
	***REMOVED***
	return exitCode
***REMOVED***
