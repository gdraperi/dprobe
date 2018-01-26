// +build freebsd darwin

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
// For example if current binary is "docker" at "/usr/bin/", then cmd.Path will
// be set to "/usr/bin/docker".
func Command(args ...string) *exec.Cmd ***REMOVED***
	return &exec.Cmd***REMOVED***
		Path: Self(),
		Args: args,
	***REMOVED***
***REMOVED***
