package system

import (
	"fmt"
	"os/exec"
	"syscall"
)

// GetExitCode returns the ExitStatus of the specified error if its type is
// exec.ExitError, returns 0 and an error otherwise.
func GetExitCode(err error) (int, error) ***REMOVED***
	exitCode := 0
	if exiterr, ok := err.(*exec.ExitError); ok ***REMOVED***
		if procExit, ok := exiterr.Sys().(syscall.WaitStatus); ok ***REMOVED***
			return procExit.ExitStatus(), nil
		***REMOVED***
	***REMOVED***
	return exitCode, fmt.Errorf("failed to get exit code")
***REMOVED***
