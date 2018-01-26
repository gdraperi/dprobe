package netlink

import (
	"fmt"
	"net"
	"strings"
	"syscall"

	"github.com/vishvananda/netlink/nl"
	"github.com/vishvananda/netns"
)

// RtAttr is shared so it is in netlink_linux.go

const (
	SCOPE_UNIVERSE Scope = syscall.RT_SCOPE_UNIVERSE
	SCOPE_SITE     Scope = syscall.RT_SCOPE_SITE
	SCOPE_LINK     Scope = syscall.RT_SCOPE_LINK
	SCOPE_HOST     Scope = syscall.RT_SCOPE_HOST
	SCOPE_NOWHERE  Scope = syscall.RT_SCOPE_NOWHERE
)

const (
	RT_FILTER_PROTOCOL uint64 = 1 << (1 + iota)
	RT_FILTER_SCOPE
	RT_FILTER_TYPE
	RT_FILTER_TOS
	RT_FILTER_IIF
	RT_FILTER_OIF
	RT_FILTER_DST
	RT_FILTER_SRC
	RT_FILTER_GW
	RT_FILTER_TABLE
)

const (
	FLAG_ONLINK    NextHopFlag = syscall.RTNH_F_ONLINK
	FLAG_PERVASIVE NextHopFlag = syscall.RTNH_F_PERVASIVE
)

var testFlags = []flagString***REMOVED***
	***REMOVED***f: FLAG_ONLINK, s: "onlink"***REMOVED***,
	***REMOVED***f: FLAG_PERVASIVE, s: "pervasive"***REMOVED***,
***REMOVED***

func listFlags(flag int) []string ***REMOVED***
	var flags []string
	for _, tf := range testFlags ***REMOVED***
		if flag&int(tf.f) != 0 ***REMOVED***
			flags = append(flags, tf.s)
		***REMOVED***
	***REMOVED***
	return flags
***REMOVED***

func (r *Route) ListFlags() []string ***REMOVED***
	return listFlags(r.Flags)
***REMOVED***

func (n *NexthopInfo) ListFlags() []string ***REMOVED***
	return listFlags(n.Flags)
***REMOVED***

type MPLSDestination struct ***REMOVED***
	Labels []int
***REMOVED***

func (d *MPLSDestination) Family() int ***REMOVED***
	return nl.FAMILY_MPLS
***REMOVED***

func (d *MPLSDestination) Decode(buf []byte) error ***REMOVED***
	d.Labels = nl.DecodeMPLSStack(buf)
	return nil
***REMOVED***

func (d *MPLSDestination) Encode() ([]byte, error) ***REMOVED***
	return nl.EncodeMPLSStack(d.Labels...), nil
***REMOVED***

func (d *MPLSDestination) String() string ***REMOVED***
	s := make([]string, 0, len(d.Labels))
	for _, l := range d.Labels ***REMOVED***
		s = append(s, fmt.Sprintf("%d", l))
	***REMOVED***
	return strings.Join(s, "/")
***REMOVED***

func (d *MPLSDestination) Equal(x Destination) bool ***REMOVED***
	o, ok := x.(*MPLSDestination)
	if !ok ***REMOVED***
		return false
	***REMOVED***
	if d == nil && o == nil ***REMOVED***
		return true
	***REMOVED***
	if d == nil || o == nil ***REMOVED***
		return false
	***REMOVED***
	if d.Labels == nil && o.Labels == nil ***REMOVED***
		return true
	***REMOVED***
	if d.Labels == nil || o.Labels == nil ***REMOVED***
		return false
	***REMOVED***
	if len(d.Labels) != len(o.Labels) ***REMOVED***
		return false
	***REMOVED***
	for i := range d.Labels ***REMOVED***
		if d.Labels[i] != o.Labels[i] ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

type MPLSEncap struct ***REMOVED***
	Labels []int
***REMOVED***

func (e *MPLSEncap) Type() int ***REMOVED***
	return nl.LWTUNNEL_ENCAP_MPLS
***REMOVED***

