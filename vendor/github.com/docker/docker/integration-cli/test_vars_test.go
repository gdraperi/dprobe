package main

// sleepCommandForDaemonPlatform is a helper function that determines what
// the command is for a sleeping container based on the daemon platform.
// The Windows busybox image does not have a `top` command.
func sleepCommandForDaemonPlatform() []string ***REMOVED***
	if testEnv.OSType == "windows" ***REMOVED***
		return []string***REMOVED***"sleep", "240"***REMOVED***
	***REMOVED***
	return []string***REMOVED***"top"***REMOVED***
***REMOVED***
