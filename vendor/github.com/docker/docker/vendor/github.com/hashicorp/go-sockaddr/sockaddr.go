package sockaddr

import (
	"fmt"
	"strings"
)

type SockAddrType int
type AttrName string

const (
	TypeUnknown SockAddrType = 0x0
	TypeUnix                 = 0x1
	TypeIPv4                 = 0x2
	TypeIPv6                 = 0x4

	// TypeIP is the union of TypeIPv4 and TypeIPv6
	TypeIP = 0x6
)

type SockAddr interface ***REMOVED***
	// CmpRFC returns 0 if SockAddr exactly matches one of the matched RFC
	// networks, -1 if the receiver is contained within the RFC network, or
	// 1 if the address is not contained within the RFC.
	CmpRFC(rfcNum uint, sa SockAddr) int

	// Contains returns true if the SockAddr arg is contained within the
	// receiver
	Contains(SockAddr) bool

	// Equal allows for the comparison of two SockAddrs
	Equal(SockAddr) bool

	DialPacketArgs() (string, string)
	DialStreamArgs() (string, string)
	ListenPacketArgs() (string, string)
	ListenStreamArgs() (string, string)

	// String returns the string representation of SockAddr
	String() string

	// Type returns the SockAddrType
	Type() SockAddrType
***REMOVED***

// sockAddrAttrMap is a map of the SockAddr type-specific attributes.
var sockAddrAttrMap map[AttrName]func(SockAddr) string
var sockAddrAttrs []AttrName

func init() ***REMOVED***
	sockAddrInit()
***REMOVED***

// New creates a new SockAddr from the string.  The order in which New()
// attempts to construct a SockAddr is: IPv4Addr, IPv6Addr, SockAddrUnix.
//
// NOTE: New() relies on the heuristic wherein if the path begins with either a
// '.'  or '/' character before creating a new UnixSock.  For UNIX sockets that
// are absolute paths or are nested within a sub-directory, this works as
// expected, however if the UNIX socket is contained in the current working
// directory, this will fail unless the path begins with "./"
// (e.g. "./my-local-socket").  Calls directly to NewUnixSock() do not suffer
// this limitation.  Invalid IP addresses such as "256.0.0.0/-1" will run afoul
// of this heuristic and be assumed to be a valid UNIX socket path (which they
// are, but it is probably not what you want and you won't realize it until you
// stat(2) the file system to discover it doesn't exist).
func NewSockAddr(s string) (SockAddr, error) ***REMOVED***
	ipv4Addr, err := NewIPv4Addr(s)
	if err == nil ***REMOVED***
		return ipv4Addr, nil
	***REMOVED***

	ipv6Addr, err := NewIPv6Addr(s)
	if err == nil ***REMOVED***
		return ipv6Addr, nil
	***REMOVED***

	// Check to make sure the string begins with either a '.' or '/', or
	// contains a '/'.
	if len(s) > 1 && (strings.IndexAny(s[0:1], "./") != -1 || strings.IndexByte(s, '/') != -1) ***REMOVED***
		unixSock, err := NewUnixSock(s)
		if err == nil ***REMOVED***
			return unixSock, nil
		***REMOVED***
	***REMOVED***

	return nil, fmt.Errorf("Unable to convert %q to an IPv4 or IPv6 address, or a UNIX Socket", s)
***REMOVED***

// ToIPAddr returns an IPAddr type or nil if the type conversion fails.
func ToIPAddr(sa SockAddr) *IPAddr ***REMOVED***
	ipa, ok := sa.(IPAddr)
	if !ok ***REMOVED***
		return nil
	***REMOVED***
	return &ipa
***REMOVED***

// ToIPv4Addr returns an IPv4Addr type or nil if the type conversion fails.
func ToIPv4Addr(sa SockAddr) *IPv4Addr ***REMOVED***
	switch v := sa.(type) ***REMOVED***
	case IPv4Addr:
		return &v
	default:
		return nil
	***REMOVED***
***REMOVED***

// ToIPv6Addr returns an IPv6Addr type or nil if the type conversion fails.
func ToIPv6Addr(sa SockAddr) *IPv6Addr ***REMOVED***
	switch v := sa.(type) ***REMOVED***
	case IPv6Addr:
		return &v
	default:
		return nil
	***REMOVED***
***REMOVED***

// ToUnixSock returns a UnixSock type or nil if the type conversion fails.
func ToUnixSock(sa SockAddr) *UnixSock ***REMOVED***
	switch v := sa.(type) ***REMOVED***
	case UnixSock:
		return &v
	default:
		return nil
	***REMOVED***
***REMOVED***

// SockAddrAttr returns a string representation of an attribute for the given
// SockAddr.
func SockAddrAttr(sa SockAddr, selector AttrName) string ***REMOVED***
	fn, found := sockAddrAttrMap[selector]
	if !found ***REMOVED***
		return ""
	***REMOVED***

	return fn(sa)
***REMOVED***

// String() for SockAddrType returns a string representation of the
// SockAddrType (e.g. "IPv4", "IPv6", "UNIX", "IP", or "unknown").
func (sat SockAddrType) String() string ***REMOVED***
	switch sat ***REMOVED***
	case TypeIPv4:
		return "IPv4"
	case TypeIPv6:
		return "IPv6"
	// There is no concrete "IP" type.  Leaving here as a reminder.
	// case TypeIP:
	// 	return "IP"
	case TypeUnix:
		return "UNIX"
	default:
		panic("unsupported type")
	***REMOVED***
***REMOVED***

// sockAddrInit is called once at init()
func sockAddrInit() ***REMOVED***
	sockAddrAttrs = []AttrName***REMOVED***
		"type", // type should be first
		"string",
	***REMOVED***

	sockAddrAttrMap = map[AttrName]func(sa SockAddr) string***REMOVED***
		"string": func(sa SockAddr) string ***REMOVED***
			return sa.String()
		***REMOVED***,
		"type": func(sa SockAddr) string ***REMOVED***
			return sa.Type().String()
		***REMOVED***,
	***REMOVED***
***REMOVED***

// UnixSockAttrs returns a list of attributes supported by the UnixSock type
func SockAddrAttrs() []AttrName ***REMOVED***
	return sockAddrAttrs
***REMOVED***
