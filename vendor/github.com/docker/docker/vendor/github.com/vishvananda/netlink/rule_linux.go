package netlink

import (
	"fmt"
	"net"
	"syscall"

	"github.com/vishvananda/netlink/nl"
)

// RuleAdd adds a rule to the system.
// Equivalent to: ip rule add
func RuleAdd(rule *Rule) error ***REMOVED***
	return pkgHandle.RuleAdd(rule)
***REMOVED***

// RuleAdd adds a rule to the system.
// Equivalent to: ip rule add
func (h *Handle) RuleAdd(rule *Rule) error ***REMOVED***
	req := h.newNetlinkRequest(syscall.RTM_NEWRULE, syscall.NLM_F_CREATE|syscall.NLM_F_EXCL|syscall.NLM_F_ACK)
	return ruleHandle(rule, req)
***REMOVED***

// RuleDel deletes a rule from the system.
// Equivalent to: ip rule del
func RuleDel(rule *Rule) error ***REMOVED***
	return pkgHandle.RuleDel(rule)
***REMOVED***

// RuleDel deletes a rule from the system.
// Equivalent to: ip rule del
func (h *Handle) RuleDel(rule *Rule) error ***REMOVED***
	req := h.newNetlinkRequest(syscall.RTM_DELRULE, syscall.NLM_F_CREATE|syscall.NLM_F_EXCL|syscall.NLM_F_ACK)
	return ruleHandle(rule, req)
***REMOVED***

func ruleHandle(rule *Rule, req *nl.NetlinkRequest) error ***REMOVED***
	msg := nl.NewRtMsg()
	msg.Family = syscall.AF_INET
	if rule.Family != 0 ***REMOVED***
		msg.Family = uint8(rule.Family)
	***REMOVED***
	var dstFamily uint8

	var rtAttrs []*nl.RtAttr
	if rule.Dst != nil && rule.Dst.IP != nil ***REMOVED***
		dstLen, _ := rule.Dst.Mask.Size()
		msg.Dst_len = uint8(dstLen)
		msg.Family = uint8(nl.GetIPFamily(rule.Dst.IP))
		dstFamily = msg.Family
		var dstData []byte
		if msg.Family == syscall.AF_INET ***REMOVED***
			dstData = rule.Dst.IP.To4()
		***REMOVED*** else ***REMOVED***
			dstData = rule.Dst.IP.To16()
		***REMOVED***
		rtAttrs = append(rtAttrs, nl.NewRtAttr(syscall.RTA_DST, dstData))
	***REMOVED***

	if rule.Src != nil && rule.Src.IP != nil ***REMOVED***
		msg.Family = uint8(nl.GetIPFamily(rule.Src.IP))
		if dstFamily != 0 && dstFamily != msg.Family ***REMOVED***
			return fmt.Errorf("source and destination ip are not the same IP family")
		***REMOVED***
		srcLen, _ := rule.Src.Mask.Size()
		msg.Src_len = uint8(srcLen)
		var srcData []byte
		if msg.Family == syscall.AF_INET ***REMOVED***
			srcData = rule.Src.IP.To4()
		***REMOVED*** else ***REMOVED***
			srcData = rule.Src.IP.To16()
		***REMOVED***
		rtAttrs = append(rtAttrs, nl.NewRtAttr(syscall.RTA_SRC, srcData))
	***REMOVED***

	if rule.Table >= 0 ***REMOVED***
		msg.Table = uint8(rule.Table)
		if rule.Table >= 256 ***REMOVED***
			msg.Table = syscall.RT_TABLE_UNSPEC
		***REMOVED***
	***REMOVED***

	req.AddData(msg)
	for i := range rtAttrs ***REMOVED***
		req.AddData(rtAttrs[i])
	***REMOVED***

	native := nl.NativeEndian()

	if rule.Priority >= 0 ***REMOVED***
		b := make([]byte, 4)
		native.PutUint32(b, uint32(rule.Priority))
		req.AddData(nl.NewRtAttr(nl.FRA_PRIORITY, b))
	***REMOVED***
	if rule.Mark >= 0 ***REMOVED***
		b := make([]byte, 4)
		native.PutUint32(b, uint32(rule.Mark))
		req.AddData(nl.NewRtAttr(nl.FRA_FWMARK, b))
	***REMOVED***
	if rule.Mask >= 0 ***REMOVED***
		b := make([]byte, 4)
		native.PutUint32(b, uint32(rule.Mask))
		req.AddData(nl.NewRtAttr(nl.FRA_FWMASK, b))
	***REMOVED***
	if rule.Flow >= 0 ***REMOVED***
		b := make([]byte, 4)
		native.PutUint32(b, uint32(rule.Flow))
		req.AddData(nl.NewRtAttr(nl.FRA_FLOW, b))
	***REMOVED***
	if rule.TunID > 0 ***REMOVED***
		b := make([]byte, 4)
		native.PutUint32(b, uint32(rule.TunID))
		req.AddData(nl.NewRtAttr(nl.FRA_TUN_ID, b))
	***REMOVED***
	if rule.Table >= 256 ***REMOVED***
		b := make([]byte, 4)
		native.PutUint32(b, uint32(rule.Table))
		req.AddData(nl.NewRtAttr(nl.FRA_TABLE, b))
	***REMOVED***
	if msg.Table > 0 ***REMOVED***
		if rule.SuppressPrefixlen >= 0 ***REMOVED***
			b := make([]byte, 4)
			native.PutUint32(b, uint32(rule.SuppressPrefixlen))
			req.AddData(nl.NewRtAttr(nl.FRA_SUPPRESS_PREFIXLEN, b))
		***REMOVED***
		if rule.SuppressIfgroup >= 0 ***REMOVED***
			b := make([]byte, 4)
			native.PutUint32(b, uint32(rule.SuppressIfgroup))
			req.AddData(nl.NewRtAttr(nl.FRA_SUPPRESS_IFGROUP, b))
		***REMOVED***
	***REMOVED***
	if rule.IifName != "" ***REMOVED***
		req.AddData(nl.NewRtAttr(nl.FRA_IIFNAME, []byte(rule.IifName)))
	***REMOVED***
	if rule.OifName != "" ***REMOVED***
		req.AddData(nl.NewRtAttr(nl.FRA_OIFNAME, []byte(rule.OifName)))
	***REMOVED***
	if rule.Goto >= 0 ***REMOVED***
		msg.Type = nl.FR_ACT_NOP
		b := make([]byte, 4)
		native.PutUint32(b, uint32(rule.Goto))
		req.AddData(nl.NewRtAttr(nl.FRA_GOTO, b))
	***REMOVED***

	_, err := req.Execute(syscall.NETLINK_ROUTE, 0)
	return err
