//+build !windows

package daemon

import (
	"testing"
)

func TestContainerTopValidatePSArgs(t *testing.T) ***REMOVED***
	tests := map[string]bool***REMOVED***
		"ae -o uid=PID":             true,
		"ae -o \"uid= PID\"":        true,  // ascii space (0x20)
		"ae -o \"uid= PID\"":        false, // unicode space (U+2003, 0xe2 0x80 0x83)
		"ae o uid=PID":              true,
		"aeo uid=PID":               true,
		"ae -O uid=PID":             true,
		"ae -o pid=PID2 -o uid=PID": true,
		"ae -o pid=PID":             false,
		"ae -o pid=PID -o uid=PIDX": true, // FIXME: we do not need to prohibit this
		"aeo pid=PID":               false,
		"ae":                        false,
		"":                          false,
	***REMOVED***
	for psArgs, errExpected := range tests ***REMOVED***
		err := validatePSArgs(psArgs)
		t.Logf("tested %q, got err=%v", psArgs, err)
		if errExpected && err == nil ***REMOVED***
			t.Fatalf("expected error, got %v (%q)", err, psArgs)
		***REMOVED***
		if !errExpected && err != nil ***REMOVED***
			t.Fatalf("expected nil, got %v (%q)", err, psArgs)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestContainerTopParsePSOutput(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		output      []byte
		pids        []uint32
		errExpected bool
	***REMOVED******REMOVED***
		***REMOVED***[]byte(`  PID COMMAND
   42 foo
   43 bar
		- -
  100 baz
`), []uint32***REMOVED***42, 43***REMOVED***, false***REMOVED***,
		***REMOVED***[]byte(`  UID COMMAND
   42 foo
   43 bar
		- -
  100 baz
`), []uint32***REMOVED***42, 43***REMOVED***, true***REMOVED***,
		// unicode space (U+2003, 0xe2 0x80 0x83)
		***REMOVED***[]byte(` PID COMMAND
   42 foo
   43 bar
		- -
  100 baz
`), []uint32***REMOVED***42, 43***REMOVED***, true***REMOVED***,
		// the first space is U+2003, the second one is ascii.
		***REMOVED***[]byte(` PID COMMAND
   42 foo
   43 bar
  100 baz
`), []uint32***REMOVED***42, 43***REMOVED***, true***REMOVED***,
	***REMOVED***

	for _, f := range tests ***REMOVED***
		_, err := parsePSOutput(f.output, f.pids)
		t.Logf("tested %q, got err=%v", string(f.output), err)
		if f.errExpected && err == nil ***REMOVED***
			t.Fatalf("expected error, got %v (%q)", err, string(f.output))
		***REMOVED***
		if !f.errExpected && err != nil ***REMOVED***
			t.Fatalf("expected nil, got %v (%q)", err, string(f.output))
		***REMOVED***
	***REMOVED***
***REMOVED***