func (e *MPLSEncap) Decode(buf []byte) error ***REMOVED***
	if len(buf) < 4 ***REMOVED***
		return fmt.Errorf("Lack of bytes")
	***REMOVED***
	native := nl.NativeEndian()
	l := native.Uint16(buf)
	if len(buf) < int(l) ***REMOVED***
		return fmt.Errorf("Lack of bytes")
	***REMOVED***
	buf = buf[:l]
	typ := native.Uint16(buf[2:])
	if typ != nl.MPLS_IPTUNNEL_DST ***REMOVED***
		return fmt.Errorf("Unknown MPLS Encap Type: %d", typ)
	***REMOVED***
	e.Labels = nl.DecodeMPLSStack(buf[4:])
	return nil
***REMOVED***

func (e *MPLSEncap) Encode() ([]byte, error) ***REMOVED***
	s := nl.EncodeMPLSStack(e.Labels...)
	native := nl.NativeEndian()
	hdr := make([]byte, 4)
	native.PutUint16(hdr, uint16(len(s)+4))
	native.PutUint16(hdr[2:], nl.MPLS_IPTUNNEL_DST)
	return append(hdr, s...), nil
***REMOVED***

func (e *MPLSEncap) String() string ***REMOVED***
	s := make([]string, 0, len(e.Labels))
	for _, l := range e.Labels ***REMOVED***
		s = append(s, fmt.Sprintf("%d", l))
	***REMOVED***
	return strings.Join(s, "/")
***REMOVED***

func (e *MPLSEncap) Equal(x Encap) bool ***REMOVED***
	o, ok := x.(*MPLSEncap)
	if !ok ***REMOVED***
		return false
	***REMOVED***
	if e == nil && o == nil ***REMOVED***
		return true
	***REMOVED***
	if e == nil || o == nil ***REMOVED***
		return false
	***REMOVED***
	if e.Labels == nil && o.Labels == nil ***REMOVED***
		return true
	***REMOVED***
	if e.Labels == nil || o.Labels == nil ***REMOVED***
		return false
	***REMOVED***
	if len(e.Labels) != len(o.Labels) ***REMOVED***
		return false
	***REMOVED***
	for i := range e.Labels ***REMOVED***
		if e.Labels[i] != o.Labels[i] ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// RouteAdd will add a route to the system.
// Equivalent to: `ip route add $route`
func RouteAdd(route *Route) error ***REMOVED***
	return pkgHandle.RouteAdd(route)
***REMOVED***

// RouteAdd will add a route to the system.
// Equivalent to: `ip route add $route`
func (h *Handle) RouteAdd(route *Route) error ***REMOVED***
	flags := syscall.NLM_F_CREATE | syscall.NLM_F_EXCL | syscall.NLM_F_ACK
	req := h.newNetlinkRequest(syscall.RTM_NEWROUTE, flags)
	return h.routeHandle(route, req, nl.NewRtMsg())
***REMOVED***

// RouteReplace will add a route to the system.
// Equivalent to: `ip route replace $route`
func RouteReplace(route *Route) error ***REMOVED***
	return pkgHandle.RouteReplace(route)
***REMOVED***

// RouteReplace will add a route to the system.
// Equivalent to: `ip route replace $route`
func (h *Handle) RouteReplace(route *Route) error ***REMOVED***
	flags := syscall.NLM_F_CREATE | syscall.NLM_F_REPLACE | syscall.NLM_F_ACK
	req := h.newNetlinkRequest(syscall.RTM_NEWROUTE, flags)
	return h.routeHandle(route, req, nl.NewRtMsg())
***REMOVED***

// RouteDel will delete a route from the system.
// Equivalent to: `ip route del $route`
func RouteDel(route *Route) error ***REMOVED***
	return pkgHandle.RouteDel(route)
***REMOVED***

// RouteDel will delete a route from the system.
// Equivalent to: `ip route del $route`
func (h *Handle) RouteDel(route *Route) error ***REMOVED***
	req := h.newNetlinkRequest(syscall.RTM_DELROUTE, syscall.NLM_F_ACK)
	return h.routeHandle(route, req, nl.NewRtDelMsg())
***REMOVED***

