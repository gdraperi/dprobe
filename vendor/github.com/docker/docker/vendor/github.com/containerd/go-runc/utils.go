package runc

import (
	"io/ioutil"
	"strconv"
	"syscall"
)

// ReadPidFile reads the pid file at the provided path and returns
// the pid or an error if the read and conversion is unsuccessful
func ReadPidFile(path string) (int, error) ***REMOVED***
	data, err := ioutil.ReadFile(path)
	if err != nil ***REMOVED***
		return -1, err
	***REMOVED***
	return strconv.Atoi(string(data))
***REMOVED***

const exitSignalOffset = 128

// exitStatus returns the correct exit status for a process based on if it
// was signaled or exited cleanly
func exitStatus(status syscall.WaitStatus) int ***REMOVED***
	if status.Signaled() ***REMOVED***
		return exitSignalOffset + int(status.Signal())
	***REMOVED***
	return status.ExitStatus()
***REMOVED***
