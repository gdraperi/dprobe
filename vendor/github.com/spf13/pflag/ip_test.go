package pflag

import (
	"fmt"
	"net"
	"os"
	"testing"
)

func setUpIP(ip *net.IP) *FlagSet ***REMOVED***
	f := NewFlagSet("test", ContinueOnError)
	f.IPVar(ip, "address", net.ParseIP("0.0.0.0"), "IP Address")
	return f
***REMOVED***

func TestIP(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		input    string
		success  bool
		expected string
	***REMOVED******REMOVED***
		***REMOVED***"0.0.0.0", true, "0.0.0.0"***REMOVED***,
		***REMOVED***" 0.0.0.0 ", true, "0.0.0.0"***REMOVED***,
		***REMOVED***"1.2.3.4", true, "1.2.3.4"***REMOVED***,
		***REMOVED***"127.0.0.1", true, "127.0.0.1"***REMOVED***,
		***REMOVED***"255.255.255.255", true, "255.255.255.255"***REMOVED***,
		***REMOVED***"", false, ""***REMOVED***,
		***REMOVED***"0", false, ""***REMOVED***,
		***REMOVED***"localhost", false, ""***REMOVED***,
		***REMOVED***"0.0.0", false, ""***REMOVED***,
		***REMOVED***"0.0.0.", false, ""***REMOVED***,
		***REMOVED***"0.0.0.0.", false, ""***REMOVED***,
		***REMOVED***"0.0.0.256", false, ""***REMOVED***,
		***REMOVED***"0 . 0 . 0 . 0", false, ""***REMOVED***,
	***REMOVED***

	devnull, _ := os.Open(os.DevNull)
	os.Stderr = devnull
	for i := range testCases ***REMOVED***
		var addr net.IP
		f := setUpIP(&addr)

		tc := &testCases[i]

		arg := fmt.Sprintf("--address=%s", tc.input)
		err := f.Parse([]string***REMOVED***arg***REMOVED***)
		if err != nil && tc.success == true ***REMOVED***
			t.Errorf("expected success, got %q", err)
			continue
		***REMOVED*** else if err == nil && tc.success == false ***REMOVED***
			t.Errorf("expected failure")
			continue
		***REMOVED*** else if tc.success ***REMOVED***
			ip, err := f.GetIP("address")
			if err != nil ***REMOVED***
				t.Errorf("Got error trying to fetch the IP flag: %v", err)
			***REMOVED***
			if ip.String() != tc.expected ***REMOVED***
				t.Errorf("expected %q, got %q", tc.expected, ip.String())
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
