package netlink

import (
	"fmt"
	"net"
)

// Neigh represents a link layer neighbor from netlink.
type Neigh struct ***REMOVED***
	LinkIndex    int
	Family       int
	State        int
	Type         int
	Flags        int
	IP           net.IP
	HardwareAddr net.HardwareAddr
	LLIPAddr     net.IP //Used in the case of NHRP
***REMOVED***

// String returns $ip/$hwaddr $label
func (neigh *Neigh) String() string ***REMOVED***
	return fmt.Sprintf("%s %s", neigh.IP, neigh.HardwareAddr)
***REMOVED***
