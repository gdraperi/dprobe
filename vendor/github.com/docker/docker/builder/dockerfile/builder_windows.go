package dockerfile

func defaultShellForOS(os string) []string ***REMOVED***
	if os == "linux" ***REMOVED***
		return []string***REMOVED***"/bin/sh", "-c"***REMOVED***
	***REMOVED***
	return []string***REMOVED***"cmd", "/S", "/C"***REMOVED***
***REMOVED***