func (h *Handle) routeHandle(route *Route, req *nl.NetlinkRequest, msg *nl.RtMsg) error ***REMOVED***
	if (route.Dst == nil || route.Dst.IP == nil) && route.Src == nil && route.Gw == nil && route.MPLSDst == nil ***REMOVED***
		return fmt.Errorf("one of Dst.IP, Src, or Gw must not be nil")
	***REMOVED***

	family := -1
	var rtAttrs []*nl.RtAttr

	if route.Dst != nil && route.Dst.IP != nil ***REMOVED***
		dstLen, _ := route.Dst.Mask.Size()
		msg.Dst_len = uint8(dstLen)
		dstFamily := nl.GetIPFamily(route.Dst.IP)
		family = dstFamily
		var dstData []byte
		if dstFamily == FAMILY_V4 ***REMOVED***
			dstData = route.Dst.IP.To4()
		***REMOVED*** else ***REMOVED***
			dstData = route.Dst.IP.To16()
		***REMOVED***
		rtAttrs = append(rtAttrs, nl.NewRtAttr(syscall.RTA_DST, dstData))
	***REMOVED*** else if route.MPLSDst != nil ***REMOVED***
		family = nl.FAMILY_MPLS
		msg.Dst_len = uint8(20)
		msg.Type = syscall.RTN_UNICAST
		rtAttrs = append(rtAttrs, nl.NewRtAttr(syscall.RTA_DST, nl.EncodeMPLSStack(*route.MPLSDst)))
	***REMOVED***

	if route.NewDst != nil ***REMOVED***
		if family != -1 && family != route.NewDst.Family() ***REMOVED***
			return fmt.Errorf("new destination and destination are not the same address family")
		***REMOVED***
		buf, err := route.NewDst.Encode()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		rtAttrs = append(rtAttrs, nl.NewRtAttr(nl.RTA_NEWDST, buf))
	***REMOVED***

	if route.Encap != nil ***REMOVED***
		buf := make([]byte, 2)
		native.PutUint16(buf, uint16(route.Encap.Type()))
		rtAttrs = append(rtAttrs, nl.NewRtAttr(nl.RTA_ENCAP_TYPE, buf))
		buf, err := route.Encap.Encode()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		rtAttrs = append(rtAttrs, nl.NewRtAttr(nl.RTA_ENCAP, buf))
	***REMOVED***

	if route.Src != nil ***REMOVED***
		srcFamily := nl.GetIPFamily(route.Src)
		if family != -1 && family != srcFamily ***REMOVED***
			return fmt.Errorf("source and destination ip are not the same IP family")
		***REMOVED***
		family = srcFamily
		var srcData []byte
		if srcFamily == FAMILY_V4 ***REMOVED***
			srcData = route.Src.To4()
		***REMOVED*** else ***REMOVED***
			srcData = route.Src.To16()
		***REMOVED***
		// The commonly used src ip for routes is actually PREFSRC
		rtAttrs = append(rtAttrs, nl.NewRtAttr(syscall.RTA_PREFSRC, srcData))
	***REMOVED***

	if route.Gw != nil ***REMOVED***
		gwFamily := nl.GetIPFamily(route.Gw)
		if family != -1 && family != gwFamily ***REMOVED***
			return fmt.Errorf("gateway, source, and destination ip are not the same IP family")
		***REMOVED***
		family = gwFamily
		var gwData []byte
		if gwFamily == FAMILY_V4 ***REMOVED***
			gwData = route.Gw.To4()
		***REMOVED*** else ***REMOVED***
			gwData = route.Gw.To16()
		***REMOVED***
		rtAttrs = append(rtAttrs, nl.NewRtAttr(syscall.RTA_GATEWAY, gwData))
	***REMOVED***

	if len(route.MultiPath) > 0 ***REMOVED***
		buf := []byte***REMOVED******REMOVED***
		for _, nh := range route.MultiPath ***REMOVED***
			rtnh := &nl.RtNexthop***REMOVED***
				RtNexthop: syscall.RtNexthop***REMOVED***
					Hops:    uint8(nh.Hops),
					Ifindex: int32(nh.LinkIndex),
					Flags:   uint8(nh.Flags),
				***REMOVED***,
			***REMOVED***
			children := []nl.NetlinkRequestData***REMOVED******REMOVED***
			if nh.Gw != nil ***REMOVED***
				gwFamily := nl.GetIPFamily(nh.Gw)
				if family != -1 && family != gwFamily ***REMOVED***
					return fmt.Errorf("gateway, source, and destination ip are not the same IP family")
				***REMOVED***
				if gwFamily == FAMILY_V4 ***REMOVED***
					children = append(children, nl.NewRtAttr(syscall.RTA_GATEWAY, []byte(nh.Gw.To4())))
				***REMOVED*** else ***REMOVED***
					children = append(children, nl.NewRtAttr(syscall.RTA_GATEWAY, []byte(nh.Gw.To16())))
				***REMOVED***
			***REMOVED***
			if nh.NewDst != nil ***REMOVED***
				if family != -1 && family != nh.NewDst.Family() ***REMOVED***
					return fmt.Errorf("new destination and destination are not the same address family")
				***REMOVED***
				buf, err := nh.NewDst.Encode()
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				children = append(children, nl.NewRtAttr(nl.RTA_NEWDST, buf))
			***REMOVED***
			if nh.Encap != nil ***REMOVED***
				buf := make([]byte, 2)
				native.PutUint16(buf, uint16(nh.Encap.Type()))
				rtAttrs = append(rtAttrs, nl.NewRtAttr(nl.RTA_ENCAP_TYPE, buf))
				buf, err := nh.Encap.Encode()
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				children = append(children, nl.NewRtAttr(nl.RTA_ENCAP, buf))
			***REMOVED***
			rtnh.Children = children
			buf = append(buf, rtnh.Serialize()...)
		***REMOVED***
		rtAttrs = append(rtAttrs, nl.NewRtAttr(syscall.RTA_MULTIPATH, buf))
	***REMOVED***

	if route.Table > 0 ***REMOVED***
		if route.Table >= 256 ***REMOVED***
			msg.Table = syscall.RT_TABLE_UNSPEC
			b := make([]byte, 4)
			native.PutUint32(b, uint32(route.Table))
			rtAttrs = append(rtAttrs, nl.NewRtAttr(syscall.RTA_TABLE, b))
		***REMOVED*** else ***REMOVED***
			msg.Table = uint8(route.Table)
		***REMOVED***
	***REMOVED***

	if route.Priority > 0 ***REMOVED***
		b := make([]byte, 4)
		native.PutUint32(b, uint32(route.Priority))
		rtAttrs = append(rtAttrs, nl.NewRtAttr(syscall.RTA_PRIORITY, b))
	***REMOVED***
	if route.Tos > 0 ***REMOVED***
		msg.Tos = uint8(route.Tos)
	***REMOVED***
	if route.Protocol > 0 ***REMOVED***
		msg.Protocol = uint8(route.Protocol)
	***REMOVED***
	if route.Type > 0 ***REMOVED***
		msg.Type = uint8(route.Type)
	***REMOVED***

	msg.Flags = uint32(route.Flags)
	msg.Scope = uint8(route.Scope)
	msg.Family = uint8(family)
	req.AddData(msg)
	for _, attr := range rtAttrs ***REMOVED***
		req.AddData(attr)
	***REMOVED***

	var (
		b      = make([]byte, 4)
		native = nl.NativeEndian()
	)
	native.PutUint32(b, uint32(route.LinkIndex))

	req.AddData(nl.NewRtAttr(syscall.RTA_OIF, b))

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

