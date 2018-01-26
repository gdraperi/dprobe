package netlink

import (
	"fmt"
	"net"
)

// Rule represents a netlink rule.
type Rule struct ***REMOVED***
	Priority          int
	Family            int
	Table             int
	Mark              int
	Mask              int
	TunID             uint
	Goto              int
	Src               *net.IPNet
	Dst               *net.IPNet
	Flow              int
	IifName           string
	OifName           string
	SuppressIfgroup   int
	SuppressPrefixlen int
***REMOVED***

func (r Rule) String() string ***REMOVED***
	return fmt.Sprintf("ip rule %d: from %s table %d", r.Priority, r.Src, r.Table)
***REMOVED***

// NewRule return empty rules.
func NewRule() *Rule ***REMOVED***
	return &Rule***REMOVED***
		SuppressIfgroup:   -1,
		SuppressPrefixlen: -1,
		Priority:          -1,
		Mark:              -1,
		Mask:              -1,
		Goto:              -1,
		Flow:              -1,
	***REMOVED***
***REMOVED***
