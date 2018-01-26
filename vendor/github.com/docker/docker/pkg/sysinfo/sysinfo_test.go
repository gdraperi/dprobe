package sysinfo

import "testing"

func TestIsCpusetListAvailable(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		provided  string
		available string
		res       bool
		err       bool
	***REMOVED******REMOVED***
		***REMOVED***"1", "0-4", true, false***REMOVED***,
		***REMOVED***"01,3", "0-4", true, false***REMOVED***,
		***REMOVED***"", "0-7", true, false***REMOVED***,
		***REMOVED***"1--42", "0-7", false, true***REMOVED***,
		***REMOVED***"1-42", "00-1,8,,9", false, true***REMOVED***,
		***REMOVED***"1,41-42", "43,45", false, false***REMOVED***,
		***REMOVED***"0-3", "", false, false***REMOVED***,
	***REMOVED***
	for _, c := range cases ***REMOVED***
		r, err := isCpusetListAvailable(c.provided, c.available)
		if (c.err && err == nil) && r != c.res ***REMOVED***
			t.Fatalf("Expected pair: %v, %v for %s, %s. Got %v, %v instead", c.res, c.err, c.provided, c.available, (c.err && err == nil), r)
		***REMOVED***
	***REMOVED***
***REMOVED***
