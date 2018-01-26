// +build static_build

package dbus

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func lookupHomeDir() string ***REMOVED***
	myUid := os.Getuid()

	f, err := os.Open("/etc/passwd")
	if err != nil ***REMOVED***
		return "/"
	***REMOVED***
	defer f.Close()

	s := bufio.NewScanner(f)

	for s.Scan() ***REMOVED***
		if err := s.Err(); err != nil ***REMOVED***
			break
		***REMOVED***

		line := strings.TrimSpace(s.Text())
		if line == "" ***REMOVED***
			continue
		***REMOVED***

		parts := strings.Split(line, ":")

		if len(parts) >= 6 ***REMOVED***
			uid, err := strconv.Atoi(parts[2])
			if err == nil && uid == myUid ***REMOVED***
				return parts[5]
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Default to / if we can't get a better value
	return "/"
***REMOVED***
