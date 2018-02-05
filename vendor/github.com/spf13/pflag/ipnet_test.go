package pflag

import (
	"fmt"
	"net"
	"os"
	"testing"
)

func setUpIPNet(ip *net.IPNet) *FlagSet ***REMOVED***
	f := NewFlagSet("test", ContinueOnError)
	_, def, _ := net.ParseCIDR("0.0.0.0/0")
	f.IPNetVar(ip, "address", *def, "IP Address")
	return f
***REMOVED***

func TestIPNet(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		input    string
		success  bool
		expected string
	***REMOVED******REMOVED***
		***REMOVED***"0.0.0.0/0", true, "0.0.0.0/0"***REMOVED***,
		***REMOVED***" 0.0.0.0/0 ", true, "0.0.0.0/0"***REMOVED***,
		***REMOVED***"1.2.3.4/8", true, "1.0.0.0/8"***REMOVED***,
		***REMOVED***"127.0.0.1/16", true, "127.0.0.0/16"***REMOVED***,
		***REMOVED***"255.255.255.255/19", true, "255.255.224.0/19"***REMOVED***,
		***REMOVED***"255.255.255.255/32", true, "255.255.255.255/32"***REMOVED***,
		***REMOVED***"", false, ""***REMOVED***,
		***REMOVED***"/0", false, ""***REMOVED***,
		***REMOVED***"0", false, ""***REMOVED***,
		***REMOVED***"0/0", false, ""***REMOVED***,
		***REMOVED***"localhost/0", false, ""***REMOVED***,
		***REMOVED***"0.0.0/4", false, ""***REMOVED***,
		***REMOVED***"0.0.0./8", false, ""***REMOVED***,
		***REMOVED***"0.0.0.0./12", false, ""***REMOVED***,
		***REMOVED***"0.0.0.256/16", false, ""***REMOVED***,
		***REMOVED***"0.0.0.0 /20", false, ""***REMOVED***,
		***REMOVED***"0.0.0.0/ 24", false, ""***REMOVED***,
		***REMOVED***"0 . 0 . 0 . 0 / 28", false, ""***REMOVED***,
		***REMOVED***"0.0.0.0/33", false, ""***REMOVED***,
	***REMOVED***

	devnull, _ := os.Open(os.DevNull)
	os.Stderr = devnull
	for i := range testCases ***REMOVED***
		var addr net.IPNet
		f := setUpIPNet(&addr)

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
			ip, err := f.GetIPNet("address")
			if err != nil ***REMOVED***
				t.Errorf("Got error trying to fetch the IP flag: %v", err)
			***REMOVED***
			if ip.String() != tc.expected ***REMOVED***
				t.Errorf("expected %q, got %q", tc.expected, ip.String())
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
