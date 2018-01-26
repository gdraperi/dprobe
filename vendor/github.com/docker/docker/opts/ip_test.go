package opts

import (
	"net"
	"testing"
)

func TestIpOptString(t *testing.T) ***REMOVED***
	addresses := []string***REMOVED***"", "0.0.0.0"***REMOVED***
	var ip net.IP

	for _, address := range addresses ***REMOVED***
		stringAddress := NewIPOpt(&ip, address).String()
		if stringAddress != address ***REMOVED***
			t.Fatalf("IpOpt string should be `%s`, not `%s`", address, stringAddress)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestNewIpOptInvalidDefaultVal(t *testing.T) ***REMOVED***
	ip := net.IPv4(127, 0, 0, 1)
	defaultVal := "Not an ip"

	ipOpt := NewIPOpt(&ip, defaultVal)

	expected := "127.0.0.1"
	if ipOpt.String() != expected ***REMOVED***
		t.Fatalf("Expected [%v], got [%v]", expected, ipOpt.String())
	***REMOVED***
***REMOVED***

func TestNewIpOptValidDefaultVal(t *testing.T) ***REMOVED***
	ip := net.IPv4(127, 0, 0, 1)
	defaultVal := "192.168.1.1"

	ipOpt := NewIPOpt(&ip, defaultVal)

	expected := "192.168.1.1"
	if ipOpt.String() != expected ***REMOVED***
		t.Fatalf("Expected [%v], got [%v]", expected, ipOpt.String())
	***REMOVED***
***REMOVED***

func TestIpOptSetInvalidVal(t *testing.T) ***REMOVED***
	ip := net.IPv4(127, 0, 0, 1)
	ipOpt := &IPOpt***REMOVED***IP: &ip***REMOVED***

	invalidIP := "invalid ip"
	expectedError := "invalid ip is not an ip address"
	err := ipOpt.Set(invalidIP)
	if err == nil || err.Error() != expectedError ***REMOVED***
		t.Fatalf("Expected an Error with [%v], got [%v]", expectedError, err.Error())
	***REMOVED***
***REMOVED***
