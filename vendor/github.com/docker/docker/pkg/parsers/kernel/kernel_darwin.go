// +build darwin

// Package kernel provides helper function to get, parse and compare kernel
// versions for different platforms.
package kernel

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/mattn/go-shellwords"
)

// GetKernelVersion gets the current kernel version.
func GetKernelVersion() (*VersionInfo, error) ***REMOVED***
	release, err := getRelease()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return ParseRelease(release)
***REMOVED***

// getRelease uses `system_profiler SPSoftwareDataType` to get OSX kernel version
func getRelease() (string, error) ***REMOVED***
	cmd := exec.Command("system_profiler", "SPSoftwareDataType")
	osName, err := cmd.Output()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	var release string
	data := strings.Split(string(osName), "\n")
	for _, line := range data ***REMOVED***
		if strings.Contains(line, "Kernel Version") ***REMOVED***
			// It has the format like '      Kernel Version: Darwin 14.5.0'
			content := strings.SplitN(line, ":", 2)
			if len(content) != 2 ***REMOVED***
				return "", fmt.Errorf("Kernel Version is invalid")
			***REMOVED***

			prettyNames, err := shellwords.Parse(content[1])
			if err != nil ***REMOVED***
				return "", fmt.Errorf("Kernel Version is invalid: %s", err.Error())
			***REMOVED***

			if len(prettyNames) != 2 ***REMOVED***
				return "", fmt.Errorf("Kernel Version needs to be 'Darwin x.x.x' ")
			***REMOVED***
			release = prettyNames[1]
		***REMOVED***
	***REMOVED***

	return release, nil
***REMOVED***
