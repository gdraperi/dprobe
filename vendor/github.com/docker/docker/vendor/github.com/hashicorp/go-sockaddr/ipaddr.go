package sockaddr

import (
	"fmt"
	"math/big"
	"net"
	"strings"
)

// Constants for the sizes of IPv3, IPv4, and IPv6 address types.
const (
	IPv3len = 6
	IPv4len = 4
	IPv6len = 16
)

// IPAddr is a generic IP address interface for IPv4 and IPv6 addresses,
// networks, and socket endpoints.
type IPAddr interface ***REMOVED***
	SockAddr
	AddressBinString() string
	AddressHexString() string
	Cmp(SockAddr) int
	CmpAddress(SockAddr) int
	CmpPort(SockAddr) int
	FirstUsable() IPAddr
	Host() IPAddr
	IPPort() IPPort
	LastUsable() IPAddr
	Maskbits() int
	NetIP() *net.IP
	NetIPMask() *net.IPMask
	NetIPNet() *net.IPNet
	Network() IPAddr
	Octets() []int
***REMOVED***

// IPPort is the type for an IP port number for the TCP and UDP IP transports.
type IPPort uint16

// IPPrefixLen is a typed integer representing the prefix length for a given
// IPAddr.
type IPPrefixLen byte

// ipAddrAttrMap is a map of the IPAddr type-specific attributes.
var ipAddrAttrMap map[AttrName]func(IPAddr) string
var ipAddrAttrs []AttrName

func init() ***REMOVED***
	ipAddrInit()
***REMOVED***

// NewIPAddr creates a new IPAddr from a string.  Returns nil if the string is
// not an IPv4 or an IPv6 address.
func NewIPAddr(addr string) (IPAddr, error) ***REMOVED***
	ipv4Addr, err := NewIPv4Addr(addr)
	if err == nil ***REMOVED***
		return ipv4Addr, nil
	***REMOVED***

	ipv6Addr, err := NewIPv6Addr(addr)
	if err == nil ***REMOVED***
		return ipv6Addr, nil
	***REMOVED***

	return nil, fmt.Errorf("invalid IPAddr %v", addr)
***REMOVED***

// IPAddrAttr returns a string representation of an attribute for the given
// IPAddr.
func IPAddrAttr(ip IPAddr, selector AttrName) string ***REMOVED***
	fn, found := ipAddrAttrMap[selector]
	if !found ***REMOVED***
		return ""
	***REMOVED***

	return fn(ip)
***REMOVED***

// IPAttrs returns a list of attributes supported by the IPAddr type
func IPAttrs() []AttrName ***REMOVED***
	return ipAddrAttrs
***REMOVED***

// MustIPAddr is a helper method that must return an IPAddr or panic on invalid
// input.
func MustIPAddr(addr string) IPAddr ***REMOVED***
	ip, err := NewIPAddr(addr)
	if err != nil ***REMOVED***
		panic(fmt.Sprintf("Unable to create an IPAddr from %+q: %v", addr, err))
	***REMOVED***
	return ip
***REMOVED***

// ipAddrInit is called once at init()
func ipAddrInit() ***REMOVED***
	// Sorted for human readability
	ipAddrAttrs = []AttrName***REMOVED***
		"host",
		"address",
		"port",
		"netmask",
		"network",
		"mask_bits",
		"binary",
		"hex",
		"first_usable",
		"last_usable",
		"octets",
	***REMOVED***

	ipAddrAttrMap = map[AttrName]func(ip IPAddr) string***REMOVED***
		"address": func(ip IPAddr) string ***REMOVED***
			return ip.NetIP().String()
		***REMOVED***,
		"binary": func(ip IPAddr) string ***REMOVED***
			return ip.AddressBinString()
		***REMOVED***,
		"first_usable": func(ip IPAddr) string ***REMOVED***
			return ip.FirstUsable().String()
		***REMOVED***,
		"hex": func(ip IPAddr) string ***REMOVED***
			return ip.AddressHexString()
		***REMOVED***,
		"host": func(ip IPAddr) string ***REMOVED***
			return ip.Host().String()
		***REMOVED***,
		"last_usable": func(ip IPAddr) string ***REMOVED***
			return ip.LastUsable().String()
		***REMOVED***,
		"mask_bits": func(ip IPAddr) string ***REMOVED***
			return fmt.Sprintf("%d", ip.Maskbits())
		***REMOVED***,
		"netmask": func(ip IPAddr) string ***REMOVED***
			switch v := ip.(type) ***REMOVED***
			case IPv4Addr:
				ipv4Mask := IPv4Addr***REMOVED***
					Address: IPv4Address(v.Mask),
					Mask:    IPv4HostMask,
				***REMOVED***
				return ipv4Mask.String()
			case IPv6Addr:
				ipv6Mask := new(big.Int)
				ipv6Mask.Set(v.Mask)
				ipv6MaskAddr := IPv6Addr***REMOVED***
					Address: IPv6Address(ipv6Mask),
					Mask:    ipv6HostMask,
				***REMOVED***
				return ipv6MaskAddr.String()
			default:
				return fmt.Sprintf("<unsupported type: %T>", ip)
			***REMOVED***
		***REMOVED***,
		"network": func(ip IPAddr) string ***REMOVED***
			return ip.Network().NetIP().String()
		***REMOVED***,
		"octets": func(ip IPAddr) string ***REMOVED***
			octets := ip.Octets()
			octetStrs := make([]string, 0, len(octets))
			for _, octet := range octets ***REMOVED***
				octetStrs = append(octetStrs, fmt.Sprintf("%d", octet))
			***REMOVED***
			return strings.Join(octetStrs, " ")
		***REMOVED***,
		"port": func(ip IPAddr) string ***REMOVED***
			return fmt.Sprintf("%d", ip.IPPort())
		***REMOVED***,
	***REMOVED***
***REMOVED***
