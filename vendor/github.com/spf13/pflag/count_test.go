package pflag

import (
	"os"
	"testing"
)

func setUpCount(c *int) *FlagSet ***REMOVED***
	f := NewFlagSet("test", ContinueOnError)
	f.CountVarP(c, "verbose", "v", "a counter")
	return f
***REMOVED***

func TestCount(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		input    []string
		success  bool
		expected int
	***REMOVED******REMOVED***
		***REMOVED***[]string***REMOVED******REMOVED***, true, 0***REMOVED***,
		***REMOVED***[]string***REMOVED***"-v"***REMOVED***, true, 1***REMOVED***,
		***REMOVED***[]string***REMOVED***"-vvv"***REMOVED***, true, 3***REMOVED***,
		***REMOVED***[]string***REMOVED***"-v", "-v", "-v"***REMOVED***, true, 3***REMOVED***,
		***REMOVED***[]string***REMOVED***"-v", "--verbose", "-v"***REMOVED***, true, 3***REMOVED***,
		***REMOVED***[]string***REMOVED***"-v=3", "-v"***REMOVED***, true, 4***REMOVED***,
		***REMOVED***[]string***REMOVED***"--verbose=0"***REMOVED***, true, 0***REMOVED***,
		***REMOVED***[]string***REMOVED***"-v=0"***REMOVED***, true, 0***REMOVED***,
		***REMOVED***[]string***REMOVED***"-v=a"***REMOVED***, false, 0***REMOVED***,
	***REMOVED***

	devnull, _ := os.Open(os.DevNull)
	os.Stderr = devnull
	for i := range testCases ***REMOVED***
		var count int
		f := setUpCount(&count)

		tc := &testCases[i]

		err := f.Parse(tc.input)
		if err != nil && tc.success == true ***REMOVED***
			t.Errorf("expected success, got %q", err)
			continue
		***REMOVED*** else if err == nil && tc.success == false ***REMOVED***
			t.Errorf("expected failure, got success")
			continue
		***REMOVED*** else if tc.success ***REMOVED***
			c, err := f.GetCount("verbose")
			if err != nil ***REMOVED***
				t.Errorf("Got error trying to fetch the counter flag")
			***REMOVED***
			if c != tc.expected ***REMOVED***
				t.Errorf("expected %d, got %d", tc.expected, c)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
