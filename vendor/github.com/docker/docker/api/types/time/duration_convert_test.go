package time

import (
	"testing"
	"time"
)

func TestDurationToSecondsString(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		in       time.Duration
		expected string
	***REMOVED******REMOVED***
		***REMOVED***0 * time.Second, "0"***REMOVED***,
		***REMOVED***1 * time.Second, "1"***REMOVED***,
		***REMOVED***1 * time.Minute, "60"***REMOVED***,
		***REMOVED***24 * time.Hour, "86400"***REMOVED***,
	***REMOVED***

	for _, c := range cases ***REMOVED***
		s := DurationToSecondsString(c.in)
		if s != c.expected ***REMOVED***
			t.Errorf("wrong value for input `%v`: expected `%s`, got `%s`", c.in, c.expected, s)
			t.Fail()
		***REMOVED***
	***REMOVED***
***REMOVED***
