package time

import (
	"fmt"
	"testing"
	"time"
)

func TestGetTimestamp(t *testing.T) ***REMOVED***
	now := time.Now().In(time.UTC)
	cases := []struct ***REMOVED***
		in, expected string
		expectedErr  bool
	***REMOVED******REMOVED***
		// Partial RFC3339 strings get parsed with second precision
		***REMOVED***"2006-01-02T15:04:05.999999999+07:00", "1136189045.999999999", false***REMOVED***,
		***REMOVED***"2006-01-02T15:04:05.999999999Z", "1136214245.999999999", false***REMOVED***,
		***REMOVED***"2006-01-02T15:04:05.999999999", "1136214245.999999999", false***REMOVED***,
		***REMOVED***"2006-01-02T15:04:05Z", "1136214245.000000000", false***REMOVED***,
		***REMOVED***"2006-01-02T15:04:05", "1136214245.000000000", false***REMOVED***,
		***REMOVED***"2006-01-02T15:04:0Z", "", true***REMOVED***,
		***REMOVED***"2006-01-02T15:04:0", "", true***REMOVED***,
		***REMOVED***"2006-01-02T15:04Z", "1136214240.000000000", false***REMOVED***,
		***REMOVED***"2006-01-02T15:04+00:00", "1136214240.000000000", false***REMOVED***,
		***REMOVED***"2006-01-02T15:04-00:00", "1136214240.000000000", false***REMOVED***,
		***REMOVED***"2006-01-02T15:04", "1136214240.000000000", false***REMOVED***,
		***REMOVED***"2006-01-02T15:0Z", "", true***REMOVED***,
		***REMOVED***"2006-01-02T15:0", "", true***REMOVED***,
		***REMOVED***"2006-01-02T15Z", "1136214000.000000000", false***REMOVED***,
		***REMOVED***"2006-01-02T15+00:00", "1136214000.000000000", false***REMOVED***,
		***REMOVED***"2006-01-02T15-00:00", "1136214000.000000000", false***REMOVED***,
		***REMOVED***"2006-01-02T15", "1136214000.000000000", false***REMOVED***,
		***REMOVED***"2006-01-02T1Z", "1136163600.000000000", false***REMOVED***,
		***REMOVED***"2006-01-02T1", "1136163600.000000000", false***REMOVED***,
		***REMOVED***"2006-01-02TZ", "", true***REMOVED***,
		***REMOVED***"2006-01-02T", "", true***REMOVED***,
		***REMOVED***"2006-01-02+00:00", "1136160000.000000000", false***REMOVED***,
		***REMOVED***"2006-01-02-00:00", "1136160000.000000000", false***REMOVED***,
		***REMOVED***"2006-01-02-00:01", "1136160060.000000000", false***REMOVED***,
		***REMOVED***"2006-01-02Z", "1136160000.000000000", false***REMOVED***,
		***REMOVED***"2006-01-02", "1136160000.000000000", false***REMOVED***,
		***REMOVED***"2015-05-13T20:39:09Z", "1431549549.000000000", false***REMOVED***,

		// unix timestamps returned as is
		***REMOVED***"1136073600", "1136073600", false***REMOVED***,
		***REMOVED***"1136073600.000000001", "1136073600.000000001", false***REMOVED***,
		// Durations
		***REMOVED***"1m", fmt.Sprintf("%d", now.Add(-1*time.Minute).Unix()), false***REMOVED***,
		***REMOVED***"1.5h", fmt.Sprintf("%d", now.Add(-90*time.Minute).Unix()), false***REMOVED***,
		***REMOVED***"1h30m", fmt.Sprintf("%d", now.Add(-90*time.Minute).Unix()), false***REMOVED***,

		// String fallback
		***REMOVED***"invalid", "invalid", false***REMOVED***,
	***REMOVED***

	for _, c := range cases ***REMOVED***
		o, err := GetTimestamp(c.in, now)
		if o != c.expected ||
			(err == nil && c.expectedErr) ||
			(err != nil && !c.expectedErr) ***REMOVED***
			t.Errorf("wrong value for '%s'. expected:'%s' got:'%s' with error: `%s`", c.in, c.expected, o, err)
			t.Fail()
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestParseTimestamps(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		in                        string
		def, expectedS, expectedN int64
		expectedErr               bool
	***REMOVED******REMOVED***
		// unix timestamps
		***REMOVED***"1136073600", 0, 1136073600, 0, false***REMOVED***,
		***REMOVED***"1136073600.000000001", 0, 1136073600, 1, false***REMOVED***,
		***REMOVED***"1136073600.0000000010", 0, 1136073600, 1, false***REMOVED***,
		***REMOVED***"1136073600.00000001", 0, 1136073600, 10, false***REMOVED***,
		***REMOVED***"foo.bar", 0, 0, 0, true***REMOVED***,
		***REMOVED***"1136073600.bar", 0, 1136073600, 0, true***REMOVED***,
		***REMOVED***"", -1, -1, 0, false***REMOVED***,
	***REMOVED***

	for _, c := range cases ***REMOVED***
		s, n, err := ParseTimestamps(c.in, c.def)
		if s != c.expectedS ||
			n != c.expectedN ||
			(err == nil && c.expectedErr) ||
			(err != nil && !c.expectedErr) ***REMOVED***
			t.Errorf("wrong values for input `%s` with default `%d` expected:'%d'seconds and `%d`nanosecond got:'%d'seconds and `%d`nanoseconds with error: `%s`", c.in, c.def, c.expectedS, c.expectedN, s, n, err)
			t.Fail()
		***REMOVED***
	***REMOVED***
***REMOVED***
