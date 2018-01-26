package netlink

import (
	"fmt"
	"net"
	"strings"
)

// Addr represents an IP address from netlink. Netlink ip addresses
// include a mask, so it stores the address as a net.IPNet.
type Addr struct ***REMOVED***
	*net.IPNet
	Label       string
	Flags       int
	Scope       int
	Peer        *net.IPNet
	Broadcast   net.IP
	PreferedLft int
	ValidLft    int
***REMOVED***

// String returns $ip/$netmask $label
func (a Addr) String() string ***REMOVED***
	return strings.TrimSpace(fmt.Sprintf("%s %s", a.IPNet, a.Label))
***REMOVED***

// ParseAddr parses the string representation of an address in the
// form $ip/$netmask $label. The label portion is optional
func ParseAddr(s string) (*Addr, error) ***REMOVED***
	label := ""
	parts := strings.Split(s, " ")
	if len(parts) > 1 ***REMOVED***
		s = parts[0]
		label = parts[1]
	***REMOVED***
	m, err := ParseIPNet(s)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &Addr***REMOVED***IPNet: m, Label: label***REMOVED***, nil
***REMOVED***

// Equal returns true if both Addrs have the same net.IPNet value.
func (a Addr) Equal(x Addr) bool ***REMOVED***
	sizea, _ := a.Mask.Size()
	sizeb, _ := x.Mask.Size()
	// ignore label for comparison
	return a.IP.Equal(x.IP) && sizea == sizeb
***REMOVED***

func (a Addr) PeerEqual(x Addr) bool ***REMOVED***
	sizea, _ := a.Peer.Mask.Size()
	sizeb, _ := x.Peer.Mask.Size()
	// ignore label for comparison
	return a.Peer.IP.Equal(x.Peer.IP) && sizea == sizeb
***REMOVED***
