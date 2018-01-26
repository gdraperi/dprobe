package idtools

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// add a user and/or group to Linux /etc/passwd, /etc/group using standard
// Linux distribution commands:
// adduser --system --shell /bin/false --disabled-login --disabled-password --no-create-home --group <username>
// useradd -r -s /bin/false <username>

var (
	once        sync.Once
	userCommand string

	cmdTemplates = map[string]string***REMOVED***
		"adduser": "--system --shell /bin/false --no-create-home --disabled-login --disabled-password --group %s",
		"useradd": "-r -s /bin/false %s",
		"usermod": "-%s %d-%d %s",
	***REMOVED***

	idOutRegexp = regexp.MustCompile(`uid=([0-9]+).*gid=([0-9]+)`)
	// default length for a UID/GID subordinate range
	defaultRangeLen   = 65536
	defaultRangeStart = 100000
	userMod           = "usermod"
)

// AddNamespaceRangesUser takes a username and uses the standard system
// utility to create a system user/group pair used to hold the
// /etc/sub***REMOVED***uid,gid***REMOVED*** ranges which will be used for user namespace
// mapping ranges in containers.
func AddNamespaceRangesUser(name string) (int, int, error) ***REMOVED***
	if err := addUser(name); err != nil ***REMOVED***
		return -1, -1, fmt.Errorf("Error adding user %q: %v", name, err)
	***REMOVED***

	// Query the system for the created uid and gid pair
	out, err := execCmd("id", name)
	if err != nil ***REMOVED***
		return -1, -1, fmt.Errorf("Error trying to find uid/gid for new user %q: %v", name, err)
	***REMOVED***
	matches := idOutRegexp.FindStringSubmatch(strings.TrimSpace(string(out)))
	if len(matches) != 3 ***REMOVED***
		return -1, -1, fmt.Errorf("Can't find uid, gid from `id` output: %q", string(out))
	***REMOVED***
	uid, err := strconv.Atoi(matches[1])
	if err != nil ***REMOVED***
		return -1, -1, fmt.Errorf("Can't convert found uid (%s) to int: %v", matches[1], err)
	***REMOVED***
	gid, err := strconv.Atoi(matches[2])
	if err != nil ***REMOVED***
		return -1, -1, fmt.Errorf("Can't convert found gid (%s) to int: %v", matches[2], err)
	***REMOVED***

	// Now we need to create the subuid/subgid ranges for our new user/group (system users
	// do not get auto-created ranges in subuid/subgid)

	if err := createSubordinateRanges(name); err != nil ***REMOVED***
		return -1, -1, fmt.Errorf("Couldn't create subordinate ID ranges: %v", err)
	***REMOVED***
	return uid, gid, nil
***REMOVED***

func addUser(userName string) error ***REMOVED***
	once.Do(func() ***REMOVED***
		// set up which commands are used for adding users/groups dependent on distro
		if _, err := resolveBinary("adduser"); err == nil ***REMOVED***
			userCommand = "adduser"
		***REMOVED*** else if _, err := resolveBinary("useradd"); err == nil ***REMOVED***
			userCommand = "useradd"
		***REMOVED***
	***REMOVED***)
	if userCommand == "" ***REMOVED***
		return fmt.Errorf("Cannot add user; no useradd/adduser binary found")
	***REMOVED***
	args := fmt.Sprintf(cmdTemplates[userCommand], userName)
	out, err := execCmd(userCommand, args)
	if err != nil ***REMOVED***
		return fmt.Errorf("Failed to add user with error: %v; output: %q", err, string(out))
	***REMOVED***
	return nil
***REMOVED***

func createSubordinateRanges(name string) error ***REMOVED***

	// first, we should verify that ranges weren't automatically created
	// by the distro tooling
	ranges, err := parseSubuid(name)
	if err != nil ***REMOVED***
		return fmt.Errorf("Error while looking for subuid ranges for user %q: %v", name, err)
	***REMOVED***
	if len(ranges) == 0 ***REMOVED***
		// no UID ranges; let's create one
		startID, err := findNextUIDRange()
		if err != nil ***REMOVED***
			return fmt.Errorf("Can't find available subuid range: %v", err)
		***REMOVED***
		out, err := execCmd(userMod, fmt.Sprintf(cmdTemplates[userMod], "v", startID, startID+defaultRangeLen-1, name))
		if err != nil ***REMOVED***
			return fmt.Errorf("Unable to add subuid range to user: %q; output: %s, err: %v", name, out, err)
		***REMOVED***
	***REMOVED***

	ranges, err = parseSubgid(name)
	if err != nil ***REMOVED***
		return fmt.Errorf("Error while looking for subgid ranges for user %q: %v", name, err)
	***REMOVED***
	if len(ranges) == 0 ***REMOVED***
		// no GID ranges; let's create one
		startID, err := findNextGIDRange()
		if err != nil ***REMOVED***
			return fmt.Errorf("Can't find available subgid range: %v", err)
		***REMOVED***
		out, err := execCmd(userMod, fmt.Sprintf(cmdTemplates[userMod], "w", startID, startID+defaultRangeLen-1, name))
		if err != nil ***REMOVED***
			return fmt.Errorf("Unable to add subgid range to user: %q; output: %s, err: %v", name, out, err)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func findNextUIDRange() (int, error) ***REMOVED***
	ranges, err := parseSubuid("ALL")
	if err != nil ***REMOVED***
		return -1, fmt.Errorf("Couldn't parse all ranges in /etc/subuid file: %v", err)
	***REMOVED***
	sort.Sort(ranges)
	return findNextRangeStart(ranges)
***REMOVED***

func findNextGIDRange() (int, error) ***REMOVED***
	ranges, err := parseSubgid("ALL")
	if err != nil ***REMOVED***
		return -1, fmt.Errorf("Couldn't parse all ranges in /etc/subgid file: %v", err)
	***REMOVED***
	sort.Sort(ranges)
	return findNextRangeStart(ranges)
***REMOVED***

func findNextRangeStart(rangeList ranges) (int, error) ***REMOVED***
	startID := defaultRangeStart
	for _, arange := range rangeList ***REMOVED***
		if wouldOverlap(arange, startID) ***REMOVED***
			startID = arange.Start + arange.Length
		***REMOVED***
	***REMOVED***
	return startID, nil
***REMOVED***

func wouldOverlap(arange subIDRange, ID int) bool ***REMOVED***
	low := ID
	high := ID + defaultRangeLen
	if (low >= arange.Start && low <= arange.Start+arange.Length) ||
		(high <= arange.Start+arange.Length && high >= arange.Start) ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***
