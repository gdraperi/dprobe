package main

import systemdDaemon "github.com/coreos/go-systemd/daemon"

// preNotifySystem sends a message to the host when the API is active, but before the daemon is
func preNotifySystem() ***REMOVED***
***REMOVED***

// notifySystem sends a message to the host when the server is ready to be used
func notifySystem() ***REMOVED***
	// Tell the init daemon we are accepting requests
	go systemdDaemon.SdNotify(false, "READY=1")
***REMOVED***