// RouteList gets a list of routes in the system.
// Equivalent to: `ip route show`.
// The list can be filtered by link and ip family.
func RouteList(link Link, family int) ([]Route, error) ***REMOVED***
	return pkgHandle.RouteList(link, family)
***REMOVED***

// RouteList gets a list of routes in the system.
// Equivalent to: `ip route show`.
// The list can be filtered by link and ip family.
func (h *Handle) RouteList(link Link, family int) ([]Route, error) ***REMOVED***
	var routeFilter *Route
	if link != nil ***REMOVED***
		routeFilter = &Route***REMOVED***
			LinkIndex: link.Attrs().Index,
		***REMOVED***
	***REMOVED***
	return h.RouteListFiltered(family, routeFilter, RT_FILTER_OIF)
***REMOVED***

// RouteListFiltered gets a list of routes in the system filtered with specified rules.
// All rules must be defined in RouteFilter struct
func RouteListFiltered(family int, filter *Route, filterMask uint64) ([]Route, error) ***REMOVED***
	return pkgHandle.RouteListFiltered(family, filter, filterMask)
***REMOVED***

// RouteListFiltered gets a list of routes in the system filtered with specified rules.
// All rules must be defined in RouteFilter struct
func (h *Handle) RouteListFiltered(family int, filter *Route, filterMask uint64) ([]Route, error) ***REMOVED***
	req := h.newNetlinkRequest(syscall.RTM_GETROUTE, syscall.NLM_F_DUMP)
	infmsg := nl.NewIfInfomsg(family)
	req.AddData(infmsg)

	msgs, err := req.Execute(syscall.NETLINK_ROUTE, syscall.RTM_NEWROUTE)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var res []Route
	for _, m := range msgs ***REMOVED***
		msg := nl.DeserializeRtMsg(m)
		if msg.Flags&syscall.RTM_F_CLONED != 0 ***REMOVED***
			// Ignore cloned routes
			continue
		***REMOVED***
		if msg.Table != syscall.RT_TABLE_MAIN ***REMOVED***
			if filter == nil || filter != nil && filterMask&RT_FILTER_TABLE == 0 ***REMOVED***
				// Ignore non-main tables
				continue
			***REMOVED***
		***REMOVED***
		route, err := deserializeRoute(m)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if filter != nil ***REMOVED***
			switch ***REMOVED***
			case filterMask&RT_FILTER_TABLE != 0 && filter.Table != syscall.RT_TABLE_UNSPEC && route.Table != filter.Table:
				continue
			case filterMask&RT_FILTER_PROTOCOL != 0 && route.Protocol != filter.Protocol:
				continue
			case filterMask&RT_FILTER_SCOPE != 0 && route.Scope != filter.Scope:
				continue
			case filterMask&RT_FILTER_TYPE != 0 && route.Type != filter.Type:
				continue
			case filterMask&RT_FILTER_TOS != 0 && route.Tos != filter.Tos:
				continue
			case filterMask&RT_FILTER_OIF != 0 && route.LinkIndex != filter.LinkIndex:
				continue
			case filterMask&RT_FILTER_IIF != 0 && route.ILinkIndex != filter.ILinkIndex:
				continue
			case filterMask&RT_FILTER_GW != 0 && !route.Gw.Equal(filter.Gw):
				continue
			case filterMask&RT_FILTER_SRC != 0 && !route.Src.Equal(filter.Src):
				continue
			case filterMask&RT_FILTER_DST != 0:
				if filter.MPLSDst == nil || route.MPLSDst == nil || (*filter.MPLSDst) != (*route.MPLSDst) ***REMOVED***
					if !ipNetEqual(route.Dst, filter.Dst) ***REMOVED***
						continue
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		res = append(res, route)
	***REMOVED***
	return res, nil
