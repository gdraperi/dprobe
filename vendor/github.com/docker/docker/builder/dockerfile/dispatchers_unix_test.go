// +build !windows

package dockerfile

import (
	"runtime"
	"testing"
)

func TestNormalizeWorkdir(t *testing.T) ***REMOVED***
	testCases := []struct***REMOVED*** current, requested, expected, expectedError string ***REMOVED******REMOVED***
		***REMOVED***``, ``, ``, `cannot normalize nothing`***REMOVED***,
		***REMOVED***``, `foo`, `/foo`, ``***REMOVED***,
		***REMOVED***``, `/foo`, `/foo`, ``***REMOVED***,
		***REMOVED***`/foo`, `bar`, `/foo/bar`, ``***REMOVED***,
		***REMOVED***`/foo`, `/bar`, `/bar`, ``***REMOVED***,
	***REMOVED***

	for _, test := range testCases ***REMOVED***
		normalized, err := normalizeWorkdir(runtime.GOOS, test.current, test.requested)

		if test.expectedError != "" && err == nil ***REMOVED***
			t.Fatalf("NormalizeWorkdir should return an error %s, got nil", test.expectedError)
		***REMOVED***

		if test.expectedError != "" && err.Error() != test.expectedError ***REMOVED***
			t.Fatalf("NormalizeWorkdir returned wrong error. Expected %s, got %s", test.expectedError, err.Error())
		***REMOVED***

		if normalized != test.expected ***REMOVED***
			t.Fatalf("NormalizeWorkdir error. Expected %s for current %s and requested %s, got %s", test.expected, test.current, test.requested, normalized)
		***REMOVED***
	***REMOVED***
***REMOVED***
