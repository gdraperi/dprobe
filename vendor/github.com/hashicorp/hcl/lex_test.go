package hcl

import (
	"testing"
)

func TestLexMode(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		Input string
		Mode  lexModeValue
	***REMOVED******REMOVED***
		***REMOVED***
			"",
			lexModeHcl,
		***REMOVED***,
		***REMOVED***
			"foo",
			lexModeHcl,
		***REMOVED***,
		***REMOVED***
			"***REMOVED******REMOVED***",
			lexModeJson,
		***REMOVED***,
		***REMOVED***
			"  ***REMOVED******REMOVED***",
			lexModeJson,
		***REMOVED***,
	***REMOVED***

	for i, tc := range cases ***REMOVED***
		actual := lexMode([]byte(tc.Input))

		if actual != tc.Mode ***REMOVED***
			t.Fatalf("%d: %#v", i, actual)
		***REMOVED***
	***REMOVED***
***REMOVED***
