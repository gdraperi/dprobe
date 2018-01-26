package netlink

import (
	"fmt"
	"syscall"

	"github.com/vishvananda/netlink/nl"
)

func LinkGetProtinfo(link Link) (Protinfo, error) ***REMOVED***
	return pkgHandle.LinkGetProtinfo(link)
***REMOVED***

func (h *Handle) LinkGetProtinfo(link Link) (Protinfo, error) ***REMOVED***
	base := link.Attrs()
	h.ensureIndex(base)
	var pi Protinfo
	req := h.newNetlinkRequest(syscall.RTM_GETLINK, syscall.NLM_F_DUMP)
	msg := nl.NewIfInfomsg(syscall.AF_BRIDGE)
	req.AddData(msg)
	msgs, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	if err != nil ***REMOVED***
		return pi, err
	***REMOVED***

	for _, m := range msgs ***REMOVED***
		ans := nl.DeserializeIfInfomsg(m)
		if int(ans.Index) != base.Index ***REMOVED***
			continue
		***REMOVED***
		attrs, err := nl.ParseRouteAttr(m[ans.Len():])
		if err != nil ***REMOVED***
			return pi, err
		***REMOVED***
		for _, attr := range attrs ***REMOVED***
			if attr.Attr.Type != syscall.IFLA_PROTINFO|syscall.NLA_F_NESTED ***REMOVED***
				continue
			***REMOVED***
			infos, err := nl.ParseRouteAttr(attr.Value)
			if err != nil ***REMOVED***
				return pi, err
			***REMOVED***
			pi = *parseProtinfo(infos)

			return pi, nil
		***REMOVED***
	***REMOVED***
	return pi, fmt.Errorf("Device with index %d not found", base.Index)
***REMOVED***

func parseProtinfo(infos []syscall.NetlinkRouteAttr) *Protinfo ***REMOVED***
	var pi Protinfo
	for _, info := range infos ***REMOVED***
		switch info.Attr.Type ***REMOVED***
		case nl.IFLA_BRPORT_MODE:
			pi.Hairpin = byteToBool(info.Value[0])
		case nl.IFLA_BRPORT_GUARD:
			pi.Guard = byteToBool(info.Value[0])
		case nl.IFLA_BRPORT_FAST_LEAVE:
			pi.FastLeave = byteToBool(info.Value[0])
		case nl.IFLA_BRPORT_PROTECT:
			pi.RootBlock = byteToBool(info.Value[0])
		case nl.IFLA_BRPORT_LEARNING:
			pi.Learning = byteToBool(info.Value[0])
		case nl.IFLA_BRPORT_UNICAST_FLOOD:
			pi.Flood = byteToBool(info.Value[0])
		case nl.IFLA_BRPORT_PROXYARP:
			pi.ProxyArp = byteToBool(info.Value[0])
		case nl.IFLA_BRPORT_PROXYARP_WIFI:
			pi.ProxyArpWiFi = byteToBool(info.Value[0])
		***REMOVED***
	***REMOVED***
	return &pi
***REMOVED***
