package system

import "os"

// IsProcessAlive returns true if process with a given pid is running.
func IsProcessAlive(pid int) bool ***REMOVED***
	_, err := os.FindProcess(pid)

	return err == nil
***REMOVED***

// KillProcess force-stops a process.
func KillProcess(pid int) ***REMOVED***
	p, err := os.FindProcess(pid)
	if err == nil ***REMOVED***
		p.Kill()
	***REMOVED***
***REMOVED***
