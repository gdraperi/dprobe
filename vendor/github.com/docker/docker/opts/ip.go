package opts

import (
	"fmt"
	"net"
)

// IPOpt holds an IP. It is used to store values from CLI flags.
type IPOpt struct ***REMOVED***
	*net.IP
***REMOVED***

// NewIPOpt creates a new IPOpt from a reference net.IP and a
// string representation of an IP. If the string is not a valid
// IP it will fallback to the specified reference.
func NewIPOpt(ref *net.IP, defaultVal string) *IPOpt ***REMOVED***
	o := &IPOpt***REMOVED***
		IP: ref,
	***REMOVED***
	o.Set(defaultVal)
	return o
***REMOVED***

// Set sets an IPv4 or IPv6 address from a given string. If the given
// string is not parsable as an IP address it returns an error.
func (o *IPOpt) Set(val string) error ***REMOVED***
	ip := net.ParseIP(val)
	if ip == nil ***REMOVED***
		return fmt.Errorf("%s is not an ip address", val)
	***REMOVED***
	*o.IP = ip
	return nil
***REMOVED***

// String returns the IP address stored in the IPOpt. If stored IP is a
// nil pointer, it returns an empty string.
func (o *IPOpt) String() string ***REMOVED***
	if *o.IP == nil ***REMOVED***
		return ""
	***REMOVED***
	return o.IP.String()
***REMOVED***

// Type returns the type of the option
func (o *IPOpt) Type() string ***REMOVED***
	return "ip"
***REMOVED***
