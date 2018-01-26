// +build !windows

package sys

import "golang.org/x/sys/unix"

// Exit is the wait4 information from an exited process
type Exit struct ***REMOVED***
	Pid    int
	Status int
***REMOVED***

// Reap reaps all child processes for the calling process and returns their
// exit information
func Reap(wait bool) (exits []Exit, err error) ***REMOVED***
	var (
		ws  unix.WaitStatus
		rus unix.Rusage
	)
	flag := unix.WNOHANG
	if wait ***REMOVED***
		flag = 0
	***REMOVED***
	for ***REMOVED***
		pid, err := unix.Wait4(-1, &ws, flag, &rus)
		if err != nil ***REMOVED***
			if err == unix.ECHILD ***REMOVED***
				return exits, nil
			***REMOVED***
			return exits, err
		***REMOVED***
		if pid <= 0 ***REMOVED***
			return exits, nil
		***REMOVED***
		exits = append(exits, Exit***REMOVED***
			Pid:    pid,
			Status: exitStatus(ws),
		***REMOVED***)
	***REMOVED***
***REMOVED***

const exitSignalOffset = 128

// exitStatus returns the correct exit status for a process based on if it
// was signaled or exited cleanly
func exitStatus(status unix.WaitStatus) int ***REMOVED***
	if status.Signaled() ***REMOVED***
		return exitSignalOffset + int(status.Signal())
	***REMOVED***
	return status.ExitStatus()
***REMOVED***