***REMOVED***

// RuleList lists rules in the system.
// Equivalent to: ip rule list
func RuleList(family int) ([]Rule, error) ***REMOVED***
	return pkgHandle.RuleList(family)
***REMOVED***

// RuleList lists rules in the system.
// Equivalent to: ip rule list
func (h *Handle) RuleList(family int) ([]Rule, error) ***REMOVED***
	req := h.newNetlinkRequest(syscall.RTM_GETRULE, syscall.NLM_F_DUMP|syscall.NLM_F_REQUEST)
	msg := nl.NewIfInfomsg(family)
	req.AddData(msg)

	msgs, err := req.Execute(syscall.NETLINK_ROUTE, syscall.RTM_NEWRULE)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	native := nl.NativeEndian()
	var res = make([]Rule, 0)
	for i := range msgs ***REMOVED***
		msg := nl.DeserializeRtMsg(msgs[i])
		attrs, err := nl.ParseRouteAttr(msgs[i][msg.Len():])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		rule := NewRule()

		for j := range attrs ***REMOVED***
			switch attrs[j].Attr.Type ***REMOVED***
			case syscall.RTA_TABLE:
				rule.Table = int(native.Uint32(attrs[j].Value[0:4]))
			case nl.FRA_SRC:
				rule.Src = &net.IPNet***REMOVED***
					IP:   attrs[j].Value,
					Mask: net.CIDRMask(int(msg.Src_len), 8*len(attrs[j].Value)),
				***REMOVED***
			case nl.FRA_DST:
				rule.Dst = &net.IPNet***REMOVED***
					IP:   attrs[j].Value,
					Mask: net.CIDRMask(int(msg.Dst_len), 8*len(attrs[j].Value)),
				***REMOVED***
			case nl.FRA_FWMARK:
				rule.Mark = int(native.Uint32(attrs[j].Value[0:4]))
			case nl.FRA_FWMASK:
				rule.Mask = int(native.Uint32(attrs[j].Value[0:4]))
			case nl.FRA_TUN_ID:
				rule.TunID = uint(native.Uint64(attrs[j].Value[0:4]))
			case nl.FRA_IIFNAME:
				rule.IifName = string(attrs[j].Value[:len(attrs[j].Value)-1])
			case nl.FRA_OIFNAME:
				rule.OifName = string(attrs[j].Value[:len(attrs[j].Value)-1])
			case nl.FRA_SUPPRESS_PREFIXLEN:
				i := native.Uint32(attrs[j].Value[0:4])
				if i != 0xffffffff ***REMOVED***
					rule.SuppressPrefixlen = int(i)
				***REMOVED***
			case nl.FRA_SUPPRESS_IFGROUP:
				i := native.Uint32(attrs[j].Value[0:4])
				if i != 0xffffffff ***REMOVED***
					rule.SuppressIfgroup = int(i)
				***REMOVED***
			case nl.FRA_FLOW:
				rule.Flow = int(native.Uint32(attrs[j].Value[0:4]))
			case nl.FRA_GOTO:
				rule.Goto = int(native.Uint32(attrs[j].Value[0:4]))
			case nl.FRA_PRIORITY:
				rule.Priority = int(native.Uint32(attrs[j].Value[0:4]))
			***REMOVED***
		***REMOVED***
		res = append(res, *rule)
	***REMOVED***

	return res, nil
***REMOVED***
