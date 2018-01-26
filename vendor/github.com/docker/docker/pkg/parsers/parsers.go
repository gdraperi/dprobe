// Package parsers provides helper functions to parse and validate different type
// of string. It can be hosts, unix addresses, tcp addresses, filters, kernel
// operating system versions.
package parsers

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseKeyValueOpt parses and validates the specified string as a key/value pair (key=value)
func ParseKeyValueOpt(opt string) (string, string, error) ***REMOVED***
	parts := strings.SplitN(opt, "=", 2)
	if len(parts) != 2 ***REMOVED***
		return "", "", fmt.Errorf("Unable to parse key/value option: %s", opt)
	***REMOVED***
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
***REMOVED***

// ParseUintList parses and validates the specified string as the value
// found in some cgroup file (e.g. `cpuset.cpus`, `cpuset.mems`), which could be
// one of the formats below. Note that duplicates are actually allowed in the
// input string. It returns a `map[int]bool` with available elements from `val`
// set to `true`.
// Supported formats:
//     7
//     1-6
//     0,3-4,7,8-10
//     0-0,0,1-7
//     03,1-3      <- this is gonna get parsed as [1,2,3]
//     3,2,1
//     0-2,3,1
func ParseUintList(val string) (map[int]bool, error) ***REMOVED***
	if val == "" ***REMOVED***
		return map[int]bool***REMOVED******REMOVED***, nil
	***REMOVED***

	availableInts := make(map[int]bool)
	split := strings.Split(val, ",")
	errInvalidFormat := fmt.Errorf("invalid format: %s", val)

	for _, r := range split ***REMOVED***
		if !strings.Contains(r, "-") ***REMOVED***
			v, err := strconv.Atoi(r)
			if err != nil ***REMOVED***
				return nil, errInvalidFormat
			***REMOVED***
			availableInts[v] = true
		***REMOVED*** else ***REMOVED***
			split := strings.SplitN(r, "-", 2)
			min, err := strconv.Atoi(split[0])
			if err != nil ***REMOVED***
				return nil, errInvalidFormat
			***REMOVED***
			max, err := strconv.Atoi(split[1])
			if err != nil ***REMOVED***
				return nil, errInvalidFormat
			***REMOVED***
			if max < min ***REMOVED***
				return nil, errInvalidFormat
			***REMOVED***
			for i := min; i <= max; i++ ***REMOVED***
				availableInts[i] = true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return availableInts, nil
***REMOVED***