***REMOVED***

// deserializeRoute decodes a binary netlink message into a Route struct
func deserializeRoute(m []byte) (Route, error) ***REMOVED***
	msg := nl.DeserializeRtMsg(m)
	attrs, err := nl.ParseRouteAttr(m[msg.Len():])
	if err != nil ***REMOVED***
		return Route***REMOVED******REMOVED***, err
	***REMOVED***
	route := Route***REMOVED***
		Scope:    Scope(msg.Scope),
		Protocol: int(msg.Protocol),
		Table:    int(msg.Table),
		Type:     int(msg.Type),
		Tos:      int(msg.Tos),
		Flags:    int(msg.Flags),
	***REMOVED***

	native := nl.NativeEndian()
	var encap, encapType syscall.NetlinkRouteAttr
	for _, attr := range attrs ***REMOVED***
		switch attr.Attr.Type ***REMOVED***
		case syscall.RTA_GATEWAY:
			route.Gw = net.IP(attr.Value)
		case syscall.RTA_PREFSRC:
			route.Src = net.IP(attr.Value)
		case syscall.RTA_DST:
			if msg.Family == nl.FAMILY_MPLS ***REMOVED***
				stack := nl.DecodeMPLSStack(attr.Value)
				if len(stack) == 0 || len(stack) > 1 ***REMOVED***
					return route, fmt.Errorf("invalid MPLS RTA_DST")
				***REMOVED***
				route.MPLSDst = &stack[0]
			***REMOVED*** else ***REMOVED***
				route.Dst = &net.IPNet***REMOVED***
					IP:   attr.Value,
					Mask: net.CIDRMask(int(msg.Dst_len), 8*len(attr.Value)),
				***REMOVED***
			***REMOVED***
		case syscall.RTA_OIF:
			route.LinkIndex = int(native.Uint32(attr.Value[0:4]))
		case syscall.RTA_IIF:
			route.ILinkIndex = int(native.Uint32(attr.Value[0:4]))
		case syscall.RTA_PRIORITY:
			route.Priority = int(native.Uint32(attr.Value[0:4]))
		case syscall.RTA_TABLE:
			route.Table = int(native.Uint32(attr.Value[0:4]))
		case syscall.RTA_MULTIPATH:
			parseRtNexthop := func(value []byte) (*NexthopInfo, []byte, error) ***REMOVED***
				if len(value) < syscall.SizeofRtNexthop ***REMOVED***
					return nil, nil, fmt.Errorf("Lack of bytes")
				***REMOVED***
				nh := nl.DeserializeRtNexthop(value)
				if len(value) < int(nh.RtNexthop.Len) ***REMOVED***
					return nil, nil, fmt.Errorf("Lack of bytes")
				***REMOVED***
				info := &NexthopInfo***REMOVED***
					LinkIndex: int(nh.RtNexthop.Ifindex),
					Hops:      int(nh.RtNexthop.Hops),
					Flags:     int(nh.RtNexthop.Flags),
				***REMOVED***
				attrs, err := nl.ParseRouteAttr(value[syscall.SizeofRtNexthop:int(nh.RtNexthop.Len)])
				if err != nil ***REMOVED***
					return nil, nil, err
				***REMOVED***
				var encap, encapType syscall.NetlinkRouteAttr
				for _, attr := range attrs ***REMOVED***
					switch attr.Attr.Type ***REMOVED***
					case syscall.RTA_GATEWAY:
						info.Gw = net.IP(attr.Value)
					case nl.RTA_NEWDST:
						var d Destination
						switch msg.Family ***REMOVED***
						case nl.FAMILY_MPLS:
							d = &MPLSDestination***REMOVED******REMOVED***
						***REMOVED***
						if err := d.Decode(attr.Value); err != nil ***REMOVED***
							return nil, nil, err
						***REMOVED***
						info.NewDst = d
					case nl.RTA_ENCAP_TYPE:
						encapType = attr
					case nl.RTA_ENCAP:
						encap = attr
					***REMOVED***
				***REMOVED***

				if len(encap.Value) != 0 && len(encapType.Value) != 0 ***REMOVED***
					typ := int(native.Uint16(encapType.Value[0:2]))
					var e Encap
					switch typ ***REMOVED***
					case nl.LWTUNNEL_ENCAP_MPLS:
						e = &MPLSEncap***REMOVED******REMOVED***
						if err := e.Decode(encap.Value); err != nil ***REMOVED***
							return nil, nil, err
						***REMOVED***
					***REMOVED***
					info.Encap = e
				***REMOVED***

				return info, value[int(nh.RtNexthop.Len):], nil
			***REMOVED***
			rest := attr.Value
			for len(rest) > 0 ***REMOVED***
				info, buf, err := parseRtNexthop(rest)
				if err != nil ***REMOVED***
					return route, err
				***REMOVED***
				route.MultiPath = append(route.MultiPath, info)
				rest = buf
			***REMOVED***
		case nl.RTA_NEWDST:
			var d Destination
			switch msg.Family ***REMOVED***
			case nl.FAMILY_MPLS:
				d = &MPLSDestination***REMOVED******REMOVED***
			***REMOVED***
			if err := d.Decode(attr.Value); err != nil ***REMOVED***
				return route, err
			***REMOVED***
			route.NewDst = d
		case nl.RTA_ENCAP_TYPE:
			encapType = attr
		case nl.RTA_ENCAP:
			encap = attr
		***REMOVED***
	***REMOVED***

	if len(encap.Value) != 0 && len(encapType.Value) != 0 ***REMOVED***
		typ := int(native.Uint16(encapType.Value[0:2]))
		var e Encap
		switch typ ***REMOVED***
		case nl.LWTUNNEL_ENCAP_MPLS:
			e = &MPLSEncap***REMOVED******REMOVED***
			if err := e.Decode(encap.Value); err != nil ***REMOVED***
				return route, err
			***REMOVED***
		***REMOVED***
		route.Encap = e
	***REMOVED***

	return route, nil
