package sockaddr

import (
	"fmt"
	"net"
)

// IfAddr is a union of a SockAddr and a net.Interface.
type IfAddr struct ***REMOVED***
	SockAddr
	net.Interface
***REMOVED***

// Attr returns the named attribute as a string
func (ifAddr IfAddr) Attr(attrName AttrName) (string, error) ***REMOVED***
	val := IfAddrAttr(ifAddr, attrName)
	if val != "" ***REMOVED***
		return val, nil
	***REMOVED***

	return Attr(ifAddr.SockAddr, attrName)
***REMOVED***

// Attr returns the named attribute as a string
func Attr(sa SockAddr, attrName AttrName) (string, error) ***REMOVED***
	switch sockType := sa.Type(); ***REMOVED***
	case sockType&TypeIP != 0:
		ip := *ToIPAddr(sa)
		attrVal := IPAddrAttr(ip, attrName)
		if attrVal != "" ***REMOVED***
			return attrVal, nil
		***REMOVED***

		if sockType == TypeIPv4 ***REMOVED***
			ipv4 := *ToIPv4Addr(sa)
			attrVal := IPv4AddrAttr(ipv4, attrName)
			if attrVal != "" ***REMOVED***
				return attrVal, nil
			***REMOVED***
		***REMOVED*** else if sockType == TypeIPv6 ***REMOVED***
			ipv6 := *ToIPv6Addr(sa)
			attrVal := IPv6AddrAttr(ipv6, attrName)
			if attrVal != "" ***REMOVED***
				return attrVal, nil
			***REMOVED***
		***REMOVED***

	case sockType == TypeUnix:
		us := *ToUnixSock(sa)
		attrVal := UnixSockAttr(us, attrName)
		if attrVal != "" ***REMOVED***
			return attrVal, nil
		***REMOVED***
	***REMOVED***

	// Non type-specific attributes
	switch attrName ***REMOVED***
	case "string":
		return sa.String(), nil
	case "type":
		return sa.Type().String(), nil
	***REMOVED***

	return "", fmt.Errorf("unsupported attribute name %q", attrName)
***REMOVED***
