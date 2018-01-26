// Package operatingsystem provides helper function to get the operating system
// name for different platforms.
package operatingsystem

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mattn/go-shellwords"
)

var (
	// file to use to detect if the daemon is running in a container
	proc1Cgroup = "/proc/1/cgroup"

	// file to check to determine Operating System
	etcOsRelease = "/etc/os-release"

	// used by stateless systems like Clear Linux
	altOsRelease = "/usr/lib/os-release"
)

// GetOperatingSystem gets the name of the current operating system.
func GetOperatingSystem() (string, error) ***REMOVED***
	osReleaseFile, err := os.Open(etcOsRelease)
	if err != nil ***REMOVED***
		if !os.IsNotExist(err) ***REMOVED***
			return "", fmt.Errorf("Error opening %s: %v", etcOsRelease, err)
		***REMOVED***
		osReleaseFile, err = os.Open(altOsRelease)
		if err != nil ***REMOVED***
			return "", fmt.Errorf("Error opening %s: %v", altOsRelease, err)
		***REMOVED***
	***REMOVED***
	defer osReleaseFile.Close()

	var prettyName string
	scanner := bufio.NewScanner(osReleaseFile)
	for scanner.Scan() ***REMOVED***
		line := scanner.Text()
		if strings.HasPrefix(line, "PRETTY_NAME=") ***REMOVED***
			data := strings.SplitN(line, "=", 2)
			prettyNames, err := shellwords.Parse(data[1])
			if err != nil ***REMOVED***
				return "", fmt.Errorf("PRETTY_NAME is invalid: %s", err.Error())
			***REMOVED***
			if len(prettyNames) != 1 ***REMOVED***
				return "", fmt.Errorf("PRETTY_NAME needs to be enclosed by quotes if they have spaces: %s", data[1])
			***REMOVED***
			prettyName = prettyNames[0]
		***REMOVED***
	***REMOVED***
	if prettyName != "" ***REMOVED***
		return prettyName, nil
	***REMOVED***
	// If not set, defaults to PRETTY_NAME="Linux"
	// c.f. http://www.freedesktop.org/software/systemd/man/os-release.html
	return "Linux", nil
***REMOVED***

// IsContainerized returns true if we are running inside a container.
func IsContainerized() (bool, error) ***REMOVED***
	b, err := ioutil.ReadFile(proc1Cgroup)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	for _, line := range bytes.Split(b, []byte***REMOVED***'\n'***REMOVED***) ***REMOVED***
		if len(line) > 0 && !bytes.HasSuffix(line, []byte***REMOVED***'/'***REMOVED***) && !bytes.HasSuffix(line, []byte("init.scope")) ***REMOVED***
			return true, nil
		***REMOVED***
	***REMOVED***
	return false, nil
***REMOVED***