***REMOVED***

// RouteGet gets a route to a specific destination from the host system.
// Equivalent to: 'ip route get'.
func RouteGet(destination net.IP) ([]Route, error) ***REMOVED***
	return pkgHandle.RouteGet(destination)
***REMOVED***

// RouteGet gets a route to a specific destination from the host system.
// Equivalent to: 'ip route get'.
func (h *Handle) RouteGet(destination net.IP) ([]Route, error) ***REMOVED***
	req := h.newNetlinkRequest(syscall.RTM_GETROUTE, syscall.NLM_F_REQUEST)
	family := nl.GetIPFamily(destination)
	var destinationData []byte
	var bitlen uint8
	if family == FAMILY_V4 ***REMOVED***
		destinationData = destination.To4()
		bitlen = 32
	***REMOVED*** else ***REMOVED***
		destinationData = destination.To16()
		bitlen = 128
	***REMOVED***
	msg := &nl.RtMsg***REMOVED******REMOVED***
	msg.Family = uint8(family)
	msg.Dst_len = bitlen
	req.AddData(msg)

	rtaDst := nl.NewRtAttr(syscall.RTA_DST, destinationData)
	req.AddData(rtaDst)

	msgs, err := req.Execute(syscall.NETLINK_ROUTE, syscall.RTM_NEWROUTE)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var res []Route
	for _, m := range msgs ***REMOVED***
		route, err := deserializeRoute(m)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		res = append(res, route)
	***REMOVED***
	return res, nil

