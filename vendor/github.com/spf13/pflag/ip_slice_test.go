package pflag

import (
	"fmt"
	"net"
	"strings"
	"testing"
)

func setUpIPSFlagSet(ipsp *[]net.IP) *FlagSet ***REMOVED***
	f := NewFlagSet("test", ContinueOnError)
	f.IPSliceVar(ipsp, "ips", []net.IP***REMOVED******REMOVED***, "Command separated list!")
	return f
***REMOVED***

func setUpIPSFlagSetWithDefault(ipsp *[]net.IP) *FlagSet ***REMOVED***
	f := NewFlagSet("test", ContinueOnError)
	f.IPSliceVar(ipsp, "ips",
		[]net.IP***REMOVED***
			net.ParseIP("192.168.1.1"),
			net.ParseIP("0:0:0:0:0:0:0:1"),
		***REMOVED***,
		"Command separated list!")
	return f
***REMOVED***

func TestEmptyIP(t *testing.T) ***REMOVED***
	var ips []net.IP
	f := setUpIPSFlagSet(&ips)
	err := f.Parse([]string***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***

	getIPS, err := f.GetIPSlice("ips")
	if err != nil ***REMOVED***
		t.Fatal("got an error from GetIPSlice():", err)
	***REMOVED***
	if len(getIPS) != 0 ***REMOVED***
		t.Fatalf("got ips %v with len=%d but expected length=0", getIPS, len(getIPS))
	***REMOVED***
***REMOVED***

func TestIPS(t *testing.T) ***REMOVED***
	var ips []net.IP
	f := setUpIPSFlagSet(&ips)

	vals := []string***REMOVED***"192.168.1.1", "10.0.0.1", "0:0:0:0:0:0:0:2"***REMOVED***
	arg := fmt.Sprintf("--ips=%s", strings.Join(vals, ","))
	err := f.Parse([]string***REMOVED***arg***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	for i, v := range ips ***REMOVED***
		if ip := net.ParseIP(vals[i]); ip == nil ***REMOVED***
			t.Fatalf("invalid string being converted to IP address: %s", vals[i])
		***REMOVED*** else if !ip.Equal(v) ***REMOVED***
			t.Fatalf("expected ips[%d] to be %s but got: %s from GetIPSlice", i, vals[i], v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestIPSDefault(t *testing.T) ***REMOVED***
	var ips []net.IP
	f := setUpIPSFlagSetWithDefault(&ips)

	vals := []string***REMOVED***"192.168.1.1", "0:0:0:0:0:0:0:1"***REMOVED***
	err := f.Parse([]string***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	for i, v := range ips ***REMOVED***
		if ip := net.ParseIP(vals[i]); ip == nil ***REMOVED***
			t.Fatalf("invalid string being converted to IP address: %s", vals[i])
		***REMOVED*** else if !ip.Equal(v) ***REMOVED***
			t.Fatalf("expected ips[%d] to be %s but got: %s", i, vals[i], v)
		***REMOVED***
	***REMOVED***

	getIPS, err := f.GetIPSlice("ips")
	if err != nil ***REMOVED***
		t.Fatal("got an error from GetIPSlice")
	***REMOVED***
	for i, v := range getIPS ***REMOVED***
		if ip := net.ParseIP(vals[i]); ip == nil ***REMOVED***
			t.Fatalf("invalid string being converted to IP address: %s", vals[i])
		***REMOVED*** else if !ip.Equal(v) ***REMOVED***
			t.Fatalf("expected ips[%d] to be %s but got: %s", i, vals[i], v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestIPSWithDefault(t *testing.T) ***REMOVED***
	var ips []net.IP
	f := setUpIPSFlagSetWithDefault(&ips)

	vals := []string***REMOVED***"192.168.1.1", "0:0:0:0:0:0:0:1"***REMOVED***
	arg := fmt.Sprintf("--ips=%s", strings.Join(vals, ","))
	err := f.Parse([]string***REMOVED***arg***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	for i, v := range ips ***REMOVED***
		if ip := net.ParseIP(vals[i]); ip == nil ***REMOVED***
			t.Fatalf("invalid string being converted to IP address: %s", vals[i])
		***REMOVED*** else if !ip.Equal(v) ***REMOVED***
			t.Fatalf("expected ips[%d] to be %s but got: %s", i, vals[i], v)
		***REMOVED***
	***REMOVED***

	getIPS, err := f.GetIPSlice("ips")
	if err != nil ***REMOVED***
		t.Fatal("got an error from GetIPSlice")
	***REMOVED***
	for i, v := range getIPS ***REMOVED***
		if ip := net.ParseIP(vals[i]); ip == nil ***REMOVED***
			t.Fatalf("invalid string being converted to IP address: %s", vals[i])
		***REMOVED*** else if !ip.Equal(v) ***REMOVED***
			t.Fatalf("expected ips[%d] to be %s but got: %s", i, vals[i], v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestIPSCalledTwice(t *testing.T) ***REMOVED***
	var ips []net.IP
	f := setUpIPSFlagSet(&ips)

	in := []string***REMOVED***"192.168.1.2,0:0:0:0:0:0:0:1", "10.0.0.1"***REMOVED***
	expected := []net.IP***REMOVED***net.ParseIP("192.168.1.2"), net.ParseIP("0:0:0:0:0:0:0:1"), net.ParseIP("10.0.0.1")***REMOVED***
	argfmt := "ips=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string***REMOVED***arg1, arg2***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("expected no error; got", err)
	***REMOVED***
	for i, v := range ips ***REMOVED***
		if !expected[i].Equal(v) ***REMOVED***
			t.Fatalf("expected ips[%d] to be %s but got: %s", i, expected[i], v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestIPSBadQuoting(t *testing.T) ***REMOVED***

	tests := []struct ***REMOVED***
		Want    []net.IP
		FlagArg []string
	***REMOVED******REMOVED***
		***REMOVED***
			Want: []net.IP***REMOVED***
				net.ParseIP("a4ab:61d:f03e:5d7d:fad7:d4c2:a1a5:568"),
				net.ParseIP("203.107.49.208"),
				net.ParseIP("14.57.204.90"),
			***REMOVED***,
			FlagArg: []string***REMOVED***
				"a4ab:61d:f03e:5d7d:fad7:d4c2:a1a5:568",
				"203.107.49.208",
				"14.57.204.90",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			Want: []net.IP***REMOVED***
				net.ParseIP("204.228.73.195"),
				net.ParseIP("86.141.15.94"),
			***REMOVED***,
			FlagArg: []string***REMOVED***
				"204.228.73.195",
				"86.141.15.94",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			Want: []net.IP***REMOVED***
				net.ParseIP("c70c:db36:3001:890f:c6ea:3f9b:7a39:cc3f"),
				net.ParseIP("4d17:1d6e:e699:bd7a:88c5:5e7e:ac6a:4472"),
			***REMOVED***,
			FlagArg: []string***REMOVED***
				"c70c:db36:3001:890f:c6ea:3f9b:7a39:cc3f",
				"4d17:1d6e:e699:bd7a:88c5:5e7e:ac6a:4472",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			Want: []net.IP***REMOVED***
				net.ParseIP("5170:f971:cfac:7be3:512a:af37:952c:bc33"),
				net.ParseIP("93.21.145.140"),
				net.ParseIP("2cac:61d3:c5ff:6caf:73e0:1b1a:c336:c1ca"),
			***REMOVED***,
			FlagArg: []string***REMOVED***
				" 5170:f971:cfac:7be3:512a:af37:952c:bc33  , 93.21.145.140     ",
				"2cac:61d3:c5ff:6caf:73e0:1b1a:c336:c1ca",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			Want: []net.IP***REMOVED***
				net.ParseIP("2e5e:66b2:6441:848:5b74:76ea:574c:3a7b"),
				net.ParseIP("2e5e:66b2:6441:848:5b74:76ea:574c:3a7b"),
				net.ParseIP("2e5e:66b2:6441:848:5b74:76ea:574c:3a7b"),
				net.ParseIP("2e5e:66b2:6441:848:5b74:76ea:574c:3a7b"),
			***REMOVED***,
			FlagArg: []string***REMOVED***
				`"2e5e:66b2:6441:848:5b74:76ea:574c:3a7b,        2e5e:66b2:6441:848:5b74:76ea:574c:3a7b,2e5e:66b2:6441:848:5b74:76ea:574c:3a7b     "`,
				" 2e5e:66b2:6441:848:5b74:76ea:574c:3a7b"***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***

		var ips []net.IP
		f := setUpIPSFlagSet(&ips)

		if err := f.Parse([]string***REMOVED***fmt.Sprintf("--ips=%s", strings.Join(test.FlagArg, ","))***REMOVED***); err != nil ***REMOVED***
			t.Fatalf("flag parsing failed with error: %s\nparsing:\t%#v\nwant:\t\t%s",
				err, test.FlagArg, test.Want[i])
		***REMOVED***

		for j, b := range ips ***REMOVED***
			if !b.Equal(test.Want[j]) ***REMOVED***
				t.Fatalf("bad value parsed for test %d on net.IP %d:\nwant:\t%s\ngot:\t%s", i, j, test.Want[j], b)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
