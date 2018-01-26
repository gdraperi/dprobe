package dbus

import (
	"errors"
	"os/exec"
)

func sessionBusPlatform() (*Conn, error) ***REMOVED***
	cmd := exec.Command("launchctl", "getenv", "DBUS_LAUNCHD_SESSION_BUS_SOCKET")
	b, err := cmd.CombinedOutput()

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if len(b) == 0 ***REMOVED***
		return nil, errors.New("dbus: couldn't determine address of session bus")
	***REMOVED***

	return Dial("unix:path=" + string(b[:len(b)-1]))
***REMOVED***
