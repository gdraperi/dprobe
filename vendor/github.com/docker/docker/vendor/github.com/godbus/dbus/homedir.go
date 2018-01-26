package dbus

import (
	"os"
	"sync"
)

var (
	homeDir     string
	homeDirLock sync.Mutex
)

func getHomeDir() string ***REMOVED***
	homeDirLock.Lock()
	defer homeDirLock.Unlock()

	if homeDir != "" ***REMOVED***
		return homeDir
	***REMOVED***

	homeDir = os.Getenv("HOME")
	if homeDir != "" ***REMOVED***
		return homeDir
	***REMOVED***

	homeDir = lookupHomeDir()
	return homeDir
***REMOVED***
