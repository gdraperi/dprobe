// Package netlink provides a simple library for netlink. Netlink is
// the interface a user-space program in linux uses to communicate with
// the kernel. It can be used to add and remove interfaces, set up ip
// addresses and routes, and confiugre ipsec. Netlink communication
// requires elevated privileges, so in most cases this code needs to
// be run as root. The low level primitives for netlink are contained
// in the nl subpackage. This package attempts to provide a high-level
// interface that is loosly modeled on the iproute2 cli.
package netlink

import (
	"errors"
	"net"
)

var (
	// ErrNotImplemented is returned when a requested feature is not implemented.
	ErrNotImplemented = errors.New("not implemented")
)

// ParseIPNet parses a string in ip/net format and returns a net.IPNet.
// This is valuable because addresses in netlink are often IPNets and
// ParseCIDR returns an IPNet with the IP part set to the base IP of the
// range.
func ParseIPNet(s string) (*net.IPNet, error) ***REMOVED***
	ip, ipNet, err := net.ParseCIDR(s)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &net.IPNet***REMOVED***IP: ip, Mask: ipNet.Mask***REMOVED***, nil
***REMOVED***

// NewIPNet generates an IPNet from an ip address using a netmask of 32 or 128.
func NewIPNet(ip net.IP) *net.IPNet ***REMOVED***
	if ip.To4() != nil ***REMOVED***
		return &net.IPNet***REMOVED***IP: ip, Mask: net.CIDRMask(32, 32)***REMOVED***
	***REMOVED***
	return &net.IPNet***REMOVED***IP: ip, Mask: net.CIDRMask(128, 128)***REMOVED***
***REMOVED***
