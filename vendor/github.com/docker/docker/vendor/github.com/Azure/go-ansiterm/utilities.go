package ansiterm

import (
	"strconv"
)

func sliceContains(bytes []byte, b byte) bool ***REMOVED***
	for _, v := range bytes ***REMOVED***
		if v == b ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

func convertBytesToInteger(bytes []byte) int ***REMOVED***
	s := string(bytes)
	i, _ := strconv.Atoi(s)
	return i
***REMOVED***
