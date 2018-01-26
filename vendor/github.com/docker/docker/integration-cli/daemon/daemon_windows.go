package daemon

import (
	"fmt"
	"strconv"

	"github.com/go-check/check"
	"golang.org/x/sys/windows"
)

// SignalDaemonDump sends a signal to the daemon to write a dump file
func SignalDaemonDump(pid int) ***REMOVED***
	ev, _ := windows.UTF16PtrFromString("Global\\docker-daemon-" + strconv.Itoa(pid))
	h2, err := windows.OpenEvent(0x0002, false, ev)
	if h2 == 0 || err != nil ***REMOVED***
		return
	***REMOVED***
	windows.PulseEvent(h2)
***REMOVED***

func signalDaemonReload(pid int) error ***REMOVED***
	return fmt.Errorf("daemon reload not supported")
***REMOVED***

func cleanupExecRoot(c *check.C, execRoot string) ***REMOVED***
***REMOVED***
