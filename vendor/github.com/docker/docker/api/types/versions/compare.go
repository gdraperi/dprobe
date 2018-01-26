package versions

import (
	"strconv"
	"strings"
)

// compare compares two version strings
// returns -1 if v1 < v2, 1 if v1 > v2, 0 otherwise.
func compare(v1, v2 string) int ***REMOVED***
	var (
		currTab  = strings.Split(v1, ".")
		otherTab = strings.Split(v2, ".")
	)

	max := len(currTab)
	if len(otherTab) > max ***REMOVED***
		max = len(otherTab)
	***REMOVED***
	for i := 0; i < max; i++ ***REMOVED***
		var currInt, otherInt int

		if len(currTab) > i ***REMOVED***
			currInt, _ = strconv.Atoi(currTab[i])
		***REMOVED***
		if len(otherTab) > i ***REMOVED***
			otherInt, _ = strconv.Atoi(otherTab[i])
		***REMOVED***
		if currInt > otherInt ***REMOVED***
			return 1
		***REMOVED***
		if otherInt > currInt ***REMOVED***
			return -1
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

// LessThan checks if a version is less than another
func LessThan(v, other string) bool ***REMOVED***
	return compare(v, other) == -1
***REMOVED***

// LessThanOrEqualTo checks if a version is less than or equal to another
func LessThanOrEqualTo(v, other string) bool ***REMOVED***
	return compare(v, other) <= 0
***REMOVED***

// GreaterThan checks if a version is greater than another
func GreaterThan(v, other string) bool ***REMOVED***
	return compare(v, other) == 1
***REMOVED***

// GreaterThanOrEqualTo checks if a version is greater than or equal to another
func GreaterThanOrEqualTo(v, other string) bool ***REMOVED***
	return compare(v, other) >= 0
***REMOVED***

// Equal checks if a version is equal to another
func Equal(v, other string) bool ***REMOVED***
	return compare(v, other) == 0
***REMOVED***
