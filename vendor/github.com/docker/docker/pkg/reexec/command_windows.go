package reexec

import (
	"os/exec"
)

// Self returns the path to the current process's binary.
// Uses os.Args[0].
func Self() string ***REMOVED***
	return naiveSelf()
***REMOVED***

// Command returns *exec.Cmd which has Path as current binary.
// For example if current binary is "docker.exe" at "C:\", then cmd.Path will
// be set to "C:\docker.exe".
func Command(args ...string) *exec.Cmd ***REMOVED***
	return &exec.Cmd***REMOVED***
		Path: Self(),
		Args: args,
	***REMOVED***
***REMOVED***