***REMOVED***

// RouteSubscribe takes a chan down which notifications will be sent
// when routes are added or deleted. Close the 'done' chan to stop subscription.
func RouteSubscribe(ch chan<- RouteUpdate, done <-chan struct***REMOVED******REMOVED***) error ***REMOVED***
	return routeSubscribeAt(netns.None(), netns.None(), ch, done, nil)
***REMOVED***

// RouteSubscribeAt works like RouteSubscribe plus it allows the caller
// to choose the network namespace in which to subscribe (ns).
func RouteSubscribeAt(ns netns.NsHandle, ch chan<- RouteUpdate, done <-chan struct***REMOVED******REMOVED***) error ***REMOVED***
	return routeSubscribeAt(ns, netns.None(), ch, done, nil)
***REMOVED***

// RouteSubscribeOptions contains a set of options to use with
// RouteSubscribeWithOptions.
type RouteSubscribeOptions struct ***REMOVED***
	Namespace     *netns.NsHandle
	ErrorCallback func(error)
***REMOVED***

// RouteSubscribeWithOptions work like RouteSubscribe but enable to
// provide additional options to modify the behavior. Currently, the
// namespace can be provided as well as an error callback.
func RouteSubscribeWithOptions(ch chan<- RouteUpdate, done <-chan struct***REMOVED******REMOVED***, options RouteSubscribeOptions) error ***REMOVED***
	if options.Namespace == nil ***REMOVED***
		none := netns.None()
		options.Namespace = &none
	***REMOVED***
	return routeSubscribeAt(*options.Namespace, netns.None(), ch, done, options.ErrorCallback)
***REMOVED***

func routeSubscribeAt(newNs, curNs netns.NsHandle, ch chan<- RouteUpdate, done <-chan struct***REMOVED******REMOVED***, cberr func(error)) error ***REMOVED***
	s, err := nl.SubscribeAt(newNs, curNs, syscall.NETLINK_ROUTE, syscall.RTNLGRP_IPV4_ROUTE, syscall.RTNLGRP_IPV6_ROUTE)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if done != nil ***REMOVED***
		go func() ***REMOVED***
			<-done
			s.Close()
		***REMOVED***()
	***REMOVED***
	go func() ***REMOVED***
		defer close(ch)
		for ***REMOVED***
			msgs, err := s.Receive()
			if err != nil ***REMOVED***
				if cberr != nil ***REMOVED***
					cberr(err)
				***REMOVED***
				return
			***REMOVED***
			for _, m := range msgs ***REMOVED***
				route, err := deserializeRoute(m.Data)
				if err != nil ***REMOVED***
					if cberr != nil ***REMOVED***
						cberr(err)
					***REMOVED***
					return
				***REMOVED***
				ch <- RouteUpdate***REMOVED***Type: m.Header.Type, Route: route***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return nil
***REMOVED***
