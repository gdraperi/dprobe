// +build freebsd darwin

package operatingsystem

import (
	"errors"
	"os/exec"
)

// GetOperatingSystem gets the name of the current operating system.
func GetOperatingSystem() (string, error) ***REMOVED***
	cmd := exec.Command("uname", "-s")
	osName, err := cmd.Output()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return string(osName), nil
***REMOVED***

// IsContainerized returns true if we are running inside a container.
// No-op on FreeBSD and Darwin, always returns false.
func IsContainerized() (bool, error) ***REMOVED***
	// TODO: Implement jail detection for freeBSD
	return false, errors.New("Cannot detect if we are in container")
***REMOVED***
