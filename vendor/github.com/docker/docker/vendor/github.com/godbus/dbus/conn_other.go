// +build !darwin

package dbus

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
)

func sessionBusPlatform() (*Conn, error) ***REMOVED***
	cmd := exec.Command("dbus-launch")
	b, err := cmd.CombinedOutput()

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	i := bytes.IndexByte(b, '=')
	j := bytes.IndexByte(b, '\n')

	if i == -1 || j == -1 ***REMOVED***
		return nil, errors.New("dbus: couldn't determine address of session bus")
	***REMOVED***

	env, addr := string(b[0:i]), string(b[i+1:j])
	os.Setenv(env, addr)

	return Dial(addr)
***REMOVED***
