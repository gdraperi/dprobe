package netlink

import (
	"fmt"
	"net"
	"strings"
)

// Scope is an enum representing a route scope.
type Scope uint8

type NextHopFlag int

type Destination interface ***REMOVED***
	Family() int
	Decode([]byte) error
	Encode() ([]byte, error)
	String() string
	Equal(Destination) bool
***REMOVED***

type Encap interface ***REMOVED***
	Type() int
	Decode([]byte) error
	Encode() ([]byte, error)
	String() string
	Equal(Encap) bool
***REMOVED***

// Route represents a netlink route.
type Route struct ***REMOVED***
	LinkIndex  int
	ILinkIndex int
	Scope      Scope
	Dst        *net.IPNet
	Src        net.IP
	Gw         net.IP
	MultiPath  []*NexthopInfo
	Protocol   int
	Priority   int
	Table      int
	Type       int
	Tos        int
	Flags      int
	MPLSDst    *int
	NewDst     Destination
	Encap      Encap
***REMOVED***

func (r Route) String() string ***REMOVED***
	elems := []string***REMOVED******REMOVED***
	if len(r.MultiPath) == 0 ***REMOVED***
		elems = append(elems, fmt.Sprintf("Ifindex: %d", r.LinkIndex))
	***REMOVED***
	if r.MPLSDst != nil ***REMOVED***
		elems = append(elems, fmt.Sprintf("Dst: %d", r.MPLSDst))
	***REMOVED*** else ***REMOVED***
		elems = append(elems, fmt.Sprintf("Dst: %s", r.Dst))
	***REMOVED***
	if r.NewDst != nil ***REMOVED***
		elems = append(elems, fmt.Sprintf("NewDst: %s", r.NewDst))
	***REMOVED***
	if r.Encap != nil ***REMOVED***
		elems = append(elems, fmt.Sprintf("Encap: %s", r.Encap))
	***REMOVED***
	elems = append(elems, fmt.Sprintf("Src: %s", r.Src))
	if len(r.MultiPath) > 0 ***REMOVED***
		elems = append(elems, fmt.Sprintf("Gw: %s", r.MultiPath))
	***REMOVED*** else ***REMOVED***
		elems = append(elems, fmt.Sprintf("Gw: %s", r.Gw))
	***REMOVED***
	elems = append(elems, fmt.Sprintf("Flags: %s", r.ListFlags()))
	elems = append(elems, fmt.Sprintf("Table: %d", r.Table))
	return fmt.Sprintf("***REMOVED***%s***REMOVED***", strings.Join(elems, " "))
***REMOVED***

func (r Route) Equal(x Route) bool ***REMOVED***
	return r.LinkIndex == x.LinkIndex &&
		r.ILinkIndex == x.ILinkIndex &&
		r.Scope == x.Scope &&
		ipNetEqual(r.Dst, x.Dst) &&
		r.Src.Equal(x.Src) &&
		r.Gw.Equal(x.Gw) &&
		nexthopInfoSlice(r.MultiPath).Equal(x.MultiPath) &&
		r.Protocol == x.Protocol &&
		r.Priority == x.Priority &&
		r.Table == x.Table &&
		r.Type == x.Type &&
		r.Tos == x.Tos &&
		r.Flags == x.Flags &&
		(r.MPLSDst == x.MPLSDst || (r.MPLSDst != nil && x.MPLSDst != nil && *r.MPLSDst == *x.MPLSDst)) &&
		(r.NewDst == x.NewDst || (r.NewDst != nil && r.NewDst.Equal(x.NewDst))) &&
		(r.Encap == x.Encap || (r.Encap != nil && r.Encap.Equal(x.Encap)))
***REMOVED***

func (r *Route) SetFlag(flag NextHopFlag) ***REMOVED***
	r.Flags |= int(flag)
***REMOVED***

func (r *Route) ClearFlag(flag NextHopFlag) ***REMOVED***
	r.Flags &^= int(flag)
***REMOVED***

type flagString struct ***REMOVED***
	f NextHopFlag
	s string
***REMOVED***

// RouteUpdate is sent when a route changes - type is RTM_NEWROUTE or RTM_DELROUTE
type RouteUpdate struct ***REMOVED***
	Type uint16
	Route
***REMOVED***

type NexthopInfo struct ***REMOVED***
	LinkIndex int
	Hops      int
	Gw        net.IP
	Flags     int
	NewDst    Destination
	Encap     Encap
***REMOVED***

func (n *NexthopInfo) String() string ***REMOVED***
	elems := []string***REMOVED******REMOVED***
	elems = append(elems, fmt.Sprintf("Ifindex: %d", n.LinkIndex))
	if n.NewDst != nil ***REMOVED***
		elems = append(elems, fmt.Sprintf("NewDst: %s", n.NewDst))
	***REMOVED***
	if n.Encap != nil ***REMOVED***
		elems = append(elems, fmt.Sprintf("Encap: %s", n.Encap))
	***REMOVED***
	elems = append(elems, fmt.Sprintf("Weight: %d", n.Hops+1))
	elems = append(elems, fmt.Sprintf("Gw: %s", n.Gw))
	elems = append(elems, fmt.Sprintf("Flags: %s", n.ListFlags()))
	return fmt.Sprintf("***REMOVED***%s***REMOVED***", strings.Join(elems, " "))
***REMOVED***

func (n NexthopInfo) Equal(x NexthopInfo) bool ***REMOVED***
	return n.LinkIndex == x.LinkIndex &&
		n.Hops == x.Hops &&
		n.Gw.Equal(x.Gw) &&
		n.Flags == x.Flags &&
		(n.NewDst == x.NewDst || (n.NewDst != nil && n.NewDst.Equal(x.NewDst))) &&
		(n.Encap == x.Encap || (n.Encap != nil && n.Encap.Equal(x.Encap)))
***REMOVED***

type nexthopInfoSlice []*NexthopInfo

func (n nexthopInfoSlice) Equal(x []*NexthopInfo) bool ***REMOVED***
	if len(n) != len(x) ***REMOVED***
		return false
	***REMOVED***
	for i := range n ***REMOVED***
		if n[i] == nil || x[i] == nil ***REMOVED***
			return false
		***REMOVED***
		if !n[i].Equal(*x[i]) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// ipNetEqual returns true iff both IPNet are equal
func ipNetEqual(ipn1 *net.IPNet, ipn2 *net.IPNet) bool ***REMOVED***
	if ipn1 == ipn2 ***REMOVED***
		return true
	***REMOVED***
	if ipn1 == nil || ipn2 == nil ***REMOVED***
		return false
	***REMOVED***
	m1, _ := ipn1.Mask.Size()
	m2, _ := ipn2.Mask.Size()
	return m1 == m2 && ipn1.IP.Equal(ipn2.IP)
***REMOVED***
