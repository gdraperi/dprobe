package cgroups

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// V1 returns all the groups in the default cgroups mountpoint in a single hierarchy
func V1() ([]Subsystem, error) ***REMOVED***
	root, err := v1MountPoint()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	subsystems, err := defaults(root)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var enabled []Subsystem
	for _, s := range pathers(subsystems) ***REMOVED***
		// check and remove the default groups that do not exist
		if _, err := os.Lstat(s.Path("/")); err == nil ***REMOVED***
			enabled = append(enabled, s)
		***REMOVED***
	***REMOVED***
	return enabled, nil
***REMOVED***

// v1MountPoint returns the mount point where the cgroup
// mountpoints are mounted in a single hiearchy
func v1MountPoint() (string, error) ***REMOVED***
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() ***REMOVED***
		if err := scanner.Err(); err != nil ***REMOVED***
			return "", err
		***REMOVED***
		var (
			text   = scanner.Text()
			fields = strings.Split(text, " ")
			// safe as mountinfo encodes mountpoints with spaces as \040.
			index               = strings.Index(text, " - ")
			postSeparatorFields = strings.Fields(text[index+3:])
			numPostFields       = len(postSeparatorFields)
		)
		// this is an error as we can't detect if the mount is for "cgroup"
		if numPostFields == 0 ***REMOVED***
			return "", fmt.Errorf("Found no fields post '-' in %q", text)
		***REMOVED***
		if postSeparatorFields[0] == "cgroup" ***REMOVED***
			// check that the mount is properly formated.
			if numPostFields < 3 ***REMOVED***
				return "", fmt.Errorf("Error found less than 3 fields post '-' in %q", text)
			***REMOVED***
			return filepath.Dir(fields[4]), nil
		***REMOVED***
	***REMOVED***
	return "", ErrMountPointNotExist
***REMOVED***
