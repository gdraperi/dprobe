// Code generated by "stringer -type RoundingMode"; DO NOT EDIT.

package number

import "fmt"

const _RoundingMode_name = "ToNearestEvenToNearestZeroToNearestAwayToPositiveInfToNegativeInfToZeroAwayFromZeronumModes"

var _RoundingMode_index = [...]uint8***REMOVED***0, 13, 26, 39, 52, 65, 71, 83, 91***REMOVED***

func (i RoundingMode) String() string ***REMOVED***
	if i >= RoundingMode(len(_RoundingMode_index)-1) ***REMOVED***
		return fmt.Sprintf("RoundingMode(%d)", i)
	***REMOVED***
	return _RoundingMode_name[_RoundingMode_index[i]:_RoundingMode_index[i+1]]
***REMOVED***