package daemon

import (
	"runtime"
	"testing"

	"github.com/docker/docker/volume"
)

func TestParseVolumesFrom(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		spec    string
		expID   string
		expMode string
		fail    bool
	***REMOVED******REMOVED***
		***REMOVED***"", "", "", true***REMOVED***,
		***REMOVED***"foobar", "foobar", "rw", false***REMOVED***,
		***REMOVED***"foobar:rw", "foobar", "rw", false***REMOVED***,
		***REMOVED***"foobar:ro", "foobar", "ro", false***REMOVED***,
		***REMOVED***"foobar:baz", "", "", true***REMOVED***,
	***REMOVED***

	parser := volume.NewParser(runtime.GOOS)

	for _, c := range cases ***REMOVED***
		id, mode, err := parser.ParseVolumesFrom(c.spec)
		if c.fail ***REMOVED***
			if err == nil ***REMOVED***
				t.Fatalf("Expected error, was nil, for spec %s\n", c.spec)
			***REMOVED***
			continue
		***REMOVED***

		if id != c.expID ***REMOVED***
			t.Fatalf("Expected id %s, was %s, for spec %s\n", c.expID, id, c.spec)
		***REMOVED***
		if mode != c.expMode ***REMOVED***
			t.Fatalf("Expected mode %s, was %s for spec %s\n", c.expMode, mode, c.spec)
		***REMOVED***
	***REMOVED***
***REMOVED***
